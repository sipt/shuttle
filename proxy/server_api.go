package proxy

import (
	"fmt"
	"github.com/sipt/shuttle/util"
	"strings"
	"sync"
	"time"
	"unsafe"
)

var serverLock = &sync.RWMutex{}
var groupLock = &sync.RWMutex{}

type ProxyExternal struct {
	Name     string        `json:"name"`
	Rtt      time.Duration `json:"rtt"`
	RttText  string        `json:"rtt_text"`
	Protocol string        `json:"protocol"`
}

type GroupExternal struct {
	Name       string           `json:"name"`
	SelectType string           `json:"select_type"`
	Servers    []*ProxyExternal `json:"servers"`
	Selected   *ProxyExternal   `json:"selected"`
}

func GetServerExternals() []*ProxyExternal {
	serverLock.RLock()
	defer serverLock.RUnlock()
	reply := make([]*ProxyExternal, 0, len(servers))
	for _, v := range servers {
		switch v.Name {
		case ProxyReject, ProxyDirect:
		default:
			reply = append(reply, &ProxyExternal{
				Name:     v.Name,
				Rtt:      v.Rtt,
				RttText:  Duration2Str(v.Rtt),
				Protocol: v.ProxyProtocol,
			})
		}
	}
	util.QuickSort2(reply, func(x, y uintptr) bool {
		return (*ProxyExternal)(unsafe.Pointer(x)).Name > (*ProxyExternal)(unsafe.Pointer(y)).Name
	})
	return reply
}

func GetGroupExternals() []*GroupExternal {
	groupLock.RLock()
	defer groupLock.RUnlock()
	reply := make([]*GroupExternal, len(groups))
	for i, v := range groups {
		reply[i] = &GroupExternal{
			Name:       v.Name,
			SelectType: v.SelectType,
		}
		selected, _ := v.Selector.Current().GetServer()
		reply[i].Selected = &ProxyExternal{
			Name:     selected.Name,
			Rtt:      selected.Rtt,
			RttText:  Duration2Str(selected.Rtt),
			Protocol: selected.ProxyProtocol,
		}
		reply[i].Servers = make([]*ProxyExternal, len(v.Servers))
		for j, v := range v.Servers {
			p := &ProxyExternal{}
			is := v.(IServer)
			p.Name = is.GetName()
			s, _ := is.GetServer()
			if s.Name != p.Name {
				p.Name = fmt.Sprintf("%s(%s)", p.Name, s.Name)
			}
			p.Rtt = s.Rtt
			p.RttText = Duration2Str(s.Rtt)
			p.Protocol = s.ProxyProtocol
			reply[i].Servers[j] = p
		}
	}
	util.QuickSort2(reply, func(x, y uintptr) bool {
		return (*ProxyExternal)(unsafe.Pointer(x)).Name > (*ProxyExternal)(unsafe.Pointer(y)).Name
	})
	return reply
}

func AddProxy(name string, vs []string) error {
	serverLock.Lock()
	defer serverLock.Unlock()
	if _, isExist := ProxyExist(name); isExist {
		return fmt.Errorf("[%s] duplicate [Proxy]", name)
	}
	groupLock.RLock()
	if _, isExist := GroupExist(name); isExist {
		groupLock.RUnlock()
		return fmt.Errorf("[%s] duplicate [ProxyGroup]", name)
	}
	groupLock.RUnlock()
	if len(vs) < 2 {
		return fmt.Errorf("[AddProxy] [%s] failed", name)
	}
	rttUrl := ""
	{ //rtt Url check
		last := vs[len(vs)-1]
		if len(last) > len("http://") {
			if strings.HasPrefix(last, "http://") || strings.HasPrefix(last, "https://") {
				rttUrl = last
				vs = vs[:len(vs)-1]
			}
		}
	}
	s, err := NewServer(name, vs)
	if err != nil {
		return err
	}
	s.RttUrl = rttUrl
	for _, v := range groups {
		if v.Name == ProxyGlobal {
			v.Servers = append(v.Servers, s)
			if err != nil {
				return err
			}
		}
	}
	servers = append(servers, s)

	return err
}
func EditProxy(name string, vs []string) error {
	serverLock.Lock()
	defer serverLock.Unlock()
	var index = -1
	for i, s := range servers {
		if s.Name == name {
			index = i
			break
		}
	}
	if index == -1 {
		return fmt.Errorf("[Proxy: %s] not found", name)
	}
	if len(vs) < 2 {
		return fmt.Errorf("[EditProxy] [%s] failed", name)
	}
	rttUrl := ""
	{ //rtt Url check
		last := vs[len(vs)-1]
		if len(last) > len("http://") {
			if strings.HasPrefix(last, "http://") || strings.HasPrefix(last, "https://") {
				rttUrl = last
				vs = vs[:len(vs)-1]
			}
		}
	}
	s, err := NewServer(name, vs)
	s.RttUrl = rttUrl
	for _, v := range groups {
		if v.Name == ProxyGlobal {
			for i, is := range v.Servers {
				if is.(IServer).GetName() == name {
					v.Servers[i] = s
				}
			}
			if err != nil {
				return err
			}
		}
	}
	servers[index] = s
	return err
}
func RemoveProxy(name string) (effects, deletes []string, err error) {
	//remove
	serverLock.Lock()
	defer serverLock.Unlock()
	isExist := false
	for i, v := range servers {
		if v.Name == name {
			servers = append(servers[:i], servers[i+1:]...)
			isExist = true
			break
		}
	}
	if !isExist {
		err = fmt.Errorf("[ProxyGroup: %s] not found", name)
		return
	}

	//remove in group
	effects = make([]string, 0, 4)
	deletes = make([]string, 0, 4)
	var g *ServerGroup
	for i := len(groups) - 1; i >= 0; i-- {
		g = groups[i]
		if g.Remove(name) {
			if g.Name == ProxyGlobal {
				continue
			}
			if len(g.Servers) == 0 {
				deletes = append(deletes, g.Name)
				groups = append(groups[:i], groups[i+1:]...)
			} else {
				effects = append(effects, g.Name)
			}
		}
	}
	return
}

func AddGroup(name string, vs []string) (err error) {
	groupLock.Lock()
	defer groupLock.Unlock()
	if _, isExist := GroupExist(name); isExist {
		return fmt.Errorf("[%s] duplicate [ProxyGroup]", name)
	}
	serverLock.RLock()
	if _, isExist := ProxyExist(name); isExist {
		serverLock.RUnlock()
		return fmt.Errorf("[%s] duplicate [Proxy]", name)
	}
	serverLock.RUnlock()

	//add group
	g := &ServerGroup{
		Name: name,
	}
	if len(vs) < 2 {
		return fmt.Errorf("[AddGroup] [%s] failed", name)
	}
	g.SelectType = vs[0]
	vs = vs[1:]
	{ //rtt Url check
		last := vs[len(vs)-1]
		if len(last) > len("http://") {
			if strings.HasPrefix(last, "http://") || strings.HasPrefix(last, "https://") {
				g.RttUrl = last
				vs = vs[:len(vs)-1]
			}
		}
	}
	g.Servers = make([]interface{}, len(vs))
	var isExist bool
	for i, v := range vs {
		isExist = false
		g.Servers[i], isExist = ProxyExist(v)
		if !isExist {
			g.Servers[i], isExist = GroupExist(v)
			if !isExist {
				return fmt.Errorf("[Proxy or ProxyGroup: %s] not found", v)
			}
		}
	}
	g.Selector, err = GetSelector(g.SelectType, g)

	for _, v := range groups {
		if v.Name == ProxyGlobal {
			v.Servers = append(v.Servers, g)
			if err != nil {
				return err
			}
		}
	}
	groups = append(groups, g)
	return
}

func EditGroup(name string, vs []string) (err error) {
	groupLock.Lock()
	defer groupLock.Unlock()
	var index = -1
	for i, s := range groups {
		if s.Name == name {
			index = i
			break
		}
	}
	if index == -1 {
		return fmt.Errorf("[ProxyGroup: %s] not found", name)
	}
	//edit group
	g := &ServerGroup{
		Name: name,
	}
	if len(vs) < 2 {
		return fmt.Errorf("[AddGroup] [%s] failed", name)
	}
	g.SelectType = vs[0]
	vs = vs[1:]
	{ //rtt Url check
		last := vs[len(vs)-1]
		if len(last) > len("http://") {
			if strings.HasPrefix(last, "http://") || strings.HasPrefix(last, "https://") {
				g.RttUrl = last
				vs = vs[:len(vs)-1]
			}
		}
	}
	g.Servers = make([]interface{}, len(vs))
	var isExist bool
	for i, v := range vs {
		isExist = false
		g.Servers[i], isExist = ProxyExist(v)
		if !isExist {
			g.Servers[i], isExist = GroupExist(v)
			if !isExist {
				return fmt.Errorf("[Proxy or ProxyGroup: %s] not found", v)
			}
		}
	}
	g.Selector, err = GetSelector(g.SelectType, g)

	for _, v := range groups {
		if v.Name == ProxyGlobal {
			for i, is := range v.Servers {
				if is.(IServer).GetName() == name {
					v.Servers[i] = g
				}
			}
			if err != nil {
				return err
			}
		}
	}
	groups = append(groups, g)
	return
}

func RemoveGroup(name string) error {
	groupLock.Lock()
	defer groupLock.Unlock()
	for i, v := range groups {
		if v.Name == ProxyGlobal {
			v.Remove(name)
		}
		if v.Name == name {
			groups = append(groups[:i], groups[i+1:]...)
			v.Selector.Destroy()
			return nil
		}
	}
	return fmt.Errorf("[ProxyGroup: %s] not found", name)
}

func ProxyExist(name string) (*Server, bool) {
	for _, v := range servers {
		if v.Name == name {
			return v, true
		}
	}
	return nil, false
}

func GroupExist(name string) (*ServerGroup, bool) {
	for _, v := range groups {
		if v.Name == name {
			return v, true
		}
	}
	return nil, false
}
