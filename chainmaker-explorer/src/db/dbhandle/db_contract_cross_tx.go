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
func BatchInsertContractCrossCallTxs(chainId string, crossContractCallTxs []*db.ContractCrossCallTransaction) error {
	if len(crossContractCallTxs) == 0 || chainId == "" {
		return nil
	}

	tableName := db.GetTableName(chainId, db.TableContractCrossCallTransaction)
	return InsertData(tableName, crossContractCallTxs)
}
