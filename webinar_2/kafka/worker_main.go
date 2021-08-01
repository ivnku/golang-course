package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

//https://gist.github.com/ik5/d8ecde700972d4378d87#gistcomment-3074524
var (
	Black   = Color("\033[1;30m%s\033[0m\n")
	Red     = Color("\033[1;31m%s\033[0m\n")
	Green   = Color("\033[1;32m%s\033[0m\n")
	Yellow  = Color("\033[1;33m%s\033[0m\n")
	Purple  = Color("\033[1;34m%s\033[0m\n")
	Magenta = Color("\033[1;35m%s\033[0m\n")
	Teal    = Color("\033[1;36m%s\033[0m\n")
	White   = Color("\033[1;37m%s\033[0m\n")

	Info = Teal
	Warn = Yellow
	Fata = Red
)

func Color(colorString string) func(string, ...interface{}) {
	sprint := func(format string, args ...interface{}) {
		fmt.Printf(colorString, fmt.Sprintf(format, args...))
	}
	return sprint
}

type CreationIn struct {
	XID      string
	MD5      string
	Filename string
	UserID   uint32
	SpamProp int
}

var (
	mysqlAddr = flag.String("mysql", "root:love@tcp(localhost:3306)/golang?&charset=utf8&interpolateParams=true", "mysql addr")
)

func main() {
	flag.Parse()
	var err error

	db, err := gorm.Open(mysql.Open(*mysqlAddr), &gorm.Config{})
	fatalOnError("cant connect to mysql", err)

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	// go func() {
	// 	select {
	// 	case <-sigterm:
	// 		log.Println("sigterm happend")
	// 		cancel()
	// 	}
	// }()

	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, config)
	if err != nil {
		log.Fatalln("Failed to start Sarama producer:", err)
	}

	var (
		addrs = []string{"localhost:9092"}
	)
	// тут мы "модерируем" запись - проставляем ей вероятность что это спам
	wg.Add(1)
	go startConsume(ctx, wg, addrs, "moderation", []string{"moderation"},
		&ModerationConsumerHandler{Log: Magenta, Next: producer})

	// тут мы создаем запись в БД
	wg.Add(1)
	go startConsume(ctx, wg, addrs, "creation", []string{"create"},
		&CreateConsumerHandler{Log: Purple, Next: producer, DB: db})

	// тут мы уведомляем рассылаем уведолмения об этой записи
	wg.Add(1)
	go startConsume(ctx, wg, addrs, "notify_email", []string{"publish"},
		&NotifyEmailConsumerHandler{Log: Yellow})
	wg.Add(1)
	go startConsume(ctx, wg, addrs, "update_search", []string{"publish"},
		&UpdateSearchConsumerHandler{Log: Green})

	select {
	case <-sigterm:
		log.Println("sigterm happend in finish")
		cancel()
	}

	// wg.Wait()
}

func startConsume(ctx context.Context, wg *sync.WaitGroup,
	addrs []string, groupName string, topics []string,
	handler sarama.ConsumerGroupHandler,
) {
	defer wg.Done()

	config := sarama.NewConfig()
	config.Version = sarama.V2_5_0_0
	config.Consumer.Return.Errors = true
	// config.Consumer.
	log.Printf("start consumer group %s, topics %+v\n", groupName, topics)

	group, err := sarama.NewConsumerGroup(addrs, groupName, config)
	if err != nil {
		log.Fatalln("cant create consumer group", err)
	}

	go func() {
		for err := range group.Errors() {
			fmt.Println("ERROR", err)
		}
	}()

	for {
		err := group.Consume(ctx, topics, handler)
		if err != nil {
			log.Fatalln("group.Consume failed", err)
		}
	}
}

// только на этапе инициализации!
func fatalOnError(msg string, err error) {
	if err != nil {
		log.Fatal(msg + " " + err.Error())
	}
}
