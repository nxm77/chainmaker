# PK模式长安链介绍

## Public适用场景说明

在长安链2.X版本中, 我们实现了三种身份权限管理模型：

* PermissionedWithCert：基于数字证书的用户标识体系、基于角色的权限控制体系;
* PermissionedWithKey：基于公钥的用户标识体系、基于角色的权限控制体系。
* Public：基于公钥的用户标识体系、基于角色的权限控制体系。

其中，PermissionedWithCert和PermissionedWithKey模式面向强权限控制场景(联盟链)，Public模式面向弱权限控制场景(公链)

### 不同链账户模式对比

| 对比项\身份模式 | PermissionWithCert                  | Public                     | PermissionWithKey           |
|----------|-------------------------------------|----------------------------|-----------------------------|
| 模式名称     | [证书模式](./Cert模式长安链介绍.html)          | [公钥模式](./PK模式长安链介绍.html)   | [公钥注册模式](./PWK模式长安链介绍.html) |
| 模式简称     | cert模式                              | pk模式                       | pwk模式                       |
| 账户类型     | 节点账户(共识节点、同步节点、轻节点), 用户账户(管理员、普通用户) | 节点账户(共识节点), 用户账户(管理员、普通用户) | 同证书模式                       | 
| 账户标识     | 数字证书                                | 公钥/地址                      | 公钥/地址                       | 
| 是否需要准入   | 是，证书需要CA签发                          | 否，普通用户可直接调用合约              | 是，账户需要管理员在链上注册              |
| 账户与组织关系  | 账户属于某个组织                            | 账户无组织概念                    | 账户属于某个组织                    |  
| 适用链类型    | 联盟链                                 | 公链                         | 联盟链                         |
| 共识算法     | TBFT、RAFT、MaxBFT                    | TBFT、DPOS                  | TBFT、RAFT                   |

在Public（**公钥**）账户模式下，一条链仅存在一个默认组织（public），在配置文件（bc.yml）的`trust roots` 字段中定义管理员公钥,在`consensus`里定义组织的共识节点列表，一条链可以有多个链管理员、多个共识节点。普通用户加入没有门槛，任意公私钥创建的链账户都可以成为链的普通用户,如下图所示：

<img loading="lazy" src="../images/Identity-UserSystem-public.png" style="zoom:70%;" />

## 账户角色与权限说明
### 角色类型

<span id="role_type"></span>

长安链中，定义了以下几种角色类型：

- 共识节点 `consensus`：有权参与区块共识流程的链上节点；
- 管理员 `admin`：可代表组织进行链上治理的用户；
- 普通用户 `client`：无权进行链上治理，但可发送和查询交易的用户。

### 权限说明
权限标识方法：采用**公钥**的标识方式，能够避免像数字证书体系那样繁琐的签发流程，使用户加入区块链网络更加简单快捷，但无法承载用户或节点的组织信息以及角色信息。

权限管理及修改可参考：[权限管理](管理PK账户模式的链.html#权限管理)


## 链配置文件说明

### 创世块配置文件说明

**配置文件：bc.yml**

- auth_type：身份模式

  public：面向弱权限控制场景，基于公钥的用户标识体系、基于角色的权限控制体系。

* consensus：共识配置

  - nodes：共识节点列表

    **TBFT共识模式需要配置共识节点列表，DPOS共识模式不需要配置。**

* trust_roots：信任根配置列表

  **默认只使用列表下第一个配置**

  - org_id：public模式配置标识（需要填写public）

  ```yaml
  org_id: "public"
  ```

  - root：链管理员公钥所在路径列表

### 创世块配置文件示例


```yaml
chain_id: chain1                      # 链标识
version: v1.0.0                       # 链版本
sequence: 0                           # 配置版本
auth_type: "public"                   # 认证类型 permissionedWithCert / permissionedWithKey / public

crypto:
  hash: SHA256

# 合约支持类型的配置
contract:
  enable_sql_support: false

# 交易、区块相关配置
block:
  tx_timestamp_verify: true # 是否需要开启交易时间戳校验
  tx_timeout: 600  # 交易时间戳的过期时间(秒)
  block_tx_capacity: 100  # 区块中最大交易数
  block_size: 10  # 区块最大限制，单位MB
  block_interval: 2000 # 出块间隔，单位:ms

# core模块
core:
  tx_scheduler_timeout: 10 #  [0, 60] 交易调度器从交易池拿到交易后, 进行调度的时间
  tx_scheduler_validate_timeout: 10 # [0, 60] 交易调度器从区块中拿到交易后, 进行验证的超时时间
  consensus_turbo_config:
    consensus_message_turbo: false # 是否开启共识报文压缩
    retry_time: 500 # 根据交易ID列表从交易池获取交易的重试次数
    retry_interval: 20 # 重试间隔，单位:ms

#共识配置
consensus:
  # 共识类型(0-SOLO,1-TBFT,2-MBFT,3-HOTSTUFF,4-RAFT,5-DPOS)
  type: 5
  ext_config: # 扩展字段，记录难度、奖励等其他类共识算法配置
    - key: aa
      value: chain01_ext11
  dpos_config: # DPoS
    #ERC20合约配置
    - key: erc20.total
      value: "10000000"
    - key: erc20.owner
      value: "6CeSsjU5M62Ee3Gx9umUX6nXJoaBkWYufQdTZqEJM5di"
    - key: erc20.decimals
      value: "18"
    - key: erc20.account:DPOS_STAKE
      value: "10000000"
    #Stake合约配置
    - key: stake.minSelfDelegation
      value: "2500000"
    - key: stake.epochValidatorNum
      value: "4"
    - key: stake.epochBlockNum
      value: "10"
    - key: stake.completionUnbondingEpochNum
      value: "1"
    - key: stake.candidate:6CeSsjU5M62Ee3Gx9umUX6nXJoaBkWYufQdTZqEJM5di
      value: "2500000"
    - key: stake.candidate:F5tJ4ca4vdbuyffpc1Szw3WHU3caGaTVAh52MRMS4qBt
      value: "2500000"
    - key: stake.candidate:FxfunVWGkKgYMjngxMtLkd4pUNYVNAHNAqiDqopg5zdw
      value: "2500000"
    - key: stake.candidate:DYt7DfcZnqKNpjgyJ6tU6GFixNfLMkkmnqdwB3NNiAP7
      value: "2500000"

    - key: stake.nodeID:6CeSsjU5M62Ee3Gx9umUX6nXJoaBkWYufQdTZqEJM5di
      value: "QmZcFcJFYYoZ3FNNGL88QaszUZwFwuBdFqYh6yPzJURc3s"
    - key: stake.nodeID:F5tJ4ca4vdbuyffpc1Szw3WHU3caGaTVAh52MRMS4qBt
      value: "QmXwtuPemSgH5ypzoKvcLdCLbd9jZ25FbpNf7VPjHF3HMS"
    - key: stake.nodeID:FxfunVWGkKgYMjngxMtLkd4pUNYVNAHNAqiDqopg5zdw
      value: "QmRmQLHJoqAYGkuLFaNY6HLzwtTNxr45UJsYpSjdKvBQw2"
    - key: stake.nodeID:DYt7DfcZnqKNpjgyJ6tU6GFixNfLMkkmnqdwB3NNiAP7
      value: "QmURUHTGsuzzjgh1Xg6s92G1Q3gK91A6JEZGPfYNWwJMiT"

# 超级管理员
trust_roots:
  - org_id: "public"
    root:
      - "../config-pk/public/admin/admin1/admin1.pem"
      - "../config-pk/public/admin/admin2/admin2.pem"
      - "../config-pk/public/admin/admin3/admin3.pem"
      - "../config-pk/public/admin/admin4/admin4.pem"
```



### 节点配置文件说明

**配置文件：chainmaker.yml**

- auth_type：身份模式

  public：面向弱权限控制场景，基于公钥的用户标识体系、基于角色的权限控制体系。

- node：节点配置

  - priv_key_file：节点私钥地址
  - cert_file：不需要配置

- net：网络配置

  - tls：TLS配置
    - priv_key_file：节点私钥地址
    - cert_file：不需要配置

**注：node和net里需要配置同一个私钥的地址**。

### 节点配置文件示例


```yaml
auth_type: "public"                                                        # permissionedWithCert / permissionedWithKey / public

log:
  config_file: ../config-pk/public/node/node1/log.yml                           # config file of logger configuration.

blockchain:
  - chainId: chain1
    genesis: ../config-pk/public/node/node1/chainconfig/bc1.yml

node:
  # 节点类型：full
  type:              full
  org_id:            wx-org1.chainmaker.org
  priv_key_file:     ../config-pk/public/node/node1/node1.key
  signer_cache_size: 1000
  cert_cache_size:   1000

net:
  provider: LibP2P
  listen_addr: /ip4/0.0.0.0/tcp/11351
  seeds:
    - "/ip4/127.0.0.1/tcp/11351/p2p/QmZcFcJFYYoZ3FNNGL88QaszUZwFwuBdFqYh6yPzJURc3s"
    - "/ip4/127.0.0.1/tcp/11352/p2p/QmXwtuPemSgH5ypzoKvcLdCLbd9jZ25FbpNf7VPjHF3HMS"
    - "/ip4/127.0.0.1/tcp/11353/p2p/QmRmQLHJoqAYGkuLFaNY6HLzwtTNxr45UJsYpSjdKvBQw2"
    - "/ip4/127.0.0.1/tcp/11354/p2p/QmURUHTGsuzzjgh1Xg6s92G1Q3gK91A6JEZGPfYNWwJMiT"
  tls:
    enabled: true
    priv_key_file: ../config-pk/public/node/node1/node1.key

txpool:
  max_txpool_size: 5120 # 普通交易池上限
  max_config_txpool_size: 10 # config交易池的上限
  full_notify_again_time: 30 # 交易池溢出后，再次通知的时间间隔(秒)

rpc:
  provider: grpc
  port: 12301
  tls:
    # TLS模式:
    #   disable - 不启用TLS
    #   oneway  - 单向认证
    #   twoway  - 双向认证
    #mode: disable
    #mode: oneway
    mode: disable

monitor:
  enabled: false
  port: 14321

pprof:
  enabled: false
  port: 24321

storage:
  store_path: ../data/node1/ledgerData1
  blockdb_config:
    provider: leveldb
    leveldb_config:
      store_path: ../data/node1/blocks
  statedb_config:
    provider: leveldb
    leveldb_config:
      store_path: ../data/node1/state
  historydb_config:
    provider: leveldb
    leveldb_config:
      store_path: ../data/node1/history
  resultdb_config:
    provider: leveldb
    leveldb_config:
      store_path: ../data/node1/result
  disable_contract_eventdb: true  #是否禁止合约事件存储功能，默认为true，如果设置为false,需要配置mysql
  contract_eventdb_config:
    provider: sql                 #如果开启contract event db 功能，需要指定provider为sql
    sqldb_config:
      sqldb_type: mysql           #contract event db 只支持mysql
      dsn: root:password@tcp(127.0.0.1:3306)/  #mysql的连接信息，包括用户名、密码、ip、port等，示例：root:admin@tcp(127.0.0.1:3306)/
debug:
  # 是否开启CLI功能，过度期间使用
  is_cli_open: true
  is_http_open: false
```
