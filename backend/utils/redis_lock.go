package utils

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"course-system/config"
	"github.com/redis/go-redis/v9"
)

// RedisLock Redis分布式锁结构
// 基于Redis SET命令的NX和EX选项实现
type RedisLock struct {
	key        string        // 锁的键名
	value      string        // 锁的唯一标识（防止误删其他进程的锁）
	expiration time.Duration // 锁的过期时间（防止死锁）
}

// NewRedisLock 创建一个新的分布式锁
// 参数:
//   - key: 锁的键名（建议格式: "lock:resource:id"）
//   - expiration: 锁的过期时间（建议5-30秒）
// 返回:
//   - *RedisLock: 锁对象
func NewRedisLock(key string, expiration time.Duration) *RedisLock {
	return &RedisLock{
		key:        key,
		value:      generateLockValue(), // 生成唯一标识
		expiration: expiration,
	}
}

// Lock 尝试获取锁（阻塞式）
// 参数:
//   - ctx: 上下文（用于超时控制）
//   - retryInterval: 重试间隔（建议50-200毫秒）
//   - maxRetries: 最大重试次数（0表示无限重试）
// 返回:
//   - error: 获取锁失败时的错误
//
// 工作原理:
//  1. 使用SET key value NX EX命令原子性地设置锁
//  2. NX: 只在键不存在时设置（保证互斥）
//  3. EX: 设置过期时间（防止死锁）
func (l *RedisLock) Lock(ctx context.Context, retryInterval time.Duration, maxRetries int) error {
	retries := 0
	ticker := time.NewTicker(retryInterval)
	defer ticker.Stop()

	for {
		// 尝试获取锁
		// SET key value NX EX expiration
		// 返回OK表示成功，nil表示key已存在
		success, err := config.RedisClient.SetNX(ctx, l.key, l.value, l.expiration).Result()
		if err != nil {
			return fmt.Errorf("Redis操作失败: %v", err)
		}

		if success {
			// 成功获取锁
			return nil
		}

		// 检查是否超过最大重试次数
		if maxRetries > 0 {
			retries++
			if retries >= maxRetries {
				return errors.New("获取锁超时：超过最大重试次数")
			}
		}

		// 等待后重试
		select {
		case <-ctx.Done():
			// 上下文取消或超时
			return ctx.Err()
		case <-ticker.C:
			// 继续下一次重试
			continue
		}
	}
}

// TryLock 尝试获取锁（非阻塞）
// 返回:
//   - bool: true表示成功获取锁，false表示锁已被占用
//   - error: Redis操作失败时的错误
func (l *RedisLock) TryLock(ctx context.Context) (bool, error) {
	success, err := config.RedisClient.SetNX(ctx, l.key, l.value, l.expiration).Result()
	if err != nil {
		return false, fmt.Errorf("Redis操作失败: %v", err)
	}
	return success, nil
}

// Unlock 释放锁
// 使用Lua脚本保证原子性：只有锁的持有者才能释放锁
//
// 为什么使用Lua脚本？
//   避免以下竞态条件：
//   1. 进程A获取锁，锁value=A
//   2. 进程A执行超时，锁自动过期
//   3. 进程B获取锁，锁value=B
//   4. 进程A尝试释放锁，如果不检查value，会错误地释放B的锁
//
// 返回:
//   - error: 释放锁失败时的错误
func (l *RedisLock) Unlock(ctx context.Context) error {
	// Lua脚本：检查value是否匹配，匹配则删除
	script := redis.NewScript(`
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("del", KEYS[1])
		else
			return 0
		end
	`)

	result, err := script.Run(ctx, config.RedisClient, []string{l.key}, l.value).Result()
	if err != nil {
		return fmt.Errorf("释放锁失败: %v", err)
	}

	// result为0表示锁不存在或已被其他进程占用
	if result == int64(0) {
		return errors.New("释放锁失败：锁不存在或已被其他进程占用")
	}

	return nil
}

// Extend 延长锁的过期时间
// 用于处理时间较长的任务
// 参数:
//   - ctx: 上下文
//   - extension: 延长的时间
// 返回:
//   - error: 延长失败时的错误
func (l *RedisLock) Extend(ctx context.Context, extension time.Duration) error {
	// Lua脚本：检查value是否匹配，匹配则延长过期时间
	script := redis.NewScript(`
		if redis.call("get", KEYS[1]) == ARGV[1] then
			return redis.call("pexpire", KEYS[1], ARGV[2])
		else
			return 0
		end
	`)

	milliseconds := int64(extension / time.Millisecond)
	result, err := script.Run(ctx, config.RedisClient, []string{l.key}, l.value, milliseconds).Result()
	if err != nil {
		return fmt.Errorf("延长锁失败: %v", err)
	}

	if result == int64(0) {
		return errors.New("延长锁失败：锁不存在或已被其他进程占用")
	}

	return nil
}

// generateLockValue 生成锁的唯一标识
// 使用随机字节保证唯一性
// 返回:
//   - string: 16字节的十六进制字符串
func generateLockValue() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// WithLock 使用分布式锁执行函数（高阶函数）
// 自动处理加锁、解锁和错误恢复
// 参数:
//   - ctx: 上下文
//   - lockKey: 锁的键名
//   - expiration: 锁的过期时间
//   - fn: 需要在锁保护下执行的函数
// 返回:
//   - error: 执行过程中的错误
//
// 使用示例:
//   err := WithLock(ctx, "lock:course:123", 10*time.Second, func() error {
//       // 执行需要加锁的业务逻辑
//       return nil
//   })
func WithLock(ctx context.Context, lockKey string, expiration time.Duration, fn func() error) error {
	lock := NewRedisLock(lockKey, expiration)

	// 尝试获取锁（最多重试20次，每次间隔100ms）
	if err := lock.Lock(ctx, 100*time.Millisecond, 20); err != nil {
		return fmt.Errorf("获取锁失败: %v", err)
	}

	// 确保锁被释放
	defer func() {
		unlockCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if err := lock.Unlock(unlockCtx); err != nil {
			// 记录日志，但不影响业务结果
			fmt.Printf("释放锁时出错: %v\n", err)
		}
	}()

	// 执行业务函数
	return fn()
}
