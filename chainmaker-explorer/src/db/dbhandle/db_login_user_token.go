// Package dbhandle 数据库操作
package dbhandle

import (
	"chainmaker_web/src/db"
	"time"
)

// InsertUserToken 添加user的token数据
func InsertUserToken(userToken *db.LoginUserToken) error {
	if userToken == nil {
		return nil
	}

	return db.GormDB.Table(db.TableLoginUserToken).Create(userToken).Error
}

// GetTokenInfo 查询token是否已经存入数据库
func GetTokenInfo(token string) (*db.LoginUserToken, error) {
	var tokenInfo *db.LoginUserToken
	where := map[string]interface{}{
		"token": token,
	}
	result := db.GormDB.Table(db.TableLoginUserToken).Where(where).Find(&tokenInfo)

	// 如果没有找到记录，result.RowsAffected将为0
	if result.RowsAffected == 0 {
		return nil, nil
	}

	// 如果查询过程中发生了错误，返回错误
	if result.Error != nil {
		return nil, result.Error
	}

	return tokenInfo, nil
}

// GetUserLoginBySign 查询验签的sign值是否已经存入数据库
func GetUserLoginBySign(sign string) (*db.LoginUserToken, error) {
	var tokenInfo *db.LoginUserToken
	where := map[string]interface{}{
		"sign": sign,
	}
	result := db.GormDB.Table(db.TableLoginUserToken).Where(where).Find(&tokenInfo)

	// 如果没有找到记录，result.RowsAffected将为0
	if result.RowsAffected == 0 {
		return nil, nil
	}

	// 如果查询过程中发生了错误，返回错误
	if result.Error != nil {
		return nil, result.Error
	}

	return tokenInfo, nil
}

// UpdateUserToken 更新sql
func UpdateUserToken(userAddr, token string, expireTime time.Time) error {
	where := map[string]interface{}{
		"userAddr": userAddr,
		"token":    token,
	}
	params := map[string]interface{}{
		"expireTime": expireTime,
	}
	err := db.GormDB.Table(db.TableLoginUserToken).Where(where).Updates(params).Error
	return err
}
