package storageRoom

import "go.mongodb.org/mongo-driver/mongo"

type RoomStorage struct {
	db *mongo.Database
}

func NewRoomStorage(db *mongo.Database) *RoomStorage {
	return &RoomStorage{
		db: db,
	}
}
