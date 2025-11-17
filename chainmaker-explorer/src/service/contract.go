/*
Package service comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package service

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/db"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"chainmaker.org/chainmaker/contract-utils/standard"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gin-gonic/gin"

	"chainmaker_web/src/config"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/entity"
)

// GetLatestContractListHandler get
type GetLatestContractListHandler struct{}

// Handle deal
func (getLatestContractListHandler *GetLatestContractListHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetLatestContractHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetLatestContractList param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//从缓存获取最新的Contract
	contractList, err := getContractListFromRedis(ctx, params.ChainId)
	if err != nil {
		log.Errorf("GetLatestContractList get redis fail err:%v", err)
	}
	count := int64(len(contractList))
	if count == 0 {
		// 获取ContractList
		contractList, err = dbhandle.GetLatestContractList(params.ChainId)
		if err != nil {
			log.Errorf("GetLatestContractList err : %s", err.Error())
			ConvergeHandleFailureResponse(ctx, err)
			return
		}
	}

	//获取创建合约账户对应的账户信息
	accountMap := GetContractAccountMap(params.ChainId, contractList)

	//数据渲染
	contractViews := arraylist.New()
	for i, contract := range contractList {
		//获取地址BNS
		senderAddrBns := GetAccountBNS(contract.CreatorAddr, accountMap)
		latestChainView := &entity.LatestContractView{
			Id:               i + 1,
			ContractName:     contract.Name,
			ContractAddr:     contract.Addr,
			ContractType:     contract.ContractType,
			Sender:           contract.CreateSender,
			SenderAddr:       contract.CreatorAddr,
			SenderAddrBNS:    senderAddrBns,
			Version:          contract.Version,
			TxNum:            contract.TxNum,
			CreateTimestamp:  contract.Timestamp,
			UpgradeTimestamp: contract.UpgradeTimestamp,
			UpgradeUser:      contract.UpgradeAddr,
			Timestamp:        contract.Timestamp,
		}
		contractViews.Add(latestChainView)
	}
	ConvergeListResponse(ctx, contractViews.Values(), count, nil)
}

// getContractListFromRedis 获取缓存数据
func getContractListFromRedis(ctx *gin.Context, chainId string) ([]*db.Contract, error) {
	contractList := make([]*db.Contract, 0)
	//从缓存获取最新的合约
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisLatestContractList, prefix, chainId)
	redisList := cache.GlobalRedisDb.ZRevRange(ctx, redisKey, 0, 9).Val()
	for _, resStr := range redisList {
		contractInfo := &db.Contract{}
		err := json.Unmarshal([]byte(resStr), contractInfo)
		if err != nil {
			log.Errorf("getContractListFromRedis json Unmarshal err : %s", err.Error())
			return contractList, err
		}
		contractList = append(contractList, contractInfo)
	}

	return contractList, nil
}

// GetContractListHandler handler
type GetContractListHandler struct{}

// Handle deal
func (getContractListHandler *GetContractListHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetContractListHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetContractList param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	var (
		senders      []string
		senderAddrs  []string
		upgraders    []string
		upgradeAddrs []string
	)

	chainId := params.ChainId
	offset := params.Offset
	limit := params.Limit
	if params.Creators != "" {
		senders = strings.Split(params.Creators, ",")
	}
	if params.CreatorAddrs != "" {
		senderAddrs = strings.Split(params.CreatorAddrs, ",")
	}
	if params.Upgraders != "" {
		upgraders = strings.Split(params.Upgraders, ",")
	}
	if params.UpgradeAddrs != "" {
		upgradeAddrs = strings.Split(params.UpgradeAddrs, ",")
	}

	contractList, count, err := dbhandle.GetContractList(chainId, offset, limit, params.Status, params.RuntimeType,
		params.ContractKey, params.ContractType, params.Order, params.OrderBy,
		senders, senderAddrs, upgraders, upgradeAddrs)
	if err != nil {
		log.Errorf("GetContractList err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	//获取创建合约账户对应的账户信息
	accountMap := GetContractAccountMap(chainId, contractList)
	contractListView := arraylist.New()
	for i, contract := range contractList {
		//获取地址BNS
		senderAddrBns := GetAccountBNS(contract.CreatorAddr, accountMap)

		//获取合约验证状态
		var verifyStatus int
		upgradeContract, _ := dbhandle.GetUpgradeContractInfo(chainId, contract.Addr, contract.Version)
		if upgradeContract != nil {
			verifyStatus = upgradeContract.VerifyStatus
		}

		listId := params.Offset*params.Limit + i + 1
		contractView := &entity.ContractListView{
			Id:               strconv.Itoa(listId),
			ContractName:     contract.Name,
			ContractSymbol:   contract.ContractSymbol,
			ContractAddr:     contract.Addr,
			ContractType:     contract.ContractType,
			Version:          contract.Version,
			Creator:          contract.CreateSender,
			CreatorAddr:      contract.CreatorAddr,
			CreatorAddrBns:   senderAddrBns,
			Upgrader:         contract.Upgrader,
			UpgradeAddr:      contract.UpgradeAddr,
			UpgradeOrgId:     contract.UpgradeOrgId,
			TxNum:            contract.TxNum,
			Status:           contract.ContractStatus,
			CreateTimestamp:  contract.Timestamp,
			UpgradeTimestamp: contract.UpgradeTimestamp,
			RuntimeType:      contract.RuntimeType,
			VerifyStatus:     verifyStatus,
		}
		contractListView.Add(contractView)
	}

	ConvergeListResponse(ctx, contractListView.Values(), count, nil)
}

// GetContractDetailHandler handler
type GetContractDetailHandler struct{}

// Handle deal
func (getContractDetailHandler *GetContractDetailHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetContractDetailHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	chainId := params.ChainId
	//获取合约
	contract, err := dbhandle.GetContractByNameOrAddr(chainId, params.ContractKey)
	if err != nil {
		log.Errorf("GetContractDetail err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	//获取创建合约账户对应的账户信息
	accountMap := GetContractAccountMap(chainId, []*db.Contract{contract})
	//获取地址BNS
	senderAddrBns := GetAccountBNS(contract.CreatorAddr, accountMap)

	var dataAssetNum int64
	//判断是否是IDA合约
	if contract.ContractType == standard.ContractStandardNameCMIDA {
		idaContractInfo, err := dbhandle.GetIDAContractByAddr(chainId, contract.Addr)
		if err != nil {
			if err != nil {
				log.Errorf("GetIDAContractByAddr err : %v", err)
				ConvergeHandleFailureResponse(ctx, err)
				return
			}
		}
		if idaContractInfo != nil {
			dataAssetNum = idaContractInfo.TotalNormalAssets
		}
	}

	//获取合约验证状态
	var verifyStatus int
	upgradeContract, _ := dbhandle.GetUpgradeContractInfo(chainId, contract.Addr, contract.Version)
	if upgradeContract != nil {
		verifyStatus = upgradeContract.VerifyStatus
	}

	contractDetailView := &entity.ContractDetailView{
		ContractName:    contract.Name,
		ContractNameBak: contract.NameBak,
		ContractAddr:    contract.Addr,
		ContractSymbol:  contract.ContractSymbol,
		ContractType:    contract.ContractType,
		Version:         contract.Version,
		ContractStatus:  contract.ContractStatus,
		TxId:            contract.CreateTxId,
		CreateSender:    contract.CreateSender,
		CreatorAddr:     contract.CreatorAddr,
		CreatorAddrBns:  senderAddrBns,
		Timestamp:       contract.Timestamp,
		DataAssetNum:    dataAssetNum,
		RuntimeType:     contract.RuntimeType,
		VerifyStatus:    verifyStatus,
	}

	ConvergeDataResponse(ctx, contractDetailView, nil)
}

// GetEventListHandler handler
type GetEventListHandler struct {
}

// Handle deal
func (getEventListHandler *GetEventListHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetEventListHandler(ctx)
	if params == nil {
		ConvergeFailureResponse(ctx, entity.NewError(entity.ErrorParamWrong, "参数绑定失败"))
		return
	}
	if err := params.Validate(); err != nil {
		// 直接返回结构化错误
		ConvergeFailureResponse(ctx, entity.NewError(entity.ErrorParamWrong, err.Error()))
		return
	}

	var contractInfo *db.Contract
	var err error
	//获取合约详情
	if params.ContractAddr != "" {
		contractInfo, err = dbhandle.GetContractByAddr(params.ChainId, params.ContractAddr)
	} else {
		contractInfo, err = dbhandle.GetContractByName(params.ChainId, params.ContractName)
	}
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
	}

	if contractInfo == nil || contractInfo.EventNum == 0 {
		ConvergeListResponse(ctx, []interface{}{}, 0, nil)
		return
	}

	var topics []string
	var totalCount int64
	if params.Topic != "" {
		topicList, errDB := dbhandle.GetContractEventTopic(params.ChainId, contractInfo.NameBak, params.Topic)
		if errDB != nil {
			log.Errorf("GetEventByTopic err : %s", errDB)
			ConvergeHandleFailureResponse(ctx, errDB)
			return
		}
		for _, eventTopic := range topicList {
			topics = append(topics, eventTopic.Topic)
		}

		if len(topics) == 0 {
			ConvergeListResponse(ctx, []interface{}{}, 0, nil)
			return
		}

		for _, eventTopic := range topicList {
			totalCount += eventTopic.TxNum
		}
	} else {
		totalCount = contractInfo.EventNum
	}

	if params.TxId != "" {
		totalCount, err = dbhandle.GetEventListCount(params.ChainId, contractInfo.NameBak, params.TxId, topics)
		if err != nil {
			log.Errorf("GetEventListCount err : %s", err.Error())
			ConvergeListResponse(ctx, []interface{}{}, 0, nil)
			return
		}
	}

	eventList, err := dbhandle.GetContractEventList(params.Offset, params.Limit, params.ChainId, contractInfo.NameBak,
		params.TxId, topics)
	if err != nil {
		log.Errorf("GetEventList err : %s", err.Error())
		ConvergeListResponse(ctx, []interface{}{}, 0, nil)
		return
	}

	eventListView := arraylist.New()
	for _, event := range eventList {
		contractEventView := &entity.ContractEventView{
			Topic:           event.Topic,
			EventInfo:       event.EventData,
			TxId:            event.TxId,
			Timestamp:       event.Timestamp,
			ContractVersion: event.ContractVersion,
		}
		eventListView.Add(contractEventView)
	}

	ConvergeListResponse(ctx, eventListView.Values(), totalCount, nil)
}

type GetContractTypesHandler struct {
}

// Handle deal
func (handler *GetContractTypesHandler) Handle(ctx *gin.Context) {
	params := entity.BindContractTypesHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	contractTypes, err := dbhandle.GetContractTypes(params.ChainId)
	if err != nil {
		log.Errorf("GetContractTypes err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}
	//组装response数据
	view := &entity.ContractTypesView{
		ContractTypes: contractTypes,
	}
	ConvergeDataResponse(ctx, view, nil)
}
