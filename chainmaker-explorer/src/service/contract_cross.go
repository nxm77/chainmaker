package service

import (
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/entity"

	"chainmaker.org/chainmaker/pb-go/v2/syscontract"
	"github.com/gin-gonic/gin"
)

// GetCrossContractCallsHandler get
type GetCrossContractCallsHandler struct {
}

// Handle GetCrossContractCallsHandler
func (handler *GetCrossContractCallsHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetCrossContractCallsHandler(ctx)
	if params == nil || !params.IsLegal() {
		ConvergeFailureResponse(ctx, entity.GetErrorMsgParams())
		return
	}

	// 根据链ID和合约地址查询合约
	contract, err := dbhandle.GetContractByAddr(params.ChainId, params.ContractAddr)
	if err != nil || contract == nil {
		newError := entity.GetErrorMsg(entity.ErrorContractNotExist)
		ConvergeHandleFailureResponse(ctx, newError)
		return
	}

	contractCrossCalls, err := dbhandle.GetContractCrossCallsByName(params.ChainId, contract.NameBak)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	targetContracts := make([]string, 0)
	invokingContracts := make([]string, 0)
	invokingContractSet := make(map[string]struct{}) // 用于去重
	targetContractSet := make(map[string]struct{})   // 用于去重
	methodTargetMap := make(map[string]entity.MethodTargetContract)

	// 遍历合约跨合约调用
	for _, invoking := range contractCrossCalls {
		if IsSystemContract(invoking.InvokingContract) || IsSystemContract(invoking.TargetContract) {
			continue
		}

		if invoking.TargetContract == contract.NameBak {
			if _, exists := invokingContractSet[invoking.InvokingContract]; !exists {
				invokingContractSet[invoking.InvokingContract] = struct{}{}
				invokingContracts = append(invokingContracts, invoking.InvokingContract)
			}
		}

		if invoking.InvokingContract == contract.NameBak {
			if _, exists := targetContractSet[invoking.TargetContract]; !exists {
				targetContractSet[invoking.TargetContract] = struct{}{}
				targetContracts = append(targetContracts, invoking.TargetContract)
			}

			if _, exists := methodTargetMap[invoking.InvokingMethod]; !exists {
				methodTargetMap[invoking.InvokingMethod] = entity.MethodTargetContract{
					Method: invoking.InvokingMethod,
				}
			}

			mtc := methodTargetMap[invoking.InvokingMethod]
			mtc.TargetContracts = append(mtc.TargetContracts, invoking.TargetContract)
			methodTargetMap[invoking.InvokingMethod] = mtc
		}
	}

	// 将方法调用目标转换为切片
	methodTargetsView := make([]entity.MethodTargetContract, 0, len(methodTargetMap))
	for _, mtc := range methodTargetMap {
		methodTargetsView = append(methodTargetsView, mtc)
	}

	crossCallsView := &entity.ContractCrossCallView{
		CurrentContract:      contract.NameBak,
		InvokingContracts:    invokingContracts,
		TargetContracts:      targetContracts,
		MethodTargetContract: methodTargetsView,
	}

	ConvergeDataResponse(ctx, crossCallsView, nil)
}

// 剔除掉系统合约
func IsSystemContract(contractName string) bool {
	return contractName == syscontract.SystemContract_CONTRACT_MANAGE.String()
}
