
#  Nodejs SDK 使用说明
## 基本概念定义

1. 整体介绍

`SDK`是业务模块与长安链交互的桥梁，支持双向`TLS`认证，提供安全可靠的加密通信信道。

提供的接口，覆盖合约管理、链配置管理、证书管理、多签收集、各类查询操作、事件订阅、数据归档等场景，满足了不同的业务场景需要。

2. 名词概念说明

- **`Node`（节点）**：代表一个链节点的基本信息，包括：节点地址、连接数、是否启用`TLS`认证等信息
- **`ChainClient`（链客户端）**：所有客户端对链节点的操作接口都来自`ChainClient`
- **压缩证书**：可以为`ChainClient`开启证书压缩功能，开启后可以减小交易包大小，提升处理性能

## 环境准备
### 软件依赖

**nodejs**

> nodejs 14.0.0+

下载地址：https://nodejs.org/dist/

若已安装，请通过命令查看版本：

```bash
$ node --version
v14.0.0
```

### 下载安装

```bash
$ git clone -b v2.0.0  --depth=1 https://git.chainmaker.org.cn/chainmaker/sdk-nodejs.git
```

## 怎么使用SDK

### 示例代码

> 注： 下方文档示例可能过时，以gitlab示例为准。
>
> evm和其他合约使用方法在构建参数时有区别。
>
> evm的可参考示例：[TestEvmContract](https://git.chainmaker.org.cn/chainmaker/sdk-java/-/blob/v2.3.1.3/src/test/java/org/chainmaker/sdk/TestEvmContract.java)

#### 创建节点

> 更多内容请参看：`sdkInit.js`

```javascript
// 创建节点
this.node = new Node(nodeConfigArray, timeout);
```

#### 参数形式创建ChainClient

> 更多内容请参看：`sdkInit.js`

```javascript
// 创建ChainClient
const ChainClient = new Sdk(chainID, orgID, userKeyPathFile, userCertPathFile, nodeConfigArray, 30000, archiveConfig);

```

#### 以配置文件形式创建ChainClient

> 更多内容请参看：`sdkinit.js`

```javascript
const ChainClient = new LoadFromYaml(path.join(__dirname, './sdk_config.yaml'));
```

#### 创建合约

> 更多内容请参看：`testUserContractMgr.js`

```javascript
  const testCreateUserContract = async (sdk, contractName, contractVersion, contractFilePath) => {
		const response = await sdk.userContractMgr.createUserContract({
			contractName,
			contractVersion,
			contractFilePath,
			runtimeType: utils.common.RuntimeType.GASM,
			params: {
				key1: 'value1',
				key2: 'value2',
			},
		});
		return response;
	};
```

#### 调用合约

> 更多内容请参看：`testUserContractMgr.js`

```javascript
  const testInvokeUserContract = async (sdk, contractName) => {
		const response = await sdk.callUserContract.invokeUserContract({
			contractName, method: 'save', params: {
				file_hash: '1234567890',
				file_name: 'test.txt',
			},
		});
		return response;
	};
```

### 更多示例和用法

> 更多示例和用法，请参看单元测试用例

安装mocha：

```bash
$ npm install -g mocha
```

使用脚本搭建chainmaker运行环境（4组织4节点），将build文件中的cryptogen复制到当前项目的test/testFile文件中

运行测试命令：

```bash
$ npm test
```

| 功能     | 单测代码            |
| -------- | ------------------- |
| 基础配置 | `sdkInit.js`        |
| 用户合约 | `userContract.js`   |
| 系统合约 | `systemContract.js` |
| 链配置   | `chainConfig.js`    |
| 证书管理 | `cert.js`           |
| 消息订阅 | `subscribe.js`      |

## 接口说明

请参看：[《chainmaker-nodejs-sdk》](https://git.chainmaker.org.cn/chainmaker/sdk-nodejs/-/blob/v2.0.0/sdk_interface.md)

## demo
参考： https://git.chainmaker.org.cn/chainmaker/sdk-nodejs-demo