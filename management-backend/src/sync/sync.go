/*
Package sync comment
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"management_backend/src/entity"
	loggers "management_backend/src/logger"
)

var (
	sdkClientPool *SdkClientPool
)

// SubscribeChain subscribe chain
func SubscribeChain(sdkConfig *entity.SdkConfig) error {
	var err error
	sdkClient, err := NewSdkClient(sdkConfig)
	loggers.WebLogger.Infof("[SDK] first the current chain: %v", sdkClient.ChainId)
	if err != nil {
		loggers.WebLogger.Error("create sdkClient failed: %v", err.Error())
		return err
	}
	_, err = sdkClient.ChainClient.GetChainConfig()
	if err != nil {
		loggers.WebLogger.Errorf("订阅链失败: %v", err.Error())
		return err
	}
	sdkClientPool = GetSdkClientPool()
	if sdkClientPool != nil {
		err = sdkClientPool.AddSdkClient(sdkClient)
		if err != nil {
			loggers.WebLogger.Error("[WEB] AddSdkClient err : ", err.Error())
			return err
		}
	} else {
		sdkClientPool = NewSdkClientPool(sdkClient)
	}

	LoadChainRefInfos(sdkClient)
	sdkClientPool.LoadChains(sdkConfig.ChainId)
	return nil
}

// GetSdkClientPool get sdk client pool
func GetSdkClientPool() *SdkClientPool {
	if sdkClientPool == nil {
		sdkClients := make(map[string]*SdkClient)
		sdkClientPool = &SdkClientPool{
			SdkClients: sdkClients,
		}
	}
	return sdkClientPool
}

// StopSubscribe stop subscribe
func StopSubscribe(chainId string) {
	sdkClientPool = GetSdkClientPool()
	if _, ok := sdkClientPool.SdkClients[chainId]; ok {
		sdkClientPool.SdkClients[chainId].LoadInfoStop <- struct{}{}
		sdkClientPool.SdkClients[chainId].SubscribeStop <- struct{}{}
		sdkClientPool.RemoveSdkClient(chainId)
	}
}
