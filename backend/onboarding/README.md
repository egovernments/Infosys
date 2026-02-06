# Property Tax Onboarding Service Documentation

## Table of Contents

1. [Overview](#1-overview)
2. [Features](#2-features)
3. [Architecture & Design](#3-architecture--design)
4. [Authentication & Authorization](#4-authentication--authorization)
5. [API Specification](#5-api-specification)
6. [Configuration](#6-configuration)
7. [Project Structure](#7-project-structure)
8. [Getting Started](#8-getting-started)
9. [Development Guidelines](#9-development-guidelines)
10. [API Documentation](#10-api-documentation)
11. [License](#11-license)
12. [Support](#12-support)


## 1. Overview

### Purpose
The Property Tax Onboarding Service is a Go-based microservice responsible for user profile management in the property tax system. It handles complete user lifecycle operations including registration, profile management, and user data persistence while integrating seamlessly with Keycloak for identity management.

### Scope
- Part of the **Property Tax Management System**
- Handles user onboarding and profile management
- Integrates with Keycloak for user creation and role assignment
- Provides RESTful APIs for user operations
- Manages user data persistence in PostgreSQL

## Technology Stack

- **Language**: Go 1.21+
- **Web Framework**: Gin
- **ORM**: GORM
- **Database Driver**: pgx/v5
- **Logging**: Logrus
- **Validation**: go-playground/validator
- **Configuration**: Environment variables with godotenv
- **UUID**: Google UUID library
- **Development**: Local development setup

## 2. Features

- **User Profile Management**: Complete CRUD operations for user profiles with profile and address management
- **Zone & Ward Mapping**: Assign users to multiple zones and wards using PostgreSQL array types
- **Advanced Filtering**: Filter users by role, status, email, username, phone number, and wards
- **User Analytics**: Get user counts by role and status for dashboard metrics
- **Keycloak Integration**: Seamless integration with Keycloak for user creation and role assignment
- **Automated User Activation**: Scheduled cron job for automatic user activation/deactivation based on date ranges
- **MDMS Integration**: Integration with Master Data Management Service for data validation and retrieval
- **RESTful API**: Well-structured REST endpoints following best practices
- **CORS Support**: Configurable CORS middleware for cross-origin requests
- **Database Integration**: PostgreSQL database with GORM ORM and pgx driver
- **Layered Architecture**: Clean separation of concerns with distinct layers (Handler → Service → Repository)
- **Transaction Support**: Atomic operations for complex updates using GORM transactions
- **Comprehensive Logging**: Structured logging with Logrus
- **Error Handling**: Robust error handling with custom error types and validation
- **Pagination Support**: Efficient pagination for list endpoints

## 3. Architecture & Design

### High-Level Architecture
```
┌─────────────────────────────────────────┐
│              API Gateway                │  ← Handles Auth/AuthZ
├─────────────────────────────────────────┤
│            HTTP Handlers                │  ← Request Processing
├─────────────────────────────────────────┤
│          Service Layer                  │  ← Business Logic
├─────────────────────────────────────────┤
│         Repository Layer                │  ← Data Access
├─────────────────────────────────────────┤
│          Database (PostgreSQL)          │  ← Data Persistence
└─────────────────────────────────────────┘
```

### Design Patterns

#### 1. Layered Architecture
The service is structured into distinct layers, each with a clear responsibility:
- **Handler Layer** (`internal/handlers/user_handler.go`): Handles HTTP requests, parses input, and returns responses.
- **Service Layer** (`internal/services/user_service.go`): Contains business logic, orchestrates workflows, and coordinates between handlers and repositories.
- **Repository Layer** (`internal/repositories/users_repo.go`): Encapsulates all database access and persistence logic.
- **Model Layer** (`internal/models/user.go`, `internal/models/database_models.go`): Defines data structures and DTOs for both API and database.
This separation improves maintainability, testability, and scalability.

#### 2. RESTful API Design
The service exposes endpoints that follow RESTful conventions:
- Uses standard HTTP methods (`GET`, `POST`, `PUT`, `DELETE`).
- Resource-oriented URLs (e.g., `/api/v1/users`).
- Returns appropriate HTTP status codes and standardized JSON responses (`pkg/response/response.go`).
- Supports filtering, pagination, and role-based access via query parameters and payloads.

#### 3. Repository Pattern
All database operations are abstracted behind repository interfaces (`internal/repositories/interfaces.go`), allowing:
- Easy swapping of database implementations.
- Centralized query logic.
- Simplified mocking for tests.

#### 4. Dependency Injection
Dependencies (e.g., repositories, services) are injected into handlers and services at runtime, typically in the application bootstrap (`cmd/server/main.go`). This enables:
- Loose coupling between components.
- Easier unit testing with mocks/stubs.

#### 5. Standardized API Response Pattern
All API responses use a consistent structure, defined in [`pkg/response/response.go`](pkg/response/response.go ):
- `success`: Boolean indicating operation result.
- `message`: Human-readable status.
- `data`: Payload (if any).
- `errors`: List of error messages (if any).
This ensures predictable client-side handling and easier error tracking.

#### 6. External Service Integration Pattern
Integration with Keycloak for user management is encapsulated in a dedicated service (`internal/services/keycloak_service.go`), following the "service gateway" pattern. This:
- Isolates external API logic.
- Allows for retries, error handling, and future replacement.

#### 7. Validation Pattern
Input validation is performed using the go-playground/validator library, with validation logic centralized in the validator package (`internal/validator/user_validation_service.go`). This ensures:
- Consistent validation rules across endpoints.
- Clear error messages for invalid input.

#### 8. Soft Delete Pattern
User records are soft-deleted using a `deleted_at` timestamp field, preserving data for audit and recovery (`internal/models/database_models.go`).

#### 9. Scheduler Pattern
The service includes a cron-based scheduler (`internal/scheduler/scheduler.go`) that runs automated tasks:
- **User Activation/Deactivation**: Automatically manages user status based on `startDate` and `endDate` fields.
- **Execution Frequency**: Configurable cron expression (currently runs every 2 minutes: `*/2 * * * *`).
- **Integration**: Coordinates with KeycloakService to enable/disable users in both the database and Keycloak.
- **Error Handling**: Logs failures and continues processing other users.

This pattern enables time-bound user access, useful for temporary agents or seasonal workers.

### Service Dependencies
- **Keycloak Server**: User creation, role management, and identity operations
- **PostgreSQL Database**: Primary data storage for users, profiles, and zone mappings
- **MDMS (Master Data Management Service)**: Master data validation and retrieval
- **API Gateway**: Authentication and authorization (external)

### Key Components
- **Handlers**: HTTP request/response processing (`internal/handlers/`)
- **Services**: Business logic and orchestration (`internal/services/`)
- **Repositories**: Data access layer (`internal/repositories/`)
- **Models**: Data structures and DTOs (`internal/models/`)
- **Config**: Configuration management (`internal/config/`)
- **Routes**: API route definitions (`internal/routes/`)
- **Scheduler**: Background cron jobs for automated tasks (`internal/scheduler/`)
- **Middleware**: CORS and other HTTP middleware (`internal/middleware/`)
- **Utils**: Helper functions for Keycloak and MDMS integration (`internal/utils/`)



## 4. Authentication & Authorization

This onboarding service does not require authentication or authorization. All API endpoints are accessible without tokens or credentials. No JWT, OAuth, or other authentication mechanisms are enforced for any endpoint.

## 5. API Specification

### Base URL
```
http://localhost:8080/api/v1
```

### API Endpoints

| Method | Endpoint                      | Description                                 |
|--------|-------------------------------|---------------------------------------------|
| GET    | /api/v1/users                 | Get all users with filters, pagination, and search |
| POST   | /api/v1/users                 | Create a new user (citizen, agent, service manager, commissioner) |
| GET    | /api/v1/users/count           | Get total count of users by role and status |
| GET    | /api/v1/users/{id}            | Retrieve a user by their ID                 |
| PUT    | /api/v1/users/{id}            | Update complete user information            |
| DELETE | /api/v1/users/{id}            | Soft delete user by ID                      |
| GET    | /api/v1/users/by-role/{role}  | Get users by specific role                  |
| GET    | /health                       | Service health check endpoint               |



## 6. Configuration

### Environment Variables
```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=property_tax_db
DB_SSL_MODE=disable

# Server Configuration
SERVER_PORT=8080
GIN_MODE=debug

# Keycloak Configuration
KEYCLOAK_BASE_URL=http://localhost:8080/auth
KEYCLOAK_REALM=property-tax-realm
KEYCLOAK_ADMIN_USERNAME=admin
KEYCLOAK_ADMIN_PASSWORD=admin
KEYCLOAK_CLIENT_ID=property-tax-client
KEYCLOAK_CLIENT_SECRET=your-client-secret
TOKEN_URL=http://localhost:8080/auth/realms/property-tax-realm/protocol/openid-connect/token
USER_URL=http://localhost:8080/auth/admin/realms/property-tax-realm/users
ASSIGN_ROLE_URL=http://localhost:8080/auth/admin/realms/property-tax-realm/users/{userId}/role-mappings/realm
ROLE_URL=http://localhost:8080/auth/admin/realms/property-tax-realm/roles/{roleName}
ROLES_URL=http://localhost:8080/auth/admin/realms/property-tax-realm/roles
DELETE_URL=http://localhost:8080/auth/admin/realms/property-tax-realm/users/{userId}

# CORS Configuration
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:3001
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Content-Type,Authorization

# MDMS Configuration
MDMS_API_URL=http://localhost:8081/mdms-v2

# Logging Configuration
LOG_LEVEL=info
LOG_FORMAT=json

# Application Configuration
APP_NAME=property-tax-onboarding
APP_VERSION=1.0.0
```

### Configuration Structure
The service uses a structured configuration approach defined in `internal/config/config.go` with sections for:
- **DatabaseConfig**: Database connection settings
- **ServerConfig**: Server and Gin settings
- **KeycloakConfig**: Keycloak integration settings
- **LogConfig**: Logging configuration
- **AppConfig**: Application metadata

### Sample Configuration File (.env)
```bash
# Copy this to .env and update values
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=property_tax_db
DB_SCHEMA=DIGIT3

SERVER_PORT=8080
GIN_MODE=release

KEYCLOAK_BASE_URL=http://keycloak:8080/auth
KEYCLOAK_REALM=property-tax
KEYCLOAK_ADMIN_USERNAME=admin
KEYCLOAK_ADMIN_PASSWORD=admin_password
KEYCLOAK_CLIENT_ID=property-tax-onboarding
KEYCLOAK_CLIENT_SECRET=client_secret
TOKEN_URL=http://keycloak:8080/auth/realms/property-tax/protocol/openid-connect/token
USER_URL=http://keycloak:8080/auth/admin/realms/property-tax/users
ASSIGN_ROLE_URL=http://keycloak:8080/auth/admin/realms/property-tax/users/{userId}/role-mappings/realm
ROLE_URL=http://keycloak:8080/auth/admin/realms/property-tax/roles/{roleName}
ROLES_URL=http://keycloak:8080/auth/admin/realms/property-tax/roles
DELETE_URL=http://keycloak:8080/auth/admin/realms/property-tax/users/{userId}

CORS_ALLOWED_ORIGINS=http://localhost:3000
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Content-Type,Authorization

MDMS_API_URL=http://localhost:8081/mdms-v2

LOG_LEVEL=info
LOG_FORMAT=json
```

## 7. Project Structure

```
property-tax-onboarding/
├── .gitignore
├── Dockerfile
├── go.mod
├── README.md
├── service-documentation 2.md
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── constants/
│   │   ├── errors_constants.go
│   │   └── logs_constants.go
│   ├── database/
│   │   └── database.go
│   ├── errors/
│   │   └── errors.go
│   ├── handlers/
│   │   └── user_handler.go
│   ├── middleware/
│   │   └── cors.go
│   ├── models/
│   │   ├── user.go
│   │   └── database_models.go
│   ├── repositories/
│   │   ├── users_repo.go
│   │   ├── transaction.go
│   │   └── ...
│   ├── routes/
│   │   └── routes.go
│   ├── scheduler/
│   │   └── scheduler.go
│   ├── services/
│   │   ├── user_service.go
│   │   └── keycloak_service.go
│   ├── utils/
│   │   ├── mdms_utils.go
│   │   └── keycloak_utils.go
│   └── validator/
│       └── user_validation_service.go
├── migrations/
│   └── README.md
├── pkg/
│   ├── logger/
│   │   └── logger.go
│   └── response/
│       └── response.go
├── postman/
│   └── new_onboarding.postman_collection.json
```

- `cmd/server/main.go`: Application entry point
- `internal/handlers/`: HTTP request handlers
- `internal/services/`: Business logic and orchestration
- `internal/repositories/`: Data access layer
- `internal/models/`: Data models and DTOs
- `internal/config/`: Configuration management
- `internal/routes/`: API route definitions
- `internal/scheduler/`: Cron job scheduler for automated user activation/deactivation
- `internal/middleware/`: Custom middleware (e.g., CORS)
- `internal/utils/`: Utility functions (e.g., MDMS, Keycloak helpers)
- `internal/validator/`: Input validation logic
- `pkg/logger/`: Logging utilities
- `pkg/response/`: Standardized API responses
- `migrations/`: Database migration scripts and docs
- `postman/`: Postman collections for API testing

For more details, see the respective files and folders.

## 8. Getting Started

Follow these steps to set up and run the Property Tax Onboarding Service for local development or in a containerized environment.

### Prerequisites
- Go 1.21 or higher
- PostgreSQL 12 or higher
- Docker (optional, for containerized setup)
- Keycloak server (for user identity management)

### Local Development Setup
1. **Clone the repository**
	```bash
	git clone <repository-url>
	cd property-tax-onboarding
	```
2. **Install Go dependencies**
	```bash
	go mod download
	```
3. **Configure environment variables**
	Copy the example environment file and update values as needed:
	```bash
	cp .env.example .env
	# Edit .env with your database and service URLs
	```
4. **Set up the PostgreSQL database**
	- Ensure PostgreSQL is running and accessible as per your .env configuration.
	- Create the database and user if not already present.
5. **Run database migrations**
	Apply schema and seed data:
	```bash
	go run migrations/*.go
	# Or use a migration tool if provided
	```
6. **Start the service**
	```bash
	go run cmd/server/main.go
	```
	The service will be available at http://localhost:8080

## 9. Development Guidelines

### Code Style
- Follow Go best practices and idiomatic Go style (gofmt, golint).
- Use clear, descriptive names for variables, functions, and files.
- Keep functions small and focused on a single responsibility.
- Use comments to explain complex logic or business rules.

### Project Structure
- Place new features in the appropriate layer (handler, service, repository, model).
- Avoid business logic in handlers; keep it in services.
- Use interfaces for repositories and services to enable mocking and testing.

### Version Control
- Use feature branches for new work (e.g., feature/your-feature-name).
- Write clear, concise commit messages.
- Rebase and squash commits before merging to main/develop.

### Testing
- Write unit tests for all business logic and data access code.
- Use table-driven tests for Go functions.
- Mock external dependencies (database, Keycloak) in tests.
- Run tests locally before pushing changes.

### Pull Requests
- Ensure all tests pass before submitting a PR.
- Request reviews from at least one other team member.
- Address review comments promptly.

### Documentation
- Update documentation (README, service-documentation) for new features or changes.
- Document public APIs and exported functions.

### Security & Secrets
- Never commit secrets, passwords, or sensitive data to the repository.
- Use environment variables or secret management tools for credentials.

### Issue Tracking
- Use the issue tracker to report bugs, request features, or track tasks.

---
For more details, refer to the team wiki or contact the maintainers.

## 10. API Documentation

Deployed Swagger UI — View and interact with the live API documentation:

## 11. License

This project is licensed under the MIT License - see the LICENSE file for details.

## 12. Support

For support and questions:
- Create an issue in the GitHub repository
- Contact the development team
- Check the API documentation in `/docs`