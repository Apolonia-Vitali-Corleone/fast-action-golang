-- ============================================================================
-- 选课系统数据库初始化脚本
-- 数据库名称: course_system
-- 字符集: utf8mb4
-- 排序规则: utf8mb4_unicode_ci
-- ============================================================================

-- 性能优化配置
SET FOREIGN_KEY_CHECKS = 0;
SET UNIQUE_CHECKS = 0;
SET AUTOCOMMIT = 0;

-- 创建数据库
CREATE DATABASE IF NOT EXISTS `course_system` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE `course_system`;

-- ============================================================================
-- 删除旧表（按依赖关系逆序删除）
-- ============================================================================
DROP TABLE IF EXISTS `captcha_codes`;
DROP TABLE IF EXISTS `sms_codes`;
DROP TABLE IF EXISTS `course_schedules`;
DROP TABLE IF EXISTS `enrollments`;
DROP TABLE IF EXISTS `courses`;
DROP TABLE IF EXISTS `students`;
DROP TABLE IF EXISTS `teachers`;

-- ============================================================================
-- 创建表
-- ============================================================================

-- 学生表
CREATE TABLE `students` (
    `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '主键，自增',
    `username` VARCHAR(100) NOT NULL UNIQUE COMMENT '用户名，唯一索引',
    `password` VARCHAR(255) NOT NULL COMMENT '密码（bcrypt加密）',
    `phone` VARCHAR(20) NOT NULL UNIQUE COMMENT '手机号，唯一索引',
    `email` VARCHAR(255) NOT NULL UNIQUE COMMENT '邮箱，唯一索引',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX `idx_username` (`username`),
    INDEX `idx_phone` (`phone`),
    INDEX `idx_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='学生表';

-- 教师表
CREATE TABLE `teachers` (
    `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '主键，自增',
    `username` VARCHAR(100) NOT NULL UNIQUE COMMENT '用户名，唯一索引',
    `password` VARCHAR(255) NOT NULL COMMENT '密码（bcrypt加密）',
    `email` VARCHAR(255) NOT NULL UNIQUE COMMENT '邮箱，唯一索引',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX `idx_username` (`username`),
    INDEX `idx_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='教师表';

-- 课程表（含乐观锁和已选人数字段）
CREATE TABLE `courses` (
    `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '主键，自增',
    `name` VARCHAR(200) NOT NULL COMMENT '课程名称',
    `description` TEXT COMMENT '课程描述',
    `teacher_id` INT NOT NULL COMMENT '教师ID（应用层关联）',
    `capacity` INT NOT NULL DEFAULT 50 COMMENT '课程容量',
    `enrolled` INT NOT NULL DEFAULT 0 COMMENT '已选人数（用于快速查询，避免COUNT）',
    `version` INT NOT NULL DEFAULT 0 COMMENT '乐观锁版本号（每次更新+1，防止并发冲突）',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX `idx_teacher_id` (`teacher_id`),
    INDEX `idx_enrolled` (`enrolled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='课程表';

-- 选课记录表
CREATE TABLE `enrollments` (
    `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '主键，自增',
    `student_id` INT NOT NULL COMMENT '学生ID（应用层关联）',
    `course_id` INT NOT NULL COMMENT '课程ID（应用层关联）',
    `enrolled_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '选课时间',
    UNIQUE INDEX `idx_student_course` (`student_id`, `course_id`) COMMENT '学生-课程联合唯一索引，防止重复选课'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='选课记录表';

-- 课程时间表（用于选课冲突检测）
-- 使用节次制：1=上午第一节, 2=上午第二节, 3=下午第一节, 4=下午第二节
CREATE TABLE `course_schedules` (
    `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '主键，自增',
    `course_id` INT NOT NULL COMMENT '课程ID（应用层关联）',
    `day_of_week` TINYINT NOT NULL COMMENT '星期几（1-5，1=周一，5=周五）',
    `time_slot` TINYINT NOT NULL COMMENT '节次（1-4: 1=上午一节, 2=上午二节, 3=下午一节, 4=下午二节）',
    `start_week` TINYINT NOT NULL COMMENT '开始周次（如第1周）',
    `end_week` TINYINT NOT NULL COMMENT '结束周次（如第16周）',
    `classroom` VARCHAR(100) DEFAULT '' COMMENT '教室',
    INDEX `idx_course_id` (`course_id`),
    INDEX `idx_day_slot` (`day_of_week`, `time_slot`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='课程时间表';

-- 短信验证码表
CREATE TABLE `sms_codes` (
    `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '主键，自增',
    `phone` VARCHAR(20) NOT NULL COMMENT '手机号',
    `code` VARCHAR(10) NOT NULL COMMENT '验证码',
    `purpose` VARCHAR(20) NOT NULL COMMENT '用途：register(注册)、login(登录)',
    `used` BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否已使用',
    `expires_at` DATETIME NOT NULL COMMENT '过期时间',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX `idx_phone` (`phone`),
    INDEX `idx_expires_at` (`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='短信验证码表';

-- 图形验证码表
CREATE TABLE `captcha_codes` (
    `id` INT AUTO_INCREMENT PRIMARY KEY COMMENT '主键，自增',
    `captcha_id` VARCHAR(100) NOT NULL UNIQUE COMMENT '验证码ID，唯一索引',
    `code` VARCHAR(10) NOT NULL COMMENT '验证码答案',
    `used` BOOLEAN NOT NULL DEFAULT FALSE COMMENT '是否已使用',
    `expires_at` DATETIME NOT NULL COMMENT '过期时间',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX `idx_expires_at` (`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='图形验证码表';

-- ============================================================================
-- 插入测试数据
-- 所有测试账户的密码都是: password123（已使用bcrypt加密）
-- bcrypt哈希: $2a$10$Qa/NKW56HT.gL8rN4jxVv.eaME.V8tx5vzR3I/uUr9e0jil3Ed0Iu
-- ============================================================================

START TRANSACTION;

-- 测试学生账户（10条数据）
INSERT INTO `students` (`username`, `password`, `phone`, `email`) VALUES
('student1', '$2a$10$Qa/NKW56HT.gL8rN4jxVv.eaME.V8tx5vzR3I/uUr9e0jil3Ed0Iu', '13800001001', 'student1@test.com'),
('student2', '$2a$10$Qa/NKW56HT.gL8rN4jxVv.eaME.V8tx5vzR3I/uUr9e0jil3Ed0Iu', '13800001002', 'student2@test.com'),
('student3', '$2a$10$Qa/NKW56HT.gL8rN4jxVv.eaME.V8tx5vzR3I/uUr9e0jil3Ed0Iu', '13800001003', 'student3@test.com'),
('student4', '$2a$10$Qa/NKW56HT.gL8rN4jxVv.eaME.V8tx5vzR3I/uUr9e0jil3Ed0Iu', '13800001004', 'student4@test.com'),
('student5', '$2a$10$Qa/NKW56HT.gL8rN4jxVv.eaME.V8tx5vzR3I/uUr9e0jil3Ed0Iu', '13800001005', 'student5@test.com'),
('student6', '$2a$10$Qa/NKW56HT.gL8rN4jxVv.eaME.V8tx5vzR3I/uUr9e0jil3Ed0Iu', '13800001006', 'student6@test.com'),
('student7', '$2a$10$Qa/NKW56HT.gL8rN4jxVv.eaME.V8tx5vzR3I/uUr9e0jil3Ed0Iu', '13800001007', 'student7@test.com'),
('student8', '$2a$10$Qa/NKW56HT.gL8rN4jxVv.eaME.V8tx5vzR3I/uUr9e0jil3Ed0Iu', '13800001008', 'student8@test.com'),
('student9', '$2a$10$Qa/NKW56HT.gL8rN4jxVv.eaME.V8tx5vzR3I/uUr9e0jil3Ed0Iu', '13800001009', 'student9@test.com'),
('student10', '$2a$10$Qa/NKW56HT.gL8rN4jxVv.eaME.V8tx5vzR3I/uUr9e0jil3Ed0Iu', '13800001010', 'student10@test.com');

-- 测试教师账户（10条数据）
INSERT INTO `teachers` (`username`, `password`, `email`) VALUES
('teacher1', '$2a$10$Qa/NKW56HT.gL8rN4jxVv.eaME.V8tx5vzR3I/uUr9e0jil3Ed0Iu', 'teacher1@test.com'),
('teacher2', '$2a$10$Qa/NKW56HT.gL8rN4jxVv.eaME.V8tx5vzR3I/uUr9e0jil3Ed0Iu', 'teacher2@test.com'),
('teacher3', '$2a$10$Qa/NKW56HT.gL8rN4jxVv.eaME.V8tx5vzR3I/uUr9e0jil3Ed0Iu', 'teacher3@test.com'),
('teacher4', '$2a$10$Qa/NKW56HT.gL8rN4jxVv.eaME.V8tx5vzR3I/uUr9e0jil3Ed0Iu', 'teacher4@test.com'),
('teacher5', '$2a$10$Qa/NKW56HT.gL8rN4jxVv.eaME.V8tx5vzR3I/uUr9e0jil3Ed0Iu', 'teacher5@test.com'),
('teacher6', '$2a$10$Qa/NKW56HT.gL8rN4jxVv.eaME.V8tx5vzR3I/uUr9e0jil3Ed0Iu', 'teacher6@test.com'),
('teacher7', '$2a$10$Qa/NKW56HT.gL8rN4jxVv.eaME.V8tx5vzR3I/uUr9e0jil3Ed0Iu', 'teacher7@test.com'),
('teacher8', '$2a$10$Qa/NKW56HT.gL8rN4jxVv.eaME.V8tx5vzR3I/uUr9e0jil3Ed0Iu', 'teacher8@test.com'),
('teacher9', '$2a$10$Qa/NKW56HT.gL8rN4jxVv.eaME.V8tx5vzR3I/uUr9e0jil3Ed0Iu', 'teacher9@test.com'),
('teacher10', '$2a$10$Qa/NKW56HT.gL8rN4jxVv.eaME.V8tx5vzR3I/uUr9e0jil3Ed0Iu', 'teacher10@test.com');

-- 测试课程（10条数据）
INSERT INTO `courses` (`name`, `description`, `teacher_id`, `capacity`, `enrolled`) VALUES
('Golang高级编程', 'Go语言并发编程与性能优化，深入理解goroutine、channel、context等核心概念', 1, 30, 2),
('微服务架构设计', 'Spring Cloud微服务实战，构建高可用分布式系统', 1, 40, 1),
('MySQL性能调优', 'MySQL从入门到精通，索引优化、查询优化、分库分表', 2, 25, 1),
('分布式系统设计', '深入理解分布式系统原理与实践，CAP定理、一致性协议、分布式事务', 3, 35, 1),
('数据结构与算法', '算法分析与设计，LeetCode刷题技巧与面试攻略', 4, 50, 1),
('云原生架构', 'Kubernetes与Docker实战，容器化部署与微服务治理', 5, 30, 1),
('前端开发实战', 'React与Vue.js开发，现代前端工程化实践', 6, 45, 1),
('人工智能基础', '机器学习与深度学习入门，TensorFlow与PyTorch实战', 7, 40, 1),
('网络安全', '网络安全原理与实践，Web安全、渗透测试、CTF竞赛', 8, 35, 1),
('软件工程', '软件开发生命周期管理，敏捷开发、DevOps、持续集成', 9, 50, 0);

-- 测试选课记录（10条数据）
INSERT INTO `enrollments` (`student_id`, `course_id`) VALUES
(1, 1),  -- student1 选了 Golang高级编程
(1, 2),  -- student1 选了 微服务架构设计
(2, 1),  -- student2 选了 Golang高级编程
(3, 3),  -- student3 选了 MySQL性能调优
(4, 4),  -- student4 选了 分布式系统设计
(5, 5),  -- student5 选了 数据结构与算法
(6, 6),  -- student6 选了 云原生架构
(7, 7),  -- student7 选了 前端开发实战
(8, 8),  -- student8 选了 人工智能基础
(9, 9);  -- student9 选了 网络安全

-- 课程时间安排（使用节次制）
-- 节次说明：1=上午第一节(08:00-10:00), 2=上午第二节(10:00-12:00), 3=下午第一节(14:00-16:00), 4=下午第二节(16:00-18:00)
INSERT INTO `course_schedules` (`course_id`, `day_of_week`, `time_slot`, `start_week`, `end_week`, `classroom`) VALUES
-- Golang高级编程（课程ID=1）：周一上午第一节，第1-16周
(1, 1, 1, 1, 16, 'A101'),
-- 微服务架构设计（课程ID=2）：周二上午第二节，第1-16周
(2, 2, 2, 1, 16, 'B201'),
-- MySQL性能调优（课程ID=3）：周三下午第一节，第1-16周
(3, 3, 3, 1, 16, 'C301'),
-- 分布式系统设计（课程ID=4）：周四上午第一节，第1-16周
(4, 4, 1, 1, 16, 'D401'),
-- 数据结构与算法（课程ID=5）：周五上午第二节，第1-16周
(5, 5, 2, 1, 16, 'A102'),
-- 云原生架构（课程ID=6）：周一下午第一节，第1-16周
(6, 1, 3, 1, 16, 'B202'),
-- 前端开发实战（课程ID=7）：周二下午第一节，第1-16周
(7, 2, 3, 1, 16, 'C302'),
-- 人工智能基础（课程ID=8）：周三上午第一节，第1-16周
(8, 3, 1, 1, 16, 'D402'),
-- 网络安全（课程ID=9）：周四下午第一节，第1-16周
(9, 4, 3, 1, 16, 'A103'),
-- 软件工程（课程ID=10）：周五上午第一节，第1-16周
(10, 5, 1, 1, 16, 'B203');

-- 提交事务
COMMIT;

-- 恢复配置
SET FOREIGN_KEY_CHECKS = 1;
SET UNIQUE_CHECKS = 1;
SET AUTOCOMMIT = 1;

-- ============================================================================
-- 初始化完成
-- ============================================================================
-- 测试账户信息：
-- 学生账户: student1~student10 (密码: password123, 手机号: 13800001001~13800001010)
-- 教师账户: teacher1~teacher10 (密码: password123)
-- ============================================================================
