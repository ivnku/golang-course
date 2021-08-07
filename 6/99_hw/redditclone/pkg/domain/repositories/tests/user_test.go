package tests

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"redditclone/pkg/domain/models"
	"redditclone/pkg/domain/repositories"
	"reflect"
	"regexp"
	"testing"

	sqlmock "gopkg.in/DATA-DOG/go-sqlmock.v1"
)

/**
 * @Description: Testing getting the list of all users
 * @param t
 */
func TestList(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repoDb, err := gorm.Open(mysql.New(mysql.Config{Conn: db}))
	repo := repositories.NewUsersRepository(repoDb)

	rows := sqlmock.NewRows([]string{"id", "name", "password"})
	expect := []*models.User{
		{1, "Someuser", "somepasswd"},
		{2, "Anotheruser", "anotherpasswd"},
		{3, "Thirduser", "thirdpasswd"},
	}
	for _, user := range expect {
		rows = rows.AddRow(user.ID, user.Name, user.Password)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users`")).WillReturnRows(rows)

	users, err := repo.List()
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if !reflect.DeepEqual(users, expect) {
		t.Errorf("results not match, want %v, have %v", expect, users)
		return
	}
}

/**
 * @Description: Testing getting of a user by id
 * @param t
 */
func TestGet(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repoDb, err := gorm.Open(mysql.New(mysql.Config{Conn: db}))
	repo := repositories.NewUsersRepository(repoDb)

	var userId int
	rows := sqlmock.NewRows([]string{"id", "name", "password"})
	expect := []*models.User{{uint(userId), "Someuser", "somepasswd"}}
	for _, user := range expect {
		rows = rows.AddRow(user.ID, user.Name, user.Password)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE `users`.`id` = ?") + " +").WillReturnRows(rows)

	user, err := repo.Get(userId)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if !reflect.DeepEqual(user, expect[0]) {
		t.Errorf("results not match, want %v, have %v", expect[0], user)
		return
	}
}

/**
 * @Description: Test getting a user by name
 * @param t
 */
func TestGetByName(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repoDb, err := gorm.Open(mysql.New(mysql.Config{Conn: db}))
	repo := repositories.NewUsersRepository(repoDb)

	rows := sqlmock.NewRows([]string{"id", "name", "password"})
	expect := []*models.User{
		{1, "Someuser", "somepasswd"},
		{2, "Targetuser", "somepasswd"},
		{3, "Anotheruser", "somepasswd"},
	}
	for _, user := range expect {
		rows = rows.AddRow(user.ID, user.Name, user.Password)
	}

	mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM `users` WHERE name = ?") + " +").WillReturnRows(rows)

	userName := "Targetuser"
	user, err := repo.GetByName(userName)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
		return
	}
	if !reflect.DeepEqual(user, expect[0]) {
		t.Errorf("results not match, want %v, have %v", expect[0], user)
		return
	}
}

/**
 * @Description: Tetsing create of a user
 * @param t
 */
func TestCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("cant create mock: %s", err)
	}
	defer db.Close()

	repoDb, err := gorm.Open(mysql.New(mysql.Config{Conn: db}))
	repo := repositories.NewUsersRepository(repoDb)

	name := "newuser"
	pass := "newpasswd"
	testUser := &models.User{
		Name:     name,
		Password: pass,
	}

	//ok query
	mock.ExpectBegin()
	mock.
		ExpectExec("INSERT INTO `users` .*").
		WithArgs(name, pass).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	id, err := repo.Create(testUser)
	if err != nil {
		t.Errorf("unexpected err: %s", err)
		return
	}
	if id != 1 {
		t.Errorf("bad id: want %v, have %v", id, 1)
		return
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	// query error
	mock.ExpectBegin()
	mock.
		ExpectExec("INSERT INTO `users`").
		WithArgs(name, pass).
		WillReturnError(fmt.Errorf("bad query"))

	_, err = repo.Create(testUser)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}

	// result error
	mock.ExpectBegin()
	mock.
		ExpectExec("INSERT INTO `users`").
		WithArgs(name, pass).
		WillReturnResult(sqlmock.NewErrorResult(fmt.Errorf("bad_result")))
	mock.ExpectCommit()

	_, err = repo.Create(testUser)
	if err == nil {
		t.Errorf("expected error, got nil")
		return
	}
}
