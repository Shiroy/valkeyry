package memory

type Cache struct {
	entries map[string]string
}

func NewCache() *Cache {
	return &Cache{entries: make(map[string]string)}
}

func (c *Cache) Set(key string, value string) {
	c.entries[key] = value
}

func (c *Cache) Get(key string) (string, bool) {
	val, ok := c.entries[key]
	return val, ok
}
