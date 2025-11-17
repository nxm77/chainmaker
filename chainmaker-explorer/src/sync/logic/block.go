/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package logic

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/sync/common"
	"chainmaker_web/src/sync/model"
	"chainmaker_web/src/utils"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	pbCommon "chainmaker.org/chainmaker/pb-go/v2/common"
)

// DealBlockInfo 处理block信息
// @param blockInfo block信息
// @param hashType   hash类型
// @return *db.Block block信息
// @return error 错误信息
func DealBlockInfo(blockInfo *pbCommon.BlockInfo, hashType string, dealResult *model.ProcessedBlockData) error {
	if blockInfo == nil {
		return fmt.Errorf("blockInfo is nil")
	}
	chainId := blockInfo.Block.Header.ChainId
	blockHeight := blockInfo.Block.Header.BlockHeight

	newUUID := uuid.New().String()
	timestamp := blockInfo.Block.Header.BlockTimestamp
	//根据block信息填充block中的信息
	modBlock := &db.Block{
		ID: newUUID,
		//nolint:gosec
		BlockHeight: int64(blockHeight),
		BlockHash:   hex.EncodeToString(blockInfo.Block.Header.BlockHash),
		//nolint:gosec
		BlockVersion:  int32(blockInfo.Block.Header.BlockVersion),
		PreBlockHash:  hex.EncodeToString(blockInfo.Block.Header.PreBlockHash),
		ConsensusArgs: utils.Base64Encode(blockInfo.Block.Header.ConsensusArgs),
		DagHash:       hex.EncodeToString(blockInfo.Block.Header.DagHash),
		Timestamp:     timestamp,
		TimestampDate: utils.GetDateFromTimestamp(timestamp),
		TxCount:       int(blockInfo.Block.Header.TxCount),
	}
	modBlock.RwSetHash = hex.EncodeToString(blockInfo.Block.Header.RwSetRoot)
	modBlock.Signature = utils.Base64Encode(blockInfo.Block.Header.Signature)
	modBlock.TxRootHash = hex.EncodeToString(blockInfo.Block.Header.TxRoot)
	member := blockInfo.Block.Header.Proposer
	//根据proposer信息填充block中的地址,id信息
	if member != nil {
		modBlock.OrgId = member.OrgId
		//根据proposer信息填充block中的地址,id信息
		getInfos, err := common.GetMemberIdAddrAndCertNew(chainId, hashType, member)
		if err != nil {
			log.Error("getMemberIdAddrAndCert Failed: " + err.Error())
			return err
		}
		modBlock.ProposerAddr = getInfos.UserAddr
		modBlock.ProposerId = getInfos.UserId
	}

	//解析block中的dag信息
	dagBytes, _ := json.Marshal(blockInfo.Block.Dag)
	modBlock.BlockDag = string(dagBytes)
	dealResult.BlockDetail = modBlock
	return nil
}

// BuildLatestBlockListCache
//
//	@Description:设置最新区块缓存列表
//	@param chainId
//	@param modBlock 区块数据
func BuildLatestBlockListCache(chainId string, modBlock *db.Block) {
	var blockList []*db.Block
	//获取最新区块列表
	blockListCache, _ := dbhandle.GetLatestBlockListCache(chainId)
	if len(blockListCache) > 0 {
		//缓存存在
		blockList = append(blockList, modBlock)
	} else {
		//缓存可能丢失
		blockList, _ = dbhandle.GetLatestBlockListCache(chainId)
	}
	if len(blockList) == 0 {
		return
	}
	// 缓存交易信息
	dbhandle.SetLatestBlockListCache(chainId, blockList)
}

// BuildOverviewMaxBlockHeightCache
//
//	@Description: 缓存最高区块高度
//	@param chainId
//	@param blockInfo
func BuildOverviewMaxBlockHeightCache(chainId string, blockInfo *db.Block) {
	maxBlockHeight := blockInfo.BlockHeight
	dbhandle.SetMaxBlockHeightCache(chainId, maxBlockHeight)
}

// DealBlockchainStatistics 处理区块链统计信息
// @param chainId 区块链id
// @param maxHeight 区块高度
// @param txList 交易列表
// @param insertAccount 账户信息
// @return *db.Statistics 区块链统计信息
func DealBlockchainStatistics(chainId string, insertAccount []*db.Account) *db.Statistics {
	//获取区块链统计信息
	statistics, err := dbhandle.GetChainStatistics(chainId)
	if err != nil {
		log.Errorf("DealStatisticsRealtime Get chain statistics err : %s", err.Error())
		return nil
	}

	statistics.TotalAccounts += int64(len(insertAccount))
	orgs, err := dbhandle.GetOrgNum(chainId)
	if err == nil {
		statistics.TotalOrgs = orgs
	}
	nodes, err := dbhandle.GetNodeNum(chainId, "")
	if err == nil {
		statistics.TotalNodes = nodes
	}

	return statistics
}

// DealStatisticsRealtime 处理区块链统计信息
// @param chainId 区块链id
// @param maxHeight 区块高度
// @param txList 交易列表
// @param insertAccount 账户信息
// @return *db.Statistics 区块链统计信息
func DealStatisticsRealtime(chainId string, blockHeight int64, dealResult *model.ProcessedBlockData) error {
	if dealResult == nil {
		log.Errorf("DealStatisticsRealtime dealResult is nil")
		return nil
	}

	//获取区块链统计信息
	statistics, err := dbhandle.GetChainStatistics(chainId)
	if err != nil {
		log.Errorf("DealStatisticsRealtime Get chain statistics err : %s", err.Error())
		return nil
	}

	if blockHeight <= statistics.BlockHeight {
		return nil
	}
	statistics.BlockHeight = blockHeight
	if len(dealResult.Transactions) > 0 {
		statistics.TotalTransactions += int64(len(dealResult.Transactions))
	}
	crossResult := dealResult.CrossChainResult
	if crossResult != nil && len(crossResult.InsertCrossTransfer) > 0 {
		statistics.TotalCrossTx += int64(len(crossResult.InsertCrossTransfer))
	}
	if len(dealResult.InsertContracts) > 0 {
		statistics.TotalContracts += int64(len(dealResult.InsertContracts))
	}

	// 更新统计信息
	err = dbhandle.UpdateStatisticsRealtime(chainId, statistics)
	return err
}
