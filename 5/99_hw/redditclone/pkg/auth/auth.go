package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"redditclone/configs"
	"strconv"
	"time"
)

type MyCustomClaims struct {
	jwt.StandardClaims
	User map[string]string `json:"user"`
}

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
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &MyCustomClaims{
		jwt.StandardClaims{
			ExpiresAt:  time.Now().Add(time.Minute * time.Duration(20)).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		map[string]string{"username": login, "id": strconv.Itoa(int(userId))},
	})

	config := configs.Conf

	tokenString, err := token.SignedString([]byte(config.Token))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

/**
 * @Description: Check validity of jwt token. And if it's valid return token data
 * @param inToken
 * @return *MyCustomClaims
 * @return error
 */
func CheckToken(inToken string) (*MyCustomClaims, error) {
	config := configs.Conf
	secret := config.Token

	hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || method.Alg() != "HS256" {
			return nil, fmt.Errorf("bad sign method")
		}
		return []byte(secret), nil
	}

	token, err := jwt.ParseWithClaims(inToken, &MyCustomClaims{}, hashSecretGetter)
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("bad token")
	}

	tokenData, ok := token.Claims.(*MyCustomClaims)
	if !ok {
		return nil, fmt.Errorf("no payload")
	}

	return tokenData, nil
}
