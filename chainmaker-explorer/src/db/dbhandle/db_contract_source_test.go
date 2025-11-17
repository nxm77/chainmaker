/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"chainmaker_web/src/db"
	"testing"

	"chainmaker.org/chainmaker/common/v2/random/uuid"
)

var VerifyId1 = "verify1"

func TestInsertContractSource(t *testing.T) {
	// Test case 1: Normal case with multiple contract source files
	sourceFile1 := []*db.ContractSourceFile{
		{
			ID:              uuid.GetUUID(),
			VerifyId:        "verify1",
			ContractAddr:    "addr1",
			ContractVersion: "v1",
		},
	}
	err1 := InsertContractSource(ChainID, sourceFile1)
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}

	err2 := InsertContractSource(ChainID, nil)
	if err2 != nil {
		t.Errorf("Test case 2 failed: %v", err2)
	}
}

func TestGetContractSourceFile(t *testing.T) {
	// Test case 1: Normal case with existing contract source files
	_, err1 := GetContractSourceFile(ChainID, VerifyId1)
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}
	_, err2 := GetContractSourceFile(ChainID, "")
	if err2 != nil {
		t.Errorf("Test case 1 failed: %v", err2)
	}
}
