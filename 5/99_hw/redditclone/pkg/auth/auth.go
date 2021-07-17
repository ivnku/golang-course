package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"redditclone/configs"
	"strconv"
	"time"
)

/**
 * @Description: Authenticate user
 * @param userId
 * @param login
 * @param password
 * @return string
 * @return error
 */
func Auth(userId uint, login, password string) (string, error) {
	passwordHash, err := HashPassword(password)

	if err != nil {
		return "", err
	}

	if !isPasswordCorrect(password, passwordHash) {
		return "", fmt.Errorf("password is incorrect")
	}

	return GenerateJWT(login, userId)
}

/**
 * @Description: Generate a hash string from a password
 * @param password
 * @return string
 * @return error
 */
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

/**
 * @Description: Check if provided password matches the password from db
 * @param password
 * @param hash
 * @return bool
 */
func isPasswordCorrect(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

/**
 * @Description: Generate jwt token
 * @param login
 * @param userId
 * @return string
 * @return error
 */
func GenerateJWT(login string, userId uint) (string, error) {
	type MyCustomClaims struct {
		User map[string]string `json:"user"`
		jwt.StandardClaims
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.MapClaims{
		"user": map[string]string{"username": login, "id": strconv.Itoa(int(userId))},
		"exp":  time.Now().Add(time.Minute * time.Duration(10)).Unix(),
		"iat":  time.Now().Unix(),
	})

	config, err := configs.LoadConfig("configs")

	if err != nil || &config == nil {
		return "", err
	}

	tokenString, err := token.SignedString([]byte(config.Token))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
