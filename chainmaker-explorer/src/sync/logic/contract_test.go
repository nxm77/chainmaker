/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package logic

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/db"
	"chainmaker_web/src/sync/common"
	"chainmaker_web/src/sync/model"
	"context"
	"encoding/json"
	"testing"
	"time"

	"chainmaker.org/chainmaker/contract-utils/standard"
	"chainmaker.org/chainmaker/pb-go/v2/accesscontrol"
	pbCommon "chainmaker.org/chainmaker/pb-go/v2/common"
	"github.com/google/go-cmp/cmp"
	"github.com/test-go/testify/assert"
)

func SetMemberInfoCache(chainId, hashType string, member *accesscontrol.Member) error {
	//缓存key
	memberKey, memberKeyErr := common.GetMemberInfoKey(chainId, hashType, int32(member.MemberType), member.MemberInfo)
	if memberKeyErr != nil {
		return memberKeyErr
	}
	userInfoJson := getUserInfoInfoTest("0_userInfoJson.json")
	userInfoJsonBytes, _ := json.Marshal(userInfoJson)
	//缓存数据
	// 设置键值对和过期时间
	ctx := context.Background()
	err := cache.GlobalRedisDb.Set(ctx, memberKey, string(userInfoJsonBytes), time.Hour).Err()
	return err
}

func TestParallelParseContract(t *testing.T) {
	blockInfo := getBlockInfoTest("1_blockInfoJsonContract.json")
	txInfo := getTxInfoInfoTest("1_txInfoJsonContract.json")
	dealResult := getDealResultTest("1_dealResultJsonContract.json")
	chainId := blockInfo.Block.Header.ChainId
	err := SetMemberInfoCache(chainId, "SHA256", txInfo.Sender.Signer)
	if err != nil {
		return
	}

	type args struct {
		blockInfo  *pbCommon.BlockInfo
		hashType   string
		dealResult *model.ProcessedBlockData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test case 1",
			args: args{
				blockInfo:  blockInfo,
				hashType:   "SHA256",
				dealResult: dealResult,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ParallelParseContract(tt.args.blockInfo, tt.args.hashType, tt.args.dealResult)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParallelParseContract() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(tt.args.dealResult.ContractWriteSetData) == 0 {
				t.Errorf("ParallelParseContract() error = %v, ContractWriteSetData is nil", err)
			}
		})
	}
}

func Test_containsAllFunctions(t *testing.T) {
	blockInfo := getBlockInfoTest("2_blockInfoJsonContractERC20.json")
	if blockInfo == nil {
		return
	}
	rwSetList := blockInfo.RwsetList[0]
	contractWriteSet, err := GetContractByWriteSet(rwSetList.TxWrites, "ea9c7f588e2bce761ae33ac5bf31092abefb1aae")
	if err != nil {
		return
	}
	if contractWriteSet.ContractResult == nil {
		return
	}

	signatures := common.ExtractFunctionSignatures(contractWriteSet.ByteCode)
	type args struct {
		evmType       string
		signatures    [][]byte
		functionNames map[string]bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test case 1",
			args: args{
				evmType:       common.ContractStandardNameEVMDFA,
				signatures:    signatures,
				functionNames: common.CopyMap(common.ERC20Functions),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = containsAllFunctions(tt.args.evmType, tt.args.signatures, tt.args.functionNames)
			//if got != tt.want {
			//	t.Errorf("containsAllFunctions() = %v, want %v", got, tt.want)
			//}
		})
	}
}

// func TestBuildContractInfo(t *testing.T) {
// 	blockInfo := getBlockInfoTest("1_blockInfoJsonContract.json")
// 	txInfo := getTxInfoInfoTest("1_txInfoJsonContract.json")
// 	userInfo := getUserInfoInfoTest("1_userInfoJsonContract.json")

// 	if blockInfo == nil || txInfo == nil || userInfo == nil {
// 		return
// 	}
// 	type args struct {
// 		i         int
// 		blockInfo *pbCommon.BlockInfo
// 		txInfo    *pbCommon.Transaction
// 		userInfo  *common.MemberAddrIdCert
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Test case 1",
// 			args: args{
// 				blockInfo: blockInfo,
// 				txInfo:    txInfo,
// 				userInfo:  userInfo,
// 			},
// 			wantErr: false,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			_, err := BuildContractInfo(tt.args.i, tt.args.blockInfo, tt.args.txInfo, tt.args.userInfo)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("BuildContractInfo() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

func TestGetContractByWriteSet(t *testing.T) {
	blockInfo := getBlockInfoTest("2_blockInfoJsonContractERC20.json")
	if blockInfo == nil {
		return
	}
	rwSetList := blockInfo.RwsetList[0]
	type args struct {
		txWriteList []*pbCommon.TxWrite
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test case 1",
			args: args{
				txWriteList: rwSetList.TxWrites,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetContractByWriteSet(tt.args.txWriteList, "ea9c7f588e2bce761ae33ac5bf31092abefb1aae")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContractByWriteSet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.ContractResult == nil || len(got.ByteCode) == 0 {
				t.Errorf("GetContractByWriteSet() got = %v", got)
			}
		})
	}
}

func Test_containsAllFunctions1(t *testing.T) {
	blockInfo := getBlockInfoTest("2_blockInfoJsonContractERC20.json")
	if blockInfo == nil {
		return
	}
	rwSetList := blockInfo.RwsetList[0]
	contractWriteSet, err := GetContractByWriteSet(rwSetList.TxWrites, "ea9c7f588e2bce761ae33ac5bf31092abefb1aae")
	if err != nil {
		return
	}
	if contractWriteSet.ContractResult == nil {
		return
	}

	signatures := common.ExtractFunctionSignatures(contractWriteSet.ByteCode)
	type args struct {
		evmType       string
		signatures    [][]byte
		functionNames map[string]bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Test case 1",
			args: args{
				evmType:       common.ContractStandardNameEVMDFA,
				signatures:    signatures,
				functionNames: common.CopyMap(common.ERC20Functions),
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = containsAllFunctions(tt.args.evmType, tt.args.signatures, tt.args.functionNames)
			//if got != tt.want {
			//	t.Errorf("containsAllFunctions() = %v, want %v", got, tt.want)
			//}
		})
	}
}

func TestGenesisBlockSystemContract(t *testing.T) {
	type args struct {
		blockInfo  *pbCommon.BlockInfo
		dealResult *model.ProcessedBlockData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "非创世区块",
			args: args{
				blockInfo: &pbCommon.BlockInfo{
					Block: &pbCommon.Block{
						Header: &pbCommon.BlockHeader{
							BlockHeight: 0,
						},
						Txs: []*pbCommon.Transaction{
							{
								Payload: &pbCommon.Payload{
									TxId: "1234",
								},
							},
						},
					},
					RwsetList: []*pbCommon.TxRWSet{
						{
							TxId: "1234",
						},
					},
				},
				dealResult: &model.ProcessedBlockData{},
			},
			wantErr: false,
		},
		// 添加其他测试用例
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GenesisBlockSystemContract(tt.args.blockInfo, tt.args.dealResult); (err != nil) != tt.wantErr {
				t.Errorf("GenesisBlockSystemContract() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_dealStandardContract(t *testing.T) {
	contract1 := &db.Contract{
		ContractType: "ERC20",
		Name:         "ContractName",
		NameBak:      "ContractName",
		Addr:         "12345",
	}
	contract2 := &db.Contract{
		ContractType: "ERC721",
		Name:         "ContractName",
		NameBak:      "ContractName",
		Addr:         "12345",
	}
	want1 := &db.FungibleContract{
		ContractType:    "ERC20",
		ContractName:    "ContractName",
		ContractNameBak: "ContractName",
		ContractAddr:    "12345",
	}
	want2 := &db.NonFungibleContract{
		ContractType:    "ERC721",
		ContractName:    "ContractName",
		ContractNameBak: "ContractName",
		ContractAddr:    "12345",
	}

	type args struct {
		contract *db.Contract
	}
	tests := []struct {
		name  string
		args  args
		want  *db.FungibleContract
		want1 *db.NonFungibleContract
	}{
		{
			name: "Test case 1",
			args: args{
				contract: contract1,
			},
			want: want1,
		},
		{
			name: "Test case 2",
			args: args{
				contract: contract2,
			},
			want1: want2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := dealStandardContract(tt.args.contract)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("dealStandardContract() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
			if !cmp.Equal(got1, tt.want1) {
				t.Errorf("dealStandardContract() got1 = %v, want1 %v\ndiff: %s", got1, tt.want1, cmp.Diff(got1, tt.want1))
			}
		})
	}
}

func TestDealIDAContractData(t *testing.T) {
	// Test case 1: Test with a fungible contract
	contract1 := &db.Contract{
		ContractType: standard.ContractStandardNameCMIDA,
		Name:         "Test Contract",
		NameBak:      "Test Contract",
		Addr:         "0x1234567890abcdef",
		Timestamp:    12344,
	}

	contract := dealIDAContractData(contract1)
	if contract == nil {
		t.Errorf("Test case 1 failed: contract is nil")
	}
}

func TestGetContractType(t *testing.T) {
	byteCodes := []byte("0x1234567890abcdef")
	runtimeType1 := common.RuntimeTypeDockerGo
	runtimeType2 := common.RuntimeTypeEVM
	// Test case 1: Test with a fungible contract
	_, err := GetContractType(db.UTchainID, "contract1", runtimeType1, byteCodes)
	if err != nil {
		t.Errorf("Test case 1 failed: contract is nil")
	}

	// Test case 2: Test with a fungible contract
	_, err = GetContractType(db.UTchainID, "contract1", runtimeType2, byteCodes)
	if err != nil {
		t.Errorf("Test case 1 failed: contract is nil")
	}
}

func TestGetContractSymbol(t *testing.T) {
	// Test case 1: Test with a valid contract
	contractType1 := common.ContractStandardNameCMDFA
	_, _ = GetContractSymbol(db.UTchainID, contractType1, "123")
	contractType2 := common.ContractStandardNameCMDFA
	_, _ = GetContractSymbol(db.UTchainID, contractType2, "123")
}

func TestGetTotalSupply(t *testing.T) {
	// Test case 1: Test with a valid contract
	contractType1 := common.ContractStandardNameCMDFA
	contractName1 := "contract12"
	_, _ = GetTotalSupply(contractType1, db.UTchainID, contractName1)
	contractType2 := common.ContractStandardNameEVMDFA
	_, _ = GetTotalSupply(contractType2, db.UTchainID, contractName1)
}

func TestGetContractDecimals(t *testing.T) {
	contractType1 := common.ContractStandardNameCMDFA
	contractName1 := "contract13"
	_, _ = GetContractDecimals(db.UTchainID, contractType1, contractName1)
	contractType2 := common.ContractStandardNameEVMDFA
	_, _ = GetContractDecimals(db.UTchainID, contractType2, contractName1)
}

func TestGetContractMapByAddrs(t *testing.T) {
	// Test case 1: Test with a valid contract
	contractAddrMap1 := map[string]string{"addr1": "addr1"}
	_, err1 := GetContractMapByAddrs(db.UTchainID, contractAddrMap1)
	assert.NoError(t, err1)
}

func TestHandleContractInsertOrUpdate(t *testing.T) {
	// Test case 1: Test with a new contract
	dealResult1 := &model.ProcessedBlockData{
		ContractWriteSetData: map[string]*model.ContractWriteSetData{
			"contract14": {
				ContractType: common.ContractStandardNameCMDFA,
				ContractName: "contract15",
				ContractAddr: "addr1",
			},
		},
	}

	err1 := HandleContractInsertOrUpdate(db.UTchainID, dealResult1)
	assert.NoError(t, err1)
}
