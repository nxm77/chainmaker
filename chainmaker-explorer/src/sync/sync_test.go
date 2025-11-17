package sync

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/sync/common"
	"testing"
	"time"

	pbconfig "chainmaker.org/chainmaker/pb-go/v2/config"
)

func TestStartSync(t *testing.T) {
	// Test case 1: Test with multiple chains
	chainInfo := BuildChainInfo(&db.Subscribe{})
	StartSync([]*config.ChainInfo{chainInfo})
	time.Sleep(time.Second * 5) // Wait for the goroutines to finish

	// Test case 2: Test with no chains
	chainList2 := []*config.ChainInfo{}
	StartSync(chainList2)
	time.Sleep(time.Second * 5) // Wait for the goroutines to finish
}

func TestSubscribeToChain(t *testing.T) {
	// Test case 4: Test with a valid chain
	chainInfo := BuildChainInfo(&db.Subscribe{})
	SubscribeToChain(chainInfo)
}

func TestBeginSubscribeChain(t *testing.T) {
	// Test case 6: Test with a valid chain
	BeginSubscribeChain(ChainId1)

	// Test case 7: Test with an invalid chain
	chainId7 := "invalidChain7"
	BeginSubscribeChain(chainId7)
}

func TestCreateSubscribeClientPool(t *testing.T) {
	// Test case 8: Test with a valid chain
	chainInfo := BuildChainInfo(&db.Subscribe{})
	CreateSubscribeClientPool(chainInfo)
}

func TestPersistChainSubscriptionInfo(t *testing.T) {
	// Test case 10: Test with a valid chain
	chainInfo := BuildChainInfo(&db.Subscribe{})
	chainConfig10 := &pbconfig.ChainConfig{}
	PersistChainSubscriptionInfo(chainInfo, chainConfig10, db.SubscribeOK)
}

func TestReStartChain(t *testing.T) {
	// Test case 12: Test with a valid chain
	chainId12 := "chain1"
	ReStartChain(chainId12)

	// Test case 13: Test with an invalid chain
	chainId13 := "invalidChain13"
	ReStartChain(chainId13)
}

func TestDockerGetContractType(t *testing.T) {
	chainId := "chainmaker_pk"
	contractName := "official_identity"
	contractType, err := common.DockerGetContractType(chainId, contractName)
	//t.Errorf("DockerGetContractType contractType: %v, err: %v", contractType, err)
	if contractType == "" {
		t.Errorf("DockerGetContractType result is empty, err: %v", err)
	}
}
