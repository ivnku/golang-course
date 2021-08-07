package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"redditclone/configs"
	"redditclone/pkg/auth"
	"redditclone/pkg/domain/models"
	"redditclone/pkg/domain/repositories"
	"redditclone/pkg/helpers"
	"strconv"
	"time"
)

type PostsHandler struct {
	PostsRepository    repositories.IPostsRepository
	CommentsRepository repositories.ICommentsRepository
	UsersRepository    repositories.IUsersRepository
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

	post, err := h.PostsRepository.Get(params["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if post == nil {
		http.Error(w, "Post doesn't exist", http.StatusNotFound)
		return
	}

	// increment views each time a user open the post
	post.Views++
	post, err = h.PostsRepository.Update(post, []primitive.E{{"Views", post.Views}})

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
	user := ctx.Value(configs.UserCtx).(auth.UserData)

	userId, err := strconv.ParseUint(user.Id, 10, 0)

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't convert userId to uint!")
		return
	}

	post.User.ID = uint(userId)
	post.User.Name = user.Username
	post.CreatedAt = time.Now()

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

	success, err := h.PostsRepository.Delete(params["id"])

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

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
	user := ctx.Value(configs.UserCtx).(auth.UserData)

	userId, err := strconv.ParseUint(user.Id, 10, 0)

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't convert userId to uint!")
		return
	}

	postComment.UserID = uint(userId)
	postComment.User.ID = uint(userId)
	postComment.User.Name = user.Username
	postComment.Body = reqComment.Comment
	postComment.Created = time.Now()
	postId, err := primitive.ObjectIDFromHex(routeParams["id"])
	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't create a primitive.ObjectIDFromHex() for postId!")
		return
	}
	postComment.PostID = postId

	postComment, err = h.CommentsRepository.Create(postComment)

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't create a comment!")
		return
	}

	post, err := h.PostsRepository.Get(routeParams["id"])

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

	_, err := h.CommentsRepository.Delete(routeParams["commentId"])

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't delete post comment!")
		return
	}

	post, err := h.PostsRepository.Get(routeParams["postId"])

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't get the post!")
		return
	}

	helpers.SerializeAndReturn(w, post)
}
