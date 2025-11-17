package sync

import (
	"chainmaker_web/src/config"
	"testing"
	"time"

	client "chainmaker_web/src/sync/clients"

	pbconfig "chainmaker.org/chainmaker/pb-go/v2/config"

	"chainmaker.org/chainmaker/pb-go/v2/discovery"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
)

func newClientTest() *client.SdkClient {
	chainInfo4 := &config.ChainInfo{
		ChainId:   "chain4",
		AuthType:  config.PUBLIC,
		HashType:  "SHA256",
		OrgId:     "org4",
		NodesList: []*config.NodeInfo{},
		UserInfo: &config.UserInfo{
			UserSignKey: "userSignKey4",
			UserSignCrt: "userSignCrt4",
			UserEncKey:  "userEncKey4",
			UserEncCrt:  "userEncCrt4",
		},
		TlsMode: config.TlsModelSingle,
	}
	newClient := client.NewSdkClient(chainInfo4, nil)
	return newClient
}

func TestPeriodicLoadStart(t *testing.T) {
	// Test case 1: Test with valid sdkClient
	sdkClient1 := newClientTest()
	go PeriodicLoadStart(sdkClient1)
	time.Sleep(2 * time.Second)

	// Test case 3: Test with sdkClient status as STOP
	sdkClient1.Status = client.STOP
	go PeriodicLoadStart(sdkClient1)
	time.Sleep(2 * time.Second)
}

func TestLoadChainRefInfos(t *testing.T) {
	// Test case 1: Test with valid sdkClient
	sdkClient1 := newClientTest()
	go loadChainRefInfos(sdkClient1)
}

func TestLoadChainUser(t *testing.T) {
	// Test case 1: Test with valid sdkClient
	sdkClient1 := newClientTest()
	loadChainUser(sdkClient1)

	// Test case 3: Test with sdkClient status as STOP
	sdkClient1.Status = client.STOP
	loadChainUser(sdkClient1)
}

func TestGetAdminUserByConfig(t *testing.T) {
	// Test case 1: Test with valid chainConfig
	chainConfig1 := &pbconfig.ChainConfig{
		ChainId: "chain1",
		TrustRoots: []*pbconfig.TrustRootConfig{
			{
				Root: []string{"root1"},
			},
		},
		Crypto: &pbconfig.CryptoConfig{
			Hash: "hash1",
		},
	}
	GetAdminUserByConfig(chainConfig1)

	// Test case 2: Test with nil chainConfig
	chainConfig2 := (*pbconfig.ChainConfig)(nil)
	GetAdminUserByConfig(chainConfig2)

	// Test case 3: Test with empty chainConfig
	chainConfig3 := &pbconfig.ChainConfig{}
	GetAdminUserByConfig(chainConfig3)
}

func TestLoadOrgInfo(t *testing.T) {
	// Test case 1: Test with valid sdkClient
	sdkClient1 := newClientTest()
	loadOrgInfo(sdkClient1)

	// Test case 3: Test with sdkClient status as STOP
	sdkClient1.Status = client.STOP
	loadOrgInfo(sdkClient1)
}

func TestLoadNodeInfo(t *testing.T) {
	// Test case 1: Test with valid sdkClient
	sdkClient1 := newClientTest()
	loadNodeInfo(sdkClient1)

	// Test case 3: Test with sdkClient status as STOP
	sdkClient1.Status = client.STOP
	loadNodeInfo(sdkClient1)
}

func TestParseAllNodeList(t *testing.T) {
	// Test case 1: Test with valid input
	nodeList1 := map[string]string{"node1": "org1", "node2": "org2"}
	consensusNodeOrgList1 := map[string]string{"node1": "org1", "node2": "org2"}
	chainConfigNodeList1 := map[string]*discovery.Node{
		"node1": {NodeAddress: "address1"},
		"node2": {NodeAddress: "address2"},
	}
	parseAllNodeList(nodeList1, consensusNodeOrgList1, chainConfigNodeList1)
}

func TestGetAllNodeList(t *testing.T) {
	// Test case 2: Test with nil chainClient
	chainClient2 := (*sdk.ChainClient)(nil)
	GetAllNodeList(ChainId1, chainClient2)
}

func TestGetChainNodeList(t *testing.T) {
	// Test case 2: Test with nil chainClient
	chainClient2 := (*sdk.ChainClient)(nil)
	GetChainNodeList(chainClient2)
}

func TestGetChainNodeData(t *testing.T) {
	// Test case 1: Test with valid sdkClient
	sdkClient1 := newClientTest()
	_, _ = GetChainNodeData(sdkClient1)

	// Test case 3: Test with sdkClient status as STOP
	sdkClient1.Status = client.STOP
	_, _ = GetChainNodeData(sdkClient1)
}
