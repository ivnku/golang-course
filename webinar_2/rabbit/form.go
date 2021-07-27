package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"text/template"

	"github.com/streadway/amqp"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var uploadFormTmpl = []byte(`
<html>
	<body>
	<form action="/upload" method="post" enctype="multipart/form-data">
		Image: <input type="file" name="my_file">
		<input type="submit" value="Upload">
	</form>
	</body>
</html>
`)

var imagesListTmpl = `
<html>
<body>
	<h1>Images <a href="/publish">Publish</a></h1>
	{{range .}}
		<div style="border:1px solid black; margin: 3px;">
		{{if eq .Status 1}}
		<img src="/images/{{.Url}}_160.jpg">
		{{else}}
		record proceesing, please wait
		{{end}}
		</div>
	{{end}}
</body>
</html>
`

type ImgResizeTask struct {
	PhotoID uint32
	Name    string
	MD5     string
}

func publishPage(w http.ResponseWriter, r *http.Request) {
	w.Write(uploadFormTmpl)
}

type Photo struct {
	ID     uint32 `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Url    string
	Status int
	UserID uint32
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	items := make([]*Photo, 0, 10)
	if err := db.Find(&items).Error; err != nil {
		http.Error(w, "db err:"+err.Error(), http.StatusInternalServerError)
		return
	}

	// не парсите шаблоны каждый раз - это мделенно. делайте это 1 раз
	tmpl, err := template.New(`example`).Parse(imagesListTmpl)
	if err != nil {
		http.Error(w, "tmpl err:"+err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, items)
	if err != nil {
		http.Error(w, "tmpl exec err:"+err.Error(), http.StatusInternalServerError)
		return
	}
}

func uploadPage(w http.ResponseWriter, r *http.Request) {
	uploadData, handler, err := r.FormFile("my_file")
	defer uploadData.Close()

	tmpName := RandStringRunes(32)

	tmpFile := "./images/" + tmpName + ".jpg"
	newFile, _ := os.Create(tmpFile)

	hasher := md5.New()
	io.Copy(newFile, io.TeeReader(uploadData, hasher))
	newFile.Sync()
	newFile.Close()

	md5Sum := hex.EncodeToString(hasher.Sum(nil))

	realFile := "./images/" + md5Sum + ".jpg"
	os.Rename(tmpFile, realFile)

	photo := &Photo{Url: md5Sum, UserID: 0}
	err = db.Create(photo).Error
	log.Println("created elem id:", photo.ID, err)

	data, _ := json.Marshal(ImgResizeTask{photo.ID, handler.Filename, md5Sum})
	fmt.Println("put task ", string(data))

	err = rabbitChan.Publish(
		"",                   // exchange
		ImageResizeQueueName, // routing key
		false,                // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         data,
		})
	if err != nil {
		http.Error(w, "publish event: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", 302)
}

const (
	ImageResizeQueueName = "image_resize"
)

var (
	rabbitAddr = flag.String("rabbit", "amqp://guest:guest@localhost:5672/", "rabbit addr")
	mysqlAddr  = flag.String("mysql", "root:love@tcp(localhost:3306)/golang?&charset=utf8&interpolateParams=true", "mysql addr")

	rabbitConn *amqp.Connection
	rabbitChan *amqp.Channel

	db *gorm.DB
)

func main() {
	flag.Parse()
	var err error

	db, err = gorm.Open(mysql.Open(*mysqlAddr), &gorm.Config{})
	fatalOnError("cant connect to mysql", err)

	rabbitConn, err = amqp.Dial(*rabbitAddr)
	fatalOnError("cant connect to rabbit", err)

	rabbitChan, err = rabbitConn.Channel()
	fatalOnError("cant open chan", err)
	defer rabbitChan.Close()

	q, err := rabbitChan.QueueDeclare(
		ImageResizeQueueName, // name
		true,                 // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	fatalOnError("cant init queue", err)
	fmt.Printf("queue %s have %d msg and %d consumers\n",
		q.Name, q.Messages, q.Consumers)

	http.HandleFunc("/", indexPage)
	http.HandleFunc("/publish", publishPage)
	http.HandleFunc("/upload", uploadPage)
	staticHandler := http.StripPrefix(
		"/images/",
		http.FileServer(http.Dir("./images")),
	)
	http.Handle("/images/", staticHandler)

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}

// только на этапе инициализации!
func fatalOnError(msg string, err error) {
	if err != nil {
		log.Fatal(msg + " " + err.Error())
	}
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
