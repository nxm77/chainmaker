package sync

import (
	"chainmaker_web/src/sync/model"
	"testing"

	pbCommon "chainmaker.org/chainmaker/pb-go/v2/common"
)

func Test_executeDataInsertTasks(t *testing.T) {
	type args struct {
		chainId    string
		dealResult model.ProcessedBlockData
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test case 1",
			args: args{
				chainId: ChainId1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := executeDataInsertTasks(tt.args.chainId, tt.args.dealResult); (err != nil) != tt.wantErr {
				t.Errorf("executeDataInsertTasks() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestProcessedBlockInfo_ProcessedBlockHandle(t *testing.T) {
	blockInfo := getBlockInfoTest("1_blockInfoJsonContract.json")
	type args struct {
		BlockInfo *pbCommon.BlockInfo
		HashType  string
	}
	tests := []struct {
		name    string
		args    args
		want    *model.ProcessedBlockData
		wantErr bool
	}{
		{
			name: "Test case 1",
			args: args{
				BlockInfo: blockInfo,
				HashType:  "123",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block := &ProcessedBlockInfo{
				BlockInfo: tt.args.BlockInfo,
				HashType:  tt.args.HashType,
			}
			_, err := block.ProcessedBlockHandle()
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessedBlockInfo.ProcessedBlockHandle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
