package shuttle

import (
	"net"
	"strings"
	"bytes"
	"github.com/miekg/dns"
	"fmt"
	"time"
	"sync"
	"github.com/sipt/shuttle/util"
)

const (
	DNSTypeStatic = "static"
	DNSTypeDirect = "direct"
	DNSTypeRemote = "remote"
)

type DomainHost struct {
	IP            net.IP
	Country       string
	DNS           net.IP
	RemoteResolve bool
}

type DNS struct {
	MatchType string
	Domain    string
	IPs       []net.IP
	DNSs      []net.IP
	Type      string
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
	_LocalDNS = localDNS
	if _CacheDNS == nil {
		_CacheDNS = NewDefaultDNSCache()
	}
	_CacheDNS.Init()
	return nil
}

//清除DNS缓存
func ClearDNSCache() {
	_CacheDNS.Clear()
}

func ResolveDomain(req *Request) error {
	//DomainHost
	temp := _LocalDNS
	for _, v := range temp {
		switch v.MatchType {
		case RuleDomainSuffix:
			if strings.HasSuffix(req.Addr, v.Domain) {
				Logger.Debug("[DNS] [Local] ", v.String())
				return localResolve(v, req)
			}
		case RuleDomain:
			if req.Addr == v.Domain {
				Logger.Debug("[DNS] [Local] ", v.String())
				return localResolve(v, req)
			}
		case RuleDomainKeyword:
			if strings.Index(req.Addr, v.Domain) >= 0 {
				Logger.Debug("[DNS] [Local] ", v.String())
				return localResolve(v, req)
			}
		}
	}
	d := _CacheDNS.Pop(req.Addr)
	if d != nil {
		Logger.Debug("[DNS] [Cache] ", d.String())
		return localResolve(d, req)
	}
	return directResolve(_DNS, req)
}

func localResolve(dns *DNS, req *Request) error {
	switch dns.Type {
	case DNSTypeStatic:
		if len(req.IP) >= len(dns.IPs[0]) {
			req.IP = req.IP[:len(dns.IPs[0])]
		} else {
			req.IP = make([]byte, len(dns.IPs[0]))
		}
		copy(req.IP, dns.IPs[0])
		return nil
	case DNSTypeDirect:
		//直连DNS解析
		return directResolve(dns.DNSs, req)
	case DNSTypeRemote:
		//remote 到proxy-server上解析
		//remoteResolve(dns, req)
		Logger.Debug("[DNS] [Remote] ", req.Addr)
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
				Logger.Error("[DNS] [%s] failed: ", err)
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
		}
	)
	for _, v := range r.Msg.Answer {
		a, ok = v.(*dns.A)
		if ok {
			cache.IPs = append(cache.IPs, net.ParseIP(a.A.String()))
		}
	}
	_CacheDNS.Push(cache)
	req.IP = cache.IPs[0]
	ip, err := util.WatchIP(req.IP.String())
	if err != nil {
		Logger.Errorf("[DNS] watch ip[%s] country failed : %v", req.IP.String(), err)
	}
	cache.Country = ip.CountryID
	req.DomainHost.Country = ip.CountryID
	req.DomainHost.DNS = r.DNS
	Logger.Debug("[DNS] ", cache.String())
	return nil
}

func resolveDomain(req *Request, s net.IP, c chan *_Reply) (err error) {
	conn, err := DirectConn(&Request{
		Cmd:  cmdUDP,
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
	d.duration = 60 * time.Second // DNS缓存: 60s
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
	temp := d.vs
	for i, v := range temp {
		if v.v.Domain == domain {
			if v.t.Before(time.Now()) {
				d.rw.Lock()
				if d.vs[i].v.Domain == v.v.Domain && d.vs[i].t.Before(time.Now()) {
					d.vs[i] = d.vs[len(d.vs)-1]
					d.vs = d.vs[:len(d.vs)-1]
				}
				d.rw.Unlock()
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
