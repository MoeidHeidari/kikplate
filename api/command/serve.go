package command

import (
	"fmt"
	"net/http"

	"github.com/kickplate/api/handler/middleware"
	"github.com/kickplate/api/handler/routes"
	"github.com/kickplate/api/lib"
	"github.com/kickplate/api/repository"
	"github.com/spf13/cobra"
)

type ServeCommand struct{}

func (c ServeCommand) Short() string {
	return "Start the server"
}

func (c ServeCommand) Setup(cmd *cobra.Command) {}

func (s *ServeCommand) Run() lib.CommandRunner {
	return func(
		env lib.Env,
		logger lib.Logger,
		handler lib.RequestHandler,
		accountRepo repository.AccountRepository,
		r routes.Routes,
	) {
		handler.Mux.Use(middleware.Authenticate(env, logger))
		handler.Mux.Use(middleware.HeaderAuth(env, accountRepo, logger))

		r.Setup()

		addr := fmt.Sprintf(":%s", env.ServerPort)
		logger.Info("Running server on port ", env.ServerPort)

		if err := http.ListenAndServe(addr, handler.Mux); err != nil {
			logger.Fatal("Server failed: ", err)
		}
	}
}

func NewServeCommand() *ServeCommand {
	return &ServeCommand{}
}
