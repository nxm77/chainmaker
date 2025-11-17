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

// GetUpgradeContractTxList
//
//	@Description: 获取版本更新交易列表
//	@param offset
//	@param limit
//	@param chainId
//	@param contractName
//	@param contractAddr
//	@param senders
//	@param runtimeType
//	@param status
//	@return []*db.UpgradeContractTransaction
//	@return int64
//	@return error
func GetUpgradeContractTxList(offset int, limit int, chainId string, contractName, contractAddr string,
	senders []string, runtimeType string, status int, startTime, endTime int64) (
	[]*db.UpgradeContractTransaction, int64, error) {
	var count int64
	txList := make([]*db.UpgradeContractTransaction, 0)
	tableName := db.GetTableName(chainId, db.TableContractUpgradeTransaction)
	conditions := []QueryCondition{
		{Field: "timestamp", Value: contractName, Condition: "desc"},
	}

	if contractName != "" {
		conditions = append(conditions, QueryCondition{Field: "contractNameBak", Value: contractName, Condition: "="})
	}
	if contractAddr != "" {
		conditions = append(conditions, QueryCondition{Field: "contractAddr", Value: contractAddr, Condition: "="})
	}
	if len(senders) > 0 {
		conditions = append(conditions, QueryCondition{Field: "sender", Value: senders, Condition: "in"})
	}

	if runtimeType != "" {
		conditions = append(conditions, QueryCondition{Field: "contractRuntimeType", Value: runtimeType, Condition: "="})
	}

	if startTime > 0 && endTime > 0 {
		conditions = append(conditions, QueryCondition{
			Field:     "timestamp",
			Value:     []int64{startTime, endTime},
			Condition: "between",
		})
	}

	query := BuildQuery(tableName, conditions)
	err := query.Count(&count).Error
	if err != nil {
		return nil, 0, fmt.Errorf("GetUpgradeContractTxList err, cause : %s", err.Error())
	}

	query = query.Offset(offset * limit).Limit(limit)
	err = query.Find(&txList).Error
	if err != nil {
		return nil, 0, fmt.Errorf("GetUpgradeContractTxList err, cause : %s", err.Error())
	}

	return txList, count, nil
}

// InsertUpgradeContractTx 新增或者更新交易
func InsertUpgradeContractTx(chainId string, transactions []*db.UpgradeContractTransaction) error {
	if len(transactions) == 0 {
		return nil
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableContractUpgradeTransaction)
	return CreateInBatchesData(tableName, transactions)
}

// UpdateUpgradeContractName 更新合约名称
func UpdateUpgradeContractName(chainId string, contract *db.Contract) error {
	if chainId == "" || contract == nil || contract.Addr == "" {
		return nil
	}

	where := map[string]interface{}{
		"contractAddr": contract.Addr,
	}
	params := map[string]interface{}{
		"contractName":    contract.Name,
		"contractNameBak": contract.NameBak,
	}
	tableName := db.GetTableName(chainId, db.TableContractUpgradeTransaction)
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	if err != nil {
		return err
	}

	return nil
}

// UpdateUpgradeContractName 更新合约名称
func GetContractVersions(chainId string, contractAddr string) ([]string, error) {
	versions := make([]string, 0)

	if chainId == "" || contractAddr == "" {
		return versions, nil
	}

	where := map[string]interface{}{
		"contractAddr": contractAddr,
	}

	tableName := db.GetTableName(chainId, db.TableContractUpgradeTransaction)
	err := db.GormDB.Table(tableName).Select("contractVersion").Where(where).
		Order("timestamp desc").Find(&versions).Error
	if err != nil {
		return versions, err
	}

	return versions, nil
}

// GetUpgradeContractInfo 获取合约升级信息
func GetUpgradeContractInfo(chainId, contractAddr, contractVersion string) (*db.UpgradeContractTransaction, error) {
	var upgradeContract *db.UpgradeContractTransaction
	if chainId == "" || contractAddr == "" {
		return nil, nil
	}

	where := map[string]interface{}{
		"contractAddr":    contractAddr,
		"contractVersion": contractVersion,
	}

	tableName := db.GetTableName(chainId, db.TableContractUpgradeTransaction)
	err := db.GormDB.Table(tableName).Where(where).First(&upgradeContract).Error
	if err != nil {
		log.Errorf("GetUpgradeContractInfo err : %s", err.Error())
		return upgradeContract, err
	}

	return upgradeContract, nil
}

// UpdateUpgradeContractVerifyStatus 更新合约验证状态
func UpdateUpgradeContractVerifyStatus(chainId, contractAddr, contractVersion string, verifyStatus int) error {
	if chainId == "" || contractAddr == "" {
		return nil
	}

	where := map[string]interface{}{
		"contractAddr":    contractAddr,
		"contractVersion": contractVersion,
	}

	params := map[string]interface{}{
		"verifyStatus": verifyStatus,
	}
	tableName := db.GetTableName(chainId, db.TableContractUpgradeTransaction)
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	if err != nil {
		return err
	}
	return nil
}
