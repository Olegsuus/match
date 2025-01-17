package handlersRoom

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strconv"
)

func (h *RoomHandler) GetMatches(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room_id")
	roomObjID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		http.Error(w, "invalid room_id", http.StatusBadRequest)
		return
	}

	matches, err := h.rmP.GetMatches(r.Context(), roomObjID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(matches)
}

func (h *RoomHandler) GetMovies(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room_id")
	pageStr := r.URL.Query().Get("page")

	if roomID == "" {
		http.Error(w, "room_id is required", http.StatusBadRequest)
		return
	}

	roomObjID, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		http.Error(w, "invalid room_id", http.StatusBadRequest)
		return
	}

	page := 1
	if pageStr != "" {
		p, err := strconv.Atoi(pageStr)
		if err == nil && p > 0 {
			page = p
		}
	}

	movies, err := h.rmP.GetMoviesForRoom(r.Context(), roomObjID, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(movies)
}
