# 选课系统 (Course Selection System)

一个基于 Vue 3 + Golang 的全栈选课系统，支持学生选课、教师管理课程等功能。

## 技术栈

### 前端
- **Vue 3**: 渐进式 JavaScript 框架
- **Vite**: 前端构建工具
- **Axios**: HTTP 客户端
- **Vue Router**: 路由管理

### 后端
- **Golang**: 编程语言
- **Gin**: Web 框架
- **GORM**: ORM（对象关系映射）
- **MySQL**: 数据库
- **bcrypt**: 密码加密
- **Session**: 用户认证

## 项目结构

```
fast-action-golang/
├── frontend/              # Vue 3 前端项目
│   ├── src/
│   │   ├── components/   # 组件
│   │   ├── views/        # 页面视图
│   │   ├── App.vue       # 根组件
│   │   └── main.js       # 入口文件
│   ├── package.json
│   └── vite.config.js
├── backend/              # Golang 后端项目
│   ├── config/           # 配置
│   │   └── database.go   # 数据库连接
│   ├── models/           # 数据模型
│   │   └── models.go
│   ├── controllers/      # 控制器
│   │   ├── student.go    # 学生控制器
│   │   └── teacher.go    # 教师控制器
│   ├── middleware/       # 中间件
│   │   └── auth.go       # 认证中间件
│   ├── main.go           # 主程序入口
│   └── go.mod            # Go 模块依赖
├── init.sql              # 数据库初始化脚本
└── README.md
```

## 功能特性

### 学生功能
- 用户注册/登录
- 浏览所有课程
- 选课/退课
- 查看我的课程

### 教师功能
- 用户注册/登录
- 创建课程
- 查看我的课程
- 删除课程
- 查看课程学生名单

## 快速开始

### 前置要求

- **Node.js**: 16.0 或更高版本
- **Go**: 1.21 或更高版本
- **MySQL**: 5.7 或更高版本

### 1. 数据库设置

首先创建数据库并导入初始数据：

```bash
# 登录 MySQL
mysql -u root -p

# 创建数据库（或直接导入 init.sql）
source init.sql
```

`init.sql` 文件会自动创建 `course_system` 数据库和必要的表结构。

### 2. 后端设置

```bash
# 进入后端目录
cd backend

# 安装依赖
go mod tidy

# 修改 main.go 中的数据库配置
# 将数据库密码改为你的 MySQL 密码

# 运行后端服务
go run main.go
```

后端服务将在 `http://localhost:8000` 启动

### 3. 前端设置

```bash
# 进入前端目录
cd frontend

# 安装依赖
npm install

# 运行开发服务器
npm run dev
```

前端服务将在 `http://localhost:5173` 启动

### 4. 访问应用

打开浏览器访问 `http://localhost:5173`

## API 接口文档

### 学生接口

| 方法 | 路径 | 说明 | 需要登录 |
|------|------|------|----------|
| POST | `/api/student/register/` | 学生注册 | 否 |
| POST | `/api/student/login/` | 学生登录 | 否 |
| GET | `/api/student/courses/` | 获取所有课程 | 是 |
| GET | `/api/student/my-courses/` | 获取我的课程 | 是 |
| POST | `/api/student/enroll/` | 选课 | 是 |
| POST | `/api/student/drop/` | 退课 | 是 |

### 教师接口

| 方法 | 路径 | 说明 | 需要登录 |
|------|------|------|----------|
| POST | `/api/teacher/register/` | 教师注册 | 否 |
| POST | `/api/teacher/login/` | 教师登录 | 否 |
| GET | `/api/teacher/courses/` | 获取我的课程 | 是 |
| POST | `/api/teacher/courses/create/` | 创建课程 | 是 |
| DELETE | `/api/teacher/courses/:id/delete/` | 删除课程 | 是 |
| GET | `/api/teacher/courses/:id/students/` | 获取课程学生列表 | 是 |

### 通用接口

| 方法 | 路径 | 说明 | 需要登录 |
|------|------|------|----------|
| GET | `/api/current-user/` | 获取当前用户信息 | 是 |
| POST | `/api/logout/` | 退出登录 | 否 |

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

## 配置说明

### 后端配置

在 `backend/main.go` 中修改数据库配置：

```go
dbConfig := config.DBConfig{
    Host:     "localhost",     // 数据库地址
    Port:     "3306",          // 数据库端口
    User:     "root",          // 用户名
    Password: "your_password", // 修改为你的密码
    DBName:   "course_system", // 数据库名
}
```

### CORS 配置

默认允许 `http://localhost:5173` 访问，如需修改请编辑 `backend/main.go` 中的 CORS 配置。

## 注意事项

1. **密码安全**: 所有密码使用 bcrypt 加密存储
2. **Session 密钥**: 生产环境请修改 `main.go` 中的 session 密钥
3. **数据库连接**: 确保 MySQL 服务已启动并正确配置
4. **端口冲突**: 确保 8000 和 5173 端口未被占用

## 开发说明

这是一个教学演示项目，生产环境使用前建议添加：
- 完善的错误处理
- 日志记录系统
- 单元测试
- 请求参数验证
- 性能优化
- 安全加固

## License

MIT
