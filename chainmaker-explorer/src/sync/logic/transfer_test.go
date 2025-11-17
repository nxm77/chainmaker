package logic

import (
	"testing"
	"time"

	"chainmaker_web/src/db"
	"chainmaker_web/src/sync/common"
)

func TestDealTransferList(t *testing.T) {
	// 测试数据准备
	validTxId := "tx1231"
	validContractName := "contract1"
	validContractAddr := "0x123"
	validFromAddr := "0xfrom"
	validToAddr := "0xto"
	validTokenId := "token123"
	validAmount := "100"
	validDecimals := 18
	validTimestamp := time.Now().Unix()

	// 测试用例
	tests := []struct {
		name                 string
		eventDataList        []*db.ContractEventData
		contractInfoMap      map[string]*db.Contract
		txInfoMap            map[string]*db.Transaction
		expectFungibleLen    int
		expectNonFungibleLen int
	}{
		{
			name:                 "Empty input",
			eventDataList:        []*db.ContractEventData{},
			contractInfoMap:      map[string]*db.Contract{},
			txInfoMap:            map[string]*db.Transaction{},
			expectFungibleLen:    0,
			expectNonFungibleLen: 0,
		},
		{
			name: "Valid fungible contract",
			eventDataList: []*db.ContractEventData{{
				Topic:        common.TopicEventDataKey["Transfer"],
				ContractName: validContractName,
				EventData: &db.TransferTopicEventData{
					FromAddress: validFromAddr,
					ToAddress:   validToAddr,
					Amount:      validAmount,
				},
				TxId:      validTxId,
				Index:     0,
				Timestamp: validTimestamp,
			}},
			contractInfoMap: map[string]*db.Contract{
				validContractName: {
					Name:         validContractName,
					Addr:         validContractAddr,
					ContractType: common.ContractStandardNameCMDFA,
					Decimals:     validDecimals,
				},
			},
			txInfoMap: map[string]*db.Transaction{
				validTxId: {
					ContractMethod: "transfer",
				},
			},
			expectFungibleLen:    1,
			expectNonFungibleLen: 0,
		},
		{
			name: "Valid non-fungible contract",
			eventDataList: []*db.ContractEventData{{
				Topic:        common.TopicEventDataKey["Transfer"],
				ContractName: validContractName,
				EventData: &db.TransferTopicEventData{
					FromAddress: validFromAddr,
					ToAddress:   validToAddr,
					TokenId:     validTokenId,
				},
				TxId:      validTxId,
				Index:     0,
				Timestamp: validTimestamp,
			}},
			contractInfoMap: map[string]*db.Contract{
				validContractName: {
					Name:         validContractName,
					Addr:         validContractAddr,
					ContractType: common.ContractStandardNameCMNFA,
				},
			},
			txInfoMap: map[string]*db.Transaction{
				validTxId: {
					ContractMethod: "transferNFT",
				},
			},
			expectFungibleLen:    0,
			expectNonFungibleLen: 1,
		},
		{
			name: "Skip invalid topic",
			eventDataList: []*db.ContractEventData{{
				Topic:        "invalid_topic",
				ContractName: validContractName,
				EventData: &db.TransferTopicEventData{
					FromAddress: validFromAddr,
					ToAddress:   validToAddr,
					Amount:      validAmount,
				},
				TxId:      validTxId,
				Index:     0,
				Timestamp: validTimestamp,
			}},
			contractInfoMap: map[string]*db.Contract{
				validContractName: {
					Name:         validContractName,
					Addr:         validContractAddr,
					ContractType: common.ContractStandardNameCMDFA,
					Decimals:     validDecimals,
				},
			},
			txInfoMap: map[string]*db.Transaction{
				validTxId: {
					ContractMethod: "transfer",
				},
			},
			expectFungibleLen:    0,
			expectNonFungibleLen: 0,
		},
		{
			name: "Skip missing contract",
			eventDataList: []*db.ContractEventData{{
				Topic:        common.TopicEventDataKey["Transfer"],
				ContractName: "missing_contract",
				EventData: &db.TransferTopicEventData{
					FromAddress: validFromAddr,
					ToAddress:   validToAddr,
					Amount:      validAmount,
				},
				TxId:      validTxId,
				Index:     0,
				Timestamp: validTimestamp,
			}},
			contractInfoMap: map[string]*db.Contract{},
			txInfoMap: map[string]*db.Transaction{
				validTxId: {
					ContractMethod: "transfer",
				},
			},
			expectFungibleLen:    0,
			expectNonFungibleLen: 0,
		},
		{
			name: "Skip missing event data",
			eventDataList: []*db.ContractEventData{{
				Topic:        common.TopicEventDataKey["Transfer"],
				ContractName: validContractName,
				EventData:    nil,
				TxId:         validTxId,
				Index:        0,
				Timestamp:    validTimestamp,
			}},
			contractInfoMap: map[string]*db.Contract{
				validContractName: {
					Name:         validContractName,
					Addr:         validContractAddr,
					ContractType: common.ContractStandardNameCMDFA,
					Decimals:     validDecimals,
				},
			},
			txInfoMap: map[string]*db.Transaction{
				validTxId: {
					ContractMethod: "transfer",
				},
			},
			expectFungibleLen:    0,
			expectNonFungibleLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _ = DealTransferList(tt.eventDataList, tt.contractInfoMap, tt.txInfoMap)
		})
	}
}
