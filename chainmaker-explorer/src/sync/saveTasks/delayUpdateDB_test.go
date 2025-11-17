/*
Package saveTasks comment： resolver delay update DB
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package saveTasks

import (
	"chainmaker_web/src/db"
	"testing"

	"github.com/test-go/testify/assert"
)

func TestInsertFungibleTransferToDB(t *testing.T) {
	transferList := []*db.FungibleTransfer{
		{ContractAddr: "transfer1"},
	}

	err := InsertFungibleTransferToDB(db.UTchainID, transferList)
	assert.Nil(t, err)
}

func TestInsertNonFungibleTransferToDB(t *testing.T) {
	transferList := []*db.NonFungibleTransfer{
		{ContractAddr: "transfer1"},
	}

	err := InsertNonFungibleTransferToDB(db.UTchainID, transferList)
	assert.Nil(t, err)
}

func TestSaveNonFungibleToken(t *testing.T) {
	tokenResult := &db.TokenResult{
		InsertUpdateToken: []*db.NonFungibleToken{
			{TokenId: "token1", ContractAddr: "addr1"},
		},
		DeleteToken: []*db.NonFungibleToken{
			{TokenId: "token1", ContractAddr: "addr1"},
		},
	}

	err := SaveNonFungibleToken(db.UTchainID, tokenResult)
	assert.Nil(t, err)
}

func TestInsertNonFungibleTokenConcurrent(t *testing.T) {
	tokenList := []*db.NonFungibleToken{
		{TokenId: "token1"},
	}

	err := InsertNonFungibleTokenConcurrent(db.UTchainID, tokenList)
	assert.Nil(t, err)
}

func TestUpdateNonFungibleToken(t *testing.T) {
	tokenList := []*db.NonFungibleToken{
		{TokenId: "token1"},
	}

	err := UpdateNonFungibleToken(db.UTchainID, tokenList)
	assert.Nil(t, err)
}

func TestSavePositionToDB(t *testing.T) {
	dbPositionOperates := &db.BlockPosition{
		InsertFungiblePosition: []*db.FungiblePosition{
			{ContractAddr: "position1"},
		},
		DeleteFungiblePosition: []*db.FungiblePosition{
			{ContractAddr: "position1"},
		},
		UpdateFungiblePosition: []*db.FungiblePosition{
			{ContractAddr: "position1"},
		},
		InsertNonFungible: []*db.NonFungiblePosition{
			{ContractAddr: "position1"},
		},
		DeleteNonFungible: []*db.NonFungiblePosition{
			{ContractAddr: "position1"},
		},
		UpdateNonFungible: []*db.NonFungiblePosition{
			{ContractAddr: "position1"},
		},
	}

	err := SavePositionToDB(db.UTchainID, dbPositionOperates)
	assert.Nil(t, err)
}

func TestInsertGasToDB(t *testing.T) {
	insertGas := []*db.Gas{
		{Address: "gas1"},
	}

	err := InsertGasToDB(db.UTchainID, insertGas)
	assert.Nil(t, err)
}

func TestUpdateGasToDB(t *testing.T) {
	updateGas := []*db.Gas{
		{Address: "gas1"},
	}

	err := UpdateGasToDB(db.UTchainID, updateGas)
	assert.Nil(t, err)
}

func TestUpdateTxBlackToDB(t *testing.T) {
	txBlockList := &db.UpdateTxBlack{
		AddTxBlack: []*db.BlackTransaction{
			{TxId: "tx1"},
		},
		DeleteTxBlack: []*db.Transaction{
			{TxId: "tx1"},
		},
	}

	err := UpdateTxBlackToDB(db.UTchainID, txBlockList)
	assert.Nil(t, err)
}

func TestUpdateContract(t *testing.T) {
	identityContract := []*db.IdentityContract{
		{ContractName: "contract1"},
	}

	err := UpdateContract(db.UTchainID, identityContract)
	assert.Nil(t, err)
}

func TestUpdateContractTxNum(t *testing.T) {
	updateContracts := []*db.Contract{
		{Addr: "contract1"},
	}

	err := UpdateContractTxNum(db.UTchainID, updateContracts)
	assert.Nil(t, err)
}

func TestSaveFungibleContractResult(t *testing.T) {
	contractResult := &db.GetContractResult{
		UpdateFungibleContract: []*db.FungibleContract{
			{ContractAddr: "contract1"},
		},
		UpdateNonFungible: []*db.NonFungibleContract{
			{ContractAddr: "contract2"},
		},
	}

	err := SaveFungibleContractResult(db.UTchainID, contractResult)
	assert.Nil(t, err)
}

func TestUpdateBlockStatusToDB(t *testing.T) {
	blockHeightList := []int64{1, 2, 3}

	err := UpdateBlockStatusToDB(db.UTchainID, blockHeightList)
	assert.Nil(t, err)
}

func TestUpdateIDAContract(t *testing.T) {
	updateIdaContract := map[string]*db.IDAContract{
		"addr1": {ContractAddr: "contract1"},
	}

	err := UpdateIDAContract(db.UTchainID, updateIdaContract)
	assert.Nil(t, err)
}

func TestSaveIDAAssetDataToDB(t *testing.T) {
	idaAssetsData := &db.IDAAssetsDataDB{
		IDAAssetDetail: []*db.IDAAssetDetail{
			{ID: "asset1"},
		},
	}

	err := SaveIDAAssetDataToDB(db.UTchainID, idaAssetsData)
	assert.Nil(t, err)
}

func TestUpdateIDAAssetDataToDB(t *testing.T) {
	updateAssetsData := &db.IDAAssetsUpdateDB{
		UpdateAssetDetails: []*db.IDAAssetDetail{
			{ID: "asset1"},
		},
	}

	err := UpdateIDAAssetDataToDB(db.UTchainID, updateAssetsData)
	assert.Nil(t, err)
}

func TestSaveEventTopicData(t *testing.T) {
	inserts := []*db.ContractEventTopic{
		{ID: "event1"},
	}
	updates := []*db.ContractEventTopic{
		{ID: "event2"},
	}

	err := SaveEventTopicData(db.UTchainID, inserts, updates)
	assert.Nil(t, err)
}

func TestUpdateDIDDataToDB(t *testing.T) {
	didSaveDate := &db.DIDSaveData{
		SaveDIDDetails: []*db.DIDDetail{
			{DID: "did1"},
		},
	}

	err := UpdateDIDDataToDB(db.UTchainID, didSaveDate)
	assert.Nil(t, err)
}

func TestSaveChainStatistics(t *testing.T) {
	statistics := &db.Statistics{
		ChainId: db.UTchainID,
	}

	err := SaveChainStatistics(db.UTchainID, statistics)
	assert.Nil(t, err)
}

func TestSaveAccountToDB(t *testing.T) {
	accountResult := &db.UpdateAccountResult{
		InsertAccount: []*db.Account{
			{Address: "account1"},
		},
		UpdateAccount: []*db.Account{
			{Address: "account2"},
		},
	}

	err := SaveAccountToDB(db.UTchainID, accountResult)
	assert.Nil(t, err)
}
