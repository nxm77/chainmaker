# Go SDK 使用说明

<span id="section_sdk"></span>

本篇介绍：

1、环境依赖

2、sdk依赖使用

3、普通合约安装、调用

4、EVM合约安装、调用

5、更多的示例及全部接口

## 长安链SDK概述

1. 整体介绍

长安链`SDK`是业务模块与长安链交互的桥梁，支持双向`TLS`认证，提供安全可靠的加密通信信道。

长安链提供了多种语言的`SDK`，包括：`Go SDK`、`Java SDK`、`Python SDK`、`Nodejs SDK`方便开发者根据需要进行选用。

提供的`SDK`接口，覆盖合约管理、链配置管理、证书管理、多签收集、各类查询操作、事件订阅等场景，满足了不同的业务场景需要。

2. 名词概念说明

- **`Node`（节点）**：代表一个链节点的基本信息，包括：节点地址、连接数、是否启用`TLS`认证等信息
- **`ChainClient`（链客户端）**：所有客户端对链节点的操作接口都来自`ChainClient`
- **压缩证书**：可以为`ChainClient`开启证书压缩功能，开启后可以减小交易包大小，提升处理性能

## 环境准备

### 软件环境依赖

**golang** ： 版本为1.16或以上

下载地址：<https://golang.org/dl/>

若已安装，请通过命令查看版本：

```bash
$ go version
go version go1.16 linux/amd64
```

### 下载安装sdk
进入您的Go项目，执行以下命令添加对sdk的引用：
```bash
go get chainmaker.org/chainmaker/sdk-go/v2@v2.3.6
```

### 长安链环境准备

创建一条证书模式的长安链，并确保相关节点网络通畅，相关教程见：[《通过命令行体验链》](../quickstart/通过命令行体验链.md)

## 怎么使用SDK

### 示例代码
#### 创建节点

设置节点信息，可用作创建与该节点连接的客户端

```go
// 创建节点
func createNode(nodeAddr string, connCnt int) *NodeConfig {
 node := NewNodeConfig(
  // 节点地址，格式：127.0.0.1:12301
  WithNodeAddr(nodeAddr),
  // 节点连接数
  WithNodeConnCnt(connCnt),
  // 节点是否启用TLS认证
  WithNodeUseTLS(true),
  // 根证书路径，支持多个
  WithNodeCAPaths(caPaths),
  // TLS Hostname
  WithNodeTLSHostName(tlsHostName),
 )

 return node
}
```

#### 以参数形式创建ChainClient

> 更多内容请参看：`sdk_client_test.go`
>
> 注：示例中证书采用路径方式去设置，也可以使用证书内容去设置，具体请参看`createClientWithCaCerts`方法

```go
// 创建ChainClient
func createClient() (*ChainClient, error) {
 if node1 == nil {
  // 创建节点1
  node1 = createNode(nodeAddr1, connCnt1)
 }

 if node2 == nil {
  // 创建节点2
  node2 = createNode(nodeAddr2, connCnt2)
 }

 chainClient, err := NewChainClient(
  // 设置归属组织
  WithChainClientOrgId(chainOrgId),
  // 设置链ID
  WithChainClientChainId(chainId),
  // 设置logger句柄，若不设置，将采用默认日志文件输出日志
  WithChainClientLogger(getDefaultLogger()),
  // 设置客户端用户私钥路径
  WithUserKeyFilePath(userKeyPath),
  // 设置客户端用户证书
  WithUserCrtFilePath(userCrtPath),
  // 添加节点1
  AddChainClientNodeConfig(node1),
  // 添加节点2
  AddChainClientNodeConfig(node2),
  )

 if err != nil {
  return nil, err
 }

 //启用证书压缩（开启证书压缩可以减小交易包大小，提升处理性能）
 err = chainClient.EnableCertHash()
 if err != nil {
  log.Fatal(err)
 }

 return chainClient, nil
}
```

#### 以配置文件形式创建ChainClient

> 注：参数形式和配置文件形式两个可以同时使用，同时配置时，以参数传入为准

```go
func createClientWithConfig() (*ChainClient, error) {

 chainClient, err := NewChainClient(
  WithConfPath("./testdata/sdk_config.yml"),
 )

 if err != nil {
  return nil, err
 }

 //启用证书压缩（开启证书压缩可以减小交易包大小，提升处理性能）
 err = chainClient.EnableCertHash()
 if err != nil {
  return nil, err
 }

 return chainClient, nil
}
```

#### 部署wasm合约

下文，将演示通过sdk部署wasm合约，

> `sdk_user_contract_claim_test.go`

```go
func testUserContractClaimCreate(t *testing.T, client *ChainClient,
 admin1, admin2, admin3, admin4 *ChainClient, withSyncResult bool, isIgnoreSameContract bool) {

 resp, err := createUserContract(client, admin1, admin2, admin3, admin4,
  claimContractName, claimVersion, claimByteCodePath, common.RuntimeType_WASMER, []*common.KeyValuePair{}, withSyncResult)
 if !isIgnoreSameContract {
  require.Nil(t, err)
 }

 fmt.Printf("CREATE claim contract resp: %+v\n", resp)
}

func createUserContract(client *ChainClient, admin1, admin2, admin3, admin4 *ChainClient,
 contractName, version, byteCodePath string, runtime common.RuntimeType, kvs []*common.KeyValuePair, withSyncResult bool) (*common.TxResponse, error) {

 payloadBytes, err := client.CreateContractCreatePayload(contractName, version, byteCodePath, runtime, kvs)
 if err != nil {
  return nil, err
 }

 // 各组织Admin权限用户签名
 signedPayloadBytes1, err := admin1.SignContractManagePayload(payloadBytes)
 if err != nil {
  return nil, err
 }

 signedPayloadBytes2, err := admin2.SignContractManagePayload(payloadBytes)
 if err != nil {
  return nil, err
 }

 signedPayloadBytes3, err := admin3.SignContractManagePayload(payloadBytes)
 if err != nil {
  return nil, err
 }

 signedPayloadBytes4, err := admin4.SignContractManagePayload(payloadBytes)
 if err != nil {
  return nil, err
 }

 // 收集并合并签名
 mergeSignedPayloadBytes, err := client.MergeContractManageSignedPayload([][]byte{signedPayloadBytes1,
  signedPayloadBytes2, signedPayloadBytes3, signedPayloadBytes4})
 if err != nil {
  return nil, err
 }

 // 发送创建合约请求
 resp, err := client.SendContractManageRequest(mergeSignedPayloadBytes, createContractTimeout, withSyncResult)
 if err != nil {
  return nil, err
 }

 err = checkProposalRequestResp(resp, true)
 if err != nil {
  return nil, err
 }

 return resp, nil
```

#### 调用wasm合约

下文，将演示通过sdk调用wasm合约，

> `sdk_user_contract_claim_test.go`

```go
func testUserContractClaimInvoke(client *ChainClient,
 method string, withSyncResult bool) (string, error) {

 curTime := fmt.Sprintf("%d", CurrentTimeMillisSeconds())
 fileHash := uuid.GetUUID()
 params := map[string]string{
  "time":      curTime,
  "file_hash": fileHash,
  "file_name": fmt.Sprintf("file_%s", curTime),
 }

 err := invokeUserContract(client, claimContractName, method, "", params, withSyncResult)
 if err != nil {
  return "", err
 }

 return fileHash, nil
}

func invokeUserContract(client *ChainClient, contractName, method, txId string, params map[string]string, withSyncResult bool) error {

 resp, err := client.InvokeContract(contractName, method, txId, params, -1, withSyncResult)
 if err != nil {
  return err
 }

 if resp.Code != common.TxStatusCode_SUCCESS {
  return fmt.Errorf("invoke contract failed, [code:%d]/[msg:%s]\n", resp.Code, resp.Message)
 }

 if !withSyncResult {
  fmt.Printf("invoke contract success, resp: [code:%d]/[msg:%s]/[txId:%s]\n", resp.Code, resp.Message, resp.ContractResult.Result)
 } else {
  fmt.Printf("invoke contract success, resp: [code:%d]/[msg:%s]/[contractResult:%s]\n", resp.Code, resp.Message, resp.ContractResult)
 }

 return nil
}
```

#### 创建及调用evm合约

> `sdk-go/examples/user_contract_evm_balance/main.go`(<https://git.chainmaker.org.cn/chainmaker/sdk-go/-/blob/master/examples/user_contract_evm_balance/main.go>)

### 更多示例和用法

> 更多示例和用法，请参看单元测试用例
<br>示例: (<https://git.chainmaker.org.cn/chainmaker/sdk-go/-/blob/master/examples>)

| 功能     | 单测代码                      |
| -------- | ----------------------------- |
| 用户合约 | `sdk_user_contract_test.go`   |
| 系统合约 | `sdk_system_contract_test.go` |
| 链配置   | `sdk_chain_config_test.go`    |
| 证书管理 | `sdk_cert_manage_test.go`     |
| 消息订阅 | `sdk_subscribe_test.go`       |

### demo

sdk-go demo参考：

[文件 · v2.3.2 · chainmaker / sdk-go-demo · ChainMaker](https://git.chainmaker.org.cn/chainmaker/sdk-go-demo/-/tree/v2.3.2)

## 接口说明

请参看：[《chainmaker-go-sdk》](https://git.chainmaker.org.cn/chainmaker/sdk-go/-/blob/v2.3.6/sdk_interface.md)
所有go-sdk接口参考： https://git.chainmaker.org.cn/chainmaker/sdk-go/-/blob/v2.3.6/sdk_interface.go
