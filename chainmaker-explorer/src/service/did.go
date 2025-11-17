/*
Package service comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package service

import (
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/entity"
	"encoding/json"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gin-gonic/gin"
)

const (
	//OperationTypeCreate did创建
	OperationTypeCreate = iota
	//OperationTypeUpdate did更新
	OperationTypeUpdate
)

// GetDIDDetailHandler get
type GetDIDDetailHandler struct {
}

// Handle deal
func (handler *GetDIDDetailHandler) Handle(ctx *gin.Context) {
	// 参数处理
	params := entity.BindGetDIDDetailHandler(ctx)
	if params == nil || !params.IsLegal() {
		// 参数错误，返回错误信息
		newError := entity.NewError(entity.ErrorParamWrong, entity.ErrorMsgParam)
		ConvergeFailureResponse(ctx, newError)
		return
	}

	// 根据参数获取DID详情
	didInfo, err := dbhandle.GetDIDDetailById(params.ChainId, params.DID)
	if err != nil {
		// 获取DID详情失败，返回错误信息
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	// 初始化IssuerService
	issuerService := make([]entity.IssuerService, 0)
	authentication := make([]string, 0)
	accountList := make([]entity.AccountData, 0)
	if didInfo.IssuerService != "" {
		// 解析IssuerService
		_ = json.Unmarshal([]byte(didInfo.IssuerService), &issuerService)
	}
	if didInfo.Authentication != "" {
		// 解析Authentication
		_ = json.Unmarshal([]byte(didInfo.Authentication), &authentication)
	}
	if didInfo.AccountJson != "" {
		// 解析AccountJson
		_ = json.Unmarshal([]byte(didInfo.AccountJson), &accountList)
	}

	// 构造DIDDetailView
	didDetailView := &entity.DIDDetailView{
		DID:            didInfo.DID,
		Document:       didInfo.Document,
		ContractName:   didInfo.ContractName,
		ContractAddr:   didInfo.ContractAddr,
		Status:         didInfo.Status,
		IsIssuer:       didInfo.IsIssuer,
		Timestamp:      didInfo.Timestamp,
		TxId:           didInfo.CreateTxId,
		IssuerService:  issuerService,
		Authentication: authentication,
		AccountList:    accountList,
	}
	// 返回DIDDetailView
	ConvergeDataResponse(ctx, didDetailView, nil)
}

// GetDIDListHandler get
type GetDIDListHandler struct {
}

// Handle GetBlockList区块链列表
// 处理获取DID列表的请求
func (handler *GetDIDListHandler) Handle(ctx *gin.Context) {
	// 绑定请求参数
	params := entity.BindGetDIDListHandler(ctx)
	// 如果参数为空或者不合法，则返回错误响应
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, entity.ErrorMsgParam)
		ConvergeFailureResponse(ctx, newError)
		return
	}

	// 从数据库中获取DID列表和总数
	didList, totalCount, err := dbhandle.GetDIDListAndCount(params.Offset, params.Limit,
		params.ChainId, params.ContractAddr, params.DID)
	// 如果获取失败，则返回错误响应
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	// 创建一个空的DID列表视图
	listView := arraylist.New()
	// 遍历DID列表，将每个DID信息转换为DID列表视图
	for _, didInfo := range didList {
		didView := &entity.DIDListView{
			DID:          didInfo.DID,
			ContractName: didInfo.ContractName,
			ContractAddr: didInfo.ContractAddr,
			IsIssuer:     didInfo.IsIssuer,
			Status:       didInfo.Status,
			Timestamp:    didInfo.Timestamp,
		}
		listView.Add(didView)
	}
	// 返回DID列表视图
	ConvergeListResponse(ctx, listView.Values(), totalCount, nil)
}

// GetDIDSetHistoryHandler get
type GetDIDSetHistoryHandler struct {
}

func (handler *GetDIDSetHistoryHandler) Handle(ctx *gin.Context) {
	// 绑定请求参数
	params := entity.BindGetDIDHistoryHandler(ctx)
	// 如果参数为空或者不合法，则返回错误响应
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, entity.ErrorMsgParam)
		ConvergeFailureResponse(ctx, newError)
		return
	}

	didDetail, err := dbhandle.GetDIDDetailById(params.ChainId, params.DID)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	// 从数据库中获取DID列表和总数
	didList, totalCount, err := dbhandle.GetDIDHistoryAndCount(params.Offset, params.Limit,
		params.ChainId, params.DID)
	// 如果获取失败，则返回错误响应
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	// 创建一个空的DID列表视图
	listView := arraylist.New()
	// 遍历DID列表，将每个DID信息转换为DID列表视图
	for _, didInfo := range didList {
		// 根据创建交易ID和当前交易ID判断操作类型
		operationType := OperationTypeUpdate
		if didInfo.TxId == didDetail.CreateTxId {
			operationType = OperationTypeCreate
		}

		// 创建DID历史视图
		didView := &entity.DIDHistoryView{
			DID:           didInfo.DID,
			Document:      didInfo.Document,
			ContractName:  didInfo.ContractName,
			ContractAddr:  didInfo.ContractAddr,
			TxId:          didInfo.TxId,
			OperationType: operationType,
			Timestamp:     didInfo.Timestamp,
		}
		listView.Add(didView)
	}
	// 返回DID列表视图
	ConvergeListResponse(ctx, listView.Values(), totalCount, nil)
}
