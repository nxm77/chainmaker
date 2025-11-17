/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package logic

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/sync/model"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/shopspring/decimal"
)

func GetNumberDecimal(amount string) decimal.Decimal {
	// 将字符串转换为 decimal.Decimal 值
	amountDecimal, _ := decimal.NewFromString(amount)
	return amountDecimal
}

func TestDealContractTotalSupply(t *testing.T) {
	contractEvents := []*db.ContractEventData{
		{
			Index:        1,
			Topic:        "mint",
			TxId:         "123456789",
			ContractName: "ERC20",
			EventData: &db.TransferTopicEventData{
				ToAddress: "123456789",
				Amount:    "10000",
			},
			Timestamp: 123456789,
		},
		{
			Index:        1,
			Topic:        "burn",
			TxId:         "1234567890",
			ContractName: "ERC20",
			EventData: &db.TransferTopicEventData{
				FromAddress: "123456789",
				Amount:      "1000",
			},
			Timestamp: 123456789,
		},
		{
			Index:        1,
			Topic:        "transfer",
			TxId:         "1234567891",
			ContractName: "ERC20",
			EventData: &db.TransferTopicEventData{
				FromAddress: "123456789",
				ToAddress:   "123456781",
				Amount:      "1000",
			},
			Timestamp: 123456789,
		},
	}

	contractMap := map[string]*db.Contract{
		"ERC20": {
			Name:           "ERC20",
			NameBak:        "ERC20",
			Addr:           "112233445566",
			Version:        "1.0.0",
			RuntimeType:    "EVM",
			ContractStatus: 0,
			ContractType:   "CMDFA",
			Decimals:       1,
			BlockHeight:    10,
		},
	}

	amountDecimal := decimal.NewFromFloat(900)
	want := map[string]decimal.Decimal{
		"112233445566": amountDecimal,
	}
	type args struct {
		contractEvents []*db.ContractEventData
		contractMap    map[string]*db.Contract
	}
	tests := []struct {
		name string
		args args
		want map[string]decimal.Decimal
	}{
		{
			name: "Test case 1: Valid blockInfo and hashType",
			args: args{
				contractEvents: contractEvents,
				contractMap:    contractMap,
			},
			want: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DealContractTotalSupply(tt.args.contractEvents, tt.args.contractMap)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("DealContractTotalSupply() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDealFungibleContractUpdateData(t *testing.T) {
	amountDecimal := decimal.NewFromFloat(900)
	holdCountMap := map[string]int64{
		"112233445566":     2,
		"1122334455667788": -1,
		"223344556677":     1,
	}
	totalSupplyMap := map[string]decimal.Decimal{
		"112233445566": amountDecimal,
	}

	contractResult := &db.GetContractResult{
		UpdateFungibleContract: nil,
	}

	fungibleContracts := map[string]*db.FungibleContract{
		"112233445566": {
			ContractName:    "ERC20",
			ContractNameBak: "ERC20",
			ContractAddr:    "112233445566",
			ContractType:    "CMDFA",
			TotalSupply:     GetNumberDecimal("1000"),
			HolderCount:     100,
			BlockHeight:     10,
		},
	}

	type args struct {
		holdCountMap      map[string]int64
		totalSupplyMap    map[string]decimal.Decimal
		delayedUpdateData *model.DelayedUpdateData
		fungibleContracts map[string]*db.FungibleContract
		minHeight         int64
	}
	tests := []struct {
		name            string
		args            args
		holdCountWant   int64
		totalSupplyWant string
	}{
		{
			name: "Test case 1: Valid blockInfo and hashType",
			args: args{
				holdCountMap:   holdCountMap,
				totalSupplyMap: totalSupplyMap,
				delayedUpdateData: &model.DelayedUpdateData{
					ContractResult: contractResult,
				},
				fungibleContracts: fungibleContracts,
				minHeight:         34,
			},
			holdCountWant:   102,
			totalSupplyWant: "1900",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DealFungibleContractUpdateData(tt.args.holdCountMap, tt.args.totalSupplyMap, tt.args.fungibleContracts, tt.args.minHeight)
		})
	}
}

func TestDealNonFungibleContractUpdateData(t *testing.T) {
	amountDecimal := decimal.NewFromFloat(900)
	holdCountMap := map[string]int64{
		"112233445566":     2,
		"1122334455667788": -1,
		"223344556677":     1,
	}
	totalSupplyMap := map[string]decimal.Decimal{
		"112233445566": amountDecimal,
	}

	contractResult := &db.GetContractResult{
		UpdateNonFungible: nil,
	}

	nonFungibleMap := map[string]*db.NonFungibleContract{
		"112233445566": {
			ContractName:    "ERC20",
			ContractNameBak: "ERC20",
			ContractAddr:    "112233445566",
			ContractType:    "CMDFA",
			TotalSupply:     GetNumberDecimal("1000"),
			HolderCount:     100,
			BlockHeight:     10,
		},
	}

	type args struct {
		holdCountMap      map[string]int64
		totalSupplyMap    map[string]decimal.Decimal
		delayedUpdateData *model.DelayedUpdateData
		nonFungibleMap    map[string]*db.NonFungibleContract
		minHeight         int64
	}
	tests := []struct {
		name            string
		args            args
		holdCountWant   int64
		totalSupplyWant string
	}{
		{
			name: "Test case 1: Valid blockInfo and hashType",
			args: args{
				holdCountMap:   holdCountMap,
				totalSupplyMap: totalSupplyMap,
				delayedUpdateData: &model.DelayedUpdateData{
					ContractResult: contractResult,
				},
				nonFungibleMap: nonFungibleMap,
				minHeight:      34,
			},
			holdCountWant:   102,
			totalSupplyWant: "1900",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DealNonFungibleContractUpdateData(tt.args.holdCountMap, tt.args.totalSupplyMap, tt.args.nonFungibleMap, tt.args.minHeight)
		})
	}
}

func Test_dealContractHoldCount(t *testing.T) {
	insertFungiblePosition := []*db.FungiblePosition{
		{
			OwnerAddr:    "123456789",
			ContractAddr: "112233445566",
			ContractName: "ERC20",
			Amount:       GetNumberDecimal("200"),
			BlockHeight:  10,
		},
		{
			OwnerAddr:    "223456789",
			ContractAddr: "112233445566",
			ContractName: "ERC20",
			Amount:       GetNumberDecimal("200"),
			BlockHeight:  10,
		},
	}
	deleteFungiblePosition := []*db.FungiblePosition{
		{
			OwnerAddr:    "123456789",
			ContractAddr: "1122334455667788",
			ContractName: "ERC20",
			Amount:       GetNumberDecimal("200"),
			BlockHeight:  10,
		},
	}
	insertNonFungible := []*db.NonFungiblePosition{
		{
			OwnerAddr:    "123456789",
			ContractAddr: "223344556677",
			ContractName: "ERC20",
			Amount:       GetNumberDecimal("200"),
			BlockHeight:  10,
		},
		{
			OwnerAddr:    "223456789",
			ContractAddr: "223344556677",
			ContractName: "ERC20",
			Amount:       GetNumberDecimal("200"),
			BlockHeight:  10,
		},
	}
	deleteNonFungible := []*db.NonFungiblePosition{
		{
			OwnerAddr:    "123456789",
			ContractAddr: "223344556677",
			ContractName: "ERC20",
			Amount:       GetNumberDecimal("200"),
			BlockHeight:  10,
		},
	}

	want := map[string]int64{
		"112233445566":     2,
		"1122334455667788": -1,
		"223344556677":     1,
	}
	type args struct {
		positionOperates *db.BlockPosition
	}
	tests := []struct {
		name string
		args args
		want map[string]int64
	}{
		{
			name: "Test case 1: Valid blockInfo and hashType",
			args: args{
				positionOperates: &db.BlockPosition{
					InsertFungiblePosition: insertFungiblePosition,
					DeleteFungiblePosition: deleteFungiblePosition,
					InsertNonFungible:      insertNonFungible,
					DeleteNonFungible:      deleteNonFungible,
				},
			},
			want: want,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DealContractHoldCount(tt.args.positionOperates)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("DealContractHoldCount() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNFTContractTransferNum(t *testing.T) {
	// Test case 1: Test with a valid contract
	transferList1 := []*db.NonFungibleTransfer{
		{
			ContractAddr: "addr1",
		},
		{
			ContractAddr: "addr1",
		},
		{
			ContractAddr: "addr2",
		},
	}
	nftContractMap1 := map[string]*db.NonFungibleContract{
		"addr1": {
			ContractAddr: "addr1",
			TransferNum:  0,
		},
		"addr2": {
			ContractAddr: "addr2",
			TransferNum:  0,
		},
	}
	minHeight1 := int64(10)

	updateNFTContractMap1 := NFTContractTransferNum(transferList1, nftContractMap1, minHeight1)
	if len(updateNFTContractMap1) == 0 {
		t.Errorf("NFTContractTransferNum() returned incorrect length: got %d, want %d", len(updateNFTContractMap1), 2)
	}
}

func TestMergeNFTContractMaps(t *testing.T) {
	// Test case 1: Test with a valid contract
	minHeight1 := int64(10)
	updateNFTContractMap1 := map[string]*db.NonFungibleContract{
		"addr1": {
			ContractAddr: "addr1",
			TransferNum:  2,
		},
		"addr2": {
			ContractAddr: "addr2",
			TransferNum:  1,
		},
	}
	nftContractTransferMap1 := map[string]*db.NonFungibleContract{
		"addr1": {
			ContractAddr: "addr1",
			TransferNum:  3,
		},
		"addr2": {
			ContractAddr: "addr2",
			TransferNum:  2,
		},
	}

	mergedNFTContractSlice1 := MergeNFTContractMaps(minHeight1, updateNFTContractMap1, nftContractTransferMap1)
	if len(mergedNFTContractSlice1) == 0 {
		t.Errorf("MergeNFTContractMaps() returned incorrect length: got %d, want %d", len(mergedNFTContractSlice1), 2)
	}
}

func TestFTContractTransferNum(t *testing.T) {
	// Test case 1: Test with a valid contract
	transferList1 := []*db.FungibleTransfer{
		{
			ContractAddr: "addr1",
		},
		{
			ContractAddr: "addr1",
		},
		{
			ContractAddr: "addr2",
		},
	}
	ftContractMap1 := map[string]*db.FungibleContract{
		"addr1": {
			ContractAddr: "addr1",
			TransferNum:  0,
		},
		"addr2": {
			ContractAddr: "addr2",
			TransferNum:  0,
		},
	}
	minHeight1 := int64(10)

	updateFTContractMap1 := FTContractTransferNum(transferList1, ftContractMap1, minHeight1)
	if len(updateFTContractMap1) == 0 {
		t.Errorf("FTContractTransferNum() returned incorrect length: got %d, want %d", len(updateFTContractMap1), 2)
	}
}

func TestMergeFTContractMaps(t *testing.T) {
	// Test case 1: Test with a valid contract
	minHeight1 := int64(10)
	updateFTContractMap1 := map[string]*db.FungibleContract{
		"addr1": {
			ContractAddr: "addr1",
			TransferNum:  2,
		},
		"addr2": {
			ContractAddr: "addr2",
			TransferNum:  1,
		},
	}
	ftContractTransferMap1 := map[string]*db.FungibleContract{
		"addr1": {
			ContractAddr: "addr1",
			TransferNum:  3,
		},
		"addr2": {
			ContractAddr: "addr2",
			TransferNum:  2,
		},
	}

	mergedFTContractSlice1 := MergeFTContractMaps(minHeight1, updateFTContractMap1, ftContractTransferMap1)
	if len(mergedFTContractSlice1) == 0 {
		t.Errorf("MergeFTContractMaps() returned incorrect length: got %d, want %d", len(mergedFTContractSlice1), 2)
	}
}
