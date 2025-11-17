/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"chainmaker_web/src/db"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// InsertAccount 插入账户
func InsertAccount(chainId string, accountList []*db.Account) error {
	if len(accountList) == 0 {
		return nil
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableAccount)
	err := CreateInBatchesData(tableName, accountList)
	if err != nil {
		return err
	}

	//更新账户缓存
	DelAccountTotalCache(chainId)
	return nil
}

// GetAccountList 获取账户列表
func GetAccountList(offset int, limit int, chainId, did, address string) ([]*db.Account, int64, error) {
	accountList := make([]*db.Account, 0)
	var count int64
	where := map[string]interface{}{}
	if did != "" {
		where["did"] = did
	}
	if address != "" {
		where["address"] = address
	}
	tableName := db.GetTableName(chainId, db.TableAccount)
	query := db.GormDB.Table(tableName).Where(where)
	err := query.Count(&count).Error
	if err != nil {
		return accountList, 0, err
	}
	query = query.Offset(offset * limit).Limit(limit)
	err = query.Find(&accountList).Error
	if err != nil {
		return accountList, 0, err
	}
	return accountList, count, nil
}

// UpdateAccount 更新账户
func UpdateAccount(chainId string, accountInfo *db.Account) error {
	if chainId == "" || accountInfo == nil {
		return db.ErrTableParams
	}

	//获取交易表名称
	tableName := db.GetTableName(chainId, db.TableAccount)
	where := map[string]interface{}{
		"address": accountInfo.Address,
	}
	params := map[string]interface{}{
		"did": accountInfo.DID,
		"bns": accountInfo.BNS,
	}

	if accountInfo.TxNum > 0 {
		params["txNum"] = accountInfo.TxNum
	}

	if accountInfo.NFTNum > 0 {
		params["nftNum"] = accountInfo.NFTNum
	}

	if accountInfo.BlockHeight > 0 {
		params["blockHeight"] = accountInfo.BlockHeight
	}

	err := db.GormDB.Table(tableName).Where(where).Updates(params).Error
	if err != nil {
		return err
	}

	//更新账户缓存
	UpdateAccountDataCache(chainId, accountInfo)
	return nil
}

// QueryAccountExists 根据addr查询
func QueryAccountExists(chainId string, addrList []string) (map[string]*db.Account, error) {
	var (
		accountList    = make([]*db.Account, 0)
		accountMap     = make(map[string]*db.Account, 0)
		selectAccounts = make([]string, 0)
	)
	for _, addr := range addrList {
		//获取账户缓存
		accountData, err := GetAccountCacheByAddr(chainId, addr)
		if err == nil && accountData != nil {
			accountMap[accountData.Address] = accountData
		} else {
			selectAccounts = append(selectAccounts, addr)
		}
	}

	if len(selectAccounts) == 0 {
		return accountMap, nil
	}

	tableName := db.GetTableName(chainId, db.TableAccount)
	err := db.GormDB.Table(tableName).Where("address in ?", selectAccounts).Find(&accountList).Error
	if err != nil {
		return nil, err
	}

	for _, v := range accountList {
		accountMap[v.Address] = v
		//设置账户缓存
		SetAccountDataCache(chainId, v)
	}
	return accountMap, nil
}

// GetAccountByAddr 根据账户地址获取账户信息
func GetAccountByAddr(chainId, address string) (*db.Account, error) {
	if address == "" || chainId == "" {
		return nil, db.ErrTableParams
	}

	account := &db.Account{}
	where := map[string]interface{}{
		"address": address,
	}
	tableName := db.GetTableName(chainId, db.TableAccount)
	err := db.GormDB.Table(tableName).Where(where).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	//设置账户缓存
	SetAccountDataCache(chainId, account)
	return account, nil
}

// GetAccountByDID
//
//	@Description: 根据DID 获取账户信息,一个DID对应多个地址
//	@param chainId
//	@param did
//	@return []*db.Account
//	@return error
func GetAccountByDID(chainId, did string) ([]*db.Account, error) {
	if did == "" || chainId == "" {
		return nil, db.ErrTableParams
	}

	accountList := make([]*db.Account, 0)
	where := map[string]interface{}{
		"did": did,
	}
	tableName := db.GetTableName(chainId, db.TableAccount)
	err := db.GormDB.Table(tableName).Where(where).Find(&accountList).Error
	if err != nil {
		return accountList, err
	}

	return accountList, nil
}

// GetAccountByBNS
//
//	@Description: 根据BNS获取账户地址，一个BNS对应一个地址
//	@param chainId
//	@param bns
//	@return *db.Account
//	@return error
func GetAccountByBNS(chainId, bns string) (*db.Account, error) {
	if bns == "" || chainId == "" {
		return nil, db.ErrTableParams
	}

	account := &db.Account{}
	where := map[string]interface{}{
		"bns": bns,
	}
	tableName := db.GetTableName(chainId, db.TableAccount)
	err := db.GormDB.Table(tableName).Where(where).First(&account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return account, nil
}

// GetAccountByBNSList 根据账户地址获取账户信息
func GetAccountByBNSList(chainId string, bnsList []string) ([]*db.Account, error) {
	accountList := make([]*db.Account, 0)
	if len(bnsList) == 0 {
		return accountList, nil
	}

	tableName := db.GetTableName(chainId, db.TableAccount)
	err := db.GormDB.Table(tableName).Where("bns in ?", bnsList).Find(&accountList).Error
	if err != nil {
		return accountList, err
	}

	return accountList, nil
}

// GetAccountByDIDList 根据账户地址获取账户信息
func GetAccountByDIDList(chainId string, didList []string) ([]*db.Account, error) {
	accountList := make([]*db.Account, 0)
	if len(didList) == 0 {
		return accountList, nil
	}

	tableName := db.GetTableName(chainId, db.TableAccount)
	err := db.GormDB.Table(tableName).Where("did in ?", didList).Find(&accountList).Error
	if err != nil {
		return accountList, err
	}

	return accountList, nil
}

// GetAccountDetail 获取账户详情
func GetAccountDetail(chainId, address, bns string) (*db.Account, error) {
	var accountInfo *db.Account
	if chainId == "" && address == "" && bns == "" {
		return nil, db.ErrTableParams
	}
	where := map[string]interface{}{}
	if address != "" {
		where["address"] = address
	}
	if bns != "" {
		where["bns"] = bns
	}
	tableName := db.GetTableName(chainId, db.TableAccount)
	err := db.GormDB.Table(tableName).Where(where).First(&accountInfo).Error
	if err != nil {
		// if errors.Is(err, gorm.ErrRecordNotFound) {
		// 	return nil, nil
		// }
		return nil, err
	}
	return accountInfo, nil
}

// GetTotalTxNumByAccount 获取账户总交易数
func GetTotalTxNumByAccount(chainId string) (int64, error) {
	tableName := db.GetTableName(chainId, db.TableAccount)
	var totalTxNum int64 // 使用 interface{} 类型
	err := db.GormDB.Table(tableName).Select("sum(txNum)").Row().Scan(&totalTxNum)
	if err != nil {
		return 0, fmt.Errorf("select sum(txNum) err, cause : %v", err)
	}

	return totalTxNum, nil
}

// GetAccountTotal 获取账户总数
func GetAccountTotal(chainId string) (int64, error) {
	count := GetAccountTotalCache(chainId)
	if count != 0 {
		return count, nil
	}

	tableName := db.GetTableName(chainId, db.TableAccount)
	err := db.GormDB.Table(tableName).Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("count user err, cause : %s", err.Error())
	}

	//设置缓存
	SetAccountTotalCache(chainId, count)
	return count, nil
}
