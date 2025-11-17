/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package logic

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/sync/common"
	"chainmaker_web/src/sync/model"
	"testing"

	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	"github.com/test-go/testify/assert"
)

func TestProcessBlackListEvent(t *testing.T) {
	topicEventResult := &model.TopicEventResult{}
	contractName := syscontract.SystemContract_TRANSACTION_MANAGER.String()
	topicAdd := common.TopicTxAddBlack
	topicDelete := common.TopicTxDeleteBlack
	eventDataAdd := []string{"chain1", "addr1", "addr2"}
	eventDataDelete := []string{"chain1", "addr3", "addr4"}

	processBlackListEvent(topicEventResult, contractName, topicAdd, eventDataAdd)
	assert.Equal(t, []string{"addr1", "addr2"}, topicEventResult.AddBlack)
	assert.Empty(t, topicEventResult.DeleteBlack)

	processBlackListEvent(topicEventResult, contractName, topicDelete, eventDataDelete)
	assert.Equal(t, []string{"addr1", "addr2"}, topicEventResult.AddBlack)
	assert.Equal(t, []string{"addr3", "addr4"}, topicEventResult.DeleteBlack)
}

func TestGetTransferEventDta(t *testing.T) {
	contractType := common.ContractStandardNameCMDFA
	topic := "transfer"
	senderUser := "sender1"
	eventData := []string{"fromAddr", "toAddr", "100"}

	topicEventData := GetTransferEventDta(contractType, topic, senderUser, eventData)
	assert.NotNil(t, topicEventData)
	assert.Equal(t, "fromAddr", topicEventData.FromAddress)
	assert.Equal(t, "toAddr", topicEventData.ToAddress)
	assert.Equal(t, "100", topicEventData.Amount)
}

func TestDealDockerDFAEventData1(t *testing.T) {
	tests := []struct {
		name      string
		topic     string
		eventData []string
		expected  *db.TransferTopicEventData
	}{
		{
			name:      "Mint",
			topic:     "mint",
			eventData: []string{"toAddr", "100"},
			expected: &db.TransferTopicEventData{
				ToAddress: "toAddr",
				Amount:    "100",
			},
		},
		{
			name:      "Transfer",
			topic:     "transfer",
			eventData: []string{"fromAddr", "toAddr", "100"},
			expected: &db.TransferTopicEventData{
				FromAddress: "fromAddr",
				ToAddress:   "toAddr",
				Amount:      "100",
			},
		},
		{
			name:      "Burn",
			topic:     "burn",
			eventData: []string{"fromAddr", "100"},
			expected: &db.TransferTopicEventData{
				FromAddress: "fromAddr",
				Amount:      "100",
			},
		},
		{
			name:      "Approve",
			topic:     "approve",
			eventData: []string{"fromAddr", "toAddr", "100"},
			expected: &db.TransferTopicEventData{
				FromAddress: "fromAddr",
				ToAddress:   "toAddr",
				Amount:      "100",
			},
		},
		{
			name:      "InvalidEventData",
			topic:     "transfer",
			eventData: []string{"fromAddr"},
			expected:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			topicEventData := DealDockerDFAEventData(tt.topic, tt.eventData)
			assert.Equal(t, tt.expected, topicEventData)
		})
	}
}

func TestDealDockerNFAEventData1(t *testing.T) {
	tests := []struct {
		name      string
		topic     string
		eventData []string
		sender    string
		expected  *db.TransferTopicEventData
	}{
		{
			name:      "Mint",
			topic:     "Mint",
			eventData: []string{"", "toAddr", "tokenId", "category", "metadata"},
			sender:    "",
			expected: &db.TransferTopicEventData{
				FromAddress:  "",
				ToAddress:    "toAddr",
				TokenId:      "tokenId",
				CategoryName: "category",
				Metadata:     "metadata",
			},
		},
		{
			name:      "TransferFrom",
			topic:     "TransferFrom",
			eventData: []string{"fromAddr", "toAddr", "tokenId"},
			sender:    "",
			expected: &db.TransferTopicEventData{
				FromAddress: "fromAddr",
				ToAddress:   "toAddr",
				TokenId:     "tokenId",
			},
		},
		{
			name:      "Burn",
			topic:     "Burn",
			eventData: []string{"tokenId"},
			sender:    "sender",
			expected: &db.TransferTopicEventData{
				FromAddress: "sender",
				TokenId:     "tokenId",
			},
		},
		{
			name:      "SetApproval",
			topic:     "SetApproval",
			eventData: []string{"fromAddr", "toAddr", "tokenId", "approval"},
			sender:    "",
			expected: &db.TransferTopicEventData{
				FromAddress: "fromAddr",
				ToAddress:   "toAddr",
				TokenId:     "tokenId",
				Approval:    "approval",
			},
		},
		{
			name:      "SetApprovalForAll",
			topic:     "SetApprovalForAll",
			eventData: []string{"fromAddr", "toAddr", "approval"},
			sender:    "",
			expected: &db.TransferTopicEventData{
				FromAddress: "fromAddr",
				ToAddress:   "toAddr",
				Approval:    "approval",
			},
		},
		{
			name:      "SetApprovalByCategory",
			topic:     "SetApprovalByCategory",
			eventData: []string{"fromAddr", "toAddr", "category", "approval"},
			sender:    "",
			expected: &db.TransferTopicEventData{
				FromAddress:  "fromAddr",
				ToAddress:    "toAddr",
				CategoryName: "category",
				Approval:     "approval",
			},
		},
		{
			name:      "InvalidEventData",
			topic:     "Mint",
			eventData: []string{"toAddr"},
			sender:    "",
			expected:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			topicEventData := DealDockerNFAEventData(tt.topic, tt.eventData, tt.sender)
			assert.Equal(t, tt.expected, topicEventData)
		})
	}
}

func TestDealEVMDFAEventData1(t *testing.T) {
	tests := []struct {
		name      string
		topic     string
		eventData []string
		expected  *db.TransferTopicEventData
	}{
		{
			name:      "Transfer",
			topic:     common.EVMEventTopicTransfer,
			eventData: []string{"0x000000000000000000000000fromAddr", "0x000000000000000000000000toAddr", "64"},
			expected: &db.TransferTopicEventData{
				FromAddress: "fromAddr",
				ToAddress:   "toAddr",
				Amount:      "100",
			},
		},
		{
			name:      "InvalidEventData",
			topic:     common.EVMEventTopicTransfer,
			eventData: []string{"0x000000000000000000000000fromAddr"},
			expected:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = DealEVMDFAEventData(tt.topic, tt.eventData)
		})
	}
}

func TestDealEVMNFAEventData1(t *testing.T) {
	tests := []struct {
		name      string
		topic     string
		eventData []string
		expected  *db.TransferTopicEventData
	}{
		{
			name:      "Transfer",
			topic:     common.EVMEventTopicTransfer,
			eventData: []string{"0x000000000000000000000000fromAddr", "0x000000000000000000000000toAddr", "64"},
			expected: &db.TransferTopicEventData{
				FromAddress: "fromAddr",
				ToAddress:   "toAddr",
				TokenId:     "100",
			},
		},
		{
			name:      "InvalidEventData",
			topic:     common.EVMEventTopicTransfer,
			eventData: []string{"0x000000000000000000000000fromAddr"},
			expected:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = DealEVMNFAEventData(tt.topic, tt.eventData)
		})
	}
}
