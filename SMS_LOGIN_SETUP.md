# 学生手机号注册登录功能说明

## 功能概述

本次更新实现了学生手机号注册和登录功能，替换了原有的用户名+密码方式。

### 注册流程
1. 用户输入用户名（唯一）
2. 用户输入手机号（唯一）
3. 接收短信验证码
4. 验证通过后注册成功

### 登录流程
1. 用户输入手机号
2. 获取并输入图形验证码
3. 通过图形验证码后，系统发送短信验证码
4. 输入短信验证码完成登录

## 代码变更

### 1. 数据模型变更（models/models.go）
- **Student表**: 添加 `phone` 字段，移除 `password` 和 `email` 字段
- **新增 SMSCode 表**: 存储短信验证码
- **新增 CaptchaCode 表**: 存储图形验证码

### 2. 新增文件
- `config/sms.go` - 阿里云短信配置
- `utils/captcha.go` - 图形验证码生成和验证
- `utils/sms.go` - 短信发送和验证
- `controllers/verification.go` - 验证码相关API

### 3. 修改的API
- `POST /api/student/register/` - 改为使用 `{username, phone, sms_code}`
- `POST /api/student/login/` - 改为使用 `{phone, sms_code}`

### 4. 新增的API
- `GET /api/captcha/` - 获取图形验证码
- `POST /api/sms/send/` - 发送短信验证码

## 安装依赖

进入 backend 目录后运行：

```bash
cd backend

# 安装图形验证码库
go get github.com/mojocn/base64Captcha

# 安装UUID库
go get github.com/google/uuid

# 安装阿里云短信SDK
go get github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi

# 或者直接运行 go mod tidy
go mod tidy
```

## 数据库迁移

运行迁移脚本更新数据库结构：

```bash
mysql -u root -p course_system < migration_sms.sql
```

**注意**: 迁移脚本会删除现有的学生数据，因为表结构已改变。

## 配置阿里云短信服务

### 方式一：使用环境变量（推荐）

```bash
export ALIYUN_ACCESS_KEY_ID="你的AccessKey ID"
export ALIYUN_ACCESS_KEY_SECRET="你的AccessKey Secret"
export ALIYUN_SMS_SIGN_NAME="你的短信签名"
export ALIYUN_SMS_TEMPLATE_CODE="你的短信模板代码"
export ALIYUN_REGION_ID="cn-hangzhou"  # 可选，默认为 cn-hangzhou
```

### 方式二：开发模式（不发送短信）

如果没有配置阿里云密钥，系统会自动进入开发模式，验证码会在控制台打印，不会实际发送短信。
这对于本地开发和测试非常方便。

### 获取阿里云密钥

1. 登录阿里云控制台
2. 进入 "访问控制" -> "用户" -> 创建用户或使用现有用户
3. 创建 AccessKey
4. 开通短信服务并配置签名和模板

## API 使用示例

### 1. 获取图形验证码
```bash
curl http://localhost:8000/api/captcha/
```

响应:
```json
{
  "captcha_id": "uuid-string",
  "image": "data:image/png;base64,iVBORw0KG..."
}
```

### 2. 发送短信验证码（注册）
```bash
curl -X POST http://localhost:8000/api/sms/send/ \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "13800138000",
    "purpose": "register"
  }'
```

### 3. 发送短信验证码（登录，需要先验证图形验证码）
```bash
curl -X POST http://localhost:8000/api/sms/send/ \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "13800138000",
    "purpose": "login",
    "captcha_id": "从步骤1获取的captcha_id",
    "captcha_code": "1234"
  }'
```

### 4. 学生注册
```bash
curl -X POST http://localhost:8000/api/student/register/ \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "phone": "13800138000",
    "sms_code": "123456"
  }'
```

响应:
```json
{
  "message": "注册成功",
  "token": "jwt-token-string",
  "user": {
    "id": 1,
    "username": "testuser",
    "phone": "13800138000",
    "role": "student"
  }
}
```

### 5. 学生登录
```bash
curl -X POST http://localhost:8000/api/student/login/ \
  -H "Content-Type: application/json" \
  -d '{
    "phone": "13800138000",
    "sms_code": "123456"
  }'
```

## 验证码有效期

- **图形验证码**: 5分钟
- **短信验证码**: 5分钟
- **短信发送频率限制**: 同一手机号1分钟内只能发送一次

## 测试数据

迁移脚本会自动插入以下测试学生：

| 用户名 | 手机号 |
|--------|--------|
| test_student1 | 13800138001 |
| test_student2 | 13800138002 |
| test_student3 | 13800138003 |
| test_student4 | 13800138004 |
| test_student5 | 13800138005 |

在开发模式下（未配置阿里云密钥），你可以使用任意6位数字作为验证码进行测试，
验证码会在服务器控制台打印出来。

## 注意事项

1. **手机号格式**: 系统要求手机号为11位数字
2. **验证码安全**: 验证码使用后会自动标记为已使用，不可重复使用
3. **过期清理**: 建议定期调用 `CleanExpiredSMSCodes()` 和 `CleanExpiredCaptcha()` 清理过期验证码
4. **生产环境**: 生产环境务必配置正确的阿里云密钥，否则无法发送短信
5. **HTTPS**: 生产环境建议使用HTTPS来保护验证码传输安全
