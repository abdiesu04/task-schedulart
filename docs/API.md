# Task Schedulart API Documentation

## Base URL

```
http://localhost:8080/api
```

## Authentication

JWT-based authentication is implemented. Include the JWT token in the Authorization header:

```http
Authorization: Bearer <your_jwt_token>
```

## Rate Limiting

Rate limiting is implemented with the following limits:
- 100 requests per minute for authenticated users
- 20 requests per minute for unauthenticated users

## Pagination

All list endpoints support pagination with the following query parameters:
- `page`: Page number (default: 1)
- `page_size`: Items per page (default: 10, max: 100)

Example:
```http
GET /api/tasks?page=2&page_size=20
```

## Endpoints

### Health Check

```http
GET /health
```

Response:
```json
{
  "status": "healthy"
}
```

### Tasks

#### List Tasks

```http
GET /tasks
```

Query Parameters:
- `status` (optional): Filter by task status (pending, running, completed, failed)
- `priority` (optional): Filter by priority (low, medium, high)
- `tags` (optional): Filter by tags (comma-separated)
- `search` (optional): Search in task name and description
- `sort_by` (optional): Field to sort by (created_at, schedule_time, priority)
- `order` (optional): Sort order (asc, desc)
- `page` (optional): Page number
- `page_size` (optional): Items per page

Response:
```json
{
  "tasks": [
    {
      "id": 1,
      "name": "Example Task",
      "description": "Task description",
      "scheduleTime": "2024-03-20T15:00:00Z",
      "priority": "high",
      "status": "pending",
      "tags": ["important", "deadline"],
      "retryCount": 0,
      "createdAt": "2024-03-19T10:00:00Z",
      "updatedAt": "2024-03-19T10:00:00Z"
    }
  ],
  "pagination": {
    "current_page": 1,
    "page_size": 10,
    "total_items": 45,
    "total_pages": 5
  }
}
```

#### Get Task by ID

```http
GET /tasks/:id
```

Response:
```json
{
  "id": 1,
  "name": "Example Task",
  "description": "Task description",
  "scheduleTime": "2024-03-20T15:00:00Z",
  "priority": "high",
  "status": "pending",
  "tags": ["important", "deadline"],
  "retryCount": 0,
  "createdAt": "2024-03-19T10:00:00Z",
  "updatedAt": "2024-03-19T10:00:00Z"
}
```

#### Create Task

```http
POST /tasks
Content-Type: application/json
```

Request Body:
```json
{
  "name": "Example Task",
  "description": "Task description",
  "scheduleTime": "2024-03-20T15:00:00Z",
  "priority": "high",
  "tags": ["important", "deadline"]
}
```

Response:
```json
{
  "id": 1,
  "name": "Example Task",
  "description": "Task description",
  "scheduleTime": "2024-03-20T15:00:00Z",
  "priority": "high",
  "status": "pending",
  "tags": ["important", "deadline"],
  "retryCount": 0,
  "createdAt": "2024-03-19T10:00:00Z",
  "updatedAt": "2024-03-19T10:00:00Z"
}
```

#### Update Task

```http
PUT /tasks/:id
Content-Type: application/json
```

Request Body:
```json
{
  "name": "Updated Task",
  "description": "Updated description",
  "scheduleTime": "2024-03-21T15:00:00Z",
  "priority": "medium",
  "tags": ["updated", "modified"]
}
```

Response:
```json
{
  "id": 1,
  "name": "Updated Task",
  "description": "Updated description",
  "scheduleTime": "2024-03-21T15:00:00Z",
  "priority": "medium",
  "status": "pending",
  "tags": ["updated", "modified"],
  "retryCount": 0,
  "createdAt": "2024-03-19T10:00:00Z",
  "updatedAt": "2024-03-19T15:00:00Z"
}
```

#### Update Task Status

```http
PUT /tasks/:id/status
Content-Type: application/json
```

Request Body:
```json
{
  "status": "completed"
}
```

Response:
```json
{
  "message": "Task status updated"
}
```

Valid status values:
- `pending`
- `running`
- `completed`
- `failed`

#### Retry Failed Task

```http
POST /tasks/:id/retry
```

Response:
```json
{
  "message": "Task scheduled for retry"
}
```

#### Delete Task

```http
DELETE /tasks/:id
```

Response:
```json
{
  "message": "Task deleted"
}
```

#### Get Tasks by Tags

```http
GET /tasks/tags?tags=important,urgent
```

Response:
```json
[
  {
    "id": 1,
    "name": "Example Task",
    "description": "Task description",
    "scheduleTime": "2024-03-20T15:00:00Z",
    "priority": "high",
    "status": "pending",
    "tags": ["important", "urgent"],
    "retryCount": 0,
    "createdAt": "2024-03-19T10:00:00Z",
    "updatedAt": "2024-03-19T10:00:00Z"
  }
]
```

#### Create Recurring Task

```http
POST /tasks/recurring
Content-Type: application/json
```

Request Body:
```json
{
  "task": {
    "name": "Recurring Task",
    "description": "Task that repeats",
    "priority": "medium",
    "tags": ["recurring", "automated"]
  },
  "pattern": {
    "type": "daily",
    "interval": 1,
    "weekdays": [1,2,3,4,5],
    "endDate": "2024-12-31T23:59:59Z",
    "cronExpr": "0 9 * * 1-5"
  }
}
```

Response:
```json
{
  "id": 1,
  "name": "Recurring Task",
  "description": "Task that repeats",
  "priority": "medium",
  "status": "pending",
  "tags": ["recurring", "automated"],
  "isRecurring": true,
  "recurringConfig": {
    "type": "daily",
    "interval": 1,
    "weekdays": [1,2,3,4,5],
    "endDate": "2024-12-31T23:59:59Z",
    "cronExpr": "0 9 * * 1-5"
  },
  "createdAt": "2024-03-19T10:00:00Z",
  "updatedAt": "2024-03-19T10:00:00Z"
}
```

### WebSocket Events

Connect to WebSocket endpoint:
```
ws://localhost:8080/ws
```

Event Types:
- `task.created`: New task created
- `task.updated`: Task details updated
- `task.deleted`: Task deleted
- `task.status`: Task status changed
- `task.progress`: Task progress updated

Example WebSocket message:
```json
{
  "event": "task.status",
  "data": {
    "id": 1,
    "status": "completed"
  }
}
```

### Metrics

```http
GET /metrics
```

Returns Prometheus-formatted metrics including:
- Total tasks created
- Total tasks completed
- Total tasks failed
- Current tasks processing
- Task duration histogram
- Tasks by status
- Tasks by priority

## Error Responses

All endpoints return error responses in the following format:

```json
{
  "error": "Error message description"
}
```

HTTP Status Codes:
- 200: Success
- 201: Created
- 400: Bad Request
- 401: Unauthorized
- 403: Forbidden
- 404: Not Found
- 429: Too Many Requests
- 500: Internal Server Error

## Best Practices

1. Always include appropriate headers:
   ```http
   Content-Type: application/json
   Authorization: Bearer <token>
   ```

2. Use pagination for large result sets
3. Include error handling for all API calls
4. Use WebSocket for real-time updates
5. Monitor rate limits
6. Handle authentication token expiration 