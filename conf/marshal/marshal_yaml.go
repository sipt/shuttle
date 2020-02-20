package marshal

import (
	"bytes"

	"github.com/pkg/errors"
	"github.com/sipt/yaml"
)

func init() {
	Register("yaml", newYamlMarshal)
}

func newYamlMarshal(_ map[string]string) (IMarshal, error) {
	return &yamlMarshal{}, nil
}

type yamlMarshal struct{}

func (t *yamlMarshal) Marshal(entity interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := yaml.NewEncoder(buf).Encode(entity); err != nil {
		return nil, errors.Wrap(err, "marshal entity failed")
	}
	return buf.Bytes(), nil
}

func (t *yamlMarshal) UnMarshal(data []byte, entity interface{}) (interface{}, error) {
	if err := yaml.Unmarshal(data, entity); err != nil {
		return nil, errors.Wrap(err, "unmarshal entity failed")
	}
	return entity, nil
}
