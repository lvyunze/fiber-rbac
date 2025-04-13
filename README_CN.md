# RBAC 后端系统

基于 Go Fiber 框架构建的全面角色基础访问控制（RBAC）后端系统。

[English Documentation](./README.md)

## 功能特点

- **用户管理**：创建、查询、更新和删除用户
- **角色管理**：创建、查询、更新和删除角色
- **权限管理**：创建、查询、更新和删除权限
- **角色-权限关联**：为角色分配权限
- **用户-角色关联**：为用户分配角色
- **JWT 认证**：使用 JWT 令牌进行安全的 API 访问
- **RESTful API**：遵循 RESTful API 设计原则
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
/rbac
├── cmd                 # 应用程序入口点
├── config              # 配置文件和包
├── docs                # 文档文件
├── internal            # 内部包
│   ├── app             # 应用程序设置（Fiber、中间件、路由）
│   ├── handler         # HTTP 处理器
│   ├── middleware      # 自定义中间件
│   ├── model           # 数据模型
│   ├── pkg             # 共享包
│   ├── repository      # 数据访问层
│   └── service         # 业务逻辑层
├── scripts             # 实用脚本
└── test                # 测试文件和工具
```

## 快速开始

### 前提条件

- Go 1.22 或更高版本
- PostgreSQL
- Git

### 安装

1. 克隆仓库：

```bash
git clone https://github.com/lvyunze/rbac.git
cd rbac
```

2. 安装依赖：

```bash
go mod download
```

3. 配置应用程序：

编辑 `config/config.yaml` 文件以匹配您的环境设置。

4. 初始化数据库：

```bash
go run scripts/init_db.go
```

5. 运行应用程序：

```bash
go run cmd/main.go
```

### 默认凭据

- **用户名**：admin
- **密码**：admin123
- **邮箱**：admin@example.com

## API 文档

API 遵循 OpenAPI 规范。您可以在 `docs/rbac.openapi.json` 文件中找到 API 文档。

### 主要端点

- **认证**：
  - POST `/api/v1/auth/login`：用户登录
  - GET `/api/v1/auth/me`：获取当前用户信息

- **用户**：
  - GET `/api/v1/users`：列出用户
  - POST `/api/v1/users`：创建用户
  - GET `/api/v1/users/:id`：获取用户详情
  - PUT `/api/v1/users/:id`：更新用户
  - DELETE `/api/v1/users/:id`：删除用户

- **角色**：
  - GET `/api/v1/roles`：列出角色
  - POST `/api/v1/roles`：创建角色
  - GET `/api/v1/roles/:id`：获取角色详情
  - PUT `/api/v1/roles/:id`：更新角色
  - DELETE `/api/v1/roles/:id`：删除角色

- **权限**：
  - GET `/api/v1/permissions`：列出权限
  - POST `/api/v1/permissions`：创建权限
  - GET `/api/v1/permissions/:id`：获取权限详情
  - PUT `/api/v1/permissions/:id`：更新权限
  - DELETE `/api/v1/permissions/:id`：删除权限

- **角色-权限管理**：
  - POST `/api/v1/roles/:id/permissions`：为角色分配权限

- **用户-角色管理**：
  - POST `/api/v1/users/:id/roles`：为用户分配角色

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
2. 创建您的功能分支（`git checkout -b feature/amazing-feature`）
3. 提交您的更改（`git commit -m 'Add some amazing feature'`）
4. 推送到分支（`git push origin feature/amazing-feature`）
5. 打开一个 Pull Request

## 许可证

本项目采用 MIT 许可证 - 详情请参阅 LICENSE 文件
