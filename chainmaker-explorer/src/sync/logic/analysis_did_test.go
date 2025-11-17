/*
Package logic comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package logic

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/sync/model"
	"testing"
)

func TestDealDIDSaveData(t *testing.T) {
	// 测试用例1：所有字段都为空
	didEventData1 := &model.DIDEventData{}
	didSaveData1 := DealDIDSaveData(didEventData1)
	if didSaveData1.SaveDIDDetails != nil || didSaveData1.InsertDIDHistory != nil || didSaveData1.SaveVCTemplates != nil || didSaveData1.InsertVCHistory != nil || len(didSaveData1.UpdateDIDStatus) != 0 || len(didSaveData1.UpdateDIDIssuer) != 0 || len(didSaveData1.UpdateVCStatus) != 0 {
		t.Errorf("Test case 1 failed")
	}

	// 测试用例2 ：所有字段都有值
	didEventData2 := &model.DIDEventData{
		DIDDetail:        []*db.DIDDetail{},
		DIDSetHistory:    []*db.DIDSetHistory{},
		VCTemplate:       []*db.VCTemplate{},
		DIDUnBinds:       []string{"DID1", "DID2"},
		DIDAddBlacks:     []string{"DID1", "DID2"},
		DIDDeleteBlacks:  []string{"DID3", "DID4"},
		DIDAddIssuers:    []string{"DID5", "DID6"},
		DIDDeleteIssuers: []string{"DID7", "DID8"},
		VCDeleteIds:      []string{"VC1", "VC2"},
	}
	_ = DealDIDSaveData(didEventData2)
}
