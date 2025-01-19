package storageUser

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	apperr "match/internal/errors"
	models "match/internal/models/user"
	storage "match/internal/storage/mongo"
)

func (s *UserStorage) Add(ctx context.Context, username, hashedPassword string) (*models.User, error) {
	_, err := s.Get(ctx, username)
	if err == nil {
		return nil, apperr.AppError{
			Status:     409,
			BErrorText: apperr.ErrUserAlreadyExists,
			UErrorText: "Пользователь с таким userName уже создан",
		}
	}

	user := models.User{
		ID:       primitive.NewObjectID(),
		Username: username,
		Password: hashedPassword,
	}

	_, err = s.db.Collection(storage.UserCollection).InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
