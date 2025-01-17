package storageRoom

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	models "match/internal/models/room"
	"match/internal/storage/mongo"
)

func (r *RoomStorage) GetRoom(ctx context.Context, roomID primitive.ObjectID) (*models.Room, error) {
	filter := bson.M{"_id": roomID}

	var room models.Room
	err := r.db.Collection(mongo.RoomCollection).FindOne(ctx, filter).Decode(&room)
	if err != nil {
		return nil, err
	}

	return &room, nil
}

func (r *RoomStorage) GetMatches(ctx context.Context, id primitive.ObjectID) ([]string, error) {
	room, err := r.GetRoom(ctx, id)
	if err != nil {
		return nil, err
	}

	return room.LikedMovies, nil
}
