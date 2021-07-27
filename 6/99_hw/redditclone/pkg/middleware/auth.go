package middleware

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"redditclone/configs"
	"redditclone/pkg/auth"
	"redditclone/pkg/helpers"
	"strings"
)

func AuthCheck(sessionManager auth.SessionManager) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					fmt.Printf("error is here: %v", err)
					http.Error(w, "Internal server error!", http.StatusInternalServerError)
				}
			}()

			authHeader := r.Header.Get("Authorization")

			if authHeader == "" {
				helpers.JsonError(w, http.StatusUnauthorized, "No authorization header")
				return
			}

			authParts := strings.Split(authHeader, " ")

			if len(authParts) != 2 {
				helpers.JsonError(w, http.StatusUnauthorized, "Invalid authorization header")
				return
			}

			if authParts[0] != "Bearer" {
				helpers.JsonError(w, http.StatusUnauthorized, "Invalid authorization header")
				return
			}

			inToken := authParts[1]

			tokenData, err := sessionManager.CheckToken(inToken)

			if err != nil {
				helpers.JsonError(w, http.StatusUnauthorized, err.Error())
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, configs.UserCtx, tokenData.User)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
