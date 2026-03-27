package handlers

import (
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(NewHelloHandler),
	fx.Provide(NewAuthHandler),
	fx.Provide(NewPlateHandler),
	fx.Provide(NewBadgeHandler),
	fx.Provide(NewOrganizationHandler),
	fx.Provide(NewConfigHandler),
)
