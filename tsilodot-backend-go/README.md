# Tsilodot Backend (Go)

A high-performance task management (todo list) backend API built with Go, featuring Redis caching, JWT authentication, and comprehensive API documentation.

## 🚀 Features

- **RESTful API:** Clean and predictable API for task management.
- **Layered Architecture:** Follows separation of concerns (Routes, Controllers, Services, Repositories).
- **Authentication:** Secure user authentication using JWT (JSON Web Tokens).
- **Caching:** Redis integration for fast task retrieval with a 10-minute TTL.
- **Database:** PostgreSQL with GORM ORM.
- **Migrations:** Robust database schema management with `golang-migrate`.
- **Validation:** Strong request validation using `go-playground/validator`.
- **Documentation:** Interactive API docs with Swagger UI and Scalar.
- **Logging:** Structured logging with `zerolog`.

## 🛠️ Tech Stack

- **Framework:** [Fiber v3](https://docs.gofiber.io/)
- **ORM:** [GORM](https://gorm.io/) with PostgreSQL
- **Cache:** [Redis](https://redis.io/)
- **Documentation:** Swagger UI (`/swagger`) & Scalar UI (`/scalar`)
- **Auth:** JWT (JSON Web Tokens)
- **Logging:** [Zerolog](https://github.com/rs/zerolog)
- **Validation:** [go-playground/validator](https://github.com/go-playground/validator)
- **Migrations:** [golang-migrate](https://github.com/golang-migrate/migrate)

## 🏗️ Architecture

The project follows a **Layered Architecture** to ensure maintainability and testability:

1.  **Routes (`/routes`):** Defines API endpoints and applies middlewares.
2.  **Controllers (`/controller`):** Handles HTTP request parsing, validation, and response formatting.
3.  **Services (`/service`):** Contains business logic. Interacts with repositories and cache.
4.  **Repositories (`/repository`):** Handles data persistence logic using GORM.
5.  **Models (`/model`):** Defines database schemas and entity structures.
6.  **DTOs (`/dto`):** Data Transfer Objects for structured API requests and responses.

## 🏁 Getting Started

### Prerequisites

- Go 1.25+
- PostgreSQL
- Redis
- `golang-migrate` CLI (optional, but recommended for manual migrations)

### ⚙️ Environment Setup

1. Clone the repository:
   ```bash
   git clone <repository-url>
   cd tsilodot-backend-go
   ```

2. Copy `.env.example` to `.env` and fill in your local configurations:
   ```bash
   cp .env.example .env
   ```

3. Update the `.env` file with your database and redis credentials:
   ```env
   APP_PORT=8080
   DB_HOST=localhost
   DB_USER=postgres
   DB_PASSWORD=your_password
   DB_NAME=tsilodot_db
   DB_PORT=5432
   DB_SSL_MODE=disable
   JWT_SECRET_KEY=your_secret_key
   REDIS_HOST=localhost
   REDIS_PORT=6379
   ```

### 🏃 Running the Application

**Run with Docker (Recommended):**
This will start the Go app, PostgreSQL, Redis, and automatically run database migrations.

```bash
docker compose up -d
```

**Run locally (Manual):**

1. Install dependencies:
   ```bash
   go mod tidy
   ```

2. Run the application:
   ```bash
   go run main.go
   ```

The server will start on the port specified in your `.env` file (default: `8080`).

> You might need to run the seeder after run the app. [See at the `🌱 Database Seeder` section](###-🌱-database-seeder)

### 🗄️ Database Migrations

Migrations are stored in `db/migrations`. The application can be updated using the `migrate` CLI:

```bash
# Run up migrations
migrate -path="./db/migrations" -database="postgres://user:pass@localhost:5432/dbname?sslmode=disable" up

# Rollback migrations
migrate -path="./db/migrations" -database="postgres://user:pass@localhost:5432/dbname?sslmode=disable" down
```

### 🌱 Database Seeder

To populate the database with initial test data (5 users, each with 4-15 tasks):

```bash
go run db/seeds/seeder.go
```

**Note:** This command will **truncate** the `tasks` and `users` tables before seeding. All existing user and task data will be lost.

### 🧪 Running Tests

```bash
go test ./...
```

The test suite includes unit tests for controllers, services, repositories, and middlewares.

## 📄 API Documentation

Once the server is running, you can access the interactive documentation at:

- **Swagger UI:** `http://localhost:8080/swagger`
- **Scalar UI:** `http://localhost:8080/scalar`

The OpenAPI specification is located at `docs/openapi.yaml`.

## 📁 Key Directories

- `controller/`: Request handlers and input validation.
- `service/`: Business logic and cache management.
- `repository/`: Data access layer (GORM).
- `model/`: Database entities and schemas.
- `dto/`: Data Transfer Objects for API contracts.
- `db/`: Database connections, Redis client, and migrations.
- `helpers/`: Utility functions (JWT, Bcrypt, Logger, Validators).
- `middlewares/`: Custom Fiber middlewares (e.g., Auth).
- `routes/`: API endpoint definitions.
- `docs/`: API documentation (OpenAPI spec).

## 📝 Coding Conventions

- **Dependency Injection:** Interfaces are used for Services and Repositories to facilitate testing and mocking.
- **Concurrency:** Uses `sync.WaitGroup` for parallel database operations where appropriate.
- **Cache-Aside Pattern:** Redis is used to cache task lookups, with invalidation on updates/deletes.
- **Clean Responses:** Standardized error and success responses via DTOs and helper functions.

## 📡 API Endpoints

All API endpoints are prefixed with `/api`. Authentication is handled via a Bearer Token in the `Authorization` header for protected routes.

### 🔐 Authentication

#### `POST /auth/register`
- **Description**: Register a new user.
- **Auth Required**: No
- **Request Body**:
  ```json
  {
    "name": "John Doe",
    "email": "john@example.com",
    "password": "securepassword",
    "confirm_password": "securepassword"
  }
  ```
- **Response (200 OK)**:
  ```json
  {
    "message": "Registration successful",
    "data": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "balance": 0,
      "created_at": "2024-03-20T10:00:00Z",
      "updated_at": "2024-03-20T10:00:00Z",
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
  }
  ```

#### `POST /auth/login`
- **Description**: Login user and get JWT token.
- **Auth Required**: No
- **Request Body**:
  ```json
  {
    "email": "john@example.com",
    "password": "securepassword"
  }
  ```
- **Response (200 OK)**:
  ```json
  {
    "message": "Login successful",
    "data": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com",
      "balance": 0,
      "created_at": "2024-03-20T10:00:00Z",
      "updated_at": "2024-03-20T10:00:00Z",
      "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
  }
  ```

### 📋 Task Management

#### `GET /tasks`
- **Description**: List all tasks for the authenticated user with pagination.
- **Auth Required**: Yes
- **Query Parameters**:
  - `page` (integer, default: `1`)
  - `limit` (integer, default: `5`)
- **Response (200 OK)**:
  ```json
  {
    "message": "Tasks fetched successfully",
    "pagination": {
      "current_page": 1,
      "total_pages": 1,
      "total_items": 1
    },
    "data": [
      {
        "id": 1,
        "user_id": 1,
        "title": "Complete project documentation",
        "description": "Finish the README file for the backend API",
        "status": "pending",
        "due_date": "2024-03-25T00:00:00Z",
        "created_at": "2024-03-20T10:00:00Z",
        "updated_at": "2024-03-20T10:00:00Z"
      }
    ]
  }
  ```

#### `POST /tasks`
- **Description**: Create a new task.
- **Auth Required**: Yes
- **Request Body**:
  ```json
  {
    "title": "Complete project documentation",
    "description": "Finish the README file for the backend API",
    "status": "pending",
    "due_date": "2024-03-25"
  }
  ```
- **Response (200 OK)**:
  ```json
  {
    "message": "Task created successfully",
    "data": {
      "id": 1,
      "user_id": 1,
      "title": "Complete project documentation",
      "description": "Finish the README file for the backend API",
      "status": "pending",
      "due_date": "2024-03-25T00:00:00Z",
      "created_at": "2024-03-20T10:00:00Z",
      "updated_at": "2024-03-20T10:00:00Z"
    }
  }
  ```

#### `GET /tasks/:id`
- **Description**: Get details of a specific task.
- **Auth Required**: Yes
- **Response (200 OK)**:
  ```json
  {
    "message": "Task fetched successfully",
    "data": {
      "id": 1,
      "user_id": 1,
      "title": "Complete project documentation",
      "description": "Finish the README file for the backend API",
      "status": "pending",
      "due_date": "2024-03-25T00:00:00Z",
      "created_at": "2024-03-20T10:00:00Z",
      "updated_at": "2024-03-20T10:00:00Z"
    }
  }
  ```

#### `PUT /tasks/:id`
- **Description**: Update an existing task.
- **Auth Required**: Yes
- **Request Body**:
  ```json
  {
    "title": "Complete project documentation",
    "description": "Finish the README and add examples",
    "status": "completed",
    "due_date": "2024-03-25"
  }
  ```
- **Response (200 OK)**:
  ```json
  {
    "message": "Task updated successfully",
    "data": {
      "id": 1,
      "user_id": 1,
      "title": "Complete project documentation",
      "description": "Finish the README and add examples",
      "status": "completed",
      "due_date": "2024-03-25T00:00:00Z",
      "created_at": "2024-03-20T10:00:00Z",
      "updated_at": "2024-03-20T10:15:00Z"
    }
  }
  ```

#### `DELETE /tasks/:id`
- **Description**: Delete a task.
- **Auth Required**: Yes
- **Response (200 OK)**:
  ```json
  {
    "message": "Task deleted successfully"
  }
  ```
