package middleware

import (
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// Recovery 异常恢复中间件
// 捕获panic并返回500错误，防止服务崩溃
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 打印错误堆栈
				log.Printf("panic recovered: %v\n", err)
				log.Printf("stack trace:\n%s", debug.Stack())

				// 返回500错误
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("服务器内部错误: %v", err),
				})

				// 终止请求
				c.Abort()
			}
		}()

		// 处理请求
		c.Next()
	}
}
