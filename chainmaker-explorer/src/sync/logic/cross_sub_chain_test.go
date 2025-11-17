package logic

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/entity_cross"
	"testing"

	tcipCommon "chainmaker.org/chainmaker/tcip-go/v2/common"

	"github.com/stretchr/testify/assert"
)

func TestGetCrossModelByExtra(t *testing.T) {
	// Test case 1: SyncData is "true"
	crossChainInfo := &tcipCommon.CrossChainInfo{
		CrossChainMsg: []*tcipCommon.CrossChainMsg{
			{ExtraData: "{\"SyncData\":\"true\"}"},
		},
	}
	result := GetCrossModelByExtra(crossChainInfo)
	assert.Equal(t, entity_cross.CrossModelSync, result)

	// Test case 2: SyncData is not "true"
	crossChainInfo = &tcipCommon.CrossChainInfo{
		CrossChainMsg: []*tcipCommon.CrossChainMsg{
			{ExtraData: "{\"SyncData\":\"false\"}"},
		},
	}
	result = GetCrossModelByExtra(crossChainInfo)
	assert.Equal(t, entity_cross.CrossModelOther, result)

	// Test case 3: ExtraData is empty
	crossChainInfo = &tcipCommon.CrossChainInfo{
		CrossChainMsg: []*tcipCommon.CrossChainMsg{
			{ExtraData: ""},
		},
	}
	result = GetCrossModelByExtra(crossChainInfo)
	assert.Equal(t, entity_cross.CrossModelOther, result)
}

func TestGetCrossChainMsgExtra(t *testing.T) {
	// Test case 1: Valid ExtraData
	crossChainInfo := &tcipCommon.CrossChainInfo{
		CrossChainMsg: []*tcipCommon.CrossChainMsg{
			{ExtraData: "{\"SyncData\":\"true\",\"SyncDataBatchCount\":\"10\",\"SyncFromBlockHeight\":\"100\",\"SyncToBlockHeight\":\"200\"}"},
		},
	}
	result := GetCrossChainMsgExtra(crossChainInfo)
	assert.Equal(t, "true", result.SyncData)
	assert.Equal(t, "10", result.SyncDataBatchCount)
	assert.Equal(t, "100", result.SyncFromBlockHeight)
	assert.Equal(t, "200", result.SyncToBlockHeight)

	// Test case 2: Invalid ExtraData
	crossChainInfo = &tcipCommon.CrossChainInfo{
		CrossChainMsg: []*tcipCommon.CrossChainMsg{
			{ExtraData: "invalid_json"},
		},
	}
	result = GetCrossChainMsgExtra(crossChainInfo)
	assert.Equal(t, "", result.SyncData)
	assert.Equal(t, "", result.SyncDataBatchCount)
	assert.Equal(t, "", result.SyncFromBlockHeight)
	assert.Equal(t, "", result.SyncToBlockHeight)

	// Test case 3: CrossChainInfo is nil
	result = GetCrossChainMsgExtra(nil)
	assert.Equal(t, "", result.SyncData)
	assert.Equal(t, "", result.SyncDataBatchCount)
	assert.Equal(t, "", result.SyncFromBlockHeight)
	assert.Equal(t, "", result.SyncToBlockHeight)
}

func TestGetMainCrossTransaction(t *testing.T) {
	// Test case 1: Valid CrossChainInfo
	crossChainInfo := &tcipCommon.CrossChainInfo{
		CrossChainId: "cross123",
		State:        tcipCommon.CrossChainStateValue_CONFIRM_END,
		CrossChainMsg: []*tcipCommon.CrossChainMsg{
			{ExtraData: "{\"SyncData\":\"true\"}"},
		},
	}
	result := GetMainCrossTransaction(100, crossChainInfo, "tx123", 1620000000)
	assert.Equal(t, "tx123", result.TxId)
	assert.Equal(t, "cross123", result.CrossId)
	assert.Equal(t, entity_cross.CrossModelSync, result.CrossModel)
	assert.Equal(t, int32(tcipCommon.CrossChainStateValue_CONFIRM_END), result.Status)
	assert.Equal(t, int64(100), result.BlockHeight)
	assert.Equal(t, int64(1620000000), result.Timestamp)

	// Test case 2: CrossChainInfo is nil
	result = GetMainCrossTransaction(100, nil, "tx123", 1620000000)
	assert.Nil(t, result)
}

func TestGetBusinessTransaction(t *testing.T) {
	// Test case 1: CrossChainInfo is nil
	result := GetBusinessTransaction("chain123", nil)
	assert.Equal(t, 0, len(result))

	// Test case 2: CrossChainInfo state is not CONFIRM_END or CANCEL_END
	crossChainInfo := &tcipCommon.CrossChainInfo{
		//State: tcipCommon.CrossChainStateValue_INIT,
	}
	result = GetBusinessTransaction("chain123", crossChainInfo)
	assert.Equal(t, 0, len(result))
}

func TestGetCrossTxTransfer(t *testing.T) {
	// Test case 1: CrossChainInfo is nil
	result := GetCrossTxTransfer("chain123", "tx123", 100, 1620000000, nil)
	assert.Nil(t, result)

	// Test case 2: CrossModel is Other
	crossChainInfo := &tcipCommon.CrossChainInfo{
		CrossChainId: "cross123",
		From:         "gateway123",
		// FirstTxContent: &tcipCommon{
		// 	TxContent: &tcipCommon.TxContent{
		// 		ChainRid:  "chain456",
		// 		GatewayId: "gateway123",
		// 	},
		// },
		CrossChainMsg: []*tcipCommon.CrossChainMsg{
			{
				ChainRid:     "chain789",
				GatewayId:    "gateway456",
				ContractName: "contract123",
				Method:       "method123",
				Parameter:    "param123",
			},
		},
	}
	_ = GetCrossTxTransfer("chain123", "tx123", 100, 1620000000, crossChainInfo)

	// Test case 3: CrossModel is Sync
	crossChainInfo = &tcipCommon.CrossChainInfo{
		CrossChainId: "cross123",
		From:         "gateway123",
		CrossChainMsg: []*tcipCommon.CrossChainMsg{
			{
				ExtraData:    "{\"SyncData\":\"true\",\"SyncToBlockHeight\":\"200\"}",
				ChainRid:     "chain789",
				GatewayId:    "gateway456",
				ContractName: "contract123",
				Method:       "method123",
			},
		},
	}
	_ = GetCrossTxTransfer("chain123", "tx123", 100, 1620000000, crossChainInfo)
}

func TestBuildExecutionTransaction(t *testing.T) {
	// Test case 1: TxContent is nil
	result := BuildExecutionTransaction(nil)
	assert.NotNil(t, result)

	// Test case 2: Valid TxContent
	txContent := &tcipCommon.TxContent{
		TxId: "tx123",
	}
	_ = BuildExecutionTransaction(txContent)
}

func TestParseCrossCycleTxTransfer(t *testing.T) {
	// Test case 1: Empty transfers
	result := ParseCrossCycleTxTransfer(nil)
	assert.Equal(t, 0, len(result))

	// Test case 2: Valid transfers
	transfers := []*db.CrossTransactionTransfer{
		{
			FromChainId: "chain123",
			ToChainId:   "chain456",
		},
		{
			FromChainId: "chain123",
			ToChainId:   "chain789",
		},
		{
			FromChainId: "chain456",
			ToChainId:   "chain123",
		},
	}
	_ = ParseCrossCycleTxTransfer(transfers)
}

func TestDealSubChainCrossChainNum(t *testing.T) {
	// Test case 1: Empty subChainIdMap
	_, _, err := DealSubChainCrossChainNum("chain123", nil, nil, 100)
	assert.Nil(t, err)

	// Test case 2: Valid subChainIdMap with new entries
	subChainIdMap := map[string]map[string]int64{
		"chain123": {"chain456": 1},
	}
	_, _, err = DealSubChainCrossChainNum("chain123", subChainIdMap, nil, 100)
	assert.Nil(t, err)

	// Test case 3: Valid subChainIdMap with existing entries
	subChainCrossDB := []*db.CrossSubChainCrossChain{
		{
			SubChainId:   "chain123",
			CrossChainId: "chain456",
			TxNum:        1,
			BlockHeight:  50,
		},
	}
	_, _, err = DealSubChainCrossChainNum("chain123", subChainIdMap, subChainCrossDB, 100)
	assert.Nil(t, err)
}

func TestDealCrossSubChainData(t *testing.T) {
	// Test case 1: Empty insertCrossTransfers
	inserts, updates, contracts := DealCrossSubChainData(nil, nil)
	assert.Equal(t, 0, len(inserts))
	assert.Equal(t, 0, len(updates))
	assert.Equal(t, 0, len(contracts))

	// Test case 2: Valid insertCrossTransfers with new entries
	insertCrossTransfers := []*db.CrossTransactionTransfer{
		{
			FromChainId:     "chain123",
			FromGatewayId:   "gateway123",
			FromIsMainChain: true,
			FromBlockHeight: 100,
			StartTime:       1620000000,
		},
	}
	_, _, _ = DealCrossSubChainData(insertCrossTransfers, nil)

	// Test case 3: Valid insertCrossTransfers with existing entries
	subChainDataDB := map[string]*db.CrossSubChainData{
		"chain123": {
			SubChainId:  "chain123",
			GatewayId:   "gateway123",
			IsMainChain: true,
			BlockHeight: 50,
			TxNum:       1,
		},
	}
	_, _, _ = DealCrossSubChainData(insertCrossTransfers, subChainDataDB)

	// Test case 4: Valid insertCrossTransfers with contracts
	insertCrossTransfers = []*db.CrossTransactionTransfer{
		{
			ToChainId:    "chain456",
			ContractName: "contract123",
		},
	}
	_, _, _ = DealCrossSubChainData(insertCrossTransfers, nil)

}

func TestGetCrossChainAndContracts(t *testing.T) {
	// Test case 1: Empty crossTxTransferList
	crossChains, contracts := GetCrossChainAndContracts(nil)
	assert.Equal(t, 0, len(crossChains))
	assert.Equal(t, 0, len(contracts))

	// Test case 2: Valid crossTxTransferList
	crossTxTransferList := []*db.CrossTransactionTransfer{
		{
			ToChainId:     "chain123",
			ToGatewayId:   "gateway123",
			ToIsMainChain: true,
			ToBlockHeight: 100,
			StartTime:     1620000000,
			ContractName:  "contract123",
		},
	}
	crossChains, contracts = GetCrossChainAndContracts(crossTxTransferList)
	assert.Equal(t, 1, len(crossChains))
	assert.Equal(t, 1, len(contracts))
}
