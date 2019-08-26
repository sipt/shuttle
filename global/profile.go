package global

import (
	"sync"

	"github.com/sipt/shuttle/group"
	"github.com/sipt/shuttle/server"

	"github.com/sipt/shuttle/conf/model"
	"github.com/sipt/shuttle/dns"
	"github.com/sipt/shuttle/rule"
)

var profileMap = make(map[string]*Profile)
var profileMapMutex = &sync.RWMutex{}

func GetProfile(key string) *Profile {
	profileMapMutex.RLock()
	defer profileMapMutex.RUnlock()
	return profileMap[key]
}

func AddProfile(key string, profile *Profile) {
	profileMapMutex.Lock()
	defer profileMapMutex.Unlock()
	profileMap[key] = profile
}

func RemoveProfile(key string) {
	profileMapMutex.Lock()
	defer profileMapMutex.Unlock()
	delete(profileMap, key)
}

func NewProfile(
	config *model.Config,
	dnsHandle dns.Handle,
	ruleHandle rule.Handle,
	group map[string]group.IGroup,
	server map[string]server.IServer) (*Profile, error) {
	return &Profile{
		uri:        config.Info.URI,
		config:     config,
		dnsHandle:  dnsHandle,
		ruleHandle: ruleHandle,
		group:      group,
		server:     server,
	}, nil
}

type Profile struct {
	uri        string
	config     *model.Config
	dnsHandle  dns.Handle
	ruleHandle rule.Handle
	group      map[string]group.IGroup
	server     map[string]server.IServer
}

func (p *Profile) URI() string {
	return p.uri
}

func (p *Profile) Config() *model.Config {
	return p.config
}

func (p *Profile) DNSHandle() dns.Handle {
	return p.dnsHandle
}

func (p *Profile) RuleHandle() rule.Handle {
	return p.ruleHandle
}

func (p *Profile) Group() map[string]group.IGroup {
	return p.group
}

func (p *Profile) Server() map[string]server.IServer {
	return p.server
}
