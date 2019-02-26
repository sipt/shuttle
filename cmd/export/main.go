package main

import (
	"encoding/json"
	"github.com/sipt/shuttle/cmd"
)
import "C"

type Config struct {
	LogMode    string `json:"log_mode"`
	LogPath    string `json:"log_path"`
	ConfigPath string `json:"config_path"`
}

//export Run
func Run(confStr string) int {
	conf := &Config{}
	err := json.Unmarshal([]byte(confStr), conf)
	if err != nil {
		return 1
	}
	go cmd.Run(conf.LogMode, conf.LogPath, conf.ConfigPath)
	return 0
}

//export Shutdown
func Shutdown() {
	cmd.Shutdown("manual")
}

//export ReloadConfig
func ReloadConfig(configPath string) int {
	return cmd.ReloadConfig(configPath)
}

func main() {
}
