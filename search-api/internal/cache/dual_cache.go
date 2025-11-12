package cache

import "time"

// DualCache combines local and distributed caches; here distributed is optional
type DualCache struct {
    local *LocalCache
    // distributed placeholder; can be implemented with Memcached
}

func NewDual(localTTL time.Duration) *DualCache {
    return &DualCache{local: NewLocalCache(localTTL)}
}

func (d *DualCache) Get(key string) (any, bool) { return d.local.Get(key) }
func (d *DualCache) Set(key string, v any)      { d.local.Set(key, v) }
func (d *DualCache) Delete(key string)          { d.local.Delete(key) }
func (d *DualCache) Clear()                     { d.local.Clear() }

