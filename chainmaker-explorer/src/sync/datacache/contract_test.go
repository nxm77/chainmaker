package datacache

import (
	"chainmaker_web/src/db"
	"testing"
)

var ChainId2 = "chain2"

func TestSetLatestContractListCache(t *testing.T) {
	// Test case 1: Insert new contracts

	blockHeight1 := int64(100)
	insertContracts1 := []*db.Contract{
		{Addr: "contract1", Version: "1.0", ContractStatus: 0, UpgradeAddr: "", UpgradeTimestamp: 0, TxNum: 0, EventNum: 0},
		{Addr: "contract2", Version: "1.0", ContractStatus: 0, UpgradeAddr: "", UpgradeTimestamp: 0, TxNum: 0, EventNum: 0},
	}
	updateContracts1 := []*db.Contract{}
	SetLatestContractListCache(ChainId2, blockHeight1, insertContracts1, updateContracts1)

	// Test case 2: Update existing contracts
	blockHeight2 := int64(100)
	insertContracts2 := []*db.Contract{}
	updateContracts2 := []*db.Contract{
		{Addr: "contract1", Version: "2.0", ContractStatus: 0, UpgradeAddr: "", UpgradeTimestamp: 0, TxNum: 0, EventNum: 0},
		{Addr: "contract2", Version: "2.0", ContractStatus: 0, UpgradeAddr: "", UpgradeTimestamp: 0, TxNum: 0, EventNum: 0},
	}
	SetLatestContractListCache(ChainId2, blockHeight2, insertContracts2, updateContracts2)

	// Test case 3: No new or updated contracts
	blockHeight3 := int64(100)
	insertContracts3 := []*db.Contract{}
	updateContracts3 := []*db.Contract{}
	SetLatestContractListCache(ChainId2, blockHeight3, insertContracts3, updateContracts3)
}

func TestUpdateLatestContractCache(t *testing.T) {
	// Test case 1: Update existing contracts
	updateContracts1 := []*db.Contract{
		{Addr: "contract1", Version: "2.0", ContractStatus: 0, UpgradeAddr: "", UpgradeTimestamp: 0, TxNum: 0, EventNum: 0},
		{Addr: "contract2", Version: "2.0", ContractStatus: 0, UpgradeAddr: "", UpgradeTimestamp: 0, TxNum: 0, EventNum: 0},
	}
	UpdateLatestContractCache(ChainId2, updateContracts1)

	// Test case 2: No existing contracts to update
	updateContracts2 := []*db.Contract{}
	UpdateLatestContractCache(ChainId2, updateContracts2)
}
