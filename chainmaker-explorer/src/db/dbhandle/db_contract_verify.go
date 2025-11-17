/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"chainmaker_web/src/db"
	"errors"

	"gorm.io/gorm"
)

// InsertContractVerifyResult 批量插入gas
func InsertContractVerifyResult(chainId string, verify *db.ContractVerifyResult) error {
	if verify == nil {
		return nil
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableContractVerifyResult)
	return CreateInBatchesData(tableName, verify)
}

// GetGasByAddrInfo 根据多个addr获取Gas余额
func GetContractVerifyResult(chainId, contractAddr, version string) (*db.ContractVerifyResult, error) {
	verifyResult := &db.ContractVerifyResult{}
	if chainId == "" || contractAddr == "" || version == "" {
		return nil, nil
	}

	tableName := db.GetTableName(chainId, db.TableContractVerifyResult)
	where := map[string]interface{}{
		"contractAddr":    contractAddr,
		"contractVersion": version,
	}
	err := db.GormDB.Table(tableName).Where(where).First(&verifyResult).Error
	if err == nil {
		return verifyResult, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return nil, err
}

// UpdateContractVerifyResult
func UpdateContractVerifyResult(chainId string, verifyResult *db.ContractVerifyResult) error {
	if chainId == "" || verifyResult == nil {
		return nil
	}

	where := map[string]interface{}{
		"verifyId": verifyResult.VerifyId,
	}
	params := map[string]interface{}{
		"verifyStatus":    verifyResult.VerifyStatus,
		"compilerPath":    verifyResult.CompilerPath,
		"byteCode":        verifyResult.ByteCode,
		"abi":             verifyResult.ABI,
		"metaData":        verifyResult.MetaData,
		"compilerVersion": verifyResult.CompilerVersion,
		"openLicenseType": verifyResult.OpenLicenseType,
		"evmVersion":      verifyResult.EvmVersion,
		"optimization":    verifyResult.Optimization,
		"runNum":          verifyResult.RunNum,
	}
	// 获取表名
	tableName := db.GetTableName(chainId, db.TableContractVerifyResult)
	err := db.GormDB.Table(tableName).Model(&db.ContractVerifyResult{}).Where(where).Updates(params).Error
	return err
}
