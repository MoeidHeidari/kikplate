package bootstrap

import (
	"github.com/kickplate/api/handler/handlers"
	"github.com/kickplate/api/handler/routes"
	"github.com/kickplate/api/lib"
	"github.com/kickplate/api/repository/postgres"
	"github.com/kickplate/api/service"
	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	lib.Module,
	service.Module,
	routes.Module,
	handlers.Module,
	postgres.Module,
)
