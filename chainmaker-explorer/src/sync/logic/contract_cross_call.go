/*
Package logic comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package logic

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/utils"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

func DealCrossContractCallTxs(chainId string, txInfo *db.Transaction) (
	map[string]*db.ContractCrossCall, []*db.ContractCrossCallTransaction) {
	crossContractCallTxs := make([]*db.ContractCrossCallTransaction, 0)
	crossContractCalls := make(map[string]*db.ContractCrossCall, 0)

	existingIDs := make(map[string]struct{}) // 用于追踪已存在的ID
	// 处理读集合
	processRWSets(txInfo, txInfo.ReadSet, existingIDs, &crossContractCalls, &crossContractCallTxs)
	// 处理写集合
	processRWSets(txInfo, txInfo.WriteSet, existingIDs, &crossContractCalls, &crossContractCallTxs)
	return crossContractCalls, crossContractCallTxs
}

func processRWSets(tx *db.Transaction, rwSets string, existingIDs map[string]struct{},
	crossContractCalls *map[string]*db.ContractCrossCall,
	crossContractCallTxs *[]*db.ContractCrossCallTransaction) {
	if rwSets == "" {
		return
	}

	sets := make([]config.RwSet, 0)
	if err := json.Unmarshal([]byte(rwSets), &sets); err != nil {
		log.Errorf("processRWSets Unmarshal failed: %v", err)
		return
	}

	for _, set := range sets {
		targetContract := set.ContractName
		// 检查合约名称有效性
		if targetContract == "" {
			continue
		}

		// 生成唯一ID
		hashStr := generateCallHash(tx, targetContract)
		// 去重检查
		if _, exists := existingIDs[hashStr]; !exists {
			*crossContractCallTxs = append(*crossContractCallTxs, &db.ContractCrossCallTransaction{
				Id:               hashStr,
				TxId:             tx.TxId,
				BlockHeight:      tx.BlockHeight,
				InvokingContract: tx.ContractNameBak,
				InvokingMethod:   tx.ContractMethod,
				TargetContract:   targetContract,
				UserAddr:         tx.UserAddr,
				IsCross:          isCrossContract(targetContract, tx),
				Timestamp:        tx.Timestamp,
			})
			existingIDs[hashStr] = struct{}{}
		}

		if tx.ContractNameBak != targetContract && tx.ContractAddr != targetContract {
			newUUID := uuid.New().String()
			callsKey := fmt.Sprintf("%s_%s_%s", tx.ContractNameBak, tx.ContractMethod, targetContract)
			(*crossContractCalls)[callsKey] = &db.ContractCrossCall{
				Id:               newUUID,
				InvokingContract: tx.ContractNameBak,
				InvokingMethod:   tx.ContractMethod,
				TargetContract:   targetContract,
			}
		}
	}
}

// 辅助函数,判断是否是跨合约调用
func isCrossContract(target string, tx *db.Transaction) bool {
	return target != tx.ContractNameBak && target != tx.ContractAddr
}

// 生成唯一ID,	用于去重
func generateCallHash(tx *db.Transaction, target string) string {
	payload := fmt.Sprintf("%s_%s_%s_%s",
		tx.TxId,
		tx.ContractNameBak,
		tx.ContractMethod,
		target,
	)

	return utils.CalculateSHA256([]byte(payload))
}
