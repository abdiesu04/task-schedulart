# Task Schedulart Docker Image

A cloud-native task scheduler and worker system built with Go.

## Quick Start

Run the container:
```bash
docker run -p 8080:8080 your-dockerhub-username/task-schedulart:latest
```

## Kubernetes Deployment

1. Save the following YAML to `deployment.yaml`:
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: task-schedulart
spec:
  replicas: 3
  selector:
    matchLabels:
      app: task-schedulart
  template:
    metadata:
      labels:
        app: task-schedulart
    spec:
      containers:
      - name: task-schedulart
        image: your-dockerhub-username/task-schedulart:latest
        ports:
        - containerPort: 8080
```

2. Deploy to Kubernetes:
```bash
kubectl apply -f deployment.yaml
```

## Environment Variables

- `PORT`: HTTP port (default: 8080)

## Health Check

The service provides a health check endpoint at:
```
GET /health
```

Expected response:
```json
{"status": "healthy"}
``` 