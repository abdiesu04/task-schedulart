# API Documentation

## Base URL

```
https://your-api-url.com/api
```

## Authentication

_Coming soon_

## Endpoints

### Tasks

#### List Tasks

```http
GET /tasks
```

Query Parameters:
- `status` (optional): Filter by task status (pending, running, completed, failed)
- `priority` (optional): Filter by priority (low, medium, high)
- `tags` (optional): Filter by tags (comma-separated)

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
    "tags": ["important", "deadline"],
    "retryCount": 0,
    "createdAt": "2024-03-19T10:00:00Z",
    "updatedAt": "2024-03-19T10:00:00Z"
  }
]
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

## Error Responses

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
- 500: Internal Server Error

## Rate Limiting

_Coming soon_

## Pagination

_Coming soon_ 