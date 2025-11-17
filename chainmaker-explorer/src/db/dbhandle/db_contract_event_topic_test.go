/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"chainmaker_web/src/db"
	"testing"

	"chainmaker.org/chainmaker/common/v2/random/uuid"
)

func TestInsertContractEventTopic(t *testing.T) {
	// Test case 1: Insert a single event topic
	inserts1 := []*db.ContractEventTopic{
		{
			ID:           uuid.GetUUID(),
			ContractName: "contract1",
			Topic:        "topic1",
			TxNum:        1,
			BlockHeight:  100,
		},
	}
	err1 := InsertContractEventTopic(ChainID, inserts1)
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}

	// Test case 2: Insert multiple event topics
	inserts2 := []*db.ContractEventTopic{
		{
			ID:           uuid.GetUUID(),
			ContractName: "contract2",
			Topic:        "topic2",
			TxNum:        2,
			BlockHeight:  200,
		},
		{
			ID:           uuid.GetUUID(),
			ContractName: "contract3",
			Topic:        "topic3",
			TxNum:        3,
			BlockHeight:  300,
		},
	}
	err2 := InsertContractEventTopic(ChainID, inserts2)
	if err2 != nil {
		t.Errorf("Test case 2 failed: %v", err2)
	}

	// Test case 3: Insert no event topics
	inserts3 := []*db.ContractEventTopic{}
	err3 := InsertContractEventTopic(ChainID, inserts3)
	if err3 != nil {
		t.Errorf("Test case 3 failed: %v", err3)
	}
}

func TestGetEventTopicByNames(t *testing.T) {
	// Test case 4: Get event topics by names
	contractNames4 := []string{"contract1", "contract2"}
	_, err4 := GetEventTopicByNames(ChainID, contractNames4)
	if err4 != nil {
		t.Errorf("Test case 4 failed: %v", err4)
	}

	// Test case 5: Get event topics by names with no matching topics
	contractNames5 := []string{"contractX", "contractY"}
	_, err5 := GetEventTopicByNames(ChainID, contractNames5)
	if err5 != nil {
		t.Errorf("Test case 5 failed: %v", err5)
	}
}

func TestUpdateContractEventTopic(t *testing.T) {
	// Test case 7: Update event topic
	update7 := &db.ContractEventTopic{
		ID:           uuid.GetUUID(),
		ContractName: "contract1",
		Topic:        "topic1",
		TxNum:        10,
		BlockHeight:  1000,
	}
	err7 := UpdateContractEventTopic(ChainID, update7)
	if err7 != nil {
		t.Errorf("Test case 7 failed: %v", err7)
	}
}

func TestGetContractEventTopic(t *testing.T) {
	// Test case 10: Get event topic by contract name and topic
	topic10 := "topic1"
	_, err10 := GetContractEventTopic(ChainID, ContractName1, topic10)
	if err10 != nil {
		t.Errorf("Test case 10 failed: %v", err10)
	}
}
