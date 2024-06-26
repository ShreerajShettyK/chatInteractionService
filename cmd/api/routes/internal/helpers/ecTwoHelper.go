package helpers

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// EC2ClientGetter is an interface that defines methods from the EC2 client we use
type EC2ClientGetter interface {
	DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
}

// getPublicIP retrieves the public IP address of the specified EC2 instance
func GetPublicIP(client EC2ClientGetter, instanceID string, region string) (string, error) {
	var publicIpAddress string

	// Describe the instance to get its public IP address
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	}

	result, err := client.DescribeInstances(context.Background(), input)
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

// EC2ClientLoader is a function type for loading AWS config
type EC2ClientLoader func(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error)

func defaultConfigLoader(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error) {
	return config.LoadDefaultConfig(ctx, optFns...)
}

// NewEC2Client creates a new EC2 client using the provided region and config loader function
func NewEC2Client(region string, loader EC2ClientLoader) (*ec2.Client, error) {
	cfg, err := loader(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS SDK config: %v", err)
	}
	return ec2.NewFromConfig(cfg), nil
}
