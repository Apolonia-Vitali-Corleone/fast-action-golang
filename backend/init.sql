-- MySQL初始化脚本
-- 数据库：course_system

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

-- 测试学生账户
INSERT INTO `students` (`username`, `password`, `email`) VALUES
('student1', '$2a$10$Qa/NKW56HT.gL8rN4jxVv.eaME.V8tx5vzR3I/uUr9e0jil3Ed0Iu', 'student1@test.com'),
('student2', '$2a$10$w551xkwmtEPKy5dMPmcQKu28my/xgg.GE1sw.g3VvGH1d/UUxWWt6', 'student2@test.com');

-- 测试教师账户
INSERT INTO `teachers` (`username`, `password`, `email`) VALUES
('teacher1', '$2a$10$ASbr94uFmKwWS/njpgPbZ.YeMCnz7RpGGbuLK7.kVBKVVYkfWxrFS', 'teacher1@test.com'),
('teacher2', '$2a$10$V3KFAQOV0iMkpMp.CslRX.FycAn3RyCmUiM.cWqwJbjGtc42gYScO', 'teacher2@test.com');

-- 测试课程（注意：enrolled字段会自动初始化为0，后面通过选课记录更新）
INSERT INTO `courses` (`name`, `description`, `teacher_id`, `capacity`, `enrolled`) VALUES
('Golang高级编程', 'Go语言并发编程与性能优化', 1, 30, 0),
('微服务架构设计', 'Spring Cloud微服务实战', 1, 40, 0),
('MySQL性能调优', 'MySQL从入门到精通', 2, 25, 0);

-- 测试选课记录
INSERT INTO `enrollments` (`student_id`, `course_id`) VALUES
(1, 1),  -- student1 选了 Golang高级编程
(1, 2),  -- student1 选了 微服务架构设计
(2, 1);  -- student2 选了 Golang高级编程

-- 更新课程的enrolled字段（根据实际选课记录统计）
UPDATE `courses` SET `enrolled` = (
  SELECT COUNT(*) FROM `enrollments` WHERE `course_id` = `courses`.`id`
) WHERE `id` IN (1, 2, 3);

-- 课程时间安排（用于测试选课冲突检测）
INSERT INTO `course_schedules` (`course_id`, `day_of_week`, `start_time`, `end_time`, `classroom`) VALUES
-- Golang高级编程（课程ID=1）：周一 08:00-10:00，周三 14:00-16:00
(1, 1, '08:00:00', '10:00:00', 'A101'),
(1, 3, '14:00:00', '16:00:00', 'A101'),

-- 微服务架构设计（课程ID=2）：周二 10:00-12:00，周四 14:00-16:00
(2, 2, '10:00:00', '12:00:00', 'B201'),
(2, 4, '14:00:00', '16:00:00', 'B201'),

-- MySQL性能调优（课程ID=3）：周三 14:00-16:00（与课程1时间冲突！）
(3, 3, '14:00:00', '16:00:00', 'C301');
