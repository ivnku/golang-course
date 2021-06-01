package main

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	// "html/template"
	"log"
	"net/http"
	"regexp"
	"strconv"

	"golang.org/x/crypto/argon2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/vk"

	"github.com/asaskevich/govalidator"
)

type User struct {
	ID    uint32
	Login string
	Email string
	Ver   int32
}

const (
	VK_APP_ID  = "7065390"
	VK_APP_KEY = "cQZe3Vvo4mHotmetUdXK"
	// куда идти с токеном за информацией
	VK_API_URL = "https://api.vk.com/method/users.get?fields=photo_50&access_token=%s&v=5.52"
	// куда идти для получения токена
	VK_AUTH_URL = "https://oauth.vk.com/authorize?client_id=7065390&redirect_uri=http://localhost:8080/user/login_oauth&response_type=code&scope=email"
)

type VKOauthResp struct {
	Response []struct {
		FirstName string `json:"first_name"`
		Photo     string `json:"photo_50"`
	}
}

type UserHandler struct {
	DB       *sql.DB
	Tmpl     *MyTemplate
	Sessions SessionManager
}

func (uh *UserHandler) hashPass(plainPassword, salt string) []byte {
	hashedPass := argon2.IDKey([]byte(plainPassword), []byte(salt), 1, 64*1024, 4, 32)
	res := []byte(salt)
	return append(res, hashedPass...)
}

var (
	errUserNotFound = errors.New("No user record found")
	errBadPass      = errors.New("Bad password")
	errUserExists   = errors.New("User Exists")

	loginRE = regexp.MustCompile(`^[\w-_\.]+$`)
)

func (uh *UserHandler) passwordIsValid(pass string, row *sql.Row) (*User, error) {
	var (
		dbPass []byte
		user   = &User{}
	)
	err := row.Scan(&user.ID, &user.Login, &user.Ver, &dbPass)
	if err == sql.ErrNoRows {
		return nil, errUserNotFound
	} else if err != nil {
		return nil, err
	}

	salt := string(dbPass[0:8])
	if !bytes.Equal(uh.hashPass(pass, salt), dbPass) {
		return nil, errBadPass
	}
	return user, nil
}

func GetUserByLogin(db *sql.DB, login string) (*User, error) {
	return parseRowToUser(db.QueryRow("SELECT id, login, email, ver FROM users WHERE login = ?", login))
}

func GetUserByID(db *sql.DB, id uint32) (*User, error) {
	return parseRowToUser(db.QueryRow("SELECT id, login, email, ver FROM users WHERE id = ?", id))
}

func parseRowToUser(row *sql.Row) (*User, error) {
	user := &User{}
	err := row.Scan(&user.ID, &user.Login, &user.Email, &user.Ver)
	if err == sql.ErrNoRows {
		return nil, errUserNotFound
	} else if err != nil {
		return nil, err
	}
	return user, nil
}

func (uh *UserHandler) checkPasswordByUserID(uid uint32, pass string) (*User, error) {
	row := uh.DB.QueryRow("SELECT id, login, ver, password FROM users WHERE id = ?", uid)
	return uh.passwordIsValid(pass, row)
}

func (uh *UserHandler) checkPasswordByLogin(login, pass string) (*User, error) {
	row := uh.DB.QueryRow("SELECT id, login, ver, password FROM users WHERE login = ?", login)
	return uh.passwordIsValid(pass, row)
}

func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		uh.Tmpl.Render(r.Context(), w, "login.html", map[string]interface{}{
			"VKAuthURL": VK_AUTH_URL,
		})
		return
	}

	login := r.FormValue("login")
	pass := r.FormValue("password")

	user, err := uh.checkPasswordByLogin(login, pass)

	switch err {
	case nil:
		// all is ok
	case errUserNotFound:
		http.Error(w, "No user", http.StatusBadRequest)
	case errBadPass:
		http.Error(w, "Bad pass", http.StatusBadRequest)
	default:
		http.Error(w, "Db err", http.StatusInternalServerError)
	}
	if err != nil {
		return
	}

	uh.Sessions.Create(w, user)
	http.Redirect(w, r, "/photos/", http.StatusFound)
}

func (uh *UserHandler) LoginOauth(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	if code == "" {
		http.Error(w, "no ouath code", http.StatusBadRequest)
		return
	}

	conf := oauth2.Config{
		ClientID:     VK_APP_ID,
		ClientSecret: VK_APP_KEY,
		RedirectURL:  "http://localhost:8080/user/login_oauth",
		Endpoint:     vk.Endpoint,
	}

	token, err := conf.Exchange(r.Context(), code)
	if err != nil {
		log.Println("exchange err", err)
		http.Error(w, "cannot get oauth token", http.StatusInternalServerError)
		return
	}

	emailRaw := token.Extra("email")
	email := ""
	okEmail := true
	if emailRaw != nil {
		email, okEmail = emailRaw.(string)
	}
	userIDraw, okID := token.Extra("user_id").(float64)
	if !okEmail || !okID {
		log.Printf("cant convert data: UID: %T %v, Email: %T %v", token.Extra("user_id"), token.Extra("user_id"), token.Extra("email"), token.Extra("email"))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	login := "vk" + strconv.Itoa(int(userIDraw))
	pass := RandStringRunes(50)
	user, err := uh.createUser(login, email, pass)
	if err != nil && err != errUserExists {
		log.Println("db err", err)
		http.Error(w, "Db err", http.StatusInternalServerError)
		return
	}

	uh.Sessions.Create(w, user)
	http.Redirect(w, r, "/photos/", http.StatusFound)
}

func (uh *UserHandler) createUser(login, email, passIn string) (*User, error) {
	salt := RandStringRunes(8)
	pass := uh.hashPass(passIn, salt)

	user := &User{
		ID:    0,
		Ver:   0,
		Email: email,
	}

	err := uh.DB.QueryRow("SELECT id, ver, login FROM users WHERE email = ? OR login = ?", email, login).
		Scan(&user.ID, &user.Ver, &user.Login)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("db error: %v", err)
	}
	if err != sql.ErrNoRows {
		return user, errUserExists
	}

	result, err := uh.DB.Exec("INSERT INTO users(login, email, password) VALUES(?, ?, ?)", login, email, pass)
	if err != nil {
		return nil, fmt.Errorf("insert error: %v", err)
	}

	affected, _ := result.RowsAffected()
	if affected == 0 {
		return nil, fmt.Errorf("no rows affected")
	}
	uid, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("LastInsertId err: %v", err)
	}
	user.ID = uint32(uid)

	return user, nil
}

func (uh *UserHandler) Reg(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		uh.Tmpl.Render(r.Context(), w, "reg.html", nil)
		return
	}

	login := r.FormValue("login")
	pass := r.FormValue("password")
	email := r.FormValue("email")

	if !govalidator.IsEmail(email) {
		http.Error(w, "Bad email", http.StatusBadRequest)
		return
	}

	if !loginRE.MatchString(login) {
		http.Error(w, "Bad login", http.StatusBadRequest)
		return
	}

	user, err := uh.createUser(login, email, pass)
	switch err {
	case nil:
		// all is ok
	case errUserExists:
		http.Error(w, "Looks like user exists", http.StatusBadRequest)
	default:
		log.Println("db err", err)
		http.Error(w, "Db err", http.StatusInternalServerError)
	}
	if err != nil {
		return
	}

	uh.Sessions.Create(w, user)
	http.Redirect(w, r, "/photos/", http.StatusFound)
}

func (uh *UserHandler) Logout(w http.ResponseWriter, r *http.Request) {
	uh.Sessions.DestroyCurrent(w, r)
	http.Redirect(w, r, "/user/login", http.StatusFound)
}

func (uh *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		uh.Tmpl.Render(r.Context(), w, "change_pass.html", nil)
		return
	}

	if r.FormValue("pass1") == "" || r.FormValue("pass1") != r.FormValue("pass2") {
		http.Error(w, "New password mistmatch", http.StatusBadRequest)
		return
	}

	sess, _ := SessionFromContext(r.Context())
	user, err := uh.checkPasswordByUserID(sess.UserID, r.FormValue("old_password"))
	if err != nil {
		http.Error(w, "Bad pass", http.StatusBadRequest)
		return
	}

	salt := RandStringRunes(8)
	pass := uh.hashPass(r.FormValue("pass1"), salt)

	_, err = uh.DB.Exec("UPDATE users SET password = ?, ver = ver + 1 WHERE id = ?",
		pass, user.ID)
	if err != nil {
		log.Println("update password error", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	user.Ver++ // во избежание рейсов лучше подгрузить из базы

	uh.Sessions.DestroyAll(w, user)
	uh.Sessions.Create(w, user)
	http.Redirect(w, r, "/photos/", http.StatusFound)
}

func (uh *UserHandler) FollowAPI(w http.ResponseWriter, r *http.Request) {
	sess, _ := SessionFromContext(r.Context())
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		RespJSONError(w, http.StatusBadRequest, nil, "bad id")
		return
	}
	folUser, err := GetUserByID(uh.DB, uint32(id))
	if err == errUserNotFound {
		RespJSONError(w, http.StatusBadRequest, nil, "no user")
		return
	}
	if err != nil {
		RespJSONError(w, http.StatusInternalServerError, fmt.Errorf("db error: %v", err), "internal")
		return
	}

	rate := 1
	if r.FormValue("unfollow") == "1" {
		rate = -1
	}

	var res sql.Result
	if rate == 1 {
		res, err = uh.DB.Exec(`INSERT IGNORE INTO user_follows(follow_id, user_id) VALUES(?, ?)`, folUser.ID, sess.UserID)
	} else {
		res, err = uh.DB.Exec(`DELETE FROM user_follows WHERE follow_id = ? AND user_id = ?`, folUser.ID, sess.UserID)
	}
	if err != nil {
		RespJSONError(w, http.StatusInternalServerError, fmt.Errorf("db error: %v", err), "internal")
		return
	}
	aff, _ := res.RowsAffected()
	// dont update rating twice
	if aff <= 0 {
		RespJSON(w, map[string]interface{}{
			"id": id,
		})
		return
	}
	_, err = uh.DB.Exec("UPDATE users SET followers_cnt = followers_cnt + ? WHERE id = ?", rate, folUser.ID)
	if err != nil {
		RespJSONError(w, http.StatusInternalServerError, fmt.Errorf("db error: %v", err), "internal")
		return
	}
	_, err = uh.DB.Exec("UPDATE users SET following_cnt = following_cnt + ? WHERE id = ?", rate, sess.UserID)
	if err != nil {
		RespJSONError(w, http.StatusInternalServerError, fmt.Errorf("db error: %v", err), "internal")
		return
	}

	RespJSON(w, map[string]interface{}{
		"id": id,
	})
	return
}

type UserResp struct {
	ID    uint32 `json:"id"`
	Login string `json:"login"`
}

func (uh *UserHandler) FollowingAPI(w http.ResponseWriter, r *http.Request) {
	sess, _ := SessionFromContext(r.Context())

	rows, err := uh.DB.Query(`SELECT users.id, users.login 
	FROM user_follows 
	LEFT JOIN users ON users.id = user_follows.follow_id
	WHERE user_follows.user_id = ?`, sess.UserID)
	if err != nil {
		RespJSONError(w, http.StatusInternalServerError, fmt.Errorf("db error: %v", err), "internal")
		return
	}
	defer rows.Close()
	result := make([]*UserResp, 0, 10)
	for rows.Next() {
		u := &UserResp{}
		err := rows.Scan(&u.ID, &u.Login)
		if err != nil {
			RespJSONError(w, http.StatusInternalServerError, fmt.Errorf("db error: %v", err), "internal")
			return
		}
		result = append(result, u)
	}
	RespJSON(w, map[string]interface{}{
		"users":    result,
		"followed": true,
	})
}

func (uh *UserHandler) RecomendsAPI(w http.ResponseWriter, r *http.Request) {
	sess, _ := SessionFromContext(r.Context())

	rows, err := uh.DB.Query(`select users.id, users.login 
	from users 
	left join user_follows on users.id = user_follows.follow_id and user_follows.user_id = ?
	where users.id != ? and user_follows.user_id is null`, sess.UserID, sess.UserID)
	if err != nil {
		RespJSONError(w, http.StatusInternalServerError, fmt.Errorf("db error: %v", err), "internal")
		return
	}
	defer rows.Close()
	result := make([]*UserResp, 0, 10)
	for rows.Next() {
		u := &UserResp{}
		err := rows.Scan(&u.ID, &u.Login)
		if err != nil {
			RespJSONError(w, http.StatusInternalServerError, fmt.Errorf("db error: %v", err), "internal")
			return
		}
		result = append(result, u)
	}
	RespJSON(w, map[string]interface{}{
		"users":    result,
		"followed": false,
	})
}
