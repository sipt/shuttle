package global

import (
	"sync"

	"github.com/sipt/shuttle/conf/model"
	"github.com/sipt/shuttle/conn/filter"
	"github.com/sipt/shuttle/conn/stream"
	"github.com/sipt/shuttle/dns"
	"github.com/sipt/shuttle/group"
	"github.com/sipt/shuttle/rule"
	"github.com/sipt/shuttle/server"
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
	dnsCache dns.ICache,
	ruleHandle rule.Handle,
	group map[interface{}]group.IGroup,
	server map[string]server.IServer,
	filter filter.FilterFunc,
	before, after stream.DecorateFunc) (*Profile, error) {
	return &Profile{
		uri:        config.Info.URI,
		config:     config,
		dnsHandle:  dnsHandle,
		dnsCache:   dnsCache,
		ruleHandle: ruleHandle,
		group:      group,
		server:     server,
		filter:     filter,
		before:     before,
		after:      after,
	}, nil
}

type Profile struct {
	uri           string
	config        *model.Config
	dnsHandle     dns.Handle
	dnsCache      dns.ICache
	ruleHandle    rule.Handle
	group         map[interface{}]group.IGroup
	server        map[string]server.IServer
	filter        filter.FilterFunc
	before, after stream.DecorateFunc
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

func (p *Profile) DNSCache() dns.ICache {
	return p.dnsCache
}

func (p *Profile) RuleHandle() rule.Handle {
	return p.ruleHandle
}

func (p *Profile) Group() map[interface{}]group.IGroup {
	return p.group
}

func (p *Profile) Server() map[string]server.IServer {
	return p.server
}

func (p *Profile) Filter() filter.FilterFunc {
	return p.filter
}

func (p *Profile) BeforeStream() stream.DecorateFunc {
	return p.before
}

func (p *Profile) AfterStream() stream.DecorateFunc {
	return p.after
}
