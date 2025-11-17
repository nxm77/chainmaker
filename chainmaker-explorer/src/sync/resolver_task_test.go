package sync

import (
	"chainmaker_web/src/db"
	"chainmaker_web/src/sync/common"
	"chainmaker_web/src/sync/model"
	"errors"
	"testing"

	"github.com/test-go/testify/assert"
)

func TestDelayParallelParseGetDB(t *testing.T) {
	type args struct {
		chainId            string
		delayedUpdateCache *model.GetRealtimeCacheData
		contractMap        map[string]*db.Contract
		topicEventResult   *model.TopicEventResult
		crossSubChainIdMap map[string]map[string]int64
		eventTopicTxNum    map[string]map[string]int64
	}
	tests := []struct {
		name    string
		args    args
		want    *model.GetDBResult
		wantErr bool
	}{
		{
			name: "Test case 1: Sample input",
			args: args{
				chainId:            "testChainId",
				delayedUpdateCache: &model.GetRealtimeCacheData{},
				contractMap:        map[string]*db.Contract{},
				topicEventResult: &model.TopicEventResult{
					IDAEventData: &model.IDAEventData{},
				},
				crossSubChainIdMap: map[string]map[string]int64{},
			},
			want:    &model.GetDBResult{}, // Fill the expected GetDBResult fields here
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := DelayParallelParseGetDB(tt.args.chainId, tt.args.delayedUpdateCache, tt.args.contractMap, tt.args.topicEventResult,
				tt.args.crossSubChainIdMap, tt.args.eventTopicTxNum)
			if (err != nil) != tt.wantErr {
				t.Errorf("DelayParallelParseGetDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestTaskFunc_Run(t *testing.T) {
	tests := []struct {
		name    string
		f       TaskFunc
		wantErr bool
	}{
		{
			name: "Test case 1: TaskFunc with no error",
			f: TaskFunc(func() error {
				return nil
			}),
			wantErr: false,
		},
		{
			name: "Test case 2: TaskFunc with error",
			f: TaskFunc(func() error {
				return errors.New("sample error")
			}),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.f.Run(); (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_getAccountByBNSListTask(t *testing.T) {
	type args struct {
		chainId         string
		getDBResult     *model.GetDBResult
		bnsUnBindDomain []string
	}
	tests := []struct {
		name string
		args args
		want TaskSync
	}{
		{
			name: "Test case 1: Sample input",
			args: args{
				chainId:         "testChainId",
				getDBResult:     &model.GetDBResult{},
				bnsUnBindDomain: []string{"sample1", "sample2"},
			},
		},
		// Add more test cases if needed
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getAccountByBNSListTask(tt.args.chainId, tt.args.getDBResult, tt.args.bnsUnBindDomain)
		})
	}
}

func Test_getAccountByDIDListTask(t *testing.T) {
	type args struct {
		chainId       string
		getDBResult   *model.GetDBResult
		didUnBindList []string
	}
	tests := []struct {
		name string
		args args
		want TaskSync
	}{
		{
			name: "Test case 1: Sample input",
			args: args{
				chainId:       "testChainId",
				getDBResult:   &model.GetDBResult{},
				didUnBindList: []string{"sample1", "sample2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getAccountByDIDListTask(tt.args.chainId, tt.args.getDBResult, tt.args.didUnBindList)
		})
	}
}

func Test_getAccountMapTask(t *testing.T) {
	type args struct {
		chainId          string
		getDBResult      *model.GetDBResult
		topicEventResult *model.TopicEventResult
		userInfoMap      map[string]*db.User
	}
	tests := []struct {
		name string
		args args
		want TaskSync
	}{
		{
			name: "Test case 1: Sample input",
			args: args{
				chainId:          "testChainId",
				getDBResult:      &model.GetDBResult{},
				topicEventResult: &model.TopicEventResult{},
				userInfoMap:      map[string]*db.User{},
			},
		},
		// Add more test cases if needed
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getAccountMapTask(tt.args.chainId, tt.args.getDBResult, tt.args.topicEventResult, tt.args.userInfoMap)
		})
	}
}

func Test_getAddBlackTxListTask(t *testing.T) {
	type args struct {
		chainId       string
		getDBResult   *model.GetDBResult
		addBlackTxIds []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test case 1: Sample input",
			args: args{
				chainId:       "testChainId",
				getDBResult:   &model.GetDBResult{},
				addBlackTxIds: []string{"sample1", "sample2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getAddBlackTxListTask(tt.args.chainId, tt.args.getDBResult, tt.args.addBlackTxIds)
		})
	}
}
func Test_getDeleteBlackTxListTask(t *testing.T) {
	type args struct {
		chainId          string
		getDBResult      *model.GetDBResult
		deleteBlackTxIds []string
	}
	tests := []struct {
		name string
		args args
		want TaskSync
	}{
		{
			name: "Test case 1: Sample input",
			args: args{
				chainId:          "testChainId",
				getDBResult:      &model.GetDBResult{},
				deleteBlackTxIds: []string{"sample1", "sample2"},
			},
		},
		// Add more test cases if needed
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getDeleteBlackTxListTask(tt.args.chainId, tt.args.getDBResult, tt.args.deleteBlackTxIds)
		})
	}
}

func Test_getFungibleContractTask(t *testing.T) {
	type args struct {
		chainId     string
		getDBResult *model.GetDBResult
		contractMap map[string]*db.Contract
	}
	tests := []struct {
		name string
		args args
		want TaskSync
	}{
		{
			name: "Test case 1: Sample input",
			args: args{
				chainId:     "testChainId",
				getDBResult: &model.GetDBResult{},
				contractMap: map[string]*db.Contract{},
			},
		},
		// Add more test cases if needed
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getFungibleContractTask(tt.args.chainId, tt.args.getDBResult, tt.args.contractMap)
		})
	}
}

func Test_getPositionMapTask(t *testing.T) {
	type args struct {
		chainId     string
		getDBResult *model.GetDBResult
		ownerAdders []string
	}
	tests := []struct {
		name string
		args args
		want TaskSync
	}{
		{
			name: "Test case 1: Sample input",
			args: args{
				chainId:     "testChainId",
				getDBResult: &model.GetDBResult{},
				ownerAdders: []string{"sample1", "sample2"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getPositionMapTask(tt.args.chainId, tt.args.getDBResult, tt.args.ownerAdders)
		})
	}
}

func Test_getGasListTask(t *testing.T) {
	type args struct {
		chainId            string
		delayedUpdateCache *model.GetRealtimeCacheData
		getDBResult        *model.GetDBResult
	}
	tests := []struct {
		name string
		args args
		//want TaskSync
	}{
		{
			name: "Test case 1: Sample input",
			args: args{
				chainId:            "testChainId",
				delayedUpdateCache: nil,
				getDBResult:        &model.GetDBResult{},
			},
			//want: getGasListTask("testChainId", &GetRealtimeCacheData{}, &GetDBResult{}),
		},
		// Add more test cases if needed
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getGasListTask(tt.args.chainId, tt.args.delayedUpdateCache, tt.args.getDBResult)
			//if got := ; !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("getGasListTask() = %v, want %v", got, tt.want)
			//}
		})
	}
}

func Test_getNonFungibleContractTask(t *testing.T) {
	type args struct {
		chainId     string
		getDBResult *model.GetDBResult
		contractMap map[string]*db.Contract
	}
	tests := []struct {
		name string
		args args
		//want TaskSync
	}{
		{
			name: "Test case 1: Sample input",
			args: args{
				chainId:     "testChainId",
				getDBResult: &model.GetDBResult{},
				contractMap: map[string]*db.Contract{
					"sample1": {
						Addr:         "sample1",
						ContractType: common.ContractStandardNameCMNFA,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getNonFungibleContractTask(tt.args.chainId, tt.args.getDBResult, tt.args.contractMap)
		})
	}
}

func Test_getNonPositionMapTask(t *testing.T) {
	type args struct {
		chainId     string
		getDBResult *model.GetDBResult
		ownerAdders []string
	}
	tests := []struct {
		name string
		args args
		//want TaskSync
	}{
		{
			name: "Test case 1: Sample input",
			args: args{
				chainId:     "testChainId",
				getDBResult: &model.GetDBResult{},
				ownerAdders: []string{"sample1", "sample2"},
			},
			//want: getNonPositionMapTask("testChainId", &GetDBResult{}, []string{"sample1", "sample2"}),
		},
		// Add more test cases if needed
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getNonPositionMapTask(tt.args.chainId, tt.args.getDBResult, tt.args.ownerAdders)
		})
	}
}

func Test_getSubChainCrossListTask(t *testing.T) {
	type args struct {
		chainId            string
		getDBResult        *model.GetDBResult
		crossSubChainIdMap map[string]map[string]int64
	}
	tests := []struct {
		name string
		args args
		//want TaskSync
	}{
		{
			name: "Test case 1: Sample input",
			args: args{
				chainId:            "testChainId",
				getDBResult:        &model.GetDBResult{},
				crossSubChainIdMap: map[string]map[string]int64{},
			},
			//want: getSubChainCrossListTask("testChainId", &GetDBResult{}, map[string]map[string]int64{}),
		},
		// Add more test cases if needed
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getSubChainCrossListTask(tt.args.chainId, tt.args.getDBResult, tt.args.crossSubChainIdMap)
		})
	}
}

func Test_getSubChainDBMapTask(t *testing.T) {
	type args struct {
		chainId            string
		getDBResult        *model.GetDBResult
		crossSubChainIdMap map[string]map[string]int64
	}
	tests := []struct {
		name string
		args args
		//want TaskSync
	}{
		{
			name: "Test case 1: Sample input",
			args: args{
				chainId:            "testChainId",
				getDBResult:        &model.GetDBResult{},
				crossSubChainIdMap: map[string]map[string]int64{},
			},
			//want: getSubChainDBMapTask("testChainId", &GetDBResult{}, map[string]map[string]int64{}),
		},
		// Add more test cases if needed
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getSubChainDBMapTask(tt.args.chainId, tt.args.getDBResult, tt.args.crossSubChainIdMap)
		})
	}
}

func TestGetContractEvents(t *testing.T) {
	txIds := []string{"tx1", "tx2"}
	_, err := GetContractEvents(db.UTchainID, txIds)
	assert.Nil(t, err)
}

func TestGetContractEvents_EmptyTxIds(t *testing.T) {
	txIds := []string{}

	_, err := GetContractEvents(db.UTchainID, txIds)
	assert.Nil(t, err)
}
