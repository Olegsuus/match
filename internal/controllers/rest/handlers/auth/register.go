package handlersAuth

import (
	"encoding/json"
	"net/http"
)

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var body reqBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := h.asP.Register(r.Context(), body.Username, body.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "user created",
		"user_id": user.ID.Hex(),
	})
}
