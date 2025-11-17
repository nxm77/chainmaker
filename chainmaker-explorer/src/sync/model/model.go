/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package model

import (
	"chainmaker_web/src/db"

	pbConfig "chainmaker.org/chainmaker/pb-go/v2/config"
)

// ProcessedBlockData 表示经过区块链同步处理后的区块数据
type ProcessedBlockData struct {
	BlockDetail          *db.Block                        // 区块详细信息
	UserList             map[string]*db.User              // 用户列表
	Transactions         map[string]*db.Transaction       // 交易数据
	UpgradeContractTx    []*db.UpgradeContractTransaction // 合约升级交易
	ChainConfigList      []*pbConfig.ChainConfig          // 链配置列表
	ContractWriteSetData map[string]*ContractWriteSetData // 合约写集数据
	InsertContracts      []*db.Contract                   // 新增合约
	UpdateContracts      []*db.Contract                   // 更新合约
	FungibleContract     []*db.FungibleContract           // 同质化合约
	NonFungibleContract  []*db.NonFungibleContract        // 非同质化合约
	InsertIDAContracts   []*db.IDAContract                // 新增 IDA 合约
	ContractByteCodeList []*db.ContractByteCode           // 合约字节码列表
	EvidenceList         []*db.EvidenceContract           // 存证合约列表
	ContractEvents       []*db.ContractEvent              // 合约事件
	GasRecordList        []*db.GasRecord                  // gas 消耗记录
	CrossChainResult     *CrossChainResult                // 跨链相关数据
	Timestamp            int64
	ContractCrossTxs     []*db.ContractCrossCallTransaction // 合约跨合约调用交易
	ContractCrossCalls   map[string]*db.ContractCrossCall   // 合约跨合约调用数据
}

// NewProcessedBlockData new一个pool
func NewProcessedBlockData() *ProcessedBlockData {
	return &ProcessedBlockData{
		BlockDetail:          &db.Block{},
		UserList:             map[string]*db.User{},
		Transactions:         map[string]*db.Transaction{},
		UpgradeContractTx:    make([]*db.UpgradeContractTransaction, 0),
		ChainConfigList:      make([]*pbConfig.ChainConfig, 0),
		ContractWriteSetData: map[string]*ContractWriteSetData{},
		InsertContracts:      make([]*db.Contract, 0),
		UpdateContracts:      make([]*db.Contract, 0),
		FungibleContract:     make([]*db.FungibleContract, 0),
		NonFungibleContract:  make([]*db.NonFungibleContract, 0),
		InsertIDAContracts:   make([]*db.IDAContract, 0),
		EvidenceList:         make([]*db.EvidenceContract, 0),
		ContractEvents:       make([]*db.ContractEvent, 0),
		GasRecordList:        make([]*db.GasRecord, 0),
		CrossChainResult:     &CrossChainResult{},
		ContractCrossCalls:   make(map[string]*db.ContractCrossCall, 0),
	}
}

// CrossChainSaveDB 跨链主子链数据
type CrossChainSaveDB struct {
	//InsertSubChainList 新增子链
	InsertSubChainList []*db.CrossSubChainData
	//UpdateSubChainList 更新子链
	UpdateSubChainList []*db.CrossSubChainData
	//SubChainBlockHeight 需要更新的子链高度列表
	SubChainBlockHeight map[string]int64
}

// GetRealtimeCacheData 区块同步缓存数据，供异步更新计算使用
type GetRealtimeCacheData struct {
	TxList         map[string]*db.Transaction
	ContractAddrs  map[string]string
	GasRecords     []*db.GasRecord
	ContractEvents []*db.ContractEvent
	//CrossTransfers   []*db.CrossTransactionTransfer
	UserInfoMap      map[string]*db.User
	CrossChainResult *CrossChainResult
}

// NewGetRealtimeCacheData 初始化GetRealtimeCacheData
func NewGetRealtimeCacheData() *GetRealtimeCacheData {
	return &GetRealtimeCacheData{
		TxList:         map[string]*db.Transaction{},
		ContractAddrs:  map[string]string{},
		GasRecords:     []*db.GasRecord{},
		ContractEvents: []*db.ContractEvent{},
		//CrossTransfers: []*db.CrossTransactionTransfer{},
		UserInfoMap: map[string]*db.User{},
		CrossChainResult: &CrossChainResult{
			CrossTransfer:       make(map[string]*db.CrossTransactionTransfer),
			InsertCrossTransfer: make([]*db.CrossTransactionTransfer, 0),
			UpdateCrossTransfer: make([]*db.CrossTransactionTransfer, 0),
			BusinessTxMap:       make(map[string]*db.CrossBusinessTransaction),
		},
	}
}

// ContractWriteSetData 读写集解析合约数据
type ContractWriteSetData struct {
	ContractName    string
	ContractNameBak string
	ContractSymbol  string
	ContractAddr    string
	ContractType    string
	Version         string
	RuntimeType     string
	ContractStatus  int32
	BlockHeight     int64
	BlockHash       string
	OrgId           string
	SenderTxId      string
	Sender          string
	SenderOrgId     string
	SenderAddr      string
	Timestamp       int64
	Decimals        int
}

// DelayedUpdateData 异步计算存储数据
type DelayedUpdateData struct {
	DelayCrossChain     *DelayCrossChain
	InsertGasList       []*db.Gas
	UpdateGasList       []*db.Gas
	UpdateTxBlack       *db.UpdateTxBlack
	ContractResult      *db.GetContractResult
	FungibleTransfer    []*db.FungibleTransfer
	NonFungibleTransfer []*db.NonFungibleTransfer
	BlockPosition       *db.BlockPosition
	UpdateAccountResult *db.UpdateAccountResult
	TokenResult         *db.TokenResult
	ContractMap         map[string]*db.Contract
	IDAInsertAssetsData *db.IDAAssetsDataDB
	IDAUpdateAssetsData *db.IDAAssetsUpdateDB
	DIDSaveDate         *db.DIDSaveData
	ChainStatistics     *db.Statistics
	ABITopicTableEvents map[string][]map[string]interface{}
}

type DelayCrossChain struct {
	InsertSubChainCross []*db.CrossSubChainCrossChain
	UpdateSubChainCross []*db.CrossSubChainCrossChain
	InsertSubChainData  []*db.CrossSubChainData
	UpdateSubChainData  []*db.CrossSubChainData
	CrossChainContracts []*db.CrossChainContract
}

// GetDBResult 需要用到的数据库数据
type GetDBResult struct {
	GasList                []*db.Gas
	PositionMapList        map[string][]*db.FungiblePosition
	NonPositionMapList     map[string][]*db.NonFungiblePosition
	FungibleContractMap    map[string]*db.FungibleContract
	NonFungibleContractMap map[string]*db.NonFungibleContract
	AddBlackTxList         []*db.Transaction
	DeleteBlackTxList      []*db.BlackTransaction
	CrossSubChainCross     []*db.CrossSubChainCrossChain
	CrossSubChainMap       map[string]*db.CrossSubChainData
	AccountBNSList         []*db.Account
	AccountDIDList         []*db.Account
	AccountDBMap           map[string]*db.Account
	IDAContractMap         map[string]*db.IDAContract
	IDAAssetDetailMap      map[string]*db.IDAAssetDetail
	EventTopicMap          map[string]map[string]*db.ContractEventTopic
}

type BatchDelayedUpdateLog struct {
	GetRealtimeCacheTime      int64
	DelayedUpdateDataTime     int64
	UpdateDataToDBTime        int64
	UpdateBlockStatusToDBTime int64
}

type TopicEventResult struct {
	AddBlack          []string
	DeleteBlack       []string
	IdentityContract  []*db.IdentityContract
	ContractEventData []*db.ContractEventData
	OwnerAdders       []string
	BNSBindEventData  []*db.BNSTopicEventData
	BNSUnBindDomain   []string
	IDAEventData      *IDAEventData
	DIDEventData      *DIDEventData
}

type DIDEventData struct {
	DIDUnBinds       []string
	Account          map[string]*db.Account
	DIDDetail        []*db.DIDDetail
	DIDAddBlacks     []string
	DIDDeleteBlacks  []string
	DIDAddIssuers    []string
	DIDDeleteIssuers []string
	DIDSetHistory    []*db.DIDSetHistory
	VCIssueHistory   []*db.VCIssueHistory
	VCTemplate       []*db.VCTemplate
	VCDeleteIds      []string
}

// 创建一个新的DIDEventData结构体
// 创建一个新的DIDEventData结构体
// 初始化DIDUnBinds字段，类型为string切片，长度为0
func NewDIDEventData() *DIDEventData {
	return &DIDEventData{
		// 初始化Account字段，类型为map[string]*db.Account，长度为0
		DIDUnBinds: make([]string, 0),
		// 初始化DIDDetail字段，类型为*db.DIDDetail切片，长度为0
		Account: make(map[string]*db.Account),
		// 初始化DIDAddBlacks字段，类型为string切片，长度为0
		DIDDetail: make([]*db.DIDDetail, 0),
		// 初始化DIDDeleteBlacks字段，类型为string切片，长度为0
		DIDAddBlacks: make([]string, 0),
		// 初始化DIDAddIssuers字段，类型为string切片，长度为0
		DIDDeleteBlacks: make([]string, 0),
		// 初始化DIDDeleteIssuers字段，类型为string切片，长度为0
		DIDAddIssuers: make([]string, 0),
		// 初始化DIDSetHistory字段，类型为*db.DIDSetHistory切片，长度为0
		DIDDeleteIssuers: make([]string, 0),
		// 初始化VCIssueHistory字段，类型为*db.VCIssueHistory切片，长度为0
		DIDSetHistory: make([]*db.DIDSetHistory, 0),
		// 初始化VCTemplate字段，类型为*db.VCTemplate切片，长度为0
		VCIssueHistory: make([]*db.VCIssueHistory, 0),
		VCTemplate:     make([]*db.VCTemplate, 0),
		VCDeleteIds:    make([]string, 0),
	}
}

type IDAEventData struct {
	IDACreatedMap     map[string][]*db.IDACreatedInfo
	IDAUpdatedMap     map[string][]*db.EventIDAUpdatedData
	IDADeletedCodeMap map[string][]string
	EventTime         int64
}

// CrossChainResult 跨链主子链数据
type CrossChainResult struct {
	//CrossMainTransaction 跨链交易主链交易
	CrossMainTransaction []*db.CrossMainTransaction
	//CrossTransfer 跨链交易流转数据
	CrossTransfer       map[string]*db.CrossTransactionTransfer
	InsertCrossTransfer []*db.CrossTransactionTransfer
	UpdateCrossTransfer []*db.CrossTransactionTransfer
	//BusinessTxMap 跨链交易-具体业务交易
	BusinessTxMap map[string]*db.CrossBusinessTransaction
}

func NewCrossChainResult() *CrossChainResult {
	return &CrossChainResult{
		// SaveSubChainList:      make(map[string]*db.CrossSubChainData, 0),
		CrossMainTransaction: make([]*db.CrossMainTransaction, 0),
		CrossTransfer:        make(map[string]*db.CrossTransactionTransfer),
		InsertCrossTransfer:  make([]*db.CrossTransactionTransfer, 0),
		UpdateCrossTransfer:  make([]*db.CrossTransactionTransfer, 0),
		BusinessTxMap:        make(map[string]*db.CrossBusinessTransaction),
		// SubChainBlockHeight:   make(map[string]int64),
		// CrossChainContractMap: make(map[string]map[string]string),
	}
}
