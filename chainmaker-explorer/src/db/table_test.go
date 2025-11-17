package db

import (
	"reflect"
	"testing"

	"github.com/test-go/testify/assert"
)

func TestGetTableName(t *testing.T) {
	type args struct {
		chainId string
		table   string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "tase case 1",
			args: args{
				chainId: "chain1",
				table:   "table",
			},
			want: "chain1_test_table",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetTableName(tt.args.chainId, tt.args.table); got != tt.want {
				t.Errorf("GetTableName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetClickhouseTableOptions(t *testing.T) {
	tests := []struct {
		name string
		want map[string]string
	}{
		{
			name: "获取Clickhouse表选项",
			want: GetClickhouseTableOptions(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetClickhouseTableOptions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetClickhouseTableOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetClickhouseDBIndexFields(t *testing.T) {
	tests := []struct {
		tableName string
		expected  []CHIndexInfo
	}{
		{TableBlock, []CHIndexInfo{{IndexType: CHIndexTypeIndex, Fields: []string{"blockHeight"}}}},
		{TableTransaction, []CHIndexInfo{
			{IndexType: CHIndexTypeRange, Fields: []string{"timestamp"}},
			{IndexType: CHIndexTypeIndex, Fields: []string{"sender", "blockHeight", "userAddr"}},
		}},
		{TableContractEvent, []CHIndexInfo{
			{IndexType: CHIndexTypeRange, Fields: []string{"timestamp"}},
			{IndexType: CHIndexTypeIndex, Fields: []string{"contractNameBak"}},
		}},
		{TableFungibleContract, []CHIndexInfo{
			{IndexType: CHIndexTypeIndex, Fields: []string{"contractNameBak", "contractAddr"}},
		}},
		{TableNonFungibleContract, []CHIndexInfo{
			{IndexType: CHIndexTypeIndex, Fields: []string{"contractNameBak", "contractAddr"}},
		}},
		{TableEvidenceContract, []CHIndexInfo{
			{IndexType: CHIndexTypeRange, Fields: []string{"timestamp"}},
			{IndexType: CHIndexTypeIndex, Fields: []string{"hash", "contractName"}},
		}},
		{TableIdentityContract, []CHIndexInfo{
			{IndexType: CHIndexTypeIndex, Fields: []string{"contractName", "contractAddr"}},
		}},
		{TableFungibleTransfer, []CHIndexInfo{
			{IndexType: CHIndexTypeRange, Fields: []string{"timestamp"}},
			{IndexType: CHIndexTypeIndex, Fields: []string{"fromAddr", "toAddr", "contractName"}},
		}},
		{TableNonFungibleTransfer, []CHIndexInfo{
			{IndexType: CHIndexTypeRange, Fields: []string{"timestamp"}},
			{IndexType: CHIndexTypeIndex, Fields: []string{"fromAddr", "toAddr", "tokenId"}},
		}},
		{TableFungiblePosition, []CHIndexInfo{
			{IndexType: CHIndexTypeIndex, Fields: []string{"ownerAddr"}},
		}},
		{TableNonFungiblePosition, []CHIndexInfo{
			{IndexType: CHIndexTypeIndex, Fields: []string{"ownerAddr"}},
		}},
		{TableNonFungibleToken, []CHIndexInfo{
			{IndexType: CHIndexTypeIndex, Fields: []string{"ownerAddr"}},
		}},
		{TableAccount, []CHIndexInfo{
			{IndexType: CHIndexTypeIndex, Fields: []string{"did", "bns"}},
		}},
	}

	for _, test := range tests {
		result := GetClickhouseDBIndexFields(test.tableName)
		assert.Equal(t, test.expected, result)
	}
}

func TestPGSqlCreateTableWithComment(t *testing.T) {
	tableInfo := &TableInfo{
		Name:        TableBlock,
		Structure:   &Block{},
		Description: "区块信息表",
	}
	_ = PGSqlCreateTableWithComment(GormDB, UTchainID, *tableInfo)
}

func TestClickHouseCreateTableWithComment(t *testing.T) {
	tableInfo := &TableInfo{
		Name:        TableBlock,
		Structure:   &Block{},
		Description: "区块信息表",
	}
	_ = ClickHouseCreateTableWithComment(UTchainID, *tableInfo)
}
