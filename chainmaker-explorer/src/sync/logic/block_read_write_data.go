package logic

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/sync/common"
	"chainmaker_web/src/sync/model"
	"encoding/base64"
	"encoding/json"
	"sync"

	pbCommon "chainmaker.org/chainmaker/pb-go/v2/common"
	pbConfig "chainmaker.org/chainmaker/pb-go/v2/config"
	tcipCommon "chainmaker.org/chainmaker/tcip-go/v2/common"
	"github.com/gogo/protobuf/proto"
	"github.com/panjf2000/ants/v2"
)

// CrossContract 跨链合约
type CrossContract struct {
	SubChainId   string
	ContractName string
}

// ParallelParseWriteSetData
// @Description: 并发解析读写集数据，获取链配置，主子链数据
// @param blockInfo 区块信息
// @param dealResult 处理结果
// @return error 错误信息
func ParallelParseWriteSetData(blockInfo *pbCommon.BlockInfo, dealResult *model.ProcessedBlockData) error {
	//主子链数据
	crossResult := model.NewCrossChainResult()
	// 使用同步互斥锁保护共享资源
	var mutx sync.Mutex
	var wg sync.WaitGroup

	// 创建一个固定大小的 goroutine 池
	goRoutinePool, err := ants.NewPool(10, ants.WithPreAlloc(false))
	if err != nil {
		log.Errorf("Failed to create goroutine pool: %v", err)
		return err
	}
	defer goRoutinePool.Release()

	chainId := blockInfo.Block.Header.ChainId
	//nolint:gosec
	blockHeight := int64(blockInfo.Block.Header.BlockHeight)
	// 用来接收并发任务的错误，减少通道的容量，避免阻塞
	errChan := make(chan error, 10)

	for i, tx := range blockInfo.Block.Txs {
		wg.Add(1)
		errSub := goRoutinePool.Submit(func(i int, blockInfo *pbCommon.BlockInfo, txInfo *pbCommon.Transaction) func() {
			return func() {
				defer wg.Done()
				// 处理跨链交易
				timestamp := txInfo.Payload.Timestamp
				txId := txInfo.Payload.TxId
				//读集
				// 确保 RwsetList 有足够的长度
				if i >= len(blockInfo.RwsetList) {
					//log.Errorf("Index out of range: i=%d, RwsetList length=%d", i, len(blockInfo.RwsetList))
					return
				}

				rwSetList := blockInfo.RwsetList[i]

				// 处理其他写集数据，包括链配置数据，主子链数据等
				if errOther := processWriteSetDataOther(&mutx, rwSetList, txInfo, dealResult); errOther != nil {
					errChan <- errOther
				}

				// 处理跨链写集数据，包括跨链交易数据，跨链合约数据等
				if errCross := ProcessCrossChainTransaction(&mutx, rwSetList, chainId, txId,
					blockHeight, timestamp, crossResult); errCross != nil {
					errChan <- errCross
				}
			}
		}(i, blockInfo, tx))
		if errSub != nil {
			log.Error("ParallelParseWriteSetData submit Failed: " + errSub.Error())
			wg.Done() // 减少 WaitGroup 计数
		}
	}

	wg.Wait()
	if len(errChan) > 0 {
		err = <-errChan
		return err
	}

	dealResult.CrossChainResult = crossResult
	return nil
}

// processWriteSetDataOther 处理其他写集数据
// @param mutx 互斥锁
// @param rwSetList 读写集列表
// @param txInfo 交易信息
// @param dealResult 处理结果
// @return error 错误信息
func processWriteSetDataOther(mutx *sync.Mutex, rwSetList *pbCommon.TxRWSet, txInfo *pbCommon.Transaction,
	dealResult *model.ProcessedBlockData) error {
	// 根据写集解析链配置数据
	chainConfig, err := GetChainConfigByWriteSet(rwSetList, txInfo)
	if err != nil {
		return err
	}

	// 锁定操作，避免并发修改
	mutx.Lock()
	defer mutx.Unlock()
	// 修改链配置
	if chainConfig != nil && chainConfig.ChainId != "" {
		dealResult.ChainConfigList = append(dealResult.ChainConfigList, chainConfig)
	}

	return nil
}

// processTransaction 处理跨链交易
// @param mutx 互斥锁
// @param rwSetList 读写集列表
// @param txInfo 交易信息
// @param chainId 链ID
// @param blockHeight 区块高度，用于版本控制
// @param crossResult 跨链结果
// @return error 错误信息
func ProcessCrossChainTransaction(mutx *sync.Mutex, rwSetList *pbCommon.TxRWSet,
	chainId, txId string, blockHeight, timestamp int64, crossResult *model.CrossChainResult) error {
	// 根据写集解析跨链交易数据
	crossChainInfo, crossErr := GetCrossChainInfoByWriteSet(rwSetList)
	if crossErr != nil {
		return crossErr
	}

	if crossChainInfo == nil {
		return nil
	}

	//跨链-主子链交易流转信息
	crossTxTransferList := GetCrossTxTransfer(chainId, txId, blockHeight, timestamp, crossChainInfo)

	//跨链-主链交易
	mainTransaction := GetMainCrossTransaction(blockHeight, crossChainInfo, txId, timestamp)

	//跨链-具体执行的交易数据
	businessTxMap := GetBusinessTransaction(chainId, crossChainInfo)
	// 锁定操作，避免并发修改
	mutx.Lock()
	defer mutx.Unlock()

	//更新跨链交易
	UpdateCrossChainResultTx(crossResult, mainTransaction, crossTxTransferList, businessTxMap)
	return nil
}

// UpdateCrossChainResultTx 更新跨链交易结果
// @param crossResult 跨链结果
// @param mainTransaction 主交易
// @param crossTxTransferList 跨链交易流转
// @param saveCycleTx 保存周期交易
// @param updateCycleTx 更新周期交易
// @param businessTxMap 业务交易
// @return error 错误信息
func UpdateCrossChainResultTx(crossResult *model.CrossChainResult, mainTransaction *db.CrossMainTransaction,
	crossTxTransferList []*db.CrossTransactionTransfer, businessTxMap map[string]*db.CrossBusinessTransaction) {
	// 跨链交易
	if mainTransaction != nil {
		crossResult.CrossMainTransaction = append(crossResult.CrossMainTransaction, mainTransaction)
	}

	// 跨链交易流转
	if len(crossTxTransferList) > 0 {
		for _, transfer := range crossTxTransferList {
			mapKey := transfer.CrossId + "_" + transfer.FromChainId + "_" + transfer.ToChainId
			if value, ok := crossResult.CrossTransfer[mapKey]; ok {
				if transfer.Status > value.Status {
					value.Status = transfer.Status
					value.EndTime = transfer.EndTime
				}
			} else {
				crossResult.CrossTransfer[mapKey] = transfer
			}
		}
	}

	// 跨链-具体执行的交易
	if len(businessTxMap) > 0 {
		for mapKey, executionTx := range businessTxMap {
			if _, ok := crossResult.BusinessTxMap[mapKey]; !ok {
				crossResult.BusinessTxMap[mapKey] = executionTx
			}
		}
	}
}

// GetChainConfigByWriteSet
//
//	@Description: 通过读写集解析链配置
//	@param txRWSet
//	@param txInfo
//	@return *pbConfig.ChainConfig
//	@return error
func GetChainConfigByWriteSet(txRWSet *pbCommon.TxRWSet, txInfo *pbCommon.Transaction) (*pbConfig.ChainConfig, error) {
	chainConfig := &pbConfig.ChainConfig{}
	//是否配置类交易
	isConfigTx := common.IsConfigTx(txInfo)
	if !isConfigTx || txRWSet == nil {
		return nil, nil
	}

	for _, write := range txRWSet.TxWrites {
		if string(write.Key) == common.TxReadWriteKeyChainConfig {
			err := proto.Unmarshal(write.Value, chainConfig)
			if err != nil {
				return nil, err
			}
			break
		}
	}
	if chainConfig.ChainId != "" {
		return chainConfig, nil
	}
	return nil, nil
}

// GetCrossChainInfoByWriteSet
//
//	@Description: 通过读写集解析跨链数据
//	@param txRWSet
//	@param txInfo
func GetCrossChainInfoByWriteSet(txRWSet *pbCommon.TxRWSet) (*tcipCommon.CrossChainInfo, error) {
	crossChainInfo := &tcipCommon.CrossChainInfo{}
	if txRWSet == nil {
		return nil, nil
	}

	for _, write := range txRWSet.TxWrites {
		//以c开头的key可以解析成CrossChainInfo
		isCrossChainTx := common.IsRelayCrossChainInfo(string(write.Key), write.ContractName)
		if !isCrossChainTx {
			continue
		}

		// Base64 解码
		decoded, err := base64.StdEncoding.DecodeString(string(write.Value))
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(decoded, crossChainInfo)
		if err != nil {
			return nil, err
		}
		break
	}

	if crossChainInfo.CrossChainId != "" {
		return crossChainInfo, nil
	}
	return nil, nil
}
