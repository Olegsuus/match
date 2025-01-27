package handlersRoom

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"match/internal/controllers/ws"
	modelsMovie "match/internal/models/movie"
	modelsRoom "match/internal/models/room"
)

type RoomHandler struct {
	rmP RoomServicesProvider
	ws  *ws.WSHandler
}

type RoomServicesProvider interface {
	Add(ctx context.Context, genre string, userIDs []primitive.ObjectID) (*modelsRoom.Room, error)
	LikeMovie(ctx context.Context, roomID primitive.ObjectID, imdbID string) error
	GetMatches(ctx context.Context, roomID primitive.ObjectID) ([]string, error)
	GetNextMovie(ctx context.Context, genre string, page int) ([]modelsMovie.Movie, error)
	GetMoviesForRoom(ctx context.Context, roomID primitive.ObjectID, page int) ([]modelsMovie.Movie, error)
}

func NewRoomHandler(rmP RoomServicesProvider, ws *ws.WSHandler) *RoomHandler {
	return &RoomHandler{
		rmP: rmP,
		ws:  ws,
	}
}
