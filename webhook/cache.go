package webhook

import (
	"sync"
	"time"
)

type entry struct {
	nonce  string
	expire time.Time
	next   *entry
}

type cache struct {
	expire time.Duration
	mu     sync.Mutex
	entry  map[string]*entry
	first  *entry
	last   *entry
}

// NewMemCache用内存实现AuthOptions.Cache，可用于单实例执行推推webhook的应用。
func NewMemCache(expire time.Duration) Cache {
	return &cache{
		expire: expire,
		entry:  make(map[string]*entry),
	}
}

func (c *cache) Set(nonce string) bool {
	now := time.Now()

	c.mu.Lock()
	defer c.mu.Unlock()

	c.evict(now)

	if e := c.entry[nonce]; e != nil {
		return false
	}

	e := &entry{nonce: nonce, expire: now.Add(c.expire)}
	c.entry[nonce] = e
	if c.last != nil {
		c.last.next = e
	} else {
		c.first = e
	}
	c.last = e
	return true
}

func (c *cache) evict(now time.Time) {
	var cur, next *entry
	for cur = c.first; cur != nil; cur = next {
		next = cur.next
		if now.After(cur.expire) {
			delete(c.entry, cur.nonce)
			cur.next = nil
			continue
		}
		break
	}

	if cur != nil {
		c.first = cur
	} else {
		c.first = nil
		c.last = nil
	}
}
