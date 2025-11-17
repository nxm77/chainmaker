# 管理PK账户模式的链

## 从零生成链账户部署PK模式链
如果对链账户的公私钥对无特殊要求则可基于[部署PK账户模式的链](../instructions/启动PK模式的链.md)搭建链，如果有特殊要求则参考此文章。

阅读本章节前，请先阅读[部署PK账户模式的链](../instructions/启动PK模式的链.md)，并确保已完成该章节相关的环境依赖准备。

下文将展示通过长安链cmc工具生成新的公私钥对，并基于此搭建新的链；如果有其他生成公私钥对的工具，用户可自行采用，搭建链的流程相似。

### 生成链账户

pk模式无证书概念，因此对于节点、管理员以及用户而言，都是公私钥的形式。 节点账户的地址需要添加到共识列表，管理员账户需要添加到管理员列表；
而普通用户账户，可以直接调用管理员部署的合约，无需注册。

#### 部署cmc命令行工具

长安链命令行工具cmc安装, 请参考[长安链命令行工具](../dev/命令行工具pk.html#编译&配置)

#### 生成节点账户
```shell
  # 生成节点账户私钥
$ ./cmc key gen -a ECC_P256 -p ./ -n consensus1.key
 
  # 导出节点账户公钥
$ ./cmc key export_pub -k consensus1.key -n consensus1.pem
 
 # 计算节点账户地址
$ ./cmc cert nid --node-pk-path=./consensus1.pem
node id : QmeqnZEgGeQYyc4qX92XV3SxafqRJqCQ9388jWp2N1oA93
```

#### 生成admin账户
```shell
# 生成管理员私钥
 ./cmc key gen -a ECC_P256 -p ./ -n admin.key

# 导出管理员公钥
 ./cmc key export_pub -k admin.key -n admin.pem
 
# 查看管理员公钥（账户）
$ cat admin.pem
-----BEGIN PUBLIC KEY-----
MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEsIYfDrQCt1/H8Yj5KKKD+uO28zz7
nTovDgim/jezoGdpmNOUp6lwrN47pxBUpnxEXIqHBwz8uaVR1z3y9kDvTg==
-----END PUBLIC KEY-----
```

重复以上步骤分别生成4个节点账户和4个管理员账户。


### 修改链配置
先将[部署PK账户模式的链](../instructions/启动PK模式的链.md)文章处所生成的相关配置文件内的链账户信息替换成上文新生成的。

**配置文件目录结构说明**
```shell
./config
└── node1
    ├── admin
    │   ├── admin1 #optional
    │   │   ├── admin1.key 
    │   │   └── admin1.pem
    │   ├── admin2
    │   │   ├── admin2.key
    │   │   └── admin2.pem
    │   ├── admin3
    │   │   ├── admin3.key
    │   │   └── admin3.pem
    │   ├── admin4
    │   │   ├── admin4.key
    │   │   └── admin4.pem
    ├── chainconfig
    │   └── bc1.yml # 更新后的链配置文件，参考上节说明
    ├── chainmaker.yml # 更新后的链配置文件，参考上节说明
    ├── log.yml
    ├── node1.key #节点1密钥
    ├── node1.nodeid #节点1的nodeid
    ├── node1.pem #节点1的公钥
    └── user #optional
        └── client1
            ├── client1.addr
            ├── client1.key
            └── client1.pem

```
#### 替换公私钥对
在启动节点之前，需要将以上部署过程中用到的各个节点公私钥、管理员公私钥替换为以上用cmc生成的公私钥对。

#### 修改配置文件
将上文获取的节点账户和管理员账户，配置到链配置文件`bc1.yml`的`trust_roots`里，并将`bc1.yml`和`chainmaker.yml`中的`nodes.node_id`替换为以上获取的节点nodeId。

**配置文件修改位置如下**

- bc1.yml（链配置文件）

```yaml
#共识配置
consensus:
  # Consensus type: 1-TBFT,5-DPOS
  type: 1
  nodes:
    - org_id: "public"
      node_id:
        - "QmeqnZEgGeQYyc4qX92XV3SxafqRJqCQ9388jWp2N1oA93"
        - "QmSXhWkujKh2PEN5tFuiTBnSbkx7vN6P9zPoa6RViLdpdA"
        - "QmfR4jNLsBK3FedeCLyTmzaCeaECK3VjDRoY6XcjiUpWYJ"
        - "QmRKZyDH89CH3zJaSC5VHn9tgcBbs7jc1mDVd9QgyzmKan"
      
trust_roots:
  - org_id: "public"
    root:
      - "../config/node1/admin/admin1/admin1.pem"
      - "../config/node1/admin/admin2/admin2.pem"
      - "../config/node1/admin/admin3/admin3.pem"
      - "../config/node1/admin/admin4/admin4.pem"
```

- chainmaker.yml （节点配置文件）

```yaml
# Network Settings
net:
  # Network provider, can be libp2p or liquid.
  # libp2p: using libp2p components to build the p2p module.
  # liquid: a new p2p module we build from 0 to 1.
  # This item must be consistent across the blockchain network.
  provider: LibP2P

  # The address and port the node listens on.
  # By default, it uses 0.0.0.0 to listen on all network interfaces.
  listen_addr: /ip4/0.0.0.0/tcp/11301

  # The seeds peer list used to join in the network when starting.
  # The connection supervisor will try to dial seed peer whenever the connection is broken.
  # Example ip format: "/ip4/127.0.0.1/tcp/11301/p2p/"+nodeid
  # Example dns format："/dns/cm-node1.org/tcp/11301/p2p/"+nodeid
  seeds:
    - "/ip4/127.0.0.1/tcp/11301/p2p/QmeqnZEgGeQYyc4qX92XV3SxafqRJqCQ9388jWp2N1oA93"
    - "/ip4/127.0.0.1/tcp/11302/p2p/QmSXhWkujKh2PEN5tFuiTBnSbkx7vN6P9zPoa6RViLdpdA"
    - "/ip4/127.0.0.1/tcp/11303/p2p/QmfR4jNLsBK3FedeCLyTmzaCeaECK3VjDRoY6XcjiUpWYJ"
    - "/ip4/127.0.0.1/tcp/11304/p2p/QmRKZyDH89CH3zJaSC5VHn9tgcBbs7jc1mDVd9QgyzmKan"

```

同时各个节点部署目录下的`bc.yml`和`chainmaker.yml`配置文件按照以上说明进行修改。


### 节点部署和链启动 

节点部署启动与[部署PK账户模式的链](../instructions/启动PK模式的链.md)文章一致，可自行参考。

按照以上方式启动org1、org2、org3和org4下的共识节点，等到所有节点建立连接，表明区块链网络部署成功，并可以对外提供区块链服务。  
可通过遗下命令查看节点日志，若看到all necessary peers connected则表示节点已经准备就绪。

```shell
[INFO]  [Net]   libp2pnet/libp2p_connection_supervisor.go:116   [ConnSupervisor] all necessary peers connected.
```

### 部署/调用智能合约
链部署后，再进行部署/调用合约测试，以验证链是否正常运行。部署合约的使用教程可详见：[部署示例合约](./部署示例合约.md)。



## 链节点增删管理

### 同步节点的增删
在pk模式下，同步节点可以自由加入和退出网络，并同步账本，但不参与共识。

#### 新增同步节点（待补充）
1、先生成新节点的节点账户，可参考[生成节点账户](#生成节点账户)

2、再基于所要加入的链的创世区块配置文件。bc.yml，启动节点，节点启动成功后将自动加入链网络。

#### 更新、删除、查询同步节点
暂不支持同步节点的更新、查询、删除功能，原则上同步节点停止运行则视为自动退出。


### 共识节点的增删
如果节点要从同步节点升级为共识节点，则需要用管理员账号进行操作。

#### 通过CMC工具管理共识节点

- 增加共识节点
```shell
./cmc client chainconfig consensusnodeid add \
--sdk-conf-path=./testdata/sdk_config_pk.yml \
--user-signkey-file-path=./testdata/crypto-config/node1/admin/admin1/admin1.key \
--admin-key-file-paths=./testdata/crypto-config/node1/admin/admin1/admin1.key,./testdata/crypto-config/node1/admin/admin2/admin2.key,./testdata/crypto-config/node1/admin/admin3/admin3.key \
--node-id=QmeqnZEgGeQYyc4qX92XV3SxafqRJqCQ9388jWp2N1oA93 \  #这里为新添加节点的nodeId
--node-org-id=public
```

- 更新共识节点Id
```shell
./cmc client chainconfig consensusnodeid update \
--sdk-conf-path=./testdata/sdk_config_pk.yml \
--user-signkey-file-path=./testdata/crypto-config/node1/admin/admin1/admin1.key \
--admin-key-file-paths=./testdata/crypto-config/node1/admin/admin1/admin1.key,./testdata/crypto-config/node1/admin/admin2/admin2.key,./testdata/crypto-config/node1/admin/admin3/admin3.key \
--node-id=QmXxeLkNTcvySPKMkv3FUqQgpVZ3t85KMo5E4cmcmrexrC \
--node-id-old=QmeqnZEgGeQYyc4qX92XV3SxafqRJqCQ9388jWp2N1oA93 \ 
--node-org-id=public
```

- 删除共识节点
```shell
./cmc client chainconfig consensusnodeid remove \
--sdk-conf-path=./testdata/sdk_config_pk.yml \
--user-signkey-file-path=./testdata/crypto-config/node1/admin/admin1/admin1.key \
--admin-key-file-paths=./testdata/crypto-config/node1/admin/admin1/admin1.key,./testdata/crypto-config/node1/admin/admin2/admin2.key,./testdata/crypto-config/node1/admin/admin3/admin3.key \
--node-id=QmeqnZEgGeQYyc4qX92XV3SxafqRJqCQ9388jWp2N1oA93 \
--node-org-id=public
```

- 查询共识节点  
  共识节点可以通过链配置查询，在`consensus`字段下会返回当前区块链网络中的共识节点列表，命令如下：
```shell
./cmc client chainconfig query \
--sdk-conf-path=./testdata/sdk_config_pk.yml

  "consensus": {
    "nodes": [
      {
        "node_id": [
          "QmeqnZEgGeQYyc4qX92XV3SxafqRJqCQ9388jWp2N1oA93",
          "QmSXhWkujKh2PEN5tFuiTBnSbkx7vN6P9zPoa6RViLdpdA",
          "QmfR4jNLsBK3FedeCLyTmzaCeaECK3VjDRoY6XcjiUpWYJ",
          "QmRKZyDH89CH3zJaSC5VHn9tgcBbs7jc1mDVd9QgyzmKan"
        ],
        "org_id": "public"
      }
    ],
    "type": 1
  },
```


## 链账户的增删

### 普通用户的增删
在长安链pk模式下，普通账户无需注册，使用任意工具生成和链保持一致的密码算法的公私钥对，即可为PK模式链的链用户，可直接调用合约。

- 示例:通过生成普通账户
```shell
# 生成普通用户私钥user.key
 ./cmc key gen -a ECC_P256 -p ./ -n user.key
```

### 链管理员账户的增删

目前公钥模式只支持批量更新管理员，新管理员列表会替换旧管理员列表，请谨慎操作。

```shell
# 生成管理员私钥
 ./cmc key gen -a ECC_P256 -p ./ -n admin.key

# 导出管理员公钥
 ./cmc key export_pub -k admin.key -n admin.pem
 
# 更新管理员账户
./cmc client chainconfig trustroot update \
--sdk-conf-path=./testdata/sdk_config_pk.yml \
--user-signkey-file-path=./testdata/crypto-config/node1/admin/admin1/admin1.key \
--admin-key-file-paths=./testdata/crypto-config/node1/admin/admin1/admin1.key,./testdata/crypto-config/node1/admin/admin2/admin2.key,./testdata/crypto-config/node1/admin/admin3/admin3.key \
--trust-root-org-id=public \
--trust-root-path=./admin.pem,./testdata/crypto-config/node1/admin/admin1/admin1.pem,./testdata/crypto-config/node1/admin/admin2/admin2.pem,./testdata/crypto-config/node1/admin/admin3/admin3.pem

# 查询管理员列表
```shell
 ./cmc client chainconfig query \
--sdk-conf-path=./testdata/sdk_config_pk.yml 

 "trust_roots": [
    {
      "org_id": "public",
      "root": [
        "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEi95yJxLXrKEeBi5ZJqjk2lEFMKfM\n4pydPq78oTbnHQgQc47eTUENVxBIAEI/mAKjsK82i32amXG0Q9dyqZUWRw==\n-----END PUBLIC KEY-----\n",
        "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEfVzW4O+RjSi+0mPl7HE80LfDup+E\n3s1mziNwP/d5r6X5D5pSdtcGhR80+9rOnIaayM2Eb61m147K72HmgH0I5A==\n-----END PUBLIC KEY-----\n",
        "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE3cFf3ISXtD+vyc6LjuohlHX8A4yG\nIHlMpbwB+H1411TYCgutRoyjXbUy9kcJrXySLS7UCb+/c/yNZ+tz0a6dmA==\n-----END PUBLIC KEY-----\n",
        "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEPuPoOk1otXmNnY0nF0B/eBhIMhEC\niE19OneK0AA4nsk6lsef5PoOG8rI5EljCNIxJQ4pOthTMhX6B0gVjlWoZw==\n-----END PUBLIC KEY-----\n"
      ]
    }
  ],

```

## 链权限管理

### 权限定义

长安链采用三段式语法定义资源的访问权限：规则 (`rule`)、组织列表 (`orgList`)、角色列表 (`roleList`)

- 规则：以关键字的形式描述了需要多少个组织的用户共同认可才可访问资源，合法的规则包括：
    - `ALL`：要求 `orgList` 列表中所有组织参与，每个组织至少提供一个符合 `roleList` 要求角色的签名；
    - `ANY`：要求 `orgList` 列表中任意一个组织提供符合 `roleList` 要求角色的签名；
    - `MAJORITY`：要求联盟链中过半数组织提供各自 `admin` 角色的签名；
    - 一个以字符串形式表达的**整数** (e.g. "3")：要求`orgList` 列表中大于或等于规定数目的组织提供符合 `roleList` 要求角色的签名；
    - 一个以字符串形式表达的**分数** (e.g. "2/3") ：要求`orgList` 列表中大于或等于规定比例的组织提供符合 `roleList` 要求角色的签名；
    - `SELF`：要求资源所属的组织提供符合 `roleList` 要求角色的签名，在此关键字下，`orgList`中的组织列表信息不生效，该规则目前只适用于修改组织根证书、修改组织共识节点地址这两个操作的权限配置；
    - `FORBIDDEN`：此规则表示禁止所有人访问，在此关键字下，`orgList`和 `roleList` 不生效。
- 组织列表：合法的组织列表集合，组织需出现在配置文件的 `trust root` 中，若为空则默认出现在 `trust root` 中的所有组织；
- 角色列表：合法的角色列表集合，若为空则默认所有角色。

示例如下：

| 权限定义                                     | 说明                                                         |
| -------------------------------------------- | ------------------------------------------------------------ |
| `ALL` `[org1, org2, org3]` `[admin, client]` | 三个组织各自提供至少一个管理员或普通用户提供签名才可访问对应资源 |
| `1/2` `[] ` `[admin]`                        | 链上所有组织中过半数组织的管理员提供签名才可访问对应资源（自定义版本的`MAJORITY`规则） |
| `SELF` `[] ` `[admin]`                       | 资源所属组织的管理员提供签名才可访问对应资源，例如组织管理员有权修改各自组织的根证书 |

### 支持的共识算法

Public模式下，目前支持的共识算法及他们对应的权限策略见下：
* [DPOS](#DPOS)
* [TBFT](#TBFT)

### 权限管理

Public模式下，由于权限开放性很高，而且用户没有组织属性，因此不支持通过配置进行自定义权限策略。在该模式下不再需要的系统合约都会被禁止，而针对不同的共识算法，权限控制策略也会有所差别。
#### 交易类型

**任意普通用户都可以调用的交易类型**：

|    交易类型     |   功能   |
| :-------------: | :------: |
| INVOKE_CONTRACT | 合约调用 |
| QUERY_CONTRACT  | 合约查询 |
|    SUBSCRIBE    |   订阅   |

**需要任何一个链管理员签名才可以调用的交易类型**：

| 交易类型 |   功能   |
| :------: | :------: |
| ARCHIVE  | 数据归档 |


<span id="DPOS"></span>

#### DPOS共识

大部分系统合约被禁止使用，仅支持如下系统合约

**需要任何一个链管理员签名才可以调用的系统合约**：

`CHAIN_CONFIG`：链配置更新系统合约

| 方法         | 功能             |
| ------------ | ---------------- |
| CORE_UPDATE  | 核心引擎配置更新 |
| BLOCK_UPDATE | 区块设置配置更新 |

`CONTRACT_MANAGE`：合约管理系统合约

| 方法              | 功能     |
| ----------------- | -------- |
| UPGRADE_CONTRACT  | 升级合约 |
| FREEZE_CONTRACT   | 冻结合约 |
| UNFREEZE_CONTRACT | 解冻合约 |
| REVOKE_CONTRACT   | 吊销合约 |

**需要半数以上链管理员签名才可以调用的系统合约**：

`CHAIN_CONFIG`：链配置更新系统合约

| 方法              | 功能         |
| ----------------- | ------------ |
| TRUST_ROOT_UPDATE | 链管理员更新 |

**任意普通用户都可以调用的系统合约**：

`CONTRACT_MANAGE`：合约管理系统合约

| 方法          | 功能       |
| ------------- | ---------- |
| INIT_CONTRACT | 初始化合约 |

`DPOS_ERC20`：DPOS ERC20系统合约所有方法。

`DPOS_STAKE`：DPOS STAKE系统合约所有方法

Public-DPOS模式默认权限列表如下：

| 合约名             | 方法名                          | 资源名                                          | 功能描述                          | 默认权限                  | 权限描述        |
|-----------------|------------------------------|----------------------------------------------|-------------------------------|-----------------------|-------------|
|                 |                              | ARCHIVE                                      | 归档                            | {[ADMIN] ANY []}      | 任一管理员签名     |
|                 |                              | INVOKE_CONTRACT                              | 执行合约                          | {[] ANY []}           | 无限制         |
|                 |                              | SUBSCRIBE                                    | 订阅                            | {[] ANY []}           | 无限制         |
|                 |                              | QUERY_CONTRACT                               | 查询合约                          | {[] ANY []}           | 无限制         |
| ACCOUNT_MANAGER | REFUND_GAS_VM                | ACCOUNT_MANAGER-REFUND_GAS_VM                | /                             | {[] FORBIDDEN []}     | 禁止          |
| ACCOUNT_MANAGER | SET_ADMIN                    | ACCOUNT_MANAGER-SET_ADMIN                    | /                             | {[] FORBIDDEN []}     | 禁止          |
| ACCOUNT_MANAGER | SET_CONTRACT_METHOD_PAYER    | ACCOUNT_MANAGER-SET_CONTRACT_METHOD_PAYER    | /                             | {[] FORBIDDEN []}     | 禁止          |
| ACCOUNT_MANAGER | CHARGE_GAS                   | ACCOUNT_MANAGER-CHARGE_GAS                   | 收取gas费用                       | {[] FORBIDDEN []}     | 禁止          |
| ACCOUNT_MANAGER | CHARGE_GAS_FOR_MULTI_ACCOUNT | ACCOUNT_MANAGER-CHARGE_GAS_FOR_MULTI_ACCOUNT | 收取多个帐户gas费用                   | {[CONSENSUS] ANY []}  | 任意节点可以操作    |
| CERT_MANAGE     | CERTS_ALIAS_DELETE           | CERT_MANAGE-CERTS_ALIAS_DELETE               | /                             | {[] FORBIDDEN []}     | 禁止          |
| CERT_MANAGE     | CERTS_DELETE                 | CERT_MANAGE-CERTS_DELETE                     | /                             | {[] FORBIDDEN []}     | 禁止          |
| CERT_MANAGE     | CERTS_FREEZE                 | CERT_MANAGE-CERTS_FREEZE                     | /                             | {[] FORBIDDEN []}     | 禁止          |
| CERT_MANAGE     | CERTS_QUERY                  | CERT_MANAGE-CERTS_QUERY                      | 查询证书                          | {[] FORBIDDEN []}     | 禁止          |
| CERT_MANAGE     | CERTS_REVOKE                 | CERT_MANAGE-CERTS_REVOKE                     | /                             | {[] FORBIDDEN []}     | 禁止          |
| CERT_MANAGE     | CERTS_UNFREEZE               | CERT_MANAGE-CERTS_UNFREEZE                   | /                             | {[] FORBIDDEN []}     | 禁止          |
| CERT_MANAGE     | CERT_ADD                     | CERT_MANAGE-CERT_ADD                         | /                             | {[] FORBIDDEN []}     | 禁止          |
| CERT_MANAGE     | CERT_ALIAS_ADD               | CERT_MANAGE-CERT_ALIAS_ADD                   | /                             | {[] FORBIDDEN []}     | 禁止          |
| CERT_MANAGE     | CERT_ALIAS_UPDATE            | CERT_MANAGE-CERT_ALIAS_UPDATE                | /                             | {[] FORBIDDEN []}     | 禁止          |
| CHAIN_CONFIG    | UPDATE_VERSION               | CHAIN_CONFIG-UPDATE_VERSION                  | 更新链版本                         | {[ADMIN] MAJORITY []} | 半数以上管理员多签   |
| CHAIN_CONFIG    | BLOCK_UPDATE                 | CHAIN_CONFIG-BLOCK_UPDATE                    | 更新区块配置                        | {[ADMIN] ANY []}      | 任一管理员签名     |
| CHAIN_CONFIG    | CONSENSUS_EXT_ADD            | CHAIN_CONFIG-CONSENSUS_EXT_ADD               | /                             | {[] FORBIDDEN []}     | 禁止          |
| CHAIN_CONFIG    | CONSENSUS_EXT_DELETE         | CHAIN_CONFIG-CONSENSUS_EXT_DELETE            | /                             | {[] FORBIDDEN []}     | 禁止          |
| CHAIN_CONFIG    | CONSENSUS_EXT_UPDATE         | CHAIN_CONFIG-CONSENSUS_EXT_UPDATE            | /                             | {[] FORBIDDEN []}     | 禁止          |
| CHAIN_CONFIG    | CORE_UPDATE                  | CHAIN_CONFIG-CORE_UPDATE                     | 更新核心模块配置                      | {[ADMIN] ANY []}      | 任一管理员签名     |
| CHAIN_CONFIG    | ENABLE_OR_DISABLE_GAS        | CHAIN_CONFIG-ENABLE_OR_DISABLE_GAS           | /                             | {[] FORBIDDEN []}     | 禁止          |
| CHAIN_CONFIG    | NODE_ID_ADD                  | CHAIN_CONFIG-NODE_ID_ADD                     | /                             | {[] FORBIDDEN []}     | 禁止          |
| CHAIN_CONFIG    | NODE_ID_DELETE               | CHAIN_CONFIG-NODE_ID_DELETE                  | /                             | {[] FORBIDDEN []}     | 禁止          |
| CHAIN_CONFIG    | NODE_ID_UPDATE               | CHAIN_CONFIG-NODE_ID_UPDATE                  | /                             | {[] FORBIDDEN []}     | 禁止          |
| CHAIN_CONFIG    | NODE_ORG_ADD                 | CHAIN_CONFIG-NODE_ORG_ADD                    | /                             | {[] FORBIDDEN []}     | 禁止          |
| CHAIN_CONFIG    | NODE_ORG_DELETE              | CHAIN_CONFIG-NODE_ORG_DELETE                 | /                             | {[] FORBIDDEN []}     | 禁止          |
| CHAIN_CONFIG    | NODE_ORG_UPDATE              | CHAIN_CONFIG-NODE_ORG_UPDATE                 | /                             | {[] FORBIDDEN []}     | 禁止          |
| CHAIN_CONFIG    | PERMISSION_ADD               | CHAIN_CONFIG-PERMISSION_ADD                  | 添加权限                          | {[ADMIN] MAJORITY []} | 半数以上管理员多签   |
| CHAIN_CONFIG    | PERMISSION_DELETE            | CHAIN_CONFIG-PERMISSION_DELETE               | 删除权限（恢复默认）                    | {[ADMIN] MAJORITY []} | 半数以上管理员多签   |
| CHAIN_CONFIG    | PERMISSION_UPDATE            | CHAIN_CONFIG-PERMISSION_UPDATE               | 更新权限                          | {[ADMIN] MAJORITY []} | 半数以上管理员多签   |
| CHAIN_CONFIG    | SET_INVOKE_BASE_GAS          | CHAIN_CONFIG-SET_INVOKE_BASE_GAS             | 设置基础扣费的Gas大小（单次调用的最少扣除的Gas数量） | {[ADMIN] MAJORITY []} | 半数以上组织管理员多签 |
| CHAIN_CONFIG    | SET_INVOKE_GAS_PRICE         | CHAIN_CONFIG-SET_INVOKE_GAS_PRICE            | 设置调用 Gas 价格                   | {[ADMIN] MAJORITY []} | 半数以上管理员多签   |
| CHAIN_CONFIG    | TRUST_MEMBER_ADD             | CHAIN_CONFIG-TRUST_MEMBER_ADD                | /                             | {[] FORBIDDEN []}     | 禁止          |
| CHAIN_CONFIG    | TRUST_MEMBER_DELETE          | CHAIN_CONFIG-TRUST_MEMBER_DELETE             | /                             | {[] FORBIDDEN []}     | 禁止          |
| CHAIN_CONFIG    | TRUST_MEMBER_UPDATE          | CHAIN_CONFIG-TRUST_MEMBER_UPDATE             | /                             | {[] FORBIDDEN []}     | 禁止          |
| CHAIN_CONFIG    | TRUST_ROOT_ADD               | CHAIN_CONFIG-TRUST_ROOT_ADD                  | 添加管理员公钥                       | {[ADMIN] MAJORITY []} | 半数以上管理员多签   |
| CHAIN_CONFIG    | TRUST_ROOT_DELETE            | CHAIN_CONFIG-TRUST_ROOT_DELETE               | 删除管理员公钥                       | {[ADMIN] MAJORITY []} | 半数以上管理员多签   |
| CHAIN_CONFIG    | TRUST_ROOT_UPDATE            | CHAIN_CONFIG-TRUST_ROOT_UPDATE               | 更新管理员公钥                       | {[ADMIN] MAJORITY []} | 半数以上管理员多签   |
| CHAIN_CONFIG    | SET_INSTALL_BASE_GAS         | CHAIN_CONFIG-SET_INSTALL_BASE_GAS            | 设置安装/升级合约花费gas                | {[ADMIN] MAJORITY []} | 半数以上组织管理员多签 |
| CHAIN_CONFIG    | SET_INSTALL_GAS_PRICE        | CHAIN_CONFIG-SET_INSTALL_GAS_PRICE           | 设置安装/升级合约花费gas/byte           | {[ADMIN] MAJORITY []} | 半数以上组织管理员多签 |
| CHAIN_CONFIG    | ENABLE_ONLY_CREATOR_UPGRADE  | CHAIN_CONFIG-ENABLE_ONLY_CREATOR_UPGRADE     | 开启只允许创建者升级合约功能                | {[ADMIN] MAJORITY []} | 半数以上组织管理员多签 |
| CHAIN_CONFIG    | DISABLE_ONLY_CREATOR_UPGRADE | CHAIN_CONFIG-DISABLE_ONLY_CREATOR_UPGRADE    | 关闭只允许创建者升级合约功能                | {[ADMIN] MAJORITY []} | 半数以上组织管理员多签 |
| CHAIN_CONFIG    | MULTI_SIGN_ENABLE_MANUAL_RUN | CHAIN_CONFIG-MULTI_SIGN_ENABLE_MANUAL_RUN    | 启用发送者的多重签名执行合约                | {[] FORBIDDEN []}     | 禁止          |
| CONTRACT_MANAGE | FREEZE_CONTRACT              | CONTRACT_MANAGE-FREEZE_CONTRACT              | 冻结合约                          | {[ADMIN] ANY []}      | 任一管理员签名     |
| CONTRACT_MANAGE | GET_DISABLED_CONTRACT_LIST   | CONTRACT_MANAGE-GET_DISABLED_CONTRACT_LIST   | 获取冻结合约列表                      | {[] ANY []}           | 无限制         |
| CONTRACT_MANAGE | GRANT_CONTRACT_ACCESS        | CONTRACT_MANAGE-GRANT_CONTRACT_ACCESS        | /                             | {[] FORBIDDEN []}     | 禁止          |
| CONTRACT_MANAGE | REVOKE_CONTRACT              | CONTRACT_MANAGE-REVOKE_CONTRACT              | 吊销合约                          | {[ADMIN] ANY []}      | 任一管理员签名     |
| CONTRACT_MANAGE | REVOKE_CONTRACT_ACCESS       | CONTRACT_MANAGE-REVOKE_CONTRACT_ACCESS       | /                             | {[] FORBIDDEN []}     | 禁止          |
| CONTRACT_MANAGE | UNFREEZE_CONTRACT            | CONTRACT_MANAGE-UNFREEZE_CONTRACT            | 解冻合约                          | {[ADMIN] ANY []}      | 任一管理员签名     |
| CONTRACT_MANAGE | UPGRADE_CONTRACT             | CONTRACT_MANAGE-UPGRADE_CONTRACT             | 升级合约                          | {[ADMIN] ANY []}      | 任一管理员签名     |
| CONTRACT_MANAGE | VERIFY_CONTRACT_ACCESS       | CONTRACT_MANAGE-VERIFY_CONTRACT_ACCESS       | /                             | {[] FORBIDDEN []}     | 禁止          |
| MULTI_SIGN      | QUERY                        | MULTI_SIGN-QUERY                             | /                             | {[] FORBIDDEN []}     | 禁止          |
| MULTI_SIGN      | REQ                          | MULTI_SIGN-REQ                               | /                             | {[] FORBIDDEN []}     | 禁止          |
| MULTI_SIGN      | VOTE                         | MULTI_SIGN-VOTE                              | /                             | {[] FORBIDDEN []}     | 禁止          |
| PRIVATE_COMPUTE | SAVE_CA_CERT                 | PRIVATE_COMPUTE-SAVE_CA_CERT                 | /                             | {[] FORBIDDEN []}     | 禁止          |
| PRIVATE_COMPUTE | SAVE_ENCLAVE_REPORT          | PRIVATE_COMPUTE-SAVE_ENCLAVE_REPORT          | /                             | {[] FORBIDDEN []}     | 禁止          |
| PUBKEY_MANAGE   | PUBKEY_ADD                   | PUBKEY_MANAGE-PUBKEY_ADD                     | /                             | {[] FORBIDDEN []}     | 禁止          |
| PUBKEY_MANAGE   | PUBKEY_DELETE                | PUBKEY_MANAGE-PUBKEY_DELETE                  | /                             | {[] FORBIDDEN []}     | 禁止          |

```
  V2.3.3版本去掉ALTER_ADDR_TYPE（修改地址类型）功能，地址类型需要在链初始化时确定，确定后不支持修改。
```
<span id="TBFT"></span>

#### TBFT共识

大部分系统合约被禁止使用，仅支持如下系统合约

**需要任何一个链管理员签名才可以调用的系统合约**：

`CONTRACT_MANAGE`：合约管理系统合约

|       方法        |   功能   |
| :---------------: | :------: |
|   INIT_CONTRACT   | 创建合约 |
| UPGRADE_CONTRACT  | 升级合约 |
|  FREEZE_CONTRACT  | 冻结合约 |
| UNFREEZE_CONTRACT | 解冻合约 |
|  REVOKE_CONTRACT  | 吊销合约 |

**需要半数以上链管理员签名才可以调用的系统合约**：

`CHAIN_CONFIG`：链配置更新系统合约

|         方法          |       功能       |
| :-------------------: | :--------------: |
|      CORE_UPDATE      | 核心引擎配置更新 |
|     BLOCK_UPDATE      | 区块设置配置更新 |
|      NODE_ID_ADD      |   共识节点增加   |
|    NODE_ID_DELETE     |   共识节点删除   |
|    NODE_ID_UPDATE     |   共识节点更新   |
|    NODE_ORG_UPDATE    | 共识节点org更新  |
| ENABLE_OR_DISABLE_GAS |   GAS开关设置    |
|   TRUST_ROOT_UPDATE   |   链管理员更新   |


`ACCOUNT_MANAGER`：审计管理系统合约

|   方法    |     功能      |
| :-------: | :-----------: |
| SET_ADMIN | GAS管理员设置 |


Public-TBFT模式默认权限列表如下：

| 合约名             | 方法名                          | 资源名                                          | 功能描述                          | 默认权限                              | 权限描述                            |
|-----------------|------------------------------|----------------------------------------------|-------------------------------|-----------------------------------|---------------------------------|
|                 |                              | ARCHIVE                                      | 归档                            | {[ADMIN] ANY []}                  | 任一管理员签名                         |
|                 |                              | INVOKE_CONTRACT                              | 执行合约                          | {[] ANY []}                       | 无限制                             |
|                 |                              | SUBSCRIBE                                    | 订阅                            | {[] ANY []}                       | 无限制                             |
|                 |                              | QUERY_CONTRACT                               | 查询合约                          | {[] ANY []}                       | 无限制                             |
| ACCOUNT_MANAGER | REFUND_GAS_VM                | ACCOUNT_MANAGER-REFUND_GAS_VM                | /                             | {[] FORBIDDEN []}                 | 禁止                              |
| ACCOUNT_MANAGER | SET_ADMIN                    | ACCOUNT_MANAGER-SET_ADMIN                    | 设置管理员地址                       | {[ADMIN] MAJORITY []}             | 半数以上管理员多签                       |
| ACCOUNT_MANAGER | SET_CONTRACT_METHOD_PAYER    | ACCOUNT_MANAGER-SET_CONTRACT_METHOD_PAYER    | 为合约方法设置代付款账户                  | {[CONSENSUS CLIENT ADMIN] ANY []} | 任意CLIENT、ADMIN、ConsensusNode可操作 |
| ACCOUNT_MANAGER | CHARGE_GAS                   | ACCOUNT_MANAGER-CHARGE_GAS                   | 收取gas费用                       | {[] FORBIDDEN []}                 | 禁止                              |
| ACCOUNT_MANAGER | CHARGE_GAS_FOR_MULTI_ACCOUNT | ACCOUNT_MANAGER-CHARGE_GAS_FOR_MULTI_ACCOUNT | 收取多个帐户gas费用                   | {[CONSENSUS] ANY []}              | 任意节点可以操作                        |
| CERT_MANAGE     | CERTS_ALIAS_DELETE           | CERT_MANAGE-CERTS_ALIAS_DELETE               | /                             | {[] FORBIDDEN []}                 | 禁止                              |
| CERT_MANAGE     | CERTS_DELETE                 | CERT_MANAGE-CERTS_DELETE                     | /                             | {[] FORBIDDEN []}                 | 禁止                              |
| CERT_MANAGE     | CERTS_FREEZE                 | CERT_MANAGE-CERTS_FREEZE                     | /                             | {[] FORBIDDEN []}                 | 禁止                              |
| CERT_MANAGE     | CERTS_QUERY                  | CERT_MANAGE-CERTS_QUERY                      | 查询证书                          | {[] FORBIDDEN []}                 | 禁止                              |
| CERT_MANAGE     | CERTS_REVOKE                 | CERT_MANAGE-CERTS_REVOKE                     | /                             | {[] FORBIDDEN []}                 | 禁止                              |
| CERT_MANAGE     | CERTS_UNFREEZE               | CERT_MANAGE-CERTS_UNFREEZE                   | /                             | {[] FORBIDDEN []}                 | 禁止                              |
| CERT_MANAGE     | CERT_ADD                     | CERT_MANAGE-CERT_ADD                         | /                             | {[] FORBIDDEN []}                 | 禁止                              |
| CERT_MANAGE     | CERT_ALIAS_ADD               | CERT_MANAGE-CERT_ALIAS_ADD                   | /                             | {[] FORBIDDEN []}                 | 禁止                              |
| CERT_MANAGE     | CERT_ALIAS_UPDATE            | CERT_MANAGE-CERT_ALIAS_UPDATE                | /                             | {[] FORBIDDEN []}                 | 禁止                              |
| CHAIN_CONFIG    | UPDATE_VERSION               | CHAIN_CONFIG-UPDATE_VERSION                  | 更新链版本                         | {[ADMIN] MAJORITY []}             | 半数以上管理员多签                       |
| CHAIN_CONFIG    | BLOCK_UPDATE                 | CHAIN_CONFIG-BLOCK_UPDATE                    | 更新区块配置                        | {[ADMIN] MAJORITY []}             | 半数以上管理员多签                       |
| CHAIN_CONFIG    | CONSENSUS_EXT_ADD            | CHAIN_CONFIG-CONSENSUS_EXT_ADD               | 添加共识扩展参数                      | {[ADMIN] MAJORITY []}             | 半数以上管理员多签                       |
| CHAIN_CONFIG    | CONSENSUS_EXT_DELETE         | CHAIN_CONFIG-CONSENSUS_EXT_DELETE            | 删除共识扩展参数                      | {[ADMIN] MAJORITY []}             | 半数以上管理员多签                       |
| CHAIN_CONFIG    | CONSENSUS_EXT_UPDATE         | CHAIN_CONFIG-CONSENSUS_EXT_UPDATE            | 更新共识扩展参数                      | {[ADMIN] MAJORITY []}             | 半数以上管理员多签                       |
| CHAIN_CONFIG    | CORE_UPDATE                  | CHAIN_CONFIG-CORE_UPDATE                     | 更新核心模块配置                      | {[ADMIN] MAJORITY []}             | 半数以上管理员多签                       |
| CHAIN_CONFIG    | ENABLE_OR_DISABLE_GAS        | CHAIN_CONFIG-ENABLE_OR_DISABLE_GAS           | 是否开启gas                       | {[ADMIN] MAJORITY []}             | 半数以上管理员多签                       |
| CHAIN_CONFIG    | NODE_ID_ADD                  | CHAIN_CONFIG-NODE_ID_ADD                     | 添加节点ID                        | {[ADMIN] MAJORITY []}             | 半数以上管理员多签                       |
| CHAIN_CONFIG    | NODE_ID_DELETE               | CHAIN_CONFIG-NODE_ID_DELETE                  | 删除节点ID                        | {[ADMIN] MAJORITY []}             | 半数以上管理员多签                       |
| CHAIN_CONFIG    | NODE_ID_UPDATE               | CHAIN_CONFIG-NODE_ID_UPDATE                  | 更新节点ID                        | {[ADMIN] MAJORITY []}             | 半数以上管理员多签                       |
| CHAIN_CONFIG    | NODE_ORG_ADD                 | CHAIN_CONFIG-NODE_ORG_ADD                    | /                             | {[] FORBIDDEN []}                 | 禁止                              |
| CHAIN_CONFIG    | NODE_ORG_DELETE              | CHAIN_CONFIG-NODE_ORG_DELETE                 | /                             | {[] FORBIDDEN []}                 | 禁止                              |
| CHAIN_CONFIG    | NODE_ORG_UPDATE              | CHAIN_CONFIG-NODE_ORG_UPDATE                 | 更新节点ID列表                      | {[ADMIN] MAJORITY []}             | 半数以上管理员多签                       |
| CHAIN_CONFIG    | PERMISSION_ADD               | CHAIN_CONFIG-PERMISSION_ADD                  | 添加权限                          | {[ADMIN] MAJORITY []}             | 半数以上管理员多签                       |
| CHAIN_CONFIG    | PERMISSION_DELETE            | CHAIN_CONFIG-PERMISSION_DELETE               | 删除权限（恢复默认）                    | {[ADMIN] MAJORITY []}             | 半数以上管理员多签                       |
| CHAIN_CONFIG    | PERMISSION_UPDATE            | CHAIN_CONFIG-PERMISSION_UPDATE               | 更新权限                          | {[ADMIN] MAJORITY []}             | 半数以上管理员多签                       |                           |
| CHAIN_CONFIG    | SET_ACCOUNT_MANAGER_ADMIN    | CHAIN_CONFIG-SET_ACCOUNT_MANAGER_ADMIN       | 设置管理员地址                       | {[ADMIN] MAJORITY []}             | 半数以上管理员多签                       |
| CHAIN_CONFIG    | SET_INVOKE_BASE_GAS          | CHAIN_CONFIG-SET_INVOKE_BASE_GAS             | 设置基础扣费的Gas大小（单次调用的最少扣除的Gas数量） | {[ADMIN] MAJORITY []}             | 半数以上组织管理员多签                     |
| CHAIN_CONFIG    | SET_INVOKE_GAS_PRICE         | CHAIN_CONFIG-SET_INVOKE_GAS_PRICE            | 设置调用 Gas 价格                   | {[ADMIN] MAJORITY []}             | 半数以上管理员多签                       |
| CHAIN_CONFIG    | TRUST_MEMBER_ADD             | CHAIN_CONFIG-TRUST_MEMBER_ADD                | /                             | {[] FORBIDDEN []}                 | 禁止                              |
| CHAIN_CONFIG    | TRUST_MEMBER_DELETE          | CHAIN_CONFIG-TRUST_MEMBER_DELETE             | /                             | {[] FORBIDDEN []}                 | 禁止                              |
| CHAIN_CONFIG    | TRUST_MEMBER_UPDATE          | CHAIN_CONFIG-TRUST_MEMBER_UPDATE             | /                             | {[] FORBIDDEN []}                 | 禁止                              |
| CHAIN_CONFIG    | TRUST_ROOT_ADD               | CHAIN_CONFIG-TRUST_ROOT_ADD                  | 添加管理员公钥                       | {[ADMIN] MAJORITY []}             | 半数以上管理员多签                       |
| CHAIN_CONFIG    | TRUST_ROOT_DELETE            | CHAIN_CONFIG-TRUST_ROOT_DELETE               | 删除管理员公钥                       | {[ADMIN] MAJORITY []}             | 半数以上管理员多签                       |
| CHAIN_CONFIG    | TRUST_ROOT_UPDATE            | CHAIN_CONFIG-TRUST_ROOT_UPDATE               | 更新管理员公钥                       | {[ADMIN] MAJORITY []}             | 半数以上管理员多签                       |
| CHAIN_CONFIG    | SET_INSTALL_BASE_GAS         | CHAIN_CONFIG-SET_INSTALL_BASE_GAS            | 设置安装/升级合约花费gas                | {[ADMIN] MAJORITY []}             | 半数以上组织管理员多签                     |
| CHAIN_CONFIG    | SET_INSTALL_GAS_PRICE        | CHAIN_CONFIG-SET_INSTALL_GAS_PRICE           | 设置安装/升级合约花费gas/byte           | {[ADMIN] MAJORITY []}             | 半数以上组织管理员多签                     |
| CHAIN_CONFIG    | ENABLE_ONLY_CREATOR_UPGRADE  | CHAIN_CONFIG-ENABLE_ONLY_CREATOR_UPGRADE     | 开启只允许创建者升级合约功能                | {[ADMIN] MAJORITY []}             | 半数以上组织管理员多签                     |
| CHAIN_CONFIG    | DISABLE_ONLY_CREATOR_UPGRADE | CHAIN_CONFIG-DISABLE_ONLY_CREATOR_UPGRADE    | 关闭只允许创建者升级合约功能                | {[ADMIN] MAJORITY []}             | 半数以上组织管理员多签                     |
| CHAIN_CONFIG    | MULTI_SIGN_ENABLE_MANUAL_RUN | CHAIN_CONFIG-MULTI_SIGN_ENABLE_MANUAL_RUN    | 启用发送者的多重签名执行合约                | {[ADMIN] MAJORITY []}             | 半数以上管理员多签                       |
| CONTRACT_MANAGE | FREEZE_CONTRACT              | CONTRACT_MANAGE-FREEZE_CONTRACT              | 冻结合约                          | {[ADMIN] ANY []}                  | 任一管理员签名                         |
| CONTRACT_MANAGE | GET_DISABLED_CONTRACT_LIST   | CONTRACT_MANAGE-GET_DISABLED_CONTRACT_LIST   | 获取冻结合约列表                      | {[] ANY []}                       | 无限制                             |
| CONTRACT_MANAGE | GRANT_CONTRACT_ACCESS        | CONTRACT_MANAGE-GRANT_CONTRACT_ACCESS        | /                             | {[] FORBIDDEN []}                 | 禁止                              |
| CONTRACT_MANAGE | INIT_CONTRACT                | CONTRACT_MANAGE-INIT_CONTRACT                | 安装合约                          | {[ADMIN] ANY []}                  | 任一管理员签名                         |
| CONTRACT_MANAGE | REVOKE_CONTRACT              | CONTRACT_MANAGE-REVOKE_CONTRACT              | 吊销合约                          | {[ADMIN] ANY []}                  | 任一管理员签名                         |
| CONTRACT_MANAGE | REVOKE_CONTRACT_ACCESS       | CONTRACT_MANAGE-REVOKE_CONTRACT_ACCESS       | /                             | {[] FORBIDDEN []}                 | 禁止                              |
| CONTRACT_MANAGE | UNFREEZE_CONTRACT            | CONTRACT_MANAGE-UNFREEZE_CONTRACT            | 解冻合约                          | {[ADMIN] ANY []}                  | 任一管理员签名                         |
| CONTRACT_MANAGE | UPGRADE_CONTRACT             | CONTRACT_MANAGE-UPGRADE_CONTRACT             | 升级合约                          | {[ADMIN] ANY []}                  | 任一管理员签名                         |
| CONTRACT_MANAGE | VERIFY_CONTRACT_ACCESS       | CONTRACT_MANAGE-VERIFY_CONTRACT_ACCESS       | /                             | {[] FORBIDDEN []}                 | 禁止                              |
| PRIVATE_COMPUTE | SAVE_CA_CERT                 | PRIVATE_COMPUTE-SAVE_CA_CERT                 | /                             | {[] FORBIDDEN []}                 | 禁止                              |
| PRIVATE_COMPUTE | SAVE_ENCLAVE_REPORT          | PRIVATE_COMPUTE-SAVE_ENCLAVE_REPORT          | /                             | {[] FORBIDDEN []}                 | 禁止                              |
| PUBKEY_MANAGE   | PUBKEY_ADD                   | PUBKEY_MANAGE-PUBKEY_ADD                     | /                             | {[] FORBIDDEN []}                 | 禁止                              |
| PUBKEY_MANAGE   | PUBKEY_DELETE                | PUBKEY_MANAGE-PUBKEY_DELETE                  | /                             | {[] FORBIDDEN []}                 | 禁止                              |

```
  V2.3.3版本去掉ALTER_ADDR_TYPE（修改地址类型）功能，地址类型需要在链初始化时确定，确定后不支持修改。
```


