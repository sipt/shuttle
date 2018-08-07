package selector

import (
	"github.com/sipt/shuttle"
)

func init() {
	shuttle.RegisterSelector("select", func(group *shuttle.ServerGroup) (shuttle.ISelector, error) {
		s := &manualSelector{
			group: group,
		}
		s.Refresh()
		return s, nil
	})
}

type manualSelector struct {
	group    *shuttle.ServerGroup
	selected shuttle.IServer
}

func (m *manualSelector) Get() (*shuttle.Server, error) {
	return m.selected.GetServer()
}
func (m *manualSelector) Select(name string) error {
	var (
		n  shuttle.IServer
		ok bool
	)
	for _, v := range m.group.Servers {
		n, ok = v.(shuttle.IServer)
		if ok && n.GetName() == name {
			m.selected = n
		}
	}
	return nil
}
func (m *manualSelector) Refresh() error {
	m.selected = m.group.Servers[0].(shuttle.IServer)
	return nil
}
func (m *manualSelector) Reset(group *shuttle.ServerGroup) error {
	m.group = group
	m.selected = m.group.Servers[0].(shuttle.IServer)
	return nil
}
func (m *manualSelector) Destroy() {}
