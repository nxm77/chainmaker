/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package logic

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/sync/common"

	"github.com/google/uuid"
)

func DealTransferList(eventDataList []*db.ContractEventData, contractInfoMap map[string]*db.Contract,
	txInfoMap map[string]*db.Transaction) (
	[]*db.FungibleTransfer, []*db.NonFungibleTransfer) {
	dbFungibleTransfer := make([]*db.FungibleTransfer, 0)
	dbNonFungibleTransfer := make([]*db.NonFungibleTransfer, 0)
	for _, event := range eventDataList {
		//只解析交易流转的topic
		if _, ok := common.TopicEventDataKey[event.Topic]; !ok {
			continue
		}

		//合约数据
		contract, ok := contractInfoMap[event.ContractName]
		if !ok || event.EventData == nil {
			continue
		}

		var contractMethod string
		//同质化合约
		if contract.ContractType == common.ContractStandardNameCMDFA ||
			contract.ContractType == common.ContractStandardNameEVMDFA {
			if value, okTx := txInfoMap[event.TxId]; okTx {
				contractMethod = value.ContractMethod
			}
			// 将字符串转换为 Decimal 值
			amountDecimal := common.StringAmountDecimal(event.EventData.Amount, contract.Decimals)
			newUUID := uuid.New().String()
			dbFungibleTransfer = append(dbFungibleTransfer, &db.FungibleTransfer{
				ID:             newUUID,
				TxId:           event.TxId,
				EventIndex:     event.Index,
				ContractName:   contract.Name,
				ContractAddr:   contract.Addr,
				ContractMethod: contractMethod,
				Topic:          event.Topic,
				FromAddr:       event.EventData.FromAddress,
				ToAddr:         event.EventData.ToAddress,
				Amount:         amountDecimal,
				Timestamp:      event.Timestamp,
			})
		} else if contract.ContractType == common.ContractStandardNameCMNFA ||
			contract.ContractType == common.ContractStandardNameEVMNFA {
			if value, okTx := txInfoMap[event.TxId]; okTx {
				contractMethod = value.ContractMethod
			}
			//非同质化合约
			newUUID := uuid.New().String()
			dbNonFungibleTransfer = append(dbNonFungibleTransfer, &db.NonFungibleTransfer{
				ID:             newUUID,
				TxId:           event.TxId,
				EventIndex:     event.Index,
				ContractName:   contract.Name,
				ContractAddr:   contract.Addr,
				ContractMethod: contractMethod,
				Topic:          event.Topic,
				FromAddr:       event.EventData.FromAddress,
				ToAddr:         event.EventData.ToAddress,
				TokenId:        event.EventData.TokenId,
				Timestamp:      event.Timestamp,
			})
		}
	}
	return dbFungibleTransfer, dbNonFungibleTransfer
}
