package service_test

import (
	"github.com/lvyunze/fiber-rbac/internal/model"
	"github.com/lvyunze/fiber-rbac/internal/pkg/errors"
	"github.com/lvyunze/fiber-rbac/internal/schema"
	"github.com/lvyunze/fiber-rbac/internal/service"
	"github.com/lvyunze/fiber-rbac/test/mocks"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// 测试权限创建功能
func TestPermissionService_Create(t *testing.T) {
	// 测试用例表
	tests := []struct {
		name           string
		request        *schema.CreatePermissionRequest
		mockSetup      func(mockPermRepo *mocks.MockPermissionRepository)
		expectedID     uint64
		expectedError  error
		expectedCalled bool
	}{
		{
			name: "成功创建权限",
			request: &schema.CreatePermissionRequest{
				Name:        "测试权限",
				Code:        "test:permission",
				Description: "这是一个测试权限",
			},
			mockSetup: func(mockPermRepo *mocks.MockPermissionRepository) {
				// 模拟权限名称不存在
				mockPermRepo.On("GetByName", "测试权限").Return(nil, nil)
				
				// 模拟权限标识不存在
				mockPermRepo.On("GetByCode", "test:permission").Return(nil, nil)
				
				// 模拟创建权限成功
				mockPermRepo.On("Create", mock.AnythingOfType("*model.Permission")).Run(func(args mock.Arguments) {
					perm := args.Get(0).(*model.Permission)
					perm.ID = 1 // 设置ID
				}).Return(nil)
			},
			expectedID:     1,
			expectedError:  nil,
			expectedCalled: true,
		},
		{
			name: "权限标识已存在",
			request: &schema.CreatePermissionRequest{
				Name:        "已存在权限",
				Code:        "existing:permission",
				Description: "这是一个已存在的权限",
			},
			mockSetup: func(mockPermRepo *mocks.MockPermissionRepository) {
				// 模拟权限名称存在
				mockPermRepo.On("GetByName", "已存在权限").Return(nil, nil)
				
				// 模拟权限标识已存在
				existingPerm := &model.Permission{ID: 1, Code: "existing:permission"}
				mockPermRepo.On("GetByCode", "existing:permission").Return(existingPerm, nil)
				
				// 模拟创建权限失败，因为权限标识已被其他权限使用，所以我们需要修改mockSetup函数来模拟这种情况
				mockPermRepo.On("Create", mock.AnythingOfType("*model.Permission")).Return(nil)
			},
			expectedID:     0,
			expectedError:  errors.ErrPermissionExists,
			expectedCalled: false,
		},
		{
			name: "数据库错误",
			request: &schema.CreatePermissionRequest{
				Name:        "测试权限",
				Code:        "test:permission",
				Description: "这是一个测试权限",
			},
			mockSetup: func(mockPermRepo *mocks.MockPermissionRepository) {
				// 模拟权限名称不存在
				mockPermRepo.On("GetByName", "测试权限").Return(nil, nil)
				
				// 模拟权限标识不存在
				mockPermRepo.On("GetByCode", "test:permission").Return(nil, nil)
				
				// 模拟创建权限失败
				mockPermRepo.On("Create", mock.AnythingOfType("*model.Permission")).Return(assert.AnError)
			},
			expectedID:     0,
			expectedError:  assert.AnError,
			expectedCalled: true,
		},
	}

	// 执行测试用例
	for _, tt := range tests {
		tt := tt // 防止闭包问题
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟对象
			mockPermRepo := new(mocks.MockPermissionRepository)
			
			// 设置模拟行为
			tt.mockSetup(mockPermRepo)
			
			// 创建服务实例
			permissionService := service.NewPermissionService(mockPermRepo)
			
			// 调用被测试的方法
			id, err := permissionService.Create(tt.request)
			
			// 验证结果
			assert.Equal(t, tt.expectedID, id)
			assert.Equal(t, tt.expectedError, err)
			
			// 验证模拟对象的方法是否被调用
			if tt.expectedCalled {
				mockPermRepo.AssertCalled(t, "Create", mock.AnythingOfType("*model.Permission"))
			} else {
				mockPermRepo.AssertNotCalled(t, "Create", mock.AnythingOfType("*model.Permission"))
			}
		})
	}
}

// 测试权限更新功能
func TestPermissionService_Update(t *testing.T) {
	// 测试用例表
	tests := []struct {
		name          string
		request       *schema.UpdatePermissionRequest
		mockSetup     func(mockPermRepo *mocks.MockPermissionRepository)
		expectedError error
	}{
		{
			name: "成功更新权限",
			request: &schema.UpdatePermissionRequest{
				ID:          1,
				Name:        "更新后的权限",
				Code:        "updated:permission",
				Description: "这是一个更新后的权限",
			},
			mockSetup: func(mockPermRepo *mocks.MockPermissionRepository) {
				// 模拟权限存在
				existingPerm := &model.Permission{ID: 1, Name: "原权限", Code: "original:permission"}
				mockPermRepo.On("GetByID", uint64(1)).Return(existingPerm, nil)
				
				// 模拟权限名称不存在
				mockPermRepo.On("GetByName", "更新后的权限").Return(nil, nil)
				
				// 模拟权限标识不存在
				mockPermRepo.On("GetByCode", "updated:permission").Return(nil, nil)
				
				// 模拟更新权限成功
				mockPermRepo.On("Update", mock.AnythingOfType("*model.Permission")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "权限不存在",
			request: &schema.UpdatePermissionRequest{
				ID:          2,
				Name:        "更新后的权限",
				Code:        "updated:permission",
				Description: "这是一个更新后的权限",
			},
			mockSetup: func(mockPermRepo *mocks.MockPermissionRepository) {
				// 模拟权限不存在
				mockPermRepo.On("GetByID", uint64(2)).Return(nil, nil)
			},
			expectedError: errors.ErrPermissionNotFound,
		},
		{
			name: "权限标识已被其他权限使用",
			request: &schema.UpdatePermissionRequest{
				ID:          1,
				Name:        "更新后的权限",
				Code:        "existing:permission",
				Description: "这是一个更新后的权限",
			},
			mockSetup: func(mockPermRepo *mocks.MockPermissionRepository) {
				// 模拟权限存在
				existingPerm := &model.Permission{ID: 1, Name: "原权限", Code: "original:permission"}
				mockPermRepo.On("GetByID", uint64(1)).Return(existingPerm, nil)
				
				// 模拟权限名称不存在
				mockPermRepo.On("GetByName", "更新后的权限").Return(nil, nil)
				
				// 模拟权限标识已被其他权限使用
				otherPerm := &model.Permission{ID: 2, Code: "existing:permission"}
				mockPermRepo.On("GetByCode", "existing:permission").Return(otherPerm, nil)
				
				// 模拟更新权限失败，因为权限标识已被其他权限使用，所以我们需要修改mockSetup函数来模拟这种情况
				mockPermRepo.On("Update", mock.AnythingOfType("*model.Permission")).Return(nil)
			},
			expectedError: errors.ErrPermissionExists,
		},
		{
			name: "权限不存在",
			request: &schema.UpdatePermissionRequest{
				ID:          2,
				Name:        "更新后的权限",
				Code:        "updated:permission",
				Description: "这是一个更新后的权限",
			},
			mockSetup: func(mockPermRepo *mocks.MockPermissionRepository) {
				// 模拟权限不存在
				mockPermRepo.On("GetByID", uint64(2)).Return(nil, nil)
			},
			expectedError: errors.ErrPermissionNotFound,
		},
	}

	// 执行测试用例
	for _, tt := range tests {
		tt := tt // 防止闭包问题
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟对象
			mockPermRepo := new(mocks.MockPermissionRepository)
			
			// 设置模拟行为
			tt.mockSetup(mockPermRepo)
			
			// 创建服务实例
			permissionService := service.NewPermissionService(mockPermRepo)
			
			// 调用被测试的方法
			err := permissionService.Update(tt.request)
			
			// 验证结果
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

// 测试权限删除功能
func TestPermissionService_Delete(t *testing.T) {
	// 测试用例表
	tests := []struct {
		name          string
		permissionID  uint64
		mockSetup     func(mockPermRepo *mocks.MockPermissionRepository)
		expectedError error
	}{
		{
			name:         "成功删除权限",
			permissionID: 1,
			mockSetup: func(mockPermRepo *mocks.MockPermissionRepository) {
				// 模拟权限存在
				perm := &model.Permission{ID: 1}
				mockPermRepo.On("GetByID", uint64(1)).Return(perm, nil)
				
				// 模拟权限没有被角色使用
				mockPermRepo.On("GetRolesByPermissionID", uint64(1)).Return([]*model.Role{}, nil)
				
				// 模拟删除权限成功
				mockPermRepo.On("Delete", uint64(1)).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:         "权限不存在",
			permissionID: 2,
			mockSetup: func(mockPermRepo *mocks.MockPermissionRepository) {
				// 模拟权限不存在
				mockPermRepo.On("GetByID", uint64(2)).Return(nil, nil)
			},
			expectedError: errors.ErrPermissionNotFound,
		},
		{
			name:         "权限正在使用中",
			permissionID: 1,
			mockSetup: func(mockPermRepo *mocks.MockPermissionRepository) {
				// 模拟权限存在
				perm := &model.Permission{ID: 1}
				mockPermRepo.On("GetByID", uint64(1)).Return(perm, nil)
				
				// 模拟权限被角色使用
				perm.Roles = []model.Role{{ID: 1, Name: "角色1"}}
				
				// 模拟删除权限失败，因为权限正在使用中，所以我们需要修改mockSetup函数来模拟这种情况
				mockPermRepo.On("Delete", uint64(1)).Return(nil)
			},
			expectedError: errors.ErrPermissionInUse,
		},
	}

	// 执行测试用例
	for _, tt := range tests {
		tt := tt // 防止闭包问题
		t.Run(tt.name, func(t *testing.T) {
			// 创建模拟对象
			mockPermRepo := new(mocks.MockPermissionRepository)
			
			// 设置模拟行为
			tt.mockSetup(mockPermRepo)
			
			// 创建服务实例
			permissionService := service.NewPermissionService(mockPermRepo)
			
			// 调用被测试的方法
			err := permissionService.Delete(tt.permissionID)
			
			// 验证结果
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

// TestPermissionService_List 分页列表测试
func TestPermissionService_List(t *testing.T) {
	mockPermRepo := new(mocks.MockPermissionRepository)
	service := service.NewPermissionService(mockPermRepo)

	tests := []struct {
		name       string
		page       int
		pageSize   int
		keyword    string
		mockSetup  func()
		expectResp *schema.ListPermissionResponse
		expectErr  error
	}{
		{
			name:     "正常分页返回",
			page:     1,
			pageSize: 2,
			keyword:  "",
			mockSetup: func() {
				perms := []*model.Permission{{ID: 1, Name: "perm1"}, {ID: 2, Name: "perm2"}}
				mockPermRepo.On("List", 1, 2, "").Return(perms, int64(2), nil)
			},
			expectResp: &schema.ListPermissionResponse{
				Total:      2,
				Page:       1,
				PageSize:   2,
				TotalPages: 1,
				Items: []schema.PermissionResponse{{ID: 1, Name: "perm1"}, {ID: 2, Name: "perm2"}},
			},
			expectErr: nil,
		},
		{
			name:     "无数据",
			page:     1,
			pageSize: 10,
			keyword:  "",
			mockSetup: func() {
				mockPermRepo.On("List", 1, 10, "").Return([]*model.Permission{}, int64(0), nil)
			},
			expectResp: &schema.ListPermissionResponse{
				Total:      0,
				Page:       1,
				PageSize:   10,
				TotalPages: 0,
				Items:      []schema.PermissionResponse{},
			},
			expectErr: nil,
		},
		{
			name:     "数据库异常",
			page:     1,
			pageSize: 10,
			keyword:  "",
			mockSetup: func() {
				mockPermRepo.On("List", 1, 10, "").Return([]*model.Permission(nil), int64(0), errors.ErrDB)
			},
			expectResp: nil,
			expectErr: errors.ErrDB,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPermRepo.ExpectedCalls = nil // 清理历史
			if tt.mockSetup != nil {
				tt.mockSetup()
			}
			resp, err := service.List(&schema.ListPermissionRequest{Page: tt.page, PageSize: tt.pageSize, Keyword: tt.keyword})
			assert.Equal(t, tt.expectErr, err)
			assert.Equal(t, tt.expectResp, resp)
		})
	}
}
