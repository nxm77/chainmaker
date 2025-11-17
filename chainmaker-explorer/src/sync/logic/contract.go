/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package logic

import (
	"chainmaker_web/src/chain"
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/sync/common"
	"chainmaker_web/src/sync/model"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/shopspring/decimal"

	"chainmaker.org/chainmaker/contract-utils/standard"
	pbCommon "chainmaker.org/chainmaker/pb-go/v2/common"
	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	"github.com/gogo/protobuf/proto"
	"github.com/panjf2000/ants/v2"
)

// ParallelParseContract
//
//	@Description: 并发处理合约数据
//	@param blockInfo 区块数据
//	@param hashType
//	@param dealResult 结果集
//	@return error
func ParallelParseContract(blockInfo *pbCommon.BlockInfo, hashType string, dealResult *model.ProcessedBlockData) error {
	// 使用同步互斥锁保护共享资源
	var mutx sync.Mutex
	var wg sync.WaitGroup

	// 创建一个固定大小的 goroutine 池
	goRoutinePool, err := ants.NewPool(10, ants.WithPreAlloc(false))
	if err != nil {
		log.Errorf("Failed to create goroutine pool: %v", err)
		return err
	}
	defer goRoutinePool.Release()

	chainId := blockInfo.Block.Header.ChainId
	//nolint:gosec
	blockHeight := int64(blockInfo.Block.Header.BlockHeight)
	// 用来接收并发任务的错误，减少通道的容量，避免阻塞
	errChan := make(chan error, 10)

	for i, tx := range blockInfo.Block.Txs {
		wg.Add(1)
		// 使用匿名函数传递参数，避免使用外部变量
		errSub := goRoutinePool.Submit(func(i int, blockInfo *pbCommon.BlockInfo, txInfo *pbCommon.Transaction) func() {
			return func() {
				defer wg.Done()
				var err error
				payload := txInfo.Payload
				//所有上链的数据都是invoke数据
				if payload.TxType != pbCommon.TxType_QUERY_CONTRACT &&
					payload.TxType != pbCommon.TxType_ARCHIVE &&
					payload.TxType != pbCommon.TxType_SUBSCRIBE &&
					payload.TxType != pbCommon.TxType_INVOKE_CONTRACT {
					return
				}

				// 计算账户信息
				userResult, err := GetSenderAndPayerUser(chainId, hashType, txInfo)
				if err != nil {
					log.Errorf("ParallelParseTransactions get User err:%v", err)
				}

				//存证合约
				evidenceList, err := DealEvidence(blockHeight, txInfo, userResult)
				if err != nil {
					log.Errorf("ParallelParseTransactions DealEvidence err:%v", err)
					errChan <- err
					return
				}

				//处理合约源码
				contractByteCodeInfo := BuildContractByteCode(txInfo)

				//处理通用合约数据
				contractWriteSetData, err := BuildContractInfo(i, blockInfo, txInfo, userResult)
				if err != nil {
					log.Errorf("ParallelParseTransactions BuildContractInfo err:%v")
					errChan <- err
					return
				}

				// 构造合约升级交易数据
				upgradeTx := buildUpgradeContractTransaction(contractWriteSetData)

				mutx.Lock()         // 锁定互斥锁
				defer mutx.Unlock() // 使用 defer 确保互斥锁被解锁
				if contractWriteSetData != nil {
					//合约修改事件】
					if dealResult.ContractWriteSetData == nil {
						dealResult.ContractWriteSetData = make(map[string]*model.ContractWriteSetData, 0)
					}
					dealResult.ContractWriteSetData[contractWriteSetData.SenderTxId] = contractWriteSetData
				}

				//处理部署，升级合约
				if upgradeTx != nil {
					dealResult.UpgradeContractTx = append(dealResult.UpgradeContractTx, upgradeTx)
				}

				//处理合约源码
				if contractByteCodeInfo != nil {
					//合约源码事件
					dealResult.ContractByteCodeList = append(dealResult.ContractByteCodeList, contractByteCodeInfo)
				}

				//存证合约
				if len(evidenceList) > 0 {
					dealResult.EvidenceList = append(dealResult.EvidenceList, evidenceList...)
				}
			}
		}(i, blockInfo, tx))
		if errSub != nil {
			log.Error("ParallelParseContract submit Failed : " + errSub.Error())
			wg.Done() // 减少 WaitGroup 计数
		}
	}

	wg.Wait()
	if len(errChan) > 0 {
		errLog := <-errChan
		return errLog
	}
	return nil
}

// GenesisBlockSystemContract
//
//	@Description: 创世区块解析系统合约
//	@param blockInfo
//	@param dealResult
//	@return error
func GenesisBlockSystemContract(blockInfo *pbCommon.BlockInfo, dealResult *model.ProcessedBlockData) error {
	blockHeight := blockInfo.Block.Header.BlockHeight
	timestamp := blockInfo.Block.Header.BlockTimestamp
	//创世区块
	if blockHeight != 0 {
		return nil
	}

	for i, txInfo := range blockInfo.Block.Txs {
		rwSetList := blockInfo.RwsetList[i]
		if rwSetList == nil {
			continue
		}

		//解析读写集
		systemContractMap, err := GenesisBlockGetContractByWriteSet(rwSetList.TxWrites)
		if err != nil {
			return err
		}

		if len(systemContractMap) == 0 {
			continue
		}

		for _, contract := range systemContractMap {
			runtimeType := pbCommon.RuntimeType_name[int32(contract.RuntimeType)]
			contractInfo := &db.Contract{
				Name:           contract.Name,
				NameBak:        contract.Name,
				Version:        contract.Version,
				RuntimeType:    runtimeType,
				CreateTxId:     txInfo.Payload.TxId,
				BlockHeight:    int64(blockHeight),
				Addr:           contract.Address,
				ContractStatus: dbhandle.SystemContractStatus,
				ContractType:   common.ContractStandardNameOTHER,
				Timestamp:      timestamp,
			}
			dealResult.InsertContracts = append(dealResult.InsertContracts, contractInfo)
		}
	}

	return nil
}

// BuildContractInfo
//
//	@Description: 构造合约数据
//	@param i
//	@param blockInfo
//	@param txInfo
//	@param userInfo
//	@return *db.Contract  合约数据
//	@return string  新增或者更新合约
//	@return error
func BuildContractInfo(i int, blockInfo *pbCommon.BlockInfo, txInfo *pbCommon.Transaction,
	userInfo *db.SenderPayerUser) (*model.ContractWriteSetData, error) {
	if blockInfo == nil || txInfo == nil || userInfo == nil {
		return nil, nil
	}

	//非合约类交易不用处理合约数据，直接返回
	isContractTx := common.IsContractTx(txInfo)
	if !isContractTx {
		return nil, nil
	}

	blockHeight := blockInfo.Block.Header.BlockHeight
	blockHash := hex.EncodeToString(blockInfo.Block.Header.BlockHash)
	rwSetList := blockInfo.RwsetList[i]
	payload := txInfo.Payload
	if rwSetList == nil {
		return nil, nil
	}

	var contractNameAddr string
	//解析合约参数，获取合约名称和合约字节码
	for _, parameter := range payload.Parameters {
		switch parameter.Key {
		case syscontract.InitContract_CONTRACT_NAME.String():
			contractNameAddr = string(parameter.Value)
		}
	}

	//解析读写集，获取合约信息
	contractWriteSet, err := GetContractByWriteSet(rwSetList.TxWrites, contractNameAddr)
	if err != nil || contractWriteSet.ContractResult == nil {
		return nil, nil
	}

	contractResult := contractWriteSet.ContractResult
	if contractResult.Address == "" {
		return nil, nil
	}

	var decimal int
	if contractWriteSet.Decimal != "" {
		decimal, _ = strconv.Atoi(contractWriteSet.Decimal)
	}

	runtimeType := contractResult.RuntimeType.String()
	if runtimeType == common.RuntimeTypeGo {
		runtimeType = common.RuntimeTypeDockerGo
	}
	var orgId string
	if contractResult.Creator != nil {
		orgId = contractResult.Creator.OrgId
	}

	//构造合约数据
	txId := common.GetContractSendTxId(txInfo)
	contractInfo := &model.ContractWriteSetData{
		ContractName:    contractResult.Name,
		ContractNameBak: contractResult.Name,
		ContractAddr:    contractResult.Address,
		ContractSymbol:  contractWriteSet.Symbol,
		Version:         contractResult.Version,
		RuntimeType:     runtimeType,
		//nolint:gosec
		ContractStatus: int32(contractResult.Status),
		//nolint:gosec
		BlockHeight: int64(blockHeight),
		BlockHash:   blockHash,
		OrgId:       orgId,
		SenderTxId:  txId,
		Sender:      userInfo.SenderUserId,
		SenderAddr:  userInfo.SenderUserAddr,
		SenderOrgId: userInfo.SenderOrgId,
		Timestamp:   payload.Timestamp,
		Decimals:    decimal,
	}

	//敏感词过滤
	_, flag := common.FilteringSensitive(contractInfo.ContractName)
	if flag {
		contractInfo.ContractName = config.ContractWarnMsg
	}
	return contractInfo, nil
}

// BuildContractByteCode
//
//	@Description: 构造合约数据
//	@param i
//	@param blockInfo
//	@param txInfo
//	@param userInfo
//	@return *db.Contract  合约数据
//	@return string  新增或者更新合约
//	@return error
func BuildContractByteCode(txInfo *pbCommon.Transaction) *db.ContractByteCode {
	if txInfo == nil {
		return nil
	}

	//非合约类交易不用处理合约数据，直接返回
	isContractTx := common.IsContractTx(txInfo)
	if !isContractTx {
		return nil
	}

	var byteCode []byte
	payload := txInfo.Payload
	//解析合约参数，获取合约名称和合约字节码
	for _, parameter := range payload.Parameters {
		switch parameter.Key {
		case syscontract.InitContract_CONTRACT_BYTECODE.String():
			byteCode = parameter.Value
		case syscontract.UpgradeContract_CONTRACT_BYTECODE.String():
			byteCode = parameter.Value
		}
	}

	txId := common.GetContractSendTxId(txInfo)
	if byteCode == nil {
		log.Warnf("BuildContractByteCode byteCode is nil, txId: %s", txId)
		return nil
	}

	//构造合约数据
	contractByteCodeInfo := &db.ContractByteCode{
		TxId:     txId,
		ByteCode: byteCode,
	}

	return contractByteCodeInfo
}

// GetContractSDKData
//
//	@Description: 从SDK获取合约类型，简称，小数
//	@param chainId
//	@param contractInfo 合约数据
//	@return string 合约类型
//	@return string 合约检查
//	@return int 合约小数
//	@return error
func GetContractSDKData(chainId string, contractInfo *db.Contract, byteCode []byte) error {
	//计算合约类型
	contractType, err := GetContractType(chainId, contractInfo.Name, contractInfo.RuntimeType, byteCode)
	if err != nil {
		log.Error("BuildContractInfo get contractType err: " + err.Error())
		return err
	}
	contractInfo.ContractType = contractType

	//只有同质化合约才有合约简称和小数
	if contractType == common.ContractStandardNameCMDFA ||
		contractType == common.ContractStandardNameEVMDFA {
		//获取合约简称
		if contractInfo.ContractSymbol == "" {
			symbol, _ := GetContractSymbol(contractType, chainId, contractInfo.Addr)

			contractInfo.ContractSymbol = symbol
			if err != nil {
				log.Warnf("GetContractSDKData GetContractSymbol err: %s", err.Error())
			}
		}

		//获取合约合约小数
		if contractInfo.Decimals == 0 {
			decimals, err := GetContractDecimals(chainId, contractType, contractInfo.Addr)
			contractInfo.Decimals = decimals
			if err != nil {
				log.Warnf("GetContractSDKData GetContractDecimals err: %s", err.Error())
			}
		}
	}
	return nil
}

// dealStandardContract
//
//	@Description: 将合约数据处理成同质化，非同质化合约数据
//	@param contract
//	@return *db.FungibleContract 同质化合约
//	@return *db.NonFungibleContract 非同质化合约
func dealStandardContract(contract *db.Contract) (*db.FungibleContract, *db.NonFungibleContract) {
	switch contract.ContractType {
	case common.ContractStandardNameCMDFA:
		fallthrough
	case common.ContractStandardNameEVMDFA:
		//同质化合约
		fungibleContract := &db.FungibleContract{
			ContractName:    contract.Name,
			ContractNameBak: contract.NameBak,
			Symbol:          contract.ContractSymbol,
			ContractAddr:    contract.Addr,
			ContractType:    contract.ContractType,
			TotalSupply:     decimal.Zero,
			Timestamp:       contract.Timestamp,
		}
		return fungibleContract, nil
	case common.ContractStandardNameCMNFA:
		fallthrough
	case common.ContractStandardNameEVMNFA:
		//非同质化合约
		nonFungibleContract := &db.NonFungibleContract{
			ContractName:    contract.Name,
			ContractNameBak: contract.NameBak,
			ContractAddr:    contract.Addr,
			ContractType:    contract.ContractType,
			TotalSupply:     decimal.Zero,
			Timestamp:       contract.Timestamp,
		}
		return nil, nonFungibleContract
	}

	return nil, nil
}

// dealIDAContractData
//
//	@Description: 将合约数据处理成同质化，非同质化合约数据
//	@param contract
//	@return *db.FungibleContract 同质化合约
//	@return *db.NonFungibleContract 非同质化合约
func dealIDAContractData(contract *db.Contract) *db.IDAContract {
	if contract.ContractType != standard.ContractStandardNameCMIDA {
		return nil
	}

	idaContract := &db.IDAContract{
		ContractName:    contract.Name,
		ContractNameBak: contract.NameBak,
		ContractAddr:    contract.Addr,
		ContractType:    contract.ContractType,
		Timestamp:       contract.Timestamp,
	}
	return idaContract
}

// GetContractType
//
//	@Description: 获取合约类型
//	@param chainId
//	@param contractName 合约名称
//	@param runtimeType
//	@param bytecode 创建合约bytecode
//	@return string
//	@return error
func GetContractType(chainId, contractName, runtimeType string, bytecode []byte) (string, error) {
	var err error
	contractType := common.ContractStandardNameOTHER
	if runtimeType == common.RuntimeTypeDockerGo {
		contractType, err = common.DockerGetContractType(chainId, contractName)
		if err != nil {
			log.Warnf("【sdk】GetContractType docker go err :%v", err)
			//失败重试一次
			contractType, err = common.DockerGetContractType(chainId, contractName)
			if err != nil {
				log.Warnf("【sdk】GetContractType docker go err :%v", err)
			}
		}
		return contractType, nil
	} else if runtimeType == common.RuntimeTypeEVM {
		if len(bytecode) == 0 {
			log.Warnf("【sdk】GetContractType EVM err, bytecode is nil")
			return common.ContractStandardNameOTHER, fmt.Errorf("GetContractType bytecode is nil ")
		}

		//获取4字节列表
		signatures := common.ExtractFunctionSignatures(bytecode)
		// 检查是否包含所有ERC20函数
		if containsAllFunctions(common.ContractStandardNameEVMDFA, signatures, common.CopyMap(common.ERC20Functions)) {
			return common.ContractStandardNameEVMDFA, nil
		}

		// 检查是否包含所有ERC721函数
		if containsAllFunctions(common.ContractStandardNameEVMNFA, signatures, common.CopyMap(common.ERC721Functions)) {
			return common.ContractStandardNameEVMNFA, nil
		}
	}

	log.Infof("GetContractType Not a standard contract，contractName：%v, runtimeType:%v",
		contractName, runtimeType)
	return contractType, nil
}

// containsAllFunctions
//
//	@Description: 判断EVM合约方法是否在标准合约方法中
//	@param evmType evm合约类型
//	@param signatures 合约方法的byte值
//	@param functionNames 标准合约方法名称
//	@return bool
func containsAllFunctions(evmType string, signatures [][]byte, functionNames map[string]bool) bool {
	// 创建一个通道用于接收找到的函数名
	foundChan := make(chan string, len(signatures))
	ercAbi := common.GetEvmAbi(evmType)
	if ercAbi == nil {
		log.Errorf("containsAllFunctions unmarshal ercAbi failed, ercAbi is null")
		return false
	}

	// 使用一个 WaitGroup 来等待所有的 goroutine 完成
	var wg sync.WaitGroup

	// 遍历签名并调用 EVMGetMethodName
	var allNameList []string
	for _, sig := range signatures {
		wg.Add(1)
		go func(sig []byte) {
			defer wg.Done()
			name, _ := common.EVMGetMethodName(ercAbi, sig)
			foundChan <- name
			allNameList = append(allNameList, name)
		}(sig)
	}

	// 等待所有的 goroutine 完成
	go func() {
		wg.Wait()
		close(foundChan)
	}()

	// 从通道中读取找到的函数名
	for range signatures {
		name := <-foundChan
		if _, found := functionNames[name]; found {
			// 如果找到匹配的函数名，则从映射中删除
			delete(functionNames, name)

			// 如果映射为空，则已找到所有函数名
			if len(functionNames) == 0 {
				return true
			}
		}
	}

	if len(functionNames) > 0 {
		allNameListJson, _ := json.Marshal(allNameList)
		functionNamesJson, _ := json.Marshal(functionNames)
		log.Infof("【sdk】EVM ContractType containsAllFunctions allNameList:%v, not have name :%v",
			string(allNameListJson), string(functionNamesJson))
	}
	// 如果映射不为空，则没有找到所有函数名
	return false
}

// GetContractSymbol
//
//	@Description: 获取合约简称
//	@param chainId
//	@param contractType 合约类型
//	@param contractAddr 合约地址
//	@return string 简称
//	@return error
func GetContractSymbol(chainId, contractType, contractAddr string) (string, error) {
	var symbolName string
	var err error
	if contractType == common.ContractStandardNameCMDFA ||
		contractType == common.ContractStandardNameCMNFA {
		symbolName, err = common.DockerGetContractSymbol(chainId, contractAddr)
	}

	if contractType == common.ContractStandardNameEVMDFA ||
		contractType == common.ContractStandardNameEVMNFA {
		symbolName, err = common.EVMGetContractSymbol(chainId, contractAddr, contractType)
	}
	return symbolName, err
}

// GetTotalSupply
//
//	@Description: 获取合约总发行量
//	@param contractType
//	@param chainId
//	@param contractName
//	@return string
//	@return error
func GetTotalSupply(contractType, chainId, contractName string) (string, error) {
	totalSupply := "0"
	var err error
	if contractType == common.ContractStandardNameCMDFA {
		totalSupply, err = common.DockerGetTotalSupply(chainId, contractName)
	}

	if contractType == common.ContractStandardNameEVMDFA {
		totalSupply, err = common.EvmGetTotalSupply(contractType, chainId, contractName)
	}

	return totalSupply, err
}

// GetContractDecimals
//
//	@Description: 获取合约小数位数
//	@param chainId
//	@param contractType 合约类型
//	@param contractName 合约名称，地址
//	@return int
//	@return error
func GetContractDecimals(chainId, contractType, contractName string) (int, error) {
	var decimals int
	var err error
	if contractType == common.ContractStandardNameCMDFA {
		decimals, err = common.DockerGetDecimals(chainId, contractName)
	}

	if contractType == common.ContractStandardNameEVMDFA {
		decimals, err = common.EvmGetDecimals(contractType, chainId, contractName)
	}

	return decimals, err
}

// GetContractByWriteSet
//
//	@Description: 根据读写接，解析合约数据
//	@param txWriteList
//	@return *db.GetContractWriteSet
//	@return error
func GetContractByWriteSet(txWriteList []*pbCommon.TxWrite, contractNameAddr string) (*db.GetContractWriteSet, error) {
	contractWriteSet := &db.GetContractWriteSet{}
	var contractResult pbCommon.Contract
	var contractResultBak pbCommon.Contract
	writeKeyContract := "Contract:" + contractNameAddr
	for _, write := range txWriteList {
		if string(write.Key) == writeKeyContract {
			err := proto.Unmarshal(write.Value, &contractResult)
			if err != nil {
				return contractWriteSet, err
			}
			contractWriteSet.ContractResult = &contractResult
		}
		if strings.HasPrefix(string(write.Key), "Contract:") {
			err := proto.Unmarshal(write.Value, &contractResultBak)
			if err != nil {
				log.Errorf("GetContractByWriteSet Unmarshal err:%v", err)
			}
		} else if strings.HasPrefix(string(write.Key), "ContractByteCode:") {
			contractWriteSet.ByteCode = write.Value
		} else if string(write.Key) == "decimal" {
			contractWriteSet.Decimal = string(write.Value)
		} else if string(write.Key) == "symbol" {
			contractWriteSet.Symbol = string(write.Value)
		}
	}

	if contractWriteSet.ContractResult == nil {
		contractWriteSet.ContractResult = &contractResultBak
	}

	return contractWriteSet, nil
}

// GenesisBlockGetContractByWriteSet
//
//	@Description: 根据读写接，解析合约数据
//	@param txWriteList
//	@return *db.GetContractWriteSet
//	@return error
func GenesisBlockGetContractByWriteSet(txWriteList []*pbCommon.TxWrite) (map[string]pbCommon.Contract, error) {
	systemContractList := make(map[string]pbCommon.Contract, 0)
	for _, write := range txWriteList {
		var contractResult pbCommon.Contract
		if strings.HasPrefix(string(write.Key), "Contract:") {
			err := proto.Unmarshal(write.Value, &contractResult)
			if err != nil {
				return systemContractList, err
			}

			if contractResult.Address != "" {
				systemContractList[contractResult.Address] = contractResult
			}
		}
	}
	return systemContractList, nil
}

// GetContractMapByAddrs 根据合约地址获取合约数据
// @param chainId
// @param contractAddrMap
// @return map[string]*db.Contract
// @return error
func GetContractMapByAddrs(chainId string, contractAddrMap map[string]string) (map[string]*db.Contract, error) {
	//获取本次涉及到的合约信息
	contractAddrs := make([]string, 0)
	for _, addr := range contractAddrMap {
		contractAddrs = append(contractAddrs, addr)
	}

	// 批量查询合约信息
	contractMap, err := dbhandle.GetContractByAddrs(chainId, contractAddrs)
	return contractMap, err
}

// HandleContractInsertOrUpdate
//
//	@Description: 顺序处理合约数据
//	@param chainId
//	@param dealResult
//	@return error
func HandleContractInsertOrUpdate(chainId string, dealResult *model.ProcessedBlockData) error {
	if len(dealResult.ContractWriteSetData) == 0 {
		return nil
	}

	// 顺序处理每个合约
	for _, contractData := range dealResult.ContractWriteSetData {
		var contractInfo *db.Contract
		// 缓存判断合约是否存在
		contractDB, err := dbhandle.GetContractByAddr(chainId, contractData.ContractAddr)
		if err != nil {
			return err
		}

		if contractDB == nil {
			txNum, _ := dbhandle.GetTxNumByContractName(chainId, contractData.ContractNameBak)
			eventNum, _ := dbhandle.GetEventNumByContractName(chainId, contractData.ContractNameBak)
			contractInfo = &db.Contract{
				Name:             contractData.ContractName,
				NameBak:          contractData.ContractNameBak,
				Addr:             contractData.ContractAddr,
				Version:          contractData.Version,
				RuntimeType:      contractData.RuntimeType,
				ContractStatus:   contractData.ContractStatus,
				ContractType:     contractData.ContractType,
				ContractSymbol:   contractData.ContractSymbol,
				Decimals:         contractData.Decimals,
				TxNum:            txNum,
				EventNum:         eventNum,
				OrgId:            contractData.OrgId,
				CreateTxId:       contractData.SenderTxId,
				CreateSender:     contractData.Sender,
				CreatorAddr:      contractData.SenderAddr,
				Timestamp:        contractData.Timestamp,
				Upgrader:         contractData.Sender,
				UpgradeAddr:      contractData.SenderAddr,
				UpgradeOrgId:     contractData.OrgId,
				UpgradeTimestamp: contractData.Timestamp,
			}

			//处理其他合约数据
			BuildContractOtherData(chainId, contractData.SenderTxId, contractInfo, dealResult)

			// 获取标准化合约数据
			dbFungibleContract, dbNonFungibleContract := dealStandardContract(contractInfo)

			// 构造IDA合约数据
			dbIDAContract := dealIDAContractData(contractInfo)

			// 更新插入数据
			dealResult.InsertContracts = append(dealResult.InsertContracts, contractInfo)
			if dbFungibleContract != nil {
				dealResult.FungibleContract = append(dealResult.FungibleContract, dbFungibleContract)
			}
			if dbNonFungibleContract != nil {
				dealResult.NonFungibleContract = append(dealResult.NonFungibleContract, dbNonFungibleContract)
			}
			if dbIDAContract != nil {
				dealResult.InsertIDAContracts = append(dealResult.InsertIDAContracts, dbIDAContract)
			}

			// 更新主链主did合约
			chain.SetMainDIDContract(chainId, contractInfo)
		} else {
			contractInfo = contractDB
			contractInfo.ContractStatus = contractData.ContractStatus
			contractInfo.Upgrader = contractData.SenderAddr
			contractInfo.UpgradeAddr = contractData.SenderAddr
			contractInfo.UpgradeOrgId = contractData.OrgId
			contractInfo.UpgradeTimestamp = contractData.Timestamp
			contractInfo.Version = contractData.Version

			// 更新更新数据
			dealResult.UpdateContracts = append(dealResult.UpdateContracts, contractInfo)
		}
	}

	// 如果没有错误，直接返回
	return nil
}

func BuildContractOtherData(chainId, senderTxId string, contractInfo *db.Contract,
	dealResult *model.ProcessedBlockData) {
	var byteCodeDB []byte
	for _, byteCode := range dealResult.ContractByteCodeList {
		if senderTxId == byteCode.TxId {
			byteCodeDB = byteCode.ByteCode
			break
		}
	}

	if byteCodeDB == nil {
		byteCodeInfo, err := dbhandle.GetContractByteCodeByTx(chainId, senderTxId)
		if err != nil {
			log.Errorf("BuildContractOtherData GetContractByteCodeByTx err: %s", err.Error())
			return
		}
		if byteCodeInfo != nil {
			byteCodeDB = byteCodeInfo.ByteCode
		}
	}

	// 计算合约类型
	if byteCodeDB == nil {
		log.Errorf("BuildContractOtherData byteCodeDB is nil, senderTxId: %s", senderTxId)
		return
	}

	//字节码不为空，获取合约类型
	err := GetContractSDKData(chainId, contractInfo, byteCodeDB)
	if err != nil {
		log.Errorf("BuildContractOtherData GetContractSDKData err: %s", err.Error())
		return
	}
}
