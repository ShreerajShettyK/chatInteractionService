package service

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

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
