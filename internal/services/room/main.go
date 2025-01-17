package servicesRoom

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log/slog"
	modelsMovie "match/internal/models/movie"
	modelsRoom "match/internal/models/room"
	storageRoom "match/internal/storage/room"
)

type RoomService struct {
	rsP *storageRoom.RoomStorage
	msP MovieStorageProvider
	l   *slog.Logger
}

type MovieStorageProvider interface {
	GetMoviesByGenre(ctx context.Context, genre string, page int) ([]modelsMovie.Movie, error)
}

type RoomStorageProvider interface {
	Add(ctx context.Context, genre string, userIDs []primitive.ObjectID) (*modelsRoom.Room, error)
	AddLike(ctx context.Context, roomID primitive.ObjectID, imdbID string) error
	GetRoom(ctx context.Context, roomID primitive.ObjectID) (*modelsRoom.Room, error)
	GetMatches(ctx context.Context, id primitive.ObjectID) ([]string, error)
}

func NewRoomService(rsP *storageRoom.RoomStorage, msP MovieStorageProvider, l *slog.Logger) *RoomService {
	return &RoomService{
		rsP: rsP,
		msP: msP,
		l:   l,
	}
}
