package main

import (
	"encoding/json"
	"fmt"
	"image/jpeg"
	"log"
	"os"
	"time"

	"github.com/Shopify/sarama"
	"github.com/nfnt/resize"
	"gorm.io/gorm"
)

type CreateConsumerHandler struct {
	Log  func(string, ...interface{})
	Next sarama.SyncProducer
	DB   *gorm.DB
}

func (h *CreateConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error {
	h.Log("Setup happend")
	return nil
}

func (h *CreateConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	h.Log("Cleanup happend")
	return nil
}

var (
	sizes = []uint{80, 160, 320}
)

type PhotoCreate struct {
	ID       uint32 `sql:"AUTO_INCREMENT" gorm:"primary_key"`
	XID      string `gorm:"column:xid"`
	Url      string
	Status   int
	UserID   uint32
	SpamProp int
}

func (i *PhotoCreate) TableName() string {
	return "photos"
}

func (h *CreateConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.Log("Processing %v/%v/%v %s %s", msg.Topic, msg.Partition, msg.Offset,
			msg.Key,
			msg.Value)

		in := &CreationIn{}
		err := json.Unmarshal(msg.Value, in)

		originalPath := fmt.Sprintf("./images/%s.jpg", in.MD5)
		for _, size := range sizes {
			time.Sleep(3 * time.Second)
			resizedPath := fmt.Sprintf("./images/%s_%d.jpg", in.MD5, size)
			err := ResizeImage(originalPath, resizedPath, size)
			if err != nil {
				fmt.Println("resize failed", err)
			}
		}

		photo := &PhotoCreate{
			XID:      in.XID,
			Url:      in.MD5,
			Status:   1,
			SpamProp: in.SpamProp,
			UserID:   in.UserID,
		}
		err = h.DB.Create(photo).Error
		log.Println("created elem id:", photo.ID, err)

		out, _ := json.Marshal(photo)
		partition, offset, err := h.Next.SendMessage(&sarama.ProducerMessage{
			Topic: "publish",
			Value: sarama.ByteEncoder(out),
		})
		h.Log("Next %v %v %v", partition, offset, err)

		sess.MarkMessage(msg, "")
	}
	return nil
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
