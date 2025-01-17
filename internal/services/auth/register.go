package servicesAuth

import (
	"context"
	"fmt"
	models "match/internal/models/user"
)

func (s *AuthService) Register(ctx context.Context, username, password string) (*models.User, error) {
	const op = "servicesAuth.Register"

	hashed, err := hashPassword(password)
	if err != nil {
		s.l.Error("ошибка при хешировании пароля", fmt.Errorf("%s: %v", op, err))
		return nil, err
	}

	user, err := s.usP.Add(ctx, username, hashed)
	if err != nil {
		s.l.Error("ошибка при регистрации", fmt.Errorf("%s: %v", op, err))
		return nil, err
	}

	s.l.Info("регистрация прошла успешно")

	return user, nil
}
