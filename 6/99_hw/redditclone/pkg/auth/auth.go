package auth

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"redditclone/configs"
	"strconv"
	"time"
)

var ctx = context.Background()

type SessionManager struct {
	User   UserData
	Rdb    *redis.Client
	config configs.Config
}

type MyCustomClaims struct {
	jwt.StandardClaims
	User UserData `json:"user"`
}

type UserData struct {
	Id       string `json:"id"`
	Username string `json:"username"`
}

/**
 * @Description: Constructor for the SessionManager
 * @receiver sm
 * @param config
 * @return *SessionManager
 */
func NewSessionManager(config configs.Config, rdb *redis.Client) *SessionManager {
	return &SessionManager{
		config: config,
		Rdb:    rdb,
	}
}

/**
 * @Description: Authenticate user
 * @param userId
 * @param login
 * @param password
 * @return string
 * @return error
 */
func (sm *SessionManager) Auth(userId uint, login, password string) (string, error) {
	passwordHash, err := sm.HashPassword(password)

	if err != nil {
		return "", err
	}

	if !sm.isPasswordCorrect(password, passwordHash) {
		return "", fmt.Errorf("password is incorrect")
	}

	token, err := sm.CreateSession(login, userId)

	if err != nil {
		return "", fmt.Errorf("couldn't create session")
	}

	return token, nil
}

/**
 * @Description: Generate a hash string from a password
 * @param password
 * @return string
 * @return error
 */
func (sm *SessionManager) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

/**
 * @Description: Check if provided password matches the password from db
 * @param password
 * @param hash
 * @return bool
 */
func (sm *SessionManager) isPasswordCorrect(password, hash string) bool {
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
func (sm *SessionManager) GenerateJWT(login string, userId uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &MyCustomClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * time.Duration(20)).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserData{strconv.Itoa(int(userId)), login},
	})

	tokenString, err := token.SignedString([]byte(sm.config.Token))

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
func (sm *SessionManager) CheckToken(inToken string) (*MyCustomClaims, error) {
	secret := sm.config.Token

	hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || method.Alg() != "HS256" {
			return nil, fmt.Errorf("bad sign method")
		}
		return []byte(secret), nil
	}

	token, err := jwt.ParseWithClaims(inToken, &MyCustomClaims{}, hashSecretGetter)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse token: " + err.Error())
	}

	tokenData, ok := token.Claims.(*MyCustomClaims)
	if !ok {
		return nil, fmt.Errorf("no payload")
	}

	if !token.Valid {
		sm.Rdb.Del(ctx, tokenData.User.Id)
		return nil, fmt.Errorf("bad token")
	}

	return tokenData, nil
}

/**
 * @Description: Generate a jwt token and create a session with its value
 * @receiver sm
 * @param login
 * @param userId
 * @return token
 * @return err
 */
func (sm *SessionManager) CreateSession(login string, userId uint) (token string, err error) {
	token, err = sm.GenerateJWT(login, userId)

	if err != nil {
		return token, fmt.Errorf("couldn't create token string")
	}

	sm.Rdb.Set(ctx, strconv.Itoa(int(userId)), token, 0)
	return
}
