package marshal

import (
	"bytes"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

func init() {
	Register("toml", newTomlMarshal)
}

func newTomlMarshal(_ map[string]string) (IMarshal, error) {
	return &tomlMarshal{}, nil
}

type tomlMarshal struct{}

func (t *tomlMarshal) Marshal(entity interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(entity); err != nil {
		return nil, errors.Wrap(err, "marshal entity failed")
	}
	return buf.Bytes(), nil
}

func (t *tomlMarshal) UnMarshal(data []byte, model interface{}) (interface{}, error) {
	if err := toml.Unmarshal(data, model); err != nil {
		return nil, errors.Wrap(err, "unmarshal entity failed")
	}
	return model, nil
}
