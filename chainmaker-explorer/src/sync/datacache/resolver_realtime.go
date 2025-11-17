/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package datacache

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/sync/model"
	"context"
	"encoding/json"

	loggers "chainmaker_web/src/logger"
)

var (
	log = loggers.GetLogger(loggers.MODULE_SYNC)
)

// SetDelayedUpdateCache
//
//	@Description: 设置数据缓存，异步计算使用
//	@param chainId
//	@param blockHeight
//	@param dealResult
func SetDelayedUpdateCache(chainId string, blockHeight int64, dealResult model.ProcessedBlockData) {
	contractAddrs := make(map[string]string, 0)
	delayedUpdateData := model.NewGetRealtimeCacheData()
	if len(dealResult.Transactions) == 0 {
		return
	}

	for _, txInfo := range dealResult.Transactions {
		if txInfo.ContractAddr == "" {
			continue
		}
		contractAddrs[txInfo.ContractAddr] = txInfo.ContractAddr
	}
	for _, event := range dealResult.ContractEvents {
		if event.ContractAddr == "" {
			continue
		}
		contractAddrs[event.ContractAddr] = event.ContractAddr
	}
	delayedUpdateData.TxList = dealResult.Transactions
	delayedUpdateData.GasRecords = dealResult.GasRecordList
	delayedUpdateData.ContractEvents = dealResult.ContractEvents
	delayedUpdateData.UserInfoMap = dealResult.UserList
	delayedUpdateData.ContractAddrs = contractAddrs
	delayedUpdateData.CrossChainResult = dealResult.CrossChainResult

	redisKey, expiration := cache.GetKeyDelayedUpdateData(chainId, blockHeight)
	retJson, _ := json.Marshal(delayedUpdateData)
	// 设置键值对和过期时间
	ctx := context.Background()
	_ = cache.GlobalRedisDb.Set(ctx, redisKey, string(retJson), expiration).Err()
}

// GetDelayedUpdateCache
//
//	@Description: 获取数据缓存，异步计算使用
//	@param chainId
//	@param blockHeight
//	@return *GetRealtimeCacheData 缓存数据
func GetDelayedUpdateCache(chainId string, blockHeight int64) *model.GetRealtimeCacheData {
	ctx := context.Background()
	redisKey, _ := cache.GetKeyDelayedUpdateData(chainId, blockHeight)
	redisRes := cache.GlobalRedisDb.Get(ctx, redisKey)
	if redisRes == nil || redisRes.Val() == "" {
		return nil
	}

	cacheResult := &model.GetRealtimeCacheData{}
	err := json.Unmarshal([]byte(redisRes.Val()), &cacheResult)
	if err != nil {
		return nil
	}
	return cacheResult
}

// GetRealtimeDataCache 获取实时数据缓存
// @param chainId 链ID
// @param height 区块高度
// @param realtimeCacheData 实时数据缓存
// @return bool 是否获取成功
func GetRealtimeDataCache(chainId string, height int64, realtimeCacheData *model.GetRealtimeCacheData) bool {
	// 获取缓存数据
	cacheResult := GetDelayedUpdateCache(chainId, height)
	if cacheResult == nil {
		return false
	}

	// 合并缓存数据到 delayedUpdateData
	for k, v := range cacheResult.TxList {
		realtimeCacheData.TxList[k] = v
	}
	// 合并缓存数据到 delayedUpdateData
	for userAddr, user := range cacheResult.UserInfoMap {
		realtimeCacheData.UserInfoMap[userAddr] = user
	}

	for _, addr := range cacheResult.ContractAddrs {
		realtimeCacheData.ContractAddrs[addr] = addr
	}
	if cacheResult.CrossChainResult != nil &&
		len(cacheResult.CrossChainResult.InsertCrossTransfer) > 0 {
		realtimeCacheData.CrossChainResult.InsertCrossTransfer = append(
			realtimeCacheData.CrossChainResult.InsertCrossTransfer,
			cacheResult.CrossChainResult.InsertCrossTransfer...,
		)
	}
	if cacheResult.CrossChainResult != nil && len(cacheResult.CrossChainResult.UpdateCrossTransfer) > 0 {
		realtimeCacheData.CrossChainResult.UpdateCrossTransfer = append(
			realtimeCacheData.CrossChainResult.UpdateCrossTransfer,
			cacheResult.CrossChainResult.UpdateCrossTransfer...,
		)
	}

	realtimeCacheData.GasRecords = append(realtimeCacheData.GasRecords, cacheResult.GasRecords...)
	realtimeCacheData.ContractEvents = append(realtimeCacheData.ContractEvents, cacheResult.ContractEvents...)
	return true
}
