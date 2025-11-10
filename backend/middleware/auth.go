package middleware

import (
	"course-system/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT认证中间件（含Token自动刷新）
// 功能:
//  1. 验证用户是否已登录（通过JWT）
//  2. 自动刷新即将过期的Token（剩余有效期<2小时时）
//  3. 通过响应头 X-New-Token 返回新Token
//
// 前端需要处理：
//   - 检查响应头中是否有 X-New-Token
//   - 如果有，则更新本地存储的Token
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ========== 步骤1: 从请求头获取Token ==========
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未提供认证token",
			})
			c.Abort()
			return
		}

		// ========== 步骤2: 解析Bearer Token ==========
		// 格式: "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "token格式错误",
			})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// ========== 步骤3: 验证Token ==========
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "无效的token",
			})
			c.Abort()
			return
		}

		// ========== 步骤4: 将用户信息存储到上下文中 ==========
		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)

		// ========== 步骤5: 检查是否需要刷新Token ==========
		// 如果Token剩余有效期 < 2小时，则自动刷新
		refreshThreshold := 2 * time.Hour
		if utils.ShouldRefreshToken(claims, refreshThreshold) {
			// 生成新Token
			newToken, err := utils.RefreshToken(claims)
			if err == nil {
				// 通过响应头返回新Token
				// 前端应监听此响应头并更新本地存储的Token
				c.Header("X-New-Token", newToken)
			}
			// 即使刷新失败也不影响当前请求（因为旧Token仍然有效）
		}

		// ========== 步骤6: 继续处理请求 ==========
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
