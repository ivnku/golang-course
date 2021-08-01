package main

import (
	"github.com/Shopify/sarama"
)

type PublishConsumerHandler struct {
	Log func(string, ...interface{})
}

func (h *PublishConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error {
	h.Log("Setup happend")
	return nil
}

func (h *PublishConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	h.Log("Cleanup happend")
	return nil
}

func (h *PublishConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.Log("Processing %v/%v/%v %s %s", msg.Topic, msg.Partition, msg.Offset,
			msg.Key,
			msg.Value)
		sess.MarkMessage(msg, "")
	}
	return nil
}

// ---

type NotifyEmailConsumerHandler struct {
	Log func(string, ...interface{})
}

func (h *NotifyEmailConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error {
	h.Log("Setup happend")
	return nil
}

func (h *NotifyEmailConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	h.Log("Cleanup happend")
	return nil
}

func (h *NotifyEmailConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.Log("Processing %v/%v/%v %s %s", msg.Topic, msg.Partition, msg.Offset,
			msg.Key,
			msg.Value)
		sess.MarkMessage(msg, "")
	}
	return nil
}

// ---

type UpdateSearchConsumerHandler struct {
	Log func(string, ...interface{})
}

func (h *UpdateSearchConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error {
	h.Log("Setup happend")
	return nil
}

func (h *UpdateSearchConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	h.Log("Cleanup happend")
	return nil
}

func (h *UpdateSearchConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		h.Log("Processing %v/%v/%v %s %s", msg.Topic, msg.Partition, msg.Offset,
			msg.Key,
			msg.Value)
		sess.MarkMessage(msg, "")
	}
	return nil
}
