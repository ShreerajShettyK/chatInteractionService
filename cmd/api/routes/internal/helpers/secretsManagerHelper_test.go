package helpers

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/stretchr/testify/assert"
)

type MockSecretsManagerClient struct {
	GetSecretValueFunc func(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error)
}

func (m *MockSecretsManagerClient) GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	return m.GetSecretValueFunc(ctx, params, optFns...)
}

func TestFetchSecrets(t *testing.T) {
	mockClient := &MockSecretsManagerClient{}

	// Positive test case
	t.Run("successful retrieval and parsing of secret", func(t *testing.T) {
		secretString := `{
			"KAFKA_TOPIC": "test-topic",
			"EC2_INSTANCE_ID": "i-1234567890abcdef0",
			"KAFKA_PORT": "9092",
			"REGION": "us-west-2",
			"USER_POOL_ID": "us-west-2_abcdef123"
		}`

		mockClient.GetSecretValueFunc = func(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
			return &secretsmanager.GetSecretValueOutput{
				SecretString: aws.String(secretString),
			}, nil
		}

		topic, instanceID, port, region, userPoolID, err := FetchSecrets(mockClient)

		assert.NoError(t, err)
		assert.Equal(t, "test-topic", topic)
		assert.Equal(t, "i-1234567890abcdef0", instanceID)
		assert.Equal(t, "9092", port)
		assert.Equal(t, "us-west-2", region)
		assert.Equal(t, "us-west-2_abcdef123", userPoolID)
	})

	// Negative test cases
	t.Run("error in retrieving secret", func(t *testing.T) {
		mockClient.GetSecretValueFunc = func(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
			return nil, fmt.Errorf("failed to get secret value")
		}

		_, _, _, _, _, err := FetchSecrets(mockClient)

		assert.Error(t, err)
	})

	t.Run("secret string is nil", func(t *testing.T) {
		mockClient.GetSecretValueFunc = func(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
			return &secretsmanager.GetSecretValueOutput{
				SecretString: nil,
			}, nil
		}

		_, _, _, _, _, err := FetchSecrets(mockClient)

		assert.Error(t, err)
		assert.Equal(t, "secret string is nil", err.Error())
	})

	t.Run("error in parsing secret", func(t *testing.T) {
		secretString := `invalid-json-format`

		mockClient.GetSecretValueFunc = func(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
			return &secretsmanager.GetSecretValueOutput{
				SecretString: aws.String(secretString),
			}, nil
		}

		_, _, _, _, _, err := FetchSecrets(mockClient)

		assert.Error(t, err)
	})

	t.Run("missing KAFKA_TOPIC in secret", func(t *testing.T) {
		secretString := `{
			"EC2_INSTANCE_ID": "i-1234567890abcdef0",
			"KAFKA_PORT": "9092",
			"REGION": "us-west-2",
			"USER_POOL_ID": "us-west-2_abcdef123"
		}`

		mockClient.GetSecretValueFunc = func(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
			return &secretsmanager.GetSecretValueOutput{
				SecretString: aws.String(secretString),
			}, nil
		}

		_, _, _, _, _, err := FetchSecrets(mockClient)

		assert.Error(t, err)
		assert.Equal(t, "KAFKA_TOPIC not found in secret data", err.Error())
	})

	t.Run("missing EC2_INSTANCE_ID in secret", func(t *testing.T) {
		secretString := `{
			"KAFKA_TOPIC": "test-topic",
			"KAFKA_PORT": "9092",
			"REGION": "us-west-2",
			"USER_POOL_ID": "us-west-2_abcdef123"
		}`

		mockClient.GetSecretValueFunc = func(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
			return &secretsmanager.GetSecretValueOutput{
				SecretString: aws.String(secretString),
			}, nil
		}

		_, _, _, _, _, err := FetchSecrets(mockClient)

		assert.Error(t, err)
		assert.Equal(t, "EC2_INSTANCE_ID not found in secret data", err.Error())
	})

	t.Run("missing KAFKA_PORT in secret", func(t *testing.T) {
		secretString := `{
			"KAFKA_TOPIC": "test-topic",
			"EC2_INSTANCE_ID": "i-1234567890abcdef0",
			"REGION": "us-west-2",
			"USER_POOL_ID": "us-west-2_abcdef123"
		}`

		mockClient.GetSecretValueFunc = func(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
			return &secretsmanager.GetSecretValueOutput{
				SecretString: aws.String(secretString),
			}, nil
		}
		_, _, _, _, _, err := FetchSecrets(mockClient)

		assert.Error(t, err)
		assert.Equal(t, "KAFKA_PORT not found in secret data", err.Error())
	})

	t.Run("missing REGION in secret", func(t *testing.T) {
		secretString := `{
			"KAFKA_TOPIC": "test-topic",
			"EC2_INSTANCE_ID": "i-1234567890abcdef0",
			"KAFKA_PORT": "9092",
			"USER_POOL_ID": "us-west-2_abcdef123"
		}`

		mockClient.GetSecretValueFunc = func(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
			return &secretsmanager.GetSecretValueOutput{
				SecretString: aws.String(secretString),
			}, nil
		}

		_, _, _, _, _, err := FetchSecrets(mockClient)

		assert.Error(t, err)
		assert.Equal(t, "REGION not found in secret data", err.Error())
	})

	t.Run("missing USER_POOL_ID in secret", func(t *testing.T) {
		secretString := `{
			"KAFKA_TOPIC": "test-topic",
			"EC2_INSTANCE_ID": "i-1234567890abcdef0",
			"KAFKA_PORT": "9092",
			"REGION": "us-west-2"
		}`

		mockClient.GetSecretValueFunc = func(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
			return &secretsmanager.GetSecretValueOutput{
				SecretString: aws.String(secretString),
			}, nil
		}

		topic, instanceID, port, region, userPoolID, err := FetchSecrets(mockClient)

		assert.Error(t, err)
		assert.Equal(t, "USER_POOL_ID not found in secret data", err.Error())
		assert.Equal(t, "test-topic", topic)
		assert.Equal(t, "i-1234567890abcdef0", instanceID)
		assert.Equal(t, "9092", port)
		assert.Equal(t, "us-west-2", region)
		assert.Equal(t, "", userPoolID)
	})
}
