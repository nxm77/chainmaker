# JavaSDK 接口说明
## 合约管理接口
### 创建合约payload
**参数说明**
   - contractName: 合约名
   - version: 版本号
   - byteCodes: 合约字节数组
   - runtimeType: 合约运行环境
   - params: 合约初始化参数
```java
Request.Payload createContractCreatePayload(
        String contractName,
        String version, byte[] byteCode,
        ContractOuterClass.RuntimeType runtime,
        Map<String, byte[]> params)
```
### 创建升级合约payload
**参数说明**
   - contractName: 合约名
   - version: 版本号
   - byteCodes: 合约字节数组
   - runtimeType: 合约运行环境
   - params: 合约初始化参数
```java
Request.Payload createContractUpgradePayload(
        String contractName, 
        String version, byte[] byteCode, 
        ContractOuterClass.RuntimeType runtime,
        Map<String, byte[]> params)
```
### 创建冻结合约payload
**参数说明**
   - contractName: 合约名
```java
Request.Payload createContractFreezePayload(String contractName) 
```
### 创建解冻合约payload
**参数说明**
   contractName: 合约名
```java
Request.Payload createContractUnFreezePayload(String contractName)
```
### 创建吊销合约payload
**参数说明**
   contractName: 合约名
```java
Request.Payload createContractRevokePayload(String contractName) 
```
### 发送合约操作请求（创建、更新、冻结、解冻、吊销）
**参数说明**
   - payload: 交易payload
   - endorsementEntries: 背书签名信息列表
   - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
   - syncResultTimeout: 同步获取执行结果超时时间，小于等于0代表不等待执行结果，直接返回(返回信息里包含交易ID)，单位：毫秒
```java
ResultOuterClass.TxResponse sendContractManageRequest(
        Request.Payload payload,
        Request.EndorsementEntry[] endorsementEntries,
        long rpcCallTimeout, long syncResultTimeout)
```   
### 合约调用
**参数说明**
   - contractName: 合约名
   - method: 方法名
   - txId: 交易id
      - 格式要求： 长度为64字节，字符在a-z0-9，可为空，若为空字符串，将自动生成txId
   - params: 执行参数
   - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
   - syncResultTimeout: 同步获取执行结果超时时间，小于等于0代表不等待执行结果，直接返回(返回信息里包含交易ID)，单位：毫秒
```java
ResultOuterClass.TxResponse invokeContract(
        String contractName, String method, String txId,
        Map<String, byte[]> params,
        long rpcCallTimeout, long syncResultTimeout)
```
### 合约查询
**参数说明**
   - contractName: 合约名
   - method: 方法名
   - txId: 交易id
      - 格式要求： 长度为64字节，字符在a-z0-9，可为空，若为空字符串，将自动生成txId
   - params: 执行参数
   - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
ResultOuterClass.TxResponse queryContract(
        String contractName, 
        String method, 
        String txId, 
        Map<String, byte[]> params, 
        long rpcCallTimeout)
```
### 合约调用的payload
**参数说明**
  - contractName: 合约名
  - method: 调用方法
  - txId: 交易id
     - 格式要求： 长度为64字节，字符在a-z0-9，可为空，若为空字符串，将自动生成txId
  - params: 请求参数
```java
Request.Payload invokeContractPayload(
        String contractName, 
        String method, 
        String txId,
        Map<String, byte[]> params)
```
### 合约查询(user)
**参数说明**
   - contractName: 合约名
   - method: 方法名
   - txId: 交易id
   - params: 执行参数
   - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
   - user: 用户
```java
ResultOuterClass.TxResponse queryContract(
        String contractName, 
        String method, 
        String txId,
        Map<String, byte[]> params, 
        long rpcCallTimeout, 
        User user)
```
## 系统合约接口
### 根据交易Id查询交易
**参数说明**
   - txId: 交易ID
   - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
 ChainmakerTransaction.TransactionInfo getTxByTxId(
         String txId, 
         long rpcCallTimeout)
```
### 根据交易Id查询包含rwset的交易
**参数说明**
   - txId: 交易ID
   - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
ChainmakerTransaction.TransactionInfoWithRWSet getTxWithRWSetByTxId(
        String txId, 
        long rpcCallTimeout)
```
### 根据区块高度查询区块
**参数说明**
   - blockHeight: 区块高度
   - withRWSet: 是否返回读写集
   - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
ChainmakerBlock.BlockInfo getBlockByHeight(
        long blockHeight, 
        boolean withRWSet, 
        long rpcCallTimeout)
```
### 根据区块哈希查询区块
**参数说明**
   - blockHash: 区块hash
   - withRWSet: 是否返回读写集
   - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
ChainmakerBlock.BlockInfo getBlockByHash(
        String blockHash, 
        boolean withRWSet, 
        long rpcCallTimeout)
```
### 根据交易Id查询区块
**参数说明**
  - txId: 交易Id
  - withRWSet: 是否返回读写集
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
ChainmakerBlock.BlockInfo getBlockByTxId(
        String txId, 
        boolean withRWSet, 
        long rpcCallTimeout)
```
### 查询最后一个配置块
**参数说明**
  - withRWSet: 是否返回读写集
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
ChainmakerBlock.BlockInfo getLastConfigBlock(
        boolean withRWSet, 
        long rpcCallTimeout)
```
### 查询节点加入的链信息
**参数说明**
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
Discovery.ChainList getNodeChainList(
        long rpcCallTimeout)
```
### 查询链信息
**参数说明**
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
Discovery.ChainInfo getChainInfo(
        long rpcCallTimeout)
```
### 根据txId查询区块高度
**参数说明**
  - txId: 交易id
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
long getBlockHeightByTxId(
        String txId, 
        long rpcCallTimeout)
```
### 根据blockHash查询区块高度
**参数说明**
  - blockHash: 区块哈希
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
long getBlockHeightByBlockHash(
        String blockHash, 
        long timeout) 
```
### 根据区块高度查询完整区块
**参数说明**
  - blockHeight: 区块高度
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
Store.BlockWithRWSet getFullBlockByHeight(
        long blockHeight, 
        long rpcCallTimeout)
```
### 查询最新区块信息
**参数说明**
  - withRWSet: 是否返回读写集
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
ChainmakerBlock.BlockInfo getLastBlock(
        boolean withRWSet, 
        long rpcCallTimeout)
```
### 查询最新区块高度
**参数说明**
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
long getCurrentBlockHeight(
        long rpcCallTimeout)
```
### 根据区块高度查询区块头
**参数说明**
  - blockHeight: 区块高度
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
ChainmakerBlock.BlockHeader getBlockHeaderByHeight(
        long blockHeight, 
        long rpcCallTimeout)
```
### 系统合约调用
**参数说明**
  - contractName: 合约名
  - method: 方法名
  - txId: 交易id
  - params: 执行参数
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
  - syncResultTimeout: 同步获取执行结果超时时间，小于等于0代表不等待执行结果，直接返回（返回信息里包含交易ID），单位：毫秒
```java
ResultOuterClass.TxResponse invokeSystemContract(
        String contractName, 
        String method, 
        String txId, 
        Map<String, byte[]> params,
        long rpcCallTimeout, 
        long syncResultTimeout)
```
### 合约查询接口
**参数说明**
  - contractName: 合约名
  - method: 方法名
  - txId: 交易id
  - params: 执行参数
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
ResultOuterClass.TxResponse querySystemContract(
        String contractName, 
        String method, 
        String txId,
        Map<String, byte[]> params, 
        long rpcCallTimeout)
```
### 根据交易Id获取Merkle路径
**参数说明**
  - txId: 交易ID
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
byte[] getMerklePathByTxId(
        String txId, 
        long rpcCallTimeout)
```
### 开放系统合约
**参数说明**
  - grantContractList: 需要开放的系统合约字符串数组
```java
Request.Payload createNativeContractAccessGrantPayload(
        String[] grantContractList)
```
### 弃用系统合约
**参数说明**
  - revokeContractList: 需要弃用的系统合约字符串数组
```java
Request.Payload createNativeContractAccessRevokePayload(
        String[] revokeContractList)
```
### 查询弃用的系统合约名单
```java
Request.Payload createGetDisabledNativeContractListPayload()
```
### 查询指定合约的信息，包括系统合约和用户合约
**参数说明**
  - contractName: 指定查询的合约名字，包括系统合约和用户合约
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
ContractOuterClass.Contract getContractInfo(
        String contractName, 
        long rpcCallTimeout)
```
### 查询所有的合约名单，包括系统合约和用户合约
**参数说明**
  - rpcCallTimeout
```java
ContractOuterClass.Contract[] getContractList(
        long rpcCallTimeout) 
```
### 查询已禁用的系统合约名单, 无禁用合约时返回null
**参数说明**
  - rpcCallTimeout
```java
ContractOuterClass.Contract[] getDisabledNativeContractList(
        long rpcCallTimeout)
```
## 链配置接口
### 查询最新链配置
**参数说明**
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
ChainConfigOuterClass.ChainConfig getChainConfig(
        long rpcCallTimeout)
```
### 根据指定区块高度查询最近链配置
**参数说明**
  - blockHeight: 区块高度
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
ChainConfigOuterClass.ChainConfig getChainConfigByBlockHeight(
        long blockHeight, 
        long rpcCallTimeout)
```
### 查询最新链配置序号Sequence
**参数说明**
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
long getChainConfigSequence(
        long rpcCallTimeout)
```
### 生成更新Core模块待签名payload
**参数说明**
  - txSchedulerTimeout: 交易调度器从交易池拿到交易后, 进行调度的时间，其值范围为[0, 60]，若无需修改，请置为-1
  - txSchedulerValidateTimeout: 交易调度器从区块中拿到交易后, 进行验证的超时时间，其值范围为[0, 60]，若无需修改，请置为-1
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
Request.Payload createPayloadOfChainConfigCoreUpdate(
        int txSchedulerTimeout, 
        int txSchedulerValidateTimeout, 
        long rpcCallTimeout)
```
### 生成更新Block模块待签名payload
**参数说明**
  - txTimestampVerify: 是否需要开启交易时间戳校验
  - (以下参数，若无需修改，请置为-1)
  - txTimeout: 交易时间戳的过期时间(秒)，其值范围为[600, +∞)
  - blockTxCapacity: 区块中最大交易数，其值范围为(0, +∞]
  - blockSize: 区块最大限制，单位MB，其值范围为(0, +∞]
  - blockInterval: 出块间隔，单位:ms，其值范围为[10, +∞]
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
Request.Payload createPayloadOfChainConfigBlockUpdate(
        boolean txTimestampVerify, 
        int txTimeout, 
        int blockTxCapacity,
        int blockSize, 
        int blockInterval, 
        int txParameterSize, 
        long rpcCallTimeout)
```
### 生成添加信任组织根证书待签名payload
**参数说明**
  - trustRootOrgId: 组织Id
  - trustRootCrt: 根证书
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
Request.Payload createPayloadOfChainConfigTrustRootAdd(
        String trustRootOrgId, 
        String[] trustRootCrt, 
        long rpcCallTimeout)
```
### 生成更新信任组织根证书待签名payload
**参数说明**
  - trustRootOrgId: 组织Id
  - trustRootCrt: 根证书
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
Request.Payload createPayloadOfChainConfigTrustRootUpdate(
        String trustRootOrgId, 
        String[] trustRootCrt, 
        long rpcCallTimeout)
```
### 生成删除信任组织根证书待签名payload
**参数说明**
  - trustRootOrgId: 组织Id
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
Request.Payload createPayloadOfChainConfigTrustRootDelete(
        String orgIdOrPKPubkeyPEM, 
        long rpcCallTimeout)
```
### 添加权限配置待签名payload生成
**参数说明**
  - permissionResourceName: 权限名
  - principle: 权限规则
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
Request.Payload createPayloadOfChainConfigPermissionAdd(
        String permissionResourceName, 
        PolicyOuterClass.Policy principal, 
        long rpcCallTimeout)
```
### 更新权限配置待签名payload生成
**参数说明**
  - permissionResourceName: 权限名
  - principle: 权限规则
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
Request.Payload createPayloadOfChainConfigPermissionUpdate(
        String permissionResourceName, 
        PolicyOuterClass.Policy principal, 
        long rpcCallTimeout)
```
### 删除权限配置待签名payload生成
**参数说明**
  - permissionResourceName: 权限名
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
Request.Payload createPayloadOfChainConfigPermissionDelete(
        String permissionResourceName, 
        long rpcCallTimeout)
```
### 添加共识节点地址待签名payload生成
**参数说明**
  - nodeOrgId: 节点组织Id
  - nodeAddresses: 节点地址
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
Request.Payload createPayloadOfChainConfigConsensusNodeAddrAdd(
        String nodeOrgId, String[] nodeAddresses, 
        long rpcCallTimeout)
```
### 更新共识节点地址待签名payload生成
**参数说明**
  - nodeOrgId: 节点组织Id
  - nodeOldAddress: 节点原地址
  - nodeNewAddress: 节点新地址
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
Request.Payload createPayloadOfChainConfigConsensusNodeAddrUpdate(
        String nodeOrgId, 
        String nodeOldAddress, 
        String nodeNewAddress, 
        long rpcCallTimeout)
```
### 删除共识节点地址待签名payload生成
**参数说明**
  - nodeOrgId: 节点组织Id
  - nodeAddress: 节点地址
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
Request.Payload createPayloadOfChainConfigConsensusNodeAddrDelete(
        String nodeOrgId, 
        String nodeAddress, 
        long rpcCallTimeout)
```
### 添加共识节点待签名payload生成
**参数说明**
  - nodeOrgId: 节点组织Id
  - nodeAddresses: 节点地址
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
Request.Payload createPayloadOfChainConfigConsensusNodeOrgAdd(
        String nodeOrgId, 
        String[] nodeAddresses, 
        long rpcCallTimeout)
```
### 更新共识节点待签名payload生成
**参数说明**
  - nodeOrgId: 节点组织Id
  - nodeAddresses: 节点地址
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
Request.Payload createPayloadOfChainConfigConsensusNodeOrgUpdate(
        String nodeOrgId, 
        String[] nodeAddresses, 
        long rpcCallTimeout)
```
### 删除共识节点待签名payload生成
**参数说明**
  - nodeOrgId: 节点组织Id
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
Request.Payload createPayloadOfChainConfigConsensusNodeOrgDelete(
        String nodeOrgId, 
        long rpcCallTimeout)
```
### 添加共识扩展字段待签名payload生成
**参数说明**
  - params: Map<String, byte[]>
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
Request.Payload createPayloadOfChainConfigConsensusExtAdd(
        Map<String, byte[]> params, 
        long rpcCallTimeout)
```
### 添加共识扩展字段待签名payload生成
**参数说明**
  - params: Map<String, byte[]>
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
Request.Payload createPayloadOfChainConfigConsensusExtUpdate(
        Map<String, byte[]> params, 
        long rpcCallTimeout)
```
### 添加共识扩展字段待签名payload生成
**参数说明**
  - keys: 待删除字段
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
Request.Payload createPayloadOfChainConfigConsensusExtDelete(
        String[] keys, 
        long rpcCallTimeout)
```
### 添加信任成员证书待签名payload生成
**参数说明**
  - trustMemberOrgId: 组织Id
  - trustMemberNodeId: 节点Id
  - trustMemberRole: 成员角色
  - trustMemberInfo: 成员信息内容
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
Request.Payload createChainConfigTrustMemberAddPayload(
        String trustMemberOrgId, 
        String trustMemberNodeId, 
        String trustMemberRole, 
        String trustMemberInfo, 
        long rpcCallTimeout)
```
### 删除信任成员证书待签名payload生成
**参数说明**
  - trustMemberInfo: 成员信息内容
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
Request.Payload createChainConfigTrustMemberDeletePayload(
        String trustMemberInfo, 
        long rpcCallTimeout)
```
### 发送链配置更新请求
**参数说明**
  - payload: 待签名payload
  - endorsementEntries: 背书实体数组
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
  - syncResultTimeout: 同步结果超时时间, 单位：毫秒
```java
ResultOuterClass.TxResponse updateChainConfig(
        Request.Payload payload, 
        Request.EndorsementEntry[] endorsementEntries, 
        long rpcCallTimeout, 
        long syncResultTimeout)
```
### 修改地址类型payload生成
**参数说明**
  - addrType: 地址类型，0-ChainMaker; 1-ZXL
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
Request.Payload createChainConfigAlterAddrTypePayload(
        String addressType, 
        long rpcCallTimeout)
```
### 开启或关闭链配置的Gas优化payload生成
**参数说明**
  - enable: 是否开启
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
Request.Payload createChainConfigOptimizeChargeGasPayload(
        Boolean enable, 
        long rpcCallTimeout)
```
### 查询最新权限配置列表
**参数说明**
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
List<ChainConfigOuterClass.ResourcePolicy> getChainConfigPermissionList(
        long rpcCallTimeout)
```
## 证书管理接口
### 用户证书添加
**参数说明**rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
ResultOuterClass.TxResponse addCert(
        long rpcCallTimeout)
```
### 用户证书删除
**参数说明**
  - payload: 合约内容
  - endorsementEntries: 带签名的合约内容
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
  - syncResultTimeout: 同步结果超时时间, 单位：毫秒
```java
ResultOuterClass.TxResponse deleteCert(
        Request.Payload payload, 
        Request.EndorsementEntry[] endorsementEntries, 
        long rpcCallTimeout, long syncResultTimeout)
```
### 用户证书查询
**参数说明**
  - certHashes: 证书Hash列表
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
ResultOuterClass.CertInfos queryCert(
        String[] certHashes, 
        long rpcCallTimeout)
```
### 证书冻结
**参数说明**
  - payload: 证书冻结的payload
  - endorsementEntries: 带签名的合约内容
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
  - syncResultTimeout: 同步结果超时时间, 单位：毫秒
```java
ResultOuterClass.TxResponse freezeCerts(
        Request.Payload payload, 
        Request.EndorsementEntry[] endorsementEntries, 
        long rpcCallTimeout, 
        long syncResultTimeout)
```
### 证书解冻
**参数说明**
  - payload: 证书解冻的payload
  - endorsementEntries: 带签名的合约内容
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
  - syncResultTimeout: 同步结果超时时间, 单位：毫秒
```java
ResultOuterClass.TxResponse unfreezeCerts(
        Request.Payload payload, 
        Request.EndorsementEntry[] endorsementEntries, 
        long rpcCallTimeout, 
        long syncResultTimeout)
```
### 证书吊销
**参数说明**
  - payload: 证书注销的payload
  - endorsementEntries: 带签名的合约内容
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
  - syncResultTimeout: 同步结果超时时间, 单位：毫秒
```java
ResultOuterClass.TxResponse revokeCerts(
        Request.Payload payload, 
        Request.EndorsementEntry[] endorsementEntries, 
        long rpcCallTimeout, 
        long syncResultTimeout)
```
### 证书操作payload生成
**参数说明**
  - method: 证书操作方法
  - params: 参数
```java
Request.Payload createCertManagePayload(
        String method, 
        Map<String, byte[]> params)
```
### 证书删除payload生成
**参数说明**
  - certHashes: 证书Hash列表
```java
Request.Payload createCertDeletePayload(
        String[] certHashes)
```
### 证书冻结payload生成
**参数说明**
  - certs: 证书内容列表
```java
Request.Payload createCertFreezePayload(
        String[] certs)
```
### 证书解冻payload生成
**参数说明**
  - certs: 证书内容列表
```java
Request.Payload createPayloadOfUnfreezeCerts(
        String[] certs)
```
### 发送证书管理请求（证书冻结、解冻、吊销）
**参数说明**
  - payload: 交易payload
  - endorsers: 背书签名信息列表
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
  - syncResultTimeout: 同步结果超时时间, 单位：毫秒
```java
ResultOuterClass.TxResponse sendCertManageRequest(
        Request.Payload payload, 
        Request.EndorsementEntry[] endorsementEntries, 
        long rpcCallTimeout, 
        long syncResultTimeout)
```
### 创建用户证书吊销payload
**参数说明**
  - certCrl: 吊销的证书列表
```java
Request.Payload createPayloadOfRevokeCerts(
        String certCrl)
```
## 消息订阅接口
### 区块订阅
**参数说明**
  - startBlock: 订阅起始区块高度，表示订阅实时最新区块
  - endBlock: 订阅结束区块高度，若为-1，表示订阅实时最新区块
  - withRwSet: 是否返回读写集
  - onlyHeader: 是否只返回区块头
  - blockStreamObserver: 区块流观察者
```java
void subscribeBlock(
        long startBlock, 
        long endBlock, 
        boolean withRwSet, 
        boolean onlyHeader, 
        StreamObserver<ResultOuterClass.SubscribeResult> blockStreamObserver)
```
### 交易订阅
**参数说明**
  - startBlock: 订阅起始区块高度，表示订阅实时最新区块
  - endBlock: 订阅结束区块高度，若为-1，表示订阅实时最新区块
  - contractName: 订阅合约名
  - txIds: 订阅txId列表，若为空，表示订阅所有txId
  - txStreamObserver: 交易流观察者
```java
void subscribeTx(
        long startBlock, 
        long endBlock, 
        String contractName, 
        String[] txIds, 
        StreamObserver<ResultOuterClass.SubscribeResult> txStreamObserver)
```
### 通过别名前缀进行交易订阅
**参数说明**
  - startBlock: 订阅起始区块高度，表示订阅实时最新区块
  - endBlock: 订阅结束区块高度，若为-1，表示订阅实时最新区块
  - preAlias: 别名前缀
  - txStreamObserver: 交易流观察者
```java
void subscribeTxByPreAlias(
        long startBlock, 
        long endBlock, 
        String preAlias, 
        StreamObserver<ResultOuterClass.SubscribeResult> txStreamObserver)
```
### 通过交易ID前缀进行交易订阅
**参数说明**
  - startBlock: 订阅起始区块高度，表示订阅实时最新区块
  - endBlock: 订阅结束区块高度，若为-1，表示订阅实时最新区块
  - preTxId: 交易ID前缀
  - txStreamObserver: 交易流观察者
```java
void subscribeTxByPreTxId(
        long startBlock, 
        long endBlock, 
        String preTxId, 
        StreamObserver<ResultOuterClass.SubscribeResult> txStreamObserver)
```
### 通过组织ID前缀进行交易订阅
**参数说明**
  - startBlock: 订阅起始区块高度，表示订阅实时最新区块
  - endBlock: 订阅结束区块高度，若为-1，表示订阅实时最新区块
  - preOrgId: 组织ID前缀
  - txStreamObserver: 交易流观察者
```java
void subscribeTxByPreOrgId(
        long startBlock, 
        long endBlock, 
        String preOrgId, 
        StreamObserver<ResultOuterClass.SubscribeResult> txStreamObserver)
```
### 事件订阅
**参数说明**
  - startBlock: 订阅起始区块高度，表示订阅实时最新区块
  - endBlock: 订阅结束区块高度，若为-1，表示订阅实时最新区块
  - topic: 订阅话题
  - contractName: 订阅合约名
  - txStreamObserver: 交易流观察者
```java
void subscribeContractEvent(
        long startBlock, 
        long endBlock, 
        String topic, 
        String contractName,
        StreamObserver<ResultOuterClass.SubscribeResult> contractEventStreamObserver)
```
## 数据归档接口
### 发送数据归档请求
**参数说明**
  - payload: 数据归档payload
  - timeout: 超时时间，单位：毫秒
```java
ResultOuterClass.TxResponse sendArchiveBlockRequest(
        Request.Payload payload, 
        long timeout)
```
### 发送归档恢复请求
**参数说明**
  - payload: 归档恢复payload
  - timeout: 超时时间，单位：毫秒
```java
ResultOuterClass.TxResponse sendRestoreBlockRequest(
        Request.Payload payload, 
        long timeout)
```
### 数据归档payload生成
**参数说明**
  - targetBlockHeight: 归档区块高度
```java
Request.Payload createArchiveBlockPayload(
        long targetBlockHeight)
```
### 归档恢复payload生成
**参数说明**
  - fullBlock: 归档恢复数据
```java
Request.Payload createRestoreBlockPayload(
        byte[] fullBlock)
```
### 获取归档数据
**参数说明**
  - blockHeight: 归档区块高度
```java
Store.BlockWithRWSet getArchivedFullBlockByHeight(
        long blockHeight)
```
### 获取归档区块信息
**参数说明**
  - blockHeight: 归档区块高度
  - withRWSet: 是否获取读写集
```java
ChainmakerBlock.BlockInfo getArchivedBlockByHeight(
        long blockHeight, 
        boolean withRWSet)
```
### 获取已归档区块高度
**参数说明**
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
long getArchivedBlockHeight(
        long rpcCallTimeout)
```
## 公钥身份类接口
### 构造添加公钥身份请求payload
**参数说明**
  - pubkey: 公钥信息
  - orgId: 组织id
  - role: 角色，支持client,light,common
```java
Request.Payload createPubkeyAddPayload(
        String pubkey, 
        String orgId, 
        String role)
```
### 构造删除公钥身份请求payload
**参数说明**
  - pubkey: 公钥信息
  - orgId: 组织id
```java
Request.Payload createPubkeyDelPayload(
        String pubkey, 
        String orgId)
```
### 构造查询公钥身份请求payload
**参数说明**
  - pubkey: 公钥信息
```java
Request.Payload createPubkeyQueryPayload(
        String pubkey)
```
### 发送公钥身份管理请求（添加、删除）
**参数说明**
  - payload: 合约内容
  - endorsementEntries: 带签名的合约内容
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
  - syncResultTimeout: 同步获取执行结果超时时间，小于等于0代表不等待执行结果，直接返回（返回信息里包含交易ID），单位：毫秒
```java
ResultOuterClass.TxResponse sendPubkeyManageRequest(
        Request.Payload payload, 
        Request.EndorsementEntry[] endorsementEntries, 
        long rpcCallTimeout, 
        long syncResultTimeout)
```
## 多签类接口
### 发起多签请求
**参数说明**
  - payload: 多签payload
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
  - withSyncResult: 是否同步获取交易执行结果
```java
ResultOuterClass.TxResponse multiSignContractReq(
        Request.Payload payload, 
        long rpcCallTimeout, 
        boolean withSyncResult)
```
### 发起带有代付者的多签请求
**参数说明**
  - payload: 多签payload
  - endorsers: 签名者列表
  - payer: 代付者
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
  - withSyncResult: 是否同步获取交易执行结果
```java
ResultOuterClass.TxResponse multiSignContractReqWithPayer(
        Request.Payload payload, 
        Request.EndorsementEntry[] endorsers, 
        Request.EndorsementEntry payer, 
        long rpcCallTimeout, 
        boolean withSyncResult)
```
### 发起多签投票
**参数说明**
  - payload: 多签payload
  - endorsementEntry: 多签信息
  - isAgree: 投票人对多签请求是否同意，true为同意，false则反对
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
  - withSyncResult: 是否同步获取交易执行结果
```java
ResultOuterClass.TxResponse multiSignContractVote(
        Request.Payload payload, 
        Request.EndorsementEntry endorsementEntry, 
        boolean isAgree, long rpcCallTimeout, 
        boolean withSyncResult)
```
### 发起带有gas限制的多签投票
**参数说明**
  - payload: 多签payload
  - endorsementEntry: 多签信息
  - isAgree: 投票人对多签请求是否同意，true为同意，false则反对
  - gasLimit: gas限制
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
  - withSyncResult: 是否同步获取交易执行结果
```java
ResultOuterClass.TxResponse multiSignContractVoteWithGasLimit(
        Request.Payload payload, 
        Request.EndorsementEntry endorsementEntry, 
        boolean isAgree, 
        long gasLimit, 
        long rpcCallTimeout, 
        boolean withSyncResult)
```
### 发起带有gas限制和代付者的多签投票
**参数说明**
  - payload: 多签payload
  - endorsementEntry: 多签信息
  - payer: 代付者
  - isAgree: 投票人对多签请求是否同意，true为同意，false则反对
  - gasLimit: gas限制
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
  - withSyncResult: 是否同步获取交易执行结果
```java
ResultOuterClass.TxResponse multiSignContractVoteWithGasLimitAndPayer(
        Request.Payload payload, 
        Request.EndorsementEntry endorsementEntry, 
        Request.EndorsementEntry payer, 
        boolean isAgree, 
        long gasLimit, 
        long rpcCallTimeout, 
        boolean withSyncResult)
```
### 触发执行多签请求
**参数说明**
  - payload: 多签payload
  - limit: gas限制
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
  - withSyncResult: 是否同步获取交易执行结果
```java
ResultOuterClass.TxResponse multiSignContractTrig(
        Request.Payload payload, 
        Request.Limit limit, 
        long rpcCallTimeout, 
        boolean withSyncResult)
```
### 触发执行带有代付者的多签请求
**参数说明**
  - payload: 多签payload
  - payer: 代付者
  - limit: gas限制
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
  - withSyncResult: 是否同步获取交易执行结果
```java
ResultOuterClass.TxResponse multiSignContractTrigWithPayer(
        Request.Payload payload, 
        Request.EndorsementEntry payer, 
        Request.Limit limit, 
        long rpcCallTimeout, 
        boolean withSyncResult)
```
### 触发执行多签请求
**参数说明**
  - txId: 交易id
  - limit: gas限制
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
  - withSyncResult: 是否同步获取交易执行结果
```java
ResultOuterClass.TxResponse multiSignContractTrig(
        String txId, 
        Request.Limit limit, 
        long rpcCallTimeout, 
        boolean withSyncResult)
```
### 触发执行带有代付者的多签请求
**参数说明**
  - txId: 交易id
  - payer: 代付者
  - limit: gas限制
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
  - withSyncResult: 是否同步获取交易执行结果
```java
ResultOuterClass.TxResponse multiSignContractTrigWithPayer(
        String txId, 
        Request.EndorsementEntry payer, 
        Request.Limit limit, 
        long rpcCallTimeout, 
        boolean withSyncResult)
```
### 多签查询
**参数说明**
  - txId: 交易id
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
ResultOuterClass.TxResponse multiSignContractQuery(
        String txId, 
        long rpcCallTimeout)
```
### 根据参数进行多签查询
**参数说明**
  - txId: 交易id
  - param: 多签参数
  - rpcCallTimeout: 调用rcp接口超时时间, 单位：毫秒
```java
ResultOuterClass.TxResponse multiSignContractQueryWithParams(
        String txId, Map<String, byte[]> params, 
        long rpcCallTimeout)
```
### 创建多签请求payload
**参数说明**
  - params: 多签参数
```java
Request.Payload createMultiSignReqPayload(
        Map<String, byte[]> params)
```
### 多签请求待签名payload生成
**参数说明**
  - params: 多签参数
  - gasLimit: gas限制
```java
Request.Payload createMultiSignReqPayloadWithGasLimit(
        Map<String, byte[]> params, 
        long gasLimit)
```
### 创建多签投票payload
**参数说明**
  - params: 多签参数
  - limit: gas限制
```java
Request.Payload createMultiSignVotePayload(
        Map<String, byte[]> params, 
        Request.Limit limit)
```
### 创建多签触发payload
**参数说明**
  - params: 多签参数
  - limit: gas限制
```java
Request.Payload createMultiSignTrigPayload(
        Map<String, byte[]> params, 
        Request.Limit limit)
```
### 创建多签查询payload
**参数说明**
  - params: 多签参数
```java
Request.Payload createMultiSignQueryPayload(
        Map<String, byte[]> params)
```
## 管理类接口
### SDK停止接口：关闭连接池连接，释放资源
```java
void stop()
```
### 获取链版本
**参数说明**
    - timeout: 调用rcp接口超时时间, 单位：毫秒
```java
String getChainMakerServerVersion(
        long timeout)
```
### 更新链配置
**参数说明**
    - nodeConfig: 节点配置信息
    - timeout: 调用rcp接口超时时间, 单位：毫秒
```java
LocalConfig.CheckNewBlockChainConfigResponse checkNewBlockChainConfig(
        NodeConfig nodeConfig, 
        long timeout)
```
## gas管理相关接口
### 构造设置gas管理员payload
**参数说明**
  - address: gas管理员的地址
```java
Request.Payload createSetGasAdminPayload(
        String address)
```
### 查询gas管理员
**参数说明**
  - rpcCallTimeout: 调用rpc接口超时时间, 单位：毫秒
```java
String getGasAdmin(
        long rpcCallTimeout)
```
### 构造充值gas账户payload
**参数说明**
  - rechargeGasList: 一个gas账户充值指定gas数量
```java
Request.Payload createRechargeGasPayload(
        AccountManager.RechargeGas[] rechargeGasList)
```
### 查询gas账户余额（根据公钥）
**参数说明**
  - address: 查询gas余额的账户地址
  - rpcCallTimeout: 调用rpc接口超时时间, 单位：毫秒
```java
long getGasBalance(
        String address, 
        long rpcCallTimeout)
```
### 构造退还gas的payload
**参数说明**
  - address: 退还gas的账户地址
  - amount: 退还gas的数量
```java
Request.Payload createRefundGasPayload(
        String address, 
        long amount)
```
### 构造冻结gas账户的payload
**参数说明**
  - address: 冻结指定gas账户的账户地址
```java
Request.Payload createFrozenGasAccountPayload(
        String address)
```
### 构造解冻指定gas账户的payload
**参数说明**
  - address: 解冻指定gas账户的账户地址
```java
Request.Payload createUnfrozenGasAccountPayload(
        String address)
```
### 查询gas账户的状态
**参数说明**
  - address: 指定gas账户的账户地址
  - rpcCallTimeout: 调用rpc接口超时时间, 单位：毫秒
  **返回值说明**
  - boolean: true表示账号未被冻结，false表示账号已被冻结
```java
boolean getGasAccountStatus(
        String address, 
        long rpcCallTimeout)
```
### 发送gas管理类请求
**参数说明**
  - payload: 交易payload
  - endorsementEntries: 背书签名信息列表
  - rpcCallTimeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
  - syncResultTimeout: 是否同步获取交易执行结果
```java
ResultOuterClass.TxResponse sendGasManageRequest(
        Request.Payload payload, 
        Request.EndorsementEntry[] endorsementEntries, 
        long rpcCallTimeout, 
        long syncResultTimeout)
```
### 为payload添加gas limit
**参数说明**
  - payload: 交易payload
  - limit: gas limit
```java
Request.Payload attachGasLimit(
        Request.Payload payload, 
        Request.Limit limit)
```
### 启用或停用Gas计费开关payload生成
**参数说明**
  - rpcCallTimeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
```java
Request.Payload createChainConfigEnableOrDisableGasPayload(
        long rpcCallTimeout)
```
### 构造配置账户基础gas消耗数量payload
**参数说明**
  - amount: 基础gas消耗数量
  - rpcCallTimeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
```java
Request.Payload createSetInvokeBaseGasPayload(
        long amount, 
        long rpcCallTimeout)
```
### 构造设置调用gas price的payload
**参数说明**
  - gasPrice: gas价格
  - rpcCallTimeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
```java
Request.Payload createSetInvokeGasPricePayload(
        String gasPrice, 
        long rpcCallTimeout)
```
### 构造配置账户基础gas消耗数量的payload
**参数说明**
  - amount: 基础gas消耗数量
  - rpcCallTimeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
```java
Request.Payload createSetInstallBaseGasPayload(
        long amount, 
        long rpcCallTimeout)
```
### 构造配置调用gas price的payload
**参数说明**
  - gasPrice: gas价格
  - rpcCallTimeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
```java
Request.Payload createSetInstallGasPricePayload(
        String gasPrice, 
        long rpcCallTimeout)
```
### 估算交易的gas消耗量
**参数说明**
  - payload: 待估算gas消耗量的交易payload
```java
long estimateGas(
        Request.Payload payload, 
        long rpcCallTimeout)
```
## 别名相关接口
### 添加别名
**参数说明**
  - rpcCallTimeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
```java
ResultOuterClass.TxResponse addAlias(
        long rpcCallTimeout)
```
### 构造更新别名payload
**参数说明**
  - alias: 要更新的别名
  - certPEM: 对应的证书
```java
Request.Payload createAliasUpdatePayload(
        String alias, 
        String certPem)
```
### 发起更新别名交易
**参数说明**
  - payload: 待签名的payload
  - endorsementEntries: 背书签名信息列表
  - rpcCallTimeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
```java
ResultOuterClass.TxResponse updateAlias(
        Request.Payload payload, 
        Request.EndorsementEntry[] endorsementEntries, 
        long rpcCallTimeout)
```
### 查询别名详情交易
**参数说明**
  - aliasList: 要查询的别名列表
  - rpcCallTimeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
```java
ResultOuterClass.AliasInfos queryAlias(
        String[] aliasList, 
        long rpcCallTimeout)
```
### 生成删除别名payload
**参数说明**
  - aliasList: 要删除的别名列表
  - rpcCallTimeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
```java
Request.Payload createAliasDeletePayload(
        String[] aliasList)
```
### 发起删除别名交易
**参数说明**
  - payload: 待签名的payload
  - endorsementEntries: 背书签名信息列表
  - rpcCallTimeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
  - syncResultTimeout: 是否同步获取交易执行结果
```java
ResultOuterClass.TxResponse deleteAlias(
        Request.Payload payload, 
        Request.EndorsementEntry[] endorsementEntries, 
        long rpcCallTimeout, 
        long syncResultTimeout)
```
## 交易池相关接口
### 获取交易池状态
**参数说明**
  - rpcCallTimeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s
```java
TransactionPool.TxPoolStatus getPoolStatus(
        long rpcCallTimeout)
```
### 获取不同交易类型和阶段中的交易Id列表。
**参数说明**
  - txType: 交易类型，在pb的txpool包中进行了定义
  - txStage: 交易阶段，在pb的txpool包中进行了定义
  - rpcCallTimeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s

**返回值说明**
  - []string: 交易Id列表
```java
List<String> getTxIdsByTypeAndStage(
        TransactionPool.TxType txType, 
        TransactionPool.TxStage txStage, 
        long rpcCallTimeout)
```
### 根据txIds获取交易池中存在的txs，并返回交易池缺失的tx的txIds
**参数说明**
  - txIds: 交易Id列表
  - rpcCallTimeout: 超时时间，单位：s，若传入-1，将使用默认超时时间：10s

**返回值说明**
  - []*common.Transaction: 交易池中存在的txs
  - []string: 交易池缺失的tx的txIds
```java
List<ChainmakerTransaction.Transaction> getTxsInPoolByTxIds(
        String[] txIds, 
        long rpcCallTimeout)
```
## 交易相关接口
### 创建交易请求
**参数说明**
  - payload: 交易payload
  - endorsementEntries: 背书签名信息列表
```java
Request.TxRequest createTxRequest(Request.Payload payload, Request.EndorsementEntry[] endorsementEntries)
```
## 归档相关接口
### 根据交易ID获取归档区块
**参数说明**
  - txId: 交易ID
  - withRWSet: 是否包含读写集
  - timeout: 超时时间
```java
ChainmakerBlock.BlockInfo getArchivedBlockByTxId(
        String txId, 
        boolean withRWSet, 
        long timeout)
```
### 根据区块哈希获取归档区块
**参数说明**
  - blockHash: 区块哈希
  - withRWSet: 是否包含读写集
  - timeout: 超时时间
```java
ChainmakerBlock.BlockInfo getArchivedBlockByHash(
        String blockHash, 
        boolean withRWSet, 
        long timeout)
```
### 根据交易ID获取归档交易
**参数说明**
  - txId: 交易ID
  - timeout: 超时时间
```java
ChainmakerTransaction.TransactionInfo getArchivedTxByTxId(
        String txId, 
        long timeout)
```

### 根据区块高度进行归档
**参数说明**
  - archiveHeight: 区块高度（归档的高度为已经归档的高度到当前区块高度）
  - notice: 通知方法
  - rpcCallTimeout: 超时时间
```java
void archiveBlocks(long archiveHeight, Notice notice, long rpcCallTimeout)
```

### 根据区块高度进行恢复
**参数说明**
  - restoreHeight: 区块高度（恢复高度为当前高度到已经归档的高度）
  - notice: 通知方法
  - rpcCallTimeout: 超时时间
```java
void restoreBlocks(long restoreHeight, Notice notice, long rpcCallTimeout)
```