package helpers

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

var (
	createTopic     = CreateTopic
	fetchSecrets    = FetchSecrets
	getPublicIP     = GetPublicIP
	sendMessage     = SendMessage
	getSaramaConfig = sarama.NewConfig
	newSyncProducer = sarama.NewSyncProducer
	newClusterAdmin = sarama.NewClusterAdmin
)

func SendMessage(producer sarama.SyncProducer, topic string, from string, to string, message string) (int32, int64, error) {
	return producer.SendMessage(&sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(fmt.Sprintf("From:%s, To:%s, Message:%s", from, to, message)),
	})
}

func NewClusterAdmin(brokers []string, config *sarama.Config) (sarama.ClusterAdmin, error) {
	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		return admin, fmt.Errorf("failed to create Kafka admin: %v", err)
	}
	return admin, err
}

// createTopic ensures that the given Kafka topic is created if it does not exist
func CreateTopic(topic string, admin sarama.ClusterAdmin) error {
	return admin.CreateTopic(topic, &sarama.TopicDetail{
		NumPartitions:     1, // Parallelism
		ReplicationFactor: 1, // Redundancy and fault tolerance
	}, false)
}

// sendToKafka publishes the message to the Kafka topic
func SendToKafka(ec2Client EC2ClientGetter, secretsClient SecretsManagerClient, from string, to string, message string, region string) error {
	// Fetch the Kafka topic, Kafka port, and EC2 instance ID from AWS Secrets Manager
	topic, instanceID, port, _, _, err := fetchSecrets(secretsClient)
	if err != nil {
		return fmt.Errorf("failed to fetch environment variables from Secrets Manager: %v", err)
	}

	publicIP, err := getPublicIP(ec2Client, instanceID, region)
	if err != nil {
		return fmt.Errorf("failed to get public IP address: %v", err)
	}

	log.Printf("Public IP address of EC2 instance %s: %s\n", instanceID, publicIP)

	// Set the Kafka broker address dynamically
	brokers := []string{fmt.Sprintf("%s:%s", publicIP, port)}

	config := getSaramaConfig()

	admin, err := newClusterAdmin(brokers, config)
	if err != nil {
		return fmt.Errorf("failed to create sarama cluster admin: %v", err)
	}

	// Create the Kafka topic if it doesn't exist
	err = createTopic(topic, admin)

	if err != nil {
		if topicErr, ok := err.(*sarama.TopicError); ok && topicErr.Err == sarama.ErrTopicAlreadyExists {
			log.Printf("Topic %s already exists\n", topic)
		} else {
			return fmt.Errorf("failed to ensure Kafka topic exists: %v", err)
		}
	}
	log.Printf("Topic %s created successfully\n", topic)

	//add code here

	// Set the necessary configuration for the SyncProducer
	config.Producer.Return.Successes = true

	// Create a new SyncProducer
	producer, err := newSyncProducer(brokers, config)
	if err != nil {
		return fmt.Errorf("failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	// Send the message to the Kafka topic
	_, _, err = sendMessage(producer, topic, from, to, message)
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %v", err)
	}

	log.Printf("Message sent to topic %s\n", topic)
	return nil
}
