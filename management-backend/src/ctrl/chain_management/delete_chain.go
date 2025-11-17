package chain_management

import (
	loggers "management_backend/src/logger"
	"management_backend/src/sync"

	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
	dbchain "management_backend/src/db/chain"
	"management_backend/src/entity"
)

// DeleteChainHandler delete chain
type DeleteChainHandler struct{}

// LoginVerify login verify
func (deleteChainHandler *DeleteChainHandler) LoginVerify() bool {
	return true
}

// Handle deal
//
//	@Description:
//	@receiver deleteChainHandler
//	@param user
//	@param ctx
func (deleteChainHandler *DeleteChainHandler) Handle(user *entity.User, ctx *gin.Context) {
	params := BindDeleteChainHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}

	err := dbchain.DeleteChain(params.ChainId)
	if err != nil {
		loggers.WebLogger.Error("DeleteChain err : " + err.Error())
		common.ConvergeFailureResponse(ctx, common.ErrorDeleteChain)
		return
	}
	sync.StopSubscribe(params.ChainId)
	common.ConvergeDataResponse(ctx, common.NewStatusResponse(), nil)
}
