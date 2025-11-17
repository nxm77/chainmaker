// Package db 数据库操作
package db

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/entity"
	"fmt"
	"sort"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MySQLHandler MySQL 数据库操作
type MySQLHandler struct {
	// DBConfig 数据库配置
	DBConfig config.DBConf
}

// ConnectDatabase 连接数据库
// @param useDatabase 是否使用数据库
// @return *gorm.DB 数据库连接
// @return error 错误信息
func (handler *MySQLHandler) ConnectDatabase(useDatabase bool) (*gorm.DB, error) {
	dsn := handler.DBConfig.ToMysqlUrl(useDatabase)
	log.Infof("ConnectDatabase ToMysqlUrl========dsn:%v", dsn)
	gormDBConn, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,
		DontSupportRenameColumn:   true,  // rename column not supported before clickhouse 20.4
		SkipInitializeWithVersion: false, // smart configure based on used version
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	sqlDB, err := gormDBConn.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(config.DbMaxIdleConns)
	sqlDB.SetMaxOpenConns(config.DbMaxOpenConns)
	return gormDBConn, nil
}

// InsertWithNativeSQL 插入数据
// @param tableName 表名
// @param records   记录列表
// @return error 错误信息
func (handler *MySQLHandler) InsertWithNativeSQL(tableName string, records []map[string]interface{}) error {
	if len(records) == 0 {
		return nil
	}

	// 1. 批量构建插入语句
	var (
		placeholders []string
		values       []interface{}
		columns      []string
	)

	// 提取字段顺序（保证多次调用一致性）
	if len(records) > 0 {
		columns = make([]string, 0, len(records[0]))
		for k := range records[0] {
			columns = append(columns, k)
		}
		sort.Strings(columns) // 排序保证字段顺序一致
	}

	// 2. 构建参数化查询
	for _, record := range records {
		var ph []string
		for _, col := range columns {
			ph = append(ph, "?")
			values = append(values, record[col])
		}
		placeholders = append(placeholders, "("+strings.Join(ph, ",")+")")
	}

	//3. 执行原生 SQL
	query := fmt.Sprintf("INSERT INTO `%s` (%s) VALUES %s",
		tableName,
		"`"+strings.Join(columns, "`,`")+"`",
		strings.Join(placeholders, ","),
	)

	return GormDB.Exec(query, values...).Error
}

func (handler *MySQLHandler) GetDecodeEventByABIAndTotal(offset, limit int, chainId, contractAddr, version, topic,
	tableName string, topicColumns []string, searchParams []entity.SearchParam) (
	[]map[string]interface{}, int64, error) {
	var total int64
	results := make([]map[string]interface{}, 0)
	if contractAddr == "" || chainId == "" || version == "" || tableName == "" {
		return results, total, ErrTableParams
	}

	tableName = GetTableName(chainId, tableName)
	where := map[string]interface{}{
		ABISystemFieldContractVer: version,
	}

	// 添加过滤条件
	for _, param := range searchParams {
		where[param.Name] = param.Value
	}

	// 构建查询条件
	query := GormDB.Table(tableName).
		Select(strings.Join(topicColumns, ", ")). // 显式指定字段
		Where(where)

	// 执行查询并获取结果
	err := query.Order(ABISystemFieldTimestamp + " desc").Offset(offset).Limit(limit).Find(&results).Error
	if err != nil {
		return results, total, err
	}

	// 计算总记录数
	err = GormDB.Table(tableName).
		Where(where).
		Count(&total).Error

	if err != nil {
		return results, total, err
	}

	return results, total, nil
}
