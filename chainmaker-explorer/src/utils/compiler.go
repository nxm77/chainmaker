package utils

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/entity"
	"chainmaker_web/src/models/compiler"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"strconv"
)

// HttpGetGoIDEVersions 获取Go编译器版本
// @return compiler.VersionsData 编译器版本
// @return error 错误信息
func HttpGetGoIDEVersions() (compiler.VersionsData, error) {
	result := compiler.VersionsData{}
	compilerUrl := config.GlobalConfig.WebConf.ContractCompileUrl
	if compilerUrl == "" {
		return result, fmt.Errorf("compilerUrl is empty")
	}

	url := compilerUrl + CmbGetGoIDEVersions
	body, err := SendHttpGetRequest(url, nil)
	if err != nil {
		log.Errorf("【http】 HttpGetCompilerVersions err:%v", err)
		return result, err
	}

	var respJson compiler.CompilerVersionsResponse
	err = json.Unmarshal(body, &respJson)
	if err != nil {
		log.Errorf("【http】 HttpGetCompilerVersions err:%v", err)
		return result, err
	}

	if respJson.Code != 200 {
		return result, fmt.Errorf("HttpGetCompilerVersions err:%s", respJson.Msg)
	}

	result = respJson.Data
	return result, nil
}

// HttpGetCompilerVersions 获取EVM编译器版本
// @return compiler.VersionsData 编译器版本
// @return error 错误信息
func HttpGetCompilerVersions() (compiler.VersionsData, error) {
	compilerUrl := config.GlobalConfig.WebConf.ContractCompileUrl
	url := compilerUrl + CmbGetCompilerVersions
	body, err := SendHttpGetRequest(url, nil)
	result := compiler.VersionsData{}
	if err != nil {
		log.Errorf("【http】 HttpGetCompilerVersions err:%v", err)
		return result, err
	}
	var respJson compiler.CompilerVersionsResponse
	err = json.Unmarshal(body, &respJson)
	if err != nil {
		log.Errorf("【http】 HttpGetCompilerVersions err:%v", err)
		return result, err
	}

	if respJson.Code != 200 {
		return result, fmt.Errorf("HttpGetCompilerVersions err:%s", respJson.Msg)
	}

	result = respJson.Data
	return result, nil
}

// HttpGetEvmVersions 获取EVM版本
// @return compiler.VersionsData EVM版本
// @return error 错误信息
func HttpGetEvmVersions() (compiler.VersionsData, error) {
	compilerUrl := config.GlobalConfig.WebConf.ContractCompileUrl
	url := compilerUrl + CmbGetEvmVersions
	body, err := SendHttpGetRequest(url, nil)
	result := compiler.VersionsData{}
	if err != nil {
		log.Errorf("【http】 HttpGetCompilerVersions err:%v", err)
		return result, err
	}
	var respJson compiler.CompilerVersionsResponse
	err = json.Unmarshal(body, &respJson)
	if err != nil {
		log.Errorf("【http】 HttpGetCompilerVersions err:%v", err)
		return result, err
	}

	if respJson.Code != 200 {
		return result, fmt.Errorf("HttpGetCompilerVersions err:%s", respJson.Msg)
	}

	result = respJson.Data
	return result, nil
}

// EVMContractCompile 发送合约编译请求，返回编译ID
// @params VerifyContractParams 合约编译参数
// @return string 编译ID
// @return error 错误信息
func SendEVMContractCompile(params *entity.VerifyContractParams) (*compiler.ContractCompileResponse, error) {
	// 将结构体转换为map
	paramMap := make(map[string]string)
	paramMap["chainId"] = params.ChainId
	paramMap["contractAddr"] = params.ContractAddr
	paramMap["contractVersion"] = params.ContractVersion
	paramMap["compilerPath"] = params.CompilerPath
	paramMap["compilerVersion"] = params.CompilerVersion
	paramMap["openLicenseType"] = params.OpenLicenseType
	paramMap["optimization"] = strconv.FormatBool(params.Optimization)
	paramMap["runs"] = strconv.Itoa(params.Runs)
	paramMap["evmVersion"] = params.EvmVersion

	// 将文件添加到一个切片中
	files := map[string]*multipart.FileHeader{
		"contractSourceFile": params.ContractSourceFile,
	}

	// 发送请求
	compilerUrl := config.GlobalConfig.WebConf.ContractCompileUrl
	url := compilerUrl + CmbContractCompile
	body, err := SendHttpPostRequestNew(url, paramMap, files)
	//log.Infof("[http] SendPostContractCompile body:%v", string(body))
	if err != nil {
		log.Errorf("【http】EVMContractCompile Failed to send  url: %v, err:%v", url, err)
		return nil, err
	}

	var respJson *compiler.ContractCompileResponse
	err = json.Unmarshal(body, &respJson)
	if err != nil {
		log.Errorf("【http】EVMContractCompile Unmarshal err:%v", err)
		return nil, err
	}

	if respJson.Code != 200 {
		return nil, fmt.Errorf("【http】EVMContractCompile err:%s", respJson.Msg)
	}

	return respJson, nil
}

// HttpGetContractCompileResult 获取合约编译结果
// @params compileID 编译ID
// @return *compiler.CompilerRersultResponse 编译结果
// @return error 错误信息
func HttpGetContractCompileResult(compileID string) (*compiler.CompilerRersultResponse, error) {
	compilerUrl := config.GlobalConfig.WebConf.ContractCompileUrl
	url := fmt.Sprintf("%s%s&compileID=%s", compilerUrl, CmbGetGetContractCompileResult, compileID)
	body, err := SendHttpGetRequest(url, nil)
	//log.Infof("[http] SendPostContractCompile body:%v", string(body))

	if err != nil {
		log.Errorf("【http】 HttpGetContractCompileResult err:%v", err)
		return nil, err
	}
	var respJson *compiler.CompilerRersultResponse
	err = json.Unmarshal(body, &respJson)
	if err != nil {
		log.Errorf("【http】 HttpGetContractCompileResult err:%v", err)
		return nil, err
	}

	if respJson.Code != 200 {
		return nil, fmt.Errorf("HttpGetContractCompileResult err:%s", respJson.Msg)
	}

	return respJson, nil
}

// SendGoContractCompile 发送合约编译请求，返回编译ID
// @params VerifyContractParams 合约编译参数
// @return string 编译ID
// @return error 错误信息
func SendGoContractCompile(params *entity.VerifyContractParams) (*compiler.ContractCompileResponse, error) {
	// 将结构体转换为map
	paramMap := make(map[string]string)
	paramMap["chainId"] = params.ChainId
	paramMap["contractAddr"] = params.ContractAddr
	paramMap["contractVersion"] = params.ContractVersion
	paramMap["compilerVersion"] = params.CompilerVersion

	files := map[string]*multipart.FileHeader{}
	if params.ContractSourceFile != nil {
		// 将文件添加到一个切片中
		files = map[string]*multipart.FileHeader{
			"contractSourceFile": params.ContractSourceFile,
		}
	}

	// 发送请求
	compilerUrl := config.GlobalConfig.WebConf.ContractCompileUrl
	url := compilerUrl + CmbGoContractCompile
	body, err := SendHttpPostRequestNew(url, paramMap, files)
	if err != nil {
		log.Errorf("【http】SendGoContractCompile Failed to send  url: %v, err:%v", url, err)
		return nil, err
	}

	var respJson *compiler.ContractCompileResponse
	err = json.Unmarshal(body, &respJson)
	if err != nil {
		log.Errorf("【http】SendGoContractCompile Unmarshal err:%v", err)
		return nil, err
	}

	if respJson.Code != 200 {
		return nil, fmt.Errorf("【http】SendGoContractCompile err:%s", respJson.Msg)
	}

	return respJson, nil
}
