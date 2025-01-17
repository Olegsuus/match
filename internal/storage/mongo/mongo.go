package mongo

import (
	"context"
	"fmt"
	"log"
	"match/internal/config"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	UserCollection = "users"
	RoomCollection = "rooms"
)

type StorageMongo struct {
	Client         *mongo.Client
	DataBase       *mongo.Database
	UserCollection *mongo.Collection
	RoomCollection *mongo.Collection
}

func NewStorageMongo(cfg *config.Config) (*StorageMongo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.Mongo.Uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к MongoDB: %w", err)
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("ошибка проверки соединения с MongoDB: %w", err)
	}

	db := client.Database(cfg.Mongo.DBName)

	userColl := db.Collection(UserCollection)
	roomColl := db.Collection(RoomCollection)

	indexModel := mongo.IndexModel{
		Keys:    bson.M{"username": 1},
		Options: options.Index().SetUnique(true).SetName("unique_username"),
	}
	_, err = userColl.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания индекса на username: %w", err)
	}

	log.Println("Подключение к MongoDB установлено")

	return &StorageMongo{
		Client:         client,
		DataBase:       db,
		UserCollection: userColl,
		RoomCollection: roomColl,
	}, nil
}

func (s *StorageMongo) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return s.Client.Disconnect(ctx)
}
