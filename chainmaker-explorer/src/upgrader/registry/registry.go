// registry/registry.go
package registry

import "sync"

var (
	handlers = make(map[string]func(string)) // 修改为接收string参数
	mu       sync.RWMutex
)

// 注册处理器（跨包可见）
func Register(version string, handler func(string)) {
	mu.Lock()
	handlers[version] = handler
	mu.Unlock()
}

// 获取处理器（跨包可见）
func Get(version string) (func(string), bool) {
	mu.RLock()
	defer mu.RUnlock()
	handler, exists := handlers[version]
	return handler, exists
}
