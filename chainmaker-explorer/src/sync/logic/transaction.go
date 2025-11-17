/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package logic

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/sync/common"
	"chainmaker_web/src/sync/model"
	"chainmaker_web/src/utils"
	"encoding/hex"
	"encoding/json"
	"strings"
	"sync"

	pbCommon "chainmaker.org/chainmaker/pb-go/v2/common"
	"github.com/panjf2000/ants/v2"
)

const MD5Str = "md5:"

// ParallelParseTransactions
//
//	@Description: 并发解析所有交易数据
//	@param blockInfo
//	@param hashType
//	@param dealResult
//	@return var
//	@return err
func ParallelParseTransactions(blockInfo *pbCommon.BlockInfo, hashType string,
	dealResult *model.ProcessedBlockData) error {
	var wg sync.WaitGroup
	// 使用同步互斥锁保护共享资源
	var mutx sync.Mutex

	// 创建一个固定大小的 goroutine 池
	goRoutinePool, err := ants.NewPool(10, ants.WithPreAlloc(false))
	if err != nil {
		log.Errorf("Failed to create goroutine pool: %v", err)
		return err
	}
	defer goRoutinePool.Release()

	chainId := blockInfo.Block.Header.ChainId
	errChan := make(chan error, 10) // 用来接收并发任务的错误，减少通道的容量，避免阻塞
	// 并发提交任务到 goroutine 池中
	for i, tx := range blockInfo.Block.Txs {
		wg.Add(1)

		// 使用闭包提交任务，避免外部变量污染
		err := goRoutinePool.Submit(func(i int, blockInfo *pbCommon.BlockInfo, txInfo *pbCommon.Transaction) func() {
			return func() {
				defer wg.Done()

				// 处理每笔交易的业务逻辑
				if err := processTransaction(i, &mutx, blockInfo, txInfo, chainId, hashType, dealResult); err != nil {
					errChan <- err
				}
			}
		}(i, blockInfo, tx))
		if err != nil {
			log.Errorf("Submit failed for transaction %d: %v", i, err)
			wg.Done()
		}
	}

	// 等待所有并发任务完成
	wg.Wait()

	// 关闭 errChan 并返回错误，如果有的话
	close(errChan)

	// 获取第一个错误
	for e := range errChan {
		if e != nil {
			return e
		}
	}

	return nil
}

// processTransaction 处理每一笔交易
// @param i 交易索引
// @param mutx 互斥锁
// @param blockInfo 区块信息
// @param txInfo 交易信息
// @param chainId 链ID
// @param hashType 哈希类型
// @param dealResult 处理结果
// @return error
func processTransaction(i int, mutx *sync.Mutex, blockInfo *pbCommon.BlockInfo, txInfo *pbCommon.Transaction,
	chainId, hashType string, dealResult *model.ProcessedBlockData) error {
	// 解析所有合约事件
	tempContractEvents := DealContractEvents(txInfo)

	// 计算账户信息
	userResult, senderErr := GetSenderAndPayerUser(chainId, hashType, txInfo)
	if senderErr != nil {
		log.Errorf("ParallelParseTransactions get User err:%v", senderErr)
	}

	// 构建 Gas 记录
	tempGasRecords, gasErr := buildGasRecord(txInfo, userResult)
	if gasErr != nil {
		return gasErr
	}

	// 构建交易数据
	transaction, tranErr := buildTransaction(i, blockInfo, txInfo, userResult, hashType)
	if tranErr != nil {
		return tranErr
	}

	//跨合约调用关系
	contractCrossCalls, contractCrossTxs := DealCrossContractCallTxs(chainId, transaction)

	// 锁定互斥锁
	mutx.Lock()
	// 使用 defer 确保互斥锁被解锁
	defer mutx.Unlock()
	// 并发写入 dealResult
	if transaction.TxId != "" {
		dealResult.Transactions[transaction.TxId] = transaction
	}

	if len(contractCrossTxs) > 0 {
		// 合约跨链调用关系
		dealResult.ContractCrossTxs = append(dealResult.ContractCrossTxs, contractCrossTxs...)
	}

	if len(contractCrossCalls) > 0 {
		// 合约跨链调用关系
		for key, crossCall := range contractCrossCalls {
			if _, ok := dealResult.ContractCrossCalls[key]; !ok {
				dealResult.ContractCrossCalls[key] = crossCall
			}
		}
	}

	// 合约事件
	if len(tempContractEvents) > 0 {
		dealResult.ContractEvents = append(dealResult.ContractEvents, tempContractEvents...)
	}

	// gas记录
	if len(tempGasRecords) > 0 {
		dealResult.GasRecordList = append(dealResult.GasRecordList, tempGasRecords...)
	}

	// 用户列表
	if userResult.SenderUserAddr != "" {
		if _, ok := dealResult.UserList[userResult.SenderUserAddr]; !ok {
			userInfo := &db.User{
				UserId:    userResult.SenderUserId,
				UserAddr:  userResult.SenderUserAddr,
				Role:      userResult.SenderRole,
				OrgId:     userResult.SenderOrgId,
				Timestamp: txInfo.Payload.Timestamp,
			}
			dealResult.UserList[userResult.SenderUserAddr] = userInfo
		}
	}

	if userResult.PayerUserAddr != "" {
		if _, ok := dealResult.UserList[userResult.PayerUserAddr]; !ok {
			userInfo := &db.User{
				UserId:    userResult.PayerUserId,
				UserAddr:  userResult.PayerUserAddr,
				Timestamp: txInfo.Payload.Timestamp,
			}
			dealResult.UserList[userResult.PayerUserAddr] = userInfo
		}
	}

	return nil
}

// 构造合约升级交易数据
func buildUpgradeContractTransaction(contractWriteSetData *model.ContractWriteSetData) *db.UpgradeContractTransaction {
	// 如果传入的合约写集数据为空，则返回nil
	if contractWriteSetData == nil {
		return nil
	}

	//构造合约升级交易数据
	upgradeContractTransaction := &db.UpgradeContractTransaction{
		TxId:                contractWriteSetData.SenderTxId,      // 交易ID
		SenderOrgId:         contractWriteSetData.SenderOrgId,     // 发送者组织ID
		Sender:              contractWriteSetData.Sender,          // 发送者
		UserAddr:            contractWriteSetData.SenderAddr,      // 发送者地址
		BlockHeight:         contractWriteSetData.BlockHeight,     // 区块高度
		BlockHash:           contractWriteSetData.BlockHash,       // 区块哈希
		Timestamp:           contractWriteSetData.Timestamp,       // 时间戳
		ContractRuntimeType: contractWriteSetData.RuntimeType,     // 合约运行时类型
		ContractName:        contractWriteSetData.ContractName,    // 合约名称
		ContractNameBak:     contractWriteSetData.ContractNameBak, // 合约名称备份
		ContractAddr:        contractWriteSetData.ContractAddr,    // 合约地址
		ContractVersion:     contractWriteSetData.Version,         // 合约版本
		ContractType:        contractWriteSetData.ContractType,    // 合约类型
	}

	return upgradeContractTransaction
}

// buildTransaction 构建交易数据
// @param i 交易索引
// @param blockInfo 区块信息
// @param txInfo 交易信息
// @param userResult 用户信息
// @return *db.Transaction 交易数据
// @return error 错误信息
func buildTransaction(i int, blockInfo *pbCommon.BlockInfo, txInfo *pbCommon.Transaction,
	userResult *db.SenderPayerUser, hashType string) (*db.Transaction, error) {
	payload := txInfo.Payload
	contractNameAddr := payload.ContractName
	chainId := payload.ChainId

	//构造交易数据
	transaction := &db.Transaction{
		TxId:    payload.TxId,
		TxIndex: i + 1,
		TxType:  payload.TxType.String(),
		//nolint:gosec
		BlockHeight:        int64(blockInfo.Block.Header.BlockHeight),
		BlockHash:          hex.EncodeToString(blockInfo.Block.Header.BlockHash),
		ContractMessage:    txInfo.Result.ContractResult.Message,
		GasUsed:            txInfo.Result.ContractResult.GasUsed,
		Sequence:           payload.Sequence,
		ContractResult:     txInfo.Result.ContractResult.Result,
		ContractResultCode: txInfo.Result.ContractResult.Code,
		ExpirationTime:     payload.ExpirationTime,
		RwSetHash:          hex.EncodeToString(txInfo.Result.RwSetHash),
		Timestamp:          payload.Timestamp,
		TimestampDate:      utils.GetDateFromTimestamp(payload.Timestamp),
		ChainTimestamp:     blockInfo.Block.Header.BlockTimestamp,
		TxStatusCode:       txInfo.Result.Code.String(),
		ContractMethod:     payload.Method,
	}

	member := blockInfo.Block.Header.Proposer
	//根据proposer信息填充block中的地址,id信息
	if member != nil {
		//根据proposer信息填充block中的地址,id信息
		getInfos, err := common.GetMemberIdAddrAndCertNew(chainId, hashType, member)
		if err != nil {
			log.Error("GetMemberIdAddrAndCertNew Failed: " + err.Error())
		}
		transaction.ProposerId = getInfos.UserId
	}

	// 创建一个新的参数列表来存储处理后的参数1
	newParameters := make([]*pbCommon.KeyValuePair, len(payload.Parameters))
	// 对每一个参数进行深拷贝
	for i, originalParam := range payload.Parameters {
		// 通过拷贝KeyValuePair的值来实现深拷贝
		copiedParam := &pbCommon.KeyValuePair{
			Key:   originalParam.Key,
			Value: append([]byte(nil), originalParam.Value...), // 深拷贝Value的字节切片
		}
		newParameters[i] = copiedParam
	}

	//对参数进行MD5处理
	for _, parameter := range newParameters {
		switch parameter.Key {
		case "CONTRACT_NAME":
			transaction.ContractName = string(parameter.Value)
			transaction.ContractNameBak = string(parameter.Value)
		case "CONTRACT_VERSION":
			transaction.ContractVersion = string(parameter.Value)
		case "CONTRACT_RUNTIME_TYPE":
			transaction.ContractRuntimeType = string(parameter.Value)
		case "CONTRACT_BYTECODE":
			parameter.Value = []byte(MD5Str + common.MD5(string(parameter.Value)))
		case "UPGRADE_CONTRACT_BYTECODE":
			parameter.Value = []byte(MD5Str + common.MD5(string(parameter.Value)))
		}
	}

	//解析合约名称
	transaction.ContractName = contractNameAddr
	transaction.ContractNameBak = contractNameAddr
	//解析参数
	parametersBytes, err := json.Marshal(newParameters)
	if err == nil {
		transaction.ContractParameters = string(parametersBytes)
	}

	//解析读写集
	transaction.ReadSet, transaction.WriteSet = buildReadWriteSet(*blockInfo.RwsetList[i])
	if userResult != nil {
		transaction.Sender = userResult.SenderUserId
		transaction.SenderOrgId = userResult.SenderOrgId
		transaction.UserAddr = userResult.SenderUserAddr
		transaction.PayerAddr = userResult.PayerUserAddr
	}

	//解析背书信息
	if len(txInfo.Endorsers) > 0 {
		endorsementBytes, _ := json.Marshal(txInfo.Endorsers)
		transaction.Endorsement = string(endorsementBytes)
	}

	//解析事件
	if len(txInfo.Result.ContractResult.ContractEvent) > 0 {
		eventList := make([]config.TxEventData, 0)
		for k, event := range txInfo.Result.ContractResult.ContractEvent {
			eventList = append(eventList, config.TxEventData{
				Index:        k,
				ContractName: event.ContractName,
				Key:          event.Topic,
				Value:        strings.Join(event.EventData, ","),
			})
		}
		eventByte, _ := json.Marshal(eventList)
		transaction.Event = string(eventByte)
	}

	return transaction, nil
}

// buildReadWriteSet 解析读写集，数据过长也不在进行截断
func buildReadWriteSet(rwSetList pbCommon.TxRWSet) (string, string) {
	readList := make([]config.RwSet, 0)
	for j, read := range rwSetList.TxReads {
		if strings.HasPrefix(string(read.Key), "ContractByteCode:") {
			read.Value = []byte(MD5Str + common.MD5(string(read.Value)))
		}
		readList = append(readList, config.RwSet{
			Index:        j,
			Key:          string(read.Key),
			Value:        read.Value,
			ContractName: read.ContractName,
		})
	}

	writeList := make([]config.RwSet, 0)
	for j, write := range rwSetList.TxWrites {
		if strings.HasPrefix(string(write.Key), "ContractByteCode:") {
			write.Value = []byte(MD5Str + common.MD5(string(write.Value)))
		}

		writeList = append(writeList, config.RwSet{
			Index:        j,
			Key:          string(write.Key),
			Value:        write.Value,
			ContractName: write.ContractName,
		})
	}
	readByte, _ := json.Marshal(readList)
	writeByte, _ := json.Marshal(writeList)
	return string(readByte), string(writeByte)
}
