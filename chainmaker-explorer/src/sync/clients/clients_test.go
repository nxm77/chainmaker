/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package client

import (
	"chainmaker_web/src/config"
	"reflect"
	"sync"
	"testing"

	pbconfig "chainmaker.org/chainmaker/pb-go/v2/config"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
)

var ChainId1 = "ChainId1"

func TestNewSingleSdkClientPool(t *testing.T) {
	type args struct {
		chainInfo       *config.ChainInfo
		systemSdkClient *SdkClient
		queryClient     *sdk.ChainClient
	}
	tests := []struct {
		name string
		args args
		want *SingleSdkClientPool
	}{
		{
			name: "Test case 1: Valid chainInfo, systemSdkClient, and queryClient",
			args: args{
				chainInfo:       &config.ChainInfo{}, // TODO: Provide a valid ChainInfo instance
				systemSdkClient: &SdkClient{},        // TODO: Provide a valid SdkClient instance
				queryClient:     &sdk.ChainClient{},  // TODO: Provide a valid ChainClient instance
			},
			want: &SingleSdkClientPool{
				chainInfo:       &config.ChainInfo{}, // TODO: Provide the expected ChainInfo instance
				systemSdkClient: &SdkClient{},        // TODO: Provide the expected SdkClient instance
				queryClient:     &sdk.ChainClient{},  // TODO: Provide the expected ChainClient instance
				sdkClients:      sync.Map{},          // TODO: Provide the expected sync.Map instance with the systemSdkClient added
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSingleSdkClientPool(tt.args.chainInfo, tt.args.systemSdkClient, tt.args.queryClient)
			// TODO: Add any necessary checks to verify the SingleSdkClientPool has been initialized correctly.
			if got == nil {
				t.Errorf("NewSingleSdkClientPool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetAllSdkClient(t *testing.T) {
	type args struct {
		chainList []*config.ChainInfo
	}
	tests := []struct {
		name string
		args args
		want []*SdkClient
	}{
		{
			name: "Test case 1: Valid chainId",
			args: args{
				chainList: []*config.ChainInfo{
					{
						ChainId: ChainId1,
					},
				},
			},
			want: []*SdkClient{}, // TODO: Provide the expected SdkClient instance
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAllSdkClient(tt.args.chainList); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllSdkClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetChainClient1(t *testing.T) {
	type args struct {
		chainId string
	}
	tests := []struct {
		name string
		args args
		want *sdk.ChainClient
	}{
		{
			name: "Test case 1: Valid chainId",
			args: args{
				chainId: ChainId1,
			},
			want: nil, // TODO: Provide the expected SdkClient instance
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetChainClient(tt.args.chainId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetChainClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetSdkClient(t *testing.T) {
	tests := []struct {
		name    string
		chainId string
		want    *SdkClient
	}{
		{
			name:    "Test case 1: Valid chainId",
			chainId: "testchain1",
			want:    &SdkClient{}, // TODO: Provide the expected SdkClient instance
		},
		{
			name:    "Test case 2: Invalid chainId",
			chainId: "nonexistentchain",
			want:    nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = GetSdkClient(tt.chainId)
		})
	}
}

func Test_getDefaultLogger(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "Test case 1: Check if logger is not nil",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getDefaultLogger()
			if (got != nil) != tt.want {
				t.Errorf("getDefaultLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateChainClient(t *testing.T) {
	// Test case 1: Create chain client with public auth type
	chainInfo1 := &config.ChainInfo{
		ChainId:  "chain1",
		AuthType: config.PUBLIC,
		HashType: "SHA256",
		OrgId:    "org1",
		NodesList: []*config.NodeInfo{
			{
				Addr: "127.0.0.1:12301",
			},
		},
		UserInfo: &config.UserInfo{
			UserSignKey: "userSignKey1",
			UserSignCrt: "userSignCrt1",
			UserEncKey:  "userEncKey1",
			UserEncCrt:  "userEncCrt1",
		},
		TlsMode: config.TlsModelSingle,
	}
	_, err1 := CreateChainClient(chainInfo1)
	if err1 == nil {
		t.Errorf("Test case 1 failed: expected error, got nil")
	}

	// Test case 3: Create chain client with invalid auth type
	chainInfo3 := &config.ChainInfo{
		ChainId:   "chain3",
		AuthType:  "invalid",
		HashType:  "SHA256",
		OrgId:     "org3",
		NodesList: []*config.NodeInfo{},
		UserInfo: &config.UserInfo{
			UserSignKey: "userSignKey3",
			UserSignCrt: "userSignCrt3",
			UserEncKey:  "userEncKey3",
			UserEncCrt:  "userEncCrt3",
		},
		TlsMode: config.TlsModelSingle,
	}
	_, err3 := CreateChainClient(chainInfo3)
	if err3 == nil {
		t.Errorf("Test case 2 failed: expected error, got nil")
	}
}

func TestNewSdkClient(t *testing.T) {
	// Test case 4: Create new SDK client
	chainInfo4 := &config.ChainInfo{
		ChainId:   "chain4",
		AuthType:  config.PUBLIC,
		HashType:  "SHA256",
		OrgId:     "org4",
		NodesList: []*config.NodeInfo{},
		UserInfo: &config.UserInfo{
			UserSignKey: "userSignKey4",
			UserSignCrt: "userSignCrt4",
			UserEncKey:  "userEncKey4",
			UserEncCrt:  "userEncCrt4",
		},
		TlsMode: config.TlsModelSingle,
	}
	client4 := &sdk.ChainClient{}
	newClient := NewSdkClient(chainInfo4, client4)
	if newClient == nil {
		t.Errorf("Test case 4 failed: expected non-nil client, got nil")
	}
}

func TestGetSingleSdkClient(t *testing.T) {
	// Test case 10: Get single SDK client by chain ID
	chainInfo10 := &config.ChainInfo{
		ChainId:   "chain10",
		AuthType:  config.PUBLIC,
		HashType:  "SHA256",
		OrgId:     "org10",
		NodesList: []*config.NodeInfo{},
		UserInfo: &config.UserInfo{
			UserSignKey: "userSignKey10",
			UserSignCrt: "userSignCrt10",
			UserEncKey:  "userEncKey10",
			UserEncCrt:  "userEncCrt10",
		},
		TlsMode: config.TlsModelSingle,
	}

	client10 := &sdk.ChainClient{}
	sdkClient10 := NewSdkClient(chainInfo10, client10)
	singleSdkClientPool := NewSingleSdkClientPool(chainInfo10, sdkClient10, client10)
	SdkClientPool.sdkClients.Store(chainInfo10.ChainId, singleSdkClientPool)
	client := GetSingleSdkClient(chainInfo10.ChainId)
	if client == nil {
		t.Errorf("Test case 10 failed")
	}

	// Test case 11: Get single SDK client with invalid chain ID
	client11 := GetSingleSdkClient("invalidChainId")
	if client11 != nil {
		t.Errorf("Test case 11 failed: expected error, got nil")
	}
}

func TestAddSdkClient(t *testing.T) {
	// Test case 14: Add SDK client to pool
	chainInfo14 := &config.ChainInfo{
		ChainId:   "chain14",
		AuthType:  config.PUBLIC,
		HashType:  "SHA256",
		OrgId:     "org14",
		NodesList: []*config.NodeInfo{},
		UserInfo: &config.UserInfo{
			UserSignKey: "userSignKey14",
			UserSignCrt: "userSignCrt14",
			UserEncKey:  "userEncKey14",
			UserEncCrt:  "userEncCrt14",
		},
		TlsMode: config.TlsModelSingle,
	}

	client14 := &sdk.ChainClient{}
	sdkClient14 := NewSdkClient(chainInfo14, client14)
	singleSdkClientPool := NewSingleSdkClientPool(chainInfo14, sdkClient14, client14)
	singleSdkClientPool.AddSdkClient(sdkClient14)
	if _, ok := singleSdkClientPool.sdkClients.Load(chainInfo14.ChainId); !ok {
		t.Errorf("Test case 14 failed: expected SDK client to be added to pool")
	}

	// Test case 15: Add SDK client to pool with invalid chain info
	chainInfo15 := &config.ChainInfo{
		ChainId:   "chain15",
		AuthType:  "invalid",
		HashType:  "SHA256",
		OrgId:     "org15",
		NodesList: []*config.NodeInfo{},
		UserInfo: &config.UserInfo{
			UserSignKey: "userSignKey15",
			UserSignCrt: "userSignCrt15",
			UserEncKey:  "userEncKey15",
			UserEncCrt:  "userEncCrt15",
		},
		TlsMode: config.TlsModelSingle,
	}

	client15 := &sdk.ChainClient{}
	sdkClient15 := NewSdkClient(chainInfo15, client15)
	singleSdkClientPool = NewSingleSdkClientPool(chainInfo15, sdkClient15, client15)
	singleSdkClientPool.AddSdkClient(sdkClient15)
	if _, ok := singleSdkClientPool.sdkClients.Load(chainInfo15.ChainId); !ok {
		t.Errorf("Test case 15 failed: expected SDK client to be added to pool")
	}
}

func TestRemoveSdkClient(t *testing.T) {
	// Test case 16: Remove SDK client from pool
	chainInfo16 := &config.ChainInfo{
		ChainId:   "chain16",
		AuthType:  config.PUBLIC,
		HashType:  "SHA256",
		OrgId:     "org16",
		NodesList: []*config.NodeInfo{},
		UserInfo: &config.UserInfo{
			UserSignKey: "userSignKey16",
			UserSignCrt: "userSignCrt16",
			UserEncKey:  "userEncKey16",
			UserEncCrt:  "userEncCrt16",
		},
		TlsMode: config.TlsModelSingle,
	}

	client16 := &sdk.ChainClient{}
	sdkClient16 := NewSdkClient(chainInfo16, client16)
	singleSdkClientPool := NewSingleSdkClientPool(chainInfo16, sdkClient16, client16)
	singleSdkClientPool.AddSdkClient(sdkClient16)
	SdkClientPool.sdkClients.Store(chainInfo16.ChainId, singleSdkClientPool)
	SdkClientPool.RemoveSdkClient(chainInfo16.ChainId)
	if _, ok := SdkClientPool.sdkClients.Load(chainInfo16.ChainId); ok {
		t.Errorf("Test case 16 failed: expected SDK client to be removed from pool")
	}

	// Test case 17: Remove SDK client from pool with invalid chain ID
	SdkClientPool.RemoveSdkClient("invalidChainId")
}

func TestUpdateChainConfig(t *testing.T) {
	// Test case 18: Update chain config
	chainInfo18 := &config.ChainInfo{
		ChainId:   "chain18",
		AuthType:  config.PUBLIC,
		HashType:  "SHA256",
		OrgId:     "org18",
		NodesList: []*config.NodeInfo{},
		UserInfo: &config.UserInfo{
			UserSignKey: "userSignKey18",
			UserSignCrt: "userSignCrt18",
			UserEncKey:  "userEncKey18",
			UserEncCrt:  "userEncCrt18",
		},
		TlsMode: config.TlsModelSingle,
	}

	client18 := &sdk.ChainClient{}
	sdkClient18 := NewSdkClient(chainInfo18, client18)
	SdkClientPool.sdkClients.Store(chainInfo18.ChainId, sdkClient18)
	newChainConfig18 := &pbconfig.ChainConfig{
		Crypto: &pbconfig.CryptoConfig{
			Hash: "SHA256",
		},
	}
	SdkClientPool.UpdateChainConfig(chainInfo18.ChainId, newChainConfig18)

	// Test case 19: Update chain config with invalid chain ID
	SdkClientPool.UpdateChainConfig("invalidChainId", newChainConfig18)
}
