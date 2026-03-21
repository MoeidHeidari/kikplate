package service

import (
	auth "github.com/kickplate/api/service/auth"
	"github.com/kickplate/api/service/plate"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(auth.NewAuthService),
	fx.Provide(plate.NewPlateService),
)
