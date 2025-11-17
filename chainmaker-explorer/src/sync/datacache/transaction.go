package datacache

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/redis/go-redis/v9"
)

// GetLatestTxListCache 获取最新的交易列表
// @param chainId 链id
// @return []*db.Transaction 交易列表
// @return error 错误信息
func GetLatestTxListCache(chainId string) ([]*db.Transaction, error) {
	ctx := context.Background()
	txList := make([]*db.Transaction, 0)
	//从缓存获取最新的block
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisLatestTransactions, prefix, chainId)
	redisTxList := cache.GlobalRedisDb.ZRevRange(ctx, redisKey, 0, 9).Val()
	for _, txString := range redisTxList {
		txInfo := &db.Transaction{}
		err := json.Unmarshal([]byte(txString), txInfo)
		if err != nil {
			log.Errorf("GetLatestTxListCache json Unmarshal err : %s, redisRes:%v", err.Error(), txString)
			return txList, err
		}
		txList = append(txList, txInfo)
	}

	return txList, nil
}

// BuildLatestTxListCache 构建最新的交易列表缓存
// @param chainId 链id
// @param txMap 交易map
func BuildLatestTxListCache(chainId string, txMap map[string]*db.Transaction) {
	// 从缓存中获取交易列表
	txList, err := GetLatestTxListCache(chainId)
	if len(txList) == 0 || err != nil {
		//缓存可能丢失
		txList, _ = dbhandle.GetLatestTxList(chainId)
		if len(txList) == 0 {
			for _, txInfo := range txMap {
				txList = append(txList, txInfo)
			}
		}
	} else {
		//缓存存在,缓存数据加入新数据
		for _, txInfo := range txMap {
			txList = append(txList, txInfo)
		}
	}

	// 根据 blockHeight 和 txInfo.TxIndex 排序交易列表
	sort.Slice(txList, func(i, j int) bool {
		if txList[i].BlockHeight == txList[j].BlockHeight {
			return txList[i].TxIndex < txList[j].TxIndex
		}
		return txList[i].BlockHeight > txList[j].BlockHeight
	})

	// 保留最新的 10 条交易数据
	if len(txList) > 10 {
		txList = txList[:10]
	}

	// 缓存交易信息
	SetLatestTxListCache(chainId, txList)
}

// SetLatestTxListCache 设置最新的交易列表缓存
// @param chainId 链id
// @param transactions 交易列表
func SetLatestTxListCache(chainId string, transactions []*db.Transaction) {
	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisLatestTransactions, prefix, chainId)

	// 删除旧的缓存数据
	cache.GlobalRedisDb.Del(ctx, redisKey)

	// 将新的交易列表存储到缓存中
	for _, txInfo := range transactions {
		txInfoJson, err := json.Marshal(txInfo)
		if err != nil {
			log.Errorf("Error marshaling txInfo: %v，redisKey：:%v", err, redisKey)
			continue
		}
		blockHeight := txInfo.BlockHeight
		cache.GlobalRedisDb.ZAdd(ctx, redisKey, redis.Z{
			Score:  float64(blockHeight*100000 + int64(txInfo.TxIndex)),
			Member: string(txInfoJson),
		})
	}

	// 设置过期时间
	cache.GlobalRedisDb.Expire(ctx, redisKey, 12*time.Hour)
}

// BuildOverviewTxTotalCache
//
//	@Description: 缓存首页交易总量
//	@param chainId
//	@param transactions
func BuildOverviewTxTotalCache(chainId string, txCount int64) {
	if txCount == 0 {
		return
	}

	txTotal, err := dbhandle.GetTotalTxNumCache(chainId)
	if err != nil {
		log.Errorf("BuildOverviewTxTotalCache GetTotalTxNumCache err :%v", err)
		return
	}

	txTotal += txCount

	// 缓存交易总量
	dbhandle.SetTotalTxNumCache(chainId, txTotal)
}
