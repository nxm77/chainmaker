/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"chainmaker_web/src/db"
)

func InsertDIDDetail(chainId string, didDetail *db.DIDDetail) error {
	if didDetail == nil {
		return nil
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableDIDDetail)
	return CreateInBatchesData(tableName, didDetail)
}

// GetDIDDetailById 根据ID获取DID详情
func GetDIDDetailById(chainId, did string) (*db.DIDDetail, error) {
	detail := &db.DIDDetail{}
	where := map[string]interface{}{
		"did": did,
	}
	tableName := db.GetTableName(chainId, db.TableDIDDetail)
	err := db.GormDB.Table(tableName).Where(where).First(&detail).Error
	if err != nil {
		return nil, err
	}

	return detail, nil
}

// UpdateAccount 更新账户
func UpdateDIDDetail(chainId string, didDetail *db.DIDDetail) error {
	if chainId == "" || didDetail == nil {
		return db.ErrTableParams
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableDIDDetail)
	where := map[string]interface{}{
		"did": didDetail.DID,
	}
	params := map[string]interface{}{
		"document":      didDetail.Document,
		"issuerService": didDetail.IssuerService,
		"accountJson":   didDetail.AccountJson,
		"contractName":  didDetail.ContractName,
		"contractAddr":  didDetail.ContractAddr,
	}

	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	return err
}

func UpdateDIDStatus(chainId, did string, status int) error {
	if chainId == "" || did == "" {
		return db.ErrTableParams
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableDIDDetail)
	where := map[string]interface{}{
		"did": did,
	}
	params := map[string]interface{}{
		"status": status,
	}
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	return err
}

func UpdateDIDIssuer(chainId, did string, isIssuer bool) error {
	if chainId == "" || did == "" {
		return db.ErrTableParams
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableDIDDetail)
	where := map[string]interface{}{
		"did": did,
	}
	params := map[string]interface{}{
		"isIssuer": isIssuer,
	}
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	return err
}

// GetDIDListAndCount 获取DID列表和总数
func GetDIDListAndCount(offset, limit int, chainId, contractAddr, did string) ([]*db.DIDDetail, int64, error) {
	didList := make([]*db.DIDDetail, 0)
	if chainId == "" || contractAddr == "" {
		return didList, 0, db.ErrTableParams
	}

	where := map[string]interface{}{
		"contractAddr": contractAddr,
	}
	if did != "" {
		where["did"] = did
	}

	tableName := db.GetTableName(chainId, db.TableDIDDetail)
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
