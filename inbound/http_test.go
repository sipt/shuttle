package inbound

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/sipt/shuttle/conn"

	"github.com/stretchr/testify/assert"
)

func TestHttp(t *testing.T) {
	ctx := context.Background()
	f, err := newHTTPInbound(":20000", map[string]string{"auth_type": "basic", "user": "foo", "password": "bar"})
	assert.NoError(t, err)
	err = f(ctx, func(conn conn.ICtxConn) {
		for {
			data := make([]byte, 2048)
			n, err := conn.Read(data)
			str := string(data[:n])
			fmt.Println(str)
			if err != nil {
				t.Fatal(err)
			}
		}
	})
	assert.NoError(t, err)

}

func Test(t *testing.T) {
	fmt.Println(string("Basic c2lwdDoxMjMxMjM="[len("Basic "):]))
	data, err := base64.StdEncoding.DecodeString("c2lwdDoxMjMxMjM=")
	assert.NoError(t, err)
	fmt.Println(string(data))
}
