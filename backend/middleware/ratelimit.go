package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// TokenBucket 令牌桶结构
type TokenBucket struct {
	capacity int           // 桶容量
	tokens   int           // 当前令牌数
	rate     time.Duration // 令牌生成速率
	mu       sync.Mutex    // 互斥锁
	ticker   *time.Ticker  // 定时器
}

// NewTokenBucket 创建令牌桶
// capacity: 桶容量
// rate: 令牌生成速率（多久生成一个令牌）
func NewTokenBucket(capacity int, rate time.Duration) *TokenBucket {
	tb := &TokenBucket{
		capacity: capacity,
		tokens:   capacity, // 初始填满
		rate:     rate,
		ticker:   time.NewTicker(rate),
	}

	// 启动令牌生成器
	go tb.generate()

	return tb
}

// generate 定期生成令牌
func (tb *TokenBucket) generate() {
	for range tb.ticker.C {
		tb.mu.Lock()
		// 如果未满，则添加一个令牌
		if tb.tokens < tb.capacity {
			tb.tokens++
		}
		tb.mu.Unlock()
	}
}

// Take 尝试获取一个令牌
// 返回true表示获取成功，false表示桶空
func (tb *TokenBucket) Take() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	if tb.tokens > 0 {
		tb.tokens--
		return true
	}
	return false
}

// Stop 停止令牌桶
func (tb *TokenBucket) Stop() {
	tb.ticker.Stop()
}

// 全局令牌桶实例
var globalBucket *TokenBucket

// InitRateLimiter 初始化限流器
// qps: 每秒允许的请求数
func InitRateLimiter(qps int) {
	// 计算生成速率：1秒 / qps = 每个令牌的生成间隔
	rate := time.Second / time.Duration(qps)
	globalBucket = NewTokenBucket(qps*2, rate) // 容量设为qps的2倍，允许一定突发
}

// RateLimit 限流中间件
// 基于令牌桶算法实现
func RateLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if globalBucket == nil {
			// 如果未初始化，使用默认配置（100 QPS）
			InitRateLimiter(100)
		}

		// 尝试获取令牌
		if !globalBucket.Take() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
