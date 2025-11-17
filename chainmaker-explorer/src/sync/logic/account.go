/*
Package logic comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package logic

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	loggers "chainmaker_web/src/logger"
	"chainmaker_web/src/sync/common"
)

var (
	log = loggers.GetLogger(loggers.MODULE_SYNC)
)

// GetAccountType 获取账户类型
func GetAccountType(chainId, address string) int {
	//判断是否是合约地址
	//contractInfo, err := dbhandle.GetContractByCacheOrAddr(chainId, address)
	contractInfo, err := dbhandle.GetContractByAddr(chainId, address)
	addrType := common.AddrTypeUser
	if contractInfo != nil && err == nil {
		addrType = common.AddrTypeContract
	}
	return addrType
}

/**
 * @description: 统计账户NFT数量
 * @param {string} chainId 链id
 * @param {int64} minHeight  本次订阅区块高度
 * @param {*} accountInsertMap 新增账户
 * @param {*} accountUpdateMap 更新账户
 * @param {map[string]*db.Account} accountMap 本次订阅涉及账户
 * @param {map[string]int64} accountNFT 账户NFT数量
 */
func dealAccountNFTNum(chainId string, minHeight int64, accountNFT map[string]int64,
	accountInsertMap, accountUpdateMap, accountMap map[string]*db.Account) {
	for address, num := range accountNFT {
		if address == "" {
			continue
		}

		if account, ok := accountInsertMap[address]; ok {
			nftNum := account.NFTNum + num
			if nftNum < 0 {
				nftNum = 0
			}
			account.NFTNum = nftNum
		} else if accountUp, okUp := accountUpdateMap[address]; okUp {
			nftNum := accountUp.NFTNum + num
			if nftNum < 0 {
				nftNum = 0
			}
			accountUp.NFTNum = nftNum
		} else if accountDB, okDB := accountMap[address]; okDB {
			nftNum := accountDB.NFTNum + num
			if nftNum < 0 {
				nftNum = 0
			}
			accountDB.NFTNum = nftNum
			accountUpdateMap[address] = accountDB
		} else {
			var nftNum int64
			if num > 0 {
				nftNum = num
			}
			//获取账户类型
			addrType := GetAccountType(chainId, address)
			accountInfo := &db.Account{
				AddrType:    addrType,
				Address:     address,
				NFTNum:      nftNum,
				BlockHeight: minHeight,
			}
			accountInsertMap[address] = accountInfo
		}
	}
}

/**
 * @description: 统计账户交易数量
 * @param {string} chainId
 * @param {int64} minHeight
 * @param {*} accountInsertMap  新增账户
 * @param {*} accountUpdateMap  更新账户
 * @param {map[string]*db.Account} accountMap 本次订阅涉及账户
 * @param {map[string]int64} accountTx 本次订阅账户交易量
 */
func dealAccountTxNum(chainId string, minHeight int64, accountTx map[string]int64,
	accountInsertMap, accountUpdateMap, accountMap map[string]*db.Account) {
	for address, num := range accountTx {
		if address == "" || num == 0 {
			continue
		}

		if account, ok := accountInsertMap[address]; ok {
			txNum := account.TxNum + num
			if txNum < 0 {
				txNum = 0
			}
			account.TxNum = txNum
		} else if accountUp, okUp := accountUpdateMap[address]; okUp {
			txNum := accountUp.TxNum + num
			if txNum < 0 {
				txNum = 0
			}
			accountUp.TxNum = txNum
		} else if accountDB, okDB := accountMap[address]; okDB {
			txNum := accountDB.TxNum + num
			if txNum < 0 {
				txNum = 0
			}
			accountDB.TxNum = txNum
			accountUpdateMap[address] = accountDB
		} else {
			var txNum int64
			if num > 0 {
				txNum = num
			}
			//获取账户类型
			addrType := GetAccountType(chainId, address)
			accountInfo := &db.Account{
				AddrType:    addrType,
				Address:     address,
				TxNum:       txNum,
				BlockHeight: minHeight,
			}
			accountInsertMap[address] = accountInfo
		}
	}
}

// processBNSAccounts
//
//	@Description: 更新账户BNS
//	@param chainId
//	@param bnsBindEventData 绑定bns
//	@param accountInsertMap
//	@param accountUpdateMap
//	@param accountMap
func processBNSAccounts(bnsBindEventData []*db.BNSTopicEventData, unBindBNSs []*db.Account, accountInsertMap,
	accountUpdateMap, accountMap map[string]*db.Account) {
	//绑定BNS
	for _, event := range bnsBindEventData {
		accountAddr := event.Value
		accountBNS := event.Domain
		if accountDB, ok := accountMap[accountAddr]; ok {
			// 数据库存在，更新BNS
			if updateAcc, okUp := accountUpdateMap[accountAddr]; okUp {
				updateAcc.BNS = accountBNS
				accountUpdateMap[accountAddr] = updateAcc
			} else {
				accountDB.BNS = accountBNS
				accountUpdateMap[accountAddr] = accountDB
			}
		} else {
			// 数据库不存在，判断是否已经在insert了
			if account, okIn := accountInsertMap[accountAddr]; okIn {
				account.BNS = accountBNS
			}
		}
	}

	//解绑BNS
	for _, account := range unBindBNSs {
		accountAddr := account.Address
		if accountDB, ok := accountMap[accountAddr]; ok {
			// 数据库存在，更新BNS
			if updateAcc, okUp := accountUpdateMap[accountAddr]; okUp {
				updateAcc.BNS = ""
				accountUpdateMap[accountAddr] = updateAcc
			} else {
				accountDB.BNS = ""
				accountUpdateMap[accountAddr] = accountDB
			}
		}
	}

}

// 处理DID账户
func processDIDAccounts(bindAccounts map[string]*db.Account, unBindAccounts []*db.Account,
	accountInsertMap, accountUpdateMap, accountMap map[string]*db.Account) {
	// 将 unBindDIDs 转换为一个包含地址的集合，用于解绑DID
	unbindDIDSet := make(map[string]bool)
	for _, account := range unBindAccounts {
		unbindDIDSet[account.Address] = true
	}

	for address, account := range bindAccounts {
		if accountDB, ok := accountMap[address]; ok {
			// 数据库存在，更新DID
			if updateAcc, okUp := accountUpdateMap[address]; okUp {
				updateAcc.DID = account.DID
				//updateAcc.AddrID = account.AddrID
			} else {
				accountDB.DID = account.DID
				//accountDB.AddrID = account.AddrID
				accountUpdateMap[address] = accountDB
			}

			// 从 unbindDIDSet 中移除已处理的地址
			delete(unbindDIDSet, address)
		} else {
			// 数据库不存在的账户地址
			if accountInsert, okIn := accountInsertMap[address]; okIn {
				accountInsert.DID = account.DID
				//accountInsert.AddrID = account.AddrID
			}
		}
	}

	// 处理 unbindDIDSet 中剩余的地址，这些数据将要解绑DID
	for accountAddr := range unbindDIDSet {
		if accountDB, ok := accountMap[accountAddr]; ok {
			// 数据库存在，更新DID
			if updateAcc, okUp := accountUpdateMap[accountAddr]; okUp {
				updateAcc.DID = ""
				//updateAcc.AddrID = ""
				accountUpdateMap[accountAddr] = updateAcc
			} else {
				accountDB.DID = ""
				//accountDB.AddrID = ""
				accountUpdateMap[accountAddr] = accountDB
			}
		}
	}
}
