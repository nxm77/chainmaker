/*
Package service comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package service

import (
	"chainmaker_web/src/auth"
	"chainmaker_web/src/config"
	"chainmaker_web/src/entity_cross"
	"fmt"
	"strconv"
	"time"

	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"

	"chainmaker_web/src/entity"
	loggers "chainmaker_web/src/logger"
	"chainmaker_web/src/monitor_prometheus"
	"chainmaker_web/src/utils"
)

var handlerMap *hashmap.Map
var (
	log      = loggers.GetLogger(loggers.MODULE_WEB)
	apiHisto *prometheus.HistogramVec
	apiTotal *prometheus.CounterVec
)

// init map for key cmb-key and value = cmb-value
func init() {
	handlerMap = hashmap.New()
	// 将对应的处理器加入

	//登录接口
	initLoginHandlers()

	//首页接口
	initOverviewHandlers()

	//交易接口
	initTransactionHandlers()

	//主子链接口
	initCrossSubChainHandlers()

	//用户列表
	handlerMap.Put(entity.GetUserList, &GetUserListHandler{})
	//节点列表
	handlerMap.Put(entity.GetNodeList, &GetNodeListHandler{})
	//组织列表
	handlerMap.Put(entity.GetOrgList, &GetOrgListHandler{})
	//账户列表
	handlerMap.Put(entity.GetAccountList, &GetAccountListHandler{})
	handlerMap.Put(entity.GetAccountDetail, &GetAccountDetailHandler{})

	//区块详情
	handlerMap.Put(entity.GetBlockDetail, &GetBlockDetailHandler{})
	//区块列表
	handlerMap.Put(entity.GetBlockList, &GetBlockListHandler{})
	//区块交易依赖关系
	handlerMap.Put(entity.GetTxDependencies, &GetTxDependenciesHandler{})

	//事件列表
	handlerMap.Put(entity.GetContractTypes, &GetContractTypesHandler{})
	//获取合约类交易列表
	handlerMap.Put(entity.GetContractVersionList, &GetContractVersionListHandler{})
	//合约列表
	handlerMap.Put(entity.GetContractList, &GetContractListHandler{})
	//合约详情
	handlerMap.Put(entity.GetContractDetail, &GetContractDetailHandler{})
	//事件列表
	handlerMap.Put(entity.GetEventList, &GetEventListHandler{})

	//合约ABI
	handlerMap.Put(entity.UploadContractABI, &UploadContractAbiHandler{})
	handlerMap.Put(entity.GetContractABIData, &GetContractABIDataHandler{})
	handlerMap.Put(entity.GetDecodeContractEvents, &GetDecodeContractEventsHandler{})
	handlerMap.Put(entity.GetContractTopics, &GetContractTopicsHandler{})

	//合约调用
	handlerMap.Put(entity.GetCrossContractCalls, &GetCrossContractCallsHandler{})

	//FT合约列表
	handlerMap.Put(entity.GetFTContractList, &GetFTContractListHandler{})
	handlerMap.Put(entity.GetFTContractDetail, &GetFTContractDetailHandler{})
	handlerMap.Put(entity.GetFTPositionList, &GetFTPositionListHandler{})
	handlerMap.Put(entity.GetUserFTPositionList, &GetUserFTPositionListHandler{})
	handlerMap.Put(entity.GetFTTransferList, &GetFungibleTransferListHandler{})

	//NFT合约列表
	handlerMap.Put(entity.GetNFTContractList, &GetNFTContractListHandler{})
	handlerMap.Put(entity.GetNFTContractDetail, &GetNFTContractDetailHandler{})
	handlerMap.Put(entity.GetNFTPositionList, &GetNonFungiblePositionListHandler{})
	handlerMap.Put(entity.GetNFTTransferList, &GetNonFungibleTransferListHandler{})
	handlerMap.Put(entity.GetNFTList, &GetNFTListHandler{})
	handlerMap.Put(entity.GetNFTDetail, &GetNFTDetailHandler{})

	//handlerMap.Put(entity.GetEvidenceContractList, &GetEvidenceContractListHandler{})
	handlerMap.Put(entity.GetEvidenceContract, &GetEvidenceContractHandler{})
	//handlerMap.Put(entity.GetIdentityContractList, &GetIdentityContractListHandler{})
	handlerMap.Put(entity.GetIdentityContract, &GetIdentityContractHandler{})

	//gas接口
	handlerMap.Put(entity.GetGasList, &GetGasListHandler{})
	handlerMap.Put(entity.GetGasRecordList, &GetGasRecordListHandler{})
	handlerMap.Put(entity.GetGasInfo, &GetGasInfoHandler{})

	//订阅链
	handlerMap.Put(entity.SubscribeChain, &SubscribeChainHandler{})
	//删除订阅
	handlerMap.Put(entity.DeleteSubscribe, &DeleteSubscribeHandler{})
	//修改订阅
	handlerMap.Put(entity.ModifySubscribe, &ModifySubscribeHandler{})
	//暂停订阅
	handlerMap.Put(entity.CancelSubscribe, &CancelSubscribeHandler{})

	//更新操作
	//交易加入黑名单
	handlerMap.Put(entity.ModifyTxBlackList, &ModifyTxBlackListHandler{})
	//handlerMap.Put(entity.DeleteTxBlackList, &DeleteTxBlackListHandler{})

	handlerMap.Put(entity.ModifyUserStatus, &ModifyUserStatusHandler{})

	//更新敏感词
	handlerMap.Put(entity.UpdateTxSensitiveWord, &UpdateTxSensitiveWordHandler{})
	handlerMap.Put(entity.UpdateEventSensitiveWord, &UpdateEventSensitiveWordHandler{})
	handlerMap.Put(entity.UpdateEvidenceSensitiveWord, &UpdateEvidenceSensitiveWordHandler{})
	handlerMap.Put(entity.UpdateNFTSensitiveWord, &UpdateNFTSensitiveWordHandler{})
	handlerMap.Put(entity.UpdateContractNameSensitiveWord, &UpdateContractNameSensitiveWordHandler{})
	handlerMap.Put(entity.GetQueryTxList, &GetQueryTxListHandler{})

	//数据要素接口
	handlerMap.Put(entity.GetIDAContractList, &GetIDAContractListHandler{})
	handlerMap.Put(entity.GetIDADataList, &GetIDADataListHandler{})
	handlerMap.Put(entity.GetIDADataDetail, &GetIDADataDetailHandler{})

	//合约验证
	//获取合约版本
	handlerMap.Put(entity.GetContractVersions, &GetContractVersionsHandler{})
	//获取编译器版本
	handlerMap.Put(entity.GetCompilerVersions, &GetCompilerVersionsHandler{})
	//获取evm版本
	handlerMap.Put(entity.GetEvmVersions, &GetEvmVersionsHandler{})
	//获取开源许可证
	handlerMap.Put(entity.GetOpenLicenseTypes, &GetOpenLicenseTypesHandler{})
	//合约验证
	handlerMap.Put(entity.VerifyContractSourceCode, &VerifyContractSourceCodeHandler{})
	//获取合约源码
	handlerMap.Put(entity.GetContractCode, &GetContractCodeHandler{})

	//获取go IDE版本
	handlerMap.Put(entity.GetGoIDEVersions, &GetGoIDEVersionsHandler{})

	//DID相关接口
	initDIDContractHandlers()

	// 打印目前加载的所有处理Handler
	keys := handlerMap.Keys()
	for _, k := range keys {
		if value, ok := handlerMap.Get(k); ok {
			fmt.Printf("Load handler[%s] -> [%T] \n", k, value)
		}
	}
	initMonitorSubName(utils.MonitorNameSpace)
}

func initLoginHandlers() {
	//PluginLogin 浏览器插件登录
	handlerMap.Put(entity.PluginLogin, &PluginLoginHandler{})
	//PluginLogin 账户登录
	handlerMap.Put(entity.AccountLogin, &AccountLoginHandler{})
	//Logout 登录退出接口
	handlerMap.Put(entity.Logout, &LogoutHandler{})
	//监测登录是否有效
	handlerMap.Put(entity.CheckLogin, &CheckLoginHandler{})
}

// initDIDContractHandlers 初始化DID合约相关接口
func initDIDContractHandlers() {
	// 获取DID列表
	handlerMap.Put(entity.GetDIDList, &GetDIDListHandler{})
	// 获取DID详情
	handlerMap.Put(entity.GetDIDDetail, &GetDIDDetailHandler{})
	// 获取DID设置历史
	handlerMap.Put(entity.GetDIDSetHistory, &GetDIDSetHistoryHandler{})
	// 获取VC发行历史
	handlerMap.Put(entity.GetVcIssueHistory, &GetVcIssueHistoryHandler{})
	// 获取VC模板列表
	handlerMap.Put(entity.GetVcTemplateList, &GetVcTemplateListHandler{})
}

// initCrossSubChainHandlers 初始化跨链相关接口
func initCrossSubChainHandlers() {
	handlerMap.Put(entity_cross.GetMainCrossConfig, &GetMainCrossConfigHandler{})
	handlerMap.Put(entity_cross.CrossSearch, &CrossSearchHandler{})
	handlerMap.Put(entity_cross.CrossOverviewData, &CrossOverviewDataHandler{})
	handlerMap.Put(entity_cross.CrossLatestTxList, &CrossLatestTxListHandler{})
	handlerMap.Put(entity_cross.CrossLatestSubChainList, &CrossLatestSubChainListHandler{})
	handlerMap.Put(entity_cross.GetCrossTxList, &GetCrossTxListHandler{})
	handlerMap.Put(entity_cross.CrossSubChainList, &CrossSubChainListHandler{})
	handlerMap.Put(entity_cross.CrossSubChainDetail, &CrossSubChainDetailHandler{})
	handlerMap.Put(entity_cross.GetCrossTxDetail, &GetCrossTxDetailHandler{})
	handlerMap.Put(entity_cross.SubChainCrossChainList, &SubChainCrossChainListHandler{})
	handlerMap.Put(entity_cross.GetCrossSyncTxDetail, &GetCrossSyncTxDetailHandler{})
}

// initTransactionHandlers 初始化交易相关接口
func initTransactionHandlers() {
	// 交易列表
	handlerMap.Put(entity.GetTxList, &GetTxListHandler{})
	// 合约交易列表
	handlerMap.Put(entity.GetContractTxList, &GetContractTxListHandler{})
	// 区块交易列表
	handlerMap.Put(entity.GetBlockTxList, &GetBlockTxListHandler{})
	// 用户交易列表
	handlerMap.Put(entity.GetUserTxList, &GetUserTxListHandler{})
	// 交易列表
	handlerMap.Put(entity.GetQueryTxList, &GetQueryTxListHandler{})
	// 交易详情
	handlerMap.Put(entity.GetTxDetail, &GetTxDetailHandler{})
}

// / initOverviewHandlers 初始化首页相关接口
func initOverviewHandlers() {
	//首页数据
	handlerMap.Put(entity.GetOverviewData, &GetOverviewDataHandler{})
	handlerMap.Put(entity.Decimal, &DecimalHandler{})
	//首页搜索
	handlerMap.Put(entity.Search, &SearchHandler{})
	// 根据时间获取交易数量
	handlerMap.Put(entity.GetTxNumByTime, &GetTxNumByTimeHandler{})
	//链列表
	handlerMap.Put(entity.GetChainList, &GetChainListHandler{})
	//链配置信息
	handlerMap.Put(entity.GetChainConfig, &GetChainConfigHandler{})
	//最新区块列表
	handlerMap.Put(entity.GetLatestBlockList, &GetLatestBlockListHandler{})

	//最新交易列表
	handlerMap.Put(entity.GetLatestTxList, &GetLatestTxListHandler{})
	//最新合约列表
	handlerMap.Put(entity.GetLatestContractList, &GetLatestContractListHandler{})
}

// Dispatcher 分发，需要判断具体的业务
func Dispatcher(ctx *gin.Context) {
	start := time.Now()
	apiPath := ctx.Request.URL.Path
	contextHandler, param := ParseUrl(ctx)
	if contextHandler == nil {
		// 返回错误信息
		err := entity.NewError(entity.ErrorAuthFailure, "can not find this API")
		ConvergeFailureResponse(ctx, err)
		return
	}

	// 根据处理器的需求动态地应用 TokenAuthMiddleware
	if authHandler, ok := contextHandler.(AuthRequiredHandler); ok && authHandler.RequiresAuth() {
		// 登录验证插件，用于需要进行登录验证的接口，获取登录用户的 userId
		authMiddleware := auth.TokenAuthMiddleware()
		authMiddleware(ctx)
		if ctx.IsAborted() {
			return
		}
	}

	contextHandler.Handle(ctx)
	statusCode := strconv.Itoa(ctx.Writer.Status())
	used := time.Since(start).Seconds()
	path := apiPath + "/" + param
	apiHisto.WithLabelValues(path, statusCode).Observe(used)
	apiTotal.WithLabelValues(path, statusCode).Inc()
}

// ParseUrl 解析Url
func ParseUrl(ctx *gin.Context) (ContextHandler, string) {
	param, ok := ctx.GetQuery(entity.CMB)
	//log.Infof("Receive http request[%s]", ctx.Request.URL.String())
	if !ok {
		return nil, param
	}
	if handler, ok := handlerMap.Get(param); ok {
		if handlerVal, ok := handler.(ContextHandler); ok {
			return handlerVal, param
		}
		return nil, param
	}
	return nil, param
}

func initMonitorSubName(subName string) {
	//apiGauage = monitor.NewGaugeVec(subName, "http_process_duration", "request process time", "path")
	apiHisto = monitor_prometheus.NewHistogramVec(subName, "http_histogram", "process consume time histogram",
		[]float64{0.01, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9, 1.0, 2.0, 3.0, 4.0, 5.0, 10.0, 20.0},
		"path", "statusCode")
	apiTotal = monitor_prometheus.NewCounterVec(subName, "http_total", "request total count", "path", "statusCode")
}

// ValidateAPIKey 函数，用于验证API密钥
func ValidateAPIKey(apiKey string) bool {
	// 在这里验证API密钥。例如，查询数据库以检查API密钥是否与已注册的用户匹配。
	// 这里仅为演示目的，我们将API密钥硬编码为"123456"
	return apiKey == config.GlobalConfig.WebConf.ManageBackendApiKey
}
