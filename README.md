# Fiber-RBAC Backend System

A comprehensive Role-Based Access Control (RBAC) backend system built with Go Fiber framework.

## Core Features

- **User Management**: Create, read, update, and delete users
- **Role Management**: Create, read, update, and delete roles
- **Permission Management**: Create, read, update, and delete permissions
- **Role-Permission Association**: Assign permissions to roles
- **User-Role Association**: Assign roles to users
- **JWT Authentication**: Secure API access with JWT tokens
- **Standardized API Design**: Follows OpenAPI specification with unified interface design
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
/fiber-rbac
├── cmd                 # Application entry points
├── config              # Configuration files and package
├── docs                # Documentation files (including OpenAPI spec)
├── internal            # Internal packages
│   ├── app             # Application setup (Fiber, middleware, router)
│   ├── handler         # HTTP handlers
│   ├── middleware      # Custom middleware
│   ├── model           # Data models
│   ├── pkg             # Shared packages
│   ├── repository      # Data access layer
│   ├── schema          # Request/response structures
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
git clone https://github.com/lvyunze/fiber-rbac.git
cd fiber-rbac
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
  - POST `/api/v1/auth/refresh`: Refresh token
  - POST `/api/v1/auth/profile`: Get current user information
  - POST `/api/v1/auth/check-permission`: Check permission

- **User Management**:
  - POST `/api/v1/users/list`: List users
  - POST `/api/v1/users/create`: Create user
  - POST `/api/v1/users/detail`: Get user details
  - POST `/api/v1/users/update`: Update user
  - POST `/api/v1/users/delete`: Delete user
  - POST `/api/v1/users/assign-roles`: Assign roles to user
  - POST `/api/v1/users/list-roles`: List user roles

- **Role Management**:
  - POST `/api/v1/roles/list`: List roles
  - POST `/api/v1/roles/create`: Create role
  - POST `/api/v1/roles/detail`: Get role details
  - POST `/api/v1/roles/update`: Update role
  - POST `/api/v1/roles/delete`: Delete role
  - POST `/api/v1/roles/assign-permissions`: Assign permissions to role
  - POST `/api/v1/roles/list-permissions`: List role permissions

- **Permission Management**:
  - POST `/api/v1/permissions/list`: List permissions
  - POST `/api/v1/permissions/create`: Create permission
  - POST `/api/v1/permissions/detail`: Get permission details
  - POST `/api/v1/permissions/update`: Update permission
  - POST `/api/v1/permissions/delete`: Delete permission

## API Design Features

- **Unified Request Method**: All endpoints use POST method, simplifying frontend calls
- **Request Body Parameter Passing**: All parameters are passed through request body instead of path parameters, enhancing security
- **Standardized Error Handling**: Unified error response format for easier frontend handling
- **Structured Request Validation**: Each endpoint has a corresponding request structure, ensuring consistent parameter validation

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
