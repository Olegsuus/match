package storageRoom

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	apperr "match/internal/errors"
	models "match/internal/models/room"
	"match/internal/storage/mongo"
)

func (r *RoomStorage) Add(ctx context.Context, genre string, userIDs []primitive.ObjectID) (*models.Room, error) {
	room := models.Room{
		ID:          primitive.NewObjectID(),
		Genre:       genre,
		UserIDs:     userIDs,
		LikedMovies: []string{},
	}

	_, err := r.db.Collection(mongo.RoomCollection).InsertOne(ctx, room)
	if err != nil {
		return nil, apperr.AppError{
			BErrorText: err.Error(),
			UErrorText: "ошибка при добавлении комнаты",
		}
	}
	return &room, nil
}

func (r *RoomStorage) AddLike(ctx context.Context, roomID primitive.ObjectID, imdbID string) error {

	filter := bson.M{"_id": roomID}

	update := bson.M{"$addToSet": bson.M{"liked_movies": imdbID}}

	_, err := r.db.Collection(mongo.RoomCollection).UpdateOne(ctx, filter, update)
	return err
}
