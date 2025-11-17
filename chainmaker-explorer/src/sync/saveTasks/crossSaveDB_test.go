package saveTasks

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/sync/model"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var TxID1 = "TxID111"

func TestMain(m *testing.M) {
	_ = config.InitConfig("", "")
	// 初始化数据库配置
	db.InitRedisContainer()
	db.InitMySQLContainer()

	// 运行其他测试
	os.Exit(m.Run())
}

func TestSaveBusinessTransaction(t *testing.T) {
	businessTxMap := map[string]*db.CrossBusinessTransaction{
		TxID1: {TxId: TxID1},
	}

	err := SaveBusinessTransaction(db.UTchainID, businessTxMap)
	assert.Nil(t, err)
}

func TestSaveCrossSubChainCrossToDB(t *testing.T) {
	inserts := []*db.CrossSubChainCrossChain{
		{ID: "cross1"},
	}
	updates := []*db.CrossSubChainCrossChain{
		{ID: "cross2"},
	}

	err := SaveCrossSubChainCrossToDB(db.UTchainID, inserts, updates)
	assert.Nil(t, err)
}

func TestTaskSaveRelayCrossChain_Execute(t *testing.T) {
	task := TaskSaveRelayCrossChain{
		ChainId: db.UTchainID,
		CrossChainResult: &model.CrossChainResult{
			CrossMainTransaction: []*db.CrossMainTransaction{{TxId: "TxID1"}},
			BusinessTxMap:        map[string]*db.CrossBusinessTransaction{"TxID1": {TxId: "TxID1"}},
		},
	}

	err := task.Execute()
	assert.Nil(t, err)
}

func TestTaskSaveRelayCrossChain_SaveRelayCrossChainToDB_NilResult(t *testing.T) {
	task := TaskSaveRelayCrossChain{
		ChainId:          db.UTchainID,
		CrossChainResult: nil,
	}

	err := task.SaveRelayCrossChainToDB()
	assert.Nil(t, err)
}

func TestUpdateCrossTransfer(t *testing.T) {
	insertTransferMap := []*db.CrossTransactionTransfer{
		{CrossId: "transfer1"},
	}
	updateTransferMap := []*db.CrossTransactionTransfer{
		{CrossId: "transfer2"},
	}

	crossChainResult := &model.CrossChainResult{
		InsertCrossTransfer: insertTransferMap,
		UpdateCrossTransfer: updateTransferMap,
	}

	err := UpdateCrossTransfer(db.UTchainID, crossChainResult)
	assert.Nil(t, err)
}

func TestGetSubChainSaveList(t *testing.T) {
	saveSubChainList := map[string]*db.CrossSubChainData{
		"sub1": {SubChainId: "sub1"},
	}

	insertList, updateList, err := GetSubChainSaveList(db.UTchainID, saveSubChainList)
	assert.Nil(t, err)
	assert.NotNil(t, insertList)
	assert.NotNil(t, updateList)
}

func TestSaveDelayCrossChainToDB(t *testing.T) {
	delayCrossChain := &model.DelayCrossChain{
		InsertSubChainCross: []*db.CrossSubChainCrossChain{{ID: "cross11"}},
		UpdateSubChainCross: []*db.CrossSubChainCrossChain{{ID: "cross22"}},
		InsertSubChainData:  []*db.CrossSubChainData{{SubChainId: "sub1"}},
		UpdateSubChainData:  []*db.CrossSubChainData{{SubChainId: "sub2"}},
		CrossChainContracts: []*db.CrossChainContract{{ContractName: "contract1"}},
	}

	err := SaveDelayCrossChainToDB(db.UTchainID, delayCrossChain)
	assert.Nil(t, err)
}

func TestSaveDelayCrossChainToDB_NilInput(t *testing.T) {
	err := SaveDelayCrossChainToDB(db.UTchainID, nil)
	assert.Nil(t, err)
}
