/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package logic

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var UTTestABIJson = `[{"anonymous":false,"inputs":[{"indexed":true,"internalType":"string","name":"address","type":"string"},{"indexed":true,"internalType":"string","name":"address2","type":"string"},{"indexed":false,"internalType":"string","name":"value","type":"string"},{"indexed":true,"internalType":"string","name":"id1","type":"string"}],"name":"approve","describe":"跨境转账","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"string","name":"address","type":"string"},{"indexed":false,"internalType":"string","name":"count","type":"string"}],"name":"mint","describe":"更新白名单数据","type":"event"}]`
var UTContractAddr = "0x1234567890abcdef1234567890abcdef12345678"
var UTContractName = "TestContract"
var UTContractVersion = "v1"
var UTContractTopic = "approve"

func TestGetTopicABIInputs(t *testing.T) {
	uuid := uuid.New().String()
	insert := &db.ContractABIFile{
		Id:              uuid,
		ContractName:    UTContractName,
		ContractAddr:    UTContractAddr,
		ContractVersion: UTContractVersion,
		ABIJson:         UTTestABIJson,
	}

	err := dbhandle.InsertContractABIFile(db.UTchainID, insert)
	assert.NoError(t, err)

	// 正常情况
	fields := GetTopicABIInputs(db.UTchainID, UTContractAddr, UTContractVersion, UTContractTopic)
	assert.Equal(t, 4, len(fields))
	assert.Equal(t, "address", fields[0].Name)
	assert.Equal(t, "string", fields[0].Type)

	// ABI不存在
	fields = GetTopicABIInputs("c1", "addr1", "v1", "topic1")
	assert.Equal(t, 0, len(fields))
}

func TestBuildEventDataByABI(t *testing.T) {
	insert := &db.ContractABITopic{
		Id:              uuid.New().String(),
		ContractName:    UTContractName,
		ContractAddr:    UTContractAddr,
		ContractVersion: UTContractVersion,
		Topic:           UTContractTopic,
		TopicTableName:  "topictablename",
	}
	err := dbhandle.InsertContractABITopic(db.UTchainID, insert)
	assert.NoError(t, err)

	events := []*db.ContractEvent{
		{
			ContractAddr:    UTContractAddr,
			Topic:           UTContractTopic,
			ContractVersion: UTContractVersion,
			TxId:            "tx1",
			EventIndex:      1,
			EventData:       `["a","b","a","b"]`,
		},
	}
	result := BuildEventDataByABI(db.UTchainID, events)
	assert.Equal(t, 1, len(result))

	// 空events
	result = BuildEventDataByABI(db.UTchainID, []*db.ContractEvent{})
	assert.Equal(t, 0, len(result))
}

func TestBuildEventDataByABIInputs(t *testing.T) {
	event := &db.ContractEvent{
		TxId:            "txid",
		EventIndex:      2,
		ContractAddr:    UTContractAddr,
		ContractVersion: UTContractVersion,
		Topic:           UTContractTopic,
		Timestamp:       123456,
		EventData:       `["val1","123","val1","123"]`,
	}
	_, err := BuildEventDataByABIInputs(db.UTchainID, event)
	assert.NoError(t, err)

	// eventData数量不匹配
	event.EventData = `["val1"]`
	_, err = BuildEventDataByABIInputs(db.UTchainID, event)
	assert.Error(t, err)

	// eventData不是json数组
	event.EventData = `invalid json`
	_, err = BuildEventDataByABIInputs(db.UTchainID, event)
	assert.Error(t, err)
}
