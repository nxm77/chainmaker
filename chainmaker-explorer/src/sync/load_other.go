/*
Package sync comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package sync

import (
	"chainmaker_web/src/db"
	client "chainmaker_web/src/sync/clients"
	"encoding/json"
	"errors"
	"time"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/common/v2/crypto/asym"
	pbconfig "chainmaker.org/chainmaker/pb-go/v2/config"
	sdk "chainmaker.org/chainmaker/sdk-go/v2"
	commonutils "chainmaker.org/chainmaker/utils/v2"

	"chainmaker_web/src/config"
	"chainmaker_web/src/db/dbhandle"

	"chainmaker.org/chainmaker/pb-go/v2/discovery"
)

// PeriodicCheckSubChainStatus
//
//	@Description: 检查子链健康状态， 1小时检查一次
//	@param sdkClient 链连接
// func PeriodicCheckSubChainStatus(sdkClient *client.SdkClient) {
// 	chainId := sdkClient.ChainId

// 	//1小时定时器
// 	ticker := time.NewTicker(time.Hour)
// 	defer ticker.Stop()

// 	for range ticker.C {
// 		//链订阅已经停止,停止定时器
// 		if sdkClient.Status == client.STOP {
// 			return
// 		}

// 		//查询子链列表
// 		crossSubChains, err := dbhandle.GetCrossSubChainAll(chainId)
// 		if err != nil {
// 			log.Errorf("[load_other] get cross sub chain failed:%v", err)
// 		}

// 		for _, subChain := range crossSubChains {
// 			//Grpc-获取子链健康状态
// 			chainOk, errGrpc := common.CheckSubChainStatus(subChain)
// 			if errGrpc != nil {
// 				subChainJson, _ := json.Marshal(subChain)
// 				log.Errorf("[load_other] CheckSubChainStatus failed, err:%v, subChainJson:%v",
// 					errGrpc, string(subChainJson))
// 			}

// 			status := dbhandle.SubChainStatusSuccess
// 			if !chainOk {
// 				status = dbhandle.SubChainStatusFail
// 			}

// 			//为获取同步区块高度合约
// 			var spvContractName string
// 			if subChain.SpvContractName == "" {
// 				subChainInfo, errRpc := utils.GetCrossSubChainInfo(subChain.SubChainId)
// 				if errRpc != nil || subChainInfo == nil {
// 					log.Errorf("【load_other】http get sub chain failed, err:%v, SubChainId:%v",
// 						err, subChain.SubChainId)
// 				} else {
// 					spvContractName = subChainInfo.SpvContractName
// 				}
// 			}

// 			//健康状态变更，更新数据库
// 			if status != subChain.Status || spvContractName != "" {
// 				err = dbhandle.UpdateCrossSubChainStatus(chainId, subChain.SubChainId, spvContractName, status)
// 				if err != nil {
// 					log.Errorf("[load_other] update cross sub chain status failed, err:%v", err)
// 				}
// 			}
// 		}

// 	}
// }

// PeriodicLoadStart
//
//	@Description:  1小时请求一次，检查是否有新增节点
//	@param sdkClient
func PeriodicLoadStart(sdkClient *client.SdkClient) {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		//查询链状态
		newClient := client.GetSdkClient(sdkClient.ChainId)
		if newClient == nil || newClient.Status == client.STOP {
			return
		}
		loadNodeInfo(newClient)
	}
}

// loadChainRefInfos
//
//	@Description: load chain and other information
//	@param sdkClient
//	@return error
func loadChainRefInfos(sdkClient *client.SdkClient) error {
	var err error
	//处理节点数据
	loadNodeInfo(sdkClient)
	if sdkClient.ChainInfo.AuthType == config.PUBLIC {
		//处理user数据
		err = loadChainUser(sdkClient)
	} else {
		//处理组织数据
		err = loadOrgInfo(sdkClient)
	}

	if err != nil {
		return err
	}

	return nil
}

// loadChainUser loadChainUser
// @Description: 更新user数据
// @param sdkClient
// @return error
func loadChainUser(sdkClient *client.SdkClient) error {
	if sdkClient == nil || sdkClient.ChainConfig == nil {
		return errors.New("sdkClient is nil")
	}

	//链配置信息
	chainConfig := sdkClient.ChainConfig
	chainId := chainConfig.ChainId
	userList := GetAdminUserByConfig(chainConfig)
	return dbhandle.BatchInsertUser(chainId, userList)
}

// GetAdminUserByConfig
//
//	@Description: 根据链配置获取user信息
//	@param chainConfig 链配置
//	@return []*db.User user列表
func GetAdminUserByConfig(chainConfig *pbconfig.ChainConfig) []*db.User {
	userList := make([]*db.User, 0)
	if chainConfig == nil || chainConfig.ChainId == "" || len(chainConfig.TrustRoots) <= 0 {
		return userList
	}

	hashType := chainConfig.Crypto.Hash
	for _, root := range chainConfig.TrustRoots[0].Root {
		publicKey, err := asym.PublicKeyFromPEM([]byte(root))
		if err != nil {
			log.Error("[SDK] get publicKey by PK err : " + err.Error())
			continue
		}
		addr, err := commonutils.PkToAddrStr(publicKey, pbconfig.AddrType_ETHEREUM, crypto.HashAlgoMap[hashType])
		if err != nil {
			log.Error("[SDK] get addr by PK err : " + err.Error())
			continue
		}
		//userId, err := helper.CreateLibp2pPeerIdWithPublicKey(publicKey)
		//if err != nil {
		//	continue
		//}

		user := &db.User{
			UserId:    addr,
			UserAddr:  addr,
			Role:      "admin",
			OrgId:     config.PUBLIC,
			Timestamp: time.Now().Unix(),
		}
		userList = append(userList, user)
	}
	return userList
}

// loadOrgInfo
//
//	@Description: 加载组织数据
//	@param sdkClient
//	@return error
func loadOrgInfo(sdkClient *client.SdkClient) error {
	if sdkClient == nil {
		return errors.New("sdkClient is nil")
	}

	chainConfig := sdkClient.ChainConfig
	return dbhandle.SaveOrgByConfig(chainConfig)
}

// loadNodeInfo
//
//	@Description: 加载节点数据
//	@param sdkClient
func loadNodeInfo(sdkClient *client.SdkClient) {
	chainId := sdkClient.ChainId
	if chainId == "" {
		return
	}

	//获取所有的链节点
	allNodeList, err1 := GetAllNodeList(chainId, sdkClient.ChainClient)
	if err1 != nil {
		log.Warnf("[sync-%v] loadNodeInfo GetAllNodeList Failed : err:%v", chainId, err1)
	}

	//获取链配置的节点数据
	chainConfigNodeList, err2 := GetChainNodeList(sdkClient.ChainClient)
	if err2 != nil {
		log.Warnf("[sync-%v] loadNodeInfo GetChainNodeList Failed : err:%v", chainId, err2)
	}

	//获取链配置的节点数据
	consensusNodeOrgList, err3 := GetAllConsensusNodeOrgId(sdkClient)
	if err3 != nil {
		log.Warnf("[sync-%v] loadNodeInfo GetChainNodeList Failed : err:%v", chainId, err3)
	}

	if err1 == nil && err2 == nil && err3 == nil {
		// 数据库节点
		nodeIds, err := dbhandle.GetNodesRef(chainId)
		log.Infof("[sync-%v] loadNodeInfo Get node DB, nodeIds:%v", chainId, nodeIds)
		if err != nil {
			log.Warnf("[sync-%v] loadNodeInfo Get node DB fail, err:%v", chainId, err)
			return
		}

		var deleteNodeIds []string
		// 链上没有的节点在数据库中删除
		for _, dbNodeId := range nodeIds {
			_, ok1 := allNodeList[dbNodeId]
			_, ok2 := chainConfigNodeList[dbNodeId]
			_, ok3 := consensusNodeOrgList[dbNodeId]
			if !ok1 && !ok2 && !ok3 {
				deleteNodeIds = append(deleteNodeIds, dbNodeId)
			}
		}

		//删除节点数据1
		if len(deleteNodeIds) > 0 {
			err := dbhandle.DeleteNodeById(chainId, deleteNodeIds)
			if err != nil {
				log.Warnf("[sync] loadNodeInfo Delete Node Info Failed : ", err)
				return
			}
		}
	}

	nodes := parseAllNodeList(allNodeList, consensusNodeOrgList, chainConfigNodeList)
	nodeJson, _ := json.Marshal(nodes)
	log.Infof("[sync-%v] loadNodeInfo parseAllNodeList, nodes:%s", chainId, string(nodeJson))

	err := dbhandle.BatchInsertNode(chainId, nodes)
	if err != nil {
		log.Warnf("[DB] Update Node Info Failed err:" + err.Error())
	}
}

// parseNodeList parse node information
func parseAllNodeList(nodeList map[string]string, consensusNodeOrgList map[string]string,
	chainConfigNodeList map[string]*discovery.Node) []*db.Node {
	nodeMap := make(map[string]*db.Node, 0)
	nodeDBList := make([]*db.Node, 0)
	for nodeId, orgId := range consensusNodeOrgList {
		var address string
		if v, ok := chainConfigNodeList[nodeId]; ok {
			address = v.GetNodeAddress()
		}
		nodeMap[nodeId] = &db.Node{
			NodeId:  nodeId,
			Role:    "consensus",
			OrgId:   orgId,
			Address: address,
		}
	}

	for nodeId, v := range chainConfigNodeList {
		if _, ok := nodeMap[nodeId]; ok {
			continue
		}
		nodeMap[nodeId] = &db.Node{
			NodeId:  nodeId,
			Role:    "common",
			OrgId:   "",
			Address: v.GetNodeAddress(),
		}
	}

	for _, nodeId := range nodeList {
		if _, ok := nodeMap[nodeId]; ok {
			continue
		}
		nodeMap[nodeId] = &db.Node{
			NodeId: nodeId,
			Role:   "common",
		}
	}

	for _, node := range nodeMap {
		nodeDBList = append(nodeDBList, node)
	}

	return nodeDBList
}

// GetCha
// GetChainNodeData
//
//	@Description:  获取链上节点和需要删除的节点
//	@param sdkClient
//	@return []*discovery.Node
//	@return []string
func GetAllConsensusNodeOrgId(sdkClient *client.SdkClient) (map[string]string, error) {
	if sdkClient == nil || sdkClient.ChainConfig == nil {
		return nil, errors.New("sdkClient is nil")
	}

	consensusNodeOrgId := make(map[string]string)
	chainConfig := sdkClient.ChainConfig
	consensusNodes := chainConfig.GetConsensus().Nodes
	log.Infof("[sync] ChainConfig consensusNodes:%v", consensusNodes)
	if len(consensusNodes) > 0 {
		authType := sdkClient.ChainInfo.AuthType
		if authType == config.PUBLIC {
			for _, nodeId := range consensusNodes[0].NodeId {
				consensusNodeOrgId[nodeId] = nodeId
			}
		} else {
			for _, node := range consensusNodes {
				for _, nodeId := range node.NodeId {
					consensusNodeOrgId[nodeId] = node.OrgId
				}
			}
		}
	}
	return consensusNodeOrgId, nil
}

// GetChainNodeData
//
//	@Description:  获取链上节点和需要删除的节点
//	@param sdkClient
//	@return []*discovery.Node
//	@return []string
func GetAllNodeList(chainId string, chainClient *sdk.ChainClient) (map[string]string, error) {
	if chainClient == nil {
		return nil, errors.New("chainClient is nil")
	}

	chainNodeIdMap := make(map[string]string)
	//获取所有的链节点
	resp, err := chainClient.GetSyncState(true)
	log.Infof("[sync-%v] ChainClient GetSyncState resp:%v", chainId, resp)
	if err != nil {
		log.Warnf("[sync-%v] ChainClient GetSyncState Failed : err:%v", chainId, err)
		return chainNodeIdMap, err
	}

	if resp == nil || len(resp.Others) == 0 {
		return chainNodeIdMap, nil
	}

	for _, value := range resp.Others {
		chainNodeIdMap[value.NodeId] = value.NodeId
	}
	return chainNodeIdMap, nil
}

// GetChainNodeData
//
//	@Description:  获取链上节点和需要删除的节点
//	@param sdkClient
//	@return []*discovery.Node
//	@return []string
func GetChainNodeList(chainClient *sdk.ChainClient) (map[string]*discovery.Node, error) {
	if chainClient == nil {
		return nil, errors.New("chainClient is nil")
	}

	chainNodeMap := make(map[string]*discovery.Node)
	//获取链信息
	chainInfo, err := chainClient.GetChainInfo()
	log.Infof("[sync] ChainClient GetChainInfo resp:%v", chainInfo)
	if err != nil {
		return nil, err
	}

	if chainInfo == nil || len(chainInfo.NodeList) == 0 {
		return nil, nil
	}

	for _, value := range chainInfo.NodeList {
		chainNodeMap[value.NodeId] = value
	}
	return chainNodeMap, nil
}

// GetChainNodeData
//
//	@Description:  获取链上节点和需要删除的节点
//	@param sdkClient
//	@return []*discovery.Node
//	@return []string
func GetChainNodeData(sdkClient *client.SdkClient) ([]*discovery.Node, []string) {
	if sdkClient == nil || sdkClient.ChainClient == nil {
		return nil, nil
	}

	deleteNodeIds := make([]string, 0)
	chainNodeIdMap := make(map[string]string)
	chainId := sdkClient.ChainId
	//链配置
	chainClient := sdkClient.ChainClient
	//获取链上节点数据
	var nodeList []*discovery.Node
	chainInfo, err := chainClient.GetChainInfo()
	if err != nil || chainInfo == nil || len(chainInfo.NodeList) == 0 {
		log.Errorf("[SDK] Get Chain Info Failed : %v", err)
		return nodeList, deleteNodeIds
	}

	log.Infof("【loadNodeInfo】chainClient GetChainInfo, chainInfo:%v", chainInfo)
	nodeList = chainInfo.NodeList
	//if len(chainInfo.NodeList) == 0 {
	//	//获取默认节点
	//	nodeList, err = DealChainConfigError(sdkClient)
	//	if err != nil || len(nodeList) == 0 {
	//		log.Infof("Get chain node Failed : %v", err)
	//		return nodeList, deleteNodeIds
	//	}
	//}

	//链上节点
	for _, node := range nodeList {
		chainNodeIdMap[node.NodeId] = node.NodeId
	}

	//数据库节点
	nodeIds, err := dbhandle.GetNodesRef(chainId)
	if err != nil {
		log.Errorf("Get nodeIds fail, err:%v", err)
		return nodeList, deleteNodeIds
	}

	//链上没有的节点在数据库中删除
	for _, dbNodeId := range nodeIds {
		if _, ok := chainNodeIdMap[dbNodeId]; !ok {
			deleteNodeIds = append(deleteNodeIds, dbNodeId)
		}
	}
	return nodeList, deleteNodeIds
}

// // parseNodeInfo parse node information
// func parseNodeInfo(node *discovery.Node) (string, []string, []string, error) {
// 	// return OrgId/Role
// 	_, rest := pem.Decode(node.GetNodeTlsCert())
// 	if rest == nil {
// 		log.Error("can not decode tls cert")
// 		return "", nil, nil, errors.New("can not decode tls cert")
// 	}
// 	cert, err := x509.ParseCertificate(rest)
// 	if err != nil {
// 		return "", nil, nil, err
// 	}
// 	return cert.Subject.CommonName, cert.Subject.Organization, cert.Subject.OrganizationalUnit, nil
// }
