package main

import (
	"course-system/config"
	"course-system/controllers"
	"course-system/middleware"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// ========== 1. 初始化数据库 ==========
	dbConfig := config.DBConfig{
		Host:     "192.168.233.136", // 数据库主机地址
		Port:     "3306",            // MySQL默认端口
		User:     "root",            // 数据库用户名（根据实际情况修改）
		Password: "1234",            // 数据库密码（根据实际情况修改）
		DBName:   "course_system",   // 数据库名称
	}

	// 连接数据库（含连接池优化）
	if err := config.InitDB(dbConfig); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}

	// ========== 2. 初始化Redis（用于分布式锁和缓存） ==========
	redisConfig := config.RedisConfig{
		Host:     "192.168.233.136", // Redis服务器地址（根据实际情况修改）
		Port:     "6379",            // Redis默认端口
		Password: "",                // Redis密码（无密码则为空字符串）
		DB:       0,                 // 使用0号数据库
	}

	// 连接Redis（含连接池优化）
	if err := config.InitRedis(redisConfig); err != nil {
		log.Fatalf("Redis初始化失败: %v", err)
	}

	// ========== 3. 初始化限流器 ==========
	// 设置为每秒1000个请求（QPS=1000）
	// 支持500+并发用户同时选课
	middleware.InitRateLimiter(1000)

	// ========== 4. 创建Gin应用（不使用默认中间件） ==========
	r := gin.New()

	// ========== 5. 配置四层中间件链（按顺序） ==========
	// 第一层：Recovery - 捕获panic，防止服务崩溃
	r.Use(middleware.Recovery())

	// 第二层：Logger - 记录每个请求的日志
	r.Use(middleware.Logger())

	// 第三层：RateLimit - 令牌桶限流
	r.Use(middleware.RateLimit())

	// 第四层：CORS - 跨域资源共享
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},                             // 允许的来源（前端地址）
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},           // 允许的HTTP方法
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"}, // 允许的请求头（包含Authorization）
		ExposeHeaders:    []string{"X-New-Token"},                                       // 允许前端读取的响应头（用于Token自动刷新）
		AllowCredentials: true,                                                          // 允许携带凭证
	}))

	// ========== 6. 配置路由 ==========
	// API基础路径组
	api := r.Group("/api")
	{
		// ---------- 验证码相关路由（公开接口） ----------
		api.GET("/captcha/", controllers.GetCaptcha)      // 获取图形验证码
		api.POST("/sms/send/", controllers.SendSMSCode)   // 发送短信验证码

		// ---------- 学生相关路由 ----------
		student := api.Group("/student")
		{
			// 公开接口（无需登录）
			student.POST("/register/", controllers.StudentRegister) // 学生注册
			student.POST("/login/", controllers.StudentLogin)       // 学生登录

			// 需要登录且是学生身份的接口
			// 使用RequireAuth中间件验证登录，RequireStudent验证学生身份
			student.GET("/courses/", middleware.RequireAuth(), middleware.RequireStudent(), controllers.GetCourses)        // 获取所有课程
			student.GET("/my-courses/", middleware.RequireAuth(), middleware.RequireStudent(), controllers.GetMyCourses)   // 获取我的课程
			student.GET("/schedule/", middleware.RequireAuth(), middleware.RequireStudent(), controllers.GetScheduleTable) // 获取课表
			student.POST("/enroll/", middleware.RequireAuth(), middleware.RequireStudent(), controllers.EnrollCourse)      // 选课
			student.POST("/drop/", middleware.RequireAuth(), middleware.RequireStudent(), controllers.DropCourse)          // 退课
		}

		// ---------- 教师相关路由 ----------
		teacher := api.Group("/teacher")
		{
			// 公开接口
			teacher.POST("/register/", controllers.TeacherRegister) // 教师注册
			teacher.POST("/login/", controllers.TeacherLogin)       // 教师登录

			// 需要登录且是教师身份的接口
			teacher.GET("/courses/", middleware.RequireAuth(), middleware.RequireTeacher(), controllers.GetTeacherCourses)              // 获取我的课程
			teacher.POST("/courses/create/", middleware.RequireAuth(), middleware.RequireTeacher(), controllers.CreateCourse)           // 创建课程
			teacher.PUT("/courses/:id/update/", middleware.RequireAuth(), middleware.RequireTeacher(), controllers.UpdateCourse)        // 修改课程
			teacher.DELETE("/courses/:id/delete/", middleware.RequireAuth(), middleware.RequireTeacher(), controllers.DeleteCourse)     // 删除课程
			teacher.GET("/courses/:id/students/", middleware.RequireAuth(), middleware.RequireTeacher(), controllers.GetCourseStudents) // 获取选课学生
		}

		// ---------- 通用路由 ----------
		// 获取当前登录用户信息（需要JWT认证）
		api.GET("/current-user/", middleware.JWTAuth(), func(c *gin.Context) {
			// 从上下文获取用户信息（由JWT中间件设置）
			userID, _ := c.Get("user_id")
			role, _ := c.Get("role")

			c.JSON(200, gin.H{
				"user": gin.H{
					"id":   userID,
					"role": role,
				},
			})
		})

		// 退出登录（JWT无需服务端处理，由前端删除token即可）
		api.POST("/logout/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "已退出",
			})
		})
	}

	// ========== 7. 启动服务器 ==========
	// 监听8000端口
	log.Println("服务器启动在 http://localhost:8000")
	if err := r.Run(":8000"); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
