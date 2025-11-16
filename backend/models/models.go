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
	Phone     string    `gorm:"type:varchar(20);uniqueIndex" json:"phone"`     // 手机号，唯一索引
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
	Enrolled    int       `gorm:"default:0" json:"enrolled"`                 // 已选人数，用于快速查询
	Version     int       `gorm:"default:0" json:"version"`                  // 乐观锁版本号，每次更新+1
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

// CourseSchedule 课程时间表模型
// 用于记录课程的上课时间，支持选课时间冲突检测
// TimeSlot: 1=上午第一节, 2=上午第二节, 3=下午第一节, 4=下午第二节
type CourseSchedule struct {
	ID        int    `gorm:"primaryKey;autoIncrement" json:"id"`       // 主键，自增
	CourseID  int    `gorm:"index" json:"course_id"`                   // 课程ID，建立索引
	DayOfWeek int    `gorm:"type:tinyint" json:"day_of_week"`          // 星期几（1-5，1表示周一，5表示周五）
	TimeSlot  int    `gorm:"type:tinyint" json:"time_slot"`            // 节次（1-4: 1=上午一节, 2=上午二节, 3=下午一节, 4=下午二节）
	StartWeek int    `gorm:"type:tinyint" json:"start_week"`           // 开始周次（如第1周）
	EndWeek   int    `gorm:"type:tinyint" json:"end_week"`             // 结束周次（如第16周）
	Classroom string `gorm:"type:varchar(100)" json:"classroom"`       // 教室
}

// TableName 指定表名
func (CourseSchedule) TableName() string {
	return "course_schedules"
}

// SMSCode 短信验证码表模型
// 用于存储发送给用户的短信验证码
type SMSCode struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`          // 主键，自增
	Phone     string    `gorm:"type:varchar(20);index" json:"phone"`         // 手机号，建立索引
	Code      string    `gorm:"type:varchar(10)" json:"code"`                // 验证码
	Purpose   string    `gorm:"type:varchar(20)" json:"purpose"`             // 用途：register(注册)、login(登录)
	Used      bool      `gorm:"default:false" json:"used"`                   // 是否已使用
	ExpiresAt time.Time `gorm:"index" json:"expires_at"`                     // 过期时间，建立索引
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`            // 创建时间，自动填充
}

// TableName 指定表名
func (SMSCode) TableName() string {
	return "sms_codes"
}

// CaptchaCode 图形验证码表模型
// 用于存储生成的图形验证码
type CaptchaCode struct {
	ID        int       `gorm:"primaryKey;autoIncrement" json:"id"`       // 主键，自增
	CaptchaID string    `gorm:"type:varchar(100);uniqueIndex" json:"captcha_id"` // 验证码ID，唯一索引
	Code      string    `gorm:"type:varchar(10)" json:"code"`             // 验证码答案
	Used      bool      `gorm:"default:false" json:"used"`                // 是否已使用
	ExpiresAt time.Time `gorm:"index" json:"expires_at"`                  // 过期时间，建立索引
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`         // 创建时间，自动填充
}

// TableName 指定表名
func (CaptchaCode) TableName() string {
	return "captcha_codes"
}
