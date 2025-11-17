package logic

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/sync/model"
	"testing"

	"chainmaker.org/chainmaker/pb-go/v2/common"
)

const ChainId = "testChainId"

func TestParallelParseTransactions(t *testing.T) {
	dealResult := &model.ProcessedBlockData{
		UserList:     map[string]*db.User{},
		Transactions: map[string]*db.Transaction{},
	}
	blockInfo := getBlockInfoTest("2_blockInfoJsonContractERC20.json")
	if blockInfo == nil {
		return
	}

	type args struct {
		blockInfo  *common.BlockInfo
		hashType   string
		dealResult *model.ProcessedBlockData
	}
	tests := []struct {
		name    string
		args    args
		want    *model.ProcessedBlockData
		wantErr bool
	}{
		{
			name: "test case 1",
			args: args{
				blockInfo:  blockInfo,
				hashType:   "1222222",
				dealResult: dealResult,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ParallelParseTransactions(tt.args.blockInfo, tt.args.hashType, tt.args.dealResult)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParallelParseTransactions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_buildReadWriteSet(t *testing.T) {
	blockInfo := getBlockInfoTest("2_blockInfoJson_1704945589.json")
	if blockInfo == nil || len(blockInfo.RwsetList) == 0 {
		return
	}
	type args struct {
		rwsetList common.TxRWSet
	}

	type testStruct struct {
		name string
		args args
	}

	var tests []testStruct
	temp := testStruct{
		name: "test case 1",
		args: args{
			rwsetList: *blockInfo.RwsetList[0],
		},
	}
	tests = append(tests, temp)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buildReadWriteSet(tt.args.rwsetList)
		})
	}
}

func Test_buildTransaction(t *testing.T) {
	blockInfo := getBlockInfoTest("2_blockInfoJsonContractERC20.json")
	if blockInfo == nil {
		return
	}
	txInfo := getTxInfoInfoTest("2_txInfoJsonContractERC20.json")
	if txInfo == nil {
		return
	}
	buildTxResult := getBuildTxInfoTest("2_buildTxResult.json")
	if buildTxResult == nil {
		return
	}
	type args struct {
		i          int
		blockInfo  *common.BlockInfo
		txInfo     *common.Transaction
		userResult *db.SenderPayerUser
	}
	tests := []struct {
		name    string
		args    args
		want    *db.Transaction
		wantErr bool
	}{
		{
			name: "test case 1",
			args: args{
				i:          0,
				blockInfo:  blockInfo,
				txInfo:     txInfo,
				userResult: nil,
			},
			want:    buildTxResult,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := buildTransaction(tt.args.i, tt.args.blockInfo, tt.args.txInfo, tt.args.userResult, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("buildTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("buildTransaction() got = %v, want %v", got, tt.want)
			// }
		})
	}
}
