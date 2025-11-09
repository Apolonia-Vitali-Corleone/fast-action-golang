-- MySQL初始化脚本
-- 数据库：course_system

-- 创建数据库
CREATE DATABASE IF NOT EXISTS `course_system` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE `course_system`;

-- 删除旧表（如果存在）
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

-- 课程表
CREATE TABLE `courses` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `name` VARCHAR(200) NOT NULL COMMENT '课程名称',
  `description` TEXT COMMENT '课程描述',
  `teacher_id` INT NOT NULL COMMENT '教师ID（应用层关联）',
  `capacity` INT NOT NULL DEFAULT 50 COMMENT '容量',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  INDEX `idx_teacher_id` (`teacher_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='课程表';

-- 选课记录表
CREATE TABLE `enrollments` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `student_id` INT NOT NULL COMMENT '学生ID（应用层关联）',
  `course_id` INT NOT NULL COMMENT '课程ID（应用层关联）',
  `enrolled_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '选课时间',
  INDEX `idx_student_course` (`student_id`, `course_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='选课记录表';

-- 插入测试数据（密码都是: password123，已经用Django的make_password加密）
-- 测试学生
INSERT INTO `students` (`username`, `password`, `email`) VALUES
('student1', 'pbkdf2_sha256$600000$yuxKmTcx5gEpqYHS8SzVUV$b2Pj8DkJj0yVYhvA2EGPrzRhVApLfTyeXqeZ7A9gWo0=', 'student1@test.com'),
('student2', 'pbkdf2_sha256$600000$aU0rTWrwRQYeY1K5ZNyeWc$cUKFdUunu7w317d0otqiiA/FD1s+GdGosHBhYy+jSjY=', 'student2@test.com');

-- 测试教师
INSERT INTO `teachers` (`username`, `password`, `email`) VALUES
('teacher1', 'pbkdf2_sha256$600000$yBWYx2XU7Fe77F5blN0ptT$6kZId4KCrTgRgDYGzUMnnrsO6pxrGPf+heinHXrlsoA=', 'teacher1@test.com'),
('teacher2', 'pbkdf2_sha256$600000$3DvDze07fSuTR0mYC66SEl$z7GkjWsZd31NQ3OZ3Ju1V0CVBCLHzXnKUcuyu6jiR+E=', 'teacher2@test.com');

-- 测试课程
INSERT INTO `courses` (`name`, `description`, `teacher_id`, `capacity`) VALUES
('Python编程', 'Python基础到进阶', 1, 30),
('Web开发', 'Django框架实战', 1, 40),
('数据库设计', 'MySQL从入门到精通', 2, 25);

-- 测试选课记录
INSERT INTO `enrollments` (`student_id`, `course_id`) VALUES
(1, 1),
(1, 2),
(2, 1);
