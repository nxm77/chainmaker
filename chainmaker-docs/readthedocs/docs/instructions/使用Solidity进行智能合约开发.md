# 使用Solidity进行智能合约开发

读者对象：本章节主要描述使用Solidity进行ChainMaker合约编写的方法，主要面向于使用Solidity进行ChainMaker的合约开发的开发者。

Solidity 是一门面向合约的、为实现智能合约而创建的高级编程语言。这门语言受到了 C++，Python 和 Javascript 语言的影响，设计的目的是能在虚拟机（EVM）上运行。

Solidity 是静态类型语言，支持继承、库和复杂的用户定义类型等特性。

**概览**
1、运行时虚拟机类型（runtime_type）： EVM
2、介绍了环境依赖
3、介绍了开发方式及sdk接口
4、提供了一个示例合约


## 环境依赖

1. 操作系统

   无特殊要求，Linux、Mac和Windows均支持。

2. 软件依赖

   无

3. 长安链环境准备

   启动长安链，以及长安链CMC工具，用于将写编写好的合约，部署到链上进行测试。相关安装教程请详见：

   - [部署长安链教程。](../quickstart/通过命令行体验链.md)
   - [部署长安链CMC工具的教程。](../dev/命令行工具.md)


## 编写Solidity智能合约

### 搭建开发环境
开发者无需自己搭建开发环境，可使用[Remix](https://remix.ethereum.org)在线IDE开发solidity合约。长安链对solidity完全兼容，使用Remix开发的或者以太坊生态内的solidity合约，可直接在长安链部署运行。

### 代码编写规则

solidity语法和代码编写规则参见solidity官方开发文档：https://docs.soliditylang.org/。

#### 长安链solidity与以太坊solidity的异同

长安链的solidity在内置接口、预编译合约和跨合约调用上，与以太坊有一些区别，具体参见本章节下，各分小节详细内容。

##### solidity 内置接口

solidity为开发者提供了一些内置接口，包括内置变量和函数，可在合约中直接使用。因为长安链为无币链，所以，与原生token相关的内置接口，多默认为0，详见以下内置接口说明。


```solidity
//指定区块的区块哈希，已经不推荐使用，由 blockhash(uint blockNumber) 代替
block.blockhash(uint blockNumber) returns (bytes32)

//address类型，当前区块的出块节点地址，即 block 的 proposer
block.coinbase

//uint类型，当前区块难度，值固定为0
block.difficulty

//uint类型，当前区块 gas 限额
block.gaslimit

//uint类型，当前区块号
block.number

//uint类型，自 unix epoch 起始当前区块以秒计的时间戳
block.timestamp

//uint256类型，剩余的 gas
gasleft() returns 

//bytes类型，完整的 calldata
msg.data 

//uint类型，剩余 gas - 自 0.4.21 版本开始已经不推荐使用，由 gesleft() 代替
msg.gas 

//address类型，消息发送者（当前调用）
msg.sender 

//bytes4类型，calldata 的前 4 字节（也就是函数标识符）
msg.sig 

//uint类型，随消息发送的 wei 的数量，值固定为0
msg.value 

//uint类型，目前区块时间戳（block.timestamp）
now 

//uint类型，交易的 gas 价格，值固定为1
tx.gasprice 

//address类型，交易发起者
tx.origin 
```

##### 预编译合约

因为EVM是基于栈的虚拟机，它根据操作的内容来计算gas，所以如果牵涉到十分复杂的计算，把运算过程放在EVM中执行可能十分地低效，同时消耗非常多的gas。

预编译合约是EVM中为了提供一些不适合写成opcode的较为复杂的库函数（多用于加密、哈希等复杂运算）而采用的一种折中方案，适用于合约逻辑简单但调用频繁，或者合约逻辑固定而计算量大的场景，与以太坊相比，长安链有五个预编译合约尚未支持，正在开发中，详见以下表格。

<span id="precompiled-contract"></span>
###### 预编译合约列表

| **预编译合约名**                                      | **地址**   | **功能**                                                     |
| ----------------------------------------------------- | ---------- | ------------------------------------------------------------ |
| **ecrecover(hash, v, r, s)**                          | **0x01**   | 根据给定签名计算地址——**目前长安链尚不支持，待后续版本实现** |
| **sha256(data)**                                      | **0x02**   | 计算SHA256哈希                                               |
| **ripemd160(data)**                                   | **0x03**   | 计算RIPEMD160哈希                                            |
| **datacopy(data)**                                    | **0x04**   | 只读拷贝数据                                                 |
| **bigModExp(base, exp, mod)**                         | **0x05**   | 计算base ^ exp % mod的结果                                   |
| **bn256Add(ax, ay, bx, by)**                          | **0x06**   | BN256曲线点加法计算，成功返回(ax,ay)+(bx,by)，表示这两点是BN256曲线上的有效点，失败返回0 |
| **bn256ScalarMul(x, y, scalar)**                      | **0x07**   | BN256曲线乘法，成功返回一个曲线点scalar*(x,y)，表示(x,y)点是BN256曲线上的有效点，失败返回0 |
| **bn256Pairing(a1, b1, a2, b2, a3, b3, ..., ak, bk)** | **0x08**   | 实现BN256曲线配对操作，进行zkSNARK验证。 |
| **blake2F(rounds, h, m, t, f)**                       | **0x09**   | 实现BLAKER2b F压缩功能。 |
| **verify(pkLen, pk, msgLen, msg, signLen, sign)**     | **0x03ef** | chainmaker扩展的预编译合约，可执行**国密验签**，参数依次为序列化公钥长度、序列化公钥、消息长度、消息、签名长度、签名 |


预编译合约调用不能像普通内置接口一样直接调用，需要借助solidity的汇编语法，有一定的复杂度，所以须谨慎调用。solidity汇编语法及汇编指令列表请参见官方文档 [Solidity汇编](https://solidity-cn.readthedocs.io/zh/develop/assembly.html)。

###### 国密验签
原生solidity合约不支持sm2国密验签，长安链通过预编译合约扩展接口，为evm虚拟机添加了sm2国密验签功能，用户通过调用地址为`0x03ef`的预编译合约即可使用国密验签功能。调用传参时，需要将公钥长度、公钥、消息长度、消息、签名长度和签名通过`abi.encodePacked`函数编码为单个`bytes`字节流传入。另外，验签的签名格式支持多种编码类型，包括签名原始字节流、16进制编码、base64编码和base58编码类型，用户无序做额外处理，预编译合约会自动识别编码类型并对其验签。

以下示例为调用国密验签预编译合约verify的过程。

```sh
// SPDX-License-Identifier: GPL-3.0

pragma solidity >=0.4.21;

contract sm2 {

    function verifyTest(bytes memory pubKey, bytes memory message, bytes memory signature) public view returns (bytes32[2] memory) {
    	//返回值缓冲区
        bytes32[2] memory output;
				
		//参数封装，预编译合约接口只接受一个bytes类型的参数，如果有多个参数，需要将它们按顺序组合为一个bytes，且必须是有序的；
        bytes memory input = abi.encodePacked(pubKey.length, pubKey, message.length, message, signature.length, signature);
        uint256 inPutLen = input.length;
        
        //汇编指令指示符，在assembly { } 内的代码为汇编代码。
        assembly {
            /*使用staticcall指令调用预编译合约，指令参数依次为：
                0: 转账额，长安链为无币链，所以此参数一般为0
                0x03ef: 国密验签预编译合约verify的地址，staticcall指令将根据地址找到verify预编译合约
                add(input, 32): 预编译合约verify的参数，因为input的前32字节为参数长度，所以用add指令越过
                inPutLen: 预编译合约verify的参数的长度
                output: 预编译合约的返回值，对当前国密验签预编译合约来说，返回1表示验签成功，返回0表示验签失败
                0x40: 预编译合约返回值的长度
            */
            let success := staticcall(0, 0x03ef, add(input, 32), inPutLen, output, 0x40)
            //staticcall指令成功返回1，失败返回0
            switch success
            case 0 {
                revert(0,0)
            }
        }
        
        //预编译合约的返回值已写入output
        return output;
    }
}
```

预编译合约需要使用`call/staticcall`指令调用，如果被调合约执行时有状态数据变更，使用`call`指令，否则使用`staticcall`指令，因为预编译合约基本都不包含状态变量，所以一般使用`staticcall`指令即可，指令成功返回1，失败返回0。

注意，这里说的返回值，指的是调用指令，即`call/staticcall`的返回值，只代表该指令调用成功与否，而非被调用预编译合约本身所执行的结果。预编译合约执行的结果将作为出参的方式传入`call/staticcall`指令，预编译合约执行结束后，会将执行结果写入出参，由主调合约读取出参获得被调用预编译合约的执行结果。

`staticcall`指令的参数按顺序释义如下：

- 转账额：调用合约转给被调用合约的代币数量，因为长安链为无币链，所以赋值为0即可。
- 合约地址：`staticcall`指令将根据地址找到对应的合约，示例中传入的 `0x03ef` 为国密验签预编译合约verify的地址。
- 合约参数：被调用的合约的参数（注意和`staticcall`指令的参数区分），该参数为一个字节流类型，但solidity的`bytes`类型在起首的32字节为`bytes`长度，所以传递给预编译合约时需要去掉这个长度，一般使用`add`指令将参数（下标）向后偏移32字节，如示例中`add(input, 32)`那样。
- 合约参数长度：因为传递给合约的参数被去掉了长度，所以evm虚拟机需要根据此参数确定`input`参数的长度。
- 合约返回值缓冲区：预编译合约的返回值为两元素的`bytes32`数组，合约执行结果将被写入该返回值缓冲区，在`staticcall`指令视角，此参数用法类似于C语言的出参。国密验签合约如果执行成功，返回1，失败返回0，但因为返回值缓冲区为两元素的`bytes32`数组，其中多余的字节将被补0。
- 合约返回值长度：返回值缓冲区的字节容量，因为返回值缓冲区是2个`byte32`，所以是64字节长度，用16进制表示就是`0x40`。

以上示例就是对国密验签预编译合约的使用方式，其他预编译合约的使用方式与之相同。都需要将该预编译合约的参数封装为一个`bytes`类型的参数，并使用汇编关键字`assembly`和`staticcall`指令调用。

###### 日志

另外，solidity合约不支持打印日志，只能用event的方式，将事件写入区块，而非将日志内容写入日志文件或者打印到控制台，所以，不太方便开发者调试合约和跟踪问题。针对该情况，长安链使用预编译合约的方式，实现了合约日志打印功能。

用户通过调用地址为`0x03ee`的预编译合约，并将被打印内容作为参数输入，就可以将其打印到日志文件或者控制台了。考虑到预编译合约调用较为复杂，长安链封装了一个具备日志打印功能的基类合约，用户在使用时，只需要继承该合约，并调用`print`方法就可以了。

```sh
// SPDX-License-Identifier: MIT
// OpenZeppelin Contracts v4.4.1 (utils/Strings.sol)

pragma solidity > 0.4.22;
contract Print {
		//字符串转换方法，可以将无符号整型转换为一个string
    function toStr(uint256 value) internal pure returns (string memory) {
        if (value == 0) {
            return "0";
        }

        uint256 temp = value;
        uint256 digits;
        while (temp != 0) {
            digits++;
            temp /= 10;
        }

        bytes memory buffer = new bytes(digits);
        while (value != 0) {
            digits -= 1;
            buffer[digits] = bytes1(uint8(48 + uint256(value % 10)));
            value /= 10;
        }

        return string(buffer);
    }

		//字符串转换方法，可以将一个布尔值转换为一个string
    function toStr(bool value) internal pure returns (string memory) {
        if(value) {
            return "true";
        } else {
            return "false";
        }
    }

		//打印方法，调用该方法可以将一个字符串写入日志文件或者打印到控制台，打印内容以 evm log>> 前缀开头
    //arg level: 日志打印级别，参数取值列表[0-debug，1-info, 2-warning, 3-error]
    //arg logs: 被打印内容
    function print(uint level, string memory logs) internal view {
        bytes memory input = bytes(string.concat(toStr(level), logs));

        bytes32[2] memory result;
        uint256 len = input.length;
        assembly{
            let ok := staticcall(0, 0x03ee, add(input, 0x20), len, result, 0x40)
            switch ok
            case 0 {
                revert(0,0)
            }
        } 
    }
}
```

如示例代码所示，日志打印分为4个级别0 ~ 3，依次表示`debug`、`info`、`warning`、`errors`。用户打印时，除提供日志级别外，还需提供一个字符串参数作为日志本身。如果要打印多个参数，可先将非`string`类型参数通过`toStr`方法转换为`string`，然后再调用`string.concat`将多个字符串拼接为一个字符串作为参数传入，如以下示例所示。

```sh
// SPDX-License-Identifier: MIT
// OpenZeppelin Contracts v4.4.1 (utils/Strings.sol)

pragma solidity > 0.4.22;

//导入Print合约所在源文件，继承Print合约的子类合约可以直接使用print方法打印日志
import "./Print.sol";

//testPrint合约继承Print合约
contract testPrint is Print {

    function test() public view {
        //打印日志级别为debug级别，参数71240和true被toStr方法转换成了string，然后被string.concat和其他参数拼接为一个string
        print(0, string.concat("test debug level: ", "print uint --", toStr(71240), ", print bool --", toStr(true)));

        //info:    1
        print(1, string.concat("test info level: ", "print uint --", toStr(7749), ", print bool --", toStr(false)));

        //warning: 2
        print(2, string.concat("test warning level: ", "print uint --", toStr(9981), ", print bool --", toStr(false)));

        //error:   3
        print(3, string.concat("test error level: ", "print uint --", toStr(10000), ", print bool --", toStr(true)));
    }
}
```



##### 跨合约调用对比

长安链对solidity的跨合约调用做了增量修改，除完全兼容和支持以太坊所支持的跨合约调用外，还支持对（长安链已有的）其他虚拟机合约跨合约调用。具体内容，参见[跨合约调用教程](./如何进行跨合约调用.md)。

### 合约示例源码展示

**Token合约**示例，实现功能ERC20

```
/*
SPDX-License-Identifier: Apache-2.0
*/
pragma solidity >0.5.11;
contract Token {

		//状态变量（注意，如果对合约升级，新合约不得改变原合约的状态变量，包括名称、类型和顺序，但是可以追加新的状态变量）
    string public name = "token";      //  token name
    string public symbol = "TK";           //  token symbol
    uint256 public decimals = 6;            //  token digit

    mapping (address => uint256) public balanceOf;
    mapping (address => mapping (address => uint256)) public allowance;

    uint256 public totalSupply = 0;
    bool public stopped = false;

    uint256 constant valueFounder = 100000000000000000;
    address owner = address(0x0);

    modifier isOwner {
        assert(owner == msg.sender);
        _;
    }

    modifier isRunning {
        assert (!stopped);
        _;
    }

    modifier validAddress {
        assert(address(0x0) != msg.sender);
        _;
    }

    constructor (address _addressFounder) {
        owner = msg.sender;
        totalSupply = valueFounder;
        balanceOf[_addressFounder] = valueFounder;
        
        emit Transfer(address(0x0), _addressFounder, valueFounder);
    }

    function transfer(address _to, uint256 _value) public isRunning validAddress returns (bool success) {
        require(balanceOf[msg.sender] >= _value);
        require(balanceOf[_to] + _value >= balanceOf[_to]);
        balanceOf[msg.sender] -= _value;
        balanceOf[_to] += _value;
        emit Transfer(msg.sender, _to, _value);
        return true;
    }

    function transferFrom(address _from, address _to, uint256 _value) public isRunning validAddress returns (bool success) {
        require(balanceOf[_from] >= _value);
        require(balanceOf[_to] + _value >= balanceOf[_to]);
        require(allowance[_from][msg.sender] >= _value);
        balanceOf[_to] += _value;
        balanceOf[_from] -= _value;
        allowance[_from][msg.sender] -= _value;
        emit Transfer(_from, _to, _value);
        return true;
    }

    function approve(address _spender, uint256 _value) public isRunning validAddress returns (bool success) {
        require(_value == 0 || allowance[msg.sender][_spender] == 0);
        allowance[msg.sender][_spender] = _value;
        emit Approval(msg.sender, _spender, _value);
        return true;
    }

    function stop() public isOwner {
        stopped = true;
    }

    function start() public isOwner {
        stopped = false;
    }

    function setName(string memory _name) public isOwner {
        name = _name;
    }

    function burn(uint256 _value) public {
        require(balanceOf[msg.sender] >= _value);
        balanceOf[msg.sender] -= _value;
        balanceOf[address(0x0)] += _value;
        emit Transfer(msg.sender, address(0x0), _value);
    }

    event Transfer(address indexed _from, address indexed _to, uint256 _value);
    event Approval(address indexed _owner, address indexed _spender, uint256 _value);
}
```

**注意：solidity合约中的状态变量是以槽的方式顺序存储的，所以，如果后续对合约升级的话，新合约不得改变原合约的状态变量，包括变量名称、类型和顺序，但是可以追加（不可插入）新的状态变量。**

### 编译合约

#### 使用Docker镜像搭建编译环境

开发者可使用ChainMaker已经打包好的Docker镜像编译solidity合约代码，ChainMaker官方已经将容器发布至 [docker hub](https://hub.docker.com/u/chainmakerofficial)

**拉取镜像**

请使用以下命令拉取用于编译solidity合约的镜像。

```sh
docker pull chainmakerofficial/chainmaker-solidity-contract:2.0.0
```

**启动镜像**

启动镜像前，需要指定本地开发目录，用于映射为docker镜像的home目录。请指定你本机的工作目录$WORK_DIR，例如/data/workspace/contract，挂载到docker容器中以方便后续进行必要的一些文件拷贝。

```sh
#启动并进入容器，$WORK_DIR即本地工作目录
docker run -it --name chainmaker-solidity-contract -v $WORK_DIR:/home chainmakerofficial/chainmaker-solidity-contract:2.0.0 bash
# 或者先后台启动
docker run -d --name chainmaker-solidity-contract -v $WORK_DIR:/home chainmakerofficial/chainmaker-solidity-contract:2.0.0 bash -c "while true; do echo hello world; sleep 5;done"
# 再进入容器
docker exec -it chainmaker-solidity-contract bash
```

#### 编译示例合约

在本地开发目录内的solidity合约，可以在docker编译镜像的home目录内直接看到，因为二者已映射。进入docker编译镜像后，切换到home目录，执行以下命令，即可编译solidity合约。

```sh
# cd /home/
# solc --abi --bin --hashes --overwrite -o . token.sol
```

solc为编译命令， --abi选项指示生成abi文件，--bin指示生成字节码文件， --hashes指示生成函数签名文件， --overwrite指示如果生成文件已存在则覆盖， -o 指示编译生成的文件存放的目录。solc命令更详细用法，可使用solc --help查看。

生成的字节码位于solc命令用 -o 指定的目录内，示例 中为当前目录：

```
/home/Token.bin
```

#### 编译说明

solc编译命令使用的是0.8.4+commit.c7e474f2.Linux.g++版本的编译器，被编译的solidity合约版本号必须大于等于0.8.4，否则有可能编译告警或报错。

如果开发者不愿意修改solidity合约以适应solc编译器的版本，那么也可以直接使用Remix编译。通过Remix编译出的字节码也可以在长安链上直接部署运行。

#### 调用参数

调用合约前，需要首先将被调方法的签名和调用参数打包成calldata, 将calldata作为参数传递给被调合约，calldata需要通过ABI接口生成。

##### calldata生成示例

长安链项目的common模块提供了abi计算包，开发者导入abi包后，可指定合约方法和参数生成calldata。

```go
import "chainmaker.org/chainmaker/common/v2/evmutils/abi"
...
//读取abi文件
abiBytes, _ := ioutil.ReadFile("xxxx.abi")
//构建abi对象
abiObj, _ := abi.JSON(strings.NewReader(string(abiBytes)))
//计算calldata
calldata, err := abiObj.Pack("methodName", big.NewInt(100))
...
```


##### 



