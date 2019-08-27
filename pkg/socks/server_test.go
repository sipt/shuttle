package socks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicAuthorization(t *testing.T) {
	b := []byte{0, 4, 'r', 'o', 'o', 't', 3, '1', '2', '3'}
	assert.False(t, BasicAuthorization("root", "123", b))
	b[0] = 1
	assert.False(t, BasicAuthorization("root", "1231", b))
	assert.False(t, BasicAuthorization("roo", "123", b))
	assert.True(t, BasicAuthorization("root", "123", b))
}
