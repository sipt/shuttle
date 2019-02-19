package conn

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPool(t *testing.T) {
	pool := NewPool()
	pool.Put(Connection{
		ClientConn: &MockConn{ID: 1, Closed: false},
		ServerConn: &MockConn{ID: 2, Closed: false},
	})
	pool.Put(Connection{
		ClientConn: &MockConn{ID: 3, Closed: false},
		ServerConn: &MockConn{ID: 4, Closed: false},
	}, "Chrome")
	pool.Put(Connection{
		ClientConn: &MockConn{ID: 5, Closed: false},
		ServerConn: &MockConn{ID: 6, Closed: false},
	}, "Chrome")
	pool.Put(Connection{
		ClientConn: &MockConn{ID: 7, Closed: false},
		ServerConn: &MockConn{ID: 8, Closed: false},
	})
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

	pool.Put(Connection{
		ClientConn: &MockConn{ID: 9, Closed: false},
		ServerConn: &MockConn{ID: 10, Closed: false},
	}, "QQ")
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
