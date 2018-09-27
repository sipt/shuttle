package config

import (
	"os"
	"path/filepath"
	"io/ioutil"
)

var HomeDir string
var ShuttleHomeDir string

func init() {
	var err error
	HomeDir, err = HomePath()
	if err != nil {
		ioutil.WriteFile("error.log", []byte(err.Error()), 0664)
		panic(err)
	} else {
		HomeDir += string(os.PathSeparator)
	}
	ShuttleHomeDir = filepath.Join(HomeDir, "Documents", "shuttle")
}
