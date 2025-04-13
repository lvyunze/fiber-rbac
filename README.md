# RBAC Backend System

A comprehensive Role-Based Access Control (RBAC) backend system built with Go Fiber framework.

[中文文档](./README_CN.md)

## Features

- **User Management**: Create, read, update, and delete users
- **Role Management**: Create, read, update, and delete roles
- **Permission Management**: Create, read, update, and delete permissions
- **Role-Permission Association**: Assign permissions to roles
- **User-Role Association**: Assign roles to users
- **JWT Authentication**: Secure API access with JWT tokens
- **RESTful API**: Follows RESTful API design principles
- **Structured Logging**: Uses Go's built-in slog for structured logging
- **Configuration Management**: Uses Viper for configuration management

## Tech Stack

- **Language**: Go 1.22+
- **Framework**: Fiber
- **ORM**: GORM
- **Database**: PostgreSQL
- **Authentication**: JWT
- **Logging**: Standard library slog
- **Configuration**: Viper + YAML

## Project Structure

```
/rbac
├── cmd                 # Application entry points
├── config              # Configuration files and package
├── docs                # Documentation files
├── internal            # Internal packages
│   ├── app             # Application setup (Fiber, middleware, router)
│   ├── handler         # HTTP handlers
│   ├── middleware      # Custom middleware
│   ├── model           # Data models
│   ├── pkg             # Shared packages
│   ├── repository      # Data access layer
│   └── service         # Business logic layer
├── scripts             # Utility scripts
└── test                # Test files and utilities
```

## Getting Started

### Prerequisites

- Go 1.22 or higher
- PostgreSQL
- Git

### Installation

1. Clone the repository:

```bash
git clone https://github.com/lvyunze/rbac.git
cd rbac
```

2. Install dependencies:

```bash
go mod download
```

3. Configure the application:

Edit the `config/config.yaml` file to match your environment settings.

4. Initialize the database:

```bash
go run scripts/init_db.go
```

5. Run the application:

```bash
go run cmd/main.go
```

### Default Credentials

- **Username**: admin
- **Password**: admin123
- **Email**: admin@example.com

## API Documentation

The API follows the OpenAPI specification. You can find the API documentation in the `docs/rbac.openapi.json` file.

### Key Endpoints

- **Authentication**:
  - POST `/api/v1/auth/login`: User login
  - GET `/api/v1/auth/me`: Get current user information

- **Users**:
  - GET `/api/v1/users`: List users
  - POST `/api/v1/users`: Create user
  - GET `/api/v1/users/:id`: Get user details
  - PUT `/api/v1/users/:id`: Update user
  - DELETE `/api/v1/users/:id`: Delete user

- **Roles**:
  - GET `/api/v1/roles`: List roles
  - POST `/api/v1/roles`: Create role
  - GET `/api/v1/roles/:id`: Get role details
  - PUT `/api/v1/roles/:id`: Update role
  - DELETE `/api/v1/roles/:id`: Delete role

- **Permissions**:
  - GET `/api/v1/permissions`: List permissions
  - POST `/api/v1/permissions`: Create permission
  - GET `/api/v1/permissions/:id`: Get permission details
  - PUT `/api/v1/permissions/:id`: Update permission
  - DELETE `/api/v1/permissions/:id`: Delete permission

- **Role-Permission Management**:
  - POST `/api/v1/roles/:id/permissions`: Assign permissions to role

- **User-Role Management**:
  - POST `/api/v1/users/:id/roles`: Assign roles to user

## Testing

Run the automated tests:

```bash
go test ./...
```

For API testing, use the provided test script:

```bash
bash test/api_test.sh
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.
