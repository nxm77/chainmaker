/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"chainmaker_web/src/db"
)

func InsertDIDHistorys(chainId string, didHistorys []*db.DIDSetHistory) error {
	if len(didHistorys) == 0 {
		return nil
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableDIDSetHistory)
	return CreateInBatchesData(tableName, didHistorys)
}

// GetDIDHistoryAndCount 获取DID历史记录
func GetDIDHistoryAndCount(offset, limit int, chainId, did string) ([]*db.DIDSetHistory, int64, error) {
	didList := make([]*db.DIDSetHistory, 0)
	if chainId == "" {
		return didList, 0, db.ErrTableParams
	}

	where := map[string]interface{}{
		"did": did,
	}

	tableName := db.GetTableName(chainId, db.TableDIDSetHistory)
	query := db.GormDB.Table(tableName).Where(where)
	// 获取总数
	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return didList, 0, err
	}

	// 获取数据
	err = query.Order("timestamp desc").Offset(offset * limit).Limit(limit).Find(&didList).Error
	if err != nil {
		return didList, 0, err
	}

	return didList, total, nil
}
