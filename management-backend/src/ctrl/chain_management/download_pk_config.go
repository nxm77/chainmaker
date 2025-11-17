package chain_management

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	loggers "management_backend/src/logger"
	"os"
	"strconv"
	"strings"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/common/v2/crypto/asym"
	hashAlo "chainmaker.org/chainmaker/common/v2/crypto/hash"
	"chainmaker.org/chainmaker/pb-go/v2/consensus"
	"github.com/mr-tron/base58/base58"
	"gopkg.in/yaml.v2"

	"management_backend/src/config"
	"management_backend/src/db/chain"
	"management_backend/src/db/chain_participant"
	"management_backend/src/db/common"
	"management_backend/src/db/relation"
	"management_backend/src/global"
	"management_backend/src/utils"
)

// PKVersion pk version
const PKVersion = "v2.3.5"

// ChainPkConfig chain pk config
type ChainPkConfig struct {
	Chain      *common.Chain
	Nodes      []*common.ChainOrgNode
	Admins     []*common.ChainUser
	AdminCerts []*common.Cert
	Seeds      []string
	Consensus  *config.ConsensusPkBcConf
	DposNode   map[string]string
}

// CreatePkConfig create pk config
//
//	@Description:
//	@param chainId
//	@param confYml
//	@return chainName
//	@return err
func CreatePkConfig(chainId, confYml string) (chainName string, err error) {
	chainName = chainId
	pkConfig := ChainPkConfig{
		Consensus: &config.ConsensusPkBcConf{},
		DposNode:  map[string]string{},
	}
	// 获取chain
	pkConfig.Chain, err = chain.GetChainByChainId(chainId)
	if err != nil {
		loggers.WebLogger.Error("GetChainByChainId err : " + err.Error())
		return
	}
	// 获取节点
	pkConfig.Nodes, err = relation.GetChainOrgByChainIdList(chainId)
	if err != nil {
		loggers.WebLogger.Error("GetChainOrgByChainIdList err : " + err.Error())
		return
	}
	// 获取用户
	pkConfig.Admins, err = relation.GetChainUserByChainId(chainId, "")
	if err != nil {
		loggers.WebLogger.Error("GetChainUser err : " + err.Error())
		return
	}
	// 配置基本信息
	err = pkConfig.DealBaseInfo()
	if err != nil {
		loggers.WebLogger.Error("GetChainByChainId err : " + err.Error())
		return
	}
	var nodePaths string
	var nodeCert *common.Cert
	// 创建node对应的文件
	for _, node := range pkConfig.Nodes {
		err = pkConfig.CreateBc(node.NodeName, confYml)
		if err != nil {
			loggers.WebLogger.Error("CreateBc err : " + err.Error())
		}
		err = pkConfig.CreateChainMaker(node, confYml)
		if err != nil {
			loggers.WebLogger.Error("CreateChainMaker err : " + err.Error())
		}
		// 获取用户
		nodeCert, err = chain_participant.GetPemCert(node.NodeName)
		if err != nil {
			loggers.WebLogger.Error("GetChainUser err : " + err.Error())
			return
		}
		err = pkConfig.createBinAndLib(chainId, node.NodeName, nodeCert.Addr, confYml)
		if err != nil {
			loggers.WebLogger.Error("createBinAndLib err : " + err.Error())
		}
		err = pkConfig.CreateCert(node)
		if err != nil {
			loggers.WebLogger.Error("CreateCert err : " + err.Error())
		}
		pv := fmt.Sprintf("./chainmaker-%v-", PKVersion)
		nodePaths = nodePaths + pv + node.NodeName + ","
	}
	nodePaths = strings.TrimRight(nodePaths, ",")
	// 创建chainmaker
	err = pkConfig.createChainmakerAndScript(chainName, confYml, nodePaths)
	if err != nil {
		loggers.WebLogger.Error("createChainmakerAndScript err : " + err.Error())
		return
	}
	return
}

// DealBaseInfo deal base info
//
//	@Description:
//	@receiver c
//	@return error
func (c *ChainPkConfig) DealBaseInfo() error {
	err := os.RemoveAll("release/")
	if err != nil {
		loggers.WebLogger.Error(err.Error())
		return err
	}
	adminMap := make(map[string]*common.Cert)
	if c.AdminCerts == nil || len(c.AdminCerts) <= 0 {
		for _, admin := range c.Admins {
			userInfo, userInfoErr := chain_participant.GetUserCertByCertUse(admin.UserName, chain_participant.PEM)
			if userInfoErr != nil {
				loggers.WebLogger.Error("GetUserTlsCert err : " + userInfoErr.Error())
			}
			adminMap[userInfo.RemarkName] = userInfo
			c.AdminCerts = append(c.AdminCerts, userInfo)
		}
	}

	nodes := &config.NodesConf{
		OrgId: "public",
	}
	for _, node := range c.Nodes {
		nodes.NodeId = append(nodes.NodeId, node.NodeId)
		var nodeIp string
		if c.Chain.Single == SINGLE {
			nodeIp = LOCAL_IP
		} else {
			nodeIp = node.NodeIp
		}
		c.Seeds = append(c.Seeds, "/ip4/"+nodeIp+"/tcp/"+strconv.Itoa(node.NodeP2pPort)+"/p2p/"+node.NodeId)
	}
	c.Consensus.Nodes = []*config.NodesConf{nodes}
	c.Consensus.DposConfig = []*config.KvConf{}
	c.Consensus.Type = consensus.ConsensusType_value[c.Chain.Consensus]
	// dpos
	if c.Consensus.Type == int32(consensus.ConsensusType_DPOS) {
		nodes.OrgId = "dpos_org_id"
		var stakes []*Stake
		nodesMap := make(map[string]string, len(c.Nodes))
		for _, node := range c.Nodes {
			nodesMap[node.NodeName] = node.NodeId
		}
		_ = json.Unmarshal([]byte(c.Chain.Stakes), &stakes)
		c.Consensus.DposConfig = append(c.Consensus.DposConfig, &config.KvConf{Key: "erc20.decimals", Value: "18"})
		c.Consensus.DposConfig = append(c.Consensus.DposConfig,
			&config.KvConf{Key: "stake.minSelfDelegation", Value: strconv.Itoa(c.Chain.StakeMinCount)})
		c.Consensus.DposConfig = append(c.Consensus.DposConfig,
			&config.KvConf{Key: "stake.epochValidatorNum", Value: strconv.Itoa(len(stakes))})
		c.Consensus.DposConfig = append(c.Consensus.DposConfig,
			&config.KvConf{Key: "stake.epochBlockNum", Value: "10"})
		c.Consensus.DposConfig = append(c.Consensus.DposConfig,
			&config.KvConf{Key: "stake.completionUnbondingEpochNum", Value: "1"})
		dposStake := 0
		var nodeIds, candidates []*config.KvConf
		var addr string
		for _, stake := range stakes {
			var admin *common.Cert
			var ok bool
			if admin, ok = adminMap[stake.RemarkName]; !ok {
				admin, err = chain_participant.GetUserCertByCertUse(stake.RemarkName, chain_participant.PEM)
				if err != nil {
					return err
				}
			}
			pk, err := asym.PublicKeyFromPEM([]byte(admin.PublicKey))
			if err != nil {
				return err
			}
			pkStr, err := pk.String()
			if err != nil {
				return err
			}
			pubkey := []byte(pkStr)
			var hashType string
			if admin.Algorithm == global.ECDSA {
				hashType = crypto.CRYPTO_ALGO_SHA256
			} else {
				hashType = crypto.CRYPTO_ALGO_SM3
			}
			var hashBz []byte
			if hashBz, err = hashAlo.GetByStrType(hashType, pubkey); err != nil {
				return err
			}
			// 赋值dpos地址
			admin.Addr = base58.Encode(hashBz[:])
			nodeIds = append(nodeIds, &config.KvConf{Key: "stake.nodeID:" + admin.Addr,
				Value: nodesMap[stake.NodeName]})
			candidates = append(candidates, &config.KvConf{Key: "stake.candidate:" + admin.Addr,
				Value: strconv.Itoa(stake.Count)})
			c.DposNode[stake.NodeName] = admin.Addr
			if addr == "" {
				addr = admin.Addr
			}
			dposStake = dposStake + stake.Count
		}
		c.Consensus.DposConfig = append(c.Consensus.DposConfig, &config.KvConf{Key: "erc20.total",
			Value: strconv.Itoa(dposStake)})
		c.Consensus.DposConfig = append(c.Consensus.DposConfig, &config.KvConf{Key: "erc20.account:DPOS_STAKE",
			Value: strconv.Itoa(dposStake)})
		c.Consensus.DposConfig = append(c.Consensus.DposConfig, nodeIds...)
		c.Consensus.DposConfig = append(c.Consensus.DposConfig, candidates...)
		c.Consensus.DposConfig = append(c.Consensus.DposConfig, &config.KvConf{Key: "erc20.owner", Value: addr})

	}

	return nil
}

// CreateBc create bc
//
//	@Description:
//	@receiver c
//	@param nodeName
//	@param confYml
//	@return error
func (c *ChainPkConfig) CreateBc(nodeName, confYml string) error {
	bcConfig := &config.PkBc{}
	bcFile, err := ioutil.ReadFile(confYml + "/config_pk_tpl/chainconfig/bc1.yml")
	if err != nil {
		loggers.WebLogger.Error(err.Error())
		return err
	}
	_ = yaml.Unmarshal(bcFile, bcConfig)
	bcConfig.Crypto.Hash = c.Chain.CryptoHash
	bcConfig.Block.TxTimeout = c.Chain.TxTimeout
	bcConfig.Block.BlockTxCapacity = c.Chain.BlockTxCapacity
	bcConfig.Block.BlockInterval = int(c.Chain.BlockInterval)
	bcConfig.Consensus = c.Consensus
	bcConfig.ChainId = c.Chain.ChainId
	trustRoots := &config.TrustRootsConf{
		OrgId: "public",
	}
	if c.Chain.CryptoHash != "" {
		bcConfig.Crypto.Hash = c.Chain.CryptoHash
	}
	trustList := []string{}
	for _, admin := range c.Admins {
		trustList = append(trustList, fmt.Sprintf("../config/%v/admin/%v/%v.pem", nodeName, admin.UserName, admin.UserName))
	}
	trustRoots.Root = trustList
	bcConfig.TrustRoots = []*config.TrustRootsConf{trustRoots}
	err = os.MkdirAll(fmt.Sprintf("release/chainmaker-%v-%v/config/%v/chainconfig",
		PKVersion, nodeName, nodeName), os.ModePerm)
	if err != nil {
		loggers.WebLogger.Error("Mkdir bc1 path err : " + err.Error())
		return err
	}
	bcBytes, _ := yaml.Marshal(bcConfig)
	err = utils.CreateAndRename("bc1.yaml", fmt.Sprintf("release/chainmaker-%v-%v/"+
		"config/%v/chainconfig/bc1.yml", PKVersion, nodeName, nodeName), string(bcBytes))
	if err != nil {
		loggers.WebLogger.Error("create and rename crt file err :", err.Error())
	}
	return err
}

// createBinAndLib
//
//	@Description:
//	@receiver c
//	@param nodeName
//	@param confYml
//	@return error
func (c *ChainPkConfig) createBinAndLib(chainId, nodeName, addr string, confYml string) error {
	versionPath := fmt.Sprintf("release/chainmaker-%v-%v", PKVersion, nodeName)
	binPath := versionPath + "/bin"
	err := os.MkdirAll(binPath, os.ModePerm)
	if err != nil {
		loggers.WebLogger.Error("Mkdir bin path err : " + err.Error())
		return err
	}

	restartPath := binPath + "/restart"
	stopPath := binPath + "/stop"
	dockerStartPath := binPath + "/docker_start"

	err = utils.CreateAndCopy(restartPath, confYml+"/bin/restart.sh", 0777)
	if err != nil {
		loggers.WebLogger.Error("create and copy restart.sh file err :", err.Error())
		return err
	}
	dockerEnable := c.Chain.DockerVm == DOCKER_VM
	replace := map[string]string{
		"{org_id}":        nodeName,
		"{docker_enable}": strconv.FormatBool(dockerEnable),
		"{chain_id}":      chainId,
		"{node_addr}":     nodeName,
	}
	err = utils.RePlaceMore(restartPath, replace)
	if err != nil {
		loggers.WebLogger.Error("rePlace bin/restart.sh err : " + err.Error())
	}
	err = utils.CreateAndCopy(stopPath, confYml+"/bin/stop.sh", 0777)
	if err != nil {
		loggers.WebLogger.Error("create and copy stop.sh file err :", err.Error())
		return err
	}
	if c.Chain.DockerVm == DOCKER_VM {
		err = utils.CreateAndCopy(dockerStartPath, confYml+"/bin/docker_start.sh", 0777)
		if err != nil {
			loggers.WebLogger.Error("create and copy stop.sh file err :", err.Error())
			return err
		}
		dockerReplace := map[string]string{
			"{org_id}":        nodeName,
			"{docker_enable}": strconv.FormatBool(dockerEnable),
			"{chain_id}":      chainId,
			"{node_addr}":     addr,
		}
		err = utils.RePlaceMore(dockerStartPath, dockerReplace)
		if err != nil {
			loggers.WebLogger.Error("rePlace bin/docker_start.sh err : " + err.Error())
		}
		err = utils.RePlaceMore(stopPath, dockerReplace)
		if err != nil {
			loggers.WebLogger.Error("rePlace bin/stop.sh err : " + err.Error())
		}
	} else {
		err = utils.RePlaceMore(stopPath, replace)
		if err != nil {
			loggers.WebLogger.Error("rePlace bin/stop.sh err : " + err.Error())
		}
	}

	libPath := versionPath + "/lib"
	err = os.MkdirAll(libPath, os.ModePerm)
	if err != nil {
		loggers.WebLogger.Error("Mkdir lib path err : " + err.Error())
		return err
	}
	err = utils.CreateAndCopy(libPath+"/libwasmer.dylib", confYml+"/lib/libwasmer.dylib", 0777)
	if err != nil {
		loggers.WebLogger.Error("create and copy libwasmer.dylib file err :", err.Error())
		return err
	}
	err = utils.CreateAndCopy(libPath+"/libwasmer.so", confYml+"/lib/libwasmer.so", 0777)
	if err != nil {
		loggers.WebLogger.Error("create and copy libwasmer.dylib file err :", err.Error())
		return err
	}
	err = utils.CreateAndCopy(libPath+"/wxdec", confYml+"/lib/wxdec", 0777)
	if err != nil {
		loggers.WebLogger.Error("create and copy wxdec file err :", err.Error())
		return err
	}
	return err
}

// CreateChainMaker create chainmaker
//
//	@Description:
//	@receiver c
//	@param node
//	@param confYml
//	@return error
func (c *ChainPkConfig) CreateChainMaker(node *common.ChainOrgNode, confYml string) error {
	chainMaker := &config.PKChainmaker{}
	yamlFile, _ := ioutil.ReadFile(confYml + "/config_pk_tpl/chainmaker.yml")
	_ = yaml.Unmarshal(yamlFile, chainMaker)
	chainMaker.NetConf.Seeds = c.Seeds
	if c.Chain.EnableHttp == 1 {
		chainMaker.RpcConf.GatewayConf.Enabled = true
	}
	chainMaker.ChainLogConf.ConfigFile = strings.Replace(chainMaker.ChainLogConf.ConfigFile,
		"{org_path}", node.NodeName, -1)
	blockchainConf := config.BlockchainConf{
		ChainId: c.Chain.ChainId,
		Genesis: fmt.Sprintf("../config/%v/chainconfig/bc1.yml", node.NodeName),
	}
	chainMaker.BlockchainConf = []*config.BlockchainConf{&blockchainConf}
	chainMaker.NodeConf = &config.NodePkConf{
		PrivKeyFile:   fmt.Sprintf("../config/%v/%v.key", node.NodeName, node.NodeName),
		CertCacheSize: chainMaker.NodeConf.CertCacheSize,
		Pkcs11:        chainMaker.NodeConf.Pkcs11,
		FastSync:      chainMaker.NodeConf.FastSync,
	}
	chainMaker.NetConf.ListenAddr = strings.Replace(chainMaker.NetConf.ListenAddr,
		"{net_port}", strconv.Itoa(node.NodeP2pPort), -1)
	chainMaker.RpcConf.Port = node.NodeRpcPort
	chainMaker.NetConf.Tls.PrivKeyFile = fmt.Sprintf("../config/%v/%v.key", node.NodeName, node.NodeName)
	if c.Chain.DockerVm == NO_DOCKER_VM {
		chainMaker.VmConf.DockerGo.Enable = false
	} else {
		chainMaker.VmConf.DockerGo.Enable = true
		chainMaker.VmConf.DockerGo.RuntimeServer.Port = node.NodeRpcPort + 20050
		chainMaker.VmConf.DockerGo.ContractEngine.Port = node.NodeRpcPort + 10050
	}
	chainmakerBytes, _ := yaml.Marshal(chainMaker)
	chainmakerStr := strings.Replace(string(chainmakerBytes), "{org_id}", node.NodeName, -1)
	err := utils.CreateAndRename("chainmaker.yml",
		fmt.Sprintf("release/chainmaker-%v-%v/config/%v/chainmaker.yml",
			PKVersion, node.NodeName, node.NodeName), chainmakerStr)
	if err != nil {
		loggers.WebLogger.Error("create and copy crt file err :", err.Error())
		return err
	}
	err = utils.CreateAndCopy(fmt.Sprintf("release/chainmaker-%v-%v/config/%v/log.yml",
		PKVersion, node.NodeName, node.NodeName), confYml+"/config_tpl/log.yml", 0)
	if err != nil {
		loggers.WebLogger.Error("create and copy crt file err :", err.Error())
	}
	return err
}

// CreateCert create cert
//
//	@Description:
//	@receiver c
//	@param node
//	@return error
func (c *ChainPkConfig) CreateCert(node *common.ChainOrgNode) error {
	path := fmt.Sprintf("release/chainmaker-%v-%v/config/%v", PKVersion, node.NodeName, node.NodeName)
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		loggers.WebLogger.Error("Mkdir org certs/node path err : " + err.Error())
	}
	err = utils.CreateAndRename(node.NodeName+".nodeid",
		fmt.Sprintf("release/chainmaker-%v-%v/config/%v/%v.nodeid",
			PKVersion, node.NodeName, node.NodeName, node.NodeName), node.NodeId)
	if err != nil {
		loggers.WebLogger.Error("create and rename crt file err :", err.Error())
		return err
	}
	nodeCert, err := chain_participant.GetUserCertByCertUse(node.NodeName, chain_participant.PEM)
	if err != nil {
		loggers.WebLogger.Error("GetNodeCert erBlockr : " + err.Error())
		return err
	}
	if nodeCert != nil {
		err = utils.CreateAndRename(node.NodeName+".pem", path+"/"+node.NodeName+".pem", nodeCert.PublicKey)
		if err != nil {
			loggers.WebLogger.Error("create and rename crt file err :", err.Error())
			return err
		}
		err = utils.CreateAndRename(node.NodeName+".key", path+"/"+node.NodeName+".key", nodeCert.PrivateKey)
		if err != nil {
			loggers.WebLogger.Error("create and rename kry file err :", err.Error())
		}
		err = utils.CreateAndRename(node.NodeName+".addr", path+"/"+node.NodeName+".addr", nodeCert.Addr)
		if err != nil {
			loggers.WebLogger.Error("create and rename addr file err :", err.Error())
		}
	}

	for _, cert := range c.AdminCerts {
		userPath := fmt.Sprintf(path+"/admin/%v", cert.RemarkName)
		err = os.MkdirAll(userPath, os.ModePerm)
		if err != nil {
			loggers.WebLogger.Error("Mkdir org certs/user path err : " + err.Error())
		}
		err = utils.CreateAndRename(cert.RemarkName+".pem", userPath+"/"+cert.RemarkName+".pem", cert.PublicKey)
		if err != nil {
			loggers.WebLogger.Error("create and rename crt file err :", err.Error())
			return err
		}
		err = utils.CreateAndRename(cert.RemarkName+".key", userPath+"/"+cert.RemarkName+".key", cert.PrivateKey)
		if err != nil {
			loggers.WebLogger.Error("create and rename kry file err :", err.Error())
		}
		err = utils.CreateAndRename(cert.RemarkName+".addr", userPath+"/"+cert.RemarkName+".addr", cert.Addr)
		if err != nil {
			loggers.WebLogger.Error("create and rename kry file err :", err.Error())
		}
	}
	return nil
}

// createChainmakerAndScript
//
//	@Description:
//	@receiver c
//	@param chainName
//	@param confYml
//	@param nodePaths
//	@return error
func (c *ChainPkConfig) createChainmakerAndScript(chainName, confYml, nodePaths string) error {
	err := utils.CreateAndCopy("release/chainmaker", confYml+"/bin/chainmaker", 0777)
	if err != nil {
		loggers.WebLogger.Error("create and copy chainmaker file err :", err.Error())
		return err
	}
	if c.Chain.Monitor == MONITOR_START {
		err = utils.CreateAndCopy("release/cmlogagentd", confYml+"/bin/cmlogagentd", 0777)
		if err != nil {
			loggers.WebLogger.Error("create and copy chainmaker file err :", err.Error())
			return err
		}
		err = utils.CreateAndCopy("release/start", confYml+"/bin/logagentd_start.sh", 0777)
		if err != nil {
			loggers.WebLogger.Error("CopyFile bin/start.sh err : " + err.Error())
		}
		err = utils.RePlace("release/start", "{node_paths}", nodePaths)
		if err != nil {
			loggers.WebLogger.Error("rePlace release/start.sh err : " + err.Error())
		}
	} else {
		err = utils.CreateAndCopy("release/start.sh", confYml+"/bin/start.sh", 0777)
		if err != nil {
			loggers.WebLogger.Error("create and copy start.sh file err :", err.Error())
			return err
		}
	}
	err = utils.CreateAndCopy("release/quick_stop.sh", confYml+"/bin/quick_stop.sh", 0777)
	if err != nil {
		loggers.WebLogger.Error("create and copy quick_stop.sh file err :", err.Error())
		return err
	}
	err = utils.Zip("release", "./chain_config/"+chainName+".zip")
	if err != nil {
		loggers.WebLogger.Error("zip file err :", err.Error())
	}
	return err
}
