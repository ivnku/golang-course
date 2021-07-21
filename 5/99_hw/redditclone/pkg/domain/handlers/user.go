package handlers

import (
	"encoding/json"
	"net/http"
	"redditclone/pkg/auth"
	"redditclone/pkg/domain/models"
	"redditclone/pkg/domain/repositories"
	"redditclone/pkg/helpers"
)

type UsersHandler struct {
	UsersRepository repositories.UsersRepository
}

/**
 * @Description: Register a new user
 * @receiver h
 * @param w
 * @param r
 */
func (h *UsersHandler) Register(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	type formData struct {
		Login    string `json:"username"`
		Password string `json:"password"`
	}

	var data formData
	err := decoder.Decode(&data)
	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't decode request data!")
		return
	}

	user, err := h.UsersRepository.GetByName(data.Login)

	if user != nil {
		helpers.JsonError(w, http.StatusBadRequest, "User with the such name already exists!")
		return
	}

	passwordHash, err := auth.HashPassword(data.Password)

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't hash the password!")
		return
	}

	userToRegister := &models.User{Name: data.Login, Password: passwordHash}
	userId, err := h.UsersRepository.Create(userToRegister)

	if err != nil {
		helpers.JsonError(w, http.StatusInternalServerError, "Couldn't register user!")
		return
	}

	tokenString, err := auth.GenerateJWT(data.Login, userId)

	if err != nil {
		helpers.JsonError(w, http.StatusInternalServerError, "Couldn't create tokenString: "+err.Error())
		return
	}

	resp, _ := json.Marshal(map[string]interface{}{"token": tokenString})

	w.Write(resp)
}

/**
 * @Description: Authenticate a user
 * @receiver h
 * @param w
 * @param r
 */
func (h *UsersHandler) Auth(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	type formData struct {
		Login    string `json:"username"`
		Password string `json:"password"`
	}

	var data formData
	err := decoder.Decode(&data)
	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't decode request data!")
		return
	}

	user, err := h.UsersRepository.GetByName(data.Login)

	if user == nil {
		helpers.JsonError(w, http.StatusBadRequest, "User doesn't exist!")
		return
	}

	tokenString, err := auth.Auth(user.ID, user.Name, user.Password)

	if err != nil {
		helpers.JsonError(w, http.StatusInternalServerError, "Couldn't Authenticate user: "+err.Error())
		return
	}

	resp, _ := json.Marshal(map[string]interface{}{"token": tokenString})

	w.Write(resp)
}

/**
 * @Description: List all users
 * @receiver h
 * @param w
 * @param r
 */
func (h *UsersHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.UsersRepository.List()

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
