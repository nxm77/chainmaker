package saveTasks

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/sync/model"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type Task interface {
	Execute() error
}

// TaskFunc 是一个类型别名，表示任务函数的类型。
type TaskFunc func(...interface{}) error

// Task 任务
type DealTask struct {
	Name     string
	Function TaskFunc
	Args     []interface{}
}

// RetrySleepTime 重试等待时间
func RetrySleepTime(retryCount int) int {
	if retryCount < 3 {
		return 0
	}

	return retryCount - 3
}

// WithRetry 执行任务，失败后重试
func WithRetry(task func() error, logFuncName string, errCh chan<- error) {
	retryCount := 0
	for {
		err := task()
		if err == nil {
			break
		}
		retryCount++
		log.Errorf("dealBlockData %v-[%d] failed, err:%v", logFuncName, retryCount, err)
		if retryCount > config.MaxRetryCount {
			errCh <- fmt.Errorf("dealBlockData Error: %v, Retry count: %d", err, retryCount)
			break
		}
		//重试总失败就先等待一下
		sleepTime := time.Duration(RetrySleepTime(retryCount))
		time.Sleep(time.Second * sleepTime)
	}
}

// ExecuteTaskWithRetry 会执行给定的任务，并在出错时重试。
// 任务将继续重试，直到成功或达到 maxRetryCount 为止。
func ExecuteTaskWithRetry(ctx context.Context, wg *sync.WaitGroup, task DealTask, retryCountMap *sync.Map,
	errCh chan<- error) {
	defer wg.Done()
	for {
		err := task.Function(task.Args...)
		if err == nil {
			return
		}
		retryCount, _ := retryCountMap.LoadOrStore(task.Name, 0)
		retryCountInt, ok := retryCount.(int)
		if !ok {
			errCh <- fmt.Errorf("ExecuteTaskWithRetry task %s failed err: %v", task.Name, err)
			return
		}

		retryCountMap.Store(task.Name, retryCountInt+1)
		if retryCountInt >= config.MaxRetryCount {
			// 将错误发送到错误通道
			errCh <- fmt.Errorf("ExecuteTaskWithRetry task %s failed err: %v", task.Name, err)
			return
		}

		select {
		case <-ctx.Done():
			// 上下文已取消
			return
		default:
			// 继续重试
			log.Errorf("ExecuteTaskWithRetry task[%v] run failed err:%v", task.Name, err)
			//重试总失败就先等待一下
			sleepTime := time.Duration(RetrySleepTime(retryCountInt))
			time.Sleep(time.Second * sleepTime)
		}
	}
}

//-------异步更新--------

// TaskUpdateTxBlackToDB 更新交易黑名单数据
func TaskUpdateTxBlackToDB(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskContractToDB: expected string for args[0], got %T", args[0])
	}
	updateTxBlack, ok := args[1].(*db.UpdateTxBlack)
	if !ok {
		return fmt.Errorf("TaskUpdateTxBlackToDB: expected string for args[1], got %T", args[1])
	}

	return UpdateTxBlackToDB(chainId, updateTxBlack)
}

// TaskUpdateContractResult 更新合约数据
func TaskUpdateContractResult(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskUpdateContractResult: expected string for args[0], got %T", args[0])
	}
	contractResult, ok := args[1].(*db.GetContractResult)
	if !ok {
		return fmt.Errorf("TaskUpdateContractResult: expected string for args[1], got %T", args[1])
	}

	//更新合约交易量
	err := UpdateContractTxNum(chainId, contractResult.UpdateContractTxEventNum)
	if err != nil {
		return err
	}

	//更新IDA合约资产数量
	err = UpdateIDAContract(chainId, contractResult.UpdateIdaContract)
	if err != nil {
		return err
	}

	//更新IDA合约资产数量
	err = SaveEventTopicData(chainId, contractResult.InsertEventTopic, contractResult.UpdateEventTopic)
	if err != nil {
		return err
	}

	return UpdateContract(chainId, contractResult.IdentityContract)
}

// TaskInsertFungibleTransferToDB 保留流转数据
func TaskInsertFungibleTransferToDB(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskInsertFungibleTransferToDB: expected string for args[0], got %T", args[0])
	}
	transferList, ok := args[1].([]*db.FungibleTransfer)
	if !ok {
		return fmt.Errorf("TaskInsertFungibleTransferToDB: expected string for args[1], got %T", args[1])
	}

	return InsertFungibleTransferToDB(chainId, transferList)
}

// TaskInsertNonFungibleTransferToDB 保留流转数据
func TaskInsertNonFungibleTransferToDB(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskInsertNonFungibleTransferToDB: expected string for args[0], got %T", args[0])
	}
	transferList, ok := args[1].([]*db.NonFungibleTransfer)
	if !ok {
		return fmt.Errorf("TaskInsertNonFungibleTransferToDB: expected string for args[1], got %T", args[1])
	}
	return InsertNonFungibleTransferToDB(chainId, transferList)
}

// TaskSaveAccountListToDB 保存账户数据
func TaskSaveAccountListToDB(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskSaveContractResult: expected string for args[0], got %T", args[0])
	}
	updateAccountResult, ok := args[1].(*db.UpdateAccountResult)
	if !ok {
		return fmt.Errorf("TaskSaveAccountListToDB: expected string for args[1], got %T", args[1])
	}

	return SaveAccountToDB(chainId, updateAccountResult)
}

// TaskSaveTokenResultToDB 保存非同质化token
func TaskSaveTokenResultToDB(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskSaveContractResult: expected string for args[0], got %T", args[0])
	}
	tokenResult, ok := args[1].(*db.TokenResult)
	if !ok {
		return fmt.Errorf("TaskSaveTokenResultToDB: expected string for args[1], got %T", args[1])
	}

	return SaveNonFungibleToken(chainId, tokenResult)
}

func TaskSaveGasToDB(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskSaveGasToDB: expected string for args[0], got %T", args[0])
	}
	insertGasList, ok := args[1].([]*db.Gas)
	if !ok {
		return fmt.Errorf("TaskSaveGasToDB: expected string for args[1], got %T", args[1])
	}
	updateGasList, ok := args[2].([]*db.Gas)
	if !ok {
		return fmt.Errorf("TaskSaveGasToDB: expected string for args[1], got %T", args[1])
	}

	err := InsertGasToDB(chainId, insertGasList)
	if err != nil {
		return err
	}
	err = UpdateGasToDB(chainId, updateGasList)
	if err != nil {
		return err
	}
	return nil
}

func TaskSaveFungibleContractResult(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskSaveFungibleContractResult: expected string for args[0], got %T", args[0])
	}
	contractResult, ok := args[1].(*db.GetContractResult)
	if !ok {
		return fmt.Errorf("TaskSaveFungibleContractResult: expected string for args[1], got %T", args[1])
	}

	return SaveFungibleContractResult(chainId, contractResult)
}

func TaskSavePositionToDB(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskSavePositionToDB: expected string for args[0], got %T", args[0])
	}
	blockPosition, ok := args[1].(*db.BlockPosition)
	if !ok {
		return fmt.Errorf("TaskSavePositionToDB: expected string for args[1], got %T", args[1])
	}

	return SavePositionToDB(chainId, blockPosition)
}

func TaskCrossSubChainCrossToDB(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskSavePositionToDB: expected string for args[0], got %T", args[0])
	}
	insertList, ok := args[1].([]*db.CrossSubChainCrossChain)
	if !ok {
		return fmt.Errorf("TaskCrossSubChainCrossToDB: expected string for args[1], got %T", args[1])
	}
	updateList, ok := args[2].([]*db.CrossSubChainCrossChain)
	if !ok {
		return fmt.Errorf("TaskCrossSubChainCrossToDB: expected string for args[2], got %T", args[2])
	}

	return SaveCrossSubChainCrossToDB(chainId, insertList, updateList)
}

func TaskDelayCrossChain(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskDelayCrossChain: expected string for args[0], got %T", args[0])
	}
	delayCrossChain, ok := args[1].(*model.DelayCrossChain)
	if !ok {
		return fmt.Errorf("TaskDelayCrossChain: expected string for args[1], got %T", args[1])
	}

	return SaveDelayCrossChainToDB(chainId, delayCrossChain)
}

func TaskSaveIDAAssetDataToDB(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskSaveIDAAssetDataToDB: expected string for args[0], got %T", args[0])
	}
	idaAssetsDataDB, ok := args[1].(*db.IDAAssetsDataDB)
	if !ok {
		return fmt.Errorf("TaskSaveIDAAssetDataToDB: expected string for args[1], got %T", args[1])
	}

	return SaveIDAAssetDataToDB(chainId, idaAssetsDataDB)
}

func TaskUpdateIDAAssetDataToDB(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskUpdateIDAAssetDataToDB: expected string for args[0], got %T", args[0])
	}
	idaAssetsDataDB, ok := args[1].(*db.IDAAssetsUpdateDB)
	if !ok {
		return fmt.Errorf("TaskUpdateIDAAssetDataToDB: expected string for args[1], got %T", args[1])
	}

	return UpdateIDAAssetDataToDB(chainId, idaAssetsDataDB)
}

func TaskUpdateDIDDataToDB(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskUpdateDIDDataToDB: expected string for args[0], got %T", args[0])
	}
	didSaveDate, ok := args[1].(*db.DIDSaveData)
	if !ok {
		return fmt.Errorf("TaskUpdateDIDDataToDB: expected string for args[1], got %T", args[1])
	}

	err := UpdateDIDDataToDB(chainId, didSaveDate)
	if err != nil {
		didSaveDateJson, _ := json.Marshal(didSaveDate)
		log.Errorf("TaskUpdateDIDDataToDB: %s didSaveDate: %v", err.Error, string(didSaveDateJson))
	}
	return err
}

func TaskSaveChainStatisticsToDB(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskSaveChainStatisticsToDB: expected string for args[0], got %T", args[0])
	}
	starts, ok := args[1].(*db.Statistics)
	if !ok {
		return fmt.Errorf("TaskSaveChainStatisticsToDB: expected string for args[1], got %T", args[1])
	}

	return SaveChainStatistics(chainId, starts)
}

func TaskSaveABITopicTableEvents(args ...interface{}) error {
	chainId, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("TaskSaveABITopicTableEvents: expected string for args[0], got %T", args[0])
	}
	topicTableEvents, ok := args[1].(map[string][]map[string]interface{})
	if !ok {
		return fmt.Errorf("TaskSaveABITopicTableEvents: expected string for args[1], got %T", args[1])
	}

	return SaveABITopicTableEvents(chainId, topicTableEvents)
}
