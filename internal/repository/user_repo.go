package repository

import (
	"errors"
	"fmt"
	"github.com/lvyunze/fiber-rbac/internal/model"
	"strings"

	"gorm.io/gorm"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	Create(user *model.User) error
	Update(user *model.User) error
	Delete(id uint64) error
	GetByID(id uint64) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	List(page, pageSize int, keyword string) ([]*model.User, int64, error)
	AddRoles(userID uint64, roleIDs []uint64) error
	RemoveRoles(userID uint64, roleIDs []uint64) error
	UpdateRoles(userID uint64, roleIDs []uint64) error
	GetUserWithRoles(userID uint64) (*model.User, error)
}

// userRepo 用户仓储实现
type userRepo struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储实例
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

// Create 创建用户
func (r *userRepo) Create(user *model.User) error {
	return r.db.Create(user).Error
}

// Update 更新用户
func (r *userRepo) Update(user *model.User) error {
	// 只更新非零值字段
	return r.db.Model(user).Updates(user).Error
}

// Delete 删除用户（软删除）
func (r *userRepo) Delete(id uint64) error {
	// 使用软删除
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("deleted_at", model.SoftDelete()).Error
}

// GetByID 根据ID获取用户
func (r *userRepo) GetByID(id uint64) (*model.User, error) {
	var user model.User
	result := r.db.Preload("Roles").Where("id = ? AND deleted_at IS NULL", id).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil // 用户不存在返回nil，而不是错误
		}
		return nil, result.Error
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (r *userRepo) GetByUsername(username string) (*model.User, error) {
	var user model.User
	result := r.db.Preload("Roles").Where("username = ? AND deleted_at IS NULL", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (r *userRepo) GetByEmail(email string) (*model.User, error) {
	var user model.User
	result := r.db.Preload("Roles").Where("email = ? AND deleted_at IS NULL", email).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

// GetUserWithRoles 获取用户及其角色和权限
func (r *userRepo) GetUserWithRoles(userID uint64) (*model.User, error) {
	var user model.User
	result := r.db.Preload("Roles.Permissions").Where("id = ? AND deleted_at IS NULL", userID).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

// List 获取用户列表
func (r *userRepo) List(page, pageSize int, keyword string) ([]*model.User, int64, error) {
	var users []*model.User
	var total int64

	// 默认分页参数
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// 构建查询
	query := r.db.Model(&model.User{}).Where("deleted_at IS NULL")

	// 添加关键词搜索
	if keyword != "" {
		keyword = fmt.Sprintf("%%%s%%", strings.ToLower(keyword))
		query = query.Where("LOWER(username) LIKE ? OR LOWER(email) LIKE ?", keyword, keyword)
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取分页数据
	offset := (page - 1) * pageSize
	if err := query.Preload("Roles").Offset(offset).Limit(pageSize).Order("id DESC").Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// AddRoles 为用户添加角色
func (r *userRepo) AddRoles(userID uint64, roleIDs []uint64) error {
	// 开启事务
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 检查用户是否存在
		var user model.User
		if err := tx.Where("id = ? AND deleted_at IS NULL", userID).First(&user).Error; err != nil {
			return err
		}

		// 添加角色关联
		for _, roleID := range roleIDs {
			// 检查角色是否存在
			var role model.Role
			if err := tx.Where("id = ? AND deleted_at IS NULL", roleID).First(&role).Error; err != nil {
				return fmt.Errorf("角色ID %d 不存在: %w", roleID, err)
			}

			// 检查关联是否已存在
			var count int64
			tx.Model(&model.UserRole{}).Where("user_id = ? AND role_id = ?", userID, roleID).Count(&count)
			if count == 0 {
				// 创建关联
				userRole := model.UserRole{
					UserID: userID,
					RoleID: roleID,
				}
				if err := tx.Create(&userRole).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// RemoveRoles 移除用户的角色
func (r *userRepo) RemoveRoles(userID uint64, roleIDs []uint64) error {
	return r.db.Where("user_id = ? AND role_id IN ?", userID, roleIDs).Delete(&model.UserRole{}).Error
}

// UpdateRoles 更新用户的角色
func (r *userRepo) UpdateRoles(userID uint64, roleIDs []uint64) error {
	// 开启事务
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 删除所有现有角色
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
			return err
		}

		// 添加新角色
		if len(roleIDs) > 0 {
			return r.AddRoles(userID, roleIDs)
		}

		return nil
	})
}
