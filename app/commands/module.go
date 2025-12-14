package commands

import "go.uber.org/fx"

func AsCommand(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(Command)),
		fx.ResultTags(`group:"commands"`),
	)
}

var Module = fx.Module("commands",
	fx.Provide(
		AsCommand(NewEcho),
		AsCommand(NewPing),
		AsCommand(NewSet),
		AsCommand(NewGet),
	),
)
