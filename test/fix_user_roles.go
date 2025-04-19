package test

import (
	"fmt"
	"os"
	"time"

	"github.com/lvyunze/fiber-rbac/config"
	"github.com/lvyunze/fiber-rbac/internal/model"
)

// FixUserRoles 修复用户角色关联
func FixUserRoles() {
	// 加载配置
	cfg, err := config.LoadConfig("../config/config.yaml")
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化数据库
	err = model.InitDB(&cfg.Database, cfg.Env)
	if err != nil {
		fmt.Printf("数据库连接失败: %v\n", err)
		os.Exit(1)
	}

	// 获取数据库连接
	db := model.GetDB()

	// 查询所有用户
	var users []model.User
	if err := db.Find(&users).Error; err != nil {
		fmt.Printf("查询用户失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("找到 %d 个用户\n", len(users))

	// 查询管理员角色
	var adminRole model.Role
	if err := db.Where("code = ?", "admin").First(&adminRole).Error; err != nil {
		fmt.Printf("查询管理员角色失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("找到管理员角色: ID=%d, 名称=%s\n", adminRole.ID, adminRole.Name)

	// 为每个用户检查并创建角色关联
	for _, user := range users {
		// 检查用户是否已经有角色关联
		var count int64
		db.Model(&model.UserRole{}).Where("user_id = ?", user.ID).Count(&count)

		if count == 0 {
			fmt.Printf("用户 %s (ID: %d) 没有角色关联，正在添加...\n", user.Username, user.ID)

			// 创建用户角色关联
			userRole := model.UserRole{
				UserID:    user.ID,
				RoleID:    adminRole.ID,
				CreatedAt: time.Now().Unix(),
			}

			if err := db.Create(&userRole).Error; err != nil {
				fmt.Printf("为用户 %s 创建角色关联失败: %v\n", user.Username, err)
				continue
			}

			fmt.Printf("成功为用户 %s 添加管理员角色\n", user.Username)
		} else {
			fmt.Printf("用户 %s (ID: %d) 已有 %d 个角色关联\n", user.Username, user.ID, count)
		}
	}

	// 验证修复结果
	var userRoleCount int64
	db.Model(&model.UserRole{}).Count(&userRoleCount)
	fmt.Printf("用户角色关联总数: %d\n", userRoleCount)

	fmt.Println("用户角色关联修复完成")
}
