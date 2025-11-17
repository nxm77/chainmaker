// Package chain provides chain Methods
package chain

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	loggers "chainmaker_web/src/logger"
	"encoding/json"

	"chainmaker.org/chainmaker/contract-utils/standard"
)

var (
	log = loggers.GetLogger(loggers.MODULE_SYNC)
)

// InitOther 初始化其他模块
func InitOther() {
	// 初始化链配置
	InitChainConfig()
	// 初始化链主DID
	InitChainMainDID()

	// 初始化链列表表
	InitChainListTable(config.SubscribeChains)
}

func InitChainListTable(chainList []*config.ChainInfo) {
	//chainList := config.SubscribeChains
	dbConfig := config.GlobalConfig.DBConf
	//初始化表
	switch dbConfig.DbProvider {
	case config.MySql:
		db.InitChainListMySqlTable(chainList)
	case config.Pgsql:
		db.InitChainListPgsqlTable(chainList)
	case config.ClickHouse:
		db.InitChainLsitClickHouseTable(chainList)
	}
}

// 初始化主DID合约
func InitChainMainDID() {
	// 获取全局配置中的合约配置
	chainConf := config.GlobalConfig.ChainConf
	// 获取合约配置
	if chainConf.ContractConfig == nil {
		log.Warnf("InitChainMainDID contract config is nil")
		return
	}

	// 遍历合约配置
	for _, v := range chainConf.ContractConfig {
		// 将合约配置中的主DID合约添加到主DID合约映射中
		config.MainDIDContractMap[v.ChainId] = v.MainDIDContract
	}

	log.Infof("InitChainMainDID contract config: %v", config.MainDIDContractMap)
	// 遍历订阅的链
	for _, v := range config.SubscribeChains {
		// 如果主DID合约映射中已经存在该链的主DID合约，则跳过
		if _, ok := config.MainDIDContractMap[v.ChainId]; !ok {
			//从数据库中获取该链的主DID合约
			chainMainDIDContract := dbhandle.GetChainMainDIDContract(v.ChainId)
			// 如果获取成功，则将主DID合约添加到主DID合约映射中
			if chainMainDIDContract != nil {
				config.MainDIDContractMap[v.ChainId] = chainMainDIDContract.NameBak
			}
		}
	}
}

// 设置主DID合约
func SetMainDIDContract(chainId string, contractInfo *db.Contract) {
	if contractInfo.ContractType != standard.ContractStandardNameCMDID {
		return
	}

	// 如果config.MainDIDContractMap中已经存在chainId，则直接返回
	if _, ok := config.MainDIDContractMap[chainId]; ok {
		return
	}
	// 否则，将chainId和contractName添加到config.MainDIDContractMap中
	config.MainDIDContractMap[chainId] = contractInfo.NameBak
	log.Infof("SetMainDIDContract chainId: %s, contractName: %s", chainId, contractInfo.NameBak)
}

// 判断是否为主DID合约
func CheckIsMainDIDContract(chainId, contractName string) bool {
	// 判断config.MainDIDContractMap中是否存在chainId
	if _, ok := config.MainDIDContractMap[chainId]; ok {
		// 如果存在，则判断config.MainDIDContractMap[chainId]是否等于contractName
		return config.MainDIDContractMap[chainId] == contractName
	}
	// 如果不存在，则返回false
	return false
}

// InitChainConfig
//
//	@Description: 初始化链订阅数据
func InitChainConfig() {
	//从数据库获取订阅数据
	subscribeChains, err := GetSubscribeChains()
	if len(subscribeChains) > 0 && err == nil {
		//将数据库订阅和配置订阅数据合并
		mergedChains := mergeChainInfo(config.SubscribeChains, subscribeChains)
		config.SubscribeChains = mergedChains
	}
}

// mergeChainInfo
//
//	@Description: 将DB订阅数据和配置额订阅数据合并， 合并配置数据和DB数据,相同数据用DB数据
//	@param configChains 配置的订阅数据
//	@param dbChains 数据库的订阅数据
//	@return []*config.ChainInfo 合并后的订阅数据
func mergeChainInfo(configChains, dbChains []*config.ChainInfo) []*config.ChainInfo {
	mergedChains := make([]*config.ChainInfo, 0)

	// 创建一个映射，用于存储 chains2 中的 ChainInfo，以便于根据 ChainId 进行查找
	dbChainsMap := make(map[string]*config.ChainInfo)
	for _, chain := range dbChains {
		dbChainsMap[chain.ChainId] = chain
	}

	// 遍历 configChains，如果存在相同的 ChainId，则使用 dbChains 中的数据
	for _, chain := range configChains {
		if chain2, ok := dbChainsMap[chain.ChainId]; ok {
			mergedChains = append(mergedChains, chain2)
			// 从 chains2Map 中删除已合并的 ChainInfo，以便于后续处理 chains2 中剩余的数据
			delete(dbChainsMap, chain.ChainId)
		} else {
			mergedChains = append(mergedChains, chain)
		}
	}

	// 将 dbChains 中剩余的 ChainInfo 添加到 mergedChains 中
	for _, chain := range dbChainsMap {
		mergedChains = append(mergedChains, chain)
	}

	return mergedChains
}

// GetSubscribeChains 获取订阅信息，没有就用配置文件
func GetSubscribeChains() ([]*config.ChainInfo, error) {
	var err error
	chainConfigs := make([]*config.ChainInfo, 0)
	//数据库获取订阅数据
	subscribeChains, err := dbhandle.GetDBSubscribeChains()
	if err != nil {
		log.Errorf("Init GetSubscribeChains err:%v", err)
		return nil, err
	}
	if len(subscribeChains) == 0 {
		return chainConfigs, nil
	}
	//数据库链数据
	for _, chainInfo := range subscribeChains {
		var nodeList []*config.NodeInfo
		if chainInfo.NodeList != "" {
			err = json.Unmarshal([]byte(chainInfo.NodeList), &nodeList)
			if err != nil {
				log.Errorf("chain node list json Unmarshal failed, err:%v", err)
				continue
			}
		}

		chain := &config.ChainInfo{
			ChainId:   chainInfo.ChainId,
			AuthType:  chainInfo.AuthType,
			OrgId:     chainInfo.OrgId,
			HashType:  chainInfo.HashType,
			NodesList: nodeList,
			TlsMode:   chainInfo.TlsMode,
			Tls:       chainInfo.Tls,
			UserInfo: &config.UserInfo{
				UserSignKey: chainInfo.UserSignKey,
				UserSignCrt: chainInfo.UserSignCrt,
				UserTlsKey:  chainInfo.UserTlsKey,
				UserTlsCrt:  chainInfo.UserTlsCrt,
				UserEncKey:  chainInfo.UserEncKey,
				UserEncCrt:  chainInfo.UserEncCrt,
			},
		}
		chainConfigs = append(chainConfigs, chain)
	}
	return chainConfigs, nil
}

// GetSubscribeByChainId 获取订阅信息，没有就用配置文件
func GetSubscribeByChainId(chainId string) (*config.ChainInfo, error) {
	var err error
	//数据库获取订阅数据
	chainInfoDB, err := dbhandle.GetSubscribeByChainId(chainId)
	if err != nil || chainInfoDB == nil {
		log.Errorf("Init GetSubscribeChains err:%v", err)
		return nil, err
	}

	var nodeList []*config.NodeInfo
	if chainInfoDB.NodeList != "" {
		err = json.Unmarshal([]byte(chainInfoDB.NodeList), &nodeList)
		if err != nil {
			log.Errorf("chain node list json Unmarshal failed, err:%v", err)
			return nil, err
		}
	}

	chainConfig := &config.ChainInfo{
		ChainId:   chainInfoDB.ChainId,
		AuthType:  chainInfoDB.AuthType,
		OrgId:     chainInfoDB.OrgId,
		HashType:  chainInfoDB.HashType,
		TlsMode:   chainInfoDB.TlsMode,
		Tls:       chainInfoDB.Tls,
		NodesList: nodeList,
		UserInfo: &config.UserInfo{
			UserSignKey: chainInfoDB.UserSignKey,
			UserSignCrt: chainInfoDB.UserSignCrt,
			UserTlsKey:  chainInfoDB.UserTlsKey,
			UserTlsCrt:  chainInfoDB.UserTlsCrt,
			UserEncKey:  chainInfoDB.UserEncKey,
			UserEncCrt:  chainInfoDB.UserEncCrt,
		},
	}
	return chainConfig, nil
}
