package compiler

type CompilerVersionsResponse struct {
	Code int          `json:"code"`
	Msg  string       `json:"msg"`
	Data VersionsData `json:"data"`
}

type VersionsData struct {
	Versions []string `json:"version"`
}

type ContractCompileResponse struct {
	Code int          `json:"code"`
	Msg  string       `json:"msg"`
	Data *CompileData `json:"data"`
}
type CompileData struct {
	CompileID string `json:"compileID"`
}

type CompilerRersultResponse struct {
	Code int             `json:"code"`
	Msg  string          `json:"msg"`
	Data *CompilerResult `json:"data"`
}

// CompilerResult 定义编译结果结构体
type CompilerResult struct {
	Bytecode string `json:"bytecode"`
	ABI      string `json:"abi"`
	MetaData string `json:"metadata"`
	Message  string `json:"message"`
	Status   int    `json:"status"`
}
