// 在项目根目录执行以下命令修复依赖问题：
// go get github.com/go-sql-driver/mysql
// go mod tidy

package logic

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/entity"
	"chainmaker_web/src/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ParamsUTTest = &entity.UploadContractAbiParams{
	ChainId:         db.UTchainID,
	ContractAddr:    "123",
	ContractVersion: "v1",
	AbiJson:         nil, // 模拟文件上传
}

var UpgradeContractUTTTest = &db.UpgradeContractTransaction{
	ContractName: "TestContract",
}

func TestValidateABIFile(t *testing.T) {
	// 正常事件
	abi := []*utils.ContractABI{
		{Type: utils.ABIEventType, Name: "evt", Inputs: []utils.ABIParamIn{{Name: "a", Type: "string"}}},
	}
	assert.NoError(t, validateABIFile(abi))

	// 非法类型
	abi = []*utils.ContractABI{
		{Type: "unknown", Name: "evt"},
	}
	assert.Error(t, validateABIFile(abi))

	// 系统字段
	abi = []*utils.ContractABI{
		{Type: utils.ABIEventType, Name: "evt", Inputs: []utils.ABIParamIn{{Name: db.ABISystemFieldTxID, Type: "string"}}},
	}
	assert.Error(t, validateABIFile(abi))
}

func TestProcessEventTables(t *testing.T) {
	parsedABI := []*utils.ContractABI{
		{Type: utils.ABIEventType, Name: "evt", Inputs: []utils.ABIParamIn{{Name: "a"}}},
		{Type: "function", Name: "f"},
	}
	err := processEventTables(ParamsUTTest, UpgradeContractUTTTest, parsedABI)
	assert.NoError(t, err)
}

func TestHandleEventVersion1(t *testing.T) {
	event := &utils.ContractABI{
		Type: utils.ABIEventType,
		Name: "TestEvent",
		Inputs: []utils.ABIParamIn{
			{Name: "param1", Type: "string"},
			{Name: "param2", Type: "int"},
		},
		StateMutability: "view",
		Describe:        "This is a test event",
	}

	// 测试场景1: 新记录（记录不存在）
	t.Run("NewRecord", func(t *testing.T) {
		err := handleEventVersion(ParamsUTTest, UpgradeContractUTTTest, event, "table2")
		assert.NoError(t, err)
	})

	// 测试场景2: 表名相同（不需要更新）
	t.Run("SameTableName", func(t *testing.T) {
		insert := &db.ContractABITopic{
			ContractName:    UpgradeContractUTTTest.ContractName,
			ContractAddr:    ParamsUTTest.ContractAddr,
			ContractVersion: ParamsUTTest.ContractVersion,
			Topic:           event.Name,
			TopicTableName:  "table1",
		}
		err := dbhandle.InsertContractABITopic(db.UTchainID, insert)
		assert.NoError(t, err)

		// 使用相同的表名再次调用
		err = handleEventVersion(ParamsUTTest, UpgradeContractUTTTest, event, insert.TopicTableName)
		assert.NoError(t, err)
	})

	// 测试场景3: 表名不同（需要更新）
	t.Run("DifferentTableName", func(t *testing.T) {
		insert := &db.ContractABITopic{
			ContractName:    UpgradeContractUTTTest.ContractName,
			ContractAddr:    ParamsUTTest.ContractAddr,
			ContractVersion: ParamsUTTest.ContractVersion,
			Topic:           event.Name,
			TopicTableName:  "table3",
		}
		err := dbhandle.InsertContractABITopic(db.UTchainID, insert)
		assert.NoError(t, err)

		// 使用相同的表名再次调用
		err = handleEventVersion(ParamsUTTest, UpgradeContractUTTTest, event, "table4")
		assert.NoError(t, err)
	})
}

func TestCreateTopicRecord(t *testing.T) {
	event := &utils.ContractABI{
		Type: utils.ABIEventType,
		Name: "TestEvent1",
		Inputs: []utils.ABIParamIn{
			{Name: "param3", Type: "string"},
			{Name: "param4", Type: "int"},
		},
		StateMutability: "view1",
		Describe:        "event",
	}
	err := createTopicRecord(ParamsUTTest, UpgradeContractUTTTest, event, "table")
	assert.NoError(t, err)
}

func TestGetContractTopics(t *testing.T) {
	// 测试获取合约主题
	topics, err := GetContractTopics(db.UTchainID, "123", "v1")
	assert.NoError(t, err)
	assert.NotNil(t, topics)

	// 测试无事件的情况
	topics, err = GetContractTopics(db.UTchainID, "456", "v1")
	assert.NoError(t, err)
	assert.Empty(t, topics)
}

func TestGetContractTopicIndexs(t *testing.T) {
	// 测试获取合约主题索引
	indexes, err := GetContractTopicIndexs(db.UTchainID, "123", "v1", "TestEvent")
	assert.NoError(t, err)
	assert.NotNil(t, indexes)

	// 测试无索引的情况
	indexes, err = GetContractTopicIndexs(db.UTchainID, "456", "v1", "UnknownEvent")
	assert.NoError(t, err)
	assert.Empty(t, indexes)
}

func TestGetDecodeContractEvents(t *testing.T) {
	// 测试解码合约事件
	events, total, err := GetDecodeContractEvents(&entity.GetDecodeContractEventsParams{
		ChainId:         db.UTchainID,
		ContractAddr:    "123",
		ContractVersion: "v1",
		Topic:           "TestEvent",
	}, []string{"param1", "param2"}, []string{"param1"})
	assert.NoError(t, err)
	assert.NotNil(t, events)
	assert.GreaterOrEqual(t, total, int64(0))
}

func TestAsyncHandleContractABIEvent(t *testing.T) {
	// 测试异步处理合约事件
	err := AsyncHandleContractABIEvent(db.UTchainID, "123", "v1")
	assert.NoError(t, err)
}
