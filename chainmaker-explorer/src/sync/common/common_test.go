package common

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"os"
	"testing"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	tcipCommon "chainmaker.org/chainmaker/tcip-go/v2/common"
	"github.com/test-go/testify/assert"
)

func TestMain(m *testing.M) {
	_ = config.InitConfig("", "")
	// 初始化数据库配置
	db.InitRedisContainer()
	db.InitMySQLContainer()

	// 运行其他测试
	os.Exit(m.Run())
}

func TestRemoveAddrPrefix(t *testing.T) {
	address := "0x1234567890abcdef"
	expected := "1234567890abcdef"
	result := RemoveAddrPrefix(address)
	if result != expected {
		t.Errorf("RemoveAddrPrefix(%s) returned %s, expected %s", address, result, expected)
	}
}

func TestIsZeroAddress(t *testing.T) {
	address := "0000000000000000000000000000000000000000"
	result := IsZeroAddress(address)
	if !result {
		t.Errorf("IsZeroAddress(%s) returned %v, expected %v", address, result, true)
	}
}

func TestStringAmountDecimal(t *testing.T) {
	amount := "12345678901234567890"
	decimals := 18
	expected := "12.34567890123456789"
	result := StringAmountDecimal(amount, decimals).String()
	if result != expected {
		t.Errorf("StringAmountDecimal(%s, %d) returned %s, expected %s", amount, decimals, result, expected)
	}
}

func TestCopyMap(t *testing.T) {
	src := map[string]bool{"a": true, "b": false}
	expected := map[string]bool{"a": true, "b": false}
	result := CopyMap(src)
	if result == nil {
		t.Errorf("CopyMap(%v) returned %v, expected %v", src, result, expected)
	}
}

func TestGetMemberInfoKey(t *testing.T) {
	chainId := "chain1"
	hashType := "hashType1"
	memberType := int32(1)
	memberBytes := []byte("memberBytes1")
	prefix := "prefix"
	config.GlobalConfig.RedisDB.Prefix = prefix
	_, err := GetMemberInfoKey(chainId, hashType, memberType, memberBytes)
	if err != nil {
		t.Errorf("GetMemberInfoKey(%s, %s, %d, %v) returned error %v", chainId, hashType, memberType, memberBytes, err)
	}
}

func TestMD5(t *testing.T) {
	str := "test1"
	expected := "5a105e8b9d40e1329780d62ea2265d8a"
	result := MD5(str)
	if result != expected {
		t.Errorf("MD5(%s) returned %s, expected %s", str, result, expected)
	}
}

func TestIsContractTx(t *testing.T) {
	txInfo := &common.Transaction{
		Payload: &common.Payload{
			ContractName: syscontract.SystemContract_CONTRACT_MANAGE.String(),
		},
		Result: &common.Result{
			ContractResult: &common.ContractResult{
				Code: 0,
			},
		},
	}

	result := IsContractTx(txInfo)
	if !result {
		t.Errorf("IsContractTx(%v) returned %v, expected %v", txInfo, result, true)
	}

	txInfo2 := &common.Transaction{
		Payload: &common.Payload{
			ContractName: "123",
		},
		Result: &common.Result{
			ContractResult: &common.ContractResult{
				Code: 1,
			},
		},
	}

	result2 := IsContractTx(txInfo2)
	if result2 {
		t.Errorf("IsContractTx(%v) returned %v, expected %v", txInfo2, result2, true)
	}
}

func TestIsConfigTx(t *testing.T) {
	txInfo := &common.Transaction{
		Payload: &common.Payload{
			ContractName: syscontract.SystemContract_CHAIN_CONFIG.String(),
		},
	}
	result := IsConfigTx(txInfo)
	if !result {
		t.Errorf("IsConfigTx(%v) returned %v, expected %v", txInfo, result, true)
	}
}

func TestIsRelayCrossChainTx(t *testing.T) {
	txInfo := &common.Transaction{
		Payload: &common.Payload{
			ContractName: syscontract.SystemContract_RELAY_CROSS.String(),
		},
	}
	result := IsRelayCrossChainTx(txInfo)
	if !result {
		t.Errorf("IsRelayCrossChainTx(%v) returned %v, expected %v", txInfo, result, true)
	}
}

func TestIsSubChainSpvContractTx(t *testing.T) {
	txInfo := &common.Transaction{
		Payload: &common.Payload{
			ContractName: SubChainSpvPrefix + "test",
		},
	}
	result, _ := IsSubChainSpvContractTx(txInfo)
	if !result {
		t.Errorf("IsSubChainSpvContractTx(%v) returned %v, expected %v", txInfo, result, true)
	}
}

func TestIsInBlockHeight(t *testing.T) {
	height := int64(10)
	heightList := []int64{5, 10, 15}
	result := IsInBlockHeight(height, heightList)
	if !result {
		t.Errorf("IsInBlockHeight(%d, %v) returned %v, expected %v", height, heightList, result, true)
	}
}

func TestGetMaxBlockHeight(t *testing.T) {
	heightList := []int64{5, 10, 15}
	expected := int64(15)
	result := GetMaxBlockHeight(heightList)
	if result != expected {
		t.Errorf("GetMaxBlockHeight(%v) returned %d, expected %d", heightList, result, expected)
	}
}

func TestGetMinBlockHeight(t *testing.T) {
	heightList := []int64{5, 10, 15}
	expected := int64(5)
	result := GetMinBlockHeight(heightList)
	if result != expected {
		t.Errorf("GetMinBlockHeight(%v) returned %d, expected %d", heightList, result, expected)
	}
}

func TestIsMainChainGateway(t *testing.T) {
	gatewayID := tcipCommon.MainGateway_MAIN_GATEWAY_ID.String()
	result := IsMainChainGateway(gatewayID)
	if !result {
		t.Errorf("IsMainChainGateway(%s) returned %v, expected %v", gatewayID, result, true)
	}
}

func TestParallelParseBatchWhere(t *testing.T) {
	wheres := []string{"where1", "where2", "where3", "where4"}
	batchSize := 2
	expected := [][]string{{"where1", "where2"}, {"where3", "where4"}}
	result := ParallelParseBatchWhere(wheres, batchSize)
	if !sliceEqual(result, expected) {
		t.Errorf("ParallelParseBatchWhere(%v, %d) returned %v, expected %v", wheres, batchSize, result, expected)
	}
}

func TestIsCrossEnd(t *testing.T) {
	status := int32(tcipCommon.CrossChainStateValue_CONFIRM_END)
	result := IsCrossEnd(status)
	if !result {
		t.Errorf("IsCrossEnd(%d) returned %v, expected %v", status, result, true)
	}
}

func TestExtractTxIdsAndContractNames(t *testing.T) {
	txInfoList := []*db.Transaction{
		{
			TxId:         "tx1",
			ContractAddr: "addr1",
		},
		{
			TxId:         "tx2",
			ContractAddr: "addr2",
		},
	}
	expectedTxIds := []string{"tx1", "tx2"}
	expectedContractAddrMap := map[string]string{
		"addr1": "addr1",
		"addr2": "addr2",
	}

	txIds, contractAddrMap, _ := ExtractTxIdsAndContractNames(txInfoList)
	if !sliceEqual([][]string{txIds}, [][]string{expectedTxIds}) {
		t.Errorf("ExtractTxIdsAndContractNames(%v) returned txIds %v, expected %v", txInfoList, txIds, expectedTxIds)
	}
	if !mapEqual(contractAddrMap, expectedContractAddrMap) {
		t.Errorf("ExtractTxIdsAndContractNames(%v) returned contractAddrMap %v, expected %v", txInfoList, contractAddrMap, expectedContractAddrMap)
	}

}

func mapEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}

func sliceEqual(a, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func TestGetEvmAbi(t *testing.T) {
	// 测试 ERC20 ABI
	abiERC20 := GetEvmAbi(ContractStandardNameEVMDFA)
	if abiERC20 == nil {
		t.Errorf("Expected non-nil ABI for ERC20, got nil")
	}

	// 确认 ABI 解析正确
	if len(abiERC20.Methods) == 0 {
		t.Errorf("Expected non-empty methods for ERC20 ABI, got empty")
	}

	// 测试 ERC721 ABI
	abiERC721 := GetEvmAbi(ContractStandardNameEVMNFA)
	if abiERC721 == nil {
		t.Errorf("Expected non-nil ABI for ERC721, got nil")
	}

	// 确认 ABI 解析正确
	if len(abiERC721.Methods) == 0 {
		t.Errorf("Expected non-empty methods for ERC721 ABI, got empty")
	}
}

func TestExtractTxIdsAndContractNames1(t *testing.T) {
	txInfoList := []*db.Transaction{
		{
			TxId:         "txId1",
			ContractAddr: "contractAddr1",
		},
		{
			TxId:         "txId2",
			ContractAddr: "contractAddr2",
		},
	}
	txIds, contractAddrMap, txInfoMap := ExtractTxIdsAndContractNames(txInfoList)
	expectedTxIds := []string{"txId1", "txId2"}
	expectedContractAddrMap := map[string]string{"contractAddr1": "contractAddr1", "contractAddr2": "contractAddr2"}
	expectedTxInfoMap := map[string]*db.Transaction{
		"txId1": txInfoList[0],
		"txId2": txInfoList[1],
	}
	assert.Equal(t, expectedTxIds, txIds)
	assert.Equal(t, expectedContractAddrMap, contractAddrMap)
	assert.Equal(t, expectedTxInfoMap, txInfoMap)
}

func TestGetContractSendTxId(t *testing.T) {
	txInfo := &common.Transaction{
		Payload: &common.Payload{
			ContractName: syscontract.SystemContract_MULTI_SIGN.String(),
			Parameters: []*common.KeyValuePair{
				{
					Key:   syscontract.MultiVote_TX_ID.String(),
					Value: []byte("txId"),
				},
			},
		},
	}
	expected := "txId"
	result := GetContractSendTxId(txInfo)
	assert.Equal(t, expected, result)
}

func TestIsMultiSignTx(t *testing.T) {
	contractName := syscontract.SystemContract_MULTI_SIGN.String()
	contractMethod := syscontract.MultiSignFunction_REQ.String()
	result := IsMultiSignTx(contractName, contractMethod)
	assert.True(t, result)
}

func TestIsContractManageTx(t *testing.T) {
	contractName := syscontract.SystemContract_CONTRACT_MANAGE.String()
	contractMethod := syscontract.ContractManageFunction_INIT_CONTRACT.String()
	result := IsContractManageTx(contractName, contractMethod)
	assert.True(t, result)
}

func TestIsContractTxByName(t *testing.T) {
	contractName := syscontract.SystemContract_CONTRACT_MANAGE.String()
	contractMethod := syscontract.ContractManageFunction_INIT_CONTRACT.String()
	result := IsContractTxByName(contractName, contractMethod)
	assert.True(t, result)
}
