package config

import (
	"fmt"
	"os"
)

// SMSConfig 短信服务配置
type SMSConfig struct {
	AccessKeyID     string // 阿里云AccessKey ID
	AccessKeySecret string // 阿里云AccessKey Secret
	SignName        string // 短信签名
	TemplateCode    string // 短信模板代码
	RegionID        string // 区域ID，默认cn-hangzhou
}

// GetSMSConfig 获取短信配置
// 从环境变量读取配置，如果没有设置则使用默认值
func GetSMSConfig() SMSConfig {
	return SMSConfig{
		AccessKeyID:     getEnv("ALIYUN_ACCESS_KEY_ID", ""),
		AccessKeySecret: getEnv("ALIYUN_ACCESS_KEY_SECRET", ""),
		SignName:        getEnv("ALIYUN_SMS_SIGN_NAME", "课程系统"),
		TemplateCode:    getEnv("ALIYUN_SMS_TEMPLATE_CODE", "SMS_154950909"), // 示例模板代码
		RegionID:        getEnv("ALIYUN_REGION_ID", "cn-hangzhou"),
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// ValidateSMSConfig 验证短信配置是否完整
func ValidateSMSConfig(config SMSConfig) error {
	if config.AccessKeyID == "" {
		return fmt.Errorf("AccessKeyID 未配置，请设置环境变量 ALIYUN_ACCESS_KEY_ID")
	}
	if config.AccessKeySecret == "" {
		return fmt.Errorf("AccessKeySecret 未配置，请设置环境变量 ALIYUN_ACCESS_KEY_SECRET")
	}
	if config.SignName == "" {
		return fmt.Errorf("SignName 未配置，请设置环境变量 ALIYUN_SMS_SIGN_NAME")
	}
	if config.TemplateCode == "" {
		return fmt.Errorf("TemplateCode 未配置，请设置环境变量 ALIYUN_SMS_TEMPLATE_CODE")
	}
	return nil
}
