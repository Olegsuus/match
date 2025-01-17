package handlers

import (
	handlersAuth "match/internal/controllers/rest/handlers/auth"
	handlersRoom "match/internal/controllers/rest/handlers/room"
)

type Handlers struct {
	Auth *handlersAuth.AuthHandler
	Room *handlersRoom.RoomHandler
}
