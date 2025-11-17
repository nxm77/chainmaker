package dbhandle

import (
	"chainmaker_web/src/db"
	"testing"
)

const (
	ContractName1 = "ContractName1"
	ContractName2 = "ContractName2"
	ContractAddr1 = "123456"
	ContractAddr2 = "223456"
	TxId1         = "123456789"
	TxId2         = "223456789"
)

func insertContractEventTest1() ([]*db.ContractEvent, error) {
	contractEvents := []*db.ContractEvent{
		{
			TxId:            TxId1,
			EventIndex:      0,
			Topic:           "123",
			ContractName:    ContractName1,
			ContractNameBak: ContractName1,
			ContractAddr:    ContractAddr1,
			Timestamp:       123456,
		},
	}
	err := InsertContractEvent(ChainId1, contractEvents)
	return contractEvents, err
}

func insertContractEventTest2() ([]*db.ContractEvent, error) {
	contractEvents := []*db.ContractEvent{
		{
			TxId:            TxId2,
			EventIndex:      0,
			Topic:           "123",
			ContractName:    ContractName2,
			ContractNameBak: ContractName2,
			ContractAddr:    ContractAddr2,
			Timestamp:       223456,
		},
	}
	err := InsertContractEvent(ChainId1, contractEvents)
	return contractEvents, err
}

func TestGetEventDataByTxIds(t *testing.T) {
	contractEvent1, err := insertContractEventTest1()
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
		want    []*db.ContractEvent
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId: ChainID,
				txIds: []string{
					TxId1,
				},
			},
			want:    contractEvent1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetEventDataByTxIds(tt.args.chainId, tt.args.txIds)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEventDataByTxIds() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestGetEventList(t *testing.T) {
	contractEvent1, err := insertContractEventTest1()
	if err != nil {
		return
	}

	type args struct {
		offset       int
		limit        int
		chainId      string
		contractName string
		txId         string
		topics       []string
	}
	tests := []struct {
		name    string
		args    args
		want    []*db.ContractEvent
		want1   int64
		wantErr bool
	}{
		{
			name: "test: case 1",
			args: args{
				chainId:      ChainID,
				offset:       0,
				limit:        10,
				contractName: ContractName1,
			},
			want:    contractEvent1,
			want1:   int64(len(contractEvent1)),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := GetEventListCount(tt.args.chainId, tt.args.contractName, tt.args.txId, tt.args.topics)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetEventList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestUpdateContractEventBak(t *testing.T) {
	_, err := insertContractEventTest2()
	if err != nil {
		return
	}

	type args struct {
		chainId       string
		contractEvent *db.ContractEvent
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
				contractEvent: &db.ContractEvent{
					TxId:       TxId2,
					EventIndex: 0,
					Topic:      "12345",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateContractEventBak(tt.args.chainId, tt.args.contractEvent); (err != nil) != tt.wantErr {
				t.Errorf("UpdateContractEventBak() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUpdateContractEventSensitiveWord(t *testing.T) {
	_, err := insertContractTest1()
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
					Addr:    contractAdder1,
					NameBak: "",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UpdateContractEventSensitiveWord(tt.args.chainId, tt.args.contract); (err != nil) != tt.wantErr {
				t.Errorf("UpdateContractEventSensitiveWord() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetContractEventList(t *testing.T) {
	// Test case 1: Normal case with valid data
	offset1 := 0
	limit1 := 10
	txId1 := "txId11"
	topics1 := []string{"topic1"}
	_, err1 := GetContractEventList(offset1, limit1, ChainID, contractName1, txId1, topics1)
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}
}

func TestGetEventDataByIds(t *testing.T) {
	// Test case 1: Normal case with valid data
	ids1 := []string{"id1", "id2"}
	_, err1 := GetEventDataByIds(ChainID, ids1)
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}
}

func TestGetEventNumByContractName(t *testing.T) {
	// Test case 1: Getting contract byte code by transaction id successfully
	_, err1 := GetEventNumByContractName(db.UTchainID, "name")
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}
}

func TestGetABIEventsByVersion(t *testing.T) {
	// Test case 1: Getting contract byte code by transaction id successfully
	_, err1 := GetABIEventsByVersion(0, 1, db.UTchainID, contractAdder1, ContractVersionUT)
	if err1 != nil {
		t.Errorf("Test case 1 failed: %v", err1)
	}
}
