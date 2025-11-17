package dbhandle

import (
	"chainmaker_web/src/db"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestInsertCrossTxTransfers tests the InsertCrossTxTransfers function
func TestInsertCrossTxTransfers(t *testing.T) {
	// Test case: empty crossTxTransfers
	err := InsertCrossTxTransfers(db.UTchainID, nil)
	assert.NoError(t, err)

	// Test case: valid crossTxTransfers
	crossTxTransfers := []*db.CrossTransactionTransfer{{
		CrossId:      "cross1",
		FromChainId:  db.UTchainID,
		ToChainId:    "chain2",
		ContractName: "contract1",
		Status:       1,
		StartTime:    123456789,
		EndTime:      123456790,
	}}
	err = InsertCrossTxTransfers(db.UTchainID, crossTxTransfers)
	assert.NoError(t, err)
}

// TestUpdateCrossTxTransfers tests the UpdateCrossTxTransfers function
func TestUpdateCrossTxTransfers(t *testing.T) {
	// Test case: empty chainId or crossId
	err := UpdateCrossTxTransfers("", "cross1", 2, 123456790)
	assert.Error(t, err)

	// Test case: valid update
	err = UpdateCrossTxTransfers(db.UTchainID, "cross1", 2, 123456790)
	assert.NoError(t, err)
}

// TestCheckCrossIdsExistenceTransfer tests the CheckCrossIdsExistenceTransfer function
func TestCheckCrossIdsExistenceTransfer(t *testing.T) {
	// Test case: empty chainId or crossIds
	_, err := CheckCrossIdsExistenceTransfer("", nil)
	assert.NoError(t, err)

	// Test case: valid crossIds
	crossIds := []string{"cross1", "cross2"}
	_, err = CheckCrossIdsExistenceTransfer(db.UTchainID, crossIds)
	assert.NoError(t, err)
}

// TestGetCrossCycleTransferById tests the GetCrossCycleTransferById function
func TestGetCrossCycleTransferById(t *testing.T) {
	// Test case: empty chainId or crossId
	_, err := GetCrossCycleTransferById("", "")
	assert.Error(t, err)

	// Test case: valid chainId and crossId
	_, err = GetCrossCycleTransferById(db.UTchainID, "cross1")
	assert.NoError(t, err)
}

// TestGetCrossCycleTransferByCrossIds tests the GetCrossCycleTransferByCrossIds function
func TestGetCrossCycleTransferByCrossIds(t *testing.T) {
	// Test case: empty chainId or crossIds
	_, err := GetCrossCycleTransferByCrossIds("", nil)
	assert.Error(t, err)

	// Test case: valid crossIds
	crossIds := []string{"cross1", "cross2"}
	_, err = GetCrossCycleTransferByCrossIds(db.UTchainID, crossIds)
	assert.NoError(t, err)
}

// TestGetCrossCycleTransferByHeight tests the GetCrossCycleTransferByHeight function
func TestGetCrossCycleTransferByHeight(t *testing.T) {
	// Test case: empty chainId or blockHeights
	_, err := GetCrossCycleTransferByHeight("", nil)
	assert.Error(t, err)

	// Test case: valid blockHeights
	blockHeights := []int64{100, 200}
	_, err = GetCrossCycleTransferByHeight(db.UTchainID, blockHeights)
	assert.NoError(t, err)
}

// TestGetCrossTransferDurationByTime tests the GetCrossTransferDurationByTime function
func TestGetCrossTransferDurationByTime(t *testing.T) {
	// Test case: empty chainId
	_, err := GetCrossTransferDurationByTime("", 123456789, 123456790)
	assert.Error(t, err)

	// Test case: valid time range
	_, err = GetCrossTransferDurationByTime(db.UTchainID, 123456789, 123456790)
	assert.NoError(t, err)
}

// TestGetCrossSubChainTransferList tests the GetCrossSubChainTransferList function
func TestGetCrossSubChainTransferList(t *testing.T) {
	// Test case: empty chainId
	_, _, err := GetCrossSubChainTransferList(0, 10, 123456789, 123456790, "", "", "", "", "")
	assert.Error(t, err)

	// Test case: valid parameters
	_, _, err = GetCrossSubChainTransferList(0, 10, 123456789, 123456790, db.UTchainID, "cross1", "subChain1", "fromChain1", "toChain1")
	assert.NoError(t, err)
}

// TestGetCrossTxTransferLatestList tests the GetCrossTxTransferLatestList function
func TestGetCrossTxTransferLatestList(t *testing.T) {
	// Test case: empty chainId
	result, err := GetCrossTxTransferLatestList("")
	assert.Error(t, err)
	assert.Empty(t, result)

	// Test case: valid chainId
	result, err = GetCrossTxTransferLatestList(db.UTchainID)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

// TestGetCrossContractByTransfer tests the GetCrossContractByTransfer function
func TestGetCrossContractByTransfer(t *testing.T) {
	// Test case: empty chainId
	result, err := GetCrossContractByTransfer("")
	assert.Error(t, err)
	assert.Empty(t, result)

	// Test case: valid chainId
	result, err = GetCrossContractByTransfer(db.UTchainID)
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}
