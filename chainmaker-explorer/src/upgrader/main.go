package main

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/config"
	"chainmaker_web/src/db"
	"chainmaker_web/src/upgrader/registry"
	_ "chainmaker_web/src/upgrader/versions/v2.3.8-v2.3.9"
	"flag"
	"fmt"
	"os"

	"log"
)

var (
	GlobalConfig *config.Config
)

// 当前版本号
const Version = "v2.3.9"

func main() {

	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s <chainId> <version> (e.g. %s chain1 2.3.8)", os.Args[0], os.Args[0])
	}
	chainId := os.Args[1]
	targetVersion := os.Args[2]

	// 组装版本范围（低版本在前，高版本在后）
	versionRange := fmt.Sprintf("%s-%s", targetVersion, Version)

	// 初始化基础服务
	InitConfig()
	InitDB()
	InitRedis()

	// 从注册包获取处理器
	handler, exists := registry.Get(versionRange)
	if !exists {
		log.Fatalf("跨版本升级不被支持！当前版本: %s，目标版本: %s\n"+
			"仅支持顺序升级（例如 %s -> %s）",
			Version, targetVersion, "v2.3.8", "v2.3.9")
	}
	handler(chainId) // 执行升级，传入chainId
}

// 初始化配置
func InitConfig() {
	configPath, env := configYml()
	GlobalConfig = config.InitConfig(configPath, env)
}

// 初始化数据库连接
func InitDB() {
	if GlobalConfig == nil {
		log.Fatal("Configuration not initialized")
	}
	db.InitDbConn(GlobalConfig.DBConf)
}

// 初始化Redis
func InitRedis() {
	if GlobalConfig == nil {
		log.Fatal("Configuration not initialized")
	}
	cache.InitRedis(GlobalConfig.RedisDB)
}

// 获取配置路径
func configYml() (string, string) {
	configPath := flag.String("config", "configs", "config.yml's path")
	env := flag.String("env", "", "yml file name")
	flag.Parse()
	return *configPath, *env
}
