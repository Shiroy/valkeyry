package memory

import "go.uber.org/fx"

var Module = fx.Module("memory", fx.Provide(fx.Annotate(NewCache, fx.As(new(Cache)))))
