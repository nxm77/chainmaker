package entity_cross

// MainCrossConfig 主子链网配置
type MainCrossConfig struct {
	ShowTag bool
}

// CrossSearchView 主子链网搜索
type CrossSearchView struct {
	Type int
	Data string
}

// OverviewDataView OverviewData
type OverviewDataView struct {
	TotalBlockHeight int64
	ShortestTime     int64
	LongestTime      int64
	AverageTime      int64
	SubChainNum      int64
	TxNum            int64
}

// LatestTxListView latest
type LatestTxListView struct {
	CrossId         string
	FromChainId     string
	FromIsMainChain bool
	ToChainId       string
	ToIsMainChain   bool
	Timestamp       int64
	Status          int32
	CrossModel      int32
	TxNum           int64
}

// LatestSubChainListView latest
type LatestSubChainListView struct {
	SubChainId       string
	SubChainName     string
	IsMainChain      bool
	BlockHeight      int64
	CrossTxNum       int64
	CrossContractNum int64
	Timestamp        int64
}

// GetTxListView latest
type GetTxListView struct {
	CrossId         string
	FromChainId     string
	FromIsMainChain bool
	ToChainId       string
	ToIsMainChain   bool
	Timestamp       int64
	Status          int32
	CrossModel      int32
	TxNum           int64
}

// GetSubChainListView get
type GetSubChainListView struct {
	SubChainId       string
	BlockHeight      int64
	IsMainChain      bool
	CrossTxNum       int64
	CrossContractNum int64
	Timestamp        int64
}

// GetCrossSubChainDetailView get
type GetCrossSubChainDetailView struct {
	SubChainId       string
	BlockHeight      int64
	CrossTxNum       int64
	CrossContractNum int64
	Timestamp        int64
	GatewayId        string
	IsMainChain      bool
}

// GetCrossTxDetailView get
type GetCrossTxDetailView struct {
	CrossId        string
	Status         int32
	CrossDuration  int64
	ContractName   string
	ContractMethod string
	Parameter      string
	ContractResult string
	CrossDirection *CrossDirection
	FromChainInfo  *TxChainInfo
	ToChainInfo    *TxChainInfo
	Timestamp      int64
}

type CrossDirection struct {
	FromChain string
	ToChain   string
}

type TxChainInfo struct {
	ChainId      string
	ContractName string
	IsMainChain  bool
	TxId         string
	TxStatus     int32
	Gas          string
}

// GetSubChainCrossView get
type GetSubChainCrossView struct {
	ChainId   string
	ChainName string
	TxNum     int64
}

type GetCrossSyncTxDetailView struct {
	CrossId        string
	Timestamp      int64
	TxId           string
	TxNum          int64
	CrossDirection *CrossDirection
	FromChainInfo  *FromChainInfo
	ToChainInfo    *ToChainInfo
}

type FromChainInfo struct {
	ChainId     string
	IsMainChain bool
}
type ToChainInfo struct {
	ChainId        string
	IsMainChain    bool
	ContractName   string
	ContractMethod string
}
