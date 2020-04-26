package main

import "C"
import (
	"encoding/json"
	"strings"

	"github.com/sipt/shuttle/cmd"
	"github.com/sipt/shuttle/dns"
	"github.com/sipt/shuttle/global/namespace"
	"github.com/sipt/shuttle/pkg/close"
)

func main() {
}

type config struct {
	ConfigFile  string `json:"config_file"`
	GeoipFile   string `json:"geoip_file"`
	RuntimeFile string `json:"runtime_file"`
}

//export start_shuttle
func start_shuttle(confStr *C.char) *C.char {
	conf := &config{}
	err := json.Unmarshal([]byte(C.GoString(confStr)), conf)
	if err != nil {
		return C.CString("error: " + err.Error())
	}
	if len(conf.ConfigFile) == 0 {
		return C.CString("error: config_path is empty")
	}
	if len(conf.GeoipFile) == 0 {
		return C.CString("error: geoip_path is empty")
	}
	if len(conf.RuntimeFile) == 0 {
		return C.CString("error: runtime_file is empty")
	}
	index := strings.LastIndex(conf.ConfigFile, ".")
	if index < 0 {
		return C.CString("error: config_path ext not found")
	}
	*cmd.Encoding = conf.ConfigFile[index+1:]
	*cmd.Path = conf.ConfigFile
	*dns.GeoipPath = conf.GeoipFile
	*cmd.RuntimePath = conf.RuntimeFile
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
