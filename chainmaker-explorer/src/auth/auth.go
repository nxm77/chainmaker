// Package auth 登录中间件
package auth

import (
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/entity"
	loggers "chainmaker_web/src/logger"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	//SecretKeyStr 	签名Key
	SecretKeyStr = "explorer-backend"
	//log 日志
	log = loggers.GetLogger(SecretKeyStr)
	//TokenBearer token 前缀
	TokenBearer = "Bearer "
)

// TokenAuthMiddleware 登录验证插件，用于需要进行登录验证的接口，根据RequiresAuth判断
func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 监测登录是否有效
		userAddr, tokenString, err := RegularlyVerifyToken(c)
		if err != nil {
			entity.ConvergeFailureResponse(c, err)
			c.Abort()
			return
		}

		// 将 userId 附加到 gin.Context 上，以便在后续的处理程序中使用
		c.Set("userAddr", userAddr)
		c.Set("token", tokenString)
		c.Next()
	}
}

func RegularlyVerifyToken(c *gin.Context) (string, string, *entity.Error) {
	//登录token信息
	authHeader := c.Request.Header.Get("Authorization")
	log.Infof("TokenAuthMiddleware authHeader:%v", authHeader)
	if authHeader == "" {
		//未登录
		newError := entity.GetErrorNotLogged()
		return "", "", newError
	}

	// 检查header是否以"Bearer "开头
	if !strings.HasPrefix(authHeader, TokenBearer) {
		newError := entity.GetErrorNoPermission(entity.ErrorMsgTokenFormat)
		return "", "", newError
	}

	// 提取token字符串
	tokenString := strings.TrimPrefix(authHeader, TokenBearer)
	if tokenString == "" {
		newError := entity.GetErrorNotLogged()
		return "", "", newError
	}

	//校验token是否有效，解析出userAddr
	userAddr, err := VerifyToken(tokenString)
	if err != "" {
		newError := entity.GetErrorNoPermission(err)
		return "", "", newError
	}

	//校验登录是否过期
	isActive := checkTokenExpire(tokenString)
	if !isActive {
		newError := entity.GetErrorNoPermission(entity.ErrorMsgTokenExpired)
		return "", "", newError
	}

	return userAddr, tokenString, nil
}

// GetUserIDAndToken 获取token解析出的userId, token, 并校验token是否有效
// @param ctx gin.Context
// @return userId userId
// @return token token
// @return isValid 是否有效
func GetUserAddrAndToken(ctx *gin.Context) (string, string, bool) {
	// 获取 userId 和 token 并进行类型断言
	userAddrInterface, exists := ctx.Get("userAddr")
	tokenInterface, exists2 := ctx.Get("token")
	userAddr, ok1 := userAddrInterface.(string)
	token, ok2 := tokenInterface.(string)
	if !exists || !exists2 || !ok1 || !ok2 || userAddr == "" || token == "" {
		return userAddr, token, false
	}
	return userAddr, token, true
}

// VerifyToken 校验token是否有效，解析出userID
func VerifyToken(tokenString string) (string, string) {
	secretKey := []byte(SecretKeyStr)
	// 解析token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 检查签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("无效的签名方法: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return "", entity.ErrorMsgTokenExpired
	}

	// 验证token是否有效
	if !token.Valid {
		return "", entity.ErrorMsgTokenInvalid
	}

	// 提取userId
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", entity.ErrorMsgTokenParseUserID
	}

	userAddr, ok := claims["userAddr"].(string)
	if !ok {
		return "", entity.ErrorMsgTokenParseUserID
	}

	return userAddr, ""
}

// CreateSignedToken  jwt根据userID创建带失效时间的token
// @param userId 用户id
// @param expireTime 失效时间
// @return string token
// @return error 错误
// CreateSignedToken 函数用于创建并签名 Token
func CreateSignedToken(userAddr string, expireTime time.Time) (string, error) {
	// 从配置文件中获取密钥
	secretKey := []byte(SecretKeyStr)
	// 创建并签名 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userAddr": userAddr,
		"exp":      expireTime.Unix(),
	})

	// 使用密钥签名 Token
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		log.Errorf("签名 Token 时出错：%v", err)
		return "", err
	}

	return signedToken, nil
}

// checkTokenExpire 根据失效时间检查token是否有效
// @param token token
// @return bool 是否有效
func checkTokenExpire(token string) bool {
	//根据token获取登录信息
	userToken, err := dbhandle.GetTokenInfo(token)
	if err != nil || userToken == nil {
		return false
	}

	//判断登录是否过期
	nowTime := time.Now().Unix()
	return userToken.ExpireTime.Unix() >= nowTime
}
