package dbhandle

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"fmt"
	"testing"
	"time"
)

const (
	contractName1 = "contract1"
	contractName2 = "contract2"
	//contractName3  = "contract3"
	contractAdder1 = "12345678"
	contractAdder2 = "22345678"
	//contractAdder3 = "32345678"
)

func insertContractTest1() (*db.Contract, error) {
	contractInfo := &db.Contract{
		Name:         contractName1,
		NameBak:      contractName1,
		Addr:         contractAdder1,
		RuntimeType:  "DOCKER_GO",
		ContractType: "CMDFA",
		TxNum:        12,
		Timestamp:    12345,
	}
	err := InsertContract(ChainID, contractInfo)
	return contractInfo, err
}

func insertContractTest2() (*db.Contract, error) {
	contractInfo := &db.Contract{
		Name:           contractName2,
		NameBak:        contractName2,
		Addr:           contractAdder2,
		RuntimeType:    "EVM",
		ContractType:   "ERC20",
		ContractStatus: 0,
		Timestamp:      123456,
	}
	err := InsertContract(ChainID, contractInfo)
	return contractInfo, err
}

func TestGetContractByAddersOrNames(t *testing.T) {
	contractInfo1, err := insertContractTest1()
	if err != nil {
		return
	}
	contractInfo2, err := insertContractTest2()
	if err != nil {
		return
	}

	type args struct {
		chainId  string
		nameList []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]*db.Contract
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				nameList: []string{
					contractName1,
					contractAdder2,
				},
			},
			want: map[string]*db.Contract{
				contractAdder2: contractInfo2,
				contractName2:  contractInfo2,
				contractAdder1: contractInfo1,
				contractName1:  contractInfo1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetContractByAddersOrNames(tt.args.chainId, tt.args.nameList)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContractByAddersOrNames() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetContractByName(t *testing.T) {
	contractInfo1, err := insertContractTest1()
	if err != nil {
		return
	}

	type args struct {
		chainId      string
		contractName string
	}
	tests := []struct {
		name    string
		args    args
		want    *db.Contract
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:      ChainID,
				contractName: contractName1,
			},
			want:    contractInfo1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetContractByName(tt.args.chainId, tt.args.contractName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContractByName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetContractByNameOrAddr(t *testing.T) {
	contractInfo1, err := insertContractTest1()
	if err != nil {
		return
	}

	type args struct {
		chainId     string
		contractKey string
	}
	tests := []struct {
		name    string
		args    args
		want    *db.Contract
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:     ChainID,
				contractKey: contractName1,
			},
			want:    contractInfo1,
			wantErr: false,
		},
		{
			name: "test: case 2",
			args: args{
				chainId:     ChainID,
				contractKey: contractAdder1,
			},
			want:    contractInfo1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetContractByNameOrAddr(tt.args.chainId, tt.args.contractKey)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContractByNameOrAddr() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetContractList(t *testing.T) {
	contractInfo1, err := insertContractTest1()
	if err != nil {
		return
	}
	contractInfo2, err := insertContractTest2()
	if err != nil {
		return
	}

	contractStatus := int32(0)

	type args struct {
		chainId      string
		offset       int
		limit        int
		status       *int32
		runtimeType  string
		contractType string
		contractKey  string
		creators     []string
		creatorAddrs []string
		upgrades     []string
		upgradeAddrs []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.Contract
		want1   int64
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				offset:  0,
				limit:   10,
			},
			want: []*db.Contract{
				contractInfo2,
				contractInfo1,
			},
			want1:   2,
			wantErr: false,
		},
		{
			name: "test: case 2",
			args: args{
				chainId:     ChainID,
				offset:      0,
				limit:       10,
				contractKey: contractAdder1,
			},
			want: []*db.Contract{
				contractInfo1,
			},
			want1:   1,
			wantErr: false,
		},
		{
			name: "test: case 3",
			args: args{
				chainId:      ChainID,
				offset:       0,
				limit:        10,
				status:       &contractStatus,
				runtimeType:  "evm",
				contractType: "native",
				creators:     []string{contractAdder1},
				creatorAddrs: []string{contractInfo1.Addr},
				upgrades:     []string{contractInfo1.Addr},
				upgradeAddrs: []string{contractInfo1.Addr},
				contractKey:  contractAdder1,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := GetContractList(tt.args.chainId, tt.args.offset, tt.args.limit, tt.args.status, tt.args.runtimeType, tt.args.contractKey,
				tt.args.contractType, "desc", "blockHeight", tt.args.creators, tt.args.creatorAddrs, tt.args.upgrades, tt.args.upgradeAddrs)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContractList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetContractNum(t *testing.T) {
	_, err := insertContractTest1()
	if err != nil {
		return
	}
	_, err = insertContractTest2()
	if err != nil {
		return
	}

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
			want:    2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetContractNum(tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContractNum() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetContractNumCache(t *testing.T) {
	_, err := insertContractTest1()
	if err != nil {
		return
	}

	_, err = insertContractTest2()
	if err != nil {
		return
	}

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
			want:    2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prefix := config.GlobalConfig.RedisDB.Prefix
			redisKey := fmt.Sprintf(cache.RedisOverviewContractCount, prefix, tt.args.chainId)
			_, _ = GetContractNumCache(redisKey)
		})
	}
}

func TestGetLatestContractList(t *testing.T) {
	contractInfo1, err := insertContractTest1()
	if err != nil {
		return
	}
	contractInfo2, err := insertContractTest2()
	if err != nil {
		return
	}

	type args struct {
		chainId string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.Contract
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
			},
			want: []*db.Contract{
				contractInfo2,
				contractInfo1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetLatestContractList(tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLatestContractList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestInsertContract(t *testing.T) {
	if db.GormDB == nil {
		return
	}
	_, err := insertContractTest1()
	if err != nil {
		t.Errorf("InsertContract() error = %v", err)
	}
}

func TestUpdateContract(t *testing.T) {
	contractInfo2, err := insertContractTest2()
	if err != nil {
		return
	}

	type args struct {
		chainId  string
		contract *db.Contract
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
				contract: &db.Contract{
					Addr:           contractInfo2.Addr,
					ContractStatus: 1,
				},
			},
			wantErr: false,
		},
		{
			name: "test: case 2",
			args: args{
				chainId: ChainID,
				contract: &db.Contract{
					Addr:             contractInfo2.Addr,
					ContractStatus:   1,
					UpgradeTimestamp: time.Now().Unix(),
					Upgrader:         AccountAddr1,
					UpgradeAddr:      AccountAddr1,
					UpgradeOrgId:     AccountAddr1,
					Version:          "v1.0.1",
					ContractSymbol:   "Symbol",
					Decimals:         18,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateContract(tt.args.chainId, tt.args.contract); (err != nil) != tt.wantErr {
				t.Errorf("UpdateContract() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateContractNameBak(t *testing.T) {
	contractInfo2, err := insertContractTest2()
	if err != nil {
		return
	}

	type args struct {
		chainId  string
		contract *db.Contract
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
				contract: &db.Contract{
					Addr:    contractInfo2.Addr,
					Name:    "123",
					NameBak: "123",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateContractNameBak(tt.args.chainId, tt.args.contract); (err != nil) != tt.wantErr {
				t.Errorf("UpdateContractNameBak() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateContractTxNum(t *testing.T) {
	contractInfo1, err := insertContractTest1()
	if err != nil {
		return
	}

	type args struct {
		chainId  string
		contract *db.Contract
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
				contract: &db.Contract{
					Addr:  contractInfo1.Addr,
					TxNum: 23,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateContractTxNum(tt.args.chainId, tt.args.contract); (err != nil) != tt.wantErr {
				t.Errorf("UpdateContractTxNum() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetContractCountByRange(t *testing.T) {
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
				chainId:   ChainID,
				startTime: 123455667,
				endTime:   123366666,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetContractCountByRange(tt.args.chainId, tt.args.startTime, tt.args.endTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetContractCountByRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetChainMainDIDContract(t *testing.T) {
	contract := GetChainMainDIDContract(db.UTchainID)
	if contract == nil {
		t.Errorf("GetChainMainDIDContract failed")
	}
}

func TestGetContractByAddr(t *testing.T) {
	contractAddr := "0x123456"
	_, err := GetContractByAddr(db.UTchainID, contractAddr)
	if err != nil {
		t.Errorf("GetContractByAddr failed: %s", err)
	}
}

func TestInsertContractByteCodes(t *testing.T) {
	// Test case 1: Inserting contract byte codes successfully
	byteCode1 := []*db.ContractByteCode{
		{TxId: "tx1", ByteCode: []byte{1, 2, 3}},
		{TxId: "tx2", ByteCode: []byte{4, 5, 6}},
	}
	err1 := InsertContractByteCodes(db.UTchainID, byteCode1)
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}
}

func TestGetContractByteCodeByTx(t *testing.T) {
	// Test case 1: Getting contract byte code by transaction id successfully
	_, err1 := GetContractByteCodeByTx(db.UTchainID, "tx1")
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}

	_, err2 := GetContractByteCodeByTx(db.UTchainID, "tx100")
	if err2 != nil {
		t.Errorf("Test case 2 failed: %v", err2)
	}
}

func TestGetContractTypes(t *testing.T) {
	_, _ = GetContractTypes(db.UTchainID)
}

func TestGetContractByAddrs(t *testing.T) {
	_, err := GetContractByAddrs(db.UTchainID, []string{contractAdder1})
	if err != nil {
		t.Errorf("GetContractByAddrs failed: %s", err)
	}
}
