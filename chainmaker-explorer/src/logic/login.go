// Package logic 逻辑处理
package logic

import (
	"chainmaker_web/src/auth"
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/entity"
	"chainmaker_web/src/utils"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"chainmaker.org/chainmaker/common/v2/crypto"
	"chainmaker.org/chainmaker/common/v2/crypto/asym"
	pbconfig "chainmaker.org/chainmaker/pb-go/v2/config"
	commonutils "chainmaker.org/chainmaker/utils/v2"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var optSM3 = &crypto.SignOpts{
	Hash: crypto.HASH_TYPE_SM3,
	UID:  crypto.CRYPTO_DEFAULT_UID,
}

// CreateAndSaveToken 根据userID生成一个token，然后存储到token表
func CreateAndSaveToken(userAddr, signStr string) (string, error) {
	//获取过期时间配置
	webConf := config.GlobalConfig.WebConf
	expireData := webConf.LoginExpireTime
	//默认过期时间7天
	if expireData == 0 {
		expireData = 7
	}

	expireTime := time.Now().Add(time.Hour * 24 * time.Duration(expireData))
	//生成token
	tokenStr, err := auth.CreateSignedToken(userAddr, expireTime)
	if err != nil {
		return tokenStr, err
	}

	newUUID := uuid.New().String()
	userToken := &db.LoginUserToken{
		Id:         newUUID,
		UserAddr:   userAddr,
		Token:      tokenStr,
		Sign:       signStr,
		ExpireTime: expireTime,
	}
	//存储token
	err = dbhandle.InsertUserToken(userToken)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

// UpdateTokenExpireTime  退出登录，将token改成已失效
func UpdateTokenExpireTime(userAddr, token string) error {
	expireTime := time.Now()
	err := dbhandle.UpdateUserToken(userAddr, token, expireTime)
	return err
}

// VerifyPluginLogin 验证加密数据
func VerifyPluginLogin(publicKey, signBase64 string) (bool, string, error) {
	//查看signBase64是否存在，一个sign值只能登录一次
	userLogin, err := dbhandle.GetUserLoginBySign(signBase64)
	if err != nil {
		return false, "", err
	}

	//已经登录一次，不在能继续登录
	if userLogin != nil {
		return false, "", fmt.Errorf("signBase64 has expired")
	}

	// verify sig
	pk, err := asym.PublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		return false, "", err
	}

	//publicKey解析pkAddr
	pkAddr, err := commonutils.PkToAddrStr(pk, pbconfig.AddrType_ETHEREUM, crypto.HASH_TYPE_SM3)
	if err != nil {
		return false, pkAddr, err
	}

	signBytes, err := base64.StdEncoding.DecodeString(signBase64)
	if err != nil {
		return false, pkAddr, err
	}

	addrBytes, _ := hex.DecodeString(pkAddr)
	//验证pkAddr和加密的signBytes是否一致
	isValid, err := pk.VerifyWithOpts(addrBytes, signBytes, optSM3)
	if err != nil {
		return false, pkAddr, err
	}

	return isValid, pkAddr, nil
}

// VerifyAccountLogin 验证账户登录数据
func VerifyAccountLogin(randomNum int64, passwordMd5 string) (bool, string, error) {
	//验证登录是否有效
	md5Hash := utils.GetMd5Hash(randomNum)
	if md5Hash != passwordMd5 {
		return false, "", fmt.Errorf(entity.ErrorPasswordError)
	}

	//查看signBase64是否存在，一个sign值只能登录一次
	userLogin, err := dbhandle.GetUserLoginBySign(passwordMd5)
	if err != nil {
		return false, "", err
	}

	//已经登录一次，不在能继续登录
	if userLogin != nil {
		return false, "", fmt.Errorf("signBase64 has expired")
	}

	pkAddr := utils.GetAccountHashStr()
	return true, pkAddr, nil
}

// CheckAdminLogin 检查管理员登录
func CheckAdminLogin(ctx *gin.Context) bool {
	//获取用户地址和token
	userAddr, _, exists := auth.GetUserAddrAndToken(ctx)
	//如果用户未登录
	if !exists {
		//登录验证,输出未登录日志
		log.Errorf("CheckAdminLogin user not login")
		return false
	}
	adminAddr := utils.GetAccountHashStr()
	log.Infof("CheckAdminLogin userAddr: %s, adminAddr: %s", userAddr, adminAddr)
	return adminAddr == userAddr
}
