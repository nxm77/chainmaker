/*
Package logic comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package logic

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/sync/common"
	"chainmaker_web/src/sync/model"
)

// AccountHandler 处理合约相关的所有业务逻辑
type AccountHandler struct {
	ChainId          string
	MinHeight        int64
	TxList           map[string]*db.Transaction
	TransferEvents   []*db.ContractEventData
	EventResults     *model.TopicEventResult
	AccountMap       map[string]*db.Account
	DelayGetDBResult *model.GetDBResult
}

// NewAccountHandler 构造函数，初始化AccountHandler
func NewAccountHandler(chainId string, minHeight int64, eventResults *model.TopicEventResult,
	delayedUpdateCache *model.GetRealtimeCacheData, delayGetDBResult *model.GetDBResult,
	transferEvents []*db.ContractEventData) *AccountHandler {
	return &AccountHandler{
		ChainId:          chainId,
		MinHeight:        minHeight,
		TxList:           delayedUpdateCache.TxList,
		TransferEvents:   transferEvents,
		EventResults:     eventResults,
		AccountMap:       delayGetDBResult.AccountDBMap,
		DelayGetDBResult: delayGetDBResult,
	}
}

// DealWithContractData 处理合约相关数据
// 处理账户数据
func (ch *AccountHandler) DealWithAccountData() *db.UpdateAccountResult {
	// 计算本次账户交易和NFT数量
	accountTxNum, accountNFTNum := ch.DealAccountTxNFTNum()

	// 计算需要新增，更新的账户信息
	insertAccountMap, updateAccountMap := ch.BuildAccountInsertOrUpdate(accountTxNum, accountNFTNum)
	// 返回结果
	result := &db.UpdateAccountResult{
		InsertAccount: insertAccountMap,
		UpdateAccount: updateAccountMap,
	}
	return result
}

func (ch *AccountHandler) BuildAccountInsertOrUpdate(accountTx, accountNFT map[string]int64) (
	[]*db.Account, []*db.Account) {
	//定义两个map，分别存储新增和更新的账户
	var (
		accountInsertMap = make(map[string]*db.Account)
		accountUpdateMap = make(map[string]*db.Account)
	)

	//获取账户map和延迟获取数据库结果
	accountMap := ch.AccountMap
	delayGetDBResult := ch.DelayGetDBResult
	topicEventResult := ch.EventResults

	//数据库不存在的账户即为新增数据
	for _, addr := range topicEventResult.OwnerAdders {
		if addr == "" {
			continue
		}
		//数据库不存在的账户地址
		if _, ok := accountMap[addr]; !ok {
			//获取账户类型
			addrType := GetAccountType(ch.ChainId, addr)
			accountInfo := &db.Account{
				AddrType:    addrType,
				Address:     addr,
				BlockHeight: ch.MinHeight,
			}
			accountInsertMap[addr] = accountInfo
		}
	}

	//处理BNS账户
	bnsBind := topicEventResult.BNSBindEventData
	bnsUnBindAccount := delayGetDBResult.AccountBNSList
	processBNSAccounts(bnsBind, bnsUnBindAccount, accountInsertMap, accountUpdateMap, accountMap)

	//处理DID账户
	if topicEventResult.DIDEventData != nil {
		didBindAccounts := topicEventResult.DIDEventData.Account
		didUnBindAccounts := delayGetDBResult.AccountDIDList
		processDIDAccounts(didBindAccounts, didUnBindAccounts, accountInsertMap, accountUpdateMap, accountMap)
	}

	//统计账户交易数量
	dealAccountTxNum(ch.ChainId, ch.MinHeight, accountTx, accountInsertMap, accountUpdateMap, accountMap)
	//统计账户NFT数量
	dealAccountNFTNum(ch.ChainId, ch.MinHeight, accountNFT, accountInsertMap, accountUpdateMap, accountMap)

	var (
		accountInsert []*db.Account
		accountUpdate []*db.Account
	)
	if len(accountInsertMap) > 0 {
		for _, account := range accountInsertMap {
			accountInsert = append(accountInsert, account)
			accountMap[account.Address] = account
		}
	}
	if len(accountUpdateMap) > 0 {
		for _, account := range accountUpdateMap {
			//已经更新过了
			if account.BlockHeight >= ch.MinHeight {
				continue
			}
			account.BlockHeight = ch.MinHeight
			accountUpdate = append(accountUpdate, account)
			accountMap[account.Address] = account
		}
	}

	ch.AccountMap = accountMap
	return accountInsert, accountUpdate
}

func (ch *AccountHandler) DealAccountTxNFTNum() (map[string]int64, map[string]int64) {
	accountTxNum := make(map[string]int64)
	accountNFTNum := make(map[string]int64)
	// 更新交易数量
	for _, tx := range ch.TxList {
		accountTxNum[tx.UserAddr]++
	}

	for _, event := range ch.TransferEvents {
		// 只解析交易流转的topic
		if _, ok := common.TopicEventDataKey[event.Topic]; !ok {
			continue
		}

		if event.EventData == nil || event.EventData.TokenId == "" {
			continue
		}

		fromAddr := event.EventData.FromAddress
		toAddr := event.EventData.ToAddress
		if fromAddr != "" {
			accountNFTNum[fromAddr]--
		}
		if toAddr != "" {
			accountNFTNum[toAddr]++
		}
	}

	return accountTxNum, accountNFTNum
}
