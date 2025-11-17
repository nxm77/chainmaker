package service

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/entity"
	"chainmaker_web/src/entity_cross"
	"chainmaker_web/src/utils"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	tcipCommon "chainmaker.org/chainmaker/tcip-go/v2/common"
	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gin-gonic/gin"
)

const (
	//SearchCrossID searchType
	SearchCrossID = iota
	// SearchSubChainID searchType
	SearchSubChainID
	//SearchCrossUnKnow searchType
	SearchCrossUnKnow = -1
)

// GetMainCrossConfigHandler get
type GetMainCrossConfigHandler struct {
}

// Handle GetMainCrossConfigHandler 主子链网配置,是否是主链
func (handler *GetMainCrossConfigHandler) Handle(ctx *gin.Context) {
	params := entity_cross.BindGetChainIdHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	mainCrossConfigView := &entity_cross.MainCrossConfig{
		ShowTag: utils.GetIsMainChain(),
	}
	//返回response
	ConvergeDataResponse(ctx, mainCrossConfigView, nil)
}

// CrossSearchHandler sub
type CrossSearchHandler struct {
}

// Handle SubscribeChainHandler
func (handler *CrossSearchHandler) Handle(ctx *gin.Context) {
	params := entity_cross.BindGetCrossSearchHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	// 根据searchType的值，执行不同的操作
	var typeView int
	var valueView string
	switch params.Type {
	case SearchCrossID:
		//跨链ID
		crossTxInfo, err := dbhandle.GetCrossCycleTransferById(params.ChainId, params.Value)
		if err == nil && len(crossTxInfo) > 0 {
			valueView = crossTxInfo[0].CrossId
			typeView = SearchCrossID
		}
	case SearchSubChainID:
		subChainInfo, err := dbhandle.GetCrossSubChainInfoById(params.ChainId, params.Value)
		if err == nil && subChainInfo != nil {
			valueView = subChainInfo.SubChainId
			typeView = SearchSubChainID
		}
	default:
		// 参数错误
		err := entity.NewError(entity.ErrorParamWrong, "Search param is wrong")
		ConvergeFailureResponse(ctx, err)
		return
	}

	if valueView == "" {
		typeView = SearchCrossUnKnow
	}
	crossSearchView := &entity_cross.CrossSearchView{
		Type: typeView,
		Data: valueView,
	}
	//返回response
	ConvergeDataResponse(ctx, crossSearchView, nil)
}

// CrossOverviewDataHandler cancel
type CrossOverviewDataHandler struct {
}

// Handle deal
func (handler *CrossOverviewDataHandler) Handle(ctx *gin.Context) {
	params := entity_cross.BindGetChainIdHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, entity.ErrorMsgParam)
		ConvergeFailureResponse(ctx, newError)
		return
	}

	chainId := params.ChainId
	//获取缓存数据
	overviewData := GetCrossOverviewDataCache(ctx, chainId)
	if overviewData != nil {
		ConvergeDataResponse(ctx, overviewData, nil)
		return
	}

	//总区块高度
	totalBlockHeight, err := dbhandle.GetAllSubChainBlockHeight(chainId)
	if err != nil {
		log.Errorf("GetAllSubChainBlockHeight err : %v", err)
	}

	//本自然月的开始，结束时间
	startTime, endTime := GetCurrentMonthStartAndEndTime()
	//周期交易最短完成时间
	transferList, err := dbhandle.GetCrossTransferDurationByTime(chainId, startTime, endTime)
	if err != nil {
		log.Errorf("GetCrossTxCycleListByTime err : %v", err)
	}

	shortestTime, longestTime, averageTime := CalculateStats(transferList)
	//子链总数
	subChainNum, err := dbhandle.GetCrossSubChainAllCount(chainId)
	if err != nil {
		log.Errorf("GetCrossSubChainAllCount err : %v", err)
	}

	//获取区块链统计信息
	statistics, err := dbhandle.GetChainStatistics(chainId)
	if err != nil {
		log.Errorf("Get chain statistics err : %s", err.Error())
		statistics = &db.Statistics{}
	}

	overviewData = &entity_cross.OverviewDataView{
		TotalBlockHeight: totalBlockHeight,
		ShortestTime:     shortestTime,
		LongestTime:      longestTime,
		AverageTime:      averageTime,
		SubChainNum:      subChainNum,
		TxNum:            statistics.TotalCrossTx,
	}

	//设置缓存
	SetCrossOverviewDataCache(ctx, chainId, *overviewData)
	//返回response
	ConvergeDataResponse(ctx, overviewData, nil)
}

// CalculateStats 计算最大值、最小值和平均值
func CalculateStats(cycleList []*db.CrossTransactionTransfer) (int64, int64, int64) {
	if len(cycleList) == 0 {
		return 0, 0, 0
	}

	minDuration := int64(cycleList[0].EndTime - cycleList[0].StartTime)
	maxDuration := int64(cycleList[0].EndTime - cycleList[0].StartTime)
	sumDuration := int64(0)

	for _, cycle := range cycleList {
		duration := int64(cycle.EndTime - cycle.StartTime)
		if duration < minDuration {
			minDuration = duration
		}
		if duration > maxDuration {
			maxDuration = duration
		}
		if duration > 0 {
			sumDuration += duration
		}
	}

	average := float64(sumDuration) / float64(len(cycleList))
	if minDuration < 0 {
		minDuration = 0
	}
	if maxDuration < 0 {
		maxDuration = 0
	}
	return minDuration, maxDuration, int64(average)
}

// GetCrossOverviewDataCache 获取首页缓存数据
func GetCrossOverviewDataCache(ctx *gin.Context, chainId string) *entity_cross.OverviewDataView {
	cacheResult := &entity_cross.OverviewDataView{}
	//从缓存获取最新的block
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisCrossOverviewData, prefix, chainId)
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

// SetCrossOverviewDataCache 缓存首页信息
func SetCrossOverviewDataCache(ctx *gin.Context, chainId string, overviewData entity_cross.OverviewDataView) {
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisCrossOverviewData, prefix, chainId)
	retJson, err := json.Marshal(overviewData)
	if err != nil {
		log.Errorf("【Redis】set cache failed, key:%v, result:%v", redisKey, retJson)
		return
	}
	// 设置键值对和过期时间(40s 过期)
	_ = cache.GlobalRedisDb.Set(ctx, redisKey, string(retJson), 40*time.Second).Err()
}

// CrossLatestTxListHandler modify
type CrossLatestTxListHandler struct {
}

// Handle CrossLatestTxListHandler 最新跨链交易
func (handler *CrossLatestTxListHandler) Handle(ctx *gin.Context) {
	params := entity_cross.BindGetChainIdHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, entity.ErrorMsgParam)
		ConvergeFailureResponse(ctx, newError)
		return
	}

	// crossTxTransferList
	crossTxTransferList, err := dbhandle.GetCrossTxTransferLatestList(params.ChainId)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	txViews := arraylist.New()
	if len(crossTxTransferList) == 0 {
		ConvergeListResponse(ctx, txViews.Values(), 0, nil)
	}

	for _, txTransfer := range crossTxTransferList {
		var status int32
		if txTransfer.Status == int32(tcipCommon.CrossChainStateValue_CONFIRM_END) {
			status = 1
		} else if txTransfer.Status == int32(tcipCommon.CrossChainStateValue_CANCEL_END) {
			status = 2
		}

		latestListView := &entity_cross.LatestTxListView{
			CrossId:         txTransfer.CrossId,
			Status:          status, //跨链状态（0:进行中，1:成功，2:失败）
			Timestamp:       txTransfer.StartTime,
			FromChainId:     txTransfer.FromChainId,
			FromIsMainChain: txTransfer.FromIsMainChain,
			ToChainId:       txTransfer.ToChainId,
			ToIsMainChain:   txTransfer.ToIsMainChain,
			CrossModel:      txTransfer.CrossModel,
			TxNum:           txTransfer.TxNum,
		}
		txViews.Add(latestListView)
	}
	ConvergeListResponse(ctx, txViews.Values(), 10, nil)
}

// CrossLatestSubChainListHandler delete
type CrossLatestSubChainListHandler struct {
}

// Handle CrossLatestSubChainListHandler
func (handler *CrossLatestSubChainListHandler) Handle(ctx *gin.Context) {
	params := entity_cross.BindGetChainIdHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, entity.ErrorMsgParam)
		ConvergeFailureResponse(ctx, newError)
		return
	}

	// crossSubChainList
	crossSubChainList, err := dbhandle.GetCrossLatestSubChainList(params.ChainId)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	sunChainListViews := arraylist.New()
	if len(crossSubChainList) == 0 {
		ConvergeListResponse(ctx, sunChainListViews.Values(), 0, nil)
	}

	for _, subChain := range crossSubChainList {
		//跨链合约数
		crossContractNum, err := dbhandle.GetCrossContractCount(params.ChainId, subChain.SubChainId)
		if err != nil {
			log.Errorf("Get CrossContract Count err : %v", err)
		}
		latestListView := &entity_cross.LatestSubChainListView{
			SubChainId:       subChain.SubChainId,
			BlockHeight:      subChain.BlockHeight,
			IsMainChain:      subChain.IsMainChain,
			Timestamp:        subChain.Timestamp,
			CrossTxNum:       subChain.TxNum,
			CrossContractNum: crossContractNum,
		}
		sunChainListViews.Add(latestListView)
	}

	ConvergeListResponse(ctx, sunChainListViews.Values(), int64(len(crossSubChainList)), nil)
}

// GetCrossTxListHandler 跨链交易列表
type GetCrossTxListHandler struct {
}

// Handle CrossLatestTxListHandler 最新跨链交易
func (handler *GetCrossTxListHandler) Handle(ctx *gin.Context) {
	params := entity_cross.BindGetCrossTxListHandler(ctx)
	if params == nil || !params.IsLegal() || !params.RangeBody.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, entity.ErrorMsgParam)
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//获取交易列表
	crossTransferList, totalCount, err := dbhandle.GetCrossSubChainTransferList(params.Offset, params.Limit,
		params.StartTime, params.EndTime, params.ChainId, params.CrossId, params.SubChainId, params.FromChainId,
		params.ToChainId)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	//构造返回值
	txListViews := buildCrossTxListView(crossTransferList)
	ConvergeListResponse(ctx, txListViews.Values(), totalCount, nil)
}

// 构造交易列表返回值
func buildCrossTxListView(crossTransferList []*db.CrossTransactionTransfer) *arraylist.List {
	txListViews := arraylist.New()
	if len(crossTransferList) == 0 {
		return txListViews
	}

	for _, transfer := range crossTransferList {
		var status int32
		if transfer.Status == int32(tcipCommon.CrossChainStateValue_CONFIRM_END) {
			status = 1
		} else if transfer.Status == int32(tcipCommon.CrossChainStateValue_CANCEL_END) {
			status = 2
		}
		txView := &entity_cross.GetTxListView{
			CrossId:         transfer.CrossId,
			Status:          status, //跨链状态（0:进行中，1:成功，2:失败）
			Timestamp:       transfer.StartTime,
			FromChainId:     transfer.FromChainId,
			FromIsMainChain: transfer.FromIsMainChain,
			ToChainId:       transfer.ToChainId,
			ToIsMainChain:   transfer.ToIsMainChain,
			CrossModel:      transfer.CrossModel,
			TxNum:           transfer.TxNum,
		}
		txListViews.Add(txView)
	}
	return txListViews
}

// CrossSubChainListHandler get
type CrossSubChainListHandler struct {
}

// Handle CrossSubChainListHandler
func (handler *CrossSubChainListHandler) Handle(ctx *gin.Context) {
	params := entity_cross.BindGetCrossSubChainListHandler(ctx)
	if params == nil || !params.IsLegal() || !params.RangeBody.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, entity.ErrorMsgParam)
		ConvergeFailureResponse(ctx, newError)
		return
	}

	crossSubChainList, totalCount, err := dbhandle.GetCrossSubChainList(params.Offset, params.Limit,
		params.ChainId, params.SubChainId)
	if err != nil {
		log.Errorf("GetCrossSubChainList err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	sunChainListViews := arraylist.New()
	if len(crossSubChainList) == 0 {
		ConvergeListResponse(ctx, sunChainListViews.Values(), 0, nil)
		return
	}

	for _, subChain := range crossSubChainList {
		//跨链合约数
		crossContractNum, err := dbhandle.GetCrossContractCount(params.ChainId, subChain.SubChainId)
		if err != nil {
			log.Errorf("Get CrossContract Count err : %v", err)
		}
		chainView := &entity_cross.GetSubChainListView{
			SubChainId:       subChain.SubChainId,
			BlockHeight:      subChain.BlockHeight,
			IsMainChain:      subChain.IsMainChain,
			Timestamp:        subChain.Timestamp,
			CrossTxNum:       subChain.TxNum,
			CrossContractNum: crossContractNum,
		}
		sunChainListViews.Add(chainView)
	}

	ConvergeListResponse(ctx, sunChainListViews.Values(), totalCount, nil)
}

// CrossSubChainDetailHandler get
type CrossSubChainDetailHandler struct {
}

// Handle CrossSubChainDetailHandler 主子链网-子链详情
func (handler *CrossSubChainDetailHandler) Handle(ctx *gin.Context) {
	params := entity_cross.BindGetCrossSubChainDetailHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, entity.ErrorMsgParam)
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//子链信息
	subChainInfo, err := dbhandle.GetCrossSubChainInfoById(params.ChainId, params.SubChainId)
	if err != nil {
		log.Errorf("GetCrossSubChainInfoById err : %v", err)
		ConvergeHandleFailureResponse(ctx, entity.ErrSelectFailed)
		return
	} else if subChainInfo == nil {
		ConvergeHandleFailureResponse(ctx, entity.ErrRecordNotFoundErr)
		return
	}

	//跨链合约数
	crossContractNum, err := dbhandle.GetCrossContractCount(params.ChainId, params.SubChainId)
	if err != nil {
		log.Errorf("Get CrossContract Count err : %v", err)
	}

	subChainView := &entity_cross.GetCrossSubChainDetailView{
		SubChainId:       subChainInfo.SubChainId,
		BlockHeight:      subChainInfo.BlockHeight,
		CrossTxNum:       subChainInfo.TxNum,
		CrossContractNum: crossContractNum,
		GatewayId:        subChainInfo.GatewayId,
		Timestamp:        subChainInfo.Timestamp,
		IsMainChain:      subChainInfo.IsMainChain,
	}

	//返回response
	ConvergeDataResponse(ctx, subChainView, nil)
}

// GetCrossTxDetailHandler get
type GetCrossTxDetailHandler struct {
}

// Handle GetCrossTxDetailHandler 主子链网-跨链交易详情
func (handler *GetCrossTxDetailHandler) Handle(ctx *gin.Context) {
	params := entity_cross.BindGetCrossTxDetailHandler(ctx)
	newError := entity.NewError(entity.ErrorParamWrong, entity.ErrorMsgParam)
	if params == nil || !params.IsLegal() {
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//跨链交易流转
	crossTransfers, err := dbhandle.GetCrossCycleTransferById(params.ChainId, params.CrossId)
	if err != nil || len(crossTransfers) == 0 {
		ConvergeHandleFailureResponse(ctx, entity.ErrSelectFailed)
		return
	}

	transferInfo := crossTransfers[0]
	businessList, err := dbhandle.GetCrossBusinessTxByCross(params.ChainId, params.CrossId)
	var (
		fromGas  = "-"
		toGas    = "-"
		fromTxId string
		toTxId   string
	)
	fromTxInfo := &db.CrossBusinessTransaction{}
	toTxInfo := &db.CrossBusinessTransaction{}
	if err == nil && len(businessList) > 0 {
		for _, tx := range businessList {
			if transferInfo.FromChainId == tx.SubChainId {
				fromTxInfo = tx
				fromGas = strconv.FormatUint(tx.GasUsed, 10)
				fromTxId = fromTxInfo.TxId
			} else if transferInfo.ToChainId == tx.SubChainId {
				toTxInfo = tx
				toGas = strconv.FormatUint(tx.GasUsed, 10)
				toTxId = toTxInfo.TxId
			}
		}
	}

	crossDirection := &entity_cross.CrossDirection{
		FromChain: transferInfo.FromChainId,
		ToChain:   transferInfo.ToChainId,
	}
	fromChainTx := &entity_cross.TxChainInfo{
		ChainId:      transferInfo.FromChainId,
		ContractName: fromTxInfo.ContractName,
		IsMainChain:  transferInfo.FromIsMainChain,
		TxId:         fromTxId,
		TxStatus:     fromTxInfo.TxStatus,
		Gas:          fromGas,
	}

	toChainTx := &entity_cross.TxChainInfo{
		ChainId:      transferInfo.ToChainId,
		ContractName: transferInfo.ContractName,
		IsMainChain:  transferInfo.ToIsMainChain,
		TxId:         toTxId,
		TxStatus:     toTxInfo.TxStatus,
		Gas:          toGas,
	}

	var crossDuration int64
	if transferInfo.EndTime > 0 && transferInfo.StartTime > 0 && transferInfo.Status >= 3 {
		crossDuration = transferInfo.EndTime - transferInfo.StartTime
	}
	txDetailView := &entity_cross.GetCrossTxDetailView{
		CrossId:        transferInfo.CrossId,
		Status:         transferInfo.Status,
		CrossDuration:  crossDuration,
		ContractName:   transferInfo.ContractName,
		ContractMethod: transferInfo.ContractMethod,
		Parameter:      transferInfo.Parameter,
		ContractResult: toTxInfo.CrossContractResult,
		CrossDirection: crossDirection,
		FromChainInfo:  fromChainTx,
		ToChainInfo:    toChainTx,
		Timestamp:      transferInfo.StartTime,
	}

	//返回response
	ConvergeDataResponse(ctx, txDetailView, nil)
}

// SubChainCrossChainListHandler get
type SubChainCrossChainListHandler struct {
}

// Handle CrossSubChainListHandler
func (handler *SubChainCrossChainListHandler) Handle(ctx *gin.Context) {
	params := entity_cross.BindGetSubChainCrossChainListHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, entity.ErrorMsgParam)
		ConvergeFailureResponse(ctx, newError)
		return
	}

	subChainCrossList, err := dbhandle.GetSubChainCrossChainList(params.ChainId, params.SubChainId)
	if err != nil {
		log.Errorf("GetSubChainCrossChainList err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	sunChainCrossListViews := arraylist.New()
	if len(subChainCrossList) == 0 {
		ConvergeListResponse(ctx, sunChainCrossListViews.Values(), 0, nil)
		return
	}

	for _, crossChain := range subChainCrossList {
		crossChainView := &entity_cross.GetSubChainCrossView{
			ChainId: crossChain.CrossChainId,
			TxNum:   crossChain.TxNum,
		}
		sunChainCrossListViews.Add(crossChainView)
	}

	ConvergeListResponse(ctx, sunChainCrossListViews.Values(), int64(len(subChainCrossList)), nil)
}

// GetCrossSyncTxDetailHandler get
type GetCrossSyncTxDetailHandler struct {
}

// Handle GetCrossSyncTxDetailHandler 主子链网-跨链交易详情
func (handler *GetCrossSyncTxDetailHandler) Handle(ctx *gin.Context) {
	params := entity_cross.BindGetCrossSyncTxDetailHandler(ctx)
	newError := entity.NewError(entity.ErrorParamWrong, entity.ErrorMsgParam)
	if params == nil || !params.IsLegal() {
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//跨链交易流转
	crossTransfers, err := dbhandle.GetCrossCycleTransferById(params.ChainId, params.CrossId)
	if err != nil || len(crossTransfers) == 0 {
		ConvergeHandleFailureResponse(ctx, entity.ErrSelectFailed)
		return
	}

	transferInfo := crossTransfers[0]
	crossDirection := &entity_cross.CrossDirection{
		FromChain: transferInfo.FromChainId,
		ToChain:   transferInfo.ToChainId,
	}
	fromChainTx := &entity_cross.FromChainInfo{
		ChainId:     transferInfo.FromChainId,
		IsMainChain: transferInfo.FromIsMainChain,
	}

	toChainTx := &entity_cross.ToChainInfo{
		ChainId:        transferInfo.ToChainId,
		ContractName:   transferInfo.ContractName,
		IsMainChain:    transferInfo.ToIsMainChain,
		ContractMethod: transferInfo.ContractMethod,
	}

	txDetailView := &entity_cross.GetCrossSyncTxDetailView{
		CrossId:        transferInfo.CrossId,
		TxId:           transferInfo.TxId,
		CrossDirection: crossDirection,
		FromChainInfo:  fromChainTx,
		ToChainInfo:    toChainTx,
		Timestamp:      transferInfo.StartTime,
		TxNum:          transferInfo.TxNum,
	}

	//返回response
	ConvergeDataResponse(ctx, txDetailView, nil)
}
