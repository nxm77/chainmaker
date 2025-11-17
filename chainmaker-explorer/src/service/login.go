/*
Package service comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package service

import (
	"chainmaker_web/src/auth"
	"chainmaker_web/src/entity"
	"chainmaker_web/src/logic"

	"github.com/gin-gonic/gin"
)

// PluginLoginHandler Handler
type PluginLoginHandler struct{}

// RequiresAuth 是否需要登录验证
func (p *PluginLoginHandler) RequiresAuth() bool {
	return false
}

// Handle 浏览器插件登录
func (handler *PluginLoginHandler) Handle(ctx *gin.Context) {
	//参数验证
	params := entity.BindPluginLoginHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, entity.ErrorMsgParam)
		ConvergeFailureResponse(ctx, newError)
		return
	}

	log.Infof("PluginLoginHandler params : %v", params)
	//验证插件的登录加密数据，公钥key和加密sign值，sign解密userAddr
	isValid, userAddr, err := logic.VerifyPluginLogin(params.PubKey, params.SignBase64)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	//验证失败，登录无效
	if !isValid {
		newError := entity.NewError(entity.ErrorParamWrong, entity.ErrorMsgParam)
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//生成token，并存储到token表
	signedToken, err := logic.CreateAndSaveToken(userAddr, params.SignBase64)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	loginView := &entity.PluginLoginView{
		Token:    signedToken,
		UserAddr: userAddr,
	}
	//返回response
	ConvergeDataResponse(ctx, loginView, nil)
}

// AccountLoginHandler Handler
type AccountLoginHandler struct{}

// RequiresAuth 是否需要登录验证1
func (p *AccountLoginHandler) RequiresAuth() bool {
	return false
}

// Handle 浏览器账户登录
func (handler *AccountLoginHandler) Handle(ctx *gin.Context) {
	//参数验证
	params := entity.BindAccountLoginHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, entity.ErrorMsgParam)
		ConvergeFailureResponse(ctx, newError)
		return
	}

	log.Infof("AccountLoginHandler params: %v", params)

	//验证插件的登录加密数据
	isValid, userAddr, err := logic.VerifyAccountLogin(params.RandomNum, params.Password)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	//验证失败，登录无效
	if !isValid {
		newError := entity.NewError(entity.ErrorParamWrong, entity.ErrorMsgParam)
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//生成token，并存储到token表
	signedToken, err := logic.CreateAndSaveToken(userAddr, params.Password)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	loginView := &entity.PluginLoginView{
		Token:    signedToken,
		UserAddr: userAddr,
	}
	//返回response
	ConvergeDataResponse(ctx, loginView, nil)
}

// LogoutHandler Handler
type LogoutHandler struct{}

// RequiresAuth 是否需要登录验证
func (p *LogoutHandler) RequiresAuth() bool {
	return true
}

// Handle 退出登录
func (handler *LogoutHandler) Handle(ctx *gin.Context) {
	// 获取 userId 和 token 并进行类型断言
	userAddr, token, exists := auth.GetUserAddrAndToken(ctx)
	if !exists {
		// userAddr 不存在或类型断言失败，处理错误
		newError := entity.GetErrorNoPermission("")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//将token过期时间改为当前时间
	err := logic.UpdateTokenExpireTime(userAddr, token)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	//返回response
	ConvergeDataResponse(ctx, "OK", nil)
}

// CheckLoginHandler Handler
type CheckLoginHandler struct{}

// RequiresAuth 是否需要登录验证
func (p *CheckLoginHandler) RequiresAuth() bool {
	return true
}

// Handle 退出登录
func (handler *CheckLoginHandler) Handle(ctx *gin.Context) {
	// 获取 userAddr 和 token 并进行类型断言
	userAddr, _, exists := auth.GetUserAddrAndToken(ctx)
	if !exists {
		// userAddr 不存在或类型断言失败，处理错误
		newError := entity.GetErrorNotLogged()
		ConvergeFailureResponse(ctx, newError)
		return
	}

	loginView := &entity.CheckLoginView{
		UserAddr: userAddr,
	}
	//返回response
	ConvergeDataResponse(ctx, loginView, nil)
}
