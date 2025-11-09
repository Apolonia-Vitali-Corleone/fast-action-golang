# 选课系统后端 - Golang + Gin + GORM

这是一个基于Golang的简单选课系统后端API，使用Gin框架和GORM ORM。

## 技术栈

- **Golang**: 编程语言
- **Gin**: Web框架
- **GORM**: ORM（对象关系映射）
- **MySQL**: 数据库
- **bcrypt**: 密码加密
- **Session**: 用户认证

## 项目结构

```
backend/
├── main.go              # 主程序入口，路由配置
├── go.mod              # Go模块依赖
├── config/
│   └── database.go     # 数据库连接配置
├── models/
│   └── models.go       # 数据模型定义
├── controllers/
│   ├── student.go      # 学生控制器
│   └── teacher.go      # 教师控制器
└── middleware/
    └── auth.go         # 认证中间件
```

## 快速开始

### 1. 前置要求

- Go 1.21 或更高版本
- MySQL 5.7 或更高版本
- 已执行项目根目录的 `init.sql` 创建数据库

### 2. 修改数据库配置

打开 `main.go` 文件，修改数据库配置：

```go
dbConfig := config.DBConfig{
    Host:     "localhost",     // 数据库地址
    Port:     "3306",          // 数据库端口
    User:     "root",          // 用户名
    Password: "password",      // 密码（修改为你的密码）
    DBName:   "course_system", // 数据库名
}
```

### 3. 安装依赖

```bash
cd backend
go mod tidy
```

### 4. 运行项目

```bash
go run main.go
```

服务器将在 `http://localhost:8000` 启动

## API接口文档

### 学生接口

#### 1. 学生注册
- **URL**: `POST /api/student/register/`
- **请求体**:
  ```json
  {
    "username": "student1",
    "password": "123456",
    "email": "student1@test.com"
  }
  ```

#### 2. 学生登录
- **URL**: `POST /api/student/login/`
- **请求体**:
  ```json
  {
    "username": "student1",
    "password": "123456"
  }
  ```

#### 3. 获取所有课程
- **URL**: `GET /api/student/courses/`
- **需要登录**: 是
- **权限**: 学生

#### 4. 获取我的课程
- **URL**: `GET /api/student/my-courses/`
- **需要登录**: 是
- **权限**: 学生

#### 5. 选课
- **URL**: `POST /api/student/enroll/`
- **需要登录**: 是
- **权限**: 学生
- **请求体**:
  ```json
  {
    "course_id": 1
  }
  ```

#### 6. 退课
- **URL**: `POST /api/student/drop/`
- **需要登录**: 是
- **权限**: 学生
- **请求体**:
  ```json
  {
    "course_id": 1
  }
  ```

### 教师接口

#### 1. 教师注册
- **URL**: `POST /api/teacher/register/`
- **请求体**:
  ```json
  {
    "username": "teacher1",
    "password": "123456",
    "email": "teacher1@test.com"
  }
  ```

#### 2. 教师登录
- **URL**: `POST /api/teacher/login/`
- **请求体**:
  ```json
  {
    "username": "teacher1",
    "password": "123456"
  }
  ```

#### 3. 获取我的课程
- **URL**: `GET /api/teacher/courses/`
- **需要登录**: 是
- **权限**: 教师

#### 4. 创建课程
- **URL**: `POST /api/teacher/courses/create/`
- **需要登录**: 是
- **权限**: 教师
- **请求体**:
  ```json
  {
    "name": "Python编程",
    "description": "Python从入门到精通",
    "capacity": 30
  }
  ```

#### 5. 删除课程
- **URL**: `DELETE /api/teacher/courses/:id/delete/`
- **需要登录**: 是
- **权限**: 教师

#### 6. 获取课程学生列表
- **URL**: `GET /api/teacher/courses/:id/students/`
- **需要登录**: 是
- **权限**: 教师

### 通用接口

#### 1. 获取当前用户
- **URL**: `GET /api/current-user/`
- **需要登录**: 是

#### 2. 退出登录
- **URL**: `POST /api/logout/`

## 数据库表结构

### students（学生表）
- `id`: 主键
- `username`: 用户名（唯一）
- `password`: 密码（bcrypt加密）
- `email`: 邮箱（唯一）
- `created_at`: 创建时间

### teachers（教师表）
- `id`: 主键
- `username`: 用户名（唯一）
- `password`: 密码（bcrypt加密）
- `email`: 邮箱（唯一）
- `created_at`: 创建时间

### courses（课程表）
- `id`: 主键
- `name`: 课程名称
- `description`: 课程描述
- `teacher_id`: 教师ID
- `capacity`: 课程容量
- `created_at`: 创建时间

### enrollments（选课记录表）
- `id`: 主键
- `student_id`: 学生ID
- `course_id`: 课程ID
- `enrolled_at`: 选课时间

## 注意事项

1. **密码安全**: 使用bcrypt加密，原SQL中的测试数据密码需要重新注册或修改
2. **Session密钥**: `main.go`中的session密钥需要在生产环境修改
3. **CORS配置**: 默认只允许`http://localhost:5173`访问，如需修改请编辑`main.go`
4. **数据库连接**: 确保MySQL服务已启动并正确配置

## 测试账号

由于密码加密方式不同，建议使用前端注册新账号，或者手动修改数据库密码为bcrypt加密后的值。

## 开发说明

这是一个极简化的教学项目，没有考虑以下内容：
- 完善的错误处理
- 日志记录
- 单元测试
- 请求验证
- 性能优化
- 生产环境配置

如需用于生产环境，请添加以上功能。
