package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

type EmailMessage struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Code    string `json:"code"`
}

type EmailProducer struct {
	writer *kafka.Writer
	queue  chan EmailMessage
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

	p := &EmailProducer{
		writer: w,
		queue:  make(chan EmailMessage, 100),
	}

	go p.worker()

	log.Println("[Producer] Connected to", broker)

	return p
}

func (p *EmailProducer) SendAsync(msg EmailMessage) {
	select {
	case p.queue <- msg:
	default:
		log.Println("[WARN] EmailProducer queue full — message dropped:", msg.Email)
	}
}

func (p *EmailProducer) worker() {
	for msg := range p.queue {
		for {
			err := p.sendToKafka(msg)
			if err != nil {
				log.Println("[Producer Worker] Kafka send failed, retrying:", err)
				time.Sleep(2 * time.Second)
				continue
			}
			break
		}
	}
}

func (p *EmailProducer) sendToKafka(msg EmailMessage) error {
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

	log.Println("[Producer] sent ->", msg.Email)
	return nil
}
