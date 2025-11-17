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

func InsertVCTemplate(chainId string, vcTemplate *db.VCTemplate) error {
	if vcTemplate == nil {
		return nil
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableVCTemplate)
	return CreateInBatchesData(tableName, vcTemplate)
}

func GetVCTemplateById(chainId, templateId, contractAddr string) (*db.VCTemplate, error) {
	template := &db.VCTemplate{}
	where := map[string]interface{}{
		"templateId":   templateId,
		"contractAddr": contractAddr,
	}
	tableName := db.GetTableName(chainId, db.TableVCTemplate)
	err := db.GormDB.Table(tableName).Where(where).First(&template).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return template, nil
}

// UpdateAccount 更新账户
func UpdateVCTemplate(chainId string, template *db.VCTemplate) error {
	if chainId == "" || template == nil {
		return db.ErrTableParams
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableVCTemplate)
	where := map[string]interface{}{
		"templateId":   template.TemplateID,
		"contractAddr": template.ContractAddr,
	}
	params := map[string]interface{}{
		"templateName": template.TemplateName,
		"shortName":    template.ShortName,
		"vcTypre":      template.VCType,
		"version":      template.Version,
		"template":     template.Template,
		"txId":         template.TxId,
	}

	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	return err
}

// GetVCTempListAndCount 获取vc模板列表
func GetVCTempListAndCount(offset, limit int, chainId, contractAddr, templateID string) (
	[]*db.VCTemplate, int64, error) {
	vcList := make([]*db.VCTemplate, 0)
	if chainId == "" {
		return vcList, 0, db.ErrTableParams
	}

	where := map[string]interface{}{}
	if contractAddr != "" {
		where["contractAddr"] = contractAddr
	}
	if templateID != "" {
		where["templateId"] = templateID
	}

	tableName := db.GetTableName(chainId, db.TableVCTemplate)
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
