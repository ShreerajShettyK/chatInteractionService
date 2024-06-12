package service

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

// createTopic ensures that the given Kafka topic is created if it does not exist
func createTopic(brokers []string, topic string, config *sarama.Config) error {
	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		return fmt.Errorf("failed to create Kafka admin: %v", err)
	}
	defer admin.Close()
	//ensures the admin client is closed when the function exits

	topicDetail := &sarama.TopicDetail{
		NumPartitions:     1, //parallelism
		ReplicationFactor: 1, //Replication ensures that data is copied to multiple brokers, providing redundancy and fault tolerance
	}

	err = admin.CreateTopic(topic, topicDetail, false)
	if err != nil {
		// Check if the error is a *sarama.TopicError and whether the topic already exists
		if topicErr, ok := err.(*sarama.TopicError); ok && topicErr.Err == sarama.ErrTopicAlreadyExists {
			log.Printf("Topic %s already exists\n", topic)
			return nil
		}
		return fmt.Errorf("failed to create Kafka topic: %v", err)
	}

	log.Printf("Topic %s created successfully\n", topic)
	return nil
}

// sendToKafka publishes the message to the Kafka topic
func sendToKafka(from string, to string, message string) error {
	// Fetch the Kafka topic, Kafka port, and EC2 instance ID from AWS Secrets Manager
	topic, instanceID, port, _, _, err := fetchSecrets()
	if err != nil {
		return fmt.Errorf("failed to fetch environment variables from Secrets Manager: %v", err)
	}

	publicIP, err := getPublicIP(instanceID)
	if err != nil {
		return fmt.Errorf("failed to get public IP address: %v", err)
	}

	log.Printf("Public IP address of EC2 instance %s: %s\n", instanceID, publicIP)

	// Set the Kafka broker address dynamically
	brokers := []string{fmt.Sprintf("%s:%s", publicIP, port)}

	config := sarama.NewConfig()

	// Create the Kafka topic if it doesn't exist
	err = createTopic(brokers, topic, config)
	if err != nil {
		return fmt.Errorf("failed to ensure Kafka topic exists: %v", err)
	}

	// Set the necessary configuration for the SyncProducer
	// config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	// Create a new SyncProducer
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return fmt.Errorf("failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	// Send the message to the Kafka topic
	_, _, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(fmt.Sprintf("From:%s, To:%s, Message:%s", from, to, message)),
	})
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %v", err)
	}

	log.Printf("Message sent to topic %s\n", topic)
	return nil
}
