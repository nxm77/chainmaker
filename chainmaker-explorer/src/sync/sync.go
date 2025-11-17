/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package sync

// nolint
import (
	"chainmaker_web/src/chain"
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	loggers "chainmaker_web/src/logger"
	client "chainmaker_web/src/sync/clients"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	pbconfig "chainmaker.org/chainmaker/pb-go/v2/config"
)

var (
	log              = loggers.GetLogger(loggers.MODULE_SYNC)
	StartSyncCancels sync.Map // 订阅锁，一个链只能同时订阅一次
)

// Start 开始链订阅任务，如果已经存在订阅任务，则取消之前的订阅任务
// @param chainList 链列表
// @return error 错误信息
func StartSync(chainList []*config.ChainInfo) {
	for _, chainInfo := range chainList {
		// 订阅链
		go func(chainInfo *config.ChainInfo) {
			// 检查是否已经存在一个 cancel 函数
			if cancel, ok := StartSyncCancels.Load(chainInfo.ChainId); ok {
				// 如果存在，调用 cancel 函数来取消上下文
				if cancelFunc, ok := cancel.(context.CancelFunc); ok {
					cancelFunc()
				} else {
					// 添加错误日志或处理逻辑
					log.Errorf("无效的 cancel 类型: %T", cancel)
				}
			}

			// 创建新的 context 和取消函数
			ctx, cancel := context.WithCancel(context.Background())
			StartSyncCancels.Store(chainInfo.ChainId, cancel)

			for {
				select {
				case <-ctx.Done():
					log.Infof("【sync】区块链【%v】订阅任务已取消", chainInfo.ChainId)
					return
				default:
					log.Errorf("【sync】区块链【%v】重新订阅中...", chainInfo.ChainId)
					err := SubscribeToChain(chainInfo)
					if err == nil {
						log.Infof("【sync】区块链【%v】订阅成功", chainInfo.ChainId)
						return
					}
					log.Errorf("【sync】区块链【%v】订阅失败, 正在尝试重新订阅, err:%v", chainInfo.ChainId, err)
					time.Sleep(time.Second * 10) // 添加一个短暂延迟，避免频繁重试
				}
			}
		}(chainInfo)
	}
}

// SubscribeToChain 订阅链
// @param chainInfo 链信息
// @return error 错误信息
func SubscribeToChain(chainInfo *config.ChainInfo) error {
	chainId := chainInfo.ChainId
	// 获取链客户端
	poolSdkClient := client.GetSdkClient(chainId)
	if poolSdkClient != nil {
		// 停止之前的订阅
		client.StopChainClient(chainId)
	}

	//订阅状态
	subscribeStatus := db.SubscribeOK
	// 创建区块链连接，并加入连接池
	chainConfig, clientErr := CreateSubscribeClientPool(chainInfo)
	if clientErr != nil {
		log.Errorf("【sync】区块链【%v】连接失败: %v, 尝试重新订阅...", chainId, clientErr)
		subscribeStatus = db.SubscribeFailed
	}

	// 开启订阅
	syncErr := BeginSubscribeChain(chainId)
	if syncErr != nil {
		subscribeStatus = db.SubscribeFailed
	}

	// 更新订阅状态
	errDB := PersistChainSubscriptionInfo(chainInfo, chainConfig, subscribeStatus)
	if errDB != nil {
		log.Errorf("【sync】 SaveSubscribeToDB failed, err:%v", errDB)
	}

	if clientErr != nil {
		return clientErr
	} else if syncErr != nil {
		return syncErr
	}
	return nil
}

// BeginSubscribeChain 开始订阅
// @param chainId 链ID
// @return error 错误信息
func BeginSubscribeChain(chainId string) error {
	log.Infof("[WEB] begin to load chain's information, [chain:%s] ", chainId)
	// 获取链客户端
	sdkClient := client.GetSdkClient(chainId)
	if sdkClient == nil {
		return fmt.Errorf("sdkClient is nil")
	}

	//处理节点，组织数据
	err := loadChainRefInfos(sdkClient)
	if err != nil {
		return err
	}

	//检查链统计信息
	checkChainStatistics(chainId)

	//订阅区块数据
	go PeriodicGetSubscribeLock(sdkClient)

	//定时处理
	//定期处理节点数据
	go PeriodicLoadStart(sdkClient)
	//定期检查子链健康状态
	//go PeriodicCheckSubChainStatus(sdkClient)
	return nil
}

// checkChainStatistics 检查链统计信息，如果不存在则创建
func checkChainStatistics(chainId string) {
	// 获取链统计信息
	statistics, err := dbhandle.GetChainStatistics(chainId)
	if err != nil {
		log.Errorf("checkChainStatistics Get chain statistics err : %s", err.Error())
	}

	if statistics == nil {
		//新增数据
		statistics = &db.Statistics{
			ChainId:           chainId,
			BlockHeight:       0,
			TotalTransactions: 0,
			TotalCrossTx:      0,
		}
		err := dbhandle.InsertStatistics(chainId, statistics)
		if err != nil {
			log.Errorf("checkChainStatistics InsertStatistics err : %s", err.Error())
		}
	}
}

// subscribe 在这里执行订阅操作，如果订阅成功，返回nil，否则返回错误
// @param chainInfo
// @return error
// subscribe函数用于订阅链信息，并返回链配置和错误信息
func CreateSubscribeClientPool(chainInfo *config.ChainInfo) (*pbconfig.ChainConfig, error) {
	// 将链信息转换为JSON格式
	chainInfoJson, _ := json.Marshal(chainInfo)
	// 打印日志，记录链ID和链信息
	log.Infof("【Sync】 chainId[%v] init sdk clients Start, chainInfoJson:%v", chainInfo.ChainId, string(chainInfoJson))

	// 创建链客户端
	chainClient, err := client.CreateChainClient(chainInfo)
	// 创建SDK客户端
	sdkClient := client.NewSdkClient(chainInfo, chainClient)
	if err != nil {
		// 打印错误日志，记录错误信息和链信息
		log.Errorf("【Sync】创建chain Client失败: err:%v, chainInfo:%v",
			err.Error(), string(chainInfoJson))
		return nil, err
	}

	//判断节点是否存活
	chainConfig, err := sdkClient.ChainClient.GetChainConfig()
	sdkClient.ChainConfig = chainConfig
	if err != nil {
		// 停止链客户端
		_ = chainClient.Stop()
		// 打印错误日志，记录错误信息和链信息
		log.Errorf("【Sync】try to connect chain failed:err:%v , chainInfo:%v",
			err.Error(), string(chainInfoJson))
		return nil, err
	}

	// 加入到连接池
	clientPool := client.NewSingleSdkClientPool(chainInfo, sdkClient, chainClient)
	client.SdkClientPool.AddSdkClientPool(clientPool)

	log.Infof("【Sync】 chainId[%v] init sdk clients success", chainInfo.ChainId)
	return chainConfig, nil
}

// PersistChainSubscriptionInfo
//
//	@Description: 将订阅信息存储数据库
//	@param chainInfo
//	@param chainConfig
//	@return error
func PersistChainSubscriptionInfo(chainInfo *config.ChainInfo, chainConfig *pbconfig.ChainConfig,
	subscribeStatus int) error {
	//更新订阅状态
	err := dbhandle.InsertOrUpdateSubscribe(chainInfo, subscribeStatus)
	log.Infof("【sync】 save Subscribe finished, err:%v, chain:%v, subscribeStatus:%v",
		err, chainInfo, subscribeStatus)
	if err != nil {
		return err
	}

	//处理链数据
	chain := &db.Chain{}
	if chainConfig == nil {
		// 如果链配置为空，则只保存链ID和认证类型
		chain.ChainId = chainInfo.ChainId
		chain.AuthType = chainInfo.AuthType
		chain.Timestamp = time.Now().Unix()
	} else {
		// 如果链配置不为空，则将链配置序列化为JSON格式
		chainConfigBytes, _ := json.Marshal(chainConfig)

		// 默认值处理，避免 nil 指针
		var enableGas bool
		var blockInterval, blockSize, blockTxCapacity, txTimeout int
		var txTimestampVerify bool
		var consensusType, hashType string

		if chainConfig.AccountConfig != nil {
			enableGas = chainConfig.AccountConfig.EnableGas
		}
		if chainConfig.Block != nil {
			blockInterval = int(chainConfig.Block.BlockInterval)
			blockSize = int(chainConfig.Block.BlockSize)
			blockTxCapacity = int(chainConfig.Block.BlockTxCapacity)
			txTimestampVerify = chainConfig.Block.TxTimestampVerify
			txTimeout = int(chainConfig.Block.TxTimeout)
		}
		if chainConfig.Consensus != nil {
			consensusType = chainConfig.Consensus.Type.String()
		}
		if chainConfig.Crypto != nil {
			hashType = chainConfig.Crypto.Hash
		}

		chain = &db.Chain{
			ChainId:           chainInfo.ChainId,
			Version:           chainConfig.Version,
			EnableGas:         enableGas,
			BlockInterval:     blockInterval,
			BlockSize:         blockSize,
			BlockTxCapacity:   blockTxCapacity,
			TxTimestampVerify: txTimestampVerify,
			TxTimeout:         txTimeout,
			Consensus:         consensusType,
			HashType:          hashType,
			AuthType:          chainInfo.AuthType,
			ChainConfig:       string(chainConfigBytes),
			Timestamp:         time.Now().Unix(),
		}
	}

	//插入，更新链信息
	err = dbhandle.InsertUpdateChainInfo(chain, subscribeStatus)
	log.Infof("【sync】 save chaininfo finished, err:%v, chain:%v", err, chain)
	if err != nil {
		return err
	}

	return nil
}

// ReStartChain 重启订阅
// @param chainId 链ID
func ReStartChain(chainId string) {
	//判断连接池还在不在，不在的话不在重启
	poolSdkClient := client.GetSdkClient(chainId)
	if poolSdkClient == nil {
		log.Infof("【ReStartChain】poolSdkClient is null, chain is cancel，chainId:%v", chainId)
		return
	}

	//重启这条链的订阅
	log.Infof("【ReStartChain】restart chain, chainId:%v", chainId)
	chainConfig, err := chain.GetSubscribeByChainId(chainId)
	if err == nil && chainConfig != nil {
		time.Sleep(time.Second * 10)
		StartSync([]*config.ChainInfo{chainConfig})
	}
}

// StopChainClient 停止链客户端
// @param chainId 链ID
func SyncStopChainClient(chainId string) {
	// 检查是否已经存在一个 cancel 函数
	if cancel, ok := StartSyncCancels.Load(chainId); ok {
		// 如果存在，调用 cancel 函数来取消上下文
		if cancelFunc, ok := cancel.(context.CancelFunc); ok {
			cancelFunc()
		} else {
			// 添加错误日志或处理逻辑
			log.Errorf("无效的 cancel 类型: %T", cancel)
		}
	}

	// 停止链客户端
	log.Infof("【StopChainClient】stop chain, chainId:%v", chainId)
	client.StopChainClient(chainId)
}
