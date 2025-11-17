package dbhandle

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func insertSubscribeTest() (*db.Subscribe, error) {
	subscribeInfo := &db.Subscribe{
		ChainId:     ChainID,
		UserSignKey: "1234",
		UserSignCrt: "1234",
	}

	err := InsertSubscribe(subscribeInfo)
	return subscribeInfo, err
}

func insertSubscribeTest2() (*db.Subscribe, error) {
	subscribeInfo := &db.Subscribe{
		ChainId:     ChainId2,
		UserSignKey: "1234",
		UserSignCrt: "1234",
	}

	err := InsertSubscribe(subscribeInfo)
	return subscribeInfo, err
}

func TestDeleteSubscribe(t *testing.T) {
	_, err := insertSubscribeTest()
	if err != nil {
		return
	}

	type args struct {
		chainId string
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
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteSubscribe(tt.args.chainId); (err != nil) != tt.wantErr {
				t.Errorf("DeleteSubscribe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetDBSubscribeChains(t *testing.T) {
	subscribeInfo, err := insertSubscribeTest()
	if err != nil {
		return
	}

	tests := []struct {
		name    string
		want    []*db.Subscribe
		wantErr bool
	}{
		{
			name: "test: case 1",
			want: []*db.Subscribe{
				subscribeInfo,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetDBSubscribeChains()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDBSubscribeChains() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetSubscribeByChainId(t *testing.T) {
	subscribeInfo, err := insertSubscribeTest()
	if err != nil {
		return
	}

	type args struct {
		chainId string
	}
	tests := []struct {
		name    string
		args    args
		want    *db.Subscribe
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
			},
			want:    subscribeInfo,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSubscribeByChainId(tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSubscribeByChainId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !cmp.Equal(got, tt.want, cmpopts.IgnoreFields(db.Subscribe{}, "CreatedAt", "UpdatedAt")) {
				t.Errorf("GetSubscribeByChainId() got = %v, want %v\ndiff: %s", got, tt.want, cmp.Diff(got, tt.want))
			}
		})
	}
}

func TestInsertOrUpdateSubscribe(t *testing.T) {
	type args struct {
		chainInfo *config.ChainInfo
		status    int
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainInfo: &config.ChainInfo{
					ChainId:  ChainId3,
					AuthType: "1234",
					HashType: "1234",
				},
			},
			wantErr: false,
		},
		{
			name: "test: case 2",
			args: args{
				chainInfo: &config.ChainInfo{
					ChainId:  ChainId3,
					AuthType: "1232",
					HashType: "12343",
					UserInfo: &config.UserInfo{
						UserSignKey: "1232324",
						UserSignCrt: "123234",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InsertOrUpdateSubscribe(tt.args.chainInfo, tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("InsertOrUpdateSubscribe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInsertSubscribe(t *testing.T) {
	subscribeInfo := &db.Subscribe{
		ChainId:     ChainID,
		UserSignKey: "1234",
		UserSignCrt: "1234",
	}
	type args struct {
		subscribe *db.Subscribe
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				subscribe: subscribeInfo,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InsertSubscribe(tt.args.subscribe); (err != nil) != tt.wantErr {
				t.Errorf("InsertSubscribe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetSubscribeStatus(t *testing.T) {
	subscribeInfo, err := insertSubscribeTest2()
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
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: subscribeInfo.ChainId,
				status:  1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetSubscribeStatus(tt.args.chainId, tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("SetSubscribeStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateSubscribe(t *testing.T) {
	subscribeInfo, err := insertSubscribeTest2()
	if err != nil {
		return
	}

	type args struct {
		subscribe *db.Subscribe
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				subscribe: subscribeInfo,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateSubscribe(tt.args.subscribe); (err != nil) != tt.wantErr {
				t.Errorf("UpdateSubscribe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetSubscribeByChainIds(t *testing.T) {
	_, err := GetSubscribeByChainIds([]string{db.UTchainID})
	if err != nil {
		t.Errorf("Test case 1 failed: %v", err)
	}
}
