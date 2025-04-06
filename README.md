# Fiber RBAC Service

[中文版](README_CN.md) | English

## Introduction
This is a Role-Based Access Control (RBAC) service built using Go Fiber. It supports SQLite, MySQL, and PostgreSQL databases.

## Features
- User management (CRUD operations)
- Role management (CRUD operations)
- Permission management (CRUD operations)
- JWT authentication
- Support for multiple databases (SQLite, MySQL, PostgreSQL)

## Getting Started

### Prerequisites
- Go 1.22 or later
- SQLite, MySQL, or PostgreSQL

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/lvyunze/fiber-rbac.git
   ```
2. Navigate to the project directory:
   ```bash
   cd fiber-rbac
   ```
3. Install dependencies:
   ```bash
   go mod tidy
   ```

### Configuration
Edit the `config.yaml` file to configure your database settings:

```yaml
db:
  driver: sqlite  # Options: sqlite, mysql, postgres
  host: localhost
  port: 3306
  user: root
  password: password
  name: rbac
```

### Running the Service
Start the server:
```bash
go run main.go
```

## API Endpoints

### Users
- `POST /api/v1/users`: Create a new user
- `GET /api/v1/users`: Get all users
- `GET /api/v1/users/:id`: Get a user by ID
- `PUT /api/v1/users/:id`: Update a user
- `DELETE /api/v1/users/:id`: Delete a user

### Roles
- `POST /api/v1/roles`: Create a new role
- `GET /api/v1/roles`: Get all roles
- `GET /api/v1/roles/:id`: Get a role by ID
- `PUT /api/v1/roles/:id`: Update a role
- `DELETE /api/v1/roles/:id`: Delete a role

### Permissions
- `POST /api/v1/permissions`: Create a new permission
- `GET /api/v1/permissions`: Get all permissions
- `GET /api/v1/permissions/:id`: Get a permission by ID
- `PUT /api/v1/permissions/:id`: Update a permission
- `DELETE /api/v1/permissions/:id`: Delete a permission

## Project Structure
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
│   │   └── auth.go
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
│       └── jwt.go
│
├── main.go
├── go.mod
├── README.md
└── README_CN.md
```

## License
This project is licensed under the MIT License.

---

# Fiber RBAC 服务

## 介绍
这是一个使用 Go Fiber 构建的基于角色的访问控制（RBAC）服务。支持 SQLite、MySQL 和 PostgreSQL 数据库。

## 功能
- 用户管理
- 角色管理
- 权限管理
- JWT 认证

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
编辑 `config.yaml` 文件以配置数据库设置。

### 运行服务
启动服务器：
```bash
   go run main.go
```

## API 端点
- `POST /api/v1/users`: 创建新用户

## 许可证
此项目使用 MIT 许可证。 