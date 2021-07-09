package user

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	Repository Repository
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.Repository.List()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	usersSerialized, err := json.Marshal(users)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(usersSerialized)
	return
}