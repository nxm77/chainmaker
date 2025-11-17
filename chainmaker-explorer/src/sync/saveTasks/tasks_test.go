package saveTasks

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/sync/model"
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRetrySleepTime(t *testing.T) {
	tests := []struct {
		retryCount int
		expected   int
	}{
		{0, 0},
		{1, 0},
		{2, 0},
		{3, 0},
		{4, 1},
		{5, 2},
		{6, 3},
	}

	for _, test := range tests {
		result := RetrySleepTime(test.retryCount)
		if result != test.expected {
			t.Errorf("RetrySleepTime(%d) = %d; expected %d", test.retryCount, result, test.expected)
		}
	}
}

func TestWithRetry(t *testing.T) {
	errCh := make(chan error)
	task := func() error {
		return errors.New("test error")
	}
	logFuncName := "testFunc"

	go WithRetry(task, logFuncName, errCh)
}

func TestExecuteTaskWithRetry(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	task := DealTask{
		Name:     "testTask",
		Function: func(...interface{}) error { return errors.New("test error") },
		Args:     []interface{}{},
	}
	retryCountMap := &sync.Map{}
	errCh := make(chan error)

	go ExecuteTaskWithRetry(ctx, wg, task, retryCountMap, errCh)
}

func TestTaskUpdateTxBlackToDB(t *testing.T) {
	updateTxBlack := &db.UpdateTxBlack{}
	err := TaskUpdateTxBlackToDB(db.UTchainID, updateTxBlack)
	assert.Nil(t, err)

	// Test with invalid data
	err = TaskUpdateTxBlackToDB(db.UTchainID, nil)
	assert.Error(t, err)

	err = TaskUpdateTxBlackToDB("", nil)
	assert.Error(t, err)
}

func TestTaskUpdateContractResult(t *testing.T) {
	contractResult := &db.GetContractResult{}
	err := TaskUpdateContractResult(db.UTchainID, contractResult)
	assert.Nil(t, err)

	// Test with invalid data
	err = TaskUpdateContractResult(db.UTchainID, nil)
	assert.Error(t, err)

	err = TaskUpdateContractResult("", nil)
	assert.Error(t, err)
}

func TestTaskInsertFungibleTransferToDB(t *testing.T) {
	transferList := []*db.FungibleTransfer{}
	err := TaskInsertFungibleTransferToDB(db.UTchainID, transferList)
	assert.Nil(t, err)

	// Test with invalid data
	err = TaskInsertFungibleTransferToDB(db.UTchainID, nil)
	assert.Error(t, err)

	err = TaskInsertFungibleTransferToDB("", nil)
	assert.Error(t, err)
}

func TestTaskInsertNonFungibleTransferToDB(t *testing.T) {
	transferList := []*db.NonFungibleTransfer{}
	err := TaskInsertNonFungibleTransferToDB(db.UTchainID, transferList)
	assert.Nil(t, err)

	// Test with invalid data
	err = TaskInsertNonFungibleTransferToDB(db.UTchainID, nil)
	assert.Error(t, err)

	err = TaskInsertNonFungibleTransferToDB("", nil)
	assert.Error(t, err)
}

func TestTaskSaveAccountListToDB(t *testing.T) {
	updateAccountResult := &db.UpdateAccountResult{}
	err := TaskSaveAccountListToDB(db.UTchainID, updateAccountResult)
	assert.Nil(t, err)

	// Test with invalid data
	err = TaskSaveAccountListToDB(db.UTchainID, nil)
	assert.Error(t, err)

	err = TaskSaveAccountListToDB("", nil)
	assert.Error(t, err)
}

func TestTaskSaveTokenResultToDB(t *testing.T) {
	tokenResult := &db.TokenResult{}
	err := TaskSaveTokenResultToDB(db.UTchainID, tokenResult)
	assert.Nil(t, err)

	// Test with invalid data
	err = TaskSaveTokenResultToDB(db.UTchainID, nil)
	assert.Error(t, err)

	err = TaskSaveTokenResultToDB("", nil)
	assert.Error(t, err)
}

func TestTaskSaveGasToDB(t *testing.T) {
	insertGasList := []*db.Gas{}
	updateGasList := []*db.Gas{}
	err := TaskSaveGasToDB(db.UTchainID, insertGasList, updateGasList)
	assert.Nil(t, err)

	// Test with invalid data
	err = TaskSaveGasToDB(db.UTchainID, nil, nil)
	assert.Error(t, err)

	err = TaskSaveGasToDB("", nil, nil)
	assert.Error(t, err)
}

func TestTaskSaveFungibleContractResult(t *testing.T) {
	contractResult := &db.GetContractResult{}
	err := TaskSaveFungibleContractResult(db.UTchainID, contractResult)
	assert.Nil(t, err)

	// Test with invalid data
	err = TaskSaveFungibleContractResult(db.UTchainID, nil)
	assert.Error(t, err)

	err = TaskSaveFungibleContractResult("", nil)
	assert.Error(t, err)
}

func TestTaskSavePositionToDB(t *testing.T) {
	blockPosition := &db.BlockPosition{}
	err := TaskSavePositionToDB(db.UTchainID, blockPosition)
	assert.Nil(t, err)

	// Test with invalid data
	err = TaskSavePositionToDB(db.UTchainID, nil)
	assert.Error(t, err)

	err = TaskSavePositionToDB("", nil)
	assert.Error(t, err)
}

func TestTaskCrossSubChainCrossToDB(t *testing.T) {
	insertList := []*db.CrossSubChainCrossChain{}
	updateList := []*db.CrossSubChainCrossChain{}
	err := TaskCrossSubChainCrossToDB(db.UTchainID, insertList, updateList)
	assert.Nil(t, err)

	// Test with invalid data
	err = TaskCrossSubChainCrossToDB(db.UTchainID, nil, nil)
	assert.Error(t, err)

	err = TaskCrossSubChainCrossToDB("", nil, nil)
	assert.Error(t, err)
}

func TestTaskSaveIDAAssetDataToDB(t *testing.T) {
	idaAssetsDataDB := &db.IDAAssetsDataDB{}
	err := TaskSaveIDAAssetDataToDB(db.UTchainID, idaAssetsDataDB)
	assert.Nil(t, err)

	// Test with invalid data
	err = TaskSaveIDAAssetDataToDB(db.UTchainID, nil)
	assert.Error(t, err)

	err = TaskSaveIDAAssetDataToDB("", nil)
	assert.Error(t, err)
}

func TestTaskUpdateIDAAssetDataToDB(t *testing.T) {
	idaAssetsDataDB := &db.IDAAssetsUpdateDB{}
	err := TaskUpdateIDAAssetDataToDB(db.UTchainID, idaAssetsDataDB)
	assert.Nil(t, err)

	// Test with invalid data
	err = TaskUpdateIDAAssetDataToDB(db.UTchainID, nil)
	assert.Error(t, err)

	err = TaskUpdateIDAAssetDataToDB("", nil)
	assert.Error(t, err)
}

func TestTaskUpdateDIDDataToDB(t *testing.T) {
	didSaveDate := &db.DIDSaveData{
		SaveVCTemplates: []*db.VCTemplate{
			{
				TemplateID:   "test",
				Template:     "test",
				ContractAddr: "test",
				ContractName: "test",
			},
		},
		UpdateDIDStatus: map[string]int{
			"test": 1,
		},
	}

	err := TaskUpdateDIDDataToDB(db.UTchainID, didSaveDate)
	assert.Nil(t, err)

	err = TaskUpdateDIDDataToDB(db.UTchainID, nil)
	assert.Error(t, err)

	err = TaskUpdateDIDDataToDB("", nil)
	assert.Error(t, err)

}

func TestTaskSaveChainStatisticsToDB(t *testing.T) {
	starts := &db.Statistics{}
	err := TaskSaveChainStatisticsToDB(db.UTchainID, starts)
	assert.Nil(t, err)

	err = TaskSaveChainStatisticsToDB("", starts)
	assert.Nil(t, err)

	err = TaskSaveChainStatisticsToDB("", nil)
	assert.Error(t, err)
}

func TestTaskSaveABITopicTableEvents(t *testing.T) {
	topicTableEvents := map[string][]map[string]interface{}{}
	err := TaskSaveABITopicTableEvents(db.UTchainID, topicTableEvents)
	assert.Nil(t, err)
	err = TaskSaveABITopicTableEvents(db.UTchainID, nil)
	assert.Error(t, err)
	err = TaskSaveABITopicTableEvents("", nil)
	assert.Error(t, err)
}

func TestTaskDelayCrossChain(t *testing.T) {
	delayCrossChain := &model.DelayCrossChain{}
	err := TaskDelayCrossChain(db.UTchainID, delayCrossChain)
	assert.Nil(t, err)
	err = TaskDelayCrossChain(db.UTchainID, nil)
	assert.Error(t, err)
	err = TaskDelayCrossChain("", nil)
	assert.Error(t, err)
}
