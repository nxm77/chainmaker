/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package logic

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/sync/model"
	"encoding/json"

	"chainmaker.org/chainmaker/contract-utils/standard"
)

func processIDAEvent(event *db.ContractEvent, eventData []string, idaEventResult *model.IDAEventData) {
	idaCreatedMap := idaEventResult.IDACreatedMap
	idaUpdatedMap := idaEventResult.IDAUpdatedMap
	idaDeletedCodeMap := idaEventResult.IDADeletedCodeMap
	idaInfos, idaUpdateData, idaDeleteIds := BuildIDAEventData(event.ContractType, event.Topic, eventData)
	if len(idaInfos) > 0 {
		createdInfos := make([]*db.IDACreatedInfo, 0)
		for _, idaInfo := range idaInfos {
			createdInfo := &db.IDACreatedInfo{
				IDAInfo:      idaInfo,
				ContractAddr: event.ContractAddr,
				EventTime:    event.Timestamp,
			}
			createdInfos = append(createdInfos, createdInfo)
		}

		// 检查并初始化切片
		if _, ok := idaCreatedMap[event.ContractAddr]; !ok {
			idaCreatedMap[event.ContractAddr] = make([]*db.IDACreatedInfo, 0)
		}
		idaCreatedMap[event.ContractAddr] = append(idaCreatedMap[event.ContractAddr], createdInfos...)
	}

	if idaUpdateData != nil {
		idaUpdateData.EventTime = event.Timestamp
		if _, ok := idaUpdatedMap[idaUpdateData.IDACode]; !ok {
			idaUpdatedMap[idaUpdateData.IDACode] = make([]*db.EventIDAUpdatedData, 0)
		}
		idaUpdatedMap[idaUpdateData.IDACode] = append(idaUpdatedMap[idaUpdateData.IDACode], idaUpdateData)
	}

	if len(idaDeleteIds) > 0 {
		// 检查并初始化切片
		if _, ok := idaDeletedCodeMap[event.ContractAddr]; !ok {
			idaDeletedCodeMap[event.ContractAddr] = make([]string, 0)
		}
		idaDeletedCodeMap[event.ContractAddr] = append(idaDeletedCodeMap[event.ContractAddr], idaDeleteIds...)
		idaEventResult.EventTime = event.Timestamp
	}
}

// BuildIDAEventData 解析IDA数据
// @param contractType 合约类型
// @param topic 事件类型
// @param eventData 事件数据
// @return idaInfoList 解析后的IDA信息
// @return idaUpdateData 解析后的IDA更新信息
// @return idaIds 解析后的IDA删除信息
func BuildIDAEventData(contractType, topic string, eventData []string) (
	[]*standard.IDAInfo, *db.EventIDAUpdatedData, []string) {
	idaIds := make([]string, 0)
	idaInfoList := make([]*standard.IDAInfo, 0)
	var idaUpdateData *db.EventIDAUpdatedData
	//判断是否是IDA合约
	if contractType != standard.ContractStandardNameCMIDA {
		return idaInfoList, idaUpdateData, idaIds
	}

	switch topic {
	case standard.EventIDACreated:
		idaInfoList = DealEventIDACreated(eventData)
	case standard.EventIDAUpdated:
		idaUpdateData = DealEventIDAUpdated(eventData)
	case standard.EventIDADeleted:
		idaIds = DealEventIDADeleted(eventData)
	}

	return idaInfoList, idaUpdateData, idaIds
}

// DealEventIDACreated 解析IDA event
func DealEventIDACreated(eventData []string) []*standard.IDAInfo {
	//standard.EventIDACreated
	idaInfoList := make([]*standard.IDAInfo, 0)
	if len(eventData) == 0 {
		return idaInfoList
	}
	err := json.Unmarshal([]byte(eventData[0]), &idaInfoList)
	if err != nil {
		log.Errorf("DealEventIDACreated json Unmarshal err, err:%v, eventData:%v", err, eventData)
	}
	return idaInfoList
}

// DealEventIDAUpdated 解析IDA event
func DealEventIDAUpdated(eventData []string) *db.EventIDAUpdatedData {
	updateData := &db.EventIDAUpdatedData{}
	if len(eventData) < 3 {
		return updateData
	}
	updateData.IDACode = eventData[0]
	updateData.Field = eventData[1]
	updateData.Update = eventData[2]
	return updateData
}

// DealEventIDADeleted 解析IDA event
func DealEventIDADeleted(eventData []string) []string {
	idaCodes := make([]string, 0)
	if len(eventData) == 0 {
		return idaCodes
	}
	idaCodes = eventData
	return idaCodes
}

// DealEventIDACreated 解析IDA event
func UnmarshalIDAUpdatedBasic(updateJson string) (*standard.Basic, error) {
	basicInfo := &standard.Basic{}
	if updateJson == "" {
		return nil, nil
	}
	err := json.Unmarshal([]byte(updateJson), &basicInfo)
	if err != nil {
		log.Errorf("UnmarshalIDAUpdatedBasic json Unmarshal err, err:%v, updateJson:%v", err, updateJson)
	}
	return basicInfo, err
}

func UnmarshalIDAUpdatedSupply(updateJson string) (*standard.Supply, error) {
	supplyInfo := &standard.Supply{}
	if updateJson == "" {
		return supplyInfo, nil
	}
	err := json.Unmarshal([]byte(updateJson), &supplyInfo)
	if err != nil {
		log.Errorf("UnmarshalIDAUpdatedSupply json Unmarshal err, err:%v, updateJson:%v", err, updateJson)
	}
	return supplyInfo, err
}

func UnmarshalIDAUpdatedDetails(updateJson string) (*standard.Details, error) {
	detailInfo := &standard.Details{}
	if updateJson == "" {
		return detailInfo, nil
	}
	err := json.Unmarshal([]byte(updateJson), &detailInfo)
	if err != nil {
		log.Errorf("UnmarshalIDAUpdatedDetails json Unmarshal err, err:%v, updateJson:%v", err, updateJson)
	}
	return detailInfo, err
}

func UnmarshalIDAUpdatedOwnership(updateJson string) (*standard.Ownership, error) {
	ownershipInfo := &standard.Ownership{}
	if updateJson == "" {
		return ownershipInfo, nil
	}
	err := json.Unmarshal([]byte(updateJson), &ownershipInfo)
	if err != nil {
		log.Errorf("UnmarshalIDAUpdatedOwnership json Unmarshal err, err:%v, updateJson:%v", err, updateJson)
	}
	return ownershipInfo, err
}

func UnmarshalIDAUpdatedColumns(updateJson string) ([]*standard.ColumnInfo, error) {
	columns := make([]*standard.ColumnInfo, 0)
	if updateJson == "" {
		return columns, nil
	}
	err := json.Unmarshal([]byte(updateJson), &columns)
	if err != nil {
		log.Errorf("UnmarshalIDAUpdatedColumns json Unmarshal err, err:%v, updateJson:%v", err, updateJson)
	}
	return columns, err
}

func UnmarshalIDAUpdatedApis(updateJson string) ([]*standard.APIInfo, error) {
	apis := make([]*standard.APIInfo, 0)
	if updateJson == "" {
		return apis, nil
	}
	err := json.Unmarshal([]byte(updateJson), &apis)
	if err != nil {
		log.Errorf("UnmarshalIDAUpdatedApis json Unmarshal err, err:%v, updateJson:%v", err, updateJson)
	}
	return apis, err
}
