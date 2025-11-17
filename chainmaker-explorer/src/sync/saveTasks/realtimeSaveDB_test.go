package saveTasks

import (
	"chainmaker_web/src/db"
	"testing"

	pbConfig "chainmaker.org/chainmaker/pb-go/v2/config"
	"github.com/test-go/testify/assert"
)

func TestTaskSaveTransactions_Execute(t *testing.T) {
	task := TaskSaveTransactions{
		ChainId: db.UTchainID,
		Transactions: map[string]*db.Transaction{
			"tx1": {TxId: "tx1"},
		},
		UpgradeContractTxs: []*db.UpgradeContractTransaction{
			{TxId: "upgradeTx1"},
		},
	}

	err := task.Execute()
	assert.Nil(t, err)
}

func TestSaveTransactionsToDB(t *testing.T) {
	transactions := map[string]*db.Transaction{
		"tx1": {TxId: "tx1"},
	}

	err := SaveTransactionsToDB(db.UTchainID, transactions)
	assert.Nil(t, err)
}

func TestSaveUpgradeContractTxToDB(t *testing.T) {
	upgradeContractTxs := []*db.UpgradeContractTransaction{
		{TxId: "upgradeTx1"},
	}

	err := SaveUpgradeContractTxToDB(db.UTchainID, upgradeContractTxs)
	assert.Nil(t, err)
}

func TestTaskInsertUser_Execute(t *testing.T) {
	task := TaskInsertUser{
		ChainId: db.UTchainID,
		UserList: map[string]*db.User{
			"user1": {UserId: "user1"},
		},
	}

	err := task.Execute()
	assert.Nil(t, err)
}

func TestTaskSaveContract_Execute(t *testing.T) {
	task := TaskSaveContract{
		ChainId: db.UTchainID,
		InsertList: []*db.Contract{
			{Addr: "contract1"},
		},
		UpdateList: []*db.Contract{
			{Addr: "contract2"},
		},
		ByteCodeList: []*db.ContractByteCode{
			{TxId: "TxId"},
		},
	}

	err := task.Execute()
	assert.Nil(t, err)
}

func TestTaskSaveStandardContract_Execute(t *testing.T) {
	task := TaskSaveStandardContract{
		ChainId: db.UTchainID,
		InsertFTContracts: []*db.FungibleContract{
			{ContractAddr: "ftContract1"},
		},
		InsertNFTContracts: []*db.NonFungibleContract{
			{ContractAddr: "nftContract1"},
		},
		InsertIDAContracts: []*db.IDAContract{
			{ContractAddr: "idaContract1"},
		},
	}

	err := task.Execute()
	assert.Nil(t, err)
}

func TestTaskEvidenceContract_Execute(t *testing.T) {
	task := TaskEvidenceContract{
		ChainId: db.UTchainID,
		EvidenceContracts: []*db.EvidenceContract{
			{Hash: "hash1"},
		},
	}

	err := task.Execute()
	assert.Nil(t, err)
}

func TestTaskInsertContractEvents_Execute(t *testing.T) {
	task := TaskInsertContractEvents{
		ChainId: db.UTchainID,
		ContractEvents: []*db.ContractEvent{
			{ID: "event1"},
		},
	}

	err := task.Execute()
	assert.Nil(t, err)
}

func TestTaskInsertGasRecord_Execute(t *testing.T) {
	task := TaskInsertGasRecord{
		ChainId: db.UTchainID,
		InsertGasRecords: []*db.GasRecord{
			{ID: "record1"},
		},
	}

	err := task.Execute()
	assert.Nil(t, err)
}

func TestTaskSaveChainConfig_Execute(t *testing.T) {
	task := TaskSaveChainConfig{
		ChainId: db.UTchainID,
		UpdateChainConfigs: []*pbConfig.ChainConfig{
			{ChainId: db.UTchainID},
		},
	}

	err := task.Execute()
	assert.Nil(t, err)
}

func TestTaskSaveContractCrossCallTxs_Execute(t *testing.T) {
	task := TaskSaveContractCrossCallTxs{
		ChainId: db.UTchainID,
		InsertCrossTxs: []*db.ContractCrossCallTransaction{
			{TxId: "crossCallTx1"},
		},
		InsertCalls: map[string]*db.ContractCrossCall{
			"call1": {Id: "call1", InvokingContract: "contract1", TargetContract: "contract2", InvokingMethod: "method1"},
		},
	}

	err := task.Execute()
	assert.Nil(t, err)
}
