// +build !release
//go:generate go run github.com/UnnoTed/fileb0x b0x.yaml

package assets

import (
	"io/ioutil"
	"net/http"
)

var HTTP http.FileSystem = http.Dir("./")

// ReadFile is adapTed from ioutil
func ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}
