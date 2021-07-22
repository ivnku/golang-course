package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"sync"
	"time"

	"github.com/nfnt/resize"
	"github.com/streadway/amqp"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ImgResizeTask struct {
	PhotoID uint32
	Name    string
	MD5     string
}

type Photo struct {
	ID     uint32 `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	Url    string
	Status int
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

	_, err = rabbitChan.QueueDeclare(
		ImageResizeQueueName, // name
		true,                 // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)
	fatalOnError("cant init queue", err)

	err = rabbitChan.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	fatalOnError("cant set QoS", err)

	tasks, err := rabbitChan.Consume(
		ImageResizeQueueName, // queue
		"",                   // consumer
		false,                // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)
	fatalOnError("cant register consumer", err)

	wg := &sync.WaitGroup{}
	wg.Add(10)

	for i := 0; i <= 10; i++ {
		go ResizeWorker(tasks)
	}

	fmt.Println("worker started")
	wg.Wait()
}

var (
	sizes = []uint{80, 160, 320}
)

func ResizeWorker(tasks <-chan amqp.Delivery) {
	for taskItem := range tasks {
		fmt.Printf("incoming task %+v\n", taskItem)

		task := &ImgResizeTask{}
		err := json.Unmarshal(taskItem.Body, task)
		if err != nil {
			fmt.Println("cant unpack json", err)
			taskItem.Ack(false)
			continue
		}

		originalPath := fmt.Sprintf("./images/%s.jpg", task.MD5)
		for _, size := range sizes {
			time.Sleep(3 * time.Second)
			resizedPath := fmt.Sprintf("./images/%s_%d.jpg", task.MD5, size)
			err := ResizeImage(originalPath, resizedPath, size)
			if err != nil {
				fmt.Println("resize failed", err)
			}
		}

		db.Model(&Photo{ID: task.PhotoID}).Updates(map[string]interface{}{
			"status": 1,
		})

		taskItem.Ack(false)
	}
}

func ResizeImage(originalPath string, resizedPath string, size uint) error {
	file, err := os.Open(originalPath)
	if err != nil {
		return fmt.Errorf("cant open file %s: %s", originalPath, err)
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		return fmt.Errorf("cant jpeg decode file %s", err)
	}
	file.Close()

	resizeImage := resize.Resize(size, 0, img, resize.Lanczos3)

	out, err := os.Create(resizedPath)
	if err != nil {
		return fmt.Errorf("cant create file %s: %s", resizedPath, err)
	}
	defer out.Close()

	jpeg.Encode(out, resizeImage, nil)

	return nil
}

// только на этапе инициализации!
func fatalOnError(msg string, err error) {
	if err != nil {
		log.Fatal(msg + " " + err.Error())
	}
}
