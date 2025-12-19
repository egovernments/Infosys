# Property Tax Enumeration Service

A robust Go microservice for managing property tax enumeration, including applications, properties, owners, GIS data, amenities, documents, and workflow integration. Part of the DIGIT3 platform.

---

## Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [API Overview](#api-overview)
- [Authentication & Authorization](#authentication--authorization)
- [Getting Started](#getting-started)
- [Configuration](#configuration)
- [Development Guidelines](#development-guidelines)
- [API Documentation](#api-documentation)

---

## Overview

The Property Tax Enumeration Service manages the end-to-end lifecycle of property tax applications, property and owner records, GIS and amenity data, and document management. It is designed for seamless integration with workflow, ID generation, and file storage microservices, providing a robust backend for the Property Tax Management System.

The service is extensible and supports integration with additional modules such as notification, audit logging, and master data management (MDMS). It is built to handle high transaction volumes and supports multi-tenancy, role-based access control, and flexible business rules. The architecture enables rapid onboarding of new property types, owner categories, and workflow states, making it suitable for evolving municipal and state-level property tax requirements.

---


## Features

### Business Features
- End-to-end property application lifecycle management (initiate, assign, verify, approve, reject, resubmit)
- Multi-owner property support with share validation and joint ownership
- Comprehensive user management with roles: CITIZEN, AGENT, SERVICE_MANAGER, COMMISSIONER
- GIS and coordinates management for spatial property data
- Amenity and construction details tracking for each property
- Document upload, verification, and audit logging for compliance
- Workflow integration for application status transitions and approvals
- Flexible business rules for property types, owner categories, and workflow states
- Multi-tenancy and support for evolving municipal/state requirements

### Technical Features
- RESTful APIs with OpenAPI 3.0 specification
- JWT-based authentication and role-based authorization
- MDMS integration for dynamic role and permission management
- PostgreSQL database with GORM ORM
- Gin framework for high-performance HTTP routing
- Docker support for containerized deployment
- Structured and configurable logging
- Health check endpoints for monitoring
- Pagination, filtering, and search for all major endpoints

---


## Architecture & Design

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
- **Layered Architecture**: Used throughout the codebase—handlers (`internal/handlers/`), services (`internal/services/`), repositories (`internal/repositories/`), and models (`internal/models/`). Each layer has a clear responsibility, improving maintainability and testability.
- **RESTful API Design**: All endpoints (see `internal/routes/routes.go`) use standard HTTP methods and resource-oriented URLs, returning consistent JSON responses. Pagination, filtering, and role-based access are supported via query parameters.
- **Repository Pattern**: Database access is abstracted behind interfaces in `internal/repositories/interfaces.go`. Each entity has a dedicated repository (e.g., `application_repository.go`), making DB logic modular and testable.
- **Dependency Injection**: Dependencies are injected at runtime in `cmd/server/main.go`, allowing handlers and services to receive repositories and other services as parameters. This enables easy mocking and unit testing.
- **Standardized API Response**: All API responses use a consistent structure from `pkg/response/response.go` (`success`, `message`, `data`, `errors`), ensuring predictable client handling and error tracking.
- **External Service Integration**: Clients for workflow, IDGen, MDMS, and notification are implemented in `internal/clients/`, encapsulating external API logic and supporting retries and error handling.
- **Validation Pattern**: Input validation is centralized in `internal/validators/` using go-playground/validator, ensuring consistent validation rules and clear error messages across all endpoints.
- **Soft Delete Pattern**: Most models (see `internal/models/models.go`) include a `deleted_at` timestamp, enabling soft deletion for audit and recovery without physical data removal.

### Key Components
- **Handlers**: Parse HTTP requests, validate input, invoke business logic, and format API responses. (`internal/handlers/`)
- **Services**: Contain business logic, orchestrate workflows, and coordinate between handlers and repositories. (`internal/services/`)
- **Repositories**: Abstract all database access and persistence logic behind interfaces for easy mocking and DB swaps. (`internal/repositories/`)
- **Models**: Define data structures for both API and database, including GORM models and DTOs. (`internal/models/`)
- **Config**: Centralized configuration management, including environment variable parsing and structured config loading. (`internal/config/`)
- **Routes**: API route definitions, grouping, and middleware setup. (`internal/routes/`)

---

## Project Structure

```
enumeration/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── clients/
│   │   ├── mdms.go
│   │   ├── notification.go
│   │   └── workflow.go
│   ├── config/
│   │   └── config.go
│   ├── constants/
│   │   ├── application.go
│   │   ├── errors.go
│   │   └── validation.go
│   ├── database/
│   │   └── database.go
│   ├── dto/
│   │   └── dto.go
│   ├── handlers/
│   │   ├── additional_property_details_handler.go
│   │   ├── amenity_handler.go
│   │   ├── application_handler.go
│   │   ├── application_log_handler.go
│   │   ├── assessment_details_handler.go
│   │   ├── construction_details_handler.go
│   │   ├── coordinates_handler.go
│   │   ├── document_handler.go
│   │   ├── floor_details_handler.go
│   │   ├── gis_handler.go
│   │   ├── igrs_handler.go
│   │   ├── property_address_handler.go
│   │   ├── property_handler.go
│   │   └── property_owner_handler.go
│   ├── middleware/
│   │   ├── auth_middleware.go
│   │   └── cors_logMiddleware.go
│   ├── models/
│   │   └── models.go
│   ├── repositories/
│   │   ├── additional_property_details_repository.go
│   │   ├── amenity_repo.go
│   │   ├── application_log_repository.go
│   │   ├── application_repository.go
│   │   ├── assessment_details_repository.go
│   │   ├── construction_details_repository.go
│   │   ├── coordinates_repo.go
│   │   ├── document_repository.go
│   │   ├── floor_details_repository.go
│   │   ├── gis_repository.go
│   │   ├── igrs_repository.go
│   │   ├── interfaces.go
│   │   ├── property_address_repository.go
│   │   ├── property_owner_repository.go
│   │   └── property_repository.go
│   ├── routes/
│   │   └── routes.go
│   ├── security/
│   │   └── security.go
│   ├── services/
│   │   ├── additional_property_details_service.go
│   │   ├── amenity_service.go
│   │   ├── application_log_service.go
│   │   ├── application_service.go
│   │   ├── assessment_details_service.go
│   │   ├── construction_details_service.go
│   │   ├── coordinates_service.go
│   │   ├── document_service.go
│   │   ├── floor_details_service.go
│   │   ├── gis_service.go
│   │   ├── igrs_service.go
│   │   ├── interfaces.go
│   │   ├── property_address_service.go
│   │   ├── property_owner_service.go
│   │   └── property_service.go
│   └── validators/
│       └── coordinates_validator.go
├── migrations/
│   ├── initial_schema.sql
│   └── insert_table.sql
├── pkg/
│   ├── logger/
│   │   └── logger.go
│   ├── response/
│   │   └── response.go
│   └── utils/
│       └── custom_date.go
├── Dockerfile
├── go.mod
├── go.sum
├── property-tax-apiv2.yml
├── README.md
├── README_UPDATED.md
├── service-documentation 2.md
└── .env
```

---




## API Overview

The Property Tax Enumeration Service exposes a comprehensive set of RESTful APIs for managing property tax applications, properties, owners, documents, GIS data, amenities, and related modules. Below is a categorized list of all major endpoints:

### Additional Property Details
- `GET    /v1/additional-property-details` — List additional property details
- `POST   /v1/additional-property-details` — Create additional property details
- `GET    /v1/additional-property-details/field/{fieldName}` — Get details by field name
- `GET    /v1/additional-property-details/property/{propertyId}` — Get details by property ID
- `GET    /v1/additional-property-details/{id}` — Get details by ID
- `PUT    /v1/additional-property-details/{id}` — Update additional property details
- `DELETE /v1/additional-property-details/{id}` — Delete additional property details

### Amenities
- `GET    /v1/amenities` — List amenities
- `POST   /v1/amenities` — Create amenity
- `GET    /v1/amenities/property/{propertyId}` — Get amenities by property ID
- `GET    /v1/amenities/{id}` — Get amenity by ID
- `PUT    /v1/amenities/{id}` — Update amenity
- `DELETE /v1/amenities/{id}` — Delete amenity

### Application Logs
- `GET    /v1/application-logs` — List application logs
- `POST   /v1/application-logs` — Create application log
- `GET    /v1/application-logs/{id}` — Get application log by ID
- `PUT    /v1/application-logs/{id}` — Update application log
- `DELETE /v1/application-logs/{id}` — Delete application log
- `GET    /v1/applications/{applicationId}/logs` — List logs by application ID

### Applications
- `GET    /v1/applications` — List applications
- `POST   /v1/applications` — Create application
- `GET    /v1/applications/by-number/{applicationNo}` — Get application by application number
- `GET    /v1/applications/search` — Search applications
- `GET    /v1/applications/{id}` — Get application by ID
- `PUT    /v1/applications/{id}` — Update application
- `DELETE /v1/applications/{id}` — Delete application
- `POST   /v1/applications/{id}/action` — Perform action on application
- `PATCH  /v1/applications/{id}/status` — Update application status

### Assessment Details
- `GET    /v1/assessment-details` — List assessment details
- `POST   /v1/assessment-details` — Create assessment details
- `GET    /v1/assessment-details/property/{propertyId}` — Get assessment details by property ID
- `GET    /v1/assessment-details/{id}` — Get assessment details by ID
- `PUT    /v1/assessment-details/{id}` — Update assessment details
- `DELETE /v1/assessment-details/{id}` — Delete assessment details

### Construction Details
- `GET    /v1/construction-details` — List construction details
- `POST   /v1/construction-details` — Create construction details
- `GET    /v1/construction-details/property/{propertyId}` — Get construction details by property ID
- `GET    /v1/construction-details/{id}` — Get construction details by ID
- `PUT    /v1/construction-details/{id}` — Update construction details
- `DELETE /v1/construction-details/{id}` — Delete construction details

### Coordinates
- `GET    /v1/coordinates` — List coordinates
- `POST   /v1/coordinates` — Create coordinates
- `POST   /v1/coordinates/batch` — Create coordinates batch
- `PUT    /v1/coordinates/gis/{gisDataId}` — Replace coordinates by GIS Data ID
- `GET    /v1/coordinates/{id}` — Get coordinates by ID
- `PUT    /v1/coordinates/{id}` — Update coordinates
- `DELETE /v1/coordinates/{id}` — Delete coordinates

### Documents
- `GET    /v1/documents` — List documents
- `POST   /v1/documents` — Create document
- `POST   /v1/documents/batch` — Create multiple documents
- `GET    /v1/documents/property/{propertyId}` — Get documents by property ID
- `GET    /v1/documents/{id}` — Get document by ID
- `PUT    /v1/documents/{id}` — Update document
- `DELETE /v1/documents/{id}` — Delete document

### Floor Details
- `GET    /v1/floor-details` — List floor details
- `POST   /v1/floor-details` — Create floor details
- `GET    /v1/floor-details/{id}` — Get floor details by ID
- `PUT    /v1/floor-details/{id}` — Update floor details
- `DELETE /v1/floor-details/{id}` — Delete floor details

### GIS Data
- `GET    /v1/gis-data` — List GIS data
- `POST   /v1/gis-data` — Create GIS data
- `GET    /v1/gis-data/property/{propertyId}` — Get GIS data by property ID
- `GET    /v1/gis-data/{id}` — Get GIS data by ID
- `PUT    /v1/gis-data/{id}` — Update GIS data
- `DELETE /v1/gis-data/{id}` — Delete GIS data

### IGRS
- `GET    /v1/igrs` — List IGRS
- `POST   /v1/igrs` — Create IGRS
- `GET    /v1/igrs/{id}` — Get IGRS by ID
- `PUT    /v1/igrs/{id}` — Update IGRS
- `DELETE /v1/igrs/{id}` — Delete IGRS

### Properties
- `GET    /v1/properties` — List all properties
- `POST   /v1/properties` — Create a new property
- `GET    /v1/properties/propertyNo/{propertyNo}` — Get property by property number
- `GET    /v1/properties/search` — Search properties
- `PUT    /v1/properties/{id}` — Update property
- `DELETE /v1/properties/{id}` — Delete property

### Property Addresses
- `GET    /v1/property-addresses` — List all property addresses
- `POST   /v1/property-addresses` — Create property address
- `GET    /v1/property-addresses/property/{propertyId}` — Get property address by property ID
- `GET    /v1/property-addresses/search` — Search property addresses
- `GET    /v1/property-addresses/{id}` — Get property address by ID
- `PUT    /v1/property-addresses/{id}` — Update property address
- `DELETE /v1/property-addresses/{id}` — Delete property address

### Property Owners
- `POST   /v1/property-owners` — Create property owner
- `POST   /v1/property-owners/batch` — Create property owners batch
- `GET    /v1/property-owners/property/{propertyId}` — Get property owners by property ID
- `PUT    /v1/property-owners/{id}` — Update property owner
- `DELETE /v1/property-owners/{id}` — Delete property owner

---




## Authentication & Authorization

All endpoints (except health checks) require authentication and authorization to ensure data security and compliance.

### JWT Authentication
- Every API request must include a valid JWT token in the `Authorization` header.
- JWT claims must include user ID, roles, and tenant information.
- Required headers:
  - `Authorization: Bearer <token>`
  - `X-Tenant-ID: <tenant-id>`
  - `X-User-ID: <user-id>`
  - `X-User-Role: <user-role>`
- JWTs are signed using the secret defined in environment variables.
- Expired or invalid tokens return a `401 Unauthorized` response.

#### Example JWT Claims
```json
{
  "sub": "user-uuid",
  "role": "CITIZEN",
  "tenant": "tenant-id",
  "exp": 1700000000,
  "iat": 1699990000
}
```

### Role-Based Access Control (RBAC)
- Endpoint access is governed by user roles, validated on every request:
  - **CITIZEN**: Can create and view own applications and properties.
  - **AGENT**: Can view assigned applications and perform field verification.
  - **SERVICE_MANAGER**: Can manage all applications within their jurisdiction.
  - **COMMISSIONER**: Can approve or reject applications.
- Role checks are enforced in middleware and handlers. Unauthorized access returns `403 Forbidden`.

### MDMS-Driven Dynamic Permissions
- Role permissions and access rules are dynamically loaded from the Master Data Management Service (MDMS).
- Supports endpoint-level and method-level access control, allowing real-time permission updates without redeployment.
- Permission checks are performed on every request to ensure compliance with the latest business rules.

---

---

## Getting Started

Follow these steps to set up and run the Property Tax Enumeration Service for local development or in a containerized environment.

### Prerequisites
- **Go** 1.24 or higher
- **PostgreSQL** 12 or higher
- **Docker** (optional, for containerized setup)

### Local Development Setup

1. **Clone the repository**
  ```bash
  git clone <repository-url>
  cd Property-tax/enumeration
  ```

2. **Install Go dependencies**
  ```bash
  go mod download
  ```

3. **Configure environment variables**
  - Copy the example environment file and update values as needed:
  ```bash
  cp .env.example .env
  # Edit .env with your database, JWT, and service URLs
  ```

4. **Set up the PostgreSQL database**
  - Ensure PostgreSQL is running and accessible as per your `.env` configuration.
  - Create the database and user if not already present.

5. **Run database migrations**
  - Apply schema and seed data:
  ```bash
  go run migrations/*.go
  # Or use a migration tool if provided
  ```

6. **Start the service**
  ```bash
  go run cmd/server/main.go
  ```
  The service will be available at [http://localhost:8080](http://localhost:8080)

### Dockerized Setup

1. **Build the Docker image**
  ```bash
  docker build -t property-enumeration-service .
  ```

2. **Run the service using Docker Compose**
  - Ensure your `docker-compose.yml` is configured (create if needed).
  ```bash
  docker-compose up -d
  ```
  This will start the service and any dependencies (e.g., PostgreSQL) as defined in the compose file.

3. **Access the service**
  - The API will be available at [http://localhost:8080](http://localhost:8080) (or as configured).

---

## Configuration

This section describes how to configure the Property Tax Enumeration Service for different environments.

### Environment Variables
Set the following environment variables in your `.env` file or deployment environment:

#### Database Configuration
```bash
DB_HOST=localhost           # Database host
DB_PORT=5432                # Database port
DB_USER=postgres            # Database user
DB_PASSWORD=your_password   # Database password
DB_NAME=property            # Database name
DB_SCHEMA=DIGIT3            # Database schema (default: DIGIT3)
```

#### Server Configuration
```bash
PORT=8080                   # Port for the HTTP server
```

#### External Service Endpoints
```bash
WORKFLOW_URL=http://localhost:8081           # Workflow service URL
AUDIT_URL=http://localhost:8082/audit/logs   # Audit logging service URL
IDGEN_URL=http://localhost:8083/idgen/v1/generate # ID generation service URL
MDMS_URL=http://localhost:8084               # Master Data Management Service URL
```

#### JWT Configuration
```bash
JWT_SECRET=your_jwt_secret      # Secret for signing JWT tokens
JWT_ISSUER=property-tax-service # JWT issuer
JWT_AUDIENCE=property-tax-users # JWT audience
```

---

### Database Schema
The service uses PostgreSQL with the `DIGIT3` schema. The main tables are:

- **applications**: Stores property tax application data and workflow status
- **properties**: Master table for property details
- **property_owners**: Owner and joint ownership information
- **property_addresses**: Address and location details
- **assessment_details**: Tax assessment records for properties
- **documents**: Uploaded documents and metadata
- **application_logs**: Audit trail and activity logs

Refer to the `migrations/` directory for schema definitions and initial data.

---


## Development Guidelines

- Follow clean, layered architecture: keep handlers, services, repositories, and models separated for maintainability and testability
- Use dependency injection for all services and repositories to enable easy unit testing
- Always use context for database and external service operations
- Implement robust error handling and return standardized API responses
- Use parameterized queries and GORM best practices to prevent SQL injection
- Validate all input data using centralized validators in `internal/validators/`
- Implement proper transaction management for multi-step operations
- Use consistent response formats as defined in [`pkg/response/response.go`](pkg/response/response.go)
- Support pagination, filtering, and searching on all list endpoints
- Provide clear and actionable error messages
- Enforce role-based access control and dynamic permissions via MDMS
- Implement rate limiting and CORS policies as needed
- Use HTTPS in production and follow JWT security best practices

---

## API Documentation

The complete API documentation is available in OpenAPI 3.0 format:
- File: `property-tax-apiv2.yml`
- Import into Swagger UI or Postman for interactive documentation
- Includes request/response schemas, authentication details, and examples

### Live Swagger UI
- [Deployed Swagger UI](http://10.232.161.103:30117/swagger/index.html#) — View and interact with the live API documentation

