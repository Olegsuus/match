package handlersRoom

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (h *RoomHandler) CreateRoom(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		Genre   string   `json:"genre"`
		UserIDs []string `json:"user_ids"`
	}

	var body reqBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userObjIDs []primitive.ObjectID
	for _, u := range body.UserIDs {
		id, err := primitive.ObjectIDFromHex(u)
		if err == nil {
			userObjIDs = append(userObjIDs, id)
		}
	}

	room, err := h.rmP.Add(r.Context(), body.Genre, userObjIDs)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(body.UserIDs) > 1 {
		friendID := body.UserIDs[1]
		fromID := body.UserIDs[0]
		err := h.ws.SendInviteMessage(friendID, fromID, room.ID.Hex())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}

	json.NewEncoder(w).Encode(room)
}
