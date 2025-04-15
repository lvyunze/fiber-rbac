package repository

import (
	"time"
	"github.com/lvyunze/fiber-rbac/internal/model"
	"gorm.io/gorm"
)

// RefreshTokenRepository 刷新令牌数据访问接口
type RefreshTokenRepository interface {
	Create(token *model.UserRefreshToken) error
	FindValid(token string) (*model.UserRefreshToken, error)
	MarkUsed(token string) error
	RevokeByUser(userID uint64) error
}

type refreshTokenRepo struct {
	db *gorm.DB
}

// NewRefreshTokenRepository 创建刷新令牌仓库
func NewRefreshTokenRepository(db *gorm.DB) RefreshTokenRepository {
	return &refreshTokenRepo{db: db}
}

func (r *refreshTokenRepo) Create(token *model.UserRefreshToken) error {
	return r.db.Create(token).Error
}

func (r *refreshTokenRepo) FindValid(token string) (*model.UserRefreshToken, error) {
	var t model.UserRefreshToken
	err := r.db.Where("token = ? AND used = false AND revoked = false AND expires_at > ?", token, time.Now()).First(&t).Error
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *refreshTokenRepo) MarkUsed(token string) error {
	return r.db.Model(&model.UserRefreshToken{}).Where("token = ?", token).Update("used", true).Error
}

func (r *refreshTokenRepo) RevokeByUser(userID uint64) error {
	return r.db.Model(&model.UserRefreshToken{}).Where("user_id = ? AND used = false AND revoked = false", userID).Update("revoked", true).Error
}
