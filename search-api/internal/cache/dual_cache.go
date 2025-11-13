package cache

import "time"

// DualCache combines local in-process cache with an optional distributed backend.
type DualCache struct {
	local       *LocalCache
	distributed *DistributedCache
	codec       Codec
}

func NewDual(localTTL time.Duration, distributed *DistributedCache, codec Codec) *DualCache {
	return &DualCache{
		local:       NewLocalCache(localTTL),
		distributed: distributed,
		codec:       codec,
	}
}

func (d *DualCache) Get(key string) (any, bool) {
	if v, ok := d.local.Get(key); ok {
		return v, true
	}
	if d.distributed == nil || d.codec == nil {
		return nil, false
	}
	bytes, ok := d.distributed.Get(key)
	if !ok {
		return nil, false
	}
	value, err := d.codec.Unmarshal(bytes)
	if err != nil {
		return nil, false
	}
	d.local.Set(key, value)
	return value, true
}

func (d *DualCache) Set(key string, v any) {
	d.local.Set(key, v)
	if d.distributed == nil || d.codec == nil {
		return
	}
	bytes, err := d.codec.Marshal(v)
	if err != nil {
		return
	}
	d.distributed.Set(key, bytes)
}

func (d *DualCache) Delete(key string) {
	d.local.Delete(key)
	if d.distributed != nil {
		d.distributed.Delete(key)
	}
}

func (d *DualCache) Clear() {
	d.local.Clear()
}

func (d *DualCache) Stats() (localEntries int, distHits, distMisses uint64) {
	if d.local != nil {
		localEntries = len(d.local.data)
	}
	if d.distributed != nil {
		distHits, distMisses = d.distributed.Stats()
	}
	return
}
