package dbhandle

import (
	"chainmaker_web/src/db"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestBatchInsertContractCrossCallTxs(t *testing.T) {
	newUUID := uuid.New().String()
	insertList := []*db.ContractCrossCallTransaction{
		{
			Id:               newUUID,
			TxId:             "tx123",
			BlockHeight:      123,
			InvokingContract: "contractC",
			InvokingMethod:   "methodC",
			TargetContract:   "contractB",
			UserAddr:         UserAddr1,
			IsCross:          true,
			Timestamp:        1234567890,
		},
	}

	err := BatchInsertContractCrossCallTxs(db.UTchainID, insertList)
	assert.NoError(t, err)

	err = BatchInsertContractCrossCallTxs(db.UTchainID, nil)
	assert.NoError(t, err)
}

func TestBatchInsertContractCrossCalls(t *testing.T) {
	newUUID := uuid.New().String()
	insertList := []*db.ContractCrossCall{
		{
			Id:               newUUID,
			InvokingContract: "contractA",
			InvokingMethod:   "methodA",
			TargetContract:   "contractB",
		},
	}

	err := BatchInsertContractCrossCalls(db.UTchainID, insertList)
	assert.NoError(t, err)
}

func TestGetContractCrossCallsByName(t *testing.T) {
	TestBatchInsertContractCrossCalls(t)

	results, err := GetContractCrossCallsByName(db.UTchainID, "contractA")
	assert.NoError(t, err)
	assert.NotNil(t, results)

	// empty chainId or contractName
	results, err = GetContractCrossCallsByName("", "contractA")
	assert.NoError(t, err)
	assert.Empty(t, results)
}
