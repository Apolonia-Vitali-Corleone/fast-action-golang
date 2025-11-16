package controllers

import (
	"course-system/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCaptcha 获取图形验证码
// GET /api/captcha/
// 返回: {captcha_id: "xxx", image: "base64编码的图片"}
func GetCaptcha(c *gin.Context) {
	captchaID, b64s, err := utils.GenerateCaptcha()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成验证码失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"captcha_id": captchaID,
		"image":      b64s,
	})
}

// SendSMSCode 发送短信验证码
// POST /api/sms/send/
// 请求体: {phone: "13800138000", purpose: "register"或"login", captcha_id: "xxx", captcha_code: "1234"}
// 注意: 发送短信前需要验证图形验证码（登录时）
func SendSMSCode(c *gin.Context) {
	var req struct {
		Phone        string `json:"phone" binding:"required"`        // 手机号
		Purpose      string `json:"purpose" binding:"required"`      // 用途: register(注册) / login(登录)
		CaptchaID    string `json:"captcha_id"`                      // 图形验证码ID（登录时必填）
		CaptchaCode  string `json:"captcha_code"`                    // 图形验证码（登录时必填）
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	// 验证purpose参数
	if req.Purpose != "register" && req.Purpose != "login" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用途参数，必须是 register 或 login"})
		return
	}

	// 验证手机号格式（简单验证）
	if len(req.Phone) != 11 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "手机号格式不正确"})
		return
	}

	// 如果是登录，需要先验证图形验证码
	if req.Purpose == "login" {
		if req.CaptchaID == "" || req.CaptchaCode == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "登录时必须提供图形验证码"})
			return
		}

		// 验证图形验证码
		if !utils.VerifyCaptcha(req.CaptchaID, req.CaptchaCode) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "图形验证码错误或已过期"})
			return
		}
	}

	// 发送短信验证码
	if err := utils.SendSMSCode(req.Phone, req.Purpose); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "验证码已发送，请查收短信",
	})
}
