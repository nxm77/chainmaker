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
	"math/big"

	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
)

// processBlackListEvent 构造交易黑名单数据
func processBlackListEvent(topicEventResult *model.TopicEventResult, contractName, topic string, eventData []string) {
	//解析黑名单交易，添加，删除黑名单
	addBlack, deleteBlack := DealBackListEventData(contractName, topic, eventData)
	if len(addBlack) > 0 {
		topicEventResult.AddBlack = append(topicEventResult.AddBlack, addBlack...)
	}
	if len(deleteBlack) > 0 {
		topicEventResult.DeleteBlack = append(topicEventResult.DeleteBlack, deleteBlack...)
	}
}

// DealBackListEventData 解析黑名单交易
func DealBackListEventData(contractName, topic string, eventData []string) ([]string, []string) {
	addBlack := make([]string, 0)
	deleteBlack := make([]string, 0)
	//第一条记录是链ID
	if contractName != syscontract.SystemContract_TRANSACTION_MANAGER.String() || len(eventData) <= 1 {
		return addBlack, deleteBlack
	}

	//加入黑名单
	if topic == common.TopicTxAddBlack {
		for i := 1; i < len(eventData); i++ {
			addBlack = append(addBlack, eventData[i])
		}
	} else if topic == common.TopicTxDeleteBlack {
		//解封黑名单
		for i := 1; i < len(eventData); i++ {
			deleteBlack = append(deleteBlack, eventData[i])
		}
	}
	return addBlack, deleteBlack
}

// BuildTransferEventData 构造流转数据
func BuildTransferEventData(topicEventResult *model.TopicEventResult, ownerAddrMap map[string]string,
	contractInfoMap map[string]*db.Contract, event *db.ContractEvent, senderUser string,
	eventData []string) map[string]string {
	//合约信息
	contract := contractInfoMap[event.ContractName]
	if contract == nil {
		return ownerAddrMap
	}

	//根据eventData解析transfer流转记录
	topicEventData := GetTransferEventDta(contract.ContractType, event.Topic, senderUser, eventData)
	if topicEventData == nil {
		return ownerAddrMap
	}

	if topicEventData.TokenId != "" || topicEventData.Amount != "" {
		//统计持仓地址
		if topicEventData.FromAddress != "" {
			if _, ok := ownerAddrMap[topicEventData.FromAddress]; !ok {
				ownerAddrMap[topicEventData.FromAddress] = topicEventData.FromAddress
			}
		}
		if topicEventData.ToAddress != "" {
			if _, ok := ownerAddrMap[topicEventData.ToAddress]; !ok {
				ownerAddrMap[topicEventData.ToAddress] = topicEventData.ToAddress
			}
		}
	}

	transferData := &db.ContractEventData{
		Topic:        event.Topic,
		Index:        event.EventIndex,
		TxId:         event.TxId,
		ContractName: event.ContractName,
		EventData:    topicEventData,
		Timestamp:    event.Timestamp,
	}
	topicEventResult.ContractEventData = append(topicEventResult.ContractEventData, transferData)
	return ownerAddrMap
}

func GetTransferEventDta(contractType, topic, senderUser string, eventData []string) *db.TransferTopicEventData {
	var topicEventData *db.TransferTopicEventData
	switch contractType {
	case common.ContractStandardNameCMDFA:
		topicEventData = DealDockerDFAEventData(topic, eventData)
	case common.ContractStandardNameCMNFA:
		topicEventData = DealDockerNFAEventData(topic, eventData, senderUser)
	case common.ContractStandardNameEVMDFA:
		topicEventData = DealEVMDFAEventData(topic, eventData)
	case common.ContractStandardNameEVMNFA:
		topicEventData = DealEVMNFAEventData(topic, eventData)
	}
	// if topicEventData != nil {
	// 	//有部分初始化会给默认FromAddress
	// 	if topicEventData.FromAddress == topicEventData.ToAddress {
	// 		topicEventData.FromAddress = ""
	// 	}
	// }

	return topicEventData
}

// DealDockerDFAEventData 同质化docker解析eventdata
func DealDockerDFAEventData(topic string, eventData []string) *db.TransferTopicEventData {
	var topicEventData *db.TransferTopicEventData
	switch topic {
	case common.TopicMintEvent:
		if len(eventData) < 2 {
			return topicEventData
		}
		topicEventData = &db.TransferTopicEventData{
			ToAddress: eventData[0],
			Amount:    eventData[1],
		}
	case common.TopicTransferEvent:
		if len(eventData) < 3 {
			return topicEventData
		}
		// 解析数据
		topicEventData = &db.TransferTopicEventData{
			FromAddress: eventData[0],
			ToAddress:   eventData[1],
			Amount:      eventData[2],
		}
	case common.TopicBurnEvent:
		if len(eventData) < 2 {
			return topicEventData
		}
		topicEventData = &db.TransferTopicEventData{
			FromAddress: eventData[0],
			Amount:      eventData[1],
		}
	case common.TopicApproveEvent:
		if len(eventData) < 3 {
			return topicEventData
		}
		// 解析数据
		topicEventData = &db.TransferTopicEventData{
			FromAddress: eventData[0],
			ToAddress:   eventData[1],
			Amount:      eventData[2],
		}
	}
	return topicEventData
}

// DealDockerNFAEventData 非同质化docker解析eventdata
func DealDockerNFAEventData(topic string, eventData []string, senderUser string) *db.TransferTopicEventData {
	var topicEventData *db.TransferTopicEventData
	switch topic {
	case "Mint":
		//发行
		if len(eventData) < 5 {
			return topicEventData
		}
		if common.IsZeroAddress(eventData[0]) {
			eventData[0] = ""
		}
		topicEventData = &db.TransferTopicEventData{
			FromAddress:  eventData[0],
			ToAddress:    eventData[1],
			TokenId:      eventData[2],
			CategoryName: eventData[3],
			Metadata:     eventData[4],
		}
	case "TransferFrom":
		//转账
		if len(eventData) < 3 {
			return topicEventData
		}
		// 解析数据
		topicEventData = &db.TransferTopicEventData{
			FromAddress: eventData[0],
			ToAddress:   eventData[1],
			TokenId:     eventData[2],
		}
	case "Burn":
		//销毁
		if len(eventData) == 0 {
			return topicEventData
		}
		topicEventData = &db.TransferTopicEventData{
			FromAddress: senderUser,
			TokenId:     eventData[0],
		}
	case "SetApproval":
		//授权
		if len(eventData) < 4 {
			return topicEventData
		}
		// 解析数据
		topicEventData = &db.TransferTopicEventData{
			FromAddress: eventData[0],
			ToAddress:   eventData[1],
			TokenId:     eventData[2],
			Approval:    eventData[3],
		}
	case "SetApprovalForAll":
		//全部授权
		if len(eventData) < 3 {
			return topicEventData
		}
		// 解析数据
		topicEventData = &db.TransferTopicEventData{
			FromAddress: eventData[0],
			ToAddress:   eventData[1],
			Approval:    eventData[2],
		}
	case "SetApprovalByCategory":
		//分类授权
		if len(eventData) < 4 {
			return topicEventData
		}
		// 解析数据
		topicEventData = &db.TransferTopicEventData{
			FromAddress:  eventData[0],
			ToAddress:    eventData[1],
			CategoryName: eventData[2],
			Approval:     eventData[3],
		}
	}

	return topicEventData
}

// DealEVMDFAEventData 解析eventdata
func DealEVMDFAEventData(topic string, eventData []string) *db.TransferTopicEventData {
	var topicEventData *db.TransferTopicEventData
	switch topic {
	case common.EVMEventTopicTransfer:
		if len(eventData) < 3 {
			return topicEventData
		}
		if len(eventData[0]) < 24 || len(eventData[1]) < 24 {
			return topicEventData
		}
		// 解析数据
		fromAddress := eventData[0][24:]
		toAddress := eventData[1][24:]
		// 将十六进制字符串转换为big.Int类型
		bigInt := new(big.Int)
		bigInt.SetString(eventData[2], 16)
		// 将big.Int类型转换为十进制字符串
		amountStr := bigInt.String()
		if common.IsZeroAddress(fromAddress) {
			fromAddress = ""
		}
		if common.IsZeroAddress(toAddress) {
			toAddress = ""
		}
		topicEventData = &db.TransferTopicEventData{
			FromAddress: fromAddress,
			ToAddress:   toAddress,
			Amount:      amountStr,
		}
	}
	return topicEventData
}

// DealEVMNFAEventData 解析eventdata
func DealEVMNFAEventData(topic string, eventData []string) *db.TransferTopicEventData {
	var topicEventData *db.TransferTopicEventData
	switch topic {
	case common.EVMEventTopicTransfer:
		if len(eventData) < 3 {
			return topicEventData
		}
		if len(eventData[0]) < 24 || len(eventData[1]) < 24 {
			return topicEventData
		}
		// 解析数据
		fromAddress := eventData[0][24:]
		toAddress := eventData[1][24:]
		// 将十六进制字符串转换为big.Int类型
		bigInt := new(big.Int)
		bigInt.SetString(eventData[2], 16)
		// 将big.Int类型转换为十进制字符串
		tokenId := bigInt.String()
		if common.IsZeroAddress(fromAddress) {
			fromAddress = ""
		}
		if common.IsZeroAddress(toAddress) {
			toAddress = ""
		}
		topicEventData = &db.TransferTopicEventData{
			FromAddress: fromAddress,
			ToAddress:   toAddress,
			TokenId:     tokenId,
		}
	}
	return topicEventData
}
