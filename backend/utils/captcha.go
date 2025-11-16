package utils

import (
	"course-system/config"
	"course-system/models"
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/mojocn/base64Captcha"
)

var captchaStore = base64Captcha.DefaultMemStore

// GenerateCaptcha 生成图形验证码
// 返回: captchaID(验证码ID), base64Image(base64编码的图片), error
func GenerateCaptcha() (string, string, error) {
	// 配置验证码参数
	driver := base64Captcha.NewDriverDigit(
		80,   // 图片高度
		240,  // 图片宽度
		4,    // 验证码长度
		0.7,  // 干扰线密度
		80,   // 最大倾斜角度
	)

	// 生成验证码
	captcha := base64Captcha.NewCaptcha(driver, captchaStore)
	captchaID, b64s, _, err := captcha.Generate()
	if err != nil {
		return "", "", fmt.Errorf("生成验证码失败: %v", err)
	}

	// 保存到数据库
	code := captchaStore.Get(captchaID, false) // 获取验证码答案但不删除
	if code == "" {
		return "", "", fmt.Errorf("获取验证码答案失败")
	}

	captchaRecord := models.CaptchaCode{
		CaptchaID: captchaID,
		Code:      code,
		Used:      false,
		ExpiresAt: time.Now().Add(5 * time.Minute), // 5分钟有效期
	}

	if err := config.DB.Create(&captchaRecord).Error; err != nil {
		return "", "", fmt.Errorf("保存验证码失败: %v", err)
	}

	return captchaID, b64s, nil
}

// GenerateColorCaptcha 生成彩色图形验证码（更美观的版本）
// 返回: captchaID(验证码ID), base64Image(base64编码的图片), error
func GenerateColorCaptcha() (string, string, error) {
	// 生成唯一ID
	captchaID := uuid.New().String()

	// 生成4位随机数字验证码
	rand.Seed(time.Now().UnixNano())
	code := fmt.Sprintf("%04d", rand.Intn(10000))

	// 配置彩色字符验证码
	driver := base64Captcha.NewDriverString(
		80,    // 图片高度
		240,   // 图片宽度
		6,     // 噪点数量
		1,     // 干扰线选项
		4,     // 验证码长度
		code,  // 验证码内容
		&color.RGBA{R: 0, G: 0, B: 0, A: 255}, // 背景色
		nil,   // 字体（使用默认）
		[]string{"0123456789"}, // 字符源
	)

	// 生成验证码图片
	item, err := driver.DrawCaptcha(code)
	if err != nil {
		return "", "", fmt.Errorf("生成验证码图片失败: %v", err)
	}

	// 转换为base64
	b64s := item.EncodeB64string()

	// 保存到数据库
	captchaRecord := models.CaptchaCode{
		CaptchaID: captchaID,
		Code:      code,
		Used:      false,
		ExpiresAt: time.Now().Add(5 * time.Minute), // 5分钟有效期
	}

	if err := config.DB.Create(&captchaRecord).Error; err != nil {
		return "", "", fmt.Errorf("保存验证码失败: %v", err)
	}

	return captchaID, b64s, nil
}

// VerifyCaptcha 验证图形验证码
// 参数: captchaID - 验证码ID, code - 用户输入的验证码
// 返回: 是否验证成功
func VerifyCaptcha(captchaID, code string) bool {
	// 从数据库查询验证码
	var captchaRecord models.CaptchaCode
	if err := config.DB.Where("captcha_id = ?", captchaID).First(&captchaRecord).Error; err != nil {
		return false // 验证码不存在
	}

	// 检查是否已使用
	if captchaRecord.Used {
		return false
	}

	// 检查是否过期
	if time.Now().After(captchaRecord.ExpiresAt) {
		return false
	}

	// 验证码是否匹配（不区分大小写）
	if captchaRecord.Code != code {
		return false
	}

	// 标记为已使用
	config.DB.Model(&captchaRecord).Update("used", true)

	return true
}

// CleanExpiredCaptcha 清理过期的图形验证码
// 建议定期调用此函数（如每小时）
func CleanExpiredCaptcha() error {
	// 删除过期的验证码记录
	result := config.DB.Where("expires_at < ?", time.Now()).Delete(&models.CaptchaCode{})
	if result.Error != nil {
		return fmt.Errorf("清理过期验证码失败: %v", result.Error)
	}
	return nil
}
