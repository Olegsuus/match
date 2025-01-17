package storageUser

import "go.mongodb.org/mongo-driver/mongo"

type UserStorage struct {
	db *mongo.Database
}

func NewUserStorage(db *mongo.Database) *UserStorage {
	return &UserStorage{
		db: db,
	}
}
