/*
Package service comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package service

import (
	"chainmaker_web/src/auth"
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/entity"
	"chainmaker_web/src/logic"
	"chainmaker_web/src/utils"

	"github.com/gin-gonic/gin"
)

// UploadContractAbiHandler get
type UploadContractAbiHandler struct {
}

// RequiresAuth 是否需要登录验证
func (p *UploadContractAbiHandler) RequiresAuth() bool {
	return true
}

// Handle UploadContractAbiHandler
func (handler *UploadContractAbiHandler) Handle(ctx *gin.Context) {
	params := entity.BindUploadContractAbiHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	// 根据链ID和合约地址查询合约
	upgradeContract, err := dbhandle.GetUpgradeContractInfo(params.ChainId, params.ContractAddr, params.ContractVersion)
	if err != nil || upgradeContract == nil {
		newError := entity.GetErrorMsg(entity.ErrorContractNotExist)
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//登录验证
	if !ABICheckLogin(ctx, upgradeContract) {
		ConvergeHandleFailureResponse(ctx, entity.GetErrorNoPermission(""))
		return
	}

	// 1. 解析并验证基础数据
	err = logic.SaveContractABI(params, upgradeContract)
	if err != nil {
		log.Errorf("SaveContractABI err : %s", err)
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	//异步处理历史数据
	go func() {
		if err := logic.AsyncHandleContractABIEvent(
			params.ChainId,
			params.ContractAddr,
			params.ContractVersion,
		); err != nil {
			// 添加错误处理逻辑
			log.Errorf("处理合约ABI事件失败: %v", err)
		}
	}()

	// 2. 返回成功
	ConvergeDataResponse(ctx, "OK", nil)
}

// ABICheckLogin 函数用于验证用户是否登录
func ABICheckLogin(ctx *gin.Context, upgradeContract *db.UpgradeContractTransaction) bool {
	//获取用户地址和token
	userAddr, _, exists := auth.GetUserAddrAndToken(ctx)
	//如果用户未登录
	if !exists {
		//登录验证,输出未登录日志
		log.Errorf("ABICheckLogin user not login")
		return false
	}
	adminAddr := utils.GetAccountHashStr()
	if adminAddr == userAddr || upgradeContract.UserAddr == userAddr {
		return true
	}

	return false
}

// GetContractTopicsHandler get
type GetContractTopicsHandler struct {
}

// Handle GetContractTopicsHandler
func (handler *GetContractTopicsHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetContractTopicsHandler(ctx)
	if params == nil || !params.IsLegal() {
		ConvergeFailureResponse(ctx, entity.GetErrorMsgParams())
		return
	}

	// 根据链ID和合约地址查询合约
	contract, err := dbhandle.GetContractByAddr(params.ChainId, params.ContractAddr)
	if err != nil || contract == nil {
		newError := entity.GetErrorMsg(entity.ErrorContractNotExist)
		ConvergeHandleFailureResponse(ctx, newError)
		return
	}

	//获取合约abi详情
	topics, err := logic.GetContractTopics(params.ChainId, params.ContractAddr, params.ContractVersion)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	abiTopicView := &entity.ABITopicView{
		ContractAddr:    contract.Addr,
		ContractVersion: params.ContractVersion,
		ContractName:    contract.NameBak,
		Topics:          topics,
	}

	ConvergeDataResponse(ctx, abiTopicView, nil)
}

// GetContractABIDataHandler get
type GetContractABIDataHandler struct{}

// Handle deal
func (handler *GetContractABIDataHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetABIDataHandler(ctx)
	if params == nil || !params.IsLegal() {
		ConvergeFailureResponse(ctx, entity.GetErrorMsgParams())
		return
	}

	// 根据链ID和合约地址查询合约
	contract, err := dbhandle.GetContractByAddr(params.ChainId, params.ContractAddr)
	if err != nil || contract == nil {
		newError := entity.GetErrorMsg(entity.ErrorContractNotExist)
		ConvergeHandleFailureResponse(ctx, newError)
		return
	}

	//获取合约abi详情
	abiDetail, err := dbhandle.GetContractABIFile(params.ChainId, params.ContractAddr, params.ContractVersion)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	abiDetailView := &entity.ABIDetailView{
		ContractAddr:    contract.Addr,
		ContractVersion: params.ContractVersion,
		ContractName:    contract.NameBak,
	}
	if abiDetail != nil {
		abiDetailView.ABIJson = abiDetail.ABIJson
	}
	ConvergeDataResponse(ctx, abiDetailView, nil)
}

// GetDecodeContractEventsHandler get
type GetDecodeContractEventsHandler struct {
}

// Handle GetDecodeContractEventsHandler
func (handler *GetDecodeContractEventsHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetDecodeContractEventsHandler(ctx)
	if params == nil || !params.IsLegal() {
		ConvergeFailureResponse(ctx, entity.GetErrorMsgParams())
		return
	}

	topicIndex, err := logic.GetContractTopicIndexs(params.ChainId, params.ContractAddr,
		params.ContractVersion, params.Topic)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	topicColumns, err := logic.GetContractTopicColumns(params.ChainId, params.ContractAddr,
		params.ContractVersion, params.Topic)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	results, total, err := logic.GetDecodeContractEvents(params, topicColumns, topicIndex)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	abiTopicView := &entity.ABIEventDataView{
		TotalCount: total,
		Indexs:     topicIndex,
		Columns:    topicColumns,
		Result:     results,
	}

	ConvergeDataResponse(ctx, abiTopicView, nil)
}
