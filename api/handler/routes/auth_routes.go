package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/kickplate/api/handler/handlers"
	"github.com/kickplate/api/handler/middleware"
	"github.com/kickplate/api/lib"
)

type AuthRoutes struct {
	logger  lib.Logger
	env     lib.Env
	handler handlers.AuthHandler
	mux     lib.RequestHandler
}

func NewAuthRoutes(
	logger lib.Logger,
	env lib.Env,
	handler handlers.AuthHandler,
	mux lib.RequestHandler,
) AuthRoutes {
	return AuthRoutes{
		logger:  logger,
		env:     env,
		handler: handler,
		mux:     mux,
	}
}

func (r AuthRoutes) Setup() {
	r.mux.Mux.Route("/auth", func(router chi.Router) {
		router.Post("/register", r.handler.Register)
		router.Get("/verify-email", r.handler.VerifyEmail)
		router.Post("/login", r.handler.LoginLocal)
		router.Post("/request-password-reset", r.handler.RequestPasswordReset)
		router.Post("/reset-password", r.handler.ResetPassword)
		router.Post("/github/app/webhook", r.handler.GitHubWebhook)
		router.Get("/github/app/callback", r.handler.GitHubCallback)
		router.Get("/{provider}/redirect", r.handler.OAuthRedirect)
		router.Get("/{provider}/callback", r.handler.OAuthCallback)
		router.Get("/providers", r.handler.Providers)
	})

	r.mux.Mux.Group(func(router chi.Router) {
		router.Use(middleware.RequireAuth)
		router.Get("/auth/github/connect", r.handler.GitHubConnect)
		router.Get("/me", r.handler.Me)
		router.Delete("/me", r.handler.DeleteMe)
		router.Patch("/me/profile", r.handler.UpdateProfile)
		router.Patch("/me/username", r.handler.SetUsername)
	})
}
