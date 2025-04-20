package errors

import "errors"

// 定义通用错误类型
var (
	// 用户相关错误
	ErrUserNotFound      = errors.New("用户不存在")
	ErrInvalidCredentials = errors.New("用户名或密码错误")
	ErrUserExists         = errors.New("用户名已存在")
	ErrEmailExists        = errors.New("邮箱已被使用")

	// 角色相关错误
	ErrRoleNotFound      = errors.New("角色不存在")
	ErrRoleExists        = errors.New("角色已存在")
	ErrRoleInUse         = errors.New("角色正在使用中，无法删除")

	// 权限相关错误
	ErrPermissionNotFound = errors.New("权限不存在")
	ErrPermissionExists   = errors.New("权限已存在")
	ErrPermissionInUse    = errors.New("权限正在使用中，无法删除")

	// 令牌相关错误
	ErrInvalidToken      = errors.New("无效的令牌")
	ErrExpiredToken      = errors.New("令牌已过期")
	ErrInvalidTokenType  = errors.New("无效的令牌类型")

	// 数据库相关错误
	ErrDB = errors.New("数据库异常")
)
