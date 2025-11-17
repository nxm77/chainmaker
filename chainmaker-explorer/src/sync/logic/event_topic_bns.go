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

	"chainmaker.org/chainmaker/contract-utils/standard"
)

// 解析BNS事件，BNS绑定，解绑
func processBNSEvent(event *db.ContractEvent, eventData []string, bnsAccountMap map[string]*db.BNSTopicEventData,
	topicEventResult *model.TopicEventResult) {
	bnsBindEventData, bnsUnBindDomain := DealUserBNSEventData(event.ContractName, event.Topic, eventData)
	if bnsBindEventData != nil {
		bnsAccountMap[bnsBindEventData.Domain] = bnsBindEventData
	}
	//bns解绑
	if bnsUnBindDomain != "" {
		//如果前面有绑定，需要删除
		delete(bnsAccountMap, bnsUnBindDomain)
		topicEventResult.BNSUnBindDomain = append(topicEventResult.BNSUnBindDomain, bnsUnBindDomain)
	}
}

// DealUserBNSEventData 解析BNS eventdata
func DealUserBNSEventData(contractName, topic string, eventData []string) (*db.BNSTopicEventData, string) {
	var topicEventData *db.BNSTopicEventData
	var unBindDomain string
	if contractName != common.PayloadContractNameBNS {
		return topicEventData, unBindDomain
	}

	switch topic {
	case standard.Topic_Bind:
		//绑定BNS
		if len(eventData) < 3 {
			return topicEventData, unBindDomain
		}

		//ResourceType, _ := strconv.ParseInt(eventData[2], 10, 64)
		////BNS解析资源类型,string "1“-链地址，”2"-DID,"3"-去中心化网站，"4“-合约，"5"-子链
		//if ResourceType > 1 {
		//	//暂时只解析链地址
		//	return topicEventData
		//}

		// 解析数据
		topicEventData = &db.BNSTopicEventData{
			Domain:       eventData[0],
			Value:        eventData[1],
			ResourceType: eventData[2],
		}
	case standard.Topic_UnBind:
		//解绑BNS
		if len(eventData) == 0 {
			return topicEventData, unBindDomain
		}
		// 解析数据
		unBindDomain = eventData[0]
	}

	return topicEventData, unBindDomain
}
