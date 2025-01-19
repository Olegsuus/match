package routers

import (
	"github.com/gorilla/mux"
	handlersAuth "match/internal/controllers/rest/handlers/auth"
	handlersRoom "match/internal/controllers/rest/handlers/room"
	"match/internal/controllers/ws"
	"net/http"
)

func RegisterRoutes(authH *handlersAuth.AuthHandler, roomH *handlersRoom.RoomHandler, wsHandler *ws.WSHandler) *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/auth/register", authH.Register).Methods(http.MethodPost)
	router.HandleFunc("/auth/login", authH.Login).Methods(http.MethodPost)

	router.HandleFunc("/room", roomH.CreateRoom).Methods(http.MethodPost)
	router.HandleFunc("/room/next", roomH.NextMovie).Methods(http.MethodGet)
	router.HandleFunc("/room/like", roomH.LikeMovie).Methods(http.MethodPost)
	router.HandleFunc("/room/matches", roomH.GetMatches).Methods(http.MethodGet)
	router.HandleFunc("/room/movies", roomH.GetMovies).Methods(http.MethodGet)

	router.HandleFunc("/ws", wsHandler.HandleWSUpgrade)

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./front/")))

	return router
}
