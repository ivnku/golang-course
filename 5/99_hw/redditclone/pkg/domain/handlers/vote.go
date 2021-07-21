package handlers

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"redditclone/configs"
	"redditclone/pkg/domain/models"
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

	post, err := h.PostsRepository.Get(uint(postId))

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't get the post!")
		return
	}

	ctx := r.Context()
	user := ctx.Value(configs.UserCtx).(map[string]string)

	userId, err := strconv.ParseUint(user["id"], 10, 0)

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't convert userId to uint!")
		return
	}

	vote := &models.Vote{
		PostId: uint(postId),
		UserId: uint(userId),
		Vote:   1,
	}

	vote, err = h.VotesRepository.Create(vote)

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't create the vote!")
		return
	}

	post.Votes = append(post.Votes, vote)
	post.UpvotePercentage = votes.CalculateUpvotePercentage(post.Votes)
	post.Score = votes.CalculateScore(post.Votes)

	post, err = h.PostsRepository.Update(post, []string{"upvote_percentage", "score"})

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't update the post!")
		return
	}

	postSerialized, err := json.Marshal(post)

	if err != nil {
		helpers.JsonError(w, http.StatusBadRequest, "Couldn't marshal the post!")
		return
	}

	w.Write(postSerialized)
}
