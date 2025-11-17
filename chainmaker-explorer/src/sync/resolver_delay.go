/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/sync/common"
	datacache "chainmaker_web/src/sync/datacache"
	"chainmaker_web/src/sync/logic"
	"chainmaker_web/src/sync/model"
	"chainmaker_web/src/sync/saveTasks"
	"context"
	"sync"
	"time"
)

// BatchDelayedUpdate
//
//	@Description:批量延迟更新
//	@param chainId 链id
//	@param blockHeights 批量处理区块高度
//	@return error
func BatchDelayedUpdate(chainId string, blockHeights []int64) error {
	startTime := time.Now()
	if len(blockHeights) == 0 {
		return nil
	}

	log.Infof("【Delay update】start block-%s[%v]", chainId, blockHeights)
	//获取缓存数据,缓存缺失从数据库查询（同步插入的交易列表，合约类别等计算需要用到的数据）
	delayedUpdateNeedCache, err := GetRealtimeDataCacheOrDB(chainId, blockHeights)
	if err != nil {
		return err
	}

	//计算所有需要更新的数据
	delayedUpdateData, err := BuildDelayedUpdateData(chainId, blockHeights, delayedUpdateNeedCache)
	if err != nil {
		return err
	}

	//并发插入，更新数据库
	err = ParallelParseUpdateDataToDB(chainId, delayedUpdateData)
	if err != nil {
		return err
	}

	//最后更新区块状态，所有数据更新结束
	err = saveTasks.UpdateBlockStatusToDB(chainId, blockHeights)
	if err != nil {
		return err
	}

	//更新首页合约缓存交易量
	datacache.UpdateLatestContractCache(chainId, delayedUpdateData.ContractResult.UpdateContractTxEventNum)
	durationTime := time.Since(startTime).Milliseconds()
	log.Infof("【Delay update】end block-%s[%v], duration_time:%vms", chainId, blockHeights, durationTime)
	return nil
}

// GetRealtimeDataCacheOrDB 根据缓存获取需要异步更新的数据
// @param chainId
// @param blockHeights
// @return *GetRealtimeCacheData 缓存数据
// @return error
func GetRealtimeDataCacheOrDB(chainId string, blockHeights []int64) (*model.GetRealtimeCacheData, error) {
	//获取同步插入缓存数据
	delayedUpdateData := model.NewGetRealtimeCacheData()

	//缓存缺失的height
	heightDB := make([]int64, 0)
	for _, height := range blockHeights {
		//获取同步插入缓存数据
		isHave := datacache.GetRealtimeDataCache(chainId, height, delayedUpdateData)
		if !isHave {
			heightDB = append(heightDB, height)
		}

		// //获取主子链缓存
		// isHaveCross := datacache.GetRealtimeCrossCache(chainId, height, delayedUpdateData)
		// if !isHaveCross {
		// 	crossHeightDB = append(crossHeightDB, height)
		// }
	}

	//缓存没有，从数据库获取
	if len(heightDB) > 0 {
		//缓存没有，数据库取数据
		err := GetDelayedUpdateByDB(chainId, heightDB, delayedUpdateData)
		if err != nil {
			return delayedUpdateData, err
		}
	}

	// //主子链缓存没有，从数据库获取
	// if len(crossHeightDB) > 0 {
	// 	//缓存没有，数据库取数据
	// 	crossCycleTransfers, err := dbhandle.GetCrossCycleTransferByHeight(chainId, crossHeightDB)
	// 	if err != nil {
	// 		return delayedUpdateData, err
	// 	}
	// 	delayedUpdateData.CrossTransfers = append(delayedUpdateData.CrossTransfers, crossCycleTransfers...)
	// }

	return delayedUpdateData, nil
}

// BuildDelayedUpdateData
//
//	@Description: 计算所有需要更新的数据
//	@param chainId
//	@param blockHeights 批量处理的区块高度列表
//	@param delayedUpdateCache 同步插入的缓存数据
//	@return DelayedUpdateData 需要更新数据库的结构化数据
func BuildDelayedUpdateData(chainId string, blockHeights []int64, delayedUpdateCache *model.GetRealtimeCacheData) (
	*model.DelayedUpdateData, error) {
	//本次批量处理的最新的区块高度，用于确定异常情况重复更新问题
	minHeight := common.GetMinBlockHeight(blockHeights)

	//获取本次涉及到的合约信息
	contractMap, err := logic.GetContractMapByAddrs(chainId, delayedUpdateCache.ContractAddrs)
	if err != nil {
		return nil, err
	}

	//解析合约event
	topicEventResult := logic.DealTopicEventData(chainId, delayedUpdateCache.ContractEvents, contractMap,
		delayedUpdateCache.TxList)

	//解析出的事件列表
	eventDataList := topicEventResult.ContractEventData

	//主子链解析transfer
	insertCrossTransfers := delayedUpdateCache.CrossChainResult.InsertCrossTransfer

	//计算主子链跨链次数
	crossSubChainIdMap := logic.ParseCrossCycleTxTransfer(insertCrossTransfers)

	eventTopicTxNum := logic.DealEventTopicTxNum(delayedUpdateCache.ContractEvents)

	//DB并发获取合约，交易，持仓等数据库数据
	delayGetDBResult, err := DelayParallelParseGetDB(chainId, delayedUpdateCache, contractMap, topicEventResult,
		crossSubChainIdMap, eventTopicTxNum)
	if err != nil {
		return nil, err
	}

	//计算IDA合约数据
	idaContractMap := delayGetDBResult.IDAContractMap
	//初始化 ContractHandler
	contractHandler := logic.NewContractHandler(chainId, minHeight, contractMap, topicEventResult, delayedUpdateCache)
	//调用 ContractHandler 处理合约数据
	contractResult, err := contractHandler.DealWithContractData(idaContractMap)
	if err != nil {
		return nil, err
	}

	insertEventTopics, updateEventTopics := logic.ProcessEventTopicTxNum(eventTopicTxNum,
		delayGetDBResult.EventTopicMap, minHeight)

	//主子链计算跨链交易数
	crossSubChainCrossDB := delayGetDBResult.CrossSubChainCross
	insertSubChainCross, updateSubChainCross, err := logic.DealSubChainCrossChainNum(chainId, crossSubChainIdMap,
		crossSubChainCrossDB, minHeight)
	if err != nil {
		return nil, err
	}

	//主子链计算子链交易数
	insertSubChainData, updateSubChainData, crossChainContracts := logic.DealCrossSubChainData(
		insertCrossTransfers, delayGetDBResult.CrossSubChainMap)

	// 获取新增，更新账户信息
	accountHandler := logic.NewAccountHandler(chainId, minHeight, topicEventResult, delayedUpdateCache,
		delayGetDBResult, eventDataList)
	updateAccountResult := accountHandler.DealWithAccountData()

	//计算新增，更新gas数据
	insertGasList, updateGasList := logic.BuildGasInfo(delayedUpdateCache.GasRecords, delayGetDBResult.GasList, minHeight)

	//统计transfer流转记录
	fungibleTransfer, nonFungibleTransfer := logic.DealTransferList(eventDataList, contractMap, delayedUpdateCache.TxList)
	fungibleTransfer = logic.BuildAccountManagerGasTransfer(chainId, delayedUpdateCache.TxList, fungibleTransfer)

	//统计token列表
	tokenResult := logic.DealNonFungibleToken(chainId, eventDataList, contractMap, accountHandler.AccountMap)

	//统计新增持仓数据
	positionList := logic.BuildPositionList(eventDataList, contractMap, accountHandler.AccountMap)
	//计算持仓数据
	positionDBMap := delayGetDBResult.PositionMapList
	nonPositionDBMap := delayGetDBResult.NonPositionMapList
	positionOperates := logic.BuildUpdatePositionData(minHeight, positionList, positionDBMap, nonPositionDBMap)

	//交易黑名单1
	updateTxBlack := &db.UpdateTxBlack{
		AddTxBlack:    make([]*db.BlackTransaction, 0),
		DeleteTxBlack: make([]*db.Transaction, 0),
	}
	for _, txInfo := range delayGetDBResult.AddBlackTxList {
		//添加黑名单
		updateTxBlack.AddTxBlack = append(updateTxBlack.AddTxBlack, (*db.BlackTransaction)(txInfo))
	}
	for _, txInfo := range delayGetDBResult.DeleteBlackTxList {
		//删除黑名单
		updateTxBlack.DeleteTxBlack = append(updateTxBlack.DeleteTxBlack, (*db.Transaction)(txInfo))
	}

	//计算持有量和发行总
	//持有人数
	holdCountMap := logic.DealContractHoldCount(positionOperates)
	//发行总量
	totalSupplyMap := logic.DealContractTotalSupply(eventDataList, contractMap)
	fungibleMap := delayGetDBResult.FungibleContractMap
	nonFungibleMap := delayGetDBResult.NonFungibleContractMap

	//计算同质化合约持有人数，发行量最终数据
	updateFTContractMap := logic.DealFungibleContractUpdateData(holdCountMap, totalSupplyMap, fungibleMap, minHeight)
	//计算FT合约交易流转数量
	ftContractTransferMap := logic.FTContractTransferNum(fungibleTransfer, fungibleMap, minHeight)
	mergedFTContract := logic.MergeFTContractMaps(minHeight, updateFTContractMap, ftContractTransferMap)

	//计算非同质化合约持有人数，发行量最终数据
	updateNFTContractMap := logic.DealNonFungibleContractUpdateData(holdCountMap, totalSupplyMap,
		nonFungibleMap, minHeight)
	//计算NFT合约交易流转数量
	nftContractTransferMap := logic.NFTContractTransferNum(nonFungibleTransfer, nonFungibleMap, minHeight)
	mergedNFTContract := logic.MergeNFTContractMaps(minHeight, updateNFTContractMap, nftContractTransferMap)

	//数据要素计算方式
	//插入IDA数据资产
	dealInsertIDAAssets := logic.DealInsertIDAAssetsData(idaContractMap, topicEventResult.IDAEventData)
	//更新IDA合约资产
	dealUpdateIDAAssets := logic.DealUpdateIDAAssetsData(topicEventResult.IDAEventData,
		delayGetDBResult.IDAAssetDetailMap)

	//合约更新数据
	contractResult.UpdateFungibleContract = mergedFTContract
	contractResult.UpdateNonFungible = mergedNFTContract
	contractResult.InsertEventTopic = insertEventTopics
	contractResult.UpdateEventTopic = updateEventTopics

	//did更新数据
	didSaveDate := logic.DealDIDSaveData(topicEventResult.DIDEventData)

	//链统计
	chainStatistics := logic.DealBlockchainStatistics(chainId, updateAccountResult.InsertAccount)

	//ABI主题表
	contractTopicList := logic.BuildEventDataByABI(chainId, delayedUpdateCache.ContractEvents)

	delayCrossChain := &model.DelayCrossChain{
		InsertSubChainData:  insertSubChainData,
		UpdateSubChainData:  updateSubChainData,
		InsertSubChainCross: insertSubChainCross,
		UpdateSubChainCross: updateSubChainCross,
		CrossChainContracts: crossChainContracts,
	}
	//构建异步更新数据
	buildDelayedUpdateData := &model.DelayedUpdateData{
		DelayCrossChain:     delayCrossChain,
		InsertGasList:       insertGasList,
		UpdateGasList:       updateGasList,
		UpdateTxBlack:       updateTxBlack,
		ContractResult:      contractResult,
		FungibleTransfer:    fungibleTransfer,
		NonFungibleTransfer: nonFungibleTransfer,
		BlockPosition:       positionOperates,
		UpdateAccountResult: updateAccountResult,
		TokenResult:         tokenResult,
		ContractMap:         contractMap,
		IDAInsertAssetsData: dealInsertIDAAssets,
		IDAUpdateAssetsData: dealUpdateIDAAssets,
		DIDSaveDate:         didSaveDate,
		ChainStatistics:     chainStatistics, //链统计
		ABITopicTableEvents: contractTopicList,
	}

	return buildDelayedUpdateData, nil
}

// GetDelayedUpdateByDB
//
//	@Description: 缓存数据如果没有，根据区块高度从数据库查
//	@param chainId
//	@param heightDB
//	@param delayedUpdateData 异步更新需要用到的数据
//	@return error
func GetDelayedUpdateByDB(chainId string, heightDB []int64, delayedUpdateData *model.GetRealtimeCacheData) error {
	//获取交易列表
	txInfoList, err := dbhandle.GetTxInfoByBlockHeight(chainId, heightDB)
	if err != nil {
		return err
	}

	//解析交易id列表，合约名称列表
	txIds, contractNameMap, txInfoMap := common.ExtractTxIdsAndContractNames(txInfoList)
	if len(txIds) == 0 {
		return nil
	}

	// 合并缓存数据到 delayedUpdateData
	for k, v := range txInfoMap {
		delayedUpdateData.TxList[k] = v
	}
	for name := range contractNameMap {
		delayedUpdateData.ContractAddrs[name] = name
	}

	errCh := make(chan error, 2)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		//获取gas记录
		gasRecords, err := logic.GetGasRecord(chainId, txIds)
		if err != nil {
			errCh <- err
			return
		}
		if len(gasRecords) > 0 {
			delayedUpdateData.GasRecords = append(delayedUpdateData.GasRecords, gasRecords...)
		}
	}()

	go func() {
		defer wg.Done()
		//获取合约event
		contractEvents, err := GetContractEvents(chainId, txIds)
		if err != nil {
			errCh <- err
			return
		}
		if len(contractEvents) > 0 {
			delayedUpdateData.ContractEvents = append(delayedUpdateData.ContractEvents, contractEvents...)
		}
	}()

	wg.Wait()
	close(errCh)
	for errDB := range errCh {
		if errDB != nil {
			// 重试多次仍未成功，停掉链，重新订阅
			log.Errorf("Error: %v", errDB)
			return errDB
		}
	}

	return nil
}

// ParallelParseUpdateDataToDB
//
//	@Description: 并发更新所有的表数据，失败会进行重试
//	@param chainId
//	@param delayedUpdateData 处理好的更新数据
//	@return error
func ParallelParseUpdateDataToDB(chainId string, delayedUpdateData *model.DelayedUpdateData) error {
	var err error
	// 数据插入
	// 初始化重试计数映射
	retryCountMap := &sync.Map{}
	// 创建一个可取消的上下文
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// 创建任务列表
	tasksList := createTasksDelayedUpdate(chainId, delayedUpdateData)
	// 创建一个错误通道
	errCh := make(chan error, len(tasksList))

	// 并发执行无依赖任务
	var wg sync.WaitGroup
	wg.Add(len(tasksList))
	for _, task := range tasksList {
		go saveTasks.ExecuteTaskWithRetry(ctx, &wg, task, retryCountMap, errCh)
	}

	wg.Wait()
	close(errCh)

	if len(errCh) > 0 {
		err = <-errCh
		if err != nil {
			// 处理错误
			log.Errorf("ParallelParseUpdateDataToDB error: %v", err)
			// 取消其他任务
			cancel()
			return err
		}
	}
	return nil
}

// createTasksDelayedUpdate
//
//	@Description: 数据插入任务列表
//	@param chainId
//	@param delayedUpdate 需要插入，更新的数据
//	@return []saveTasks.Task 任务列表
func createTasksDelayedUpdate(chainId string, delayedUpdate *model.DelayedUpdateData) []saveTasks.DealTask {
	// 定义任务列表
	tasksList := []saveTasks.DealTask{
		{
			Name:     "TaskUpdateContractResult",
			Function: saveTasks.TaskUpdateContractResult,
			Args:     []interface{}{chainId, delayedUpdate.ContractResult},
		},
		{
			Name:     "TaskInsertFungibleTransferToDB",
			Function: saveTasks.TaskInsertFungibleTransferToDB,
			Args:     []interface{}{chainId, delayedUpdate.FungibleTransfer},
		},
		{
			Name:     "TaskInsertNonFungibleTransferToDB",
			Function: saveTasks.TaskInsertNonFungibleTransferToDB,
			Args:     []interface{}{chainId, delayedUpdate.NonFungibleTransfer},
		},
		{
			Name:     "TaskSaveAccountListToDB",
			Function: saveTasks.TaskSaveAccountListToDB,
			Args:     []interface{}{chainId, delayedUpdate.UpdateAccountResult},
		},
		{
			Name:     "TaskSaveTokenResultToDB",
			Function: saveTasks.TaskSaveTokenResultToDB,
			Args:     []interface{}{chainId, delayedUpdate.TokenResult},
		},
		{
			Name:     "TaskUpdateTxBlackToDB",
			Function: saveTasks.TaskUpdateTxBlackToDB,
			Args:     []interface{}{chainId, delayedUpdate.UpdateTxBlack},
		},
		{
			Name:     "TaskSaveGasToDB",
			Function: saveTasks.TaskSaveGasToDB,
			Args:     []interface{}{chainId, delayedUpdate.InsertGasList, delayedUpdate.UpdateGasList},
		},
		{
			Name:     "TaskSaveFungibleContractResult",
			Function: saveTasks.TaskSaveFungibleContractResult,
			Args:     []interface{}{chainId, delayedUpdate.ContractResult},
		},
		{
			Name:     "TaskSavePositionToDB",
			Function: saveTasks.TaskSavePositionToDB,
			Args:     []interface{}{chainId, delayedUpdate.BlockPosition},
		},
		{
			Name:     "TaskSaveIDAAssetDataToDB",
			Function: saveTasks.TaskSaveIDAAssetDataToDB,
			Args:     []interface{}{chainId, delayedUpdate.IDAInsertAssetsData},
		},
		{
			Name:     "TaskDelayCrossChain",
			Function: saveTasks.TaskDelayCrossChain,
			Args:     []interface{}{chainId, delayedUpdate.DelayCrossChain},
		},
		{
			Name:     "TaskUpdateIDAAssetDataToDB",
			Function: saveTasks.TaskUpdateIDAAssetDataToDB,
			Args:     []interface{}{chainId, delayedUpdate.IDAUpdateAssetsData},
		},
		{
			Name:     "TaskUpdateDIDDataToDB",
			Function: saveTasks.TaskUpdateDIDDataToDB,
			Args:     []interface{}{chainId, delayedUpdate.DIDSaveDate},
		},
		{
			Name:     "TaskSaveChainStatisticsToDB",
			Function: saveTasks.TaskSaveChainStatisticsToDB,
			Args:     []interface{}{chainId, delayedUpdate.ChainStatistics},
		},
		{
			Name:     "TaskSaveABITopicTableEvents",
			Function: saveTasks.TaskSaveABITopicTableEvents,
			Args:     []interface{}{chainId, delayedUpdate.ABITopicTableEvents},
		},
	}

	return tasksList
}
