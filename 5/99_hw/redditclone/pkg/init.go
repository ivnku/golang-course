package pkg

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"redditclone/pkg/domain/post"
	"redditclone/pkg/domain/user"
	"time"
)

func InitApp() {
	db := InitDb()

	usersRepo := user.Repository{DB: db}
	postsRepo := post.Repository{DB: db}

	usersHandler := user.Handler{Repository: usersRepo}
	postsHandler := post.Handler{Repository: postsRepo}

	router := mux.NewRouter()
	router.HandleFunc("/hello", usersHandler.List)
	router.HandleFunc("/hello/posts", postsHandler.List)
	router.HandleFunc("/api/posts/", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("Hello world!")) })
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
