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
	f, err := HTTP.Open(path)
	if err != nil {
		return nil, err
	}
	return ioutil.ReadAll(f)
}
