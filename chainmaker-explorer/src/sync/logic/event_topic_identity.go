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

	"github.com/google/uuid"
)

// BuildIdentityEventData 解析身份认证合约
func BuildIdentityEventData(topicEventResult *model.TopicEventResult, contractInfoMap map[string]*db.Contract,
	event *db.ContractEvent, eventData []string) {
	//合约信息
	contract := contractInfoMap[event.ContractName]
	if contract == nil || contract.ContractType != common.ContractStandardNameCMID {
		return
	}

	//解析身份认证合约
	identityEventData := DealIdentityEventData(event.Topic, eventData)
	if identityEventData == nil {
		return
	}

	newUUID := uuid.New().String()
	tempIdentity := &db.IdentityContract{
		ID:           newUUID,
		TxId:         event.TxId,
		EventIndex:   event.EventIndex,
		ContractName: contract.Name,
		ContractAddr: contract.Addr,
		UserAddr:     identityEventData.UserAddr,
		Level:        identityEventData.Level,
		PkPem:        identityEventData.PkPem,
	}
	topicEventResult.IdentityContract = append(topicEventResult.IdentityContract, tempIdentity)
}

// DealIdentityEventData 解析身份认证数据
func DealIdentityEventData(topic string, eventData []string) *db.IdentityEventData {
	var identityEvent *db.IdentityEventData
	if topic == "setIdentity" {
		if len(eventData) < 3 {
			return identityEvent
		}
		identityEvent = &db.IdentityEventData{
			UserAddr: eventData[0],
			Level:    eventData[1],
			PkPem:    eventData[2],
		}
	}
	return identityEvent
}
