package config

import (
	"github.com/sipt/shuttle/log"
	"os"
	"path/filepath"
)

var HomeDir string
var ShuttleHomeDir string

func init() {
	var err error
	HomeDir, err = HomePath()
	if err != nil {
		log.Logger.Errorf("[Extension-Config] get home path failed: %s", err.Error())
	} else {
		HomeDir += string(os.PathSeparator)
	}
	ShuttleHomeDir = filepath.Join(HomeDir, "Documents", "shuttle")
}
