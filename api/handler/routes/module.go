package routes

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewRoutes),
	fx.Provide(NewAuthRoutes),
	fx.Provide(NewHelloRoutes),
	fx.Provide(NewPlateRoutes),
	fx.Provide(NewBadgeRoutes),
	fx.Provide(NewOrganizationRoutes),
	fx.Provide(NewConfigRoutes),
)

type Route interface {
	Setup()
}

type Routes []Route

func NewRoutes(
	helloRoutes HelloRoutes,
	authRoutes AuthRoutes,
	plate PlateRoutes,
	badge BadgeRoutes,
	org OrganizationRoutes,
	config ConfigRoutes,
) Routes {
	return Routes{
		helloRoutes,
		authRoutes,
		plate,
		badge,
		org,
		config,
	}
}

func (r Routes) Setup() {
	for _, route := range r {
		route.Setup()
	}
}
