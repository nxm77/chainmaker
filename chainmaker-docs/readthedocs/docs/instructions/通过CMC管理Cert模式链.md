# 管理Cert模式链

在长安链中，不同的账户一般绑定不同的角色，具有不同的权限。 为了提高安全性，长安链默认设置了许多权限，部分操作需要多个管理员多签才能完成。 不同组织的管理员证书账户生成，参考[生成admin测试账户](基于CA服务搭建长安链.html#admin)
## 组织管理
组织管理包含共识节点组织管理和组织根证书管理, 以org5组织管理为例：
#### 新增一个组织到区块链网络
```shell
./cmc client chainconfig consensusnodeorg add \
--sdk-conf-path=./testdata/sdk_config.yml \
--org-id=wx-org1.chainmaker.org \
--user-tlscrt-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.tls.crt \
--user-tlskey-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.tls.key \
--user-signcrt-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.sign.crt \
--user-signkey-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.sign.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--node-ids=QmcQHCuAXaFkbcsPUj7e37hXXfZ9DdN7bozseo5oX4qiC4 \
--node-org-id=wx-org5.chainmaker.org
```
注：新增组织时，需要为组织至少指定一个共识节点

#### 将组织移除区块链网络
```shell
./cmc client chainconfig consensusnodeorg remove \
--sdk-conf-path=./testdata/sdk_config.yml \
--org-id=wx-org1.chainmaker.org \
--user-tlscrt-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.tls.crt \
--user-tlskey-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.tls.key \
--user-signcrt-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.sign.crt \
--user-signkey-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.sign.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--node-org-id=wx-org5.chainmaker.org
```
注：组织被移除后，该组织下的**所有共识节点**会自动降级为**同步节点**，只能同步账本数据，不能参与共识。

#### 更新区块链网络到组织信息
```shell
./cmc client chainconfig consensusnodeorg update \
--sdk-conf-path=./testdata/sdk_config.yml \
--org-id=wx-org1.chainmaker.org \
--user-tlscrt-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.tls.crt \
--user-tlskey-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.tls.key \
--user-signcrt-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.sign.crt \
--user-signkey-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.sign.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--node-ids=QmcQHCuAXaFkbcsPUj7e37hXXfZ9DdN7bozseo5oX4qiC4,QmaWrR72CbT51nFVpNDS8NaqUZjVuD4Ezf8xcHcFW9SJWF \
--node-org-id=wx-org5.chainmaker.org
```
注：可以通过更新组织的方式，可以批量更新共识节点，已有节点会被覆盖，请谨慎操作。单个共识节点的增删，请参考：[节点管理](#节点管理)

#### 组织根证书管理
- 申请组织根证书
```shell
# 申请组织org5根证书
$ curl --location --request POST 'http://localhost:8096/api/ca/gencert' \
--header 'Content-Type: application/json' \
--data-raw '{
    "orgId": "wx-org5.chainmaker.org",
    "userId": "",
    "userType": "ca",
    "certUsage": "sign",
    "country": "CN",
    "locality": "BeiJing",
    "province": "BeiJing"
}' | jq
#response:
{
    "code": 200,
    "msg": "The request service returned successfully",
    "data": {
        "certSn": 4526784566558937862,
        "issueCertSn": 3682789430329867930,
        "cert": "-----BEGIN CERTIFICATE-----\nMIICUDCCAfagAwIBAgIIPtJelFB3UwYwCgYIKoZIzj0EAwIwYjELMAkGA1UEBhMC\nQ04xEDAOBgNVBAgTB0JlaWppbmcxEDAOBgNVBAcTB0JlaWppbmcxETAPBgNVBAoT\nCG9yZy1yb290MQ0wCwYDVQQLEwRyb290MQ0wCwYDVQQDEwRyb290MB4XDTIyMDgy\nNDA3MjMzMVoXDTIzMDIyMDA3MjMzMVowbDELMAkGA1UEBhMCQ04xEDAOBgNVBAgT\nB0JlaUppbmcxEDAOBgNVBAcTB0JlaUppbmcxHzAdBgNVBAoTFnd4LW9yZzUuY2hh\naW5tYWtlci5vcmcxCzAJBgNVBAsTAmNhMQswCQYDVQQDEwJjYTBZMBMGByqGSM49\nAgEGCCqGSM49AwEHA0IABHQI7/AWTxbh9nU/x33I/WhU/6YxZm6nvASZQmAycx8q\nT2kQjLXUFYjl3i4Ku/B0X0IK7OP6NrqgWlDz82VPNoejgYswgYgwDgYDVR0PAQH/\nBAQDAgEGMA8GA1UdEwEB/wQFMAMBAf8wKQYDVR0OBCIEIKsQDmt6hb17p46/Mybw\nS+afVxm8Gk/SHdkglefR/XYIMCsGA1UdIwQkMCKAID2QQznj2uOhDBY0zQDZs0Ps\nixvqNqTx5WLTT7JErCJMMA0GA1UdEQQGMASCAIIAMAoGCCqGSM49BAMCA0gAMEUC\nIDIzqzySsueRMYzgOpvI1SkqzwPtLSgmt2yiEV6QwNOxAiEA2yhlEmr07t6YezwQ\nRBTkHdZEOZjygjUmieU77XpWgTI=\n-----END CERTIFICATE-----\n",
        "privateKey": "-----BEGIN EC PRIVATE KEY-----\nMHcCAQEEIGsZnFxjdglimqodGaaQu+nLbFnzr0KO/gjomPwCaoyloAoGCCqGSM49\nAwEHoUQDQgAEdAjv8BZPFuH2dT/Hfcj9aFT/pjFmbqe8BJlCYDJzHypPaRCMtdQV\niOXeLgq78HRfQgrs4/o2uqBaUPPzZU82hw==\n-----END EC PRIVATE KEY-----\n"
    }
}

# 保存根证书和私钥
$ echo -e "-----BEGIN CERTIFICATE-----\nMIICUDCCAfagAwIBAgIIPtJelFB3UwYwCgYIKoZIzj0EAwIwYjELMAkGA1UEBhMC\nQ04xEDAOBgNVBAgTB0JlaWppbmcxEDAOBgNVBAcTB0JlaWppbmcxETAPBgNVBAoT\nCG9yZy1yb290MQ0wCwYDVQQLEwRyb290MQ0wCwYDVQQDEwRyb290MB4XDTIyMDgy\nNDA3MjMzMVoXDTIzMDIyMDA3MjMzMVowbDELMAkGA1UEBhMCQ04xEDAOBgNVBAgT\nB0JlaUppbmcxEDAOBgNVBAcTB0JlaUppbmcxHzAdBgNVBAoTFnd4LW9yZzUuY2hh\naW5tYWtlci5vcmcxCzAJBgNVBAsTAmNhMQswCQYDVQQDEwJjYTBZMBMGByqGSM49\nAgEGCCqGSM49AwEHA0IABHQI7/AWTxbh9nU/x33I/WhU/6YxZm6nvASZQmAycx8q\nT2kQjLXUFYjl3i4Ku/B0X0IK7OP6NrqgWlDz82VPNoejgYswgYgwDgYDVR0PAQH/\nBAQDAgEGMA8GA1UdEwEB/wQFMAMBAf8wKQYDVR0OBCIEIKsQDmt6hb17p46/Mybw\nS+afVxm8Gk/SHdkglefR/XYIMCsGA1UdIwQkMCKAID2QQznj2uOhDBY0zQDZs0Ps\nixvqNqTx5WLTT7JErCJMMA0GA1UdEQQGMASCAIIAMAoGCCqGSM49BAMCA0gAMEUC\nIDIzqzySsueRMYzgOpvI1SkqzwPtLSgmt2yiEV6QwNOxAiEA2yhlEmr07t6YezwQ\nRBTkHdZEOZjygjUmieU77XpWgTI=\n-----END CERTIFICATE-----\n" > ca.crt
$ echo -e "-----BEGIN EC PRIVATE KEY-----\nMHcCAQEEIGsZnFxjdglimqodGaaQu+nLbFnzr0KO/gjomPwCaoyloAoGCCqGSM49\nAwEHoUQDQgAEdAjv8BZPFuH2dT/Hfcj9aFT/pjFmbqe8BJlCYDJzHypPaRCMtdQV\niOXeLgq78HRfQgrs4/o2uqBaUPPzZU82hw==\n-----END EC PRIVATE KEY-----\n" > ca.key

# 查看org5组织根证书
$ openssl x509 -in ca.crt -text -noout
Certificate:
    Data:
        Version: 3 (0x2)
        Serial Number: 4526784566558937862 (0x3ed25e9450775306)
        Signature Algorithm: ecdsa-with-SHA256
        Issuer: C = CN, ST = Beijing, L = Beijing, O = org-root, OU = root, CN = root
        Validity
            Not Before: Aug 24 07:23:31 2022 GMT
            Not After : Feb 20 07:23:31 2023 GMT
        Subject: C = CN, ST = BeiJing, L = BeiJing, O = wx-org5.chainmaker.org, OU = ca, CN = ca
        Subject Public Key Info:
            Public Key Algorithm: id-ecPublicKey
                Public-Key: (256 bit)
                pub:
                    04:74:08:ef:f0:16:4f:16:e1:f6:75:3f:c7:7d:c8:
                    fd:68:54:ff:a6:31:66:6e:a7:bc:04:99:42:60:32:
                    73:1f:2a:4f:69:10:8c:b5:d4:15:88:e5:de:2e:0a:
                    bb:f0:74:5f:42:0a:ec:e3:fa:36:ba:a0:5a:50:f3:
                    f3:65:4f:36:87
                ASN1 OID: prime256v1
                NIST CURVE: P-256
        X509v3 extensions:
            X509v3 Key Usage: critical
                Certificate Sign, CRL Sign
            X509v3 Basic Constraints: critical
                CA:TRUE
            X509v3 Subject Key Identifier:
                AB:10:0E:6B:7A:85:BD:7B:A7:8E:BF:33:26:F0:4B:E6:9F:57:19:BC:1A:4F:D2:1D:D9:20:95:E7:D1:FD:76:08
            X509v3 Authority Key Identifier:
                keyid:3D:90:43:39:E3:DA:E3:A1:0C:16:34:CD:00:D9:B3:43:EC:8B:1B:EA:36:A4:F1:E5:62:D3:4F:B2:44:AC:22:4C

            X509v3 Subject Alternative Name:
                DNS:, DNS:
    Signature Algorithm: ecdsa-with-SHA256
         30:45:02:20:32:33:ab:3c:92:b2:e7:91:31:8c:e0:3a:9b:c8:
         d5:29:2a:cf:03:ed:2d:28:26:b7:6c:a2:11:5e:90:c0:d3:b1:
         02:21:00:db:28:65:12:6a:f4:ee:de:98:7b:3c:10:44:14:e4:
         1d:d6:44:39:98:f2:82:35:26:89:e5:3b:ed:7a:56:81:32
```
从证书Subject字段内容上我们可以看出，该证书代表组织`wx-org5.chainmaker.org`, 证书类型为`CA`，表示为根证书

- 增加组织根证书
```shell
./cmc client chainconfig trustroot add \
   --sdk-conf-path=./testdata/sdk_config.yml \
   --org-id=wx-org1.chainmaker.org \
   --user-tlscrt-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.tls.crt \
   --user-tlskey-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.tls.key \
   --user-signcrt-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt \
   --user-signkey-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key \
   --admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
   --admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
   --trust-root-org-id=wx-org5.chainmaker.org \  # 组织id
   --trust-root-path=./testdata/crypto-config/wx-org5.chainmaker.org/ca/ca.crt # 组织根证书
```

- 删除组织根证书  
```shell
./cmc client chainconfig trustroot remove \
--sdk-conf-path=./testdata/sdk_config.yml \
--org-id=wx-org1.chainmaker.org \
--user-tlscrt-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.tls.crt \
--user-tlskey-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.tls.key \
--user-signcrt-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt \
--user-signkey-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--trust-root-org-id=wx-org5.chainmaker.org #组织id
```

- 更新组织根证书  
```shell
./cmc client chainconfig trustroot update \
--sdk-conf-path=./testdata/sdk_config.yml \
--org-id=wx-org1.chainmaker.org \
--user-tlscrt-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.tls.crt \
--user-tlskey-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.tls.key \
--user-signcrt-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt \
--user-signkey-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org5.chainmaker.org/user/admin1/admin1.sign.crt \
--admin-key-file-paths=./testdata/crypto-config/wx-org5.chainmaker.org/user/admin1/admin1.sign.key \
--trust-root-org-id=wx-org5.chainmaker.org \ #组织id
--trust-root-path=./testdata/crypto-config/wx-org5.chainmaker.org/ca/ca.crt # 组织根证书
```

#### 组织信息查询
组织信息可以通过链配置查询，在`trust_roots`字段下会返回当前区块链网络中的组织列表以及各个组织下的组织根证书，同时
在`consensus`下，会返回各个组织下的共识节点列表，命令如下：
```shell
./cmc client chainconfig query \
--sdk-conf-path=./testdata/sdk_config.yml
```

## 节点管理
以组织org5增加共识节点为例, 组织的节点证书生成，参考[生成节点证书](基于CA服务搭建长安链.html#生成节点证书)
#### 共识节点管理
- 增加共识节点
```shell
./cmc client chainconfig consensusnodeid add \
--sdk-conf-path=./testdata/sdk_config.yml \
--org-id=wx-org1.chainmaker.org \
--user-tlscrt-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.tls.crt \
--user-tlskey-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.tls.key \
--user-signcrt-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.sign.crt \
--user-signkey-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.sign.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--node-id=QmcQHCuAXaFkbcsPUj7e37hXXfZ9DdN7bozseo5oX4qiC4 \ #这里为新添加节点的nodeId
--node-org-id=wx-org5.chainmaker.org
```

- 更新共识节点Id
```shell
./cmc client chainconfig consensusnodeid update \
--sdk-conf-path=./testdata/sdk_config.yml \
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
```shell
./cmc client chainconfig consensusnodeid remove \
--sdk-conf-path=./testdata/sdk_config.yml \
--org-id=wx-org1.chainmaker.org \
--user-tlscrt-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.tls.crt \
--user-tlskey-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.tls.key \
--user-signcrt-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.sign.crt \
--user-signkey-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.sign.key \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
--node-id=QmcQHCuAXaFkbcsPUj7e37hXXfZ9DdN7bozseo5oX4qiC4 \
--node-org-id=wx-org5.chainmaker.org
```
- 查询共识节点  
共识节点可以通过链配置查询，在`consensus`字段下会返回当前区块链网络中的共识节点列表，命令如下：
```shell
./cmc client chainconfig query \
--sdk-conf-path=./testdata/sdk_config.yml

#...
  "consensus": {
    "nodes": [
      {
        "node_id": [
          "QmSx7PhNBPtSN9R8t5DRgKM2iTZNbpyKtBXWGMH6V2CK53"
        ],
        "org_id": "wx-org1.chainmaker.org"
      },
      {
        "node_id": [
          "QmVviBSVY4xK2161hFWh2v4Wh5ThGAgBiPtXV8XjzVbzPW"
        ],
        "org_id": "wx-org2.chainmaker.org"
      },
      {
        "node_id": [
          "QmbBhed1jeFMkFazYnvVJiqp9RAxnnjA6wxiPVkgAdbeDT"
        ],
        "org_id": "wx-org3.chainmaker.org"
      },
      {
        "node_id": [
          "QmdyLwr6ahCQeSixfw72E17rqn8s4vLewzJgYR64eEeQvD"
        ],
        "org_id": "wx-org4.chainmaker.org"
      }
    ],
    "type": 1
  },
```


#### 同步节点管理
- 添加同步节点  
通过CA服务生成同步节点的账户与共识节点类似，请参考[生成节点证书](基于CA服务搭建长安链.html#生成节点证书)，主要不同点在证书生成过程中`userType`字段
使用`common`  
通过CA服务生成同步节点sign和tls证书如下（示例）：
```shell
# 生成同步节点sign证书
$ curl --location --request POST 'http://localhost:8096/api/ca/gencert' \
--header 'Content-Type: application/json' \
--data-raw '{
"orgId": "wx-org1.chainmaker.org",
"userId": "org1.common1.com",
"userType": "common",
"certUsage": "sign",
"country": "CN",
"locality": "BeiJing",
"province": "BeiJing"
}'

# 生成同步节点tls证书
$ curl --location --request POST 'http://localhost:8096/api/ca/gencert' \
--header 'Content-Type: application/json' \
--data-raw '{
"orgId": "wx-org1.chainmaker.org",
"userId": "org1.common1.com",
"userType": "common",
"certUsage": "tls",
"country": "CN",
"locality": "BeiJing",
"province": "BeiJing"
}'
```
使用上面生成的同步节点sign私钥/证书和tls私钥/证书启动同步节点后，启动方式请参考[基于生成的组织和节点证书创建链](基于CA服务搭建长安链.html#基于生成的组织和节点证书创建链)。
同步节点启动后，会自动同步区块链账本数据；同时同步节点可以通过[节点管理-增加共识节点](#节点管理)的方式切换为共识节点，并参与共识。

- 更新同步节点  
  暂不支持同步节点更新相关操作

- 删除同步节点  
  同步节点的删除可以通过删除同步节点组织根证书（注：会影响组织下的所有同步节点和共识节点，请谨慎操作）； 或者通过[证书管理-冻结证书](#证书管理)方式冻结同步节点证书。
  同步节点删除后，将不能够继续同步账本数据


<span id="addConsensusNodeoperation"></span>
## 实践案例：构造新节点加入网络


假设现在存在一个由4个共识节点组成的链chain1，我们期望构造一个新的节点加入到网络中，以下是这些步骤以及相关命令行示例

### 准备节点运行的目录、配置、和执行文件

通过执行以下几个步骤，可以生成新节点的运行目录，目录结构为：

```sh
chainmaker-go/build/release]# ll
drwxr-xr-x 7 root root     4096 Jun 28 20:14 chainmaker-v2.2.1-wx-org1.chainmaker.org
drwxr-xr-x 7 root root     4096 Jun 28 20:14 chainmaker-v2.2.1-wx-org2.chainmaker.org
drwxr-xr-x 7 root root     4096 Jun 28 20:14 chainmaker-v2.2.1-wx-org3.chainmaker.org
drwxr-xr-x 7 root root     4096 Jun 28 20:14 chainmaker-v2.2.1-wx-org4.chainmaker.org
drwxr-xr-x 7 root root     4096 Jun 28 20:14 chainmaker-v2.2.1-wx-org5.chainmaker.org
chainmaker-go/build/release]# tree -L 3 chainmaker-v2.2.1-wx-org5.chainmaker.org
chainmaker-v2.2.1-wx-org5.chainmaker.org
├── bin
│   ├── chainmaker
│   ├── chainmaker.service
│   ├── docker-vm-standalone-start.sh
│   ├── docker-vm-standalone-stop.sh
│   ├── init.sh
│   ├── panic.log
│   ├── restart.sh
│   ├── run.sh
│   ├── start.sh
│   └── stop.sh
├── config
│   └── wx-org5.chainmaker.org
│       ├── certs
│       ├── chainconfig
│       ├── chainmaker.yml
│       └── log.yml
├── data
│   └── wx-org5.chainmaker.org
│       ├── block
│       ├── history
│       ├── ledgerData1
│       ├── result
│       └── state
├── lib
│   ├── libwasmer.so
│   └── wxdec
└── log
```

具体操作步骤如下：

- 基于目前的链上已存在的组织和节点扩展出新的证书，参考[【扩展证书】](../instructions/证书生成工具.md)，示例脚本和说明如下：

  ```sh
  current_version="v2.2.1"
  chainmaker_path="/home/projects/test/chainmaker-go"
  cryptogen_path="/home/projects/test/chainmaker-cryptogen"
  gen_new_cert(){
  echo "------------在chainmaker-cryptogen下生成org5的节点证书：-------------------"
  cd "${cryptogen_path}/bin"
  # 1.备份旧的证书文件
  rm -rf crypto-config-bk
  mv crypto-config crypto-config-bk
  # 2.复制链上所有节点的证书到bin目录下，适用于在现有组织证书的基础上扩展节点证书
  cp -r "${chainmaker_path}/build/crypto-config" ./
  # 3. 修改配置文件
  # 如果增加一个共识节点，node: -type: consensus 下面的count: 本来是1，改为2就等于增加一个共识节点，适用于在现有组织证书的基础上扩展节点证书
  # 如果想增加一个新的组织：修改host_name: wx-org下的count，适用于增加一个新的组织及其节点
  vim ../config/crypto_config_template.yml
  # 请手动编辑后保存。
  # 4. 执行命令生成新的证书
  ./chainmaker-cryptogen extend -c ../config/crypto_config_template.yml
  }
  ```

  方法执行后会在chainmaker-cryptogen/bin目录下生成新节点的证书：

  ```shell
  chainmaker-cryptogen/bin]# tree -L 1 crypto-config
  crypto-config
  ├── wx-org1.chainmaker.org
  ├── wx-org2.chainmaker.org
  ├── wx-org3.chainmaker.org
  ├── wx-org4.chainmaker.org
  ├── wx-org5.chainmaker.org
  ```

  

- 构建节点的执行目录

  ```sh
  copy_node_file(){
  echo "------------复制org5节点的包文件-------------------"
  cd "${chainmaker_path}/build/release"
  #1. 复制org1的节点目录
  cp -r "chainmaker-$current_version-wx-org1.chainmaker.org" "chainmaker-$current_version-wx-org5.chainmaker.org"
  #2. 把chainmaker-*-wx-org5.chainmaker.org/bin下所有的.sh脚本中所有wx-org1.chainmaker.org替换为wx-org5.chainmaker.org"
  cd chainmaker-*-wx-org5.chainmaker.org/
  sed -i 's/org1/org5/g' bin/start.sh
  sed -i 's/org1/org5/g' bin/stop.sh
  sed -i 's/org1/org5/g' bin/restart.sh
  sed -i 's/org1/org5/g' bin/run.sh
  #3.删除data和log
  rm -rf data/*
  rm -rf log/*
  # 4. config/下的org1目录重命名为org5
  mv config/wx-org1.chainmaker.org config/wx-org5.chainmaker.org
  }
  ```

  以上方法执行后，会在build/release下生成新的节点目录：

  ```sh
  chainmaker-go/build/release]# ll
  drwxr-xr-x 7 root root     4096 Jun 28 20:14 chainmaker-v2.2.1-wx-org1.chainmaker.org
  drwxr-xr-x 7 root root     4096 Jun 28 20:14 chainmaker-v2.2.1-wx-org2.chainmaker.org
  drwxr-xr-x 7 root root     4096 Jun 28 20:14 chainmaker-v2.2.1-wx-org3.chainmaker.org
  drwxr-xr-x 7 root root     4096 Jun 28 20:14 chainmaker-v2.2.1-wx-org4.chainmaker.org
  drwxr-xr-x 7 root root     4096 Jun 28 20:14 chainmaker-v2.2.1-wx-org5.chainmaker.org
  ```

  

- 更新节点目录下的证书:

  在第2步中，我们只对节点目录中对应的名称进行了替换，未将org5的证书放在目录 中，因此需要将在第1步生成的org5的证书复制过来

    ```sh
    update_cert(){
    # 使用chainmaker-cryptogen生成的wx-org5.chainmaker.org下的node和user分别覆盖掉wx-org5.chainmaker.org/config/wx-org5.chainmaker.org/certs下的node和user"
    crypto_cert=${cryptogen_path}/bin/crypto-config/wx-org5.chainmaker.org
    chainmaker_cert=${chainmaker_path}/build/release/chainmaker-*-wx-org5.chainmaker.org/config/wx-org5.chainmaker.org/certs
    cp -rf $crypto_cert/user $chainmaker_cert
    cp -rf $crypto_cert/node $chainmaker_cert
    }
    ```


- 修改org5的节点配置文件和链配置文件
  ```sh
  update_config(){
  echo "----------------修改org5的chainmaker.yml----------------"
  # 1. 把wx-org5.chainmaker.org/config/wx-org5.chainmaker.org/chainmaker.yml中所有wx-org1.chainmaker.org替换为wx-org5.chainmaker.org
  sed -i 's/org1/org5/g' config/wx-org5.chainmaker.org/chainmaker.yml
  # 2. 修改net模块，把 listen_addr: /ip4/0.0.0.0/tcp/11301 修改为 listen_addr: /ip4/0.0.0.0/tcp/11305"
  old_config="/ip4/0.0.0.0/tcp/11301"
  new_config="/ip4/0.0.0.0/tcp/11305"
  sed -i "s%${old_config}\+%${new_config}%g" config/wx-org5.chainmaker.org/chainmaker.yml
  sed -i 's/12301/12305/g' config/wx-org5.chainmaker.org/chainmaker.yml
  sed -i 's/14321/14325/g' config/wx-org5.chainmaker.org/chainmaker.yml
  sed -i 's/24321/24325/g' config/wx-org5.chainmaker.org/chainmaker.yml
  # 3. 如果是新增同步节点,将consensus证书改成common证书，如果新增的是共识节点，这一步略过。
  # sed -i 's/consensus1/common1/g' config/wx-org5.chainmaker.org/chainmaker.yml
  echo "----------------修改org5的bc1.yml----------------"
  # 4. 修改chainmaker-*-wx-org5.chainmaker.org/config/wx-org5.chainmaker.org/chainconfig/bc1.yml中的trust_roots模块
  # 把所有 ../config/wx-org1.chainmaker.org 修改为 ../config/wx-org5.chainmaker.org
  sed -i "s/wx-org1.chainmaker.org\/certs\/ca/wx-org5.chainmaker.org\/certs\/ca/g" config/wx-org5.chainmaker.org/chainconfig/bc1.yml
  }
  ```

- 如果使用cmc命令行需要更新相应的证书文件
  
  ```sh
  update_cmc_cert(){
  crypto_cert=${cryptogen_path}/bin/crypto-config
  cmc_cert=${chainmaker_path}/tools/cmc/testdata/
  ls $crypto_cert
  ls $cmc_cert
  cp -rf $crypto_cert $cmc_cert
  }
  ```
  
- 调用以上方法完成节点目录的生成
  
  ```sh
  gen_new_cert
  copy_node_file
  update_cert
  update_config
  update_cmc_cert
  ```
  

### 添加组织根证书 

如果新的节点证书是基于新的组织根证书（例如org5的ca证书）签发的，则需要添加组织根证书。

如果新的节点证书是基于旧的组织根证书（例如org1的ca证书）签发的，则不需要这一步。
[使用cmc添加组织根证书](../dev/命令行工具.html#chainConfig.addOrgRootCA)


  ```shell
  ./cmc client chainconfig trustroot add \
  --sdk-conf-path=./testdata/sdk_config.yml \
  --org-id=wx-org1.chainmaker.org \
  --user-tlscrt-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.tls.crt \
  --user-tlskey-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.tls.key \
  --user-signcrt-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.sign.crt \
  --user-signkey-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.sign.key \
  --admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
  --admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
  --trust-root-org-id=wx-org5.chainmaker.org \
  --sync-result=true \
  --trust-root-path=./testdata/crypto-config/wx-org5.chainmaker.org/ca/ca.crt
  ```

  

### 启动节点

添加trustroot之后，节点就可以做为同步节点加入网络，启动节点使区块同步到最新高度，就可以使用后续的操作步骤加入共识：

```sh
chainmaker-go/build/release/chainmaker-v2.2.1-wx-org5.chainmaker.org/bin]# ./start.sh
```



### 添加共识org：

[使用cmc添加共识节点](../dev/命令行工具.html#chainConfig.addConsensusNodeOrg) （为新加入的组织添加共识节点）
  添加共识org的作用：

  1. 添加参与共识的组织（ 即将org5添加到共识组织列表中）
  2. 添加org5组织下的节点nodeid

如果是添加了新的组织org5，要将组织以及组织下的节点（org5_node1) 添加到共识当中，则需要执行这一步

如果是在原有组织or1-org4的基础上，扩展新的共识节点（例如org1原来有1个共识节点org1_node1，需要加入org1_node2也做为共识节点，则不需要这一步，直接执行下面步骤。

  ```sh
  ./cmc client chainconfig consensusnodeorg add \
  --sdk-conf-path=./testdata/sdk_config.yml \
  --org-id=wx-org1.chainmaker.org \
  --user-tlscrt-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.tls.crt \
  --user-tlskey-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.tls.key \
  --user-signcrt-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.sign.crt \
  --user-signkey-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.sign.key \
  --admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
  --admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
  --node-ids=QmaGii9BS7tWpkXTb9Z5Bbvw43bjt1Jnc4C2U68NnAm28b \
  --sync-result=true \
  --node-org-id=wx-org5.chainmaker.org
  ```

  添加成功后，可查看log/system.log如下关键信息

  ```sh
  cat system.log |grep "validator.go"
  cat system.log |grep "new validator set"
  ```

### 添加共识节点ID 
为已经加入共识的组织，例如or1-org4, 以及执行过和3步的组织，添加更多的共识节点

  ```sh
  ./cmc client chainconfig consensusnodeid add \
  --sdk-conf-path=./testdata/sdk_config.yml \
  --org-id=wx-org1.chainmaker.org \
  --user-tlscrt-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.tls.crt \
  --user-tlskey-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.tls.key \
  --user-signcrt-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.sign.crt \
  --user-signkey-file-path=./testdata/crypto-config/wx-org1.chainmaker.org/user/client1/client1.sign.key \
  --admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
  --admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
  --node-id=QmaGii9BS7tWpkXTb9Z5Bbvw43bjt1Jnc4C2U68NnAm28b \
  --sync-result=true \
  --node-org-id=wx-org5.chainmaker.org
  ```

  添加成功后，可查看log/system.log如下关键信息

  ```sh
  cat system.log*|grep "addedValidators"
  ```



## 用户管理
#### 管理员用户管理
- 新增管理员  
 新增管理员方式，请参考 [生成admin测试账户](基于CA服务搭建长安链.html#admin)章节中管理员账户的生成，管理员证书生成后，
 管理员可以使用CA颁发的admin证书直接访问链上信息。

- 删除用户证书

暂不支持直接删除链上管理员账户，可以通过**证书管理功能**来冻结或吊销链上证书的方式，进而限制用户对链的访问，详情请参考[证书管理](#证书管理)

- 查询管理员列表  
可以通过CA服务提供的证书列表查询方式，间接查询管理员列表。 以查询组织org5下查询管理员为例，`userType`需要指定为`admin`, 如下：

```shell
curl --location --request POST 'http://localhost:8096/api/ca/querycerts' \
--header 'Content-Type: application/json' \
--data-raw '{
    "orgId": "wx-org5.chainmaker.org",
    "userType": "admin",
    "certUsage": "sign"
}

{
    "code": 200,
    "msg": "The request service returned successfully",
    "data": [
        {
            "userId": "admin1",
            "orgId": "wx-org5.chainmaker.org",
            "userType": "admin",
            "certUsage": "sign",
            "certSn": 476309406125711642,
            "issuerSn": 4526784566558937862,
            "certContent": "-----BEGIN CERTIFICATE-----\nMIICTTCCAfSgAwIBAgIIBpww4ZtM4RowCgYIKoZIzj0EAwIwbDELMAkGA1UEBhMC\nQ04xEDAOBgNVBAgTB0JlaUppbmcxEDAOBgNVBAcTB0JlaUppbmcxHzAdBgNVBAoT\nFnd4LW9yZzUuY2hhaW5tYWtlci5vcmcxCzAJBgNVBAsTAmNhMQswCQYDVQQDEwJj\nYTAeFw0yMjA4MzAwMjE4MjhaFw0yMzAyMjYwMjE4MjhaMHMxCzAJBgNVBAYTAkNO\nMRAwDgYDVQQIEwdCZWlKaW5nMRAwDgYDVQQHEwdCZWlKaW5nMR8wHQYDVQQKExZ3\neC1vcmc1LmNoYWlubWFrZXIub3JnMQ4wDAYDVQQLEwVhZG1pbjEPMA0GA1UEAxMG\nYWRtaW4xMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEG+e33zUe7tYz10yVXrTK\n+Vk9cHPRjz65J2q1fWitFdRp63vYlF2hGOGPwohHXMQ6mofNiR9UZqyFGOAjuZKF\n7KN5MHcwDgYDVR0PAQH/BAQDAgbAMCkGA1UdDgQiBCDbNi9Kuq2YskAOwJDIfTFa\n6xdsaD8CV366t1vEtEKWODArBgNVHSMEJDAigCCrEA5reoW9e6eOvzMm8Evmn1cZ\nvBpP0h3ZIJXn0f12CDANBgNVHREEBjAEggCCADAKBggqhkjOPQQDAgNHADBEAiA9\nVoR9iRV1rYwbNRjewMs/NZPrQp6l6gZsAOKch/W0uAIgSuRdqYPikrC+4WgX5HQg\nM6nwJPUxSq2FTUd6SvVtJUY=\n-----END CERTIFICATE-----\n",
            "expirationDate": 1677377908,
            "isRevoked": false
        },
        {
            "userId": "admin2",
            "orgId": "wx-org5.chainmaker.org",
            "userType": "admin",
            "certUsage": "sign",
            "certSn": 2655863384365393409,
            "issuerSn": 4526784566558937862,
            "certContent": "-----BEGIN CERTIFICATE-----\nMIICTzCCAfSgAwIBAgIIJNuFrT0Q5gEwCgYIKoZIzj0EAwIwbDELMAkGA1UEBhMC\nQ04xEDAOBgNVBAgTB0JlaUppbmcxEDAOBgNVBAcTB0JlaUppbmcxHzAdBgNVBAoT\nFnd4LW9yZzUuY2hhaW5tYWtlci5vcmcxCzAJBgNVBAsTAmNhMQswCQYDVQQDEwJj\nYTAeFw0yMjA4MzAwMjE4MzFaFw0yMzAyMjYwMjE4MzFaMHMxCzAJBgNVBAYTAkNO\nMRAwDgYDVQQIEwdCZWlKaW5nMRAwDgYDVQQHEwdCZWlKaW5nMR8wHQYDVQQKExZ3\neC1vcmc1LmNoYWlubWFrZXIub3JnMQ4wDAYDVQQLEwVhZG1pbjEPMA0GA1UEAxMG\nYWRtaW4yMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEiS86NLQero79bxS3Lse7\n46Rp+oZndblFqPGRkuiztcFwEuUyVYwVk/syEHU6yTIY8kaJelhHrq2/KJ1eNp0Q\npaN5MHcwDgYDVR0PAQH/BAQDAgbAMCkGA1UdDgQiBCDA1fzr67GDcwuiaOXzUzuW\np+8ux6DrgAfy5UHzIZyM5jArBgNVHSMEJDAigCCrEA5reoW9e6eOvzMm8Evmn1cZ\nvBpP0h3ZIJXn0f12CDANBgNVHREEBjAEggCCADAKBggqhkjOPQQDAgNJADBGAiEA\nlX9uFrgJlUHk8DwsPzhaZ3E9xOa/H8SugruS153qG3YCIQDUBpMv7HD9H3wVKfU1\nOVpBq2bCdHFSO5yPrYIBy3zF2A==\n-----END CERTIFICATE-----\n",
            "expirationDate": 1677377911,
            "isRevoked": false
        }
    ]
}
````
从查询结果看，组织org5下有两个管理员admin1和admin2

#### 普通用户管理
- 生成用户证书（client）
```shell
$ curl --location --request POST 'http://localhost:8096/api/ca/gencert' \
--header 'Content-Type: application/json' \
--data-raw '{
    "orgId": "wx-org1.chainmaker.org",
    "userId": "client1",
    "userType": "client",
    "certUsage": "sign",
    "country": "CN",
    "locality": "BeiJing",
    "province": "BeiJing"
}' | jq

{
  "code": 200,
  "msg": "The request service returned successfully",
  "data": {
    "certSn": 1877624678000075300,
    "issueCertSn": 1696861425238169000,
    "cert": "-----BEGIN CERTIFICATE-----\nMIICaDCCAg6gAwIBAgIIGg6pslGzoaYwCgYIKoZIzj0EAwIwgYMxCzAJBgNVBAYT\nAkNOMRAwDgYDVQQIEwdCZWlqaW5nMRAwDgYDVQQHEwdCZWlqaW5nMR8wHQYDVQQK\nExZ3eC1vcmcxLmNoYWlubWFrZXIub3JnMQswCQYDVQQLEwJjYTEiMCAGA1UEAxMZ\nY2Etd3gtb3JnMS5jaGFpbm1ha2VyLm9yZzAeFw0yMjA4MjQwNzA2NThaFw0yMzAy\nMjAwNzA2NThaMHUxCzAJBgNVBAYTAkNOMRAwDgYDVQQIEwdCZWlKaW5nMRAwDgYD\nVQQHEwdCZWlKaW5nMR8wHQYDVQQKExZ3eC1vcmcxLmNoYWlubWFrZXIub3JnMQ8w\nDQYDVQQLEwZjbGllbnQxEDAOBgNVBAMTB2NsaWVudDEwWTATBgcqhkjOPQIBBggq\nhkjOPQMBBwNCAASzGFb7efLCJcHbk1SQ9iYWr5gH70O/v5j53mAi6OdYLqgiNzaE\nSudb343MWUPMPKDJcFv6WWvgvuwKulPwt4ljo3kwdzAOBgNVHQ8BAf8EBAMCBsAw\nKQYDVR0OBCIEIPlN3wdwbv+AmN3BPybsylivb4c/UgTeZ//dFTnoejuhMCsGA1Ud\nIwQkMCKAIFQUtaMGHEYWLDGcfWBR0SQZjuWdI31bW7tCAyrplWpuMA0GA1UdEQQG\nMASCAIIAMAoGCCqGSM49BAMCA0gAMEUCIQCdMobpCCJX4hx3f8uNpAsK8EUrmzgU\no1Mup+v9WbpLwAIgFT/Tf4nlhchXMFIZ4W4t9s6Z2NHDDdt5fz+i2CyVOQ0=\n-----END CERTIFICATE-----\n",
    "privateKey": "-----BEGIN EC PRIVATE KEY-----\nMHcCAQEEIAkp0yFG1amB2SSuhdFMslVGOj8BH5eymvvdZOEdEivXoAoGCCqGSM49\nAwEHoUQDQgAEsxhW+3nywiXB25NUkPYmFq+YB+9Dv7+Y+d5gIujnWC6oIjc2hErn\nW9+NzFlDzDygyXBb+llr4L7sCrpT8LeJYw==\n-----END EC PRIVATE KEY-----\n"
  }
}

# 保存用户sign证书和私钥
$ echo -e "-----BEGIN CERTIFICATE-----\nMIICfDCCAiOgAwIBAgIINKV9XwIOzHEwCgYIKoZIzj0EAwIwgYMxCzAJBgNVBAYT\nAkNOMRAwDgYDVQQIEwdCZWlqaW5nMRAwDgYDVQQHEwdCZWlqaW5nMR8wHQYDVQQK\nExZ3eC1vcmcxLmNoYWlubWFrZXIub3JnMQswCQYDVQQLEwJjYTEiMCAGA1UEAxMZ\nY2Etd3gtb3JnMS5jaGFpbm1ha2VyLm9yZzAeFw0yMjA4MjQwNTUzMTRaFw0yMzAy\nMjAwNTUzMTRaMHMxCzAJBgNVBAYTAkNOMRAwDgYDVQQIEwdCZWlKaW5nMRAwDgYD\nVQQHEwdCZWlKaW5nMR8wHQYDVQQKExZ3eC1vcmcxLmNoYWlubWFrZXIub3JnMQ4w\nDAYDVQQLEwVhZG1pbjEPMA0GA1UEAxMGYWRtaW4xMFkwEwYHKoZIzj0CAQYIKoZI\nzj0DAQcDQgAE4lJgyz1A/GledacE3HUCwT2MrWHqDzHHpy81RTY22Ap6uj/ahO1h\n1+o0GSGS8xdF0ygw07BhtBOe9dkuYvptQ6OBjzCBjDAOBgNVHQ8BAf8EBAMCA/gw\nEwYDVR0lBAwwCgYIKwYBBQUHAwIwKQYDVR0OBCIEIJ9OT40+LKdBATiWmfeIXmXO\n2RkxhbG1+Ai5vt+l5LxsMCsGA1UdIwQkMCKAIFQUtaMGHEYWLDGcfWBR0SQZjuWd\nI31bW7tCAyrplWpuMA0GA1UdEQQGMASCAIIAMAoGCCqGSM49BAMCA0cAMEQCIDbU\nvsm/C2OPw5HbhZCBzZJNbJ5x3QW+I8kOKNZd4X7jAiAq12qTcqZcNSBKsqM9HRTj\nqt/2rzh38uEgu/2oUvpIUw==\n-----END CERTIFICATE-----\n" > admin1.tls.crt
$ echo -e "-----BEGIN EC PRIVATE KEY-----\nMHcCAQEEIIA9k6+tvHrCa0Ee3P+IlGO0AQFW3GR6VSYtD33lK+PDoAoGCCqGSM49\nAwEHoUQDQgAE4lJgyz1A/GledacE3HUCwT2MrWHqDzHHpy81RTY22Ap6uj/ahO1h\n1+o0GSGS8xdF0ygw07BhtBOe9dkuYvptQw==\n-----END EC PRIVATE KEY-----\n" > admin1.tls.key
```

client用户tls证书和私钥的生成、测试与[生成admin测试账户](基于CA服务搭建长安链.html#admin)章节中管理员账户的生成和测试类似，这里不在赘述。主要区别：
生成证书接口参数userType为client；client账户链上权限与admin相比权限较低。权限相关介绍参考[权限管理](#权限管理)

- 删除用户证书

暂不支持直接删除链上账户，可以通过**证书管理功能**来冻结或吊销证书的方式，进而限制用户对链的访问，详情请参考[证书管理](#证书管理)

- 查询用户证书  
可以通过CA服务提供的证书列表查询方式，间接查询普通用户列表。 以查询组织org5下查询普通账户为例，`userType`需要指定为`client`，如下：

```shell
$ curl --location --request POST 'http://localhost:8096/api/ca/querycerts' \
--header 'Content-Type: application/json' \
--data-raw '{
    "orgId": "wx-org5.chainmaker.org",
    "userType": "client",
    "certUsage": "sign"
}

{
    "code": 200,
    "msg": "The request service returned successfully",
    "data": [
        {
            "userId": "client1",
            "orgId": "wx-org5.chainmaker.org",
            "userType": "client",
            "certUsage": "sign",
            "certSn": 3649421597944817067,
            "issuerSn": 4526784566558937862,
            "certContent": "-----BEGIN CERTIFICATE-----\nMIICUDCCAfagAwIBAgIIMqVZm50jpaswCgYIKoZIzj0EAwIwbDELMAkGA1UEBhMC\nQ04xEDAOBgNVBAgTB0JlaUppbmcxEDAOBgNVBAcTB0JlaUppbmcxHzAdBgNVBAoT\nFnd4LW9yZzUuY2hhaW5tYWtlci5vcmcxCzAJBgNVBAsTAmNhMQswCQYDVQQDEwJj\nYTAeFw0yMjA4MzAwMjI0MjBaFw0yMzAyMjYwMjI0MjBaMHUxCzAJBgNVBAYTAkNO\nMRAwDgYDVQQIEwdCZWlKaW5nMRAwDgYDVQQHEwdCZWlKaW5nMR8wHQYDVQQKExZ3\neC1vcmc1LmNoYWlubWFrZXIub3JnMQ8wDQYDVQQLEwZjbGllbnQxEDAOBgNVBAMT\nB2NsaWVudDEwWTATBgcqhkjOPQIBBggqhkjOPQMBBwNCAAT8FzELiNwtT4xavmP6\nS3iTNYAKtMIXbEKDglELwGImTz01l3Lc8tZsos550VwKgjmLPAvBgMyMCIrsbIT2\n8bZMo3kwdzAOBgNVHQ8BAf8EBAMCBsAwKQYDVR0OBCIEIK6x0q+mj9DLfYT7uIxl\nowVLNx6xiuPDvG3JGe4jjEzIMCsGA1UdIwQkMCKAIKsQDmt6hb17p46/MybwS+af\nVxm8Gk/SHdkglefR/XYIMA0GA1UdEQQGMASCAIIAMAoGCCqGSM49BAMCA0gAMEUC\nIQCiAU3LEGJcr1pbdlpn+dENvzkqwuwq0BI5PeWfKn/07wIgKGpqzzJLb8qf8IWp\nJPEjn7a1TMWgAPIa6jktnRMf4gc=\n-----END CERTIFICATE-----\n",
            "expirationDate": 1677378260,
            "isRevoked": false
        }
    ]
}
```
从查询结果看，组织org5下有一个普通用户client1


## 权限管理
<span id="权限管理"></span>

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


### 权限资源定义

长安链中，资源名称的定义采用 **[合约名称]-[方法名称]** 的规则。

例如:
* 修改链上配置的系统合约 **CHAIN_CONFIG** ，该合约包含添加根证书的方法 **TRUST_ROOT_ADD**，如果要修改该方法的权限，对应的资源名称为：CHAIN_CONFIG-TRUST_ROOT_ADD。
* CHAIN_CONFIG-TRUST_ROOT_ADD对应的默认权限定义是`{MAJORITY [] [ADMIN]}`则表示**添加根证书**操作需要客户端交易满足**半数以上组织管理员多签**，才能验证通过。

目前长安链PermissionWithCert模式内部默认权限列表如下：

| 合约名             | 方法名                               | 资源名                                           | 功能描述                          | 默认权限                                           | 权限描述                            |
|-----------------|-----------------------------------|-----------------------------------------------|-------------------------------|------------------------------------------------|---------------------------------|
|                 |                                   | ARCHIVE                                       | 归档                            | {[ADMIN] {[ADMIN] ANY []}                      | 任一管理员签名                         |
|                 |                                   | INVOKE_CONTRACT                               | 调用合约                          | {[CONSENSUS COMMON CLIENT ADMIN LIGHT] ANY []} | 任一组织共识节点、同步节点、轻节点、普通用户和管理员签名    |
|                 |                                   | QUERY_CONTRACT                                | 查询合约                          | {[CONSENSUS COMMON CLIENT ADMIN LIGHT] ANY []} | 任一组织共识节点、同步节点、轻节点、普通用户和管理员签名    |
|                 |                                   | SUBSCRIBE                                     | 订阅                            | {[LIGHT CLIENT ADMIN] ANY []}                  | 任一组织轻节点、普通用户和管理员签名              |
| ACCOUNT_MANAGER | SET_ADMIN                         | ACCOUNT_MANAGER-SET_ADMIN                     | 设置管理员地址                       | {[ADMIN] MAJORITY []}                          | 半数以上组织管理员多签                     |
| ACCOUNT_MANAGER | SET_CONTRACT_METHOD_PAYER         | ACCOUNT_MANAGER-SET_CONTRACT_METHOD_PAYER     | 设置gas代扣                       | {[CONSENSUS CLIENT ADMIN] ANY []}              | 任意CLIENT、ADMIN、ConsensusNode可操作 || CERT_MANAGE  | CERTS_ALIAS_DELETE        | CERT_MANAGE-CERTS_ALIAS_DELETE             | 删除证书别名                        | {[ADMIN] ANY []}                               | 任一组织管理员签名                    |
| ACCOUNT_MANAGER | REFUND_GAS_VM                     | ACCOUNT_MANAGER-REFUND_GAS_VM                 | 退gas                          | {[] ANY []}                                    | 无限制                             |
| ACCOUNT_MANAGER | CHARGE_GAS                        | ACCOUNT_MANAGER-CHARGE_GAS                    | 收取gas费用                       | {[] FORBIDDEN []}                              | 禁止                              |
| ACCOUNT_MANAGER | CHARGE_GAS_FOR_MULTI_ACCOUNT      | ACCOUNT_MANAGER-CHARGE_GAS_FOR_MULTI_ACCOUNT  | 统一扣除区块内所有交易的gas               | {[CONSENSUS] ANY []}                           | 任意节点可以操作                        |
| CERT_MANAGE     | CERTS_ALIAS_QUERY                 | CERT_MANAGE-CERTS_ALIAS_QUERY                 | 根据别名查询证书                      | {[CONSENSUS COMMON CLIENT ADMIN LIGHT] ANY []} | 任一共识节点、同步节点、普通用户、管理员和轻节点签名      |
| CERT_MANAGE     | CERTS_DELETE                      | CERT_MANAGE-CERTS_DELETE                      | 删除证书哈希                        | {[ADMIN] ANY []}                               | 任一组织管理员签名                       |
| CERT_MANAGE     | CERTS_FREEZE                      | CERT_MANAGE-CERTS_FREEZE                      | 冻结证书                          | {[ADMIN] ANY []}                               | 任一组织管理员签名                       |
| CERT_MANAGE     | CERTS_QUERY                       | CERT_MANAGE-CERTS_QUERY                       | 查询证书                          | {[CONSENSUS COMMON CLIENT ADMIN LIGHT] ANY []} | 任一共识节点、同步节点、普通用户、管理员和轻节点签名      |
| CERT_MANAGE     | CERTS_REVOKE                      | CERT_MANAGE-CERTS_REVOKE                      | 注销证书                          | {[ADMIN] ANY []}                               | 任一组织管理员签名                       |
| CERT_MANAGE     | CERTS_UNFREEZE                    | CERT_MANAGE-CERTS_UNFREEZE                    | 解冻证书                          | {[ADMIN] ANY []}                               | 任一组织管理员签名                       |
| CERT_MANAGE     | CERT_ADD                          | CERT_MANAGE-CERT_ADD                          | 上链证书哈希                        | {[CLIENT ADMIN LIGHT] ANY []}                  | 任一组织普通用户、管理员和轻节点签名              |
| CERT_MANAGE     | CERT_ALIAS_ADD                    | CERT_MANAGE-CERT_ALIAS_ADD                    | 添加证书别名                        | {[CLIENT ADMIN LIGHT] ANY []}                  | 任一组织普通用户、管理员和轻节点签名              |
| CERT_MANAGE     | CERT_ALIAS_UPDATE                 | CERT_MANAGE-CERT_ALIAS_UPDATE                 | 更新证书别名                        | {[ADMIN] ANY []}                               | 任一组织管理员签名                       |
| CHAIN_CONFIG    | UPDATE_VERSION                    | CHAIN_CONFIG-UPDATE_VERSION                   | 更新链版本                         | {[ADMIN] MAJORITY []}                          | 半数以上管理员多签                       |
| CHAIN_CONFIG    | BLOCK_UPDATE                      | CHAIN_CONFIG-BLOCK_UPDATE                     | 更新出块配置                        | {[ADMIN] MAJORITY []}                          | 半数以上组织管理员多签                     |
| CHAIN_CONFIG    | CONSENSUS_EXT_ADD                 | CHAIN_CONFIG-CONSENSUS_EXT_ADD                | 添加共识扩展字段                      | {[ADMIN] MAJORITY []}                          | 半数以上组织管理员多签                     |
| CHAIN_CONFIG    | CONSENSUS_EXT_DELETE              | CHAIN_CONFIG-CONSENSUS_EXT_DELETE             | 删除共识扩展字段                      | {[ADMIN] MAJORITY []}                          | 半数以上组织管理员多签                     |
| CHAIN_CONFIG    | CONSENSUS_EXT_UPDATE              | CHAIN_CONFIG-CONSENSUS_EXT_UPDATE             | 更新共识扩展字段                      | {[ADMIN] MAJORITY []}                          | 半数以上组织管理员多签                     |
| CHAIN_CONFIG    | CORE_UPDATE                       | CHAIN_CONFIG-CORE_UPDATE                      | 核心模块配置更新                      | {[ADMIN] MAJORITY []}                          | 半数以上组织管理员多签                     |
| CHAIN_CONFIG    | ENABLE_OR_DISABLE_GAS             | CHAIN_CONFIG-ENABLE_OR_DISABLE_GAS            | 是否开启全局gas                     | {[ADMIN] MAJORITY []}                          | 半数以上组织管理员多签                     |
| CHAIN_CONFIG    | GET_CHAIN_CONFIG                  | CHAIN_CONFIG-GET_CHAIN_CONFIG                 | 获取链配置                         | {[CONSENSUS COMMON CLIENT ADMIN LIGHT] ANY []} | 任一组织共识节点、同步节点、普通用户、管理员或轻节点签名    |
| CHAIN_CONFIG    | NODE_ID_ADD                       | CHAIN_CONFIG-NODE_ID_ADD                      | 添加节点ID                        | {[ADMIN] MAJORITY []}                          | 半数以上组织管理员多签                     |
| CHAIN_CONFIG    | NODE_ID_DELETE                    | CHAIN_CONFIG-NODE_ID_DELETE                   | 删除节点ID                        | {[ADMIN] MAJORITY []}                          | 半数以上组织管理员多签                     |
| CHAIN_CONFIG    | NODE_ID_UPDATE                    | CHAIN_CONFIG-NODE_ID_UPDATE                   | 更新节点ID                        | {[ADMIN] SELF []}                              | 本组织管理员签名                        |
| CHAIN_CONFIG    | NODE_ORG_ADD                      | CHAIN_CONFIG-NODE_ORG_ADD                     | 添加共识组织及节点                     | {[ADMIN] MAJORITY []}                          | 半数以上组织管理员多签                     |
| CHAIN_CONFIG    | NODE_ORG_DELETE                   | CHAIN_CONFIG-NODE_ORG_DELETE                  | 删除共识组织                        | {[ADMIN] MAJORITY []}                          | 半数以上组织管理员多签                     |
| CHAIN_CONFIG    | NODE_ORG_UPDATE                   | CHAIN_CONFIG-NODE_ORG_UPDATE                  | 更新共识组织                        | {[ADMIN] MAJORITY []}                          | 半数以上组织管理员多签                     |
| CHAIN_CONFIG    | PERMISSION_ADD                    | CHAIN_CONFIG-PERMISSION_ADD                   | 添加权限                          | {[ADMIN] MAJORITY []}                          | 半数以上组织管理员多签                     |
| CHAIN_CONFIG    | PERMISSION_DELETE                 | CHAIN_CONFIG-PERMISSION_DELETE                | 删除权限（恢复默认）                    | {[ADMIN] MAJORITY []}                          | 半数以上组织管理员多签                     |
| CHAIN_CONFIG    | PERMISSION_UPDATE                 | CHAIN_CONFIG-PERMISSION_UPDATE                | 更新权限                          | {[ADMIN] MAJORITY []}                          | 半数以上组织管理员多签                     |
| CHAIN_CONFIG    | SET_ACCOUNT_MANAGER_ADMIN         | CHAIN_CONFIG-SET_ACCOUNT_MANAGER_ADMIN        | 设置管理员地址                       | {[ADMIN] MAJORITY []}                          | 半数以上组织管理员多签                     |
| CHAIN_CONFIG    | SET_INVOKE_BASE_GAS               | CHAIN_CONFIG-SET_INVOKE_BASE_GAS              | 设置基础扣费的Gas大小（单次调用的最少扣除的Gas数量） | {[ADMIN] MAJORITY []}                          | 半数以上组织管理员多签                     |
| CHAIN_CONFIG    | SET_INVOKE_GAS_PRICE              | CHAIN_CONFIG-SET_INVOKE_GAS_PRICE             | 设置调用 Gas 价格                   | {[ADMIN] MAJORITY []}                          | 半数以上管理员多签                       |
| CHAIN_CONFIG    | TRUST_MEMBER_ADD                  | CHAIN_CONFIG-TRUST_MEMBER_ADD                 | 添加第三方用户                       | {[ADMIN] MAJORITY []}                          | 半数以上组织管理员多签                     |
| CHAIN_CONFIG    | TRUST_MEMBER_DELETE               | CHAIN_CONFIG-TRUST_MEMBER_DELETE              | 删除第三方用户                       | {[ADMIN] MAJORITY []}                          | 半数以上组织管理员多签                     |
| CHAIN_CONFIG    | TRUST_MEMBER_UPDATE               | CHAIN_CONFIG-TRUST_MEMBER_UPDATE              | 更新第三方用户                       | {[ADMIN] MAJORITY []}                          | 半数以上组织管理员多签                     |
| CHAIN_CONFIG    | TRUST_ROOT_ADD                    | CHAIN_CONFIG-TRUST_ROOT_ADD                   | 添加信任根证书                       | {[ADMIN] MAJORITY []}                          | 半数以上组织管理员多签                     |
| CHAIN_CONFIG    | TRUST_ROOT_DELETE                 | CHAIN_CONFIG-TRUST_ROOT_DELETE                | 删除信任根证书                       | {[ADMIN] MAJORITY []}                          | 半数以上组织管理员多签                     |
| CHAIN_CONFIG    | TRUST_ROOT_UPDATE                 | CHAIN_CONFIG-TRUST_ROOT_UPDATE                | 更新信任根证书                       | {[ADMIN] SELF []}                              | 本组织管理员签名                        |
| CHAIN_CONFIG    | SET_INSTALL_BASE_GAS              | CHAIN_CONFIG-SET_INSTALL_BASE_GAS             | 设置安装/升级合约花费gas                | {[ADMIN] MAJORITY []}                          | 半数以上组织管理员多签                     |
| CHAIN_CONFIG    | SET_INSTALL_GAS_PRICE             | CHAIN_CONFIG-SET_INSTALL_GAS_PRICE            | 设置安装/升级合约花费gas/byte           | {[ADMIN] MAJORITY []}                          | 半数以上组织管理员多签                     |
| CHAIN_CONFIG    | MULTI_SIGN_ENABLE_MANUAL_RUN      | CHAIN_CONFIG-MULTI_SIGN_ENABLE_MANUAL_RUN     | 启用发送者的多重签名执行合约                | {[ADMIN] MAJORITY []}                          | 半数以上管理员多签                       |
| CHAIN_QUERY     | GET_ARCHIVED_BLOCK_HEIGHT         | CHAIN_QUERY-GET_ARCHIVED_BLOCK_HEIGHT         | 查询已归档的区块高度                    | {[CONSENSUS COMMON CLIENT ADMIN LIGHT] ANY []} | 任一共识节点、同步节点、普通用户、管理员和轻节点签名      |
| CHAIN_QUERY     | GET_BLOCK_BY_HASH                 | CHAIN_QUERY-GET_BLOCK_BY_HASH                 | 通过哈希查询区块                      | {[CONSENSUS COMMON CLIENT ADMIN LIGHT] ANY []} | 任一共识节点、同步节点、普通用户、管理员和轻节点签名      |
| CHAIN_QUERY     | GET_BLOCK_BY_HEIGHT               | CHAIN_QUERY-GET_BLOCK_BY_HEIGHT               | 通过高度查询区块                      | {[CONSENSUS COMMON CLIENT ADMIN LIGHT] ANY []} | 任一共识节点、同步节点、普通用户、管理员和轻节点签名      |
| CHAIN_QUERY     | GET_BLOCK_BY_TX_ID                | CHAIN_QUERY-GET_BLOCK_BY_TX_ID                | 通过交易ID查询区块                    | {[CONSENSUS COMMON CLIENT ADMIN LIGHT] ANY []} | 任一共识节点、同步节点、普通用户、管理员和轻节点签名      |
| CHAIN_QUERY     | GET_BLOCK_HEADER_BY_HEIGHT        | CHAIN_QUERY-GET_BLOCK_HEADER_BY_HEIGHT        | 通过高度查询区块头                     | {[CONSENSUS COMMON CLIENT ADMIN LIGHT] ANY []} | 任一共识节点、同步节点、普通用户、管理员和轻节点签名      |
| CHAIN_QUERY     | GET_BLOCK_HEIGHT_BY_HASH          | CHAIN_QUERY-GET_BLOCK_HEIGHT_BY_HASH          | 通过哈希查询区块头                     | {[CONSENSUS COMMON CLIENT ADMIN LIGHT] ANY []} | 任一共识节点、同步节点、普通用户、管理员和轻节点签名      |
| CHAIN_QUERY     | GET_BLOCK_HEIGHT_BY_TX_ID         | CHAIN_QUERY-GET_BLOCK_HEIGHT_BY_TX_ID         | 通过交易ID查询区块头                   | {[CONSENSUS COMMON CLIENT ADMIN LIGHT] ANY []} | 任一共识节点、同步节点、普通用户、管理员和轻节点签名      |
| CHAIN_QUERY     | GET_BLOCK_WITH_TXRWSETS_BY_HASH   | CHAIN_QUERY-GET_BLOCK_WITH_TXRWSETS_BY_HASH   | 根据区块哈希查询区块和读写集                | {[CONSENSUS COMMON CLIENT ADMIN LIGHT] ANY []} | 任一共识节点、同步节点、普通用户、管理员和轻节点签名      |
| CHAIN_QUERY     | GET_BLOCK_WITH_TXRWSETS_BY_HEIGHT | CHAIN_QUERY-GET_BLOCK_WITH_TXRWSETS_BY_HEIGHT | 根据区块高度查询区块和读写集                | {[CONSENSUS COMMON CLIENT ADMIN LIGHT] ANY []} | 任一共识节点、同步节点、普通用户、管理员和轻节点签名      |
| CHAIN_QUERY     | GET_FULL_BLOCK_BY_HEIGHT          | CHAIN_QUERY-GET_FULL_BLOCK_BY_HEIGHT          | 根据高度查询区块（区块、读写集、事件）           | {[CONSENSUS COMMON CLIENT ADMIN LIGHT] ANY []} | 任一共识节点、同步节点、普通用户、管理员和轻节点签名      |
| CHAIN_QUERY     | GET_LAST_BLOCK                    | CHAIN_QUERY-GET_LAST_BLOCK                    | 查询最新区块                        | {[CONSENSUS COMMON CLIENT ADMIN LIGHT] ANY []} | 任一共识节点、同步节点、普通用户、管理员和轻节点签名      |
| CHAIN_QUERY     | GET_LAST_CONFIG_BLOCK             | CHAIN_QUERY-GET_LAST_CONFIG_BLOCK             | 查询最新配置区块                      | {[CONSENSUS COMMON CLIENT ADMIN LIGHT] ANY []} | 任一共识节点、同步节点、普通用户、管理员和轻节点签名      |
| CHAIN_QUERY     | GET_TX_BY_TX_ID                   | CHAIN_QUERY-GET_TX_BY_TX_ID                   | 根据交易ID查询交易                    | {[CONSENSUS COMMON CLIENT ADMIN LIGHT] ANY []} | 任一共识节点、同步节点、普通用户、管理员和轻节点签名      |
| CONTRACT_MANAGE | FREEZE_CONTRACT                   | CONTRACT_MANAGE-FREEZE_CONTRACT               | 冻结合约                          | {[ADMIN] MAJORITY []}                          | 半数以上组织管理员多签                     |
| CONTRACT_MANAGE | GET_DISABLED_CONTRACT_LIST        | CONTRACT_MANAGE-GET_DISABLED_CONTRACT_LIST    | 获取冻结合约列表                      | {[] ANY []}                                    | 无限制                             |
| CONTRACT_MANAGE | GRANT_CONTRACT_ACCESS             | CONTRACT_MANAGE-GRANT_CONTRACT_ACCESS         | 授予对原生合约的访问权限                  | { {[ADMIN] MAJORITY []}                        | 半数以上组织管理员多签                     |
| CONTRACT_MANAGE | INIT_CONTRACT                     | CONTRACT_MANAGE-INIT_CONTRACT                 | 安装合约                          | {[ADMIN] ANY []}                               | 任一管理员签名                         |
| CONTRACT_MANAGE | REVOKE_CONTRACT                   | CONTRACT_MANAGE-REVOKE_CONTRACT               | 吊销合约                          | {[ADMIN] MAJORITY []}                          | 半数以上组织管理员多签                     |
| CONTRACT_MANAGE | UNFREEZE_CONTRACT                 | CONTRACT_MANAGE-UNFREEZE_CONTRACT             | 解冻合约                          | {[ADMIN] MAJORITY []}                          | 半数以上组织管理员多签                     |
| CONTRACT_MANAGE | UPGRADE_CONTRACT                  | CONTRACT_MANAGE-UPGRADE_CONTRACT              | 升级合约                          | {[ADMIN] MAJORITY []}                          | 半数以上组织管理员多签                     |
| PRIVATE_COMPUTE | SAVE_CA_CERT                      | PRIVATE_COMPUTE-SAVE_CA_CERT                  | 保存隐私合约的根证书中                   | {[ADMIN] MAJORITY []}                          | 半数以上组织管理员多签                     |
| PRIVATE_COMPUTE | SAVE_ENCLAVE_REPORT               | PRIVATE_COMPUTE-SAVE_ENCLAVE_REPORT           | 保存隐私合约Enclave的Report          | {[ADMIN] MAJORITY []}                          | 半数以上组织管理员多签                     |
| PUBKEY_MANAGE   | PUBKEY_ADD                        | PUBKEY_MANAGE-PUBKEY_ADD                      | 添加公钥                          | {[] FORBIDDEN []}                              | 禁止                              |
| PUBKEY_MANAGE   | PUBKEY_DELETE                     | PUBKEY_MANAGE-PUBKEY_DELETE                   | 删除公钥                          | {[] FORBIDDEN []}                              | 禁止                              |

```
  V2.3.3版本去掉ALTER_ADDR_TYPE（修改地址类型）功能，地址类型需要在链初始化时确定，确定后不支持修改。
```

可以通过修改链上配置中的权限定义部分，来自定义或者修改用户合约和系统合约中的某些方法的权限。 

参考如下示例：

```yaml
  - resource_name: CHAIN_CONFIG-TRUST_ROOT_ADD
    policy:
      rule: MAJORITY
      org_list:
      role_list:
        - admin
```

### 权限查询及修改

长安链目前支持资源级别的权限管理，可以通过cmc命令行工具或者sdk来查询、新增、修改以及删除资源权限。  

长安链证书模式下账户权限与绑定的角色相关，资源的访问权限也与角色相关。 一般情况下修改账户权限需要对资源权限进行修改。

注：长安链支持资源级别的细粒度修改，暂不支持账户级别的细粒度修改。


- 权限列表查询
```shell
  ./cmc client chainconfig permission list \
  --sdk-conf-path=./testdata/sdk_config.yml
```

- 设置账户权限   
权限修改相关的操作一般需要**多数管理员多签**授权，假如我们有个资源名叫：TEST_SUM, 需要设置为"任一用户可以访问"， 使用cmc命令设置权限权限如下：
```shell
  ./cmc client chainconfig permission add \
  --sdk-conf-path=./testdata/sdk_config.yml \
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
```shell
./cmc client chainconfig permission update \
  --sdk-conf-path=./testdata/sdk_config.yml \
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
```shell
./cmc client chainconfig permission delete \
  --sdk-conf-path=./testdata/sdk_config.yml \
  --admin-org-ids=wx-org1.chainmaker.org,wx-org2.chainmaker.org,wx-org3.chainmaker.org \
  --admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.key,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.key \
  --admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org2.chainmaker.org/user/admin1/admin1.sign.crt,./testdata/crypto-config/wx-org3.chainmaker.org/user/admin1/admin1.sign.crt \
  --sync-result=true \
  --permission-resource-name="TEST_SUM"
```

注：除了自定义资源的权限设置外，长安链也支持默认权限的修改，比如升级合约需要**MAJORITY多签**操作，可以通过`修改账户权限`改为**ANY单签**权限。此外，对于默认权限虽然可以进行修改，但是不能删除，删除操作仅删除到系统默认权限。

## 证书管理
### 冻结/解冻证书
长安链证书模式下，支持证书的冻结和解冻操作，在节点或者用户账户不再安全的情况下，比如说密钥丢失或被盗用，可以通过冻结证书的方式禁止该账户继续发起交易。
- 冻结证书  
```shell
$ ./cmc client certmanage freeze  \
--sdk-conf-path=./testdata/sdk_config.yml \
--cert-file-paths=./client1.sign.crt,./client2.sign.crt \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key
```
通过参数`cert-file-paths`指定要冻结的证书文件路径，如果需要冻结多个证书账户，证书文件名之间用`，`隔开。 

- 解冻证书  
证书冻结后，需要通过证书解冻才能继续发起交易，证书解冻方式如下:
```shell
$ ./cmc client certmanage unfreeze  \
--sdk-conf-path=./testdata/sdk_config.yml \
--cert-file-paths=./client1.sign.crt,./client2.sign.crt \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key
```
通过参数`cert-file-paths`指定要解冻的证书文件路径，如果需要解冻多个证书账户，证书文件名之间用`，`隔开。

注：根证书不支持冻结和解冻操作

### 撤销证书
当用户账户不安全时，除了通过证书冻结方式保护链上信息，还可以通过撤销证书的方式。 与证书冻结的方式不同，证书一旦撤销不可恢复，请谨慎操作。证书撤销需要用到`crl证书撤销列表`，
crl一般来自CA服务。
- 生成证书撤销列表crl
```shell
# 撤销证书，用证书sn指定
curl --location --request POST 'http://localhost:8096/api/ca/revokecert' \
--header 'Content-Type: application/json' \
--data-raw '{
    "revokedCertSn": 1877624678000075174,
    "issuerCertSn": 1696861425238169023,
    "reason":"密钥丢失，无法找回"
}' | jq

{
  "code": 200,
  "msg": "The request service returned successfully",
  "data": "-----BEGIN X509 CRL-----\nMIIBXDCCAQECAQEwCgYIKoZIzj0EAwIwgYMxCzAJBgNVBAYTAkNOMRAwDgYDVQQI\nEwdCZWlqaW5nMRAwDgYDVQQHEwdCZWlqaW5nMR8wHQYDVQQKExZ3eC1vcmcxLmNo\nYWlubWFrZXIub3JnMQswCQYDVQQLEwJjYTEiMCAGA1UEAxMZY2Etd3gtb3JnMS5j\naGFpbm1ha2VyLm9yZxcNMjIwODI1MDcxNjA5WhcNMjIwODI2MDcxNjA5WjAbMBkC\nCBoOqbJRs6GmFw0yMzAyMjAwNzA2NThaoC8wLTArBgNVHSMEJDAigCBUFLWjBhxG\nFiwxnH1gUdEkGY7lnSN9W1u7QgMq6ZVqbjAKBggqhkjOPQQDAgNJADBGAiEA+y1U\ntVHVH1vW1L/N9lGSz5RVXs89irzN44RFoXP0VwACIQDrllN4otO4qtPidEBEJZKv\ngkWAUWBVxn91TCM2OxeIKg==\n-----END X509 CRL-----\n"
}

# 将返回的data保存为crl文件
echo -e "-----BEGIN X509 CRL-----\nMIIBWzCCAQECAQEwCgYIKoZIzj0EAwIwgYMxCzAJBgNVBAYTAkNOMRAwDgYDVQQI\nEwdCZWlqaW5nMRAwDgYDVQQHEwdCZWlqaW5nMR8wHQYDVQQKExZ3eC1vcmcxLmNo\nYWlubWFrZXIub3JnMQswCQYDVQQLEwJjYTEiMCAGA1UEAxMZY2Etd3gtb3JnMS5j\naGFpbm1ha2VyLm9yZxcNMjIwODI1MDcxOTEwWhcNMjIwODI2MDcxOTEwWjAbMBkC\nCBoOqbJRs6GmFw0yMzAyMjAwNzA2NThaoC8wLTArBgNVHSMEJDAigCBUFLWjBhxG\nFiwxnH1gUdEkGY7lnSN9W1u7QgMq6ZVqbjAKBggqhkjOPQQDAgNIADBFAiBg6Mdn\nNa2aXRPDJC6/ukDWQnnDGtrDTp09oVVaGV0b4wIhALtTXy9AAlsn7UL6HMP+UHk4\nlH2sGzq3lsYslqhHKXwG\n-----END X509 CRL-----\n" > ./client1.crl


# 如果要查询某CA下的所有撤销证书，可以通过以下方式查询。
 curl --location --request POST 'http://localhost:8096/api/ca/gencrl' \
--header 'Content-Type: application/json' \
--data-raw '{
    "issuerCertSn":1696861425238169023,
    "token":""
}' | jq

{
  "code": 200,
  "msg": "The request service returned successfully",
  "data": "-----BEGIN X509 CRL-----\nMIIBWzCCAQECAQEwCgYIKoZIzj0EAwIwgYMxCzAJBgNVBAYTAkNOMRAwDgYDVQQI\nEwdCZWlqaW5nMRAwDgYDVQQHEwdCZWlqaW5nMR8wHQYDVQQKExZ3eC1vcmcxLmNo\nYWlubWFrZXIub3JnMQswCQYDVQQLEwJjYTEiMCAGA1UEAxMZY2Etd3gtb3JnMS5j\naGFpbm1ha2VyLm9yZxcNMjIwODI1MDcxOTEwWhcNMjIwODI2MDcxOTEwWjAbMBkC\nCBoOqbJRs6GmFw0yMzAyMjAwNzA2NThaoC8wLTArBgNVHSMEJDAigCBUFLWjBhxG\nFiwxnH1gUdEkGY7lnSN9W1u7QgMq6ZVqbjAKBggqhkjOPQQDAgNIADBFAiBg6Mdn\nNa2aXRPDJC6/ukDWQnnDGtrDTp09oVVaGV0b4wIhALtTXy9AAlsn7UL6HMP+UHk4\nlH2sGzq3lsYslqhHKXwG\n-----END X509 CRL-----\n"
}
```

- 链上证书撤销  
```shell
 ./cmc client certmanage revoke \
--sdk-conf-path=./testdata/sdk_config.yml \
--cert-crl-path=./client1.crl \
--admin-crt-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.crt \
--admin-key-file-paths=./testdata/crypto-config/wx-org1.chainmaker.org/user/admin1/admin1.sign.key
```
通过参数`cert-crl-path`指定crl，并发起交易，交易执行成功后，该crl中包含的证书将不能够再发起交易或访问链上信息。

<span id="renewal"></span>
## 如何给快过期的证书续期

本章节将要说明如何给快到期的证书续期、更新链上信息，此处包括：节点证书、trustRoot证书、用户证书

> tips: 操作前，最好先备份data和config目录


### 证书替换流程

1. 续期各个证书：ca、 node、admin、client证书
2. 校验生成的证书正确性
3. 停止当前链
4. 备份data config
5. 替换node、admin、client证书
6. ——————
7. 启动链
8. 发交易删除过期组织1的共识节点ID
9. 发交易删除过期组织1的共识节点
10. 发交易删除过期组织1的trustRoot（ca）
11. 发交易添加续期组织1的trustRoot（ca）
12. 发交易添加续期组织1的共识节点
13. 8-12循环执行直到替换掉所有的组织

#### 用户证书

- 普通用户证书无需特殊操作。直接使用新续期的证书发交易即可。

- 若使用证书别名上链，需要先执行别名更新或者别名删除后重新别名上链。参考：长安链CMC命令行工具的：[证书别名章节](../dev/命令行工具.html#alias)

> tips: 新旧证书的地址是一样的，可通过cmc命令验证： `./cmc address cert-to-addr /client.sign.crt`

#### 节点证书

因证书续期node id、key不会更改，故只需要将`chainmaker.yml`中的node.cert_file， net.tls.cert_file， rpc.tls.cert_file3个证书文件替换。然后重启该节点即可。

同步节点、共识节点均一样的操作。

> tips: 新旧节点证书nodeid查看方式： ./cmc cert nid --node-cert-path xxx.tls.crt
>
> 注意：nodeid指的是tls证书： xxx.tls.crt

#### trustRoot ca证书

直接使用旧的管理员签名，发送根证书更新交易即可更新，需要一个组织一个组织更新。也可以将组织的CA删除掉再重新添加。

如果当前链上仅有一个组织，也是使用旧的管理员签名即可更新。不可删除再新增。

参考：[证书根证书更新](../dev/命令行工具.html#chainConfig.updateOrgRootCA)

### 证书续期的方式

#### 通过openssl续期

cat renew.sh

```sh
#!/bin/bash

inkey=$1
incrt=$2
cakey=$3
cacrt=$4

# 根据证书和key计算csr
echo "[get csr]"
openssl x509 -x509toreq -sha256 -in $incrt -out tmp.csr -signkey $inkey

# 根据csr使用ca key重新签发新证书
echo
echo "[renew]"
openssl x509 -req -sha256 -days 365 -in tmp.csr -out new$incrt -CAkey $cakey -CA $cacrt -CAcreateserial

# 校验新证书
echo
echo "[verify]"
openssl verify -CAfile $cacrt new$incrt
```

使用方式如下：
新证书文件名为： "new" + crtName
```sh
# 节点证书
./renew.sh common1.tls.key common1.tls.crt ../ca/ca.key ../ca/ca.crt
# 用户证书
./renew.sh client1.tls.key client1.tls.crt ca.key ca.crt

# trustRoot ca证书
./renew.sh ca.key ca.crt ca.key ca.crt
# 注意最后一步校验需要换成使用新CA证书，校验旧节点/用户证书
openssl verify -CAfile ca.new.crt common1.sign.new.crt
```

#### 通过CA-service续期

直接调用延期证书接口即可。`/api/ca/renewcert`

参考：[ca-service](CA证书服务.html#renewcert)