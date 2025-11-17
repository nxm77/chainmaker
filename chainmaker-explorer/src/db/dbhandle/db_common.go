/*
Package dbhandle comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package dbhandle

import (
	"chainmaker_web/src/db"
	loggers "chainmaker_web/src/logger"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-sql-driver/mysql"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	SQLDesc = "desc"
	SQLAsc  = "asc"
)
const (
	//DelayUpdateFail 异步更新失败
	DelayUpdateFail = 0
	//DelayUpdateSuccess 异步更新成功
	DelayUpdateSuccess = 1
)

const (
	//ContractResultSuccess 合约执行成功
	ContractResultSuccess = 0
	//ContractResultFailed 合约执行失败
	ContractResultFailed = 1
)

const (
	MysqlErrPrimaryKeyDuplication = 1062
	MysqlErrUniqueDuplication     = 1869
)

const (
	// SystemContractStatus 系统合约状态
	SystemContractStatus = -1
	//ContractStatusOK 合约正常
	ContractStatusOK = 0
	//ContractStatusFrozen 合约冻结
	ContractStatusFrozen = 1
	//ContractStatusCancelled 合约取消
	ContractStatusCancelled = 2
)

type QueryCondition struct {
	Field     string
	Value     interface{}
	Condition string
	Operator  string
}

// SelectFile SelectFile
type SelectFile struct {
	Where        map[string]interface{}
	CurrentWhere map[string]interface{}
	NotWhere     map[string]interface{}
	WhereIn      map[string]interface{}
	NotWhereIn   map[string]interface{}
	WhereOr      map[string]interface{}
	NotNull      []string
	StartTime    int64
	EndTime      int64
}

var (
	// log 日志
	log = loggers.GetLogger(loggers.MODULE_WEB)
)

// InsertData
//
//	@Description: 单个插入数据
//	@param tableName
//	@param data
//	@return error
func InsertData(tableName string, data interface{}) error {
	err := db.GormDB.Table(tableName).Create(data).Error
	if err == nil {
		return nil
	}

	//检查错误是否为主键冲突或唯一索引冲突
	isDuplicate := isDuplicateKeyError(err)
	if isDuplicate {
		return nil
	}

	return err
}

// CreateInBatchesData
//
//	@Description: 批量插入数据
//	@param tableName
//	@param data
//	@return error
func CreateInBatchesData(tableName string, data interface{}) error {
	err := db.GormDB.Table(tableName).Create(data).Error
	if err == nil {
		return nil
	}
	//检查错误是否为主键冲突或唯一索引冲突
	isDuplicate := isDuplicateKeyError(err)
	if isDuplicate {
		//改成单个插入
		v := reflect.ValueOf(data)
		if v.Kind() != reflect.Slice {
			log.Errorf("【sql】 data must be of type []interface{}")
			return nil
		}

		//如果只有一条，直接退出
		if v.Len() == 1 {
			return nil
		}

		for i := 0; i < v.Len(); i++ {
			err = db.GormDB.Table(tableName).Create(v.Index(i).Interface()).Error
			if isDuplicateKeyError(err) {
				continue
			} else if err != nil {
				log.Errorf("【sql】 CreateInBatchesData error:%v", err)
				return nil
			}
		}
		return nil
	}

	return nil
}

// isDuplicateKeyError 检查错误是否为主键冲突或唯一索引冲突
func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	// 在这里，您可以检查错误是否为主键冲突或唯一索引冲突
	// 例如，对于MySQL，可以检查错误代码是否为1062（Duplicate entry）
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		if mysqlErr.Number == MysqlErrPrimaryKeyDuplication ||
			mysqlErr.Number == MysqlErrUniqueDuplication {
			return true
		}
	}

	// Check for PostgreSQL duplicate key error
	if strings.Contains(err.Error(), "(SQLSTATE 23505)") {
		return true
	}

	return false
}

// BuildParamsQuery 构造Gorm select参数
func BuildParamsQuery(tableName string, selectFile *SelectFile) *gorm.DB {
	query := db.GormDB.Table(tableName)
	// 添加 Where 条件
	for k, v := range selectFile.Where {
		query = query.Where(fmt.Sprintf("%s = ?", k), v)
	}
	for k, v := range selectFile.CurrentWhere {
		query = query.Where(k, v)
	}

	// 添加 NotWhere 条件
	for key, value := range selectFile.NotWhere {
		query = query.Not(fmt.Sprintf("%s = ?", key), value)
	}

	// 添加 where in 条件
	for k, v := range selectFile.WhereIn {
		data := ToInterfaceSlice(v)
		query = query.Where(fmt.Sprintf("%s IN (?)", k), data)
	}

	// 添加 NotWhereIn 条件
	if len(selectFile.NotWhereIn) > 0 {
		query = query.Not(selectFile.NotWhereIn)
	}

	// 添加 NotNull 条件
	for _, v := range selectFile.NotNull {
		query = query.Where(fmt.Sprintf("%s IS NOT NULL", v))
	}

	// 添加时间范围条件
	if selectFile.StartTime > 0 && selectFile.EndTime > 0 {
		query = query.Where("timestamp BETWEEN ? AND ?", selectFile.StartTime, selectFile.EndTime)
	}

	// 添加 whereOr 条件
	if len(selectFile.WhereOr) > 0 {
		orConditions := make([]string, 0)
		orValues := make([]interface{}, 0)
		for key, value := range selectFile.WhereOr {
			orConditions = append(orConditions, fmt.Sprintf("%s = ?", key))
			orValues = append(orValues, value)
		}
		query = query.Where(gorm.Expr(strings.Join(orConditions, " OR "), orValues...))
	}

	return query
}

// BuildParamsQueryNew
//
//	@Description: 构造查询条件
//	@param query
//	@param selectFile
//	@return *gorm.DB
func BuildParamsQueryNew(query *gorm.DB, selectFile *SelectFile) *gorm.DB {
	//query := db.GormDB.Table(tableName).Select("txId")
	// 添加 Where 条件
	for k, v := range selectFile.Where {
		query = query.Where(fmt.Sprintf("%s = ?", k), v)
	}
	for k, v := range selectFile.CurrentWhere {
		query = query.Where(k, v)
	}

	// 添加 NotWhere 条件
	for key, value := range selectFile.NotWhere {
		query = query.Not(fmt.Sprintf("%s = ?", key), value)
	}

	// 添加 where in 条件
	for k, v := range selectFile.WhereIn {
		data := ToInterfaceSlice(v)
		query = query.Where(fmt.Sprintf("%s IN (?)", k), data)
	}

	// 添加 NotWhereIn 条件
	if len(selectFile.NotWhereIn) > 0 {
		query = query.Not(selectFile.NotWhereIn)
	}

	// 添加 NotNull 条件
	for _, v := range selectFile.NotNull {
		query = query.Where(fmt.Sprintf("%s IS NOT NULL", v))
	}

	// 添加时间范围条件
	if selectFile.StartTime > 0 && selectFile.EndTime > 0 {
		query = query.Where("timestamp BETWEEN ? AND ?", selectFile.StartTime, selectFile.EndTime)
	}

	// 添加 whereOr 条件
	if len(selectFile.WhereOr) > 0 {
		orConditions := make([]string, 0)
		orValues := make([]interface{}, 0)
		for key, value := range selectFile.WhereOr {
			orConditions = append(orConditions, fmt.Sprintf("%s = ?", key))
			orValues = append(orValues, value)
		}
		query = query.Where(gorm.Expr(strings.Join(orConditions, " OR "), orValues...))
	}

	return query
}

// ToInterfaceSlice interface转换
func ToInterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		return []interface{}{slice} // 如果输入不是切片，返回一个包含该值的新切片
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

// GenerateRedisKey
//
//	@Description: 根据查询条件生成缓存key
//	@param selectFile
//	@return string
func GenerateRedisKey(selectFile interface{}) string {
	// 将selectFile转换为JSON字符串
	jsonData, err := json.Marshal(selectFile)
	if err != nil {
		return ""
	}

	// 计算JSON字符串的SHA-256哈希值
	hash := sha256.Sum256(jsonData)

	// 将哈希值转换为十六进制字符串
	key := hex.EncodeToString(hash[:])

	return key
}

func BuildQuery(tableName string, conditions []QueryCondition) *gorm.DB {
	query := db.GormDB.Table(tableName)
	for _, condition := range conditions {
		// 将条件转换为小写以兼容大写和小写
		switch strings.ToLower(condition.Condition) {
		case "=":
			query = query.Where(clause.Eq{
				Column: clause.Column{Name: condition.Field},
				Value:  condition.Value,
			})
		case ">":
			query = query.Where(clause.Expr{
				SQL:  "? > ?",
				Vars: []interface{}{clause.Column{Name: condition.Field}, condition.Value},
			})
		case "<":
			query = query.Where(clause.Expr{
				SQL:  "? < ?",
				Vars: []interface{}{clause.Column{Name: condition.Field}, condition.Value},
			})
		case "like":
			if condition.Field == "contractKey" {
				query = query.Where(clause.Expr{
					SQL: "? LIKE ? OR ? LIKE ?",
					Vars: []interface{}{
						clause.Column{Name: "nameBak"}, condition.Value,
						clause.Column{Name: "addr"}, condition.Value,
					},
				})
			} else {
				query = query.Where(clause.Expr{
					SQL:  "? LIKE ?",
					Vars: []interface{}{clause.Column{Name: condition.Field}, condition.Value},
				})
			}
		case "in":
			query = query.Where(clause.Expr{
				SQL:  "? IN ?",
				Vars: []interface{}{clause.Column{Name: condition.Field}, condition.Value},
			})
		case "between":
			// 确保 condition.Value 是一个包含两个值的切片
			if values, ok := condition.Value.([]int64); ok && len(values) == 2 {
				query = query.Where(clause.Expr{
					SQL:  "? BETWEEN ? AND ?",
					Vars: []interface{}{clause.Column{Name: condition.Field}, values[0], values[1]},
				})
			}
		case SQLDesc:
			query = query.Order(clause.OrderByColumn{
				Column: clause.Column{Name: condition.Field},
				Desc:   true,
			})
		case SQLAsc:
			query = query.Order(clause.OrderByColumn{
				Column: clause.Column{Name: condition.Field},
			})
		}
	}

	return query
}

func BuildQueryNew(query *gorm.DB, conditions []QueryCondition) *gorm.DB {
	// 用于存储"or"条件的子查询
	var orConditions []string
	var orValues []interface{}

	// 用于存储"and"条件的子查询
	var andConditions []string
	var andValues []interface{}

	for _, condition := range conditions {
		switch strings.ToLower(condition.Condition) {
		case SQLDesc:
			query = query.Order(fmt.Sprintf("%s desc", condition.Field))
		case SQLAsc:
			query = query.Order(fmt.Sprintf("%s asc", condition.Field))
		default:
			if condition.Operator == "or" || condition.Operator == "OR" {
				orConditions = append(orConditions, fmt.Sprintf("%s %s ?", condition.Field, condition.Condition))
				orValues = append(orValues, condition.Value)
			} else { // 默认为"and"
				andConditions = append(andConditions, fmt.Sprintf("%s %s ?", condition.Field, condition.Condition))
				andValues = append(andValues, condition.Value)
			}
		}
	}

	// 添加"and"条件
	if len(andConditions) > 0 {
		for i, andCondition := range andConditions {
			query = query.Where(andCondition, andValues[i])
		}
	}

	// 添加 whereOr 条件
	if len(orConditions) > 0 {
		query = query.Where(gorm.Expr(strings.Join(orConditions, " OR "), orValues...))
	}
	return query
}
