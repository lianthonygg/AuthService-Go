package transport

import (
	"encoding/json"
	"net/http"

	"auth-service/internal/features/user/model"
	"auth-service/internal/features/user/service"
)

type UserHandler struct {
	service service.UserService
}

func New(s service.UserService) *UserHandler {
	return &UserHandler{service: s}
}

func (h *UserHandler) UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		users, err := h.service.GetAllUsers()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)

	case http.MethodPost:
		var user model.User

		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		created, err := h.service.CreateUser(&user)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		w.WriteHeader(http.StatusCreated)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(created)
	}
}
