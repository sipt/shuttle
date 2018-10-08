// +build windows

package config

import (
	"os/user"
	"os"
	"errors"
)

func HomePath() (string, error) {
	u, err := user.Current()
	if err == nil {
		return u.HomeDir, nil
	}
	var homeDir string
	homeDrive := os.Getenv("HOMEDRIVE")
	homePath := os.Getenv("HOMEPATH")
	if len(homeDrive) == 0 || len(homePath) == 0 {
		homeDir = os.Getenv("USERPROFILE")
	} else {
		homeDir = homeDir + homePath
	}
	if len(homeDir) == 0 {
		return "", errors.New("HomeDir is empty")
	}
	return homeDir, nil
}
