package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStorage(t *testing.T) {
	s := NewMemoryStorage(5)
	ms := s.(*memoryStorage)
	s.SetLength(6)
	assert.Equal(t, ms.maxLength, 6)

	var key1, key2 = "localhost", "remote"
	s.Put(key1, Record{ID: 1})
	s.Put(key1, Record{ID: 2})
	s.Put(key2, Record{ID: 3})
	assert.Equal(t, s.Count(key1), 2)
	rs := s.Get(key1)
	assert.EqualValues(t, rs[0].ID, 1)
	assert.EqualValues(t, rs[1].ID, 2)

	ks := s.Keys()
	if ks[0] != key1 && ks[0] != key2 {
		t.Errorf("[%s] not in [%s, %s]", ks[0], key1, key2)
	}
	if ks[1] != key1 && ks[1] != key2 {
		t.Errorf("[%s] not in [%s, %s]", ks[1], key1, key2)
	}

	s.Put(key1, Record{ID: 3})
	s.Put(key1, Record{ID: 4})
	s.Put(key1, Record{ID: 5})
	s.Put(key1, Record{ID: 6})
	s.Put(key1, Record{ID: 7})
	assert.EqualValues(t, len(s.Get(key1)), 6)
	assert.EqualValues(t, s.Get(key1)[0].ID, 2)

	s.Clear(key1)
	assert.Equal(t, len(s.Get(key1)), 0)

	s.Clear()
	assert.Equal(t, len(s.Keys()), 0)

}
