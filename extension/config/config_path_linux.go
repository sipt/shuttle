// +build linux

package config

import (
	"os/user"
)

func HomePath() (string, error) {
	u, err := user.Current()
	if err == nil {
		return u.HomeDir, nil
	}
	return "", err
}
