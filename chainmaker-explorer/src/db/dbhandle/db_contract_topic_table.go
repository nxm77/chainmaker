/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"chainmaker_web/src/db"
	"fmt"
)

// 修改后的 InsertDecodeEventByABI
func InsertDecodeEventByABI(chainId, tableName string, events []map[string]interface{}) error {
	tableName = db.GetTableName(chainId, tableName)

	err := db.DBHandler.InsertWithNativeSQL(tableName, events)
	//检查错误是否为主键冲突或唯一索引冲突
	isDuplicate := isDuplicateKeyError(err)
	if isDuplicate {
		return nil
	}

	return err
}

// DeleteTopicDataRecord 删除指定版本的数据记录
func DeleteTopicDataRecord(chainId, version, tableName string) {
	// 如果版本号为空，则返回nil
	if version == "" {
		return
	}
	type EmptyModel struct{} // 空结构体
	tableName = db.GetTableName(chainId, tableName)
	// 定义where条件，指定版本号
	where := map[string]interface{}{
		db.ABISystemFieldContractVer: version,
	}

	// 使用GormDB删除指定表中的数据记录
	err := db.GormDB.Table(tableName).Where(where).Delete(EmptyModel{}).Error
	if err != nil {
		log.Errorf("DeleteTopicDataRecord Error deleting records from table %s: %v", tableName, err)
		return
	}

	// 检查表是否还有数据
	var count int64
	err = db.GormDB.Table(tableName).Count(&count).Error
	if err != nil {
		log.Errorf("DeleteTopicDataRecord Error checking table %s: %v", tableName, err)
		return
	}

	// 如果表数据为空，则删除表
	if count == 0 {
		err = db.GormDB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)).Error
		log.Errorf("DeleteTopicDataRecord Error dropping table %s: %v", tableName, err)
	}
}
