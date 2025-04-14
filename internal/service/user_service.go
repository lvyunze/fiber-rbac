package service

import (
	"log/slog"
	"github.com/lvyunze/fiber-rbac/config"
	"github.com/lvyunze/fiber-rbac/internal/model"
	"github.com/lvyunze/fiber-rbac/internal/pkg/errors"
	"github.com/lvyunze/fiber-rbac/internal/pkg/hash"
	"github.com/lvyunze/fiber-rbac/internal/pkg/jwt"
	"github.com/lvyunze/fiber-rbac/internal/repository"
	"github.com/lvyunze/fiber-rbac/internal/schema"
)

// UserService 用户服务接口
type UserService interface {
	Login(req *schema.LoginRequest) (*schema.LoginResponse, error)
	RefreshToken(token string) (*schema.LoginResponse, error)
	CheckPermission(userID uint64, permission string) (bool, error)
	GetProfile(userID uint64) (*schema.UserResponse, error)
	Create(req *schema.CreateUserRequest) (uint64, error)
	Update(req *schema.UpdateUserRequest) error
	Delete(id uint64) error
	GetByID(id uint64) (*schema.UserResponse, error)
	List(req *schema.ListUserRequest) (*schema.ListUserResponse, error)
	AssignRole(userID uint64, roleIDs []uint64) error
	GetRoles(userID uint64) ([]schema.RoleResponse, error)
}

// userService 用户服务实现
type userService struct {
	userRepo       repository.UserRepository
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
	tokenService   *jwt.TokenService
}

// NewUserService 创建用户服务实例
func NewUserService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	permissionRepo repository.PermissionRepository,
	jwtConfig *config.JWTConfig,
) UserService {
	return &userService{
		userRepo:       userRepo,
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
		tokenService:   jwt.NewTokenService(jwtConfig),
	}
}

// Login 用户登录
func (s *userService) Login(req *schema.LoginRequest) (*schema.LoginResponse, error) {
	// 根据用户名查找用户
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		slog.Error("查询用户失败", "error", err)
		return nil, err
	}

	// 用户不存在
	if user == nil {
		return nil, errors.ErrInvalidCredentials
	}

	// 验证密码
	valid, err := hash.VerifyPassword(req.Password, user.Password)
	if err != nil {
		slog.Error("验证密码失败", "error", err)
		return nil, err
	}

	if !valid {
		return nil, errors.ErrInvalidCredentials
	}

	// 生成JWT令牌
	accessToken, refreshToken, err := s.tokenService.GenerateTokenPair(user.ID, user.Username)
	if err != nil {
		slog.Error("生成令牌失败", "error", err)
		return nil, err
	}

	return &schema.LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.tokenService.Config.Expire,
	}, nil
}

// RefreshToken 刷新令牌
func (s *userService) RefreshToken(token string) (*schema.LoginResponse, error) {
	// 验证刷新令牌
	claims, err := s.tokenService.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	// 检查令牌类型
	if claims.TokenType != "refresh" {
		return nil, errors.ErrInvalidTokenType
	}

	// 检查用户是否存在
	user, err := s.userRepo.GetByID(claims.UserID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	// 生成新令牌
	accessToken, _, err := s.tokenService.GenerateTokenPair(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	return &schema.LoginResponse{
		Token:     accessToken,
		ExpiresIn: s.tokenService.Config.Expire,
	}, nil
}

// CheckPermission 检查用户权限
func (s *userService) CheckPermission(userID uint64, permission string) (bool, error) {
	// 获取用户信息
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return false, err
	}

	if user == nil {
		return false, errors.ErrUserNotFound
	}

	// 检查用户角色和权限
	for _, role := range user.Roles {
		for _, perm := range role.Permissions {
			if perm.Code == permission {
				return true, nil
			}
		}
	}

	return false, nil
}

// GetProfile 获取用户个人信息
func (s *userService) GetProfile(userID uint64) (*schema.UserResponse, error) {
	// 获取用户信息
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	return s.convertToUserResponse(user), nil
}

// Create 创建用户
func (s *userService) Create(req *schema.CreateUserRequest) (uint64, error) {
	// 检查用户名是否已存在
	existingUser, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return 0, err
	}

	if existingUser != nil {
		return 0, errors.ErrUserExists
	}

	// 检查邮箱是否已存在
	existingEmail, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return 0, err
	}

	if existingEmail != nil {
		return 0, errors.ErrEmailExists
	}

	// 生成密码哈希
	hashedPassword, err := hash.GeneratePassword(req.Password)
	if err != nil {
		slog.Error("生成密码哈希失败", "error", err)
		return 0, err
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := s.userRepo.Create(user); err != nil {
		return 0, err
	}

	// 如果有指定角色，添加角色关联
	if len(req.RoleIDs) > 0 {
		if err := s.userRepo.AddRoles(user.ID, req.RoleIDs); err != nil {
			slog.Error("添加用户角色失败", "error", err)
			// 不返回错误，继续执行
		}
	}

	return user.ID, nil
}

// Update 更新用户
func (s *userService) Update(req *schema.UpdateUserRequest) error {
	// 检查用户是否存在
	existingUser, err := s.userRepo.GetByID(req.ID)
	if err != nil {
		return err
	}

	if existingUser == nil {
		return errors.ErrUserNotFound
	}

	// 检查用户名是否已被其他用户使用
	if req.Username != existingUser.Username {
		user, err := s.userRepo.GetByUsername(req.Username)
		if err != nil {
			return err
		}

		if user != nil && user.ID != req.ID {
			return errors.ErrUserExists
		}
	}

	// 检查邮箱是否已被其他用户使用
	if req.Email != existingUser.Email {
		user, err := s.userRepo.GetByEmail(req.Email)
		if err != nil {
			return err
		}

		if user != nil && user.ID != req.ID {
			return errors.ErrEmailExists
		}
	}

	// 更新用户信息
	updatedUser := &model.User{
		ID:       req.ID,
		Username: req.Username,
		Email:    req.Email,
	}

	// 如果提供了密码，更新密码
	if req.Password != "" {
		hashedPassword, err := hash.GeneratePassword(req.Password)
		if err != nil {
			return err
		}
		updatedUser.Password = hashedPassword
	}

	if err := s.userRepo.Update(updatedUser); err != nil {
		return err
	}

	// 如果提供了角色ID，更新用户角色
	if req.RoleIDs != nil {
		if err := s.userRepo.UpdateRoles(req.ID, req.RoleIDs); err != nil {
			slog.Error("更新用户角色失败", "error", err)
			// 不返回错误，继续执行
		}
	}

	return nil
}

// Delete 删除用户
func (s *userService) Delete(id uint64) error {
	// 检查用户是否存在
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.ErrUserNotFound
	}

	// 删除用户
	return s.userRepo.Delete(id)
}

// GetByID 根据ID获取用户
func (s *userService) GetByID(id uint64) (*schema.UserResponse, error) {
	// 获取用户信息
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	return s.convertToUserResponse(user), nil
}

// List 获取用户列表
func (s *userService) List(req *schema.ListUserRequest) (*schema.ListUserResponse, error) {
	// 获取用户列表
	users, total, err := s.userRepo.List(req.Page, req.PageSize, req.Keyword)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	items := make([]schema.UserResponse, 0, len(users))
	for _, user := range users {
		items = append(items, *s.convertToUserResponse(user))
	}

	return &schema.ListUserResponse{
		Total: total,
		Items: items,
	}, nil
}

// AssignRole 分配角色给用户
func (s *userService) AssignRole(userID uint64, roleIDs []uint64) error {
	// 检查用户是否存在
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	if user == nil {
		return errors.ErrUserNotFound
	}

	// 检查所有角色是否存在
	for _, roleID := range roleIDs {
		role, err := s.roleRepo.GetByID(roleID)
		if err != nil {
			return err
		}
		if role == nil {
			return errors.ErrRoleNotFound
		}
	}

	// 更新用户角色
	return s.userRepo.UpdateRoles(userID, roleIDs)
}

// GetRoles 获取用户的角色列表
func (s *userService) GetRoles(userID uint64) ([]schema.RoleResponse, error) {
	// 检查用户是否存在
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.ErrUserNotFound
	}

	// 获取并加载角色
	user, err = s.userRepo.GetUserWithRoles(userID)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	roles := make([]schema.RoleResponse, 0, len(user.Roles))
	for _, role := range user.Roles {
		roles = append(roles, schema.RoleResponse{
			ID:          role.ID,
			Code:        role.Code,
			Name:        role.Name,
			Description: role.Description,
			CreatedAt:   role.CreatedAt,
		})
	}

	return roles, nil
}

// convertToUserResponse 将用户模型转换为响应结构
func (s *userService) convertToUserResponse(user *model.User) *schema.UserResponse {
	response := &schema.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		Roles:     make([]schema.RoleSimple, 0, len(user.Roles)),
	}

	// 添加角色信息
	for _, role := range user.Roles {
		response.Roles = append(response.Roles, schema.RoleSimple{
			ID:   role.ID,
			Code: role.Code,
			Name: role.Name,
		})
	}

	return response
}
