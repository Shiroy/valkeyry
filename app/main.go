package main

import (
	"github.com/codecrafters-io/redis-starter-go/app/commands"
	"github.com/codecrafters-io/redis-starter-go/app/memory"
	"github.com/codecrafters-io/redis-starter-go/app/network"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	fx.New(
		fx.Provide(
			network.NewServer,
			zap.NewExample,
		),
		commands.Module,
		memory.Module,
		fx.Invoke(func(server *network.Server) {}),
	).Run()
}
