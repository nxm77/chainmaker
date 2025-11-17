// nolint
package db

import (
	"chainmaker_web/src/cache"
	"chainmaker_web/src/config"
	"fmt"
	"os"
	"regexp"
)

const UTchainID = "chainmaker_pk"

var nodesList = []*config.NodeInfo{
	{
		Addr: "pre-chain1.cnbn.org.cn:12391",
		//Addr: "192.168.3.170:12391",
	},
}
var ChainListConfigTest = &config.ChainInfo{
	ChainId:   "chainmaker_pk",
	AuthType:  "public",
	OrgId:     "",
	HashType:  "SM3",
	TlsMode:   0,
	Tls:       false,
	NodesList: nodesList,
	UserInfo: &config.UserInfo{
		UserSignKey: "-----BEGIN EC PRIVATE KEY-----\nMHcCAQEEIIHaiq+GBx42Lw7YdBgfOvYIWtmqCtk/XRdAzbDQvICgoAoGCCqGSM49\nAwEHoUQDQgAEBtUSf7SDTxemXSHKgIrblrzQM2xx3mqoAA4vDTYm3txZ5lfnAB7D\nBGyAX5Qbap9QLcCrcCN56WGO5iGYN7Splg==\n-----END EC PRIVATE KEY-----\n",
	},
}

func InitRedisContainer() {
	hostUrl := os.Getenv("UT_REDIS_URL")
	if hostUrl != "" {
		redisCfg := &config.RedisConfig{
			Type:     "node",
			Host:     []string{hostUrl},
			Password: "",
			Username: "",
		}

		config.GlobalConfig.RedisDB = redisCfg
	}

	log.Infof("=========redisCfg:%v=", config.GlobalConfig.RedisDB)
	cache.InitRedis(config.GlobalConfig.RedisDB)
}

func InitMySQLContainer() {
	config.SubscribeChains = []*config.ChainInfo{
		ChainListConfigTest,
	}

	hostUrl := os.Getenv("UT_MYSQL_DB_URL")
	//hostUrl = "root:123456@tcp(127.0.0.1:33061)/chainmaker_explorer_dev"
	if hostUrl != "" {
		mysqlCfg, _ := parseDBURL(hostUrl)
		config.GlobalConfig.DBConf = mysqlCfg
	}

	log.Infof("=========mysqlCfg:%v====hostUrl:%v", config.GlobalConfig.DBConf, hostUrl)

	InitDbConn(config.GlobalConfig.DBConf)

	InitChainListMySqlTable(config.SubscribeChains)
}

func parseDBURL(dbURL string) (*config.DBConf, error) {
	// 使用正则表达式解析数据库连接字符串
	re := regexp.MustCompile(`^(?P<username>[^:]+):(?P<password>[^@]+)@` +
		`tcp\((?P<host>[^:]+):(?P<port>[0-9]+)\)(?:/(?P<database>.*))?$`)
	match := re.FindStringSubmatch(dbURL)
	if match == nil {
		return nil, fmt.Errorf("invalid database URL format")
	}

	// 创建一个 map 来存储解析后的值
	paramsMap := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			paramsMap[name] = match[i]
		}
	}

	// 填充配置结构体
	dbConf := &config.DBConf{
		Host:       paramsMap["host"],
		Port:       paramsMap["port"],
		Username:   paramsMap["username"],
		Password:   paramsMap["password"],
		Database:   paramsMap["database"],
		Prefix:     "test_",
		DbProvider: "Mysql",
	}
	return dbConf, nil
}
