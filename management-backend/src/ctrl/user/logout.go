/*
Package user comment
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package user

import (
	loggers "management_backend/src/logger"

	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
	"management_backend/src/entity"
	"management_backend/src/session"
)

// LogoutHandler logout
type LogoutHandler struct{}

// LoginVerify login verify
func (handler *LogoutHandler) LoginVerify() bool {
	return true
}

// Handle deal
func (handler *LogoutHandler) Handle(user *entity.User, ctx *gin.Context) {
	// 将session清空
	err := session.UserStoreClear(ctx)
	if err != nil {
		loggers.SessionLogger.Debugf("clean captcha failed")
	}
	common.ConvergeDataResponse(ctx, common.NewStatusResponse(), nil)
}
