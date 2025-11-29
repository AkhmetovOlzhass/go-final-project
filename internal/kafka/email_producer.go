package kafka

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel"
)

type EmailMessage struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Code    string `json:"code"`
}

type queuedMessage struct {
	ctx context.Context
	msg EmailMessage
}

type EmailProducer struct {
	writer *kafka.Writer
	queue  chan queuedMessage
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
		queue:  make(chan queuedMessage, 100),
	}

	go p.worker()

	log.Println("[Producer] Connected to", broker)

	return p
}

func (p *EmailProducer) SendAsync(msg EmailMessage) {
	detachedCtx := context.Background()

	select {
	case p.queue <- queuedMessage{ctx: detachedCtx , msg: msg}:
	default:
		log.Println("[WARN] EmailProducer queue full — dropped:", msg.Email)
	}
}

func (p *EmailProducer) worker() {
	for qm := range p.queue {
		for {
			err := p.sendToKafka(qm.ctx, qm.msg)
			if err != nil {
				log.Println("[Producer Worker] Kafka send failed, retrying:", err)
				time.Sleep(2 * time.Second)
				continue
			}
			break
		}
	}
}

func (p *EmailProducer) sendToKafka(ctx context.Context, msg EmailMessage) error {
	ctx, span := otel.Tracer("kafka").Start(ctx, "Producer.WriteMessage")
	defer span.End()

	body, _ := json.Marshal(msg)

	err := p.writer.WriteMessages(
		ctx,
		kafka.Message{
			Key:   []byte(msg.Email),
			Value: body,
			Time:  time.Now(),
		},
	)

	if err != nil {
		span.RecordError(err)
		return err
	}
	return nil
}
