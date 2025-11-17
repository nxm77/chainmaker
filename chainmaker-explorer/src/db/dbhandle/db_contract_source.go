/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"chainmaker_web/src/db"
)

// InsertContractSource 批量插入gas
func InsertContractSource(chainId string, sourceFile []*db.ContractSourceFile) error {
	if len(sourceFile) == 0 {
		return nil
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableContractSourceFile)
	return CreateInBatchesData(tableName, sourceFile)
}

// GetContractSourceFile 获取合约源码文件
func GetContractSourceFile(chainId string, verifyId string) ([]*db.ContractSourceFile, error) {
	sourceFiles := make([]*db.ContractSourceFile, 0)
	if chainId == "" || verifyId == "" {
		return sourceFiles, nil
	}

	where := map[string]interface{}{
		"verifyId": verifyId,
	}

	tableName := db.GetTableName(chainId, db.TableContractSourceFile)
	err := db.GormDB.Table(tableName).Where(where).Find(&sourceFiles).Error
	if err != nil {
		return sourceFiles, err
	}

	return sourceFiles, nil
}
