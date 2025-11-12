package cache

import (
    "sync"
    "time"
)

// item with TTL
type item struct {
    v       any
    expires time.Time
}

type LocalCache struct {
    mu   sync.RWMutex
    data map[string]item
    ttl  time.Duration
}

func NewLocalCache(ttl time.Duration) *LocalCache {
    return &LocalCache{data: make(map[string]item), ttl: ttl}
}

func (c *LocalCache) Get(key string) (any, bool) {
    c.mu.RLock()
    it, ok := c.data[key]
    c.mu.RUnlock()
    if !ok || time.Now().After(it.expires) {
        if ok {
            c.mu.Lock()
            delete(c.data, key)
            c.mu.Unlock()
        }
        return nil, false
    }
    return it.v, true
}

func (c *LocalCache) Set(key string, v any) {
    c.mu.Lock()
    c.data[key] = item{v: v, expires: time.Now().Add(c.ttl)}
    c.mu.Unlock()
}

func (c *LocalCache) Delete(key string) {
    c.mu.Lock()
    delete(c.data, key)
    c.mu.Unlock()
}

func (c *LocalCache) Clear() {
    c.mu.Lock()
    c.data = make(map[string]item)
    c.mu.Unlock()
}

