# Fiber RBAC Service

[中文版](README_CN.md) | English

## Introduction
This is a Role-Based Access Control (RBAC) service built using Go Fiber. It supports SQLite, MySQL, and PostgreSQL databases.

## Features
- User management
- Role management
- Permission management
- JWT authentication
- Support for multiple databases (SQLite, MySQL, PostgreSQL)
- IP access control (whitelist/blacklist)

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
Edit the `config.yaml` file to configure your database settings and IP restrictions:

```yaml
db:
  driver: sqlite  # Options: sqlite, mysql, postgres
  host: localhost
  port: 3306
  user: root
  password: password
  name: rbac

# IP restriction configuration
ip_limit:
  enabled: true  # Whether to enable IP restrictions
  whitelist_mode: false  # true: only allow IPs in whitelist, false: block IPs in blacklist
  
  # IP whitelist
  whitelist:
    - 127.0.0.1
    - 192.168.1.100
    
  # IP blacklist
  blacklist:
    - 10.0.0.5
    - 192.168.1.200
  
  # Allowed IP networks (CIDR format)
  allowed_networks:
    - 192.168.0.0/16  # 192.168.0.0 - 192.168.255.255
    - 10.0.0.0/8      # 10.0.0.0 - 10.255.255.255
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

## License
This project is licensed under the MIT License.

## Upcoming Features
- JWT authentication implementation
- User-Role association management
- Role-Permission association management
- ~~Request rate limiting~~ (IP access control implemented)
- Logging functionality
- ~~Unit and integration tests~~ (Basic tests implemented)
- Docker containerization

## Contributing
Contributions are welcome! If you'd like to contribute, please:
1. Fork the repository
2. Create a feature branch
3. Submit a pull request

Feel free to suggest new features or improvements by opening an issue.