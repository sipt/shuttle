// +build windows

package network

import (
	"os/exec"
	"bytes"
	"net"
	"github.com/sipt/shuttle/log"
)

//reg add "HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings" /v ProxyEnable /t REG_DWORD /d 1 /f
//reg add "HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings" /v ProxyOverride /t REG_SZ /d "<local>" /f
//reg add "HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings" /v ProxyServer /t REG_SZ /d "atlproxy.test.com:8080" /f
const (
	reg           = "reg"
	add           = "add"
	settingkey    = `HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`
	proxyenable   = "ProxyEnable"
	proxyoverride = "ProxyOverride"
	proxyserver   = "ProxyServer"
	regdword      = "REG_DWORD"
	regsz         = "REG_SZ"
)

func EnableSystemProxy(host, port string) {
	err := WebProxySwitch(true, host, port)
	if err != nil {
		log.Logger.Errorf("Enable WebProxy failed: %v", err)
	}
}

func DisableSystemProxy() {
	err := WebProxySwitch(false)
	if err != nil {
		log.Logger.Errorf("Disable WebProxy failed: %v", err)
	}
}

func WebProxySwitch(status bool, args ...string) (err error) {
	if status {
		// turn on web proxy
		_, err = Command(reg, add, settingkey, "/v", proxyenable, "/t", regdword, "/d", "1", "/f")
		if err != nil {
			return err
		}
		if len(args) >= 2 {
			_, err = Command(reg, add, settingkey, "/v", proxyserver, "/t", regsz, "/d", net.JoinHostPort(args[0], args[1]), "/f")
		}
	} else {
		// turn of web proxy
		_, err = Command(reg, add, settingkey, "/v", proxyenable, "/t", regdword, "/d", "0", "/f")
	}
	return
}

func SecureWebProxySwitch(status bool, args ...string) (err error) {
	return
}

func SocksProxySwitch(status bool, args ...string) (err error) {
	// not support
	return
}

func Command(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args ...)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	return out.String(), err
}
