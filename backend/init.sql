-- MySQL初始化脚本
-- 数据库：course_system
-- ==========================================================================
-- 性能优化配置
-- ==========================================================================

-- 禁用外键检查（提高批量插入性能）
SET FOREIGN_KEY_CHECKS = 0;

-- 禁用唯一性检查（提高批量插入性能）
SET UNIQUE_CHECKS = 0;

-- 关闭自动提交（使用事务批量提交）
SET AUTOCOMMIT = 0;

-- 创建数据库
CREATE DATABASE IF NOT EXISTS `course_system` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE `course_system`;

-- 删除旧表（如果存在）
-- 注意：按照外键依赖关系的逆序删除
DROP TABLE IF EXISTS `course_schedules`;
DROP TABLE IF EXISTS `enrollments`;
DROP TABLE IF EXISTS `courses`;
DROP TABLE IF EXISTS `students`;
DROP TABLE IF EXISTS `teachers`;

-- 学生表
CREATE TABLE `students` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `username` VARCHAR(100) NOT NULL UNIQUE COMMENT '用户名',
  `password` VARCHAR(255) NOT NULL COMMENT '密码',
  `email` VARCHAR(255) NOT NULL UNIQUE COMMENT '邮箱',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  INDEX `idx_username` (`username`),
  INDEX `idx_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='学生表';

-- 教师表
CREATE TABLE `teachers` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `username` VARCHAR(100) NOT NULL UNIQUE COMMENT '用户名',
  `password` VARCHAR(255) NOT NULL COMMENT '密码',
  `email` VARCHAR(255) NOT NULL UNIQUE COMMENT '邮箱',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  INDEX `idx_username` (`username`),
  INDEX `idx_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='教师表';

-- 课程表（含乐观锁和已选人数字段）
CREATE TABLE `courses` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
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
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `student_id` INT NOT NULL COMMENT '学生ID（应用层关联）',
  `course_id` INT NOT NULL COMMENT '课程ID（应用层关联）',
  `enrolled_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '选课时间',
  INDEX `idx_student_course` (`student_id`, `course_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='选课记录表';

-- 课程时间表（用于选课冲突检测）
CREATE TABLE `course_schedules` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `course_id` INT NOT NULL COMMENT '课程ID（应用层关联）',
  `day_of_week` TINYINT NOT NULL COMMENT '星期几（1-7，1=周一）',
  `start_time` TIME NOT NULL COMMENT '开始时间',
  `end_time` TIME NOT NULL COMMENT '结束时间',
  `classroom` VARCHAR(100) DEFAULT '' COMMENT '教室',
  INDEX `idx_course_id` (`course_id`),
  INDEX `idx_day_time` (`day_of_week`, `start_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='课程时间表';

-- ==========================================================================
-- 插入测试数据
-- 所有测试账户的密码都是: password123（已使用bcrypt加密）
-- ==========================================================================

-- 开始事务（批量插入优化）
START TRANSACTION;

-- 测试学生账户（10条数据）
INSERT INTO `students` (`username`, `password`, `email`) VALUES
('student1', '$2a$10$Qa/NKW56HT.gL8rN4jxVv.eaME.V8tx5vzR3I/uUr9e0jil3Ed0Iu', 'student1@test.com'),
('student2', '$2a$10$w551xkwmtEPKy5dMPmcQKu28my/xgg.GE1sw.g3VvGH1d/UUxWWt6', 'student2@test.com'),
('student3', '$2a$10$Qa/NKW56HT.gL8rN4jxVv.eaME.V8tx5vzR3I/uUr9e0jil3Ed0Iu', 'student3@test.com'),
('student4', '$2a$10$Qa/NKW56HT.gL8rN4jxVv.eaME.V8tx5vzR3I/uUr9e0jil3Ed0Iu', 'student4@test.com'),
('student5', '$2a$10$Qa/NKW56HT.gL8rN4jxVv.eaME.V8tx5vzR3I/uUr9e0jil3Ed0Iu', 'student5@test.com'),
('student6', '$2a$10$Qa/NKW56HT.gL8rN4jxVv.eaME.V8tx5vzR3I/uUr9e0jil3Ed0Iu', 'student6@test.com'),
('student7', '$2a$10$Qa/NKW56HT.gL8rN4jxVv.eaME.V8tx5vzR3I/uUr9e0jil3Ed0Iu', 'student7@test.com'),
('student8', '$2a$10$Qa/NKW56HT.gL8rN4jxVv.eaME.V8tx5vzR3I/uUr9e0jil3Ed0Iu', 'student8@test.com'),
('student9', '$2a$10$Qa/NKW56HT.gL8rN4jxVv.eaME.V8tx5vzR3I/uUr9e0jil3Ed0Iu', 'student9@test.com'),
('student10', '$2a$10$Qa/NKW56HT.gL8rN4jxVv.eaME.V8tx5vzR3I/uUr9e0jil3Ed0Iu', 'student10@test.com');

-- 测试教师账户（10条数据）
INSERT INTO `teachers` (`username`, `password`, `email`) VALUES
('teacher1', '$2a$10$ASbr94uFmKwWS/njpgPbZ.YeMCnz7RpGGbuLK7.kVBKVVYkfWxrFS', 'teacher1@test.com'),
('teacher2', '$2a$10$V3KFAQOV0iMkpMp.CslRX.FycAn3RyCmUiM.cWqwJbjGtc42gYScO', 'teacher2@test.com'),
('teacher3', '$2a$10$ASbr94uFmKwWS/njpgPbZ.YeMCnz7RpGGbuLK7.kVBKVVYkfWxrFS', 'teacher3@test.com'),
('teacher4', '$2a$10$ASbr94uFmKwWS/njpgPbZ.YeMCnz7RpGGbuLK7.kVBKVVYkfWxrFS', 'teacher4@test.com'),
('teacher5', '$2a$10$ASbr94uFmKwWS/njpgPbZ.YeMCnz7RpGGbuLK7.kVBKVVYkfWxrFS', 'teacher5@test.com'),
('teacher6', '$2a$10$ASbr94uFmKwWS/njpgPbZ.YeMCnz7RpGGbuLK7.kVBKVVYkfWxrFS', 'teacher6@test.com'),
('teacher7', '$2a$10$ASbr94uFmKwWS/njpgPbZ.YeMCnz7RpGGbuLK7.kVBKVVYkfWxrFS', 'teacher7@test.com'),
('teacher8', '$2a$10$ASbr94uFmKwWS/njpgPbZ.YeMCnz7RpGGbuLK7.kVBKVVYkfWxrFS', 'teacher8@test.com'),
('teacher9', '$2a$10$ASbr94uFmKwWS/njpgPbZ.YeMCnz7RpGGbuLK7.kVBKVVYkfWxrFS', 'teacher9@test.com'),
('teacher10', '$2a$10$ASbr94uFmKwWS/njpgPbZ.YeMCnz7RpGGbuLK7.kVBKVVYkfWxrFS', 'teacher10@test.com');

-- 测试课程（10条数据，注意：enrolled字段会自动初始化为0，后面通过选课记录更新）
INSERT INTO `courses` (`name`, `description`, `teacher_id`, `capacity`, `enrolled`) VALUES
('Golang高级编程', 'Go语言并发编程与性能优化', 1, 30, 0),
('微服务架构设计', 'Spring Cloud微服务实战', 1, 40, 0),
('MySQL性能调优', 'MySQL从入门到精通', 2, 25, 0),
('分布式系统设计', '深入理解分布式系统原理与实践', 3, 35, 0),
('数据结构与算法', '算法分析与设计', 4, 50, 0),
('云原生架构', 'Kubernetes与Docker实战', 5, 30, 0),
('前端开发实战', 'React与Vue.js开发', 6, 45, 0),
('人工智能基础', '机器学习与深度学习入门', 7, 40, 0),
('网络安全', '网络安全原理与实践', 8, 35, 0),
('软件工程', '软件开发生命周期管理', 9, 50, 0);

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

-- 更新课程的enrolled字段（根据实际选课记录统计）
UPDATE `courses` SET `enrolled` = (
  SELECT COUNT(*) FROM `enrollments` WHERE `course_id` = `courses`.`id`
);

-- 课程时间安排（10条数据，用于测试选课冲突检测）
INSERT INTO `course_schedules` (`course_id`, `day_of_week`, `start_time`, `end_time`, `classroom`) VALUES
-- Golang高级编程（课程ID=1）：周一 08:00-10:00
(1, 1, '08:00:00', '10:00:00', 'A101'),
-- 微服务架构设计（课程ID=2）：周二 10:00-12:00
(2, 2, '10:00:00', '12:00:00', 'B201'),
-- MySQL性能调优（课程ID=3）：周三 14:00-16:00
(3, 3, '14:00:00', '16:00:00', 'C301'),
-- 分布式系统设计（课程ID=4）：周四 08:00-10:00
(4, 4, '08:00:00', '10:00:00', 'D401'),
-- 数据结构与算法（课程ID=5）：周五 10:00-12:00
(5, 5, '10:00:00', '12:00:00', 'A102'),
-- 云原生架构（课程ID=6）：周一 14:00-16:00
(6, 1, '14:00:00', '16:00:00', 'B202'),
-- 前端开发实战（课程ID=7）：周二 14:00-16:00
(7, 2, '14:00:00', '16:00:00', 'C302'),
-- 人工智能基础（课程ID=8）：周三 08:00-10:00
(8, 3, '08:00:00', '10:00:00', 'D402'),
-- 网络安全（课程ID=9）：周四 14:00-16:00
(9, 4, '14:00:00', '16:00:00', 'A103'),
-- 软件工程（课程ID=10）：周五 08:00-10:00
(10, 5, '08:00:00', '10:00:00', 'B203');

-- ==========================================================================
-- 提交事务并恢复配置
-- ==========================================================================

-- 提交所有更改
COMMIT;

-- 恢复外键检查
SET FOREIGN_KEY_CHECKS = 1;

-- 恢复唯一性检查
SET UNIQUE_CHECKS = 1;

-- 恢复自动提交
SET AUTOCOMMIT = 1;
