/*
Package logic comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package logic

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/sync/model"
	"testing"
)

func TestNewAccountHandler(t *testing.T) {
	chainId := db.UTchainID
	minHeight := int64(10)
	eventResults := &model.TopicEventResult{}
	delayedUpdateCache := &model.GetRealtimeCacheData{
		TxList: map[string]*db.Transaction{
			"tx1": {UserAddr: "addr1"},
		},
	}
	delayGetDBResult := &model.GetDBResult{
		AccountDBMap: map[string]*db.Account{
			"addr1": {Address: "addr1"},
		},
	}
	transferEvents := []*db.ContractEventData{
		{Topic: "transfer"},
	}
	handler := NewAccountHandler(chainId, minHeight, eventResults, delayedUpdateCache, delayGetDBResult, transferEvents)
	if handler.ChainId != chainId || handler.MinHeight != minHeight {
		t.Errorf("NewAccountHandler basic fields not set correctly")
	}
	if len(handler.TxList) != 1 || len(handler.TransferEvents) != 1 {
		t.Errorf("NewAccountHandler TxList or TransferEvents not set correctly")
	}
	if handler.AccountMap["addr1"].Address != "addr1" {
		t.Errorf("NewAccountHandler AccountMap not set correctly")
	}
}

func TestDealAccountTxNFTNum(t *testing.T) {
	handler := &AccountHandler{
		TxList: map[string]*db.Transaction{
			"tx1": {UserAddr: "addr1"},
			"tx2": {UserAddr: "addr2"},
			"tx3": {UserAddr: "addr1"},
		},
		TransferEvents: []*db.ContractEventData{
			{
				Topic: "transfer",
				EventData: &db.TransferTopicEventData{
					TokenId:     "token1",
					FromAddress: "addr1",
					ToAddress:   "addr2",
				},
			},
			{
				Topic: "not_in_topic_map",
				EventData: &db.TransferTopicEventData{
					TokenId:     "token2",
					FromAddress: "addr2",
					ToAddress:   "addr3",
				},
			},
			{
				Topic:     "transfer",
				EventData: nil,
			},
			{
				Topic: "transfer",
				EventData: &db.TransferTopicEventData{
					TokenId: "",
				},
			},
		},
	}
	txNum, nftNum := handler.DealAccountTxNFTNum()
	if txNum["addr1"] != 2 || txNum["addr2"] != 1 {
		t.Errorf("unexpected txNum: %+v", txNum)
	}
	if nftNum["addr1"] != -1 || nftNum["addr2"] != 1 {
		t.Errorf("unexpected nftNum: %+v", nftNum)
	}
}

func TestDealWithAccountData(t *testing.T) {
	handler := &AccountHandler{
		ChainId:   db.UTchainID,
		MinHeight: 100,
		TxList: map[string]*db.Transaction{
			"tx1": {UserAddr: "addr1"},
		},
		TransferEvents: []*db.ContractEventData{},
		EventResults: &model.TopicEventResult{
			OwnerAdders: []string{"addr1", "addr2"},
		},
		AccountMap: map[string]*db.Account{
			"addr1": {Address: "addr1"},
		},
		DelayGetDBResult: &model.GetDBResult{
			AccountDBMap: map[string]*db.Account{
				"addr1": {Address: "addr1"},
			},
		},
	}
	result := handler.DealWithAccountData()
	if len(result.InsertAccount) == 0 && len(result.UpdateAccount) == 0 {
		t.Errorf("expected insert or update accounts, got none")
	}
}

func TestBuildAccountInsertOrUpdate(t *testing.T) {
	accountMap := map[string]*db.Account{
		"addr1": {Address: "addr1", BlockHeight: 1},
	}
	delayGetDBResult := &model.GetDBResult{
		AccountDBMap:   accountMap,
		AccountBNSList: []*db.Account{},
		AccountDIDList: []*db.Account{},
	}
	eventResults := &model.TopicEventResult{
		OwnerAdders:      []string{"addr2", "addr3", ""},
		BNSBindEventData: nil,
		DIDEventData:     nil,
	}
	handler := &AccountHandler{
		ChainId:          db.UTchainID,
		MinHeight:        10,
		AccountMap:       accountMap,
		DelayGetDBResult: delayGetDBResult,
		EventResults:     eventResults,
	}
	accountTx := map[string]int64{"addr2": 5}
	accountNFT := map[string]int64{"addr2": 2}
	insert, _ := handler.BuildAccountInsertOrUpdate(accountTx, accountNFT)
	if len(insert) == 0 {
		t.Errorf("expected insert accounts, got none")
	}
	// 测试更新逻辑
	handler.AccountMap["addr2"] = &db.Account{Address: "addr2", BlockHeight: 5}
	accountTx = map[string]int64{"addr2": 10}
	accountNFT = map[string]int64{"addr2": 3}
	_, update := handler.BuildAccountInsertOrUpdate(accountTx, accountNFT)
	if len(update) == 0 {
		t.Errorf("expected update accounts, got none")
	}
}

func TestBuildAccountInsertOrUpdate_Empty(t *testing.T) {
	handler := &AccountHandler{
		ChainId:          db.UTchainID,
		MinHeight:        1,
		AccountMap:       map[string]*db.Account{},
		DelayGetDBResult: &model.GetDBResult{},
		EventResults:     &model.TopicEventResult{},
	}
	accountTx := map[string]int64{}
	accountNFT := map[string]int64{}
	insert, _ := handler.BuildAccountInsertOrUpdate(accountTx, accountNFT)
	if len(insert) != 0 {
		t.Errorf("expected no insert accounts, got insert: %d", len(insert))
	}
}

func TestDealAccountTxNFTNum_Empty(t *testing.T) {
	handler := &AccountHandler{
		TxList:         map[string]*db.Transaction{},
		TransferEvents: []*db.ContractEventData{},
	}
	txNum, nftNum := handler.DealAccountTxNFTNum()
	if len(txNum) != 0 || len(nftNum) != 0 {
		t.Errorf("expected empty txNum and nftNum, got txNum: %+v, nftNum: %+v", txNum, nftNum)
	}
}
