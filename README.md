# Task Schedulart

<div align="center">

[![Go Report Card](https://goreportcard.com/badge/github.com/abdiesu04/task-schedulart)](https://goreportcard.com/report/github.com/abdiesu04/task-schedulart)
[![GoDoc](https://godoc.org/github.com/abdiesu04/task-schedulart?status.svg)](https://godoc.org/github.com/abdiesu04/task-schedulart)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Contributions Welcome](https://img.shields.io/badge/contributions-welcome-brightgreen.svg?style=flat)](CONTRIBUTING.md)
[![GitHub Stars](https://img.shields.io/github/stars/abdiesu04/task-schedulart.svg)](https://github.com/abdiesu04/task-schedulart/stargazers)
[![GitHub Issues](https://img.shields.io/github/issues/abdiesu04/task-schedulart.svg)](https://github.com/abdiesu04/task-schedulart/issues)

<br />
<div align="center">
  <a href="https://github.com/abdiesu04/task-schedulart">
    <img src="https://images.unsplash.com/photo-1611224923853-80b023f02d71?w=600&auto=format&fit=crop&q=60&ixlib=rb-4.0.3" alt="Task Schedulart" width="200" height="200" style="border-radius: 20px;">
  </a>

  <h3 align="center">Task Schedulart</h3>

  <p align="center">
    A modern, cloud-native task scheduling and worker system built with Go
    <br />
    <a href="#documentation"><strong>Explore the docs »</strong></a>
    <br />
    <br />
    <a href="#quick-start">Quick Start</a>
    ·
    <a href="https://github.com/abdiesu04/task-schedulart/issues">Report Bug</a>
    ·
    <a href="https://github.com/abdiesu04/task-schedulart/issues">Request Feature</a>
  </p>
</div>

</div>

## 📋 Table of Contents

- [About](#about)
- [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
- [Architecture](#architecture)
- [API Documentation](#api-documentation)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)
- [Acknowledgments](#acknowledgments)

## 🚀 About

Task Schedulart is a robust, scalable task scheduling system designed for modern cloud environments. Built with Go, it provides real-time task management, monitoring, and execution capabilities with a focus on reliability and performance.

### Why Task Schedulart?

- **🎯 Simple yet Powerful**: Easy to use but packed with advanced features
- **⚡ Real-time Updates**: WebSocket support for instant task status updates
- **🔒 Secure**: Built-in authentication and authorization
- **📊 Observable**: Comprehensive metrics and monitoring
- **🌐 Cloud-Native**: Designed for containerized environments

## ✨ Features

### Core Features

- **🔄 Advanced Task Management**
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

## 🚀 Getting Started

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

## 📖 Documentation

### API Examples

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
  <img src="https://images.unsplash.com/photo-1633356122544-f134324a6cee?w=800&auto=format&fit=crop&q=60&ixlib=rb-4.0.3" alt="Architecture Diagram" width="600" style="border-radius: 10px;">
</div>

The system consists of:
- RESTful API Server
- PostgreSQL Database
- Task Scheduler
- Worker Nodes
- WebSocket Server
- Metrics Collector

## 🤝 Contributing

We love contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📞 Contact

Abdirahman Yusuf - [@abdiesu04](https://twitter.com/abdiesu04)

Project Link: [https://github.com/abdiesu04/task-schedulart](https://github.com/abdiesu04/task-schedulart)

## 🙏 Acknowledgments

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io)
- [Zap Logger](https://github.com/uber-go/zap)
- [Prometheus](https://prometheus.io)
- [WebSocket](https://github.com/gorilla/websocket)

## 🗺️ Roadmap

- [x] Basic Task Management
- [x] Docker Support
- [x] WebSocket Real-time Updates
- [x] Task Dependencies
- [x] Recurring Tasks
- [x] User Authentication
- [ ] Mobile App
- [ ] Email Notifications
- [ ] Slack Integration
- [ ] Analytics Dashboard
- [ ] Task Import/Export
- [ ] API Rate Limiting
- [ ] Task Comments & Collaboration
- [ ] SLA Monitoring
- [ ] Multi-tenant Support