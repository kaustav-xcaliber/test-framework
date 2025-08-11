# Development Guide

This document provides detailed information for developers working on the API Test Framework.

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- Docker and Docker Compose
- Git
- Your preferred IDE (VS Code, GoLand, Vim, etc.)

### Initial Setup

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd test-framework-prototype
   ```

2. **Install dependencies**

   ```bash
   make deps
   # or manually:
   go mod tidy
   go mod download
   ```

3. **Set up environment**

   ```bash
   cp env.example .env
   # Edit .env with your configuration
   ```

4. **Start development services**

   ```bash
   make docker-up
   # or manually:
   docker-compose -f docker-compose.dev.yml up -d
   ```

5. **Run the application**
   ```bash
   make run
   # or manually:
   go run cmd/api-server/main.go
   ```

## 🏗️ Project Structure

```
test-framework-prototype/
├── cmd/                    # Application entry points
│   └── api-server/        # Main API server
├── internal/               # Private application code
│   ├── config/            # Configuration management
│   ├── database/          # Database connections and migrations
│   ├── handlers/          # HTTP request handlers
│   ├── models/            # Data models and database schema
│   ├── services/          # Business logic layer
│   ├── testrunner/        # Test execution engine
│   └── utils/             # Utility functions
├── scripts/                # Utility scripts and tools
├── docker-compose.dev.yml  # Development environment
├── Dockerfile             # Application container
├── Makefile               # Build and development commands
├── go.mod                 # Go module definition
└── README.md              # Project documentation
```

## 🔧 Development Workflow

### 1. Code Changes

1. **Create a feature branch**

   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes**

   - Follow Go coding standards
   - Add tests for new functionality
   - Update documentation as needed

3. **Test your changes**

   ```bash
   make test          # Run all tests
   make test-cover    # Run tests with coverage
   make build         # Ensure code compiles
   ```

4. **Format and lint**

   ```bash
   make format        # Format code
   make lint          # Run linter
   ```

5. **Commit your changes**

   ```bash
   git add .
   git commit -m "feat: add your feature description"
   ```

6. **Push and create PR**
   ```bash
   git push origin feature/your-feature-name
   # Create Pull Request on GitHub
   ```

### 2. Testing Strategy

#### Unit Tests

- Test individual functions and methods
- Mock external dependencies
- Aim for >80% code coverage

#### Integration Tests

- Test service interactions
- Use test database
- Test API endpoints

#### Example Test Structure

```go
func TestService_CreateService(t *testing.T) {
    // Arrange
    db := setupTestDB(t)
    service := NewServiceService(db)

    // Act
    result, err := service.CreateService(&models.Service{
        Name: "test-service",
        BaseURL: "http://localhost:8080",
    })

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "test-service", result.Name)
}
```

### 3. Database Changes

#### Adding New Models

1. **Update models.go**

   ```go
   type NewModel struct {
       ID        string    `json:"id" gorm:"primarykey;type:uuid;default:gen_random_uuid()"`
       Name      string    `json:"name" gorm:"not null"`
       CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
       UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
   }
   ```

2. **Add to main.go migration**

   ```go
   if err := db.AutoMigrate(
       &models.Service{},
       &models.TestCase{},
       &models.TestRun{},
       &models.TestResult{},
       &models.NewModel{}, // Add new model here
   ); err != nil {
       log.Fatalf("Failed to migrate database: %v", err)
   }
   ```

3. **Create migration script** (for production)
   ```sql
   -- migrations/001_add_new_model.sql
   CREATE TABLE new_models (
       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
       name VARCHAR(255) NOT NULL,
       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
   );
   ```

### 4. API Changes

#### Adding New Endpoints

1. **Create handler method**

   ```go
   func (h *Handler) NewEndpoint(c *gin.Context) {
       // Implementation
   }
   ```

2. **Add route in main.go**

   ```go
   api := router.Group("/api/v1")
   {
       // Existing routes...

       // New routes
       api.GET("/new-endpoint", handler.NewEndpoint)
   }
   ```

3. **Add tests**
   ```go
   func TestHandler_NewEndpoint(t *testing.T) {
       // Test implementation
   }
   ```

## 🐳 Docker Development

### Development Environment

```bash
# Start all services
make docker-up

# View logs
docker logs -f api-test-framework-server-dev

# Access database
docker exec -it api-test-framework-postgres-dev psql -U user -d api_test_framework

# Access Redis
docker exec -it api-test-framework-redis-dev redis-cli

# Stop services
make docker-down
```

### Container Management

```bash
# Rebuild and restart specific service
docker-compose -f docker-compose.dev.yml up -d --build api-server

# View service status
docker-compose -f docker-compose.dev.yml ps

# View service logs
docker-compose -f docker-compose.dev.yml logs api-server
```

## 📊 Monitoring and Debugging

### Logging

The application uses structured logging with different levels:

```go
// Debug level (development)
log.Printf("Debug: %+v", data)

// Info level (general information)
log.Printf("Info: Service started on port %s", port)

// Error level (errors)
log.Printf("Error: Failed to connect to database: %v", err)
```

### Health Checks

```bash
# Application health
curl http://localhost:8080/health

# Database health
curl http://localhost:8080/health/db

# Redis health
curl http://localhost:8080/health/redis
```

### Debugging

1. **Set log level to debug**

   ```env
   LOG_LEVEL=debug
   ```

2. **Use Go debugger (delve)**

   ```bash
   dlv debug cmd/api-server/main.go
   ```

3. **Profile the application**
   ```bash
   go run -cpuprofile=cpu.prof cmd/api-server/main.go
   go tool pprof cpu.prof
   ```

## 🔒 Security Considerations

### Authentication

- Store sensitive data in environment variables
- Use secure authentication methods
- Implement proper input validation
- Sanitize user inputs

### Data Protection

- Encrypt sensitive data at rest
- Use HTTPS in production
- Implement proper access controls
- Log security events

### Code Security

- Keep dependencies updated
- Use security linters (gosec)
- Follow OWASP guidelines
- Regular security audits

## 📈 Performance Optimization

### Database Optimization

1. **Connection Pooling**

   ```go
   sqlDB, err := db.DB()
   if err != nil {
       return err
   }

   sqlDB.SetMaxOpenConns(25)
   sqlDB.SetMaxIdleConns(5)
   sqlDB.SetConnMaxLifetime(5 * time.Minute)
   ```

2. **Query Optimization**
   - Use database indexes
   - Avoid N+1 queries
   - Use pagination for large datasets
   - Implement query caching

### Caching Strategy

1. **Redis Caching**

   ```go
   // Cache frequently accessed data
   cacheKey := fmt.Sprintf("service:%s", serviceID)
   err := redisClient.Set(ctx, cacheKey, data, time.Hour).Err()
   ```

2. **In-Memory Caching**
   - Cache configuration data
   - Cache test results
   - Implement TTL for cache entries

## 🧪 Testing Best Practices

### Test Organization

```
internal/
├── services/
│   ├── service_service.go
│   └── service_service_test.go
├── handlers/
│   ├── service_handler.go
│   └── service_handler_test.go
```

### Test Utilities

```go
// testutils/db.go
func SetupTestDB(t *testing.T) *gorm.DB {
    // Setup test database
}

func CleanupTestDB(t *testing.T, db *gorm.DB) {
    // Cleanup test data
}

// testutils/redis.go
func SetupTestRedis(t *testing.T) *redis.Client {
    // Setup test Redis
}
```

### Mocking

```go
// Use interfaces for testability
type ServiceRepository interface {
    Create(service *models.Service) error
    GetByID(id string) (*models.Service, error)
}

// Mock implementation
type MockServiceRepository struct {
    services map[string]*models.Service
}

func (m *MockServiceRepository) Create(service *models.Service) error {
    // Mock implementation
}
```

## 📚 Documentation

### Code Documentation

- Use Go doc comments for exported functions
- Document complex business logic
- Include examples in documentation
- Keep README updated

### API Documentation

- Document all endpoints
- Include request/response examples
- Document error codes
- Use OpenAPI/Swagger if possible

### Architecture Documentation

- Document design decisions
- Include sequence diagrams
- Document data flow
- Keep architecture decisions record (ADR)

## 🚀 Deployment

### Development Deployment

```bash
# Build and run locally
make build
./bin/api-test-framework

# Run with Docker
make docker-build
make docker-up
```

### Production Deployment

```bash
# Build production binary
make prod-build

# Build production Docker image
docker build -t api-test-framework:latest .

# Deploy with docker-compose
docker-compose -f docker-compose.prod.yml up -d
```

### Environment Configuration

```env
# Production environment
ENVIRONMENT=production
LOG_LEVEL=info
DB_SSL_MODE=require
REDIS_PASSWORD=secure_password
```

## 🤝 Contributing Guidelines

### Code Review Process

1. **Self-review before PR**

   - Run all tests
   - Check code formatting
   - Verify no sensitive data
   - Update documentation

2. **PR Requirements**

   - Clear description of changes
   - Link to related issues
   - Include tests
   - Update documentation

3. **Review Checklist**
   - Code follows Go standards
   - Tests are comprehensive
   - Documentation is updated
   - No security issues
   - Performance considerations

### Communication

- Use GitHub Issues for bugs
- Use GitHub Discussions for questions
- Use GitHub Projects for planning
- Regular team sync meetings

## 📞 Getting Help

### Internal Resources

- **Code Comments**: Inline documentation
- **README.md**: Project overview
- **DEVELOPMENT.md**: This guide
- **Team Members**: Direct communication

### External Resources

- **Go Documentation**: https://golang.org/doc/
- **Gin Framework**: https://gin-gonic.com/docs/
- **GORM**: https://gorm.io/docs/
- **Go Best Practices**: https://github.com/golang/go/wiki/CodeReviewComments

### Troubleshooting

1. **Check logs first**

   ```bash
   docker logs api-test-framework-server-dev
   ```

2. **Verify environment**

   ```bash
   cat .env
   ```

3. **Check service status**

   ```bash
   docker-compose -f docker-compose.dev.yml ps
   ```

4. **Ask for help**
   - Create GitHub issue
   - Ask in team chat
   - Schedule pair programming session

---

**Happy Coding! 🚀**
