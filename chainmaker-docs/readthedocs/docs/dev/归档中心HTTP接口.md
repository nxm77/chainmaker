

# 归档中心

> v1.0.0

# 一、服务返回信息说明
- 服务返回code字段(integer)；>0 代表接口返回错误信息；接口成功返回数据则该字段不返回
- 服务返回errorMsg字段(string)；非空则代表错误信息；接口成功返回数据则该字段不返回 
- 服务返回data字段(json对象)；接口成功返回的具体数据 

# 二、查询归档中心信息

## 1 POST 压缩链归档数据接口

POST /admin/compress_under_height

> Body 请求参数

```json
{
  "chain_genesis_hash": "string",
  "height": 0
}
```

### 请求参数

|名称|位置|类型|必选|说明|
| :-: | :-: | :-: | :-: | :-: |
|x-token|header|string| 是 |none|
|body|body|object| 否 |none|
|» chain_genesis_hash|body|string| 是 |block hash的hex编码|
|» height|body|integer| 是 |none|

> 返回示例

### 返回结果

|状态码|状态码含义|说明|数据模型|
| :-: | :-: | :-: | :-: |  
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|  

### data 对应的数据结构
```go
map[string]int {
    "StartCompressHeight":0,
    "EndCompressHeight":  100,
}
```

## 2 POST 获取归档状态

POST /get_archive_status

> Body 请求参数

```json
{
  "chain_genesis_hash": "string"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» chain_genesis_hash|body|string| 是 |block hash的hex编码|

> 返回示例

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### data返回数据结构
```go
map[string]interface{}{
    "height":    100, // 归档中心已经归档的最大高度
    "inArchive": true, // 链上有节点正在归档数据  
}
```

## 3 POST 获取已经归档高度

POST /get_archived_height

> Body 请求参数

```json
{
  "chain_genesis_hash": "string"
}
```

### 请求参数

| 名称                 | 位置 | 类型   | 必选 | 说明                |
| -------------------- | ---- | ------ | ---- | ------------------- |
| body                 | body | object | 否   | none                |
| » chain_genesis_hash | body | string | 是   | block hash的hex编码 |

> 返回示例

### 返回结果

| 状态码 | 状态码含义                                              | 说明 | 数据模型 |
| ------ | ------------------------------------------------------- | ---- | -------- |
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | 成功 | Inline   |

### data返回数据结构

```go
100 // 获取归档中心已经归档的最大高度 
```



## 4 POST 获取链数据压缩状态

POST /get_chain_compress

> Body 请求参数

```json
{
  "chain_genesis_hash": "string"
}
```

### 请求参数

| 名称                 | 位置 | 类型   | 必选 | 说明                |
| -------------------- | ---- | ------ | ---- | ------------------- |
| body                 | body | object | 否   | none                |
| » chain_genesis_hash | body | string | 是   | block hash的hex编码 |

> 返回示例

### 返回结果

| 状态码 | 状态码含义                                              | 说明 | 数据模型 |
| ------ | ------------------------------------------------------- | ---- | -------- |
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | 成功 | Inline   |

### data返回数据结构

```go
map[string]interface{}{
    "CompressedHegiht": 100,// 当前压缩的最大高度
    "IsCompressing": true, // 归档中心的链正在被压缩 
}

```

## 5 POST 增加ca证书

POST /admin/add_ca

> Body 请求参数

```yaml
ca_name: string

```

### 请求参数

| 名称      | 位置   | 类型           | 必选 | 说明 |
| --------- | ------ | -------------- | ---- | ---- |
| x-token   | header | string         | 是   | none |
| body      | body   | object         | 否   | none |
| » ca_name | body   | string(binary) | 是   | none |

> 返回示例

### 返回结果

| 状态码 | 状态码含义                                              | 说明 | 数据模型 |
| ------ | ------------------------------------------------------- | ---- | -------- |
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | 成功 | Inline   |

### data返回数据结构

```go
"server add ca successful" //成功上传ca文件,data返回为字符串 
```

## 6 POST 获取当前归档中心所有链的信息

POST /get_chains_infos

> Body 请求参数



### 请求参数



> 返回示例

### 返回结果

|状态码|状态码含义|说明|数据模型|
| :-: | :-: | :-: | :-: |
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|  

### data 对应的数据结构  

```go
[]ChainStatus // 返回的各链的状态
```
## 7 POST 根据区块hash的base64编码计算hash的hex编码 
POST /get_hashhex_by_hashbyte

> Body 请求参数

```json
{
  "block_hash": "rYYRMGHFgCke/f40Bchd/wkvmKRHcXkA68S6n2B+k70="
}
```

> 返回示例

### 返回结果

|状态码|状态码含义|说明|数据模型|       
|---|---|---|---|  
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|   

### data 对应的数据结构
```go
 bdaeaa2b435ad98d553019d1455234c4dd6341671018fd571b5498d68eacd7b5// 根据调用接口传入的blockhash的base64编码计算的hex编码
```
# 三、查询区块、交易信息

## 1 POST 根据高度查全区块信息

POST /get_full_block_by_height

> Body 请求参数

```json
{
  "chain_genesis_hash": "string",
  "height": 0
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» chain_genesis_hash|body|string| 是 |block hash的hex编码|
|» height|body|integer| 是 |none|

> 返回示例

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### data返回数据结构
```go
参考下方数据模型：BlockWithRWSet

如：
{
    "data": {
        "block": {
            "header": {
                "block_height": 110,
                "block_hash": "jfg04985ter"
            }
        },
        "txRWSets": {
            
        }
    }
}
```

## 2 POST 根据hash查询带有读写集区块

POST /get_block_info_by_hash

> Body 请求参数

```json
{
  "chain_genesis_hash": "string",
  "block_hash": "string"
}
```

### 请求参数

| 名称                 | 位置 | 类型   | 必选 | 说明                             |
| -------------------- | ---- | ------ | ---- | -------------------------------- |
| body                 | body | object | 否   | none                             |
| » chain_genesis_hash | body | string | 是   | block hash的hex编码              |
| » block_hash         | body | string | 是   | block hash 的hex编码或base64编码 |

> 返回示例

### 返回结果

| 状态码 | 状态码含义                                              | 说明 | 数据模型 |
| ------ | ------------------------------------------------------- | ---- | -------- |
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | 成功 | Inline   |

### data返回数据结构

```go
参考下方数据模型：BlockInfo
```



## 3 POST 根据txid查询带有读写集的区块

POST /get_block_info_by_txid

> Body 请求参数

```json
{
  "chain_genesis_hash": "string",
  "tx_id": "string"
}
```

### 请求参数

| 名称                 | 位置 | 类型   | 必选 | 说明                |
| -------------------- | ---- | ------ | ---- | ------------------- |
| body                 | body | object | 否   | none                |
| » chain_genesis_hash | body | string | 是   | block hash的hex编码 |
| » tx_id              | body | string | 是   | none                |

> 返回示例

### 返回结果

| 状态码 | 状态码含义                                              | 说明 | 数据模型 |
| ------ | ------------------------------------------------------- | ---- | -------- |
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | 成功 | Inline   |

### data返回数据结构

```go
参考下方数据模型：BlockInfo
```



## 4 POST 根据txid查询带有读写集的事务 

POST /get_full_transaction_info_by_txid

> Body 请求参数

```json
{
  "chain_genesis_hash": "string",
  "tx_id": "string"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» chain_genesis_hash|body|string| 是 |block hash的hex编码|
|» tx_id|body|string| 是 |none|

> 返回示例

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### data返回数据结构
```go
参考下方数据模型：TransactionInfoWithRWSet

```

## 5 POST 根据txid查询裁剪后的交易信息

POST /get_truncate_tx_by_txid

> Body 请求参数

```json
{
  "chain_genesis_hash": "string",
  "tx_id": "string",
  "with_rwset": true,
  "truncate_length": 0,
  "truncate_model": "string"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» chain_genesis_hash|body|string| 是 |block hash的hex编码|
|» tx_id|body|string| 是 |txid|
|» with_rwset|body|boolean| 否 |是否带有rwset|
|» truncate_length|body|integer| 否 |裁剪长度|
|» truncate_model|body|string| 否 |可以为hash/truncate/empty|

> 返回示例

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### data返回数据结构

```go
参考下方数据模型：TransactionInfoWithRWSet
```

## 6 POST 根据txid计算merklepath

POST /get_merklepath_by_txid

> Body 请求参数

```json
{
  "chain_genesis_hash": "string",
  "tx_id": "string"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» chain_genesis_hash|body|string| 是 |block hash的hex编码|
|» tx_id|body|string| 是 |none|

> 返回示例

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### data返回数据结构

```go
参考下方数据模型：[]byte
```

## 7 POST 根据高度查询裁剪后的区块信息

POST /get_truncate_block_by_height

> Body 请求参数

```json
{
  "chain_genesis_hash": "string",
  "height": 0,
  "with_rwset": true,
  "truncate_length": 0,
  "truncate_model": "string"
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» chain_genesis_hash|body|string| 是 |block hash的hex编码|
|» height|body|integer| 是 |none|
|» with_rwset|body|boolean| 否 |是否带有rwset|
|» truncate_length|body|integer| 否 |裁剪长度|
|» truncate_model|body|string| 否 |可以为hash/truncate/empty|

> 返回示例

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|成功|Inline|

### data返回数据结构

```go
参考下方数据模型：BlockInfo
```


# 数据模型



## 接口返回的数据结构参考  

- BlockInfo 

```go
type BlockInfo struct {
	// block
	Block *Block `protobuf:"bytes,1,opt,name=block,proto3" json:"block,omitempty"`
	// The read/write set list corresponding to the transaction included in the block
	RwsetList []*TxRWSet `protobuf:"bytes,2,rep,name=rwset_list,json=rwsetList,proto3" json:"rwset_list,omitempty"`
}
```

- Block

```go
type Block struct {
	// header of the block
	Header *BlockHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	// execution sequence of intra block transactions, generated by proposer
	Dag *DAG `protobuf:"bytes,2,opt,name=dag,proto3" json:"dag,omitempty"`
	// transaction list in this block
	Txs []*Transaction `protobuf:"bytes,3,rep,name=txs,proto3" json:"txs,omitempty"`
	// stores the voting information of the current block
	// not included in block hash value calculation
	AdditionalData *AdditionalData `protobuf:"bytes,4,opt,name=additional_data,json=additionalData,proto3" json:"additional_data,omitempty"`
}
```

- BlockHeader 
```go
type BlockHeader struct {
	// block version
	BlockVersion uint32 `protobuf:"varint,1,opt,name=block_version,json=blockVersion,proto3" json:"block_version,omitempty"`
	// config block or normal block or other else
	BlockType BlockType `protobuf:"varint,2,opt,name=block_type,json=blockType,proto3,enum=common.BlockType" json:"block_type,omitempty"`
	// blockchain identifier
	ChainId string `protobuf:"bytes,3,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	// block height
	BlockHeight uint64 `protobuf:"varint,4,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
	// block hash (block identifier)
	BlockHash []byte `protobuf:"bytes,5,opt,name=block_hash,json=blockHash,proto3" json:"block_hash,omitempty"`
	// previous block hash
	PreBlockHash []byte `protobuf:"bytes,6,opt,name=pre_block_hash,json=preBlockHash,proto3" json:"pre_block_hash,omitempty"`
	// previous config block height, used to trace and check if chain config is valid
	PreConfHeight uint64 `protobuf:"varint,7,opt,name=pre_conf_height,json=preConfHeight,proto3" json:"pre_conf_height,omitempty"`
	// count of transactions
	TxCount uint32 `protobuf:"varint,8,opt,name=tx_count,json=txCount,proto3" json:"tx_count,omitempty"`
	// merkle root of transactions
	// used to verify the existence of this transactions
	TxRoot []byte `protobuf:"bytes,9,opt,name=tx_root,json=txRoot,proto3" json:"tx_root,omitempty"`
	// Save the DAG feature summary, and hash the DAG after Pb serialization
	// hash of serialized DAG
	DagHash []byte `protobuf:"bytes,10,opt,name=dag_hash,json=dagHash,proto3" json:"dag_hash,omitempty"`
	// The root hash of Merkle tree generated by read_write_set_digest in the result of each transaction in the block
	// used to verify the read-write set of the block
	RwSetRoot []byte `protobuf:"bytes,11,opt,name=rw_set_root,json=rwSetRoot,proto3" json:"rw_set_root,omitempty"`
	// the time stamp of the block
	BlockTimestamp int64 `protobuf:"varint,12,opt,name=block_timestamp,json=blockTimestamp,proto3" json:"block_timestamp,omitempty"`
	// consensus parameters
	// used to store information, include in block hash calculation
	ConsensusArgs []byte `protobuf:"bytes,13,opt,name=consensus_args,json=consensusArgs,proto3" json:"consensus_args,omitempty"`
	// proposal node identifier
	Proposer *accesscontrol.Member `protobuf:"bytes,14,opt,name=proposer,proto3" json:"proposer,omitempty"`
	// signature of proposer
	Signature []byte `protobuf:"bytes,15,opt,name=signature,proto3" json:"signature,omitempty"`
}
```
- DAG_Neighbor
```go
type DAG_Neighbor struct {
	Neighbors []uint32 `protobuf:"varint,1,rep,packed,name=neighbors,proto3" json:"neighbors,omitempty"`
}
```
- DAG
```go
type DAG struct {
	// sequence number of transaction topological sort
	// the sequence number of the transaction topological sort associated with the transaction
	Vertexes []*DAG_Neighbor `protobuf:"bytes,2,rep,name=vertexes,proto3" json:"vertexes,omitempty"`
}
```
- KeyValuePair
```go
type KeyValuePair struct {
	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}
```
- Payload
```go
type Payload struct {
	// blockchain identifier
	ChainId string `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	// transaction type
	TxType TxType `protobuf:"varint,2,opt,name=tx_type,json=txType,proto3,enum=common.TxType" json:"tx_type,omitempty"`
	// transaction id set by sender, should be unique
	TxId string `protobuf:"bytes,3,opt,name=tx_id,json=txId,proto3" json:"tx_id,omitempty"`
	// transaction timestamp, in unix timestamp format, seconds
	Timestamp int64 `protobuf:"varint,4,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	// expiration timestamp in unix timestamp format
	// after that the transaction is invalid if it is not included in block yet
	ExpirationTime int64 `protobuf:"varint,5,opt,name=expiration_time,json=expirationTime,proto3" json:"expiration_time,omitempty"`
	// smart contract name
	ContractName string `protobuf:"bytes,6,opt,name=contract_name,json=contractName,proto3" json:"contract_name,omitempty"`
	// invoke method
	Method string `protobuf:"bytes,7,opt,name=method,proto3" json:"method,omitempty"`
	// invoke parameters in k-v format
	Parameters []*KeyValuePair `protobuf:"bytes,8,rep,name=parameters,proto3" json:"parameters,omitempty"`
	// sequence number, default is 0
	Sequence uint64 `protobuf:"varint,9,opt,name=sequence,proto3" json:"sequence,omitempty"`
	// transaction limitation
	Limit *Limit `protobuf:"bytes,10,opt,name=limit,proto3" json:"limit,omitempty"`
}
```
- Member 
```go
type Member struct {
	// organization identifier of the member
	OrgId string `protobuf:"bytes,1,opt,name=org_id,json=orgId,proto3" json:"org_id,omitempty"`
	// member type
	MemberType MemberType `protobuf:"varint,2,opt,name=member_type,json=memberType,proto3,enum=accesscontrol.MemberType" json:"member_type,omitempty"`
	// member identity related info bytes
	MemberInfo []byte `protobuf:"bytes,3,opt,name=member_info,json=memberInfo,proto3" json:"member_info,omitempty"`
}
```
- EndorsementEntry
```go
type EndorsementEntry struct {
	// signer
	Signer *accesscontrol.Member `protobuf:"bytes,1,opt,name=signer,proto3" json:"signer,omitempty"`
	// signature
	Signature []byte `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
}
```
- ContractResult
```go
type ContractResult struct {
	// user contract defined return code, 0-ok, >0 user define error code. for example, insufficient balance in token transfer
	Code uint32 `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	// user contract defined result
	Result []byte `protobuf:"bytes,2,opt,name=result,proto3" json:"result,omitempty"`
	// user contract defined result message
	Message string `protobuf:"bytes,3,opt,name=message,proto3" json:"message,omitempty"`
	// gas used by current contract(include contract call)
	GasUsed uint64 `protobuf:"varint,4,opt,name=gas_used,json=gasUsed,proto3" json:"gas_used,omitempty"`
	// contract events
	ContractEvent []*ContractEvent `protobuf:"bytes,5,rep,name=contract_event,json=contractEvent,proto3" json:"contract_event,omitempty"`
}
```

- Result
```go
type Result struct {
	// response code
	Code TxStatusCode `protobuf:"varint,1,opt,name=code,proto3,enum=common.TxStatusCode" json:"code,omitempty"`
	// returned data, set in smart contract
	ContractResult *ContractResult `protobuf:"bytes,2,opt,name=contract_result,json=contractResult,proto3" json:"contract_result,omitempty"`
	// hash of the transaction's read-write set
	RwSetHash []byte `protobuf:"bytes,3,opt,name=rw_set_hash,json=rwSetHash,proto3" json:"rw_set_hash,omitempty"`
	Message   string `protobuf:"bytes,4,opt,name=message,proto3" json:"message,omitempty"`
}
```
- Transaction
```go
type Transaction struct {
	// payload
	Payload *Payload `protobuf:"bytes,1,opt,name=payload,proto3" json:"payload,omitempty"`
	// sender account and signature
	Sender *EndorsementEntry `protobuf:"bytes,2,opt,name=sender,proto3" json:"sender,omitempty"`
	// endorser accounts and signatures
	Endorsers []*EndorsementEntry `protobuf:"bytes,3,rep,name=endorsers,proto3" json:"endorsers,omitempty"`
	// result of the transaction
	Result *Result `protobuf:"bytes,4,opt,name=result,proto3" json:"result,omitempty"`
}

```
- AdditionalData 
```go
type AdditionalData struct {
	// extra data, with map type, excluded in hash calculation
	ExtraData map[string][]byte `protobuf:"bytes,1,rep,name=extra_data,json=extraData,proto3" json:"extra_data,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

```
- 
- BlockWithRWSet
```go
type BlockWithRWSet struct {
	// block data
	Block *common.Block `protobuf:"bytes,1,opt,name=block,proto3" json:"block,omitempty"`
	// transaction read/write set of blocks
	TxRWSets []*common.TxRWSet `protobuf:"bytes,2,rep,name=txRWSets,proto3" json:"txRWSets,omitempty"`
	// contract event info
	ContractEvents []*common.ContractEvent `protobuf:"bytes,3,rep,name=contract_events,json=contractEvents,proto3" json:"contract_events,omitempty"`
}
```
- ContractEvent
```go
type ContractEvent struct {
	Topic           string   `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
	TxId            string   `protobuf:"bytes,2,opt,name=tx_id,json=txId,proto3" json:"tx_id,omitempty"`
	ContractName    string   `protobuf:"bytes,3,opt,name=contract_name,json=contractName,proto3" json:"contract_name,omitempty"`
	ContractVersion string   `protobuf:"bytes,4,opt,name=contract_version,json=contractVersion,proto3" json:"contract_version,omitempty"`
	EventData       []string `protobuf:"bytes,5,rep,name=event_data,json=eventData,proto3" json:"event_data,omitempty"`
}
```
- TxRWSet
```go
type TxRWSet struct {
	// transaction identifier
	TxId string `protobuf:"bytes,1,opt,name=tx_id,json=txId,proto3" json:"tx_id,omitempty"`
	// read set
	TxReads []*TxRead `protobuf:"bytes,2,rep,name=tx_reads,json=txReads,proto3" json:"tx_reads,omitempty"`
	// write set
	TxWrites []*TxWrite `protobuf:"bytes,3,rep,name=tx_writes,json=txWrites,proto3" json:"tx_writes,omitempty"`
}
```
- TxRead
```go
type TxRead struct {
	// read key
	Key []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// the value of the key
	Value []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	// contract name, used in cross-contract calls
	// set to null if only the contract in transaction request is called
	ContractName string `protobuf:"bytes,3,opt,name=contract_name,json=contractName,proto3" json:"contract_name,omitempty"`
	// read key version
	Version *KeyVersion `protobuf:"bytes,4,opt,name=version,proto3" json:"version,omitempty"`
}
```
- KeyVersion  
```go
type KeyVersion struct {
	// the transaction identifier that last modified the key
	RefTxId string `protobuf:"bytes,3,opt,name=ref_tx_id,json=refTxId,proto3" json:"ref_tx_id,omitempty"`
	// the offset of the key in the write set of the transaction, starts from 0
	RefOffset int32 `protobuf:"varint,4,opt,name=ref_offset,json=refOffset,proto3" json:"ref_offset,omitempty"`
}
```
- TxWrite  
```go
// TxRead describes a write/delete operation on a key
type TxWrite struct {
	// write key
	Key []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	// write value
	Value []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	// contract name, used in cross-contract calls
	// set to null if only the contract in transaction request is called
	ContractName string `protobuf:"bytes,3,opt,name=contract_name,json=contractName,proto3" json:"contract_name,omitempty"`
}
```
- ChainConfig 
```go
type ChainConfig struct {
	// blockchain identifier
	ChainId string `protobuf:"bytes,1,opt,name=chain_id,json=chainId,proto3" json:"chain_id,omitempty"`
	// blockchain version
	Version string `protobuf:"bytes,2,opt,name=version,proto3" json:"version,omitempty"`
	// member type
	AuthType string `protobuf:"bytes,3,opt,name=auth_type,json=authType,proto3" json:"auth_type,omitempty"`
	// config sequence
	Sequence uint64 `protobuf:"varint,4,opt,name=sequence,proto3" json:"sequence,omitempty"`
	// encryption algorithm related configuration
	Crypto *CryptoConfig `protobuf:"bytes,5,opt,name=crypto,proto3" json:"crypto,omitempty"`
	// block related configuration
	Block *BlockConfig `protobuf:"bytes,6,opt,name=block,proto3" json:"block,omitempty"`
	// core module related configuration
	Core *CoreConfig `protobuf:"bytes,7,opt,name=core,proto3" json:"core,omitempty"`
	// consensus related configuration
	Consensus *ConsensusConfig `protobuf:"bytes,8,opt,name=consensus,proto3" json:"consensus,omitempty"`
	// trusted root related configuration
	// for alliance members, the initial member's root info of the consortium; for public chain, there is no need to configure
	// Key: node_id; value: address, node public key / CA certificate
	TrustRoots   []*TrustRootConfig   `protobuf:"bytes,9,rep,name=trust_roots,json=trustRoots,proto3" json:"trust_roots,omitempty"`
	TrustMembers []*TrustMemberConfig `protobuf:"bytes,10,rep,name=trust_members,json=trustMembers,proto3" json:"trust_members,omitempty"`
	// permission related configuration
	ResourcePolicies []*ResourcePolicy `protobuf:"bytes,11,rep,name=resource_policies,json=resourcePolicies,proto3" json:"resource_policies,omitempty"`
	Contract         *ContractConfig   `protobuf:"bytes,12,opt,name=contract,proto3" json:"contract,omitempty"`
	// snapshot module related configuration
	Snapshot *SnapshotConfig `protobuf:"bytes,13,opt,name=snapshot,proto3" json:"snapshot,omitempty"`
	// scheduler module related configuration
	Scheduler *SchedulerConfig `protobuf:"bytes,14,opt,name=scheduler,proto3" json:"scheduler,omitempty"`
	// tx sim context module related configuration
	Context *ContextConfig `protobuf:"bytes,15,opt,name=context,proto3" json:"context,omitempty"`
	// disabled native contracts list for permission control purposes
	DisabledNativeContract []string `protobuf:"bytes,16,rep,name=disabled_native_contract,json=disabledNativeContract,proto3" json:"disabled_native_contract,omitempty"`
	// gas account config
	AccountConfig *GasAccountConfig `protobuf:"bytes,18,opt,name=account_config,json=accountConfig,proto3" json:"account_config,omitempty"`
	// vm config
	Vm *Vm `protobuf:"bytes,17,opt,name=vm,proto3" json:"vm,omitempty"`
}
```
- CryptoConfig
```go
type CryptoConfig struct {
	// enable Transaction timestamp verification or Not
	Hash string `protobuf:"bytes,1,opt,name=hash,proto3" json:"hash,omitempty"`
}
```
- BlockConfig
```go
type BlockConfig struct {
	// enable transaction timestamp verification or Not
	TxTimestampVerify bool `protobuf:"varint,1,opt,name=tx_timestamp_verify,json=txTimestampVerify,proto3" json:"tx_timestamp_verify,omitempty"`
	// expiration time of transaction timestamp (seconds)
	TxTimeout uint32 `protobuf:"varint,2,opt,name=tx_timeout,json=txTimeout,proto3" json:"tx_timeout,omitempty"`
	// maximum number of transactions in a block
	BlockTxCapacity uint32 `protobuf:"varint,3,opt,name=block_tx_capacity,json=blockTxCapacity,proto3" json:"block_tx_capacity,omitempty"`
	// maximum block size, in MB
	BlockSize uint32 `protobuf:"varint,4,opt,name=block_size,json=blockSize,proto3" json:"block_size,omitempty"`
	// block proposing interval, in ms
	BlockInterval uint32 `protobuf:"varint,5,opt,name=block_interval,json=blockInterval,proto3" json:"block_interval,omitempty"`
	// maximum size of transaction's parameter, in MB
	TxParameterSize uint32 `protobuf:"varint,6,opt,name=tx_parameter_size,json=txParameterSize,proto3" json:"tx_parameter_size,omitempty"`
}
```
- CoreConfig
```go
type CoreConfig struct {
	// [0, 60], the time when the transaction scheduler gets the transaction from the transaction pool to schedule
	TxSchedulerTimeout uint64 `protobuf:"varint,1,opt,name=tx_scheduler_timeout,json=txSchedulerTimeout,proto3" json:"tx_scheduler_timeout,omitempty"`
	// [0, 60], the time-out for verification after the transaction scheduler obtains the transaction from the block
	TxSchedulerValidateTimeout uint64 `protobuf:"varint,2,opt,name=tx_scheduler_validate_timeout,json=txSchedulerValidateTimeout,proto3" json:"tx_scheduler_validate_timeout,omitempty"`
	// the configuration of consensus message turbo
	ConsensusTurboConfig *ConsensusTurboConfig `protobuf:"bytes,3,opt,name=consensus_turbo_config,json=consensusTurboConfig,proto3" json:"consensus_turbo_config,omitempty"`
	// enable sender group, used for handling txs with sender conflicts efficiently
	EnableSenderGroup bool `protobuf:"varint,4,opt,name=enable_sender_group,json=enableSenderGroup,proto3" json:"enable_sender_group,omitempty"`
	// enable conflicts bit window, used for dynamic tuning the capacity of tx execution goroutine pool
	EnableConflictsBitWindow bool `protobuf:"varint,5,opt,name=enable_conflicts_bit_window,json=enableConflictsBitWindow,proto3" json:"enable_conflicts_bit_window,omitempty"`
	// enable optimize charge gas for the same account transactions
	EnableOptimizeChargeGas bool `protobuf:"varint,6,opt,name=enable_optimize_charge_gas,json=enableOptimizeChargeGas,proto3" json:"enable_optimize_charge_gas,omitempty"`
}
```
- ConsensusTurboConfig
```go
type ConsensusTurboConfig struct {
	// switch of consensus message turbo
	ConsensusMessageTurbo bool `protobuf:"varint,1,opt,name=consensus_message_turbo,json=consensusMessageTurbo,proto3" json:"consensus_message_turbo,omitempty"`
	// retry time of get tx by txIds from txpool
	RetryTime uint64 `protobuf:"varint,2,opt,name=retry_time,json=retryTime,proto3" json:"retry_time,omitempty"`
	// the interval of retry get tx by txIds from txpool(ms)
	RetryInterval uint64 `protobuf:"varint,3,opt,name=retry_interval,json=retryInterval,proto3" json:"retry_interval,omitempty"`
}
```

- ConsensusConfig
```go
type ConsensusConfig struct {
	// consensus type
	Type consensus.ConsensusType `protobuf:"varint,1,opt,name=type,proto3,enum=consensus.ConsensusType" json:"type,omitempty"`
	// organization list of nodes
	Nodes []*OrgConfig `protobuf:"bytes,2,rep,name=nodes,proto3" json:"nodes,omitempty"`
	// expand the field, record the difficulty, reward and other consensus algorithm configuration
	ExtConfig []*ConfigKeyValue `protobuf:"bytes,3,rep,name=ext_config,json=extConfig,proto3" json:"ext_config,omitempty"`
	// Initialize the configuration of DPOS
	DposConfig []*ConfigKeyValue `protobuf:"bytes,4,rep,name=dpos_config,json=dposConfig,proto3" json:"dpos_config,omitempty"`
}
```
- OrgConfig 
```go
type OrgConfig struct {
	// organization identifier
	OrgId string `protobuf:"bytes,1,opt,name=org_id,json=orgId,proto3" json:"org_id,omitempty"`
	// address list owned by the organization
	// Deprecated , replace by node_id
	Address []string `protobuf:"bytes,2,rep,name=address,proto3" json:"address,omitempty"`
	// node id list owned by the organization
	NodeId []string `protobuf:"bytes,3,rep,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
}
```
- ConfigKeyValue 
```go
type ConfigKeyValue struct {
	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value string `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}
```

- TrustRootConfig
```go
type TrustRootConfig struct {
	// oranization ideftifier
	OrgId string `protobuf:"bytes,1,opt,name=org_id,json=orgId,proto3" json:"org_id,omitempty"`
	// root certificate / public key
	Root []string `protobuf:"bytes,2,rep,name=root,proto3" json:"root,omitempty"`
}
```
- TrustMemberConfig
```go
type TrustMemberConfig struct {
	// member info
	MemberInfo string `protobuf:"bytes,1,opt,name=member_info,json=memberInfo,proto3" json:"member_info,omitempty"`
	// oranization ideftifier
	OrgId  string `protobuf:"bytes,2,opt,name=org_id,json=orgId,proto3" json:"org_id,omitempty"`
	Role   string `protobuf:"bytes,3,opt,name=role,proto3" json:"role,omitempty"`
	NodeId string `protobuf:"bytes,4,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
}
```

- Policy
```go
type Policy struct {
	// rule keywords, e.g., ANY/MAJORITY/ALL/SELF/a number/a rate
	Rule string `protobuf:"bytes,1,opt,name=rule,proto3" json:"rule,omitempty"`
	// org_list describes the organization set included in the authentication
	OrgList []string `protobuf:"bytes,2,rep,name=org_list,json=orgList,proto3" json:"org_list,omitempty"`
	// role_list describes the role set included in the authentication
	// e.g., admin/client/consensus/common
	RoleList []string `protobuf:"bytes,3,rep,name=role_list,json=roleList,proto3" json:"role_list,omitempty"`
}
```
- ResourcePolicy 
```go
type ResourcePolicy struct {
	// resource name
	ResourceName string `protobuf:"bytes,1,opt,name=resource_name,json=resourceName,proto3" json:"resource_name,omitempty"`
	// policy(permission)
	Policy *accesscontrol.Policy `protobuf:"bytes,2,opt,name=policy,proto3" json:"policy,omitempty"`
}
```

- ContractConfig
```go
type ContractConfig struct {
	EnableSqlSupport bool `protobuf:"varint,1,opt,name=enable_sql_support,json=enableSqlSupport,proto3" json:"enable_sql_support,omitempty"`
	// disabled native contracts list for permission control purposes
	DisabledNativeContract []string `protobuf:"bytes,2,rep,name=disabled_native_contract,json=disabledNativeContract,proto3" json:"disabled_native_contract,omitempty"`
}
```

- SnapshotConfig
```go
type SnapshotConfig struct {
	// for the evidence contract
	EnableEvidence bool `protobuf:"varint,1,opt,name=enable_evidence,json=enableEvidence,proto3" json:"enable_evidence,omitempty"`
}
```
- SchedulerConfig
```go
type SchedulerConfig struct {
	// for evidence contract
	EnableEvidence bool `protobuf:"varint,1,opt,name=enable_evidence,json=enableEvidence,proto3" json:"enable_evidence,omitempty"`
}
```
- ContextConfig
```go
type ContextConfig struct {
	// for the evidence contract
	EnableEvidence bool `protobuf:"varint,1,opt,name=enable_evidence,json=enableEvidence,proto3" json:"enable_evidence,omitempty"`
}
```
- GasAccountConfig
```go
type GasAccountConfig struct {
	// for admin address
	GasAdminAddress string `protobuf:"bytes,1,opt,name=gas_admin_address,json=gasAdminAddress,proto3" json:"gas_admin_address,omitempty"`
	// for admin gas count
	GasCount uint32 `protobuf:"varint,2,opt,name=gas_count,json=gasCount,proto3" json:"gas_count,omitempty"`
	// for gas manager
	EnableGas bool `protobuf:"varint,3,opt,name=enable_gas,json=enableGas,proto3" json:"enable_gas,omitempty"`
	// by default, gas value per transaction
	DefaultGas uint64 `protobuf:"varint,4,opt,name=default_gas,json=defaultGas,proto3" json:"default_gas,omitempty"`
}
```
- Vm
```go
type Vm struct {
	SupportList []string `protobuf:"bytes,1,rep,name=support_list,json=supportList,proto3" json:"support_list,omitempty"`
	AddrType    AddrType `protobuf:"varint,2,opt,name=addr_type,json=addrType,proto3,enum=config.AddrType" json:"addr_type,omitempty"`
}
```

- TransactionInfoWithRWSet
```go
type TransactionInfoWithRWSet struct {
	// transaction raw data
	Transaction *Transaction `protobuf:"bytes,1,opt,name=transaction,proto3" json:"transaction,omitempty"`
	// block height
	BlockHeight uint64 `protobuf:"varint,2,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
	// Deprecated, block hash
	BlockHash []byte `protobuf:"bytes,3,opt,name=block_hash,json=blockHash,proto3" json:"block_hash,omitempty"`
	// transaction index in block
	TxIndex uint32 `protobuf:"varint,4,opt,name=tx_index,json=txIndex,proto3" json:"tx_index,omitempty"`
	// block header timestamp
	BlockTimestamp int64    `protobuf:"varint,5,opt,name=block_timestamp,json=blockTimestamp,proto3" json:"block_timestamp,omitempty"`
	RwSet          *TxRWSet `protobuf:"bytes,6,opt,name=rw_set,json=rwSet,proto3" json:"rw_set,omitempty"`
}
```
- TransactionStoreInfo
```go
type TransactionStoreInfo struct {
	// transaction raw data
	Transaction *common.Transaction `protobuf:"bytes,1,opt,name=transaction,proto3" json:"transaction,omitempty"`
	// block height
	BlockHeight uint64 `protobuf:"varint,2,opt,name=block_height,json=blockHeight,proto3" json:"block_height,omitempty"`
	// Deprecated, block hash
	BlockHash []byte `protobuf:"bytes,3,opt,name=block_hash,json=blockHash,proto3" json:"block_hash,omitempty"`
	// transaction index in block
	TxIndex uint32 `protobuf:"varint,4,opt,name=tx_index,json=txIndex,proto3" json:"tx_index,omitempty"`
	// block header timestamp
	BlockTimestamp int64 `protobuf:"varint,5,opt,name=block_timestamp,json=blockTimestamp,proto3" json:"block_timestamp,omitempty"`
	// transaction offset index in file
	TransactionStoreInfo *StoreInfo `protobuf:"bytes,6,opt,name=transaction_store_info,json=transactionStoreInfo,proto3" json:"transaction_store_info,omitempty"`
}
```
- StoreInfo
```go
type StoreInfo struct {
	//store type
	StoreType DataStoreType `protobuf:"varint,1,opt,name=store_type,json=storeType,proto3,enum=store.DataStoreType" json:"store_type,omitempty"`
	// file name
	FileName string `protobuf:"bytes,2,opt,name=file_name,json=fileName,proto3" json:"file_name,omitempty"`
	// offset in file
	Offset uint64 `protobuf:"varint,3,opt,name=offset,proto3" json:"offset,omitempty"`
	// data length
	ByteLen uint64 `protobuf:"varint,4,opt,name=byte_len,json=byteLen,proto3" json:"byte_len,omitempty"`
}
```

- SerializedBlock
```go
type SerializedBlock struct {
	// header of block
	Header *common.BlockHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	// transaction execution sequence of the block, described by DAG
	Dag *common.DAG `protobuf:"bytes,2,opt,name=dag,proto3" json:"dag,omitempty"`
	// transaction id list within the block
	TxIds []string `protobuf:"bytes,3,rep,name=tx_ids,json=txIds,proto3" json:"tx_ids,omitempty"`
	// block additional data, not included in block hash calculation
	AdditionalData *common.AdditionalData `protobuf:"bytes,4,opt,name=additional_data,json=additionalData,proto3" json:"additional_data,omitempty"`
}
```

- ChainStatus
```go
type ChainStatus struct {
	ChainId        string `json:"chainId"` // 链id
	GenesisHashStr string `json:"genesisHashStr"` // 链的genesisblock的hash的base64编码
	GenesisHashHex string `json:"genesisHashHex"` // 链的genesisblock的hash的hex编码
	ArchivedHeight uint64 `json:"archivedHeight"` // 链的当前归档的高度
	InArchive      bool   `json:"inArchive"` // 当前是否有链的节点在向归档中心归档数据 
}
```
