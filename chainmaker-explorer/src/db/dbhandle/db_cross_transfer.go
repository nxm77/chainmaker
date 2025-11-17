package dbhandle

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// InsertCrossTxTransfers 插入跨链交易流转
func InsertCrossTxTransfers(chainId string, crossTxTransfers []*db.CrossTransactionTransfer) error {
	if len(crossTxTransfers) == 0 {
		return nil
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableCrossTransactionTransfer)
	err := CreateInBatchesData(tableName, crossTxTransfers)
	if err != nil {
		return err
	}
	return nil
}

// UpdateAccount 更新账户
func UpdateCrossTxTransfers(chainId, crossId string, status int32, timestamp int64) error {
	if chainId == "" || crossId == "" {
		return db.ErrTableParams
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableCrossTransactionTransfer)
	where := map[string]interface{}{
		"crossId": crossId,
	}
	params := map[string]interface{}{
		"status":  status,
		"endTime": timestamp,
	}

	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	return err
}

// CheckCrossIdsExistenceTransfer 查询数据是否已经保存
func CheckCrossIdsExistenceTransfer(chainId string, crossIds []string) (map[string]bool, error) {
	txTransferMap := make(map[string]bool, 0)
	if chainId == "" || len(crossIds) == 0 {
		return txTransferMap, nil
	}

	tableName := db.GetTableName(chainId, db.TableCrossTransactionTransfer)
	// 查询与 crossIds 匹配的唯一 CrossId
	var foundIds []string
	err := db.GormDB.Table(tableName).Select("crossId").Where("crossId IN ?", crossIds).Find(&foundIds).Error

	if err != nil {
		return txTransferMap, err
	}

	// 将查询结果保存到 map 中
	for _, crossId := range foundIds {
		txTransferMap[crossId] = true
	}

	return txTransferMap, nil
}

// GetCrossCycleTransferById 根据Id获取交易流转
func GetCrossCycleTransferById(chainId, crossId string) ([]*db.CrossTransactionTransfer, error) {
	tableName := db.GetTableName(chainId, db.TableCrossTransactionTransfer)
	cycleTransfers := make([]*db.CrossTransactionTransfer, 0)
	where := map[string]interface{}{
		"crossId": crossId,
	}
	err := db.GormDB.Table(tableName).Where(where).Find(&cycleTransfers).Error
	if err != nil {
		return cycleTransfers, err
	}

	return cycleTransfers, nil
}

func GetCrossCycleTransferByCrossIds(chainId string, crossIds []string) ([]*db.CrossTransactionTransfer, error) {
	tableName := db.GetTableName(chainId, db.TableCrossTransactionTransfer)
	cycleTransfers := make([]*db.CrossTransactionTransfer, 0)
	conditions := []QueryCondition{
		{Field: "crossId", Value: crossIds, Condition: "in"},
	}
	query := BuildQuery(tableName, conditions)
	err := query.Find(&cycleTransfers).Error
	if err != nil {
		return cycleTransfers, err
	}

	return cycleTransfers, nil
}

// GetCrossCycleTransferByHeight 根据height获取交易流转
func GetCrossCycleTransferByHeight(chainId string, blockHeights []int64) ([]*db.CrossTransactionTransfer, error) {
	tableName := db.GetTableName(chainId, db.TableCrossTransactionTransfer)
	cycleTransfers := make([]*db.CrossTransactionTransfer, 0)
	conditions := []QueryCondition{
		{Field: "blockHeight", Value: blockHeights, Condition: "in"},
	}
	query := BuildQuery(tableName, conditions)
	err := query.Find(&cycleTransfers).Error
	if err != nil {
		return cycleTransfers, err
	}

	return cycleTransfers, nil
}

func GetCrossTransferDurationByTime(chainId string, startTime, endTime int64) ([]*db.CrossTransactionTransfer, error) {
	result := make([]*db.CrossTransactionTransfer, 0)
	tableName := db.GetTableName(chainId, db.TableCrossTransactionTransfer)
	conditions := []QueryCondition{
		{Field: "startTime", Value: startTime, Condition: ">="},
		{Field: "endTime", Value: endTime, Condition: "<="},
		{Field: "status", Value: 3, Condition: ">="},
	}
	query := db.GormDB.Table(tableName)
	query = BuildQueryNew(query, conditions)
	//query := BuildQueryNew(tableName, conditions)
	err := query.Find(&result).Error
	if err != nil {
		return result, err
	}

	return result, nil
}

func GetCrossSubChainTransferList(offset, limit int, startTime, endTime int64, chainId, crossId,
	subChainId, fromChainId, toChainId string) ([]*db.CrossTransactionTransfer, int64, error) {
	transferList := make([]*db.CrossTransactionTransfer, 0)
	tableName := db.GetTableName(chainId, db.TableCrossTransactionTransfer)
	conditions := []QueryCondition{}
	if crossId != "" {
		conditions = append(conditions, QueryCondition{Field: "crossId", Value: crossId, Condition: "="})
	}
	if startTime != 0 && endTime != 0 {
		conditions = append(conditions, QueryCondition{Field: "startTime", Value: startTime, Condition: ">="})
		conditions = append(conditions, QueryCondition{Field: "startTime", Value: endTime, Condition: "<="})
	}
	if fromChainId != "" {
		conditions = append(conditions, QueryCondition{Field: "fromChainId", Value: fromChainId, Condition: "="})
	}
	if toChainId != "" {
		conditions = append(conditions, QueryCondition{Field: "toChainId", Value: toChainId, Condition: "="})
	}
	if subChainId != "" {
		conditions = append(conditions, QueryCondition{
			Field:     "fromChainId",
			Value:     subChainId,
			Condition: "=",
			Operator:  "OR",
		})
		conditions = append(conditions, QueryCondition{
			Field:     "toChainId",
			Value:     subChainId,
			Condition: "=",
			Operator:  "OR"})
	}

	var totalCount int64
	query := db.GormDB.Table(tableName)
	query = BuildQueryNew(query, conditions)
	err := query.Count(&totalCount).Error
	if err != nil {
		return transferList, 0, err
	}

	query = query.Order("startTime DESC").Offset(offset * limit).Limit(limit)
	err = query.Find(&transferList).Error
	if err != nil {
		return nil, 0, fmt.Errorf("GetCrossSubChainTransferList err, cause : %s", err.Error())
	}

	return transferList, totalCount, nil
}

func GetCrossTxCount(chainId string) (int64, error) {
	var totalCount int64
	tableName := db.GetTableName(chainId, db.TableCrossTransactionTransfer)
	err := db.GormDB.Table(tableName).Count(&totalCount).Error
	if err != nil {
		return 0, err
	}
	return totalCount, nil
}

// GetCrossLatestCycleTxList 获取最后10个交易列表
func GetCrossTxTransferLatestList(chainId string) ([]*db.CrossTransactionTransfer, error) {
	cacheResults := GetCrossLatestCycleTxCache(chainId)
	if cacheResults != nil {
		return cacheResults, nil
	}

	transferList := make([]*db.CrossTransactionTransfer, 0)
	tableName := db.GetTableName(chainId, db.TableCrossTransactionTransfer)
	query := db.GormDB.Table(tableName).Order("startTime DESC").Limit(10)
	err := query.Find(&transferList).Error
	if err != nil {
		return nil, fmt.Errorf("GetLatestTxList By chainId err, cause : %s", err.Error())
	}

	SetCrossLatestCycleTxCache(chainId, transferList)
	return transferList, nil
}

// GetCrossLatestCycleTxCache 获取最后10个交易列表
func GetCrossLatestCycleTxCache(chainId string) []*db.CrossTransactionTransfer {
	ctx := context.Background()
	cacheResult := make([]*db.CrossTransactionTransfer, 0)
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisCrossLatestTransactions, prefix, chainId)
	redisRes := cache.GlobalRedisDb.Get(ctx, redisKey)
	if redisRes == nil || redisRes.Val() == "" {
		return nil
	}

	err := json.Unmarshal([]byte(redisRes.Val()), &cacheResult)
	if err != nil {
		log.Errorf("【Redis】get cache failed, key:%v, result:%v", redisKey, redisRes)
		return nil
	}
	return cacheResult
}

// SetCrossLatestCycleTxCache 缓存最后10个交易列表
func SetCrossLatestCycleTxCache(chainId string, crossCycleTxs []*db.CrossTransactionTransfer) {
	if len(crossCycleTxs) == 0 {
		return
	}

	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisCrossLatestTransactions, prefix, chainId)
	retJson, err := json.Marshal(crossCycleTxs)
	if err != nil {
		log.Errorf("【Redis】set cache failed, key:%v, result:%v", redisKey, retJson)
		return
	}
	// 设置键值对和过期时间(1h 过期)
	_ = cache.GlobalRedisDb.Set(ctx, redisKey, string(retJson), time.Minute).Err()
}

func GetCrossContractByTransfer(chainId string) ([]*db.CrossContractByTransfer, error) {
	if chainId == "" {
		return nil, db.ErrTableParams
	}

	tableName := db.GetTableName(chainId, db.TableCrossTransactionTransfer)
	contracts := make([]*db.CrossContractByTransfer, 0)
	err := db.GormDB.Table(tableName).Select("toChainId, ContractName").
		Group("toChainId, ContractName").Find(&contracts).Error
	if err != nil {
		return nil, err
	}
	return contracts, nil
}
