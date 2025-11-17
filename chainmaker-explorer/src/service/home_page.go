/*
Package service comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package service

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/entity"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gin-gonic/gin"
)

const (
	//BLOCK searchType
	BLOCK = iota
	// TRANSACTION searchType
	TRANSACTION
	// CONTRACT searchType
	CONTRACT
	//ACCOUNT searchType
	ACCOUNT
	//DID searchType
	DID
	// UNKNOWN searchType
	UNKNOWN = -1
)

const (
	SearchBlockHash      = "1"
	SearchBlockHeight    = "2"
	SearchTransaction    = "3"
	SearchContractName   = "4"
	SearchContractAddr   = "5"
	SearchAccountAddress = "6"
	SearchAccountBNS     = "7"
	SearchDID            = "8"
	SearchVCID           = "9"
)

type GetOverviewData struct {
	BlockHeight   int64
	UserCount     int64
	ContractCount int64
	TxCount       int64
	OrgCount      int64
	RunningNode   int64
	CommonNode    int64
	ConsensusNode int64
}

// DecimalHandler dec
type DecimalHandler struct{}

// Handle deal
func (decimalHandler *DecimalHandler) Handle(ctx *gin.Context) {
	params := entity.BindChainOverviewDataHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetOverviewData param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}
	chainId := params.ChainId
	overviewData := getOverviewData(ctx, chainId)

	var timeObj = time.Now()
	recentTime := timeObj.Add(-(24 * time.Hour))
	recenBlocktCount, err := dbhandle.GetBlockListByRange(chainId, recentTime.Unix(), timeObj.Unix())
	if err != nil {
		log.Errorf("decimalHandler GetBlockListByRange err : %v", err)
	}
	recentTxCount, err := dbhandle.GetTxCountByRange(chainId, recentTime.Unix(), timeObj.Unix())
	if err != nil {
		log.Errorf("decimalHandler GetTxCountByRange err : %v", err)
	}
	recentContractCount, err := dbhandle.GetContractCountByRange(chainId, recentTime.Unix(), timeObj.Unix())
	if err != nil {
		log.Errorf("decimalHandler GetcontractCountByRange err : %v", err)
	}

	recentUserCount, err := dbhandle.GetUserCountByRange(chainId, recentTime.Unix(), timeObj.Unix())
	if err != nil {
		log.Errorf("decimalHandler GetUserCountByRange err : %v", err)
	}

	decimalViewData := &entity.DecimalView{
		BlockHeight:       overviewData.BlockHeight,
		TxNum:             overviewData.TxCount,
		ContractNum:       overviewData.ContractCount,
		UserNum:           overviewData.UserCount,
		NodeNum:           overviewData.RunningNode,
		OrgNum:            overviewData.OrgCount,
		RecentBlockHeight: recenBlocktCount,
		RecentTxNum:       recentTxCount,
		RecentContractNum: recentContractCount,
		RecentUserNum:     recentUserCount,
	}

	//返回response
	ConvergeDataResponse(ctx, decimalViewData, nil)
}

// GetOverviewDataHandler Handler
type GetOverviewDataHandler struct{}

// Handle 首页概览数据
func (GetOverviewDataHandler *GetOverviewDataHandler) Handle(ctx *gin.Context) {
	params := entity.BindChainOverviewDataHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "GetOverviewData param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}
	chainId := params.ChainId
	//获取缓存数据
	overviewData := GetOverviewDataCache(ctx, chainId)
	if overviewData != nil {
		ConvergeDataResponse(ctx, overviewData, nil)
		return
	}

	resultData := getOverviewData(ctx, chainId)
	overviewData = &entity.OverviewDataView{
		ChainId:       chainId,
		BlockHeight:   resultData.BlockHeight,
		UserCount:     resultData.UserCount,
		ContractCount: resultData.ContractCount,
		OrgCount:      resultData.OrgCount,
		TxCount:       resultData.TxCount,
		RunningNode:   resultData.RunningNode,
		CommonNode:    resultData.CommonNode,
		ConsensusNode: resultData.ConsensusNode,
	}

	//返回response
	ConvergeDataResponse(ctx, overviewData, nil)
}

func getOverviewData(ctx *gin.Context, chainId string) *entity.OverviewDataView {
	//获取缓存数据
	overviewData := GetOverviewDataCache(ctx, chainId)
	if overviewData != nil {
		return overviewData
	}
	// //区块高度
	// blockHeight := dbhandle.GetMaxBlockHeight(chainId)
	// //组织数量
	// orgCount, err := dbhandle.GetOrgNum(chainId)
	// if err != nil {
	// 	log.Errorf("Count recent org num err : %s", err.Error())
	// }

	// //交易数量
	// txCount, err := dbhandle.GetTotalTxNumByAccount(chainId)
	// if err != nil {
	// 	log.Errorf("Count tx num err : %s", err.Error())
	// }

	// //合约数量
	// contractCount, err := dbhandle.GetContractNum(chainId)
	// if err != nil {
	// 	log.Errorf("Count contract num err : %s", err.Error())
	// }

	// //账户数量
	// accountCount, err := dbhandle.GetAccountTotal(chainId)
	// if err != nil {
	// 	log.Errorf("Count user num err : %s", err.Error())
	// }

	//获取区块链统计信息
	statistics, err := dbhandle.GetChainStatistics(chainId)
	if err != nil {
		log.Errorf("Get chain statistics err : %s", err.Error())
		statistics = &db.Statistics{}
	}
	runningNode, err := dbhandle.GetNodeNum(chainId, "")
	if err != nil {
		log.Errorf("Count recent node num err : %s", err.Error())
	}

	consensusNodeCount, err := dbhandle.GetNodeNum(chainId, "consensus")
	if err != nil {
		log.Errorf("Count recent node num err : %s", err.Error())
	}
	commonNodeCount, err := dbhandle.GetNodeNum(chainId, "common")
	if err != nil {
		log.Errorf("Count recent node num err : %s", err.Error())
	}
	overviewData = &entity.OverviewDataView{
		ChainId:       chainId,
		BlockHeight:   statistics.BlockHeight,
		UserCount:     statistics.TotalAccounts,
		ContractCount: statistics.TotalContracts,
		OrgCount:      statistics.TotalOrgs,
		TxCount:       statistics.TotalTransactions,
		RunningNode:   runningNode,
		CommonNode:    commonNodeCount,
		ConsensusNode: consensusNodeCount,
	}

	//设置缓存
	SetOverviewDataCache(ctx, chainId, *overviewData)
	return overviewData
}

// SetOverviewDataCache 缓存首页信息
func SetOverviewDataCache(ctx *gin.Context, chainId string, overviewData entity.OverviewDataView) {
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisHomeOverviewData, prefix, chainId)
	retJson, err := json.Marshal(overviewData)
	if err != nil {
		log.Errorf("【Redis】set cache failed, key:%v, result:%v", redisKey, retJson)
		return
	}
	// 设置键值对和过期时间(1s 过期)
	_ = cache.GlobalRedisDb.Set(ctx, redisKey, string(retJson), time.Second).Err()
}

// GetOverviewDataCache 获取首页缓存数据
func GetOverviewDataCache(ctx *gin.Context, chainId string) *entity.OverviewDataView {
	cacheResult := &entity.OverviewDataView{}
	//从缓存获取最新的block
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisHomeOverviewData, prefix, chainId)
	redisRes := cache.GlobalRedisDb.Get(ctx, redisKey)
	if redisRes == nil || redisRes.Val() == "" {
		return nil
	}

	err := json.Unmarshal([]byte(redisRes.Val()), &cacheResult)
	if err != nil {
		log.Errorf("【Redis】get cache failed, key:%v, result:%v", redisKey, redisRes)
		return nil
	}
	return cacheResult
}

// SearchHandler handler
type SearchHandler struct{}

// Handle 首页搜索
func (searchHandler *SearchHandler) Handle(ctx *gin.Context) {
	// 绑定请求参数
	params := entity.BindGetDetailHandler(ctx)
	// 如果参数为空或者不合法，则返回错误
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "Search param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	// 获取链ID
	chainId := params.ChainId
	// 获取搜索值
	viewType, viewValue, contractType, err := searchHandler.getSearchValue(chainId, params.Type, params.Value)
	// 如果搜索值为空或者出错，则设置viewType为UNKNOWN
	if err != nil || viewValue == "" {
		//未找到
		viewType = UNKNOWN
	}

	// 创建搜索视图
	searchView := &entity.SearchView{
		Type:         viewType,
		Data:         viewValue,
		ChainId:      chainId,
		ContractType: contractType,
	}

	//返回response
	ConvergeDataResponse(ctx, searchView, nil)
}

// 根据传入的chainId、searchType和searchValue，获取搜索结果
func (handle *SearchHandler) getSearchValue(chainId, searchType, searchValue string) (int, string, string, error) {
	// 定义返回值
	var (
		viewType     int
		viewValue    string
		contractType string
		err          error
	)

	// 根据searchType的值，执行不同的操作
	switch searchType {
	case SearchBlockHash, SearchBlockHeight:
		// 获取区块信息
		viewType, viewValue, err = getBlockInfoValue(chainId, searchType, searchValue)
	case SearchTransaction:
		// 获取交易信息
		var txInfo *db.Transaction
		txInfo, err = dbhandle.GetTransactionByTxId(searchValue, chainId)
		if err == nil && txInfo != nil {
			viewType = TRANSACTION
			viewValue = txInfo.TxId
		}
	case SearchContractName, SearchContractAddr:
		// 获取合约信息
		viewType, viewValue, contractType, err = getContractValue(chainId, searchType, searchValue)
	case SearchAccountAddress, SearchAccountBNS:
		// 获取账户信息
		viewType, viewValue, err = getAccountValue(chainId, searchType, searchValue)
	case SearchDID:
		// 获取DID信息
		viewType, viewValue, err = getDIDValue(chainId, searchValue)
	case SearchVCID:
		// 获取VCID信息
		viewType, viewValue, err = getDIDByVCId(chainId, searchValue)
	default:
		// 参数错误
		err = entity.NewError(entity.ErrorParamWrong, "Search param is wrong")
	}

	// 返回结果
	return viewType, viewValue, contractType, err
}

// 根据chainId、paramsType和paramsValue获取区块信息
func getBlockInfoValue(chainId, paramsType, paramsValue string) (int, string, error) {
	var (
		err       error
		blockInfo *db.Block
	)

	// 如果paramsType是SearchBlockHash，则根据paramsValue获取区块信息
	if paramsType == SearchBlockHash {
		blockInfo, err = dbhandle.GetBlockByHash(paramsValue, chainId)
	} else {
		// 否则，将paramsValue转换为int64类型，表示区块高度，根据chainId和区块高度获取区块信息
		blockHeight, _ := strconv.ParseInt(paramsValue, 10, 64)
		//说明id是数字类型的，那就表示区块高度，其他的都是字符串
		blockInfo, err = dbhandle.GetBlockByHeight(chainId, blockHeight)
	}
	// 如果没有错误且blockInfo不为空，则返回BLOCK和blockInfo的BlockHash
	if err == nil && blockInfo != nil {
		return BLOCK, blockInfo.BlockHash, nil
	}

	// 否则，返回UNKNOWN和空字符串，并返回错误
	return UNKNOWN, "", err
}

// 根据链ID、参数类型和参数值获取合约信息
func getContractValue(chainId, paramsType, paramsValue string) (int, string, string, error) {
	// 定义错误变量和合约信息变量
	var (
		err          error
		contractInfo *db.Contract
	)

	// 根据参数类型获取合约信息
	if paramsType == SearchContractName {
		contractInfo, err = dbhandle.GetContractByName(chainId, paramsValue)
	} else {
		contractInfo, err = dbhandle.GetContractByAddr(chainId, paramsValue)
	}

	// 如果没有错误且合约信息不为空，则返回合约信息
	if err == nil && contractInfo != nil {
		return CONTRACT, contractInfo.Addr, contractInfo.ContractType, nil
	}
	// 否则返回未知和错误信息
	return UNKNOWN, "", "", err
}

// 根据chainId、paramsType和paramsValue获取账户信息
func getAccountValue(chainId, paramsType, paramsValue string) (int, string, error) {
	// 定义错误变量和账户信息变量
	var (
		err         error
		accountInfo *db.Account
	)

	// 根据paramsType判断获取账户信息的方式
	if paramsType == SearchAccountAddress {
		// 根据地址获取账户信息
		accountInfo, err = dbhandle.GetAccountByAddr(chainId, paramsValue)
	} else {
		// 根据BNS获取账户信息
		accountInfo, err = dbhandle.GetAccountByBNS(chainId, paramsValue)
	}
	// 如果没有错误且账户信息不为空，则返回账户信息和地址
	if err == nil && accountInfo != nil {
		return ACCOUNT, accountInfo.Address, nil
	}

	// 否则返回未知和错误信息
	return UNKNOWN, "", err
}

// 根据chainId和paramsValue获取DID的值
func getDIDValue(chainId, paramsValue string) (int, string, error) {
	// 根据chainId和paramsValue从数据库中获取DID的详细信息
	didInfo, err := dbhandle.GetDIDDetailById(chainId, paramsValue)
	// 如果没有错误并且didInfo不为空，则返回DID的值
	if err == nil && didInfo != nil {
		return DID, didInfo.DID, nil
	}

	// 否则返回UNKNOWN和空字符串，并返回错误
	return UNKNOWN, "", err
}

func getDIDByVCId(chainId, vcId string) (int, string, error) {
	// 根据chainId和paramsValue从数据库中获取DID的详细信息
	vcInfo, err := dbhandle.GetVCInfoById(chainId, vcId)
	// 如果没有错误并且didInfo不为空，则返回DID的值
	if vcInfo == nil {
		return UNKNOWN, "", err
	}

	didInfo, err := dbhandle.GetDIDDetailById(chainId, vcInfo.HolderDID)
	if err == nil && didInfo != nil {
		return DID, didInfo.DID, nil
	}

	// 否则返回UNKNOWN和空字符串，并返回错误
	return UNKNOWN, "", err
}

// GetNodeListHandler get
type GetNodeListHandler struct{}

// Handle deal
func (getNodeListHandler *GetNodeListHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetNodeListHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	// get node
	nodes, totalCount, err := dbhandle.GetNodeList(params.ChainId,
		params.OrgId, params.NodeId, params.Offset, params.Limit)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	// view
	nodeViews := arraylist.New()
	for _, node := range nodes {
		nodeView := &entity.NodesView{
			NodeId:      node.NodeId,
			NodeAddress: node.Address,
			Role:        node.Role,
			OrgId:       node.OrgId,
			Status:      node.Status,
			Timestamp:   node.CreatedAt.Unix(),
		}
		nodeViews.Add(nodeView)
	}
	ConvergeListResponse(ctx, nodeViews.Values(), totalCount, nil)
}

type OrgNodeUserData struct {
	NodeCount int64
	UserCount int64
}

// GetOrgListHandler get
type GetOrgListHandler struct{}

// Handle deal
func (handler *GetOrgListHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetOrgListHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}
	chainId := params.ChainId
	//组织列表
	orgList, totalCount, err := dbhandle.GetOrgList(chainId, params.OrgId, params.Offset, params.Limit)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	// 创建一个互斥锁来保护对 resultMap 的访问
	var mu sync.Mutex
	resultCountMap := make(map[string]OrgNodeUserData)
	// 创建一个 WaitGroup 以等待所有 goroutines 完成
	var wg sync.WaitGroup
	wg.Add(len(orgList))

	for _, org := range orgList {
		go func(orgId string) {
			defer wg.Done() // 在 goroutine 结束时调用
			// 获取节点数和用户数（用你自己的函数替换这里的示例函数）
			nodeCount, nodeErr := dbhandle.GetNodeNumByOrg(chainId, orgId)
			userCount, userErr := dbhandle.GetUserNum(chainId, orgId)
			// 处理错误（根据需要修改这里的错误处理逻辑）
			if nodeErr != nil || userErr != nil {
				newError := entity.NewError(entity.ErrorParamWrong, "get nodeCode userCode err")
				ConvergeFailureResponse(ctx, newError)
				return
			}

			// 使用互斥锁保护对 resultMap 的访问
			mu.Lock()
			resultCountMap[orgId] = OrgNodeUserData{
				NodeCount: nodeCount,
				UserCount: userCount,
			}
			mu.Unlock()
		}(org.OrgId)

	}
	// 等待所有 goroutines 完成
	wg.Wait()

	orgViews := arraylist.New()
	for _, org := range orgList {
		var nodeCount int64
		var userCount int64
		if value, ok := resultCountMap[org.OrgId]; ok {
			nodeCount = value.NodeCount
			userCount = value.UserCount
		}
		orgView := &entity.OrgView{
			ChainId:   chainId,
			OrgId:     org.OrgId,
			Status:    org.Status,
			UserCount: userCount,
			NodeCount: nodeCount,
		}
		orgViews.Add(orgView)
	}

	ConvergeListResponse(ctx, orgViews.Values(), totalCount, nil)
}
