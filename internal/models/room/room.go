package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Room struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"ID"`
	Genre       string               `bson:"genre"         json:"Genre"`
	UserIDs     []primitive.ObjectID `bson:"user_ids"      json:"UserIDs"`
	LikedMovies []string             `bson:"liked_movies"  json:"LikedMovies"`
}
