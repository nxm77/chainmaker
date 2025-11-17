/*
Package logic comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package logic

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/sync/model"
)

// DealDIDSaveData 函数用于处理DID保存数据
func DealDIDSaveData(didEventData *model.DIDEventData) *db.DIDSaveData {
	// 创建一个新的DID保存数据对象
	didSaveData := db.NewDIDSaveData()
	// 将DID详细信息保存到DID保存数据对象中
	didSaveData.SaveDIDDetails = didEventData.DIDDetail
	// 将DID设置历史保存到DID保存数据对象中
	didSaveData.InsertDIDHistory = didEventData.DIDSetHistory
	// 将VCTemplate保存到DID保存数据对象中
	didSaveData.SaveVCTemplates = didEventData.VCTemplate
	// 将VCIssueHistory保存到DID保存数据对象中
	didSaveData.InsertVCHistory = didEventData.VCIssueHistory
	// 遍历DIDAddBlacks，将DID状态设置为黑名单
	for _, did := range didEventData.DIDAddBlacks {
		didSaveData.UpdateDIDStatus[did] = db.DIDStatusBlack
	}
	// 遍历DIDDeleteBlacks，将DID状态设置为成功
	for _, did := range didEventData.DIDDeleteBlacks {
		didSaveData.UpdateDIDStatus[did] = db.DIDStatusSuccess
	}

	// 遍历DIDAddIssuers，将DID状态设置为发行者
	for _, did := range didEventData.DIDAddIssuers {
		didSaveData.UpdateDIDIssuer[did] = db.DIDIsIssuer
	}
	// 遍历DIDDeleteIssuers，将DID状态设置为非发行者
	for _, did := range didEventData.DIDDeleteIssuers {
		didSaveData.UpdateDIDIssuer[did] = db.DIDNotIssuer
	}
	// 遍历VCDeleteIds，将VC状态设置为删除
	for _, vcId := range didEventData.VCDeleteIds {
		didSaveData.UpdateVCStatus[vcId] = db.VCStatusDelete
	}

	return didSaveData
}
