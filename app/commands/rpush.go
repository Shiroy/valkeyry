package commands

import (
	"container/list"

	"github.com/codecrafters-io/redis-starter-go/app/client"
	"github.com/codecrafters-io/redis-starter-go/app/memory"
	"github.com/codecrafters-io/redis-starter-go/app/memory/values"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type RPushParams struct {
	fx.In

	Cache *memory.Cache
	Log   *zap.Logger
}

type RPush struct {
	cache *memory.Cache
	log   *zap.Logger
}

func NewRPush(p RPushParams) *RPush {
	return &RPush{
		cache: p.Cache,
		log:   p.Log,
	}
}

func (r *RPush) Handle(c *client.Session, command []string) error {
	if len(command) < 2 {
		r.log.Info("Not enough argument")
		return c.SendErrorString("Not enough argument")
	}

	key := command[1]
	elements := command[2:]

	list := list.New()
	for _, value := range elements {
		list.PushBack(value)
	}

	r.cache.Set(key, values.NewValueList(list))

	return c.SendInteger(list.Len())
}

func (r *RPush) Mnemonic() string {
	return "RPUSH"
}
