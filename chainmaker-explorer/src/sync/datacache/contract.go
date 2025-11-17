package datacache

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// SetLatestContractListCache
//
//	@Description: 缓存最新合约列表
//	@param chainId
//	@param blockHeight
//	@param insertContracts
//	@param updateContracts
func SetLatestContractListCache(chainId string, blockHeight int64, insertContracts, updateContracts []*db.Contract) {
	if len(insertContracts) == 0 && len(updateContracts) == 0 {
		return
	}

	ctx := context.Background()
	//添加缓存合约信息
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisLatestContractList, prefix, chainId)
	//新增合约缓存
	for i, contract := range insertContracts {
		contractJson, err := json.Marshal(contract)
		if err != nil {
			log.Errorf("Error Marshal contract err: %v，redisKey：:%v", err, redisKey)
		}
		cache.GlobalRedisDb.ZAdd(ctx, redisKey, redis.Z{
			Score:  float64(blockHeight*10000 + int64(i)),
			Member: string(contractJson),
		})
	}

	// 保留最新的 10 条区块数据
	cache.GlobalRedisDb.ZRemRangeByRank(ctx, redisKey, 0, -11)
	//更新合约版本
	UpdateLatestContractCache(chainId, updateContracts)
}

// UpdateLatestContractCache 最新合约列表
func UpdateLatestContractCache(chainId string, updateContracts []*db.Contract) {
	if len(updateContracts) == 0 {
		return
	}

	//如果缓存内的合约版本更新了，也需要实时更新
	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisLatestContractList, prefix, chainId)
	// 获取缓存中的合约列表及其 Score
	contractList, err := cache.GlobalRedisDb.ZRangeWithScores(ctx, redisKey, 0, -1).Result()
	if err != nil {
		log.Errorf("Error ZRangeWithScores contract err: %v，redisKey：:%v", err, redisKey)
		return
	}

	updatedContractMap := make(map[string]*db.Contract, 0)
	for _, contract := range updateContracts {
		updatedContractMap[contract.Addr] = contract
	}

	// 遍历合约列表，找到需要更新的合约
	for _, contractWithScore := range contractList {
		if contractWithScore.Member == nil {
			continue
		}

		member, ok := contractWithScore.Member.(string)
		if !ok {
			continue
		}

		var contract db.Contract
		err = json.Unmarshal([]byte(member), &contract)
		if err != nil {
			log.Errorf("Error Unmarshal contract err: %v，redisKey：:%v, redisRes:%v", err, redisKey,
				contractWithScore.Member)
			return
		}

		// 检查是否是需要更新的合约
		if _, ok := updatedContractMap[contract.Addr]; !ok {
			continue
		}

		// 从缓存中删除旧的合约数据
		cache.GlobalRedisDb.ZRem(ctx, redisKey, contractWithScore.Member)

		// 将更新后的合约数据添加到缓存中，使用原来的 Score
		updateContract := updatedContractMap[contract.Addr]
		redisContract := contract
		redisContract.Version = updateContract.Version
		redisContract.ContractStatus = updateContract.ContractStatus
		redisContract.UpgradeAddr = updateContract.UpgradeAddr
		redisContract.UpgradeTimestamp = updateContract.UpgradeTimestamp
		if updateContract.TxNum > 0 {
			redisContract.TxNum = updateContract.TxNum
		}
		if updateContract.EventNum > 0 {
			redisContract.EventNum = updateContract.EventNum
		}
		updatedContractJson, err := json.Marshal(redisContract)
		if err != nil {
			log.Errorf("Error Marshal contract err: %v，redisContract：:%v", err, redisContract)
			return
		}

		cache.GlobalRedisDb.ZAdd(ctx, redisKey, redis.Z{
			Score:  contractWithScore.Score,
			Member: string(updatedContractJson),
		})
	}
}
