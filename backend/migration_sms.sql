-- MySQL迁移脚本
-- 功能：添加手机号登录和短信验证码功能
-- 日期：2024
-- ==========================================================================

USE `course_system`;

-- ==========================================================================
-- 第一步：修改 students 表，添加手机号字段，移除密码和邮箱字段
-- ==========================================================================

-- 添加手机号字段（临时允许NULL，后续会设置为NOT NULL）
ALTER TABLE `students` ADD COLUMN `phone` VARCHAR(20) NULL COMMENT '手机号' AFTER `username`;

-- 为phone字段创建唯一索引
CREATE UNIQUE INDEX `idx_phone` ON `students`(`phone`);

-- 删除password和email字段
ALTER TABLE `students` DROP COLUMN `password`;
ALTER TABLE `students` DROP COLUMN `email`;

-- ==========================================================================
-- 第二步：创建短信验证码表
-- ==========================================================================

CREATE TABLE IF NOT EXISTS `sms_codes` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `phone` VARCHAR(20) NOT NULL COMMENT '手机号',
    `code` VARCHAR(10) NOT NULL COMMENT '验证码',
    `purpose` VARCHAR(20) NOT NULL COMMENT '用途：register(注册)、login(登录)',
    `used` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否已使用',
    `expires_at` DATETIME NOT NULL COMMENT '过期时间',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    INDEX `idx_phone` (`phone`),
    INDEX `idx_expires_at` (`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='短信验证码表';

-- ==========================================================================
-- 第三步：创建图形验证码表
-- ==========================================================================

CREATE TABLE IF NOT EXISTS `captcha_codes` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `captcha_id` VARCHAR(100) NOT NULL UNIQUE COMMENT '验证码ID',
    `code` VARCHAR(10) NOT NULL COMMENT '验证码答案',
    `used` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '是否已使用',
    `expires_at` DATETIME NOT NULL COMMENT '过期时间',
    `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    UNIQUE INDEX `idx_captcha_id` (`captcha_id`),
    INDEX `idx_expires_at` (`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='图形验证码表';

-- ==========================================================================
-- 第四步：清空旧的学生数据（因为数据结构已改变）
-- ==========================================================================

-- 先删除选课记录（外键依赖）
DELETE FROM `enrollments`;

-- 再删除学生数据
DELETE FROM `students`;

-- ==========================================================================
-- 第五步：插入测试数据
-- ==========================================================================

-- 插入测试学生（使用手机号）
INSERT INTO `students` (`username`, `phone`) VALUES
('test_student1', '13800138001'),
('test_student2', '13800138002'),
('test_student3', '13800138003'),
('test_student4', '13800138004'),
('test_student5', '13800138005');

-- 完成
SELECT 'Migration completed successfully!' AS status;
