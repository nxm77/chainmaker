package utils

const (
	// ABIEventType 事件类型
	ABIEventType = "event"
	// ABIFunctionType 函数类型
	ABIFunctionType = "function"
	//ABIConstructorType 构造函数类型
	ABIConstructorType = "constructor"
)

// ContractABI 表示合约的ABI结构
type ContractABI struct {
	Type            string        `json:"type"`
	Name            string        `json:"name,omitempty"`
	Inputs          []ABIParamIn  `json:"inputs,omitempty"`
	Outputs         []ABIParamOut `json:"outputs,omitempty"`
	StateMutability string        `json:"stateMutability,omitempty"`
	Anonymous       bool          `json:"anonymous,omitempty"`
	Describe        string        `json:"describe,omitempty"`
}

// ABIParam 表示ABI中的参数
type ABIParamIn struct {
	InternalType string `json:"internalType"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Indexed      bool   `json:"indexed,omitempty"`
}

// ABIParam 表示ABI中的参数
type ABIParamOut struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
