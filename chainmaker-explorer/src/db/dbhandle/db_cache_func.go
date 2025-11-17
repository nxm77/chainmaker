package dbhandle

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func SetTotalTxNumCache(chainId string, num int64) {
	if num == 0 {
		return
	}

	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisOverviewTxTotal, prefix, chainId)
	// 设置键值对和过期时间
	ctx := context.Background()
	_ = cache.GlobalRedisDb.Set(ctx, redisKey, num, time.Second).Err()
}

// GetTotalTxNumCache
//
//	@Description: 获取首页交易总量缓存
//	@param chainId
//	@return int64 交易总量
//	@return error
func GetTotalTxNumCache(chainId string) (int64, error) {
	var count int64
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisOverviewTxTotal, prefix, chainId)
	ctx := context.Background()
	result, err := cache.GlobalRedisDb.Get(ctx, redisKey).Result()
	if err != nil {
		return 0, err
	}
	count, err = strconv.ParseInt(result, 10, 64)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// SetMaxBlockHeightCache
//
//	@Description: 缓存最高区块高度
//	@param chainId
//	@param num
func SetMaxBlockHeightCache(chainId string, blockHeight int64) {
	if blockHeight == 0 {
		return
	}

	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisOverviewMaxBlockHeight, prefix, chainId)
	// 设置键值对和过期时间
	ctx := context.Background()
	_ = cache.GlobalRedisDb.Set(ctx, redisKey, blockHeight, 24*time.Hour).Err()
}

// GetMaxBlockHeightCache
//
//	@Description: 获取最高区块高度缓存
//	@param chainId
//	@return int64 交易总量
//	@return error
func GetMaxBlockHeightCache(chainId string) (int64, error) {
	var blockHeight int64
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisOverviewMaxBlockHeight, prefix, chainId)
	ctx := context.Background()
	result, err := cache.GlobalRedisDb.Get(ctx, redisKey).Result()
	if err != nil {
		return 0, err
	}
	blockHeight, err = strconv.ParseInt(result, 10, 64)
	if err != nil {
		return 0, err
	}

	return blockHeight, nil
}

// UpdateContractNumCache 更新合约缓存
// @param chainId 链id
// @param contractInfo 合约信息
// func UpdateContractNumCache(chainId string, contractInfo *db.Contract) {
// 	// contractCache, _ := GetContractInfoCache(chainId, contractInfo.Addr)
// 	// if contractCache == nil {
// 	// 	return
// 	// }

// 	//交易量只会越来越大，交易量减少了说明合约更新覆盖交易量
// 	if contractInfo.TxNum < contractCache.TxNum ||
// 		contractInfo.EventNum < contractCache.EventNum {
// 		return
// 	}

// 	contractCache.TxNum = contractInfo.TxNum
// 	contractCache.EventNum = contractInfo.EventNum

// 	//写入缓存
// 	//SetContractInfoCache(chainId, contractCache)
// }

// SetContractInfoCache 设置合约缓存，合约缓存实现影响交易数据
// func SetContractInfoCache(chainId string, contractInfo *db.Contract) {
// 	if contractInfo == nil {
// 		return
// 	}

// 	keyInfoName := cache.GetKeyContractInfo(chainId, contractInfo.NameBak)
// 	keyInfoAddr := cache.GetKeyContractInfo(chainId, contractInfo.Addr)
// 	retJson, err := json.Marshal(contractInfo)
// 	if err == nil {
// 		// 设置键值对和过期时间
// 		ctx := context.Background()
// 		cache.GlobalRedisDb.Set(ctx, keyInfoName.Key, string(retJson), keyInfoName.Duration)
// 		cache.GlobalRedisDb.Set(ctx, keyInfoAddr.Key, string(retJson), keyInfoAddr.Duration)
// 		log.Infof("SetContractInfoCache success, keyInfoName:%s, retJson:%s",
// 			keyInfoName.Key, string(retJson))
// 	}
// }

// GetContractInfoCache 获取合约缓存，合约缓存实现影响交易数据
// @param chainId 链id
// @param contractKey 合约名称或者地址
// @return 合约信息
// func GetContractInfoCache(chainId, contractKey string) (*db.Contract, error) {
// 	var result *db.Contract
// 	ctx := context.Background()
// 	cacheKeyInfo := cache.GetKeyContractInfo(chainId, contractKey)
// 	redisRes := cache.GlobalRedisDb.Get(ctx, cacheKeyInfo.Key)
// 	if redisRes != nil && redisRes.Val() != "" {
// 		err := json.Unmarshal([]byte(redisRes.Val()), &result)
// 		if err == nil {
// 			return result, nil
// 		}
// 		log.Errorf("GetContractCache json Unmarshal err : %s, redisRes:%v", err.Error(), redisRes)
// 		return nil, err
// 	}

// 	// 当缓存未命中时，返回 nil 和 nil 错误
// 	return nil, nil
// }

// DelContractCountCache del Cache 删除合约总数量缓存数据
func DelContractCountCache(chainId string) {
	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisOverviewContractCount, prefix, chainId)
	_ = cache.GlobalRedisDb.Del(ctx, redisKey).Err()
}

// DelContractInfoCache del Cache 删除合约缓存数据
// func DelContractInfoCache(chainId, contractName, contractAddr string) {
// 	ctx := context.Background()
// 	cacheKeyAddr := cache.GetKeyContractInfo(chainId, contractAddr)
// 	cacheKeyName := cache.GetKeyContractInfo(chainId, contractName)
// 	// 设置键值对和过期时间
// 	delKey := []string{cacheKeyAddr.Key, cacheKeyName.Key}
// 	_ = cache.GlobalRedisDb.Del(ctx, delKey...).Err()
// }

// GetAccountCacheByAddr
//
//	@Description: 获取账户信息缓存
//	@param chainId
//	@param address 账户地址
//	@return *db.Account 账户信息
//	@return error
func GetAccountCacheByAddr(chainId, address string) (*db.Account, error) {
	var result *db.Account
	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisDBAccountData, prefix, chainId, address)
	redisRes := cache.GlobalRedisDb.Get(ctx, redisKey)
	if redisRes != nil && redisRes.Val() != "" {
		err := json.Unmarshal([]byte(redisRes.Val()), &result)
		if err != nil {
			log.Errorf("GetAccountCacheByAddr json Unmarshal err : %s, redisRes:%v", err.Error(), redisRes)
			return nil, err
		}
	}

	return nil, nil
}

// SetAccountDataCache
//
//	@Description: 设置账户信息缓存
//	@param chainId
//	@param accountData 账户信息
func SetAccountDataCache(chainId string, accountData *db.Account) {
	if accountData == nil {
		return
	}

	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisDBAccountData, prefix, chainId, accountData.Address)
	retJson, err := json.Marshal(accountData)
	if err == nil {
		// 设置键值对和过期时间
		ctx := context.Background()
		_ = cache.GlobalRedisDb.Set(ctx, redisKey, string(retJson), 24*time.Hour).Err()
	}
}

// UpdateAccountDataCache
//
//	@Description: 更新账户缓存
//	@param chainId
//	@param accountInfo
func UpdateAccountDataCache(chainId string, accountInfo *db.Account) {
	accountCache, err := GetAccountCacheByAddr(chainId, accountInfo.Address)
	if accountCache == nil || err != nil {
		return
	}

	if accountCache.TxNum > accountInfo.TxNum ||
		accountCache.NFTNum > accountInfo.NFTNum ||
		accountCache.BlockHeight > accountInfo.BlockHeight {
		//交易量只会越来越大，交易量减少了说明合约更新覆盖交易量
		return
	}

	SetAccountDataCache(chainId, accountInfo)
}

// SetLatestBlockListCache
//
//	@Description: 缓存最新区块列表
//	@param chainId
//	@param blockList
func SetLatestBlockListCache(chainId string, blockList []*db.Block) {
	var ctx *gin.Context
	if len(blockList) == 0 {
		return
	}
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisLatestBlockList, prefix, chainId)
	for _, blockInfo := range blockList {
		blockHeight := blockInfo.BlockHeight
		//缓存区块信息
		blockJson, _ := json.Marshal(blockInfo)
		cache.GlobalRedisDb.ZAdd(ctx, redisKey, redis.Z{
			Score:  float64(blockHeight),
			Member: string(blockJson),
		})
	}

	// 保留最新的 10 条区块数据
	cache.GlobalRedisDb.ZRemRangeByRank(ctx, redisKey, 0, -11)

	// 设置过期时间
	cache.GlobalRedisDb.Expire(ctx, redisKey, 12*time.Hour)
}

// GetLatestBlockListCache
//
//	@Description: 获取最新区块列表
//	@param chainId 链Id
//	@return []*db.Block 区块列表
//	@return error
func GetLatestBlockListCache(chainId string) ([]*db.Block, error) {
	var ctx *gin.Context
	blockList := make([]*db.Block, 0)
	//从缓存获取最新的block
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisLatestBlockList, prefix, chainId)
	redisBlockList := cache.GlobalRedisDb.ZRevRange(ctx, redisKey, 0, 9).Val()
	for _, blockStr := range redisBlockList {
		blockInfo := &db.Block{}
		err := json.Unmarshal([]byte(blockStr), blockInfo)
		if err != nil {
			log.Errorf("getBlockListFromRedis json Unmarshal err : %s redisRes :%v", err.Error(), blockStr)
			return blockList, err
		}
		blockList = append(blockList, blockInfo)
	}

	return blockList, nil
}

// GetTransactionListCountCache
//
//	@Description: 获取交易列表缓存数据
//	@param chainId
//	@param selectFile 查询条件
//	@return int64 交易数量
//	@return error
func GetTransactionListCountCache(chainId string, selectFile *SelectFile) (int64, error) {
	redisKeySelect := GenerateRedisKey(selectFile)
	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisTransactionListCount, prefix, chainId, redisKeySelect)
	result, err := cache.GlobalRedisDb.Get(ctx, redisKey).Result()
	if err != nil {
		return 0, err
	}
	count, err := strconv.ParseInt(result, 10, 64)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// SetTransactionListCountCache
//
//	@Description: 设置交易列表交易总数缓存
//	@param chainId
//	@param selectFile 查询条件
//	@param count 交易数量
func SetTransactionListCountCache(chainId string, selectFile *SelectFile, count int64) {
	ctx := context.Background()
	redisKeySelect := GenerateRedisKey(selectFile)
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisTransactionListCount, prefix, chainId, redisKeySelect)
	// 设置键值对和过期时间(30min 过期)
	_ = cache.GlobalRedisDb.Set(ctx, redisKey, count, 3*time.Minute).Err()
}

// GetNFTContractDataCache
//
//	@Description: 获取同质化合约
//	@param chainId
//	@param address
//	@return *db.Gas
func GetNFTContractDataCache(chainId, contractKey string) *db.NonFungibleContract {
	var result *db.NonFungibleContract
	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisNFTContractData, prefix, chainId, contractKey)
	redisRes := cache.GlobalRedisDb.Get(ctx, redisKey)
	if redisRes != nil && redisRes.Val() != "" {
		err := json.Unmarshal([]byte(redisRes.Val()), &result)
		if err == nil {
			return result
		}
		log.Errorf("GetNFTContractDataCache json Unmarshal err : %s, redisRes:%v", err.Error(), redisRes)
	}

	return nil
}

// SetNFTContractDataCache
//
//	@Description: 设置同质化合约
//	@param chainId
//	@param gasInfo
func SetNFTContractDataCache(chainId string, contractInfo *db.NonFungibleContract) {
	if contractInfo == nil {
		return
	}

	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKeyName := fmt.Sprintf(cache.RedisNFTContractData, prefix, chainId, contractInfo.ContractNameBak)
	redisKeyAddr := fmt.Sprintf(cache.RedisNFTContractData, prefix, chainId, contractInfo.ContractAddr)
	retJson, err := json.Marshal(contractInfo)
	if err == nil {
		// 设置键值对和过期时间
		ctx := context.Background()
		_ = cache.GlobalRedisDb.Set(ctx, redisKeyName, string(retJson), 30*time.Minute).Err()
		_ = cache.GlobalRedisDb.Set(ctx, redisKeyAddr, string(retJson), 30*time.Minute).Err()
	}
}

// GetFTContractDataCache
//
//	@Description: 获取同质化合约
//	@param chainId
//	@param address
//	@return *db.Gas
func GetFTContractDataCache(chainId, address string) *db.FungibleContract {
	var result *db.FungibleContract
	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisFTContractData, prefix, chainId, address)
	redisRes := cache.GlobalRedisDb.Get(ctx, redisKey)
	if redisRes != nil && redisRes.Val() != "" {
		err := json.Unmarshal([]byte(redisRes.Val()), &result)
		if err == nil {
			return result
		}
		log.Errorf("GetFTContractDataCache json Unmarshal err : %s, redisRes:%v", err.Error(), redisRes)
	}

	return nil
}

// SetFTContractDataCache
//
//	@Description: 设置同质化合约
//	@param chainId
//	@param gasInfo
func SetFTContractDataCache(chainId string, contractInfo *db.FungibleContract) {
	if contractInfo == nil {
		return
	}

	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKeyName := fmt.Sprintf(cache.RedisFTContractData, prefix, chainId, contractInfo.ContractNameBak)
	redisKeyAddr := fmt.Sprintf(cache.RedisFTContractData, prefix, chainId, contractInfo.ContractAddr)
	retJson, err := json.Marshal(contractInfo)
	if err == nil {
		// 设置键值对和过期时间
		ctx := context.Background()
		_ = cache.GlobalRedisDb.Set(ctx, redisKeyName, string(retJson), 30*time.Minute).Err()
		_ = cache.GlobalRedisDb.Set(ctx, redisKeyAddr, string(retJson), 30*time.Minute).Err()
	}
}

// SetMaxBlockHeightCache
//
//	@Description: 缓存最高区块高度
//	@param chainId
//	@param num
func SetAccountTotalCache(chainId string, total int64) {
	if total == 0 {
		return
	}

	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisDBAccountTotal, prefix, chainId)
	// 设置键值对和过期时间
	ctx := context.Background()
	_ = cache.GlobalRedisDb.Set(ctx, redisKey, total, time.Minute*5).Err()
}

// GetMaxBlockHeightCache
//
//	@Description: 获取最高区块高度缓存
//	@param chainId
//	@return int64 交易总量
//	@return error
func GetAccountTotalCache(chainId string) int64 {
	var total int64
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisDBAccountTotal, prefix, chainId)
	ctx := context.Background()
	result, err := cache.GlobalRedisDb.Get(ctx, redisKey).Result()
	if err != nil {
		return 0
	}
	total, err = strconv.ParseInt(result, 10, 64)
	if err != nil {
		return 0
	}

	return total
}

// DelAccountTotalCache delete cache
func DelAccountTotalCache(chainId string) {
	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisDBAccountTotal, prefix, chainId)
	_ = cache.GlobalRedisDb.Del(ctx, redisKey).Err()
}
