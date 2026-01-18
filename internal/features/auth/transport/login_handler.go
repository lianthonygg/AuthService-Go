package transport

import (
	"encoding/json"
	"net/http"

	"auth-service/internal/features/auth/validator"
	"auth-service/internal/shared"
)

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user validator.LoginRequest

	if err := validator.DecodeAndValidateJSON(r, &user); err != nil {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	login, err := h.authService.Login(*user.Email, *user.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
	}

	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(login)
}
