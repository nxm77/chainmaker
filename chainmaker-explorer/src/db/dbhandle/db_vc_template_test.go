/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"chainmaker_web/src/db"
	"testing"
)

func insertVCTemplateTest() {
	vcTemplate1 := &db.VCTemplate{
		TemplateID:   "template1",
		TemplateName: "template1",
		ShortName:    "short1",
		VCType:       "type1",
		Version:      "1.0",
		Template:     "template1",
		TxId:         "tx1",
		CreateDID:    "did1",
	}
	_ = InsertVCTemplate(db.UTchainID, vcTemplate1)
}

func TestGetVCTemplateById(t *testing.T) {
	insertVCTemplateTest()
	// Test case 1: Get a VC template by ID
	templateId1 := "template1"
	_, err1 := GetVCTemplateById(db.UTchainID, templateId1, "")
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}

	// Test case 2: Get a VC template by non-existing ID
	templateId2 := "non-existing-template"
	_, err2 := GetVCTemplateById(db.UTchainID, templateId2, "")
	if err2 != nil {
		t.Errorf("Test case 2 failed: %v", err2)
	}
}

func TestUpdateVCTemplate(t *testing.T) {
	insertVCTemplateTest()
	// Test case 1: Update an existing VC template
	template1 := &db.VCTemplate{TemplateID: "template1", TemplateName: "updated-template1", ShortName: "updated-short1", VCType: "updated-type1", Version: "2.0", Template: "updated-template1", TxId: "updated-tx1", CreateDID: "updated-did1"}
	err1 := UpdateVCTemplate(db.UTchainID, template1)
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}
}

func TestGetVCTempListAndCount(t *testing.T) {
	insertVCTemplateTest()
	// Test case 1: Get VC template list and count
	offset1 := 0
	limit1 := 10
	contractAddr1 := ""
	templateID1 := ""
	_, _, err1 := GetVCTempListAndCount(offset1, limit1, db.UTchainID, contractAddr1, templateID1)
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}

	// Test case 2: Get VC template list and count with non-existing contract address
	offset2 := 0
	limit2 := 10
	contractAddr2 := "non-existing-address"
	templateID2 := ""
	_, _, err2 := GetVCTempListAndCount(offset2, limit2, db.UTchainID, contractAddr2, templateID2)
	if err2 != nil {
		t.Errorf("Test case 2 failed: %v", err2)
	}
}
