package handlersAuth

import (
	"encoding/json"
	"net/http"
)

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	type reqBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var body reqBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.asP.Login(r.Context(), body.Username, body.Password)
	if err != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"token": token,
	})
}
