// Package dbhandle 数据库操作
package dbhandle

import (
	"chainmaker_web/src/db"
	"testing"
	"time"

	"chainmaker.org/chainmaker/common/v2/random/uuid"
	"github.com/stretchr/testify/assert"
)

var UTToken = "1234567890"
var UTSign = "1234567890"

func TestInsertUserToken(t *testing.T) {
	uuid := uuid.GetUUID()
	userToken := &db.LoginUserToken{
		Id:         uuid,
		UserAddr:   UserAddr1,
		Token:      UTToken,
		Sign:       UTSign,
		ExpireTime: time.Now().Add(time.Hour * 24 * 30),
	}

	err := InsertUserToken(userToken)
	assert.NoError(t, err)
}

func TestGetTokenInfo(t *testing.T) {
	TestInsertUserToken(t)

	userToken, err := GetTokenInfo(UTToken)
	assert.NoError(t, err)
	assert.NotNil(t, userToken)
}

func TestGetUserLoginBySign(t *testing.T) {
	TestInsertUserToken(t)

	userToken, err := GetUserLoginBySign(UTSign)
	assert.NoError(t, err)
	assert.NotNil(t, userToken)
}

func TestUpdateUserToken(t *testing.T) {
	TestInsertUserToken(t)
	err := UpdateUserToken(UTToken, UTSign, time.Now().Add(time.Hour*24*30))
	assert.NoError(t, err)
}
