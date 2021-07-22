package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"

	"github.com/Shopify/sarama"
	"github.com/rs/xid"
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

type ModerationOut struct {
	XID      string
	MD5      string
	Filename string
	UserID   uint32
}

func publishPage(w http.ResponseWriter, r *http.Request) {
	w.Write(uploadFormTmpl)
}

type PhotoView struct {
	ID     uint32 `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Url    string
	Status int
}

func (i *PhotoView) TableName() string {
	return "photos"
}

func indexPage(w http.ResponseWriter, r *http.Request) {
	items := make([]*PhotoView, 0, 10)
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
	uploadData, handler, _ := r.FormFile("my_file")
	defer uploadData.Close()

	xid := xid.New().String()

	tmpFile := "./images/" + xid + ".jpg"
	newFile, _ := os.Create(tmpFile)

	hasher := md5.New()
	io.Copy(newFile, io.TeeReader(uploadData, hasher))
	newFile.Sync()
	newFile.Close()

	md5Sum := hex.EncodeToString(hasher.Sum(nil))

	realFile := "./images/" + md5Sum + ".jpg"
	os.Rename(tmpFile, realFile)

	data, _ := json.Marshal(ModerationOut{
		XID:      xid,
		Filename: handler.Filename,
		MD5:      md5Sum,
		UserID:   1,
	})

	partition, offset, err := moderationProducer.SendMessage(&sarama.ProducerMessage{
		Topic: "moderation",
		Value: sarama.ByteEncoder(data),
	})
	log.Println(string(data), partition, offset, err)

	fmt.Fprintf(w, "You image will be processed shortly")
}

var (
	mysqlAddr = flag.String("mysql", "root:love@tcp(localhost:3306)/golang?&charset=utf8&interpolateParams=true", "mysql addr")

	db *gorm.DB

	moderationProducer sarama.SyncProducer
)

func main() {
	flag.Parse()
	var err error

	db, err = gorm.Open(mysql.Open(*mysqlAddr), &gorm.Config{})
	fatalOnError("cant connect to mysql", err)

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	moderationProducer, err = sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

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
