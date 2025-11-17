package dbhandle

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/entity"
	"reflect"
	"testing"
)

const (
	txId1      = "123456"
	txId2      = "223456"
	timestamp1 = 12345000
	timestamp2 = 12445000
)

func insertTxTest() ([]*db.Transaction, error) {
	insertList := []*db.Transaction{
		{
			TxId:         txId1,
			Sender:       "123",
			UserAddr:     ContractTxUser1,
			BlockHeight:  12,
			Timestamp:    timestamp1,
			ContractName: ContractName1,
		},
		{
			TxId:         txId2,
			Sender:       "123",
			UserAddr:     ContractTxUser2,
			BlockHeight:  12,
			Timestamp:    timestamp2,
			ContractName: ContractName2,
		},
	}
	err := InsertTransactions(ChainID, insertList)
	return insertList, err
}

func insertBlackTransactionsTest() ([]*db.BlackTransaction, error) {
	insertList := []*db.BlackTransaction{
		{
			TxId:         txId1,
			Sender:       "123",
			UserAddr:     ContractTxUser1,
			BlockHeight:  12,
			Timestamp:    timestamp1,
			ContractName: ContractName1,
		},
		{
			TxId:         txId2,
			Sender:       "123",
			UserAddr:     ContractTxUser1,
			BlockHeight:  12,
			Timestamp:    timestamp2,
			ContractName: ContractName2,
		},
	}
	err := InsertBlackTransactions(ChainID, insertList)
	return insertList, err
}

func TestBatchQueryBlackTxList(t *testing.T) {
	txList, err := insertBlackTransactionsTest()
	if err != nil {
		return
	}

	type args struct {
		chainId string
		txIds   []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.BlackTransaction
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				txIds: []string{
					txId1,
					txId2,
				},
			},
			want:    txList,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := BatchQueryBlackTxList(tt.args.chainId, tt.args.txIds)
			if (err != nil) != tt.wantErr {
				t.Errorf("BatchQueryBlackTxList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}

func TestBatchQueryTxList(t *testing.T) {
	txList, err := insertTxTest()
	if err != nil {
		return
	}

	type args struct {
		chainId string
		txIds   []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.Transaction
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				txIds: []string{
					txId1,
					txId2,
				},
			},
			want:    txList,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _ = BatchQueryTxList(tt.args.chainId, tt.args.txIds)
		})
	}
}

func TestDeleteBlackTransaction(t *testing.T) {
	_, err := insertBlackTransactionsTest()
	if err != nil {
		return
	}

	insertList := []*db.Transaction{
		{
			TxId:        txId1,
			Sender:      "123",
			UserAddr:    ContractTxUser1,
			BlockHeight: 12,
		},
		{
			TxId:        txId2,
			Sender:      "123",
			UserAddr:    ContractTxUser1,
			BlockHeight: 12,
		},
	}

	type args struct {
		chainId      string
		transactions []*db.Transaction
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
			if err := DeleteBlackTransaction(tt.args.chainId, tt.args.transactions); (err != nil) != tt.wantErr {
				t.Errorf("DeleteBlackTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteTransactionByTxId(t *testing.T) {
	_, err := insertTxTest()
	if err != nil {
		return
	}

	type args struct {
		chainId string
		txIds   []string
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
				txIds: []string{
					txId1,
					txId2,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteTransactionByTxId(tt.args.chainId, tt.args.txIds); (err != nil) != tt.wantErr {
				t.Errorf("DeleteTransactionByTxId() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetBlackTxInfoByTxId(t *testing.T) {
	insertList, err := insertBlackTransactionsTest()
	if err != nil {
		return
	}

	type args struct {
		chainId string
		txId    string
	}
	tests := []struct {
		name    string
		args    args
		want    *db.BlackTransaction
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				txId:    txId1,
			},
			want:    insertList[0],
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _ = GetBlackTxInfoByTxId(tt.args.chainId, tt.args.txId)
		})
	}
}

func TestGetLatestTxList(t *testing.T) {
	insertList, err := insertTxTest()
	if err != nil {
		return
	}

	type args struct {
		chainId string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.Transaction
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
			},
			want: []*db.Transaction{
				insertList[1],
				insertList[0],
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetLatestTxList(tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLatestTxList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetSafeWordTransactionList(t *testing.T) {
	_, err := insertTxTest()
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
		want    []*db.Transaction
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
			},
			want:    []*db.Transaction{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSafeWordTransactionList(tt.args.chainId, tt.args.startTime, tt.args.endTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSafeWordTransactionList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSafeWordTransactionList() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTransactionByTxId(t *testing.T) {
	insertList, err := insertTxTest()
	if err != nil {
		return
	}

	type args struct {
		txId    string
		chainId string
	}
	tests := []struct {
		name    string
		args    args
		want    *db.Transaction
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				txId:    txId1,
			},
			want:    insertList[0],
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetTransactionByTxId(tt.args.txId, tt.args.chainId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTransactionByTxId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetTransactionList(t *testing.T) {
	insertList, err := insertTxTest()
	if err != nil {
		return
	}

	type args struct {
		chainId      string
		offset       int
		limit        int
		txStatus     int
		contractName string
		blockHash    string
		startTime    int64
		endTime      int64
		txId         string
		senders      []string
		userAddrs    []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.Transaction
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
			want: []*db.Transaction{
				insertList[1],
				insertList[0],
			},
			want1:   int64(len(insertList)),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetTransactionListCount(tt.args.chainId, tt.args.txId, tt.args.contractName, tt.args.blockHash,
				tt.args.startTime, tt.args.endTime, tt.args.txStatus, tt.args.senders, tt.args.userAddrs)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTransactionList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetTransactionNumByRange(t *testing.T) {
	_, err := insertTxTest()
	if err != nil {
		return
	}

	type args struct {
		chainId   string
		userAddr  string
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
				chainId:  ChainID,
				userAddr: ContractTxUser2,
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetTransactionNumByRange(tt.args.chainId, tt.args.userAddr, tt.args.startTime, tt.args.endTime)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTransactionNumByRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetTxInfoByBlockHeight(t *testing.T) {
	insertList, err := insertTxTest()
	if err != nil {
		return
	}

	type args struct {
		chainId     string
		blockHeight []int64
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.Transaction
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:     ChainID,
				blockHeight: []int64{12, 13},
			},
			want:    insertList,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetTxInfoByBlockHeight(tt.args.chainId, tt.args.blockHeight)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTxInfoByBlockHeight() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetTxListNumByRange(t *testing.T) {
	_, err := insertTxTest()
	if err != nil {
		return
	}

	type args struct {
		chainId   string
		startTime int64
		endTime   int64
		interval  int64
	}
	tests := []struct {
		name    string
		args    args
		want    map[int64]int64
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:   ChainID,
				startTime: timestamp1 - 3600,
				endTime:   timestamp2 + 3600,
				interval:  3600,
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTxListNumByRange(tt.args.chainId, tt.args.startTime, tt.args.endTime, tt.args.interval)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTxListNumByRange() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 {
				t.Errorf("GetTxListNumByRange() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTxNumByContractName(t *testing.T) {
	_, err := insertTxTest()
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
		want    int64
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:      ChainID,
				contractName: contractName1,
			},
			want:    0,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTxNumByContractName(tt.args.chainId, tt.args.contractName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetTxNumByContractName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetTxNumByContractName() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInsertBlackTransactions(t *testing.T) {
	insertList := []*db.BlackTransaction{
		{
			TxId:         txId1,
			Sender:       "123",
			UserAddr:     ContractTxUser1,
			BlockHeight:  12,
			Timestamp:    timestamp1,
			ContractName: ContractName1,
		},
		{
			TxId:         txId2,
			Sender:       "123",
			UserAddr:     ContractTxUser1,
			BlockHeight:  12,
			Timestamp:    timestamp2,
			ContractName: ContractName2,
		},
	}

	type args struct {
		chainId      string
		transactions []*db.BlackTransaction
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
			if err := InsertBlackTransactions(tt.args.chainId, tt.args.transactions); (err != nil) != tt.wantErr {
				t.Errorf("InsertBlackTransactions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInsertTransactions(t *testing.T) {
	insertList := []*db.Transaction{
		{
			TxId:         txId1,
			Sender:       "123",
			UserAddr:     ContractTxUser1,
			BlockHeight:  12,
			Timestamp:    timestamp1,
			ContractName: ContractName1,
		},
		{
			TxId:         txId2,
			Sender:       "123",
			UserAddr:     ContractTxUser2,
			BlockHeight:  12,
			Timestamp:    timestamp2,
			ContractName: ContractName2,
		},
	}

	type args struct {
		chainId      string
		transactions []*db.Transaction
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
			if err := InsertTransactions(tt.args.chainId, tt.args.transactions); (err != nil) != tt.wantErr {
				t.Errorf("InsertTransactions() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateTransactionBak(t *testing.T) {
	insertList, err := insertTxTest()
	if err != nil {
		return
	}

	type args struct {
		chainId     string
		transaction *db.Transaction
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
				transaction: insertList[0],
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateTransactionBak(tt.args.chainId, tt.args.transaction); (err != nil) != tt.wantErr {
				t.Errorf("UpdateTransactionBak() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateTransactionContractName(t *testing.T) {
	contractInfo1, err := insertContractTest1()
	if err != nil {
		return
	}

	_, err = insertTxTest()
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
			if err := UpdateTransactionContractName(tt.args.chainId, tt.args.contract); (err != nil) != tt.wantErr {
				t.Errorf("UpdateTransactionContractName() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetTransactionIDList(t *testing.T) {
	// Test case 1: Normal case with valid data
	blockHash1 := "blockHash123"
	offset1 := 0
	limit1 := 10
	startTime1 := int64(0)
	endTime1 := int64(0)
	txId1 := "txId12"
	txStatus1 := 0
	senders1 := []string{"sender1"}
	userAddrs1 := []string{"userAddr1"}
	_, err1 := GetTransactionIDList(ChainID, ContractName1, blockHash1, offset1, limit1,
		startTime1, endTime1, txId1, txStatus1, senders1, userAddrs1)
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}
}

func TestGetBlockTransactionList(t *testing.T) {
	// Test case 1: Normal case with valid data
	txIds1 := []string{"txId12"}
	_, err1 := GetBlockTransactionList(ChainID, txIds1)
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}
}

func TestGetContractTransactionList(t *testing.T) {
	// Test case 1: Normal case with valid data
	txIds1 := []string{"txId13"}
	_, err1 := GetContractTransactionList(ChainID, txIds1)
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}
}

func TestGetBlockTxIDList(t *testing.T) {
	// Test case 1: Normal case with valid data
	blockHash1 := "blockHash1"
	offset1 := 0
	limit1 := 10
	_, err1 := GetBlockTxIDList(ChainID, blockHash1, offset1, limit1)
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}
}

func TestGetContractTxIDList(t *testing.T) {
	// Test case 1: Normal case with valid data
	offset1 := 0
	limit1 := 10
	contractMethod1 := "contractMethod1"
	userAddrs1 := []string{"userAddr1"}
	txStatus1 := 0
	_, err1 := GetContractTxIDList(offset1, limit1, ChainID, ContractAddr1, contractMethod1, userAddrs1, txStatus1)
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}
}

func TestGetContractTxCount(t *testing.T) {
	// Test case 1: Normal case with valid data
	contractMethod1 := "contractMethod2"
	userAddrs1 := []string{"userAddr1"}
	txStatus1 := 0
	_, err1 := GetContractTxCount(ChainID, ContractAddr1, contractMethod1, userAddrs1, txStatus1)
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}
}

func TestGetUserTxIDList(t *testing.T) {
	_, err := GetUserTxIDList(db.UTchainID, []string{"userAddr1"}, 0, 10)
	if err != nil {
		t.Errorf("Test case 1 failed: %v", err)
	}
}

func TestGetTxCountByRange(t *testing.T) {
	_, err := GetTxCountByRange(db.UTchainID, 0, 10)
	if err != nil {
		t.Errorf("Test case 1 failed: %v", err)
	}
}

func TestGetQueryTxIDList(t *testing.T) {
	params := &entity.GetQueryTxListParams{
		TxStatus:       1,
		StartTime:      int64(122333),
		EndTime:        int64(434344),
		ChainId:        db.UTchainID,
		ContractName:   "contract1",
		ContractAddr:   "contractAddr1",
		ContractMethod: "method1",
		UserAddr:       "userAddr1",
		Operator:       "and",
		TxId:           "tx1",
	}
	_, _, err := GetQueryTxIDList(params)
	if err != nil {
		t.Errorf("GetQueryTxIDList failed: %v", err)
	}
}

func TestGetBlockTxListByHash(t *testing.T) {
	_, err := GetBlockTxListByHash(db.UTchainID, "blockHash1", 0, 1)
	if err != nil {
		t.Errorf("GetBlockTxListByHash failed: %v", err)
	}
}

func TestGetBlockTxIdsByHeight(t *testing.T) {
	_, err := GetBlockTxIdsByHeight(db.UTchainID, 1)
	if err != nil {
		t.Errorf("GetBlockTxIdsByHeight failed: %v", err)
	}
}
