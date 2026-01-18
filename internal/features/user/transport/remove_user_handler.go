package transport

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (h *UserHandler) RemoveUserHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")

	if id == "" {
		http.Error(w, "Id is Required", http.StatusBadRequest)
	}

	err := h.service.RemoveUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode("User Removed Success")
}
