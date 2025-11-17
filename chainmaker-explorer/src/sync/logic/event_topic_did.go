/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package logic

import (
	"chainmaker_web/src/chain"
	"chainmaker_web/src/db"
	"chainmaker_web/src/sync/model"
	"chainmaker_web/src/utils"
	"encoding/json"

	"chainmaker.org/chainmaker/contract-utils/standard"
	"chainweaver.org.cn/chainweaver/did/core"
	"github.com/google/uuid"
)

// processDIDEvent 函数用于处理DID事件
// @param chainId 链ID
// @param event 合约事件
// @param eventData 合约事件数据
// @param didEventData DID事件数据
func processDIDEvent(chainId string, event *db.ContractEvent, eventData []string, didEventData *model.DIDEventData) {
	//判断是否是did合约
	if !utils.CheckIsDIDContact(event.ContractType) {
		return
	}

	//解析DID事件
	switch event.Topic {
	case standard.Topic_SetDidDocument:
		//处理设置DID文档事件
		DealTopicSetDidDocument(chainId, didEventData, eventData, event)
	case standard.Topic_AddBlackList:
		//处理添加黑名单事件
		DealTopicAddBlockList(didEventData, eventData)
	case standard.Topic_DeleteBlackList:
		//处理删除黑名单事件
		DealTopicDeleteBlockList(didEventData, eventData)
	case standard.Topic_AddTrustIssuer:
		//处理添加信任颁发者事件
		DealTopicAddTrustIssuer(didEventData, eventData)
	case standard.Topic_DeleteTrustIssuer:
		//处理删除信任颁发者事件
		DealTopicDeleteTrustIssuer(didEventData, eventData)
	case standard.Topic_SetVcTemplate:
		//处理设置VC模板事件
		DealTopicSetVCTemplate(didEventData, eventData, event)
	case standard.Topic_VcIssueLog:
		//处理VC发行日志事件
		DealTopicSetVCTLog(didEventData, eventData, event)
	case standard.Topic_RevokeVc:
		//处理VC撤销事件
		DealTopicRevokeVC(didEventData, eventData)
	}
}

// DealTopicSetDidDocument 处理设置DID文档事件
// @param chainId 链ID
// @param didEventData DID事件数据
// @param eventData 事件数据
// @param event 事件
func DealTopicSetDidDocument(chainId string, didEventData *model.DIDEventData, eventData []string,
	event *db.ContractEvent) {
	//绑定，解绑DID
	if len(eventData) < 2 {
		return
	}

	//获取DID和DIDDocument
	did := eventData[0]
	didDocument := &core.Document{}
	err := json.Unmarshal([]byte(eventData[1]), &didDocument)
	if err != nil {
		log.Errorf("DealUserDIDEventData json Unmarshal err, err:%v, eventData[1]:%v", err, eventData[1])
	}

	//判断是否是主did合约
	if chain.CheckIsMainDIDContract(chainId, event.ContractNameBak) {
		//如果是主did合约，则将DID添加到didEventData.DIDUnBinds中
		didEventData.DIDUnBinds = append(didEventData.DIDUnBinds, did)
		//将DIDDocument中的VerificationMethod添加到didEventData.Account中
		for _, value := range didDocument.VerificationMethod {
			didEventData.Account[value.Address] = &db.Account{
				Address: value.Address,
				DID:     did,
			}
		}
	} else {
		//如果不是主did合约，则将DIDDocument中的VerificationMethod添加到didEventData.Account中
		for _, value := range didDocument.VerificationMethod {
			if _, ok := didEventData.Account[value.Address]; !ok {
				didEventData.Account[value.Address] = &db.Account{
					Address: value.Address,
				}
			}
		}
	}

	//生成issuerService和accountJson
	var (
		issuerService  string
		authentication string
		accountJson    string
	)
	if len(didDocument.Service) > 0 {
		//将DIDDocument中的Service转换为json字符串
		serviceJsonByte, _ := json.Marshal(didDocument.Service)
		issuerService = string(serviceJsonByte)
	}
	if len(didDocument.Authentication) > 0 {
		//将DIDDocument中的Authentication转换为json字符串
		authenticationByte, _ := json.Marshal(didDocument.Authentication)
		authentication = string(authenticationByte)
	}
	if len(didDocument.VerificationMethod) > 0 {
		//将DIDDocument中的VerificationMethod转换为json字符串
		accountList := make([]*db.DIDAccountDate, 0)
		for _, value := range didDocument.VerificationMethod {
			accountList = append(accountList, &db.DIDAccountDate{
				ID:      value.Id,
				Address: value.Address,
			})
		}
		accountJsonByte, _ := json.Marshal(accountList)
		accountJson = string(accountJsonByte)
	}

	//生成新的UUID
	newUUID := uuid.New().String()
	//生成DIDSetHistory
	didhistory := &db.DIDSetHistory{
		ID:             newUUID,
		DID:            did,
		TxId:           event.TxId,
		ContractName:   event.ContractNameBak,
		ContractAddr:   event.ContractAddr,
		Document:       eventData[1],
		IssuerService:  issuerService,
		Authentication: authentication,
		AccountJson:    accountJson,
		Timestamp:      event.Timestamp,
	}
	//将DIDSetHistory添加到didEventData中
	didEventData.DIDSetHistory = append(didEventData.DIDSetHistory, didhistory)

	//生成DIDDetail
	newUUID = uuid.New().String()
	didDetail := &db.DIDDetail{
		ID:             newUUID,
		DID:            did,
		ContractName:   event.ContractNameBak,
		ContractAddr:   event.ContractAddr,
		Document:       eventData[1],
		IssuerService:  issuerService,
		Authentication: authentication,
		AccountJson:    accountJson,
		CreateTxId:     event.TxId,
		Timestamp:      event.Timestamp,
	}

	//将DIDDetail添加到didEventData中
	didEventData.DIDDetail = append(didEventData.DIDDetail, didDetail)
}

// DealTopicAddBlockList 函数用于处理DID事件数据，将eventData中的DID添加到黑名单中
// @param didEventData DID事件数据
// @param eventData 事件数据
func DealTopicAddBlockList(didEventData *model.DIDEventData, eventData []string) {
	// 如果eventData为空，则直接返回
	if len(eventData) == 0 {
		return
	}

	var issuerDid []string
	_ = json.Unmarshal([]byte(eventData[0]), &issuerDid)
	if len(issuerDid) == 0 {
		return
	}

	for _, value := range issuerDid {
		if value == "" {
			continue
		}
		didEventData.DIDAddBlacks = append(didEventData.DIDAddBlacks, value)
	}
}

// DealTopicDeleteBlockList 函数用于处理DID事件数据，将eventData中的DID从黑名单中删除
// @param didEventData DID事件数据
// @param eventData 事件数据
func DealTopicDeleteBlockList(didEventData *model.DIDEventData, eventData []string) {
	// 如果eventData为空，则直接返回
	if len(eventData) == 0 {
		return
	}

	var issuerDid []string
	_ = json.Unmarshal([]byte(eventData[0]), &issuerDid)
	if len(issuerDid) == 0 {
		return
	}

	for _, value := range issuerDid {
		if value == "" {
			continue
		}
		didEventData.DIDDeleteBlacks = append(didEventData.DIDDeleteBlacks, value)
	}
}

// DealTopicAddTrustIssuer 处理添加信任颁发者事件
// @param didEventData DID事件数据
// @param eventData 事件数据
func DealTopicAddTrustIssuer(didEventData *model.DIDEventData, eventData []string) {
	// 如果eventData为空，则直接返回
	if len(eventData) == 0 {
		return
	}

	var issuerDid []string
	_ = json.Unmarshal([]byte(eventData[0]), &issuerDid)
	if len(issuerDid) == 0 {
		return
	}
	didEventData.DIDAddIssuers = append(didEventData.DIDAddIssuers, issuerDid...)
}

// DealTopicDeleteTrustIssuer 处理删除信任颁发者事件
// @param didEventData DID事件数据
// @param eventData 事件数据
func DealTopicDeleteTrustIssuer(didEventData *model.DIDEventData, eventData []string) {
	// 如果eventData为空，则直接返回
	if len(eventData) == 0 {
		return
	}

	var issuerDid []string
	_ = json.Unmarshal([]byte(eventData[0]), &issuerDid)
	if len(issuerDid) == 0 {
		return
	}
	didEventData.DIDDeleteIssuers = append(didEventData.DIDDeleteIssuers, issuerDid...)
}

// DealTopicSetVCTemplate 处理设置VC模板事件
// @param didEventData DID事件数据
// @param eventData 事件数据
// @param event 合约事件
func DealTopicSetVCTemplate(didEventData *model.DIDEventData, eventData []string, event *db.ContractEvent) {
	// 如果eventData为空，则直接返回
	if len(eventData) == 0 || len(eventData) < 5 {
		return
	}

	newUUID := uuid.New().String()
	template := &db.VCTemplate{
		ID:           newUUID,
		ContractName: event.ContractNameBak,
		ContractAddr: event.ContractAddr,
		TemplateID:   eventData[0],
		TemplateName: eventData[1],
		VCType:       eventData[2],
		Version:      eventData[3],
		Template:     eventData[4],
		TxId:         event.TxId,
		Timestamp:    event.Timestamp,
	}

	// 如果eventData长度大于6，则将eventData[5]和eventData[6]分别赋值给template.ShortName和template.CreateDID
	if len(eventData) > 6 {
		template.ShortName = eventData[5]
		template.CreateDID = eventData[6]
	}
	didEventData.VCTemplate = append(didEventData.VCTemplate, template)
}

// DealTopicSetVCTLog 处理VC发行日志事件
// @param didEventData DID事件数据
// @param eventData 事件数据
// @param event 合约事件
func DealTopicSetVCTLog(didEventData *model.DIDEventData, eventData []string, event *db.ContractEvent) {
	// 如果eventData为空，则直接返回
	if len(eventData) == 0 || len(eventData) < 4 {
		return
	}

	newUUID := uuid.New().String()
	vcLog := &db.VCIssueHistory{
		ID:           newUUID,
		ContractName: event.ContractNameBak,
		ContractAddr: event.ContractAddr,
		IssuerDID:    eventData[0],
		HolderDID:    eventData[1],
		TemplateID:   eventData[2],
		VCID:         eventData[3],
		Timestamp:    event.Timestamp,
		TxId:         event.TxId,
	}
	didEventData.VCIssueHistory = append(didEventData.VCIssueHistory, vcLog)
}

// DealTopicRevokeVC 处理主题撤销VC
// @param didEventData DID事件数据
// @param eventData 事件数据
func DealTopicRevokeVC(didEventData *model.DIDEventData, eventData []string) {
	// 如果eventData长度为0，则直接返回
	if len(eventData) == 0 {
		return
	}
	didEventData.VCDeleteIds = append(didEventData.VCDeleteIds, eventData...)
}
