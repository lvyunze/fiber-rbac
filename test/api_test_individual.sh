#!/bin/bash
# RBAC API 独立测试命令脚本
# 作者：AI助手
# 创建时间：2025-04-13

# 设置基础URL
BASE_URL="http://localhost:8080/api/v1"

# 登录并获取token
echo "# 登录接口"
echo "curl -X POST \"${BASE_URL}/auth/login\" \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -d '{\"username\": \"admin\", \"password\": \"admin123\"}'"
echo ""

# 假设我们已经获取到了token
TOKEN="这里替换为实际的token"

echo "# 获取用户个人信息"
echo "curl -X GET \"${BASE_URL}/auth/profile\" \\"
echo "  -H \"Authorization: Bearer ${TOKEN}\""
echo ""

echo "# 检查权限"
echo "curl -X POST \"${BASE_URL}/auth/check-permission\" \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -H \"Authorization: Bearer ${TOKEN}\" \\"
echo "  -d '{\"permission\": \"user:list\"}'"
echo ""

echo "# 创建用户"
echo "curl -X POST \"${BASE_URL}/users\" \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -H \"Authorization: Bearer ${TOKEN}\" \\"
echo "  -d '{\"username\": \"testuser\", \"email\": \"test@example.com\", \"password\": \"password123\"}'"
echo ""

echo "# 获取用户列表"
echo "curl -X GET \"${BASE_URL}/users?page=1&page_size=10\" \\"
echo "  -H \"Authorization: Bearer ${TOKEN}\""
echo ""

echo "# 获取用户详情（替换{user_id}为实际ID）"
echo "curl -X GET \"${BASE_URL}/users/{user_id}\" \\"
echo "  -H \"Authorization: Bearer ${TOKEN}\""
echo ""

echo "# 更新用户信息（替换{user_id}为实际ID）"
echo "curl -X PUT \"${BASE_URL}/users/{user_id}\" \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -H \"Authorization: Bearer ${TOKEN}\" \\"
echo "  -d '{\"email\": \"updated@example.com\"}'"
echo ""

echo "# 创建角色"
echo "curl -X POST \"${BASE_URL}/roles\" \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -H \"Authorization: Bearer ${TOKEN}\" \\"
echo "  -d '{\"name\": \"测试角色\", \"code\": \"test:role\", \"description\": \"这是一个测试角色\"}'"
echo ""

echo "# 获取角色列表"
echo "curl -X GET \"${BASE_URL}/roles?page=1&page_size=10\" \\"
echo "  -H \"Authorization: Bearer ${TOKEN}\""
echo ""

echo "# 获取角色详情（替换{role_id}为实际ID）"
echo "curl -X GET \"${BASE_URL}/roles/{role_id}\" \\"
echo "  -H \"Authorization: Bearer ${TOKEN}\""
echo ""

echo "# 创建权限"
echo "curl -X POST \"${BASE_URL}/permissions\" \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -H \"Authorization: Bearer ${TOKEN}\" \\"
echo "  -d '{\"name\": \"测试权限\", \"code\": \"test:permission\", \"description\": \"这是一个测试权限\"}'"
echo ""

echo "# 获取权限列表"
echo "curl -X GET \"${BASE_URL}/permissions?page=1&page_size=10\" \\"
echo "  -H \"Authorization: Bearer ${TOKEN}\""
echo ""

echo "# 为角色分配权限（替换{role_id}和{permission_id}为实际ID）"
echo "curl -X POST \"${BASE_URL}/roles/{role_id}/permissions\" \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -H \"Authorization: Bearer ${TOKEN}\" \\"
echo "  -d '{\"permission_ids\": [{permission_id}]}'"
echo ""

echo "# 为用户分配角色（替换{user_id}和{role_id}为实际ID）"
echo "curl -X POST \"${BASE_URL}/users/{user_id}/roles\" \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -H \"Authorization: Bearer ${TOKEN}\" \\"
echo "  -d '{\"role_ids\": [{role_id}]}'"
echo ""

echo "# 删除权限（替换{permission_id}为实际ID）"
echo "curl -X DELETE \"${BASE_URL}/permissions/{permission_id}\" \\"
echo "  -H \"Authorization: Bearer ${TOKEN}\""
echo ""

echo "# 删除角色（替换{role_id}为实际ID）"
echo "curl -X DELETE \"${BASE_URL}/roles/{role_id}\" \\"
echo "  -H \"Authorization: Bearer ${TOKEN}\""
echo ""

echo "# 删除用户（替换{user_id}为实际ID）"
echo "curl -X DELETE \"${BASE_URL}/users/{user_id}\" \\"
echo "  -H \"Authorization: Bearer ${TOKEN}\""
echo ""
