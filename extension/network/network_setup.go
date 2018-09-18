// +build !windows
// +build !darwin

package network

func EnableSystemProxy(host, port string) {
}

func DisableSystemProxy() {
}

func WebProxySwitch(status bool, args ...string) error {
	return nil
}

func SecureWebProxySwitch(status bool, args ...string) error {
	return nil
}

func SocksProxySwitch(status bool, args ...string) error {
	return nil
}
