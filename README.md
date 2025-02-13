# Task Schedulart 📅

[![Go Report Card](https://goreportcard.com/badge/github.com/abdiesu04/task-schedulart)](https://goreportcard.com/report/github.com/abdiesu04/task-schedulart)
[![GoDoc](https://godoc.org/github.com/abdiesu04/task-schedulart?status.svg)](https://godoc.org/github.com/abdiesu04/task-schedulart)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

<div align="center">
    <img src="docs/images/logo.png" alt="Task Schedulart Logo" width="200"/>
    <p>A modern, cloud-native task scheduling and worker system built with Go</p>
</div>

## 🌟 Features

- **🔄 Advanced Task Management**
  - Create, schedule, and manage tasks
  - Priority-based scheduling
  - Task retry mechanism with exponential backoff
  - Smart tag-based organization
  - Task dependencies and workflows
  - Recurring tasks (cron-style scheduling)
  - Task templates for quick creation

- **☁️ Enterprise-Ready**
  - Built with modern Go practices
  - Docker and Kubernetes ready
  - Horizontally scalable
  - Rate limiting and throttling
  - Circuit breaker pattern
  - Distributed task locking

- **📊 Advanced Features**
  - RESTful API with OpenAPI/Swagger
  - Real-time WebSocket updates
  - Flexible filtering and sorting
  - Rich metadata support
  - Task progress tracking
  - Batch operations
  - Export to CSV/Excel

- **🔒 Security & Monitoring**
  - JWT Authentication
  - Role-based access control
  - Audit logging
  - Prometheus metrics
  - Grafana dashboards
  - Health checks
  - Performance analytics

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL 12+
- Docker (optional)
- Kubernetes (optional)

### Installation

1. Clone the repository:
```bash
git clone https://github.com/abdiesu04/task-schedulart.git
cd task-schedulart
```

2. Set up environment variables:
```bash
export DB_HOST=localhost
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=task_schedulart
export DB_PORT=5432
```

3. Run the application:
```bash
go run main.go
```

### Docker Deployment

```bash
# Build the image
docker build -t task-schedulart:latest .

# Run the container
docker run -p 8080:8080 \
  -e DB_HOST=host.docker.internal \
  -e DB_USER=postgres \
  -e DB_PASSWORD=postgres \
  task-schedulart:latest
```

## 📖 API Documentation

### Task Endpoints

#### Create Task
```http
POST /api/tasks
Content-Type: application/json

{
  "name": "Example Task",
  "description": "Task description",
  "scheduleTime": "2024-03-20T15:00:00Z",
  "priority": "high",
  "tags": ["important", "deadline"]
}
```

#### List Tasks
```http
GET /api/tasks?status=pending&priority=high
```

See [API Documentation](docs/API.md) for complete details.

## 🏗️ Architecture

<div align="center">
    <img src="docs/images/architecture.png" alt="Architecture Diagram" width="600"/>
</div>

The system consists of:
- RESTful API Server
- PostgreSQL Database
- Task Scheduler
- Worker Nodes

## 🧪 Testing

Run the test suite:
```bash
go test ./... -v
```

## 📊 Monitoring

The application exposes metrics at `/metrics` for Prometheus integration.

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

See [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## 📜 License

This project is licensed under the MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io)
- [Zap Logger](https://github.com/uber-go/zap)

## 📞 Support

- Create an issue
- Join our [Discord community](https://discord.gg/yourdiscord)
- Email: support@yourproject.com

## 🗺️ Roadmap

- [x] Basic Task Management
- [x] Docker Support
- [ ] WebSocket Real-time Updates
- [ ] Task Dependencies
- [ ] Recurring Tasks
- [ ] User Authentication
- [ ] Task Templates
- [ ] Mobile App
- [ ] Email Notifications
- [ ] Slack Integration
- [ ] Analytics Dashboard
- [ ] Task Import/Export
- [ ] API Rate Limiting
- [ ] Task Comments & Collaboration
- [ ] SLA Monitoring
- [ ] Multi-tenant Support 