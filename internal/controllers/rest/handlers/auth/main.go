package handlersAuth

import (
	"context"
	modelsUser "match/internal/models/user"
)

type AuthHandler struct {
	asP AuthServicesProvider
}

type AuthServicesProvider interface {
	Register(ctx context.Context, username, password string) (*modelsUser.User, error)
	Login(ctx context.Context, username, password string) (string, error)
}

func NewAuthHandler(asP AuthServicesProvider) *AuthHandler {
	return &AuthHandler{
		asP: asP,
	}
}
