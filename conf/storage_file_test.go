package conf

import (
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestFileStorage(t *testing.T) {
	file := "config.json"
	s, err := newFileStorage(map[string]string{
		"path": file,
	})
	assert.NoError(t, err)

	// test save
	testData := []byte(`{"data": "test_data"}`)
	err = s.Save(testData)
	assert.NoError(t, err)

	// test load
	data, err := s.Load()
	assert.NoError(t, err)
	assert.EqualValues(t, string(data), testData)

	// test notify
	testData2 := []byte(`{"data": "test_data_2"}`)
	ctx, cancel := context.WithCancel(context.Background())
	err = s.RegisterNotify(ctx, func() {
		data, err := s.Load()
		assert.NoError(t, err)
		assert.EqualValues(t, string(data), testData2)
	})
	assert.NoError(t, err)

	err = s.Save(testData2)
	assert.NoError(t, err)
	cancel() // end notify

	// rm test file
	err = os.Remove(file)
	assert.NoError(t, err)
}
