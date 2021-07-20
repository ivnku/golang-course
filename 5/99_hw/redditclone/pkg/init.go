package pkg

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"redditclone/pkg/domain/comment"
	"redditclone/pkg/domain/post"
	"redditclone/pkg/domain/user"
	"redditclone/pkg/middleware"
	"time"
)

func InitApp() {
	db := InitDb()

	usersRepo := user.NewRepository(db)
	postsRepo := post.NewRepository(db)
	commentsRepo := comment.NewRepository(db)

	usersHandler := user.Handler{Repository: usersRepo}
	postsHandler := post.Handler{Repository: postsRepo, CommentsRepo: commentsRepo, UsersRepo: usersRepo}

	router := mux.NewRouter()
	authRouter := router.PathPrefix("/").Subrouter()

	// Log in & Register routes
	router.HandleFunc("/api/register", usersHandler.Register).Methods("POST")
	router.HandleFunc("/api/login", usersHandler.Auth).Methods("POST")

	// Posts routes
	router.HandleFunc("/api/posts/", postsHandler.List).Methods("GET")
	router.HandleFunc("/api/posts/{categoryName}", postsHandler.CategoryList).Methods("GET")
	router.HandleFunc("/api/user/{userName}", postsHandler.UserList).Methods("GET")
	authRouter.HandleFunc("/api/posts", postsHandler.Create).Methods("POST")
	router.HandleFunc("/api/post/{id}", postsHandler.Get).Methods("GET")
	authRouter.HandleFunc("/api/post/{id}", postsHandler.Delete).Methods("DELETE")

	// Comments routes
	authRouter.HandleFunc("/api/post/{id}", postsHandler.Comment).Methods("POST")
	authRouter.HandleFunc("/api/post/{postId}/{commentId}", postsHandler.DeleteComment).Methods("DELETE")

	//authRouter.HandleFunc("/hello/posts", postsHandler.List).Methods("GET")

	authRouter.Use(middleware.AuthCheck)

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./template")))

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func InitDb() *gorm.DB {
	dsn := "root@tcp(127.0.0.1:3306)/redditclone?"
	dsn += "&charset=utf8"
	dsn += "&interpolateParams=true"

	db, err := gorm.Open(mysql.Open(dsn))

	if err != nil {
		panic(err)
	}

	return db
}
