package common

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	// JWT配置
	JWT_SECRET = "your-secret-key-here" // 建议替换为环境变量
	JWT_EXPIRE = 7 * 24 * time.Hour     // JWT过期时间7天
)

// ==================== JWT自定义Claims ====================
type UserClaims struct {
	jwt.RegisteredClaims
	UserID   int64  `json:"user_id"`
	Username string `json:"username"`
}

// ==================== JWT工具函数 ====================
// GenerateJWT 生成JWT令牌
func GenerateJWT(userID int64, username string) (string, error) {
	expirationTime := time.Now().Add(JWT_EXPIRE)
	claims := &UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "dist_storage",
		},
		UserID:   userID,
		Username: username,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(JWT_SECRET))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ValidateJWT 验证JWT令牌
func ValidateJWT(tokenString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWT_SECRET), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid JWT token")
}
