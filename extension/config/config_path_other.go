// +build !darwin,!linux,!windows

package config

func HomePath() (string, error) {
	return ".", nil
}
