package webhook

import (
	"testing"
	"time"
)

func TestMemCache(t *testing.T) {
	c := NewMemCache(10 * time.Millisecond)
	ok := c.Set("abc")
	if !ok {
		t.Fatal("set 'abc' failed")
	}

	ok = c.Set("abc")
	if ok {
		t.Fatal("set 'abc' success, should fail")
	}

	time.Sleep(10 * time.Millisecond)
	ok = c.Set("abc")
	if !ok {
		t.Fatal("set 'abc' failed after 10ms")
	}
}
