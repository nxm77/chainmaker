/*
Package explorer comment
Copyright (C) BABEC. All rights reserved.
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package explorer

import (
	"fmt"
	"management_backend/src/ctrl/ca"
	"management_backend/src/db"
	"management_backend/src/global"
	loggers "management_backend/src/logger"
	"strings"

	"github.com/emirpasic/gods/lists/arraylist"
	"github.com/gin-gonic/gin"

	"management_backend/src/ctrl/common"
	"management_backend/src/db/chain"
	dbcommon "management_backend/src/db/common"
	"management_backend/src/db/contract"
	contractRecord "management_backend/src/db/contract"
	"management_backend/src/entity"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// GetTxListHandler getTxListHandler
type GetTxListHandler struct {
}

// LoginVerify login verify
func (handler *GetTxListHandler) LoginVerify() bool {
	return true
}

// Handle deal
func (handler *GetTxListHandler) Handle(user *entity.User, ctx *gin.Context) {
	var (
		txList     []*dbcommon.Transaction
		totalCount int64
		offset     int
		limit      int
		addr       string
	)
	params := BindGetTxListHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}
	offset = params.PageNum * params.PageSize
	limit = params.PageSize
	if params.ContractName != "" {
		contractInfo, err := contract.GetContractByName(params.ChainId, params.ContractName)
		if err != nil {
			common.ConvergeHandleFailureResponse(ctx, err)
			return
		}
		addr = contractInfo.EvmAddress
	}

	totalCount, txList, err := chain.GetTxList(params.ChainId, offset, limit, params.BlockHeight,
		params.ContractName, addr)
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		return
	}
	txInfos := convertToTxViews(txList)
	common.ConvergeListResponse(ctx, txInfos, totalCount, nil)
}

func convertToTxViews(txList []*dbcommon.Transaction) []interface{} {
	txViews := arraylist.New()
	for _, tx := range txList {
		contract, err := contractRecord.GetContractByNameOrEvmAddress(tx.ChainId, tx.ContractName, tx.ContractName)
		if err != nil {
			loggers.WebLogger.Errorf("GetContractByNameOrEvmAddress err: %v", err.Error())
		} else {
			if contract != nil {
				tx.ContractName = contract.Name
			}
		}
		record, err := contractRecord.GetInvokeRecordByTxId(tx.TxId)
		if err != nil {
			loggers.WebLogger.Errorf("GetInvokeRecordByTxId error: %v", err.Error())
		} else {
			if record != nil {
				tx.ContractMethod = record.Method
			}
		}
		txView := NewTransactionListView(tx)
		txViews.Add(txView)
	}
	return txViews.Values()
}

// GetTxDetailHandler getTxDetailHandler
type GetTxDetailHandler struct {
}

// LoginVerify login verify
func (handler *GetTxDetailHandler) LoginVerify() bool {
	return true
}

// Handle deal
func (handler *GetTxDetailHandler) Handle(user *entity.User, ctx *gin.Context) {
	params := BindGetTxDetailHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}

	tx, err := handler.getTransaction(params)
	if err != nil {
		common.ConvergeHandleFailureResponse(ctx, err)
		return
	}

	// 检查和处理 EVM 类型的合约
	err = handler.processEvmContract(tx)
	if err != nil {
		loggers.WebLogger.Errorf("processEvmContract error: %v", err.Error())
	}

	txView := NewTransactionView(tx)
	common.ConvergeDataResponse(ctx, txView, nil)
}

func (handler *GetTxDetailHandler) getTransaction(params *GetTxDetailParams) (*dbcommon.Transaction, error) {
	var (
		tx  *dbcommon.Transaction
		err error
	)

	if params.Id != 0 {
		tx, err = chain.GetTxById(params.Id)
	} else if params.TxId != "" {
		tx, err = chain.GetTxByTxId(params.ChainId, params.TxId)
	}
	return tx, err
}

func (handler *GetTxDetailHandler) processEvmContract(tx *dbcommon.Transaction) error {
	contract, err := contractRecord.GetContractByNameOrEvmAddress(tx.ChainId, tx.ContractName, tx.ContractName)
	if err != nil {
		loggers.WebLogger.Errorf("GetContractByNameOrEvmAddress error: %v", err.Error())
		return err
	}

	if contract != nil && contract.RuntimeType == global.EVM {
		id, userId, hash, err := ca.ResolveUploadKey(contract.EvmAbiSaveKey)
		if err != nil {
			loggers.WebLogger.Errorf("ResolveUploadKey error: %v", err.Error())
			return err
		}

		upload, err := db.GetUploadByIdAndUserIdAndHash(id, userId, hash)
		if err != nil {
			loggers.WebLogger.Errorf("GetUploadByIdAndUserIdAndHash error: %v", err.Error())
			return err
		}

		if upload != nil && upload.Content != nil && len(upload.Content) > 0 {
			err = handler.unpackEvmContractResult(upload.Content, tx)
			if err != nil {
				loggers.WebLogger.Errorf("unpackEvmContractResult error: %v", err.Error())
			}
		} else {
			loggers.WebLogger.Warn("evm upload content is nil")
		}
	}

	return nil
}

func (handler *GetTxDetailHandler) unpackEvmContractResult(content []byte, tx *dbcommon.Transaction) error {
	contractAbi, err := abi.JSON(strings.NewReader(string(content)))
	if err != nil {
		return fmt.Errorf("abi json error: %w", err)
	}

	record, err := contractRecord.GetInvokeRecordByTxId(tx.TxId)
	if err != nil {
		return fmt.Errorf("GetInvokeRecordByTxId error: %w", err)
	}

	if record != nil && record.Method != "" && tx.ContractResult != nil && len(tx.ContractResult) > 0 {
		unpackedData, err := contractAbi.Unpack(record.Method, tx.ContractResult)
		if err != nil {
			return fmt.Errorf("abi unpack error: %w, method: %v", err, record.Method)
		}
		tx.ContractResult = []byte(fmt.Sprintf("%v", unpackedData))
	}

	return nil
}
