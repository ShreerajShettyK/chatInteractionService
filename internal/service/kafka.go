package service

import (
	"fmt"

	"github.com/IBM/sarama"
)

// sendToKafka publishes the message to the Kafka topic
func sendToKafka(from string, to string, message string) error {
	producer, err := sarama.NewSyncProducer([]string{"localhost:9092"}, nil)
	if err != nil {
		return fmt.Errorf("failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: "chat-topic-2",
		Value: sarama.StringEncoder(fmt.Sprintf("From:%s, To:%s, Message:%s", from, to, message)),
	}

	_, _, err = producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %v", err)
	}

	return nil
}
