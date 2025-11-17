package logic

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/entity_cross"
	"chainmaker_web/src/sync/common"
	"encoding/hex"
	"encoding/json"
	"strconv"

	"github.com/google/uuid"

	pbCommon "chainmaker.org/chainmaker/pb-go/v2/common"
	tcipCommon "chainmaker.org/chainmaker/tcip-go/v2/common"
	"github.com/gogo/protobuf/proto"
)

// CrossChainMsgExtra 跨链交易额外信息
type CrossChainMsgExtra struct {
	SyncData            string `json:"SyncData"`
	SyncDataBatchCount  string `json:"SyncDataBatchCount"`
	SyncFromBlockHeight string `json:"SyncFromBlockHeight"`
	SyncToBlockHeight   string `json:"SyncToBlockHeight"`
}

// GetCrossModelByExtra 根据extraData获取跨链模型
func GetCrossModelByExtra(crossChainInfo *tcipCommon.CrossChainInfo) int32 {
	crossModel := entity_cross.CrossModelOther
	chainMsgExtra := GetCrossChainMsgExtra(crossChainInfo)
	if chainMsgExtra.SyncData == "true" {
		crossModel = entity_cross.CrossModelSync
	}
	return crossModel
}

// GetCrossChainMsgExtra 根据跨链信息获取跨链交易额外信息
func GetCrossChainMsgExtra(crossChainInfo *tcipCommon.CrossChainInfo) *CrossChainMsgExtra {
	chainMsgExtra := &CrossChainMsgExtra{}
	if crossChainInfo == nil {
		return chainMsgExtra
	}
	for _, chainMsg := range crossChainInfo.CrossChainMsg {
		if chainMsg.ExtraData != "" {
			err := json.Unmarshal([]byte(chainMsg.ExtraData), &chainMsgExtra)
			if err != nil {
				log.Warnf("[GetMainCrossTransaction] json.Unmarshal ExtraData failed, err:%v", err)
			}
		}
	}
	return chainMsgExtra
}

// GetMainCrossTransaction 获取主链跨链交易
func GetMainCrossTransaction(blockHeight int64, crossChainInfo *tcipCommon.CrossChainInfo,
	txId string, timestamp int64) *db.CrossMainTransaction {
	if crossChainInfo == nil {
		return nil
	}

	crossModel := GetCrossModelByExtra(crossChainInfo)
	crossChainMsg, _ := json.Marshal(crossChainInfo.CrossChainMsg)
	//跨链交易
	crossTransaction := &db.CrossMainTransaction{
		TxId:        txId,
		CrossId:     crossChainInfo.CrossChainId,
		ChainMsg:    string(crossChainMsg),
		CrossModel:  crossModel,
		Status:      int32(crossChainInfo.State),
		BlockHeight: blockHeight,
		Timestamp:   timestamp,
	}
	return crossTransaction
}

// GetBusinessTransaction 跨链-具体执行的交易
func GetBusinessTransaction(chainId string,
	crossChainInfo *tcipCommon.CrossChainInfo) map[string]*db.CrossBusinessTransaction {
	executionTxMap := make(map[string]*db.CrossBusinessTransaction, 0)
	if crossChainInfo == nil {
		return executionTxMap
	}

	//跨链周期结束在解析业务交易数据，防止重复解析
	if crossChainInfo.State != tcipCommon.CrossChainStateValue_CONFIRM_END &&
		crossChainInfo.State != tcipCommon.CrossChainStateValue_CANCEL_END {
		return executionTxMap
	}

	crossModel := GetCrossModelByExtra(crossChainInfo)
	if crossModel == entity_cross.CrossModelSync {
		return executionTxMap
	}

	//子链交易
	crossId := crossChainInfo.CrossChainId
	if crossChainInfo.FirstTxContent != nil && crossChainInfo.FirstTxContent.TxContent != nil {
		txContent := crossChainInfo.FirstTxContent.TxContent
		if txContent.TxId != "" {
			//解析交易数据
			executionTx := BuildExecutionTransaction(txContent)
			isMainChain := common.IsMainChainGateway(txContent.GatewayId)
			subChainId := txContent.ChainRid
			executionTx.IsMainChain = isMainChain
			executionTx.CrossId = crossId
			executionTx.SubChainId = subChainId
			executionTx.GatewayId = txContent.GatewayId
			executionTx.TxId = txContent.TxId
			executionTx.TxStatus = int32(txContent.TxResult)
			mapKey := crossId + "_" + subChainId
			if _, ok := executionTxMap[mapKey]; !ok {
				executionTxMap[executionTx.TxId] = executionTx
			}
		}
	}
	if len(crossChainInfo.CrossChainTxContent) > 0 {
		for _, txContent := range crossChainInfo.CrossChainTxContent {
			if txContent.TxContent == nil || txContent.TxContent.TxId == "" {
				continue
			}

			//解析交易数据
			executionTx := BuildExecutionTransaction(txContent.TxContent)
			isMainChain := common.IsMainChainGateway(txContent.TxContent.GatewayId)
			subChainId := txContent.TxContent.ChainRid
			executionTx.IsMainChain = isMainChain
			executionTx.CrossId = crossId
			executionTx.SubChainId = subChainId
			executionTx.GatewayId = txContent.TxContent.GatewayId

			contractRes, _ := json.Marshal(txContent.TryResult)
			executionTx.CrossContractResult = string(contractRes)
			executionTx.TxId = txContent.TxContent.TxId
			executionTx.TxStatus = int32(txContent.TxContent.TxResult)
			mapKey := crossId + "_" + subChainId
			if _, ok := executionTxMap[mapKey]; !ok {
				executionTxMap[mapKey] = executionTx
			}
		}
	}

	return executionTxMap
}

// GetCrossTxTransfer
//
//	@Description: 跨链-交易流转方向
//	@param chainId 主链id
//	@param blockHeight 主链高度
//	@param crossChainInfo 跨链详情
//	@return []*db.CrossTransactionTransfer 跨链交易流转方向
func GetCrossTxTransfer(chainId, txId string, blockHeight, timestamp int64,
	crossChainInfo *tcipCommon.CrossChainInfo) []*db.CrossTransactionTransfer {
	if crossChainInfo == nil {
		return nil
	}

	//跨链流转交易
	var transferList []*db.CrossTransactionTransfer
	crossModel := GetCrossModelByExtra(crossChainInfo)
	if crossModel == entity_cross.CrossModelOther {
		transferList = buildCrossModelOtherTransferData(txId, crossChainInfo, blockHeight, timestamp)
	} else {
		transferList = buildCrossModelSyncTransferData(txId, crossChainInfo, blockHeight, timestamp)
	}

	return transferList
}

func buildCrossModelSyncTransferData(txId string, crossChainInfo *tcipCommon.CrossChainInfo, blockHeight,
	timestamp int64) []*db.CrossTransactionTransfer {
	transferList := make([]*db.CrossTransactionTransfer, 0)
	if crossChainInfo == nil ||
		len(crossChainInfo.CrossChainMsg) == 0 {
		return transferList
	}

	var fromSubChainId string
	if crossChainInfo.FirstTxContent != nil &&
		crossChainInfo.FirstTxContent.TxContent != nil {
		fromSubChainId = crossChainInfo.FirstTxContent.TxContent.ChainRid
	}

	isMainChain := common.IsMainChainGateway(crossChainInfo.From)
	newUUID := uuid.New().String()
	transfer := &db.CrossTransactionTransfer{
		ID:              newUUID,
		BlockHeight:     blockHeight,
		CrossId:         crossChainInfo.CrossChainId,
		TxId:            txId,
		FromChainId:     fromSubChainId,
		FromIsMainChain: isMainChain,
		FromGatewayId:   crossChainInfo.From,
	}

	for _, chainMsg := range crossChainInfo.CrossChainMsg {
		chainMsgExtra := &CrossChainMsgExtra{}
		if chainMsg.ExtraData != "" {
			err := json.Unmarshal([]byte(chainMsg.ExtraData), &chainMsgExtra)
			if err != nil {
				log.Warnf("[GetMainCrossTransaction] json.Unmarshal ExtraData failed, err:%v", err)
			}
		}

		var txNum int64
		var syncToBlockHeight int64
		if chainMsgExtra.SyncDataBatchCount != "" {
			txNum, _ = strconv.ParseInt(chainMsgExtra.SyncDataBatchCount, 10, 64)
		}
		if chainMsgExtra.SyncToBlockHeight != "" {
			syncToBlockHeight, _ = strconv.ParseInt(chainMsgExtra.SyncToBlockHeight, 10, 64)
		}

		transfer.FromBlockHeight = syncToBlockHeight
		isMainChain = common.IsMainChainGateway(chainMsg.GatewayId)
		transfer.ToChainId = chainMsg.ChainRid
		transfer.ToIsMainChain = isMainChain
		transfer.ToGatewayId = chainMsg.GatewayId
		transfer.ContractName = chainMsg.ContractName
		transfer.ContractMethod = chainMsg.Method
		transfer.TxNum = txNum
		transfer.StartTime = timestamp
		transfer.EndTime = timestamp
		transfer.Status = int32(crossChainInfo.State)
		transfer.CrossModel = entity_cross.CrossModelSync
		transferList = append(transferList, transfer)
	}

	return transferList

}

// buildCrossModelOtherTransferData 普通模式跨链交易数据
func buildCrossModelOtherTransferData(txId string, crossChainInfo *tcipCommon.CrossChainInfo, blockHeight,
	timestamp int64) []*db.CrossTransactionTransfer {
	transferList := make([]*db.CrossTransactionTransfer, 0)
	if crossChainInfo == nil ||
		crossChainInfo.FirstTxContent == nil ||
		crossChainInfo.FirstTxContent.TxContent == nil {
		return transferList
	}

	fromTxContent := crossChainInfo.FirstTxContent.TxContent
	crossChainTxContent := crossChainInfo.CrossChainTxContent
	isMainChain := common.IsMainChainGateway(fromTxContent.GatewayId)
	fromSubChainId := fromTxContent.ChainRid
	newUUID := uuid.New().String()
	transfer := &db.CrossTransactionTransfer{
		ID:              newUUID,
		BlockHeight:     blockHeight,
		CrossId:         crossChainInfo.CrossChainId,
		TxId:            txId,
		FromChainId:     fromSubChainId,
		FromIsMainChain: isMainChain,
		FromGatewayId:   fromTxContent.GatewayId,
		FromBlockHeight: fromTxContent.BlockHeight,
	}

	for i := 0; i < len(crossChainInfo.CrossChainMsg); i++ {
		var toBlockHeight int64
		if i < len(crossChainTxContent) && crossChainTxContent[i] != nil {
			toBlockHeight = crossChainTxContent[i].TxContent.BlockHeight
		}

		chainMsg := crossChainInfo.CrossChainMsg[i]
		isMainChain = common.IsMainChainGateway(chainMsg.GatewayId)
		transfer.ToChainId = chainMsg.ChainRid
		transfer.ToIsMainChain = isMainChain
		transfer.ToGatewayId = chainMsg.GatewayId
		transfer.ToBlockHeight = toBlockHeight
		transfer.ContractName = chainMsg.ContractName
		transfer.ContractMethod = chainMsg.Method
		transfer.Parameter = chainMsg.Parameter
		transfer.TxNum = 1
		transfer.StartTime = timestamp
		transfer.EndTime = timestamp
		transfer.Status = int32(crossChainInfo.State)
		transfer.CrossModel = entity_cross.CrossModelOther

		transferList = append(transferList, transfer)
	}
	return transferList
}

// BuildExecutionTransaction
//
//	@Description: 将跨链交易内容解析成交易信息
//	@param txContent
//	@return *db.CrossBusinessTransaction
//
// 构建执行交易
func BuildExecutionTransaction(txContent *tcipCommon.TxContent) *db.CrossBusinessTransaction {
	// 定义交易信息
	var txInfo pbCommon.Transaction
	// 生成新的UUID
	newUUID := uuid.New().String()
	// 初始化执行交易
	executionTx := &db.CrossBusinessTransaction{
		ID: newUUID,
	}
	// 如果交易内容为空，直接返回执行交易
	if txContent == nil || len(txContent.Tx) == 0 {
		return executionTx
	}

	// 解析交易内容
	err := proto.Unmarshal(txContent.Tx, &txInfo)
	// 获取交易负载
	payload := txInfo.Payload
	// 如果解析失败，记录错误并返回执行交易
	if err != nil || payload == nil {
		log.Errorf("BuildExecutionTransaction txContent json Unmarshal failed, err:%v", err)
		return executionTx
	}

	//构造交易数据
	executionTx = &db.CrossBusinessTransaction{
		ID:                 newUUID,
		TxId:               payload.TxId,
		TxType:             payload.TxType.String(),
		ContractMessage:    txInfo.Result.ContractResult.Message,
		GasUsed:            txInfo.Result.ContractResult.GasUsed,
		ContractResult:     txInfo.Result.ContractResult.Result,
		ContractResultCode: txInfo.Result.ContractResult.Code,
		RwSetHash:          hex.EncodeToString(txInfo.Result.RwSetHash),
		Timestamp:          payload.Timestamp,
		TxStatusCode:       txInfo.Result.Code.String(),
		ContractMethod:     payload.Method,
		ContractName:       payload.ContractName,
	}
	// 将交易参数转换为JSON字符串
	parametersBytes, err := json.Marshal(payload.Parameters)
	// 如果转换成功，将参数赋值给执行交易
	if err == nil {
		executionTx.ContractParameters = string(parametersBytes)
	}

	return executionTx
}

// ParseCrossCycleTxTransfer 解析跨链交易转移
func ParseCrossCycleTxTransfer(transfers []*db.CrossTransactionTransfer) map[string]map[string]int64 {
	// 创建一个map，用于存储子链ID和对应的转移次数
	subChainIdMap := make(map[string]map[string]int64, 0)
	// 如果转移列表为空，直接返回空map
	if len(transfers) == 0 {
		return subChainIdMap
	}

	// 遍历转移列表
	for _, transfer := range transfers {
		// 如果转移的源链ID不为空
		if transfer.FromChainId != "" {
			// 如果子链IDMap中不存在源链ID，则创建一个新的map
			if _, ok := subChainIdMap[transfer.FromChainId]; !ok {
				subChainIdMap[transfer.FromChainId] = make(map[string]int64, 0)
			}
			// 如果转移的目标链ID不为空，则将源链ID和目标链ID对应的转移次数加1
			if transfer.ToChainId != "" {
				subChainIdMap[transfer.FromChainId][transfer.ToChainId]++
			}
		}

		// 如果转移的目标链ID不为空
		if transfer.ToChainId != "" {
			// 如果子链IDMap中不存在目标链ID，则创建一个新的map
			if _, ok := subChainIdMap[transfer.ToChainId]; !ok {
				subChainIdMap[transfer.ToChainId] = make(map[string]int64, 0)
			}
			if transfer.FromChainId != "" {
				subChainIdMap[transfer.ToChainId][transfer.FromChainId]++
			}
		}
	}

	return subChainIdMap
}

// DealSubChainCrossChainNum
//
//	@Description: 根据子链跨链流转计算子链跨链交易数量明细
//	@param chainId
//	@param subChainIdMap 本次子链跨链数据
//	@param subChainCrossDB 数据库子链跨链数据
//	@param minHeight 本次批量处理最低区块高度
//	@return []*db.CrossSubChainCrossChain 新增跨链数据
//	@return []*db.CrossSubChainCrossChain  更新跨链数据
//	@return error
func DealSubChainCrossChainNum(chainId string, subChainIdMap map[string]map[string]int64,
	subChainCrossDB []*db.CrossSubChainCrossChain, minHeight int64) ([]*db.CrossSubChainCrossChain,
	[]*db.CrossSubChainCrossChain, error) {
	insertSubChainCross := make([]*db.CrossSubChainCrossChain, 0)
	updateSubChainCross := make([]*db.CrossSubChainCrossChain, 0)
	if len(subChainIdMap) == 0 {
		return insertSubChainCross, updateSubChainCross, nil
	}

	crossSubChainDBMap := make(map[string]map[string]*db.CrossSubChainCrossChain, 0)
	for _, subChain := range subChainCrossDB {
		if _, ok := crossSubChainDBMap[subChain.SubChainId]; !ok {
			crossSubChainDBMap[subChain.SubChainId] = make(map[string]*db.CrossSubChainCrossChain, 0)
		}
		crossSubChainDBMap[subChain.SubChainId][subChain.CrossChainId] = subChain
	}

	for subChainId, relatedChains := range subChainIdMap {
		for relatedChainId, txNum := range relatedChains {
			subDBMap, subChainExists := crossSubChainDBMap[subChainId]
			subDB, relatedChainExists := subDBMap[relatedChainId]
			if !subChainExists || !relatedChainExists {
				newUUID := uuid.New().String()
				crossSubChain := &db.CrossSubChainCrossChain{
					ID:           newUUID,
					SubChainId:   subChainId,
					CrossChainId: relatedChainId,
					TxNum:        txNum,
					BlockHeight:  minHeight,
				}
				insertSubChainCross = append(insertSubChainCross, crossSubChain)
			} else if subDB.BlockHeight < minHeight {
				//如果数据库中BlockHeight，大于等于BlockHeight最小值，说明已经更新过了
				// 更新数据
				subDB.TxNum += txNum
				subDB.BlockHeight = minHeight
				updateSubChainCross = append(updateSubChainCross, subDB)
			}
		}
	}

	return insertSubChainCross, updateSubChainCross, nil
}

type ChainTransfer struct {
	TxCount     int64
	BlockHeight int64
	Timestamp   int64
}

// DealCrossSubChainData 处理跨链子链数据
func DealCrossSubChainData(insertCrossTransfers []*db.CrossTransactionTransfer,
	subChainDataDB map[string]*db.CrossSubChainData) (
	inserts []*db.CrossSubChainData,
	updates []*db.CrossSubChainData,
	contracts []*db.CrossChainContract,
) {

	// 0. 预计算链ID的交易次数（消除O(n²)循环）
	chainTxCounts := make(map[string]*ChainTransfer)
	for _, transfer := range insertCrossTransfers {
		if transfer.FromChainId != "" {
			if _, exists := chainTxCounts[transfer.FromChainId]; !exists {
				chainTxCounts[transfer.FromChainId] = &ChainTransfer{}
			}
			chainTxCounts[transfer.FromChainId].TxCount++
			if transfer.BlockHeight > chainTxCounts[transfer.FromChainId].BlockHeight {
				chainTxCounts[transfer.FromChainId].BlockHeight = transfer.BlockHeight
				chainTxCounts[transfer.FromChainId].Timestamp = transfer.StartTime
			}
		}
		if transfer.ToChainId != "" {
			if _, exists := chainTxCounts[transfer.ToChainId]; !exists {
				chainTxCounts[transfer.ToChainId] = &ChainTransfer{}
			}
			chainTxCounts[transfer.ToChainId].TxCount++
			if transfer.BlockHeight > chainTxCounts[transfer.ToChainId].BlockHeight {
				chainTxCounts[transfer.ToChainId].BlockHeight = transfer.BlockHeight
				chainTxCounts[transfer.ToChainId].Timestamp = transfer.StartTime
			}
		}

		// 1. 独立处理合约逻辑（避免嵌套）
		if transfer.ToChainId != "" && transfer.ContractName != "" {
			contracts = append(contracts, &db.CrossChainContract{
				ID:           uuid.New().String(),
				SubChainId:   transfer.ToChainId,
				ContractName: transfer.ContractName,
			})
		}
	}

	inserts, updates = BuildCrossSubChainByTransfer(insertCrossTransfers, subChainDataDB, chainTxCounts)
	return
}

// BuildCrossSubChainByTransfer 根据转账构建跨链子链数据
func BuildCrossSubChainByTransfer(insertTransfers []*db.CrossTransactionTransfer,
	subChainDataDB map[string]*db.CrossSubChainData,
	chainTxCounts map[string]*ChainTransfer) (inserts []*db.CrossSubChainData, updates []*db.CrossSubChainData) {
	// 已处理过的链
	processedChains := make(map[string]bool)
	for _, transfer := range insertTransfers {
		// 如果转账的源链ID不为空
		if transfer.FromChainId != "" {
			// 如果源链已经被处理过，则跳过
			if processedChains[transfer.FromChainId] {
				continue
			}
			// 标记源链已经被处理过
			processedChains[transfer.FromChainId] = true

			// 3. 分离数据构建逻辑
			// 如果源链数据已经存在
			if subChain, exists := subChainDataDB[transfer.FromChainId]; exists {
				// 如果源链的交易数量存在
				if value, ok := chainTxCounts[transfer.FromChainId]; ok {
					// 源链的交易数量增加
					subChain.TxNum += value.TxCount
					// 如果源链的区块高度小于转账的源链区块高度
					if subChain.BlockHeight < transfer.FromBlockHeight {
						// 源链的区块高度更新为转账的源链区块高度
						subChain.BlockHeight = transfer.FromBlockHeight
						// 源链的时间戳更新为转账的时间戳
						subChain.Timestamp = transfer.StartTime
					}
					// 将更新后的源链数据添加到更新列表中
					updates = append(updates, subChain)
				}
			} else {
				// 如果源链数据不存在
				// 如果源链的交易数量存在
				if value, ok := chainTxCounts[transfer.FromChainId]; ok {
					// 创建一个新的源链数据
					insert := &db.CrossSubChainData{
						SubChainId: transfer.FromChainId,
						TxNum:      value.TxCount,
						Timestamp:  transfer.StartTime,
					}
					// 源链的网关ID更新为转账的源链网关ID
					insert.GatewayId = transfer.FromGatewayId
					insert.IsMainChain = transfer.FromIsMainChain
					insert.BlockHeight = transfer.FromBlockHeight
					inserts = append(inserts, insert)
				}
			}
		}
		if transfer.ToChainId != "" {
			if processedChains[transfer.ToChainId] {
				continue
			}
			processedChains[transfer.ToChainId] = true

			// 3. 分离数据构建逻辑
			if subChain, exists := subChainDataDB[transfer.ToChainId]; exists {
				if value, ok := chainTxCounts[transfer.ToChainId]; ok {
					subChain.TxNum += value.TxCount
					if subChain.BlockHeight < transfer.ToBlockHeight {
						subChain.BlockHeight = transfer.ToBlockHeight
						subChain.Timestamp = transfer.StartTime
					}
					updates = append(updates, subChain)
				}
			} else {
				if value, ok := chainTxCounts[transfer.ToChainId]; ok {
					insert := &db.CrossSubChainData{
						SubChainId: transfer.ToChainId,
						TxNum:      value.TxCount,
						Timestamp:  transfer.StartTime,
					}
					insert.GatewayId = transfer.ToGatewayId
					insert.IsMainChain = transfer.ToIsMainChain
					insert.BlockHeight = transfer.ToBlockHeight
					inserts = append(inserts, insert)
				}
			}
		}
	}
	return inserts, updates
}

// GetCrossChainAndContracts 获取跨链子链和合约
// 根据跨链交易列表获取跨链合约和跨链子链数据
func GetCrossChainAndContracts(crossTxTransferList []*db.CrossTransactionTransfer) (
	map[string]*db.CrossSubChainData, []*db.CrossChainContract) {
	// 初始化跨链合约列表
	crossChainContracts := make([]*db.CrossChainContract, 0)
	// 初始化跨链子链数据
	crossChains := make(map[string]*db.CrossSubChainData, 0)
	//跨链合约
	for _, transfer := range crossTxTransferList {
		crossChainContracts = append(crossChainContracts, &db.CrossChainContract{
			ID:           uuid.New().String(),
			SubChainId:   transfer.ToChainId,
			ContractName: transfer.ContractName,
		})

		if transfer.ToChainId != "" {
			if _, ok := crossChains[transfer.ToChainId]; !ok {
				crossChains[transfer.ToChainId] = &db.CrossSubChainData{
					SubChainId:  transfer.ToChainId,
					GatewayId:   transfer.ToGatewayId,
					BlockHeight: transfer.ToBlockHeight,
					IsMainChain: transfer.ToIsMainChain,
					Timestamp:   transfer.StartTime,
					TxNum:       1,
				}
			} else {
				crossChains[transfer.ToChainId].TxNum++
				if crossChains[transfer.ToChainId].BlockHeight < transfer.ToBlockHeight {
					crossChains[transfer.ToChainId].BlockHeight = transfer.ToBlockHeight
					crossChains[transfer.ToChainId].Timestamp = transfer.StartTime
				}
			}
		}

		if transfer.FromChainId != "" {
			if _, ok := crossChains[transfer.FromChainId]; !ok {
				crossChains[transfer.FromChainId] = &db.CrossSubChainData{
					SubChainId:  transfer.FromChainId,
					GatewayId:   transfer.FromGatewayId,
					BlockHeight: transfer.FromBlockHeight,
					IsMainChain: transfer.FromIsMainChain,
					Timestamp:   transfer.StartTime,
					TxNum:       1,
				}
			} else {
				crossChains[transfer.FromChainId].TxNum++
				if crossChains[transfer.FromChainId].BlockHeight < transfer.FromBlockHeight {
					crossChains[transfer.FromChainId].BlockHeight = transfer.FromBlockHeight
					crossChains[transfer.FromChainId].Timestamp = transfer.StartTime
				}
			}
		}
	}

	return crossChains, crossChainContracts
}
