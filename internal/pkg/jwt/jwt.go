package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/lvyunze/fiber-rbac/config"
)

// 定义错误类型
var (
	ErrInvalidToken = errors.New("无效的令牌")
	ErrExpiredToken = errors.New("令牌已过期")
)

// Claims 自定义JWT声明结构
type Claims struct {
	UserID    uint64 `json:"user_id"`
	Username  string `json:"username"`
	TokenType string `json:"token_type"` // access 或 refresh
	jwt.RegisteredClaims
}

// TokenService JWT令牌服务
type TokenService struct {
	Config *config.JWTConfig
}

// NewTokenService 创建JWT服务实例
func NewTokenService(cfg *config.JWTConfig) *TokenService {
	return &TokenService{Config: cfg}
}

// GenerateToken 生成JWT令牌
func (s *TokenService) GenerateToken(userID uint64, username string, tokenType string) (string, error) {
	// 确定过期时间
	var expiry time.Duration
	if tokenType == "refresh" {
		expiry = time.Duration(s.Config.RefreshExpire) * time.Second
	} else {
		expiry = time.Duration(s.Config.Expire) * time.Second
	}

	// 创建JWT声明
	claims := &Claims{
		UserID:    userID,
		Username:  username,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
			Issuer:    "rbac-system",
		},
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	tokenString, err := token.SignedString([]byte(s.Config.Secret))
	if err != nil {
		return "", fmt.Errorf("生成令牌失败: %w", err)
	}

	return tokenString, nil
}

// ValidateToken 验证JWT令牌
func (s *TokenService) ValidateToken(tokenString string) (*Claims, error) {
	// 解析令牌
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("意外的签名方法: %v", token.Header["alg"])
		}
		return []byte(s.Config.Secret), nil
	})

	// 处理解析错误
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	// 提取声明
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	return claims, nil
}

// GenerateTokenPair 生成访问令牌和刷新令牌对
func (s *TokenService) GenerateTokenPair(userID uint64, username string) (accessToken string, refreshToken string, err error) {
	// 生成访问令牌
	accessToken, err = s.GenerateToken(userID, username, "access")
	if err != nil {
		return "", "", err
	}

	// 生成刷新令牌
	refreshToken, err = s.GenerateToken(userID, username, "refresh")
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
