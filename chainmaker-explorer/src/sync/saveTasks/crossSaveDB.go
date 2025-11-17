package saveTasks

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/sync/model"
)

// TaskSaveRelayCrossChain
type TaskSaveRelayCrossChain struct {
	ChainId          string
	CrossChainResult *model.CrossChainResult
}

// Execute
// @description: 执行任务-保存合约数据
// @return error 错误信息
// SaveGasRecordToDB 保存gasR
func (e TaskSaveRelayCrossChain) Execute() error {
	return e.SaveRelayCrossChainToDB()
}

// SaveRelayCrossChainToDB 存储主子链数据
func (e TaskSaveRelayCrossChain) SaveRelayCrossChainToDB() error {
	if e.CrossChainResult == nil {
		return nil
	}

	chainId := e.ChainId
	crossChainResult := e.CrossChainResult
	var err error

	err = UpdateCrossTransfer(e.ChainId, e.CrossChainResult)
	if err != nil {
		return err
	}

	//存储跨链交易
	err = dbhandle.InsertCrossSubTransaction(chainId, crossChainResult.CrossMainTransaction)
	if err != nil {
		log.Errorf("SaveRelayCrossChainToDB insert tx failed, tx:%v",
			crossChainResult.CrossMainTransaction)
		return err
	}

	//业务交易数据
	err = SaveBusinessTransaction(chainId, crossChainResult.BusinessTxMap)
	if err != nil {
		return err
	}

	return nil
}

func UpdateCrossTransfer(chainId string, crossChainResult *model.CrossChainResult) error {
	insertTransferMap := crossChainResult.InsertCrossTransfer
	updateTransferMap := crossChainResult.UpdateCrossTransfer
	// 6. 批量插入CrossTransfer数据（预分配内存优化）
	transfers := make([]*db.CrossTransactionTransfer, 0, len(insertTransferMap))
	transfers = append(transfers, insertTransferMap...)
	if err := dbhandle.InsertCrossTxTransfers(chainId, transfers); err != nil {
		log.Errorf("Batch insert failed: %v", insertTransferMap)
		return err
	}

	// 5. 更新数据库中已存在记录的Status
	for _, transfer := range updateTransferMap {
		// 调用DAO层方法批量更新数据库状态
		err := dbhandle.UpdateCrossTxTransfers(chainId, transfer.CrossId, transfer.Status, transfer.EndTime)
		if err != nil {
			log.Errorf("UpdateCrossTransfer failed: %v", err)
		}
	}

	return nil
}

// GetSubChainSaveList 获取子链插入更新数据
func GetSubChainSaveList(chainId string, saveSubChainList map[string]*db.CrossSubChainData) (
	[]*db.CrossSubChainData, []*db.CrossSubChainData, error) {
	insertSubChainList := make([]*db.CrossSubChainData, 0)
	updateSubChainList := make([]*db.CrossSubChainData, 0)
	subChainIds := make([]string, 0)
	for _, subChainData := range saveSubChainList {
		subChainIds = append(subChainIds, subChainData.SubChainId)
	}

	crossSubChainDBMap, err := dbhandle.GetCrossSubChainById(chainId, subChainIds)
	if err != nil {
		return insertSubChainList, updateSubChainList, err
	}

	for _, subChainData := range saveSubChainList {
		if value, exists := crossSubChainDBMap[subChainData.SubChainId]; exists {
			subChainData.TxNum += value.TxNum
			updateSubChainList = append(updateSubChainList, subChainData)
		} else {
			insertSubChainList = append(insertSubChainList, subChainData)
		}
	}

	return insertSubChainList, updateSubChainList, nil
}

// SaveBusinessTransaction 保存主子链业务交易数据
func SaveBusinessTransaction(chainId string, businessTxMap map[string]*db.CrossBusinessTransaction) error {
	if len(businessTxMap) == 0 {
		return nil
	}
	insertTxList := make([]*db.CrossBusinessTransaction, 0)
	for _, txInfo := range businessTxMap {
		insertTxList = append(insertTxList, txInfo)
	}
	err := dbhandle.InsertCrossBusinessTransaction(chainId, insertTxList)
	if err != nil {
		log.Errorf("SaveBusinessTransaction failed, txlist:%v", insertTxList)
		return err
	}
	return nil
}

// SaveCrossSubChainCrossToDB 保存子链跨链数据
func SaveCrossSubChainCrossToDB(chainId string, inserts []*db.CrossSubChainCrossChain,
	updates []*db.CrossSubChainCrossChain) error {
	if len(inserts) <= 0 && len(updates) <= 0 {
		return nil
	}

	err := dbhandle.InsertCrossSubChainCross(chainId, inserts)
	if err != nil {
		log.Errorf("SaveCrossSubChainCrossToDB InsertCrossSubChainCross failed,data:%v ", inserts)
		return err
	}
	for _, subChainCross := range updates {
		err = dbhandle.UpdateCrossSubChainCross(chainId, subChainCross)
		if err != nil {
			log.Errorf("SaveCrossSubChainCrossToDB UpdateCrossSubChainCross failed,data:%v ", subChainCross)
			return err
		}
	}

	return nil
}

func SaveDelayCrossChainToDB(chainId string, delayCrossChain *model.DelayCrossChain) error {
	if delayCrossChain == nil {
		return nil
	}

	err := dbhandle.InsertCrossSubChainCross(chainId, delayCrossChain.InsertSubChainCross)
	if err != nil {
		log.Errorf("InsertCrossSubChainCross failed,data:%v ", delayCrossChain.InsertSubChainCross)
		return err
	}
	for _, subChainCross := range delayCrossChain.UpdateSubChainCross {
		err = dbhandle.UpdateCrossSubChainCross(chainId, subChainCross)
		if err != nil {
			log.Errorf("UpdateCrossSubChainCross failed,data:%v ", subChainCross)
			return err
		}
	}

	err = dbhandle.InsertCrossSubChain(chainId, delayCrossChain.InsertSubChainData)
	if err != nil {
		log.Errorf("InsertCrossSubChain failed, SubChain:%v",
			delayCrossChain.InsertSubChainData)
		return err
	}

	//更新子链网关数据
	for _, subChain := range delayCrossChain.UpdateSubChainData {
		err = dbhandle.UpdateCrossSubChainById(chainId, subChain)
		if err != nil {
			log.Errorf("UpdateCrossSubChainById failed subChain:%v ", subChain)
			return err
		}
	}

	//跨链合约
	err = dbhandle.InsertCrossContract(chainId, delayCrossChain.CrossChainContracts)
	if err != nil {
		log.Errorf("Failed to insert cross contract: %v", err)
		return err
	}

	return nil
}
