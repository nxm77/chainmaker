package chain_management

import (
	"management_backend/src/config"
	"testing"
)

// TestGetResourcesPolicies 测试获取默认策略
func TestGetResourcesPolicies(t *testing.T) {
	config.ConfEnvPath = "../../../dependence/"
	getDefaultResourcePolicies()
}
