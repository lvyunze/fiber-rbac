package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
)

// JWTClaims 定义JWT的载荷
type JWTClaims struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// 常见的JWT错误
var (
	ErrTokenExpired     = errors.New("token已过期")
	ErrTokenInvalid     = errors.New("token无效")
	ErrTokenNotProvided = errors.New("未提供token")
)

// GenerateToken 生成JWT令牌
func GenerateToken(userID uint, username string) (string, error) {
	// 获取JWT密钥
	secretKey := viper.GetString("server.jwt_secret")
	if secretKey == "" {
		secretKey = "default-secret-key-for-jwt" // 默认密钥，生产环境应配置
	}

	// 设置过期时间，默认24小时
	expirationTime := time.Now().Add(24 * time.Hour)

	// 创建载荷
	claims := &JWTClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "fiber-rbac",
			Subject:   username,
		},
	}

	// 创建JWT令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken 验证JWT令牌
func ValidateToken(tokenString string) (*JWTClaims, error) {
	if tokenString == "" {
		return nil, ErrTokenNotProvided
	}

	// 获取JWT密钥
	secretKey := viper.GetString("server.jwt_secret")
	if secretKey == "" {
		secretKey = "default-secret-key-for-jwt" // 默认密钥，生产环境应配置
	}

	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	// 验证令牌是否有效
	if !token.Valid {
		return nil, ErrTokenInvalid
	}

	// 获取载荷
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}

// ExtractClaims 从令牌中提取Claims，不验证令牌是否有效
func ExtractClaims(tokenString string) (*JWTClaims, error) {
	if tokenString == "" {
		return nil, ErrTokenNotProvided
	}

	// 获取JWT密钥
	secretKey := viper.GetString("server.jwt_secret")
	if secretKey == "" {
		secretKey = "default-secret-key-for-jwt" // 默认密钥，生产环境应配置
	}

	// 解析令牌，但不验证签名
	token, _ := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})

	// 获取载荷
	if token != nil {
		if claims, ok := token.Claims.(*JWTClaims); ok {
			return claims, nil
		}
	}

	return nil, ErrTokenInvalid
}

// RefreshToken 刷新JWT令牌
func RefreshToken(tokenString string) (string, error) {
	// 验证原始令牌
	claims, err := ValidateToken(tokenString)
	if err != nil && !errors.Is(err, ErrTokenExpired) {
		// 如果是其他错误（不是过期错误），则返回错误
		return "", err
	}

	// 重新生成令牌
	return GenerateToken(claims.UserID, claims.Username)
}
