package main

import (
	"log"
	"os"

	"learning-platform/internal/kafka"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load(".env")

	broker := os.Getenv("KAFKA_BROKER_INTERNAL")
	if broker == "" {
		broker = "kafka:29092"
	}

	log.Println("[Consumer starting...]")

	consumer := kafka.NewEmailConsumer(broker)
	consumer.Start()
}
