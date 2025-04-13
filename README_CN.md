# Fiber-RBAC 后端系统

基于 Go Fiber 框架构建的完整角色权限控制（RBAC）后端系统。

[English Documentation](./README.md)

## 核心特性

- **用户管理**：创建、查询、更新和删除用户
- **角色管理**：创建、查询、更新和删除角色
- **权限管理**：创建、查询、更新和删除权限
- **角色-权限关联**：为角色分配权限
- **用户-角色关联**：为用户分配角色
- **JWT认证**：使用 JWT 令牌保障 API 访问安全
- **标准化API设计**：遵循 OpenAPI 规范的统一接口设计
- **结构化日志**：使用 Go 内置的 slog 进行结构化日志记录
- **配置管理**：使用 Viper 进行配置管理

## 技术栈

- **语言**：Go 1.22+
- **框架**：Fiber
- **ORM**：GORM
- **数据库**：PostgreSQL
- **认证**：JWT
- **日志**：标准库 slog
- **配置**：Viper + YAML

## 项目结构

```
/fiber-rbac
├── cmd                 # 应用入口点
├── config              # 配置文件和包
├── docs                # 文档文件（包含 OpenAPI 规范）
├── internal            # 内部包
│   ├── app             # 应用设置（Fiber、中间件、路由）
│   ├── handler         # HTTP 处理器
│   ├── middleware      # 自定义中间件
│   ├── model           # 数据模型
│   ├── pkg             # 共享包
│   ├── repository      # 数据访问层
│   ├── schema          # 请求/响应结构体
│   └── service         # 业务逻辑层
├── scripts             # 实用脚本
└── test                # 测试文件和工具
```

## 快速开始

### 前置条件

- Go 1.22 或更高版本
- PostgreSQL
- Git

### 安装步骤

1. 克隆仓库：

```bash
git clone https://github.com/lvyunze/fiber-rbac.git
cd fiber-rbac
```

2. 安装依赖：

```bash
go mod download
```

3. 配置应用：

编辑 `config/config.yaml` 文件以匹配您的环境设置。

4. 初始化数据库：

```bash
go run scripts/init_db.go
```

5. 运行应用：

```bash
go run cmd/main.go
```

### 默认凭据

- **用户名**：admin
- **密码**：admin123
- **邮箱**：admin@example.com

## API 文档

API 遵循 OpenAPI 规范。您可以在 `docs/rbac.openapi.json` 文件中找到 API 文档。

### 核心接口

- **认证**：
  - POST `/api/v1/auth/login`：用户登录
  - POST `/api/v1/auth/refresh`：刷新令牌
  - POST `/api/v1/auth/profile`：获取当前用户信息
  - POST `/api/v1/auth/check-permission`：检查权限

- **用户管理**：
  - POST `/api/v1/users/list`：列出用户
  - POST `/api/v1/users/create`：创建用户
  - POST `/api/v1/users/detail`：获取用户详情
  - POST `/api/v1/users/update`：更新用户
  - POST `/api/v1/users/delete`：删除用户
  - POST `/api/v1/users/assign-roles`：为用户分配角色
  - POST `/api/v1/users/list-roles`：列出用户角色

- **角色管理**：
  - POST `/api/v1/roles/list`：列出角色
  - POST `/api/v1/roles/create`：创建角色
  - POST `/api/v1/roles/detail`：获取角色详情
  - POST `/api/v1/roles/update`：更新角色
  - POST `/api/v1/roles/delete`：删除角色
  - POST `/api/v1/roles/assign-permissions`：为角色分配权限
  - POST `/api/v1/roles/list-permissions`：列出角色权限

- **权限管理**：
  - POST `/api/v1/permissions/list`：列出权限
  - POST `/api/v1/permissions/create`：创建权限
  - POST `/api/v1/permissions/detail`：获取权限详情
  - POST `/api/v1/permissions/update`：更新权限
  - POST `/api/v1/permissions/delete`：删除权限

## API 设计特点

- **统一的请求方法**：所有接口均使用 POST 方法，简化前端调用
- **请求体参数传递**：所有参数通过请求体传递，而非路径参数，提高安全性
- **标准化错误处理**：统一的错误响应格式，便于前端处理
- **结构化请求验证**：每个接口都有对应的请求结构体，确保参数验证一致性

## 测试

运行自动化测试：

```bash
go test ./...
```

对于 API 测试，使用提供的测试脚本：

```bash
bash test/api_test.sh
```

## 贡献

1. Fork 仓库
2. 创建功能分支 (`git checkout -b feature/amazing-feature`)
3. 提交更改 (`git commit -m '添加一些惊人的功能'`)
4. 推送到分支 (`git push origin feature/amazing-feature`)
5. 打开 Pull Request

## 许可证

本项目采用 MIT 许可证 - 详情请参阅 LICENSE 文件。
