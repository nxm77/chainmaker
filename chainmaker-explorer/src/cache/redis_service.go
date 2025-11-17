package cache

// import (
// 	"chainmaker_web/src/config"
// 	"context"
// 	"encoding/json"
// 	"errors"
// 	"fmt"
// 	"strconv"
// 	"time"

// 	lru "github.com/hashicorp/golang-lru"
// 	"github.com/redis/go-redis/v9"
// )

// // 在初始化时设置默认过期时间
// const DefaultExpireDuration = 1 * time.Hour
// const RedisTypeCluster = "cluster"

// // 统一缓存服务封装
// type RedisService struct {
// 	Client  redis.Cmdable
// 	Prefix  string
// 	ChainID string
// 	Expire  time.Duration
// }

// var GlobalRedisService *RedisService // 全局实例

// // 初始化Redis服务
// func InitRedisService(cfg *config.RedisConfig) {
// 	var client redis.Cmdable

// 	if cfg.Type == RedisTypeCluster {
// 		client = redis.NewClusterClient(&redis.ClusterOptions{
// 			Addrs:    cfg.Host,
// 			Password: cfg.Password,
// 			Username: cfg.Username,
// 		})
// 	} else {
// 		client = redis.NewClient(&redis.Options{
// 			Addr:     cfg.Host[0],
// 			Password: cfg.Password,
// 			Username: cfg.Username,
// 		})
// 	}

// 	// 健康检查
// 	if err := client.Ping(context.Background()).Err(); err != nil {
// 		panic(fmt.Sprintf("Redis init failed: %v", err))
// 	}

// 	GlobalCacheInstance, _ = lru.New(1024) // 缓存大小设置为1024

// 	GlobalRedisService = &RedisService{
// 		Client: client,
// 		Prefix: cfg.Prefix,
// 	}
// }

// // 链式调用方法 ---------------------------------------------------

// func (r *RedisService) WithChainID(chainID string) *RedisService {
// 	return &RedisService{
// 		Client:  r.Client,
// 		Prefix:  r.Prefix,
// 		ChainID: chainID,
// 		Expire:  DefaultExpireDuration, // 默认过期时间
// 	}
// }

// func (r *RedisService) WithExpire(duration time.Duration) *RedisService {
// 	return &RedisService{
// 		Client:  r.Client,
// 		Prefix:  r.Prefix,
// 		ChainID: r.ChainID,
// 		Expire:  duration,
// 	}
// }

// // 核心操作方法 ---------------------------------------------------

// // 生成带前缀的key（强制包含chainID）
// func (r *RedisService) buildKey(keyFormat string, parts ...interface{}) string {
// 	// 新格式：prefix_chainID_ + 业务键格式
// 	baseKey := fmt.Sprintf("%s_%s_%s", r.ChainID, r.Prefix, keyFormat)
// 	return fmt.Sprintf(baseKey, parts...)
// }

// // Set 通用设置（自动处理JSON序列化）
// // @param keyFormat 业务键格式（如"key1"）
// // @param value 要存储的值（可以是任意类型）
// // @param parts 业务键的其他部分（如"key2"）
// func (r *RedisService) Set(keyFormat string, value interface{}, parts ...interface{}) error {
// 	if r.ChainID == "" {
// 		return errors.New("chainID must be set via WithChainID()")
// 	}

// 	// 处理过期时间
// 	var expire time.Duration
// 	switch {
// 	case r.Expire > 0:
// 		expire = r.Expire
// 	case r.Expire < 0:
// 		return errors.New("invalid negative expiration duration")
// 	default:
// 		expire = 0 // 保持redis原生语义
// 	}

// 	// 序列化逻辑
// 	var val string
// 	switch v := value.(type) {
// 	case string:
// 		val = v
// 	case []byte:
// 		val = string(v)
// 	default:
// 		b, err := json.Marshal(value)
// 		if err != nil {
// 			return fmt.Errorf("marshal error: %w", err)
// 		}
// 		val = string(b)
// 	}

// 	return r.Client.Set(context.Background(), r.buildKey(keyFormat, parts...), val, expire).Err()
// }

// // Get 通用获取（自动处理JSON反序列化）
// func (r *RedisService) Get(keyFormat string, out interface{}, parts ...interface{}) error {
// 	key := r.buildKey(keyFormat, parts...)

// 	val, err := r.Client.Get(context.Background(), key).Result()
// 	if err != nil {
// 		return err
// 	}

// 	if out == nil {
// 		return nil
// 	}

// 	switch o := out.(type) {
// 	case *string:
// 		*o = val
// 		return nil
// 	case *[]byte:
// 		*o = []byte(val)
// 		return nil
// 	default:
// 		return json.Unmarshal([]byte(val), out)
// 	}
// }

// // GetInt 获取int值
// func (r *RedisService) GetInt(keyFormat string, parts ...interface{}) (int64, error) {
// 	key := r.buildKey(keyFormat, parts...)
// 	val, err := r.Client.Get(context.Background(), key).Result()
// 	if err != nil {
// 		return 0, err
// 	}
// 	return strconv.ParseInt(val, 10, 64)
// }

// // MustGetInt 获取int值（带默认值）
// func (r *RedisService) MustGetInt(keyFormat string, parts ...interface{}) int64 {
// 	val, _ := r.GetInt(keyFormat, parts...)
// 	return val
// }

// // 常用键定义（保持原有常量）
// const (
// 	KeyDelayedUpdateData = "elay_update_data_%s" // 自动包含前缀和chainID
// 	KeyDBAccountTotal    = "db_account_total"
// 	// 其他键定义...
// )

// // 使用示例 ---------------------------------------------------

// // 设置异步更新数据（链式调用）
// // func SetDelayedData(chainID string, blockHeight int64, data interface{}) error {
// // 	return GlobalRedisService.
// // 		WithChainID(chainID).
// // 		WithExpire(CacheExpirationMin10).
// // 		Set(context.Background(),
// // 			KeyDelayedUpdateData,
// // 			data,
// // 			strconv.FormatInt(blockHeight, 10))
// // }

// // // 获取账户总数（改进版）

// // func GetAccountTotal(chainID string) int64 {
// // 	return GlobalRedisService.
// // 		WithChainID(chainID).
// // 		MustGetInt(context.Background(), KeyDBAccountTotal)
// // }

// // // 实际业务中的调用示例
// // func BusinessExample() {
// // 	// 设置数据
// // 	err := GlobalRedisService.
// // 		WithChainID("chain123").
// // 		WithExpire(time.Hour).
// // 		Set(context.Background(), KeyContractInfo, &ContractInfo{
// // 			Name: "MyContract",
// // 			Addr: "0x123...",
// // 		}, "contractKey123")

// // 	// 获取数据
// // 	var info ContractInfo
// // 	err = GlobalRedisService.
// // 		WithChainID("chain123").
// // 		Get(context.Background(), KeyContractInfo, &info, "contractKey123")
// // }
