package service

import (
	"github.com/kickplate/api/events"
	auth "github.com/kickplate/api/service/auth"
	"github.com/kickplate/api/service/email"
	"github.com/kickplate/api/service/generator"
	"github.com/kickplate/api/service/githubapp"
	"github.com/kickplate/api/service/organization"
	"github.com/kickplate/api/service/plate"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(events.NewEventEmitter),
	fx.Provide(email.NewService),
	fx.Provide(email.NewListener),
	fx.Invoke(func(l *email.Listener, e *events.EventEmitter) {
		l.Register(e)
	}),
	fx.Provide(githubapp.NewService),
	fx.Provide(auth.NewAuthService),
	fx.Provide(plate.NewPlateService),
	fx.Provide(organization.NewOrganizationService),
	fx.Provide(generator.NewGeneratorService),
)
