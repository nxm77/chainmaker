package saveTasks

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestCalculateDataSize(t *testing.T) {
	data := map[string]string{"key": "value"}
	expectedSize := utf8.RuneCount([]byte(`{"key":"value"}`))
	size := calculateDataSize(data)
	assert.Equal(t, expectedSize, size)
}

func TestBatchTransactions(t *testing.T) {
	config.MaxDBByteSize = 50
	transactions := map[string]*db.Transaction{
		"tx1": {TxId: "tx1"},
		"tx2": {TxId: "tx2"},
	}

	batches := batchTransactions(transactions)
	if len(batches) == 0 {
		t.Errorf("Expected at least one batch, but got none")
	}
}

func TestBatchUsers(t *testing.T) {
	users := map[string]*db.User{
		"user1": {UserId: "user1"},
		"user2": {UserId: "user2"},
	}

	batches := batchUsers(users)
	if len(batches) == 0 {
		t.Errorf("Expected at least one batch, but got none")
	}
}

func TestBatchContractEvents(t *testing.T) {
	config.MaxDBByteSize = 50
	contractEvents := []*db.ContractEvent{
		{ID: "event1"},
		{ID: "event2"},
	}

	batches := batchContractEvents(contractEvents)
	if len(batches) == 0 {
		t.Errorf("Expected at least one batch, but got none")
	}
}

func TestBatchGasRecords(t *testing.T) {
	config.MaxDBByteSize = 50
	gasRecords := []*db.GasRecord{
		{ID: "record1"},
		{ID: "record2"},
	}

	batches := batchGasRecords(gasRecords)
	if len(batches) == 0 {
		t.Errorf("Expected at least one batch, but got none")
	}
}

func TestBatchFungibleTransfers(t *testing.T) {
	config.MaxDBByteSize = 50
	transfers := []*db.FungibleTransfer{
		{ID: "transfer1"},
		{ID: "transfer2"},
	}

	batches := batchFungibleTransfers(transfers)
	if len(batches) == 0 {
		t.Errorf("Expected at least one batch, but got none")
	}
}

func TestBatchNonFungibleTransfers(t *testing.T) {
	config.MaxDBByteSize = 50
	transfers := []*db.NonFungibleTransfer{
		{ID: "transfer1"},
		{ID: "transfer2"},
	}

	batches := batchNonFungibleTransfers(transfers)
	if len(batches) == 0 {
		t.Errorf("Expected at least one batch, but got none")
	}
}

func TestBatchNonFungibleToken(t *testing.T) {
	config.MaxDBByteSize = 50
	tokens := []*db.NonFungibleToken{
		{TokenId: "token1"},
		{TokenId: "token2"},
	}

	batches := batchNonFungibleToken(tokens)
	if len(batches) == 0 {
		t.Errorf("Expected at least one batch, but got none")
	}
}
