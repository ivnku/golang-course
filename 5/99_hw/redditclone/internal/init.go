package internal

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/http"
	"redditclone/internal/handlers"
	"time"
)

var DB *gorm.DB

func InitApp() {
	InitDb()
	router := mux.NewRouter()
	router.HandleFunc("/hello", (&handlers.UserHandler{}).List)
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

func InitDb() {
	dsn := "root@tcp(127.0.0.1:3306)/redditclone?"
	dsn += "&charset=utf8"
	dsn += "&interpolateParams=true"

	db, err := gorm.Open(mysql.Open(dsn))

	if err != nil {
		panic(err)
	}

	DB = db
}
