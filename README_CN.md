# Fiber RBAC 服务

中文 | [English](README.md)

## 介绍
这是一个使用 Go Fiber 构建的基于角色的访问控制（RBAC）服务。支持 SQLite、MySQL 和 PostgreSQL 数据库。

## 功能
- 用户管理
- 角色管理
- 权限管理
- JWT 认证
- 支持多种数据库（SQLite、MySQL、PostgreSQL）
- IP 访问控制（白名单/黑名单）

## 快速开始

### 前提条件
- Go 1.22 或更高版本
- SQLite、MySQL 或 PostgreSQL

### 安装
1. 克隆仓库：
   ```bash
   git clone https://github.com/lvyunze/fiber-rbac.git
   ```
2. 进入项目目录：
   ```bash
   cd fiber-rbac
   ```
3. 安装依赖：
   ```bash
   go mod tidy
   ```

### 配置
编辑 `config.yaml` 文件以配置数据库设置和IP限制：

```yaml
db:
  driver: sqlite  # 选项: sqlite, mysql, postgres
  host: localhost
  port: 3306
  user: root
  password: password
  name: rbac

# IP限制配置
ip_limit:
  enabled: true  # 是否启用IP限制
  whitelist_mode: false  # true: 只允许白名单中的IP访问, false: 阻止黑名单中的IP访问
  
  # IP白名单列表
  whitelist:
    - 127.0.0.1
    - 192.168.1.100
    
  # IP黑名单列表
  blacklist:
    - 10.0.0.5
    - 192.168.1.200
  
  # 允许的IP网段 (CIDR格式)
  allowed_networks:
    - 192.168.0.0/16  # 192.168.0.0 - 192.168.255.255
    - 10.0.0.0/8      # 10.0.0.0 - 10.255.255.255
```

### 运行服务
启动服务器：
```bash
go run main.go
```

## API 端点

### 用户
- `POST /api/v1/users`：创建新用户
- `GET /api/v1/users`：获取所有用户
- `GET /api/v1/users/:id`：根据 ID 获取用户
- `PUT /api/v1/users/:id`：更新用户
- `DELETE /api/v1/users/:id`：删除用户

### 角色
- `POST /api/v1/roles`：创建新角色
- `GET /api/v1/roles`：获取所有角色
- `GET /api/v1/roles/:id`：根据 ID 获取角色
- `PUT /api/v1/roles/:id`：更新角色
- `DELETE /api/v1/roles/:id`：删除角色

### 权限
- `POST /api/v1/permissions`：创建新权限
- `GET /api/v1/permissions`：获取所有权限
- `GET /api/v1/permissions/:id`：根据 ID 获取权限
- `PUT /api/v1/permissions/:id`：更新权限
- `DELETE /api/v1/permissions/:id`：删除权限

## 项目结构
```
fiber-rbac/
│
├── api/
│   └── v1/
│       ├── user_handlers.go
│       ├── role_handlers.go
│       └── permission_handlers.go
│
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── middleware/
│   │   ├── auth.go
│   │   └── ip_limit.go
│   ├── models/
│   │   ├── user.go
│   │   ├── role.go
│   │   └── permission.go
│   ├── repository/
│   │   ├── user_repository.go
│   │   ├── role_repository.go
│   │   └── permission_repository.go
│   ├── service/
│   │   ├── user_service.go
│   │   ├── role_service.go
│   │   └── permission_service.go
│   └── utils/
│       ├── jwt.go
│       └── response.go
│
├── test/
│   ├── handler/
│   │   ├── user_handler_test.go
│   │   ├── role_handler_test.go
│   │   └── permission_handler_test.go
│   ├── middleware/
│   │   └── ip_limit_test.go
│   ├── service/
│   │   ├── user_service_test.go
│   │   ├── role_service_test.go
│   │   └── permission_service_test.go
│   └── repository/
│       ├── user_repository_test.go
│       ├── role_repository_test.go
│       └── permission_repository_test.go
│
├── main.go
├── go.mod
├── config.yaml
├── README.md
└── README_CN.md
```

## 许可证
此项目使用 MIT 许可证。 

## 待实现功能
- JWT认证机制
- 用户-角色关联管理
- 角色-权限关联管理
- 日志记录功能
- Docker容器化部署

## 贡献
欢迎贡献！如果您想参与贡献，请：
1. Fork 本仓库
2. 创建特性分支
3. 提交 Pull Request

也欢迎通过创建 Issue 来提出新功能建议或改进意见。