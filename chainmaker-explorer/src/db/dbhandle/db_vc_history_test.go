/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"testing"

	"chainmaker_web/src/db"
)

var VCTest = "vc1"

func insertVCHistorysTest() {
	// Test case 1: Normal case with valid input
	chainId1 := "chain1"
	vcHistorys1 := []*db.VCIssueHistory{
		{
			ID:         "1",
			VCID:       VCTest,
			IssuerDID:  "issuer1",
			HolderDID:  "holder1",
			TemplateID: "template1",
			Timestamp:  1234567890,
			Status:     0,
		},
		{
			ID:         "2",
			VCID:       "vc2",
			IssuerDID:  "issuer2",
			HolderDID:  "holder2",
			TemplateID: "template2",
			Timestamp:  1234567890,
			Status:     0,
		},
	}
	_ = InsertVCHistorys(chainId1, vcHistorys1)
}

func TestUpdateVCIssuerStatus(t *testing.T) {
	insertVCHistorysTest()
	// Test case 1: Normal case with valid input
	status1 := 1
	err1 := UpdateVCIssuerStatus(db.UTchainID, VCTest, status1)
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}

	// Test case 2: Invalid chainId
	chainId2 := ""
	status2 := 1
	err2 := UpdateVCIssuerStatus(chainId2, VCTest, status2)
	if err2 != db.ErrTableParams {
		t.Errorf("Test case 2 failed: %v", err2)
	}
}

func TestGetVCListAndCount(t *testing.T) {
	// Test case 1: Normal case with valid input
	insertVCHistorysTest()
	offset1 := 0
	limit1 := 10
	issuerDID1 := "issuer12"
	holderDID1 := "holder12"
	templateID1 := "template12"
	_, _, err1 := GetVCListAndCount(offset1, limit1, db.UTchainID, issuerDID1, holderDID1, templateID1, VCTest, "")
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}

	// Test case 2: Invalid chainId
	offset2 := 0
	limit2 := 10
	chainId2 := ""
	issuerDID2 := "issuer13"
	holderDID2 := "holder13"
	templateID2 := "template13"
	_, _, err2 := GetVCListAndCount(offset2, limit2, chainId2, issuerDID2, holderDID2, templateID2, VCTest, "")
	if err2 != db.ErrTableParams {
		t.Errorf("Test case 2 failed: %v", err2)
	}
}

func TestGetVCInfoById(t *testing.T) {
	// Test case 1: Normal case with valid input
	vcId1 := "vc2222221"
	_, err1 := GetVCInfoById(db.UTchainID, vcId1)
	if err1 == nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}

	// Test case 2: Invalid chainId
	chainId2 := ""
	_, err2 := GetVCInfoById(chainId2, VCTest)
	if err2 != db.ErrTableParams {
		t.Errorf("Test case 2 failed: %v", err2)
	}
}
