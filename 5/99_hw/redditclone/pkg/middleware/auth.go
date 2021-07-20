package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"io"
	"net/http"
	"redditclone/configs"
	"strings"
)

func AuthCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("error is here: %v", err)
				http.Error(w, "Internal server error!", http.StatusInternalServerError)
			}
		}()

		config, err := configs.LoadConfig("configs")

		secret := config.Token

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			jsonError(w, http.StatusUnauthorized, "No authorization header")
			return
		}

		authParts := strings.Split(authHeader, " ")

		if len(authParts) != 2 {
			jsonError(w, http.StatusUnauthorized, "Invalid authorization header")
			return
		}

		if authParts[0] != "Bearer" {
			jsonError(w, http.StatusUnauthorized, "Invalid authorization header")
			return
		}

		inToken := authParts[1]

		hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
			method, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok || method.Alg() != "HS256" {
				return nil, fmt.Errorf("bad sign method")
			}
			return []byte(secret), nil
		}

		token, err := jwt.Parse(inToken, hashSecretGetter)
		if err != nil || !token.Valid {
			jsonError(w, http.StatusUnauthorized, "bad token")
			return
		}

		tokenData, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			jsonError(w, http.StatusUnauthorized, "no payload")
		}

		userData := make(map[string]string)
		for key, value := range tokenData["user"].(map[string]interface{}) {
			strKey := fmt.Sprintf("%v", key)
			strValue := fmt.Sprintf("%v", value)

			userData[strKey] = strValue
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, configs.UserCtx, userData)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func jsonError(w io.Writer, status int, msg string) {
	resp, _ := json.Marshal(map[string]interface{}{
		"status": status,
		"error":  msg,
	})
	w.Write(resp)
}
