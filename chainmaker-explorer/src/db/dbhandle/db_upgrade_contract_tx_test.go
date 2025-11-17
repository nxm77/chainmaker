package dbhandle

import (
	"chainmaker_web/src/db"
	"testing"
)

const (
	ContractTxId1   = "123456"
	ContractTxUser1 = "123456789"
	ContractTxUser2 = "223456789"
)

func insertUpgradeContractTxTest() ([]*db.UpgradeContractTransaction, error) {
	insertList := []*db.UpgradeContractTransaction{
		{
			TxId:        ContractTxId1,
			Sender:      "123",
			UserAddr:    ContractTxUser1,
			BlockHeight: 12,
		},
	}
	err := InsertUpgradeContractTx(ChainID, insertList)
	return insertList, err
}

func TestGetUpgradeContractTxList(t *testing.T) {
	_, err := insertUpgradeContractTxTest()
	if err != nil {
		return
	}

	type args struct {
		offset       int
		limit        int
		chainId      string
		contractName string
		contractAddr string
		senders      []string
		runtimeType  string
		status       int
		startTime    int64
		endTime      int64
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.UpgradeContractTransaction
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
			want1:   1,
			wantErr: false,
		},
		{
			name: "test: case 2",
			args: args{
				chainId:      ChainID,
				offset:       0,
				limit:        10,
				contractName: ContractName1,
				contractAddr: ContractAddr1,
				senders:      []string{ContractTxUser1},
				runtimeType:  "wasm",
				status:       1,
				startTime:    0,
				endTime:      0,
			},
			want1:   1,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := GetUpgradeContractTxList(tt.args.offset, tt.args.limit, tt.args.chainId, tt.args.contractName, tt.args.contractAddr,
				tt.args.senders, tt.args.runtimeType, tt.args.status, tt.args.startTime, tt.args.endTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUpgradeContractTxList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestInsertUpgradeContractTx(t *testing.T) {
	insertList := []*db.UpgradeContractTransaction{
		{
			TxId:        ContractTxId1,
			Sender:      "123",
			UserAddr:    ContractTxUser1,
			BlockHeight: 12,
		},
	}

	type args struct {
		chainId      string
		transactions []*db.UpgradeContractTransaction
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
				transactions: insertList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InsertUpgradeContractTx(tt.args.chainId, tt.args.transactions); (err != nil) != tt.wantErr {
				t.Errorf("InsertUpgradeContractTx() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateUpgradeContractName(t *testing.T) {
	_, err := insertUpgradeContractTxTest()
	if err != nil {
		return
	}

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
				chainId:  ChainID,
				contract: contractInfo1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateUpgradeContractName(tt.args.chainId, tt.args.contract); (err != nil) != tt.wantErr {
				t.Errorf("UpdateUpgradeContractName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetContractVersions(t *testing.T) {
	// Test case 1: Normal case with valid data
	_, err1 := GetContractVersions(ChainID, ContractAddr1)
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}
}

func TestGetUpgradeContractInfo(t *testing.T) {
	_, err := insertUpgradeContractTxTest()
	if err != nil {
		return
	}

	_, _ = GetUpgradeContractInfo(ChainID, ContractAddr1, "")
}

func TestUpdateUpgradeContractVerifyStatus(t *testing.T) {
	_ = UpdateUpgradeContractVerifyStatus(ChainID, ContractAddr1, ContractVersionUT, 1)
}
