package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"learning-platform/internal/email"

	"github.com/segmentio/kafka-go"
)

type EmailConsumer struct {
	reader *kafka.Reader
	sender *email.SMTPSender
}

func NewEmailConsumer(broker string) *EmailConsumer {
	WaitForKafka(broker)

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   "email_send",
		GroupID: "email_service",
	})

	return &EmailConsumer{
		reader: reader,
		sender: email.NewSMTPSender(),
	}
}

func (c *EmailConsumer) Start() {
	log.Println("Email consumer started")

	for {
		msg, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			log.Println("Kafka read error:", err)
			continue
		}

		var em EmailMessage
		if err := json.Unmarshal(msg.Value, &em); err != nil {
			log.Println("JSON error:", err)
			continue
		}

		body := fmt.Sprintf("Your code: %s", em.Code)

		if err := c.sender.SendEmail(em.Email, em.Subject, body); err != nil {
			log.Println("Email error:", err)
			continue
		}

		log.Println("Sent to:", em.Email)
	}
}
