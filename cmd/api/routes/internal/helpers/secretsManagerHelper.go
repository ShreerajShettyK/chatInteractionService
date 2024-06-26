package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

type SecretsManagerClient interface {
	GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

var cfg aws.Config

func init() {
	var err error
	cfg, err = config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalf("Error loading AWS SDK config: %v", err)
	}
}

func getSecret(client SecretsManagerClient, secretName string) (string, error) {
	input := &secretsmanager.GetSecretValueInput{
		SecretId: &secretName,
	}

	result, err := client.GetSecretValue(context.Background(), input)
	if err != nil {
		return "", fmt.Errorf("failed to get secret value: %v", err)
	}

	if result.SecretString == nil {
		return "", fmt.Errorf("secret string is nil")
	}

	return aws.ToString(result.SecretString), nil
}

func FetchSecrets(client SecretsManagerClient) (string, string, string, string, string, error) {
	topicSecret, err := getSecret(client, "myApp/mongo-db-credentials")
	if err != nil {
		return "", "", "", "", "", err
	}

	var secretData map[string]string
	err = json.Unmarshal([]byte(topicSecret), &secretData)
	if err != nil {
		return "", "", "", "", "", fmt.Errorf("error parsing secret string: %v", err)
	}

	topic, ok := secretData["KAFKA_TOPIC"]
	if !ok {
		return "", "", "", "", "", fmt.Errorf("KAFKA_TOPIC not found in secret data")
	}

	instanceID, ok := secretData["EC2_INSTANCE_ID"]
	if !ok {
		return topic, "", "", "", "", fmt.Errorf("EC2_INSTANCE_ID not found in secret data")
	}

	port, ok := secretData["KAFKA_PORT"]
	if !ok {
		return topic, instanceID, "", "", "", fmt.Errorf("KAFKA_PORT not found in secret data")
	}

	region, ok := secretData["REGION"]
	if !ok {
		return topic, instanceID, port, "", "", fmt.Errorf("REGION not found in secret data")
	}

	userPoolID, ok := secretData["USER_POOL_ID"]
	if !ok {
		return topic, instanceID, port, region, "", fmt.Errorf("USER_POOL_ID not found in secret data")
	}

	return topic, instanceID, port, region, userPoolID, nil
}
