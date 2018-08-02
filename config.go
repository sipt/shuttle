package shuttle

func InitConfig() error {
	ss := []*Server{
		{
			Name:     "🇯🇵Linode_b",
			Host:     "jp.b.cloudss.win",
			Port:     "13819",
			Method:   "rc4-md5",
			Password: "07071818w",
		}, {
			Name: PolicyDirect,
		}, {
			Name: PolicyReject,
		},
	}
	gs := []*ServerGroup{
		{
			//🇯🇵Linode_b = custom, jp.b.cloudss.win, 13819, rc4-md5, 07071818w, (null)
			Servers: []interface{}{
				ss[0],
			},
			Name:       "JP",
			SelectType: "select",
		},
	}
	err := InitServers(gs, ss)
	if err != nil {
		Logger.Error("InitServer failed: ", err)
	}
	rules := []*Rule{
		{
			Type:    RuleDomainSuffix,
			Value:   "google.com",
			Policy:  "JP",
			Options: nil,
			Comment: "",
		},
		{
			Type:    RuleDomainSuffix,
			Value:   "facebook.com",
			Policy:  "JP",
			Options: nil,
			Comment: "",
		},
	}
	InitRule(rules)
	return nil
}
