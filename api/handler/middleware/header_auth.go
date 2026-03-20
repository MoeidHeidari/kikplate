package middleware

import (
	"net/http"

	"github.com/kickplate/api/lib"
	"github.com/kickplate/api/repository"
)

func HeaderAuth(
	env lib.Env,
	accountRepo repository.AccountRepository,
	logger lib.Logger,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if env.AuthHeader == "" {
				next.ServeHTTP(w, r)
				return
			}

			headerValue := r.Header.Get(env.AuthHeader)
			if headerValue == "" {
				next.ServeHTTP(w, r)
				return
			}

			account, err := accountRepo.GetByProvider(r.Context(), "header", headerValue)
			if err != nil {
				logger.Errorf("header auth lookup failed: %v", err)
				next.ServeHTTP(w, r)
				return
			}
			if account == nil {
				next.ServeHTTP(w, r)
				return
			}

			ctx := SetAccountID(r.Context(), account.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
