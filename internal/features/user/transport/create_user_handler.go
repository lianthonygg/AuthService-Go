package transport

import (
	"encoding/json"
	"net/http"

	"auth-service/internal/features/user/validate"
	"auth-service/internal/shared"
)

func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var user validate.CreateUserRequest

	if err := validate.DecodeAndValidateJSON(r, &user); err != nil {
		shared.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	created, err := h.service.CreateUser(ctx, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}
