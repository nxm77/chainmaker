/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package common

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/sync/model"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"os"
	"regexp"
	"strings"

	"chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	tcipCommon "chainmaker.org/chainmaker/tcip-go/v2/common"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/shopspring/decimal"
)

const (
	// ContractStandardNameOTHER - Contract Standard - other Contract
	ContractStandardNameOTHER = "OTHER"
	//ContractStandardNameCMBC - ChainMaker - 默认值
	ContractStandardNameCMBC = "CMBC"
	// ContractStandardNameCMDFA ChainMaker - Contract Standard - Digital Fungible Assets
	ContractStandardNameCMDFA = "CMDFA"
	// ContractStandardNameCMNFA ChainMaker - Contract Standard - Digital Non-Fungible Assets
	ContractStandardNameCMNFA = "CMNFA"
	// ContractStandardNameCMID ChainMaker - Contract Standard - Identity
	ContractStandardNameCMID = "CMID"
	// ContractStandardNameCMEVI  ChainMaker - Contract Standard - Evidence
	ContractStandardNameCMEVI = "CMEVI"

	// ContractStandardNameEVMDFA  - Contract Standard - ERC-20 Fungible Assets
	ContractStandardNameEVMDFA = "ERC20"
	// ContractStandardNameEVMNFA  - Contract Standard - ERC-721 Non-Fungible Assets
	ContractStandardNameEVMNFA = "ERC721"
)
const (
	GoRoutinePoolErr = "new ants pool error: "
)

const (
	//PayloadContractNameBNS BNS合约名称
	PayloadContractNameBNS = "official_bns"
	//SubChainSpvPrefix 主子链-跨链高度合约前缀合约名称
	SubChainSpvPrefix = "official_spv"
	//SubChainInfoRWSetPrefix 子链信息读写集前缀
	SubChainInfoRWSetPrefix = "c"
	//SubChainGatewayRWSetPrefix
	SubChainGatewayRWSetPrefix = "g"
)

const (
	//PayloadMethodEvidence 存证合约方法
	PayloadMethodEvidence = "Evidence"
	//PayloadMethodEvidenceBatch 批量存证合约方法
	PayloadMethodEvidenceBatch = "EvidenceBatch"
)

const (
	//TxReadWriteKeyChainConfig 读写集key，链配置
	TxReadWriteKeyChainConfig = "CHAIN_CONFIG"
)

//const warnMsg = "上链内容违反相关法律规定，内容已屏蔽"
//
//var DecimalMsg = fmt.Errorf("无法将字符串转换为 Decimal 值")

const (
	//AddrTypeUser 节点地址
	AddrTypeUser = 0
	//AddrTypeContract 合约地址
	AddrTypeContract = 1
)

const (
	//RuntimeTypeDockerGo 合约类型
	RuntimeTypeDockerGo = "DOCKER_GO"
	//RuntimeTypeGo RuntimeTypeGo
	RuntimeTypeGo = "GO"
	//RuntimeTypeEVM 合约类型
	RuntimeTypeEVM = "EVM"
)

const (
	IDADataCategoryData = 1
	IDADataCategoryAPI  = 2
)

const (
	IDADataScaleTypeNum = 1
	IDADataScaleTypeM   = 2
	IDADataScaleTypeG   = 3
)

// 使用对象类别，1: 政府用户, 2: 企业用户, 3: 个人用户, 4: 无限制用户
const (
	IDADataUserGovernment   = 1
	IDADataUserEnterprise   = 2
	IDADataUserIndividual   = 3
	IDADataUserUnrestricted = 4
)

const (
	IDAUpdateCycleMinute = 1
	IDAUpdateCycleHour   = 2
	IDAUpdateCycleday    = 3
)

// UpdateCycleType constants
const (
	UpdateCycleTypeStatic   = 1 // 静态
	UpdateCycleTypeRealTime = 2 // 实时
	UpdateCycleTypePeriodic = 3 // 周期
	UpdateCycleTypeOther    = 4 // 其他
)

var UserCategoryDescriptions = map[int]string{
	IDADataUserGovernment:   "政府用户",
	IDADataUserEnterprise:   "企业用户",
	IDADataUserIndividual:   "个人用户",
	IDADataUserUnrestricted: "无限制用户",
}

// ERC20Functions 验证ERC20合约的方法列表
var ERC20Functions = map[string]bool{
	"transfer":     true,
	"transferFrom": true,
	"approve":      true,
	"balanceOf":    true,
	"allowance":    true,
	"totalSupply":  true,
}

const (
	//TopicTransferEvent  Transfer Event
	TopicTransferEvent = "transfer"
	//TopicApproveEvent Approve Event
	TopicApproveEvent = "approve"
	//TopicMintEvent Mint Event
	TopicMintEvent = "mint"
	//TopicBurnEvent Burn Event
	TopicBurnEvent = "burn"
)

// ERC721Functions 验证ERC721合约的方法列表
var ERC721Functions = map[string]bool{
	"balanceOf":         true,
	"ownerOf":           true,
	"safeTransferFrom":  true,
	"transferFrom":      true,
	"approve":           true,
	"setApprovalForAll": true,
	"getApproved":       true,
}

// EVMEventTopicTransfer EVM topic transfer
const EVMEventTopicTransfer = "ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"

// TopicEventDataKey event topic
var TopicEventDataKey = map[string]string{
	"mint":         "CMDFA",
	"transfer":     "CMDFA",
	"burn":         "CMDFA",
	"Mint":         "CMNFA",
	"TransferFrom": "CMNFA",
	"Burn":         "CMNFA",
	"ddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef": "EVM",
}

const (
	//BNSBindEvent 绑定BNS
	BNSBindEvent = "Bind"
	//BNSUnBindEvent解绑BNS
	BNSUnBindEvent = "UnBind"

	//DIDSetDidDocument 设置DID
	DIDSetDidDocument = "SetDidDocument"

	//TopicTxAddBlack 加入交易黑名单
	TopicTxAddBlack = "100"
	//TopicTxDeleteBlack 删除交易黑名单
	TopicTxDeleteBlack = "101"
)

const (
	KeyUserInfo           = "UserInfo"
	KeyIDABasic           = "IDABasic"
	KeyIDAOwnership       = "IDAOwnership"
	KeyIDASource          = "IDASource"
	KeyIDAScenarios       = "IDAScenarios"
	KeyIDASupply          = "IDASupply"
	KeyIDADetails         = "IDADetails"
	KeyIDAPrivacy         = "IDAPrivacy"
	KeyIDAStatus          = "IDAStatus"
	KeyIDAColumns         = "IDAColumns"
	KeyIDAAPI             = "IDAApi"
	KeyIDACertifications  = "IDACertifications"
	KeyIDADataSet         = "IDADataSet"
	KeyIDAContractVersion = "IDAContractVersion"
	KeyIDAEnName          = "IDAEnName"
	KeyRegisterCount      = "RegisterCount"
	KeyPlatformInfo       = "PlatformInfo"
	KeyPlatformCount      = "PlatformCount"
)

var (
	GlobalAbiERC20  *abi.ABI
	GlobalAbiERC721 *abi.ABI
)

// RemoveAddrPrefix 删除地址的0x开头
func RemoveAddrPrefix(address string) string {
	pattern := `^0x`
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return address
	}
	return regex.ReplaceAllString(address, "")
}

// isZeroAddress 是否是空地址
func IsZeroAddress(address string) bool {
	zeroAddr := "0000000000000000000000000000000000000000"
	return address == zeroAddr
}

// StringAmountDecimal  string转decimal
func StringAmountDecimal(amount string, decimals int) decimal.Decimal {
	// 将字符串转换为 decimal.Decimal 值
	amountDecimal, _ := decimal.NewFromString(amount)
	// 创建一个新的 decimal.Decimal 值，表示 10^decimals
	divisor := decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(decimals)))
	// 使用 Div 方法将 amountDecimal 除以 divisor
	//nolint:gosec
	resultDecimal := amountDecimal.DivRound(divisor, int32(decimals))
	return resultDecimal
}

// copyMap 复制变量
func CopyMap(src map[string]bool) map[string]bool {
	dst := make(map[string]bool, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}

// getMemberInfoKey 根据用户信息构造key
func GetMemberInfoKey(chainId, hashType string, memberType int32, memberBytes []byte) (string, error) {
	mHash := md5.New()
	prefix := config.GlobalConfig.RedisDB.Prefix
	redisKey := fmt.Sprintf(cache.RedisUserMemberInfoKey, prefix, chainId, hashType, memberType)
	_, err := mHash.Write([]byte(redisKey))
	if err != nil {
		return "", err
	}
	_, err = mHash.Write(memberBytes)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", mHash.Sum(nil)), nil
}

// MD5 md5
func MD5(v string) string {
	d := []byte(v)
	m := md5.New()
	m.Write(d)
	return hex.EncodeToString(m.Sum(nil))
}

// IsContractTx 是否合约类交易
func IsContractTx(txInfo *common.Transaction) bool {
	if txInfo == nil {
		return false
	}
	payload := txInfo.Payload
	if txInfo.Result.ContractResult.Code != 0 {
		return false
	}

	if payload.ContractName == syscontract.SystemContract_CONTRACT_MANAGE.String() {
		return true
	}

	if payload.ContractName == syscontract.SystemContract_MULTI_SIGN.String() {
		return true
	}

	return false
}

// IsContractTxByName 根据合约名和方法名判断是否是合约类交易, 目前只判断合约管理类和多重签名类
// @param contractName 合约名
// @param contractMethod 合约方法名
// @return 是否是合约类交易
func IsContractTxByName(contractName, contractMethod string) bool {
	//判断是否是合约管理类交易,目前只判断合约管理类和多重签名类
	if contractName == syscontract.SystemContract_CONTRACT_MANAGE.String() &&
		(contractMethod == syscontract.ContractManageFunction_INIT_CONTRACT.String() ||
			contractMethod == syscontract.ContractManageFunction_UPGRADE_CONTRACT.String()) {
		return true
	}

	//判断是否是多重签名类交易
	if contractName == syscontract.SystemContract_MULTI_SIGN.String() &&
		contractMethod == syscontract.MultiSignFunction_REQ.String() {
		return true
	}

	return false
}

// IsContractManageTx 是否是合约管理类交易
// @param contractName 合约名
// @param contractMethod 合约方法名
// @return bool 是否是合约管理类交易
func IsContractManageTx(contractName, contractMethod string) bool {
	// 判断是否是合约管理类交易
	if contractName != syscontract.SystemContract_CONTRACT_MANAGE.String() {
		return false
	}
	if contractMethod == syscontract.ContractManageFunction_INIT_CONTRACT.String() ||
		contractMethod == syscontract.ContractManageFunction_UPGRADE_CONTRACT.String() {
		return true
	}
	return false
}

// IsMultiSignTx 是否是多重签名类交易
// @param contractName 合约名
// @param contractMethod 合约方法名
// @return bool 是否是多重签名类交易
func IsMultiSignTx(contractName, contractMethod string) bool {
	// 判断是否是多重签名类交易
	return contractName == syscontract.SystemContract_MULTI_SIGN.String()
}

// IsConfigTx 是否修改配置类交易
func IsConfigTx(txInfo *common.Transaction) bool {
	if txInfo == nil || txInfo.Payload == nil {
		return false
	}

	return txInfo.Payload.ContractName == syscontract.SystemContract_CHAIN_CONFIG.String()
}

// GetContractSendTxId 获取合约交易发送方交易ID
// @param txInfo 交易信息
// @return 交易ID
func GetContractSendTxId(txInfo *common.Transaction) string {
	txId := txInfo.Payload.TxId
	isMultiSignTx := IsMultiSignTx(txInfo.Payload.ContractName, txInfo.Payload.Method)
	if isMultiSignTx {
		for _, parameter := range txInfo.Payload.Parameters {
			switch parameter.Key {
			case syscontract.MultiVote_TX_ID.String():
				txId = string(parameter.Value)
			}
		}
	}
	return txId
}

// IsRelayCrossChainTx 是否是跨链类交易
func IsRelayCrossChainTx(txInfo *common.Transaction) bool {
	if txInfo == nil || txInfo.Payload == nil {
		return false
	}

	return txInfo.Payload.ContractName == syscontract.SystemContract_RELAY_CROSS.String()
}

// 判断是否为跨链信息
func IsRelayCrossChainInfo(writeKey, contractName string) bool {
	// 如果writeKey或contractName为空，则返回false
	if writeKey == "" || contractName == "" {
		return false
	}
	// 如果contractName为RELAY_CROSS，且writeKey以SubChainInfoRWSetPrefix开头，则返回true
	if contractName == syscontract.SystemContract_RELAY_CROSS.String() &&
		strings.HasPrefix(writeKey, SubChainInfoRWSetPrefix) {
		return true
	}

	// 否则返回false
	return false
}

func IsRelayCrossGatewayId(writeKey, contractName string) bool {
	// 如果writeKey或contractName为空，则返回false
	if writeKey == "" || contractName == "" {
		return false
	}
	// 如果contractName为RELAY_CROSS，且writeKey以SubChainInfoRWSetPrefix开头，则返回true
	if contractName == syscontract.SystemContract_RELAY_CROSS.String() &&
		strings.HasPrefix(writeKey, SubChainGatewayRWSetPrefix) {
		return true
	}

	// 否则返回false
	return false
}

// IsSubChainSpvContractTx 是否是跨链子链同步区块高度类交易
func IsSubChainSpvContractTx(txInfo *common.Transaction) (bool, string) {
	if txInfo == nil || txInfo.Payload == nil {
		return false, ""
	}

	contractName := txInfo.Payload.ContractName
	if strings.HasPrefix(contractName, SubChainSpvPrefix) {
		return true, contractName
	}
	return false, ""
}

func IsInBlockHeight(height int64, heightList []int64) bool {
	for _, value := range heightList {
		if height == value {
			return true
		}
	}

	return false
}

func GetMaxBlockHeight(heightList []int64) int64 {
	var maxHeight int64
	if len(heightList) == 0 {
		return maxHeight
	}

	for _, height := range heightList {
		if height > maxHeight {
			maxHeight = height
		}
	}

	return maxHeight
}

func GetMinBlockHeight(heightList []int64) int64 {
	if len(heightList) == 0 {
		return 0
	}

	minHeight := heightList[0]
	for _, height := range heightList {
		if height < minHeight {
			minHeight = height
		}
	}

	return minHeight
}

// IsMainChainGateway 是否是主链
func IsMainChainGateway(gatewayID string) bool {
	return strings.HasPrefix(gatewayID, tcipCommon.MainGateway_MAIN_GATEWAY_ID.String())
}

// SaveJsonFile 保存测试数据:		SaveJsonFile("cross_1_txInfoJson", txInfo)
// func SaveJsonFile(fix string, valueJson interface{}) {
// 	nowTime := time.Now().Unix()
// 	fileName := fmt.Sprintf("%s_%d.json", fix, nowTime)

// 	log.Errorf("=11111111=======fileName=====:%v", fileName)
// 	// 创建一个文件
// 	file, _ := os.Create(fileName)

// 	// 创建一个 json.Encoder 并设置缩进
// 	encoder := json.NewEncoder(file)
// 	encoder.SetIndent("", "  ")

// 	// 将 txInfo 写入到文件中 todo _1111111111
// 	//_ = encoder.Encode(valueJson)
// }

// ParallelParseBatchWhere 将交易分割为大小为batchSize的批次
func ParallelParseBatchWhere(wheres []string, batchSize int) [][]string {
	batches := make([][]string, 0)
	batch := make([]string, 0)

	for _, where := range wheres {
		batch = append(batch, where)
		if len(batch) == batchSize {
			batches = append(batches, batch)
			batch = make([]string, 0)
		}
	}

	if len(batch) > 0 {
		batches = append(batches, batch)
	}
	return batches
}

// GetEvmAbi 获取 EVM ABI，如果全局变量中不存在则从文件读取
func GetEvmAbi(evmType string) *abi.ABI {
	var ercAbi *abi.ABI
	// 根据类型选择 ABI
	if evmType == ContractStandardNameEVMDFA {
		if GlobalAbiERC20 == nil {
			GlobalAbiERC20 = getERC20AbiJson()
		}
		ercAbi = GlobalAbiERC20
	} else if evmType == ContractStandardNameEVMNFA {
		if GlobalAbiERC721 == nil {
			GlobalAbiERC721 = getERC721AbiJson()
		}
		ercAbi = GlobalAbiERC721
	}

	return ercAbi
}

// getERC20AbiJson 解析 EVM ERC20 ABI
func getERC20AbiJson() *abi.ABI {
	ercAbiJson, err := os.Open("./erc20_abi.json")
	if err != nil {
		dir, _ := os.Getwd()
		log.Errorf("Failed to open erc20_abi.json, dir:%v", dir)
		return nil
	}
	defer ercAbiJson.Close() // 确保文件在使用后关闭

	ercAbi, err := abi.JSON(ercAbiJson)
	if err != nil {
		log.Errorf("Failed to read erc20_abi: %v", err)
		return nil
	}
	return &ercAbi
}

// getERC721AbiJson 解析 EVM ERC721 ABI
func getERC721AbiJson() *abi.ABI {
	ercAbiJson, err := os.Open("./erc721_abi.json")
	if err != nil {
		dir, _ := os.Getwd()
		log.Errorf("Failed to open erc721_abi.json, dir:%v", dir)
		return nil
	}
	defer ercAbiJson.Close() // 确保文件在使用后关闭

	ercAbi, err := abi.JSON(ercAbiJson)
	if err != nil {
		log.Errorf("Failed to read erc721_abi: %v", err)
		return nil
	}
	return &ercAbi
}

// IsCrossEnd 主子链跨链交易是否结束
func IsCrossEnd(status int32) bool {
	if status == int32(tcipCommon.CrossChainStateValue_CONFIRM_END) ||
		status == int32(tcipCommon.CrossChainStateValue_CANCEL_END) {
		return true
	}
	return false
}

// ExtractTxIdsAndContractNames
//
//	@Description:  获取交易列表涉及到的所有的合约名称
//	@param txInfoList 交易列表
//	@return []string 交易id列表
//	@return map[string]string 合约列表
//	@return map[string]*db.Transaction  交易Map
func ExtractTxIdsAndContractNames(txInfoList []*db.Transaction) ([]string, map[string]string,
	map[string]*db.Transaction) {
	txInfoMap := make(map[string]*db.Transaction)
	txIds := make([]string, 0)
	contractAddrMap := make(map[string]string)
	for _, txInfo := range txInfoList {
		txInfoMap[txInfo.TxId] = txInfo
		txIds = append(txIds, txInfo.TxId)
		if txInfo.ContractAddr == "" {
			continue
		}

		contractAddrMap[txInfo.ContractAddr] = txInfo.ContractAddr
	}

	return txIds, contractAddrMap, txInfoMap
}

// checkGasEnabled
//
//	@Description: 检查是否开启Gas，未开启的话设置gas消耗为空
//	@param chainId
//	@param dealResult
//	@return error
func CheckAndDisableGasIfNotEnabled(chainId string, dealResult *model.ProcessedBlockData) error {
	chainInfo, err := dbhandle.GetChainInfoById(chainId)
	if err != nil {
		return err
	}
	if chainInfo != nil && !chainInfo.EnableGas {
		//未启用gas
		dealResult.GasRecordList = make([]*db.GasRecord, 0)
	}
	return nil
}
