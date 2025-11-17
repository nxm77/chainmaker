package db

import (
	"chainmaker_web/src/config"
	"os"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	_ = config.InitConfig("", "")
	// 初始化数据库配置
	InitRedisContainer()
	InitMySQLContainer()

	// 运行其他测试
	os.Exit(m.Run())
}

type mockDBHandler struct{}

func (m *mockDBHandler) ConnectDatabase(useDatabase bool) (*gorm.DB, error) {
	return &gorm.DB{}, nil
}
func (m *mockDBHandler) InsertWithNativeSQL(tableName string, records []map[string]interface{}) error {
	return nil
}
func (m *mockDBHandler) GetDecodeEventByABIAndTotal(offset, limit int, chainId, contractAddr, version, topic,
	tableName string, topicColumns []string, searchParams []interface{}) ([]map[string]interface{}, int64, error) {
	return nil, 0, nil
}

func TestGetDatabaseHandler(t *testing.T) {
	old := config.GlobalConfig.DBConf.DbProvider
	defer func() { config.GlobalConfig.DBConf.DbProvider = old }()
	config.GlobalConfig.DBConf.DbProvider = config.MySql
	h, err := GetDatabaseHandler()
	assert.NoError(t, err)
	assert.NotNil(t, h)
	config.GlobalConfig.DBConf.DbProvider = config.Pgsql
	h, err = GetDatabaseHandler()
	assert.NoError(t, err)
	assert.NotNil(t, h)
	config.GlobalConfig.DBConf.DbProvider = "not_supported"
	h, err = GetDatabaseHandler()
	assert.Error(t, err)
	assert.Nil(t, h)
}

func TestCapitalizeFirstLetter(t *testing.T) {
	assert.Equal(t, "Hello", CapitalizeFirstLetter("hello"))
	assert.Equal(t, "Hello", CapitalizeFirstLetter("Hello"))
	assert.Equal(t, "", CapitalizeFirstLetter(""))
	assert.Equal(t, "1abc", CapitalizeFirstLetter("1abc"))
}

func TestCreateDynamicStructWithSystemFields(t *testing.T) {
	fields := []*DynamicStructField{
		{Name: "field1", Indexed: true},
		{Name: "field2", Indexed: false},
	}
	typ := CreateDynamicStructWithSystemFields(fields)
	assert.Equal(t, reflect.Struct, typ.Kind())
	assert.Equal(t, 3, typ.NumField())
}

var testChainId = "test_chain_id"

func TestDeleteTablesByChainID(t *testing.T) {
	chainList := []*config.ChainInfo{
		{ChainId: testChainId},
	}
	InitChainListMySqlTable(chainList)

	err := DeleteTablesByChainID(testChainId)
	assert.NoError(t, err)
}

func TestCreateTopicTable(t *testing.T) {
	type tableStructure struct {
		Field1 string `json:"field1"`
	}
	err := CreateTopicTable(UTchainID, "tablename123", tableStructure{})
	assert.NoError(t, err)
}

func TestInitDBTable(t *testing.T) {
	dbConfig := &config.DBConf{DbProvider: config.MySql}
	InitDBTable(dbConfig, []*config.ChainInfo{{ChainId: "c1"}})
}

func TestCreateDatabase(t *testing.T) {
	CreateDatabase(GormDB, "testdb", config.MySql)
}

func TestInitPgsqlTable(t *testing.T) {
	InitPgsqlTable([]*config.ChainInfo{{ChainId: UTchainID}})
}

func TestInitChainListPgsqlTable(t *testing.T) {
	InitChainListPgsqlTable([]*config.ChainInfo{{ChainId: UTchainID}})
}
func TestInitClickHouseTable(t *testing.T) {
	InitClickHouseTable([]*config.ChainInfo{{ChainId: UTchainID}})
}

func TestInitChainLsitClickHouseTable(t *testing.T) {
	InitChainLsitClickHouseTable([]*config.ChainInfo{{ChainId: UTchainID}})
}
