/*
Package service comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package service

import (
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/entity"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gin-gonic/gin"
)

// GetVcIssueHistoryHandler get
type GetVcIssueHistoryHandler struct {
}

// GetVcIssueHistory handle
func (handler *GetVcIssueHistoryHandler) Handle(ctx *gin.Context) {
	// 绑定请求参数
	params := entity.BindGetVCListHandler(ctx)
	// 如果参数为空或者不合法，则返回错误响应
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, entity.ErrorMsgParam)
		ConvergeFailureResponse(ctx, newError)
		return
	}

	// 从数据库中获取DID列表和总数
	vcList, totalCount, err := dbhandle.GetVCListAndCount(params.Offset, params.Limit,
		params.ChainId, params.IssuerDID, params.HolderDID, params.TemplateID, params.VcID, params.ContractAddr)
	// 如果获取失败，则返回错误响应
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	// 创建一个空的DID列表视图
	listView := arraylist.New()
	// 遍历DID列表，将每个DID信息转换为DID列表视图
	for _, vcInfo := range vcList {
		vcView := &entity.VCListView{
			VcID:       vcInfo.VCID,
			IssuerDID:  vcInfo.IssuerDID,
			HolderDID:  vcInfo.HolderDID,
			TemplateID: vcInfo.TemplateID,
			Status:     vcInfo.Status,
			Timestamp:  vcInfo.Timestamp,
		}
		listView.Add(vcView)
	}
	// 返回DID列表视图
	ConvergeListResponse(ctx, listView.Values(), totalCount, nil)
}

// GetVcTemplateListHandler get
type GetVcTemplateListHandler struct {
}

// GetVcTemplateListHandler handle
func (handler *GetVcTemplateListHandler) Handle(ctx *gin.Context) {
	// 绑定请求参数
	params := entity.BindGetVCTemplateHandler(ctx)
	// 如果参数为空或者不合法，则返回错误响应
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, entity.ErrorMsgParam)
		ConvergeFailureResponse(ctx, newError)
		return
	}

	// 从数据库中获取DID列表和总数
	temolateList, totalCount, err := dbhandle.GetVCTempListAndCount(params.Offset, params.Limit,
		params.ChainId, params.ContractAddr, params.TemplateID)
	// 如果获取失败，则返回错误响应
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	// 创建一个空的DID列表视图
	listView := arraylist.New()
	// 遍历DID列表，将每个DID信息转换为DID列表视图
	for _, value := range temolateList {
		templateView := &entity.TemplateView{
			TemplateID:   value.TemplateID,
			TemplateName: value.TemplateName,
			ContractName: value.ContractName,
			ContractAddr: value.ContractAddr,
			Version:      value.Version,
			TxId:         value.TxId,
			TemplateJson: value.Template,
			Timestamp:    value.Timestamp,
		}
		listView.Add(templateView)
	}
	// 返回DID列表视图
	ConvergeListResponse(ctx, listView.Values(), totalCount, nil)
}
