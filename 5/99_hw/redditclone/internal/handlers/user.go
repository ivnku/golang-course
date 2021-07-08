package handlers

import (
	"encoding/json"
	"net/http"
	"redditclone/internal/repositories"
)

type UserHandler struct {}

func (uh *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	usersRepository := &repositories.UsersRepository{}

	users, err := usersRepository.List()

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