package config

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
)

// RedisClient 全局Redis客户端实例
// 用于分布式锁、缓存等场景
var RedisClient *redis.Client

// RedisConfig Redis连接配置
type RedisConfig struct {
	Host     string // Redis服务器地址
	Port     string // Redis端口
	Password string // Redis密码（无密码时为空字符串）
	DB       int    // 数据库编号（0-15）
}

// InitRedis 初始化Redis连接
// 参数:
//   - cfg: Redis配置信息
//
// 返回:
//   - error: 连接失败时的错误信息
func InitRedis(cfg RedisConfig) error {
	// 创建Redis客户端
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port), // Redis地址
		Password: cfg.Password,                             // 密码
		DB:       cfg.DB,                                   // 使用的数据库

		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled, // 显式禁用（核心配置）
			// 其他可选字段（禁用后无需配置，默认即可）：
			// Interval: 10 * time.Second, // 通知查询间隔（禁用后无效）
		},

		// 连接池配置
		PoolSize:     100,             // 连接池最大连接数（支持高并发）
		MinIdleConns: 10,              // 最小空闲连接数
		MaxRetries:   3,               // 最大重试次数
		DialTimeout:  5 * time.Second, // 连接超时
		ReadTimeout:  3 * time.Second, // 读超时
		WriteTimeout: 3 * time.Second, // 写超时
		PoolTimeout:  4 * time.Second, // 连接池超时
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RedisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("Redis连接失败: %v", err)
	}

	fmt.Println("✓ Redis连接成功")
	return nil
}

// CloseRedis 关闭Redis连接
// 在应用退出时调用
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}
