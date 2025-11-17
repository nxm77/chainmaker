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

	"chainmaker.org/chainmaker/contract-utils/standard"
)

func TestDealTopicDeleteBlockList(t *testing.T) {
	// 创建一个TopicEventResult对象
	topicEventResult := newTopicEventResult()
	//IDA数据
	topicEventResult.IDAEventData = newIDAEventResult()
	//DID数据
	topicEventResult.DIDEventData = model.NewDIDEventData()
	type args struct {
		didEventData *model.DIDEventData
		eventData    []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test",
			args: args{
				didEventData: topicEventResult.DIDEventData,
				eventData:    []string{"[\"did:cnbn:f479f90b0665ae89ecf5e23ddf3738b83aae2222\"]"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			DealTopicDeleteBlockList(tt.args.didEventData, tt.args.eventData)
		})
	}
}
func TestProcessDIDEvent(t *testing.T) {
	chainId := "chain1"
	event := &db.ContractEvent{
		ContractType:    standard.ContractStandardNameCMDID,
		ContractNameBak: "DID",
		ContractAddr:    "addr1",
		TxId:            "tx1",
		Timestamp:       1234567890,
	}
	eventData := []string{"did:cnbn:f479f90b0665ae89ecf5e23ddf3738b83aae2222", `{"id":"did:cnbn:f479f90b0665ae89ecf5e23ddf3738b83aae2222"}`}
	didEventData := model.NewDIDEventData()

	tests := []struct {
		name      string
		topic     string
		eventData []string
	}{
		{
			name:      "SetDidDocument",
			topic:     standard.Topic_SetDidDocument,
			eventData: eventData,
		},
		{
			name:      "AddBlackList",
			topic:     standard.Topic_AddBlackList,
			eventData: []string{`["did:cnbn:f479f90b0665ae89ecf5e23ddf3738b83aae2222"]`},
		},
		{
			name:      "DeleteBlackList",
			topic:     standard.Topic_DeleteBlackList,
			eventData: []string{`["did:cnbn:f479f90b0665ae89ecf5e23ddf3738b83aae2222"]`},
		},
		{
			name:      "AddTrustIssuer",
			topic:     standard.Topic_AddTrustIssuer,
			eventData: []string{`["did:cnbn:f479f90b0665ae89ecf5e23ddf3738b83aae2222"]`},
		},
		{
			name:      "DeleteTrustIssuer",
			topic:     standard.Topic_DeleteTrustIssuer,
			eventData: []string{`["did:cnbn:f479f90b0665ae89ecf5e23ddf3738b83aae2222"]`},
		},
		{
			name:      "SetVcTemplate",
			topic:     standard.Topic_SetVcTemplate,
			eventData: []string{"templateID", "templateName", "vcType", "version", "template"},
		},
		{
			name:      "VcIssueLog",
			topic:     standard.Topic_VcIssueLog,
			eventData: []string{"issuerDID", "holderDID", "templateID", "vcID"},
		},
		{
			name:      "RevokeVc",
			topic:     standard.Topic_RevokeVc,
			eventData: []string{"vcID"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event.Topic = tt.topic
			processDIDEvent(chainId, event, tt.eventData, didEventData)
		})
	}
}
