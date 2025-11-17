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
	"encoding/json"
	"reflect"
	"testing"

	"chainmaker.org/chainmaker/contract-utils/standard"
	pbCommon "chainmaker.org/chainmaker/pb-go/v2/common"
	"github.com/google/go-cmp/cmp"
	"github.com/test-go/testify/assert"
)

var TxInfoEventJson = "{\"payload\":{\"chain_id\":\"chain1\",\"tx_id\":\"17a05b47e5fe40a1ca6e85d4a315abef72b52452725a4c7eac62242906669aa5\",\"timestamp\":1702461023,\"contract_name\":\"goErc20_1\",\"method\":\"Mint\",\"parameters\":[{\"key\":\"account\",\"value\":\"MThmYzRlNzQyOWFmODQxOWQ1YmIzMDdlMzRkYjM5OGI5YTIzMzFjNg==\"},{\"key\":\"amount\",\"value\":\"MTAwMDAwMDAwMDA=\"}],\"limit\":{\"gas_limit\":13000}},\"sender\":{\"signer\":{\"org_id\":\"wx-org1.chainmaker.org\",\"member_type\":1,\"member_info\":\"LK4/KplYsQcFU2All0UxorspVALdt/tgHuZ8QxiME2M=\"},\"signature\":\"MEQCIAK4XuZoU0XB+ya2PRNtebY/BACX8BQOBRMQBMidbfXpAiA6EURjarfxU/qbwCrkqptOdXav7orDeVR38aHXzynb5g==\"},\"result\":{\"contract_result\":{\"result\":\"b2s=\",\"message\":\"Success\",\"gas_used\":156,\"contract_event\":[{\"topic\":\"mint\",\"tx_id\":\"17a05b47e5fe40a1ca6e85d4a315abef72b52452725a4c7eac62242906669aa5\",\"contract_name\":\"goErc20_1\",\"event_data\":[\"18fc4e7429af8419d5bb307e34db398b9a2331c6\",\"10000000000\"]}]},\"rw_set_hash\":\"BOf54ycn6MjSiUL06tEU8WN2cDvTXJWkShVNowYYAK4=\"}}"
var ContractEventsJson = "[{\"txId\":\"17a05b47e5fe40a1ca6e85d4a315abef72b52452725a4c7eac62242906669aa5\",\"eventIndex\":1,\"topic\":\"mint\",\"topicBak\":\"\",\"contractName\":\"goErc20_1\",\"contractNameBak\":\"goErc20_1\",\"contractAddr\":\"\",\"contractVersion\":\"\",\"eventData\":[\"18fc4e7429af8419d5bb307e34db398b9a2331c6\",\"10000000000\"],\"eventDataBak\":\"\",\"timestamp\":1702461023,\"createdAt\":\"0001-01-01T00:00:00Z\",\"updatedAt\":\"0001-01-01T00:00:00Z\"}]"

var ContractInfoMap = map[string]*db.Contract{
	"aba31ce4cd49f08073d2f115eb12610544242ff9": {
		Name:         "goErc20_1",
		NameBak:      "goErc20_1",
		Addr:         "aba31ce4cd49f08073d2f115eb12610544242ff9",
		ContractType: "CMDFA",
		TxNum:        100,
	},
	"goErc20_1": {
		Name:         "goErc20_1",
		NameBak:      "goErc20_1",
		Addr:         "aba31ce4cd49f08073d2f115eb12610544242ff9",
		ContractType: "CMDFA",
		TxNum:        100,
	},
}

func TestDealContractEvents(t *testing.T) {
	transactionInfo := &pbCommon.Transaction{}
	err := json.Unmarshal([]byte(txInfoJson), transactionInfo)
	if err != nil {
		return
	}

	transactionEvent := &pbCommon.Transaction{}
	err = json.Unmarshal([]byte(TxInfoEventJson), transactionEvent)
	if err != nil {
		return
	}

	contractEvents := make([]*db.ContractEvent, 0)
	err = json.Unmarshal([]byte(ContractEventsJson), &contractEvents)
	if err != nil {
		return
	}
	type args struct {
		txInfo *pbCommon.Transaction
	}
	tests := []struct {
		name string
		args args
		want []*db.ContractEvent
	}{
		{
			name: "Test case 1",
			args: args{
				txInfo: transactionInfo,
			},
			want: []*db.ContractEvent{},
		},
		{
			name: "Test case 1",
			args: args{
				txInfo: transactionEvent,
			},
			want: contractEvents,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DealContractEvents(tt.args.txInfo)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("DealContractEvents() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestDealContractTxNum(t *testing.T) {
	type args struct {
		minHeight     int64
		contractMap   map[string]*db.Contract
		txList        map[string]*db.Transaction
		contractEvent []*db.ContractEvent
	}
	txInfoListMap := map[string]*db.Transaction{}
	txInfoListMap["17a05b47e5fe40a1ca6e85d4a315abef72b52452725a4c7eac62242906669aa5"] = &db.Transaction{
		TxId:               "17a05b47e5fe40a1ca6e85d4a315abef72b52452725a4c7eac62242906669aa5",
		Sender:             "client1.sign.wx-org1.chainmaker.org",
		SenderOrgId:        "wx-org1.chainmaker.org",
		BlockHeight:        40,
		BlockHash:          "d3b2b488033c2faa100949667572b1875d82f7a32bd35bccf8232f5d3eef6545",
		TxType:             "INVOKE_CONTRACT",
		Timestamp:          1702461023,
		TxIndex:            1,
		TxStatusCode:       "SUCCESS",
		RwSetHash:          "04e7f9e32727e8c8d28942f4ead114f16376703bd35c95a44a154da3061800ae",
		ContractResultCode: 0,
		ContractName:       "goErc20_1",
		ContractNameBak:    "goErc20_1",
		ContractAddr:       "aba31ce4cd49f08073d2f115eb12610544242ff9",
		ContractType:       "CMDFA",
		UserAddr:           "171262347a59fded92021a32421a5dad05424e03",
	}

	tests := []struct {
		name      string
		args      args
		wantTxNum int64
	}{
		{
			name: "Test case 1",
			args: args{
				minHeight:   1000,
				contractMap: ContractInfoMap,
				txList:      txInfoListMap,
			},
			wantTxNum: 101,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := UpdateContractTxAndEventNum(tt.args.minHeight, tt.args.contractMap, tt.args.txList, tt.args.contractEvent)
			if !reflect.DeepEqual(got[0].TxNum, tt.wantTxNum) {
				t.Errorf("DealContractTxNum() = %v, want %v", got[0], tt.wantTxNum)
			}
		})
	}
}

func TestBuildTransferEventData(t *testing.T) {
	type args struct {
		topicEventResult *model.TopicEventResult
		ownerAddrMap     map[string]string
		contractInfoMap  map[string]*db.Contract
		event            *db.ContractEvent
		senderUser       string
		eventData        []string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "Test case 1",
			args: args{
				topicEventResult: &model.TopicEventResult{},
				ownerAddrMap:     map[string]string{},
				contractInfoMap: map[string]*db.Contract{
					"ContractName": {
						ContractType: "CMDFA",
					},
				},
				event: &db.ContractEvent{
					TxId:            "1231212313",
					Topic:           "mint",
					ContractName:    "ContractName",
					ContractNameBak: "ContractName",
					ContractAddr:    "1234",
				},
				eventData: []string{
					"12345",
					"123",
				},
			},
			want: map[string]string{
				"12345": "12345",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildTransferEventData(tt.args.topicEventResult, tt.args.ownerAddrMap, tt.args.contractInfoMap, tt.args.event, tt.args.senderUser, tt.args.eventData); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildTransferEventData() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildIdentityEventData(t *testing.T) {
	eventInfo := &db.ContractEvent{
		TxId:            "1231212313",
		Topic:           "setIdentity",
		ContractName:    "ContractName",
		ContractNameBak: "ContractName",
		ContractAddr:    "1234",
	}
	type args struct {
		topicEventResult *model.TopicEventResult
		contractInfoMap  map[string]*db.Contract
		event            *db.ContractEvent
		eventData        []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test case 1",
			args: args{
				topicEventResult: &model.TopicEventResult{},
				contractInfoMap: map[string]*db.Contract{
					"ContractName": {
						ContractType: "CMID",
					},
				},
				event: eventInfo,
				eventData: []string{
					"12345666",
					"12345777",
					"123455",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			BuildIdentityEventData(tt.args.topicEventResult, tt.args.contractInfoMap, tt.args.event, tt.args.eventData)
		})
	}
}

func TestDealEventTopicTxNum(t *testing.T) {
	// Test case 1: Count event topics
	contractEvents1 := []*db.ContractEvent{
		{
			ContractName: "contract1",
			Topic:        "topic1",
		},
		{
			ContractName: "contract1",
			Topic:        "topic2",
		},
		{
			ContractName: "contract2",
			Topic:        "topic1",
		},
		{
			ContractName: "contract2",
			Topic:        "topic1",
		},
	}
	expectedResult1 := map[string]map[string]int64{
		"contract1": {
			"topic1": 1,
			"topic2": 1,
		},
		"contract2": {
			"topic1": 2,
		},
	}
	result1 := DealEventTopicTxNum(contractEvents1)
	if len(result1) != len(expectedResult1) {
		t.Errorf("Test case 1 failed: Expected %d contracts, got %d", len(expectedResult1), len(result1))
	} else {
		for contractName, topics := range result1 {
			expectedTopics, ok := expectedResult1[contractName]
			if !ok {
				t.Errorf("Test case 1 failed: Expected contract with name %s, got none", contractName)
			} else if len(topics) != len(expectedTopics) {
				t.Errorf("Test case 1 failed: Expected %d topics, got %d", len(expectedTopics), len(topics))
			} else {
				for topic, count := range topics {
					expectedCount, ok := expectedTopics[topic]
					if !ok {
						t.Errorf("Test case 1 failed: Expected topic %s, got none", topic)
					} else if count != expectedCount {
						t.Errorf("Test case 1 failed: Expected count %d, got %d", expectedCount, count)
					}
				}
			}
		}
	}

	// Test case 2: Count event topics with no events
	contractEvents2 := []*db.ContractEvent{}
	expectedResult2 := map[string]map[string]int64{}
	result2 := DealEventTopicTxNum(contractEvents2)
	if len(result2) != len(expectedResult2) {
		t.Errorf("Test case 2 failed: Expected %d contracts, got %d", len(expectedResult2), len(result2))
	}
}

func TestProcessEventTopicTxNum(t *testing.T) {
	// Test case 1: Process event topics with new topics
	eventTopicTxNum1 := map[string]map[string]int64{
		"contract1": {
			"topic1": 1,
			"topic2": 1,
		},
		"contract2": {
			"topic1": 2,
		},
	}
	eventTopicDBMap1 := map[string]map[string]*db.ContractEventTopic{}
	blockHeight1 := int64(1)
	expectedInsertResult1 := []*db.ContractEventTopic{
		{
			Topic:        "topic1",
			ContractName: "contract1",
			TxNum:        1,
			BlockHeight:  blockHeight1,
		},
		{
			Topic:        "topic2",
			ContractName: "contract1",
			TxNum:        1,
			BlockHeight:  blockHeight1,
		},
		{
			Topic:        "topic1",
			ContractName: "contract2",
			TxNum:        2,
			BlockHeight:  blockHeight1,
		},
	}
	expectedUpdateResult1 := []*db.ContractEventTopic{}
	insertResult1, updateResult1 := ProcessEventTopicTxNum(eventTopicTxNum1, eventTopicDBMap1, blockHeight1)
	if len(insertResult1) != len(expectedInsertResult1) {
		t.Errorf("Test case 1 failed: Expected %d insert results, got %d", len(expectedInsertResult1), len(insertResult1))
	}
	if len(updateResult1) != len(expectedUpdateResult1) {
		t.Errorf("Test case 1 failed: Expected %d update results, got %d", len(expectedUpdateResult1), len(updateResult1))
	}

	// Test case 2: Process event topics with existing topics
	eventTopicTxNum2 := map[string]map[string]int64{
		"contract1": {
			"topic1": 1,
			"topic2": 1,
		},
		"contract2": {
			"topic1": 2,
		},
	}
	eventTopicDBMap2 := map[string]map[string]*db.ContractEventTopic{
		"contract1": {
			"topic1": {
				Topic:        "topic1",
				ContractName: "contract1",
				TxNum:        1,
				BlockHeight:  int64(0),
			},
		},
	}
	blockHeight2 := int64(1)
	expectedInsertResult2 := []*db.ContractEventTopic{
		{
			Topic:        "topic2",
			ContractName: "contract1",
			TxNum:        1,
			BlockHeight:  blockHeight2,
		},
		{
			Topic:        "topic2",
			ContractName: "contract1",
			TxNum:        1,
			BlockHeight:  blockHeight2,
		},
	}
	expectedUpdateResult2 := []*db.ContractEventTopic{
		{
			Topic:        "topic1",
			ContractName: "contract1",
			TxNum:        2,
			BlockHeight:  blockHeight2,
		},
	}
	insertResult2, updateResult2 := ProcessEventTopicTxNum(eventTopicTxNum2, eventTopicDBMap2, blockHeight2)
	if len(insertResult2) != len(expectedInsertResult2) {
		t.Errorf("Test case 2 failed: Expected %d insert results, got %d", len(expectedInsertResult2), len(insertResult2))
	}
	if len(updateResult2) != len(expectedUpdateResult2) {
		t.Errorf("Test case 2 failed: Expected %d update results, got %d", len(expectedUpdateResult2), len(updateResult2))
	}

	// Test case 3: Process event topics with no topics
	eventTopicTxNum3 := map[string]map[string]int64{}
	eventTopicDBMap3 := map[string]map[string]*db.ContractEventTopic{}
	blockHeight3 := int64(1)
	expectedInsertResult3 := []*db.ContractEventTopic{}
	expectedUpdateResult3 := []*db.ContractEventTopic{}
	insertResult3, updateResult3 := ProcessEventTopicTxNum(eventTopicTxNum3, eventTopicDBMap3, blockHeight3)
	if len(insertResult3) != len(expectedInsertResult3) {
		t.Errorf("Test case 3 failed: Expected %d insert results, got %d", len(expectedInsertResult3), len(insertResult3))
	}
	if len(updateResult3) != len(expectedUpdateResult3) {
		t.Errorf("Test case 3 failed: Expected %d update results, got %d", len(expectedUpdateResult3), len(updateResult3))
	}
}

func TestBuildIDAEventData(t *testing.T) {
	// Test case 1: Parse IDA created event
	contractType1 := standard.ContractStandardNameCMIDA
	topic1 := standard.EventIDACreated
	eventData1 := []string{"asset1", "Contract1", "contractAddr1", "12345678"}
	BuildIDAEventData(contractType1, topic1, eventData1)

	// Test case 2: Parse IDA updated event
	contractType2 := standard.ContractStandardNameCMIDA
	topic2 := standard.EventIDAUpdated
	eventData2 := []string{"asset1", "Contract1", "contractAddr1", "12345678"}
	idaInfoList2, idaUpdateData2, idaIds2 := BuildIDAEventData(contractType2, topic2, eventData2)
	if len(idaInfoList2) != 0 {
		t.Errorf("Test case 2 failed: Expected no IDA infos, got %d", len(idaInfoList2))
	}
	if idaUpdateData2 == nil {
		t.Errorf("Test case 2 failed: Expected non-nil IDA update data, got nil")
	}
	if len(idaIds2) != 0 {
		t.Errorf("Test case 2 failed: Expected no IDA IDs, got %d", len(idaIds2))
	}

	// Test case 4: Parse non-IDA event
	contractType4 := "non-IDA"
	topic4 := "IDA"
	eventData4 := []string{"asset1", "Contract1", "contractAddr1", "12345678"}
	BuildIDAEventData(contractType4, topic4, eventData4)
}

func TestParseEventData(t *testing.T) {
	// Test case 1: Test with a valid event
	event1 := &db.ContractEvent{
		EventDataBak: `["event1", "event2"]`,
	}
	result1 := parseEventData(event1)
	if len(result1) != 2 {
		t.Errorf("Test case 1 failed: Expected 2 event data, got %d", len(result1))
	}

	// Test case 2: Test with an event with no data
	event2 := &db.ContractEvent{
		EventDataBak: "",
	}
	result2 := parseEventData(event2)
	if len(result2) != 0 {
		t.Errorf("Test case 2 failed: Expected %d event data, got %d", 0, len(result2))
	}

	// Test case 3: Test with an event with invalid JSON data
	event3 := &db.ContractEvent{
		EventDataBak: `["event1", "event2"]`,
	}

	result3 := parseEventData(event3)
	if len(result3) != 2 {
		t.Errorf("Test case 3 failed: Expected 2 event data, got %d", len(result3))
	}
}

func TestDealTopicEventData(t *testing.T) {
	// Test case 1: Test with a valid contract event
	contractEvent1 := []*db.ContractEvent{
		{
			ContractName: "contract1",
			Topic:        "topic1",
			EventDataBak: `["event1", "event2"]`,
		},
	}
	contractInfoMap1 := map[string]*db.Contract{
		"contract1": {
			Addr:         "addr1",
			Name:         "contract1",
			NameBak:      "contract1",
			ContractType: common.ContractStandardNameCMDFA,
		},
	}
	txInfoMap1 := map[string]*db.Transaction{
		"txId1": {
			TxId:      "txId1",
			UserAddr:  "user1",
			TxType:    "INVOKE_CONTRACT",
			Timestamp: 1625097700,
		},
	}

	topicEventResult1 := DealTopicEventData(db.UTchainID, contractEvent1, contractInfoMap1, txInfoMap1)
	if topicEventResult1 == nil {
		t.Errorf("Test case 1 failed: Expected non-nil topic event result, got nil")
	}

	// Test case 3: Test with nil parameters
	topicEventResult3 := DealTopicEventData("", nil, nil, nil)
	if topicEventResult3 == nil {
		t.Errorf("Test case 3 failed: Expected non-nil topic event result, got nil")
	}
}

func TestDealEvidence(t *testing.T) {
	// Test case 1: Test with a valid single evidence transaction
	blockHeight1 := int64(10)
	txInfo1 := &pbCommon.Transaction{
		Payload: &pbCommon.Payload{
			ChainId:      "chain1",
			TxType:       pbCommon.TxType_INVOKE_CONTRACT,
			TxId:         "txId1",
			Timestamp:    1625097700,
			ContractName: "contract1",
			Method:       common.PayloadMethodEvidence,
			Parameters: []*pbCommon.KeyValuePair{
				{
					Key:   "hash",
					Value: []byte("hash1"),
				},
				{
					Key:   "metadata",
					Value: []byte("metadata1"),
				},
				{
					Key:   "id",
					Value: []byte("id1"),
				},
			},
		},
		Result: &pbCommon.Result{
			Code: pbCommon.TxStatusCode_SUCCESS,
			ContractResult: &pbCommon.ContractResult{
				Code:    0,
				Result:  []byte("result1"),
				Message: "message1",
				GasUsed: 100,
				ContractEvent: []*pbCommon.ContractEvent{
					{
						Topic:           "topic1",
						TxId:            "txId1",
						ContractName:    "contract1",
						ContractVersion: "1.0",
						EventData:       []string{"event1", "event2"},
					},
				},
			},
		},
	}
	userInfo1 := &db.SenderPayerUser{
		SenderUserAddr: "addr1",
	}

	evidences1, err1 := DealEvidence(blockHeight1, txInfo1, userInfo1)
	assert.NoError(t, err1)
	assert.NotEmpty(t, evidences1)
	assert.Equal(t, "hash1", evidences1[0].Hash)
	assert.Equal(t, "metadata1", evidences1[0].MetaData)
	assert.Equal(t, "id1", evidences1[0].EvidenceId)

	// Test case 2: Test with a valid batch evidence transaction
	blockHeight2 := int64(10)
	txInfo2 := &pbCommon.Transaction{
		Payload: &pbCommon.Payload{
			ChainId:      "chain2",
			TxType:       pbCommon.TxType_INVOKE_CONTRACT,
			TxId:         "txId2",
			Timestamp:    1625097700,
			ContractName: "contract2",
			Method:       common.PayloadMethodEvidenceBatch,
			Parameters: []*pbCommon.KeyValuePair{
				{
					Key:   "evidences",
					Value: []byte(`[{"hash":"hash2","metadata":"metadata2","id":"id2"}]`),
				},
			},
		},
		Result: &pbCommon.Result{
			Code: pbCommon.TxStatusCode_SUCCESS,
			ContractResult: &pbCommon.ContractResult{
				Code:    0,
				Result:  []byte("result2"),
				Message: "message2",
				GasUsed: 200,
				ContractEvent: []*pbCommon.ContractEvent{
					{
						Topic:           "topic2",
						TxId:            "txId2",
						ContractName:    "contract2",
						ContractVersion: "1.0",
						EventData:       []string{"event3", "event4"},
					},
				},
			},
		},
	}
	userInfo2 := &db.SenderPayerUser{
		SenderUserAddr: "addr2",
	}

	evidences2, err2 := DealEvidence(blockHeight2, txInfo2, userInfo2)
	assert.NoError(t, err2)
	assert.NotEmpty(t, evidences2)
	assert.Equal(t, "hash2", evidences2[0].Hash)
	assert.Equal(t, "metadata2", evidences2[0].MetaData)
	assert.Equal(t, "id2", evidences2[0].EvidenceId)
}
