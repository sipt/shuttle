package dns

import "fmt"

func ReadHosts() {
	readHosts()
	fmt.Println("path", hosts.path)
	fmt.Println("byAddr", hosts.byAddr)
	fmt.Println("byName", hosts.byName)
	fmt.Println("expire", hosts.expire)
	fmt.Println("mtime", hosts.mtime)
	fmt.Println("size", hosts.size)
}
