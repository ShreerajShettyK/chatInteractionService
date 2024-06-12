package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

var cfg aws.Config

func init() {
	// Load AWS SDK configuration
	var err error
	cfg, err = config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Error loading AWS SDK config: %v", err)
	}
}

// Fetches the value of a secret from AWS Secrets Manager
func getSecret(secretName string) (string, error) {
	var topicSecret string
	client := secretsmanager.NewFromConfig(cfg)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	}

	result, err := client.GetSecretValue(context.Background(), input)
	if err != nil {
		return topicSecret, fmt.Errorf("failed to get secret value: %v", err)
	}

	if result.SecretString == nil {
		return topicSecret, fmt.Errorf("secret string is nil")
	}
	topicSecret = aws.ToString(result.SecretString)

	return topicSecret, nil
}

// Fetches Kafka topic, EC2 instance ID, and Kafka port from Secrets Manager
func fetchSecrets() (string, string, string, string, string, error) {
	var secretData map[string]string
	topicSecret, err := getSecret("myApp/mongo-db-credentials")
	if err != nil {
		return "", "", "", "", "", err
	}
	log.Println("Secret retrieved from AWS Secrets Manager")

	// log.Printf("Secret retrieved from AWS Secrets Manager: %s\n", topicSecret)

	// Parsing the JSON format of the secret string
	err = json.Unmarshal([]byte(topicSecret), &secretData)
	if err != nil {
		log.Printf("Error parsing secret string: %v\n", err)
		return "", "", "", "", "", err
	}

	// Extracting values from the parsed secret data
	topic, ok := secretData["KAFKA_TOPIC"]
	if !ok {
		return "", "", "", "", "", fmt.Errorf("KAFKA_TOPIC not found in secret data")
	}

	instanceID, ok := secretData["EC2_INSTANCE_ID"]
	if !ok {
		return "", "", "", "", "", fmt.Errorf("EC2_INSTANCE_ID not found in secret data")
	}

	port, ok := secretData["KAFKA_PORT"]
	if !ok {
		return "", "", "", "", "", fmt.Errorf("KAFKA_PORT not found in secret data")
	}

	region, ok := secretData["REGION"]
	if !ok {
		return "", "", "", "", "", fmt.Errorf("REGION not found in secret data")
	}

	userPoolID, ok := secretData["USER_POOL_ID"]
	if !ok {
		return "", "", "", "", "", fmt.Errorf("USER_POOL_ID not found in secret data")
	}

	// log.Printf("Parsed values: Topic=%s, InstanceID=%s, Port=%s, Region=%s, UserPoolID=%s\n", topic, instanceID, port, region, userPoolID)

	return topic, instanceID, port, region, userPoolID, nil
}
