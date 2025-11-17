// Package db 数据库操作
package db

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/entity"
	"fmt"
	"sort"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// PostgreSQLHandler PostgreSQL 数据库操作
type PostgreSQLHandler struct {
	DBConfig config.DBConf
}

// ConnectDatabase 连接数据库
func (handler *PostgreSQLHandler) ConnectDatabase(useDatabase bool) (*gorm.DB, error) {
	dsn := handler.DBConfig.ToPgsqlUrl(useDatabase)
	log.Infof("ConnectDatabase ToPgsqlUrl========dsn:%v", dsn)
	gormDBConn, err := gorm.Open(postgres.Open(dsn),
		&gorm.Config{
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
func (handler *PostgreSQLHandler) InsertWithNativeSQL(tableName string, records []map[string]interface{}) error {
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

	// 3. 执行原生 SQL
	query := fmt.Sprintf(`INSERT INTO "%s" ("%s") VALUES %s`, // 使用双引号
		tableName,
		strings.Join(columns, `","`), // 正确拼接字段
		strings.Join(placeholders, ","),
	)
	return GormDB.Exec(query, values...).Error
}

func (handler *PostgreSQLHandler) GetDecodeEventByABIAndTotal(offset, limit int, chainId, contractAddr, version, topic,
	tableName string, topicColumns []string, searchParams []entity.SearchParam) (
	[]map[string]interface{}, int64, error) {
	var total int64
	results := make([]map[string]interface{}, 0)
	if contractAddr == "" || chainId == "" || version == "" || tableName == "" {
		return results, total, ErrTableParams
	}

	tableName = GetTableName(chainId, tableName)
	// 获取字段列表并过滤系统字段
	quotedColumns := make([]string, 0, len(topicColumns)) // 带引号的字段名（用于 SELECT）
	for _, column := range topicColumns {
		quotedColumns = append(quotedColumns, fmt.Sprintf("\"%s\"", column)) // 带引号字段名（用于查询）
	}

	where := map[string]interface{}{
		ABISystemFieldContractVer: version,
	}

	// 添加过滤条件
	for _, param := range searchParams {
		where[param.Name] = param.Value
	}

	// 构建查询条件
	query := GormDB.Table(tableName).
		Select(strings.Join(quotedColumns, ", ")). // 显式指定字段
		Where(where)

	// 执行查询并获取结果
	err := query.Offset(offset).Limit(limit).Find(&results).Error
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
