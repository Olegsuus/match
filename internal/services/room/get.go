package servicesRoom

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	models "match/internal/models/movie"
)

func (s *RoomService) GetNextMovie(ctx context.Context, genre string, page int) ([]models.Movie, error) {
	const op = "servicesRoom.GetNextMovie"

	movies, err := s.msP.GetMoviesByGenre(ctx, genre, page)
	if err != nil {
		s.l.Error("ошибка при получении следующего фильма по жанру", fmt.Errorf("%s: %s", op, err))
		return nil, err
	}

	s.l.Info("список фильмов успешно получен")

	return movies, nil
}

func (s *RoomService) GetMatches(ctx context.Context, roomID primitive.ObjectID) ([]string, error) {
	const op = "servicesRoom.GetMatches"

	matches, err := s.rsP.GetMatches(ctx, roomID)
	if err != nil {
		s.l.Error("ошибка при получении списка лайкнутых фильмов", fmt.Errorf("%s: %v", op, err))
	}

	s.l.Info("список лайкнутых фильмов успешно получен")

	return matches, nil
}

func (s *RoomService) GetMoviesForRoom(ctx context.Context, roomID primitive.ObjectID, page int) ([]models.Movie, error) {
	const op = "servicesRoom.GetMoviesForRoom"

	roomData, err := s.rsP.GetRoom(ctx, roomID)
	if err != nil {
		s.l.Error("не удалось получить комнату", fmt.Errorf("%s: %v", op, err))
		return nil, err
	}

	genre := roomData.Genre
	movies, err := s.msP.GetMoviesByGenre(ctx, genre, page)
	if err != nil {
		s.l.Error("ошибка при получении списка фильмов по жанру", fmt.Errorf("%s: %v", op, err))
		return nil, err
	}

	s.l.Info("список фильмов успешно получен для комнаты", "roomID", roomID.Hex(), "genre", genre)
	return movies, nil
}
