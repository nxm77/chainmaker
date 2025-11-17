/*
Package service comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package service

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/entity"
	"encoding/json"

	pbCommon "chainmaker.org/chainmaker/pb-go/v2/common"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gin-gonic/gin"
)

// GetBlockDetailHandler get
type GetBlockDetailHandler struct {
}

// Handle deal
func (handler *GetBlockDetailHandler) Handle(ctx *gin.Context) {
	var (
		block *db.Block
		err   error
	)
	// 参数处理
	params := entity.BindGetBlockDetailHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetBlockDetail param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	if params.BlockHash != "" {
		block, err = dbhandle.GetBlockByHash(params.BlockHash, params.ChainId)
	} else if params.BlockHeight != nil {
		block, err = dbhandle.GetBlockByHeight(params.ChainId, *params.BlockHeight)
	} else {
		newError := entity.NewError(entity.ErrorParamWrong, "GetBlockDetail param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	if err != nil || block == nil {
		log.Errorf("GetBlockDetail err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	// view
	blockDetailView := &entity.BlockDetailView{
		BlockHash:        block.BlockHash,
		PreBlockHash:     block.PreBlockHash,
		RwSetHash:        block.RwSetHash,
		Timestamp:        block.Timestamp,
		BlockHeight:      block.BlockHeight,
		TxCount:          block.TxCount,
		ProposalNodeId:   block.ProposerId,
		TxRootHash:       block.TxRootHash,
		Dag:              block.DagHash,
		OrgId:            block.OrgId,
		ProposalNodeAddr: block.ProposerAddr,
	}

	ConvergeDataResponse(ctx, blockDetailView, nil)
}

// GetLatestBlockListHandler get
type GetLatestBlockListHandler struct {
}

// Handle deal
func (getLatestBlockListHandler *GetLatestBlockListHandler) Handle(ctx *gin.Context) {
	//参数处理
	params := entity.BindGetBlockLatestListHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetLatestBlockList param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//从缓存获取最新的block
	blockList, err := dbhandle.GetLatestBlockListCache(params.ChainId)
	if err != nil {
		log.Errorf("GetLatestBlockList get redis fail err:%v", err)
	}

	if len(blockList) == 0 {
		// 获取block
		blockList, err = dbhandle.GetLatestBlockList(params.ChainId)
		if err != nil {
			ConvergeHandleFailureResponse(ctx, err)
			return
		}
	}

	blockViews := arraylist.New()
	for i, block := range blockList {
		latestBlockListView := &entity.LatestBlockListView{
			Id:               int64(i + 1),
			BlockHash:        block.BlockHash,
			BlockHeight:      block.BlockHeight,
			TxNum:            block.TxCount,
			ProposalNodeId:   block.ProposerId,
			ProposalNodeAddr: block.ProposerAddr,
			Timestamp:        block.Timestamp,
		}
		blockViews.Add(latestBlockListView)
	}
	ConvergeListResponse(ctx, blockViews.Values(), int64(len(blockList)), nil)
}

// GetBlockListHandler get
type GetBlockListHandler struct {
}

// Handle GetBlockList区块链列表
func (getBlockListHandler *GetBlockListHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetBlockListHandler(ctx)
	if params == nil {
		ConvergeFailureResponse(ctx, entity.NewError(entity.ErrorParamWrong, "参数绑定失败"))
		return
	}

	if err := params.Validate(); err != nil {
		// 直接返回结构化错误
		ConvergeFailureResponse(ctx, entity.NewError(entity.ErrorParamWrong, err.Error()))
		return
	}

	// 获取block数量
	totalCount, err := dbhandle.GetBlockListCount(params.ChainId, params.BlockKey, params.NodeId)
	if err != nil {
		log.Errorf("GetBlockList err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	blocks, err := dbhandle.GetBlockList(params.Offset, params.Limit, params.ChainId, params.BlockKey, params.NodeId)
	if err != nil {
		log.Errorf("GetBlockList err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	// view
	blocksView := arraylist.New()
	for _, block := range blocks {
		blockListView := &entity.BlockListView{
			BlockHeight:      block.BlockHeight,
			BlockHash:        block.BlockHash,
			TxNum:            block.TxCount,
			ProposalNodeId:   block.ProposerId,
			ProposalNodeAddr: block.ProposerAddr,
			Timestamp:        block.Timestamp,
		}
		blocksView.Add(blockListView)
	}
	ConvergeListResponse(ctx, blocksView.Values(), totalCount, nil)
}

// GetTxDependenciesHandler get
type GetTxDependenciesHandler struct {
}

// Handle deal
func (handler *GetTxDependenciesHandler) Handle(ctx *gin.Context) {
	// 参数处理
	params := entity.BindGetTxDependenciesHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetTxDependencies param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	block, err := dbhandle.GetBlockByHeight(params.ChainId, params.BlockHeight)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	if block == nil {
		ConvergeHandleFailureResponse(ctx, entity.ErrSelectFailed)
		return
	}

	dag := &pbCommon.DAG{}
	if block.BlockDag != "" {
		err = json.Unmarshal([]byte(block.BlockDag), &dag)
		if err != nil {
			log.Errorf("GetTxDependencies err : %s", err.Error())
			ConvergeHandleFailureResponse(ctx, err)
			return
		}
	}
	txIds, err := dbhandle.GetBlockTxIdsByHeight(params.ChainId, params.BlockHeight)
	if err != nil {
		log.Errorf("GetTxDependencies err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	// view
	blockDetailView := &entity.BlockTxDependenciesView{
		TxIds: txIds,
		Dag:   dag.Vertexes,
	}

	ConvergeDataResponse(ctx, blockDetailView, nil)
}
