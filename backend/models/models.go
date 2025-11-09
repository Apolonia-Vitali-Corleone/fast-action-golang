package models

import (
	"time"
)

// Student 学生表模型
// 对应数据库中的students表
type Student struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`           // 主键，自增
	Username  string    `gorm:"type:varchar(100);uniqueIndex" json:"username"` // 用户名，唯一索引
	Password  string    `gorm:"type:varchar(255)" json:"-"`                    // 密码，json序列化时忽略（安全）
	Email     string    `gorm:"type:varchar(255);uniqueIndex" json:"email"`    // 邮箱，唯一索引
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`              // 创建时间，自动填充
}

// TableName 指定表名
func (Student) TableName() string {
	return "students"
}

// Teacher 教师表模型
// 对应数据库中的teachers表
type Teacher struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`           // 主键，自增
	Username  string    `gorm:"type:varchar(100);uniqueIndex" json:"username"` // 用户名，唯一索引
	Password  string    `gorm:"type:varchar(255)" json:"-"`                    // 密码，json序列化时忽略（安全）
	Email     string    `gorm:"type:varchar(255);uniqueIndex" json:"email"`    // 邮箱，唯一索引
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`              // 创建时间，自动填充
}

// TableName 指定表名
func (Teacher) TableName() string {
	return "teachers"
}

// Course 课程表模型
// 对应数据库中的courses表
type Course struct {
	ID          int       `gorm:"primaryKey;autoIncrement" json:"id"`        // 主键，自增
	Name        string    `gorm:"type:varchar(200)" json:"name"`             // 课程名称
	Description string    `gorm:"type:text" json:"description"`              // 课程描述
	TeacherID   int       `gorm:"index" json:"teacher_id"`                   // 教师ID，建立索引
	Capacity    int       `gorm:"default:50" json:"capacity"`                // 课程容量，默认50
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`          // 创建时间，自动填充
}

// TableName 指定表名
func (Course) TableName() string {
	return "courses"
}

// Enrollment 选课记录表模型
// 对应数据库中的enrollments表
type Enrollment struct {
	ID         int       `gorm:"primaryKey;autoIncrement" json:"id"`  // 主键，自增
	StudentID  int       `gorm:"index:idx_student_course" json:"student_id"` // 学生ID，联合索引的一部分
	CourseID   int       `gorm:"index:idx_student_course" json:"course_id"`  // 课程ID，联合索引的一部分
	EnrolledAt time.Time `gorm:"autoCreateTime" json:"enrolled_at"`   // 选课时间，自动填充
}

// TableName 指定表名
func (Enrollment) TableName() string {
	return "enrollments"
}
