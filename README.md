# Fiber RBAC Service

[中文版](README_CN.md) | English

## Introduction
This is a Role-Based Access Control (RBAC) service built using Go Fiber. It supports SQLite, MySQL, and PostgreSQL databases.

## Features
- User management
- Role management
- Permission management
- JWT authentication and token refresh
- User registration and login
- Support for multiple databases (SQLite, MySQL, PostgreSQL)
- IP access control (whitelist/blacklist)
- Unified API response format

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
Edit the `config.yaml` file to configure your database settings, IP restrictions, and JWT secret:

```yaml
server:
  port: 3000
  jwt_secret: your-jwt-secret-key # JWT secret key for authentication

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

### Authentication
- `POST /api/v1/auth/register`: Register a new user
- `POST /api/v1/auth/login`: Login with username and password
- `POST /api/v1/auth/refresh`: Refresh JWT token

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

## Authentication Flow

### Registration
To register a new user, send a POST request to `/api/v1/auth/register` with the following JSON body:
```json
{
  "username": "john_doe",
  "password": "secure_password",
  "email": "john@example.com"
}
```

Upon successful registration, the API returns a JWT token which can be used for authentication.

### Login
To login, send a POST request to `/api/v1/auth/login` with:
```json
{
  "username": "john_doe",
  "password": "secure_password"
}
```

### Token Refresh
To refresh an expired JWT token, send a POST request to `/api/v1/auth/refresh` with:
```json
{
  "token": "your-expired-or-valid-jwt-token"
}
```

### Using JWT in Requests
For protected endpoints, include the JWT token in the Authorization header:
```
Authorization: Bearer your-jwt-token
```

## Response Format
All API responses follow a unified format:
```json
{
  "status": "success|error",
  "code": 1000,
  "message": "operation successful",
  "data": { ... }
}
```

Status codes:
- 1000: Success
- 1001-1099: General errors
- 1100-1199: User-related errors
- 1200-1299: Role-related errors
- 1300-1399: Permission-related errors

## Project Structure
```
fiber-rbac/
│
├── api/
│   └── v1/
│       ├── auth_handlers.go
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
│   ├── api/
│   │   └── v1/
│   │       └── auth_handlers_test.go
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
│   ├── repository/
│   │   ├── user_repository_test.go
│   │   ├── role_repository_test.go
│   │   └── permission_repository_test.go
│   └── utils/
│       └── jwt_test.go
│
├── main.go
├── go.mod
├── go.sum
├── config.yaml
├── README.md
└── README_CN.md
```

## License
This project is licensed under the MIT License.

## Upcoming Features
- ~~JWT authentication implementation~~ (Implemented)
- User-Role association management
- Role-Permission association management
- ~~Request rate limiting~~ (IP access control implemented)
- Logging functionality
- ~~Unit and integration tests~~ (Comprehensive tests implemented)
- Docker containerization

## Contributing
Contributions are welcome! If you'd like to contribute, please:
1. Fork the repository
2. Create a feature branch
3. Submit a pull request

Feel free to suggest new features or improvements by opening an issue.