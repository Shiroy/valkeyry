package commands

import (
	"github.com/codecrafters-io/redis-starter-go/app/client"
	"github.com/codecrafters-io/redis-starter-go/app/memory"
	"github.com/codecrafters-io/redis-starter-go/app/memory/values"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type GetParams struct {
	fx.In

	Cache *memory.Cache
	Log   *zap.Logger
}

type Get struct {
	cache *memory.Cache
	log   *zap.Logger
}

func (g *Get) Handle(c *client.Session, command []string) error {
	if len(command) < 1 {
		g.log.Info("Not enough argument")
		return c.SendErrorString("Not enough argument")
	}

	key := command[1]
	value, ok := g.cache.Get(key)

	if !ok {
		g.log.Debug("Not found", zap.String("key", key))
		return c.SendNullBulkString()
	}

	switch v := value.(type) {
	case values.ValueString:
		return c.SendString(v.Data)
	default:
		g.log.Debug("Unknown entry type.", zap.String("type", v.Kind().String()))
		return c.SendNullBulkString()
	}
}

func (g *Get) Mnemonic() string {
	return "GET"
}

func NewGet(p GetParams) *Get {
	return &Get{
		cache: p.Cache,
		log:   p.Log,
	}
}
