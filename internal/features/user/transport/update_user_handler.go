package transport

import (
	"encoding/json"
	"net/http"
	"strings"

	"auth-service/internal/features/user/model"
)

func (h *UserHandler) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")

	if id == "" {
		http.Error(w, "Id is Required", http.StatusBadRequest)
	}

	var user model.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	updated, err := h.service.UpdateUser(id, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}
