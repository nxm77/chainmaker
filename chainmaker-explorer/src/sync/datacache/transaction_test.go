package datacache

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"os"
	"testing"
)

const ChainId = "testChainId"

func TestMain(m *testing.M) {
	_ = config.InitConfig("", "")
	// 初始化数据库配置
	db.InitRedisContainer()
	db.InitMySQLContainer()

	// 运行其他测试
	os.Exit(m.Run())
}

func TestBuildLatestTxListCache(t *testing.T) {
	// 准备测试数据
	chainId := ChainId
	txMap := map[string]*db.Transaction{
		"tx1": {
			TxId:        "tx1",
			BlockHeight: 5,
			TxIndex:     2,
		},
		"tx2": {
			TxId:        "tx2",
			BlockHeight: 4,
			TxIndex:     1,
		},
		"tx3": {
			TxId:        "tx3",
			BlockHeight: 6,
			TxIndex:     3,
		},
	}

	// 调用 BuildLatestTxListCache 函数
	BuildLatestTxListCache(chainId, txMap)

	// 从缓存中获取交易列表
	txList, err := GetLatestTxListCache(chainId)
	if err != nil {
		return
	}

	// 检查交易列表长度是否为 3
	if len(txList) != 3 {
		t.Errorf("Expected txList length to be 3, got %d", len(txList))
	}

	// 检查交易列表是否按照预期排序
	expectedOrder := []string{"tx3", "tx1", "tx2"}
	for i, tx := range txList {
		if tx.TxId != expectedOrder[i] {
			t.Errorf("Expected txList order to be %v, got %v", expectedOrder, txList)
		}
	}
}

func TestGetLatestTxListCache(t *testing.T) {
	// 准备测试数据
	chainId := ChainId
	txMap := map[string]*db.Transaction{
		"tx1": {
			TxId:        "tx1",
			BlockHeight: 5,
			TxIndex:     2,
		},
		"tx2": {
			TxId:        "tx2",
			BlockHeight: 4,
			TxIndex:     1,
		},
		"tx3": {
			TxId:        "tx3",
			BlockHeight: 6,
			TxIndex:     3,
		},
	}

	// 调用 BuildLatestTxListCache 函数
	BuildLatestTxListCache(chainId, txMap)

	// 从缓存中获取交易列表
	txList, err := GetLatestTxListCache(chainId)
	if err != nil {
		return
	}

	// 检查交易列表长度是否为 3
	if len(txList) != 3 {
		t.Errorf("Expected txList length to be 3, got %d", len(txList))
	}

	// 检查交易列表是否按照预期排序
	expectedOrder := []string{"tx3", "tx1", "tx2"}
	for i, tx := range txList {
		if tx.TxId != expectedOrder[i] {
			t.Errorf("Expected txList order to be %v, got %v", expectedOrder, txList)
		}
	}
}

func TestSetLatestTxListCache(t *testing.T) {
	// 准备测试数据
	chainId := ChainId
	txList := []*db.Transaction{
		{
			TxId:        "tx1",
			BlockHeight: 5,
			TxIndex:     2,
		},
		{
			TxId:        "tx2",
			BlockHeight: 4,
			TxIndex:     1,
		},
		{
			TxId:        "tx3",
			BlockHeight: 6,
			TxIndex:     3,
		},
	}

	// 调用 BuildLatestTxListCache 函数
	SetLatestTxListCache(chainId, txList)

	// 从缓存中获取交易列表
	txListCache, err := GetLatestTxListCache(chainId)
	if err != nil {
		return
	}

	// 检查交易列表长度是否为 3
	if len(txListCache) != 3 {
		t.Errorf("Expected txList length to be 3, got %d", len(txList))
	}

	// 检查交易列表是否按照预期排序
	expectedOrder := []string{"tx3", "tx1", "tx2"}
	for i, tx := range txListCache {
		if tx.TxId != expectedOrder[i] {
			t.Errorf("Expected txList order to be %v, got %v", expectedOrder, txList)
		}
	}
}

func TestGetLatestTxListCache1(t *testing.T) {
	// 准备测试数据
	chainId := ChainId
	txMap := map[string]*db.Transaction{
		"tx1": {
			TxId:        "tx1",
			BlockHeight: 5,
			TxIndex:     2,
		},
		"tx2": {
			TxId:        "tx2",
			BlockHeight: 4,
			TxIndex:     1,
		},
		"tx3": {
			TxId:        "tx3",
			BlockHeight: 6,
			TxIndex:     3,
		},
	}

	// 调用 BuildLatestTxListCache 函数
	BuildLatestTxListCache(chainId, txMap)

	// 从缓存中获取交易列表
	txList, err := GetLatestTxListCache(chainId)
	if err != nil {
		return
	}

	// 检查交易列表长度是否为 3
	if len(txList) != 3 {
		t.Errorf("Expected txList length to be 3, got %d", len(txList))
	}

	// 检查交易列表是否按照预期排序
	expectedOrder := []string{"tx3", "tx1", "tx2"}
	for i, tx := range txList {
		if tx.TxId != expectedOrder[i] {
			t.Errorf("Expected txList order to be %v, got %v", expectedOrder, txList)
		}
	}
}

func TestBuildOverviewTxTotalCache(t *testing.T) {
	BuildOverviewTxTotalCache(db.UTchainID, 100)
}
