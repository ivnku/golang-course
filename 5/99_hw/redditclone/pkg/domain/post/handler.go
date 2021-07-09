package post

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	Repository Repository
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	posts, err := h.Repository.List()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	postsSerialized, err := json.Marshal(posts)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _ = w.Write(postsSerialized)
	return
}