package group

import (
	"context"
	"testing"
	"time"

	"github.com/sipt/shuttle/server"
	"github.com/stretchr/testify/assert"
)

func TestRttGroup(t *testing.T) {
	ctx := context.Background()
	_, err := Get(ctx, TypRTT, "test", map[string]string{
		ParamsKeyTestURI:  "www.baidu.com",
		ParamsKeyInterval: "1q",
	})
	assert.Error(t, err)
	_, err = Get(ctx, TypRTT, "test", map[string]string{
		ParamsKeyTestURI:  "http://www.baidu.com",
		ParamsKeyInterval: "1q",
	})
	assert.Error(t, err)
	group, err := Get(ctx, TypRTT, "test", map[string]string{
		ParamsKeyTestURI:  "http://www.baidu.com",
		ParamsKeyInterval: "10s",
	})
	assert.NoError(t, err)

	direct, err := server.Get(server.Direct, "", "", 0, nil, nil)
	assert.NoError(t, err)
	direct = server.NewRttServer(direct, map[string]string{})
	reject, err := server.Get(server.Reject, "", "", 0, nil, nil)
	assert.NoError(t, err)
	group.Append([]IServerX{
		&serverx{direct},
		&serverx{reject},
	})
	time.Sleep(time.Second)
	assert.EqualValues(t, group.Server().Name(), server.Direct)
}
