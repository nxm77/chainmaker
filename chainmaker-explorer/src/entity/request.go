/*
Package entity comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package entity

import (
	"errors"
	"mime/multipart"
)

// nolint
const (
	OffsetDefault = 0
	OffsetMin     = 0
	LimitDefault  = 10
	LimitMax      = 200
	LimitMaxSpec  = 1000
	LimitMin      = 0
	ListTotalMax  = 200000

	Project = "chainmaker"
	CMB     = "cmb"

	GetOverviewData = "GetOverviewData"
	Decimal         = "Decimal"
	GetChainConfig  = "GetChainConfig"
	GetChainList    = "GetChainList"
	Search          = "Search"

	GetNodeList      = "GetNodeList"
	GetOrgList       = "GetOrgList"
	GetAccountList   = "GetAccountList"
	GetAccountDetail = "GetAccountDetail"

	GetTxDetail       = "GetTxDetail"
	GetLatestTxList   = "GetLatestTxList"
	GetTxList         = "GetTxList"
	GetContractTxList = "GetContractTxList"
	GetBlockTxList    = "GetBlockTxList"
	GetUserTxList     = "GetUserTxList"
	GetQueryTxList    = "GetQueryTxList"
	GetTxNumByTime    = "GetTxNumByTime"
	GetTxDependencies = "GetTxDependencies"

	GetContractVersionList = "GetContractVersionList"
	GetUserList            = "GetUserList"

	GetBlockDetail       = "GetBlockDetail"
	GetLatestBlockList   = "GetLatestBlockList"
	GetBlockList         = "GetBlockList"
	GetContractEventList = "GetContractEventList"
	GetContractCode      = "GetContractCode"
	GetContractTypes     = "GetContractTypes"

	GetLatestContractList = "GetLatestContractList"
	GetContractDetail     = "GetContractDetail"
	GetContractList       = "GetContractList"
	GetEventList          = "GetEventList"

	GetFTContractList       = "GetFTContractList"
	GetFTContractDetail     = "GetFTContractDetail"
	GetNFTContractList      = "GetNFTContractList"
	GetNFTContractDetail    = "GetNFTContractDetail"
	GetEvidenceContractList = "GetEvidenceContractList"
	GetEvidenceContract     = "GetEvidenceContract"
	GetIdentityContractList = "GetIdentityContractList"
	GetIdentityContract     = "GetIdentityContract"

	GetFTTransferList  = "GetFTTransferList"
	GetNFTTransferList = "GetNFTTransferList"
	GetNFTList         = "GetNFTList"
	GetNFTDetail       = "GetNFTDetail"

	GetFTPositionList     = "GetFTPositionList"
	GetUserFTPositionList = "GetUserFTPositionList"
	GetNFTPositionList    = "GetNFTPositionList"

	GetGasRecordList = "GetGasRecordList"
	GetGasList       = "GetGasList"
	GetGasInfo       = "GetGasInfo"

	// 订阅接口访问
	SubscribeChain  = "SubscribeChain"
	ModifySubscribe = "ModifySubscribe"
	DeleteSubscribe = "DeleteSubscribe"
	CancelSubscribe = "CancelSubscribe"

	//ModifyTxBlackList 更新操作
	ModifyTxBlackList = "ModifyTxBlackList"
	ModifyUserStatus  = "ModifyUserStatus"

	UpdateTxSensitiveWord           = "UpdateTxSensitiveWord"
	UpdateEventSensitiveWord        = "UpdateEventSensitiveWord"
	UpdateEvidenceSensitiveWord     = "UpdateEvidenceSensitiveWord"
	RecoverEvidenceSensitiveWord    = "RecoverEvidenceSensitiveWord"
	UpdateNFTSensitiveWord          = "UpdateNFTSensitiveWord"
	UpdateContractNameSensitiveWord = "UpdateContractNameSensitiveWord"

	//数据要素相关
	//GetIDAContractList 获取IDA合约列表
	GetIDAContractList = "GetIDAContractList"
	//GetIDADataList ida资产列表
	GetIDADataList = "GetIDADataList"
	//GetIDADataDetail IDA资产详情
	GetIDADataDetail = "GetIDADataDetail"

	//GetContractVersions 合约版本列表
	GetContractVersions = "GetContractVersions"
	//GetCompilerVersions 编译器版本列表
	GetCompilerVersions = "GetCompilerVersions"
	//GetOpenLicenseTypes 开源版本列表
	GetOpenLicenseTypes = "GetOpenLicenseTypes"
	//GetEvmVersions EVM版本列表
	GetEvmVersions = "GetEvmVersions"
	//GetContractSourceCode 合约源码
	GetContractSourceCode = "GetContractSourceCode"
	//VerifyContractSourceCode 验证合约源码
	VerifyContractSourceCode = "VerifyContractSourceCode"
	//GetGoIDEVersions GoIDE版本列表
	GetGoIDEVersions = "GetGoIDEVersions"

	GetDIDList       = "GetDIDList"
	GetDIDDetail     = "GetDIDDetail"
	GetDIDSetHistory = "GetDIDSetHistory"

	GetVcIssueHistory = "GetVcIssueHistory"
	GetVcTemplateList = "GetVcTemplateList"

	PluginLogin  = "PluginLogin"
	AccountLogin = "AccountLogin"
	Logout       = "Logout"
	CheckLogin   = "CheckLogin"

	UploadContractABI       = "UploadContractABI"
	GetContractABIData      = "GetContractABIData"
	GetContractTopics       = "GetContractTopics"
	GetDecodeContractEvents = "GetDecodeContractEvents"

	GetCrossContractCalls = "GetCrossContractCalls"
)

var ContractSortField = map[string]string{
	"CreateTimestamp": "timestamp",
	"Timestamp":       "timestamp",
	"TxNum":           "txNum",
}

const (
	OrderByDefault = "timestamp"
	OrderDesc      = "desc"
	OrderAsc       = "asc"
)

// RequestBody body
type RequestBody interface {
	// IsLegal 是否合法
	IsLegal() bool
}

// 新增扩展接口
type RequestBodyValidatable interface {
	Validate() error
}

// ChainBody chain
type ChainBody struct {
	ChainId string
}

func (rangeBody *RangeBody) Validate() error {
	if rangeBody.Limit > LimitMax || rangeBody.Limit < LimitMin || rangeBody.Offset < OffsetMin {
		return errors.New("limit or offset is invalid")
	}

	if rangeBody.Limit*rangeBody.Offset > ListTotalMax {
		return errors.New("请求失败，请求数据量过多")
	}

	if rangeBody.Limit == 0 {
		rangeBody.Limit = LimitDefault
	}

	return nil
}

// IsLegal legal
func (chainBody *ChainBody) IsLegal() bool {
	// 不为空即合法
	return chainBody.ChainId != ""
}

// RangeBody range
type RangeBody struct {
	Offset int
	Limit  int
}

// IsLegal legal
func (rangeBody *RangeBody) IsLegal() bool {
	if rangeBody.Limit > LimitMax || rangeBody.Limit < LimitMin || rangeBody.Offset < OffsetMin {
		return false
	}

	if rangeBody.Limit*rangeBody.Offset > ListTotalMax {
		return false
	}

	if rangeBody.Limit == 0 {
		rangeBody.Limit = LimitDefault
	}

	return true
}

// NewRangeBody new
func NewRangeBody() *RangeBody {
	return &RangeBody{
		Offset: OffsetDefault,
		Limit:  LimitDefault,
	}
}

// GetChainIdParams get
type GetChainIdParams struct {
	ChainId string
}

// IsLegal legal
func (params *GetChainIdParams) IsLegal() bool {
	return len(params.ChainId) > 0
}

// GetUserListParams get
type GetUserListParams struct {
	ChainId   string
	UserIds   string
	UserAddrs string
	OrgId     string
	RangeBody
}

// IsLegal legal
func (params *GetUserListParams) IsLegal() bool {
	// 调用RangeBody的合法性校验
	if !params.RangeBody.IsLegal() {
		return false
	}

	return len(params.ChainId) > 0
}

// GetAccountListParams get
type GetAccountListParams struct {
	ChainId string
	Address string
	DID     string
	RangeBody
}

// Validate legal
func (params *GetAccountListParams) Validate() error {
	// 调用RangeBody的合法性校验
	if err := params.RangeBody.Validate(); err != nil {
		return err
	}

	if params.ChainId == "" {
		return errors.New("ChainId: 不可为空")
	}

	return nil
}

// GetAccountDetailParams get
type GetAccountDetailParams struct {
	ChainId string
	Address string
	BNS     string
}

// IsLegal legal
func (params *GetAccountDetailParams) IsLegal() bool {
	return len(params.ChainId) > 0
}

// GetTransactionNumByTimeParams get
type GetTransactionNumByTimeParams struct {
	ChainId   string
	SortType  int
	StartTime int64
	EndTime   int64
	Interval  int64
}

// IsLegal legal
func (params *GetTransactionNumByTimeParams) IsLegal() bool {
	return len(params.ChainId) > 0
}

// GetEventListParams get
type GetEventListParams struct {
	RangeBody
	ChainId      string
	ContractName string
	ContractAddr string
	TxId         string
	Topic        string
}

// Validate legal
func (params *GetEventListParams) Validate() error {
	// 调用RangeBody的合法性校验
	if err := params.RangeBody.Validate(); err != nil {
		return err
	}

	if params.ChainId == "" {
		return errors.New("ChainId: 不可为空")
	}
	return nil
}

// GetLatestContractParams get
type GetLatestContractParams struct {
	ChainId string
}

// IsLegal legal
func (params *GetLatestContractParams) IsLegal() bool {
	return len(params.ChainId) > 0
}

// GetLatestChainParams get
type GetLatestChainParams struct {
	Number int
}

// IsLegal legal
func (params *GetLatestChainParams) IsLegal() bool {
	return params.Number > 0
}

// GetLatestBlockParams get
type GetLatestBlockParams struct {
	ChainId string
}

// IsLegal legal
func (params *GetLatestBlockParams) IsLegal() bool {
	return true
}

// GetBlockDetailParams get
type GetBlockDetailParams struct {
	ChainId     string
	BlockHash   string
	BlockHeight *int64
}

// IsLegal legal
func (params *GetBlockDetailParams) IsLegal() bool {
	if params.BlockHash == "" && params.BlockHeight == nil {
		return false
	}
	return params.ChainId != ""

}

// GetTxDependenciesParams get
type GetTxDependenciesParams struct {
	ChainId     string
	BlockHeight int64
}

// IsLegal legal
func (params *GetTxDependenciesParams) IsLegal() bool {
	if params.BlockHeight < 0 {
		return false
	}
	return params.ChainId != ""

}

// GetContractDetailParams get
type GetContractDetailParams struct {
	ChainId     string
	ContractKey string
}

// IsLegal legal
func (params *GetContractDetailParams) IsLegal() bool {
	return len(params.ChainId) > 0 && len(params.ContractKey) > 0
}

// GetContractCodeParams get
type GetContractCodeParams struct {
	ChainId         string
	ContractAddr    string
	ContractVersion string
}

// IsLegal legal
func (params *GetContractCodeParams) IsLegal() bool {
	return params.ChainId != "" && params.ContractAddr != "" && params.ContractVersion != ""
}

// GetContractListParams get
type GetContractListParams struct {
	RangeBody
	ChainId      string
	Creators     string
	CreatorAddrs string
	Upgraders    string
	UpgradeAddrs string
	RuntimeType  string
	ContractKey  string
	ContractType string
	Status       *int32 //合约状态 -1:全部合约 0：正常 1：已冻结 2：已注销
	Order        string //排序方式，desc:降序 asc:升序
	OrderBy      string //排序字段，创建时间：timestamp，交易数量：txNum
}

// IsLegal legal
func (params *GetContractListParams) IsLegal() bool {
	// 调用RangeBody的合法性校验
	if !params.RangeBody.IsLegal() {
		return false
	}

	if params.Order == "" {
		params.Order = OrderDesc // 默认按降序排序
	}

	if params.OrderBy == "" {
		params.OrderBy = OrderByDefault // 默认按创建时间排序
	} else {
		// 检查 params.OrderBy 是否在 ContractSortField 的键中
		if value, exists := ContractSortField[params.OrderBy]; exists {
			params.OrderBy = value
		}
	}

	return params.ChainId != ""

}

// GetFTContractListParams get
type GetFTContractListParams struct {
	RangeBody
	ChainId     string
	ContractKey string
	Order       string //排序方式，desc:降序 asc:升序
	OrderBy     string //排序字段，创建时间：timestamp，交易数量：txNum
}

// IsLegal legal
func (params *GetFTContractListParams) IsLegal() bool {
	// 调用RangeBody的合法性校验
	if !params.RangeBody.IsLegal() {
		return false
	}

	if params.Order == "" {
		params.Order = OrderDesc // 默认按降序排序
	}

	if params.OrderBy == "" {
		params.OrderBy = OrderByDefault // 默认按创建时间排序
	} else {
		// 检查 params.OrderBy 是否在 ContractSortField 的键中
		if value, exists := ContractSortField[params.OrderBy]; exists {
			params.OrderBy = value
		}
	}

	return params.ChainId != ""

}

// GetFungibleContractParams get
type GetFungibleContractParams struct {
	ChainId      string
	ContractAddr string
}

// IsLegal legal
func (params *GetFungibleContractParams) IsLegal() bool {
	return len(params.ChainId) > 0 && len(params.ContractAddr) > 0
}

// GetNFTContractListParams get
type GetNFTContractListParams struct {
	RangeBody
	ChainId     string
	ContractKey string
	Order       string //排序方式，desc:降序 asc:升序
	OrderBy     string //排序字段，创建时间：timestamp，交易数量：txNum
}

// IsLegal legal
func (params *GetNFTContractListParams) IsLegal() bool {
	// 调用RangeBody的合法性校验
	if !params.RangeBody.IsLegal() {
		return false
	}

	if params.Order == "" {
		params.Order = OrderDesc // 默认按降序排序
	}

	if params.OrderBy == "" {
		params.OrderBy = OrderByDefault // 默认按创建时间排序
	} else {
		// 检查 params.OrderBy 是否在 ContractSortField 的键中
		if value, exists := ContractSortField[params.OrderBy]; exists {
			params.OrderBy = value
		}
	}

	return params.ChainId != ""

}

// GetNonFungibleContractParams get
type GetNonFungibleContractParams struct {
	ChainId      string
	ContractAddr string
}

// IsLegal legal
func (params *GetNonFungibleContractParams) IsLegal() bool {
	return len(params.ChainId) > 0 && len(params.ContractAddr) > 0
}

// GetEvidenceContractListParams get
type GetEvidenceContractListParams struct {
	RangeBody
	ChainId     string
	ContractKey string
}

// IsLegal legal
func (params *GetEvidenceContractListParams) IsLegal() bool {
	// 调用RangeBody的合法性校验
	if !params.RangeBody.IsLegal() {
		return false
	}

	return params.ChainId != ""

}

// GetEvidenceContractParams get
type GetEvidenceContractParams struct {
	RangeBody
	ChainId      string
	ContractName string
	TxId         string
	SenderAddrs  string
	Hashs        string
	Code         int
	Search       string
}

// IsLegal legal
func (params *GetEvidenceContractParams) IsLegal() bool {
	return len(params.ChainId) > 0
}

// GetIdentityContractListParams get
type GetIdentityContractListParams struct {
	RangeBody
	ChainId     string
	ContractKey string
}

// IsLegal legal
func (params *GetIdentityContractListParams) IsLegal() bool {
	// 调用RangeBody的合法性校验
	if !params.RangeBody.IsLegal() {
		return false
	}

	return params.ChainId != ""

}

// GetIdentityContractParams get
type GetIdentityContractParams struct {
	RangeBody
	ChainId      string
	ContractAddr string
	UserAddrs    []string
}

// IsLegal legal
func (params *GetIdentityContractParams) IsLegal() bool {
	return len(params.ChainId) > 0 && len(params.ContractAddr) > 0
}

// ChainOverviewDataParams param
type ChainOverviewDataParams struct {
	ChainId string
}

// IsLegal legal
func (params *ChainOverviewDataParams) IsLegal() bool {
	return params.ChainId != ""
}

// GetChainListParams get
type GetChainListParams struct {
	ChainId   string
	StartTime int64
	EndTime   int64
	RangeBody
}

// IsLegal legal
func (params *GetChainListParams) IsLegal() bool {
	return true
}

// GetOrgListParams get
type GetOrgListParams struct {
	ChainId string
	OrgId   string
	RangeBody
}

// IsLegal legal
func (params *GetOrgListParams) IsLegal() bool {
	// 调用RangeBody的合法性校验
	if !params.RangeBody.IsLegal() {
		return false
	}

	return params.ChainId != ""
}

// GetTxDetailParams get
type GetTxDetailParams struct {
	ChainId string
	TxId    string
}

// IsLegal legal
func (params *GetTxDetailParams) IsLegal() bool {
	return len(params.TxId) > 0
}

// GetBlockLatestListParams get
type GetBlockLatestListParams struct {
	ChainId string
}

// IsLegal legal
func (params *GetBlockLatestListParams) IsLegal() bool {
	return len(params.ChainId) > 0
}

// GetTxLatestListParams get
type GetTxLatestListParams struct {
	ChainId string
}

// IsLegal legal
func (params *GetTxLatestListParams) IsLegal() bool {
	return len(params.ChainId) > 0
}

// GetContractVersionListParams get
type GetContractVersionListParams struct {
	RangeBody
	ChainId      string
	ContractName string
	ContractAddr string
	Senders      string
	Status       *int
	RuntimeType  string
	StartTime    int64
	EndTime      int64
}

// IsLegal legal
func (params *GetContractVersionListParams) IsLegal() bool {
	// 调用RangeBody的合法性校验
	if !params.RangeBody.IsLegal() {
		return false
	}

	if params.Status == nil {
		defaultValue := -1
		params.Status = &defaultValue
	}

	return params.ChainId != ""
}

// GetContractTxListParams get
type GetContractTxListParams struct {
	RangeBody
	ChainId        string
	ContractName   string
	ContractAddr   string
	UserAddrs      string
	ContractMethod string
	TxStatus       *int
}

// Validate legal
func (params *GetContractTxListParams) Validate() error {
	// 调用RangeBody的合法性校验
	if err := params.RangeBody.Validate(); err != nil {
		return err
	}

	if params.ChainId == "" {
		return errors.New("ChainId: 不可为空")
	}

	if params.TxStatus == nil {
		defaultValue := -1
		params.TxStatus = &defaultValue
	}
	return nil
}

// IsLegal legal
// func (params *GetContractTxListParams) IsLegal() bool {
// 	// 调用RangeBody的合法性校验
// 	if !params.RangeBody.IsLegal() {
// 		return false
// 	}

// 	if params.TxStatus == nil {
// 		defaultValue := -1
// 		params.TxStatus = &defaultValue
// 	}

// 	return params.ChainId != ""
// }

// InnerGetContractVersionListParams get
type InnerGetContractVersionListParams struct {
	RangeBody
	ChainId      string
	ContractName string
	Creator      string
	RuntimeType  string
	Status       int64
	StartTime    int64
	EndTime      int64
}

// IsLegal legal
func (params *InnerGetContractVersionListParams) IsLegal() bool {
	return len(params.ChainId) > 0
}

// InnerGetTxListParams get
type InnerGetTxListParams struct {
	RangeBody
	ChainId      string
	ContractName string
	Creator      string
	TxId         string
	TxStatus     int
	StartTime    int64
	EndTime      int64
	UserAddr     string
	UserAddrs    []string
}

// IsLegal legal
func (params *InnerGetTxListParams) IsLegal() bool {
	return len(params.ChainId) > 0
}

// InnerModifyTxShowStatusParams inner
type InnerModifyTxShowStatusParams struct {
	TxId    string
	ChainId string
}

// IsLegal legal
func (params *InnerModifyTxShowStatusParams) IsLegal() bool {
	return len(params.TxId) > 0
}

// GetBlockListParams param
type GetBlockListParams struct {
	RangeBody
	ChainId   string
	BlockKey  string // Height or block Hash
	NodeId    string
	StartTime int64
	EndTime   int64
}

// Validate legal
func (params *GetBlockListParams) Validate() error {
	// 调用RangeBody的合法性校验
	if err := params.RangeBody.Validate(); err != nil {
		return err
	}

	if params.ChainId == "" {
		return errors.New("ChainId: 不可为空")
	}

	return nil
}

// GetTxListParams param
type GetTxListParams struct {
	RangeBody
	TxId           string
	ChainId        string
	ContractName   string
	ContractAddr   string
	ContractMethod string
	BlockHash      string
	UserAddrs      string
	Senders        string
	TxStatus       *int
	StartTime      int64
	EndTime        int64
}

// Validate legal
func (params *GetTxListParams) Validate() error {
	// 调用RangeBody的合法性校验
	if err := params.RangeBody.Validate(); err != nil {
		return err
	}

	if params.ChainId == "" {
		return errors.New("ChainId: 不可为空")
	}

	if params.TxStatus == nil {
		defaultValue := -1
		params.TxStatus = &defaultValue
	}
	return nil
}

// IsLegal legal
// func (params *GetTxListParams) IsLegal() bool {
// 	// 调用RangeBody的合法性校验
// 	if !params.RangeBody.IsLegal() {
// 		return false
// 	}

// 	if params.TxStatus == nil {
// 		//默认全部交易
// 		defaultValue := -1
// 		params.TxStatus = &defaultValue
// 	}

// 	return params.ChainId != ""
// }

// GetQueryTxListParams param
type GetQueryTxListParams struct {
	RangeBody
	TxId           string
	ChainId        string
	ContractName   string
	ContractAddr   string
	ContractMethod string
	UserAddr       string
	TxStatus       int
	StartTime      int64
	EndTime        int64
	Operator       string // "and" 或 "or"
}

// IsLegal legal
func (params *GetQueryTxListParams) IsLegal() bool {
	// 调用RangeBody的合法性校验
	if !params.RangeBody.IsLegal() {
		return false
	}

	if params.Operator != "and" && params.Operator != "or" {
		return false
	}
	return params.ChainId != ""
}

// GetBlockTxListParams param
type GetBlockTxListParams struct {
	RangeBody
	ChainId   string
	BlockHash string
}

func (params *GetBlockTxListParams) Validate() error {
	// 调用RangeBody的合法性校验
	if err := params.RangeBody.Validate(); err != nil {
		return err
	}

	if params.ChainId == "" {
		return errors.New("ChainId: 不可为空")
	}
	if params.BlockHash == "" {
		return errors.New("BlockHash: 不可为空")
	}

	return nil
}

// GetUserTxListParams param
type GetUserTxListParams struct {
	RangeBody
	ChainId   string
	UserAddrs string
}

func (params *GetUserTxListParams) Validate() error {
	// 调用RangeBody的合法性校验
	if err := params.RangeBody.Validate(); err != nil {
		return err
	}

	if params.ChainId == "" {
		return errors.New("ChainId: 不可为空")
	}
	if params.UserAddrs == "" {
		return errors.New("UserAddrs: 不可为空")
	}

	return nil
}

// SearchParams search
type SearchParams struct {
	Type    string
	Value   string
	ChainId string
}

// IsLegal legal
func (params *SearchParams) IsLegal() bool {
	if params.Value == "" || params.ChainId == "" {
		return false
	}
	return true
}

// GetDetailParams get
type GetDetailParams struct {
	Id      string
	ChainId string
}

// IsLegal legal
func (params *GetDetailParams) IsLegal() bool {
	return true
}

// ChainNodesParams param
type ChainNodesParams struct {
	RangeBody
	ChainId string
	OrgId   string
	NodeId  string
}

// IsLegal legal
func (params *ChainNodesParams) IsLegal() bool {
	// 调用RangeBody的合法性校验
	if !params.RangeBody.IsLegal() {
		return false
	}

	return params.ChainId != ""
}

// SubscribeChainParams 订阅链相关
type SubscribeChainParams struct {
	ChainId     string
	OrgId       string
	TlsMode     int
	Tls         bool
	UserSignKey string
	UserSignCrt string
	UserTlsCrt  string
	UserTlsKey  string
	UserEncCrt  string
	UserEncKey  string
	AuthType    string
	HashType    int
	NodeList    []SubscribeNode
}

type SubscribeNode struct {
	Addr        string
	OrgCA       string
	TLSHostName string
}

// IsLegal legal
func (params *SubscribeChainParams) IsLegal() bool {
	if params.ChainId == "" || params.HashType < 0 {
		return false
	}
	return true
}

// CancelSubscribeParams cancel
type CancelSubscribeParams struct {
	ChainId string
	Status  int
}

// IsLegal legal
func (params *CancelSubscribeParams) IsLegal() bool {
	return params.ChainId != ""
}

// ModifySubscribeParams modify
type ModifySubscribeParams struct {
	ChainId     string
	OrgId       string
	TlsMode     int
	Tls         bool
	UserSignKey string
	UserSignCrt string
	UserTlsCrt  string
	UserTlsKey  string
	UserEncCrt  string
	UserEncKey  string
	AuthType    string
	HashType    int
	NodeList    []SubscribeNode
}

// IsLegal legal
func (params *ModifySubscribeParams) IsLegal() bool {
	if params.ChainId == "" || len(params.NodeList) <= 0 {
		return false
	}
	return true
}

// DeleteSubscribeParams delete
type DeleteSubscribeParams struct {
	ChainId string
}

// IsLegal legal
func (params *DeleteSubscribeParams) IsLegal() bool {
	return params.ChainId != ""
}

// GetGasListParams get
type GetGasListParams struct {
	ChainId   string
	UserAddrs string
	RangeBody
}

// IsLegal legal
func (params *GetGasListParams) IsLegal() bool {
	// 调用RangeBody的合法性校验
	if !params.RangeBody.IsLegal() {
		return false
	}

	return len(params.ChainId) > 0
}

// GetGasRecordListParams get
type GetGasRecordListParams struct {
	ChainId      string
	BusinessType int
	UserAddrs    string
	StartTime    int64
	EndTime      int64
	RangeBody
}

// IsLegal legal
func (params *GetGasRecordListParams) IsLegal() bool {
	// 调用RangeBody的合法性校验
	if !params.RangeBody.IsLegal() {
		return false
	}

	return len(params.ChainId) > 0
}

// GetGasInfoParams get
type GetGasInfoParams struct {
	ChainId   string
	UserAddrs string
}

// IsLegal legal
func (params *GetGasInfoParams) IsLegal() bool {
	return len(params.ChainId) > 0
}

// InnerGetChainInfoParams get
type InnerGetChainInfoParams struct {
	ChainId   string
	Address   string
	Addresses []string
}

// IsLegal legal
func (params *InnerGetChainInfoParams) IsLegal() bool {
	return len(params.ChainId) > 0
}

// ModifyUserStatusParams get
type ModifyUserStatusParams struct {
	ChainId string
	Address string
	Status  int
}

// IsLegal legal
func (params *ModifyUserStatusParams) IsLegal() bool {
	return len(params.ChainId) > 0 && len(params.Address) > 0
}

// InnerEvidenceListParams get
type InnerEvidenceListParams struct {
	ChainId     string
	SenderAddr  string
	TxId        string
	Code        int
	Hash        string
	SenderAddrs []string
	RangeBody
}

// IsLegal legal
func (params *InnerEvidenceListParams) IsLegal() bool {
	if params.Limit > LimitMaxSpec || params.Offset < OffsetMin || params.Limit < LimitMin {
		return false
	}
	return len(params.ChainId) > 0
}

// GetNFTListParams get
type GetNFTListParams struct {
	RangeBody
	ChainId     string
	TokenId     string
	ContractKey string
	OwnerAddrs  string
}

// IsLegal legal
func (params *GetNFTListParams) IsLegal() bool {
	// 调用RangeBody的合法性校验
	if !params.RangeBody.IsLegal() {
		return false
	}

	return params.ChainId != ""
}

// GetNFTDetailParams get
type GetNFTDetailParams struct {
	ChainId      string
	TokenId      string
	ContractAddr string
}

// IsLegal legal
func (params *GetNFTDetailParams) IsLegal() bool {
	return params.ChainId != ""
}

// GetTransferListParams get
type GetTransferListParams struct {
	RangeBody
	ChainId      string
	TokenId      string
	ContractName string
}

// IsLegal legal
func (params *GetTransferListParams) IsLegal() bool {
	// 调用RangeBody的合法性校验
	if !params.RangeBody.IsLegal() {
		return false
	}

	return params.ChainId != ""
}

// GetFungibleTransferListParams get
type GetFungibleTransferListParams struct {
	RangeBody
	ChainId      string
	ContractAddr string
	UserAddr     string
}

// IsLegal legal
func (params *GetFungibleTransferListParams) IsLegal() bool {
	// 调用RangeBody的合法性校验
	if !params.RangeBody.IsLegal() {
		return false
	}

	return params.ChainId != ""
}

// GetNonFungibleTransferListParams get
type GetNonFungibleTransferListParams struct {
	RangeBody
	ChainId      string
	ContractAddr string
	UserAddr     string
	TokenId      string
}

// IsLegal legal
func (params *GetNonFungibleTransferListParams) IsLegal() bool {
	// 调用RangeBody的合法性校验
	if !params.RangeBody.IsLegal() {
		return false
	}

	return params.ChainId != ""
}

// GetFungiblePositionListParams get
type GetFungiblePositionListParams struct {
	RangeBody
	ChainId      string
	ContractAddr string
	OwnerAddr    string
}

// IsLegal legal
func (params *GetFungiblePositionListParams) IsLegal() bool {
	// 调用RangeBody的合法性校验
	if !params.RangeBody.IsLegal() {
		return false
	}

	return params.ChainId != ""
}

// GetUserFTPositionListParams get
type GetUserFTPositionListParams struct {
	RangeBody
	ChainId   string
	OwnerAddr string
}

// IsLegal legal
func (params *GetUserFTPositionListParams) IsLegal() bool {
	// 调用RangeBody的合法性校验
	if !params.RangeBody.IsLegal() {
		return false
	}

	return params.ChainId != "" && params.OwnerAddr != ""
}

// GetNonFungiblePositionListParams get
type GetNonFungiblePositionListParams struct {
	RangeBody
	ChainId      string
	ContractAddr string
	OwnerAddr    string
}

// IsLegal legal
func (params *GetNonFungiblePositionListParams) IsLegal() bool {
	// 调用RangeBody的合法性校验
	if !params.RangeBody.IsLegal() {
		return false
	}

	return params.ChainId != "" && params.ContractAddr != ""
}

// ModifyTxBlackListParams update
type ModifyTxBlackListParams struct {
	ChainId string
	TxId    string
	Status  *int
}

// IsLegal legal
func (params *ModifyTxBlackListParams) IsLegal() bool {
	if params.TxId == "" {
		return false
	}
	return params.ChainId != ""
}

// DeleteTxBlackListParams update
type DeleteTxBlackListParams struct {
	ChainId string
	TxId    string
}

// IsLegal legal
func (params *DeleteTxBlackListParams) IsLegal() bool {
	if params.TxId == "" {
		return false
	}
	return params.ChainId != ""
}

// UpdateTxSensitiveWordParams update
type UpdateTxSensitiveWordParams struct {
	ChainId string
	TxId    string
	Column  string
	Status  int
	WarnMsg string
}

// IsLegal legal
func (params *UpdateTxSensitiveWordParams) IsLegal() bool {
	if params.ChainId == "" || params.TxId == "" {
		return false
	}
	return true
}

// UpdateEventSensitiveWordParams update
type UpdateEventSensitiveWordParams struct {
	ChainId string
	TxId    string
	Index   int
	Column  string
	Status  *int
	WarnMsg string
}

// IsLegal legal
func (params *UpdateEventSensitiveWordParams) IsLegal() bool {
	if params.Column == "" || params.ChainId == "" || params.TxId == "" || params.Status == nil {
		return false
	}
	return true
}

// EvidenceSensitiveWordParams update
type EvidenceSensitiveWordParams struct {
	ChainId string
	Hash    string
	Column  string
	Status  *int
	WarnMsg string
}

// IsLegal legal
func (params *EvidenceSensitiveWordParams) IsLegal() bool {
	if params.Column == "" || params.ChainId == "" || params.Hash == "" || params.Status == nil {
		return false
	}
	return true
}

// NFTSensitiveWordParams update
type NFTSensitiveWordParams struct {
	ChainId      string
	TokenId      string
	ContractAddr string
	Column       string
	Status       *int
	WarnMsg      string
}

// IsLegal legal
func (params *NFTSensitiveWordParams) IsLegal() bool {
	if params.Column == "" || params.ChainId == "" || params.TokenId == "" || params.Status == nil {
		return false
	}
	return true
}

// UpdateContractNameSWParams update
type UpdateContractNameSWParams struct {
	ChainId      string
	ContractName string
	Status       *int
	WarnMsg      string
}

// IsLegal legal
func (params *UpdateContractNameSWParams) IsLegal() bool {
	if params.ChainId == "" || params.ContractName == "" || params.Status == nil {
		return false
	}
	return true
}

// GetIDAContractListParams ida合约列表
type GetIDAContractListParams struct {
	ChainId     string
	ContractKey string
	Order       string //排序方式，desc:降序 asc:升序
	OrderBy     string //排序字段，创建时间：timestamp，交易数量：txNum
	RangeBody
}

// IsLegal legal
func (params *GetIDAContractListParams) IsLegal() bool {
	// 调用RangeBody的合法性校验
	if !params.RangeBody.IsLegal() {
		return false
	}

	if params.Order == "" {
		params.Order = OrderDesc // 默认按降序排序
	}

	if params.OrderBy == "" {
		params.OrderBy = OrderByDefault // 默认按创建时间排序
	} else {
		// 检查 params.OrderBy 是否在 ContractSortField 的键中
		if value, exists := ContractSortField[params.OrderBy]; exists {
			params.OrderBy = value
		}
	}

	return params.ChainId != ""
}

// GetIDADataListParams ida资产列表
type GetIDADataListParams struct {
	ChainId      string
	AssetCode    string
	ContractAddr string
	RangeBody
}

// IsLegal legal
func (params *GetIDADataListParams) IsLegal() bool {
	// 调用RangeBody的合法性校验
	if !params.RangeBody.IsLegal() {
		return false
	}

	return params.ChainId != "" && params.ContractAddr != ""
}

// GetIDADataDetailParams ida资产详情
type GetIDADataDetailParams struct {
	ChainId      string
	AssetCode    string
	ContractAddr string
}

// IsLegal legal
func (params *GetIDADataDetailParams) IsLegal() bool {
	return params.ChainId != "" && params.AssetCode != "" && params.ContractAddr != ""
}

// GetContractVersionsParams get
type GetContractVersionsParams struct {
	ChainId      string
	ContractAddr string
}

// IsLegal legal
func (params *GetContractVersionsParams) IsLegal() bool {
	return params.ChainId != "" && params.ContractAddr != ""
}

// 定义一个结构体VerifyContractParams，用于验证合约参数
type VerifyContractParams struct {
	// 链ID
	ChainId string `json:"ChainId" form:"ChainId"`
	// 合约地址
	ContractAddr string `json:"ContractAddr" form:"ContractAddr"`
	// 合约版本
	ContractVersion string `json:"ContractVersion" form:"ContractVersion"`
	// 编译器路径
	CompilerPath string `json:"CompilerPath" form:"CompilerPath" `
	// 编译器版本
	CompilerVersion string `json:"CompilerVersion" form:"CompilerVersion"`
	// 开源许可证类型
	OpenLicenseType string `json:"OpenLicenseType" form:"OpenLicenseType"`
	// 合约源文件
	ContractSourceFile *multipart.FileHeader `form:"ContractSourceFile"`
	// 是否优化
	Optimization bool `json:"Optimization" form:"Optimization"`
	// 运行次数
	Runs int `json:"Runs" form:"Runs"`
	// EVM版本
	EvmVersion string `json:"EvmVersion" form:"EvmVersion"`
}

// IsLegal legal
func (params *VerifyContractParams) IsLegal() bool {
	return params.ChainId != "" && params.ContractAddr != "" &&
		params.ContractVersion != "" && params.CompilerVersion != ""
}

// GetDqueryChartParams
type GetDqueryChartParams struct {
	ChainId      string
	ChartType    string
	BusinessType string
}

// IsLegal legal
func (params *GetDqueryChartParams) IsLegal() bool {
	return params.ChainId != ""
}

// SaveDqueryChartParams
type SaveDqueryChartParams struct {
	ChainId      string
	ChartType    string
	BusinessType string
	TaskId       string
}

// IsLegal legal
func (params *SaveDqueryChartParams) IsLegal() bool {
	return params.ChainId != "" && params.ChartType != "" && params.TaskId != ""
}

// GetDIDListParams param
type GetDIDListParams struct {
	RangeBody
	ChainId      string
	ContractAddr string
	DID          string
}

// IsLegal legal
func (params *GetDIDListParams) IsLegal() bool {
	// 调用RangeBody的合法性校验
	if !params.RangeBody.IsLegal() {
		return false
	}

	return params.ChainId != "" && params.ContractAddr != ""
}

// GetDIDDetailParams
type GetDIDDetailParams struct {
	ChainId string
	DID     string
}

// IsLegal legal
func (params *GetDIDDetailParams) IsLegal() bool {
	return params.ChainId != "" && params.DID != ""
}

// GetDIDHistoryParams param
type GetDIDHistoryParams struct {
	RangeBody
	ChainId string
	DID     string
}

// IsLegal legal
func (params *GetDIDHistoryParams) IsLegal() bool {
	// 调用RangeBody的合法性校验
	if !params.RangeBody.IsLegal() {
		return false
	}

	return params.ChainId != "" && params.DID != ""
}

// GetVCListParams param
type GetVCListParams struct {
	RangeBody
	ChainId      string
	IssuerDID    string
	HolderDID    string
	TemplateID   string
	ContractAddr string
	VcID         string
}

// IsLegal legal
func (params *GetVCListParams) IsLegal() bool {
	// 调用RangeBody的合法性校验
	if !params.RangeBody.IsLegal() {
		return false
	}

	return params.ChainId != ""
}

// GetVCTemplateParams param
type GetVCTemplateParams struct {
	RangeBody
	ChainId      string
	ContractAddr string
	TemplateID   string
}

// IsLegal legal
func (params *GetVCTemplateParams) IsLegal() bool {
	// 调用RangeBody的合法性校验
	if !params.RangeBody.IsLegal() {
		return false
	}

	return params.ChainId != ""
}

// PluginLoginParams post
type PluginLoginParams struct {
	PubKey     string `json:"pubKey"`
	SignBase64 string `json:"signBase64"`
}

// IsLegal legal
func (params *PluginLoginParams) IsLegal() bool {
	return params.PubKey != "" && params.SignBase64 != ""
}

// AccountLoginParams post
type AccountLoginParams struct {
	RandomNum int64  `json:"randomNum"` //Sign string `json:"sign"`
	Password  string `json:"password"`  //Address string `json:"address"`
}

// IsLegal legal
func (params *AccountLoginParams) IsLegal() bool {
	return params.RandomNum != 0 && params.Password != ""
}

// GetDIDDetailParams
type GetContractTypesParams struct {
	ChainId string
}

// IsLegal legal
func (params *GetContractTypesParams) IsLegal() bool {
	return params.ChainId != ""
}

// UploadContractAbiParams post
type UploadContractAbiParams struct {
	ChainId         string                `json:"ChainId" form:"ChainId"`
	ContractAddr    string                `json:"ContractAddr" form:"ContractAddr"`
	ContractVersion string                `json:"ContractVersion" form:"ContractVersion"`
	AbiJson         *multipart.FileHeader `json:"AbiJson" form:"AbiJson"`
}

// IsLegal legal
func (params *UploadContractAbiParams) IsLegal() bool {
	return params.ChainId != "" && params.ContractAddr != "" && params.ContractVersion != ""
}

type GetABIDataParams struct {
	ChainId         string `json:"ChainId" form:"ChainId"`
	ContractAddr    string `json:"ContractAddr" form:"ContractAddr"`
	ContractVersion string `json:"ContractVersion" form:"ContractVersion"`
}

// IsLegal legal
func (params *GetABIDataParams) IsLegal() bool {
	return params.ChainId != "" && params.ContractAddr != "" && params.ContractVersion != ""
}

type GetContractTopicsParams struct {
	ChainId         string `json:"ChainId" form:"ChainId"`
	ContractAddr    string `json:"ContractAddr" form:"ContractAddr"`
	ContractVersion string `json:"ContractVersion" form:"ContractVersion"`
}

// IsLegal legal
func (params *GetContractTopicsParams) IsLegal() bool {
	return params.ChainId != "" && params.ContractAddr != "" && params.ContractVersion != ""
}

type GetDecodeContractEventsParams struct {
	RangeBody
	ChainId         string        `json:"ChainId" form:"ChainId"`
	ContractAddr    string        `json:"ContractAddr" form:"ContractAddr"`
	ContractVersion string        `json:"ContractVersion" form:"ContractVersion"`
	Topic           string        `json:"Topic" form:"Topic"`
	SearchParams    []SearchParam `json:"SearchParams" form:"SearchParams"`
}

type SearchParam struct {
	Name  string `json:"name" form:"name"`
	Value string `json:"value" form:"value"`
}

// IsLegal legal
func (params *GetDecodeContractEventsParams) IsLegal() bool {
	// 调用RangeBody的合法性校验
	if !params.RangeBody.IsLegal() {
		return false
	}

	return params.ChainId != "" && params.ContractAddr != "" && params.ContractVersion != ""
}

type GetCrossContractCallsParams struct {
	ChainId      string `json:"ChainId" form:"ChainId"`
	ContractAddr string `json:"ContractAddr" form:"ContractAddr"`
}

// IsLegal legal
func (params *GetCrossContractCallsParams) IsLegal() bool {
	return params.ChainId != "" && params.ContractAddr != ""
}
