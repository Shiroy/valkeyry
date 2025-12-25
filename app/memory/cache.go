package memory

import "github.com/codecrafters-io/redis-starter-go/app/memory/values"

type Cache interface {
	Set(key string, value values.Value)
	SetWithExpiration(key string, value values.Value, expireAt int64)
	Get(key string) (values.Value, bool)
}
