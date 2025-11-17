//nolint:lll
/*
Package db comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package db

import (
	"time"

	"github.com/shopspring/decimal"
)

const (
	// SubscribeOK ok
	SubscribeOK = 0
	// SubscribeFailed fail
	SubscribeFailed = 1
	// SubscribeCanceled cancel
	SubscribeCanceled = 2
	//// SubscribeDeleting delete
	//SubscribeDeleting = 3
)

const (
	ABISystemFieldID          = "sysId"
	ABISystemFieldTxID        = "sysTxId"
	ABISystemFieldContractVer = "sysContractVersion"
	ABISystemFieldTimestamp   = "sysTimestamp"
)

// 检查是否包含系统字段
var ABISystemFields = map[string]bool{
	ABISystemFieldID:          true,
	ABISystemFieldTxID:        true,
	ABISystemFieldContractVer: true,
	ABISystemFieldTimestamp:   true,
}

type ABIEventSystemFields struct {
	SysID          string `gorm:"primaryKey;column:sysId;type:varchar(128);comment:主键ID"`
	SysTxID        string `gorm:"column:sysTxId;type:varchar(128);comment:交易ID"`
	SysContractVer string `gorm:"column:sysContractVersion;type:varchar(128);index:,composite:version_time;comment:合约版本号"`
	SysTimestamp   int64  `gorm:"column:sysTimestamp;index:,composite:version_time;comment:上链时间"`
}

// CommonIntField common
type CommonIntField struct {
	CreatedAt time.Time `json:"createdAt" gorm:"column:createdAt;autoCreateTime;comment:创建时间"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"column:updatedAt;autoUpdateTime;comment:更新时间"`
}

// Statistics 统计数据结构
type Statistics struct {
	ChainId           string `json:"chainId" gorm:"column:chainId;type:varchar(128);primaryKey;comment:chainID"`
	BlockHeight       int64  `json:"blockHeight" gorm:"column:blockHeight;comment:区块高度"`
	TotalAccounts     int64  `json:"totalAccounts" gorm:"column:totalAccounts;comment:账户总数"`
	TotalTransactions int64  `json:"totalTransactions" gorm:"column:totalTransactions;comment:交易总数"`
	TotalContracts    int64  `json:"totalContracts" gorm:"column:totalContracts;comment:合约总数"`
	TotalOrgs         int64  `json:"totalOrgs" gorm:"column:totalOrgs;comment:组织总数"`
	TotalNodes        int64  `json:"totalNodes" gorm:"column:totalNodes;comment:节点总数"`
	TotalCrossTx      int64  `json:"totalCrossTx" gorm:"column:totalCrossTx;comment:跨链交易总数"`
	CommonIntField
}

// TableName table
func (t *Statistics) TableName() string {
	return TableStatistics
}

// Block
// @Description: 区块数据
type Block struct {
	ID                string    `json:"id" gorm:"primaryKey;comment:主键ID"`
	BlockHeight       int64     `json:"blockHeight" gorm:"column:blockHeight;comment:区块高度;uniqueIndex"`
	PreBlockHash      string    `json:"preBlockHash" gorm:"column:preBlockHash;comment:上个区块的哈希"`
	BlockHash         string    `json:"blockHash" gorm:"column:blockHash;comment:本区块的哈希;type:varchar(128);index;"`
	BlockVersion      int32     `json:"blockVersion" gorm:"column:blockVersion;comment:区块版本号"`
	OrgId             string    `json:"orgId" gorm:"column:orgId;comment:组织ID"`
	BlockDag          string    `json:"blockDag" gorm:"column:blockDag;comment:区块交易的执行依赖顺序"`
	DagHash           string    `json:"dagHash" gorm:"column:dagHash;comment:Dag的hash值"`
	TxCount           int       `json:"txCount" gorm:"column:txCount;comment:区块交易量"`
	Signature         string    `json:"signature" gorm:"column:signature;comment:区块生成者的签名"`
	RwSetHash         string    `json:"rwSetHash" gorm:"column:rwSetHash;comment:交易结果的读写集哈希"`
	TxRootHash        string    `json:"txRootHash" gorm:"column:txRootHash;comment:区块交易的Merkle Root哈希"`
	ProposerId        string    `json:"proposerId" gorm:"column:proposerId;index;comment:区块的生成者ID"`
	ProposerAddr      string    `json:"proposerAddr" gorm:"column:proposerAddr;comment:区块的生成者地址"`
	ConsensusArgs     string    `json:"consensusArgs" gorm:"column:consensusArgs;comment:共识参数"`
	DelayUpdateStatus int       `json:"delayUpdateStatus" gorm:"column:delayUpdateStatus;index;comment:订阅异步更新状态(0:未更新,1:更新成功)"`
	Timestamp         int64     `json:"timestamp" gorm:"column:timestamp;comment:上链时间;index;"`
	TimestampDate     time.Time `json:"timestampDate" gorm:"column:timestampDate;index;comment:上链时间"`
	CommonIntField
}

// TableName table
func (t *Block) TableName() string {
	return TableBlock
}

// UpgradeContractTransaction
// @Description: 合约创建，升级交易列表
type UpgradeContractTransaction struct {
	TxId                string `json:"txId" gorm:"column:txId;type:varchar(128);primaryKey;comment:交易ID"`
	SenderOrgId         string `json:"senderOrgId" gorm:"column:senderOrgId;varchar(128);comment:交易发起者组织"` //发起交易的组织
	Sender              string `json:"sender" gorm:"column:sender;comment:交易发起者ID"`
	UserAddr            string `json:"userAddr" gorm:"column:userAddr;comment:交易发起者地址"`
	BlockHeight         int64  `json:"blockHeight" gorm:"column:blockHeight;comment:区块高度"` //区块高度
	BlockHash           string `json:"blockHash" gorm:"column:blockHash;comment:区块哈希"`
	Timestamp           int64  `json:"timestamp" gorm:"column:timestamp;index;comment:上链时间"`                                                       //交易时间戳
	ContractRuntimeType string `json:"contractRuntimeType" gorm:"column:contractRuntimeType;comment:合约类型"`                                         //合约运行时版本
	ContractName        string `json:"contractName" gorm:"column:contractName;comment:合约名称"`                                                       //合约名称
	ContractNameBak     string `json:"contractNameBak" gorm:"column:contractNameBak;type:varchar(128);index;comment:合约名称"`                         //合约名称
	ContractAddr        string `json:"contractAddr" gorm:"column:contractAddr;type:varchar(128);index:,composite:addr_version;comment:合约地址"`       //合约名称
	ContractVersion     string `json:"contractVersion" gorm:"column:contractVersion;type:varchar(128);index:,composite:addr_version;comment:合约版本"` //合约版本
	ContractType        string `json:"contractType" gorm:"column:contractType;comment:合约类型"`
	VerifyStatus        int    `json:"verifyStatus" gorm:"column:verifyStatus;comment:验证结果，0：未验证，1:验成功。2:验证失败"`
	CommonIntField
}

// TableName table
func (t *UpgradeContractTransaction) TableName() string {
	return TableContractUpgradeTransaction
}

type ContractByteCode struct {
	TxId     string `json:"txId" gorm:"column:txId;type:varchar(128);primaryKey;comment:交易ID"`
	ByteCode []byte `json:"byteCode" gorm:"column:byteCode;comment:合约执行文件"`
	CommonIntField
}

// TableName table
func (t *ContractByteCode) TableName() string {
	return TableContractByteCode
}

// ContractVerifyResult
// @Description: 合约编译结果
type ContractVerifyResult struct {
	VerifyId            string `json:"verifyId" gorm:"column:verifyId;primaryKey;comment:主键ID"`
	VerifyStatus        int    `json:"verifyStatus" gorm:"column:verifyStatus;comment:验证结果，0：未验证，1:验成功。2:验证失败"`
	ContractName        string `json:"contractName" gorm:"column:contractName;comment:合约名称"`
	ContractAddr        string `json:"contractAddr" gorm:"column:contractAddr;type:varchar(128);index;comment:合约地址"`
	ContractVersion     string `json:"contractVersion" gorm:"column:contractVersion;comment:合约版本"`
	ContractRuntimeType string `json:"contractRuntimeType" gorm:"column:contractRuntimeType;comment:合约类型"`
	CompilerPath        string `json:"compilerPath" gorm:"column:compilerPath;comment:编译合约路径"`
	ByteCode            []byte `json:"byteCode" gorm:"column:byteCode;comment:合约编译结果"`
	ABI                 string `json:"abi" gorm:"column:abi;comment:编译合约ABI文件"`
	MetaData            string `json:"metaData" gorm:"column:metaData;comment:编译合约MetaData文件"`
	CompilerVersion     string `json:"compilerVersion" gorm:"column:compilerVersion;comment:合约编译器版本"`
	OpenLicenseType     string `json:"openLicenseType" gorm:"column:openLicenseType;comment:开源许可证"`
	EvmVersion          string `json:"evmVersion" gorm:"column:evmVersion;comment:EVM版本"`
	Optimization        bool   `json:"optimization" gorm:"column:optimization;comment:是否优化"`
	RunNum              int    `json:"runNum" gorm:"column:runNum;comment:优化次数"`
	CommonIntField
}

// TableName table
func (t *ContractVerifyResult) TableName() string {
	return TableContractVerifyResult
}

// ContractSourceFile
// @Description: 合约源码文件
type ContractSourceFile struct {
	ID              string `json:"id" gorm:"primaryKey;comment:主键ID"`
	VerifyId        string `json:"verifyId" gorm:"column:verifyId;index;comment:合约编译Id"`
	ContractAddr    string `json:"contractAddr" gorm:"column:contractAddr;type:varchar(128);index;comment:合约地址"`
	ContractVersion string `json:"contractVersion" gorm:"column:contractVersion;comment:合约版本"`
	SourcePath      string `json:"sourcePath" gorm:"column:sourcePath;comment:合约源码路径"`
	SourceCode      []byte `json:"sourceCode" gorm:"column:sourceCode;comment:合约源码"`
}

// TableName tableping
func (t *ContractSourceFile) TableName() string {
	return TableContractSourceFile
}

// Transaction 交易
type Transaction struct {
	TxId                  string    `json:"txId" gorm:"column:txId;type:varchar(128);primaryKey;comment:交易ID"`
	TimestampDate         time.Time `json:"timestampDate" gorm:"column:timestampDate;index;index:,composite:proposerId_timeDate_resultCode;index:,composite:user_timeDate_resultCode;comment:交易发起时间"`
	Sender                string    `json:"sender" gorm:"column:sender;type:varchar(128);index:,composite:sender_time;comment:交易发起者ID"`
	SenderOrgId           string    `json:"senderOrgId" gorm:"column:senderOrgId;varchar(128);comment:交易发起者组织"` //组织id
	UserAddr              string    `json:"userAddr" gorm:"column:userAddr;type:varchar(128);index:,composite:user_height_time;index:,composite:user_timeDate_resultCode;comment:交易发起者地址;"`
	ContractName          string    `json:"contractName" gorm:"column:contractName;comment:合约名称"`
	ContractNameBak       string    `json:"contractNameBak" gorm:"column:contractNameBak;type:varchar(128);index:,composite:contract_height_time;index;comment:合约名称;"`
	ContractAddr          string    `json:"contractAddr" gorm:"column:contractAddr;type:varchar(128);index;comment:合约地址;"`
	BlockHeight           int64     `json:"blockHeight" gorm:"column:blockHeight;index;index:,composite:height_time;index:,composite:user_height_time;index:,composite:contract_height_time;comment:区块高度"`
	BlockHash             string    `json:"blockHash" gorm:"column:blockHash;type:varchar(128);index;comment:区块哈希;"`
	TxType                string    `json:"txType" gorm:"column:txType;comment:交易类型"`
	ExpirationTime        int64     `json:"expirationTime" gorm:"column:expirationTime;comment:交易的到期的unix时间"`
	TxIndex               int       `json:"txIndex" gorm:"column:txIndex;comment:同一个区块中交易排序"`
	TxStatusCode          string    `json:"txStatusCode" gorm:"column:txStatusCode;comment:交易执行结果"`
	RwSetHash             string    `json:"rwSetHash" gorm:"column:rwSetHash;comment:交易结果的读写集哈希"`
	ContractRuntimeType   string    `json:"contractRuntimeType" gorm:"column:contractRuntimeType;comment:合约类型"`
	ContractType          string    `json:"contractType" gorm:"column:contractType;comment:合约类型"`
	ContractMethod        string    `json:"contractMethod" gorm:"column:contractMethod;type:varchar(512);index;comment:被调用的合约方法名;"`
	ContractParameters    string    `json:"contractParameters" gorm:"column:contractParameters;comment:调用合约方法时传入的参数列表"`
	ContractParametersBak string    `json:"contractParametersBak" gorm:"column:contractParametersBak;comment:调用合约方法时传入的参数列表"`
	ContractVersion       string    `json:"contractVersion" gorm:"column:contractVersion;comment:合约版本号"`
	Endorsement           string    `json:"endorsement" gorm:"column:endorsement;comment:交易发起者签名"`
	Sequence              uint64    `json:"sequence" gorm:"column:sequence;comment:交易的顺序号"`
	ReadSet               string    `json:"readSet" gorm:"column:readSet;comment:交易的读集列表"`
	ReadSetBak            string    `json:"readSetBak" gorm:"column:readSetBak;comment:交易的读集列表"`
	WriteSet              string    `json:"writeSet" gorm:"column:writeSet;comment:交易写集列表"`
	WriteSetBak           string    `json:"writeSetBak" gorm:"column:writeSetBak;comment:交易写集列表"`
	PayerAddr             string    `json:"payerAddr" gorm:"column:payerAddr;comment:代付gas费的用户地址"`
	GasUsed               uint64    `json:"gasUsed" gorm:"column:gasUsed;comment:消耗的gas费"`
	Event                 string    `json:"event" gorm:"column:event;comment:合约执行产生的事件日志"`
	ProposerId            string    `json:"proposerId" gorm:"column:proposerId;index:,composite:proposerId_time;index:,composite:proposerId_timeDate_resultCode;comment:区块的生成者ID"`
	Timestamp             int64     `json:"timestamp" gorm:"column:timestamp;index:,composite:time_resultcode;index:,composite:proposerId_time;index:,composite:sender_time;index:,composite:height_time;index:,composite:contract_height_time;index:,composite:user_height_time;comment:交易发起时间;"`
	ChainTimestamp        int64     `json:"chainTimestamp" gorm:"column:chainTimestamp;comment:交易上链时间;"`
	ContractResultCode    uint32    `json:"contractResultCode" gorm:"column:contractResultCode;index:,composite:time_resultcode;index:,composite:user_timeDate_resultCode;comment:合约执行结果"`
	ContractResult        []byte    `json:"contractResult" gorm:"column:contractResult;comment:合约执行结果"`
	ContractResultBak     []byte    `json:"contractResultBak" gorm:"column:contractResultBak;comment:合约执行结果"`
	ContractMessage       string    `json:"contractMessage" gorm:"column:contractMessage;comment:合约外返回的错误信息"`
	ContractMessageBak    string    `json:"contractMessageBak" gorm:"column:contractMessageBak;comment:合约外返回的错误信息"`
	CommonIntField
}

// TableName table
func (t *Transaction) TableName() string {
	return TableTransaction
}

// BlackTransaction 黑名单交易
type BlackTransaction Transaction

// TableName table
func (t *BlackTransaction) TableName() string {
	return TableBlackTransaction
}

// ContractEvent
// @Description: 合约事件数据
type ContractEvent struct {
	ID              string `json:"id" gorm:"primaryKey;comment:主键ID;"`
	TxId            string `json:"txId" gorm:"column:txId;type:varchar(128);uniqueIndex:,composite:contract_event_txId;comment:交易ID;"`
	EventIndex      int    `json:"eventIndex" gorm:"column:eventIndex;uniqueIndex:,composite:contract_event_txId;comment:同一个交易中的事件排序;"`
	Topic           string `json:"topic" gorm:"column:topic;comment:合约事件topic"`
	TopicBak        string `json:"topicBak" gorm:"column:topicBak;comment:合约事件topic脱敏"`
	ContractName    string `json:"contractName" gorm:"column:contractName;type:varchar(128);comment:合约名称"`
	ContractNameBak string `json:"contractNameBak" gorm:"column:contractNameBak;type:varchar(128);index:,composite:contract_name_time;comment:合约名称;"`
	ContractAddr    string `json:"contractAddr" gorm:"column:contractAddr;type:varchar(128);index:,composite:contract_version;comment:合约地址;"`
	ContractVersion string `json:"contractVersion" gorm:"column:contractVersion;type:varchar(128);index:,composite:contract_version;comment:合约版本号"`
	ContractType    string `json:"contractType" gorm:"column:contractType;comment:合约类型"`
	EventData       string `json:"eventData" gorm:"column:eventData;comment:合约执行产生的事件日志"`
	EventDataBak    string `json:"eventDataBak" gorm:"column:eventDataBak;comment:合约执行产生的事件日志脱敏"`
	Timestamp       int64  `json:"timestamp" gorm:"column:timestamp;index:,composite:contract_name_time;comment:上链时间;"`
	CommonIntField
}

// TableName table
func (t *ContractEvent) TableName() string {
	return TableContractEvent
}

// ContractEventTopic 合约topic列表
type ContractEventTopic struct {
	ID           string `json:"id" gorm:"primaryKey;comment:主键ID;"`
	Topic        string `json:"topic" gorm:"column:topic;comment:合约事件topic"`
	ContractName string `json:"contractName" gorm:"column:contractName;type:varchar(128);index;comment:合约名称"`
	TxNum        int64  `json:"txNum" gorm:"column:txNum;comment:合约topic交易总量"` // 交易数据量
	BlockHeight  int64  `json:"blockHeight" gorm:"column:blockHeight;comment:区块高度"`
	CommonIntField
}

// TableName table
func (t *ContractEventTopic) TableName() string {
	return TableContractEventTopic
}

// Node
// @Description: 节点数据
type Node struct {
	NodeId  string `json:"nodeId" gorm:"column:nodeId;primaryKey;comment:节点ID"` // 节点ID
	OrgId   string `json:"orgId" gorm:"column:orgId;comment:节点所属组织ID"`          // 所属组织ID
	Role    string `json:"role" gorm:"column:role;comment:节点所属角色"`              // 角色
	Address string `json:"address" gorm:"column:address;comment:节点地址"`          // 地址
	Status  int    `json:"status" gorm:"column:status;comment:节点状态"`            // 状态
	CommonIntField
}

// TableName table
func (t *Node) TableName() string {
	return TableNode
}

// Gas
// @Description: 账户gas数据，Todo 可以和account账户表合并的
type Gas struct {
	Address     string `json:"address" gorm:"column:address;type:varchar(128);primaryKey;comment:账户地址;"` // 账户地址
	GasBalance  int64  `json:"gasBalance" gorm:"column:gasBalance;comment:账户gas余额"`                      // 账户余额
	GasTotal    int64  `json:"gasTotal" gorm:"column:gasTotal;comment:账户获得的gas总数"`                       // 账户总获得
	GasUsed     int64  `json:"gasUsed" gorm:"column:gasUsed;comment:账户gas消耗数"`                           // 账户消耗
	BlockHeight int64  `json:"blockHeight" gorm:"column:blockHeight;comment:区块高度，只用作更新的版本号控制"`           // 区块高度
	CommonIntField
}

// TableName table
func (t *Gas) TableName() string {
	return TableGas
}

// GasRecord
// @Description: gas消耗详情
type GasRecord struct {
	ID           string `json:"id" gorm:"primaryKey;comment:主键ID"`
	TxId         string `json:"txId" gorm:"column:txId;type:varchar(128);uniqueIndex:,composite:gas_index_txId;comment:交易ID;"`
	GasIndex     int    `json:"gasIndex" gorm:"column:gasIndex;uniqueIndex:,composite:gas_index_txId;comment:交易中gas充值的排序"`          //批量消耗gas的切片key
	Address      string `json:"address" gorm:"column:address;type:varchar(128);index;comment:gas支付地址;"`                             // 支付地址
	GasAmount    int64  `json:"gasAmount" gorm:"column:gas;comment:gas金额"`                                                          // gas金额
	BusinessType int    `json:"businessType" gorm:"column:businessType;index:,composite:gas_business_time;comment:操作类型 1:领取 2:消耗;"` //
	Timestamp    int64  `json:"timestamp" gorm:"column:timestamp;index:,composite:gas_business_time;comment:上链时间;"`                 //交易时间戳
	CommonIntField
}

// TableName table
func (t *GasRecord) TableName() string {
	return TableGasRecord
}

// Subscribe 订阅链信息
type Subscribe struct {
	ChainId     string `json:"chainId" gorm:"column:chainId;type:varchar(128);primaryKey;comment:chainID"`
	OrgId       string `json:"orgId" gorm:"column:orgId;comment:组织id"`
	UserSignKey string `json:"userSignKey" gorm:"column:userSignKey;type:text;comment:用户签名证书私钥"`
	UserSignCrt string `json:"userSignCrt" gorm:"column:userSignCrt;type:text;comment:用户签名证书"`
	UserTlsKey  string `json:"userTlsKey" gorm:"column:userTlsKey;type:text;comment:用户加密证书私钥"`
	UserTlsCrt  string `json:"userTlsCrt" gorm:"column:userTlsCrt;type:text;comment:用户加密证书"`
	UserEncKey  string `json:"userEncKey" gorm:"column:userEncKey;type:text;comment:用户加密证书私钥"`
	UserEncCrt  string `json:"userEncCrt" gorm:"column:userEncCrt;type:text;comment:用户加密证书"`
	NodeList    string `json:"nodeList" gorm:"column:nodeList;type:text;comment:订阅节点列表"`
	Status      int    `json:"status" gorm:"column:status;comment:订阅状态，0:正常 1: 已暂停"`
	AuthType    string `json:"authType" gorm:"column:authType;comment:链类型"`
	HashType    string `json:"hashType" gorm:"column:hashType;comment:哈希类型"`
	TlsMode     int    `json:"tlsMode" gorm:"column:tlsMode;comment:0：单证书模式，1：双证书模式"`
	Tls         bool   `json:"tls" gorm:"column:tls;comment:是否开启tls"`
	CommonIntField
}

// TableName table
func (*Subscribe) TableName() string {
	return TableSubscribe
}

// Chain 链信息
type Chain struct {
	ChainId           string `json:"chainId" gorm:"column:chainId;type:varchar(128);primaryKey"`
	Version           string `json:"version" gorm:"column:version"`
	ChainName         string `json:"chainName" gorm:"column:chainName"`
	EnableGas         bool   `json:"enableGas" gorm:"column:enableGas"`
	Consensus         string `json:"consensus" gorm:"column:consensus"`
	TxTimestampVerify bool   `json:"txTimestampVerify" gorm:"column:txTimestampVerify"`
	TxTimeout         int    `json:"txTimeout" gorm:"column:txTimeout"`
	BlockTxCapacity   int    `json:"blockTxCapacity" gorm:"column:blockTxCapacity"`
	BlockSize         int    `json:"blockSize" gorm:"column:blockSize"`
	BlockInterval     int    `json:"blockInterval" gorm:"column:blockInterval"`
	HashType          string `json:"hashType" gorm:"column:hashType"`
	AuthType          string `json:"authType" gorm:"column:authType"`
	ChainConfig       string `json:"chainConfig" gorm:"column:chainConfig"`
	Timestamp         int64  `json:"timestamp" gorm:"column:timestamp"` //订阅时间
	CommonIntField
}

// TableName table
func (*Chain) TableName() string {
	return TableChain
}

// Org
// @Description: 组织数据
type Org struct {
	OrgId  string `json:"orgId" gorm:"type:varchar(128);column:orgId;primaryKey;comment:组织ID"`
	Status int    `json:"status" gorm:"column:status;comment:组织状态(0:正常 1: 已删掉)"` // 0:正常 1: 已删掉
	CommonIntField
}

// TableName table
func (t *Org) TableName() string {
	return TableOrg
}

// User
// @Description: 用户数据
type User struct {
	UserId    string `json:"userId" gorm:"column:userId;comment:用户ID"`
	UserAddr  string `json:"userAddr" gorm:"column:userAddr;type:varchar(128);primaryKey;comment:用户地址"`
	Role      string `json:"role" gorm:"column:role;comment:用户角色"`
	OrgId     string `json:"orgId" gorm:"column:orgId;comment:用户所属组织ID"`
	Timestamp int64  `json:"timestamp" gorm:"column:timestamp;comment:用户加入时间"`
	Status    int    `json:"status" gorm:"column:status;comment:用户状态(0:正常 1:已删除 2:禁用)"`
	CommonIntField
}

// TableName table
func (t *User) TableName() string {
	return TableUser
}

// Contract 通用合约
type Contract struct {
	Name             string `json:"name" gorm:"column:name;comment:合约名称(脱敏)"`                                         //合约名称
	NameBak          string `json:"nameBak" gorm:"column:nameBak;type:varchar(128);index;comment:合约真实名称;"`            //合约真实名称
	Addr             string `json:"addr" gorm:"column:addr;type:varchar(128);primaryKey;comment:合约地址;"`               //合约地址
	Version          string `json:"version" gorm:"column:version;comment:合约版本号"`                                      //合约版本
	RuntimeType      string `json:"runtimeType" gorm:"column:runtimeType;comment:合约运行时版本"`                            //运行时版本
	ContractStatus   int32  `json:"contractStatus" gorm:"column:contractStatus;comment:合约状态(-1:系统合约,0:正常,1:冻结,2:注销)"` // 合约状态，-1：系统合约；0：正常；1：冻结；2：注销
	ContractType     string `json:"contractType" gorm:"column:contractType;comment:合约类型"`
	ContractSymbol   string `json:"contractSymbol" gorm:"column:contractSymbol;comment:合约简称"`
	Decimals         int    `json:"decimals" gorm:"column:decimals;comment:合约小数位数"`                 //小数位数‘
	BlockHeight      int64  `json:"blockHeight" gorm:"column:blockHeight;comment:区块高度，只用作更新的版本号控制"` // 区块高度
	TxNum            int64  `json:"txNum" gorm:"column:txNum;comment:合约交易总量"`                       // 交易数据量
	EventNum         int64  `json:"eventNum" gorm:"column:eventNum;comment:合约事件总量"`                 // 合约事件总数
	OrgId            string `json:"orgId" gorm:"column:orgId;comment:合约部署用户所属组织"`                   //合约的发起组织
	CreateTxId       string `json:"createTxId" gorm:"column:createTxId;comment:创建合约的交易ID"`
	CreateSender     string `json:"createSender" gorm:"column:createSender;comment:创建合约的用户ID"`   //创建用户id
	CreatorAddr      string `json:"creatorAddr" gorm:"column:creatorAddr;comment:创建合约的用户地址"`     //创建用户地址
	UpgradeOrgId     string `json:"UpgradeOrgId" gorm:"column:UpgradeOrgId;comment:更新合约的用户所属组织"` //更新用户组织
	Upgrader         string `json:"upgrader" gorm:"column:upgrader;comment:更新合约的用户ID"`           //更新用户id
	UpgradeAddr      string `json:"upgradeAddr" gorm:"column:upgradeAddr;comment:更新合约的用户地址"`     //更新用户地址
	UpgradeTimestamp int64  `json:"upgradeTimestamp" gorm:"column:upgradeTimestamp;comment:合约更新时间"`
	Timestamp        int64  `json:"timestamp" gorm:"column:timestamp;index;comment:合约创建时间"`
	CommonIntField
}

// TableName table
func (t *Contract) TableName() string {
	return TableContract
}

// FungibleContract 同质化合约
type FungibleContract struct {
	Symbol          string          `json:"symbol" gorm:"column:symbol;comment:合约简称"`
	ContractName    string          `json:"contractName" gorm:"column:contractName;comment:合约名称(脱敏)"`
	ContractNameBak string          `json:"contractNameBak" gorm:"column:contractNameBak;type:varchar(128);index;comment:合约真实名称;"` //合约真实名称
	ContractAddr    string          `json:"contractAddr" gorm:"column:contractAddr;type:varchar(128);primaryKey;comment:合约地址;"`    //合约地址
	ContractType    string          `json:"contractType" gorm:"column:contractType;comment:合约类型"`                                  //合约类型
	TotalSupply     decimal.Decimal `json:"totalSupply" gorm:"column:totalSupply;type:decimal(50,18);comment:合约总发行量"`              //发行总量
	HolderCount     int64           `json:"holderCount" gorm:"column:holderCount;comment:合约总持有人数"`                                 //持仓人数
	TransferNum     int64           `json:"transferNum" gorm:"column:transferNum;comment:合约流转总数"`                                  //流转次数
	BlockHeight     int64           `json:"blockHeight" gorm:"column:blockHeight;comment:区块高度，只用作更新的版本号控制"`                        // 区块高度
	Timestamp       int64           `json:"timestamp" gorm:"column:timestamp;index;comment:合约创建时间;"`
	CommonIntField
}

// TableName table
func (c *FungibleContract) TableName() string {
	return TableFungibleContract
}

// NonFungibleContract 非同质化合约
type NonFungibleContract struct {
	ContractName    string          `json:"contractName" gorm:"column:contractName;comment:合约名称(脱敏)"`
	ContractNameBak string          `json:"contractNameBak" gorm:"column:contractNameBak;type:varchar(128);index;comment:合约真实名称;"` //合约真实名称
	ContractAddr    string          `json:"contractAddr" gorm:"column:contractAddr;type:varchar(128);primaryKey;comment:合约地址;"`    //合约地址
	ContractType    string          `json:"contractType" gorm:"column:contractType;comment:合约类型"`                                  //合约类型
	TotalSupply     decimal.Decimal `json:"totalSupply" gorm:"column:totalSupply;type:decimal(30,0);comment:合约总发行量"`               //发行总量
	HolderCount     int64           `json:"holderCount" gorm:"column:holderCount;comment:合约总持有人数"`                                 //持仓人数
	TransferNum     int64           `json:"transferNum" gorm:"column:transferNum;comment:合约流转总数"`                                  //流转次数
	BlockHeight     int64           `json:"blockHeight" gorm:"column:blockHeight;comment:区块高度，只用作更新的版本号控制"`                        // 区块高度
	Timestamp       int64           `json:"timestamp" gorm:"column:timestamp;index;comment:合约创建时间;"`
	CommonIntField
}

// TableName table
func (c *NonFungibleContract) TableName() string {
	return TableNonFungibleContract
}

// EvidenceContract 存证合约
type EvidenceContract struct {
	ID           string `json:"id" gorm:"primaryKey;comment:主键ID;"`
	TxId         string `json:"txId" gorm:"column:txId;type:varchar(128);uniqueIndex:,composite:txId_hash;comment:交易ID"`
	SenderAddr   string `json:"senderAddr" gorm:"column:senderAddr;comment:合约创建用户地址"`
	ContractName string `json:"contractName" gorm:"column:contractName;index;comment:合约名称;"`
	EvidenceId   string `json:"evidenceId" gorm:"column:evidenceId;comment:存证合约ID"`
	Hash         string `json:"hash" gorm:"column:hash;type:varchar(128);index;uniqueIndex:,composite:txId_hash;comment:存证合约哈希"`
	//Hash               string `json:"hash" gorm:"column:hash;index;comment:存证合约哈希"`
	MetaData           string `json:"metaData" gorm:"column:metaData;comment:存证合约数据"`
	MetaDataBak        string `json:"metaDataBak" gorm:"column:metaDataBak;comment:存证合约数据"`
	BlockHeight        int64  `json:"blockHeight" gorm:"column:blockHeight;comment:区块高度"`                 // 区块高度
	ContractResult     []byte `json:"contractResult" gorm:"column:contractResult;comment:合约执行结果"`         //合约结果
	ContractResultCode uint32 `json:"contractResultCode" gorm:"column:contractResultCode;comment:合约执行结果"` //合约结果码
	Timestamp          int64  `json:"timestamp" gorm:"column:timestamp;index;comment:上链时间;"`
	CommonIntField
}

// TableName table
func (t *EvidenceContract) TableName() string {
	return TableEvidenceContract
}

// IdentityContract 身份认证合约
type IdentityContract struct {
	ID           string `json:"id" gorm:"primaryKey;autoIncrement;comment:主键ID"`
	TxId         string `json:"txId" gorm:"column:txId;type:varchar(128);uniqueIndex:,composite:identity_txId_index;comment:交易ID;"`
	EventIndex   int    `json:"eventIndex" gorm:"column:eventIndex;uniqueIndex:,composite:identity_txId_index;comment:同一个交易内的合约顺序;"`
	ContractName string `json:"contractName" gorm:"column:contractName;type:varchar(128);index;comment:合约名称;"`
	ContractAddr string `json:"contractAddr" gorm:"column:contractAddr;type:varchar(128);index;comment:合约地址;"` //合约地址
	UserAddr     string `json:"userAddr" gorm:"column:userAddr;comment:合约创建者地址"`                               //人员地址
	Level        string `json:"level" gorm:"column:level;comment:身份信息"`                                        //身份信息
	PkPem        string `json:"pkPem" gorm:"column:pkPem;comment:公钥"`                                          //公钥
	CommonIntField
}

// TableName table
func (t *IdentityContract) TableName() string {
	return TableIdentityContract
}

// FungibleTransfer 同质化流转
type FungibleTransfer struct {
	ID             string          `json:"id" gorm:"primaryKey;comment:主键ID"`
	TxId           string          `json:"txId" gorm:"column:txId;type:varchar(128);uniqueIndex:,composite:ft_transfer_txId_index;comment:交易ID"`
	EventIndex     int             `json:"eventIndex" gorm:"column:eventIndex;uniqueIndex:,composite:ft_transfer_txId_index;comment:同一个交易内流转顺序"`
	ContractName   string          `json:"contractName" gorm:"column:contractName;type:varchar(128);index;comment:合约名称;"`
	ContractAddr   string          `json:"contractAddr" gorm:"column:contractAddr;type:varchar(128);index:,composite:ft_transfer_addr_time;index:,composite:ft_transfer_contract_from;index:,composite:ft_transfer_contract_to;comment:合约地址"`
	ContractMethod string          `json:"contractMethod" gorm:"column:contractMethod;type:varchar(512);comment:合约执行方法;"`
	Topic          string          `json:"topic" gorm:"column:topic;comment:合约事件topic"`
	FromAddr       string          `json:"fromAddr" gorm:"column:fromAddr;type:varchar(128);index:,composite:ft_transfer_contract_from;comment:交易流转from地址"`
	ToAddr         string          `json:"toAddr" gorm:"column:toAddr;type:varchar(128);index:,composite:ft_transfer_contract_to;comment:交易流转to地址"`
	Amount         decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(50,18);comment:交易流转数量"` //持有数量
	Timestamp      int64           `json:"timestamp" gorm:"column:timestamp;index:,composite:ft_transfer_addr_time;comment:上链时间"`
	CommonIntField
}

// TableName table
func (t *FungibleTransfer) TableName() string {
	return TableFungibleTransfer
}

// NonFungibleTransfer 非同质化交易流转
type NonFungibleTransfer struct {
	ID             string `json:"id" gorm:"primaryKey;comment:主键ID"`
	TxId           string `json:"txId" gorm:"column:txId;type:varchar(128);uniqueIndex:,composite:nft_transfer_txId;comment:交易ID"`
	EventIndex     int    `json:"eventIndex" gorm:"column:eventIndex;uniqueIndex:,composite:nft_transfer_txId;comment:同一个交易内流转顺序"`
	ContractName   string `json:"contractName" gorm:"column:contractName;type:varchar(128);comment:合约名称"`
	ContractAddr   string `json:"contractAddr" gorm:"column:contractAddr;type:varchar(128);index:,composite:nft_transfer_addr_time;index:,composite:nft_transfer_contract_from;index:,composite:nft_transfer_contract_to;comment:合约地址"`
	ContractMethod string `json:"contractMethod" gorm:"column:contractMethod;type:varchar(512);comment:合约执行方法;"`
	Topic          string `json:"topic" gorm:"column:topic;comment:合约事件topic"`
	FromAddr       string `json:"fromAddr" gorm:"column:fromAddr;type:varchar(128);index:,composite:nft_transfer_contract_from;comment:交易流转from地址"`
	ToAddr         string `json:"toAddr" gorm:"column:toAddr;type:varchar(128);index:,composite:nft_transfer_contract_to;comment:交易流转to地址"`
	TokenId        string `json:"tokenId" gorm:"column:tokenId;type:varchar(128);index;comment:交易流转token;"`
	Timestamp      int64  `json:"timestamp" gorm:"column:timestamp;index:,composite:nft_transfer_addr_time;comment:上链时间;"`
	CommonIntField
}

// TableName table
func (t *NonFungibleTransfer) TableName() string {
	return TableNonFungibleTransfer
}

// FungiblePosition 同质化持仓信息
type FungiblePosition struct {
	ID           string          `json:"id" gorm:"primaryKey;comment:主键ID"`
	ContractAddr string          `json:"contractAddr" gorm:"column:contractAddr;type:varchar(128);uniqueIndex:,composite:contract_owner;index:,composite:addr_amount_index;comment:持有合约地址"` //持有合约地址
	ContractName string          `json:"contractName" gorm:"column:contractName;comment:持有合约名称"`                                                                                            //持有合约名称
	OwnerAddr    string          `json:"ownerAddr" gorm:"column:ownerAddr;type:varchar(128);uniqueIndex:,composite:contract_owner;index;comment:持仓地址;"`                                     //持仓地址
	Symbol       string          `json:"symbol" gorm:"column:symbol;comment:合约简称"`                                                                                                          //合约简称
	Amount       decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(50,18);index:,composite:addr_amount_index;comment:持有数量;"`                                                  //持有数量
	BlockHeight  int64           `json:"blockHeight" gorm:"column:blockHeight;comment:区块高度，只用作更新的版本号控制"`                                                                                    // 区块高度
	CommonIntField
}

// TableName table
func (p *FungiblePosition) TableName() string {
	return TableFungiblePosition
}

// NonFungiblePosition 非同质化持仓信息
type NonFungiblePosition struct {
	ID           string          `json:"id" gorm:"primaryKey;comment:主键ID"`
	ContractAddr string          `json:"contractAddr" gorm:"column:contractAddr;type:varchar(128);uniqueIndex:,composite:contract_owner;index:,composite:addr_amount;comment:持有合约地址"` //持有合约地址
	ContractName string          `json:"contractName" gorm:"column:contractName;comment:持有合约名称"`                                                                                      //持有合约名称
	OwnerAddr    string          `json:"ownerAddr" gorm:"column:ownerAddr;type:varchar(128);uniqueIndex:,composite:contract_owner;index;comment:持仓地址;"`                               //持仓地址
	Amount       decimal.Decimal `json:"amount" gorm:"column:amount;type:decimal(30,0);index:,composite:addr_amount;comment:持有数量"`                                                    //持有数量
	BlockHeight  int64           `json:"blockHeight" gorm:"column:blockHeight;comment:区块高度，只用作更新的版本号控制"`                                                                              // 区块高度
	CommonIntField
}

// TableName table
func (p *NonFungiblePosition) TableName() string {
	return TableNonFungiblePosition
}

// NonFungibleToken 非同质化Token列表
type NonFungibleToken struct {
	ID           string `json:"id" gorm:"primaryKey;comment:主键ID;"`
	TokenId      string `json:"tokenId" gorm:"column:tokenId;type:varchar(128);uniqueIndex:,composite:tokenId_contractAddr;comment:token编号"`                                                                                       //持有token
	OwnerAddr    string `json:"ownerAddr" gorm:"column:ownerAddr;type:varchar(128);index:,composite:owner_time;index:,composite:owner_contractAddr_time;comment:持有token的账户;"`                                                      //持仓地址
	ContractAddr string `json:"contractAddr" gorm:"column:contractAddr;type:varchar(128);uniqueIndex:,composite:tokenId_contractAddr;index:,composite:contract_time;index:,composite:owner_contractAddr_time;comment:token所属合约地址"` //持有合约地址
	ContractName string `json:"contractName" gorm:"column:contractName;comment:token所属合约名称"`                                                                                                                                       //持有合约名称
	MetaData     string `json:"metaData" gorm:"column:metaData;comment:token附属信息"`                                                                                                                                                 //metaData
	MetaDataBak  string `json:"metaDataBak" gorm:"column:metaDataBak;comment:token附属信息"`
	CategoryName string `json:"categoryName" gorm:"column:categoryName;comment:token分类信息"`
	Timestamp    int64  `json:"timestamp" gorm:"column:timestamp;index:,composite:owner_time;index:,composite:contract_time;index:,composite:owner_contractAddr_time;comment:上链时间"`
	CommonIntField
}

// TableName table
func (t *NonFungibleToken) TableName() string {
	return TableNonFungibleToken
}

// Account 账户列表
type Account struct {
	//AddrID      string `json:"addrId" gorm:"column:addrId;comment:document中设置的id;"`
	Address     string `json:"address" gorm:"column:address;type:varchar(128);primaryKey;comment:账户地址"` //账户地址
	AddrType    int    `json:"addrType" gorm:"column:addrType;comment:账户类型,0:节点地址, 1:合约地址"`             //账户类型User/Contract
	TxNum       int64  `json:"txNum" gorm:"column:txNum;comment:账户发起的交易总量"`                             // 交易数据量
	NFTNum      int64  `json:"nftNum" gorm:"column:nftNum;comment:账户持有nft数量"`                           // 持有的nft数量
	BlockHeight int64  `json:"blockHeight" gorm:"column:blockHeight;comment:区块高度，只用作更新的版本号控制"`          // 区块高度
	DID         string `json:"did" gorm:"column:did;type:varchar(128);index;comment:账户DID;"`            //账户地址
	BNS         string `json:"bns" gorm:"column:bns;type:varchar(128);index;comment:账户BNS;"`            //账户地址
	CommonIntField
}

// TableName table
func (t *Account) TableName() string {
	return TableAccount
}

// CrossSubChainData 子链数据
type CrossSubChainData struct {
	SubChainId  string `json:"subChainId" gorm:"column:subChainId;primaryKey;comment:子链id标识"`             //子链ID(自动生成)
	GatewayId   string `json:"gatewayId" gorm:"column:gatewayId;type:varchar(128);index;comment:跨链网关ID;"` //跨链网关id
	BlockHeight int64  `json:"blockHeight" gorm:"column:blockHeight;comment:子链区块高度"`
	IsMainChain bool   `json:"isMainChain" gorm:"column:isMainChain;comment:是否为主链"`
	Timestamp   int64  `json:"timestamp" gorm:"column:timestamp;index;comment:上链时间;"` //上链时间
	TxNum       int64  `json:"txNum" gorm:"column:txNum;comment:子链交易总量"`              // 交易数据量
	CommonIntField
}

// TableName table
func (t *CrossSubChainData) TableName() string {
	return TableCrossSubChainData
}

// CrossChainContract 跨链合约
type CrossChainContract struct {
	ID           string `json:"id" gorm:"primaryKey;comment:主键ID;"`
	SubChainId   string `json:"subChainId" gorm:"column:subChainId;type:varchar(128);uniqueIndex:,composite:contract_sub_contract;comment:子链id标识"`       //子链ID
	ContractName string `json:"contractName" gorm:"column:contractName;type:varchar(128);uniqueIndex:,composite:contract_sub_contract;comment:子链跨链合约名称"` //合约名称
	CommonIntField
}

// TableName table
func (t *CrossChainContract) TableName() string {
	return TableCrossChainContract
}

// CrossSubChainCrossChain 子链跨链数据统计
type CrossSubChainCrossChain struct {
	ID           string `json:"id" gorm:"primaryKey;comment:主键ID"`
	SubChainId   string `json:"subChainId" gorm:"column:subChainId;type:varchar(128);index;comment:子链ID;"`
	CrossChainId string `json:"crossChainId" gorm:"column:crossChainId;comment:跨链ID;"`
	TxNum        int64  `json:"txNum" gorm:"column:txNum;index;comment:发生的跨链交易数;"`
	BlockHeight  int64  `json:"blockHeight" gorm:"column:blockHeight;comment:区块高度，只用作更新的版本号控制"`
	CommonIntField
}

// TableName table
func (t *CrossSubChainCrossChain) TableName() string {
	return TableCrossSubChainCrossChain
}

// CrossMainTransaction 跨链交易主链交易
type CrossMainTransaction struct {
	TxId        string `json:"txId" gorm:"column:txId;type:varchar(128);primaryKey;comment:主链交易ID"`
	CrossId     string `json:"crossId" gorm:"column:crossId;type:varchar(128);index;comment:跨链ID;"`      //跨链ID
	ChainMsg    string `json:"chainMsg" gorm:"column:chainMsg;comment:跨链内容"`                             // 跨链内容
	Status      int32  `json:"status" gorm:"column:status;comment:跨链状态(0:新建,1:待执行,2:待提交,3:确认结束,4:回滚结束)"` //跨链状态（0:新建，1：待执行，2:待提交，3:确认结束，4:回滚结束）
	CrossModel  int32  `json:"crossModel" gorm:"column:crossModel;comment:跨链类型（1:sync，0:other）"`         //跨链类型
	Timestamp   int64  `json:"timestamp" gorm:"column:timestamp;index;comment:上链时间;"`                    //交易时间
	BlockHeight int64  `json:"blockHeight" gorm:"column:blockHeight;comment:区块高度"`
	CommonIntField
}

// TableName table
func (t *CrossMainTransaction) TableName() string {
	return TableCrossMainTransaction
}

// CrossTransactionTransfer 跨链交易流转
type CrossTransactionTransfer struct {
	ID              string `json:"id" gorm:"primaryKey;comment:主键ID"`
	CrossId         string `json:"crossId" gorm:"column:crossId;type:varchar(128);uniqueIndex;comment:跨链ID;"`      //跨链ID
	TxId            string `json:"txId" gorm:"column:txId;type:varchar(128);index;comment:跨链交易ID;"`                //跨链交易ID
	FromGatewayId   string `json:"fromGatewayId" gorm:"column:fromGatewayId;comment:跨链发起网关"`                       //跨链发起网关
	FromChainId     string `json:"fromChainId" gorm:"column:fromChainId;type:varchar(128);index;comment:跨链发起链ID;"` //跨链发起链ID
	FromIsMainChain bool   `json:"fromIsMainChain" gorm:"column:fromIsMainChain;comment:是否是主链"`                    //是否是主链（true:主链，false:子链）
	FromBlockHeight int64  `json:"fromBlockHeight" gorm:"column:fromBlockHeight;comment:链高度"`                      //主链高度
	ToGatewayId     string `json:"toGatewayId" gorm:"column:toGatewayId;comment:跨链目标网关ID"`                         //跨链目标网关
	ToChainId       string `json:"toChainId" gorm:"column:toChainId;type:varchar(128);index;comment:跨链目标链ID;"`     //跨链目标链ID
	ToIsMainChain   bool   `json:"toIsMainChain" gorm:"column:toIsMainChain;comment:是否是主链"`                        //是否是主链（true:主链，false:子链）
	ToBlockHeight   int64  `json:"toBlockHeight" gorm:"column:toBlockHeight;comment:链高度"`                          //主链高度
	BlockHeight     int64  `json:"blockHeight" gorm:"column:blockHeight;comment:主链高度"`                             //主链高度
	ContractName    string `json:"contractName" gorm:"column:contractName;comment:合约名称"`                           //合约名称
	ContractMethod  string `json:"contractMethod" gorm:"column:contractMethod;type:varchar(512);comment:合约方法;"`    //合约方法
	Parameter       string `json:"parameter" gorm:"column:parameter;comment:合约参数"`                                 //合约参数
	TxNum           int64  `json:"txNum" gorm:"column:txNum;index;comment:发生的跨链交易数;"`
	CrossModel      int32  `json:"crossModel" gorm:"column:crossModel;comment:跨链类型（1:sync，0:other）"`                    //跨链类型
	Status          int32  `json:"status" gorm:"column:status;comment:跨链状态(0:新建,1:待执行,2:待提交,3:确认结束,4:回滚结束)"`            //跨链状态（0:新建，1：待执行，2:待提交，3:确认结束，4:回滚结束）
	StartTime       int64  `json:"startTime" gorm:"column:startTime;index:,composite:start_end_index;comment:跨链周期开始时间"` //跨链周期开始时间
	EndTime         int64  `json:"endTime" gorm:"column:endTime;index:,composite:start_end_index;comment:跨链周期结束时间"`     //跨链周期结束时间
	CommonIntField
}

// TableName table
func (t *CrossTransactionTransfer) TableName() string {
	return TableCrossTransactionTransfer
}

// CrossBusinessTransaction 跨链执行的业务交易
type CrossBusinessTransaction struct {
	ID                  string `json:"id" gorm:"primaryKey;comment:主键ID"`
	TxId                string `json:"txId" gorm:"column:txId;type:varchar(128);comment:主链/子链交易ID;"`
	CrossId             string `json:"crossId" gorm:"column:crossId;type:varchar(128);uniqueIndex:,composite:cross_business_tx_id_sub;comment:跨链ID;"`                  //跨链ID
	SubChainId          string `json:"subChainId" gorm:"column:subChainId;type:varchar(128);uniqueIndex:,composite:cross_business_tx_id_sub;index;comment:跨链主链/子链id;"` // 跨链主，子链id
	IsMainChain         bool   `json:"isMainChain" gorm:"column:isMainChain;comment:是否是主链"`                                                                            //是否是主链（true:主链，false:子链）
	GatewayId           string `json:"gatewayId" gorm:"column:gatewayId;comment:目标网关id"`                                                                               // 目标网关id
	TxStatus            int32  `json:"txStatus" gorm:"column:txStatus;comment:跨链交易执行结果"`                                                                               //跨链解析结果
	CrossContractResult string `json:"crossContractResult" gorm:"column:crossContractResult;comment:跨链合约执行结果"`                                                         //合约结果
	TxType              string `json:"txType" gorm:"column:txType;comment:交易类型"`
	Timestamp           int64  `json:"timestamp" gorm:"column:timestamp;index;comment:上链时间;"`
	TxStatusCode        string `json:"txStatusCode" gorm:"column:txStatusCode;comment:跨链交易执行结果"` //交易解析结果
	RwSetHash           string `json:"rwSetHash" gorm:"column:rwSetHash;comment:交易读写集"`
	ContractResultCode  uint32 `json:"contractResultCode" gorm:"column:contractResultCode;comment:交易合约执行结果"`
	ContractResult      []byte `json:"contractResult" gorm:"column:contractResult;comment:交易合约执行结果"`
	ContractMessage     string `json:"contractMessage" gorm:"column:contractMessage;comment:合约外返回的错误信息"`
	ContractName        string `json:"contractName" gorm:"column:contractName;type:varchar(128);index;comment:交易合约名称;"`
	ContractMethod      string `json:"contractMethod" gorm:"column:contractMethod;comment:交易执行合约方法"`
	ContractParameters  string `json:"contractParameters" gorm:"column:contractParameters;comment:交易执行合约参数"`
	GasUsed             uint64 `json:"gasUsed" gorm:"column:gasUsed;comment:交易花费gas数"`
	CommonIntField
}

// TableName table
func (t *CrossBusinessTransaction) TableName() string {
	return TableCrossBusinessTransaction
}

// IDAContract IDA合约
type IDAContract struct {
	ContractName      string `json:"contractName" gorm:"column:contractName;comment:IDA合约名称"`
	ContractNameBak   string `json:"contractNameBak" gorm:"column:contractNameBak;type:varchar(128);index;comment:IDA合约真实名称;"` //合约真实名称
	ContractAddr      string `json:"contractAddr" gorm:"column:contractAddr;type:varchar(128);primaryKey;comment:合约地址;"`       //合约地址
	ContractType      string `json:"contractType" gorm:"column:contractType;comment:合约类型"`                                     //合约类型
	TotalNormalAssets int64  `json:"totalNormalAssets" gorm:"column:totalNormalAssets;comment:正常数据资产总量"`
	TotalAssets       int64  `json:"totalAssets" gorm:"column:totalAssets;comment:数据资产总量"`
	BlockHeight       int64  `json:"blockHeight" gorm:"column:blockHeight;comment:区块高度，只用作更新的版本号控制"` //区块高度
	Timestamp         int64  `json:"timestamp" gorm:"column:timestamp;index;comment:创建时间;"`          //创建时间
	CommonIntField
}

// TableName table
func (t *IDAContract) TableName() string {
	return TableIDAContract
}

// IDAAssetDetail IDA资产详情
type IDAAssetDetail struct {
	ID                string    `json:"id" gorm:"primaryKey;comment:主键ID"`
	AssetCode         string    `json:"assetCode" gorm:"column:assetCode;primaryKey;comment:资产编号"` //资产编号
	ContractName      string    `json:"contractName" gorm:"column:contractName;comment:IDA合约名称"`
	ContractAddr      string    `json:"contractAddr" gorm:"column:contractAddr;comment:IDA合约地址;"`                               //合约地址
	AssetName         string    `json:"assetName" gorm:"column:assetName;comment:资产名称;"`                                        //资产名称
	AssetEnName       string    `json:"assetEnName" gorm:"column:assetEnName;comment:资产英文名称;"`                                  //资产英文名称
	Category          int       `json:"category" gorm:"column:category;comment:资产类型(1:数据集,2:API服务,3:数据报告,4:数据应用,5:计算模型)"`       //资产类型、1: 数据集, 2: API服务, 3: 数据报告, 4: 数据应用, 5: 计算模型
	ImmediatelySupply bool      `json:"immediatelySupply" gorm:"column:immediatelySupply;comment:资产供应方式(true:及时供应,false:延迟供应)"` //供应方式：0:及时供应，1：延迟供应
	SupplyTime        time.Time `json:"supplyTime" gorm:"column:supplyTime;default:'1000-01-01 00:00:00.000';comment:延迟供应时间"`   //延迟供应时间
	DataScale         string    `json:"dataScale" gorm:"column:dataScale;comment:数据规模(1M,1G,1条)"`                               //数据规模：1M，1G，1条
	IndustryTitle     string    `json:"industryTitle" gorm:"column:industryTitle;comment:行业标题"`                                 //行业标题
	Summary           string    `json:"summary" gorm:"column:summary;comment:资产介绍"`                                             //资产介绍
	Creator           string    `json:"creator" gorm:"column:creator;comment:资产创建人"`                                            //创建人
	Holder            string    `json:"holder" gorm:"column:holder;comment:资产持有人"`                                              //持有人
	TxID              string    `json:"txId" gorm:"column:txId;comment:交易ID"`                                                   //交易ID
	UserCategories    string    `json:"userCategories" gorm:"column:userCategories;comment:资产使用对象"`                             //使用对象
	UpdateCycleType   int       `json:"updateCycleType" gorm:"column:updateCycleType;comment:更新周期类型(1:静态,2:实时,3:周期,4:其他)"`      //更新周期类型、1: 静态, 2: 实时, 3: 周期, 4：其他
	UpdateTimeSpan    string    `json:"updateTimeSpan" gorm:"column:updateTimeSpan;comment:更新周期时间跨度(1分钟,1天)"`                   //更新周期时间跨度、1分钟，1天
	CreatedTime       int64     `json:"createdTime" gorm:"column:createdTime;index;comment:资产创建时间;"`
	UpdatedTime       int64     `json:"updatedTime" gorm:"column:updatedTime;index;comment:资产更新时间;"`
	IsDeleted         bool      `json:"isDeleted" gorm:"column:isDeleted;comment:资产状态(0:正常,1:已删除)"` //资产状态：0:正常，1:已删除
	CommonIntField
}

// TableName table
func (t *IDAAssetDetail) TableName() string {
	return TableIDAAssetDetail
}

// IDAAssetAttachment IDA资产附件信息
type IDAAssetAttachment struct {
	ID          string `json:"id" gorm:"primaryKey;comment:主键ID"`
	AssetCode   string `json:"assetCode" gorm:"column:assetCode;comment:资产编号"`                                      //资产编号
	Url         string `json:"url" gorm:"column:url;comment:材料url;"`                                                //材料url
	ContextType int    `json:"contextType" gorm:"column:contextType;comment:材料类型(1:图片,2:合规证明材料,3:估值证明材料,4:其他相关附件)"` //材料类型，1: 图片, 2: 合规证明材料, 3: 估值证明材料, 4: 其他相关附件
	CommonIntField
}

// TableName table
func (t *IDAAssetAttachment) TableName() string {
	return TableIDAAssetAttachment
}

// IDADataAsset data资产数据
type IDADataAsset struct {
	ID           string `json:"id" gorm:"primaryKey;comment:主键ID"`
	AssetCode    string `json:"assetCode" gorm:"column:assetCode;comment:资产编号"`          //资产编号
	FieldName    string `json:"fieldName" gorm:"column:fieldName;comment:字段名称;"`         //字段名称
	FieldType    string `json:"fieldType" gorm:"column:fieldType;comment:字段类型;"`         //字段类型
	FieldLength  int    `json:"fieldLength" gorm:"column:fieldLength;comment:字段长度;"`     //字段长度
	IsPrimaryKey int    `json:"isPrimaryKey" gorm:"column:isPrimaryKey;comment:是否主键;"`   // 是否主键
	IsNotNull    int    `json:"isNotNull" gorm:"column:isNotNull;comment:是否非空;"`         // 是否非空
	PrivacyQuery int    `json:"privacyQuery" gorm:"column:privacyQuery;comment:是否隐私计算;"` // 是否隐私计算
	CommonIntField
}

// TableName table
func (t *IDADataAsset) TableName() string {
	return TableIDADataAsset
}

// IDAApiAsset api资产数据
type IDAApiAsset struct {
	ID           string `json:"id" gorm:"primaryKey;comment:主键ID"`
	AssetCode    string `json:"assetCode" gorm:"column:assetCode;comment:资产编号"`        //资产编号
	Header       string `json:"header" gorm:"column:header;comment:api-header;"`       //header
	Url          string `json:"url" gorm:"column:url;comment:api的url地址;"`              //url
	Params       string `json:"params" gorm:"column:params;comment:请求参数;"`             //请求参数
	Response     string `json:"response" gorm:"column:response;comment:返回参数;"`         //返回参数
	Method       string `json:"method" gorm:"column:method;comment:请求类型;"`             //请求类型
	ResponseType string `json:"responseType" gorm:"column:responseType;comment:格式类型;"` //格式类型
	CommonIntField
}

// TableName table
func (t *IDAApiAsset) TableName() string {
	return TableIDAApiAsset
}

// DIDDetail	did信息
type DIDDetail struct {
	ID             string `json:"id" gorm:"primaryKey;comment:主键ID"`
	ContractName   string `json:"contractName" gorm:"column:contractName;type:varchar(128);comment:合约名称"`
	ContractAddr   string `json:"contractAddr" gorm:"column:contractAddr;type:varchar(128);uniqueIndex:,composite:did_contract;index:,composite:addr_time;comment:合约地址"`
	DID            string `json:"did" gorm:"column:did;type:varchar(128);uniqueIndex:,composite:did_contract;index;comment:did"`
	Status         int    `json:"status" gorm:"column:status;comment:状态，0：正常，1:冻结"`
	Document       string `json:"document" gorm:"column:document;comment:did文档"`
	IsIssuer       bool   `json:"isIssuer" gorm:"column:isIssuer;comment:是否颁证机构"`
	IssuerService  string `json:"issuerService" gorm:"column:issuerService;comment:issuerService;"`
	Authentication string `json:"authentication" gorm:"column:authentication;comment:authentication;"`
	AccountJson    string `json:"accountJson" gorm:"column:accountJson;comment:绑定账户;"`
	CreateTxId     string `json:"createTxId" gorm:"column:createTxId;comment:创建交易ID"`
	Timestamp      int64  `json:"timestamp" gorm:"column:timestamp;comment:上链时间;index;index:,composite:addr_time;"`
	CommonIntField
}

// TableName table
func (t *DIDDetail) TableName() string {
	return TableDIDDetail
}

// DIDSetHistory did变更历史
type DIDSetHistory struct {
	ID             string `json:"id" gorm:"primaryKey;comment:主键ID"`
	DID            string `json:"did" gorm:"column:did;index:,composite:did_time;comment:did"`
	ContractName   string `json:"contractName" gorm:"column:contractName;type:varchar(128);comment:合约名称"`
	ContractAddr   string `json:"contractAddr" gorm:"column:contractAddr;type:varchar(128);comment:合约地址"`
	Document       string `json:"document" gorm:"column:document;comment:did文档"`
	IssuerService  string `json:"issuerService" gorm:"column:issuerService;comment:issuerService;"`
	Authentication string `json:"authentication" gorm:"column:authentication;comment:authentication;"`
	AccountJson    string `json:"accountJson" gorm:"column:accountJson;comment:绑定账户;"`
	Timestamp      int64  `json:"timestamp" gorm:"column:timestamp;index:,composite:did_time;comment:上链时间;index;"`
	TxId           string `json:"txId" gorm:"column:txId;comment:交易ID"`
	CommonIntField
}

// TableName table
func (t *DIDSetHistory) TableName() string {
	return TableDIDSetHistory
}

// VCIssueHistory vc变更历史
type VCIssueHistory struct {
	ID           string `json:"id" gorm:"primaryKey;comment:主键ID"`
	ContractName string `json:"contractName" gorm:"column:contractName;comment:合约名称"`
	ContractAddr string `json:"contractAddr" gorm:"column:contractAddr;type:varchar(128);uniqueIndex:,composite:vc_contract;comment:合约地址"`
	VCID         string `json:"vcId" gorm:"column:vcId;type:varchar(256);uniqueIndex:,composite:vc_contract;comment:vcId"`
	IssuerDID    string `json:"issuerDID" gorm:"column:issuerDID;index:,composite:issuer_time;comment:颁发者did"`
	HolderDID    string `json:"holderDID" gorm:"column:holderDID;index:,composite:hold_time;comment:持证者did"`
	TemplateID   string `json:"templateId" gorm:"column:templateId;index:,composite:template_time;comment:模板id"`
	Timestamp    int64  `json:"timestamp" gorm:"column:timestamp;index:,composite:issuer_time;index:,composite:hold_time;index:,composite:template_time;comment:上链时间;"`
	Status       int    `json:"status" gorm:"column:status;comment:状态，0：正常，1:冻结"`
	TxId         string `json:"txId" gorm:"column:txId;comment:交易ID"`
	CommonIntField
}

// TableName table
func (t *VCIssueHistory) TableName() string {
	return TableVCIssueHistory
}

// VCTemplate vc模板
type VCTemplate struct {
	ID           string `json:"id" gorm:"primaryKey;comment:主键ID"`
	TemplateID   string `json:"templateId" gorm:"column:templateId;type:varchar(256);uniqueIndex:,composite:template_version_contract;comment:模板id"`
	TemplateName string `json:"templateName" gorm:"column:templateName;comment:模板名称"`
	ShortName    string `json:"shortName" gorm:"column:shortName;comment:简称"`
	ContractName string `json:"contractName" gorm:"column:contractName;type:varchar(128);index;comment:合约名称"`
	ContractAddr string `json:"contractAddr" gorm:"column:contractAddr;type:varchar(128);uniqueIndex:,composite:template_version_contract;comment:合约地址"`
	VCType       string `json:"vcTypre" gorm:"column:vcTypre;comment:vc类型"`
	Version      string `json:"version" gorm:"column:version;type:varchar(256);uniqueIndex:,composite:template_version_contract;comment:版本号"`
	Template     string `json:"template" gorm:"column:template;comment:vc模板内容"`
	Timestamp    int64  `json:"timestamp" gorm:"column:timestamp;comment:上链时间;index;"`
	TxId         string `json:"txId" gorm:"column:txId;comment:交易ID"`
	CreateDID    string `json:"createDID" gorm:"column:createDID;comment:创建did"`
	CommonIntField
}

// TableName table
func (t *VCTemplate) TableName() string {
	return TableVCTemplate
}

// UserToken 用户token表
type LoginUserToken struct {
	Id         string    `json:"id" gorm:"column:id;primaryKey"`
	UserAddr   string    `json:"userAddr" gorm:"column:userAddr;"`
	Token      string    `json:"token" gorm:"column:token;"`
	Sign       string    `json:"sign" gorm:"column:sign;"`
	ExpireTime time.Time `json:"expireTime" gorm:"column:expireTime"`
	CommonIntField
}

// TableName table
func (t *LoginUserToken) TableName() string {
	return TableLoginUserToken
}

type ContractABIFile struct {
	Id              string `json:"id" gorm:"column:id;primaryKey"`
	ContractName    string `json:"contractName" gorm:"column:contractName;type:varchar(128);comment:合约名称"`
	ContractAddr    string `json:"contractAddr" gorm:"column:contractAddr;type:varchar(128);uniqueIndex:,composite:addr_version;comment:合约地址"`
	ContractVersion string `json:"contractVersion" gorm:"column:contractVersion;type:varchar(128);uniqueIndex:,composite:addr_version;comment:合约版本号"`
	ABIJson         string `json:"abiJson" gorm:"column:abiJson;comment:合约abi"`
	CommonIntField
}

// TableName table
func (t *ContractABIFile) TableName() string {
	return TableContractABIFile
}

type ContractABITopic struct {
	Id              string `json:"id" gorm:"column:id;primaryKey"`
	ContractName    string `json:"contractName" gorm:"column:contractName;type:varchar(128);comment:合约名称"`
	ContractAddr    string `json:"contractAddr" gorm:"column:contractAddr;type:varchar(128);uniqueIndex:,composite:addr_version_topic;comment:合约地址"`
	ContractVersion string `json:"contractVersion" gorm:"column:contractVersion;type:varchar(128);uniqueIndex:,composite:addr_version_topic;comment:合约版本号"`
	Topic           string `json:"topic" gorm:"column:topic;type:varchar(128);uniqueIndex:,composite:addr_version_topic;comment:合约topic"`
	TopicTableName  string `json:"topicTableName" gorm:"column:topicTableName;comment:合约topic解析表名"`
	CommonIntField
}

// TableName table
func (t *ContractABITopic) TableName() string {
	return TableContractABITopic
}

type ContractCrossCallTransaction struct {
	Id               string `json:"id" gorm:"column:id;primaryKey"`
	TxId             string `json:"txId" gorm:"column:txId;comment:交易ID"`
	BlockHeight      int64  `json:"blockHeight" gorm:"column:blockHeight;comment:区块高度"`
	InvokingContract string `json:"invokingContract" gorm:"column:invokingContract;type:varchar(128);index:,composite:invoking_target_isCross;index:,composite:invoking_isCross;comment:调用合约名称"`
	InvokingMethod   string `json:"invokingMethod" gorm:"column:invokingMethod;comment:调用合约方法"`
	TargetContract   string `json:"targetContract" gorm:"column:targetContract;index;index:,composite:invoking_target_isCross;index:,composite:target_isCross;comment:被调用合约名称"`
	UserAddr         string `json:"userAddr" gorm:"column:userAddr;index;comment:调用者地址"`
	IsCross          bool   `json:"isCross" gorm:"column:isCross;index:,composite:invoking_target_isCross;index:,composite:invoking_isCross;index:,composite:target_isCross;comment:是否跨合约调用"`
	Timestamp        int64  `json:"timestamp" gorm:"column:timestamp;index;comment:上链时间;"`
	CommonIntField
}

// TableName table
func (t *ContractCrossCallTransaction) TableName() string {
	return TableContractCrossCallTransaction
}

type ContractCrossCall struct {
	Id               string `json:"id" gorm:"column:id;primaryKey"`
	InvokingContract string `json:"invokingContract" gorm:"column:invokingContract;type:varchar(128);index;uniqueIndex:,composite:invok_method_target;comment:调用合约名称"`
	InvokingMethod   string `json:"invokingMethod" gorm:"column:invokingMethod;type:varchar(128);index;uniqueIndex:,composite:invok_method_target;comment:调用合约方法"`
	TargetContract   string `json:"targetContract" gorm:"column:targetContract;type:varchar(128);index;uniqueIndex:,composite:invok_method_target;comment:被调用合约名称"`
	CommonIntField
}

// TableName table
func (t *ContractCrossCall) TableName() string {
	return TableContractCrossCall
}
