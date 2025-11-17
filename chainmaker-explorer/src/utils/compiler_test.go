package utils

import (
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/entity"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
)

var ContractCompileUrl = "http://test.com"

func TestHttpGetGoIDEVersions(t *testing.T) {
	// Test case 2: Response code is not 200
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()
	config.GlobalConfig.WebConf.ContractCompileUrl = ts.URL
	go HttpGetGoIDEVersions()
}

func TestHttpGetCompilerVersions(t *testing.T) {
	// Test case 1: Successful request to get EVM compiler versions
	config.GlobalConfig.WebConf.ContractCompileUrl = ContractCompileUrl
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code": 200, "msg": "success", "data": {"versions": ["1.0", "2.0"]}}`))
	}))
	defer server.Close()

	go HttpGetCompilerVersions()
}

func TestHttpGetEvmVersions(t *testing.T) {
	// Test case 1: Successful request to get EVM versions
	config.GlobalConfig.WebConf.ContractCompileUrl = ContractCompileUrl
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code": 200, "msg": "success", "data": {"versions": ["1.0", "2.0"]}}`))
	}))
	defer server.Close()

	go HttpGetEvmVersions()
}

func TestSendEVMContractCompile(t *testing.T) {
	// Test case 1: Successful request to compile EVM contract
	config.GlobalConfig.WebConf.ContractCompileUrl = ContractCompileUrl
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code": 200, "msg": "success", "data": {"compileID": "12345"}}`))
	}))
	defer server.Close()

	params := &entity.VerifyContractParams{
		ChainId:            db.UTchainID,
		ContractAddr:       "contract1",
		ContractVersion:    "1.0",
		CompilerPath:       "/path/to/compiler",
		CompilerVersion:    "1.0",
		OpenLicenseType:    "type1",
		Optimization:       true,
		Runs:               100,
		EvmVersion:         "1.0",
		ContractSourceFile: &multipart.FileHeader{},
	}

	go SendEVMContractCompile(params)
}
func TestHttpGetContractCompileResult(t *testing.T) {
	// Test case 1: Successful request to get contract compile result
	config.GlobalConfig.WebConf.ContractCompileUrl = ContractCompileUrl
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code": 200, "msg": "success", "data": {"result": "success"}}`))
	}))
	defer server.Close()

	compileID := "12345"
	go HttpGetContractCompileResult(compileID)
}

func TestSendGoContractCompile(t *testing.T) {
	// Test case 1: Successful request to compile Go contract
	config.GlobalConfig.WebConf.ContractCompileUrl = ContractCompileUrl
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code": 200, "msg": "success", "data": {"compileID": "12345"}}`))
	}))
	defer server.Close()

	params := &entity.VerifyContractParams{
		ChainId:            db.UTchainID,
		ContractAddr:       "contract1",
		ContractVersion:    "1.0",
		CompilerVersion:    "1.0",
		ContractSourceFile: &multipart.FileHeader{},
	}

	go SendGoContractCompile(params)
}
