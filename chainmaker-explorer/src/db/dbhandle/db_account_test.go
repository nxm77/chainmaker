package dbhandle

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"os"
	"testing"
)

const (
	ChainID      = "chainmaker_pk"
	AccountAddr1 = "123456789"
	AccountAddr2 = "223456789"
	AccountAddr3 = "323456789"
	AccountAddr4 = "423456789"
	BNSAddr1     = "bns.com"
	BNSAddr2     = "bns.2com"
	BNSAddr3     = "bns.com2"
	BNSAddr4     = "bns.com3"
	DIDAddr1     = "did:12345"
	DIDAddr2     = "did:22345"
	DIDAddr3     = "did:42345"
	DIDAddr4     = "did:42345"
)

func insertAccountTest() ([]*db.Account, error) {
	accountList := []*db.Account{
		{
			AddrType: 0,
			Address:  AccountAddr1,
			DID:      DIDAddr1,
			BNS:      BNSAddr1,
			NFTNum:   2,
			TxNum:    2,
		},
		{
			AddrType: 0,
			Address:  AccountAddr2,
			DID:      DIDAddr2,
			BNS:      BNSAddr2,
		},
	}
	err := InsertAccount(ChainID, accountList)
	return accountList, err
}

func insertAccountTest2() ([]*db.Account, error) {
	accountList := []*db.Account{
		{
			AddrType: 1,
			Address:  AccountAddr3,
			DID:      DIDAddr3,
			BNS:      BNSAddr3,
		},
		{
			AddrType: 1,
			Address:  AccountAddr4,
			DID:      DIDAddr4,
			BNS:      BNSAddr4,
		},
	}
	err := InsertAccount(ChainID, accountList)
	return accountList, err
}

func TestMain(m *testing.M) {
	_ = config.InitConfig("", "")
	// 初始化数据库配置
	db.InitRedisContainer()
	db.InitMySQLContainer()

	// 运行其他测试
	os.Exit(m.Run())
}

func TestGetAccountByAddr(t *testing.T) {
	_, err := insertAccountTest()
	if err != nil {
		return
	}

	gotWant := &db.Account{
		AddrType: 0,
		Address:  AccountAddr1,
		DID:      DIDAddr1,
		BNS:      BNSAddr1,
	}
	type args struct {
		chainId string
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    *db.Account
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				address: AccountAddr1,
			},
			wantErr: false,
			want:    gotWant,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetAccountByAddr(tt.args.chainId, tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccountByAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetAccountByBNS(t *testing.T) {
	_, err := insertAccountTest()
	if err != nil {
		return
	}

	gotWant := &db.Account{
		AddrType: 0,
		Address:  AccountAddr1,
		DID:      "did:12345",
		BNS:      "bns.com",
	}

	type args struct {
		chainId string
		bns     string
	}
	tests := []struct {
		name    string
		args    args
		want    *db.Account
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				bns:     "bns.com",
			},
			wantErr: false,
			want:    gotWant,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetAccountByBNS(tt.args.chainId, tt.args.bns)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccountByBNS() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetAccountByDID(t *testing.T) {
	_, err := insertAccountTest()
	if err != nil {
		return
	}

	gotWant := []*db.Account{
		{
			AddrType: 0,
			Address:  AccountAddr1,
			DID:      "did:12345",
			BNS:      "bns.com",
		},
	}

	type args struct {
		chainId string
		did     string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.Account
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				did:     "did:12345",
			},
			wantErr: false,
			want:    gotWant,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetAccountByDID(tt.args.chainId, tt.args.did)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccountByDID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetAccountDetail(t *testing.T) {
	_, err := insertAccountTest()
	if err != nil {
		return
	}

	type args struct {
		chainId string
		address string
		bns     string
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
				address: AccountAddr1,
			},
			wantErr: false,
		},
		{
			name: "test: case 2",
			args: args{
				chainId: ChainID,
				address: AccountAddr1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetAccountDetail(tt.args.chainId, tt.args.address, tt.args.bns)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccountDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestInsertAccount(t *testing.T) {
	accountList := []*db.Account{
		{
			AddrType: 1,
			Address:  "12345678933333",
			DID:      "did:12345333",
			BNS:      "bns.com33",
		},
		{
			AddrType: 1,
			Address:  "2234567893333",
			DID:      "did:22345333",
			BNS:      "bns.2com333",
		},
	}
	type args struct {
		chainId     string
		accountList []*db.Account
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:     "chainmaker_pk",
				accountList: accountList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InsertAccount(tt.args.chainId, tt.args.accountList); (err != nil) != tt.wantErr {
				t.Errorf("InsertAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestQueryAccountExists(t *testing.T) {
	addrListDB, err := insertAccountTest()
	if err != nil {
		return
	}

	wantMap := map[string]*db.Account{}
	for _, value := range addrListDB {
		wantMap[value.Address] = value
	}

	type args struct {
		chainId  string
		addrList []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*db.Account
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: "chainmaker_pk",
				addrList: []string{
					AccountAddr1,
					AccountAddr2,
					"343234345",
				},
			},
			wantErr: false,
			want:    wantMap,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := QueryAccountExists(tt.args.chainId, tt.args.addrList)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryAccountExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUpdateAccount(t *testing.T) {
	_, err := insertAccountTest2()
	if err != nil {
		return
	}

	accountInfo := &db.Account{
		AddrType: 1,
		Address:  AccountAddr3,
		DID:      "did:12345",
		BNS:      "bns.com",
	}

	type args struct {
		chainId     string
		accountInfo *db.Account
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:     ChainID,
				accountInfo: accountInfo,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateAccount(tt.args.chainId, tt.args.accountInfo); (err != nil) != tt.wantErr {
				t.Errorf("UpdateAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetAccountByBNSList(t *testing.T) {
	type args struct {
		chainId string
		bnsList []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.Account
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				bnsList: []string{
					"BNS:123",
				},
			},
			want:    make([]*db.Account, 0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetAccountByBNSList(tt.args.chainId, tt.args.bnsList)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccountByBNSList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetAccountByDIDList(t *testing.T) {
	type args struct {
		chainId string
		didList []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.Account
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				didList: []string{
					"BNS:123",
				},
			},
			want:    make([]*db.Account, 0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetAccountByDIDList(tt.args.chainId, tt.args.didList)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAccountByDIDList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetAccountList(t *testing.T) {
	insertAccountTest()
	// Test case 1: Get account list with valid parameters
	_, count1, err1 := GetAccountList(0, 10, db.UTchainID, "", "")
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}
	if count1 == 0 {
		t.Errorf("Test case 1 failed: No accounts found")
	}

	// Test case 2: Get account list with invalid chain ID
	_, _, err2 := GetAccountList(0, 10, db.UTchainID, "", AccountAddr1)
	if err2 != nil {
		t.Errorf("Test case 2 failed: Expected error not returned")
	}
}

func TestGetTotalTxNumByAccount(t *testing.T) {
	insertAccountTest()
	// Test case 1: Get total transaction number by valid chain ID
	totalTxNum1, err1 := GetTotalTxNumByAccount(db.UTchainID)
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}
	if totalTxNum1 == 0 {
		t.Errorf("Test case 1 failed: No transactions found")
	}

	// Test case 2: Get total transaction number by non-existing chain ID
	_, err2 := GetTotalTxNumByAccount("chain999")
	if err2 == nil {
		t.Errorf("Test case 2 failed: Expected error not returned")
	}

	// Test case 3: Get total transaction number with invalid parameters
	_, err3 := GetTotalTxNumByAccount("")
	if err3 == nil {
		t.Errorf("Test case 3 failed: Expected error not returned")
	}
}

func TestGetAccountTotal(t *testing.T) {
	insertAccountTest()
	// Test case 1: Get account total with valid parameters
	_, err1 := GetAccountTotal(db.UTchainID)
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}
}
