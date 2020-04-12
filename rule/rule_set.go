package rule

import (
	"bufio"
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sipt/shuttle/dns"
	"github.com/sirupsen/logrus"
)

const (
	KeyRuleSet        = "RULE-SET"
	ParamsKeyInterval = "interval"

	defaultInterval = time.Hour * 24
)

func init() {
	Register(KeyRuleSet, ruleSetHandle)
}

func ruleSetHandle(ctx context.Context, rule *Rule, next Handle, dnsHandle dns.Handle) (handle Handle, err error) {
	var (
		etag      string
		modified  bool
		rules     []*Rule
		subHandle = next
		interval  time.Duration
	)
	if intervalStr, ok := rule.Params[ParamsKeyInterval]; ok && len(intervalStr) > 0 {
		interval, err = time.ParseDuration(intervalStr)
		if err != nil {
			err = errors.Wrapf(err, "[%s] params [%s:%s] is invalid", KeyRuleSet, ParamsKeyInterval, intervalStr)
		}
	} else {
		interval = defaultInterval
	}
	handle = func(ctx context.Context, info RequestInfo) *Rule {
		return subHandle(ctx, info)
	}
	go func() {
		var timer *time.Timer
		for {
			if timer == nil {
				timer = time.NewTimer(interval)
			} else {
				timer.Reset(interval)
				select {
				case <-timer.C:
				case <-ctx.Done():
					return
				}
			}
			rules, modified, etag, err = downloadRuleSet(ctx, rule, etag)
			if err != nil {
				logrus.WithError(err).WithField("url", rule.Value).
					Error("download rule set failed")
				continue
			}
			if !modified {
				logrus.WithField("url", rule.Value).Info("rule_set not change")
			}
			var reply, req Handle
			req = next
			for i := len(rules) - 1; i >= 0; i-- {
				r := rules[i]
				//TODO check Proxy
				//if !proxyName[rule.Proxy] {
				//	err = errors.Errorf("rule:[%s, %s, %s, %v], proxy:[%s] not found",
				//		rule.Typ, rule.Value, rule.Proxy, rule.Params, rule.Proxy)
				//	return
				//}
				reply, err = Get(ctx, r.Typ, r, req, dnsHandle)
				if err != nil {
					logrus.WithError(err).WithField("type", KeyRuleSet).WithField("url", rule.Value).
						Error("init rule set failed")
				}
				if reply != nil {
					req = reply
				}
			}
			subHandle = req
			logrus.WithField("url", rule.Value).Info("rule_set update success")
		}
	}()
	return
}

func downloadRuleSet(ctx context.Context, rule *Rule, reqETag string) (rules []*Rule, modified bool, etag string, err error) {
	url := rule.Value
	modified = true
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		err = errors.Wrapf(err, "download rule set failed: %s", url)
		return
	}
	if len(reqETag) > 0 {
		req.Header.Set("If-None-Match", etag)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		err = errors.Wrapf(err, "download rule set failed: %s", url)
		return
	}
	if resp.StatusCode == 304 {
		modified = false
		etag = reqETag
		return
	} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
		err = errors.Errorf("download rule set failed[%d]: %s", resp.StatusCode, url)
		return
	}
	defer resp.Body.Close()
	r := bufio.NewReader(resp.Body)
	var (
		bs   []byte
		lnum int
		line string
	)
	rules = make([]*Rule, 0, 64)
	for {
		lnum++
		bs, _, err = r.ReadLine()
		if err != nil {
			if err == io.EOF {
				err = nil
				break
			}
			err = errors.Errorf("read line failed, rule set failed: %s:[%d]", url, lnum)
			return
		}
		line = string(bs)
		line = strings.TrimSpace(line)
		if len(line) == 0 || line[0] == '#' {
			continue
		}
		cells := strings.Split(line, ",")
		if len(cells) >= 2 {
			rules = append(rules, &Rule{
				Profile: rule.Profile,
				Proxy:   rule.Proxy,
				Typ:     cells[0],
				Value:   cells[1],
			})
		} else {
			logrus.WithField("line", lnum).WithField("url", url).
				Errorf("rule set parse line failed: %s", line)
		}
	}
	etag = resp.Header.Get("ETag")
	return
}
