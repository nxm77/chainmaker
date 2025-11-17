/*
Package config comment
Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
SPDX-License-Identifier: Apache-2.0
*/
package config

import "time"

// PUBLIC public模式
const PUBLIC = "public"

// TlsModelSingle 国密单证书模式
const TlsModelSingle = 0

// TlsModelDouble 国密双证书模式
const TlsModelDouble = 1

// RoleClient RoleClient
const RoleClient = "client"

const (
	// MySql mysql
	MySql = "Mysql"
	//ClickHouse ClickHouse数据库
	ClickHouse = "ClickHouse"
	//Pgsql 人大金仓数据库
	Pgsql = "Pgsql"
	// MysqlDefaultConf db default config
	MysqlDefaultConf = "?charset=utf8mb4&parseTime=True&loc=Local"
	//MysqlDatabaseConf mysql utf8
	MysqlDatabaseConf = " CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;"
	// DbMaxIdleConns db max idle config
	DbMaxIdleConns = 20
	// DbMaxOpenConns db max open config
	DbMaxOpenConns = 150
)

const (
	//BlockInsertWorkerCount 区块同步数据处理并发数量
	BlockInsertWorkerCount = 5
	//BlockWaitUpdateWorkerCount 区块异步最大等待处理数量
	BlockWaitUpdateWorkerCount = 20
	//BlockUpdateWorkerCount 区块异步数据处理并发数量
	BlockUpdateWorkerCount = 5
)

const (
	//ContractCompilationTimeout 合约编译超时时间
	ContractCompilationTimeout = 180 * time.Second
)

var OpenLicenseTypes = []string{
	"No License (None)",
	"The Unlicense (Unlicense)",
	"MIT License (MIT)",
	"GNU General Public License v2.0 (GNU GPLv2)",
	"GNU General Public License v3.0 (GNU GPLv3)",
	"GNU Lesser General Public License v2.1 (GNU LGPLv2.1)",
	"GNU Lesser General Public License v3.0 (GNU LGPLv3)",
	"BSD 2-clause 'Simplified' license (BSD-2-Clause)",
	"BSD 3-clause 'New' Or 'Revised' license* (BSD-3-Clause)",
	"Mozilla Public License 2.0 (MPL-2.0)",
	"Open Software License 3.0 (OSL-3.0)",
	"Apache 2.0 (Apache-2.0)",
	"GNU Affero General Public License (GNU AGPLv3)",
	"Business Source License (BSL 1.1)",
}

const (
	// SM2 国密
	SM2 = 0
	// ECDSA 非国密
	ECDSA = 1
)

const (
	// StatusNormal normal
	StatusNormal = 0
	// StatusDeleted delete
	StatusDeleted = 1
)

// ContractWarnMsg 合约敏感词异常提醒
const ContractWarnMsg = "合约名称违规"

// OtherWarnMsg 敏感词违规
const OtherWarnMsg = "上链内容违反相关法律规定，内容已屏蔽"

// ContractResultMsg ContractResultMsg
var ContractResultMsg = []byte(OtherWarnMsg)

// MaxRetryCount 失败最大重试次数
const MaxRetryCount = 100

// WebConf Http配置
type WebConf struct {
	Address             string `mapstructure:"address"`
	Port                int    `mapstructure:"port"`
	CrossDomain         bool   `mapstructure:"cross_domain"`
	ThirdApplyUrl       string `mapstructure:"third_apply_url"`
	RelayCrossChainUrl  string `mapstructure:"relay_cross_chain_url"`
	ContractCompileUrl  string `mapstructure:"contract_compile_url"`
	MonitorPort         int    `mapstructure:"monitor_port"`
	ManageBackendApiKey string `mapstructure:"manage_backend_api_key"`
	LoginExpireTime     int    `mapstructure:"login_expire_time"`
	AdminPassword       string `mapstructure:"admin_password"`
}

// RwSet 读写集
type RwSet struct {
	Index        int    `json:"index"`
	Key          string `json:"key"`
	Value        []byte `json:"value"`
	ContractName string `json:"contractName"`
}

// RwSet 读写集
type TxEventData struct {
	Index        int    `json:"index"`
	Key          string `json:"key"`
	Value        string `json:"value"`
	ContractName string `json:"contractName"`
}

// UserKeyConf 认证配置
type UserKeyConf struct {
	UserSignKeyFilePath string `mapstructure:"user_sign_key_file_path"`
	UserSignCrtFilePath string `mapstructure:"user_sign_crt_file_path"`
	UserTlsKeyFilePath  string `mapstructure:"user_tls_key_file_path"`
	UserTlsCrtFilePath  string `mapstructure:"user_tls_crt_file_path"`
	UserEncKeyFilePath  string `mapstructure:"user_enc_key_file_path"`
	UserEncCrtFilePath  string `mapstructure:"user_enc_crt_file_path"`
}

// DBConf 数据库配置
type DBConf struct {
	Host        string `mapstructure:"host"`
	Port        string `mapstructure:"port"`
	Database    string `mapstructure:"database"`
	Username    string `mapstructure:"username"`
	Password    string `mapstructure:"password"`
	Prefix      string `mapstructure:"prefix"`
	DbProvider  string `mapstructure:"db_provider"`
	MaxByteSize int    `mapstructure:"max_byte_size"`
	MaxPoolSize int    `mapstructure:"max_pool_size"`
}

// LogConf 日志配置
type LogConf struct {
	LogLevelDefault string            `mapstructure:"log_level_default"`
	LogLevels       map[string]string `mapstructure:"log_levels"`
	FilePath        string            `mapstructure:"file_path"`
	MaxAge          int               `mapstructure:"max_age"`
	RotationTime    int               `mapstructure:"rotation_time"`
	LogInConsole    bool              `mapstructure:"log_in_console"`
	ShowColor       bool              `mapstructure:"show_color"`
}

// AlarmerConf 告警配置
type AlarmerConf struct {
	DingEnable        bool   `mapstructure:"ding_enable"`
	WechatEnable      bool   `mapstructure:"wechat_enable"`
	DingAccessToken   string `mapstructure:"ding_access_token"`   // token
	WechatAccessToken string `mapstructure:"wechat_access_token"` // token
	Prefix            string `mapstructure:"prefix"`              // token
}

// ChainConf 链基础配置
type ChainConf struct {
	Subscribe      bool              `mapstructure:"subscribe"` // 是否开启订阅
	ShowConfig     bool              `mapstructure:"show_config"`
	IsMainChain    bool              `mapstructure:"is_main_chain"`
	MainChainName  string            `mapstructure:"main_chain_name"`
	MainChainId    string            `mapstructure:"main_chain_id"`
	ContractConfig []*ContractConfig `mapstructure:"contract_config"`
}

// ContractConfig 合约配置
type ContractConfig struct {
	ChainId         string `mapstructure:"chain_id"`               // 链id
	MainDIDContract string `mapstructure:"main_did_contract_name"` // 主DID合约名称
}

// SensitiveConf 敏感词配置
type SensitiveConf struct {
	Enable    bool   `mapstructure:"enable"`
	SecretId  string `mapstructure:"secret_id"`
	SecretKey string `mapstructure:"secret_ey"`
}

// Config 整体配置
type Config struct {
	WebConf       *WebConf       `mapstructure:"web"`
	ChainsConfig  []*ChainConfig `mapstructure:"chains"`
	DBConf        *DBConf        `mapstructure:"db"`
	PProf         *PProf         `mapstructure:"pprof"`
	LogConf       *LogConf       `mapstructure:"log"`
	ChainConf     *ChainConf     `mapstructure:"chain"`
	AlarmerConf   *AlarmerConf   `mapstructure:"alarmer"`
	MonitorConf   *MonitorConf   `mapstructure:"monitor"`
	SensitiveConf *SensitiveConf `mapstructure:"sensitive"`
	RedisDB       *RedisConfig   `mapstructure:"db_redis"`
}

type PProf struct {
	IsOpen bool   `mapstructure:"is_open"`
	Port   string `mapstructure:"port"`
}

type RedisConfig struct {
	Type             string   `mapstructure:"type"` // redis类型： cluster：集群模式/node：普通模式
	Host             []string `mapstructure:"host"` // 服务器地址:端口
	Username         string   `mapstructure:"username"`
	Password         string   `mapstructure:"password"`           // 密码
	Prefix           string   `mapstructure:"prefix"`             //缓存前缀
	PositionRankTime int      `mapstructure:"position_rank_time"` // #持仓列表缓存过期时间，默认10min，单位s
}

// ChainConfig 链配置数据
type ChainConfig struct {
	ChainId     string       `mapstructure:"chain_id"`
	AuthType    string       `mapstructure:"auth_type"`
	OrgId       string       `mapstructure:"org_id"`
	TlsModel    int          `mapstructure:"tls_model"`
	Tls         bool         `mapstructure:"tls"`
	HashType    string       `mapstructure:"hash_type"`
	NodesConfig []*NodeConf  `mapstructure:"nodes"`
	UserConf    *UserKeyConf `mapstructure:"user"`
}

// NodeConf 节点配置
type NodeConf struct {
	TlsHost string `mapstructure:"tls_host"`
	CaPaths string `mapstructure:"ca_paths"`
	Remotes string `mapstructure:"remotes"`
}

// ChainInfo 链配置数据
type ChainInfo struct {
	ChainId   string
	AuthType  string
	OrgId     string
	HashType  string
	TlsMode   int
	Tls       bool
	NodesList []*NodeInfo
	UserInfo  *UserInfo
}

// UserInfo 认证配置
type UserInfo struct {
	UserSignKey string
	UserSignCrt string
	UserTlsKey  string
	UserTlsCrt  string
	UserEncKey  string
	UserEncCrt  string
}

// NodeInfo 节点数据
type NodeInfo struct {
	Addr        string
	OrgCA       string
	TLSHostName string
}

// MonitorTxConf 监控交易配置
type MonitorTxConf struct {
	MaxTxNum int `mapstructure:"max_tx_num"`
	TxLimit  int `mapstructure:"tx_limit"`
}

// MonitorConf 监控配置
type MonitorConf struct {
	Enable            bool           `mapstructure:"enable"`
	MonitorTxConf     *MonitorTxConf `mapstructure:"monitor_tx"`
	SafeWordLimit     int            `mapstructure:"safe_word_limit"`
	MaximumHeightDiff int64          `mapstructure:"max_height_diff"`
	Interval          int            `mapstructure:"interval"`
	TryConnNum        int            `mapstructure:"try_conn_num"`
	ChainsConfig      []*ChainConfig `mapstructure:"chains"`
}
