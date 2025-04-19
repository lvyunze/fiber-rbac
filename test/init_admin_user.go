package main

import (
	"fmt"
	"os"
	"time"

	"github.com/lvyunze/fiber-rbac/config"
	"github.com/lvyunze/fiber-rbac/internal/model"
	"github.com/lvyunze/fiber-rbac/internal/pkg/hash"

	"gorm.io/gorm"
)

// main 入口，执行管理员初始化
func main() {
	InitAdminUser()
}

// InitAdminUser 初始化管理员用户和相关角色、权限
func InitAdminUser() {
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

	// 检查 admin 用户名或邮箱是否已存在，避免重复插入和脏数据修正
	var admin model.User
	hashedPwd, _ := hash.GeneratePassword("admin123")
	err = db.Where("username = ? OR email = ?", "admin", "admin@example.com").First(&admin).Error
	if err == nil {
		// 已存在则强制修正 username/email 并重置密码
		admin.Username = "admin"
		admin.Email = "admin@example.com"
		admin.Password = hashedPwd
		admin.UpdatedAt = time.Now().Unix()
		if err := db.Save(&admin).Error; err != nil {
			fmt.Printf("修正 admin 用户失败: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("admin 用户已存在，已修正并重置密码为：admin123")
		return
	} else if err != gorm.ErrRecordNotFound {
		fmt.Printf("查询 admin 用户失败: %v\n", err)
		os.Exit(1)
	}

	// 不存在则插入 admin 用户
	adminUser := &model.User{
		Username:  "admin",
		Email:     "admin@example.com",
		Password:  hashedPwd,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}
	if err := db.Create(adminUser).Error; err != nil {
		fmt.Printf("创建 admin 用户失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("admin 用户初始化成功，用户名：admin，密码：admin123")

	// 初始化默认数据
	if err := initDefaultData(db); err != nil {
		fmt.Printf("初始化默认数据失败: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("初始化管理员用户成功！")
}

// initDefaultData 初始化默认数据
func initDefaultData(db *gorm.DB) error {
	fmt.Println("开始初始化默认数据...")

	// 检查是否已存在管理员角色
	var adminRoleCount int64
	db.Model(&model.Role{}).Where("code = ?", "admin").Count(&adminRoleCount)

	var adminRole *model.Role

	if adminRoleCount == 0 {
		// 创建默认角色
		adminRole = &model.Role{
			Name:        "管理员",
			Code:        "admin",
			Description: "系统管理员，拥有所有权限",
			CreatedAt:   time.Now().Unix(),
		}

		userRole := &model.Role{
			Name:        "普通用户",
			Code:        "user",
			Description: "普通用户，拥有基本权限",
			CreatedAt:   time.Now().Unix(),
		}

		if err := db.Create(adminRole).Error; err != nil {
			fmt.Printf("创建管理员角色失败: %v\n", err)
			return err
		}

		if err := db.Create(userRole).Error; err != nil {
			fmt.Printf("创建普通用户角色失败: %v\n", err)
			return err
		}

		// 创建默认权限
		permissions := []model.Permission{
			{Code: "user:list", Name: "用户列表", Description: "查看用户列表", CreatedAt: time.Now().Unix()},
			{Code: "user:create", Name: "创建用户", Description: "创建新用户", CreatedAt: time.Now().Unix()},
			{Code: "user:update", Name: "更新用户", Description: "更新用户信息", CreatedAt: time.Now().Unix()},
			{Code: "user:delete", Name: "删除用户", Description: "删除用户", CreatedAt: time.Now().Unix()},
			{Code: "role:list", Name: "角色列表", Description: "查看角色列表", CreatedAt: time.Now().Unix()},
			{Code: "role:create", Name: "创建角色", Description: "创建新角色", CreatedAt: time.Now().Unix()},
			{Code: "role:update", Name: "更新角色", Description: "更新角色信息", CreatedAt: time.Now().Unix()},
			{Code: "role:delete", Name: "删除角色", Description: "删除角色", CreatedAt: time.Now().Unix()},
			{Code: "permission:list", Name: "权限列表", Description: "查看权限列表", CreatedAt: time.Now().Unix()},
			{Code: "permission:create", Name: "创建权限", Description: "创建新权限", CreatedAt: time.Now().Unix()},
			{Code: "permission:update", Name: "更新权限", Description: "更新权限信息", CreatedAt: time.Now().Unix()},
			{Code: "permission:delete", Name: "删除权限", Description: "删除权限", CreatedAt: time.Now().Unix()},
		}

		for _, p := range permissions {
			if err := db.Create(&p).Error; err != nil {
				fmt.Printf("创建权限失败 %s: %v\n", p.Code, err)
				return err
			}
		}

		// 为管理员角色分配所有权限
		var allPermissions []model.Permission
		if err := db.Find(&allPermissions).Error; err != nil {
			fmt.Printf("查询所有权限失败: %v\n", err)
			return err
		}

		// 重新查询管理员角色以获取ID
		if err := db.First(&adminRole, "code = ?", "admin").Error; err != nil {
			fmt.Printf("查询管理员角色失败: %v\n", err)
			return err
		}

		// 为管理员角色分配所有权限
		if err := db.Model(&adminRole).Association("Permissions").Append(&allPermissions); err != nil {
			fmt.Printf("为管理员分配权限失败: %v\n", err)
			return err
		}

		// 为普通用户角色分配基本权限
		var basicPermissions []model.Permission
		if err := db.Where("code IN ?", []string{"user:list"}).Find(&basicPermissions).Error; err != nil {
			fmt.Printf("查询基本权限失败: %v\n", err)
			return err
		}

		// 重新查询普通用户角色以获取ID
		if err := db.First(&userRole, "code = ?", "user").Error; err != nil {
			fmt.Printf("查询普通用户角色失败: %v\n", err)
			return err
		}

		// 为普通用户角色分配基本权限
		if err := db.Model(&userRole).Association("Permissions").Append(&basicPermissions); err != nil {
			fmt.Printf("为普通用户分配权限失败: %v\n", err)
			return err
		}
	} else {
		// 如果管理员角色已存在，获取它
		adminRole = &model.Role{}
		if err := db.First(adminRole, "code = ?", "admin").Error; err != nil {
			fmt.Printf("查询管理员角色失败: %v\n", err)
			return err
		}
	}

	// 检查是否已存在管理员用户
	var adminUserCount int64
	db.Model(&model.User{}).Where("username = ?", "admin").Count(&adminUserCount)

	if adminUserCount == 0 {
		// 生成密码哈希
		hashedPassword, err := hash.GeneratePassword("admin123")
		if err != nil {
			fmt.Printf("生成密码哈希失败: %v\n", err)
			return err
		}

		// 创建管理员用户
		adminUser := &model.User{
			Username:  "admin",
			Email:     "admin@example.com",
			Password:  hashedPassword,
			CreatedAt: time.Now().Unix(),
		}

		// 创建用户
		if err := db.Create(adminUser).Error; err != nil {
			fmt.Printf("创建管理员用户失败: %v\n", err)
			return err
		}

		// 为用户分配管理员角色
		if err := db.Model(adminUser).Association("Roles").Append(adminRole); err != nil {
			fmt.Printf("为管理员用户分配角色失败: %v\n", err)
			return err
		}

		fmt.Printf("创建管理员用户成功: %s (密码: admin123)\n", adminUser.Username)
	} else {
		fmt.Println("管理员用户已存在，跳过创建")
	}

	fmt.Println("默认数据初始化完成")
	return nil
}
