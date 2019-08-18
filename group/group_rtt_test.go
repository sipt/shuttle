package group

import (
	"context"
	"testing"
	"time"

	"github.com/sipt/shuttle/conf/logger"

	"github.com/sirupsen/logrus"

	"github.com/sipt/shuttle/server"

	"github.com/stretchr/testify/assert"
)

func TestRttGroup(t *testing.T) {
	logger.ConfigLogger()
	logrus.Debug("hello world")
	ctx := context.Background()
	_, err := Get(ctx, TypRTT, "test", map[string]string{
		ParamsKeyTestURL:  "www.baidu.com",
		ParamsKeyInterval: "1q",
	})
	assert.Error(t, err)
	_, err = Get(ctx, TypRTT, "test", map[string]string{
		ParamsKeyTestURL:  "http://www.baidu.com",
		ParamsKeyInterval: "1q",
	})
	assert.Error(t, err)
	group, err := Get(ctx, TypRTT, "test", map[string]string{
		ParamsKeyTestURL:  "http://www.baidu.com",
		ParamsKeyInterval: "10s",
	})
	assert.NoError(t, err)

	direct, err := server.Get(server.Direct, "", "", "", nil)
	assert.NoError(t, err)
	reject, err := server.Get(server.Reject, "", "", "", nil)
	assert.NoError(t, err)
	group.Append([]IServerX{
		&serverx{direct},
		&serverx{reject},
	})
	time.Sleep(time.Second)
	assert.EqualValues(t, group.Server().Name(), server.Direct)
}
