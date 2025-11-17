/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"chainmaker_web/src/db"

	"gorm.io/gorm"
)

// InsertContractABI 插入合约ABI
func InsertContractABITopic(chainId string, contractTopic *db.ContractABITopic) error {
	if contractTopic == nil || chainId == "" {
		return nil
	}

	tableName := db.GetTableName(chainId, db.TableContractABITopic)
	return InsertData(tableName, contractTopic)
}

// GetContractABIJson 获取合约ABI
func GetContractABITopic(chainId, contractAddr, version, topic string) (*db.ContractABITopic, error) {
	if contractAddr == "" || chainId == "" || version == "" || topic == "" {
		return nil, db.ErrTableParams
	}

	tableName := db.GetTableName(chainId, db.TableContractABITopic)
	abiTopic := &db.ContractABITopic{}
	where := map[string]interface{}{
		"contractAddr":    contractAddr,
		"contractVersion": version,
		"topic":           topic,
	}
	err := db.GormDB.Table(tableName).Where(where).First(&abiTopic).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return abiTopic, nil
}

func GetContractABITopics(chainId, contractAddr, topic string) ([]*db.ContractABITopic, error) {
	if contractAddr == "" || chainId == "" {
		return nil, db.ErrTableParams
	}

	tableName := db.GetTableName(chainId, db.TableContractABITopic)
	abiTopics := make([]*db.ContractABITopic, 0)
	where := map[string]interface{}{
		"contractAddr": contractAddr,
	}
	if topic != "" {
		where["topic"] = topic
	}
	err := db.GormDB.Table(tableName).Where(where).Find(&abiTopics).Error
	if err != nil {
		return nil, err
	}

	return abiTopics, nil
}

func GetContractABITopicByAddrs(chainId string, contractAddrs []string) ([]*db.ContractABITopic, error) {
	if len(contractAddrs) == 0 || chainId == "" {
		return nil, nil
	}

	tableName := db.GetTableName(chainId, db.TableContractABITopic)
	abiTopics := make([]*db.ContractABITopic, 0)
	conditions := []QueryCondition{
		{Field: "contractAddr", Value: contractAddrs, Condition: "in"},
	}
	query := BuildQuery(tableName, conditions)
	err := query.Find(&abiTopics).Error
	if err != nil {
		return nil, err
	}
	return abiTopics, nil
}

// UpdateContractABI 更新更新状态
func UpdateContractABITopic(chainId string, contractABI *db.ContractABITopic) error {
	tableName := db.GetTableName(chainId, db.TableContractABITopic)
	where := map[string]interface{}{
		"contractAddr":    contractABI.ContractAddr,
		"contractVersion": contractABI.ContractVersion,
		"topic":           contractABI.Topic,
	}
	params := map[string]interface{}{
		"topicTableName": contractABI.TopicTableName,
	}
	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	return err
}
