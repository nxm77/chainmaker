package service

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/db/dbhandle"
	"chainmaker_web/src/entity"
	"chainmaker_web/src/models/compiler"
	"chainmaker_web/src/sync/common"
	"chainmaker_web/src/utils"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	//VerifyStatusSuccess 合约编译中
	VerifyStatusDefault = 0
	//VerifyStatusSuccess 验证成功
	VerifyStatusSuccess = 1
	//VerifyStatusFailed 验证失败
	VerifyStatusFailed = 2
)

// GetGoIDEVersionsHandler handler
type GetGoIDEVersionsHandler struct{}

// Handle deal
func (handler *GetGoIDEVersionsHandler) Handle(ctx *gin.Context) {
	//获取编译器版本列表
	versionResult, err := utils.HttpGetGoIDEVersions()
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	//组装response数据
	view := &entity.VersionView{
		Versions: versionResult.Versions,
	}
	ConvergeDataResponse(ctx, view, nil)
}

// VerifyContractSourceCodeHandler handler
type VerifyContractSourceCodeHandler struct{}

// Handle deal
func (handler *VerifyContractSourceCodeHandler) Handle(ctx *gin.Context) {
	params := entity.BindVerifyContractCodeHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, entity.ErrorMsgParam)
		ConvergeFailureResponse(ctx, newError)
		return
	}

	chainId := params.ChainId
	contractAddr := params.ContractAddr
	version := params.ContractVersion
	//数据库获取链上合约源码
	upgradeContract, err := dbhandle.GetUpgradeContractInfo(chainId, contractAddr, version)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, entity.ErrContractNotFound)
		return
	}

	runtimeType := upgradeContract.ContractRuntimeType
	//已经验证过了
	if upgradeContract.VerifyStatus == VerifyStatusSuccess {
		returnVerfiyContractView(ctx, VerifyStatusSuccess, "", "")
		return
	}

	//执行合约编译
	cmpileID, err := ExecuteContractCompiler(params, runtimeType)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	//如果编译信息中有错误信息或者编译数据为空，则返回错误信息
	if cmpileID == "" {
		returnVerfiyContractView(ctx, VerifyStatusFailed, entity.ErrorContractCompilerFailed, "")
		return
	}

	//根据编译ID获取合约编译结果
	compilerResult, err := GetContractCompileResult(cmpileID)
	if err != nil {
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	//获取合约字节码信息
	byteCodeinfo, err := dbhandle.GetContractByteCodeByTx(chainId, upgradeContract.TxId)
	if err != nil || byteCodeinfo == nil {
		ConvergeHandleFailureResponse(ctx, entity.ErrContractNotFound)
		return
	}

	//确认合约验证结果
	verifyStatus, message := validateCompilerResult(compilerResult, runtimeType, byteCodeinfo.ByteCode)

	//保存合约验证结果
	verifyResult, err := SaveContractVerifyResult(params, verifyStatus, compilerResult, upgradeContract)
	if err != nil {
		log.Errorf("SaveContractVerifyResult err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, entity.ErrContractCompilerFailed)
		return
	}

	//更新合约验证结果
	err = dbhandle.UpdateUpgradeContractVerifyStatus(chainId, contractAddr, version, verifyStatus)
	if err != nil {
		log.Errorf("UpdateUpgradeContractVerifyStatus err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, entity.ErrContractCompilerFailed)
		return
	}

	//保存合约源码
	if verifyStatus == VerifyStatusSuccess {
		// 解压合约源码文件
		sourceFileList, err := utils.UnzipFileContractSource(params.ContractSourceFile)
		if err != nil {
			log.Errorf("UnzipFileContractSource err : %s", err.Error())
			ConvergeHandleFailureResponse(ctx, entity.ErrContractCompilerFailed)
			return
		}

		// 创建 ContractSourceFile 实例列表
		insertList := make([]*db.ContractSourceFile, 0)
		for _, file := range sourceFileList {
			newUUID := uuid.New().String()
			// 创建 ContractSourceFile 实例
			sourceFile := &db.ContractSourceFile{
				ID:              newUUID,
				VerifyId:        verifyResult.VerifyId,
				ContractAddr:    verifyResult.ContractAddr,
				ContractVersion: verifyResult.ContractVersion,
				SourcePath:      file.SourcePath,
				SourceCode:      file.SourceCode,
			}
			// 将实例添加到列表中
			insertList = append(insertList, sourceFile)
		}

		// 将 ContractSourceFile 实例保存到数据库
		err = dbhandle.InsertContractSource(chainId, insertList)
		if err != nil {
			log.Errorf("InsertContractSource err : %s", err.Error())
			// 处理失败响应
			ConvergeHandleFailureResponse(ctx, entity.ErrContractCompilerFailed)
			return
		}
	}

	// 根据编译结果返回不同的信息
	returnVerfiyContractView(ctx, verifyStatus, message, compilerResult.MetaData)
}

// returnVerfiyContractView 返回合约验证视图
// @param ctx 上下文
// @param verifyStatus 验证状态
// @param message 消息
// @param metaData 元数据
func returnVerfiyContractView(ctx *gin.Context, verifyStatus int, message, metaData string) {
	view := &entity.VerifyContractView{
		VerifyStatus: verifyStatus,
		Message:      message,
		MetaData:     metaData,
	}
	ConvergeDataResponse(ctx, view, nil)
}

// ExecuteContractCompiler 执行合约编译
// @param params 合约编译参数
// @param runtimeType 运行时类型
// @return compileID 编译ID
// @return err 错误信息
// ExecuteContractCompiler 函数用于执行合约编译器
func ExecuteContractCompiler(params *entity.VerifyContractParams, runtimeType string) (string, error) {
	var err error
	// 定义编译器信息变量
	var compilerInfo *compiler.ContractCompileResponse
	// 根据运行时类型选择编译器
	switch runtimeType {
	case common.RuntimeTypeDockerGo:
		// 发送EVM合约编译请求
		compilerInfo, err = utils.SendGoContractCompile(params)
	case common.RuntimeTypeEVM:
		// 发送Go合约编译请求
		compilerInfo, err = utils.SendEVMContractCompile(params)
	default:
		// 运行时类型错误，返回错误信息
		return "", fmt.Errorf(entity.ErrorCompilerType)
	}

	// 如果编译器信息为空，返回错误信息
	if err != nil {
		log.Errorf("ExecuteContractCompiler err : %v", err)
		return "", fmt.Errorf(entity.ErrorContractCompilerFailed)
	}

	if compilerInfo == nil || compilerInfo.Data == nil {
		return "", fmt.Errorf(entity.ErrorContractCompilerFailed)
	}

	// 返回编译ID
	return compilerInfo.Data.CompileID, nil
}

// validateCompilerResult 验证编译结果
// @param compilerResult 编译结果
// @param runtimeType 运行时类型
// @param byteCodeDB 编译结果
// @return verifyStatus 验证状态
// @return message 消息
func validateCompilerResult(compilerResult *compiler.CompilerResult, runtimeType string,
	byteCodeDB []byte) (int, string) {
	// 定义消息变量
	var (
		err                error
		message            string
		contractByteCodeDB string
	)
	// 定义验证状态变量
	verifyStatus := VerifyStatusFailed
	if runtimeType == common.RuntimeTypeDockerGo {
		//将DB的编译文件解压缩，转成16进制sha256字符串
		contractByteCodeDB, err = utils.ExtractAndHash7zFile(byteCodeDB)
		if err != nil {
			message = entity.ExtractAndHash7zFailed
			log.Errorf("ExtractAndHash7zFile err : %s", err.Error())
			return verifyStatus, message
		}
	} else if runtimeType == common.RuntimeTypeEVM {
		contractByteCodeDB = hex.EncodeToString(byteCodeDB)
	} else {
		message = entity.ErrorCompilerType
		return verifyStatus, message
	}

	// 验证成功
	if compilerResult.Bytecode == contractByteCodeDB {
		verifyStatus = VerifyStatusSuccess
	} else if compilerResult.Bytecode == "" {
		message = entity.ErrorContractCompilerFailed
		// 如果编译结果中有消息
		if compilerResult.Message != "" {
			// 将消息赋值给message变量
			message = compilerResult.Message
		}
	} else {
		// 如果编译结果与数据库中的合约字节码不同
		log.Infof("contractByteCodeDB:%s, verfity:%s,len-byteCodeDB:%d", contractByteCodeDB,
			compilerResult.Bytecode, len(byteCodeDB))
		message = entity.ErrorContractCompilerDiff
	}

	// 返回验证状态和消息
	return verifyStatus, message
}

// SaveContractVerifyResult 保存合约验证结果
// @param params 验证参数
// @param verifyStatus 验证状态
// @param compilerResult 编译结果
// @param contractName 合约名称
// @return verifyResult 验证结果
// @return err 错误信息
func SaveContractVerifyResult(params *entity.VerifyContractParams, verifyStatus int,
	compilerResult *compiler.CompilerResult, upgradeContract *db.UpgradeContractTransaction) (
	*db.ContractVerifyResult, error) {
	if compilerResult == nil {
		return nil, nil
	}

	// 查询合约验证结果
	verifyResult, err := dbhandle.GetContractVerifyResult(params.ChainId, params.ContractAddr, params.ContractVersion)
	if err != nil {
		return nil, err
	}

	// 将 ABI 字段转换为 JSON 字符串
	var abiStr string
	if compilerResult.ABI != "" {
		abiJSON, _ := json.Marshal(compilerResult.ABI)
		abiStr = string(abiJSON)
	}

	//新增
	if verifyResult == nil {
		newUUID := uuid.New().String()
		// 创建 ContractSourceFile 实例
		verifyInfo := &db.ContractVerifyResult{
			VerifyId:            newUUID,
			VerifyStatus:        verifyStatus,
			ContractAddr:        upgradeContract.ContractAddr,
			ContractName:        upgradeContract.ContractNameBak,
			ContractVersion:     params.ContractVersion,
			ContractRuntimeType: upgradeContract.ContractRuntimeType,
			CompilerVersion:     params.CompilerVersion,
			CompilerPath:        params.CompilerPath,
			ByteCode:            []byte(compilerResult.Bytecode),
			ABI:                 abiStr,
			MetaData:            compilerResult.MetaData,
			OpenLicenseType:     params.OpenLicenseType,
			EvmVersion:          params.EvmVersion,
			Optimization:        params.Optimization,
			RunNum:              params.Runs,
		}

		// 将 ContractVerify 实例保存到数据库
		err = dbhandle.InsertContractVerifyResult(params.ChainId, verifyInfo)
		return verifyInfo, err
	}

	//已经是编译成功的数据不在重新编译
	if verifyResult.VerifyStatus == VerifyStatusSuccess {
		return verifyResult, nil

	}

	//更新
	verifyResult.VerifyStatus = verifyStatus
	verifyResult.CompilerPath = params.CompilerPath
	verifyResult.ByteCode = []byte(compilerResult.Bytecode)
	verifyResult.ABI = abiStr
	verifyResult.MetaData = compilerResult.MetaData
	verifyResult.CompilerVersion = params.CompilerVersion
	verifyResult.OpenLicenseType = params.OpenLicenseType
	verifyResult.EvmVersion = params.EvmVersion
	verifyResult.Optimization = params.Optimization
	verifyResult.RunNum = params.Runs
	// 将 ContractVerify 实例保存到数据库
	err = dbhandle.UpdateContractVerifyResult(params.ChainId, verifyResult)
	return verifyResult, err
}

// GetContractCompileResult 获取合约编译结果
// @params params 	合约验证参数
// @return *compiler.CompilerResult 	编译结果
// @return error 					错误信息
func GetContractCompileResult(compileID string) (*compiler.CompilerResult, error) {
	//每秒轮询获取编译结果，直到状态不等于0结束循环
	//编译时间超过timeout后自动结束查询
	var compilerResult *compiler.CompilerResult
	result := &compiler.CompilerResult{}
	timeout := config.ContractCompilationTimeout

	// 创建一个定时器，用于超时控制
	timer := time.NewTimer(timeout)
	defer timer.Stop() // 确保在函数结束时停止定时器

	for {
		// 检查是否超时
		select {
		case <-timer.C:
			result.Message = entity.ErrorContractCompilerTimeout
			return result, nil
		default:
		}

		// 获取编译结果
		compilerResponse, err := utils.HttpGetContractCompileResult(compileID)
		if err != nil {
			log.Errorf("GetContractCompileResult err : %s", err.Error())
			newError := entity.NewError(entity.ErrorSystem, entity.GetCompilerResultFailed)
			return nil, newError
		}

		// 如果编译结果为 nil，继续轮询
		if compilerResponse == nil || compilerResponse.Data == nil {
			time.Sleep(1 * time.Second)
			continue
		}

		if compilerResponse.Code != http.StatusOK {
			result.Message = compilerResponse.Msg
			return result, nil
		}

		// 检查编译状态，是否编译结束
		if compilerResponse.Data.Status != VerifyStatusDefault {
			compilerResult = compilerResponse.Data
			return compilerResult, nil
		}

		// 每次轮询后等待1秒
		time.Sleep(1 * time.Second)
	}
}

// GetContractVersionsHandler handler
type GetContractVersionsHandler struct{}

// Handle deal
func (handler *GetContractVersionsHandler) Handle(ctx *gin.Context) {
	params := entity.BindContractVersionsHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//用户列表
	contractVersions, err := dbhandle.GetContractVersions(params.ChainId, params.ContractAddr)
	if err != nil {
		log.Errorf("GetContractVersions err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	//组装response数据
	view := &entity.VersionView{
		Versions: contractVersions,
	}
	ConvergeDataResponse(ctx, view, nil)
}

// GetCompilerVersionsHandler handler
type GetCompilerVersionsHandler struct{}

// Handle deal
func (handler *GetCompilerVersionsHandler) Handle(ctx *gin.Context) {
	//获取编译器版本列表
	versionResult, err := utils.HttpGetCompilerVersions()
	if err != nil {
		log.Errorf("HttpGetCompilerVersions err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	//组装response数据
	view := &entity.VersionView{
		Versions: versionResult.Versions,
	}
	ConvergeDataResponse(ctx, view, nil)
}

// GetEvmVersionsHandler handler
type GetEvmVersionsHandler struct{}

// Handle deal
func (handler *GetEvmVersionsHandler) Handle(ctx *gin.Context) {
	//获取编译器版本列表
	versionResult, err := utils.HttpGetEvmVersions()
	if err != nil {
		log.Errorf("HttpGetEvmVersions err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	//组装response数据
	view := &entity.VersionView{
		Versions: versionResult.Versions,
	}
	ConvergeDataResponse(ctx, view, nil)
}

// GetOpenLicenseTypesHandler handler
type GetOpenLicenseTypesHandler struct{}

// Handle deal
func (handler *GetOpenLicenseTypesHandler) Handle(ctx *gin.Context) {
	//组装response数据
	view := &entity.OpenLicenseTypeView{
		LicenseTypes: config.OpenLicenseTypes,
	}
	ConvergeDataResponse(ctx, view, nil)
}

// GetContractCodeHandler handler
type GetContractCodeHandler struct {
}

// Handle deal
func (getContractCodeHandler *GetContractCodeHandler) Handle(ctx *gin.Context) {
	params := entity.BindGetContractCodeHandler(ctx)
	if params == nil || !params.IsLegal() {
		newError := entity.NewError(entity.ErrorParamWrong, "param is wrong")
		ConvergeFailureResponse(ctx, newError)
		return
	}

	//获取合约验证结果
	verifyResult, err := dbhandle.GetContractVerifyResult(params.ChainId, params.ContractAddr, params.ContractVersion)
	if err != nil {
		log.Errorf("GetContractVerifyResult err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	if verifyResult == nil {
		ConvergeHandleFailureResponse(ctx, fmt.Errorf("合约还未验证"))
		return
	}

	//获取合约源码
	sourceFiles, err := dbhandle.GetContractSourceFile(params.ChainId, verifyResult.VerifyId)
	if err != nil {
		log.Errorf("GetContractSourceFile err : %s", err.Error())
		ConvergeHandleFailureResponse(ctx, err)
		return
	}

	//组装源码文件
	sourceCodes := make([]*entity.SourceCode, 0, len(sourceFiles))
	for _, file := range sourceFiles {
		sourceCodes = append(sourceCodes, &entity.SourceCode{
			SourcePath: file.SourcePath,
			SourceCode: string(file.SourceCode),
		})
	}

	contractCodeView := &entity.ContractCodeView{
		VerifyStatus:    verifyResult.VerifyStatus,
		ContractAbi:     verifyResult.ABI,
		SourceCodes:     sourceCodes,
		MetaData:        verifyResult.MetaData,
		ContractName:    verifyResult.ContractName,
		ContractAddr:    verifyResult.ContractAddr,
		ContractVersion: verifyResult.ContractVersion,
		RuntimeType:     verifyResult.ContractRuntimeType,
		CompilerVersion: verifyResult.CompilerVersion,
		EvmVersion:      verifyResult.EvmVersion,
		Optimization:    verifyResult.Optimization,
		Runs:            verifyResult.RunNum,
		OpenLicenseType: verifyResult.OpenLicenseType,
	}

	ConvergeDataResponse(ctx, contractCodeView, nil)
}
