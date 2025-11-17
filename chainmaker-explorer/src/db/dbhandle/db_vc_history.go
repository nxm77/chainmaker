/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"chainmaker_web/src/db"
)

func InsertVCHistorys(chainId string, vcHistorys []*db.VCIssueHistory) error {
	if len(vcHistorys) == 0 {
		return nil
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableVCIssueHistory)
	return CreateInBatchesData(tableName, vcHistorys)
}

func UpdateVCIssuerStatus(chainId, vcId string, status int) error {
	if chainId == "" || vcId == "" {
		return db.ErrTableParams
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableVCIssueHistory)
	where := map[string]interface{}{
		"vcId": vcId,
	}
	params := map[string]interface{}{
		"status": status,
	}

	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	return err
}

// GetVCListAndCount 获取vc列表
func GetVCListAndCount(offset, limit int, chainId, issuerDID, holderDID, templateID, vcID, contractAddr string) (
	[]*db.VCIssueHistory, int64, error) {
	vcList := make([]*db.VCIssueHistory, 0)
	if chainId == "" {
		return vcList, 0, db.ErrTableParams
	}

	where := map[string]interface{}{}
	if issuerDID != "" {
		where["issuerDID"] = issuerDID
	}
	if holderDID != "" {
		where["holderDID"] = holderDID
	}
	if templateID != "" {
		where["templateId"] = templateID
	}
	if vcID != "" {
		where["vcId"] = vcID
	}
	if contractAddr != "" {
		where["contractAddr"] = contractAddr
	}

	tableName := db.GetTableName(chainId, db.TableVCIssueHistory)
	query := db.GormDB.Table(tableName).Where(where)
	// 获取总数
	var total int64
	err := query.Count(&total).Error
	if err != nil {
		return vcList, 0, err
	}

	// 获取数据
	err = query.Order("timestamp desc").Offset(offset * limit).Limit(limit).Find(&vcList).Error
	if err != nil {
		return vcList, 0, err
	}

	return vcList, total, nil
}

// GetVCListAndCount 获取vc列表
func GetVCInfoById(chainId, vcId string) (*db.VCIssueHistory, error) {
	vcInfo := &db.VCIssueHistory{}
	if chainId == "" {
		return vcInfo, db.ErrTableParams
	}

	where := map[string]interface{}{
		"vcId": vcId,
	}
	tableName := db.GetTableName(chainId, db.TableVCIssueHistory)
	err := db.GormDB.Table(tableName).Where(where).First(&vcInfo).Error
	return vcInfo, err
}
