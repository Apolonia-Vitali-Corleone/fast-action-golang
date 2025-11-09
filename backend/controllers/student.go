package controllers

import (
	"course-system/config"
	"course-system/models"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// StudentRegister 学生注册
// POST /api/student/register/
// 请求体: {username, password, email}
func StudentRegister(c *gin.Context) {
	// 定义接收JSON数据的结构体
	var req struct {
		Username string `json:"username" binding:"required"` // 用户名，必填
		Password string `json:"password" binding:"required"` // 密码，必填
		Email    string `json:"email" binding:"required"`    // 邮箱，必填
	}

	// 绑定JSON数据到结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 检查用户名是否已存在
	var existingStudent models.Student
	if err := config.DB.Where("username = ?", req.Username).First(&existingStudent).Error; err == nil {
		// 找到了用户，说明用户名已存在
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	// 检查邮箱是否已存在
	if err := config.DB.Where("email = ?", req.Email).First(&existingStudent).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱已被使用"})
		return
	}

	// 使用bcrypt加密密码
	// bcrypt.DefaultCost是默认的加密强度
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 创建新学生记录
	student := models.Student{
		Username: req.Username,
		Password: string(hashedPassword), // 存储加密后的密码
		Email:    req.Email,
	}

	// 保存到数据库
	if err := config.DB.Create(&student).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	}

	// 注册成功
	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
	})
}

// StudentLogin 学生登录
// POST /api/student/login/
// 请求体: {username, password}
func StudentLogin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 查找学生
	var student models.Student
	if err := config.DB.Where("username = ?", req.Username).First(&student).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 验证密码
	// CompareHashAndPassword比较加密后的密码和明文密码
	if err := bcrypt.CompareHashAndPassword([]byte(student.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 登录成功，创建session
	session := sessions.Default(c)
	session.Set("user_id", student.ID)       // 保存用户ID
	session.Set("role", "student")            // 保存角色
	session.Set("username", student.Username) // 保存用户名
	session.Save()                            // 保存session到cookie

	// 返回用户信息
	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"user": gin.H{
			"id":       student.ID,
			"username": student.Username,
			"email":    student.Email,
			"role":     "student",
		},
	})
}

// GetCourses 获取所有课程（学生视角）
// GET /api/student/courses/
// 返回课程列表，包含是否已选、是否满员等信息
func GetCourses(c *gin.Context) {
	// 获取当前登录学生的ID
	session := sessions.Default(c)
	studentID := session.Get("user_id")

	// 查询所有课程
	var courses []models.Course
	if err := config.DB.Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取课程失败"})
		return
	}

	// 构建返回的课程列表
	var result []gin.H
	for _, course := range courses {
		// 查询该课程的教师信息
		var teacher models.Teacher
		config.DB.First(&teacher, course.TeacherID)

		// 统计已选人数
		var enrolledCount int64
		config.DB.Model(&models.Enrollment{}).Where("course_id = ?", course.ID).Count(&enrolledCount)

		// 检查当前学生是否已选这门课
		var enrollment models.Enrollment
		isEnrolled := config.DB.Where("student_id = ? AND course_id = ?", studentID, course.ID).First(&enrollment).Error == nil

		// 判断是否已满
		isFull := int(enrolledCount) >= course.Capacity

		// 添加到结果列表
		result = append(result, gin.H{
			"id":          course.ID,
			"name":        course.Name,
			"description": course.Description,
			"teacher":     teacher.Username,      // 教师名称
			"teacher_id":  course.TeacherID,
			"capacity":    course.Capacity,
			"enrolled":    enrolledCount,         // 已选人数
			"is_enrolled": isEnrolled,            // 是否已选
			"is_full":     isFull,                // 是否已满
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
	session := sessions.Default(c)
	studentID := session.Get("user_id")

	// 查询该学生的所有选课记录
	var enrollments []models.Enrollment
	if err := config.DB.Where("student_id = ?", studentID).Find(&enrollments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取课程失败"})
		return
	}

	// 构建课程详情列表
	var result []gin.H
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

// EnrollCourse 选课
// POST /api/student/enroll/
// 请求体: {course_id}
func EnrollCourse(c *gin.Context) {
	var req struct {
		CourseID int `json:"course_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 获取当前学生ID
	session := sessions.Default(c)
	studentID := session.Get("user_id").(int)

	// 检查课程是否存在
	var course models.Course
	if err := config.DB.First(&course, req.CourseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "课程不存在"})
		return
	}

	// 检查是否已选过该课程
	var existingEnrollment models.Enrollment
	if err := config.DB.Where("student_id = ? AND course_id = ?", studentID, req.CourseID).First(&existingEnrollment).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已经选过该课程"})
		return
	}

	// 检查课程是否已满
	var enrolledCount int64
	config.DB.Model(&models.Enrollment{}).Where("course_id = ?", req.CourseID).Count(&enrolledCount)
	if int(enrolledCount) >= course.Capacity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "课程已满"})
		return
	}

	// 创建选课记录
	enrollment := models.Enrollment{
		StudentID: studentID,
		CourseID:  req.CourseID,
	}

	if err := config.DB.Create(&enrollment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "选课失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "选课成功",
	})
}

// DropCourse 退课
// POST /api/student/drop/
// 请求体: {course_id}
func DropCourse(c *gin.Context) {
	var req struct {
		CourseID int `json:"course_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 获取当前学生ID
	session := sessions.Default(c)
	studentID := session.Get("user_id")

	// 查找选课记录
	var enrollment models.Enrollment
	if err := config.DB.Where("student_id = ? AND course_id = ?", studentID, req.CourseID).First(&enrollment).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未找到选课记录"})
		return
	}

	// 删除选课记录
	if err := config.DB.Delete(&enrollment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "退课失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "退课成功",
	})
}
