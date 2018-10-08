package shuttle

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
	"github.com/sipt/shuttle/log"
)

const (
	DNSTypeStatic = "static"
	DNSTypeDirect = "direct"
	DNSTypeRemote = "remote"
	DNSTypeCache  = "cache"
)

type DomainHost struct {
	IP            net.IP
	Port          uint16
	Country       string
	DNS           net.IP
	RemoteResolve bool
}

type DNS struct {
	MatchType string `json:",omitempty"`
	Domain    string
	IPs       []net.IP
	DNSs      []net.IP
	Port      uint16
	Type      string `json:",omitempty"`
	Country   string
}

func (d *DNS) String() string {
	buffer := bytes.NewBufferString(d.Domain)
	buffer.WriteString(" IPs:[")
	for i, v := range d.IPs {
		if i > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(v.String())
	}
	buffer.WriteString("] DNSs:[")
	for i, v := range d.DNSs {
		if i > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(v.String())
	}
	buffer.WriteString("] Country:[" + d.Country + "]")
	return buffer.String()
}

type _Reply struct {
	DNS net.IP
	Msg *dns.Msg
}

var (
	_DNS      []net.IP
	_LocalDNS []*DNS
	_CacheDNS IDNSCache
)

func InitDNS(dns []net.IP, localDNS []*DNS) error {
	_DNS = dns
	cdns := &DNS{
		MatchType: RuleDomain,
		Domain:    ControllerDomain,
		IPs:       []net.IP{{127, 0, 0, 1}},
		Type:      DNSTypeStatic,
	}
	cdns.Port, _ = strToUint16(controllerPort)
	_LocalDNS = append([]*DNS{cdns}, localDNS...)
	if _CacheDNS == nil {
		_CacheDNS = NewDefaultDNSCache()
	}
	_CacheDNS.Init()
	for _, v := range _LocalDNS {
		if v.Type == DNSTypeStatic {
			_CacheDNS.Push(v)
		}
	}
	return nil
}

//清除DNS缓存
func ClearDNSCache() {
	_CacheDNS.Clear()
	for _, v := range _LocalDNS {
		if v.Type == DNSTypeStatic {
			_CacheDNS.Push(v)
		}
	}
}

//DNS缓存列表
func DNSCacheList() []*DNS {
	return _CacheDNS.List()
}

func ResolveDomain(req *Request) error {
	if req.IP = net.ParseIP(req.Addr); len(req.IP) != 0 {
		return nil
	}
	//DomainHost
	temp := _LocalDNS
	for _, v := range temp {
		switch v.MatchType {
		case RuleDomainSuffix:
			if strings.HasSuffix(req.Addr, v.Domain) {
				log.Logger.Debug("[DNS] [Local] ", v.String())
				return localResolve(v, req)
			}
		case RuleDomain:
			if req.Addr == v.Domain {
				log.Logger.Debug("[DNS] [Local] ", v.String())
				return localResolve(v, req)
			}
		case RuleDomainKeyword:
			if strings.Index(req.Addr, v.Domain) >= 0 {
				log.Logger.Debug("[DNS] [Local] ", v.String())
				return localResolve(v, req)
			}
		}
	}
	d := _CacheDNS.Pop(req.Addr)
	if d != nil {
		log.Logger.Debug("[DNS] [Cache] ", d.String())
		return localResolve(d, req)
	}
	return directResolve(_DNS, req)
}

func localResolve(dns *DNS, req *Request) error {
	switch dns.Type {
	case DNSTypeStatic, DNSTypeCache:
		if len(req.IP) >= len(dns.IPs[0]) {
			req.IP = req.IP[:len(dns.IPs[0])]
		} else {
			req.IP = make([]byte, len(dns.IPs[0]))
		}
		copy(req.IP, dns.IPs[0])
		req.DomainHost.Country = dns.Country
		if len(dns.DNSs) > 0 {
			req.DomainHost.DNS = dns.DNSs[0]
		}
		if req.Port == 0 && dns.Port > 0 {
			req.Port = dns.Port
		}
		return nil
	case DNSTypeDirect:
		//直连DNS解析
		return directResolve(dns.DNSs, req)
	case DNSTypeRemote:
		//remote 到proxy-server上解析
		//remoteResolve(dns, req)
		log.Logger.Debug("[DNS] [Remote] ", req.Addr)
		return nil
	default:
		return fmt.Errorf("not support DNSType [%s]", dns.Type)
	}
}

func directResolve(servers []net.IP, req *Request) error {
	reply := make(chan *_Reply, 1)
	for _, v := range servers {
		go func(v net.IP) {
			err := resolveDomain(req, v, reply)
			if err != nil {
				log.Logger.Errorf("[DNS] [%s] failed: %v", req.Addr, err)
			}
		}(v)
	}
	r := <-reply
	var (
		a     *dns.A
		ok    bool
		cache = &DNS{
			Domain: req.Addr,
			IPs:    make([]net.IP, 0),
			DNSs:   []net.IP{r.DNS},
			Type:   DNSTypeCache,
		}
	)
	for _, v := range r.Msg.Answer {
		a, ok = v.(*dns.A)
		if ok {
			cache.IPs = append(cache.IPs, net.ParseIP(a.A.String()))
		}
	}
	if len(cache.IPs) == 0 {
		log.Logger.Errorf("[DNS] resolve ip is empty")
		return nil
	}
	req.IP = cache.IPs[0]
	_CacheDNS.Push(cache)
	country := GeoLookUp(req.IP)
	cache.Country = country
	req.DomainHost.Country = country
	req.DomainHost.DNS = r.DNS
	log.Logger.Debug("[DNS] ", cache.String())
	return nil
}

func resolveDomain(req *Request, s net.IP, c chan *_Reply) (err error) {
	conn, err := DirectConn(&Request{
		Cmd:  CmdUDP,
		IP:   s,
		Port: 53,
	})
	if err != nil {
		return err
	}
	m := &dns.Msg{}
	m.SetQuestion(dns.Fqdn(req.Addr), dns.TypeA)
	m.RecursionDesired = true
	conn, err = RealTimeDecorate(conn)
	if err != nil {
		return err
	}
	co := &dns.Conn{Conn: conn} // c is your net.Conn
	co.WriteMsg(m)
	r, err := co.ReadMsg()
	co.Close()
	if err != nil {
		return err
	}
	if r == nil || r.Rcode != dns.RcodeSuccess {
		return fmt.Errorf("resolve domain [%s] failed", req.Addr)
	}
	select {
	case c <- &_Reply{DNS: s, Msg: r}:
	default:
	}
	return
}

func NewDefaultDNSCache() IDNSCache {
	r := &DefaultDNSCache{}
	return r
}

type IDNSCache interface {
	Init()
	Push(*DNS)
	Pop(string) *DNS
	Clear()
	List() []*DNS
}

type CacheNode struct {
	v *DNS
	t time.Time
}

type DefaultDNSCache struct {
	vs       []*CacheNode
	rw       *sync.RWMutex
	duration time.Duration
}

func (d *DefaultDNSCache) Init() {
	d.rw = &sync.RWMutex{}
	d.vs = make([]*CacheNode, 0, 16)
	d.duration = 10 * time.Minute // DNS缓存: 60s
}

func (d *DefaultDNSCache) Push(dns *DNS) {
	temp := d.Pop(dns.Domain)
	if temp == nil {
		d.rw.Lock()
		d.vs = append(d.vs, &CacheNode{dns, time.Now().Add(d.duration)})
		d.rw.Unlock()
	}
}

func (d *DefaultDNSCache) Pop(domain string) *DNS {
	d.rw.Lock()
	defer d.rw.Unlock()
	temp := d.vs
	for i, v := range temp {
		if v.v.Domain == domain {
			if v.t.Before(time.Now()) {
				if d.vs[i].v.Domain == v.v.Domain && d.vs[i].t.Before(time.Now()) {
					d.vs[i] = d.vs[len(d.vs)-1]
					d.vs = d.vs[:len(d.vs)-1]
				}
				return nil
			}
			return temp[i].v
		}
	}
	return nil
}

func (d *DefaultDNSCache) Clear() {
	d.rw.Lock()
	d.vs = d.vs[:0]
	d.rw.Unlock()
}
func (d *DefaultDNSCache) List() []*DNS {
	reply := make([]*DNS, len(d.vs))
	d.rw.RLock()
	for i := range d.vs {
		v := *(d.vs[i].v)
		reply[i] = &v
	}
	d.rw.RUnlock()
	return reply
}
