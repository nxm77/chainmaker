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
	"github.com/stretchr/testify/assert"
)

const (
	ContractVersionUT = "v1.0"
)

func TestInsertContractABIFile(t *testing.T) {
	newUUID := uuid.New().String()
	insertList := &db.ContractABIFile{
		Id:              newUUID,
		ContractName:    ContractName1,
		ContractAddr:    ContractAddr1,
		ContractVersion: ContractVersionUT,
		ABIJson:         "",
	}

	err := InsertContractABIFile(db.UTchainID, insertList)
	assert.NoError(t, err)
}

func TestGetContractABIFile(t *testing.T) {
	TestInsertContractABIFile(t)
	contractABIFile, err := GetContractABIFile(db.UTchainID, ContractAddr1, ContractVersionUT)
	assert.NoError(t, err)
	assert.Equal(t, ContractName1, contractABIFile.ContractName)
	assert.Equal(t, ContractAddr1, contractABIFile.ContractAddr)
}

func TestGetContractABIFiles(t *testing.T) {
	TestInsertContractABIFile(t)
	contractABIFiles, err := GetContractABIFiles(db.UTchainID, ContractAddr1)
	assert.NoError(t, err)
	assert.Greater(t, len(contractABIFiles), 0, "Expected contractABIFiles to have more than 0 items")
}

func TestUpdateContractABIFile(t *testing.T) {
	TestInsertContractABIFile(t)
	contractABIFile := &db.ContractABIFile{
		ContractName:    ContractName1,
		ContractAddr:    ContractAddr1,
		ContractVersion: ContractVersionUT,
		ABIJson:         "",
	}
	err := UpdateContractABIFile(db.UTchainID, contractABIFile)
	assert.NoError(t, err)
}
