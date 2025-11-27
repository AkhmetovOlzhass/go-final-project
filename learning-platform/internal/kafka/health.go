package kafka

import (
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

func WaitForKafka(broker string) {
	log.Println("Waiting for Kafka...")

	for {
		conn, err := kafka.Dial("tcp", broker)
		if err == nil {
			conn.Close()
			log.Println("Kafka READY:", broker)
			return
		}

		log.Println("Kafka not ready:", err)
		time.Sleep(2 * time.Second)
	}
}
