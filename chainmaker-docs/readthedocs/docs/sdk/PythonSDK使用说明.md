#  Python SDK 使用说明

## 概述
目前Python SDK 除了**不支持隐私计算API**、**Hibe**及**国密SSL通信**外，支持python SDK所有接口，
主要概念包括：

- ChainClient: 链客户端对象
- User: 客户端用户对象
- Node: 客户端连接节点对象

几乎所有操作需要使用ChainClient对象完成，同时utils提供了文件、证书、EVM等一些实用方法。

## 环境准备

### 环境依赖

**Python**

版本为Python3.9.0以上

下载地址：https://www.python.org/downloads/

若已安装，请通过命令查看版本：

```bash
$ python3 --version
Python 3.9.17
```

### 安装方式
1. 使用pip命令，在线安装(需要安装Git)
```bash
$ pip3 install git+https://git.chainmaker.org.cn/chainmaker/sdk-python.git
```
2. 克隆或下载[sdk-python](https://git.chainmaker.org.cn/chainmaker/sdk-python)源码，在项目根目录下运行
```bash
$ python3 setup.py install
```

## 怎么使用SDK

### 示例代码

> 注： 下方文档示例可能过时，以gitlab示例为准。
>
> evm和其他合约使用方法在构建参数时有区别。
>
> evm的可参考示例：[示例Contract](https://git.chainmaker.org.cn/chainmaker/sdk-python/-/tree/v2.2.0/#创建合约)

#### 创建节点

> 注️：需要拷贝目标链chainmaker-go/build/crypto-config到脚本所在目录

```python
from chainmaker.node import Node

# 创建节点
node = Node(
    node_addr='127.0.0.1:12301',
    conn_cnt=10,
    enable_tls=True,
    cas=['./testdata/crypto-config/wx-org1.chainmaker.org/ca', './testdata/crypto-config/wx-org2.chainmaker.org/ca'],
    tls_host_name='chainmaker.org'
)
```

#### 以参数形式创建ChainClient

> 更多内容请参考：`tests/test_chain_client.py`
>
> 注：示例中证书采用路径方式去设置，也可以使用证书内容去设置，具体请参考`createClientWithCaCerts`方法

```python
from chainmaker.chain_client import ChainClient
from chainmaker.node import Node
from chainmaker.user import User
from chainmaker.utils import file_utils
    
user = User('wx-org1.chainmaker.org',
            sign_key_bytes=file_utils.read_file_bytes('./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.tls.key'),
            sign_cert_bytes=file_utils.read_file_bytes('./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.tls.crt'),
            tls_key_bytes=file_utils.read_file_bytes('./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.sign.key'),
            tls_cert_bytes=file_utils.read_file_bytes('./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.sign.crt')
            )

node = Node(
    node_addr='127.0.0.1:12301',
    conn_cnt=1,
    enable_tls=True,
    trust_cas=[
        file_utils.read_file_bytes('./testdata/crypto-config/wx-org1.chainmaker.org/ca/ca.crt'),
        file_utils.read_file_bytes('./testdata/crypto-config/wx-org2.chainmaker.org/ca/ca.crt')
    ],
    tls_host_name='chainmaker.org'
)

cc = ChainClient(chain_id='chain1', user=user, nodes=[node])
print(cc.get_chainmaker_server_version())
```

#### 以配置文件形式创建ChainClient

> 注：参数形式和配置文件形式两个可以同时使用，同时配置时，以参数传入为准

```python
from chainmaker.chain_client import ChainClient

# ./testdata/sdk_config.yml 中私钥/证书等如果使用相对路径应相对于当前运行起始目录
cc = ChainClient.from_conf('./testdata/sdk_config.yml')
```

> [配置文件 sdk_config.yml 格式参考](../recovery/配置文件一览.html#sdk-config-yml)

#### 创建合约

```python
from google.protobuf import json_format
from chainmaker.chain_client import ChainClient
from chainmaker.utils.evm_utils import calc_evm_contract_name
from chainmaker.keys import RuntimeType

endorsers_config = [{'org_id': 'wx-org1.chainmaker.org',
                  'user_sign_crt_file_path': './testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt',
                  'user_sign_key_file_path': './testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key'},
                  {'org_id': 'wx-org2.chainmaker.org',
                  'user_sign_crt_file_path': './testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt',
                  'user_sign_key_file_path': './testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key'},
                  {'org_id': 'wx-org3.chainmaker.org',
                  'user_sign_crt_file_path': './testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt',
                  'user_sign_key_file_path': './testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key'},
                  ]

cc = ChainClient.from_conf('./testdata/sdk_config.yml')

def create_contract(contract_name: str, version: str, byte_code_path: str, runtime_type: RuntimeType, params: dict = None, 
                    with_sync_result=True) -> dict:
    """创建合约"""
    # 创建请求payload
    payload = cc.create_contract_create_payload(contract_name, version, byte_code_path, runtime_type, params)
    # 创建背书
    endorsers = cc.create_endorsers(payload, endorsers_config)
    # 携带背书发送请求
    res = cc.send_request_with_sync_result(payload, with_sync_result=with_sync_result, endorsers=endorsers)
    # 交易响应结构体转为字典格式
    return json_format.MessageToDict(res)

# 创建WASM合约，本地合约文件./testdata/claim-wasm-demo/rust-fact-2.0.0.wasm应存在
result1 = create_contract('fact', '1.0', './testdata/claim-wasm-demo/rust-fact-2.0.0.wasm', RuntimeType.WASMER, {})
print(result1)

# 创建EVM合约，本地合约文件./testdata/balance-evm-demo/ledger_balance.bin应存在

contract_name = calc_evm_contract_name('balance001')
result2 = create_contract(contract_name, '1.0', './testdata/balance-evm-demo/ledger_balance.bin', RuntimeType.EVM)
print(result2)
```

#### 调用合约

```python
from google.protobuf import json_format
from chainmaker.chain_client import ChainClient
from chainmaker.utils.evm_utils import calc_evm_contract_name, calc_evm_method_params

# 创建客户端
cc = ChainClient.from_conf('./testdata/sdk_config.yml')

# 调用WASM合约
res1 = cc.invoke_contract('fact', 'save', {"file_name":"name007","file_hash":"ab3456df5799b87c77e7f88","time":"6543234"},
                          with_sync_result=True)
# 交易响应结构体转为字典格式
print(json_format.MessageToDict(res1))

# 调用EVM合约
evm_contract_name = calc_evm_contract_name('balance001')
evm_method, evm_params = calc_evm_method_params('updateBalance', [{"uint256": "10000"}, {"address": "0xa166c92f4c8118905ad984919dc683a7bdb295c1"}])
res2 = cc.invoke_contract(evm_contract_name, evm_method, evm_params, with_sync_result=True)
# 交易响应结构体转为字典格式
print(json_format.MessageToDict(res2))
```

#### 更多示例和用法

> 更多示例和用法，请参考单元测试用例
>
> [tests · v2.2.0 · chainmaker / sdk-python · GitLab](https://git.chainmaker.org.cn/chainmaker/sdk-python/-/tree/v2.2.0/tests)

| 功能     | 单测代码                      |
| -------- | ----------------------------- |
| 用户合约 | `tests/test_user_contract.py`   |
| 系统合约 | `tests/test_system_contract.py` |
| 链配置   | `tests/test_chain_config.py`    |
| 证书管理 | `tests/test_cert_manage.py`     |
| 消息订阅 | `tests/test_user_contract.py`   |

### demo

demo示例参考：

[文件 · v2.3.2 · chainmaker / sdk-python-demo · ChainMaker](https://git.chainmaker.org.cn/chainmaker/sdk-python-demo/-/tree/v2.3.2)

## 接口说明

请参看：[《chainmaker-python-sdk》](https://git.chainmaker.org.cn/chainmaker/sdk-python/-/blob/master/chainmaker-python-sdk.md)