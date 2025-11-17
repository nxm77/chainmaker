/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package logic

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/utils"
	"encoding/json"
	"fmt"
)

func GetTopicABIInputs(chainId, contractAddr, contractVersion, topic string) []*db.DynamicStructField {
	var fields []*db.DynamicStructField
	// 1. 获取合约的ABI
	contractABI, err := dbhandle.GetContractABIFile(chainId, contractAddr, contractVersion)
	if err != nil || contractABI == nil {
		return fields
	}

	// 解析JSON内容到ContractABI结构体切片
	abiList := make([]*utils.ContractABI, 0)
	err = json.Unmarshal([]byte(contractABI.ABIJson), &abiList)
	if err != nil {
		return fields
	}

	// 2. 根据topic获取inputs
	for _, abi := range abiList {
		if abi.Type != utils.ABIEventType {
			continue
		}

		if abi.Name == topic {
			fields = make([]*db.DynamicStructField, len(abi.Inputs))
			for i, input := range abi.Inputs {
				fields[i] = &db.DynamicStructField{
					Name: input.Name,
					Type: input.Type,
				}
			}
			break
		}
	}
	return fields
}

// 根据ABI构建事件数据
func BuildEventDataByABI(chainId string, contractEvents []*db.ContractEvent) map[string][]map[string]interface{} {
	result := make(map[string][]map[string]interface{})

	// 如果contractEvents为空，直接返回空map
	if len(contractEvents) == 0 {
		return result
	}

	// 批量获取表名映射（优化版）
	// 创建一个 map 用于去重
	uniqueAddrs := make(map[string]struct{})
	for _, e := range contractEvents {
		uniqueAddrs[e.ContractAddr] = struct{}{}
	}

	// 将去重后的地址存入 slice
	contractAddrs := make([]string, 0, len(uniqueAddrs))
	for addr := range uniqueAddrs {
		contractAddrs = append(contractAddrs, addr)
	}

	topicInfos, err := dbhandle.GetContractABITopicByAddrs(chainId, contractAddrs)
	if err != nil {
		log.Errorf("BuildEventDataByABI 获取ABI主题信息失败: %v", err)
		return nil
	}

	// 构建表名快速映射
	tableNameMap := make(map[string]string) // key: contractAddr|topic|version
	for _, info := range topicInfos {
		key := fmt.Sprintf("%s_%s_%s", info.ContractAddr, info.Topic, info.ContractVersion)
		tableNameMap[key] = info.TopicTableName
	}

	// 处理每个事件
	for _, event := range contractEvents {
		// 获取表名（快速匹配）
		key := fmt.Sprintf("%s_%s_%s",
			event.ContractAddr,
			event.Topic,
			event.ContractVersion,
		)
		tableName := tableNameMap[key]
		if tableName == "" {
			continue
		}

		// 生成结构体实例
		record, err := BuildEventDataByABIInputs(chainId, event)
		if err != nil {
			log.Warnf("事件[%s]数据处理失败: table:%v, err: %v", event.TxId, tableName, err)
		} else {
			// 聚合结果
			result[tableName] = append(result[tableName], record)
		}
	}

	return result
}

func BuildEventDataByABIInputs(chainId string, contractEvent *db.ContractEvent) (map[string]interface{}, error) {
	// 获取事件字段列表
	eventInputs := GetTopicABIInputs(chainId, contractEvent.ContractAddr,
		contractEvent.ContractVersion, contractEvent.Topic)

	// 校验数据完整性
	var eventDataList []string
	if err := json.Unmarshal([]byte(contractEvent.EventData), &eventDataList); err != nil {
		log.Errorf("解析事件数据失败: %v", err)
		return nil, fmt.Errorf("解析事件数据失败: %v", err)
	}
	if len(eventInputs) != len(eventDataList) {
		return nil, fmt.Errorf("输入字段数(%d)与数据项数(%d)不匹配", len(eventInputs), len(eventDataList))
	}

	// 构建字段映射（显式类型转换）
	data := make(map[string]interface{})
	uniqueString := fmt.Sprintf("%s-%d", contractEvent.TxId, contractEvent.EventIndex)
	uniqueID := utils.CalculateSHA256([]byte(uniqueString))
	data[db.ABISystemFieldID] = uniqueID
	data[db.ABISystemFieldTxID] = contractEvent.TxId
	data[db.ABISystemFieldContractVer] = contractEvent.ContractVersion
	data[db.ABISystemFieldTimestamp] = contractEvent.Timestamp

	// 动态字段处理（示例：字符串类型）
	for i, input := range eventInputs {
		data[input.Name] = eventDataList[i]
	}
	return data, nil
}
