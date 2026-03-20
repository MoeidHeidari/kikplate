package lib

import (
	"github.com/joho/godotenv"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewEnv),
	fx.Provide(GetLogger),
	fx.Provide(NewRequestHandler),
	fx.Provide(NewDatabase),
)

func init() {
	godotenv.Load()
}
