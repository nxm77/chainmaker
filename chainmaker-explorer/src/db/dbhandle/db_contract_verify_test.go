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

func TestInsertContractVerifyResult(t *testing.T) {
	// Test case 1: Insert a new contract verify result
	verifyResult1 := &db.ContractVerifyResult{
		VerifyId:        "verify1",
		ContractAddr:    "addr1",
		ContractVersion: "v1",
		CompilerPath:    "/path/to/compiler",
		ABI:             "abi1",
		MetaData:        "metadata1",
		CompilerVersion: "v1",
		OpenLicenseType: "type1",
		EvmVersion:      "v1",
		RunNum:          1,
	}
	err1 := InsertContractVerifyResult(ChainID, verifyResult1)
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}
}

func TestGetContractVerifyResult(t *testing.T) {
	// Test case 1: Get a contract verify result by chainId, contractAddr and version
	_, err1 := GetContractVerifyResult(ChainID, ContractAddr1, "v2")
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}
}

func TestUpdateContractVerifyResult(t *testing.T) {
	// Test case 1: Update a contract verify result
	verifyResult1 := &db.ContractVerifyResult{
		VerifyId:        "verify1",
		ContractAddr:    "addr1",
		ContractVersion: "v1",
		CompilerPath:    "/path/to/compiler",
		ABI:             "abi1",
		MetaData:        "metadata1",
		CompilerVersion: "v1",
		OpenLicenseType: "type1",
		EvmVersion:      "v1",
		RunNum:          1,
	}
	err1 := UpdateContractVerifyResult(ChainID, verifyResult1)
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}
}
