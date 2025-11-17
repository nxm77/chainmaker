# 搭建IBC模式账户体系

## 搭建长安链网络

本章节我们将详细介绍如何搭建完整的长安链PermissionWithIBC模式的账户体系，以及如何管理该模式下的链账户，包括组织账户、节点账户、用户账户的增删等。关于长安链 IBC 身份体系的技术细节，请参阅[IBC 技术文档](../tech/IBC技术文档.md)

本文示例说明：

- 创建新链：4组织，1共识节点/每组织，1管理员/每组织
- 链管理：组织管理、节点管理、用户管理

### 源码下载

- 下载ibc版本的`chainmaker-go`源码到本地

```bash
$ git clone -b v2.3.0_ibc --depth=1 https://git.chainmaker.org.cn/chainmaker/chainmaker-go.git
```

- 下载ibc版本的`证书生成工具`源码到本地

```bash
$ git clone -b v2.3.0_ibc  --depth=1 https://git.chainmaker.org.cn/chainmaker/chainmaker-cryptogen.git
```

### 源码编译

- 编译证书生成工具

```bash
$ cd chainmaker-cryptogen
$ make
```

### 生成配置文件

- 将编译好的`chainmaker-cryptogen`，软连接到`chainmaker-go/tools`目录

```bash
# 进入工具目录
$ cd chainmaker-go/tools

# 软连接chainmaker-cryptogen到tools目录下
$ ln -s ../../chainmaker-cryptogen/ .
```

进入`chainmaker-go/scripts`目录，执行`prepare_ibc.sh`脚本生成单链4节点集群配置，存于路径`chainmaker-go/build`中

```bash
# 进入脚本目录
$ cd ../scripts

# 生成单链4节点集群的证书和配置
$ ./prepare_ibc.sh 4 1
begin check params...
begin generate certs, cnt: 4
input consensus type (1-TBFT(default),3-MAXBFT,4-RAFT):
input log level (DEBUG|INFO(default)|WARN|ERROR):
enable vm go (YES|NO(default))
config node total        4
begin generate node1 config...
begin node1 chain1 cert config...
begin node1 trust config...
begin generate node2 config...
begin node2 chain1 cert config...
begin node2 trust config...
begin generate node3 config...
begin node3 chain1 cert config...
begin node3 trust config...
begin generate node4 config...
begin node4 chain1 cert config...
begin node4 trust config...
```

### 编译及安装包制作

- 生成证书（prepare_ibc.sh脚本）后执行`build_release.sh`脚本，将编译`chainmaker-go`模块，并打包生成安装，存于路径`chainmaker-go/build/release`中

```bash
$ ./build_release.sh
$ tree ../build/release/
../build/release/
├── chainmaker-v2.3.0-wx-org1.chainmaker.org-20221228155928-x86_64.tar.gz
├── chainmaker-v2.3.0-wx-org2.chainmaker.org-20221228155928-x86_64.tar.gz
├── chainmaker-v2.3.0-wx-org3.chainmaker.org-20221228155928-x86_64.tar.gz
├── chainmaker-v2.3.0-wx-org4.chainmaker.org-20221228155928-x86_64.tar.gz
└── crypto-config-20221228155928.tar.gz
```

### 启动节点集群

- 执行`cluster_quick_start.sh`脚本，会解压各个安装包，调用`bin`目录中的`start.sh`脚本，启动`chainmaker`节点

```bash
$ ./cluster_quick_start.sh normal
```

> 若需要关闭集群，使用脚本：
>
> ```bash
> $ ./cluster_quick_stop.sh
> ```

### 查看节点启动是否正常

- 查看进程是否存在

```bash
$ ps -ef|grep chainmaker | grep -v grep
25261  2146  4 19:55 pts/20   00:00:01 ./chainmaker start -c ../config/wx-org1.chainmaker.org/chainmaker.yml
25286  2146  4 19:55 pts/20   00:00:01 ./chainmaker start -c ../config/wx-org2.chainmaker.org/chainmaker.yml
25309  2146  4 19:55 pts/20   00:00:01 ./chainmaker start -c ../config/wx-org3.chainmaker.org/chainmaker.yml
25335  2146  4 19:55 pts/20   00:00:01 ./chainmaker start -c ../config/wx-org4.chainmaker.org/chainmaker.yml
```

- 查看端口是否监听

```bash
$ netstat -lptn | grep 1230
tcp6       0      0 :::12301                :::*                    LISTEN      25261/./chainmaker  
tcp6       0      0 :::12302                :::*                    LISTEN      25286/./chainmaker  
tcp6       0      0 :::12303                :::*                    LISTEN      25309/./chainmaker  
tcp6       0      0 :::12304                :::*                    LISTEN      25335/./chainmaker 
```

- 检查节点是否有`ERROR`日志

```bash
$ cat ../build/release/*/bin/panic.log
$ cat ../build/release/*/log/system.log
$ cat ../build/release/*/log/system.log |grep "ERROR\|put block\|all necessary"
//若看到all necessary peers connected则表示节点已经准备就绪。
```

### 部署/调用智能合约

链部署后，在进行部署/调用合约测试，以验证链是否正常运行。


### 编译

cmc工具的编译&运行方式如下：

> 创建工作目录 $WORKDIR 比如 ~/chainmaker<br>
> 启动测试链 [在工作目录下 使用脚本搭建](../quickstart/通过命令行体验链.md)<br>

```bash
# 编译cmc
$ cd $WORKDIR/chainmaker-go/tools/cmc
$ go build
# 配置测试数据
$ cp -rf ../../build/crypto-config ../../tools/cmc/testdata/ # 使用chainmaker-cryptogen生成的测试链的证书
# 查看help
$ cd ../../chainmaker-go/tools/cmc
$ ./cmc --help
```

#### 部署示例合约

- 创建wasm合约

```bash
$ ./cmc client contract user create \
--contract-name=fact \
--runtime-type=WASMER \
--byte-code-path=./testdata/claim-wasm-demo/rust-fact-2.0.0.wasm \
--version=1.0 \
--sdk-conf-path=./testdata/sdk_config_ibc.yml \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--sync-result=true \
--params="{}"
```

- 调用wasm合约

```bash
$ ./cmc client contract user invoke \
--contract-name=fact \
--method=save \
--sdk-conf-path=./testdata/sdk_config_ibc.yml \
--params="{\"file_name\":\"name007\",\"file_hash\":\"ab3456df5799b87c77e7f88\",\"time\":\"6543234\"}" \
--sync-result=true
```

- 查询合约

```bash
$ ./cmc client contract user get \
--contract-name=fact \
--method=find_by_file_hash \
--sdk-conf-path=./testdata/sdk_config_ibc.yml \
--params="{\"file_hash\":\"ab3456df5799b87c77e7f88\"}"
```

## 链管理

在长安链中，不同的账户一般绑定不同的角色，具有不同的权限。 为了提高安全性，长安链默认设置了许多权限，部分操作需要多个管理员多签才能完成。

### 组织管理

#### 共识节点Org管理

- 添加共识节点Org

```
./cmc client chainconfig consensusnodeorg add \
--sdk-conf-path=./testdata/sdk_config_ibc.yml \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--node-ids=QmcQHCuAXaFkbcsPUj7e37hXXfZ9DdN7bozseo5oX4qiC4,QmaWrR72CbT51nFVpNDS8NaqUZjVuD4Ezf8xcHcFW9SJWF \
--node-org-id=wx-org4.chainmaker.org
```

- 删除共识节点Org

```
./cmc client chainconfig consensusnodeorg remove \
--sdk-conf-path=./testdata/sdk_config_ibc.yml \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--node-org-id=wx-org4.chainmaker.org
```

- 更新共识节点Org

```
./cmc client chainconfig consensusnodeorg update \
--sdk-conf-path=./testdata/sdk_config_ibc.yml \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--node-ids=QmcQHCuAXaFkbcsPUj7e37hXXfZ9DdN7bozseo5oX4qiC4 \
--node-org-id=wx-org4.chainmaker.org
```

#### 组织主公钥管理

- 增加组织主公钥

可以用 common 包工具生成新的主公钥。我们假设生成好的主公钥路径为chainmaker-go/tools/cmc/master_public.key

```
./cmc client chainconfig masterkey add \
--sdk-conf-path=./testdata/sdk_config_ibc.yml \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--master-key-org-id=wx-org5.chainmaker.org \
--master-key-path=./master_public.key
```

- 删除组织主公钥

```
./cmc client chainconfig masterkey remove \
--sdk-conf-path=./testdata/sdk_config_ibc.yml \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--master-key-org-id=wx-org5.chainmaker.org
```

- 更新组织主公钥

```
./cmc client chainconfig masterkey update \
--sdk-conf-path=./testdata/sdk_config_ibc.yml \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org4.chainmaker.org/user/admin1/admin1.sign.crt \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org4.chainmaker.org/user/admin1/admin1.sign.key \
--master-key-org-id=wx-org4.chainmaker.org \
--master-key-path=./master_public.key
```

#### 组织信息查询

组织信息可以通过链配置查询，在`master_keys`字段下会返回当前区块链网络中的组织列表以及各个组织下的组织主密钥，同时 在`consensus`下，会返回各个组织下的共识节点列表，命令如下：

```
./cmc client chainconfig query \
--sdk-conf-path=./testdata/sdk_config_ibc.yml
```

### 节点管理

- 增加共识节点

```
./cmc client chainconfig consensusnodeid add \
--sdk-conf-path=./testdata/sdk_config_ibc.yml \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--node-id=QmcQHCuAXaFkbcsPUj7e37hXXfZ9DdN7bozseo5oX4qiC4 \
--node-org-id=wx-org4.chainmaker.org
```

- 更新共识节点Id

```
./cmc client chainconfig consensusnodeid update \
--sdk-conf-path=./testdata/sdk_config_ibc.yml \
--org-id=wx-org1.chainmaker.org \
--user-tlscrt-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.tls.crt \
--user-tlskey-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.tls.key \
--user-signcrt-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.sign.crt \
--user-signkey-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.sign.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--node-id=QmXxeLkNTcvySPKMkv3FUqQgpVZ3t85KMo5E4cmcmrexrC \
--node-id-old=QmcQHCuAXaFkbcsPUj7e37hXXfZ9DdN7bozseo5oX4qiC4 \
--node-org-id=wx-org5.chainmaker.org
```

- 删除共识节点

```
./cmc client chainconfig consensusnodeid remove \
--sdk-conf-path=./testdata/sdk_config_ibc.yml \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--node-id=QmZJuj1SyGKyU5FffzVHBwoz1QN8StjaBhpAEYFwK1wYEu \
--node-org-id=wx-org4.chainmaker.org
```

- 查询共识节点
  共识节点可以通过链配置查询，在`consensus`字段下会返回当前区块链网络中的共识节点列表，命令如下：

```
./cmc client chainconfig query \
--sdk-conf-path=./testdata/sdk_config_ibc.yml
```

### 权限管理

长安链目前支持资源级别的权限管理，可以通过cmc命令行工具或者sdk来查询、新增、修改以及删除资源权限。

- 权限列表查询

```
./cmc client chainconfig permission list \
--sdk-conf-path=./testdata/sdk_config_ibc.yml
```

- 设置账户权限
  权限修改相关的操作一般需要**多数管理员多签**授权，假如我们有个资源名叫：TEST_SUM, 需要设置为”任一用户可以访问”， 使用cmc命令设置权限权限如下：

```
./cmc client chainconfig permission add \
--sdk-conf-path=./testdata/sdk_config_ibc.yml \
--admin-org-ids=wx-org1.chainmaker.org,wx-org2.chainmaker.org,wx-org3.chainmaker.org \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--sync-result=true \
--permission-resource-name="TEST_SUM" \
--permission-resource-policy-rule=ANY \
--permission-resource-policy-roleList=CLIENT
```

- 修改账户权限
  使用cmc命令修改TEST_SUM资源权限为多数管理员多签操作

```
./cmc client chainconfig permission update \
--sdk-conf-path=./testdata/sdk_config_ibc.yml \
--admin-org-ids=wx-org1.chainmaker.org,wx-org2.chainmaker.org,wx-org3.chainmaker.org \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--sync-result=true \
--permission-resource-name="TEST_SUM" \
--permission-resource-policy-rule=MAJORITY \
--permission-resource-policy-roleList=ADMIN
```

- 删除账户权限
  使用cmc命令删除TEST_SUM资源的权限限制，删除后节点权限校验模块不会对该资源进行权限检查

```
./cmc client chainconfig permission delete \
--sdk-conf-path=./testdata/sdk_config_ibc.yml \
--admin-org-ids=wx-org1.chainmaker.org,wx-org2.chainmaker.org,wx-org3.chainmaker.org \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--sync-result=true \
--permission-resource-name="TEST_SUM"
```
