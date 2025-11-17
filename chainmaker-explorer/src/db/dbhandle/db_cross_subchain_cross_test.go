package dbhandle

import (
	"chainmaker_web/src/db"
	"testing"
)

func TestGetCrossSubChainCrossNum(t *testing.T) {
	type args struct {
		chainId     string
		subChainIds []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.CrossSubChainCrossChain
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				subChainIds: []string{
					"1234",
				},
			},
			want: make([]*db.CrossSubChainCrossChain, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetCrossSubChainCrossNum(tt.args.chainId, tt.args.subChainIds)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossSubChainCrossNum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetSubChainCrossChainList(t *testing.T) {
	type args struct {
		chainId    string
		subChainId string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.CrossSubChainCrossChain
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:    ChainID,
				subChainId: "1234",
			},
			want: make([]*db.CrossSubChainCrossChain, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetSubChainCrossChainList(tt.args.chainId, tt.args.subChainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSubChainCrossChainList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestInsertCrossSubChainCross(t *testing.T) {
	crossChainList := []*db.CrossSubChainCrossChain{
		{
			ID:         "1234",
			SubChainId: "1234",
		},
	}
	SetCrossSubChainCrossCache(db.UTchainID, "1234", crossChainList)
	insertList := []*db.CrossSubChainCrossChain{
		{
			ID:         "1234",
			SubChainId: "1234",
		},
	}
	type args struct {
		chainId    string
		insertList []*db.CrossSubChainCrossChain
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:    ChainID,
				insertList: insertList,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InsertCrossSubChainCross(tt.args.chainId, tt.args.insertList); (err != nil) != tt.wantErr {
				t.Errorf("InsertCrossSubChainCross() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateCrossSubChainCross(t *testing.T) {
	type args struct {
		chainId       string
		subChainCross *db.CrossSubChainCrossChain
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				subChainCross: &db.CrossSubChainCrossChain{
					SubChainId: "123",
					TxNum:      23,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateCrossSubChainCross(tt.args.chainId, tt.args.subChainCross); (err != nil) != tt.wantErr {
				t.Errorf("UpdateCrossSubChainCross() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetCrossSubChainCrossCache(t *testing.T) {
	// Test case 1: Getting contract byte code by transaction id successfully
	crossChainList := []*db.CrossSubChainCrossChain{
		{
			ID:         "1234",
			SubChainId: "1234",
		},
	}
	SetCrossSubChainCrossCache(db.UTchainID, "22", crossChainList)
}
