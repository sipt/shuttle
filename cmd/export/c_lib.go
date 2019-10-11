package main

import "C"
import (
	"github.com/sipt/shuttle/cmd"
	"github.com/sipt/shuttle/dns"
	"github.com/sipt/shuttle/global/namespace"
	"github.com/sipt/shuttle/pkg/close"
)

func main() {
}

//export start_shuttle
func start_shuttle(configPath, geoipPath *C.char) *C.char {
	*cmd.Encoding = "toml"
	*cmd.Path = C.GoString(configPath)
	*dns.GeoipPath = C.GoString(geoipPath)
	if err := cmd.Start(); err != nil {
		return C.CString("error: " + err.Error())
	}
	return C.CString(namespace.NamespaceWithName().Profile().Config().Controller.Addr)
}

//export stop_shuttle
func stop_shuttle() *C.char {
	if err := close.Close(true); err != nil {
		return C.CString(err.Error())
	}
	return C.CString("success")
}
