package dbhandle

import (
	"chainmaker_web/src/db"
	"testing"

	"github.com/google/uuid"
)

var crossId1 = "crossId1"
var crossId2 = "crossId2"

func insertCrossTxTransfersTest1() ([]*db.CrossTransactionTransfer, error) {
	newUUID := uuid.New().String()
	insertList := []*db.CrossTransactionTransfer{
		{
			ID:           newUUID,
			CrossId:      crossId1,
			FromChainId:  crossId1,
			ToChainId:    crossId2,
			ContractName: subContractName1,
			BlockHeight:  12,
		},
	}
	err := InsertCrossTxTransfers(ChainID, insertList)
	return insertList, err
}

func insertCrossBusinessTransactionTest1() ([]*db.CrossBusinessTransaction, error) {
	insertList := []*db.CrossBusinessTransaction{
		{
			TxId:         "1234",
			CrossId:      crossId1,
			SubChainId:   subChainId1,
			ContractName: subContractName1,
		},
	}
	err := InsertCrossBusinessTransaction(ChainID, insertList)
	return insertList, err
}

// func TestCheckCrossIdsExistenceTransfer(t *testing.T) {
// 	_, err := insertCrossTxTransfersTest1()
// 	if err != nil {
// 		return
// 	}

// 	type args struct {
// 		chainId  string
// 		crossIds []string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    map[string]bool
// 		wantErr bool
// 	}{
// 		{
// 			name: "test: case 1",
// 			args: args{
// 				chainId: ChainID,
// 				crossIds: []string{
// 					crossId1,
// 					crossId2,
// 				},
// 			},
// 			want: map[string]bool{
// 				crossId1: true,
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			_, err := CheckCrossIdsExistenceTransfer(tt.args.chainId, tt.args.crossIds)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("CheckCrossIdsExistenceTransfer() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

func TestGetCrossBusinessTxByCross(t *testing.T) {
	insertList, err := insertCrossBusinessTransactionTest1()
	if err != nil {
		return
	}

	type args struct {
		chainId string
		crossId string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.CrossBusinessTransaction
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				crossId: crossId2,
			},
			want:    []*db.CrossBusinessTransaction{},
			wantErr: false,
		},
		{
			name: "test: case 2",
			args: args{
				chainId: ChainID,
				crossId: crossId1,
			},
			want:    insertList,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetCrossBusinessTxByCross(tt.args.chainId, tt.args.crossId)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossBusinessTxByCross() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetCrossBusinessTxById(t *testing.T) {
	insertList, err := insertCrossBusinessTransactionTest1()
	if err != nil {
		return
	}

	type args struct {
		chainId string
		txList  []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.CrossBusinessTransaction
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				txList: []string{
					"1234",
					"22222",
				},
			},
			want:    insertList,
			wantErr: false,
		},
		{
			name: "test: case 2",
			args: args{
				chainId: ChainID,
				txList: []string{
					"22222",
				},
			},
			want:    make([]*db.CrossBusinessTransaction, 0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetCrossBusinessTxById(tt.args.chainId, tt.args.txList)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCrossBusinessTxById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

// func TestGetCrossCycleTransferByHeight(t *testing.T) {
// 	_, err := insertCrossTxTransfersTest1()
// 	if err != nil {
// 		return
// 	}

// 	type args struct {
// 		chainId     string
// 		blockHeight []int64
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "test: case 1",
// 			args: args{
// 				chainId: ChainID,
// 				blockHeight: []int64{
// 					12,
// 					14,
// 				},
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "test: case 2",
// 			args: args{
// 				chainId: ChainID,
// 				blockHeight: []int64{
// 					222222,
// 				},
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			_, err := GetCrossCycleTransferByHeight(tt.args.chainId, tt.args.blockHeight)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("GetCrossCycleTransferByHeight() 1 error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestGetCrossCycleTransferById(t *testing.T) {
// 	_, err := insertCrossTxTransfersTest1()
// 	if err != nil {
// 		return
// 	}

// 	type args struct {
// 		chainId string
// 		crossId string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "test: case 1",
// 			args: args{
// 				chainId: ChainID,
// 				crossId: crossId1,
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "test: case 2",
// 			args: args{
// 				chainId: ChainID,
// 				crossId: crossId2,
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			_, err := GetCrossCycleTransferById(tt.args.chainId, tt.args.crossId)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("GetCrossCycleTransferById() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

func TestInsertCrossBusinessTransaction(t *testing.T) {
	insertList := []*db.CrossBusinessTransaction{
		{
			TxId:         "1234",
			CrossId:      crossId1,
			SubChainId:   subChainId1,
			ContractName: subContractName1,
		},
	}

	type args struct {
		chainId  string
		crossTxs []*db.CrossBusinessTransaction
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
				crossTxs: insertList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InsertCrossBusinessTransaction(tt.args.chainId, tt.args.crossTxs); (err != nil) != tt.wantErr {
				t.Errorf("InsertCrossBusinessTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestInsertCrossSubTransaction(t *testing.T) {
	insertList := []*db.CrossMainTransaction{
		{
			TxId:    "1234",
			CrossId: crossId1,
		},
	}

	type args struct {
		chainId  string
		crossTxs []*db.CrossMainTransaction
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
				crossTxs: insertList,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := InsertCrossSubTransaction(tt.args.chainId, tt.args.crossTxs); (err != nil) != tt.wantErr {
				t.Errorf("InsertCrossSubTransaction() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// func TestInsertCrossTxTransfers(t *testing.T) {
// 	newUUID := uuid.New().String()
// 	insertList := []*db.CrossTransactionTransfer{
// 		{
// 			ID:           newUUID,
// 			CrossId:      crossId1,
// 			FromChainId:  crossId1,
// 			ToChainId:    crossId2,
// 			ContractName: subContractName1,
// 			BlockHeight:  12,
// 		},
// 	}

// 	type args struct {
// 		chainId          string
// 		crossTxTransfers []*db.CrossTransactionTransfer
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "test: case 1",
// 			args: args{
// 				chainId:          ChainID,
// 				crossTxTransfers: insertList,
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if err := InsertCrossTxTransfers(tt.args.chainId, tt.args.crossTxTransfers); (err != nil) != tt.wantErr {
// 				t.Errorf("InsertCrossTxTransfers() error = %v, wantErr %v", err, tt.wantErr)
// 			}
// 		})
// 	}
// }

// func TestGetCrossCycleTransferByCrossIds(t *testing.T) {
// 	// Test case 1: Getting contract byte code by transaction id successfully
// 	crossIds := []string{"121", "23"}
// 	_, err1 := GetCrossCycleTransferByCrossIds(db.UTchainID, crossIds)
// 	if err1 != nil {
// 		t.Errorf("Test case 1 failed: %v", err1)
// 	}
// }
