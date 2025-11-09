package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 日志中间件
// 记录每个请求的方法、路径、状态码和耗时
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 请求路径
		path := c.Request.URL.Path
		// 请求方法
		method := c.Request.Method

		// 处理请求
		c.Next()

		// 结束时间
		endTime := time.Now()
		// 执行时间
		latency := endTime.Sub(startTime)

		// 状态码
		statusCode := c.Writer.Status()

		// 客户端IP
		clientIP := c.ClientIP()

		// 记录日志
		log.Printf("[GIN] %s | %3d | %13v | %15s | %-7s %s",
			endTime.Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
	}
}
