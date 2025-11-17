package dbhandle

import (
	"chainmaker_web/src/db"
	"reflect"
	"testing"
)

func TestGetAllSubChainBlockHeight(t *testing.T) {
	type args struct {
		chainId string
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetAllSubChainBlockHeight(tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllSubChainBlockHeight() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetCrossLatestSubChainList(t *testing.T) {
	type args struct {
		chainId string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.CrossSubChainData
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
			},
			want:    make([]*db.CrossSubChainData, 0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetCrossLatestSubChainList(tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossLatestSubChainList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetCrossSubChainAll(t *testing.T) {
	type args struct {
		chainId string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.CrossSubChainData
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
			},
			want:    make([]*db.CrossSubChainData, 0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetCrossSubChainAll(tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossSubChainAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetCrossSubChainAllCount(t *testing.T) {
	type args struct {
		chainId string
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetCrossSubChainAllCount(tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossSubChainAllCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetCrossSubChainById(t *testing.T) {
	type args struct {
		chainId     string
		subChainIds []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*db.CrossSubChainData
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				subChainIds: []string{
					"123",
				},
			},
			want:    make(map[string]*db.CrossSubChainData, 0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCrossSubChainById(tt.args.chainId, tt.args.subChainIds)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossSubChainById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCrossSubChainById() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCrossSubChainInfoById(t *testing.T) {
	type args struct {
		chainId    string
		subChainId string
	}
	tests := []struct {
		name    string
		args    args
		want    *db.CrossSubChainData
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:    ChainID,
				subChainId: "123",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetCrossSubChainInfoById(tt.args.chainId, tt.args.subChainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossSubChainInfoById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetCrossSubChainInfoCache(t *testing.T) {
	type args struct {
		chainId    string
		subChainId string
	}
	tests := []struct {
		name string
		args args
		want *db.CrossSubChainData
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:    ChainID,
				subChainId: "123",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetCrossSubChainInfoCache(tt.args.chainId, tt.args.subChainId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCrossSubChainInfoCache() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetCrossSubChainList(t *testing.T) {
	type args struct {
		offset     int
		limit      int
		chainId    string
		subChainId string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.CrossSubChainData
		want1   int64
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:    ChainID,
				offset:     0,
				limit:      10,
				subChainId: "123",
			},
			want: make([]*db.CrossSubChainData, 0),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := GetCrossSubChainList(tt.args.offset, tt.args.limit, tt.args.chainId, tt.args.subChainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossSubChainList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetCrossSubChainListCache(t *testing.T) {
	type args struct {
		chainId string
	}
	tests := []struct {
		name string
		args args
		want []*db.CrossSubChainData
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
			_ = GetCrossSubChainListCache(tt.args.chainId)
		})
	}
}

func TestInsertCrossSubChain(t *testing.T) {
	insertList := []*db.CrossSubChainData{
		{
			SubChainId: "6666",
			TxNum:      23,
		},
	}
	type args struct {
		chainId      string
		subChainList []*db.CrossSubChainData
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
				subChainList: insertList,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InsertCrossSubChain(tt.args.chainId, tt.args.subChainList); (err != nil) != tt.wantErr {
				t.Errorf("InsertCrossSubChain() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetCrossSubChainListCache(t *testing.T) {
	type args struct {
		chainId      string
		subChainList []*db.CrossSubChainData
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				subChainList: []*db.CrossSubChainData{
					{
						SubChainId: "123",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetCrossSubChainListCache(tt.args.chainId, tt.args.subChainList)
		})
	}
}

func TestSetCrossSubChainNameCache(t *testing.T) {
	type args struct {
		chainId      string
		subChainId   string
		subChainName string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:      ChainID,
				subChainId:   "123",
				subChainName: "123",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetCrossSubChainNameCache(tt.args.chainId, tt.args.subChainId, tt.args.subChainName)
		})
	}
}

func TestUpdateCrossSubChainById(t *testing.T) {
	type args struct {
		chainId      string
		subChainInfo *db.CrossSubChainData
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
				subChainInfo: &db.CrossSubChainData{
					SubChainId: "123",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateCrossSubChainById(tt.args.chainId, tt.args.subChainInfo); (err != nil) != tt.wantErr {
				t.Errorf("UpdateCrossSubChainById() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// func TestUpdateCrossSubChainStatus(t *testing.T) {
// 	type args struct {
// 		chainId    string
// 		subChainId string
// 		status     int32
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "test: case 1",
// 			args: args{
// 				chainId:    ChainID,
// 				subChainId: "123",
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if err := UpdateCrossSubChainStatus(tt.args.chainId, tt.args.subChainId, "", tt.args.status); (err != nil) != tt.wantErr {
// 				t.Errorf("UpdateCrossSubChainStatus() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }
