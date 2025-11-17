/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"chainmaker_web/src/db"

	"gorm.io/gorm"
)

// InsertContractABI 插入合约ABI
func InsertContractABIFile(chainId string, contractABI *db.ContractABIFile) error {
	if contractABI == nil || chainId == "" {
		return nil
	}

	tableName := db.GetTableName(chainId, db.TableContractABIFile)
	return InsertData(tableName, contractABI)
}

// GetContractABIJson 获取合约ABI
func GetContractABIFile(chainId, contractAddr, version string) (*db.ContractABIFile, error) {
	if contractAddr == "" || chainId == "" || version == "" {
		return nil, db.ErrTableParams
	}

	tableName := db.GetTableName(chainId, db.TableContractABIFile)
	abiFile := &db.ContractABIFile{}
	where := map[string]interface{}{
		"contractAddr":    contractAddr,
		"contractVersion": version,
	}
	err := db.GormDB.Table(tableName).Where(where).First(&abiFile).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return abiFile, nil
}

func GetContractABIFiles(chainId, contractAddr string) ([]*db.ContractABIFile, error) {
	if contractAddr == "" || chainId == "" {
		return nil, db.ErrTableParams
	}

	tableName := db.GetTableName(chainId, db.TableContractABIFile)
	abiFiles := make([]*db.ContractABIFile, 0)
	where := map[string]interface{}{
		"contractAddr": contractAddr,
	}
	err := db.GormDB.Table(tableName).Where(where).Find(&abiFiles).Error
	if err != nil {
		return nil, err
	}

	return abiFiles, nil
}

// UpdateContractABI 更新更新状态
func UpdateContractABIFile(chainId string, contractABI *db.ContractABIFile) error {
	tableName := db.GetTableName(chainId, db.TableContractABIFile)
	where := map[string]interface{}{
		"contractAddr":    contractABI.ContractAddr,
		"contractVersion": contractABI.ContractVersion,
	}
	params := map[string]interface{}{
		"abiJson": contractABI.ABIJson,
	}
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	return err
}
