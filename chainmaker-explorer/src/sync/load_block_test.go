package sync

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	client "chainmaker_web/src/sync/clients"
	"chainmaker_web/src/sync/common"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"

	pbCommon "chainmaker.org/chainmaker/pb-go/v2/common"
)

func TestMain(m *testing.M) {
	_ = config.InitConfig("", "")
	// 初始化数据库配置
	db.InitRedisContainer()
	db.InitMySQLContainer()

	for _, chainInfo := range config.SubscribeChains {
		// 创建区块链连接，并加入连接池
		_, _ = CreateSubscribeClientPool(chainInfo)
	}

	// 运行其他测试
	os.Exit(m.Run())
}

func geFileData(fileName string) (*json.Decoder, error) {
	file, err := os.Open("./testData/" + fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}

	// 解码 JSON 文件内容到 blockInfo 结构体
	decoder := json.NewDecoder(file)
	return decoder, err
}

func getBlockInfoTest(fileName string) *pbCommon.BlockInfo {
	var blockInfo *pbCommon.BlockInfo
	// 打开 JSON 文件
	//file, err := os.Open("../testData/1_blockInfoJsonContract.json")

	decoder, err := geFileData(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return blockInfo
	}

	// 解码 JSON 文件内容到 blockInfo 结构体
	err = decoder.Decode(&blockInfo)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return blockInfo
	}
	return blockInfo
}

func TestBuildChainInfo(t *testing.T) {
	nodeList := []*config.NodeInfo{
		{
			Addr:        "123456",
			OrgCA:       "12345678",
			TLSHostName: "1234569",
		},
	}

	nodeJson, _ := json.Marshal(nodeList)
	subscribeChain := &db.Subscribe{
		ChainId:     "chain1",
		OrgId:       "123",
		UserSignKey: "345",
		UserSignCrt: "567",
		AuthType:    "567",
		HashType:    "5678",
		Tls:         true,
		NodeList:    string(nodeJson),
		Status:      0,
	}

	type args struct {
		subscribeChain *db.Subscribe
	}
	tests := []struct {
		name string
		args args
		want *config.ChainInfo
	}{
		{
			name: "test case 1",
			args: args{
				subscribeChain: subscribeChain,
			},
			want: &config.ChainInfo{
				ChainId:   "chain1",
				AuthType:  "567",
				OrgId:     "123",
				HashType:  "5678",
				NodesList: nodeList,
				Tls:       true,
				UserInfo: &config.UserInfo{
					UserSignKey: "345",
					UserSignCrt: "567",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildChainInfo(tt.args.subscribeChain); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BuildChainInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDelayUpdateOperation(t *testing.T) {
	ctx := context.Background()
	//blockWaitUpdateCh 创建一个容量为 20 的通道来存储处理完成，等待异步更新的区块数据
	blockWaitUpdateCh := make(chan *BlockWaitUpdate, config.BlockWaitUpdateWorkerCount)
	//blockListenErrCh 创建一个错误通道来接收 子线程 的错误
	blockListenErrCh := make(chan error)
	// 将处理完成的结果写入 blockWaitUpdateCh
	resultData := &BlockWaitUpdate{
		ChainId:     db.UTchainID,
		BlockHeight: 10,
	}
	blockWaitUpdateCh <- resultData

	type args struct {
		blockWaitUpdateCh chan *BlockWaitUpdate
		errCh             chan<- error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test case 1",
			args: args{
				blockWaitUpdateCh: blockWaitUpdateCh,
				errCh:             blockListenErrCh,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go DelayUpdateOperation(ctx, resultData.ChainId, tt.args.blockWaitUpdateCh, tt.args.errCh)
			// 使用 select 语句等待错误或上下文取消
			select {
			case errCh := <-blockListenErrCh:
				t.Errorf("DelayUpdateOperation() error = %v, wantErr %v", errCh, errCh)
				close(blockWaitUpdateCh)
				close(blockListenErrCh)
			case <-time.After(2 * time.Second):
				// 等待5秒后关闭通道
				close(blockWaitUpdateCh)
				close(blockListenErrCh)
			}
		})
	}
}

func TestParallelParseBlockWork(t *testing.T) {
	blockInfo := getBlockInfoTest("2_blockInfoJson_1704945589.json")
	if blockInfo == nil || len(blockInfo.RwsetList) == 0 {
		return
	}
	chainId := blockInfo.Block.Header.ChainId
	//nolint:gosec
	blockHeight := int64(blockInfo.Block.Header.BlockHeight)
	common.SetMaxHeight(chainId, blockHeight)

	hashType := "your_hash_type"

	// 创建通道和 WaitGroup
	blockInfoCh := make(chan *pbCommon.BlockInfo, 1)
	dataSaveCh := make(chan *DataSaveToDB, 1)
	errCh := make(chan error, 1)
	wg := &sync.WaitGroup{}

	// 创建一个带有超时的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	type args struct {
		ctx         context.Context
		wg          *sync.WaitGroup
		hashType    string
		blockInfoCh chan *pbCommon.BlockInfo
		dataSaveCh  chan *DataSaveToDB
		errCh       chan<- error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test case 1",
			args: args{
				ctx:         ctx,
				wg:          wg,
				hashType:    hashType,
				blockInfoCh: blockInfoCh,
				dataSaveCh:  dataSaveCh,
				errCh:       errCh,
			},
		},
	}
	go func() {
		time.Sleep(10 * time.Second)
		cancel()
		close(blockInfoCh)
		close(dataSaveCh)
	}()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 在一个新的 Goroutine 中运行 ParallelParseBlockWork
			wg.Add(1)
			go ParallelParseBlockWork(tt.args.ctx, tt.args.wg, tt.args.hashType, tt.args.blockInfoCh, tt.args.dataSaveCh, tt.args.errCh)
			// 发送测试数据到 blockInfoCh
			blockInfoCh <- blockInfo
			// 检查结果
			select {
			case err := <-errCh:
				t.Errorf("ParallelParseBlockWork() error = %v", err)
			case dataSave := <-dataSaveCh:
				if dataSave == nil {
					// 验证 dataSave 的内容
					t.Errorf("ParallelParseBlockWork() dataSave = %v", dataSave)
				}
			}
		})
	}
}

func Test_startSubscribeLockTicker(t *testing.T) {
	type args struct {
		sdkClient *client.SdkClient
		lockKey   string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test case 1",
			args: args{
				sdkClient: &client.SdkClient{
					ChainId: ChainId1,
				},
				lockKey: "1234",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			go startSubscribeLockTicker(ctx, tt.args.sdkClient, tt.args.lockKey)
			time.Sleep(1000)
			// 在测试结束时取消上下文
			defer cancel()
		})
	}
}

func TestPeriodicGetSubscribeLock(t *testing.T) {
	// Test case 1: Test with valid sdkClient
	sdkClient1 := &client.SdkClient{
		ChainId: "chain1",
		Ctx:     context.Background(),
	}
	go PeriodicGetSubscribeLock(sdkClient1)
}

func TestStartSubscribeLockTicker(t *testing.T) {
	// Test case 1: Test with valid sdkClient and lockKey
	sdkClient1 := &client.SdkClient{
		ChainId: ChainId1,
		Ctx:     context.Background(),
	}
	lockKey1 := "lockKey1"
	go startSubscribeLockTicker(sdkClient1.Ctx, sdkClient1, lockKey1)
	time.Sleep(2 * time.Second)

	sdkClient1.Status = client.STOP
	go startSubscribeLockTicker(sdkClient1.Ctx, sdkClient1, lockKey1)
}

func TestBlockListen(t *testing.T) {
	// Test case 1: Test with valid sdkClient
	sdkClient1 := &client.SdkClient{
		ChainId: ChainId1,
		Ctx:     context.Background(),
	}
	go blockListen(sdkClient1)
}

func TestOrderedSaveBlockData(t *testing.T) {
	// Test case 1: Test with valid dataSaveCh and blockWaitUpdateCh
	ctx1 := context.Background()
	dataSaveCh1 := make(chan *DataSaveToDB)
	blockWaitUpdateCh1 := make(chan *BlockWaitUpdate)
	errCh1 := make(chan error)
	go OrderedSaveBlockData(ctx1, dataSaveCh1, blockWaitUpdateCh1, errCh1)
}

func TestWaitUpdateChFailedData(t *testing.T) {
	// Test case 1: Test with valid chainId and blockWaitUpdateCh
	ctx1 := context.Background()

	blockWaitUpdateCh1 := make(chan *BlockWaitUpdate)
	go waitUpdateChFailedData(ctx1, ChainId1, blockWaitUpdateCh1)

	// Test case 2: Test with nil chainId and blockWaitUpdateCh
	ctx2 := context.Background()
	chainId2 := ""
	blockWaitUpdateCh2 := make(chan *BlockWaitUpdate)
	go waitUpdateChFailedData(ctx2, chainId2, blockWaitUpdateCh2)
}

func TestRealtimeInsertOperation(t *testing.T) {
	// Test case 1: Test with valid ctx, hash, blockInfoCh, dataSaveCh and errCh
	ctx1 := context.Background()
	hash1 := "hash1"
	blockInfoCh1 := make(chan *pbCommon.BlockInfo)
	dataSaveCh1 := make(chan *DataSaveToDB)
	errCh1 := make(chan error)
	go RealtimeInsertOperation(ctx1, hash1, blockInfoCh1, dataSaveCh1, errCh1)
}

func TestParallelParseBlockWork1(t *testing.T) {
	// Test case 1: Test with valid ctx, wg, hashType, blockInfoCh, dataSaveCh and errCh
	ctx1 := context.Background()
	wg1 := &sync.WaitGroup{}
	hashType1 := "hashType1"
	blockInfoCh1 := make(chan *pbCommon.BlockInfo)
	dataSaveCh1 := make(chan *DataSaveToDB)
	errCh1 := make(chan error)
	go ParallelParseBlockWork(ctx1, wg1, hashType1, blockInfoCh1, dataSaveCh1, errCh1)
}

// func TestDelayUpdateOperation1(t *testing.T) {
// 	// Test case 1: Test with valid blockWaitUpdateCh and errCh
// 	ctx1 := context.Background()
// 	blockWaitUpdateCh1 := make(chan *BlockWaitUpdate)
// 	errCh1 := make(chan error)
// 	go DelayUpdateOperation(ctx1, blockWaitUpdateCh1, errCh1)
// }

func TestBuildChainInfo1(t *testing.T) {
	// Test case 1: Test with valid subscribeChain
	subscribeChain1 := &db.Subscribe{
		ChainId:     ChainId1,
		AuthType:    "authType1",
		OrgId:       "orgId1",
		HashType:    "hashType1",
		TlsMode:     0,
		Tls:         true,
		UserSignKey: "userSignKey1",
		UserSignCrt: "userSignCrt1",
		UserEncKey:  "userEncKey1",
		UserEncCrt:  "userEncCrt1",
		NodeList:    `[{"nodeId":"nodeId1","orgId":"orgId1","nodeAddr":"nodeAddr1","clientKey":"clientKey1","clientCert":"clientCert1","tlsCert":"tlsCert1"}]`,
	}
	BuildChainInfo(subscribeChain1)
}
