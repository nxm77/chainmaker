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

func TestInsertIDAContract(t *testing.T) {
	// Test case 1: Normal case with valid data
	idaContracts1 := []*db.IDAContract{{ContractAddr: "addr1", ContractName: "name1"}}
	err1 := InsertIDAContract(ChainID, idaContracts1)
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}

	// Test case 3: Empty idaContracts
	idaContracts3 := []*db.IDAContract{}
	err3 := InsertIDAContract(ChainID, idaContracts3)
	if err3 != nil {
		t.Errorf("Test case 3 failed: %v", err3)
	}
}

func TestGetIDAContractMapByAddrs(t *testing.T) {
	// Test case 1: Normal case with valid data
	addrList1 := []string{"addr1", "addr2"}
	_, err1 := GetIDAContractMapByAddrs(ChainID, addrList1)
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}

	// Test case 2: Empty addrList
	addrList2 := []string{}
	_, err2 := GetIDAContractMapByAddrs(ChainID, addrList2)
	if err2 != nil {
		t.Errorf("Test case 2 failed: %v", err2)
	}
}

func TestUpdateIDAContractByAddr(t *testing.T) {
	// Test case 1: Normal case with valid data
	idaContract1 := &db.IDAContract{TotalNormalAssets: 100, TotalAssets: 200, BlockHeight: 10}
	err1 := UpdateIDAContractByAddr(ChainID, ContractAddr1, idaContract1)
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}

	// Test case 2: Empty params
	idaContract2 := &db.IDAContract{}
	err2 := UpdateIDAContractByAddr(ChainID, ContractAddr1, idaContract2)
	if err2 != nil {
		t.Errorf("Test case 2 failed: %v", err2)
	}
}

func TestGetIDAContractList(t *testing.T) {
	// Test case 1: Normal case with valid data
	offset1 := 0
	limit1 := 10
	_, _, err1 := GetIDAContractList(offset1, limit1, ChainID, "123", "desc", "timestamp")
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}
}

func TestGetIDAContractByAddr(t *testing.T) {
	// Test case 1: Normal case with valid data
	_, err1 := GetIDAContractByAddr(ChainID, ContractAddr1)
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}
}
