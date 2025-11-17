/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"chainmaker_web/src/db"
)

// 插入合约跨合约调用
func BatchInsertContractCrossCalls(chainId string, crossContractCalls []*db.ContractCrossCall) error {
	tableName := db.GetTableName(chainId, db.TableContractCrossCall)
	return CreateInBatchesData(tableName, crossContractCalls)
}

func GetContractCrossCallsByName(chainId, contractName string) ([]*db.ContractCrossCall, error) {
	results := make([]*db.ContractCrossCall, 0)
	if chainId == "" || contractName == "" {
		return results, nil
	}

	conditions := []QueryCondition{
		{
			Field:     "invokingContract",
			Value:     contractName,
			Condition: "=",
			Operator:  "or",
		},
		{
			Field:     "targetContract",
			Value:     contractName,
			Condition: "=",
			Operator:  "or",
		},
	}
	tableName := db.GetTableName(chainId, db.TableContractCrossCall)
	query := db.GormDB.Table(tableName)
	query = BuildQueryNew(query, conditions)
	err := query.Find(&results).Error
	return results, err
}
