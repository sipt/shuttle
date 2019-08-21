package main

import (
	"fmt"
	"net/url"
)

type A struct {
	Name string
}

func main() {
	u, err := url.Parse("udp://8.8.8.8")
	if err != nil {
		panic(err)
	}
	fmt.Println(u.Scheme, u.Hostname(), u.Port())
}
