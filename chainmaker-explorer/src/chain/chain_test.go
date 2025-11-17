// Package chain provides chain Methods
package chain

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	_ = config.InitConfig("", "")
	// 初始化数据库配置
	db.InitRedisContainer()
	db.InitMySQLContainer()

	// 运行其他测试
	os.Exit(m.Run())
}

func insertSubscribeTest() (*db.Subscribe, error) {
	subscribeInfo := &db.Subscribe{
		ChainId:     db.UTchainID,
		UserSignKey: "1234",
		UserSignCrt: "1234",
	}

	err := dbhandle.InsertSubscribe(subscribeInfo)
	return subscribeInfo, err
}

func TestInitChainConfig(t *testing.T) {
	insertSubscribeTest()
	InitChainConfig()
}

func TestMergeChainInfo(t *testing.T) {
	configChains := []*config.ChainInfo{
		{ChainId: "c1", AuthType: "a"},
		{ChainId: "c2", AuthType: "b"},
	}
	dbChains := []*config.ChainInfo{
		{ChainId: "c2", AuthType: "bb"},
		{ChainId: "c3", AuthType: "c"},
	}
	merged := mergeChainInfo(configChains, dbChains)
	assert.Equal(t, 3, len(merged))
	assert.Equal(t, "a", merged[0].AuthType)
	assert.Equal(t, "bb", merged[1].AuthType)
	assert.Equal(t, "c", merged[2].AuthType)
}

func TestInitChainListTable(t *testing.T) {
	// 只测试分支覆盖
	oldConf := config.GlobalConfig.DBConf
	defer func() { config.GlobalConfig.DBConf = oldConf }()
	config.GlobalConfig.DBConf.DbProvider = config.MySql
	InitChainListTable([]*config.ChainInfo{})
	config.GlobalConfig.DBConf.DbProvider = config.Pgsql
	InitChainListTable([]*config.ChainInfo{})
	config.GlobalConfig.DBConf.DbProvider = config.ClickHouse
	InitChainListTable([]*config.ChainInfo{})
}

func TestSetMainDIDContract(t *testing.T) {
	chainId := "testchain"
	contract := &db.Contract{
		ContractType: "CMDID",
		NameBak:      "mainDID",
	}
	config.MainDIDContractMap = map[string]string{}
	SetMainDIDContract(chainId, contract)
	assert.Equal(t, "mainDID", config.MainDIDContractMap[chainId])

	// 已存在不覆盖
	SetMainDIDContract(chainId, contract)
	assert.Equal(t, "mainDID", config.MainDIDContractMap[chainId])

	// 非CMDID类型不设置
	contract.ContractType = "OTHER"
	SetMainDIDContract("otherchain", contract)
	_, ok := config.MainDIDContractMap["otherchain"]
	assert.False(t, ok)
}

func TestCheckIsMainDIDContract(t *testing.T) {
	chainId := "c1"
	contractName := "main"
	config.MainDIDContractMap = map[string]string{chainId: contractName}
	assert.True(t, CheckIsMainDIDContract(chainId, contractName))
	assert.False(t, CheckIsMainDIDContract(chainId, "other"))
	assert.False(t, CheckIsMainDIDContract("notfound", contractName))
}

func TestGetSubscribeChains(t *testing.T) {
	insertSubscribeTest()
	subscribeChains, err := GetSubscribeChains()
	if err != nil {
		t.Errorf("Test case 1 failed: %v", err)
	}
	if len(subscribeChains) == 0 {
		t.Errorf("Test case 1 failed: expected 0 chains, got %d", len(subscribeChains))
	}
}

func TestGetSubscribeByChainId(t *testing.T) {
	insertSubscribeTest()
	chainConfig, err := GetSubscribeByChainId(db.UTchainID)
	if err != nil {
		t.Errorf("Test case 1 failed: %v", err)
	}
	if chainConfig == nil {
		t.Errorf("Test case 1 failed: expected nil chain config, got %v", chainConfig)
	}
}

func TestInitChainMainDID(t *testing.T) {
	InitChainMainDID()
}
