package marshal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYamlMarshal_Marshal(t *testing.T) {
	m, _ := newYamlMarshal(nil)
	data, err := m.Marshal(config)
	assert.NoError(t, err)

	config2, err := m.UnMarshal(data)
	assert.NoError(t, err)
	assert.EqualValues(t, config, config2)
}
