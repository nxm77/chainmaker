package logic

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/sync/common"
	"chainmaker_web/src/sync/model"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sync"
	"testing"

	"chainmaker.org/chainmaker/contract-utils/standard"
	pbCommon "chainmaker.org/chainmaker/pb-go/v2/common"
	pbConfig "chainmaker.org/chainmaker/pb-go/v2/config"
	tcipCommon "chainmaker.org/chainmaker/tcip-go/v2/common"
)

type JsonData struct {
	UpdatePositionData map[string]*db.PositionData `json:"UpdatePositionData"`
	ResultPositionData ResultPositionData          `json:"ResultPositionData"`
}

type ResultPositionData struct {
	FungiblePosition    []*db.FungiblePosition    `json:"FungiblePosition"`
	NonFungiblePosition []*db.NonFungiblePosition `json:"NonFungiblePosition"`
}

func geFileData(fileName string) (*json.Decoder, error) {
	file, err := os.Open("../testData/" + fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}

	// 解码 JSON 文件内容到 blockInfo 结构体
	decoder := json.NewDecoder(file)
	return decoder, err
}

func getUpdatePositionDataTest(fileName string) (map[string]*db.PositionData, []*db.FungiblePosition,
	[]*db.NonFungiblePosition) {
	// 打开 JSON 文件
	decoder, err := geFileData(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, nil, nil
	}

	// 解码 JSON 文件内容到 blockInfo 结构体
	var jsonData JsonData
	err = decoder.Decode(&jsonData)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil, nil, nil
	}
	return jsonData.UpdatePositionData, jsonData.ResultPositionData.FungiblePosition,
		jsonData.ResultPositionData.NonFungiblePosition
}

func getPositionDBJsonTest(fileName string) map[string][]*db.FungiblePosition {
	resultValue := make(map[string][]*db.FungiblePosition, 0)
	// 打开 JSON 文件
	decoder, err := geFileData(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return resultValue
	}

	// 解码 JSON 文件内容到 blockInfo 结构体
	err = decoder.Decode(&resultValue)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return resultValue
	}
	return resultValue
}

func getNonPositionDBJsonTest(fileName string) map[string][]*db.NonFungiblePosition {
	resultValue := make(map[string][]*db.NonFungiblePosition, 0)
	// 打开 JSON 文件
	decoder, err := geFileData(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return resultValue
	}
	err = decoder.Decode(&resultValue)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return resultValue
	}
	return resultValue
}

func getGotTokenResultTest(fileName string) *db.TokenResult {
	resultValue := &db.TokenResult{}
	// 打开 JSON 文件
	decoder, err := geFileData(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return resultValue
	}
	err = decoder.Decode(&resultValue)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return resultValue
	}
	return resultValue
}

func getChainListConfigTest(fileName string) []*config.ChainInfo {
	resultValue := make([]*config.ChainInfo, 0)
	// 打开 JSON 文件
	decoder, err := geFileData(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return resultValue
	}
	err = decoder.Decode(&resultValue)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return resultValue
	}
	return resultValue
}

func getChainTxRWSetTest(fileName string) *pbCommon.TxRWSet {
	var resultValue *pbCommon.TxRWSet
	// 打开 JSON 文件
	decoder, err := geFileData(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return resultValue
	}
	err = decoder.Decode(&resultValue)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return resultValue
	}
	return resultValue
}

func getBlockInfoTest(fileName string) *pbCommon.BlockInfo {
	var blockInfo *pbCommon.BlockInfo
	// 打开 JSON 文件

	decoder, err := geFileData(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return blockInfo
	}

	// 解码 JSON 文件内容到 blockInfo 结构体
	err = decoder.Decode(&blockInfo)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return blockInfo
	}
	return blockInfo
}

func getDealResultTest(fileName string) *model.ProcessedBlockData {
	var dealResult *model.ProcessedBlockData
	// 打开 JSON 文件
	decoder, err := geFileData(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return dealResult
	}
	// 解码 JSON 文件内容到 blockInfo 结构体
	err = decoder.Decode(&dealResult)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return dealResult
	}
	return dealResult
}

func getTxInfoInfoTest(fileName string) *pbCommon.Transaction {
	var txInfo *pbCommon.Transaction
	// 打开 JSON 文件
	decoder, err := geFileData(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return txInfo
	}
	// 解码 JSON 文件内容到 blockInfo 结构体
	err = decoder.Decode(&txInfo)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return txInfo
	}
	return txInfo
}

func getBuildTxInfoTest(fileName string) *db.Transaction {
	var txInfo *db.Transaction
	// 打开 JSON 文件
	decoder, err := geFileData(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return txInfo
	}
	// 解码 JSON 文件内容到 blockInfo 结构体
	err = decoder.Decode(&txInfo)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return txInfo
	}
	return txInfo
}

func getUserInfoInfoTest(fileName string) *common.MemberAddrIdCert {
	var userInfo *common.MemberAddrIdCert
	// 打开 JSON 文件
	decoder, err := geFileData(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return userInfo
	}
	// 解码 JSON 文件内容到 blockInfo 结构体
	err = decoder.Decode(&userInfo)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return userInfo
	}
	return userInfo
}

func getCrossChainInfoTest(fileName string) *tcipCommon.CrossChainInfo {
	var resultValue *tcipCommon.CrossChainInfo
	// 打开 JSON 文件
	decoder, err := geFileData(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return resultValue
	}
	// 解码 JSON 文件内容到 blockInfo 结构体
	err = decoder.Decode(&resultValue)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return resultValue
	}
	return resultValue
}

func getContractEventTest(fileName string) []*db.ContractEventData {
	contractEvents := make([]*db.ContractEventData, 0)
	// 打开 JSON 文件
	//file, err := os.Open("../testData/1_blockInfoJsonContract.json")
	decoder, err := geFileData(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return contractEvents
	}
	err = decoder.Decode(&contractEvents)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return contractEvents
	}
	return contractEvents
}

func getContractInfoMapTest(fileName string) map[string]*db.Contract {
	resultValue := make(map[string]*db.Contract, 0)
	// 打开 JSON 文件
	decoder, err := geFileData(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return resultValue
	}

	err = decoder.Decode(&resultValue)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return resultValue
	}
	return resultValue
}

func getAccountMapTest(fileName string) map[string]*db.Account {
	resultValue := make(map[string]*db.Account, 0)
	// 打开 JSON 文件
	decoder, err := geFileData(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return resultValue
	}
	err = decoder.Decode(&resultValue)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return resultValue
	}
	return resultValue
}

func getPositionListJsonTest(fileName string) map[string]*db.PositionData {
	resultValue := make(map[string]*db.PositionData, 0)
	// 打开 JSON 文件
	decoder, err := geFileData(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return resultValue
	}
	err = decoder.Decode(&resultValue)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return resultValue
	}
	return resultValue
}

func getIdaAssetTest(fileName string) *standard.IDAInfo {
	var idaInfo *standard.IDAInfo
	// 打开 JSON 文件
	decoder, err := geFileData(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return idaInfo
	}
	// 解码 JSON 文件内容到 blockInfo 结构体
	err = decoder.Decode(&idaInfo)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return idaInfo
	}
	return idaInfo
}

func TestGetChainConfigByWriteSet(t *testing.T) {
	type args struct {
		txRWSet *pbCommon.TxRWSet
		txInfo  *pbCommon.Transaction
	}
	txInfo := getTxInfoInfoTest("1_txInfoJsonContract.json")
	tests := []struct {
		name    string
		args    args
		want    *pbConfig.ChainConfig
		wantErr bool
	}{
		{
			name: "test case 1",
			args: args{
				txRWSet: nil,
				txInfo:  txInfo,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetChainConfigByWriteSet(tt.args.txRWSet, tt.args.txInfo)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetChainConfigByWriteSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetChainConfigByWriteSet() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCrossChainInfoByWriteSet(t *testing.T) {
	txRWSet := getChainTxRWSetTest("cross_1_txRWSet.json")
	crossChainInfo := getCrossChainInfoTest("cross_1_ChainInfoJson.json")

	type args struct {
		txRWSet *pbCommon.TxRWSet
	}
	tests := []struct {
		name    string
		args    args
		want    *tcipCommon.CrossChainInfo
		wantErr bool
	}{
		{
			name: "test case 1",
			args: args{
				txRWSet: txRWSet,
			},
			want:    crossChainInfo,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCrossChainInfoByWriteSet(tt.args.txRWSet)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossChainInfoByWriteSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCrossChainInfoByWriteSet() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParallelParseWriteSetData(t *testing.T) {
	blockInfo := getBlockInfoTest("2_blockInfoJson_1704945589.json")
	if blockInfo == nil || len(blockInfo.RwsetList) == 0 {
		return
	}

	// Test case 1: Normal case with multiple transactions
	dealResult1 := &model.ProcessedBlockData{}
	err1 := ParallelParseWriteSetData(blockInfo, dealResult1)
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}

	// Test case 2: Empty block
	blockInfo2 := &pbCommon.BlockInfo{
		Block: &pbCommon.Block{
			Header: &pbCommon.BlockHeader{
				ChainId:     ChainId,
				BlockHeight: 1,
			},
			Txs: []*pbCommon.Transaction{},
		},
	}
	dealResult2 := &model.ProcessedBlockData{}
	err2 := ParallelParseWriteSetData(blockInfo2, dealResult2)
	if err2 != nil {
		t.Errorf("Test case 2 failed: %v", err2)
	}

	// Test case 3: Block with only one transaction
	blockInfo3 := &pbCommon.BlockInfo{
		Block: &pbCommon.Block{
			Header: &pbCommon.BlockHeader{
				ChainId:     ChainId,
				BlockHeight: 1,
			},
			Txs: []*pbCommon.Transaction{
				{
					Payload: &pbCommon.Payload{
						TxType:    pbCommon.TxType_INVOKE_CONTRACT,
						Timestamp: 1625097700,
					},
				},
			},
		},
	}
	dealResult3 := &model.ProcessedBlockData{}
	err3 := ParallelParseWriteSetData(blockInfo3, dealResult3)
	if err3 != nil {
		t.Errorf("Test case 3 failed: %v", err3)
	}
}

func TestProcessWriteSetDataOther(t *testing.T) {
	var mutx sync.Mutex
	// Test case 1: Normal case with valid write set
	rwSetList1 := &pbCommon.TxRWSet{
		TxWrites: []*pbCommon.TxWrite{
			{
				Key:   []byte("chain12"),
				Value: []byte("value1"),
			},
		},
	}
	txInfo1 := &pbCommon.Transaction{
		Payload: &pbCommon.Payload{
			TxType:    pbCommon.TxType_INVOKE_CONTRACT,
			Timestamp: 1625097700,
		},
	}
	dealResult1 := &model.ProcessedBlockData{}
	err1 := processWriteSetDataOther(&mutx, rwSetList1, txInfo1, dealResult1)
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}

	// Test case 2: Empty write set
	rwSetList2 := &pbCommon.TxRWSet{
		TxWrites: []*pbCommon.TxWrite{},
	}
	txInfo2 := &pbCommon.Transaction{
		Payload: &pbCommon.Payload{
			TxType:    pbCommon.TxType_INVOKE_CONTRACT,
			Timestamp: 1625097700,
		},
	}
	dealResult2 := &model.ProcessedBlockData{}
	err2 := processWriteSetDataOther(&mutx, rwSetList2, txInfo2, dealResult2)
	if err2 != nil {
		t.Errorf("Test case 2 failed: %v", err2)
	}
}

func TestUpdateCrossChainResultTx(t *testing.T) {
	// 初始化 CrossChainResult
	crossResult := &model.CrossChainResult{
		CrossMainTransaction: []*db.CrossMainTransaction{},
		CrossTransfer:        make(map[string]*db.CrossTransactionTransfer),
		BusinessTxMap:        make(map[string]*db.CrossBusinessTransaction),
	}
	mainTx := &db.CrossMainTransaction{TxId: "tx123"}
	crossTxTransfer := &db.CrossTransactionTransfer{CrossId: "cross123"}
	// 模拟业务交易
	businessTxMap := map[string]*db.CrossBusinessTransaction{
		"biz1": {TxId: "biz1"},
	}
	transferList := []*db.CrossTransactionTransfer{crossTxTransfer}

	UpdateCrossChainResultTx(crossResult, mainTx, transferList, businessTxMap)
}
