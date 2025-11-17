package dbhandle

import (
	"chainmaker_web/src/db"
	"testing"
)

const (
	UserId1   = "123456"
	UserAddr1 = "123456789"
)

func insertUserTest1() ([]*db.User, error) {
	insertList := []*db.User{
		{
			UserId:    UserId1,
			UserAddr:  UserAddr1,
			Timestamp: 123456,
		},
	}
	err := BatchInsertUser(db.UTchainID, insertList)
	return insertList, err
}

func TestBatchInsertUser(t *testing.T) {
	insertList := []*db.User{
		{
			UserId:    UserId1,
			UserAddr:  UserAddr1,
			Timestamp: 123456,
		},
	}

	type args struct {
		chainId  string
		userList []*db.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:  db.UTchainID,
				userList: insertList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := BatchInsertUser(tt.args.chainId, tt.args.userList); (err != nil) != tt.wantErr {
				t.Errorf("BatchInsertUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetUserList(t *testing.T) {
	_, err := insertUserTest1()
	if err != nil {
		return
	}
	type args struct {
		offset    int
		limit     int
		chainId   string
		orgId     string
		userIds   []string
		userAddrs []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.User
		want1   int64
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: db.UTchainID,
				offset:  0,
				limit:   10,
			},
			want1:   1,
			wantErr: false,
		},
		{
			name: "test: case 2",
			args: args{
				chainId:   db.UTchainID,
				offset:    0,
				limit:     10,
				orgId:     "org1",
				userIds:   []string{UserId1},
				userAddrs: []string{UserAddr1},
			},
			want1:   1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := GetUserList(tt.args.offset, tt.args.limit, tt.args.chainId, tt.args.orgId, tt.args.userIds, tt.args.userAddrs)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetUserListByAdder(t *testing.T) {
	insertList, err := insertUserTest1()
	if err != nil {
		return
	}

	type args struct {
		chainId string
		adders  []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*db.User
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: db.UTchainID,
				adders: []string{
					UserAddr1,
				},
			},
			want: map[string]*db.User{
				insertList[0].UserAddr: insertList[0],
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetUserListByAdder(tt.args.chainId, tt.args.adders)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserListByAdder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetUserNum(t *testing.T) {
	type args struct {
		chainId string
		orgId   string
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
				chainId: db.UTchainID,
				orgId:   "1",
			},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUserNum(tt.args.chainId, tt.args.orgId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUserNum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetUserNum() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUpdateUserStatus(t *testing.T) {
	_, err := insertUserTest1()
	if err != nil {
		return
	}

	type args struct {
		address string
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
				chainId: db.UTchainID,
				address: UserAddr1,
				status:  1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateUserStatus(tt.args.address, tt.args.chainId, tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("UpdateUserStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetUserCountByRange(t *testing.T) {
	_, err := insertUserTest1()
	if err != nil {
		return
	}

	type args struct {
		chainId   string
		startTime int64
		endTime   int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:   db.UTchainID,
				startTime: 123456789,
				endTime:   123456790,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := GetUserCountByRange(tt.args.chainId, tt.args.startTime, tt.args.endTime); (err != nil) != tt.wantErr {
				t.Errorf("GetUserCountByRange() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
