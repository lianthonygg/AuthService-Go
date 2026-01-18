package transport

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (h *UserHandler) GetByIdUserHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")

	if id == "" {
		http.Error(w, "Id is Required", http.StatusBadRequest)
	}

	user, err := h.service.GetUserById(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
