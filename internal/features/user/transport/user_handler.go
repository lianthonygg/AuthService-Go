package transport

import (
	"encoding/json"
	"net/http"
	"strings"

	"auth-service/internal/features/user/model"
	"auth-service/internal/features/user/service"
	"auth-service/internal/features/user/validate"
	"auth-service/internal/shared"
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
		var user validate.CreateUserRequest

		if err := validate.DecodeAndValidateJSON(r, &user); err != nil {
			shared.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
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

func (h *UserHandler) UserHandlerById(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")

	if id == "" {
		http.Error(w, "Id is Required", http.StatusBadRequest)
	}

	switch r.Method {
	case http.MethodGet:
		user, err := h.service.GetUserById(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)

	case http.MethodPut:
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

	case http.MethodDelete:
		err := h.service.RemoveUser(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode("User Removed Success")
	}
}
