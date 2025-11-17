package utils

const (
	// MonitorNameSpace MonitorNameSpace
	MonitorNameSpace = "explorer_backend"

	// CrossSubChainInfoUrl 根据子链id获取子链信息
	CrossSubChainInfoUrl = "/mainChildChain/getSubChainInfo"

	// CrossGateWayIdUrl 根据子链网关，获取子链信息
	CrossGateWayIdUrl = "/mainChildChain/getGatewayInfo"
)

const (
	// CmbContractCompiler 合约编译
	CmbContractCompile = "/contract/compiler?cmb=ContractCompile"
	// CmbGetCompilerVersions 获取编译器版本列表
	CmbGetCompilerVersions = "/contract/compiler?cmb=GetCompilerVersions"
	// CmbGetEvmVersions 获取EVM版本列表
	CmbGetEvmVersions = "/contract/compiler?cmb=GetEvmVersions"
	// CmbGetGetContractCompileResult 获取合约编译结果
	CmbGetGetContractCompileResult = "/contract/compiler?cmb=GetContractCompileResult"
	// CmbGetCompilerVersions 获取go编译器版本列表
	CmbGetGoIDEVersions = "/contract/compiler?cmb=GetGoIDEVersions"
	// CmbGoContractCompile 合约编译
	CmbGoContractCompile = "/contract/compiler?cmb=GoContractCompile"
)

const (
	// CmbGetQueryTaskResult 获取查询任务结果
	CmbGetQueryTaskResult = "/dquery/sqlQuery?cmb=GetQueryTaskResult"
)
