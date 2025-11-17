// Package code 定义错误码
package entity

import (
	"errors"
)

const (
	// ErrorAuthFailure auth
	ErrorAuthFailure = "AuthFailure"
	// ErrorHandleFailure handel
	ErrorHandleFailure = "HandleFailure"
	// ErrorParamWrong param
	ErrorParamWrong = "ParamError"
	// ErrorSdkClient sdk
	ErrorSdkClient = "SdkClientError"
	// ErrorSubscribe sub
	ErrorSubscribe = "SubscribeError"
	// ErrorSystem sub
	ErrorSystem = "SystemError"
)
const (
	//ErrorMsgParam 参数错误
	ErrorMsgParam = "param is wrong"
	//ErrorMsgDataSelect 数据采集失败
	ErrorMsgDataSelect = "data acquisition failed"
	//ErrorMsgDataUpdate 数据更新失败
	ErrorMsgDataUpdate = "data update failed"
	//ErrorMsgParam 参数错误
	ErrorSystemError = "config is wrong"
	//ErrorChainShowError chain is not show
	ErrorChainShowError = "chain is not show"
)

const (
	// ErrorContractCompilerNotFound
	ErrorContractCompilerNotFound = "合约编译器未找到"
	//ErrorContractNotFound	链上合约信息未找到
	ErrorContractNotFound = "链上合约信息未找到"
	// ErrorContractCompilerTimeout
	ErrorContractCompilerTimeout = "源码编译超时"
	// ErrorContractCompilerFailed
	ErrorContractCompilerFailed = "源码编译失败"
	//合约验证失败
	// ErrorContractCompilerResultFailed 合约验证失败
	ErrorContractCompilerResultFailed = "合约验证失败"
	// ExtractAndHash7zFailed
	ExtractAndHash7zFailed = "源码编译文件解压缩失败"
	// ErrorContractCompilerDiff
	ErrorContractCompilerDiff = "合约编译文件比对失败, 请检查文件是否一致"
	//ErrorCompilerType
	ErrorCompilerType = "合约类型错误，仅支持 EVM 合约和 DockerGo 合约验证"
	//GetCompilerResultFailed 获取编译器编译结果失败
	GetCompilerResultFailed = "获取编译器编译结果失败"
)

var (
	// NewContractNotFound 合约未找到
	ErrContractNotFound = errors.New(ErrorContractNotFound)
	//合约验证失败
	ErrContractCompilerFailed = errors.New(ErrorContractCompilerResultFailed)
)

var (
	// ErrRecordNotFoundErr 未查询到数据
	ErrRecordNotFoundErr = errors.New("record not found")
	// ErrUpdateFail 数据更新失败
	ErrUpdateFail = errors.New("data update failed")
	// ErrSelectFailed 数据采集失败
	ErrSelectFailed = errors.New("data acquisition failed")
)

const (
	// ErrorCodeFailure  Failure code
	ErrorCodeFailure = "1001"
	//ErrorMessageFailure Failure
	ErrorMessageFailure = "Failure"

	//ErrorCodeInvalidParameters Parameter error code
	ErrorCodeInvalidParameters = "2001"
	//ErrorMessageParameters 参数错误
	ErrorMessageParameters = "invalid parameters"

	// ErrorNotLoggedIn Not logged in code
	ErrorNotLoggedIn = "401"
	//ErrorMsgNotLoggedIn 未登录
	ErrorMsgNotLoggedIn = "Not logged in"

	// ErrorCodeNoPermission No permission code
	ErrorCodeNoPermission = "1002"
	//ErrorMessageNoPermission 没有权限
	ErrorMessageNoPermission = "no permission"
	//ErrorMsgTokenExpired token失效
	// nosec G101 - false positive (error message)
	ErrorMsgTokenExpired = "Token已过期，请重新登录" // #nosec G101 - false positive (error message)
	//ErrorMsgTokenFormat token格式错误
	ErrorMsgTokenFormat = "the token does not start with 'Bearer '"
	//ErrorMsgTokenInvalid token错误
	ErrorMsgTokenInvalid = "invalid token"
	//ErrorMsgTokenParseUserID 无法解析id
	// nosec G101 - false positive (error message)
	ErrorMsgTokenParseUserID = "无法解析用户地址" // #nosec G101 - false positive (error message)

	//ErrorPasswordError 密码错误
	ErrorPasswordError = "password error"
	//不是admin账户
	ErrorNotAdmin = "not admin account"
)

const (
	// ErrorContractNotExist 合约不存在
	ErrorContractNotExist = "contract not exist"
	//ErrorGenerateTableNameFailed 生成表名称失败
	ErrorGenerateTableNameFailed = "generate table name failed"
)

const (
	// 合约abi文件中的函数type只能是function，event，constructor
	ErrorContractAbiFunctionType = "abi function type must be function, event or constructor"
	// 合约abi文件中的event事件中的inputs字段名称不能包含系统字段（sysId,sysTxId,sysContractVersion,sysTimestamp
	ErrorContractAbiEventInputs = "abi event inputs can not contain " +
		"system field(sysId, sysTxId, sysContractVersion, sysTimestamp)"
)

// GetErrorMsgParams 参数错误
func GetErrorMsgParams() *Error {
	newError := NewError(
		ErrorCodeInvalidParameters,
		ErrorMessageParameters,
	)
	return newError
}

// GetErrorNotLogged 没有登录
func GetErrorNotLogged() *Error {
	newError := NewError(
		ErrorNotLoggedIn,
		ErrorMsgNotLoggedIn,
	)
	return newError
}

// GetErrorNoPermission 没有权限
func GetErrorNoPermission(msg string) *Error {
	if msg == "" {
		msg = ErrorMessageNoPermission
	}
	newError := NewError(
		ErrorCodeNoPermission,
		msg,
	)
	return newError
}

func GetErrorMsg(message string) *Error {
	newError := NewError(
		ErrorCodeFailure,
		message,
	)
	return newError
}
