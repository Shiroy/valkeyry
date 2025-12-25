package memory

import (
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/memory/values"
	"go.uber.org/zap"
)

type cacheEntry struct {
	value    values.Value
	expireAt int64
}

type CacheImpl struct {
	entries map[string]cacheEntry
	log     *zap.Logger
}

func NewCache(log *zap.Logger) *CacheImpl {
	return &CacheImpl{
		entries: make(map[string]cacheEntry),
		log:     log.With(zap.String("component", "Cache")),
	}
}

func (c *CacheImpl) Set(key string, value values.Value) {
	c.entries[key] = cacheEntry{
		value:    value,
		expireAt: 0,
	}
}

func (c *CacheImpl) SetWithExpiration(key string, value values.Value, expireAt int64) {
	c.entries[key] = cacheEntry{
		value:    value,
		expireAt: expireAt,
	}
}

func (c *CacheImpl) Get(key string) (values.Value, bool) {
	val, ok := c.entries[key]

	if ok {
		if val.expireAt > 0 && time.Now().UnixMilli() > val.expireAt {
			c.log.Debug("Expiring key", zap.String("key", key))
			delete(c.entries, key)
			return nil, false
		} else {
			return val.value, true
		}
	} else {
		return nil, false
	}
}

