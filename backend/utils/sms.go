package utils

import (
	"course-system/config"
	"course-system/models"
	"fmt"
	"math/rand"
	"time"

	dysmsapi "github.com/aliyun/alibaba-cloud-sdk-go/services/dysmsapi"
)

// SendSMSCode 发送短信验证码
// 参数: phone - 手机号, purpose - 用途(register/login)
// 返回: error
func SendSMSCode(phone, purpose string) error {
	// 生成6位随机验证码
	rand.Seed(time.Now().UnixNano())
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	// 获取短信配置
	smsConfig := config.GetSMSConfig()

	// 验证配置
	if err := config.ValidateSMSConfig(smsConfig); err != nil {
		// 如果配置不完整，仅记录到数据库，不实际发送（开发模式）
		return saveSMSCodeOnly(phone, code, purpose)
	}

	// 创建阿里云客户端
	client, err := dysmsapi.NewClientWithAccessKey(
		smsConfig.RegionID,
		smsConfig.AccessKeyID,
		smsConfig.AccessKeySecret,
	)
	if err != nil {
		return fmt.Errorf("创建阿里云客户端失败: %v", err)
	}

	// 构建短信请求
	request := dysmsapi.CreateSendSmsRequest()
	request.Scheme = "https"
	request.PhoneNumbers = phone
	request.SignName = smsConfig.SignName
	request.TemplateCode = smsConfig.TemplateCode
	request.TemplateParam = fmt.Sprintf(`{"code":"%s"}`, code)

	// 发送短信
	response, err := client.SendSms(request)
	if err != nil {
		return fmt.Errorf("发送短信失败: %v", err)
	}

	// 检查响应
	if response.Code != "OK" {
		return fmt.Errorf("短信发送失败: %s - %s", response.Code, response.Message)
	}

	// 保存验证码到数据库
	return saveSMSCode(phone, code, purpose)
}

// saveSMSCodeOnly 仅保存验证码到数据库（开发模式，不实际发送短信）
func saveSMSCodeOnly(phone, code, purpose string) error {
	fmt.Printf("[开发模式] 短信验证码: %s (手机号: %s, 用途: %s)\n", code, phone, purpose)
	return saveSMSCode(phone, code, purpose)
}

// saveSMSCode 保存验证码到数据库
func saveSMSCode(phone, code, purpose string) error {
	// 检查是否在1分钟内已发送过验证码（防止频繁发送）
	var lastCode models.SMSCode
	oneMinuteAgo := time.Now().Add(-1 * time.Minute)
	if err := config.DB.Where("phone = ? AND created_at > ?", phone, oneMinuteAgo).
		Order("created_at DESC").First(&lastCode).Error; err == nil {
		return fmt.Errorf("验证码发送过于频繁，请1分钟后再试")
	}

	// 创建验证码记录
	smsCode := models.SMSCode{
		Phone:     phone,
		Code:      code,
		Purpose:   purpose,
		Used:      false,
		ExpiresAt: time.Now().Add(5 * time.Minute), // 5分钟有效期
	}

	if err := config.DB.Create(&smsCode).Error; err != nil {
		return fmt.Errorf("保存验证码失败: %v", err)
	}

	return nil
}

// VerifySMSCode 验证短信验证码
// 参数: phone - 手机号, code - 验证码, purpose - 用途
// 返回: 是否验证成功
func VerifySMSCode(phone, code, purpose string) bool {
	// 查询最新的未使用验证码
	var smsCode models.SMSCode
	if err := config.DB.Where("phone = ? AND purpose = ? AND used = ?", phone, purpose, false).
		Order("created_at DESC").First(&smsCode).Error; err != nil {
		return false // 验证码不存在
	}

	// 检查是否过期
	if time.Now().After(smsCode.ExpiresAt) {
		return false
	}

	// 验证码是否匹配
	if smsCode.Code != code {
		return false
	}

	// 标记为已使用
	config.DB.Model(&smsCode).Update("used", true)

	return true
}

// CleanExpiredSMSCodes 清理过期的短信验证码
// 建议定期调用此函数（如每小时）
func CleanExpiredSMSCodes() error {
	// 删除过期的验证码记录（保留24小时内的记录用于审计）
	twentyFourHoursAgo := time.Now().Add(-24 * time.Hour)
	result := config.DB.Where("expires_at < ?", twentyFourHoursAgo).Delete(&models.SMSCode{})
	if result.Error != nil {
		return fmt.Errorf("清理过期验证码失败: %v", result.Error)
	}
	return nil
}
