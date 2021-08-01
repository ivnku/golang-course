package pkg

import (
	"context"
	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"redditclone/configs"
	"redditclone/pkg/auth"
	"redditclone/pkg/domain/handlers"
	"redditclone/pkg/domain/repositories"
	"redditclone/pkg/middleware"
	"time"
)

func InitApp() {
	db := InitDb()

	rdb := redis.NewClient(&redis.Options{Addr: "localhost:6379", Password: "", DB: 0})
	defer func() {
		if err := rdb.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongodb, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := mongodb.Disconnect(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	config, err := configs.LoadConfig("configs")
	sessionManager := auth.NewSessionManager(config, rdb)

	if err != nil {
		log.Fatal(err)
	}

	usersRepo := repositories.NewUsersRepository(db)
	postsRepo := repositories.NewPostsRepository(mongodb)
	commentsRepo := repositories.NewCommentsRepository(mongodb)
	votesRepo := repositories.NewVotesRepository(mongodb)

	usersHandler := handlers.UsersHandler{UsersRepository: usersRepo, Config: config, SessionManager: *sessionManager}
	postsHandler := handlers.PostsHandler{PostsRepository: postsRepo, CommentsRepository: commentsRepo, UsersRepository: usersRepo, Config: config}
	votesHandler := handlers.VotesHandler{VotesRepository: votesRepo, PostsRepository: postsRepo, Config: config}

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

	// Votes routes
	authRouter.HandleFunc("/api/post/{id}/upvote", votesHandler.Upvote).Methods("GET")
	authRouter.HandleFunc("/api/post/{id}/downvote", votesHandler.Downvote).Methods("GET")
	authRouter.HandleFunc("/api/post/{id}/unvote", votesHandler.Unvote).Methods("GET")

	authRouter.Use(middleware.AuthCheck(*sessionManager))

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
