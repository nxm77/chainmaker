/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package logic

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/sync/common"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/google/uuid"

	pbCommon "chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	"chainmaker.org/chainmaker/sdk-go/v2/utils"
	"github.com/panjf2000/ants/v2"
)

const (
	//BusinessTypeRecharge gas充值
	BusinessTypeRecharge = 1
	//BusinessTypeConsume gas消费
	BusinessTypeConsume = 2
)

// buildGasRecord
//
//	@Description: 根据交易数据构造gas消耗列表
//	@param txInfo 交易列表
//	@param userResult 交易用户
//	@return []*db.GasRecord gas消耗
//	@return error
func buildGasRecord(txInfo *pbCommon.Transaction, userResult *db.SenderPayerUser) ([]*db.GasRecord, error) {
	gasRecords := make([]*db.GasRecord, 0)
	payload := txInfo.Payload
	txId := payload.TxId

	//gas充值
	if payload.Method == syscontract.GasAccountFunction_RECHARGE_GAS.String() {
		if txInfo.Result.Code != pbCommon.TxStatusCode_SUCCESS {
			return gasRecords, nil
		}

		req := &syscontract.RechargeGasReq{}
		for _, parameter := range payload.Parameters {
			if parameter.Key == utils.KeyGasBatchRecharge {
				err := req.Unmarshal(parameter.Value)
				if err != nil {
					return gasRecords, err
				}
				break
			}
		}

		if len(req.BatchRechargeGas) > 0 {
			for i, gasReq := range req.BatchRechargeGas {
				newUUID := uuid.New().String()
				gasInfo := &db.GasRecord{
					ID:           newUUID,
					Address:      gasReq.Address,
					GasAmount:    gasReq.GasAmount,
					BusinessType: BusinessTypeRecharge,
					Timestamp:    payload.Timestamp,
					TxId:         txId,
					GasIndex:     i + 1,
				}
				gasRecords = append(gasRecords, gasInfo)
			}
		}
	} else {
		//gas消费
		newUUID := uuid.New().String()
		gasInfo := &db.GasRecord{
			ID:       newUUID,
			TxId:     txId,
			GasIndex: 1,
			//nolint:gosec
			GasAmount:    int64(txInfo.Result.ContractResult.GasUsed),
			BusinessType: BusinessTypeConsume,
			Timestamp:    payload.Timestamp,
		}

		if userResult != nil {
			if userResult.PayerUserAddr != "" {
				//PayerUserAddr代付
				gasInfo.Address = userResult.PayerUserAddr
			} else {
				gasInfo.Address = userResult.SenderUserAddr
			}
		}
		if gasInfo.GasAmount == 0 || gasInfo.Address == "" {
			return gasRecords, nil
		}
		gasRecords = append(gasRecords, gasInfo)
	}

	return gasRecords, nil
}

// buildGasInfo
//
//	@Description: 根据gas消耗列表，和DB中gas余额计算新的gas余额
//	@param gasRecords gas消耗
//	@param gasInfoList DB中gas余额
//	@param minHeight 批量处理最小高度，用作版本号，避免重复计算
//	@return []*db.Gas 新增gas用户
//	@return []*db.Gas 更新gas用户
func BuildGasInfo(gasRecords []*db.GasRecord, gasInfoList []*db.Gas, minHeight int64) ([]*db.Gas, []*db.Gas) {
	var (
		gasUseAmount   = make(map[string]int64)
		gasTotalAmount = make(map[string]int64)
		addrMap        = make(map[string]string, 0)
		insertGas      = make([]*db.Gas, 0)
		updateGas      = make([]*db.Gas, 0)
	)
	for _, gasInfo := range gasRecords {
		addr := gasInfo.Address
		if _, ok := addrMap[addr]; !ok {
			addrMap[addr] = addr
		}

		if gasInfo.BusinessType == BusinessTypeRecharge {
			//gas充值
			if amount, ok := gasTotalAmount[addr]; ok {
				gasTotalAmount[addr] = amount + gasInfo.GasAmount
			} else {
				gasTotalAmount[addr] = gasInfo.GasAmount
			}
		} else {
			//gas消耗
			if amount, ok := gasUseAmount[addr]; ok {
				gasUseAmount[addr] = amount + gasInfo.GasAmount
			} else {
				gasUseAmount[addr] = gasInfo.GasAmount
			}
		}
	}

	if len(addrMap) == 0 {
		return insertGas, updateGas
	}

	// 创建一个映射来存储gasInfoList中的地址
	gasInfoMap := make(map[string]*db.Gas)
	for _, gas := range gasInfoList {
		gasInfoMap[gas.Address] = gas
	}

	for _, addr := range addrMap {
		gas, okMap := gasInfoMap[addr]
		if okMap {
			//数据库存在
			if gas.BlockHeight >= minHeight {
				//如果数据库中BlockHeight，大于等于BlockHeight最小值，说明已经更新过了
				continue
			}

			if amount, ok := gasUseAmount[gas.Address]; ok {
				gas.GasUsed = gas.GasUsed + amount
			}
			if amount, ok := gasTotalAmount[gas.Address]; ok {
				gas.GasTotal = gas.GasTotal + amount
			}
			gas.GasBalance = gas.GasTotal - gas.GasUsed
			gas.BlockHeight = minHeight
			updateGas = append(updateGas, gas)
		} else {
			// 将不存在的地址添加到InsertGas
			gas = &db.Gas{
				Address:     addr,
				BlockHeight: minHeight,
			}
			if amount, ok := gasUseAmount[gas.Address]; ok {
				gas.GasUsed = amount
			}
			if amount, ok := gasTotalAmount[gas.Address]; ok {
				gas.GasTotal = amount
			}
			gas.GasBalance = gas.GasTotal - gas.GasUsed
			insertGas = append(insertGas, gas)
		}
	}

	return insertGas, updateGas
}

// GetGasRecord
//
//	@Description: 数据库获取gas记录
//	@param chainId
//	@param txIds 交易ID列表
//	@return []*db.GasRecord  gas消耗
//	@return error
func GetGasRecord(chainId string, txIds []string) ([]*db.GasRecord, error) {
	gasRecords := make([]*db.GasRecord, 0)
	if len(txIds) == 0 {
		return gasRecords, nil
	}
	var (
		goRoutinePool *ants.Pool
		mutx          sync.Mutex
		err           error
	)
	// 将交易分割为大小为10的批次
	batches := common.ParallelParseBatchWhere(txIds, 100)
	errChan := make(chan error, 10)
	if goRoutinePool, err = ants.NewPool(10, ants.WithPreAlloc(false)); err != nil {
		log.Error(common.GoRoutinePoolErr + err.Error())
		return gasRecords, err
	}

	defer goRoutinePool.Release()
	var wg sync.WaitGroup
	for _, batch := range batches {
		wg.Add(1)
		// 使用匿名函数传递参数，避免使用外部变量
		// 并发插入每个批次
		errSub := goRoutinePool.Submit(func(batch []string) func() {
			return func() {
				defer wg.Done()
				//查询数据
				gasList, eventErr := dbhandle.GetGasRecordByTxIds(chainId, txIds)
				if eventErr != nil {
					errChan <- eventErr
				}
				mutx.Lock()         // 锁定互斥锁
				defer mutx.Unlock() // 使用 defer 确保互斥锁被解锁
				if len(gasList) > 0 {
					gasRecords = append(gasRecords, gasList...)
				}
			}
		}(batch))
		if errSub != nil {
			log.Error("GetGasRecord submit Failed : " + err.Error())
			wg.Done() // 减少 WaitGroup 计数
		}
	}

	wg.Wait()
	if len(errChan) > 0 {
		err = <-errChan
		return gasRecords, err
	}

	return gasRecords, nil

}

// buildGasAddrList
//
//	@Description: 计算需要更新的GAS数据
//	@param gasRecords 本次gas消耗列表
//	@return []string 计算后的gas余额
func BuildGasAddrList(gasRecords []*db.GasRecord) []string {
	//获取gas余额
	addrMap := make(map[string]string, 0)
	addrList := make([]string, 0)
	for _, gasInfo := range gasRecords {
		addr := gasInfo.Address
		if _, ok := addrMap[addr]; !ok {
			addrMap[addr] = addr
		}
	}
	for _, addr := range addrMap {
		addrList = append(addrList, addr)
	}

	return addrList
}

// BuildAccountManagerGasTransfer 函数用于构建账户管理器燃料转账信息
func BuildAccountManagerGasTransfer(chainId string, txInfoMap map[string]*db.Transaction,
	fungibleTransfers []*db.FungibleTransfer) []*db.FungibleTransfer {
	// 如果交易信息为空，则直接返回
	if len(txInfoMap) == 0 {
		return fungibleTransfers
	}

	// 遍历交易信息
	for _, txInfo := range txInfoMap {
		transferInfo := BuildGasTransfer(txInfo)
		if transferInfo != nil {
			// 将转账信息添加到转账列表中
			fungibleTransfers = append(fungibleTransfers, transferInfo)
		}

		// 构建gas充值记录
		gasRecords := BuildGasRecordTransfer(txInfo)
		if len(gasRecords) > 0 {
			fungibleTransfers = append(fungibleTransfers, gasRecords...)
		}
	}

	return fungibleTransfers
}

func BuildGasRecordTransfer(txInfo *db.Transaction) []*db.FungibleTransfer {
	result := make([]*db.FungibleTransfer, 0)
	// 如果交易结果码不为成功，则跳过
	if txInfo.ContractResultCode != dbhandle.ContractResultSuccess {
		return result
	}

	// gas充值
	if txInfo.ContractName != syscontract.SystemContract_ACCOUNT_MANAGER.String() ||
		txInfo.ContractMethod != syscontract.GasAccountFunction_RECHARGE_GAS.String() {
		return result
	}

	// 解析合约参数
	var parameters []*pbCommon.KeyValuePair
	if txInfo.ContractParameters != "" {
		_ = json.Unmarshal([]byte(txInfo.ContractParameters), &parameters)
	}

	// 解析合约参数
	req := &syscontract.RechargeGasReq{}
	for _, parameter := range parameters {
		switch parameter.Key {
		case utils.KeyGasBatchRecharge:
			err := req.Unmarshal(parameter.Value)
			if err != nil {
				log.Error("BuildAccountManagerGasTransfer Unmarshal failed: " + err.Error())
				return result
			}
		}
		// 解析参数中的转账地址和金额
		if len(req.BatchRechargeGas) > 0 {
			for _, gasReq := range req.BatchRechargeGas {
				if gasReq.Address == "" || gasReq.GasAmount == 0 {
					log.Error("BuildAccountManagerGasTransfer toAddr or amount is empty")
					continue
				}
				// 将金额转换为小数形式
				amountStr := fmt.Sprintf("%d", gasReq.GasAmount)
				amountDecimal := common.StringAmountDecimal(amountStr, 8)
				// 构建转账信息
				transferInfo := &db.FungibleTransfer{
					ID:             uuid.New().String(),
					TxId:           txInfo.TxId,
					ContractName:   txInfo.ContractName,
					ContractAddr:   txInfo.ContractAddr,
					ContractMethod: txInfo.ContractMethod,
					Topic:          txInfo.ContractMethod,
					ToAddr:         gasReq.Address,
					Amount:         amountDecimal,
					Timestamp:      txInfo.Timestamp,
				}
				result = append(result, transferInfo)
			}
		}
	}

	return result
}

func BuildGasTransfer(txInfo *db.Transaction) *db.FungibleTransfer {
	// 如果交易结果码不为成功，则跳过
	if txInfo.ContractResultCode != dbhandle.ContractResultSuccess {
		return nil
	}

	// gas转账
	if txInfo.ContractName != syscontract.SystemContract_ACCOUNT_MANAGER.String() ||
		txInfo.ContractMethod != "TRANSFER_GAS" {
		return nil
	}

	// 解析合约参数
	var parameters []*pbCommon.KeyValuePair
	if txInfo.ContractParameters != "" {
		_ = json.Unmarshal([]byte(txInfo.ContractParameters), &parameters)
	}

	// 解析参数中的转账地址和金额
	var toAddr string
	var amount string
	for _, parameter := range parameters {
		switch parameter.Key {
		case "amount":
			amount = string(parameter.Value)
		case "to":
			toAddr = string(parameter.Value)
		}
	}
	// 如果转账地址或金额为空，则跳过
	if toAddr == "" || amount == "" {
		log.Error("BuildAccountManagerGasTransfer toAddr or amount is empty")
		return nil
	}

	// 将金额转换为小数形式
	amountDecimal := common.StringAmountDecimal(amount, 8)
	// 构建转账信息
	transferInfo := &db.FungibleTransfer{
		ID:             uuid.New().String(),
		TxId:           txInfo.TxId,
		ContractName:   txInfo.ContractName,
		ContractAddr:   txInfo.ContractAddr,
		ContractMethod: txInfo.ContractMethod,
		Topic:          txInfo.ContractMethod,
		FromAddr:       txInfo.UserAddr,
		ToAddr:         toAddr,
		Amount:         amountDecimal,
		Timestamp:      txInfo.Timestamp,
	}
	return transferInfo
}
