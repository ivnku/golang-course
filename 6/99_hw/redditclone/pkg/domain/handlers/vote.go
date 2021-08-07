package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
	"redditclone/configs"
	"redditclone/pkg/auth"
	"redditclone/pkg/domain/repositories"
	"redditclone/pkg/domain/services/votes"
	"redditclone/pkg/helpers"
	"strconv"
)

type VotesHandler struct {
	VotesRepository repositories.IVotesRepository
	PostsRepository repositories.IPostsRepository
	Config          configs.Config
}

/**
 * @Description: Upvote the specific post
 * @receiver h
 * @param w
 * @param r
 */
func (h *VotesHandler) Upvote(w http.ResponseWriter, r *http.Request) {
	routeParams := mux.Vars(r)

	ctx := r.Context()
	user := ctx.Value(configs.UserCtx).(auth.UserData)

	userId, err := strconv.ParseUint(user.Id, 10, 0)

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't convert userId to uint!")
		return
	}

	post, err := votes.ApplyVote(h.PostsRepository, h.VotesRepository, routeParams["id"], uint(userId), 1)

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

	ctx := r.Context()
	user := ctx.Value(configs.UserCtx).(auth.UserData)

	userId, err := strconv.ParseUint(user.Id, 10, 0)

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't convert userId to uint!")
		return
	}

	post, err := votes.ApplyVote(h.PostsRepository, h.VotesRepository, routeParams["id"], uint(userId), -1)

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

	ctx := r.Context()
	user := ctx.Value(configs.UserCtx).(auth.UserData)

	userId, err := strconv.ParseUint(user.Id, 10, 0)

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't convert userId to uint!")
		return
	}

	post, err := votes.Unvote(h.PostsRepository, h.VotesRepository, uint(userId), routeParams["id"])

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't apply the vote!")
		return
	}

	helpers.SerializeAndReturn(w, post)
}
