# Go Project: Clean Architecture & Automation

A professionally structured Go application implementing Clean Architecture principles with comprehensive testing automation and development workflows.

## 🏗️ Architecture

The project follows a **feature-based modular architecture** with clear separation of concerns:

```
.
├── core/                # Shared components across features
│   ├── config/          # Configuration management
│   ├── nats/            # NATS client connectivity
│   ├── redis/           # Redis client implementation
│   └── utils/           # Common utility functions
├── features/            # Feature modules
│   ├── root/            # Main business feature
│   │   ├── application/ # Business types and DTOs
│   │   ├── business/    # Use cases implementation
│   │   ├── communication/ # Entry points (HTTP, Cron)
│   │   └── infrastructure/ # Redis storage implementation
│   └── side_car/        # Supporting features
│       ├── health/      # Health checks endpoint
│       └── side_car.go  # Sidecar configuration
├── docs/                # API documentation
│   ├── docs.go         # Swagger annotations
│   ├── swagger.json    # OpenAPI specification
│   └── swagger.yaml    # Swagger YAML definition
├── main.go             # Application entry point
└── test/               # Testing infrastructure
    └── unit/
        └── specify.go  # TestContainer specialization
```

### Architectural Highlights:
- **Feature-based modular architecture** with clear boundaries
- **Domain-driven design** principles
- **Separation of concerns** between core, features, and infrastructure
- **TestContainer integration** for reliable unit testing
- **Centralized configuration** management
- **Message bus integration** via NATS
- **Caching layer** with Redis

## 🚀 Quick Start

### Prerequisites
- Go 1.21+
- Docker and Docker Compose
- Taskfile (optional)

### Installation
```bash
# Clone repository
git clone <repository-url>
cd <project-name>

# Install dependencies
go mod download
```

### Running with Taskfile
```bash
# Start all services
task up
```

## 📁 Project Structure

### Core Module
Shared functionality across features:
- **Config** - Centralized configuration management
- **NATS** - Message bus client implementation
- **Redis** - Caching layer client
- **Utils** - Common utility functions

### Feature Modules

#### Root Feature (Main Business Logic)
- **Application** - Business types and data transfer objects
- **Business** - Use cases implementation (core business logic)
- **Communication** - Entry points (HTTP handlers, Cron jobs)
- **Infrastructure** - Redis-based storage implementation

#### SideCar Feature (Supporting Services)
- **Health** - Health check endpoints
- **Future-ready** for metrics and monitoring

### Use Cases
Business logic encapsulated in feature-specific Use Cases:
```go
type MatchUseCase interface {
    FindMatch(ctx context.Context, criteria MatchCriteria) (*Match, error)
    CreateMatch(ctx context.Context, match *entities.Match) error
}
```

## 🧪 Testing

### Unit Tests with TestContainers
```bash
# Run all tests
task test

# Run tests with coverage
task test:coverage

# Integration tests
task test:integration
```

### Testing Features:
- **TestContainers** for isolated infrastructure testing
- **Specialized test configuration** in `test/unit/specify.go`
- **Comprehensive feature coverage**
- **Integration testing** with real dependencies

## 🔧 Automation with Taskfile

### Core Tasks
```bash
task --list  # Show all available tasks

# Development
task up            # Run application in development mode
task build         # Build the application
task lint          # Run code quality checks

# Testing
task test          # Execute unit tests
task test:cover    # Run tests with coverage report

# Documentation
task docs          # Generate API documentation
task docs:serve    # Serve documentation locally
```

### Taskfile Benefits:
- **Standardized development workflow**
- **Consistent commands across team**
- **Automated quality checks**
- **Simplified CI/CD integration**

## 🏭 Configuration

Application configuration via environment variables. See [.env.example](.env.example) for reference:

```bash
# Core Services
NATS_URL=nats://localhost:4222
REDIS_URL=redis://localhost:6379

# Application
HTTP_PORT=8080
LOG_LEVEL=info
ENVIRONMENT=development
```

## 📊 Monitoring & Logging

- **Structured JSON logging**
- **Health checks** available at `/health/check`
- **Ready for Prometheus metrics** integration
- **Future metrics endpoint** planned in sidecar feature

## 🐳 Docker

### Container Management
```bash
# Build and run services
task docker:up

# Build application image
task docker:build

# Stop all services
task docker:down

# Clean up resources
task docker:clean
```

## 🔮 Roadmap

- [ ] Notify about creating new matches via NATS
- [ ] Advanced monitoring dashboard integration
- [ ] Kubernetes deployment manifests
- [ ] Enhanced dynamic support for bots search
- [ ] Metrics collection and visualization
- [ ] Additional sidecar features for observability

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

**Developer Note**: All development and deployment tasks are standardized through Taskfile to ensure consistency across environments. The modular feature-based architecture allows for scalable and maintainable codebase evolution.