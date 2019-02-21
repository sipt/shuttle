package conn

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {
	pool := NewPool()
	pool.Put(&MockConn{ID: 1, Closed: false}, &MockConn{ID: 2, Closed: false})
	pool.Put(&MockConn{ID: 3, Closed: false}, &MockConn{ID: 4, Closed: false}, "Chrome")
	pool.Put(&MockConn{ID: 5, Closed: false}, &MockConn{ID: 6, Closed: false}, "Chrome")
	pool.Put(&MockConn{ID: 7, Closed: false}, &MockConn{ID: 8, Closed: false})
	cs := pool.List("Chrome")
	assert.Equal(t, len(cs), 2)
	assert.Equal(t, len(pool.List()), 2)

	assert.EqualValues(t, cs[0].ClientConn.GetID(), 3)
	assert.EqualValues(t, cs[0].ServerConn.GetID(), 4)

	assert.EqualValues(t, cs[1].ClientConn.GetID(), 5)
	assert.EqualValues(t, cs[1].ServerConn.GetID(), 6)

	pool.Close(5, 6, "Chrome")
	assert.Equal(t, len(pool.List("Chrome")), 1)
	assert.Equal(t, len(pool.List()), 2)

	pool.Close(7, 8)
	assert.Equal(t, len(pool.List("Chrome")), 1)
	assert.Equal(t, len(pool.List()), 1)

	pool.Clear("Chrome")
	assert.Equal(t, len(pool.List("Chrome")), 0)
	assert.Equal(t, len(pool.List()), 1)

	pool.Replace(1, &MockConn{ID: 12, Closed: false})
	cs = pool.List()
	assert.EqualValues(t, cs[0].ClientConn.GetID(), 1)
	assert.EqualValues(t, cs[0].ServerConn.GetID(), 12)

	pool.Put(&MockConn{ID: 9, Closed: false}, &MockConn{ID: 10, Closed: false}, "QQ")
	pool.Clear()
	assert.Equal(t, len(pool.List("Chrome")), 0)
	assert.Equal(t, len(pool.List("QQ")), 0)
	assert.Equal(t, len(pool.List()), 0)
}

type MockConn struct {
	IConn
	ID     int64
	Closed bool
}

func (mc *MockConn) GetID() int64 {
	return mc.ID
}

func (mc *MockConn) Close() error {
	mc.Closed = true
	return nil
}

func (mc *MockConn) LocalAddr() net.Addr {
	return &net.TCPAddr{
		IP:   []byte{127, 0, 0, 1},
		Port: 8080,
	}
}

func (mc *MockConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{
		IP:   []byte{127, 0, 0, 2},
		Port: 8081,
	}
}
