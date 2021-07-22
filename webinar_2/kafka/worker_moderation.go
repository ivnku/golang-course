package main

import (
	"encoding/json"
	"math/rand"

	"github.com/Shopify/sarama"
)

type ModerationConsumerHandler struct {
	Log  func(string, ...interface{})
	Next sarama.SyncProducer
}

func (h *ModerationConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error {
	h.Log("Setup happend")
	return nil
}

func (h *ModerationConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	h.Log("Cleanup happend")
	return nil
}

func (h *ModerationConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.Log("Processing %v/%v/%v %s %s", msg.Topic, msg.Partition, msg.Offset,
			msg.Key,
			msg.Value)

		in := &CreationIn{}
		err := json.Unmarshal(msg.Value, in)

		in.SpamProp = rand.Intn(100)

		out, _ := json.Marshal(in)
		partition, offset, err := h.Next.SendMessage(&sarama.ProducerMessage{
			Topic: "create",
			Value: sarama.ByteEncoder(out),
		})
		h.Log("Next %v %v %v", partition, offset, err)

		sess.MarkMessage(msg, "")
	}
	return nil
}
