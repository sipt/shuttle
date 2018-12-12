package selector

import (
	"fmt"
	"github.com/sipt/shuttle/proxy"
)

func init() {
	proxy.RegisterSelector("select", func(group *proxy.ServerGroup) (proxy.ISelector, error) {
		s := &manualSelector{
			group: group,
		}
		s.Refresh()
		return s, nil
	})
}

type manualSelector struct {
	group    *proxy.ServerGroup
	selected proxy.IServer
}

func (m *manualSelector) Get() (*proxy.Server, error) {
	return m.selected.GetServer()
}
func (m *manualSelector) Select(name string) error {
	var (
		n  proxy.IServer
		ok bool
	)
	for _, v := range m.group.Servers {
		n, ok = v.(proxy.IServer)
		if ok && n.GetName() == name {
			m.selected = n
			return nil
		}
	}
	return fmt.Errorf("server[%s] is not exist", name)
}
func (m *manualSelector) Refresh() error {
	m.selected = m.group.Servers[0].(proxy.IServer)
	return nil
}
func (m *manualSelector) Reset(group *proxy.ServerGroup) error {
	m.group = group
	m.selected = m.group.Servers[0].(proxy.IServer)
	return nil
}
func (m *manualSelector) Destroy() {}
func (m *manualSelector) Current() proxy.IServer {
	return m.selected
}
