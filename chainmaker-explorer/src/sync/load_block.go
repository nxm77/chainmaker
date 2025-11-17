/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	client "chainmaker_web/src/sync/clients"
	"chainmaker_web/src/sync/common"
	"chainmaker_web/src/sync/model"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	pbCommon "chainmaker.org/chainmaker/pb-go/v2/common"
)

// BlockWaitUpdate
// @Description: 异步更新channel数据
type BlockWaitUpdate struct {
	ChainId     string
	BlockHeight int64
}

// DataSaveToDB
// @Description: 处理完的订阅数据,放入存储channel,顺序插入数据库
type DataSaveToDB struct {
	ChainId     string
	BlockHeight int64
	DealResult  *model.ProcessedBlockData
}

// PeriodicGetSubscribeLock
//
//	@Description:  10分钟请求一次,获取订阅锁,获取到锁才进行订阅
//	@param sdkClient
func PeriodicGetSubscribeLock(sdkClient *client.SdkClient) {
	ctx := sdkClient.Ctx
	chainId := sdkClient.ChainId
	prefix := config.GlobalConfig.RedisDB.Prefix
	lockKey := fmt.Sprintf(cache.RedisSubscribeLockKey, prefix, chainId)
	// 尝试获取分布式锁（第一次尝试）
	lock := cache.GlobalRedisDb.SetNX(ctx, lockKey, chainId, 3*time.Minute)
	if lock.Val() {
		log.Infof("【load】Periodic Get Subscribe Lock (first attempt)【true】, LockKey:%s", lockKey)
		// 获取到锁,说明其他节点订阅失败了,启动订阅
		err := blockListen(sdkClient)
		if err != nil {
			//重启链
			log.Errorf("【load】ReStartChain PeriodicGetSubscribeLock chainId:%s, err:%s", chainId, err.Error)
			ReStartChain(sdkClient.ChainId)
		}
	} else {
		log.Infof("【load】Periodic Get Subscribe Lock (first attempt)【false】, LockKey:%s", lockKey)
		// 如果第一次尝试失败,启动定时器
		startSubscribeLockTicker(ctx, sdkClient, lockKey)
	}
}

// startSubscribeLockTicker
//
//	@Description: 当获取锁失败,有其他机器订阅时,启动订阅定时器,知道订阅成功,停止定时器
//	@param sdkClient
//	@param LockKey
func startSubscribeLockTicker(ctx context.Context, sdkClient *client.SdkClient, lockKey string) {
	//1分钟定时器,获取订阅锁
	ticker := time.NewTicker(time.Second * 60)
	defer ticker.Stop() // 在函数返回时停止定时器
	chainId := sdkClient.ChainId

	for {
		select {
		case <-ticker.C:
			//链订阅已经停止,停止定时器
			if sdkClient.Status == client.STOP {
				return
			}

			// 尝试获取分布式锁
			lock := cache.GlobalRedisDb.SetNX(ctx, lockKey, chainId, 3*time.Minute)
			log.Infof("【load】Periodic Get Subscribe Lock【%v】, LockKey:%s", lock.Val(), lockKey)
			if lock.Val() {
				// 获取到锁,说明其他节点订阅失败了,启动订阅
				// 订阅链
				err := blockListen(sdkClient)
				if err != nil {
					//重启链
					log.Errorf("【load】ReStartChain startSubscribeLockTicker chainId:%s, err:%s", chainId, err.Error)
					ReStartChain(chainId)
					return
				}
			}
		case <-ctx.Done():
			// 当接收到上下文取消的通知时,返回函数,停止定时器
			return
		}
	}
}

// blockListen
//
//	@Description: 订阅区块数据
//	@param sdkClient 链连接
//	@return error
func blockListen(sdkClient *client.SdkClient) error {
	log.Infof("【load】 begin to subscribe [chain:%s] ", sdkClient.ChainId)
	//判断连接池还在不在,不在的话说明订阅被取消了，不在重启链
	poolSdkClient := client.GetSdkClient(sdkClient.ChainId)
	if poolSdkClient == nil {
		log.Infof("【ReStartChain】poolSdkClient is null, chain is cancel,chainId:%v", sdkClient.ChainId)
		return nil
	}

	var (
		//err error
		err error
		//chainId chainId
		chainId = sdkClient.ChainId
		// 使用sdkClient的Context
		ctx = sdkClient.Ctx
		//blockInfoCh 创建一个缓冲通道来存储格式化后的订阅数据 blockInfo
		blockInfoCh = make(chan *pbCommon.BlockInfo, config.BlockInsertWorkerCount)
		//blockWaitUpdateCh 创建一个容量为 20 的通道来存储处理完成,等待异步更新的区块数据
		blockWaitUpdateCh = make(chan *BlockWaitUpdate, config.BlockWaitUpdateWorkerCount)
		//dataSaveCh 创建一个无缓冲的通道,用来顺序执行插入操作 ,因为需要确保区块是顺序插入的,所以采用无缓存通道。
		dataSaveCh = make(chan *DataSaveToDB)
		//blockListenErrCh 创建一个错误通道来接收 子线程 的错误
		blockListenErrCh = make(chan error)
	)

	defer func() {
		log.Infof("【load】 end to blockListen [chain:%s] ", chainId)
		sdkClient.Cancel()
		close(blockInfoCh)
		close(blockWaitUpdateCh)
		close(dataSaveCh)
		close(blockListenErrCh)
	}()

	blockHeightList := GetUpdateBlockHeightListCache(chainId)
	if errUp := BatchDelayedUpdate(chainId, blockHeightList); errUp != nil {
		log.Errorf("【load】BatchDelayedUpdate err:%v", errUp)
		return errUp
	}

	//消费blockWaitUpdateCh队列,
	//启动异步更新协成,计算交易数量,持仓信息等区块数据
	go DelayUpdateOperation(ctx, chainId, blockWaitUpdateCh, blockListenErrCh)

	//写入blockWaitUpdateCh队列,
	//将上一次未更新数据写入通道,继续更新操作
	err = waitUpdateChFailedData(ctx, chainId, blockWaitUpdateCh)
	if err != nil {
		log.Errorf("【load】waitUpdateChFailedData err:%v", err)
		return err
	}

	//消费dataSaveCh队列,写入blockWaitUpdateCh队列
	//按顺序插入DB,区块数据,需要确保区块是顺序插入的,写入异步处理队列blockWaitUpdateCh
	go OrderedSaveBlockData(ctx, dataSaveCh, blockWaitUpdateCh, blockListenErrCh)

	//写入blockInfoCh队列
	//订阅数据,将订阅数据处理成结构化数据
	go SubscribeBlockSetToBlockInfoCh(ctx, sdkClient, blockInfoCh, blockListenErrCh)

	hash := sdkClient.GetChainHashType()
	//消费blockInfoCh队列,写入dataSaveCh队列
	//将格式话的区块数据,处理成可以插入DB的数据格式
	go RealtimeInsertOperation(ctx, hash, blockInfoCh, dataSaveCh, blockListenErrCh)

	// 使用 select 语句等待错误或上下文取消
	select {
	case errCh := <-blockListenErrCh:
		// 接收区块错误,无法处理成结构化数据,重启链
		log.Errorf("【sync block】Subscribe block failed, err:%v", errCh)
		return errCh
	case <-ctx.Done():
		log.Errorf("【sync block】Subscribe block failed, context cancel")
		// 上下文已取消,停止监听
		return ctx.Err()
	}
}

// OrderedSaveBlockData
//
//	@Description:  将订阅处理完的数据,按顺序插入DB,采用channel确保区块是顺序插入的
//	@param ctx
//	@param dataSaveCh 按照blockHeight顺序写入dataSaveCh,保证插入顺序
//	@param blockWaitUpdateCh 插入完成后写入blockWaitUpdateCh通道,等待更新数据
//	@param errCh
func OrderedSaveBlockData(ctx context.Context, dataSaveCh chan *DataSaveToDB, blockWaitUpdateCh chan *BlockWaitUpdate,
	errCh chan<- error) {
	for data := range dataSaveCh {
		chainId := data.ChainId
		blockHeight := data.BlockHeight
		log.Infof("【Realtime insert】start block-%s[%d]", chainId, blockHeight)
		startTime := time.Now()
		err := RealtimeDataSaveToDB(chainId, blockHeight, data.DealResult)
		log.Infof("【Realtime insert】end block-%s[%d], duration_time(ms):%v", chainId, blockHeight,
			time.Since(startTime).Milliseconds())
		if err != nil {
			errCh <- fmt.Errorf("【Realtime insert】err block-%s[%d] failed, err:%v", chainId, blockHeight, err)
			// 如果处理失败,取消上下文,会重启链
			return
		}

		// 将处理完成的结果写入 blockWaitUpdateCh
		resultData := &BlockWaitUpdate{
			ChainId:     chainId,
			BlockHeight: blockHeight,
		}
		select {
		case <-ctx.Done():
			// 上下文已取消,不要发送数据到通道
			return
		case blockWaitUpdateCh <- resultData:
		}
	}
}

// waitUpdateChFailedData
//
//	@Description: 程序异常重启后,先从数据库获取未更新的数据
//	@param ctx
//	@param chainId
//	@param blockWaitUpdateCh
//	@return error
func waitUpdateChFailedData(ctx context.Context, chainId string, blockWaitUpdateCh chan *BlockWaitUpdate) error {
	blockList, err := dbhandle.GetBlockByStatus(chainId, dbhandle.DelayUpdateFail)
	if err != nil {
		return err
	}

	//将未更新数据写入异步更新队列blockWaitUpdateCh
	for _, block := range blockList {
		log.Infof("【load】waitUpdateChFailedData, chainId:%s, blockHeight:%d", chainId, block.BlockHeight)
		// 将处理完成的结果写入 blockWaitUpdateCh
		resultData := &BlockWaitUpdate{
			ChainId:     chainId,
			BlockHeight: block.BlockHeight,
		}
		select {
		case <-ctx.Done():
			// 上下文已取消,不要发送数据到通道
			return nil
		case blockWaitUpdateCh <- resultData:
			//blockList的长度可能是blockWaitUpdateCh长度的2倍
		}
	}

	log.Infof("【load】waitUpdateChFailedData success, chainId:%s, blockList:%v", chainId, blockList)
	return nil
}

// RealtimeInsertOperation
//
//	@Description: BlockInsertWorkerCount个线程并发处理解析区块数据,存储到dataSaveCh通道,等待入库
//	@param ctx
//	@param hash hash值
//	@param blockInfoCh 订阅区块通道
//	@param dataSaveCh 保存区块数据通道
//	@param errCh
func RealtimeInsertOperation(ctx context.Context, hash string, blockInfoCh chan *pbCommon.BlockInfo,
	dataSaveCh chan *DataSaveToDB, errCh chan<- error) {
	workerCount := config.BlockInsertWorkerCount
	// 使用 sync.WaitGroup 来等待所有 worker 协程完成
	var wg sync.WaitGroup
	wg.Add(workerCount)
	// 启动 worker 协程,订阅blockInfoCh队列,并发解析区块数据,写入dataSaveCh通道
	for i := 0; i < workerCount; i++ {
		go ParallelParseBlockWork(ctx, &wg, hash, blockInfoCh, dataSaveCh, errCh)
	}
	wg.Wait()
}

// ParallelParseBlockWork
//
//	@Description: 消费blockInfoCh通道数据,解析成格式化的DB数据,存储到dataSaveCh,等待存储DB
//	@param ctx
//	@param wg
//	@param hashType
//	@param blockInfoCh 订阅区块通道
//	@param dataSaveCh 保存区块数据通道
//	@param errCh
func ParallelParseBlockWork(ctx context.Context, wg *sync.WaitGroup, hashType string,
	blockInfoCh chan *pbCommon.BlockInfo, dataSaveCh chan *DataSaveToDB, errCh chan<- error) {
	defer wg.Done()
	//blockInfoCh 阻塞持续等待blockInfoCh
	for blockInfo := range blockInfoCh {
		if blockInfo == nil {
			log.Errorf("blockInfoCh blockInfo failed.\n")
			continue
		}

		chainId := blockInfo.Block.Header.ChainId
		//nolint:gosec
		blockHeight := int64(blockInfo.Block.Header.BlockHeight)
		startTime := time.Now()
		// 处理区块数据
		log.Infof("【Realtime deal】start block-%s[%d]", chainId, blockHeight)

		// 创建 ProcessedBlockInfo 实例
		processedBlock := &ProcessedBlockInfo{
			BlockInfo: blockInfo,
			HashType:  hashType,
		}
		// 调用 ProcessedBlockHandle 方法处理区块数据
		dealResult, err := processedBlock.ProcessedBlockHandle()
		//dealResult, txTimeLog, err := RealtimeDataHandle(blockInfo, hashType)
		log.Infof("【Realtime deal】end block-%s[%d] duration_time(ms):%v",
			chainId, blockHeight, time.Since(startTime).Milliseconds())
		if err != nil {
			errCh <- fmt.Errorf("【Realtime deal】err block-%s[%d] failed, err:%v", chainId, blockHeight, err)
			// 如果处理失败,取消上下文,会重启链
			return
		}

		dataToSave := &DataSaveToDB{
			ChainId:     chainId,
			BlockHeight: blockHeight,
			DealResult:  dealResult,
		}

		safeSend(ctx, chainId, blockHeight, dataSaveCh, dataToSave)
	}
}

// safeSend 函数用于安全地向通道发送数据，防止通道关闭导致的 panic。
func safeSend(ctx context.Context, chainId string, blockHeight int64, dataSaveCh chan *DataSaveToDB,
	dataToSave *DataSaveToDB) {
	// 捕获 panic 异常
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("send panic, channel closed")
		}
	}()

	// 等待当前区块高度等于最大高度时,将数据写入 dataCh,否则需要等待一下
	done := false
	for !done {
		select {
		case <-ctx.Done():
			// 上下文已取消,跳出循环
			return
		default:
			// 如果当前区块高度等于最大高度,则将数据写入 dataCh 通道
			if blockHeight == common.GetMaxHeight(chainId) {
				dataSaveCh <- dataToSave
				//只有写入dataSaveCh队列后才会加1
				common.SetMaxHeight(chainId, blockHeight+1)
				done = true
			} else {
				//不sleep的话,高并发会占满cpu
				time.Sleep(20 * time.Millisecond)
			}
		}
	}
}

// DelayUpdateOperation
//
//	@Description: 异步数据计算更新,订阅blockWaitUpdateCh通道数据
//	@param blockWaitUpdateCh 等到异步计算的通道
//	@param errCh
func DelayUpdateOperation(ctx context.Context, chainId string, blockWaitUpdateCh chan *BlockWaitUpdate,
	errCh chan<- error) {
	for {
		//var blockWaitUpdates []*BlockWaitUpdate
		var blockHeightList []int64
		maxBatchSize := config.BlockUpdateWorkerCount

		// 从通道中读取数据，直到达到最大批量大小或通道为空
		exitLoop := false
		for len(blockHeightList) < maxBatchSize && !exitLoop {
			select {
			case blockInfo, ok := <-blockWaitUpdateCh:
				if !ok {
					// 通道关闭，退出操作
					log.Errorf("blockWaitUpdateCh 通道已关闭，终止操作。")
					return
				}
				blockHeightList = append(blockHeightList, blockInfo.BlockHeight)
			case <-ctx.Done():
				// 上下文取消，退出操作
				log.Infof("操作被取消，退出函数。")
				return
			default:
				if len(blockHeightList) > 0 && len(blockWaitUpdateCh) == 0 {
					// 设置标志位退出外层循环
					exitLoop = true
					break
				}
				time.Sleep(time.Millisecond * 100) // 适当延时，防止过度占用CPU
			}
		}

		if len(blockHeightList) == 0 {
			time.Sleep(time.Second)
			continue
		}

		// 执行批量更新操作
		SetUpdateBlockHeightListCache(chainId, blockHeightList)
		if err := BatchDelayedUpdate(chainId, blockHeightList); err != nil {
			log.Errorf("【Delay update】BatchDelayedUpdate failed, err: %v", err)
			errCh <- err
			return
		}
	}
}

func GetUpdateBlockHeightListCache(chainId string) []int64 {
	var blockHeightList []int64
	ctx := context.Background()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisUpdateBlockHeightList, prefix, chainId)
	redisRes := cache.GlobalRedisDb.Get(ctx, redisKey)
	if redisRes != nil && redisRes.Val() != "" {
		err := json.Unmarshal([]byte(redisRes.Val()), &blockHeightList)
		if err == nil {
			return blockHeightList
		}
	}

	return nil
}

func SetUpdateBlockHeightListCache(chainId string, blockHeightList []int64) {
	if len(blockHeightList) == 0 {
		return
	}

	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisUpdateBlockHeightList, prefix, chainId)
	retJson, err := json.Marshal(blockHeightList)
	if err == nil {
		// 设置键值对和过期时间
		ctx := context.Background()
		_ = cache.GlobalRedisDb.Set(ctx, redisKey, string(retJson), 24*time.Hour).Err()
	}
}

// SubscribeBlockSetToBlockInfoCh
//
//	@Description: 将订阅的区块数据Json解析成common.BlockInfo结构,写入blockInfoCh通道
//	@param ctx
//	@param sdkClient 链连接
//	@param blockInfoCh 区块处理通道
//	@param errCh
func SubscribeBlockSetToBlockInfoCh(ctx context.Context, sdkClient *client.SdkClient,
	blockInfoCh chan *pbCommon.BlockInfo, errCh chan<- error) {
	log.Info("【SubscribeBlockSetToBlockInfoCh】 start")
	chainId := sdkClient.ChainId
	chainClient := sdkClient.ChainClient
	isTest := false

	var maxBlockHeight int64
	if isTest {
		maxBlockHeight = int64(193200)
	} else {
		// get max block height for this chain
		maxBlockHeight = dbhandle.GetMaxBlockHeight(chainId)
	}

	log.Infof("【sync load】 begin to subscribe block-%s[%d] ", chainId, maxBlockHeight)
	if maxBlockHeight > 0 {
		common.SetMaxHeight(sdkClient.ChainId, maxBlockHeight+1)
	} else {
		common.SetMaxHeight(sdkClient.ChainId, 0)
	}

	//订阅区块
	c, err := chainClient.SubscribeBlock(ctx, maxBlockHeight, -1, true, false)
	if err != nil {
		common.SubscribeFail.WithLabelValues(chainId).Inc()
		errCh <- fmt.Errorf("【Sync Block】 Get Block By SDK failed:, err: %v", err)
		return
	}

	// 创建一个定时器用于刷新锁的过期时间
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for {
		//链订阅已经停止,停止定时器
		if sdkClient.Status == client.STOP {
			return
		}

		select {
		case block, ok := <-c:
			if !ok {
				common.SubscribeFail.WithLabelValues(chainId).Inc()
				errCh <- fmt.Errorf("【Sync Block】 SubscribeBlock- Chan Is Closed, chainId:%v, ok:%v",
					chainId, ok)
				return
			}

			blockInfo, ok := block.(*pbCommon.BlockInfo)
			if !ok {
				common.SubscribeFail.WithLabelValues(chainId).Inc()
				errCh <- fmt.Errorf("【Sync Block】 SubscribeBlock- The Data Type Error, chainId:%v", chainId)
				return
			}

			//根据区块高度获取区块信息
			//nolint:gosec
			height := int64(blockInfo.Block.Header.BlockHeight)
			blockDB, _ := dbhandle.GetBlockByHeight(chainId, height)
			if blockDB != nil && blockDB.BlockHash != "" && !isTest {
				//数据库已经存在
				common.SetMaxHeight(chainId, height+1)
				log.Infof("【Sync Block】block is existed, chainId:%v, block height:%v \n", chainId, height)
			} else {
				select {
				case <-ctx.Done():
					// 上下文已取消,不要发送数据到通道
					return
				case blockInfoCh <- blockInfo:
					// 成功发送数据到通道
				}
			}
		case <-ticker.C:
			// 定期刷新分布式锁的过期时间
			prefix := config.GlobalConfig.RedisDB.Prefix
			lockKey := fmt.Sprintf(cache.RedisSubscribeLockKey, prefix, chainId)
			cache.GlobalRedisDb.Expire(ctx, lockKey, 3*time.Minute)
			log.Infof("【Sync Block】Redis Set Lock, LockKey:%s", lockKey)
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

// BuildChainInfo
//
//	@Description: 构造链信息
//	@param subscribeChain 链订阅信息
//	@return *config.ChainInfo 链配置结构
func BuildChainInfo(subscribeChain *db.Subscribe) *config.ChainInfo {
	chainInfo := &config.ChainInfo{
		ChainId:  subscribeChain.ChainId,
		AuthType: subscribeChain.AuthType,
		OrgId:    subscribeChain.OrgId,
		HashType: subscribeChain.HashType,
		Tls:      subscribeChain.Tls,
		TlsMode:  subscribeChain.TlsMode,
		UserInfo: &config.UserInfo{
			UserSignKey: subscribeChain.UserSignKey,
			UserSignCrt: subscribeChain.UserSignCrt,
			UserTlsKey:  subscribeChain.UserTlsKey,
			UserTlsCrt:  subscribeChain.UserTlsCrt,
			UserEncKey:  subscribeChain.UserEncKey,
			UserEncCrt:  subscribeChain.UserEncCrt,
		},
	}
	var nodeList []*config.NodeInfo
	_ = json.Unmarshal([]byte(subscribeChain.NodeList), &nodeList)
	chainInfo.NodesList = nodeList
	return chainInfo
}
