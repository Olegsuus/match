package servicesRoom

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	models "match/internal/models/room"
)

func (s *RoomService) Add(ctx context.Context, genre string, userIDs []primitive.ObjectID) (*models.Room, error) {
	const op = "servicesRoom.Add"

	room, err := s.rsP.Add(ctx, genre, userIDs)
	if err != nil {
		s.l.Error("ошибка при добавлении комнаты", fmt.Errorf("%s: %v", op, err))
		return nil, err
	}

	s.l.Info("комната успешно добавлена")

	return room, nil
}

func (s *RoomService) LikeMovie(ctx context.Context, roomID primitive.ObjectID, imdbID string) error {
	const op = "servicesRoom.LikeMovie"

	err := s.rsP.AddLike(ctx, roomID, imdbID)
	if err != nil {
		s.l.Error("ошибка при добавлении лайка на фильм", fmt.Errorf("%s: %v", op, err))
		return err
	}

	s.l.Info("лайк на фильм успешно добавлен")

	return nil
}
