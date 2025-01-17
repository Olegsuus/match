package servicesAuth

import (
	"context"
	"fmt"
)

func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	const op = "servicesAuth.Login"

	user, err := s.usP.Get(ctx, username)
	if err != nil {
		s.l.Error("ошибка при получении пользователя", fmt.Errorf("%s: %v", op, err))
		return "", err
	}

	if err := checkPassword(user.Password, password); err != nil {
		s.l.Error("ошибка при проверки пароля", fmt.Errorf("%s: %v", op, err))
		return "", err
	}

	tokenString, err := s.generateJWT(user.ID)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
