package controllers

import (
	"course-system/config"
	"course-system/models"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// TeacherRegister 教师注册
// POST /api/teacher/register/
// 请求体: {username, password, email}
func TeacherRegister(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
		Email    string `json:"email" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 检查用户名是否已存在
	var existingTeacher models.Teacher
	if err := config.DB.Where("username = ?", req.Username).First(&existingTeacher).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	// 检查邮箱是否已存在
	if err := config.DB.Where("email = ?", req.Email).First(&existingTeacher).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱已被使用"})
		return
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "密码加密失败"})
		return
	}

	// 创建教师记录
	teacher := models.Teacher{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	}

	if err := config.DB.Create(&teacher).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "注册失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
	})
}

// TeacherLogin 教师登录
// POST /api/teacher/login/
// 请求体: {username, password}
func TeacherLogin(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 查找教师
	var teacher models.Teacher
	if err := config.DB.Where("username = ?", req.Username).First(&teacher).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(teacher.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	// 创建session
	session := sessions.Default(c)
	session.Set("user_id", teacher.ID)
	session.Set("role", "teacher")
	session.Set("username", teacher.Username)
	session.Save()

	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"user": gin.H{
			"id":       teacher.ID,
			"username": teacher.Username,
			"email":    teacher.Email,
			"role":     "teacher",
		},
	})
}

// GetTeacherCourses 获取教师的所有课程
// GET /api/teacher/courses/
// 返回当前教师创建的所有课程
func GetTeacherCourses(c *gin.Context) {
	// 获取当前教师ID
	session := sessions.Default(c)
	teacherID := session.Get("user_id")

	// 查询该教师的所有课程
	var courses []models.Course
	if err := config.DB.Where("teacher_id = ?", teacherID).Find(&courses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取课程失败"})
		return
	}

	// 构建课程列表，包含选课人数
	var result []gin.H
	for _, course := range courses {
		// 统计选课人数
		var enrolledCount int64
		config.DB.Model(&models.Enrollment{}).Where("course_id = ?", course.ID).Count(&enrolledCount)

		result = append(result, gin.H{
			"id":          course.ID,
			"name":        course.Name,
			"description": course.Description,
			"capacity":    course.Capacity,
			"enrolled":    enrolledCount, // 已选人数
			"created_at":  course.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"courses": result,
	})
}

// CreateCourse 创建课程
// POST /api/teacher/courses/create/
// 请求体: {name, description, capacity}
func CreateCourse(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`        // 课程名称，必填
		Description string `json:"description"`                    // 课程描述，可选
		Capacity    int    `json:"capacity" binding:"required,gt=0"` // 容量，必填且大于0
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 获取当前教师ID
	session := sessions.Default(c)
	teacherID := session.Get("user_id").(int)

	// 创建课程记录
	course := models.Course{
		Name:        req.Name,
		Description: req.Description,
		TeacherID:   teacherID,
		Capacity:    req.Capacity,
	}

	if err := config.DB.Create(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建课程失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功",
		"course": gin.H{
			"id":          course.ID,
			"name":        course.Name,
			"description": course.Description,
			"capacity":    course.Capacity,
		},
	})
}

// DeleteCourse 删除课程
// DELETE /api/teacher/courses/:id/delete/
// 删除指定ID的课程（必须是该教师创建的）
func DeleteCourse(c *gin.Context) {
	// 从URL参数中获取课程ID
	courseID := c.Param("id")

	// 获取当前教师ID
	session := sessions.Default(c)
	teacherID := session.Get("user_id")

	// 查找课程
	var course models.Course
	if err := config.DB.First(&course, courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "课程不存在"})
		return
	}

	// 检查课程是否属于当前教师
	if course.TeacherID != teacherID.(int) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权删除此课程"})
		return
	}

	// 先删除所有选课记录
	config.DB.Where("course_id = ?", courseID).Delete(&models.Enrollment{})

	// 删除课程
	if err := config.DB.Delete(&course).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "删除成功",
	})
}

// GetCourseStudents 获取课程的选课学生列表
// GET /api/teacher/courses/:id/students/
// 返回指定课程的所有选课学生信息
func GetCourseStudents(c *gin.Context) {
	// 从URL参数中获取课程ID
	courseID := c.Param("id")

	// 获取当前教师ID
	session := sessions.Default(c)
	teacherID := session.Get("user_id")

	// 查找课程并验证权限
	var course models.Course
	if err := config.DB.First(&course, courseID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "课程不存在"})
		return
	}

	// 检查课程是否属于当前教师
	if course.TeacherID != teacherID.(int) {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权查看此课程"})
		return
	}

	// 查询所有选课记录
	var enrollments []models.Enrollment
	config.DB.Where("course_id = ?", courseID).Find(&enrollments)

	// 构建学生列表
	var students []gin.H
	for _, enrollment := range enrollments {
		// 查询学生信息
		var student models.Student
		if err := config.DB.First(&student, enrollment.StudentID).Error; err != nil {
			continue // 学生不存在，跳过
		}

		students = append(students, gin.H{
			"id":          student.ID,
			"username":    student.Username,
			"email":       student.Email,
			"enrolled_at": enrollment.EnrolledAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"course": gin.H{
			"id":   course.ID,
			"name": course.Name,
		},
		"students": students,
		"total":    len(students), // 总人数
	})
}
