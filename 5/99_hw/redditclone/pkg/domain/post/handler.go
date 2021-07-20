package post

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"redditclone/configs"
	"redditclone/pkg/domain/comment"
	"redditclone/pkg/helpers"
	"strconv"
	"time"
)

type Handler struct {
	Repository Repository
	CommentsRepo comment.Repository
}

/**
 * @Description: Get one post by its id
 * @receiver h
 * @param w
 * @param r
 */
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 0)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	post, err := h.Repository.Get(uint(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	postSerialized, err := json.Marshal(post)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(postSerialized)
}

/**
 * @Description: Get the list of all posts
 * @receiver h
 * @param w
 * @param r
 */
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

	w.Write(postsSerialized)
}

/**
 * @Description: Create a post
 * @receiver h
 * @param w
 * @param r
 */
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	post := &Post{}

	err := decoder.Decode(post)
	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't decode request data!")
		return
	}

	ctx := r.Context()
	user := ctx.Value(configs.UserCtx).(map[string]string)

	userId, err := strconv.ParseUint(user["id"], 10, 0)

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't convert userId to uint!")
		return
	}

	post.UserID = uint(userId)
	post.User.ID = uint(userId)
	post.User.Name = user["username"]
	post.CreatedAt = time.Now().Format("2006-01-02 15:04:05")

	post, err = h.Repository.Create(post)

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't create the post!")
		return
	}

	response, err := json.Marshal(post)

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't serialize the post!")
		return
	}

	w.Write(response)
}

/**
 * @Description: Get the list of all posts
 * @receiver h
 * @param w
 * @param r
 */
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 0)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	success, err := h.Repository.Delete(uint(id))

	type Message struct {
		Message string `json:"message"`
	}

	var message string
	if success {
		message = "success"
	} else {
		message = "error"
	}

	response, err := json.Marshal(&Message{message})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(response)
}

/**
 * @Description: Create comment for a post
 * @receiver h
 * @param w
 * @param r
 */
func (h *Handler) Comment(w http.ResponseWriter, r *http.Request) {
	routeParams := mux.Vars(r)
	decoder := json.NewDecoder(r.Body)

	postComment := &comment.Comment{}

	reqComment := &struct{
		Comment string `json:"comment"`
	}{}

	err := decoder.Decode(reqComment)
	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't decode request data!")
		return
	}

	ctx := r.Context()
	user := ctx.Value(configs.UserCtx).(map[string]string)

	userId, err := strconv.ParseUint(user["id"], 10, 0)

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't convert userId to uint!")
		return
	}

	postId, err := strconv.ParseUint(routeParams["id"], 10, 0)

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't convert postId to uint!")
		return
	}

	postComment.UserID = uint(userId)
	postComment.Body = reqComment.Comment
	postComment.Created = time.Now().Format("2006-01-02 15:04:05")
	postComment.PostID = uint(postId)

	postComment, err = h.CommentsRepo.Create(postComment)

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't create a comment!")
		return
	}

	post, err := h.Repository.Get(uint(postId))

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't get the post!")
		return
	}

	postSerialized, err := json.Marshal(post)

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't serialize the post!")
		return
	}

	w.Write(postSerialized)
}
