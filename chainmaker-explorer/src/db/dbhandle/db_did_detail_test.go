/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"chainmaker_web/src/db"
	"testing"

	"github.com/google/uuid"
)

var DID1Test = "did1"

func TestInsertDIDDetail(t *testing.T) {
	// Test case 1: Insert a new DID detail
	didDetail1 := &db.DIDDetail{DID: DID1Test, Document: "document1", IssuerService: "service1", AccountJson: "json1", ContractName: "contract1", ContractAddr: "addr1", Status: 1, IsIssuer: true}
	err1 := InsertDIDDetail(db.UTchainID, didDetail1)
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}

	// Test case 2: Insert a DID detail with empty fields
	didDetail2 := &db.DIDDetail{}
	err2 := InsertDIDDetail(db.UTchainID, didDetail2)
	if err2 != nil {
		t.Errorf("Test case 2 failed: %v", err2)
	}

	// Test case 3: Insert a DID detail with nil pointer
	err3 := InsertDIDDetail(db.UTchainID, nil)
	if err3 != nil {
		t.Errorf("Test case 3 failed: %v", err3)
	}
}

func TestGetDIDDetailById(t *testing.T) {
	didDetail1 := &db.DIDDetail{DID: DID1Test, Document: "document1", IssuerService: "service1", AccountJson: "json1", ContractName: "contract1", ContractAddr: "addr1", Status: 1, IsIssuer: true}
	InsertDIDDetail(db.UTchainID, didDetail1)
	// Test case 4: Get DID detail by ID

	_, err4 := GetDIDDetailById(db.UTchainID, DID1Test)
	if err4 != nil {
		t.Errorf("Test case 4 failed: %v", err4)
	}

	// Test case 5: Get DID detail by non-existing ID
	did5 := "non-existing-did"
	_, err5 := GetDIDDetailById(db.UTchainID, did5)
	if err5 == nil {
		t.Errorf("Test case 5 failed: expected error, got nil")
	}
}

func TestUpdateDIDDetail(t *testing.T) {
	// Test case 6: Update existing DID detail
	didDetail6 := &db.DIDDetail{DID: DID1Test, Document: "updated-document", IssuerService: "updated-service", AccountJson: "updated-json", ContractName: "updated-contract", ContractAddr: "updated-addr", Status: 2, IsIssuer: false}
	err6 := UpdateDIDDetail(db.UTchainID, didDetail6)
	if err6 != nil {
		t.Errorf("Test case 6 failed: %v", err6)
	}

	// Test case 8: Update DID detail with empty fields
	didDetail8 := &db.DIDDetail{DID: DID1Test}
	err8 := UpdateDIDDetail(db.UTchainID, didDetail8)
	if err8 != nil {
		t.Errorf("Test case 8 failed: %v", err8)
	}
}

func TestUpdateDIDStatus(t *testing.T) {
	// Test case 10: Update status of existing DID
	status10 := 3
	err10 := UpdateDIDStatus(db.UTchainID, DID1Test, status10)
	if err10 != nil {
		t.Errorf("Test case 10 failed: %v", err10)
	}

	// Test case 13: Update status with empty DID
	did13 := ""
	status13 := 3
	err13 := UpdateDIDStatus(db.UTchainID, did13, status13)
	if err13 == nil {
		t.Errorf("Test case 13 failed: expected error, got nil")
	}
}

func TestUpdateDIDIssuer(t *testing.T) {
	// Test case 14: Update issuer status of existing DID
	isIssuer14 := false
	err14 := UpdateDIDIssuer(db.UTchainID, DID1Test, isIssuer14)
	if err14 != nil {
		t.Errorf("Test case 14 failed: %v", err14)
	}

	// Test case 17: Update issuer status with empty DID
	did17 := ""
	isIssuer17 := false
	err17 := UpdateDIDIssuer(db.UTchainID, did17, isIssuer17)
	if err17 == nil {
		t.Errorf("Test case 17 failed: expected error, got nil")
	}
}

func TestGetDIDListAndCount(t *testing.T) {
	// Test case 18: Get DID list and count
	contractAddr18 := "addr1"
	offset18 := 0
	limit18 := 10
	did18 := ""
	_, _, err18 := GetDIDListAndCount(offset18, limit18, db.UTchainID, contractAddr18, did18)
	if err18 != nil {
		t.Errorf("Test case 18 failed: %v", err18)
	}
}

func TestGetDIDHistoryAndCount(t *testing.T) {
	_, _, err4 := GetDIDHistoryAndCount(0, 2, db.UTchainID, "didd")
	if err4 != nil {
		t.Errorf("Test case 4 failed: %v", err4)
	}
}

func TestInsertDIDHistorys(t *testing.T) {
	newUUID := uuid.New().String()
	didHistorys := []*db.DIDSetHistory{
		{
			ID:           newUUID,
			DID:          newUUID,
			ContractName: "contractName",
			ContractAddr: "contractAddr",
		},
	}
	err := InsertDIDHistorys(db.UTchainID, didHistorys)
	if err != nil {
		t.Errorf("Test case 1 failed: %v", err)
	}
}
