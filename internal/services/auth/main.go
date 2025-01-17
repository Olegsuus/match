package servicesAuth

import (
	"context"
	"log/slog"
	models "match/internal/models/user"
)

type AuthService struct {
	usP       UserStorageProvider
	jwtSecret string
	l         *slog.Logger
}

type UserStorageProvider interface {
	Add(ctx context.Context, username, hashedPassword string) (*models.User, error)
	Get(ctx context.Context, username string) (*models.User, error)
}

func NewAuthService(usP UserStorageProvider, secret string, l *slog.Logger) *AuthService {
	return &AuthService{
		usP:       usP,
		jwtSecret: secret,
		l:         l,
	}
}
