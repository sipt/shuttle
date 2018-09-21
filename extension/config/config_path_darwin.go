// +build darwin

package config

import (
	"os/user"
	"os"
)

func HomePath() (string, error) {
	u, err := user.Current()
	if err == nil {
		return u.HomeDir, nil
	}
	home := os.Getenv("HOME")
	return home, nil
}
