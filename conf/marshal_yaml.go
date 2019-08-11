package conf

import (
	"bytes"
	"github.com/pkg/errors"
	"github.com/sipt/yaml"
)

func init() {
	RegisterMarshal("yaml", newYamlMarshal)
}

func newYamlMarshal(_ map[string]string) (IMarshal, error) {
	return &yamlMarshal{}, nil
}

type yamlMarshal struct{}

func (t *yamlMarshal) Marshal(config *Config) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := yaml.NewEncoder(buf).Encode(config); err != nil {
		return nil, errors.Wrap(err, "marshal config failed")
	}
	return buf.Bytes(), nil
}

func (t *yamlMarshal) UnMarshal(data []byte) (*Config, error) {
	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, errors.Wrap(err, "unmarshal config failed")
	}
	return config, nil
}
