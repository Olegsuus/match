package handlersRoom

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

func (h *RoomHandler) LikeMovie(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		RoomID string `json:"room_id"`
		ImdbID string `json:"imdb_id"`
	}
	var body reqBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	roomObjID, err := primitive.ObjectIDFromHex(body.RoomID)
	if err != nil {
		http.Error(w, "invalid room_id", http.StatusBadRequest)
		return
	}

	if err := h.rmP.LikeMovie(r.Context(), roomObjID, body.ImdbID); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte(`{"message": "liked"}`))
}
