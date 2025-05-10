# 🧬 Social Go API

A RESTful API built in Go that allows users to create and interact with posts, follow other users, and manage their own content. Designed with modular middleware, authentication, rate limiting, and modern development practices.

---

## 🚀 Features

- User registration and authentication (JWT)
- Basic Auth for debug endpoints
- CRUD operations on posts
- Commenting and user feeds
- Follow/unfollow functionality
- Swagger docs for API exploration
- CORS, timeout, and structured logging middleware
- Rate limiting middleware
- Docker & CI/CD ready

---

## 🔧 Tech Stack

- **Language**: Go
- **Router**: [Chi](https://github.com/go-chi/chi)
- **Database**: PostgreSQL (`pq` driver)
- **Cache**: Redis
- **Auth**: JWT & Basic Auth
- **Dev Tooling**: Docker Compose, Air (live reload), GitHub Actions

---

## 📡 API Endpoints Overview

- GET /v1/health → Public health check
- GET /v1/swagger/\* → Swagger UI
- POST /v1/authentication/user → Register user
- POST /v1/authentication/token → Generate JWT
- POST /v1/users → Create user
- PUT /v1/users/activate/{token} → Activate user
- GET /v1/users/{id} → Get user by ID (auth required)
- PUT /v1/users/{id}/follow → Follow user (auth required)
- PUT /v1/users/{id}/unfollow → Unfollow user (auth required)
- GET /v1/users/{id}/feeds → User feed (auth required)
- POST /v1/posts → Create post (auth required)
- GET /v1/posts → Get all posts (auth required)
- GET /v1/posts/{id} → Get post by ID (auth required)
- PATCH /v1/posts/{id} → Update post (moderator access)
- DELETE /v1/posts/{id} → Delete post (admin access)
- GET /v1/debug/vars → Go expvar (Basic Auth required)

---

## 🛠️ Local Development

### Prerequisites

- Go 1.23+
- PostgreSQL
- Redis
- Docker (optional)

### Using [Air](https://github.com/cosmtrek/air)

`air`

### Without Air

`go build -v ./...`

### Running Test

`go test -race ./...`

### Authentication

- **JWT**: Used for protected user and post endpoints
- **Basic**: Used for internal debug route v1/debug/vars

### Authentication

- **PORT**: port to expose api --> 5100
- **FRONTEND_URL**: to connect with frontend
- **SWAGGER_API_URL**: Swagger URL --> localhost:5100
- **SENDGRID_API_KEY**: Sengrid Key for Email
- **FROM_EMAIL**: required for sendgrid
- **DB_URI**: database connection --> postgres
- **AUTH_HEADER_USERNAME**:
- **AUTH_HEADER_PASSWORD**:
- **AUTH_TOKEN_SECRET**:
- **REDIS_ADDR**: localhost:6379
- **REDIS_PW**: e.g 0 / 1
- **CORS_ALLOWED_ORIGIN**:

### CORS Configuration

```
r.Use(cors.Handler(cors.Options{
  AllowedOrigins:   []string{envData.CORS_ALLOWED_ORIGIN},
  AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
  AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
  ExposedHeaders:   []string{"Link"},
  AllowCredentials: true,
  MaxAge:           300,
}))

```

### Rate Limiting

Rate limiting is configurable via the application config (requests per timeframe, timeframe in seconds, enabled/disabled).

## Docker & CI/CD

A GitHub Actions workflow handles build/test/lint/audit:

```
name: Audit
on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  audit:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Code
        uses: actions/checkout@v2

      - name: Set up Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.23.1

      - name: Verify Depedencies
        run: go mod verify

      - name: Build
        run: go build -v ./...

      - name: Run go vet
        run: go vet ./...

      - name: Install staticcheck
        run: go install honnef.co/go/tools/cmd/staticcheck@latest

      - name: Run staticcheck
        run: staticcheck ./...

      - name: Run Tests
        run: go test -race ./...

```

## License

MIT

## Author

Built by Michael. Feel free to update or fork this project.
