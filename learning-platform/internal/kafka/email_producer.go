package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

type EmailProducer struct {
	writer *kafka.Writer
}

func NewEmailProducer() *EmailProducer {
	broker := os.Getenv("KAFKA_BROKER_EXTERNAL")
	if broker == "" {
		broker = "localhost:9092"
	}

	w := &kafka.Writer{
		Addr:                   kafka.TCP(broker),
		Topic:                  "email_send",
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}

	log.Println("[Producer] Connected to", broker)

	return &EmailProducer{writer: w}
}

type EmailMessage struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Code    string `json:"code"`
}

func (p *EmailProducer) Send(msg EmailMessage) error {
	body, _ := json.Marshal(msg)

	err := p.writer.WriteMessages(
		context.Background(),
		kafka.Message{
			Key:   []byte(msg.Email),
			Value: body,
			Time:  time.Now(),
		},
	)

	if err != nil {
		log.Println("[Producer ERROR]", err)
		return err
	}

	log.Println("[Producer] sent →", msg.Email)
	return nil
}
