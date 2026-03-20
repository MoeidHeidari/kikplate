package routes

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewRoutes),
	fx.Provide(NewAuthRoutes),
	fx.Provide(NewHelloRoutes),
)

type Route interface {
	Setup()
}

type Routes []Route

func NewRoutes(
	helloRoutes HelloRoutes,
	authRoutes AuthRoutes,
) Routes {
	return Routes{
		helloRoutes,
		authRoutes,
	}
}

func (r Routes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
