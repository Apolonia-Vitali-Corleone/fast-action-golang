package middleware

import (
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// RequireAuth 认证中间件
// 用于验证用户是否已登录
// 从session中读取用户ID和角色信息
func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前请求的session
		session := sessions.Default(c)

		// 从session中获取用户ID
		userID := session.Get("user_id")
		if userID == nil {
			// 如果session中没有用户ID，说明未登录
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "未登录",
			})
			c.Abort() // 终止请求处理
			return
		}

		// 从session中获取用户角色（student或teacher）
		role := session.Get("role")
		if role == nil {
			// 如果没有角色信息，返回错误
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "角色信息缺失",
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中，供后续处理函数使用
		c.Set("user_id", userID)
		c.Set("role", role)

		// 继续处理请求
		c.Next()
	}
}

// RequireStudent 学生权限中间件
// 确保当前登录用户是学生
func RequireStudent() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 先调用基础认证检查
		session := sessions.Default(c)
		role := session.Get("role")

		// 检查角色是否为student
		if role != "student" {
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
		// 先调用基础认证检查
		session := sessions.Default(c)
		role := session.Get("role")

		// 检查角色是否为teacher
		if role != "teacher" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "需要教师权限",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
