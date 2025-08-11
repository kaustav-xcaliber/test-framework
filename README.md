# API Test Framework

A comprehensive HTTP API testing framework built with Go, featuring test management, execution, and detailed reporting capabilities.

## 🚀 Features

- **Service Management**: Add, edit, and manage microservices for testing
- **Test Case Management**: Create and manage HTTP test cases with JSON specifications
- **Test Execution**: Execute tests using httpexpect with comprehensive assertions
- **Result Tracking**: Track test runs and individual test results with detailed metrics
- **REST API**: Complete REST API for frontend integration
- **Redis Integration**: Job queue and caching support
- **PostgreSQL**: Persistent storage for all test data
- **Detailed Reporting**: Comprehensive test execution reports and analytics
- **Curl Import**: Create tests directly from curl commands
- **Authentication Support**: Multiple authentication methods (Bearer, API Key, Basic, OAuth2)

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    API Test Framework                      │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐        │
│  │   Service   │  │    Test     │  │   Test      │        │
│  │ Management  │  │ Management  │  │ Execution   │        │
│  │  Service    │  │  Service    │  │  Service    │        │
│  └─────────────┘  └─────────────┘  └─────────────┘        │
├─────────────────────────────────────────────────────────────┤
│                    Data Layer                              │
│  ┌─────────────┐  ┌─────────────┐                         │
│  │ PostgreSQL  │  │    Redis    │                         │
│  │  Database   │  │   Cache     │                         │
│  └─────────────┘  └─────────────┘                         │
└─────────────────────────────────────────────────────────────┘
```

## 📊 Reporting & Analytics

### Test Execution Reports

The framework provides comprehensive reporting capabilities for test execution:

#### 1. Test Run Summary

- **Execution Status**: Running, Completed, Failed
- **Test Counts**: Total, Passed, Failed, Skipped
- **Performance Metrics**: Execution time, Response times
- **Timestamps**: Start time, completion time, duration

#### 2. Individual Test Results

- **Test Case Details**: Name, description, service
- **Execution Status**: Passed, Failed, Skipped
- **Performance Data**: Response time, request/response size
- **Error Details**: Error messages, stack traces
- **Response Data**: Full response payload for analysis

#### 3. Historical Analytics

- **Trend Analysis**: Test success rates over time
- **Performance Trends**: Response time patterns
- **Failure Analysis**: Common failure patterns
- **Service Health**: Service reliability metrics

#### 4. Export Capabilities

- **JSON Reports**: Machine-readable test results
- **CSV Export**: Spreadsheet-compatible data
- **HTML Reports**: Human-readable formatted reports
- **Integration APIs**: Webhook notifications

### Metrics & KPIs

- **Test Success Rate**: Percentage of passed tests
- **Average Response Time**: Mean execution time
- **Failure Rate**: Percentage of failed tests
- **Test Coverage**: Number of endpoints tested
- **Service Uptime**: Service availability metrics

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Docker and Docker Compose
- PostgreSQL (via Docker)
- Redis (via Docker)

### 1. Clone and Setup

```bash
git clone <repository-url>
cd test-framework-prototype
```

### 2. Environment Configuration

Copy the example environment file:

```bash
cp env.example .env
```

Edit `.env` file with your configuration:

```env
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=0.0.0.0

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=user
DB_PASSWORD=password
DB_NAME=api_test_framework
DB_SSL_MODE=disable

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# Application Configuration
LOG_LEVEL=debug
ENVIRONMENT=development
```

### 3. Start Services

```bash
# Start all services (PostgreSQL, Redis, pgAdmin, API Server)
docker-compose -f docker-compose.dev.yml up -d

# Or start only the database services
docker-compose -f docker-compose.dev.yml up -d postgres redis pgadmin
```

### 4. Run the Application

```bash
# Install dependencies
go mod tidy

# Run the application
go run cmd/api-server/main.go
```

The API will be available at `http://localhost:8080`

### 5. Test the API

Import the provided Postman collection: `API_Test_Framework.postman_collection.json`

## 📡 API Endpoints

### Health Check

- `GET /health` - Health check endpoint

### Service Management

- `GET /api/v1/services` - List all services
- `POST /api/v1/services` - Create a new service
- `GET /api/v1/services/{id}` - Get service by ID
- `PUT /api/v1/services/{id}` - Update service
- `DELETE /api/v1/services/{id}` - Delete service

### Test Management

- `GET /api/v1/tests` - List all tests
- `POST /api/v1/tests` - Create a new test
- `POST /api/v1/tests/from-curl` - Create test from curl command
- `GET /api/v1/tests/{id}` - Get test by ID
- `PUT /api/v1/tests/{id}` - Update test
- `DELETE /api/v1/tests/{id}` - Delete test

### Test Execution & Reporting

- `POST /api/v1/test-runs` - Start a test run
- `GET /api/v1/test-runs/{id}` - Get test run status and summary
- `GET /api/v1/test-runs/{id}/results` - Get detailed test results
- `GET /api/v1/test-runs` - List all test runs with pagination

## 📋 Test Specification Format

Tests are defined using JSON specifications:

```json
{
  "name": "Get User Test",
  "description": "Test to get user by ID",
  "service_name": "user-service",
  "request": {
    "method": "GET",
    "url": "/users/1",
    "headers": {
      "Content-Type": "application/json"
    }
  },
  "assertions": [
    {
      "type": "status_code",
      "expected": 200
    },
    {
      "type": "json_path",
      "path": "$.id",
      "matcher": "exists"
    },
    {
      "type": "json_path",
      "path": "$.name",
      "matcher": "equals",
      "expected": "Leanne Graham"
    }
  ]
}
```

### Supported Assertions

1. **Status Code**: Verify HTTP status code

   ```json
   {
     "type": "status_code",
     "expected": 200
   }
   ```

2. **JSON Path**: Verify JSON response values

   ```json
   {
     "type": "json_path",
     "path": "$.user.id",
     "matcher": "exists"
   }
   ```

3. **JSON Path with Value**: Compare JSON values

   ```json
   {
     "type": "json_path",
     "path": "$.user.name",
     "matcher": "equals",
     "expected": "John Doe"
   }
   ```

4. **Response Time**: Verify response time constraints
   ```json
   {
     "type": "response_time",
     "matcher": "less_than",
     "expected": 2000
   }
   ```

## 🗄️ Database Schema

### Services Table

```sql
CREATE TABLE services (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    base_url VARCHAR(500) NOT NULL,
    auth_config JSONB DEFAULT '{}',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT true
);
```

### Test Cases Table

```sql
CREATE TABLE test_cases (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    service_id UUID NOT NULL REFERENCES services(id) ON DELETE CASCADE,
    name VARCHAR(200) NOT NULL,
    description TEXT,
    test_spec JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT true
);
```

### Test Runs Table

```sql
CREATE TABLE test_runs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(200),
    status VARCHAR(20) CHECK (status IN ('running', 'completed', 'failed')),
    total_tests INTEGER DEFAULT 0,
    passed_tests INTEGER DEFAULT 0,
    failed_tests INTEGER DEFAULT 0,
    execution_time_ms BIGINT,
    started_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    completed_at TIMESTAMP
);
```

### Test Results Table

```sql
CREATE TABLE test_results (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    test_run_id UUID NOT NULL REFERENCES test_runs(id) ON DELETE CASCADE,
    test_case_id UUID NOT NULL REFERENCES test_cases(id),
    status VARCHAR(20) CHECK (status IN ('passed', 'failed', 'skipped')),
    execution_time_ms INTEGER,
    error_message TEXT,
    response_data JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 🔐 Authentication Configuration

The framework supports different authentication methods for testing APIs. Each service can have its own authentication configuration stored in the `auth_config` field.

### AuthConfig Structure

```json
{
  "type": "bearer",
  "token": "your-jwt-token-here"
}
```

### Supported Authentication Types

#### 1. Bearer Token

```json
{
  "type": "bearer",
  "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### 2. API Key

```json
{
  "type": "api_key",
  "key_name": "X-API-Key",
  "key_value": "your-api-key-here"
}
```

#### 3. Basic Authentication

```json
{
  "type": "basic",
  "username": "user",
  "password": "password"
}
```

#### 4. OAuth2

```json
{
  "type": "oauth2",
  "client_id": "your-client-id",
  "client_secret": "your-client-secret",
  "token_url": "https://auth.example.com/token"
}
```

### Usage in Tests

When creating test cases, you can reference the service's authentication configuration. The test runner will automatically apply the appropriate authentication headers based on the service's `auth_config`.

## 🌀 Creating Tests from Curl Commands

The framework supports creating tests directly from curl commands, making it easy to convert existing API calls into automated tests.

### Endpoint

```
POST /api/v1/tests/from-curl
```

### Request Format

```json
{
  "service_id": "service-uuid",
  "name": "Test Name",
  "description": "Test Description",
  "curl_command": "curl -X GET 'https://api.example.com/users/1' -H 'Authorization: Bearer token'",
  "assertions": [
    {
      "type": "status_code",
      "expected": 200
    },
    {
      "type": "json_path",
      "path": "$.id",
      "matcher": "exists"
    }
  ]
}
```

### Supported Curl Features

The curl parser supports the following curl options:

- **HTTP Methods**: `-X`, `--request` (GET, POST, PUT, DELETE, PATCH)
- **Headers**: `-H`, `--header`
- **Data**: `-d`, `--data`, `--data-raw`, `--data-binary`
- **Form Data**: `-F`, `--form`
- **Authentication**: `-u`, `--user` (Basic Auth)
- **Cookies**: `-b`, `--cookie`
- **Query Parameters**: Automatically extracted from URL
- **Path Variables**: Automatically detected (e.g., `{id}`)

### Example Curl Commands

#### 1. Simple GET Request

```bash
curl -X GET 'https://jsonplaceholder.typicode.com/users/1'
```

#### 2. POST with JSON Body

```bash
curl -X POST 'https://jsonplaceholder.typicode.com/users' \
  -H 'Content-Type: application/json' \
  -d '{"name": "John Doe", "email": "john@example.com"}'
```

#### 3. FHIR Export Request

```bash
curl --location 'https://staging-fhir.ecwcloud.com/fhir/r4/FFBJCD/Group/123/$export?since=2024-01-01' \
  --header 'Accept: application/fhir+json' \
  --header 'Prefer: respond-async' \
  --header 'Authorization: Bearer <bearer-token>'
```

#### 4. Form Data

```bash
curl -X POST 'https://api.example.com/upload' \
  -F 'file=@document.pdf' \
  -F 'description=Test document'
```

#### 5. Healthcare API with Bearer Token

```bash
curl --location 'localhost:3001/get_encounter_state_by_alternate_visit_number/execute' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9...' \
--data '{
    "alternate_visit_number":"096076481"
}'
```

### Response Format

```json
{
  "data": {
    "id": "test-uuid",
    "service_id": "service-uuid",
    "name": "Test Name",
    "description": "Test Description",
    "test_spec": "{\"name\":\"Test Name\",\"request\":{...},\"assertions\":[...]}",
    "created_at": "2025-08-08T11:42:31.853Z",
    "updated_at": "2025-08-08T11:42:31.853Z",
    "is_active": true
  },
  "parsed_curl": {
    "method": "GET",
    "url": "https://api.example.com/users/1",
    "headers": {
      "Authorization": "Bearer token"
    },
    "body": "",
    "queryParams": {},
    "pathVariables": {},
    "requestType": "Fetch",
    "rawCommand": "curl -X GET 'https://api.example.com/users/1' -H 'Authorization: Bearer token'"
  }
}
```

### Default Assertions

If no custom assertions are provided, the framework automatically adds:

```json
{
  "type": "status_code",
  "expected": 200
}
```

### Custom Assertions

You can provide custom assertions to validate the response:

```json
{
  "assertions": [
    {
      "type": "status_code",
      "expected": 201
    },
    {
      "type": "json_path",
      "path": "$.id",
      "matcher": "exists"
    },
    {
      "type": "json_path",
      "path": "$.name",
      "matcher": "equals",
      "expected": "John Doe"
    },
    {
      "type": "response_time",
      "matcher": "less_than",
      "expected": 2000
    }
  ]
}
```

### Healthcare API Testing Examples

For healthcare APIs, you can create comprehensive tests with specific assertions:

#### Encounter API Test

```json
{
  "service_id": "healthcare-service-uuid",
  "name": "Get Encounter by Alternate Visit ID",
  "description": "Test to get encounter details by alternate visit number",
  "curl_command": "curl --location 'localhost:3001/get_encounter_state_by_alternate_visit_number/execute' --header 'Content-Type: application/json' --header 'Authorization: Bearer <token>' --data '{\"alternate_visit_number\":\"096076481\"}'",
  "assertions": [
    {
      "type": "status_code",
      "expected": 200
    },
    {
      "type": "json_path",
      "path": "$[0].encounter_id",
      "matcher": "exists"
    },
    {
      "type": "json_path",
      "path": "$[0].alternate_visit_id",
      "matcher": "equals",
      "expected": "096076481"
    },
    {
      "type": "json_path",
      "path": "$[0].patient_class",
      "matcher": "equals",
      "expected": "OUTPATIENT"
    },
    {
      "type": "json_path",
      "path": "$[0].type",
      "matcher": "equals",
      "expected": "A08"
    },
    {
      "type": "json_path",
      "path": "$[0].lineage",
      "matcher": "equals",
      "expected": "HL7"
    },
    {
      "type": "json_path",
      "path": "$[0].patient_id",
      "matcher": "exists"
    },
    {
      "type": "json_path",
      "path": "$[0].insurance",
      "matcher": "exists"
    },
    {
      "type": "json_path",
      "path": "$[0].providers",
      "matcher": "exists"
    },
    {
      "type": "json_path",
      "path": "$[0].facilities",
      "matcher": "exists"
    },
    {
      "type": "response_time",
      "matcher": "less_than",
      "expected": 5000
    }
  ]
}
```

#### FHIR Resource Test

```json
{
  "service_id": "fhir-service-uuid",
  "name": "Get Patient by ID",
  "description": "Test FHIR patient resource retrieval",
  "curl_command": "curl -X GET 'https://fhir.example.com/Patient/123' -H 'Accept: application/fhir+json' -H 'Authorization: Bearer <token>'",
  "assertions": [
    {
      "type": "status_code",
      "expected": 200
    },
    {
      "type": "json_path",
      "path": "$.resourceType",
      "matcher": "equals",
      "expected": "Patient"
    },
    {
      "type": "json_path",
      "path": "$.id",
      "matcher": "exists"
    },
    {
      "type": "json_path",
      "path": "$.identifier",
      "matcher": "exists"
    }
  ]
}
```

## 📈 Advanced Reporting Features

### 1. Real-time Test Execution Monitoring

- **Live Status Updates**: Real-time test execution progress
- **WebSocket Notifications**: Instant updates on test completion
- **Progress Indicators**: Visual progress bars and status indicators
- **Execution Logs**: Detailed logs during test execution

### 2. Performance Analytics

- **Response Time Distribution**: Histograms and percentiles
- **Throughput Analysis**: Requests per second metrics
- **Resource Utilization**: CPU, memory, and network usage
- **Bottleneck Identification**: Performance bottleneck detection

### 3. Failure Analysis

- **Error Categorization**: Group similar failures together
- **Root Cause Analysis**: Identify common failure patterns
- **Trend Analysis**: Track failure rates over time
- **Alerting**: Configure alerts for critical failures

### 4. Custom Dashboards

- **Service Health Overview**: Service status and metrics
- **Test Coverage Reports**: Endpoint coverage analysis
- **Performance Trends**: Historical performance data
- **Custom Metrics**: Business-specific KPIs

### 5. Export and Integration

- **Multiple Formats**: JSON, CSV, HTML, PDF reports
- **API Integration**: Webhook notifications
- **Email Reports**: Scheduled email reports
- **Third-party Tools**: Integration with monitoring tools

## 🛠️ Development

### Project Structure

```
test-framework-prototype/
├── cmd/
│   └── api-server/          # Main application entry point
├── internal/
│   ├── config/              # Configuration management
│   ├── database/            # Database connections
│   ├── models/              # Data models
│   ├── handlers/            # HTTP handlers
│   ├── services/            # Business logic
│   ├── testrunner/          # Test execution engine
│   └── utils/               # Utility functions
├── scripts/                  # Utility scripts and tools
├── docker-compose.dev.yml   # Development environment
├── Dockerfile               # Application container
└── README.md
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./internal/services

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...
```

### Building for Production

```bash
# Build the application
go build -o api-test-framework cmd/api-server/main.go

# Build with optimizations
go build -ldflags="-s -w" -o api-test-framework cmd/api-server/main.go

# Build Docker image
docker build -t api-test-framework .

# Build multi-platform Docker image
docker buildx build --platform linux/amd64,linux/arm64 -t api-test-framework .
```

### Code Quality

```bash
# Run linter
golangci-lint run

# Run formatter
go fmt ./...

# Run vet
go vet ./...

# Check for security issues
gosec ./...
```

## 📊 Monitoring and Logs

### Application Logs

```bash
# View application logs
docker logs api-test-framework-server-dev

# Follow logs in real-time
docker logs -f api-test-framework-server-dev

# View logs for specific service
docker logs api-test-framework-postgres-dev
```

### Database Access

- **pgAdmin**: http://localhost:7530
  - Email: kaustav.karan@xcaliber.health
  - Password: Ka@060604

### Redis CLI

```bash
# Connect to Redis CLI
docker exec -it api-test-framework-redis-dev redis-cli

# Monitor Redis commands
docker exec -it api-test-framework-redis-dev redis-cli monitor
```

### Health Checks

```bash
# Application health
curl http://localhost:8080/health

# Database health
docker exec api-test-framework-postgres-dev pg_isready -U user -d api_test_framework

# Redis health
docker exec api-test-framework-redis-dev redis-cli ping
```

## 🔧 Configuration

### Environment Variables

| Variable         | Description             | Default            | Required |
| ---------------- | ----------------------- | ------------------ | -------- |
| `SERVER_PORT`    | HTTP server port        | 8080               | No       |
| `SERVER_HOST`    | HTTP server host        | 0.0.0.0            | No       |
| `DB_HOST`        | PostgreSQL host         | localhost          | Yes      |
| `DB_PORT`        | PostgreSQL port         | 5432               | No       |
| `DB_USER`        | PostgreSQL username     | -                  | Yes      |
| `DB_PASSWORD`    | PostgreSQL password     | -                  | Yes      |
| `DB_NAME`        | PostgreSQL database     | api_test_framework | No       |
| `DB_SSL_MODE`    | PostgreSQL SSL mode     | disable            | No       |
| `REDIS_HOST`     | Redis host              | localhost          | Yes      |
| `REDIS_PORT`     | Redis port              | 6379               | No       |
| `REDIS_PASSWORD` | Redis password          | -                  | No       |
| `REDIS_DB`       | Redis database          | 0                  | No       |
| `LOG_LEVEL`      | Logging level           | debug              | No       |
| `ENVIRONMENT`    | Application environment | development        | No       |

### Logging Configuration

The framework supports multiple log levels:

- **debug**: Detailed debugging information
- **info**: General information messages
- **warn**: Warning messages
- **error**: Error messages only

### Performance Tuning

```env
# Database connection pool settings
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_CONN_MAX_LIFETIME=300s

# Redis connection settings
REDIS_POOL_SIZE=10
REDIS_MIN_IDLE_CONNS=5

# Server settings
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s
SERVER_IDLE_TIMEOUT=60s
```

## 🚨 Troubleshooting

### Common Issues

1. **Database Connection Failed**

   - Ensure PostgreSQL container is running
   - Check database credentials in `.env` file
   - Verify network connectivity
   - Check container logs: `docker logs api-test-framework-postgres-dev`

2. **Redis Connection Failed**

   - Ensure Redis container is running
   - Check Redis configuration in `.env` file
   - Verify network connectivity
   - Check container logs: `docker logs api-test-framework-redis-dev`

3. **Test Execution Fails**

   - Verify target service is accessible
   - Check test specification JSON format
   - Review assertion syntax
   - Check application logs for detailed error messages

4. **Performance Issues**
   - Monitor database connection pool usage
   - Check Redis memory usage
   - Review test execution logs
   - Monitor system resources

### Debug Mode

Set `LOG_LEVEL=debug` in `.env` file for detailed logging.

### Health Check Endpoints

```bash
# Application health
curl http://localhost:8080/health

# Database health
curl http://localhost:8080/health/db

# Redis health
curl http://localhost:8080/health/redis
```

### Log Analysis

```bash
# Search for errors
docker logs api-test-framework-server-dev | grep -i error

# Search for specific test runs
docker logs api-test-framework-server-dev | grep "test run"

# Monitor real-time logs
docker logs -f api-test-framework-server-dev | grep -E "(error|warn|test)"
```

## 🔒 Security Considerations

### Authentication & Authorization

- **Service-level Authentication**: Each service can have its own auth config
- **Token Management**: Secure storage of authentication tokens
- **Access Control**: Role-based access control for API endpoints
- **Audit Logging**: Track all API access and modifications

### Data Protection

- **Encryption at Rest**: Database encryption for sensitive data
- **Encryption in Transit**: HTTPS/TLS for all communications
- **Data Masking**: Mask sensitive information in logs
- **Access Logging**: Log all data access and modifications

### Network Security

- **Container Isolation**: Docker network isolation
- **Port Restrictions**: Only necessary ports exposed
- **Health Checks**: Regular security health checks
- **Vulnerability Scanning**: Regular security scans

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass (`go test ./...`)
6. Run linting (`golangci-lint run`)
7. Commit your changes (`git commit -m 'Add amazing feature'`)
8. Push to the branch (`git push origin feature/amazing-feature`)
9. Open a Pull Request

### Development Guidelines

- **Code Style**: Follow Go coding standards
- **Testing**: Maintain >80% test coverage
- **Documentation**: Update README and code comments
- **Error Handling**: Proper error handling and logging
- **Security**: Follow security best practices

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **Go Community**: For the excellent Go ecosystem
- **Gin Framework**: For the fast HTTP web framework
- **GORM**: For the powerful ORM library
- **HTTPExpect**: For the comprehensive HTTP testing library
- **PostgreSQL**: For the reliable database
- **Redis**: For the fast caching and queuing

## 📞 Support

For support and questions:

- **Issues**: Create an issue on GitHub
- **Discussions**: Use GitHub Discussions
- **Documentation**: Check this README and inline code comments
- **Community**: Join our community channels

---

**Built with ❤️ by the API Test Framework Team**
