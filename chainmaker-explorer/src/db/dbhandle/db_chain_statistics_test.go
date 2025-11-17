package dbhandle

import (
	"chainmaker_web/src/db"
	"testing"
)

// 初始化测试获取链统计列表函数
func insertStatisticsTest(t *testing.T) {
	stats1 := &db.Statistics{
		ChainId:           db.UTchainID,
		BlockHeight:       100,
		TotalAccounts:     50,
		TotalContracts:    30,
		TotalOrgs:         10,
		TotalNodes:        5,
		TotalTransactions: 200,
	}
	err := InsertStatistics(db.UTchainID, stats1)
	if err != nil {
		t.Errorf("insertStatisticsTest failed: %v", err)
	}
}

func TestInsertStatistics(t *testing.T) {
	// Test case 1: Normal case
	stats1 := &db.Statistics{
		ChainId:           db.UTchainID,
		BlockHeight:       100,
		TotalAccounts:     50,
		TotalContracts:    30,
		TotalOrgs:         10,
		TotalNodes:        5,
		TotalTransactions: 200,
	}
	err1 := InsertStatistics(db.UTchainID, stats1)
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}

	// Test case 2: Empty chainId
	chainId2 := ""
	stats2 := &db.Statistics{
		BlockHeight:       100,
		TotalAccounts:     50,
		TotalContracts:    30,
		TotalOrgs:         10,
		TotalNodes:        5,
		TotalTransactions: 200,
	}
	err2 := InsertStatistics(chainId2, stats2)
	if err2 != nil {
		t.Errorf("Test case 2 failed: %v", err2)
	}

	// Test case 3: Nil stats
	err3 := InsertStatistics(db.UTchainID, nil)
	if err3 != nil {
		t.Errorf("Test case 3 failed: %v", err3)
	}
}

func TestGetChainStatistics(t *testing.T) {
	insertStatisticsTest(t)

	// Test case 1: Normal case
	stats, err := GetChainStatistics(db.UTchainID)
	if err != nil {
		t.Errorf("Test case 1 failed: %v", err)
	}
	if stats == nil {
		t.Error("Test case 1 failed: stats is nil")
	}

	// Test case 2: Empty chainId
	_, err = GetChainStatistics("")
	if err == nil {
		t.Error("Test case 2 failed: expected error for empty chainId")
	}
}

func TestUpdateStatisticsDelay(t *testing.T) {
	insertStatisticsTest(t)

	// Test case 1: Normal case
	stats := &db.Statistics{
		ChainId:       db.UTchainID,
		TotalAccounts: 60,
		TotalOrgs:     15,
		TotalNodes:    8,
	}
	err := UpdateStatisticsDelay(db.UTchainID, stats)
	if err != nil {
		t.Errorf("Test case 1 failed: %v", err)
	}

	// Test case 2: Empty chainId
	err = UpdateStatisticsDelay("", stats)
	if err != nil {
		t.Errorf("Test case 2 failed: %v", err)
	}

	// Test case 3: Nil stats
	err = UpdateStatisticsDelay(db.UTchainID, nil)
	if err != nil {
		t.Errorf("Test case 3 failed: %v", err)
	}
}

func TestUpdateStatisticsRealtime(t *testing.T) {
	insertStatisticsTest(t)

	// Test case 1: Normal case
	stats := &db.Statistics{
		ChainId:           db.UTchainID,
		BlockHeight:       150,
		TotalTransactions: 250,
		TotalCrossTx:      50,
		TotalContracts:    40,
	}
	err := UpdateStatisticsRealtime(db.UTchainID, stats)
	if err != nil {
		t.Errorf("Test case 1 failed: %v", err)
	}

	// Test case 2: Empty chainId
	err = UpdateStatisticsRealtime("", stats)
	if err != nil {
		t.Errorf("Test case 2 failed: %v", err)
	}

	// Test case 3: Nil stats
	err = UpdateStatisticsRealtime(db.UTchainID, nil)
	if err != nil {
		t.Errorf("Test case 3 failed: %v", err)
	}
}

func TestDeleteStatistics(t *testing.T) {
	insertStatisticsTest(t)

	// Test case 1: Normal case
	err := DeleteStatistics(db.UTchainID)
	if err != nil {
		t.Errorf("Test case 1 failed: %v", err)
	}

	// Test case 2: Empty chainId
	err = DeleteStatistics("")
	if err != nil {
		t.Errorf("Test case 2 failed: %v", err)
	}
}
