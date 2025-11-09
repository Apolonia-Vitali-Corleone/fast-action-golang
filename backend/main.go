package main

import (
	"course-system/config"
	"course-system/controllers"
	"course-system/middleware"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func main() {
	// ========== 1. 初始化数据库 ==========
	dbConfig := config.DBConfig{
		Host:     "localhost", // 数据库主机地址
		Port:     "3306",      // MySQL默认端口
		User:     "root",      // 数据库用户名（根据实际情况修改）
		Password: "password",  // 数据库密码（根据实际情况修改）
		DBName:   "course_system", // 数据库名称
	}

	// 连接数据库
	if err := config.InitDB(dbConfig); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// ========== 2. 创建Gin应用 ==========
	// gin.Default()会自动添加Logger和Recovery中间件
	r := gin.Default()

	// ========== 3. 配置CORS跨域 ==========
	// 允许前端（http://localhost:5173）访问后端API
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // 允许的来源（前端地址）
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // 允许的HTTP方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"}, // 允许的请求头
		AllowCredentials: true, // 允许携带Cookie（session需要）
	}))

	// ========== 4. 配置Session ==========
	// 使用Cookie存储session，密钥用于加密cookie
	store := cookie.NewStore([]byte("secret-key-change-in-production")) // 生产环境请修改密钥
	r.Use(sessions.Sessions("session", store)) // session名称为"session"

	// ========== 5. 配置路由 ==========
	// API基础路径组
	api := r.Group("/api")
	{
		// ---------- 学生相关路由 ----------
		student := api.Group("/student")
		{
			// 公开接口（无需登录）
			student.POST("/register/", controllers.StudentRegister) // 学生注册
			student.POST("/login/", controllers.StudentLogin)       // 学生登录

			// 需要登录且是学生身份的接口
			// 使用RequireAuth中间件验证登录，RequireStudent验证学生身份
			student.GET("/courses/", middleware.RequireAuth(), middleware.RequireStudent(), controllers.GetCourses) // 获取所有课程
			student.GET("/my-courses/", middleware.RequireAuth(), middleware.RequireStudent(), controllers.GetMyCourses) // 获取我的课程
			student.POST("/enroll/", middleware.RequireAuth(), middleware.RequireStudent(), controllers.EnrollCourse)   // 选课
			student.POST("/drop/", middleware.RequireAuth(), middleware.RequireStudent(), controllers.DropCourse)       // 退课
		}

		// ---------- 教师相关路由 ----------
		teacher := api.Group("/teacher")
		{
			// 公开接口
			teacher.POST("/register/", controllers.TeacherRegister) // 教师注册
			teacher.POST("/login/", controllers.TeacherLogin)       // 教师登录

			// 需要登录且是教师身份的接口
			teacher.GET("/courses/", middleware.RequireAuth(), middleware.RequireTeacher(), controllers.GetTeacherCourses) // 获取我的课程
			teacher.POST("/courses/create/", middleware.RequireAuth(), middleware.RequireTeacher(), controllers.CreateCourse) // 创建课程
			teacher.DELETE("/courses/:id/delete/", middleware.RequireAuth(), middleware.RequireTeacher(), controllers.DeleteCourse) // 删除课程
			teacher.GET("/courses/:id/students/", middleware.RequireAuth(), middleware.RequireTeacher(), controllers.GetCourseStudents) // 获取选课学生
		}

		// ---------- 通用路由 ----------
		// 获取当前登录用户信息
		api.GET("/current-user/", func(c *gin.Context) {
			// 获取session
			session := sessions.Default(c)
			userID := session.Get("user_id")

			// 如果未登录
			if userID == nil {
				c.JSON(401, gin.H{"error": "未登录"})
				return
			}

			// 获取用户信息
			role := session.Get("role")
			username := session.Get("username")

			c.JSON(200, gin.H{
				"user": gin.H{
					"id":       userID,
					"username": username,
					"role":     role,
				},
			})
		})

		// 退出登录
		api.POST("/logout/", func(c *gin.Context) {
			// 清空session
			session := sessions.Default(c)
			session.Clear()        // 清除所有session数据
			session.Save()         // 保存更改

			c.JSON(200, gin.H{
				"message": "已退出",
			})
		})
	}

	// ========== 6. 启动服务器 ==========
	// 监听8000端口
	log.Println("服务器启动在 http://localhost:8000")
	if err := r.Run(":8000"); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
