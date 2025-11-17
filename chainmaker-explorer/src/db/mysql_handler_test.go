package db

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/entity"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var TestTableName = "test_table"

func TestMySQLHandler_ConnectDatabase(t *testing.T) {
	config.SubscribeChains = []*config.ChainInfo{
		ChainListConfigTest,
	}

	hostUrl := os.Getenv("UT_MYSQL_DB_URL")
	//hostUrl = "root:123456@tcp(127.0.0.1:33061)/chainmaker_explorer_dev"
	if hostUrl != "" {
		mysqlCfg, _ := parseDBURL(hostUrl)
		config.GlobalConfig.DBConf = mysqlCfg
	}

	tests := []struct {
		name        string
		useDatabase bool
		expectError bool
	}{
		{
			name:        "Connect with database",
			useDatabase: true,
			expectError: false,
		},
		{
			name:        "Connect without database",
			useDatabase: false,
			expectError: false,
		},
	}

	handler := &MySQLHandler{
		DBConfig: *config.GlobalConfig.DBConf,
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := handler.ConnectDatabase(tt.useDatabase)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, db)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, db)
			}
		})
	}
}

func CreateTableTest() error {
	fields := []*DynamicStructField{
		{Name: "field1", Indexed: false},
		{Name: "field2", Indexed: false},
	}
	// 需要创建新表
	structType := CreateDynamicStructWithSystemFields(fields)
	dynamicStruct := reflect.New(structType).Interface()
	return CreateTopicTable(UTchainID, TestTableName, dynamicStruct)
}

func TestMySQLHandler_InsertWithNativeSQL(t *testing.T) {
	err := CreateTableTest()
	if err != nil {
		return
	}

	tests := []struct {
		name        string
		tableName   string
		records     []map[string]interface{}
		expectError bool
	}{
		{
			name:      "Insert with single record",
			tableName: UTchainID + "_test_" + TestTableName,
			records: []map[string]interface{}{
				{"sysId": "300", "sysTxId": "tx1", "sysContractVersion": "v1", "sysTimestamp": 1234567890, "field1": "value1", "field2": 123},
			},
		},
		{
			name:      "Insert with multiple records",
			tableName: UTchainID + "_test_" + TestTableName,
			records: []map[string]interface{}{
				{"sysId": "400", "sysTxId": "tx1", "sysContractVersion": "v1", "sysTimestamp": 1234567890, "field1": "value1", "field2": 123},
				{"sysId": "500", "sysTxId": "tx2", "sysContractVersion": "v1", "sysTimestamp": 1234567890, "field1": "value2", "field2": 456},
			},
		},
	}

	handler := &MySQLHandler{DBConfig: config.DBConf{}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.InsertWithNativeSQL(tt.tableName, tt.records)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMySQLHandler_GetDecodeEventByABIAndTotal(t *testing.T) {
	config.SubscribeChains = []*config.ChainInfo{
		ChainListConfigTest,
	}

	hostUrl := os.Getenv("UT_MYSQL_DB_URL")
	//hostUrl = "root:123456@tcp(127.0.0.1:33061)/chainmaker_explorer_dev"
	if hostUrl != "" {
		mysqlCfg, _ := parseDBURL(hostUrl)
		config.GlobalConfig.DBConf = mysqlCfg
	}

	handler := &MySQLHandler{
		DBConfig: *config.GlobalConfig.DBConf,
	}

	tests := []struct {
		name         string
		offset       int
		limit        int
		chainId      string
		contractAddr string
		version      string
		topic        string
		tableName    string
		topicColumns []string
		searchParams []entity.SearchParam
		expectError  bool
	}{
		{
			name:         "Missing required parameters",
			offset:       0,
			limit:        10,
			chainId:      "",
			contractAddr: "",
			version:      "",
			tableName:    "",
			expectError:  true,
		},
		{
			name:         "Valid parameters",
			offset:       0,
			limit:        10,
			chainId:      UTchainID,
			contractAddr: "addr1",
			version:      "v1",
			tableName:    TestTableName,
			topicColumns: []string{"field1"},
			searchParams: []entity.SearchParam{
				{Name: "field1", Value: "value1"},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := handler.GetDecodeEventByABIAndTotal(
				tt.offset, tt.limit, tt.chainId, tt.contractAddr, tt.version, tt.topic,
				tt.tableName, tt.topicColumns, tt.searchParams,
			)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
