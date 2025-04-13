package model

import (
	"log/slog"
	"time"

	"gorm.io/gorm"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(db *gorm.DB) error {
	slog.Info("开始数据库迁移")

	// 迁移用户、角色和权限表
	err := db.AutoMigrate(
		&User{},
		&Role{},
		&Permission{},
	)

	if err != nil {
		slog.Error("数据库迁移失败", "error", err)
		return err
	}

	slog.Info("数据库迁移完成")
	return nil
}

// InitDefaultData 初始化默认数据
func InitDefaultData(db *gorm.DB) error {
	slog.Info("开始初始化默认数据")

	// 检查是否已存在管理员角色
	var adminRoleCount int64
	db.Model(&Role{}).Where("code = ?", "admin").Count(&adminRoleCount)

	if adminRoleCount == 0 {
		// 创建默认角色
		adminRole := &Role{
			Name:        "管理员",
			Code:        "admin",
			Description: "系统管理员，拥有所有权限",
		}

		userRole := &Role{
			Name:        "普通用户",
			Code:        "user",
			Description: "普通用户，拥有基本权限",
		}

		if err := db.Create(adminRole).Error; err != nil {
			slog.Error("创建管理员角色失败", "error", err)
			return err
		}

		if err := db.Create(userRole).Error; err != nil {
			slog.Error("创建普通用户角色失败", "error", err)
			return err
		}

		// 创建默认权限
		permissions := []Permission{
			{Code: "user:list", Name: "用户列表", Description: "查看用户列表"},
			{Code: "user:create", Name: "创建用户", Description: "创建新用户"},
			{Code: "user:update", Name: "更新用户", Description: "更新用户信息"},
			{Code: "user:delete", Name: "删除用户", Description: "删除用户"},
			{Code: "role:list", Name: "角色列表", Description: "查看角色列表"},
			{Code: "role:create", Name: "创建角色", Description: "创建新角色"},
			{Code: "role:update", Name: "更新角色", Description: "更新角色信息"},
			{Code: "role:delete", Name: "删除角色", Description: "删除角色"},
			{Code: "permission:list", Name: "权限列表", Description: "查看权限列表"},
			{Code: "permission:create", Name: "创建权限", Description: "创建新权限"},
			{Code: "permission:update", Name: "更新权限", Description: "更新权限信息"},
			{Code: "permission:delete", Name: "删除权限", Description: "删除权限"},
		}

		for _, p := range permissions {
			if err := db.Create(&p).Error; err != nil {
				slog.Error("创建权限失败", "code", p.Code, "error", err)
				return err
			}
		}

		// 为管理员角色分配所有权限
		var allPermissions []Permission
		if err := db.Find(&allPermissions).Error; err != nil {
			slog.Error("查询所有权限失败", "error", err)
			return err
		}

		// 重新查询管理员角色以获取ID
		if err := db.First(&adminRole, "code = ?", "admin").Error; err != nil {
			slog.Error("查询管理员角色失败", "error", err)
			return err
		}

		// 为管理员角色分配所有权限
		if err := db.Model(&adminRole).Association("Permissions").Append(&allPermissions); err != nil {
			slog.Error("为管理员分配权限失败", "error", err)
			return err
		}

		// 为普通用户角色分配基本权限
		var basicPermissions []Permission
		if err := db.Where("code IN ?", []string{"user:list"}).Find(&basicPermissions).Error; err != nil {
			slog.Error("查询基本权限失败", "error", err)
			return err
		}

		// 重新查询普通用户角色以获取ID
		if err := db.First(&userRole, "code = ?", "user").Error; err != nil {
			slog.Error("查询普通用户角色失败", "error", err)
			return err
		}

		// 为普通用户角色分配基本权限
		if err := db.Model(&userRole).Association("Permissions").Append(&basicPermissions); err != nil {
			slog.Error("为普通用户分配权限失败", "error", err)
			return err
		}

		// 创建默认管理员用户
		var adminUserCount int64
		db.Model(&User{}).Where("username = ?", "admin").Count(&adminUserCount)

		if adminUserCount == 0 {
			adminUser := &User{
				Username:  "admin",
				Email:     "admin@example.com",
				Password:  "$argon2id$v=19$m=65536,t=3,p=2$NTJHbkVHYmFrWUJBQUFBQQ$TzdNbW9IMXBBb1ZaVE5xUXBtVDJRZz09", // 密码: admin123
				CreatedAt: time.Now().Unix(),
			}

			if err := db.Create(adminUser).Error; err != nil {
				slog.Error("创建管理员用户失败", "error", err)
				return err
			}

			// 创建用户角色关联
			userRole := UserRole{
				UserID:    adminUser.ID,
				RoleID:    adminRole.ID,
				CreatedAt: time.Now().Unix(),
			}

			if err := db.Create(&userRole).Error; err != nil {
				slog.Error("为管理员用户分配角色失败", "error", err)
				return err
			}
		}
	}

	slog.Info("默认数据初始化完成")
	return nil
}
