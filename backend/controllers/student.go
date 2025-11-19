package controllers

import (
	"context"
	"course-system/config"
	"course-system/models"
	"course-system/utils"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// StudentRegister 学生注册
// POST /api/student/register/
// 请求体: {username, phone, sms_code}
func StudentRegister(c *gin.Context) {
	// 定义接收JSON数据的结构体
	var req struct {
		Username string `json:"username" binding:"required"` // 用户名，必填
		Phone    string `json:"phone" binding:"required"`    // 手机号，必填
		SMSCode  string `json:"sms_code" binding:"required"` // 短信验证码，必填
	}

	// 绑定JSON数据到结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 验证手机号格式
	if len(req.Phone) != 11 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "手机号格式不正确"})
		return
	}

	// 验证短信验证码
	if !utils.VerifySMSCode(req.Phone, req.SMSCode, "register") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "短信验证码错误或已过期"})
		return
	}

	// 检查用户名是否已存在
	var existingStudent models.Student
	if err := config.DB.Where("username = ?", req.Username).First(&existingStudent).Error; err == nil {
		// 找到了用户，说明用户名已存在
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	// 检查手机号是否已存在
	if err := config.DB.Where("phone = ?", req.Phone).First(&existingStudent).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "手机号已被注册"})
		return
	}

	// 创建新学生记录
	student := models.Student{
		Username: req.Username,
		Phone:    req.Phone,
	}

	// 保存到数据库
	if err := config.DB.Create(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	}

	// 注册成功，生成JWT Token
	token, err := utils.GenerateToken(student.ID, "student")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	// 返回token和用户信息
	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
		"token":   token,
		"user": gin.H{
			"id":       student.ID,
			"username": student.Username,
			"phone":    student.Phone,
			"role":     "student",
		},
	})
}

// StudentLogin 学生登录
// POST /api/student/login/
// 请求体: {phone, sms_code}
// 注意: 登录前需要先通过图形验证码验证，然后才能获取短信验证码
func StudentLogin(c *gin.Context) {
	var req struct {
		Phone   string `json:"phone" binding:"required"`    // 手机号
		SMSCode string `json:"sms_code" binding:"required"` // 短信验证码
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 验证手机号格式
	if len(req.Phone) != 11 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "手机号格式不正确"})
		return
	}

	// 验证短信验证码
	if !utils.VerifySMSCode(req.Phone, req.SMSCode, "login") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "短信验证码错误或已过期"})
		return
	}

	// 查找学生
	var student models.Student
	if err := config.DB.Where("phone = ?", req.Phone).First(&student).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "该手机号未注册"})
		return
	}

	// 登录成功，生成JWT Token
	token, err := utils.GenerateToken(student.ID, "student")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	// 返回token和用户信息
	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   token,
		"user": gin.H{
			"id":       student.ID,
			"username": student.Username,
			"phone":    student.Phone,
			"role":     "student",
		},
	})
}

// GetCourses 获取所有课程（学生视角）
// GET /api/student/courses/
// 返回课程列表，包含是否已选、是否满员等信息
//
// 性能优化：
//  1. 使用enrolled字段代替COUNT查询，减少数据库负载
//  2. 使用Preload预加载教师信息，避免N+1查询问题
func GetCourses(c *gin.Context) {
	// 获取当前登录学生的ID
	studentIDInterface, _ := c.Get("user_id")
	studentID := studentIDInterface.(int)

	// 查询所有课程（预加载教师信息，优化性能）
	// 注意：这里为了简化，暂时逐个查询教师信息
	// 在实际生产环境中，可以考虑一次性查询所有教师，然后在内存中关联
	var courses []models.Course
	if err := config.DB.Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取课程失败"})
		return
	}

	// 一次性查询当前学生的所有选课记录（优化性能）
	var enrollments []models.Enrollment
	config.DB.Where("student_id = ?", studentID).Find(&enrollments)

	// 将选课记录转换为map，便于快速查找
	enrolledCourses := make(map[int]bool)
	for _, enrollment := range enrollments {
		enrolledCourses[enrollment.CourseID] = true
	}

	// 构建返回的课程列表
	result := []gin.H{}
	for _, course := range courses {
		// 查询该课程的教师信息
		var teacher models.Teacher
		config.DB.First(&teacher, course.TeacherID)

		// 使用enrolled字段（避免COUNT查询，提升性能）
		enrolledCount := course.Enrolled

		// 检查当前学生是否已选这门课（从map中查找，O(1)时间复杂度）
		isEnrolled := enrolledCourses[course.ID]

		// 判断是否已满
		isFull := enrolledCount >= course.Capacity

		// 添加到结果列表
		result = append(result, gin.H{
			"id":          course.ID,
			"name":        course.Name,
			"description": course.Description,
			"teacher":     teacher.Username, // 教师名称
			"teacher_id":  course.TeacherID,
			"capacity":    course.Capacity, // 课程容量
			"enrolled":    enrolledCount,   // 已选人数（直接使用字段，避免COUNT）
			"is_enrolled": isEnrolled,      // 是否已选
			"is_full":     isFull,          // 是否已满
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"courses": result,
	})
}

// GetMyCourses 获取我的课程
// GET /api/student/my-courses/
// 返回当前学生已选的课程列表
func GetMyCourses(c *gin.Context) {
	// 获取当前学生ID
	studentID, _ := c.Get("user_id")

	// 查询该学生的所有选课记录
	var enrollments []models.Enrollment
	if err := config.DB.Where("student_id = ?", studentID).Find(&enrollments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取课程失败"})
		return
	}

	// 构建课程详情列表
	result := []gin.H{}
	for _, enrollment := range enrollments {
		// 查询课程信息
		var course models.Course
		if err := config.DB.First(&course, enrollment.CourseID).Error; err != nil {
			continue // 课程不存在，跳过
		}

		// 查询教师信息
		var teacher models.Teacher
		config.DB.First(&teacher, course.TeacherID)

		// 添加到结果
		result = append(result, gin.H{
			"course_id":   course.ID,
			"course_name": course.Name,
			"description": course.Description,
			"teacher":     teacher.Username,
			"enrolled_at": enrollment.EnrolledAt.Format("2006-01-02 15:04:05"), // 格式化时间
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"courses": result,
	})
}

// EnrollCourse 选课（高并发优化版）
// POST /api/student/enroll/
// 请求体: {course_id}
//
// 并发控制策略：
//  1. Redis分布式锁：防止同一课程的并发选课冲突
//  2. 乐观锁（Version字段）：防止超卖，确保库存一致性
//  3. 数据库事务：保证选课记录和课程enrolled字段的原子性更新
func EnrollCourse(c *gin.Context) {
	var req struct {
		CourseID int `json:"course_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 获取当前学生ID
	studentIDInterface, _ := c.Get("user_id")
	studentID := studentIDInterface.(int)

	// ============ 步骤1: 基础数据验证（无需加锁） ============

	// 检查课程是否存在
	var course models.Course
	if err := config.DB.First(&course, req.CourseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "课程不存在"})
		return
	}

	// 检查是否已选过该课程（防止重复选课）
	var existingEnrollment models.Enrollment
	if err := config.DB.Where("student_id = ? AND course_id = ?", studentID, req.CourseID).
		First(&existingEnrollment).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已经选过该课程"})
		return
	}

	// 检查选课时间冲突
	// 查询新课程和学生已选课程的上课时间，判断是否有时间重叠
	hasConflict, conflictMsg, err := utils.CheckScheduleConflict(studentID, req.CourseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("检测时间冲突失败: %v", err)})
		return
	}
	if hasConflict {
		c.JSON(http.StatusBadRequest, gin.H{"error": conflictMsg})
		return
	}

	// ============ 步骤2: 使用Redis分布式锁保护选课操作 ============

	// 创建分布式锁，锁的key为 "lock:course:{课程ID}"
	// 锁的超时时间设置为10秒，防止死锁
	lockKey := fmt.Sprintf("lock:course:%d", req.CourseID)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// 使用高阶函数WithLock自动处理加锁和解锁
	err = utils.WithLock(ctx, lockKey, 10*time.Second, func() error {
		// ============ 步骤3: 在锁保护下执行选课逻辑 ============

		// 开启数据库事务（保证数据一致性）
		return config.DB.Transaction(func(tx *gorm.DB) error {
			// 3.1 使用乐观锁重新查询课程信息（FOR UPDATE确保行锁）
			var currentCourse models.Course
			if err := tx.Clauses().First(&currentCourse, req.CourseID).Error; err != nil {
				return fmt.Errorf("课程不存在")
			}

			// 3.2 检查课程容量（使用enrolled字段，避免COUNT查询）
			if currentCourse.Enrolled >= currentCourse.Capacity {
				return fmt.Errorf("课程已满")
			}

			// 3.3 创建选课记录
			enrollment := models.Enrollment{
				StudentID: studentID,
				CourseID:  req.CourseID,
			}
			if err := tx.Create(&enrollment).Error; err != nil {
				return fmt.Errorf("创建选课记录失败: %v", err)
			}

			// 3.4 使用乐观锁更新课程的enrolled字段和version
			// SQL: UPDATE courses SET enrolled = enrolled + 1, version = version + 1
			//      WHERE id = ? AND version = ?
			// 如果version不匹配，说明有其他进程修改了数据，更新失败
			result := tx.Model(&models.Course{}).
				Where("id = ? AND version = ?", req.CourseID, currentCourse.Version).
				Updates(map[string]interface{}{
					"enrolled": gorm.Expr("enrolled + ?", 1),
					"version":  gorm.Expr("version + ?", 1),
				})

			if result.Error != nil {
				return fmt.Errorf("更新课程信息失败: %v", result.Error)
			}

			// 检查是否真的更新了（乐观锁校验）
			if result.RowsAffected == 0 {
				// version不匹配，说明有并发冲突
				return fmt.Errorf("选课失败，请重试（并发冲突）")
			}

			// 选课成功
			return nil
		})
	})

	// ============ 步骤4: 处理结果 ============

	if err != nil {
		// 根据错误类型返回不同的HTTP状态码
		if err.Error() == "课程已满" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else if err.Error() == "课程不存在" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "选课成功",
	})
}

// DropCourse 退课（事务优化版）
// POST /api/student/drop/
// 请求体: {course_id}
//
// 并发控制策略：
//  1. Redis分布式锁：防止同一课程的并发退课冲突
//  2. 数据库事务：保证选课记录删除和课程enrolled字段的原子性更新
func DropCourse(c *gin.Context) {
	var req struct {
		CourseID int `json:"course_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 获取当前学生ID
	studentIDInterface, _ := c.Get("user_id")
	studentID := studentIDInterface.(int)

	// ============ 步骤1: 基础数据验证 ============

	// 查找选课记录
	var enrollment models.Enrollment
	if err := config.DB.Where("student_id = ? AND course_id = ?", studentID, req.CourseID).
		First(&enrollment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到选课记录"})
		return
	}

	// ============ 步骤2: 使用Redis分布式锁保护退课操作 ============

	lockKey := fmt.Sprintf("lock:course:%d", req.CourseID)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := utils.WithLock(ctx, lockKey, 10*time.Second, func() error {
		// ============ 步骤3: 在锁保护下执行退课逻辑 ============

		// 开启数据库事务
		return config.DB.Transaction(func(tx *gorm.DB) error {
			// 3.1 删除选课记录
			if err := tx.Delete(&enrollment).Error; err != nil {
				return fmt.Errorf("删除选课记录失败: %v", err)
			}

			// 3.2 更新课程的enrolled字段（减1）
			// SQL: UPDATE courses SET enrolled = enrolled - 1, version = version + 1
			//      WHERE id = ?
			result := tx.Model(&models.Course{}).
				Where("id = ?", req.CourseID).
				Updates(map[string]interface{}{
					"enrolled": gorm.Expr("enrolled - ?", 1),
					"version":  gorm.Expr("version + ?", 1),
				})

			if result.Error != nil {
				return fmt.Errorf("更新课程信息失败: %v", result.Error)
			}

			// 确保enrolled不会变成负数
			tx.Model(&models.Course{}).
				Where("id = ? AND enrolled < 0", req.CourseID).
				Update("enrolled", 0)

			return nil
		})
	})

	// ============ 步骤4: 处理结果 ============

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "退课成功",
	})
}

// GetScheduleTable 获取学生课表（二维数组）
// GET /api/student/schedule/
// 可选参数: ?week=1 (当前周次，用于过滤课程)
func GetScheduleTable(c *gin.Context) {
	// 获取当前学生ID
	studentIDInterface, _ := c.Get("user_id")
	studentID := studentIDInterface.(int)

	// 获取可选的周次参数
	var currentWeek int
	weekParam := c.Query("week")
	if weekParam != "" {
		fmt.Sscanf(weekParam, "%d", &currentWeek)
	}

	// 调用工具函数获取课表
	schedule, err := utils.GetStudentScheduleTable(studentID, currentWeek)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"schedule": schedule,
	})
}
