package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWT密钥，生产环境应该从环境变量或配置文件读取
var jwtSecret = []byte("your-secret-key-change-in-production")

// Claims JWT载荷结构
type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"` // student 或 teacher
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT Token
// 参数: userID - 用户ID, role - 用户角色
// 返回: token字符串和错误信息
func GenerateToken(userID int, role string) (string, error) {
	// 设置过期时间（24小时）
	expirationTime := time.Now().Add(24 * time.Hour)

	// 创建Claims
	claims := &Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// 创建token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名token
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析JWT Token
// 参数: tokenString - JWT token字符串
// 返回: Claims和错误信息
func ParseToken(tokenString string) (*Claims, error) {
	// 解析token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	// 验证token
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("无效的token")
}
