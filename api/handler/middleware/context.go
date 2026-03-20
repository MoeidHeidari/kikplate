package middleware

import (
	"context"

	"github.com/google/uuid"
)

type contextKey string

const accountIDKey contextKey = "account_id"

func SetAccountID(ctx context.Context, id uuid.UUID) context.Context {
	return context.WithValue(ctx, accountIDKey, id)
}

func GetAccountID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(accountIDKey).(uuid.UUID)
	return id, ok
}
