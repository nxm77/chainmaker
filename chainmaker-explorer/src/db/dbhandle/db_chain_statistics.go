/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"chainmaker_web/src/db"
)

func InsertStatistics(chainId string, stats *db.Statistics) error {
	if stats == nil {
		return nil
	}

	//获取交易表名称
	return CreateInBatchesData(db.TableStatistics, stats)
}

func GetChainStatistics(chainId string) (*db.Statistics, error) {
	var statistics *db.Statistics
	if chainId == "" {
		return nil, db.ErrTableParams
	}

	where := map[string]interface{}{
		"chainId": chainId,
	}
	err := db.GormDB.Table(db.TableStatistics).Where(where).First(&statistics).Error
	if err != nil {
		return nil, err
	}
	return statistics, nil
}

// UpdateStatistics 更新统计数据
func UpdateStatisticsDelay(chainId string, stats *db.Statistics) error {
	if chainId == "" || stats == nil {
		return nil
	}

	where := map[string]interface{}{
		"chainId": stats.ChainId,
	}
	params := map[string]interface{}{}
	if stats.TotalAccounts > 0 {
		params["totalAccounts"] = stats.TotalAccounts
	}
	if stats.TotalOrgs > 0 {
		params["totalOrgs"] = stats.TotalOrgs
	}
	if stats.TotalNodes > 0 {
		params["totalNodes"] = stats.TotalNodes
	}

	// 获取表名
	err := db.GormDB.Table(db.TableStatistics).Where(where).Updates(params).Error
	return err
}

func UpdateStatisticsRealtime(chainId string, stats *db.Statistics) error {
	if chainId == "" || stats == nil {
		return nil
	}

	where := map[string]interface{}{
		"chainId": stats.ChainId,
	}
	params := map[string]interface{}{}
	if stats.BlockHeight > 0 {
		params["blockHeight"] = stats.BlockHeight
	}
	if stats.TotalTransactions > 0 {
		params["totalTransactions"] = stats.TotalTransactions
	}
	if stats.TotalCrossTx > 0 {
		params["totalCrossTx"] = stats.TotalCrossTx
	}
	if stats.TotalContracts > 0 {
		params["totalContracts"] = stats.TotalContracts
	}

	// 获取表名
	err := db.GormDB.Table(db.TableStatistics).Where(where).Updates(params).Error
	return err
}

// DeleteSubscribe
//
//	@Description: 删除订阅
//	@param chainId
//	@return error
func DeleteStatistics(chainId string) error {
	statistics := &db.Statistics{}
	where := map[string]interface{}{
		"chainId": chainId,
	}
	err := db.GormDB.Table(db.TableStatistics).Where(where).Delete(&statistics).Error
	if err != nil {
		log.Error("[DB] Delete DeleteStatistics Failed: " + err.Error())
	}
	return err
}
