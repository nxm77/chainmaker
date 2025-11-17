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

// InsertContractByteCode 插入合约字节码
// @param chainId 链id
// @param byteCode 合约字节码
// @return error
func InsertContractByteCodes(chainId string, byteCode []*db.ContractByteCode) error {
	if len(byteCode) == 0 {
		return nil
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableContractByteCode)
	return CreateInBatchesData(tableName, byteCode)
}

// GetContractByteCodeByTx 根据交易id获取合约字节码
// @param chainId 链id
// @param txId 交易id
// @return *db.ContractByteCode 合约字节码
// @return error
func GetContractByteCodeByTx(chainId, txId string) (*db.ContractByteCode, error) {
	byteCode := &db.ContractByteCode{}
	if chainId == "" || txId == "" {
		return nil, nil
	}

	tableName := db.GetTableName(chainId, db.TableContractByteCode)
	where := map[string]interface{}{
		"txId": txId,
	}
	err := db.GormDB.Table(tableName).Where(where).First(&byteCode).Error
	if err == nil {
		return byteCode, nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return nil, err
}
