package memory

import (
	"time"

	"go.uber.org/zap"
)

type cacheEntry struct {
	value    string
	expireAt int64
}

type Cache struct {
	entries map[string]cacheEntry
	log     *zap.Logger
}

func NewCache(log *zap.Logger) *Cache {
	return &Cache{
		entries: make(map[string]cacheEntry),
		log:     log.With(zap.String("component", "Cache")),
	}
}

func (c *Cache) Set(key string, value string) {
	c.entries[key] = cacheEntry{
		value:    value,
		expireAt: 0,
	}
}

func (c *Cache) SetWithExpiration(key string, value string, expireAt int64) {
	c.entries[key] = cacheEntry{
		value:    value,
		expireAt: expireAt,
	}
}

func (c *Cache) Get(key string) (string, bool) {
	val, ok := c.entries[key]

	if ok {
		if val.expireAt > 0 && time.Now().UnixMilli() > val.expireAt {
			c.log.Debug("Expiring key", zap.String("key", key))
			delete(c.entries, key)
			return "", false
		} else {
			return val.value, true
		}
	} else {
		return "", false
	}
}
