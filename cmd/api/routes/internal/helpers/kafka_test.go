package helpers

import (
	"testing"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"
)

// type MockEC2Client struct{}
// type MockSecretsManagerClient struct{}

// func (m *MockEC2Client) DescribeInstances() {}

type MockSaramaSyncProducer struct{}

type mockClusterAdmin struct{}

// AlterClientQuotas implements sarama.ClusterAdmin.
func (m mockClusterAdmin) AlterClientQuotas(entity []sarama.QuotaEntityComponent, op sarama.ClientQuotasOp, validateOnly bool) error {
	return nil
}

// AlterConfig implements sarama.ClusterAdmin.
func (m mockClusterAdmin) AlterConfig(resourceType sarama.ConfigResourceType, name string, entries map[string]*string, validateOnly bool) error {
	return nil
}

// AlterPartitionReassignments implements sarama.ClusterAdmin.
func (m mockClusterAdmin) AlterPartitionReassignments(topic string, assignment [][]int32) error {
	return nil
}

// Close implements sarama.ClusterAdmin.
func (m mockClusterAdmin) Close() error {
	return nil
}

// Controller implements sarama.ClusterAdmin.
func (m mockClusterAdmin) Controller() (*sarama.Broker, error) {
	return nil, nil
}

// CreateACL implements sarama.ClusterAdmin.
func (m mockClusterAdmin) CreateACL(resource sarama.Resource, acl sarama.Acl) error {
	return nil
}

// CreateACLs implements sarama.ClusterAdmin.
func (m mockClusterAdmin) CreateACLs([]*sarama.ResourceAcls) error {
	return nil
}

// CreatePartitions implements sarama.ClusterAdmin.
func (m mockClusterAdmin) CreatePartitions(topic string, count int32, assignment [][]int32, validateOnly bool) error {
	return nil
}

// CreateTopic implements sarama.ClusterAdmin.
func (m mockClusterAdmin) CreateTopic(topic string, detail *sarama.TopicDetail, validateOnly bool) error {
	return nil
}

// DeleteACL implements sarama.ClusterAdmin.
func (m mockClusterAdmin) DeleteACL(filter sarama.AclFilter, validateOnly bool) ([]sarama.MatchingAcl, error) {
	return nil, nil
}

// DeleteConsumerGroup implements sarama.ClusterAdmin.
func (m mockClusterAdmin) DeleteConsumerGroup(group string) error {
	return nil
}

// DeleteConsumerGroupOffset implements sarama.ClusterAdmin.
func (m mockClusterAdmin) DeleteConsumerGroupOffset(group string, topic string, partition int32) error {
	return nil
}

// DeleteRecords implements sarama.ClusterAdmin.
func (m mockClusterAdmin) DeleteRecords(topic string, partitionOffsets map[int32]int64) error {
	return nil
}

// DeleteTopic implements sarama.ClusterAdmin.
func (m mockClusterAdmin) DeleteTopic(topic string) error {
	return nil
}

// DeleteUserScramCredentials implements sarama.ClusterAdmin.
func (m mockClusterAdmin) DeleteUserScramCredentials(delete []sarama.AlterUserScramCredentialsDelete) ([]*sarama.AlterUserScramCredentialsResult, error) {
	return nil, nil
}

// DescribeClientQuotas implements sarama.ClusterAdmin.
func (m mockClusterAdmin) DescribeClientQuotas(components []sarama.QuotaFilterComponent, strict bool) ([]sarama.DescribeClientQuotasEntry, error) {
	return nil, nil
}

// DescribeCluster implements sarama.ClusterAdmin.
func (m mockClusterAdmin) DescribeCluster() (brokers []*sarama.Broker, controllerID int32, err error) {
	return nil, int32(1), nil
}

// DescribeConfig implements sarama.ClusterAdmin.
func (m mockClusterAdmin) DescribeConfig(resource sarama.ConfigResource) ([]sarama.ConfigEntry, error) {
	return nil, nil
}

// DescribeConsumerGroups implements sarama.ClusterAdmin.
func (m mockClusterAdmin) DescribeConsumerGroups(groups []string) ([]*sarama.GroupDescription, error) {
	return nil, nil
}

// DescribeLogDirs implements sarama.ClusterAdmin.
func (m mockClusterAdmin) DescribeLogDirs(brokers []int32) (map[int32][]sarama.DescribeLogDirsResponseDirMetadata, error) {
	return nil, nil
}

// DescribeTopics implements sarama.ClusterAdmin.
func (m mockClusterAdmin) DescribeTopics(topics []string) (metadata []*sarama.TopicMetadata, err error) {
	return nil, nil
}

// DescribeUserScramCredentials implements sarama.ClusterAdmin.
func (m mockClusterAdmin) DescribeUserScramCredentials(users []string) ([]*sarama.DescribeUserScramCredentialsResult, error) {
	return nil, nil
}

// IncrementalAlterConfig implements sarama.ClusterAdmin.
func (m mockClusterAdmin) IncrementalAlterConfig(resourceType sarama.ConfigResourceType, name string, entries map[string]sarama.IncrementalAlterConfigsEntry, validateOnly bool) error {
	return nil
}

// ListAcls implements sarama.ClusterAdmin.
func (m mockClusterAdmin) ListAcls(filter sarama.AclFilter) ([]sarama.ResourceAcls, error) {
	return nil, nil
}

// ListConsumerGroupOffsets implements sarama.ClusterAdmin.
func (m mockClusterAdmin) ListConsumerGroupOffsets(group string, topicPartitions map[string][]int32) (*sarama.OffsetFetchResponse, error) {
	return nil, nil
}

// ListConsumerGroups implements sarama.ClusterAdmin.
func (m mockClusterAdmin) ListConsumerGroups() (map[string]string, error) {
	return nil, nil
}

// ListPartitionReassignments implements sarama.ClusterAdmin.
func (m mockClusterAdmin) ListPartitionReassignments(topics string, partitions []int32) (topicStatus map[string]map[int32]*sarama.PartitionReplicaReassignmentsStatus, err error) {
	return nil, nil
}

// ListTopics implements sarama.ClusterAdmin.
func (m mockClusterAdmin) ListTopics() (map[string]sarama.TopicDetail, error) {
	return nil, nil
}

// RemoveMemberFromConsumerGroup implements sarama.ClusterAdmin.
func (m mockClusterAdmin) RemoveMemberFromConsumerGroup(groupId string, groupInstanceIds []string) (*sarama.LeaveGroupResponse, error) {
	return nil, nil
}

// UpsertUserScramCredentials implements sarama.ClusterAdmin.
func (m mockClusterAdmin) UpsertUserScramCredentials(upsert []sarama.AlterUserScramCredentialsUpsert) ([]*sarama.AlterUserScramCredentialsResult, error) {
	return nil, nil
}

// Mock fields to track state for transactional methods
var (
	isTransactional          bool
	currentTxnStatus         sarama.ProducerTxnStatusFlag
	txnMessageCount          int
	txnOffsets               map[string][]*sarama.PartitionOffsetMetadata
	txnGroupId               string
	txnCommitCalled          bool
	txnAbortCalled           bool
	txnAddMessageToTxnCalled bool
	txnAddOffsetsToTxnCalled bool
)

// Reset mock transactional state before each test
func resetMockTxnState() {
	isTransactional = false
	currentTxnStatus = sarama.ProducerTxnStatusFlag(0)
	txnMessageCount = 0
	txnOffsets = make(map[string][]*sarama.PartitionOffsetMetadata)
	txnGroupId = ""
	txnCommitCalled = false
	txnAbortCalled = false
	txnAddMessageToTxnCalled = false
	txnAddOffsetsToTxnCalled = false
}

// AbortTxn implements sarama.SyncProducer.
func (m *MockSaramaSyncProducer) AbortTxn() error {
	txnAbortCalled = true
	currentTxnStatus = sarama.ProducerTxnStatusFlag(0)
	return nil
}

// AddMessageToTxn implements sarama.SyncProducer.
func (m *MockSaramaSyncProducer) AddMessageToTxn(msg *sarama.ConsumerMessage, groupId string, metadata *string) error {
	txnAddMessageToTxnCalled = true
	txnGroupId = groupId
	txnMessageCount++
	return nil
}

// AddOffsetsToTxn implements sarama.SyncProducer.
func (m *MockSaramaSyncProducer) AddOffsetsToTxn(offsets map[string][]*sarama.PartitionOffsetMetadata, groupId string) error {
	txnAddOffsetsToTxnCalled = true
	txnGroupId = groupId
	txnOffsets = offsets
	return nil
}

// BeginTxn implements sarama.SyncProducer.
func (m *MockSaramaSyncProducer) BeginTxn() error {
	isTransactional = true
	currentTxnStatus = sarama.ProducerTxnStatusFlag(1)
	return nil
}

// CommitTxn implements sarama.SyncProducer.
func (m *MockSaramaSyncProducer) CommitTxn() error {
	txnCommitCalled = true
	currentTxnStatus = sarama.ProducerTxnStatusFlag(0)
	return nil
}

// IsTransactional implements sarama.SyncProducer.
func (m *MockSaramaSyncProducer) IsTransactional() bool {
	return isTransactional
}

// TxnStatus implements sarama.SyncProducer.
func (m *MockSaramaSyncProducer) TxnStatus() sarama.ProducerTxnStatusFlag {
	return currentTxnStatus
}

// SendMessage is a mock implementation of the SyncProducer's SendMessage method
func (m *MockSaramaSyncProducer) SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	return 0, 0, nil
}

// SendMessages is a mock implementation of the SyncProducer's SendMessages method
func (m *MockSaramaSyncProducer) SendMessages(msgs []*sarama.ProducerMessage) error {
	return nil
}

// Close is a mock implementation of the SyncProducer's Close method
func (m *MockSaramaSyncProducer) Close() error {
	return nil
}

var mockCreateTopic func(brokers []string, topic string, config *sarama.Config) error

func TestSendMessage(t *testing.T) {
	t.Run("test for send message", func(t *testing.T) {
		value1, value2, err := SendMessage(&MockSaramaSyncProducer{}, "test", "test", "test", "test")
		assert.Nil(t, err)
		assert.GreaterOrEqual(t, int32(1), value1)
		assert.GreaterOrEqual(t, int64(1), value2)
	})
}

// Test cases for createTopic function
func TestCreateTopic(t *testing.T) {
	t.Run("test for successful topic creation", func(t *testing.T) {
		err := createTopic("test-topic", mockClusterAdmin{})
		assert.Nil(t, err)
	})

	// t.Run("test for topic creation failure", func(t *testing.T) {
	// 	newClusterAdmin = func(addrs []string, conf *sarama.Config) (sarama.ClusterAdmin, error) {
	// 		return nil, assert.AnError
	// 	}
	// 	err := createTopic("test-topic", mockClusterAdmin{})
	// 	assert.NotNil(t, err)
	// 	assert.Equal(t, "failed to create Kafka topic", err.Error())
	// })

	// t.Run("test for topic already exists", func(t *testing.T) {
	// 	mockCreateTopic = func(brokers []string, topic string, config *sarama.Config) error {
	// 		return nil
	// 	}
	// 	newClusterAdmin = func(addrs []string, conf *sarama.Config) (sarama.ClusterAdmin, error) {
	// 		return mockClusterAdmin{}, nil
	// 	}

	// 	err := createTopic("test-topic", mockClusterAdmin{})
	// 	assert.Nil(t, err)
	// })
}

// // Additional test cases for SendToKafka with mock createTopic
// func TestSendToKafkaWithCreateTopic(t *testing.T) {
// 	// t.Run("test for sendToKafka with successful topic creation", func(t *testing.T) {
// 	// 	mockCreateTopic = func(brokers []string, topic string, config *sarama.Config) error {
// 	// 		return nil
// 	// 	}
// 	// 	fetchSecrets = func(client SecretsManagerClient) (string, string, string, string, string, error) {
// 	// 		return "test-topic", "instance-id", "9092", "", "", nil
// 	// 	}
// 	// 	getPublicIP = func(client EC2ClientGetter, instanceID, region string) (string, error) {
// 	// 		return "127.0.0.1", nil
// 	// 	}
// 	// 	sendMessage = func(producer sarama.SyncProducer, topic, from, to, message string) (int32, int64, error) {
// 	// 		return 0, 0, nil
// 	// 	}
// 	// 	newSyncProducer = func(addrs []string, config *sarama.Config) (sarama.SyncProducer, error) {
// 	// 		return &MockSaramaSyncProducer{}, nil
// 	// 	}
// 	// 	err := SendToKafka(&MockEC2Client{}, &MockSecretsManagerClient{}, "from", "to", "message", "us-east-1")
// 	// 	assert.Nil(t, err)
// 	// })

// 	t.Run("test for sendToKafka with topic creation failure", func(t *testing.T) {
// 		mockCreateTopic = func(brokers []string, topic string, config *sarama.Config) error {
// 			return errors.New("failed to create Kafka topic")
// 		}
// 		fetchSecrets = func(client SecretsManagerClient) (string, string, string, string, string, error) {
// 			return "test-topic", "instance-id", "9092", "", "", nil
// 		}
// 		getPublicIP = func(client EC2ClientGetter, instanceID, region string) (string, error) {
// 			return "127.0.0.1", nil
// 		}
// 		sendMessage = func(producer sarama.SyncProducer, topic, from, to, message string) (int32, int64, error) {
// 			return 0, 0, nil
// 		}
// 		newSyncProducer = func(addrs []string, config *sarama.Config) (sarama.SyncProducer, error) {
// 			return &MockSaramaSyncProducer{}, nil
// 		}
// 		err := SendToKafka(&MockEC2Client{}, &MockSecretsManagerClient{}, "from", "to", "message", "us-east-1")
// 		assert.NotNil(t, err)
// 		assert.Equal(t, "failed to ensure Kafka topic exists: failed to create Kafka topic", err.Error())
// 	})
// }

func TestSendToKafka(t *testing.T) {
	t.Run("test for sending message to kafka with error fetching secrets", func(t *testing.T) {
		resetMockTxnState()
		fetchSecrets = func(client SecretsManagerClient) (string, string, string, string, string, error) {
			return "", "", "", "", "", assert.AnError
		}
		getPublicIP = func(client EC2ClientGetter, instanceID, region string) (string, error) {
			return "test", nil
		}
		newClusterAdmin = func(addrs []string, conf *sarama.Config) (sarama.ClusterAdmin, error) {
			return mockClusterAdmin{}, nil
		}
		createTopic = func(topic string, admin sarama.ClusterAdmin) error {
			return nil
		}
		sendMessage = func(producer sarama.SyncProducer, topic, from, to, message string) (int32, int64, error) {
			return int32(0), int64(0), nil
		}
		newSyncProducer = func(addrs []string, config *sarama.Config) (sarama.SyncProducer, error) {
			return &MockSaramaSyncProducer{}, nil
		}
		err := SendToKafka(&MockEC2Client{}, &MockSecretsManagerClient{}, "test", "test", "hello world", "us-east-1")
		assert.NotNil(t, err)
	})

	t.Run("test for sending message to kafka with error getting public ip", func(t *testing.T) {
		resetMockTxnState()
		fetchSecrets = func(client SecretsManagerClient) (string, string, string, string, string, error) {
			return "test", "test", "9092", "test", "test", nil
		}
		getPublicIP = func(client EC2ClientGetter, instanceID, region string) (string, error) {
			return "", assert.AnError
		}
		newClusterAdmin = func(addrs []string, conf *sarama.Config) (sarama.ClusterAdmin, error) {
			return mockClusterAdmin{}, nil
		}
		createTopic = func(topic string, admin sarama.ClusterAdmin) error {
			return nil
		}
		sendMessage = func(producer sarama.SyncProducer, topic, from, to, message string) (int32, int64, error) {
			return int32(0), int64(0), nil
		}
		newSyncProducer = func(addrs []string, config *sarama.Config) (sarama.SyncProducer, error) {
			return &MockSaramaSyncProducer{}, nil
		}
		err := SendToKafka(&MockEC2Client{}, &MockSecretsManagerClient{}, "test", "test", "hello world", "us-east-1")
		assert.NotNil(t, err)
	})

	t.Run("test for failed to create Kafka producer", func(t *testing.T) {
		resetMockTxnState()
		fetchSecrets = func(client SecretsManagerClient) (string, string, string, string, string, error) {
			return "test", "test", "9092", "test", "test", nil
		}
		getPublicIP = func(client EC2ClientGetter, instanceID, region string) (string, error) {
			return "test", nil
		}
		newClusterAdmin = func(addrs []string, conf *sarama.Config) (sarama.ClusterAdmin, error) {
			return mockClusterAdmin{}, assert.AnError
		}
		createTopic = func(topic string, admin sarama.ClusterAdmin) error {
			return nil
		}
		sendMessage = func(producer sarama.SyncProducer, topic, from, to, message string) (int32, int64, error) {
			return 0, 0, nil
		}
		newSyncProducer = func(addrs []string, config *sarama.Config) (sarama.SyncProducer, error) {
			return &MockSaramaSyncProducer{}, nil
		}
		err := SendToKafka(&MockEC2Client{}, &MockSecretsManagerClient{}, "test", "test", "hello world", "us-east-1")
		assert.NotNil(t, err)
	})

	t.Run("test for sending message to kafka with error creating topic", func(t *testing.T) {
		resetMockTxnState()
		fetchSecrets = func(client SecretsManagerClient) (string, string, string, string, string, error) {
			return "test", "test", "9092", "test", "test", nil
		}
		getPublicIP = func(client EC2ClientGetter, instanceID, region string) (string, error) {
			return "test", nil
		}
		newClusterAdmin = func(addrs []string, conf *sarama.Config) (sarama.ClusterAdmin, error) {
			return mockClusterAdmin{}, nil
		}
		createTopic = func(topic string, admin sarama.ClusterAdmin) error {
			return assert.AnError
		}
		sendMessage = func(producer sarama.SyncProducer, topic, from, to, message string) (int32, int64, error) {
			return int32(0), int64(0), nil
		}
		newSyncProducer = func(addrs []string, config *sarama.Config) (sarama.SyncProducer, error) {
			return &MockSaramaSyncProducer{}, nil
		}
		err := SendToKafka(&MockEC2Client{}, &MockSecretsManagerClient{}, "test", "test", "hello world", "us-east-1")
		assert.NotNil(t, err)
	})

	t.Run("test for sending message to kafka with error in sending msg", func(t *testing.T) {
		resetMockTxnState()
		fetchSecrets = func(client SecretsManagerClient) (string, string, string, string, string, error) {
			return "test", "test", "9092", "test", "test", nil
		}
		getPublicIP = func(client EC2ClientGetter, instanceID, region string) (string, error) {
			return "test", nil
		}
		newClusterAdmin = func(addrs []string, conf *sarama.Config) (sarama.ClusterAdmin, error) {
			return mockClusterAdmin{}, nil
		}
		createTopic = func(topic string, admin sarama.ClusterAdmin) error {
			return nil
		}
		sendMessage = func(producer sarama.SyncProducer, topic, from, to, message string) (int32, int64, error) {
			return 0, 0, assert.AnError
		}
		newSyncProducer = func(addrs []string, config *sarama.Config) (sarama.SyncProducer, error) {
			return &MockSaramaSyncProducer{}, nil
		}
		err := SendToKafka(&MockEC2Client{}, &MockSecretsManagerClient{}, "test", "test", "hello world", "us-east-1")
		assert.NotNil(t, err)
	})

	t.Run("test for failed to create Kafka producer", func(t *testing.T) {
		resetMockTxnState()
		fetchSecrets = func(client SecretsManagerClient) (string, string, string, string, string, error) {
			return "test", "test", "9092", "test", "test", nil
		}
		getPublicIP = func(client EC2ClientGetter, instanceID, region string) (string, error) {
			return "test", nil
		}
		newClusterAdmin = func(addrs []string, conf *sarama.Config) (sarama.ClusterAdmin, error) {
			return mockClusterAdmin{}, nil
		}
		createTopic = func(topic string, admin sarama.ClusterAdmin) error {
			return nil
		}
		sendMessage = func(producer sarama.SyncProducer, topic, from, to, message string) (int32, int64, error) {
			return 0, 0, nil
		}
		newSyncProducer = func(addrs []string, config *sarama.Config) (sarama.SyncProducer, error) {
			return &MockSaramaSyncProducer{}, assert.AnError
		}
		err := SendToKafka(&MockEC2Client{}, &MockSecretsManagerClient{}, "test", "test", "hello world", "us-east-1")
		assert.NotNil(t, err)
	})

	t.Run("test for successful sending of message to kafka", func(t *testing.T) {
		resetMockTxnState()
		fetchSecrets = func(client SecretsManagerClient) (string, string, string, string, string, error) {
			return "test", "test", "9092", "test", "test", nil
		}
		getPublicIP = func(client EC2ClientGetter, instanceID, region string) (string, error) {
			return "test", nil
		}
		newClusterAdmin = func(addrs []string, conf *sarama.Config) (sarama.ClusterAdmin, error) {
			return mockClusterAdmin{}, nil
		}
		createTopic = func(topic string, admin sarama.ClusterAdmin) error {
			return nil
		}
		sendMessage = func(producer sarama.SyncProducer, topic, from, to, message string) (int32, int64, error) {
			return 0, 0, nil
		}
		newSyncProducer = func(addrs []string, config *sarama.Config) (sarama.SyncProducer, error) {
			return &MockSaramaSyncProducer{}, nil
		}
		err := SendToKafka(&MockEC2Client{}, &MockSecretsManagerClient{}, "test", "test", "hello world", "us-east-1")
		assert.Nil(t, err)
	})
}
