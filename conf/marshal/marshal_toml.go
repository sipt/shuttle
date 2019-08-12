package marshal

import (
	"bytes"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"github.com/sipt/shuttle/conf/model"
)

func init() {
	Register("toml", newTomlMarshal)
}

func newTomlMarshal(_ map[string]string) (IMarshal, error) {
	return &tomlMarshal{}, nil
}

type tomlMarshal struct{}

func (t *tomlMarshal) Marshal(config *model.Config) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := toml.NewEncoder(buf).Encode(config); err != nil {
		return nil, errors.Wrap(err, "marshal config failed")
	}
	return buf.Bytes(), nil
}

func (t *tomlMarshal) UnMarshal(data []byte) (*model.Config, error) {
	config := &model.Config{}
	if err := toml.Unmarshal(data, config); err != nil {
		return nil, errors.Wrap(err, "unmarshal config failed")
	}
	return config, nil
}
