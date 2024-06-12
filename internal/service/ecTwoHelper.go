package service

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// getPublicIP retrieves the public IP address of the specified EC2 instance
func getPublicIP(instanceID string, region string) (string, error) {
	var publicIpAddress string

	// Load the AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		log.Printf("failed to load AWS SDK config: %v", err)
		return publicIpAddress, err
	}

	// Create an EC2 client
	svc := ec2.NewFromConfig(cfg)

	// Describe the instance to get its public IP address
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	}

	result, err := svc.DescribeInstances(context.Background(), input)
	if err != nil {
		log.Printf("failed to describe EC2 instance: %v", err)
		return publicIpAddress, err
	}

	if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
		log.Printf("no instances found for instance ID %s", instanceID)
		return publicIpAddress, fmt.Errorf("no instances found for instance ID %s", instanceID)
	}

	// Check if the instance has a public IP address
	if result.Reservations[0].Instances[0].PublicIpAddress == nil {
		log.Printf("instance %s does not have a public IP address", instanceID)
		return publicIpAddress, fmt.Errorf("instance %s does not have a public IP address", instanceID)
	}

	// Retrieve the public IP address
	publicIpAddress = aws.ToString(result.Reservations[0].Instances[0].PublicIpAddress)
	return publicIpAddress, nil
}
