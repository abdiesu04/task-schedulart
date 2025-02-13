# Task Schedulart API Documentation

## Base URL

```
http://localhost:8080/api
```

## Authentication

JWT-based authentication is required for most endpoints. Include the JWT token in the Authorization header:

```http
Authorization: Bearer <your_jwt_token>
```

### Authentication Endpoints

#### Register User

```http
POST /auth/register
Content-Type: application/json
```

Request Body:
```json
{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "secure_password123"
}
```

Response:
```json
{
  "id": 1,
  "username": "john_doe",
  "email": "john@example.com",
  "role": "user",
  "createdAt": "2024-03-19T10:00:00Z"
}
```

#### Login

```http
POST /auth/login
Content-Type: application/json
```

Request Body:
```json
{
  "username": "john_doe",
  "password": "secure_password123"
}
```

Response:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
  "token_type": "Bearer",
  "expires_in": 86400
}
```

#### Refresh Token

```http
POST /auth/refresh
Content-Type: application/json
Authorization: Bearer <refresh_token>
```

Response:
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "token_type": "Bearer",
  "expires_in": 86400
}
```

## Rate Limiting

Rate limiting is implemented using a token bucket algorithm with the following limits:
- Authenticated users: 100 requests per minute
- Unauthenticated users: 20 requests per minute

Rate limit headers are included in all responses:
```http
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 95
X-RateLimit-Reset: 1679233922
```

When rate limit is exceeded:
```json
{
  "error": "Rate limit exceeded",
  "retry_after": "60s"
}
```

## Role-Based Access Control

The following roles are available:
- `admin`: Full access to all endpoints
- `user`: Access to own tasks and limited operations
- `viewer`: Read-only access to tasks

Required roles are specified in the documentation for each endpoint.

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
5. Monitor rate limits using response headers
6. Handle authentication token expiration
7. Implement token refresh before expiration
8. Use appropriate roles for different operations
9. Store tokens securely
10. Log out when finished

## Security Considerations

1. **Token Storage**:
   - Store tokens in secure HTTP-only cookies or secure local storage
   - Never store tokens in plain text
   - Clear tokens on logout

2. **Password Requirements**:
   - Minimum 8 characters
   - Mix of uppercase and lowercase letters
   - Include numbers and special characters
   - No common dictionary words

3. **Rate Limiting**:
   - Implement exponential backoff when retrying
   - Monitor rate limit headers
   - Pre-emptively wait when approaching limits

4. **Error Handling**:
   - Never expose sensitive information in error messages
   - Log security-related errors
   - Implement proper validation

5. **HTTPS**:
   - Always use HTTPS in production
   - Implement HSTS
   - Use secure cookies

### AI and Analytics

#### Get Task Analytics

```http
GET /api/analytics/tasks
```

Query Parameters:
- `start_date`: Start date for analysis (ISO 8601)
- `end_date`: End date for analysis (ISO 8601)

Response:
```json
{
  "completionRate": 85.5,
  "averageCompletion": 120.5,
  "productivityHours": {
    "9": 0.95,
    "10": 0.87,
    "14": 0.92
  },
  "tagPerformance": {
    "urgent": 78.5,
    "feature": 92.3
  },
  "priorityDistribution": {
    "high": 45,
    "medium": 30,
    "low": 25
  },
  "trendAnalysis": [
    {
      "date": "2024-03-01T00:00:00Z",
      "value": 82.5
    }
  ]
}
```

#### Optimize Task Schedule

```http
POST /api/tasks/optimize
```

Response:
```json
{
  "optimized": true,
  "changes": [
    {
      "taskId": 1,
      "oldSchedule": "2024-03-19T10:00:00Z",
      "newSchedule": "2024-03-19T14:00:00Z",
      "reason": "Better productivity hour match"
    }
  ]
}
```

### Team Collaboration

#### Create Team

```http
POST /api/teams
Content-Type: application/json
```

Request Body:
```json
{
  "name": "Development Team",
  "description": "Main development team"
}
```

Response:
```json
{
  "id": 1,
  "name": "Development Team",
  "description": "Main development team",
  "createdAt": "2024-03-19T10:00:00Z",
  "updatedAt": "2024-03-19T10:00:00Z"
}
```

#### Invite Team Member

```http
POST /api/teams/:id/invite
Content-Type: application/json
```

Request Body:
```json
{
  "userId": 2,
  "role": "member"
}
```

Response:
```json
{
  "message": "Invitation sent successfully"
}
```

#### Add Task Comment

```http
POST /api/tasks/:id/comments
Content-Type: application/json
```

Request Body:
```json
{
  "content": "Great progress! @john please review",
  "mentions": [2]
}
```

Response:
```json
{
  "id": 1,
  "taskId": 123,
  "userId": 1,
  "content": "Great progress! @john please review",
  "mentions": [2],
  "createdAt": "2024-03-19T10:00:00Z",
  "updatedAt": "2024-03-19T10:00:00Z"
}
```

#### Get Task Comments

```http
GET /api/tasks/:id/comments
```

Response:
```json
[
  {
    "id": 1,
    "taskId": 123,
    "userId": 1,
    "content": "Great progress! @john please review",
    "mentions": [2],
    "createdAt": "2024-03-19T10:00:00Z",
    "updatedAt": "2024-03-19T10:00:00Z"
  }
]
```

### Notifications

#### Configure Notification Channel

```http
POST /api/notifications/channels
Content-Type: application/json
```

Request Body:
```json
{
  "type": "slack",
  "config": {
    "webhookUrl": "https://hooks.slack.com/...",
    "channel": "#tasks"
  },
  "enabled": true
}
```

Response:
```json
{
  "message": "Channel configured successfully"
}
```

#### Create Notification Template

```http
POST /api/notifications/templates
Content-Type: application/json
```

Request Body:
```json
{
  "type": "task_completed",
  "subject": "Task Completed: {{task.name}}",
  "template": "Task {{task.name}} was completed by {{user.name}} at {{timestamp}}"
}
```

Response:
```json
{
  "id": 1,
  "type": "task_completed",
  "subject": "Task Completed: {{task.name}}",
  "template": "Task {{task.name}} was completed by {{user.name}} at {{timestamp}}"
}
```

### Activity Tracking

#### Get Task Activity

```http
GET /api/tasks/:id/activity
```

Response:
```json
[
  {
    "id": 1,
    "taskId": 123,
    "userId": 1,
    "action": "status_changed",
    "details": "Status changed from pending to completed",
    "timestamp": "2024-03-19T10:00:00Z"
  }
]
``` 