# Cert模式长安链介绍

## 适用场景说明

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
| 共识算法     | TBFT、RAFT、MaxBFT                    | TBFT、DPOS                  | TBFT、RAFT                   |
| 适用场景     | 联盟链                                 | 公链                         | 联盟链                         |

在PermissionWithCert（**数字证书**）账户模式下，一条链由多个组织构成，每个组织在创世块（bc.yml)中的`trust roots` 字段定义组织根证书。拥有某根组织证书颁发的合法节点证书的节点，即可加入网络，拥有某根组织颁发的合法用户证书的用户，则被视为该链的用户。一个组织可以有多个节点和用户，如下图所示：

<img loading="lazy" src="../images/Identity-UserSystem.png" style="zoom:70%;" />


## 账户角色与权限说明
### 角色类型

<span id="role_type"></span>

长安链中，定义了以下几种角色类型：

- 共识节点 `consensus`：有权参与区块共识流程的链上节点；
- 同步节点`common`：无权参与区块共识流程，但可在链上同步数据的节点；
- 管理员 `admin`：可代表组织进行链上治理的用户；
- 普通用户 `client`：无权进行链上治理，但可发送和查询交易的用户；
- 轻节点用户```light```：无权进行链上治理，无权发送交易，只可查询、订阅自己组织的区块、交易数据，属于SPV轻节点用户（详情见[轻节点](../tech/SPV轻节点.md)）。

### 权限说明

长安链2.X中，PermissionedWithCert模式采用了基于角色的权限控制模型，如下图所示。具体来说，需将用户绑定到某一组织的一个或多个角色上，再将链上资源（即链上需要进行权限管理的操作）绑定至特定的组织、角色，即可完成资源和用户的关联。

<img loading="lazy" src="../images/Identity-PermissionControlModel.png" style="width:900px;" />

权限管理及修改可参考：[权限管理](通过CMC管理Cert模式链.html#权限管理)

##  链配置文件说明
### 创世块配置文件说明

**配置文件：bc.yml**

- auth_type：身份模式

  permissionedWithCert：面向强权限控制场景，基于数字证书的用户标识体系、基于角色的权限控制体系。

* consensus：共识配置

  - nodes：共识节点列表

    - org_id：组织id

    - node_id：共识节点id列表，是由共识节点TLS证书里的公钥算出的唯一id

* trust_roots：信任根配置列表

  - org_id：组织id

  - root：组织CA证书所在路径列表

### 创世块配置文件示例

```yaml
chain_id: chain1                    # 链标识
version: v1.0.0                     # 链版本
sequence: 0                         # 配置版本
auth_type: "permissionedWithCert"   # 认证类型

crypto:
  hash: SHA256

# 合约支持类型的配置
contract:
  enable_sql_support: false

# 虚拟机配置
vm:
  # 虚拟机支持列表
  support_list:
    - "wasmer"
    - "gasm"
    - "evm"
    - "wxvm"

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
    consensus_message_turbo: true # 是否开启共识报文压缩
    retry_time: 500 # 根据交易ID列表从交易池获取交易的重试次数
    retry_interval: 20 # 重试间隔，单位:ms

#共识配置
consensus:
  # 共识类型(0-SOLO,1-TBFT,2-MBFT,3-HOTSTUFF,4-RAFT,5-DPOS)
  type: 1
  # 共识节点列表，组织必须出现在trust_roots的org_id中，每个组织可配置多个共识节点，节点地址采用libp2p格式
  nodes:
    - org_id: "wx-org1.chainmaker.org"
      node_id:
        - "QmcQHCuAXaFkbcsPUj7e37hXXfZ9DdN7bozseo5oX4qiC4"
    - org_id: "wx-org2.chainmaker.org"
      node_id:
        - "QmeyNRs2DwWjcHTpcVHoUSaDAAif4VQZ2wQDQAUNDP33gH"
    - org_id: "wx-org3.chainmaker.org"
      node_id:
        - "QmXf6mnQDBR9aHauRmViKzSuZgpumkn7x6rNxw1oqqRr45"
    - org_id: "wx-org4.chainmaker.org"
      node_id:
        - "QmRRWXJpAVdhFsFtd9ah5F4LDQWFFBDVKpECAF8hssqj6H"
        
  ext_config: # 扩展字段，记录难度、奖励等其他类共识算法配置
    - key: aa
      value: chain01_ext11

# 信任组织和根证书
trust_roots:
  - org_id: "wx-org1.chainmaker.org"
    root:
      - "../config/wx-org1/certs/ca/wx-org1.chainmaker.org/ca.crt"
  - org_id: "wx-org2.chainmaker.org"
    root:
      - "../config/wx-org2/certs/ca/wx-org2.chainmaker.org/ca.crt"
  - org_id: "wx-org3.chainmaker.org"
    root:
      - "../config/wx-org3/certs/ca/wx-org3.chainmaker.org/ca.crt"
  - org_id: "wx-org4.chainmaker.org"
    root:
      - "../config/wx-org4/certs/ca/wx-org4.chainmaker.org/ca.crt"

# 权限配置（只能整体添加、修改、删除）
resource_policies:
  - resource_name: CHAIN_CONFIG-NODE_ID_UPDATE
    policy:
      rule: SELF # 规则（ANY，MAJORITY...，全部大写，自动转大写）
      org_list: # 组织名称（组织名称，区分大小写）
      role_list: # 角色名称（role，自动转大写）
        - admin
  - resource_name: CHAIN_CONFIG-TRUST_ROOT_ADD
    policy:
      rule: MAJORITY
      org_list:
      role_list:
        - admin
  - resource_name: CHAIN_CONFIG-CERTS_FREEZE
    policy:
      rule: ANY
      org_list:
      role_list:
        - admin
        
disabled_native_contract:
 # - CONTRACT_NAME  # 通过指定系统合约名字来禁用系统合约
```


### 节点配置文件说明

**配置文件：chainmaker.yml**

- auth_type：身份模式

  permissionedWithCert：面向强权限控制场景，基于数字证书的用户标识体系、基于角色的权限控制体系。

- node：节点配置

  - priv_key_file：节点SIGN证书的私钥地址
  - cert_file：节点SIGN证书的地址

- net：网络配置

  - tls：TLS配置
    - priv_key_file：节点TLS证书的私钥地址
    - cert_file：节点TLS证书的地址



### 节点配置文件示例
```yaml
auth_type: "permissionedWithCert"      # permissionedWithCert / permissionedWithKey / public

log:
  config_file: ../config/wx-org1/log.yml          # config file of logger configuration.

blockchain:
  - chainId: chain1
    genesis: ../config/wx-org1/chainconfig/bc1.yml

node:
  # 节点类型：full
  type:              full
  org_id:            wx-org1.chainmaker.org
  priv_key_file:     ../config/wx-org1/certs/node/consensus1/consensus1.sign.key
  cert_file:         ../config/wx-org1/certs/node/consensus1/consensus1.sign.crt
  signer_cache_size: 1000
  cert_cache_size:   1000

net:
  provider: LibP2P
  listen_addr: /ip4/0.0.0.0/tcp/11301
  seeds:
    - "/ip4/127.0.0.1/tcp/11301/p2p/QmcQHCuAXaFkbcsPUj7e37hXXfZ9DdN7bozseo5oX4qiC4"
    - "/ip4/127.0.0.1/tcp/11302/p2p/QmNdgWgD2QSu769yCwFCnwhVKhjWzyjA3PgSXL7ZJBrUoA"
    - "/ip4/127.0.0.1/tcp/11303/p2p/QmXf6mnQDBR9aHauRmViKzSuZgpumkn7x6rNxw1oqqRr45"
    - "/ip4/127.0.0.1/tcp/11304/p2p/QmRRWXJpAVdhFsFtd9ah5F4LDQWFFBDVKpECAF8hssqj6H"
  tls:
    enabled: true
    priv_key_file: ../config/wx-org1/certs/node/consensus1/consensus1.tls.key
    cert_file:     ../config/wx-org1/certs/node/consensus1/consensus1.tls.crt

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
    mode:           twoway
    priv_key_file:  ../config/wx-org1/certs/node/consensus1/consensus1.tls.key
    cert_file:      ../config/wx-org1/certs/node/consensus1/consensus1.tls.crt

monitor:
  enabled: false
  port: 14321

pprof:
  enabled: false
  port: 24321

storage:
  store_path: ../data/org1/ledgerData1
  blockdb_config:
    provider: leveldb
    leveldb_config:
      store_path: ../data/org1/blocks
      write_buffer_size: 1024
      block_write_buffer_size: 1024
  statedb_config:
    provider: leveldb
    leveldb_config:
      store_path: ../data/org1/state
  historydb_config:
    provider: leveldb
    leveldb_config:
      store_path: ../data/org1/history
  resultdb_config:
    provider: leveldb
    leveldb_config:
      store_path: ../data/org1/result
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

