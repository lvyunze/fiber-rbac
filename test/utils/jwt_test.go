package utils_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lvyunze/fiber-rbac/internal/utils"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAndValidateToken(t *testing.T) {
	// 设置测试用的JWT密钥
	viper.Set("server.jwt_secret", "test-secret-key")

	// 测试数据
	userID := uint(1)
	username := "testuser"

	// 生成令牌
	token, err := utils.GenerateToken(userID, username)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// 验证令牌
	claims, err := utils.ValidateToken(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
	assert.Equal(t, "fiber-rbac", claims.Issuer)
	assert.Equal(t, username, claims.Subject)
}

func TestExpiredToken(t *testing.T) {
	// 设置测试用的JWT密钥
	secretKey := "test-secret-key"
	viper.Set("server.jwt_secret", secretKey)

	// 创建一个已过期的令牌
	expirationTime := time.Now().Add(-1 * time.Hour) // 1小时前过期
	claims := &utils.JWTClaims{
		UserID:   uint(1),
		Username: "testuser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)), // 2小时前签发
			Issuer:    "fiber-rbac",
			Subject:   "testuser",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secretKey))

	// 验证已过期的令牌
	_, err := utils.ValidateToken(tokenString)
	assert.Error(t, err)
	assert.Equal(t, utils.ErrTokenExpired, err)
}

func TestInvalidToken(t *testing.T) {
	// 设置测试用的JWT密钥
	viper.Set("server.jwt_secret", "test-secret-key")

	// 测试空令牌
	_, err := utils.ValidateToken("")
	assert.Error(t, err)
	assert.Equal(t, utils.ErrTokenNotProvided, err)

	// 测试无效令牌
	_, err = utils.ValidateToken("invalid.token.string")
	assert.Error(t, err)
	assert.Equal(t, utils.ErrTokenInvalid, err)
}

func TestRefreshToken(t *testing.T) {
	// 设置测试用的JWT密钥
	viper.Set("server.jwt_secret", "test-secret-key")

	// 测试数据
	userID := uint(1)
	username := "testuser"

	// 生成令牌
	token, err := utils.GenerateToken(userID, username)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// 刷新令牌
	newToken, err := utils.RefreshToken(token)
	assert.NoError(t, err)
	assert.NotEmpty(t, newToken)

	// 确保新令牌有效
	claims, err := utils.ValidateToken(newToken)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, username, claims.Username)
}
