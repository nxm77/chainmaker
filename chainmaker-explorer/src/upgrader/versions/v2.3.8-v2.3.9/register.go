package v238_239

import (
	"chainmaker_web/src/upgrader/registry"
	"log"
)

func init() {
	log.Printf("Registering upgrade handler for version v2.3.8-v2.3.9")
	registry.Register("v2.3.8-v2.3.9", Upgrade) // 直接注册 Upgrade 函数
}
