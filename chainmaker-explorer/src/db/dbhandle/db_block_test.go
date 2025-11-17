package dbhandle

import (
	"chainmaker_web/src/db"
	"strconv"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

const (
	BlockHeight1 int64 = 1
	BlockHeight2 int64 = 2
	BlockHeight3 int64 = 3
	BlockHash1         = "q123ww45"
	BlockHash2         = "q22qqq345"
	BlockHash3         = "q323eee45"
)

func insertBlockTest1() (*db.Block, error) {
	newUUID := uuid.New().String()
	blockInfo := &db.Block{
		ID:                newUUID,
		BlockHeight:       BlockHeight1,
		PreBlockHash:      "12345",
		BlockHash:         BlockHash1,
		DelayUpdateStatus: 1,
		Timestamp:         12345,
		TimestampDate:     time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
	}
	err := InsertBlock(ChainID, blockInfo)
	return blockInfo, err
}

func insertBlockTest2() (*db.Block, error) {
	newUUID := uuid.New().String()
	blockInfo := &db.Block{
		ID:            newUUID,
		BlockHeight:   BlockHeight2,
		PreBlockHash:  "12345",
		BlockHash:     BlockHash2,
		Timestamp:     123456,
		TimestampDate: time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
	}
	err := InsertBlock(ChainID, blockInfo)
	return blockInfo, err
}

func TestGetBlockByHash(t *testing.T) {
	blockInfo, err := insertBlockTest1()
	if err != nil {
		return
	}

	type args struct {
		blockHash string
		chainId   string
	}
	tests := []struct {
		name    string
		args    args
		want    *db.Block
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:   ChainID,
				blockHash: BlockHash1,
			},
			want:    blockInfo,
			wantErr: false,
		},
		{
			name: "test: case 2",
			args: args{
				chainId:   ChainID,
				blockHash: BlockHash3,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetBlockByHash(tt.args.blockHash, tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBlockByHash() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetBlockByHeight(t *testing.T) {
	blockInfo, err := insertBlockTest1()
	if err != nil {
		t.Fatal(err)
	}

	type args struct {
		chainId     string
		blockHeight int64
	}
	tests := []struct {
		name    string
		args    args
		want    *db.Block
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:     ChainID,
				blockHeight: BlockHeight1,
			},
			want:    blockInfo,
			wantErr: false,
		},
		{
			name: "test: case 2",
			args: args{
				chainId:     ChainID,
				blockHeight: BlockHeight3,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetBlockByHeight(tt.args.chainId, tt.args.blockHeight)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBlockByHeight() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetBlockByStatus(t *testing.T) {
	blockInfo, err := insertBlockTest1()
	if err != nil {
		return
	}

	type args struct {
		chainId string
		status  int
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.Block
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				status:  1,
			},
			want: []*db.Block{
				blockInfo,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetBlockByStatus(tt.args.chainId, tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBlockByStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetBlockList(t *testing.T) {
	blockInfo, err := insertBlockTest1()
	if err != nil {
		return
	}

	type args struct {
		offset   int
		limit    int
		chainId  string
		blockKey string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.Block
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:  ChainID,
				offset:   0,
				limit:    10,
				blockKey: BlockHash1,
			},
			want: []*db.Block{
				blockInfo,
			},
			wantErr: false,
		},
		{
			name: "test: case 1",
			args: args{
				chainId:  ChainID,
				offset:   0,
				limit:    10,
				blockKey: strconv.FormatInt(BlockHeight1, 10),
			},
			want: []*db.Block{
				blockInfo,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBlockList(tt.args.offset, tt.args.limit, tt.args.chainId, tt.args.blockKey, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBlockList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Block{}, "CreatedAt", "UpdatedAt", "ID")) {
				t.Errorf("GetLatestBlockList() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestGetLatestBlockList(t *testing.T) {
	blockInfo, err := insertBlockTest1()
	if err != nil {
		return
	}
	blockInfo2, err := insertBlockTest2()
	if err != nil {
		return
	}

	type args struct {
		chainId string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.Block
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
			},
			want: []*db.Block{
				blockInfo2,
				blockInfo,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetLatestBlockList(tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLatestBlockList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetMaxBlockHeight(t *testing.T) {
	_, err := insertBlockTest1()
	if err != nil {
		return
	}
	_, err = insertBlockTest2()
	if err != nil {
		return
	}

	type args struct {
		chainId string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetMaxBlockHeight(tt.args.chainId)
			if got == 0 {
				t.Errorf("GetMaxBlockHeight() = %v", got)
			}
		})
	}
}

func TestInsertBlock(t *testing.T) {
	blockInfo := &db.Block{
		BlockHeight:       BlockHeight1,
		PreBlockHash:      "12345",
		BlockHash:         BlockHash1,
		DelayUpdateStatus: 1,
		TimestampDate:     time.Date(2023, 10, 1, 0, 0, 0, 0, time.UTC),
	}
	type args struct {
		chainId   string
		blockInfo *db.Block
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:   ChainID,
				blockInfo: blockInfo,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InsertBlock(tt.args.chainId, tt.args.blockInfo); (err != nil) != tt.wantErr {
				t.Errorf("InsertBlock() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateBlockUpdateStatus(t *testing.T) {
	_, err := insertBlockTest1()
	if err != nil {
		return
	}

	type args struct {
		chainId      string
		blockHeight  int64
		updateStatus int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:      ChainID,
				blockHeight:  BlockHeight1,
				updateStatus: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateBlockUpdateStatus(tt.args.chainId, tt.args.blockHeight, tt.args.updateStatus); (err != nil) != tt.wantErr {
				t.Errorf("UpdateBlockUpdateStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetBlockListByRange(t *testing.T) {
	type args struct {
		chainId   string
		startTime int64
		endTime   int64
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
			},
			wantErr: false,
		},
		{
			name: "test: case 1",
			args: args{
				chainId:   ChainID,
				startTime: 12345678,
				endTime:   22345678,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetBlockListByRange(tt.args.chainId, tt.args.startTime, tt.args.endTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBlockListByRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetBlockListCount(t *testing.T) {
	type args struct {
		chainId  string
		blockKey string
		nodeId   string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:  ChainID,
				blockKey: "chainame",
			},
			wantErr: false,
		},
		{
			name: "test: case 2",
			args: args{
				chainId:  "",
				blockKey: "chainame",
			},
			wantErr: true,
		},
		{
			name: "test: case 3",
			args: args{
				chainId:  ChainID,
				blockKey: "123",
			},
			wantErr: false,
		},
		{
			name: "test: case 4",
			args: args{
				chainId:  ChainID,
				blockKey: "123",
				nodeId:   "node1",
			},
			wantErr: false,
		},
		{
			name: "test: case 5",
			args: args{
				chainId: ChainID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetBlockListCount(tt.args.chainId, tt.args.blockKey, tt.args.nodeId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBlockListCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
