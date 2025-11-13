package cache

import (
	"errors"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

// DistributedCache wraps a Memcached client with TTL and hit/miss stats.
type DistributedCache struct {
	client *memcache.Client
	ttl    time.Duration

	hits   uint64
	misses uint64
}

func NewDistributedCache(servers []string, ttl time.Duration) (*DistributedCache, error) {
	if len(servers) == 0 {
		return nil, errors.New("memcached: no servers configured")
	}
	c := memcache.New(servers...)
	return &DistributedCache{client: c, ttl: ttl}, nil
}

func (c *DistributedCache) Get(key string) ([]byte, bool) {
	if c == nil || c.client == nil {
		return nil, false
	}
	item, err := c.client.Get(key)
	if err != nil {
		c.misses++
		return nil, false
	}
	c.hits++
	return item.Value, true
}

func (c *DistributedCache) Set(key string, value []byte) {
	if c == nil || c.client == nil {
		return
	}
	_ = c.client.Set(&memcache.Item{
		Key:        key,
		Value:      value,
		Expiration: int32(c.ttl.Seconds()),
	})
}

func (c *DistributedCache) Delete(key string) {
	if c == nil || c.client == nil {
		return
	}
	_ = c.client.Delete(key)
}

func (c *DistributedCache) Stats() (hits, misses uint64) {
	return c.hits, c.misses
}
