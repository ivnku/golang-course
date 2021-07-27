package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"redditclone/configs"
	"redditclone/pkg/domain/models"
	"redditclone/pkg/domain/repositories"
	"redditclone/pkg/helpers"
	"strconv"
	"time"
)

type PostsHandler struct {
	PostsRepository    repositories.PostsRepository
	CommentsRepository repositories.CommentsRepository
	UsersRepository    repositories.UsersRepository
	Config             configs.Config
}

/**
 * @Description: Get one post by its id
 * @receiver h
 * @param w
 * @param r
 */
func (h *PostsHandler) Get(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 0)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	post, err := h.PostsRepository.Get(uint(id))

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// increment views each time a user open the post
	post.Views++
	post, err = h.PostsRepository.Update(post, []string{"views"})

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helpers.SerializeAndReturn(w, post)
}

/**
 * @Description: Get the list of all posts
 * @receiver h
 * @param w
 * @param r
 */
func (h *PostsHandler) List(w http.ResponseWriter, r *http.Request) {
	posts, err := h.PostsRepository.List()

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helpers.SerializeAndReturn(w, posts)
}

/**
 * @Description: Get posts within a certain category
 * @receiver h
 * @param w
 * @param r
 */
func (h *PostsHandler) CategoryList(w http.ResponseWriter, r *http.Request) {
	routeParams := mux.Vars(r)

	posts, err := h.PostsRepository.CategoryList(routeParams["categoryName"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helpers.SerializeAndReturn(w, posts)
}

/**
 * @Description: Get posts of a certain user
 * @receiver h
 * @param w
 * @param r
 */
func (h *PostsHandler) UserList(w http.ResponseWriter, r *http.Request) {
	routeParams := mux.Vars(r)

	user, err := h.UsersRepository.GetByName(routeParams["userName"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	posts, err := h.PostsRepository.UserList(user.ID)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	helpers.SerializeAndReturn(w, posts)
}

/**
 * @Description: Create a post
 * @receiver h
 * @param w
 * @param r
 */
func (h *PostsHandler) Create(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	post := &models.Post{}

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

	post, err = h.PostsRepository.Create(post)

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't create the post!")
		return
	}

	helpers.SerializeAndReturn(w, post)
}

/**
 * @Description: Get the list of all posts
 * @receiver h
 * @param w
 * @param r
 */
func (h *PostsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.ParseUint(params["id"], 10, 0)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	success, err := h.PostsRepository.Delete(uint(id))

	type Message struct {
		Message string `json:"message"`
	}

	var message string
	if success {
		message = "success"
	} else {
		message = "error"
	}

	helpers.SerializeAndReturn(w, &Message{message})
}

/**
 * @Description: Create comment for a post
 * @receiver h
 * @param w
 * @param r
 */
func (h *PostsHandler) Comment(w http.ResponseWriter, r *http.Request) {
	routeParams := mux.Vars(r)
	decoder := json.NewDecoder(r.Body)

	postComment := &models.Comment{}

	reqComment := &struct {
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

	postComment, err = h.CommentsRepository.Create(postComment)

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't create a comment!")
		return
	}

	post, err := h.PostsRepository.Get(uint(postId))

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't get the post!")
		return
	}

	helpers.SerializeAndReturn(w, post)
}

/**
 * @Description: Delete a post's comment
 * @receiver h
 * @param w
 * @param r
 */
func (h *PostsHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	routeParams := mux.Vars(r)

	commentId, err := strconv.ParseUint(routeParams["commentId"], 10, 0)

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't convert commentId to uint!")
		return
	}

	postId, err := strconv.ParseUint(routeParams["postId"], 10, 0)

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't convert postId to uint!")
		return
	}

	_, err = h.CommentsRepository.Delete(uint(commentId))

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't delete post comment!")
		return
	}

	post, err := h.PostsRepository.Get(uint(postId))

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't get the post!")
		return
	}

	helpers.SerializeAndReturn(w, post)
}
