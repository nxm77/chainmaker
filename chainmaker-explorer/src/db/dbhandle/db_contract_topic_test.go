package dbhandle

import (
	"chainmaker_web/src/db"
	"testing"

	"chainmaker.org/chainmaker/common/v2/random/uuid"
	"github.com/stretchr/testify/assert"
)

var TopicName = "Transfer"

// 测试插入合约 ABI 主题
func TestInsertContractABITopic(t *testing.T) {
	// 准备测试数据
	uuid := uuid.GetUUID()
	topic := &db.ContractABITopic{
		Id:              uuid,
		ContractName:    ContractName1,
		ContractAddr:    contractAdder1,
		ContractVersion: ContractVersionUT,
		TopicTableName:  "topic_table",
		Topic:           TopicName,
	}

	err := InsertContractABITopic(db.UTchainID, topic)
	assert.NoError(t, err)
}

// 测试获取单个合约 ABI 主题
func TestGetContractABITopic(t *testing.T) {
	TestInsertContractABITopic(t)

	// 执行测试
	result, err := GetContractABITopic(db.UTchainID, contractAdder1, ContractVersionUT, TopicName)
	assert.NoError(t, err)
	assert.NotNil(t, result)
}

// 测试获取不存在的合约 ABI 主题
func TestGetContractABITopic_NotFound(t *testing.T) {
	// 执行测试
	result, err := GetContractABITopic(db.UTchainID, contractAdder1, ContractVersionUT, "123")
	assert.NoError(t, err)
	assert.Nil(t, result)
}

// 测试获取合约的多个 ABI 主题
func TestGetContractABITopics(t *testing.T) {
	TestInsertContractABITopic(t)

	// 执行测试
	_, err := GetContractABITopics(db.UTchainID, contractAdder1, TopicName)
	assert.NoError(t, err)
}

// 测试根据合约地址列表获取 ABI 主题
func TestGetContractABITopicByAddrs(t *testing.T) {
	// 执行测试
	_, err := GetContractABITopicByAddrs(db.UTchainID, []string{contractAdder1, contractAdder2})
	assert.NoError(t, err)
}

// 测试更新合约 ABI 主题
func TestUpdateContractABITopic(t *testing.T) {
	TestInsertContractABITopic(t)

	topic := &db.ContractABITopic{
		ContractAddr:    contractAdder1,
		ContractVersion: "v1.0",
		Topic:           TopicName,
		TopicTableName:  "new_topic_table",
	}

	err := UpdateContractABITopic(db.UTchainID, topic)
	assert.NoError(t, err)
}
