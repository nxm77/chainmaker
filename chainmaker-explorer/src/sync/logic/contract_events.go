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

	"github.com/google/uuid"

	"chainmaker.org/chainmaker/contract-utils/standard"
	pbCommon "chainmaker.org/chainmaker/pb-go/v2/common"
)

// 创建一个新的TopicEventResult结构体
func newTopicEventResult() *model.TopicEventResult {
	return &model.TopicEventResult{
		//AddBlack 添加黑名单交易
		AddBlack: make([]string, 0),
		//DeleteBlack 删除黑名单交易
		DeleteBlack: make([]string, 0),
		//IdentityContract 身份合约
		IdentityContract: make([]*db.IdentityContract, 0),
		//ContractEventData 合约事件
		ContractEventData: make([]*db.ContractEventData, 0),
		//OwnerAdders 合约拥有者
		OwnerAdders: make([]string, 0),
		//BNSBindEventData BNS绑定事件
		BNSBindEventData: make([]*db.BNSTopicEventData, 0),
		//BNSUnBindDomain BNS解绑域名
		BNSUnBindDomain: make([]string, 0),
	}
}

// 创建一个新的IDAEventResult结构体
func newIDAEventResult() *model.IDAEventData {
	return &model.IDAEventData{
		//IDACreatedMap 创建IDA合约
		IDACreatedMap: make(map[string][]*db.IDACreatedInfo),
		//IDAUpdatedMap 更新IDA合约
		IDAUpdatedMap: make(map[string][]*db.EventIDAUpdatedData),
		//IDADeletedCodeMap 删除IDA合约
		IDADeletedCodeMap: make(map[string][]string),
	}
}

// DealContractEvents 解析所有合约事件
// @param txInfo 交易信息
// @return 合约事件
func DealContractEvents(txInfo *pbCommon.Transaction) []*db.ContractEvent {
	contractEvents := make([]*db.ContractEvent, 0)
	//失败的操作不处理
	if txInfo.Result.ContractResult == nil ||
		txInfo.Result.ContractResult.Code != 0 {
		return contractEvents
	}

	resEvent := txInfo.Result.ContractResult.ContractEvent
	// 处理合约交易事件
	for i, event := range resEvent {
		eventDataJson, _ := json.Marshal(event.EventData)
		newUUID := uuid.New().String()
		contractEvent := &db.ContractEvent{
			ID:              newUUID,
			Topic:           event.Topic,
			EventIndex:      i + 1,
			TxId:            event.TxId,
			ContractName:    event.ContractName,
			ContractNameBak: event.ContractName,
			ContractVersion: event.ContractVersion,
			EventData:       string(eventDataJson),
			Timestamp:       txInfo.Payload.Timestamp,
		}
		contractEvents = append(contractEvents, contractEvent)
	}
	return contractEvents
}

// parseEventData 解析eventDate
// @param event 事件
// @return 解析后的数据
func parseEventData(event *db.ContractEvent) []string {
	var eventData []string
	if event.EventDataBak != "" {
		_ = json.Unmarshal([]byte(event.EventDataBak), &eventData)
	} else if event.EventData != "" {
		_ = json.Unmarshal([]byte(event.EventData), &eventData)
	}
	return eventData
}

// DealTopicEventData 解析eventDate
// 处理合约事件数据
func DealTopicEventData(chainId string, contractEvent []*db.ContractEvent, contractInfoMap map[string]*db.Contract,
	txInfoMap map[string]*db.Transaction) *model.TopicEventResult {
	// 创建一个map，用于存储合约拥有者的地址
	ownerAddrMap := make(map[string]string, 0)
	// 创建一个map，用于存储BNS主题事件数据
	bnsAccountMap := make(map[string]*db.BNSTopicEventData, 0)
	// 创建一个TopicEventResult对象
	topicEventResult := newTopicEventResult()
	//IDA数据
	topicEventResult.IDAEventData = newIDAEventResult()
	//DID数据
	topicEventResult.DIDEventData = model.NewDIDEventData()

	if len(contractEvent) == 0 {
		return topicEventResult
	}

	for _, event := range contractEvent {
		eventData := parseEventData(event)
		if len(eventData) == 0 {
			continue
		}

		//解析IDA数据
		processIDAEvent(event, eventData, topicEventResult.IDAEventData)
		//解析BNS事件，BNS绑定，解绑
		processBNSEvent(event, eventData, bnsAccountMap, topicEventResult)
		//解析DID事件，设置DID
		processDIDEvent(chainId, event, eventData, topicEventResult.DIDEventData)
		//解析黑名单交易，添加，删除黑名单
		processBlackListEvent(topicEventResult, event.ContractName, event.Topic, eventData)

		//合约信息
		contract := contractInfoMap[event.ContractName]
		if contract == nil {
			continue
		}

		//解析身份认证合约
		BuildIdentityEventData(topicEventResult, contractInfoMap, event, eventData)

		//流转数据解析
		var senderUser string
		if txInfo, ok := txInfoMap[event.TxId]; ok {
			senderUser = txInfo.UserAddr
		}
		//根据eventData解析transfer流转记录
		ownerAddrMap = BuildTransferEventData(topicEventResult, ownerAddrMap, contractInfoMap, event,
			senderUser, eventData)
	}

	//持仓地址列表
	ownerAdders := make([]string, 0)
	for _, addr := range ownerAddrMap {
		ownerAdders = append(ownerAdders, addr)
	}
	topicEventResult.OwnerAdders = ownerAdders

	for _, value := range bnsAccountMap {
		topicEventResult.BNSBindEventData = append(topicEventResult.BNSBindEventData, value)
	}

	return topicEventResult
}

// UpdateContractTxAndEventNum 更新合约交易数和事件数
func UpdateContractTxAndEventNum(minHeight int64, contractMap map[string]*db.Contract,
	txList map[string]*db.Transaction, contractEvent []*db.ContractEvent) []*db.Contract {
	contractTxNumMap := make(map[string]int64, 0)
	contractEventNumMap := make(map[string]int64, 0)
	updateContractMap := make(map[string]*db.Contract, 0)
	updateContractNum := make([]*db.Contract, 0)
	if len(contractMap) == 0 || len(txList) == 0 {
		return updateContractNum
	}

	//统计本次交易数据量
	for _, txInfo := range txList {
		if contract, ok := contractMap[txInfo.ContractNameBak]; ok {
			//说明已经更新过了
			if contract.BlockHeight >= minHeight {
				continue
			}
			contractTxNumMap[contract.Addr]++
		}
	}

	//统计本次合约事件数据量
	for _, event := range contractEvent {
		if contract, ok := contractMap[event.ContractNameBak]; ok {
			//说明已经更新过了
			if contract.BlockHeight >= minHeight {
				continue
			}
			contractEventNumMap[contract.Addr]++
		}
	}

	for addr, txNum := range contractTxNumMap {
		var contractInfo *db.Contract
		if contract, ok := updateContractMap[addr]; ok {
			contractInfo = contract
		} else if contract, ok = contractMap[addr]; ok {
			contractInfo = contract
		} else {
			continue
		}
		contractInfo.TxNum = contractInfo.TxNum + txNum
		contractInfo.BlockHeight = minHeight
		updateContractMap[addr] = contractInfo
	}

	for addr, eventNum := range contractEventNumMap {
		var contractInfo *db.Contract
		if contract, ok := updateContractMap[addr]; ok {
			contractInfo = contract
		} else if contract, ok = contractMap[addr]; ok {
			contractInfo = contract
		} else {
			continue
		}
		contractInfo.EventNum = contractInfo.EventNum + eventNum
		contractInfo.BlockHeight = minHeight
		updateContractMap[addr] = contractInfo
	}

	for _, contract := range updateContractMap {
		updateContractNum = append(updateContractNum, contract)
	}

	return updateContractNum
}

// DealEvidence 处理存证合约
// @param blockHeight 区块高度
// @param txInfo 交易信息
// @param userInfo 用户信息
// @return evidences 存证信息
// @return err 错误信息
func DealEvidence(blockHeight int64, txInfo *pbCommon.Transaction, userInfo *db.SenderPayerUser) (
	evidences []*db.EvidenceContract, err error) {
	evidences = make([]*db.EvidenceContract, 0)
	//判断是否是存证合约，如果是存证合约，则解析存证信息
	if txInfo.Payload.Method != common.PayloadMethodEvidence &&
		txInfo.Payload.Method != common.PayloadMethodEvidenceBatch {
		return evidences, nil
	}

	//失败的操作不处理
	if txInfo.Result.ContractResult == nil ||
		txInfo.Result.ContractResult.Code != 0 {
		return evidences, nil
	}

	contractName := txInfo.Payload.ContractName
	tempEvidence := &db.EvidenceContract{
		ContractName:       contractName,
		TxId:               txInfo.Payload.TxId,
		SenderAddr:         userInfo.SenderUserAddr,
		Timestamp:          txInfo.Payload.Timestamp,
		BlockHeight:        blockHeight,
		ContractResult:     txInfo.Result.ContractResult.Result,
		ContractResultCode: txInfo.Result.ContractResult.Code,
	}
	//判断是单条存证还是批量存证，如果是单条存证，则解析单条存证信息，如果是批量存证，则解析批量存证信息
	namespace := uuid.NameSpaceOID // 固定命名空间
	if txInfo.Payload.Method == common.PayloadMethodEvidence {
		newUUID := uuid.New().String()
		//单条存证
		for _, parameter := range txInfo.Payload.Parameters {
			if parameter.Key == "hash" {
				tempEvidence.Hash = string(parameter.Value)
				newUUID = uuid.NewSHA1(namespace, parameter.Value).String()
			}
			if parameter.Key == "metadata" {
				tempEvidence.MetaData = string(parameter.Value)
			}
			if parameter.Key == "id" {
				tempEvidence.EvidenceId = string(parameter.Value)
			}
		}
		tempEvidence.ID = newUUID
		evidences = append(evidences, tempEvidence)
	} else if txInfo.Payload.Method == common.PayloadMethodEvidenceBatch {
		//批量存证
		for _, parameter := range txInfo.Payload.Parameters {
			if parameter.Key == "evidences" {
				standardEvidences := make([]standard.Evidence, 0)
				err = json.Unmarshal(parameter.Value, &standardEvidences)
				if err != nil {
					return evidences, err
				}
				for _, e := range standardEvidences {
					newUUID := uuid.NewSHA1(namespace, []byte(e.Hash)).String()
					tempEvidence.ID = newUUID
					tempEvidence.EvidenceId = e.Id
					tempEvidence.Hash = e.Hash
					tempEvidence.MetaData = e.Metadata
					evidences = append(evidences, tempEvidence)
				}
			}
		}
	}
	return evidences, err
}

// DealEventTopicTxNum 统计每个Topic在当前区块中出现的次数
// @param contractEvents 合约事件
// @return idaInfoList 解析后的IDA信息
func DealEventTopicTxNum(contractEvents []*db.ContractEvent) map[string]map[string]int64 {
	// 创建一个嵌套映射，外层映射以 ContractName 为键，内层映射以 Topic 为键
	counts := make(map[string]map[string]int64)

	for _, event := range contractEvents {
		// 初始化内层映射
		if counts[event.ContractName] == nil {
			counts[event.ContractName] = make(map[string]int64)
		}
		// 统计 Topic 的出现次数
		counts[event.ContractName][event.Topic]++
	}
	return counts
}

// ProcessEventTopicTxNum 统计每个Topic在当前区块中出现的次数
// @param eventTopicTxNum 每个Topic在当前区块中出现的次数
// @param eventTopicDBMap 每个Topic在数据库中的最新数据
// @param blockHeight 当前区块高度
// @return insertEventTpic 需要插入的Topic数据
// ProcessEventTopicTxNum 函数用于处理事件主题的交易数量
func ProcessEventTopicTxNum(eventTopicTxNum map[string]map[string]int64,
	eventTopicDBMap map[string]map[string]*db.ContractEventTopic, blockHeight int64) (
	[]*db.ContractEventTopic, []*db.ContractEventTopic) {
	// 定义插入和更新事件主题的切片
	insertEventTpic := make([]*db.ContractEventTopic, 0)
	updateEventTpic := make([]*db.ContractEventTopic, 0)
	// 遍历事件主题的交易数量
	for contractName, topics := range eventTopicTxNum {
		for topic, num := range topics {
			// 如果事件主题在数据库中存在
			if eventMapDB, exists := eventTopicDBMap[contractName]; exists {
				if eventDB, ok := eventMapDB[topic]; ok {
					// 如果当前区块高度小于等于数据库中的区块高度，则跳过
					if blockHeight <= eventDB.BlockHeight {
						continue
					}

					// 更新交易数量
					txNum := eventDB.TxNum + num
					updateEventTpic = append(updateEventTpic, &db.ContractEventTopic{
						Topic:        topic,
						ContractName: contractName,
						BlockHeight:  blockHeight,
						TxNum:        txNum,
					})
					continue
				}
			}

			// 生成新的UUID
			newUUID := uuid.New().String()
			// 插入新的事件主题
			insert := &db.ContractEventTopic{
				ID:           newUUID,
				Topic:        topic,
				ContractName: contractName,
				TxNum:        num,
				BlockHeight:  blockHeight,
			}
			insertEventTpic = append(insertEventTpic, insert)
		}
	}
	// 返回插入和更新的事件主题
	return insertEventTpic, updateEventTpic
}
