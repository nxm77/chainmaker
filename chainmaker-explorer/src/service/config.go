package service

import (
	"chainmaker_web/src/utils"

	"github.com/gin-gonic/gin"
)

// GetChainConfigHandler get
type GetChainConfigHandler struct {
}

// Handle GetChainConfigHandler是否展示订阅按钮
func (getChainConfigHandler *GetChainConfigHandler) Handle(ctx *gin.Context) {
	ConvergeDataResponse(ctx, utils.GetConfigShow(), nil)
}
