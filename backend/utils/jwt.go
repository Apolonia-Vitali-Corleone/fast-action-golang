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
// 参数:
//   - tokenString: JWT token字符串
// 返回:
//   - *Claims: 解析出的载荷信息
//   - error: 解析或验证失败时的错误
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

// ShouldRefreshToken 判断Token是否需要刷新
// 如果Token剩余有效期小于阈值，则返回true
// 参数:
//   - claims: Token的载荷信息
//   - refreshThreshold: 刷新阈值（默认建议2小时）
// 返回:
//   - bool: true表示需要刷新，false表示不需要
//
// 工作原理：
//   如果Token的剩余有效期 < refreshThreshold，则自动刷新
//   例如：Token有效期24小时，阈值2小时
//   - Token颁发后0-22小时：不刷新
//   - Token颁发后22-24小时：自动刷新
func ShouldRefreshToken(claims *Claims, refreshThreshold time.Duration) bool {
	if claims == nil || claims.ExpiresAt == nil {
		return false
	}

	// 计算剩余有效时间
	timeUntilExpiry := time.Until(claims.ExpiresAt.Time)

	// 如果剩余时间小于阈值，则需要刷新
	return timeUntilExpiry < refreshThreshold && timeUntilExpiry > 0
}

// RefreshToken 刷新Token（生成新的Token）
// 保持userID和role不变，重新设置过期时间
// 参数:
//   - claims: 旧Token的载荷信息
// 返回:
//   - string: 新的Token字符串
//   - error: 生成失败时的错误
func RefreshToken(claims *Claims) (string, error) {
	// 使用旧Token的用户信息生成新Token
	return GenerateToken(claims.UserID, claims.Role)
}
