package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"log/slog"
	"match/internal/config"
	handlersAuth "match/internal/controllers/rest/handlers/auth"
	handlersRoom "match/internal/controllers/rest/handlers/room"
	"match/internal/controllers/rest/routers"
	"match/internal/controllers/ws"
	"match/internal/logs"
	servicesAuth "match/internal/services/auth"
	servicesRoom "match/internal/services/room"
	"match/internal/storage/mongo"
	storageMovie "match/internal/storage/movie"
	storageRoom "match/internal/storage/room"
	storageUser "match/internal/storage/user"
	"net/http"
	"os"
)

type App struct {
	Config  *config.Config
	Log     *slog.Logger
	Mongo   *mongo.StorageMongo
	Router  *mux.Router
	Server  *http.Server
	LogFile *os.File
}

func NewApp(cfg *config.Config) (*App, error) {
	db, err := mongo.NewStorageMongo(cfg)
	if err != nil {
		slog.Error("ошибка при инициализации MongoDB", "error", err)
		return nil, err
	}

	logFile, err := logs.InitLogger(cfg.Env, cfg.Log.LogFile)
	if err != nil {
		slog.Error("ошибка при создании лог файла", "error", err)
		return nil, err
	}

	l := slog.Default()

	userStorage := storageUser.NewUserStorage(db.DataBase)
	roomStorage := storageRoom.NewRoomStorage(db.DataBase)
	movieStorage := storageMovie.NewMovieStorage(cfg.MoviesApi.Key, cfg.MoviesApi.Url)

	authServices := servicesAuth.NewAuthService(userStorage, cfg.JWT.Secret, l)
	roomServices := servicesRoom.NewRoomService(roomStorage, movieStorage, l)

	wsHandler := ws.NewWSHandler(roomStorage)
	authHandlers := handlersAuth.NewAuthHandler(authServices)
	roomHandlers := handlersRoom.NewRoomHandler(roomServices, wsHandler)

	router := routers.RegisterRoutes(authHandlers, roomHandlers, wsHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	app := &App{
		Config:  cfg,
		Log:     l,
		Mongo:   db,
		Router:  router,
		Server:  server,
		LogFile: logFile,
	}

	return app, nil
}

func (a *App) Run() error {
	a.Log.Info("Запуск сервера", "addr", a.Server.Addr)
	return a.Server.ListenAndServe()
}

func (a *App) Stop() error {
	a.Log.Info("Останавливаем Mongo соединение")
	if err := a.Mongo.Close(); err != nil {
		a.Log.Error("Ошибка при закрытии Mongo", "error", err)
		return err
	}
	if a.LogFile != nil {
		a.Log.Info("Закрываем лог файл")
		a.LogFile.Close()
	}
	return nil
}
