package commands

import (
	"github.com/codecrafters-io/redis-starter-go/app/client"
	"github.com/codecrafters-io/redis-starter-go/app/memory"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type SetParams struct {
	fx.In
	Cache *memory.Cache
	Log   *zap.Logger
}

type Set struct {
	cache *memory.Cache
	log   *zap.Logger
}

func (s *Set) Handle(c *client.Session, command []string) error {
	if len(command) < 3 {
		s.log.Info("Not enough argument")
		return c.SendErrorString("Not enough argument")
	}

	key := command[1]
	value := command[2]

	s.cache.Set(key, value)

	return c.SendSimpleString("OK")
}

func (s *Set) Mnemonic() string {
	return "SET"
}

func NewSet(p SetParams) *Set {
	return &Set{
		log:   p.Log.With(zap.String("command", "SET")),
		cache: p.Cache,
	}
}
