/*
Package config comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/hokaccha/go-prettyjson"

	"github.com/spf13/viper"
)

const (
	DefaultMaxByteSize      = 256
	DefaultMaxPoolCount     = 50
	DefaultPositionRankTime = 60
)

var (
	// GlobalConfig global
	GlobalConfig *Config
	//gConfPath gConfPath
	gConfPath string
	//MaxDBByteSize 数据库批量插入最大字节,单位byte
	MaxDBByteSize int
	MaxDBPoolSize int
)

var MainDIDContractMap = make(map[string]string)

// Chain 链基本数据
type Chain struct {
	ChainId  string
	AuthType string
}

// SubscribeChains 订阅的所有链数据
var SubscribeChains = make([]*ChainInfo, 0)

// MonitorChains 监控链
var MonitorChains = make([]*ChainInfo, 0)

// printLog 输出日志
func (c *Config) printLog(env string) {
	if env == "" {
		return
	}

	json, err := prettyjson.Marshal(c)
	if err != nil {
		log.Fatalf("marshal alarm config failed, %s", err.Error())
	}
	fmt.Println(string(json))
}

// InitConfig init
// 初始化配置
func InitConfig(confPath, env string) *Config {
	// 设置全局配置路径
	gConfPath = confPath
	// 如果配置路径为空，则获取默认配置路径
	if gConfPath == "" {
		gConfPath = GetConfigDirPath()
	}
	// 初始化viper配置
	webViper, err := initCMViper(gConfPath, env)
	if err != nil {
		// 如果加载配置失败，则输出错误信息并退出
		fmt.Println("can not load config.yml, exit")
		panic(err)
	}
	// 创建配置结构体
	browserConfig := &Config{}
	// 将viper配置解析到结构体中
	if err = webViper.Unmarshal(&browserConfig); err != nil {
		// 如果解析失败，则输出错误信息并退出
		log.Fatal("Unmarshal config failed, ", err)
	}
	// 打印配置日志
	browserConfig.printLog(env)

	//数据库批量插入最大字节，kb转换成byte1
	MaxDBByteSize = DefaultMaxByteSize * 1024
	MaxDBPoolSize = DefaultMaxPoolCount
	// 如果配置中有数据库配置，则设置最大字节和连接池大小
	if browserConfig.DBConf != nil {
		if browserConfig.DBConf.MaxByteSize > DefaultMaxByteSize {
			MaxDBByteSize = browserConfig.DBConf.MaxByteSize * 1024
		}
		if browserConfig.DBConf.MaxPoolSize < DefaultMaxPoolCount {
			MaxDBPoolSize = browserConfig.DBConf.MaxPoolSize
		}
	}

	// 如果配置中有Redis配置，且PositionRankTime小于等于0，则设置默认值
	if browserConfig.RedisDB != nil &&
		browserConfig.RedisDB.PositionRankTime <= 0 {
		browserConfig.RedisDB.PositionRankTime = DefaultPositionRankTime
	}

	//配置全局变量
	GlobalConfig = browserConfig

	//获取配置中所有订阅数据
	SubscribeChains = GetChainsInfoAll()
	//获取监控链
	MonitorChains = getMonitorChainsConfig()
	// 返回配置结构体
	return browserConfig
}

// GetConfigDirPath 绝对路径
func GetConfigDirPath() string {
	_, currentFilePath, _, _ := runtime.Caller(0)
	configDir := filepath.Join(filepath.Dir(currentFilePath), "../../", "configs")
	return configDir
}

// initCMViper
func initCMViper(gConfPath, env string) (*viper.Viper, error) {
	cmViper := viper.New()
	// 使用 env 参数构建配置文件名
	configFilePath := GetConfigFilePath(gConfPath, env)
	cmViper.SetConfigFile(configFilePath)
	if err := cmViper.ReadInConfig(); err != nil {
		return nil, err
	}
	return cmViper, nil
}

// GetConfigFilePath 配置文件路径
func GetConfigFilePath(gConfPath, env string) string {
	configFile := "config.yml"
	if env != "" {
		configFile = fmt.Sprintf("config.%s.yml", env)
	}
	return gConfPath + "/" + configFile
}

// GetChainsInfoAll 获取订阅链信息
func GetChainsInfoAll() []*ChainInfo {
	chainList := make([]*ChainInfo, 0)
	//如果没有订阅信息使用配置信息
	if len(GlobalConfig.ChainsConfig) == 0 {
		return chainList
	}

	for _, chain := range GlobalConfig.ChainsConfig {
		chainInfo := buildChainInfo(chain)
		chainList = append(chainList, chainInfo)
	}

	return chainList
}

// getMonitorChainsConfig 获取监控链配置信息
func getMonitorChainsConfig() []*ChainInfo {
	chainList := make([]*ChainInfo, 0)
	if GlobalConfig == nil || GlobalConfig.MonitorConf == nil ||
		GlobalConfig.MonitorConf.ChainsConfig == nil {
		return chainList
	}

	if !GlobalConfig.MonitorConf.Enable {
		return chainList
	}

	//处理配置监控链数据
	for _, chain := range GlobalConfig.MonitorConf.ChainsConfig {
		chainInfo := buildChainInfo(chain)
		chainList = append(chainList, chainInfo)
	}
	return chainList
}

// buildChainInfo buildChainInfo
func buildChainInfo(chain *ChainConfig) *ChainInfo {
	if chain == nil {
		return &ChainInfo{}
	}

	chainInfo := &ChainInfo{
		ChainId:  chain.ChainId,
		AuthType: chain.AuthType,
		OrgId:    chain.OrgId,
		HashType: chain.HashType,
		TlsMode:  chain.TlsModel,
		Tls:      chain.Tls,
	}
	nodesList := make([]*NodeInfo, 0)

	for _, node := range chain.NodesConfig {
		orgCA := ReadAbsPathFile(node.CaPaths)
		nodesList = append(nodesList, &NodeInfo{
			Addr:        node.Remotes,
			OrgCA:       orgCA,
			TLSHostName: node.TlsHost,
		})
	}
	chainInfo.NodesList = nodesList

	//根据路径读取密钥文件
	userKeyFile := ReadAbsPathFile(chain.UserConf.UserSignKeyFilePath)
	userCrtFile := ReadAbsPathFile(chain.UserConf.UserSignCrtFilePath)
	tlsUserKeyFile := ReadAbsPathFile(chain.UserConf.UserTlsKeyFilePath)
	tlsUserCrtFile := ReadAbsPathFile(chain.UserConf.UserTlsCrtFilePath)
	encUserKeyFile := ReadAbsPathFile(chain.UserConf.UserEncKeyFilePath)
	encUserCrtFile := ReadAbsPathFile(chain.UserConf.UserEncCrtFilePath)
	userInfo := &UserInfo{
		UserSignKey: userKeyFile,
		UserSignCrt: userCrtFile,
		UserTlsKey:  tlsUserKeyFile,
		UserTlsCrt:  tlsUserCrtFile,
		UserEncKey:  encUserKeyFile,
		UserEncCrt:  encUserCrtFile,
	}

	chainInfo.UserInfo = userInfo
	return chainInfo
}

func ReadAbsPathFile(path string) string {
	if path == "" {
		return ""
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		// 处理错误
		fmt.Println("Error getting absolute path:", err)
		return ""
	}
	pathBytes, err := os.ReadFile(absPath)
	if err != nil {
		fmt.Println("Error ReadFile path: err:", absPath, err)
		return ""
	}

	return string(pathBytes)
}

// ToClickHouseUrl to
func (dbConfig *DBConf) ToClickHouseUrl(useDataBase bool) string {
	database := "default"
	if useDataBase && dbConfig.Database != "" {
		database = dbConfig.Database
	}
	connStr := fmt.Sprintf("clickhouse://%s:%s@%s:%s/%s",
		dbConfig.Username, dbConfig.Password, dbConfig.Host, dbConfig.Port, database)
	return connStr
}

// ToMysqlUrl to
func (dbConfig *DBConf) ToMysqlUrl(useDataBase bool) string {
	var url string
	if useDataBase {
		url = fmt.Sprintf("tcp(%s:%s)/%s", dbConfig.Host, dbConfig.Port, dbConfig.Database)
	} else {
		url = fmt.Sprintf("tcp(%s:%s)/", dbConfig.Host, dbConfig.Port)
	}
	return dbConfig.Username + ":" + dbConfig.Password + "@" + url + MysqlDefaultConf
}

// ToPgsqlUrl to
func (dbConfig *DBConf) ToPgsqlUrl(useDataBase bool) string {
	var url string
	if useDataBase {
		url = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable client_encoding=UTF8",
			dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.Password, dbConfig.Database)
	} else {
		url = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable client_encoding=UTF8",
			dbConfig.Host, dbConfig.Port, dbConfig.Username, dbConfig.Password, "template1")
	}
	return url
}

// ToUrl to
func (webConfig *WebConf) ToUrl() string {
	return webConfig.Address + ":" + strconv.Itoa(webConfig.Port)
}
