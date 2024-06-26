package helpers

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEC2Client is a mock implementation of EC2ClientGetter
type MockEC2Client struct {
	mock.Mock
}

// DescribeInstances mocks the DescribeInstances method of the EC2 client
func (m *MockEC2Client) DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	args := m.Called(ctx, params)
	if args.Error(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ec2.DescribeInstancesOutput), nil
}

// Mock for config.LoadDefaultConfig to simulate failures
func mockLoadDefaultConfig(ctx context.Context, optFns ...func(*config.LoadOptions) error) (aws.Config, error) {
	return aws.Config{}, errors.New("failed to load AWS SDK config")
}

func TestGetPublicIP(t *testing.T) {
	instanceID := "i-07a5d5ffeb7a8422d" // Valid instance ID format
	region := "us-east-1"

	// Mock EC2 client
	mockSvc := new(MockEC2Client)

	// Mock response for DescribeInstances call failure
	mockSvc.On("DescribeInstances", mock.Anything, mock.Anything).Return(nil, errors.New("describe instances error")).Once()

	// Test DescribeInstances call failure case
	ip, err := GetPublicIP(mockSvc, instanceID, region)
	expectedErr := "describe instances error"
	assert.NotNil(t, err)
	assert.Equal(t, expectedErr, err.Error())
	assert.Empty(t, ip)

	// Mock response for instance not found case (empty response)
	mockRespNotFound := &ec2.DescribeInstancesOutput{
		Reservations: []types.Reservation{},
	}
	mockSvc.On("DescribeInstances", mock.Anything, mock.Anything).Return(mockRespNotFound, nil).Once()

	// Test instance not found case
	ip, err = GetPublicIP(mockSvc, instanceID, region)
	expectedErr = "no instances found for instance ID " + instanceID
	assert.NotNil(t, err)
	assert.Equal(t, expectedErr, err.Error())
	assert.Empty(t, ip)

	// Mock response for positive case
	mockResp := &ec2.DescribeInstancesOutput{
		Reservations: []types.Reservation{
			{
				Instances: []types.Instance{
					{
						InstanceId:      aws.String(instanceID),
						PublicIpAddress: aws.String("1.2.3.4"),
					},
				},
			},
		},
	}
	mockSvc.On("DescribeInstances", mock.Anything, mock.Anything).Return(mockResp, nil).Once()

	// Test positive case
	ip, err = GetPublicIP(mockSvc, instanceID, region)
	assert.Nil(t, err)
	assert.Equal(t, "1.2.3.4", ip)

	// Test instance without public IP
	mockRespNoPublicIP := &ec2.DescribeInstancesOutput{
		Reservations: []types.Reservation{
			{
				Instances: []types.Instance{
					{
						InstanceId: aws.String(instanceID),
					},
				},
			},
		},
	}
	mockSvc.On("DescribeInstances", mock.Anything, mock.Anything).Return(mockRespNoPublicIP, nil).Once()

	ip, err = GetPublicIP(mockSvc, instanceID, region)
	expectedErr = "instance " + instanceID + " does not have a public IP address"
	assert.NotNil(t, err)
	assert.Equal(t, expectedErr, err.Error())
	assert.Empty(t, ip)
}

func TestNewEC2Client(t *testing.T) {
	region := "us-east-1"

	// Test successful creation of EC2 client with default loader
	client, err := NewEC2Client(region, defaultConfigLoader)
	assert.Nil(t, err)
	assert.NotNil(t, client)

	// Test failure to load AWS SDK config
	client, err = NewEC2Client(region, mockLoadDefaultConfig)
	expectedErr := "failed to load AWS SDK config: failed to load AWS SDK config"
	assert.NotNil(t, err)
	assert.Equal(t, expectedErr, err.Error())
	assert.Nil(t, client)
}
