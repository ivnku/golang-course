package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"redditclone/configs"
	"redditclone/pkg/domain/repositories"
	"redditclone/pkg/domain/services/votes"
	"redditclone/pkg/helpers"
	"strconv"
)

type VotesHandler struct {
	VotesRepository repositories.VotesRepository
	PostsRepository repositories.PostsRepository
}

/**
 * @Description: Upvote the specific post
 * @receiver h
 * @param w
 * @param r
 */
func (h *VotesHandler) Upvote(w http.ResponseWriter, r *http.Request) {
	routeParams := mux.Vars(r)

	postId, err := strconv.ParseUint(routeParams["id"], 10, 0)

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't convert postId to uint!")
		return
	}

	ctx := r.Context()
	user := ctx.Value(configs.UserCtx).(map[string]string)

	userId, err := strconv.ParseUint(user["id"], 10, 0)

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't convert userId to uint!")
		return
	}

	post, err := votes.ApplyVote(h.PostsRepository, h.VotesRepository, uint(postId), uint(userId), 1)

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't apply the vote!")
		return
	}

	helpers.SerializeAndReturn(w, post)
}

/**
 * @Description: Downvote the specific post
 * @receiver h
 * @param w
 * @param r
 */
func (h *VotesHandler) Downvote(w http.ResponseWriter, r *http.Request) {
	routeParams := mux.Vars(r)

	postId, err := strconv.ParseUint(routeParams["id"], 10, 0)

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't convert postId to uint!")
		return
	}

	ctx := r.Context()
	user := ctx.Value(configs.UserCtx).(map[string]string)

	userId, err := strconv.ParseUint(user["id"], 10, 0)

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't convert userId to uint!")
		return
	}

	post, err := votes.ApplyVote(h.PostsRepository, h.VotesRepository, uint(postId), uint(userId), -1)

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't apply the vote!")
		return
	}

	helpers.SerializeAndReturn(w, post)
}

/**
 * @Description: Unvote for the specific post
 * @receiver h
 * @param w
 * @param r
 */
func (h *VotesHandler) Unvote(w http.ResponseWriter, r *http.Request) {
	routeParams := mux.Vars(r)

	postId, err := strconv.ParseUint(routeParams["id"], 10, 0)

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't convert postId to uint!")
		return
	}

	ctx := r.Context()
	user := ctx.Value(configs.UserCtx).(map[string]string)

	userId, err := strconv.ParseUint(user["id"], 10, 0)

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't convert userId to uint!")
		return
	}

	post, err := votes.Unvote(h.PostsRepository, h.VotesRepository, uint(userId), uint(postId))

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't apply the vote!")
		return
	}

	helpers.SerializeAndReturn(w, post)
}