package storageUser

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	apperr "match/internal/errors"
	models "match/internal/models/user"
	storage "match/internal/storage/mongo"
)

func (s *UserStorage) Get(ctx context.Context, username string) (*models.User, error) {
	filter := bson.M{"username": username}

	var user models.User
	err := s.db.Collection(storage.UserCollection).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, apperr.AppError{
				Status:     404,
				BErrorText: err.Error(),
				UErrorText: "нет пользователя с таким userName",
			}
		}
		return nil, err
	}
	return &user, nil
}
