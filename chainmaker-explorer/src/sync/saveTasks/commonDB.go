package saveTasks

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	loggers "chainmaker_web/src/logger"
	"encoding/json"
	"unicode/utf8"
)

const (
	SaveBatchSize    = 500
	GoRoutinePoolErr = "new ants pool error: "
)

var (
	log = loggers.GetLogger(loggers.MODULE_WEB)
)

// calculateDataSize
//
//	@Description: 计算存续数据的字节数
//	@param data 存储数据
//	@return int 字节数
func calculateDataSize(data interface{}) int {
	bytes, _ := json.Marshal(data)
	return utf8.RuneCount(bytes)
}

// batchTransactions
//
//	@Description: 将交易分割为大小为batchSize的批次
//	@param transactions 交易列表
//	@return [][]*db.Transaction 交易批次
func batchTransactions(transactions map[string]*db.Transaction) [][]*db.Transaction {
	batches := make([][]*db.Transaction, 0)
	batch := make([]*db.Transaction, 0)
	batchSize := 0
	for _, transaction := range transactions {
		transactionSize := calculateDataSize(transaction)
		if batchSize+transactionSize > config.MaxDBByteSize {
			batches = append(batches, batch)
			batch = make([]*db.Transaction, 0)
			batchSize = 0
		}

		batch = append(batch, transaction)
		batchSize += transactionSize
	}

	if len(batch) > 0 {
		batches = append(batches, batch)
	}

	return batches
}

// 将交易分割为大小为batchSize的批次
func batchUsers(users map[string]*db.User) [][]*db.User {
	batches := make([][]*db.User, 0)
	batch := make([]*db.User, 0)

	for _, user := range users {
		batch = append(batch, user)
		if len(batch) == SaveBatchSize {
			batches = append(batches, batch)
			batch = make([]*db.User, 0)
		}
	}

	if len(batch) > 0 {
		batches = append(batches, batch)
	}

	return batches
}

// 将交易分割为大小为batchSize的批次
func batchContractEvents(contractEvents []*db.ContractEvent) [][]*db.ContractEvent {
	batches := make([][]*db.ContractEvent, 0)
	batch := make([]*db.ContractEvent, 0)
	batchSize := 0

	for _, contractEvent := range contractEvents {
		//计算内存大小
		contractEventSize := calculateDataSize(contractEvent)
		if batchSize+contractEventSize > config.MaxDBByteSize {
			batches = append(batches, batch)
			batch = make([]*db.ContractEvent, 0)
			batchSize = 0
		}

		batch = append(batch, contractEvent)
		batchSize += contractEventSize
	}

	if len(batch) > 0 {
		batches = append(batches, batch)
	}

	return batches
}

// 将交易分割为大小为batchSize的批次
func batchGasRecords(gasRecords []*db.GasRecord) [][]*db.GasRecord {
	batches := make([][]*db.GasRecord, 0)
	batch := make([]*db.GasRecord, 0)
	batchSize := 0

	for _, gasRecord := range gasRecords {
		//计算内存大小
		gasRecordSize := calculateDataSize(gasRecord)
		if batchSize+gasRecordSize > config.MaxDBByteSize {
			batches = append(batches, batch)
			batch = make([]*db.GasRecord, 0)
			batchSize = 0
		}

		batch = append(batch, gasRecord)
		batchSize += gasRecordSize
	}

	if len(batch) > 0 {
		batches = append(batches, batch)
	}

	return batches
}

// 将交易分割为大小为batchSize的批次
func batchFungibleTransfers(transferList []*db.FungibleTransfer) [][]*db.FungibleTransfer {
	batches := make([][]*db.FungibleTransfer, 0)
	batch := make([]*db.FungibleTransfer, 0)
	batchSize := 0

	for _, transfer := range transferList {
		//计算内存大小
		transferSize := calculateDataSize(transfer)
		if batchSize+transferSize > config.MaxDBByteSize {
			batches = append(batches, batch)
			batch = make([]*db.FungibleTransfer, 0)
			batchSize = 0
		}

		batch = append(batch, transfer)
		batchSize += transferSize
	}

	if len(batch) > 0 {
		batches = append(batches, batch)
	}

	return batches
}

// 将交易分割为大小为batchSize的批次
func batchNonFungibleTransfers(transferList []*db.NonFungibleTransfer) [][]*db.NonFungibleTransfer {
	batches := make([][]*db.NonFungibleTransfer, 0)
	batch := make([]*db.NonFungibleTransfer, 0)
	batchSize := 0

	for _, transfer := range transferList {
		//计算内存大小
		transferSize := calculateDataSize(transfer)
		if batchSize+transferSize > config.MaxDBByteSize {
			batches = append(batches, batch)
			batch = make([]*db.NonFungibleTransfer, 0)
			batchSize = 0
		}

		batch = append(batch, transfer)
		batchSize += transferSize
	}

	if len(batch) > 0 {
		batches = append(batches, batch)
	}

	return batches
}

func batchNonFungibleToken(tokenList []*db.NonFungibleToken) [][]*db.NonFungibleToken {
	batches := make([][]*db.NonFungibleToken, 0)
	batch := make([]*db.NonFungibleToken, 0)
	batchSize := 0

	for _, token := range tokenList {
		//计算内存大小
		tokenSize := calculateDataSize(token)
		if batchSize+tokenSize > config.MaxDBByteSize {
			batches = append(batches, batch)
			batch = make([]*db.NonFungibleToken, 0)
			batchSize = 0
		}

		batch = append(batch, token)
		batchSize += tokenSize
	}

	if len(batch) > 0 {
		batches = append(batches, batch)
	}

	return batches
}

// batchContractCallTxs函数用于将传入的ContractCrossCallTransaction切片按照SaveBatchSize大小进行分批
func batchContractCallTxs(inserts []*db.ContractCrossCallTransaction) [][]*db.ContractCrossCallTransaction {
	// 创建一个空的二维切片，用于存储分批后的ContractCrossCallTransaction切片
	batches := make([][]*db.ContractCrossCallTransaction, 0)
	// 创建一个空的切片，用于存储当前批次的ContractCrossCallTransaction
	batch := make([]*db.ContractCrossCallTransaction, 0)

	// 遍历传入的ContractCrossCallTransaction切片
	for _, tx := range inserts {
		// 将当前ContractCrossCallTransaction添加到当前批次的切片中
		batch = append(batch, tx)
		// 如果当前批次的切片大小达到了SaveBatchSize，则将该批次添加到二维切片中，并清空当前批次的切片
		if len(batch) == SaveBatchSize {
			batches = append(batches, batch)
			batch = make([]*db.ContractCrossCallTransaction, 0)
		}
	}

	// 如果当前批次的切片不为空，则将该批次添加到二维切片中
	if len(batch) > 0 {
		batches = append(batches, batch)
	}

	// 返回分批后的二维切片
	return batches
}
