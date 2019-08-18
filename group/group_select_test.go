package group

import (
	"context"
	"testing"

	"github.com/sipt/shuttle/server"
	"github.com/stretchr/testify/assert"
)

func TestSelectGroup(t *testing.T) {
	ctx := context.Background()
	group, err := Get(ctx, TypSelect, "test", nil)
	assert.NoError(t, err)

	group.Append([]IServerX{
		&serverx{&server.DirectServer{}},
		&serverx{&server.RejectServer{}},
	})
	assert.EqualValues(t, group.Server().Name(), server.Direct)

	err = group.Select(server.Reject)
	assert.NoError(t, err)
	assert.EqualValues(t, group.Server().Name(), server.Reject)

	err = group.Select("hahaha")
	assert.Error(t, err)
}
