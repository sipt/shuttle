// +build darwin

package network

import (
	"bytes"
	"github.com/sipt/shuttle/log"
	"os/exec"
	"strings"
)

const (
	on  = "on"
	off = "off"

	networksetup = "networksetup"

	setwebproxystate           = "-setwebproxystate"
	setwebproxy                = "-setwebproxy"
	getwebproxy                = "-getwebproxy"
	setsecurewebproxystate     = "-setsecurewebproxystate"
	getsecurewebproxy          = "-getsecurewebproxy"
	setsecurewebproxy          = "-setsecurewebproxy"
	setsocksfirewallproxystate = "-setsocksfirewallproxystate"
	getsocksfirewallproxy      = "-getsocksfirewallproxy"
	setsocksfirewallproxy      = "-setsocksfirewallproxy"

	listallnetworkservices = "-listallnetworkservices"
)

type networkSetupFunc func(name string) error

var networkServices = []string{"Wi-Fi", "Thunderbolt Bridge", "Thunderbolt Ethernet"}

func EnableSystemProxy(host, port string) {
	err := WebProxySwitch(true, host, port)
	if err != nil {
		log.Logger.Errorf("Enable WebProxy failed: %v", err)
	}
	err = SecureWebProxySwitch(true, host, port)
	if err != nil {
		log.Logger.Errorf("Enable SecureWeb failed: %v", err)
	}
	err = SocksProxySwitch(true, host, port)
	if err != nil {
		log.Logger.Errorf("Enable SocksProxy failed: %v", err)
	}
}

func DisableSystemProxy() {
	err := WebProxySwitch(false)
	if err != nil {
		log.Logger.Errorf("Disable WebProxy failed: %v", err)
	}
	err = SecureWebProxySwitch(false)
	if err != nil {
		log.Logger.Errorf("Disable SecureWebProxy failed: %v", err)
	}
	err = SocksProxySwitch(false)
	if err != nil {
		log.Logger.Errorf("Disable SocksProxy failed: %v", err)
	}
}

func listServices(callback networkSetupFunc) error {
	out, err := Command(networksetup, listallnetworkservices)
	if err != nil {
		return err
	}
	out = strings.TrimSpace(out)
	services := strings.Split(out, "\n")
	for _, v := range services {
		if !InServers(v) {
			continue
		}
		err = callback(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func Command(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args ...)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	return out.String(), err
}

func WebProxySwitch(status bool, args ...string) error {
	return listServices(
		func(name string) (err error) {
			if status {
				// turn on web proxy
				_, err = Command(networksetup, setwebproxystate, name, off)
				if err != nil {
					return err
				}
				if len(args) >= 2 {
					_, err = Command(networksetup, setwebproxy, name, args[0], args[1])
				}
			} else {
				// turn of web proxy
				_, err = Command(networksetup, setwebproxystate, name, off)
			}
			return err
		})
}

func SecureWebProxySwitch(status bool, args ...string) error {
	return listServices(
		func(name string) (err error) {
			if status {
				// turn on web proxy
				_, err = Command(networksetup, setsecurewebproxystate, name, off)
				if err != nil {
					return err
				}
				if len(args) >= 2 {
					_, err = Command(networksetup, setsecurewebproxy, name, args[0], args[1])
				}
			} else {
				// turn of web proxy
				_, err = Command(networksetup, setsecurewebproxystate, name, off)
			}
			return err
		})
}

func SocksProxySwitch(status bool, args ...string) error {
	return listServices(
		func(name string) (err error) {
			if status {
				// turn on web proxy
				_, err = Command(networksetup, setsocksfirewallproxystate, name, off)
				if err != nil {
					return err
				}
				if len(args) >= 2 {
					_, err = Command(networksetup, setsocksfirewallproxy, name, args[0], args[1])
				}
			} else {
				// turn of web proxy
				_, err = Command(networksetup, setsocksfirewallproxystate, name, off)
			}
			return err
		})
}

func InServers(v string) bool {
	for _, s := range networkServices {
		if v == s {
			return true
		}
	}
	return false
}
