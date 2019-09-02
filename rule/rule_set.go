package rule

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/pkg/errors"
	"github.com/sipt/shuttle/dns"
)

const (
	KeyRuleSet   = "RULE-SET"
	ParamsKeyURL = "url"
)

func init() {
	Register(KeyRuleSet, ruleSetHandle)
}

func ruleSetHandle(rule *Rule, next Handle, dnsHandle dns.Handle) (Handle, error) {
	url := rule.Params[ParamsKeyURL]
	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "download rule set failed: %s", url)
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, errors.Errorf("download rule set failed[%d]: %s", resp.StatusCode, url)
	}
	defer resp.Body.Close()
	r := bufio.NewReader(resp.Body)
	var (
		bs    []byte
		lnum  int
		line  string
		rules = make([]*Rule, 0, 64)
	)
	for {
		lnum++
		bs, _, err = r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, errors.Errorf("read line failed, rule set failed: %s:[%d]", url, lnum)
		}
		line = string(bs)
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		cells := strings.Split(line, ",")
		r := *rule
		for i, v := range cells {
			if i == 0 {
				r.Typ = v
			}
			switch fmt.Sprintf("$%d", i) {
			case r.Value:
				r.Value = v
			case r.Proxy:
				r.Proxy = v
			default:
			}
		}
		fmt.Println(lnum, (&r).String())
		rules = append(rules, &r)
	}
	handle := next
	for i := len(rules) - 1; i >= 0; i-- {
		r := rules[i]
		//TODO check Proxy
		//if !proxyName[rule.Proxy] {
		//	err = errors.Errorf("rule:[%s, %s, %s, %v], proxy:[%s] not found",
		//		rule.Typ, rule.Value, rule.Proxy, rule.Params, rule.Proxy)
		//	return
		//}
		handle, err = Get(r.Typ, r, handle, dnsHandle)
		if err != nil {
			logrus.WithError(err).WithField("type", KeyRuleSet).WithField("url", url).
				Error("init rule set failed")
		}
	}
	return handle, err
}
