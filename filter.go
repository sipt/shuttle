package shuttle

import (
	"errors"
	"github.com/sipt/shuttle/dns"
	"github.com/sipt/shuttle/log"
	"github.com/sipt/shuttle/proxy"
	"github.com/sipt/shuttle/rule"
)

type IRequest interface {
	Network() string
	Domain() string
	IP() string
	Port() string
	Answer() *dns.Answer
	SetAnswer(*dns.Answer)

	ID() int64    //return request id
	Host() string //return [domain/ip]:[port]
	Addr() string //return domain!=""?domain:ip
}

func FilterByReq(req IRequest) (r *rule.Rule, s *proxy.Server, err error) {
	//DNS
	var answer *dns.Answer
	if len(req.IP()) == 0 {
		answer, err = dns.ResolveDomainByCache(req.Domain())
	} else {
		answer, err = dns.ResolveIP(req.IP())
	}
	if err != nil {
		// skip error
		log.Logger.Errorf("[FilterByReq] %s", err.Error())
		return
	}
	req.SetAnswer(answer)
	//Rules RuleFilter
	r, err = rule.RuleFilter(req)
	if err != nil {
		return
	}
	if r == rule.RejectRule {
		s, _ = proxy.GetServer(r.Policy)
		err = ErrorReject
		return
	}
	if r == nil {
		log.Logger.Infof("[RULE] [ID:%d] [%s] rule: [%v]", req.ID(), req.Host(), rule.PolicyDirect)
		s, err = proxy.GetServer(rule.PolicyDirect) // 没有匹配规则，直连
	} else {
		country := ""
		if req.Answer() != nil {
			country = req.Answer().Country
		}
		log.Logger.Infof("[RULE] [ID:%d] [%s, %s, %s] rule: [%s, %s, %s]", req.ID(), req.Host(), req.Addr(),
			country, r.Type, r.Value, r.Policy)
		//Select proxy server
		s, err = proxy.GetServer(r.Policy)
		if err != nil {
			err = errors.New(err.Error() + ":" + r.Policy)
			return
		}
		log.Logger.Debugf("[RULE] [ID:%d] Get server by policy [%s] => [%s]", req.ID(), r.Policy, s.Name)
	}
	return
}
