package middleware

import (
	"course-system/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件
// 用于验证用户是否已登录（通过JWT）
// 从Authorization header中读取token并解析
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头获取Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未提供认证token",
			})
			c.Abort()
			return
		}

		// 解析Bearer Token
		// 格式: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "token格式错误",
			})
			c.Abort()
			return
		}

		// 解析token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "无效的token",
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		// 继续处理请求
		c.Next()
	}
}

// RequireAuth 认证中间件（兼容旧代码，实际使用JWTAuth）
func RequireAuth() gin.HandlerFunc {
	return JWTAuth()
}

// RequireStudent 学生权限中间件
// 确保当前登录用户是学生
func RequireStudent() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取角色
		roleInterface, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未认证",
			})
			c.Abort()
			return
		}

		role, ok := roleInterface.(string)
		if !ok || role != "student" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "需要学生权限",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireTeacher 教师权限中间件
// 确保当前登录用户是教师
func RequireTeacher() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取角色
		roleInterface, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未认证",
			})
			c.Abort()
			return
		}

		role, ok := roleInterface.(string)
		if !ok || role != "teacher" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "需要教师权限",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
