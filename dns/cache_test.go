package dns

import (
	"testing"
	"time"
)

func TestCacheManager(t *testing.T) {
	m := NewCacheManager()
	m.Run()
	defer m.Stop()
	m.Push("a", time.Second)
	m.Push("b", 3*time.Second)
	m.Push("c", -1*time.Second)
	m.Range(func(data interface{}) (breaked bool) {
		s := data.(string)
		if s != "a" && s != "b" {
			t.Errorf("value failed.")
		}
		return false
	})
	time.Sleep(1100 * time.Millisecond)
	m.Range(func(data interface{}) (breaked bool) {
		s := data.(string)
		if s != "b" {
			t.Errorf("value failed. %s", s)
		}
		return false
	})
	time.Sleep(time.Second)
	m.Range(func(data interface{}) (breaked bool) {
		s := data.(string)
		if s != "b" {
			t.Errorf("value failed.")
		}
		return false
	})
	time.Sleep(time.Second)
	m.Range(func(data interface{}) (breaked bool) {
		t.Errorf("value failed.")
		return false
	})
}
