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

type ProxyExternal2 struct {
	ProxyExternal
	SubName    string `json:"sub_name"`
	IsSelected bool   `json:"is_selected"`
}

type GroupExternal struct {
	Name       string            `json:"name"`
	SelectType string            `json:"select_type"`
	Servers    []*ProxyExternal2 `json:"servers"`
	Selected   *ProxyExternal2   `json:"selected"`
	RttUrl     string            `json:"rtt_url"`
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
		return (*ProxyExternal)(unsafe.Pointer(x)).Name < (*ProxyExternal)(unsafe.Pointer(y)).Name
	})
	return reply
}

func GetGroupExternals(names ...string) []*GroupExternal {
	groupLock.RLock()
	defer groupLock.RUnlock()
	reply := make([]*GroupExternal, 0, len(groups))
	for _, v := range groups {
		if len(names) > 0 {
			isExist := false
			for _, name := range names {
				if v.Name == name {
					isExist = true
					break
				}
			}
			if !isExist {
				continue
			}
		}
		g := &GroupExternal{
			Name:       v.Name,
			SelectType: v.SelectType,
			RttUrl:     v.RttUrl,
		}
		if len(v.RttUrl) == 0 {
			g.RttUrl = globalRttUrl
		}
		selected, _ := v.Selector.Current().GetServer()
		g.Selected = &ProxyExternal2{
			ProxyExternal: ProxyExternal{
				Name:     selected.Name,
				Rtt:      selected.Rtt,
				RttText:  Duration2Str(selected.Rtt),
				Protocol: selected.ProxyProtocol,
			},
		}
		g.Servers = make([]*ProxyExternal2, len(v.Servers))
		for j, x := range v.Servers {
			p := &ProxyExternal2{}
			is := x.(IServer)
			p.Name = is.GetName()
			s, _ := is.GetServer()
			if s.Name != p.Name {
				p.SubName = s.Name
			}
			p.Rtt = s.Rtt
			p.RttText = Duration2Str(s.Rtt)
			p.Protocol = s.ProxyProtocol
			p.IsSelected = v.Selector.Current().GetName() == p.Name
			g.Servers[j] = p
		}
		reply = append(reply, g)
	}
	util.QuickSort2(reply, func(x, y uintptr) bool {
		return (*ProxyExternal2)(unsafe.Pointer(x)).Name < (*ProxyExternal2)(unsafe.Pointer(y)).Name
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
	if err != nil {
		return err
	}
	s.RttUrl = rttUrl
	for _, v := range groups {
		if v.Name == ProxyGlobal {
			for i, is := range v.Servers {
				if is.(IServer).GetName() == name {
					v.Servers[i] = s
				}
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
	if err != nil {
		return
	}
	for _, v := range groups {
		if v.Name == ProxyGlobal {
			v.Servers = append(v.Servers, g)
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
	if err != nil {
		return
	}

	for i, v := range groups {
		if v.Name == name {
			groups[i] = g
		} else if v.Name == ProxyGlobal {
			for i, is := range v.Servers {
				if is.(IServer).GetName() == name {
					v.Servers[i] = g
				}
			}
		}
	}
	return
}

func RemoveGroup(name string) (effects, deletes []string, err error) {
	groupLock.Lock()
	defer groupLock.Unlock()
	effects = make([]string, 0, 4)
	deletes = make([]string, 0, 4)
	var g *ServerGroup
	for i := len(groups) - 1; i >= 0; i-- {
		g = groups[i]
		if g.Name == name {
			groups = append(groups[:i], groups[i+1:]...)
			g.Selector.Destroy()
			continue
		}
		fmt.Println("------>", name, g.Name, len(g.Servers))
		if g.Remove(name) {
			fmt.Println("------>", name, g.Name, len(g.Servers))
			if g.Name == ProxyGlobal {
				continue
			}
			fmt.Println("------>", name, g.Name, len(g.Servers))
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
