#!/bin/bash
# RBAC API 测试脚本
# 作者：AI助手
# 创建时间：2025-04-13

# 设置基础URL
BASE_URL="http://localhost:8080/api/v1"
# 保存token的变量
TOKEN=""
# 保存测试用户ID
USER_ID=""
# 保存测试角色ID
ROLE_ID=""
# 保存测试权限ID
PERMISSION_ID=""

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # 无颜色

# 打印带颜色的消息
print_success() {
    echo -e "${GREEN}[成功]${NC} $1"
}

print_error() {
    echo -e "${RED}[错误]${NC} $1"
}

print_info() {
    echo -e "${YELLOW}[信息]${NC} $1"
}

# 登录并获取token
login() {
    print_info "测试登录接口..."
    response=$(curl -s -X POST "${BASE_URL}/auth/login" \
        -H "Content-Type: application/json" \
        -d '{
            "username": "admin",
            "password": "admin123"
        }')
    
    # 提取token
    TOKEN=$(echo $response | grep -o '"token":"[^"]*' | sed 's/"token":"//')
    
    if [ -n "$TOKEN" ]; then
        print_success "登录成功，获取到token: ${TOKEN:0:20}..."
    else
        print_error "登录失败: $response"
        exit 1
    fi
}

# 获取用户个人信息
get_profile() {
    print_info "测试获取个人信息接口..."
    response=$(curl -s -X GET "${BASE_URL}/auth/profile" \
        -H "Authorization: Bearer $TOKEN")
    
    echo $response | json_pp
    
    if [[ $response == *"id"* ]]; then
        print_success "获取个人信息成功"
    else
        print_error "获取个人信息失败: $response"
    fi
}

# 检查权限
check_permission() {
    local permission=$1
    print_info "测试权限检查接口 ($permission)..."
    response=$(curl -s -X POST "${BASE_URL}/auth/check-permission" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d "{
            \"permission\": \"$permission\"
        }")
    
    echo $response | json_pp
    
    if [[ $response == *"true"* ]]; then
        print_success "权限检查成功，用户拥有权限: $permission"
    else
        print_info "用户没有权限: $permission"
    fi
}

# 创建用户
create_user() {
    print_info "测试创建用户接口..."
    response=$(curl -s -X POST "${BASE_URL}/users" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d '{
            "username": "testuser",
            "email": "test@example.com",
            "password": "password123"
        }')
    
    echo $response | json_pp
    
    # 提取用户ID
    USER_ID=$(echo $response | grep -o '"id":[0-9]*' | sed 's/"id"://')
    
    if [ -n "$USER_ID" ]; then
        print_success "创建用户成功，用户ID: $USER_ID"
    else
        print_error "创建用户失败: $response"
    fi
}

# 获取用户列表
list_users() {
    print_info "测试获取用户列表接口..."
    response=$(curl -s -X GET "${BASE_URL}/users?page=1&page_size=10" \
        -H "Authorization: Bearer $TOKEN")
    
    echo $response | json_pp
    
    if [[ $response == *"items"* ]]; then
        print_success "获取用户列表成功"
    else
        print_error "获取用户列表失败: $response"
    fi
}

# 获取用户详情
get_user() {
    if [ -z "$USER_ID" ]; then
        print_error "用户ID为空，无法获取用户详情"
        return
    fi
    
    print_info "测试获取用户详情接口..."
    response=$(curl -s -X GET "${BASE_URL}/users/$USER_ID" \
        -H "Authorization: Bearer $TOKEN")
    
    echo $response | json_pp
    
    if [[ $response == *"id"* ]]; then
        print_success "获取用户详情成功"
    else
        print_error "获取用户详情失败: $response"
    fi
}

# 更新用户
update_user() {
    if [ -z "$USER_ID" ]; then
        print_error "用户ID为空，无法更新用户"
        return
    fi
    
    print_info "测试更新用户接口..."
    response=$(curl -s -X PUT "${BASE_URL}/users/$USER_ID" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d '{
            "email": "updated@example.com"
        }')
    
    echo $response | json_pp
    
    if [[ $response == *"success"* ]]; then
        print_success "更新用户成功"
    else
        print_error "更新用户失败: $response"
    fi
}

# 创建角色
create_role() {
    print_info "测试创建角色接口..."
    response=$(curl -s -X POST "${BASE_URL}/roles" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d '{
            "name": "测试角色",
            "code": "test:role",
            "description": "这是一个测试角色"
        }')
    
    echo $response | json_pp
    
    # 提取角色ID
    ROLE_ID=$(echo $response | grep -o '"id":[0-9]*' | sed 's/"id"://')
    
    if [ -n "$ROLE_ID" ]; then
        print_success "创建角色成功，角色ID: $ROLE_ID"
    else
        print_error "创建角色失败: $response"
    fi
}

# 获取角色列表
list_roles() {
    print_info "测试获取角色列表接口..."
    response=$(curl -s -X GET "${BASE_URL}/roles?page=1&page_size=10" \
        -H "Authorization: Bearer $TOKEN")
    
    echo $response | json_pp
    
    if [[ $response == *"items"* ]]; then
        print_success "获取角色列表成功"
    else
        print_error "获取角色列表失败: $response"
    fi
}

# 获取角色详情
get_role() {
    if [ -z "$ROLE_ID" ]; then
        print_error "角色ID为空，无法获取角色详情"
        return
    fi
    
    print_info "测试获取角色详情接口..."
    response=$(curl -s -X GET "${BASE_URL}/roles/$ROLE_ID" \
        -H "Authorization: Bearer $TOKEN")
    
    echo $response | json_pp
    
    if [[ $response == *"id"* ]]; then
        print_success "获取角色详情成功"
    else
        print_error "获取角色详情失败: $response"
    fi
}

# 创建权限
create_permission() {
    print_info "测试创建权限接口..."
    response=$(curl -s -X POST "${BASE_URL}/permissions" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d '{
            "name": "测试权限",
            "code": "test:permission",
            "description": "这是一个测试权限"
        }')
    
    echo $response | json_pp
    
    # 提取权限ID
    PERMISSION_ID=$(echo $response | grep -o '"id":[0-9]*' | sed 's/"id"://')
    
    if [ -n "$PERMISSION_ID" ]; then
        print_success "创建权限成功，权限ID: $PERMISSION_ID"
    else
        print_error "创建权限失败: $response"
    fi
}

# 获取权限列表
list_permissions() {
    print_info "测试获取权限列表接口..."
    response=$(curl -s -X GET "${BASE_URL}/permissions?page=1&page_size=10" \
        -H "Authorization: Bearer $TOKEN")
    
    echo $response | json_pp
    
    if [[ $response == *"items"* ]]; then
        print_success "获取权限列表成功"
    else
        print_error "获取权限列表失败: $response"
    fi
}

# 为角色分配权限
assign_permission_to_role() {
    if [ -z "$ROLE_ID" ] || [ -z "$PERMISSION_ID" ]; then
        print_error "角色ID或权限ID为空，无法分配权限"
        return
    fi
    
    print_info "测试为角色分配权限接口..."
    response=$(curl -s -X POST "${BASE_URL}/roles/$ROLE_ID/permissions" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d "{
            \"permission_ids\": [$PERMISSION_ID]
        }")
    
    echo $response | json_pp
    
    if [[ $response == *"success"* ]]; then
        print_success "为角色分配权限成功"
    else
        print_error "为角色分配权限失败: $response"
    fi
}

# 为用户分配角色
assign_role_to_user() {
    if [ -z "$USER_ID" ] || [ -z "$ROLE_ID" ]; then
        print_error "用户ID或角色ID为空，无法分配角色"
        return
    fi
    
    print_info "测试为用户分配角色接口..."
    response=$(curl -s -X POST "${BASE_URL}/users/$USER_ID/roles" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $TOKEN" \
        -d "{
            \"role_ids\": [$ROLE_ID]
        }")
    
    echo $response | json_pp
    
    if [[ $response == *"success"* ]]; then
        print_success "为用户分配角色成功"
    else
        print_error "为用户分配角色失败: $response"
    fi
}

# 删除权限
delete_permission() {
    if [ -z "$PERMISSION_ID" ]; then
        print_error "权限ID为空，无法删除权限"
        return
    fi
    
    print_info "测试删除权限接口..."
    response=$(curl -s -X DELETE "${BASE_URL}/permissions/$PERMISSION_ID" \
        -H "Authorization: Bearer $TOKEN")
    
    echo $response | json_pp
    
    if [[ $response == *"success"* ]]; then
        print_success "删除权限成功"
    else
        print_error "删除权限失败: $response"
    fi
}

# 删除角色
delete_role() {
    if [ -z "$ROLE_ID" ]; then
        print_error "角色ID为空，无法删除角色"
        return
    fi
    
    print_info "测试删除角色接口..."
    response=$(curl -s -X DELETE "${BASE_URL}/roles/$ROLE_ID" \
        -H "Authorization: Bearer $TOKEN")
    
    echo $response | json_pp
    
    if [[ $response == *"success"* ]]; then
        print_success "删除角色成功"
    else
        print_error "删除角色失败: $response"
    fi
}

# 删除用户
delete_user() {
    if [ -z "$USER_ID" ]; then
        print_error "用户ID为空，无法删除用户"
        return
    fi
    
    print_info "测试删除用户接口..."
    response=$(curl -s -X DELETE "${BASE_URL}/users/$USER_ID" \
        -H "Authorization: Bearer $TOKEN")
    
    echo $response | json_pp
    
    if [[ $response == *"success"* ]]; then
        print_success "删除用户成功"
    else
        print_error "删除用户失败: $response"
    fi
}

# 主函数
main() {
    print_info "开始测试RBAC API..."
    
    # 认证相关
    login
    get_profile
    check_permission "user:list"
    
    # 用户管理
    create_user
    list_users
    get_user
    update_user
    
    # 角色管理
    create_role
    list_roles
    get_role
    
    # 权限管理
    create_permission
    list_permissions
    
    # 关联操作
    assign_permission_to_role
    assign_role_to_user
    
    # 检查新分配的权限
    check_permission "test:permission"
    
    # 清理测试数据
    delete_user
    delete_role
    delete_permission
    
    print_success "RBAC API测试完成！"
}

# 执行主函数
main
