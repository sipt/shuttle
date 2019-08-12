package marshal

import (
	"bytes"

	"github.com/pkg/errors"
	"github.com/sipt/shuttle/conf/model"
	"github.com/sipt/yaml"
)

func init() {
	Register("yaml", newYamlMarshal)
}

func newYamlMarshal(_ map[string]string) (IMarshal, error) {
	return &yamlMarshal{}, nil
}

type yamlMarshal struct{}

func (t *yamlMarshal) Marshal(config *model.Config) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := yaml.NewEncoder(buf).Encode(config); err != nil {
		return nil, errors.Wrap(err, "marshal config failed")
	}
	return buf.Bytes(), nil
}

func (t *yamlMarshal) UnMarshal(data []byte) (*model.Config, error) {
	config := &model.Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, errors.Wrap(err, "unmarshal config failed")
	}
	return config, nil
}
