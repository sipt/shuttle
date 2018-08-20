package util

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"encoding/json"
	"net"
	"context"
)

var ipCache = make(map[string]*IP)

type IPInfo struct {
	Code int `json:"code"`
	Data IP  `json:"data`
}

type IP struct {
	IP        string `json:"ip"`
	Country   string `json:"country"`
	Area      string `json:"area"`
	Region    string `json:"region"`
	City      string `json:"city"`
	Isp       string `json:"isp"`
	CountryID string `json:"country_id"`
	AreaID    string `json:"area_id"`
	RegionID  string `json:"region_id"`
	CityID    string `json:"city_id"`
	IspID     string `json:"isp_id"`
}

func tabaoAPI(ip string) (*IPInfo, error) {
	tr := &http.Transport{
		DialContext: func(_ context.Context, network, addr string) (net.Conn, error) {
			return net.Dial(network, addr)
		},
	}
	client := &http.Client{Transport: tr}
	url := fmt.Sprintf("http://ip.taobao.com/service/getIpInfo.php?ip=%s", ip)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result IPInfo
	if err := json.Unmarshal(out, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func WatchIP(addr string) (*IP, error) {
	if ip, ok := ipCache[addr]; ok {
		return ip, nil
	}
	reply, err := tabaoAPI(addr)
	if err != nil {
		return nil, err
	}
	if reply != nil {
		ipCache[addr] = &(reply.Data)
		return &(reply.Data), nil
	}
	return nil, nil
}
