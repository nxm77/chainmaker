/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package logic

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/sync/model"
)

// ContractHandler 处理合约相关的所有业务逻辑
type ContractHandler struct {
	ChainId        string
	MinHeight      int64
	TxList         map[string]*db.Transaction
	ContractMap    map[string]*db.Contract
	ContractEvents []*db.ContractEvent
	EventResults   *model.TopicEventResult
}

// NewContractHandler 构造函数，初始化ContractHandler
func NewContractHandler(chainId string, minHeight int64, contractMap map[string]*db.Contract,
	eventResults *model.TopicEventResult, delayedUpdateCache *model.GetRealtimeCacheData) *ContractHandler {
	return &ContractHandler{
		ChainId:        chainId,
		MinHeight:      minHeight,
		TxList:         delayedUpdateCache.TxList,
		ContractMap:    contractMap,
		ContractEvents: delayedUpdateCache.ContractEvents,
		EventResults:   eventResults,
	}
}

// DealWithContractData 处理合约相关数据
func (ch *ContractHandler) DealWithContractData(idaContractMap map[string]*db.IDAContract) (
	*db.GetContractResult, error) {
	// 这里处理所有与合约相关的业务逻辑，计算合约交易量、合约持有人数等
	updateContractTxEventNum := ch.UpdateContractTxAndEventNum()

	//更新IDA合约数据
	updateIDAContractData := ch.DealIDAContractUpdateData(idaContractMap)

	identityContract := make([]*db.IdentityContract, 0)
	if ch.EventResults != nil {
		identityContract = ch.EventResults.IdentityContract
	}
	result := &db.GetContractResult{
		UpdateContractTxEventNum: updateContractTxEventNum,
		IdentityContract:         identityContract,
		UpdateIdaContract:        updateIDAContractData,
	}
	return result, nil
}

// DealIDAContractUpdateData 更新IDA合约数据
func (ch *ContractHandler) DealIDAContractUpdateData(
	idaContractMap map[string]*db.IDAContract) map[string]*db.IDAContract {
	idaEventData := ch.EventResults.IDAEventData
	minHeight := ch.MinHeight

	//更新数据
	idaContractUpdate := make(map[string]*db.IDAContract, 0)
	if idaEventData == nil {
		return idaContractUpdate
	}

	//新增数据资产
	idaCreateds := idaEventData.IDACreatedMap
	//删除数据资产
	idaDeleteIds := idaEventData.IDADeletedCodeMap
	for contractAddr, idaInfos := range idaCreateds {
		idaContract, ok := idaContractMap[contractAddr]
		if !ok || idaContract.BlockHeight >= minHeight {
			//如果数据库中BlockHeight，大于等于BlockHeight最小值，说明已经更新过了
			continue
		}

		//更新IDA合约数据资产数量
		totalNormalAssets := idaContract.TotalNormalAssets + int64(len(idaInfos))
		totalAssets := idaContract.TotalAssets + int64(len(idaInfos))
		idaContractUpdate[contractAddr] = &db.IDAContract{
			ContractAddr:      contractAddr,
			TotalNormalAssets: totalNormalAssets,
			TotalAssets:       totalAssets,
			BlockHeight:       minHeight,
		}
	}

	//删除数据资产
	for contractAddr, deleteIds := range idaDeleteIds {
		idaContract, ok := idaContractMap[contractAddr]
		if !ok || idaContract.BlockHeight >= minHeight {
			//如果数据库中BlockHeight，大于等于BlockHeight最小值，说明已经更新过了
			continue
		}

		//更新IDA合约数据资产数量
		if _, ok := idaContractUpdate[contractAddr]; ok {
			idaContractUpdate[contractAddr].TotalNormalAssets -= int64(len(deleteIds))
		} else {
			assetNum := idaContract.TotalNormalAssets - int64(len(deleteIds))
			idaContractUpdate[contractAddr] = &db.IDAContract{
				ContractAddr:      contractAddr,
				TotalNormalAssets: assetNum,
				BlockHeight:       minHeight,
			}
		}
	}

	return idaContractUpdate
}

// UpdateContractTxAndEventNum 更新合约交易数和事件数
func (ch *ContractHandler) UpdateContractTxAndEventNum() []*db.Contract {
	contractTxNumMap := make(map[string]int64, 0)
	contractEventNumMap := make(map[string]int64, 0)
	updateContractMap := make(map[string]*db.Contract, 0)
	updateContractNum := make([]*db.Contract, 0)

	contractMap := ch.ContractMap
	minHeight := ch.MinHeight
	txList := ch.TxList
	contractEvent := ch.ContractEvents
	//统计本次交易数据量
	for _, txInfo := range txList {
		if contract, ok := contractMap[txInfo.ContractNameBak]; ok {
			//说明已经更新过了
			if contract.BlockHeight >= minHeight {
				continue
			}
			contractTxNumMap[contract.Addr]++
		}
	}

	//统计本次合约事件数据量
	for _, event := range contractEvent {
		if contract, ok := contractMap[event.ContractNameBak]; ok {
			//说明已经更新过了
			if contract.BlockHeight >= minHeight {
				continue
			}
			contractEventNumMap[contract.Addr]++
		}
	}

	//更新合约交易和事件数量
	for addr, txNum := range contractTxNumMap {
		var contractInfo *db.Contract
		if contract, ok := updateContractMap[addr]; ok {
			contractInfo = contract
		} else if contract, ok = contractMap[addr]; ok {
			contractInfo = contract
		} else {
			continue
		}
		contractInfo.TxNum = contractInfo.TxNum + txNum
		contractInfo.BlockHeight = minHeight
		updateContractMap[addr] = contractInfo
	}

	//统计本次合约事件数据量
	for addr, eventNum := range contractEventNumMap {
		var contractInfo *db.Contract
		if contract, ok := updateContractMap[addr]; ok {
			contractInfo = contract
		} else if contract, ok = contractMap[addr]; ok {
			contractInfo = contract
		} else {
			continue
		}
		contractInfo.EventNum = contractInfo.EventNum + eventNum
		contractInfo.BlockHeight = minHeight
		updateContractMap[addr] = contractInfo
	}

	//更新合约交易和事件数量
	for _, contract := range updateContractMap {
		updateContractNum = append(updateContractNum, contract)
	}

	return updateContractNum
}
