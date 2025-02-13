# Task Schedulart

A cloud-native task scheduler and worker system built with Go, Docker, and Kubernetes.

## Prerequisites

- Go 1.21 or later
- Docker
- Kubernetes (for deployment)

## Local Development

1. Build and run locally:
```bash
go mod download
go run main.go
```

2. Build and run with Docker:
```bash
# Build the Docker image
docker build -t task-schedulart:latest .

# Run the container
docker run -p 8080:8080 task-schedulart:latest
```

## Testing the API

Once the application is running, you can test the health check endpoint:

```bash
curl http://localhost:8080/health
```

You should receive a response like:
```json
{"status":"healthy"}
```

## Features (Coming Soon)

- Task Submission API
- Task Queueing & Processing
- Worker Nodes
- Kubernetes Deployment
- Monitoring & Observability
- Authentication & Authorization 