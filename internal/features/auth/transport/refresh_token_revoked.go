package transport

import (
	"encoding/json"
	"io"
	"net/http"
)

func (h *AuthHandler) RefreshTokenRevokedHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		RefreshToken *string `json:"refresh_token"`
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	defer r.Body.Close()

	if !json.Valid(body) {
		http.Error(w, "invalid JSON format", http.StatusUnauthorized)
	}

	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	response, err := h.authService.Logout(*request.RefreshToken)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(response)
}
