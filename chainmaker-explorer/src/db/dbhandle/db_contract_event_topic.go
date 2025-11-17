/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"chainmaker_web/src/db"
)

// InsertContractEventTopic 批量保存event
func InsertContractEventTopic(chainId string, inserts []*db.ContractEventTopic) error {
	if len(inserts) == 0 {
		return nil
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableContractEventTopic)
	return CreateInBatchesData(tableName, inserts)
}

// GetEventTopicByNames 根据ID获取event
func GetEventTopicByNames(chainId string, contractNames []string) ([]*db.ContractEventTopic, error) {
	eventTopicList := make([]*db.ContractEventTopic, 0)
	if len(contractNames) == 0 {
		return eventTopicList, nil
	}
	tableName := db.GetTableName(chainId, db.TableContractEventTopic)
	conditions := []QueryCondition{
		{Field: "contractName", Value: contractNames, Condition: "in"},
	}
	query := BuildQuery(tableName, conditions)
	err := query.Find(&eventTopicList).Error
	return eventTopicList, err
}

// UpdateContractEventTopic 更新更新状态
func UpdateContractEventTopic(chainId string, update *db.ContractEventTopic) error {
	if update == nil {
		return nil
	}

	where := map[string]interface{}{
		"contractName": update.ContractName,
		"topic":        update.Topic,
	}

	params := map[string]interface{}{}
	if update.TxNum > 0 {
		params["txNum"] = update.TxNum
	}
	if update.BlockHeight > 0 {
		params["blockHeight"] = update.BlockHeight
	}

	if len(params) == 0 {
		return nil
	}

	tableName := db.GetTableName(chainId, db.TableContractEventTopic)
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	return err
}

// GetContractEventTopic 根据ID获取event
func GetContractEventTopic(chainId, contractName, topic string) ([]*db.ContractEventTopic, error) {
	eventTopicList := make([]*db.ContractEventTopic, 0)
	if contractName == "" || chainId == "" {
		return eventTopicList, nil
	}
	tableName := db.GetTableName(chainId, db.TableContractEventTopic)
	conditions := []QueryCondition{
		{Field: "contractName", Value: contractName, Condition: "="},
	}
	if topic != "" {
		conditions = append(conditions, QueryCondition{Field: "topic", Value: "%" + topic + "%", Condition: "LIKE"})
	}
	query := BuildQuery(tableName, conditions)
	err := query.Find(&eventTopicList).Error
	return eventTopicList, err
}
