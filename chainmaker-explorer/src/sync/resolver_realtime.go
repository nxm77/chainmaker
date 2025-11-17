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
	"fmt"
	"sync"

	pbCommon "chainmaker.org/chainmaker/pb-go/v2/common"
)

type ProcessedBlockInfo struct {
	BlockInfo *pbCommon.BlockInfo
	HashType  string
}

// RealtimeDataHandle
//
//	@Description: 订阅解析区块数据
//	@param blockInfo 订阅区块
//	@param hashType
//	@return *RealtimeDealResult 解析格式化数据
//	@return *TxTimeLog 耗时日志
//	@return error
func (block *ProcessedBlockInfo) ProcessedBlockHandle() (*model.ProcessedBlockData, error) {
	blockInfo := block.BlockInfo
	hashType := block.HashType
	dealResult := model.NewProcessedBlockData()
	dealResult.Timestamp = blockInfo.Block.Header.BlockTimestamp
	//处理创世区块
	err := logic.GenesisBlockSystemContract(blockInfo, dealResult)
	if err != nil {
		return dealResult, err
	}

	errCh := make(chan error, 1) // 使用一个错误通道
	var wg sync.WaitGroup

	// 任务集合，包含所有的并发任务
	tasks := []struct {
		name string
		task func() error
	}{
		{
			//交易解析合约事件，gas，交易数据，合约升级，用户列表
			name: "ParallelParseTransactions",
			task: func() error {
				return logic.ParallelParseTransactions(blockInfo, hashType, dealResult)
			},
		},
		{
			//读写集解析合约数据，通用合约，存证合约
			name: "ParallelParseContract",
			task: func() error {
				return logic.ParallelParseContract(blockInfo, hashType, dealResult)
			},
		},
		{
			//解析区块数据
			name: "DealBlockInfo",
			task: func() error {
				return logic.DealBlockInfo(blockInfo, hashType, dealResult)
			},
		},
		{
			//读写集解析其他数据，合约配置，主子链
			name: "ParallelParseWriteSetData",
			task: func() error {
				return logic.ParallelParseWriteSetData(blockInfo, dealResult)
			},
		},
	}

	// 启动所有并发任务
	for _, task := range tasks {
		wg.Add(1)
		go func(taskName string, taskFunc func() error) {
			defer wg.Done()
			//失败重试方式
			saveTasks.WithRetry(taskFunc, taskName, errCh)
		}(task.name, task.task)
	}

	// 等待所有任务完成
	wg.Wait()

	// ---- 处理通道中的错误 ----
	close(errCh) // 关闭错误通道
	for err := range errCh {
		if err != nil {
			// 如果有错误，停止链并重新订阅
			log.Errorf("Error in processing block: %v", err)
			return dealResult, err
		}
	}

	// ---- 敏感词过滤 ----
	_ = common.FilterTxAndEvent(dealResult.Transactions, dealResult.ContractEvents)

	return dealResult, nil
}

// RealtimeDataSaveToDB
//
//	@Description:  顺序插入同步处理数据
//	@param chainId 链ID
//	@param blockHeight 区块高度
//	@param dealResult 格式化后待存储DB的数据
//	@param txTimeLog 耗时日志
//	@return error
func RealtimeDataSaveToDB(chainId string, blockHeight int64, dealResult *model.ProcessedBlockData) error {
	var err error
	//处理合约数据，判断是新增合约还是更新合约
	err = logic.HandleContractInsertOrUpdate(chainId, dealResult)
	if err != nil {
		return err
	}

	//将合约数据写入交易信息
	SetTransactionContract(chainId, dealResult)

	// 检查是否启用了 gas，如果没有启用 Gas，就会禁用或清空 Gas 相关记录
	err = common.CheckAndDisableGasIfNotEnabled(chainId, dealResult)
	if err != nil {
		return err
	}

	//检查跨链交易数据
	err = CheckCrossChainTransactionData(chainId, dealResult.CrossChainResult)
	if err != nil {
		return err
	}

	// 设置浏览器统计数据
	err = logic.DealStatisticsRealtime(chainId, blockHeight, dealResult)
	if err != nil {
		return err
	}

	// 执行数据插入任务
	err = executeDataInsertTasks(chainId, *dealResult)
	if err != nil {
		return err
	}

	//最后保存block数据
	err = dbhandle.InsertBlock(chainId, dealResult.BlockDetail)
	if err != nil {
		return err
	}

	//异步更新使用
	//设置缓存数据，
	datacache.SetDelayedUpdateCache(chainId, blockHeight, *dealResult)

	// //缓存主子链流转数据
	// datacache.SetCrossSubChainCrossCache(chainId, blockHeight, *dealResult)

	//浏览器首页使用
	//缓存最新交易列表
	datacache.BuildLatestTxListCache(chainId, dealResult.Transactions)
	//缓存首页交易总量
	//datacache.BuildOverviewTxTotalCache(chainId, int64(len(dealResult.Transactions)))
	//最新合约缓存
	datacache.SetLatestContractListCache(chainId, blockHeight, dealResult.InsertContracts, dealResult.UpdateContracts)
	//缓存最新区块高度
	logic.BuildOverviewMaxBlockHeightCache(chainId, dealResult.BlockDetail)
	//最新区块缓存
	logic.BuildLatestBlockListCache(chainId, dealResult.BlockDetail)

	return nil
}

// CheckCrossChainTransactionData 检查跨链交易数据
func CheckCrossChainTransactionData(chainId string, crossChainResult *model.CrossChainResult) error {
	saveTransferMap := crossChainResult.CrossTransfer
	var crossIds []string
	for _, transfer := range saveTransferMap {
		crossIds = append(crossIds, transfer.CrossId)
	}

	existsMap, err := dbhandle.CheckCrossIdsExistenceTransfer(chainId, crossIds)
	if err != nil {
		return err
	}

	// 分类处理：存在则放入UpdateCrossTransfer，不存在则放入insertCrossTransfer
	for _, transfer := range saveTransferMap {
		if _, exists := existsMap[transfer.CrossId]; exists {
			crossChainResult.UpdateCrossTransfer = append(crossChainResult.UpdateCrossTransfer, transfer)
		} else {
			crossChainResult.InsertCrossTransfer = append(crossChainResult.InsertCrossTransfer, transfer)
		}
	}

	return nil
}

// SetTransactionContract
//
//	@Description: 将合约数据写入交易表
//	@param chainId
//	@param transactionMap
//
// SetTransactionContract 函数用于将交易信息写入合约数据
func SetTransactionContract(chainId string, dealResult *model.ProcessedBlockData) {
	// 获取交易信息
	transactionMap := dealResult.Transactions
	// 获取合约数据
	contractWriteSetData := dealResult.ContractWriteSetData
	//交易信息写入合约数据
	for _, transaction := range transactionMap {
		// 根据合约名称获取合约信息
		contractInfo, err := fetchContractDetails(chainId, transaction.ContractNameBak, dealResult.InsertContracts)
		if err != nil {
			log.Warnf("fetchContractDetails err: %v", err)
			continue
		}

		// 如果合约信息不为空且没有错误
		if contractInfo != nil {
			// 将合约信息写入交易信息
			transaction.ContractName = contractInfo.Name
			transaction.ContractNameBak = contractInfo.NameBak
			transaction.ContractAddr = contractInfo.Addr
			transaction.ContractRuntimeType = contractInfo.RuntimeType
			transaction.ContractType = contractInfo.ContractType
		}
	}

	// 遍历升级合约交易
	for _, contractTx := range dealResult.UpgradeContractTx {
		// 如果合约数据中存在该交易
		if writeSetData, ok := contractWriteSetData[contractTx.TxId]; ok {
			// 将合约数据写入升级合约交易
			contractTx.ContractName = writeSetData.ContractName
			contractTx.ContractNameBak = writeSetData.ContractNameBak
			contractTx.ContractAddr = writeSetData.ContractAddr
			contractTx.ContractRuntimeType = writeSetData.RuntimeType
			contractTx.ContractVersion = writeSetData.Version
			contractTx.ContractType = writeSetData.ContractType
			// 如果交易信息中存在该交易
		} else if txInfo, ok := transactionMap[contractTx.TxId]; ok {
			// 将交易信息写入升级合约交易
			contractTx.ContractName = txInfo.ContractName
			contractTx.ContractNameBak = txInfo.ContractNameBak
			contractTx.ContractAddr = txInfo.ContractAddr
			contractTx.ContractRuntimeType = txInfo.ContractRuntimeType
			contractTx.ContractVersion = txInfo.ContractVersion
			contractTx.ContractType = txInfo.ContractType
		}
	}

	// 遍历合约事件
	for _, contractEvent := range dealResult.ContractEvents {
		// 根据合约名称获取合约信息
		contractInfo, err := fetchContractDetails(chainId, contractEvent.ContractNameBak, dealResult.InsertContracts)
		if err != nil {
			log.Errorf("fetchContractDetails err: %v, contractName: %s", err, contractEvent.ContractNameBak)
			continue
		}

		contractEvent.ContractName = contractInfo.NameBak
		contractEvent.ContractNameBak = contractInfo.NameBak
		contractEvent.ContractAddr = contractInfo.Addr
		contractEvent.ContractVersion = contractInfo.Version
		contractEvent.ContractType = contractInfo.ContractType
	}
}

// fetchContractDetails 根据合约名称获取合约信息
func fetchContractDetails(chainId, contractName string, insertContracts []*db.Contract) (*db.Contract, error) {
	for _, contract := range insertContracts {
		if contract.Name == contractName || contract.Addr == contractName {
			return contract, nil
		}
	}

	//DB获取合约信息
	contractInfo, err := dbhandle.GetContractByNameOrAddr(chainId, contractName)
	return contractInfo, err
}

// executeDataInsertTasks
//
//	@Description: 执行数据并发插入任务
//	@param chainId 链ID
//	@param dealResult 待插入数据
//	@return error 错误信息
func executeDataInsertTasks(chainId string, dealResult model.ProcessedBlockData) error {
	// 创建一个可取消的上下文
	_, cancel := context.WithCancel(context.Background())
	defer cancel()
	// 创建任务列表
	tasksList := createTasks(chainId, dealResult)
	// 创建一个错误通道
	errCh := make(chan error, len(tasksList))

	// 并发执行无依赖任务
	var wg sync.WaitGroup
	wg.Add(len(tasksList))
	for _, task := range tasksList {
		go func(t saveTasks.Task) {
			defer wg.Done()
			if err := t.Execute(); err != nil {
				errCh <- fmt.Errorf("task failed err: %v", err)
			}
		}(task)
	}

	wg.Wait()
	close(errCh)
	if len(errCh) > 0 {
		err := <-errCh
		if err != nil {
			// 处理错误
			log.Errorf("executeDataInsertTasks error: %v", err)
			// 取消其他任务
			cancel()
			return err
		}
	}
	return nil
}

// createTasks
// @Description:  数据插入任务列表
// @param chainId 链ID
// @param dealResult 待插入数据
// @return []saveTasks.Task 任务列表
func createTasks(chainId string, dealResult model.ProcessedBlockData) []saveTasks.Task {
	// 创建任务列表
	tasksList := []saveTasks.Task{
		// 保存交易任务
		saveTasks.TaskSaveTransactions{
			ChainId:            chainId,                      // 链ID
			Transactions:       dealResult.Transactions,      // 交易列表
			UpgradeContractTxs: dealResult.UpgradeContractTx, // 合约升级交易列表
		},
		// 插入用户任务
		saveTasks.TaskInsertUser{
			ChainId:  chainId,             // 链ID
			UserList: dealResult.UserList, // 用户列表
		},
		// 保存合约任务
		saveTasks.TaskSaveContract{
			ChainId:      chainId,                         // 链ID
			InsertList:   dealResult.InsertContracts,      // 插入合约列表
			UpdateList:   dealResult.UpdateContracts,      // 更新合约列表
			ByteCodeList: dealResult.ContractByteCodeList, // 合约字节码列表
		},
		// 保存标准合约任务
		saveTasks.TaskSaveStandardContract{
			ChainId:            chainId,                        // 链ID
			InsertFTContracts:  dealResult.FungibleContract,    // 可替换代币合约列表
			InsertNFTContracts: dealResult.NonFungibleContract, // 不可替换代币合约列表
			InsertIDAContracts: dealResult.InsertIDAContracts,  // IDA合约列表
		},
		// 保存证据合约任务
		saveTasks.TaskEvidenceContract{
			ChainId:           chainId,                 // 链ID
			EvidenceContracts: dealResult.EvidenceList, // 证据合约列表
		},
		// 插入合约事件任务
		saveTasks.TaskInsertContractEvents{
			ChainId:        chainId,                   // 链ID
			ContractEvents: dealResult.ContractEvents, // 合约事件列表
		},
		// 插入燃气记录任务
		saveTasks.TaskInsertGasRecord{
			ChainId:          chainId,                  // 链ID
			InsertGasRecords: dealResult.GasRecordList, // 燃气记录列表
		},
		// 保存链配置任务
		saveTasks.TaskSaveChainConfig{
			ChainId:            chainId,                    // 链ID
			UpdateChainConfigs: dealResult.ChainConfigList, // 链配置列表
		},
		// 保存跨链任务
		saveTasks.TaskSaveRelayCrossChain{
			ChainId:          chainId,                     // 链ID
			CrossChainResult: dealResult.CrossChainResult, // 跨链结果
		},
		// 保存跨链任务
		saveTasks.TaskSaveContractCrossCallTxs{
			ChainId:        chainId,                       // 链ID
			InsertCrossTxs: dealResult.ContractCrossTxs,   // 跨链结果
			InsertCalls:    dealResult.ContractCrossCalls, // 跨链结果
		},
	}

	return tasksList
}
