package logic

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/sync/model"
	"testing"

	"chainmaker.org/chainmaker/contract-utils/standard"
	"github.com/stretchr/testify/assert"
)

func TestNewContractHandler(t *testing.T) {
	minHeight := int64(100)
	contractMap := make(map[string]*db.Contract)
	eventResults := &model.TopicEventResult{}
	delayedUpdateCache := &model.GetRealtimeCacheData{
		TxList:         make(map[string]*db.Transaction),
		ContractEvents: []*db.ContractEvent{},
	}

	handler := NewContractHandler(db.UTchainID, minHeight, contractMap, eventResults, delayedUpdateCache)

	assert.Equal(t, db.UTchainID, handler.ChainId)
	assert.Equal(t, minHeight, handler.MinHeight)
	assert.Equal(t, delayedUpdateCache.TxList, handler.TxList)
	assert.Equal(t, contractMap, handler.ContractMap)
	assert.Equal(t, delayedUpdateCache.ContractEvents, handler.ContractEvents)
	assert.Equal(t, eventResults, handler.EventResults)
}

func TestDealWithContractData(t *testing.T) {
	minHeight := int64(100)
	contractMap := make(map[string]*db.Contract)
	eventResults := &model.TopicEventResult{
		IdentityContract: []*db.IdentityContract{
			{ContractAddr: "addr1"},
		},
	}
	delayedUpdateCache := &model.GetRealtimeCacheData{
		TxList:         make(map[string]*db.Transaction),
		ContractEvents: []*db.ContractEvent{},
	}
	handler := NewContractHandler(db.UTchainID, minHeight, contractMap, eventResults, delayedUpdateCache)

	idaContractMap := make(map[string]*db.IDAContract)
	result, err := handler.DealWithContractData(idaContractMap)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, eventResults.IdentityContract, result.IdentityContract)
}

func TestDealIDAContractUpdateData(t *testing.T) {
	minHeight := int64(100)
	contractMap := make(map[string]*db.Contract)
	eventResults := &model.TopicEventResult{
		IDAEventData: &model.IDAEventData{
			IDACreatedMap: map[string][]*db.IDACreatedInfo{
				"addr1": {
					{IDAInfo: &standard.IDAInfo{
						Basic: standard.Basic{
							ID:   "addr1",
							Name: "test",
						},
					}},
				},
			},
			IDADeletedCodeMap: map[string][]string{
				"addr2": {"id2"},
			},
		},
	}
	delayedUpdateCache := &model.GetRealtimeCacheData{
		TxList:         make(map[string]*db.Transaction),
		ContractEvents: []*db.ContractEvent{},
	}
	handler := NewContractHandler(db.UTchainID, minHeight, contractMap, eventResults, delayedUpdateCache)

	idaContractMap := map[string]*db.IDAContract{
		"addr1": {ContractAddr: "addr1", TotalNormalAssets: 1, TotalAssets: 1, BlockHeight: 50},
		"addr2": {ContractAddr: "addr2", TotalNormalAssets: 2, TotalAssets: 2, BlockHeight: 50},
	}
	result := handler.DealIDAContractUpdateData(idaContractMap)
	assert.NotNil(t, result)
}

func TestUpdateContractTxAndEventNum(t *testing.T) {
	minHeight := int64(100)
	contractMap := map[string]*db.Contract{
		"contract1": {Addr: "contract1", BlockHeight: 50},
		"contract2": {Addr: "contract2", BlockHeight: 50},
	}
	eventResults := &model.TopicEventResult{}
	delayedUpdateCache := &model.GetRealtimeCacheData{
		TxList: map[string]*db.Transaction{
			"tx1": {ContractNameBak: "contract1"},
			"tx2": {ContractNameBak: "contract2"},
		},
		ContractEvents: []*db.ContractEvent{
			{ContractNameBak: "contract1"},
			{ContractNameBak: "contract2"},
		},
	}
	handler := NewContractHandler(db.UTchainID, minHeight, contractMap, eventResults, delayedUpdateCache)
	result := handler.UpdateContractTxAndEventNum()
	assert.NotNil(t, result)
}
