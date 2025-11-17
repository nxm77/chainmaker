/*
Package db comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package db

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/entity"
	loggers "chainmaker_web/src/logger"
	"database/sql"
	"fmt"
	"reflect"
	"unicode"

	"gorm.io/driver/clickhouse"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var (
	// SqlDB DB client
	SqlDB *sql.DB
	//GormDB DB client
	GormDB    *gorm.DB
	log       = loggers.GetLogger(loggers.MODULE_WEB)
	DBHandler DatabaseHandler
)

type DynamicStructField struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Indexed bool   `json:"indexed,omitempty"`
}

// DatabaseHandler 数据库处理器
type DatabaseHandler interface {
	// ConnectDatabase 连接数据库
	ConnectDatabase(useDatabase bool) (*gorm.DB, error)
	InsertWithNativeSQL(tableName string, records []map[string]interface{}) error
	GetDecodeEventByABIAndTotal(offset, limit int, chainId, contractAddr, version, topic,
		tableName string, topicColumns []string, searchParams []entity.SearchParam) (
		[]map[string]interface{}, int64, error)
}

func GetDatabaseHandler() (DatabaseHandler, error) {
	dbConfig := config.GlobalConfig.DBConf
	switch dbConfig.DbProvider {
	case config.MySql:
		return &MySQLHandler{DBConfig: *dbConfig}, nil
	case config.Pgsql:
		return &PostgreSQLHandler{DBConfig: *dbConfig}, nil
	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s", dbConfig.DbProvider)
	}
}

// InitDbConn init database connection
func InitDbConn(dbConfig *config.DBConf) {
	// 根据数据库类型选择处理器
	dbHandler, err := GetDatabaseHandler()
	if err != nil {
		panic(err)
	}

	// 设置全局数据库处理器
	if dbHandler != nil {
		DBHandler = dbHandler
	}

	// 创建 MySQL 和 ClickHouse 数据库连接
	GormDB, err = ConnectDatabase(dbConfig, true)
	if err != nil {
		//创建数据库
		GormDB, err = ConnectDatabase(dbConfig, false)
		if err != nil {
			log.Errorf("failed to connect database: %v", err)
			panic(err)
		}
		CreateDatabase(GormDB, dbConfig.Database, dbConfig.DbProvider)
		//重新连接数据库
		GormDB, err = ConnectDatabase(dbConfig, true)
		if err != nil {
			panic(err)
		}
	}

	//初始化表结构
	//InitTableName(dbConfig.Prefix)
	//初始化数据库
	InitDBTable(dbConfig, config.SubscribeChains)
}

// InitDBTable 初始化数据库表
// @param dbConfig 数据库配置
// @param chainList 链列表
func InitDBTable(dbConfig *config.DBConf, chainList []*config.ChainInfo) {
	if len(chainList) == 0 {
		return
	}
	//初始化表
	switch dbConfig.DbProvider {
	case config.MySql:
		InitMysqlTable(chainList)
	case config.ClickHouse:
		InitClickHouseTable(chainList)
	case config.Pgsql:
		InitPgsqlTable(chainList)
	}
}

// ConnectDatabase 连接数据库
func ConnectDatabase(dbConfig *config.DBConf, useDataBase bool) (*gorm.DB, error) {
	var err error
	switch dbConfig.DbProvider {
	case config.MySql:
		dsn := dbConfig.ToMysqlUrl(useDataBase)
		GormDB, err = gorm.Open(mysql.New(mysql.Config{
			DSN:                       dsn,
			DontSupportRenameColumn:   true,  // rename column not supported before clickhouse 20.4
			SkipInitializeWithVersion: false, // smart configure based on used version
		}), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		})
	case config.ClickHouse:
		dsn := dbConfig.ToClickHouseUrl(useDataBase)
		GormDB, err = gorm.Open(clickhouse.New(clickhouse.Config{
			DSN:                          dsn,
			DontSupportRenameColumn:      true,  // rename column not supported before clickhouse 20.4
			DontSupportEmptyDefaultValue: false, // do not consider empty strings as valid default values
			SkipInitializeWithVersion:    false, // smart configure based on used version
		}), &gorm.Config{
			//Logger: logger.Default.LogMode(logger.Info),
		})
	case config.Pgsql:
		dsn := dbConfig.ToPgsqlUrl(useDataBase)
		GormDB, err = gorm.Open(postgres.Open(dsn),
			&gorm.Config{
				Logger: logger.Default.LogMode(logger.Info),
			})

		GormDB.NamingStrategy = schema.NamingStrategy{
			NoLowerCase: true, // 关键配置
		}
	}

	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	sqlDB, _ := GormDB.DB()
	sqlDB.SetMaxIdleConns(config.DbMaxIdleConns)
	sqlDB.SetMaxOpenConns(config.DbMaxOpenConns)
	return GormDB, nil
}

// CreateDatabase 创建数据库
func CreateDatabase(db *gorm.DB, database, dbProvider string) {
	var createDatabaseQuery string
	// 创建数据库
	switch dbProvider {
	case config.MySql:
		// 如果数据库不存在，则创建数据库
		createDatabaseQuery = fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s %s",
			database, config.MysqlDatabaseConf)
	case config.ClickHouse:
		// 如果数据库不存在，则创建数据库
		createDatabaseQuery = fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", database)
	case config.Pgsql:
		// 创建数据库
		createDatabaseQuery = fmt.Sprintf("CREATE DATABASE %s", database)
	default:
		// 如果数据库类型不匹配，则返回
		return
	}

	err := db.Exec(createDatabaseQuery).Error
	log.Infof("CREATE DATABASE %v", database)
	if err != nil {
		log.Errorf("CREATE DATABASE failed, err:%v", err)
	}
}

// InitMysqlTable 初始化数据库
// @param chainList 链信息列表
func InitMysqlTable(chainList []*config.ChainInfo) {
	err := GormDB.AutoMigrate(
		&Chain{},
		&Subscribe{},
		&LoginUserToken{},
		&Statistics{},
	)
	if err != nil {
		log.Errorf("AutoMigrate failed, err:%v", err)
	}
}

func InitChainListMySqlTable(chainList []*config.ChainInfo) {
	//其他表按链ID分表
	blockTableNames := GetBlockTableNames()
	for _, chainInfo := range chainList {
		for _, tableInfo := range blockTableNames {
			err := SqlCreateTableWithComment(GormDB, chainInfo.ChainId, tableInfo)
			if err != nil {
				log.Errorf("SqlCreateTableWithComment failed, err:%v", err)
			}
		}
	}
}

// InitPgsqlTable 初始化数据库
// @param chainList 链信息列表
func InitPgsqlTable(chainList []*config.ChainInfo) {
	err := GormDB.AutoMigrate(
		&Chain{},
		&Subscribe{},
		&LoginUserToken{},
		&Statistics{},
	)
	if err != nil {
		log.Errorf("AutoMigrate failed, err:%v", err)
	}
}

func InitChainListPgsqlTable(chainList []*config.ChainInfo) {
	//其他表按链ID分表
	blockTableNames := GetBlockTableNames()
	for _, chainInfo := range chainList {
		for _, tableInfo := range blockTableNames {
			err := PGSqlCreateTableWithComment(GormDB, chainInfo.ChainId, tableInfo)
			if err != nil {
				log.Errorf("PGSqlCreateTableWithComment failed, err:%v", err)
			}
		}
	}
}

// InitClickHouseTable 初始化数据库
func InitClickHouseTable(chainList []*config.ChainInfo) {
	err := GormDB.AutoMigrate(
		&Chain{},
		&Subscribe{},
		&LoginUserToken{},
		&Statistics{},
	)
	if err != nil {
		panic(err)
	}
}

func InitChainLsitClickHouseTable(chainList []*config.ChainInfo) {
	//其他表按链ID分表
	blockTableNames := GetBlockTableNames()
	for _, chainInfo := range chainList {
		for _, tableInfo := range blockTableNames {
			err := ClickHouseCreateTableWithComment(chainInfo.ChainId, tableInfo)
			if err != nil {
				panic(err)
			}
		}
	}
}

// DeleteTablesByChainID 根据链ID删除相关表
func DeleteTablesByChainID(chainId string) error {
	// 获取区块表名
	blockTableNames := GetBlockTableNames()
	// 遍历区块表名
	for _, tableInfo := range blockTableNames {
		// 根据链ID和表名获取表名
		tableName := GetTableName(chainId, tableInfo.Name)
		// 执行删除表的SQL语句
		err := GormDB.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)).Error
		// 如果有错误，返回错误
		if err != nil {
			return err
		}
	}

	// 没有错误，返回nil
	return nil
}

func CreateTopicTable(chainId, tableName string, tableStructure interface{}) error {
	if chainId == "" || tableName == "" {
		return fmt.Errorf("tableName is empty")
	}

	dbConfig := config.GlobalConfig.DBConf
	tableInfo := &TableInfo{
		Name:        tableName,
		Structure:   tableStructure,
		Description: "合约事件信息解析表",
	}

	switch dbConfig.DbProvider {
	case config.MySql:
		err := SqlCreateTableWithComment(GormDB, chainId, *tableInfo)
		log.Errorf("SqlCreateTableWithComment chainId:%s, tableName:%s, err:%v", chainId, tableName, err)
	case config.ClickHouse:
		err := ClickHouseCreateTableWithComment(chainId, *tableInfo)
		log.Errorf("ClickHouseCreateTableWithComment chainId:%s, tableName:%s, err:%v", chainId, tableName, err)
	case config.Pgsql:
		err := PGSqlCreateTableWithComment(GormDB, chainId, *tableInfo)
		log.Errorf("PGSqlCreateTableWithComment chainId:%s, tableName:%s, err:%v", chainId, tableName, err)
	}
	return nil
}

func CreateDynamicStructWithSystemFields(eventFields []*DynamicStructField) reflect.Type {
	// 固定系统字段
	structFields := []reflect.StructField{
		{
			Name: "System",
			Type: reflect.TypeOf(ABIEventSystemFields{}),
			Tag:  `gorm:"embedded"`, // 嵌入式结构体
		},
	}

	// 动态事件字段
	for _, field := range eventFields {
		var tags string
		if field.Indexed {
			tags = fmt.Sprintf(`gorm:"column:%s;type:varchar(255);index;"`, field.Name)
		} else {
			tags = fmt.Sprintf(`gorm:"column:%s"`, field.Name)
		}
		structFields = append(structFields, reflect.StructField{
			Name: CapitalizeFirstLetter(field.Name),
			Type: reflect.TypeOf(""), // 默认类型为 string
			Tag:  reflect.StructTag(tags),
		})
	}

	return reflect.StructOf(structFields)
}

// CapitalizeFirstLetter 首字母大写
// @param s string 字符串
// @return string 首字母大写的字符串
func CapitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}

	// 将首字母转为大写，并确保其他字符不受影响
	r := []rune(s) // 将字符串转换为 rune 切片，处理多字节字符
	if unicode.IsLower(r[0]) {
		r[0] = unicode.ToUpper(r[0]) // 将首字母大写
	}

	return string(r)
}
