//////kafka topic,public ip from env

package service

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/IBM/sarama"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/joho/godotenv"
)

func init() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

// getPublicIP retrieves the public IP address of the specified EC2 instance
func getPublicIP(instanceID string) (string, error) {
	// Load the AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion("us-east-1"))
	if err != nil {
		return "", fmt.Errorf("failed to load AWS SDK config: %v", err)
	}

	// Create an EC2 client
	svc := ec2.NewFromConfig(cfg)

	// Describe the instance to get its public IP address
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	}

	result, err := svc.DescribeInstances(context.Background(), input)
	if err != nil {
		return "", fmt.Errorf("failed to describe EC2 instance: %v", err)
	}

	if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
		return "", fmt.Errorf("no instances found for instance ID %s", instanceID)
	}

	instance := result.Reservations[0].Instances[0]
	if instance.PublicIpAddress == nil {
		return "", fmt.Errorf("instance %s does not have a public IP address", instanceID)
	}

	return *instance.PublicIpAddress, nil
}

// createTopic ensures that the given Kafka topic is created if it does not exist
func createTopic(brokers []string, topic string) error {
	config := sarama.NewConfig()
	admin, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		return fmt.Errorf("failed to create Kafka admin: %v", err)
	}
	defer admin.Close()

	topicDetail := &sarama.TopicDetail{
		NumPartitions:     1,
		ReplicationFactor: 1,
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
	brokers := []string{"3.85.126.195:9092"}

	// Fetch the topic from environment variables
	topic := os.Getenv("KAFKA_TOPIC")
	if topic == "" {
		return fmt.Errorf("KAFKA_TOPIC environment variable is not set")
	}

	// Create the topic if it doesn't exist
	err := createTopic(brokers, topic)
	if err != nil {
		return fmt.Errorf("failed to ensure Kafka topic exists: %v", err)
	}

	// Fetch the instance ID from environment variables
	instanceID := os.Getenv("EC2_INSTANCE_ID")
	if instanceID == "" {
		return fmt.Errorf("EC2_INSTANCE_ID environment variable is not set")
	}

	// Fetch the public IP address of the EC2 instance
	publicIP, err := getPublicIP(instanceID)
	if err != nil {
		return fmt.Errorf("failed to get public IP address: %v", err)
	}

	log.Printf("Public IP address of EC2 instance %s: %s\n", instanceID, publicIP)

	// Set the necessary configuration for the SyncProducer
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return fmt.Errorf("failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(fmt.Sprintf("From:%s, To:%s, Message:%s", from, to, message)),
	}

	_, _, err = producer.SendMessage(msg)
	if err != nil {
		return fmt.Errorf("failed to send message to Kafka: %v", err)
	}

	log.Printf("Message sent to topic %s\n", topic)
	return nil
}
