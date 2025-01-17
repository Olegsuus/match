package handlersRoom

import (
	"encoding/json"
	"net/http"
)

func (h *RoomHandler) NextMovie(w http.ResponseWriter, r *http.Request) {
	genre := r.URL.Query().Get("genre")
	page := 1
	jsonMovies, err := h.rmP.GetNextMovie(r.Context(), genre, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(jsonMovies)
}
