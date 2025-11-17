/*
 * @Author: dongxuliang dongxuliang@tencent.com
 * @Date: 2024-07-15 14:22:41
 * @LastEditors: dongxuliang dongxuliang@tencent.com
 * @LastEditTime: 2024-07-15 15:16:54
 * @FilePath: /chainmaker-explorer-backend/src/sync/account_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package logic

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/sync/common"
	"chainmaker_web/src/sync/model"
	"testing"
)

func Test_processBNSAccounts(t *testing.T) {
	type args struct {
		bnsBindEventData []*db.BNSTopicEventData
		unBindBNSs       []*db.Account
		accountInsertMap map[string]*db.Account
		accountUpdateMap map[string]*db.Account
		accountMap       map[string]*db.Account
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "测试处理BNS帐户",
			args: args{
				bnsBindEventData: []*db.BNSTopicEventData{
					{
						Domain: "1234",
						Value:  "12345",
					},
				},
				unBindBNSs: []*db.Account{
					{
						Address: "123456",
						BNS:     "1234",
					},
				},
				accountInsertMap: map[string]*db.Account{
					"12345": {
						Address: "123456",
						BNS:     "1234",
						DID:     "did:32345",
					},
				},
				accountUpdateMap: map[string]*db.Account{
					"67890": {
						Address: "67890",
						BNS:     "1234",
						DID:     "did:32345",
					},
				},
				accountMap: map[string]*db.Account{
					"12345": {
						Address: "123456",
						BNS:     "1234",
						DID:     "did:32345",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processBNSAccounts(tt.args.bnsBindEventData, tt.args.unBindBNSs, tt.args.accountInsertMap, tt.args.accountUpdateMap, tt.args.accountMap)
		})
	}
}

// func Test_processDIDAccounts(t *testing.T) {
// 	type args struct {
// 		didAccount       map[string][]string
// 		unBindDIDs       []*db.Account
// 		accountInsertMap map[string]*db.Account
// 		accountUpdateMap map[string]*db.Account
// 		accountMap       map[string]*db.Account
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 	}{
// 		{
// 			name: "测试处理BNS帐户",
// 			args: args{
// 				didAccount: map[string][]string{
// 					"did:12345": {
// 						"12345",
// 						"22345",
// 					},
// 					"did:22345": {
// 						"67890",
// 						"6789000",
// 					},
// 				},
// 				unBindDIDs: []*db.Account{
// 					{
// 						Address: "123456",
// 						BNS:     "1234",
// 						DID:     "did:32345",
// 					},
// 				},
// 				accountInsertMap: map[string]*db.Account{
// 					"12345": {
// 						Address: "123456",
// 						BNS:     "1234",
// 						DID:     "did:32345",
// 					},
// 				},
// 				accountUpdateMap: map[string]*db.Account{
// 					"67890": {
// 						Address: "67890",
// 						BNS:     "1234",
// 						DID:     "did:32345",
// 					},
// 				},
// 				accountMap: map[string]*db.Account{
// 					"12345": {
// 						Address: "123456",
// 						BNS:     "1234",
// 						DID:     "did:32345",
// 					},
// 				},
// 			},
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			processDIDAccounts(tt.args.didAccount, tt.args.unBindDIDs, tt.args.accountInsertMap, tt.args.accountUpdateMap, tt.args.accountMap)
// 		})
// 	}
// }

func TestGetAccountType2(t *testing.T) {
	type args struct {
		chainId string
		address string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "test case 1",
			args: args{
				chainId: db.UTchainID,
				address: "123456789",
			},
			want: common.AddrTypeUser,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetAccountType(tt.args.chainId, tt.args.address); got != tt.want {
				t.Errorf("GetAccountType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dealAccountNFTNum(t *testing.T) {
	type args struct {
		chainId          string
		minHeight        int64
		accountNFT       map[string]int64
		accountInsertMap map[string]*db.Account
		accountUpdateMap map[string]*db.Account
		accountMap       map[string]*db.Account
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test case 1",
			args: args{
				chainId: db.UTchainID,
				accountNFT: map[string]int64{
					"123456789": 12,
				},
				accountInsertMap: map[string]*db.Account{},
				accountUpdateMap: map[string]*db.Account{},
				accountMap:       map[string]*db.Account{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dealAccountNFTNum(tt.args.chainId, tt.args.minHeight, tt.args.accountNFT, tt.args.accountInsertMap, tt.args.accountUpdateMap, tt.args.accountMap)
		})
	}
}

func Test_dealAccountTxNum(t *testing.T) {
	type args struct {
		chainId          string
		minHeight        int64
		accountTx        map[string]int64
		accountInsertMap map[string]*db.Account
		accountUpdateMap map[string]*db.Account
		accountMap       map[string]*db.Account
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test case 1",
			args: args{
				chainId: db.UTchainID,
				accountTx: map[string]int64{
					"123456789": 12,
				},
				accountInsertMap: map[string]*db.Account{},
				accountUpdateMap: map[string]*db.Account{},
				accountMap:       map[string]*db.Account{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dealAccountTxNum(tt.args.chainId, tt.args.minHeight, tt.args.accountTx, tt.args.accountInsertMap, tt.args.accountUpdateMap, tt.args.accountMap)
		})
	}
}

func TestAccountHandler(t *testing.T) {
	// Test case 1: Test the construction of AccountHandler
	minHeight := int64(100)
	eventResults := &model.TopicEventResult{
		OwnerAdders: []string{
			"123456789",
		},
	}

	delayedUpdateCache := &model.GetRealtimeCacheData{}
	delayGetDBResult := &model.GetDBResult{
		AccountDBMap: map[string]*db.Account{
			"123456789": {
				Address: "123456789",
			},
		},
	}
	transferEvents := []*db.ContractEventData{
		{
			Index: 1,
			Topic: "transfer",
			TxId:  "1212",
			EventData: &db.TransferTopicEventData{
				FromAddress: "222",
				ToAddress:   "23223",
				TokenId:     "1111",
			},
		},
	}
	accountHandler := NewAccountHandler(
		db.UTchainID,
		minHeight,
		eventResults,
		delayedUpdateCache,
		delayGetDBResult,
		transferEvents,
	)
	if accountHandler.ChainId != db.UTchainID || accountHandler.MinHeight != minHeight {
		t.Errorf("Test case 1 failed")
	}

	// Test case 2: Test the DealWithAccountData method
	accountHandler.AccountMap = make(map[string]*db.Account)
	accountHandler.DelayGetDBResult = &model.GetDBResult{}
	accountHandler.EventResults = &model.TopicEventResult{}
	accountHandler.TransferEvents = []*db.ContractEventData{}
	accountHandler.TxList = make(map[string]*db.Transaction)
	result := accountHandler.DealWithAccountData()
	if result == nil {
		t.Errorf("Test case 2 failed")
	}

	// Test case 3: Test the BuildAccountInsertOrUpdate method
	accountHandler.AccountMap = make(map[string]*db.Account)
	accountHandler.DelayGetDBResult = &model.GetDBResult{}
	accountHandler.EventResults = &model.TopicEventResult{}
	accountHandler.TransferEvents = []*db.ContractEventData{}
	accountHandler.TxList = make(map[string]*db.Transaction)
	_, _ = accountHandler.BuildAccountInsertOrUpdate(make(map[string]int64), make(map[string]int64))

	// Test case 4: Test the DealAccountTxNFTNum method
	accountHandler.AccountMap = make(map[string]*db.Account)
	accountHandler.DelayGetDBResult = &model.GetDBResult{}
	accountHandler.EventResults = &model.TopicEventResult{}
	accountHandler.TransferEvents = []*db.ContractEventData{}
	accountHandler.TxList = make(map[string]*db.Transaction)
	accountTxNum, accountNFTNum := accountHandler.DealAccountTxNFTNum()
	if accountTxNum == nil || accountNFTNum == nil {
		t.Errorf("Test case 4 failed")
	}
}

func TestProcessDIDAccounts1(t *testing.T) {
	// Test case 1: Account exists in accountMap
	bindAccounts1 := map[string]*db.Account{"address1": {Address: "address1", DID: "did1"}}
	unBindAccounts1 := []*db.Account{}
	accountInsertMap1 := make(map[string]*db.Account)
	accountUpdateMap1 := make(map[string]*db.Account)
	accountMap1 := map[string]*db.Account{"address1": {Address: "address1", DID: "oldDID"}}
	processDIDAccounts(bindAccounts1, unBindAccounts1, accountInsertMap1, accountUpdateMap1, accountMap1)
	expected1 := map[string]*db.Account{"address1": {Address: "address1", DID: "did1"}}
	if len(accountUpdateMap1) != len(expected1) {
		t.Errorf("Test case 1 failed: expected %v, got %v", expected1, accountUpdateMap1)
	} else {
		for address, account := range accountUpdateMap1 {
			if account.DID != expected1[address].DID {
				t.Errorf("Test case 1 failed: expected %s, got %s", expected1[address].DID, account.DID)
			}
		}
	}
}

func TestGetAccountType(t *testing.T) {
	// Test case 1: Address is a user address
	address1 := "userAddress1"
	expected1 := common.AddrTypeUser
	result1 := GetAccountType(db.UTchainID, address1)
	if result1 != expected1 {
		t.Errorf("Test case 1 failed: expected %d, got %d", expected1, result1)
	}
}

func TestDealAccountNFTNum(t *testing.T) {
	// Test case 1: Account exists in accountInsertMap
	minHeight1 := int64(100)
	accountNFT1 := map[string]int64{"address1": 10}
	accountInsertMap1 := map[string]*db.Account{"address1": {Address: "address1", NFTNum: 5}}
	accountUpdateMap1 := make(map[string]*db.Account)
	accountMap1 := make(map[string]*db.Account)
	dealAccountNFTNum(db.UTchainID, minHeight1, accountNFT1, accountInsertMap1, accountUpdateMap1, accountMap1)

	// Test case 2: Account exists in accountUpdateMap
	minHeight2 := int64(100)
	accountNFT2 := map[string]int64{"address2": 10}
	accountInsertMap2 := make(map[string]*db.Account)
	accountUpdateMap2 := map[string]*db.Account{"address2": {Address: "address2", NFTNum: 5}}
	accountMap2 := make(map[string]*db.Account)
	dealAccountNFTNum(db.UTchainID, minHeight2, accountNFT2, accountInsertMap2, accountUpdateMap2, accountMap2)

	// Test case 3: Account exists in accountMap
	minHeight3 := int64(100)
	accountNFT3 := map[string]int64{"address3": 10}
	accountInsertMap3 := make(map[string]*db.Account)
	accountUpdateMap3 := make(map[string]*db.Account)
	accountMap3 := map[string]*db.Account{"address3": {Address: "address3", NFTNum: 5}}
	dealAccountNFTNum(db.UTchainID, minHeight3, accountNFT3, accountInsertMap3, accountUpdateMap3, accountMap3)

	// Test case 4: Account does not exist in any map
	minHeight4 := int64(100)
	accountNFT4 := map[string]int64{"address4": 10}
	accountInsertMap4 := make(map[string]*db.Account)
	accountUpdateMap4 := make(map[string]*db.Account)
	accountMap4 := make(map[string]*db.Account)
	dealAccountNFTNum(db.UTchainID, minHeight4, accountNFT4, accountInsertMap4, accountUpdateMap4, accountMap4)
}

func TestDealAccountTxNum(t *testing.T) {
	// Test case 1: Account exists in accountInsertMap
	minHeight1 := int64(100)
	accountTx1 := map[string]int64{"address1": 10}
	accountInsertMap1 := map[string]*db.Account{"address1": {Address: "address1", TxNum: 5}}
	accountUpdateMap1 := make(map[string]*db.Account)
	accountMap1 := make(map[string]*db.Account)
	dealAccountTxNum(db.UTchainID, minHeight1, accountTx1, accountInsertMap1, accountUpdateMap1, accountMap1)

	// Test case 2: Account exists in accountUpdateMap
	chainId2 := "chain2"
	minHeight2 := int64(100)
	accountTx2 := map[string]int64{"address2": 10}
	accountInsertMap2 := make(map[string]*db.Account)
	accountUpdateMap2 := map[string]*db.Account{"address2": {Address: "address2", TxNum: 5}}
	accountMap2 := make(map[string]*db.Account)
	dealAccountTxNum(chainId2, minHeight2, accountTx2, accountInsertMap2, accountUpdateMap2, accountMap2)

	// Test case 3: Account exists in accountMap
	chainId3 := "chain3"
	minHeight3 := int64(100)
	accountTx3 := map[string]int64{"address3": 10}
	accountInsertMap3 := make(map[string]*db.Account)
	accountUpdateMap3 := make(map[string]*db.Account)
	accountMap3 := map[string]*db.Account{"address3": {Address: "address3", TxNum: 5}}
	dealAccountTxNum(chainId3, minHeight3, accountTx3, accountInsertMap3, accountUpdateMap3, accountMap3)

	// Test case 4: Account does not exist in any map
	chainId4 := "chain4"
	minHeight4 := int64(100)
	accountTx4 := map[string]int64{"address4": 10}
	accountInsertMap4 := make(map[string]*db.Account)
	accountUpdateMap4 := make(map[string]*db.Account)
	accountMap4 := make(map[string]*db.Account)
	dealAccountTxNum(chainId4, minHeight4, accountTx4, accountInsertMap4, accountUpdateMap4, accountMap4)
}

func TestProcessBNSAccounts(t *testing.T) {
	// Test case 1: Account exists in accountMap
	bnsBindEventData1 := []*db.BNSTopicEventData{{Value: "address1", Domain: "bns1"}}
	unBindBNSs1 := []*db.Account{}
	accountInsertMap1 := make(map[string]*db.Account)
	accountUpdateMap1 := make(map[string]*db.Account)
	accountMap1 := map[string]*db.Account{"address1": {Address: "address1", BNS: "oldBNS"}}
	processBNSAccounts(bnsBindEventData1, unBindBNSs1, accountInsertMap1, accountUpdateMap1, accountMap1)

	// Test case 2: Account does not exist in accountMap
	bnsBindEventData2 := []*db.BNSTopicEventData{{Value: "address2", Domain: "bns2"}}
	unBindBNSs2 := []*db.Account{}
	accountInsertMap2 := make(map[string]*db.Account)
	accountUpdateMap2 := make(map[string]*db.Account)
	accountMap2 := make(map[string]*db.Account)
	processBNSAccounts(bnsBindEventData2, unBindBNSs2, accountInsertMap2, accountUpdateMap2, accountMap2)

	// Test case 3: Account exists in accountMap and BNS is unbound
	bnsBindEventData3 := []*db.BNSTopicEventData{}
	unBindBNSs3 := []*db.Account{{Address: "address3"}}
	accountInsertMap3 := make(map[string]*db.Account)
	accountUpdateMap3 := make(map[string]*db.Account)
	accountMap3 := map[string]*db.Account{"address3": {Address: "address3", BNS: "bns3"}}
	processBNSAccounts(bnsBindEventData3, unBindBNSs3, accountInsertMap3, accountUpdateMap3, accountMap3)
}

func TestProcessDIDAccounts(t *testing.T) {
	// Test case 1: Account exists in accountMap
	bindAccounts1 := map[string]*db.Account{"address1": {Address: "address1", DID: "did1"}}
	unBindAccounts1 := []*db.Account{}
	accountInsertMap1 := make(map[string]*db.Account)
	accountUpdateMap1 := make(map[string]*db.Account)
	accountMap1 := map[string]*db.Account{"address1": {Address: "address1", DID: "oldDID"}}
	processDIDAccounts(bindAccounts1, unBindAccounts1, accountInsertMap1, accountUpdateMap1, accountMap1)

	// Test case 2: Account does not exist in accountMap
	bindAccounts2 := map[string]*db.Account{"address2": {Address: "address2", DID: "did2"}}
	unBindAccounts2 := []*db.Account{}
	accountInsertMap2 := make(map[string]*db.Account)
	accountUpdateMap2 := make(map[string]*db.Account)
	accountMap2 := make(map[string]*db.Account)
	processDIDAccounts(bindAccounts2, unBindAccounts2, accountInsertMap2, accountUpdateMap2, accountMap2)

	// Test case 3: Account exists in accountMap and DID is unbound
	bindAccounts3 := map[string]*db.Account{}
	unBindAccounts3 := []*db.Account{{Address: "address3"}}
	accountInsertMap3 := make(map[string]*db.Account)
	accountUpdateMap3 := make(map[string]*db.Account)
	accountMap3 := map[string]*db.Account{"address3": {Address: "address3", DID: "did3"}}
	processDIDAccounts(bindAccounts3, unBindAccounts3, accountInsertMap3, accountUpdateMap3, accountMap3)
	expected3 := map[string]*db.Account{"address3": {Address: "address3", DID: ""}}
	if len(accountUpdateMap3) != len(expected3) {
		t.Errorf("Test case 3 failed: expected %v, got %v", expected3, accountUpdateMap3)
	}
}
