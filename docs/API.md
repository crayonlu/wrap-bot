# Bot Admin API Documentation

## Overview

This document describes all REST API endpoints exposed by the bot's admin backend server. The API is built using the Echo framework and provides functionality for managing the bot, plugins, tasks, configuration, and AI features.

**Base URL**: `http://localhost:8080` (configurable via `SERVER_PORT`)

**Authentication**: Most endpoints require JWT authentication (except `/api/auth/login`)

---

## Authentication

### POST /api/auth/login

Authenticate with the admin panel and receive a JWT token.

**Request Body**:
```json
{
  "username": "admin",
  "password": "your_password"
}
```

**Response** (200 OK):
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Response** (401 Unauthorized):
```json
{
  "error": "invalid credentials"
}
```

**Notes**:
- Default username is `admin` (configurable via `ADMIN_USERNAME`)
- Password must be set via `ADMIN_PASSWORD` environment variable
- The returned token must be included in the `Authorization` header for subsequent requests: `Authorization: Bearer <token>`

---

## Bot Status

### GET /api/status

Get the current status of the bot engine.

**Authentication**: Required

**Response** (200 OK):
```json
{
  "running": true,
  "uptime": 3600,
  "version": "1.0.0",
  "go_version": "go1.23.3"
}
```

**Response** (503 Service Unavailable):
```json
{
  "error": "Engine not available"
}
```

**Fields**:
- `running`: Boolean indicating if the bot is currently running
- `uptime`: Bot uptime in seconds
- `version`: Bot version
- `go_version`: Go runtime version

---

## Plugin Management

### GET /api/plugins

Get all registered plugins and their status.

**Authentication**: Required

**Response** (200 OK):
```json
[
  {
    "name": "ping",
    "enabled": true,
    "description": "Simple ping-pong command"
  },
  {
    "name": "ai_chat",
    "enabled": true,
    "description": "AI conversation plugin"
  },
  {
    "name": "tech_push",
    "enabled": false,
    "description": "Tech news push service"
  }
]
```

**Response** (503 Service Unavailable):
```json
{
  "error": "engine not available"
}
```

### POST /api/plugins/:name/toggle

Enable or disable a specific plugin.

**Authentication**: Required

**URL Parameters**:
- `name` (path parameter): The name of the plugin to toggle

**Response** (200 OK):
```json
{
  "name": "ai_chat",
  "enabled": false
}
```

**Response** (404 Not Found):
```json
{
  "error": "plugin not found"
}
```

**Notes**:
- This endpoint toggles the enabled/disabled state of the plugin
- WebSocket clients will receive a broadcast with updated plugin status

---

## Task Management

### GET /api/tasks

Get all registered scheduled tasks and their status.

**Authentication**: Required

**Response** (200 OK):
```json
[
  {
    "id": "tech_push",
    "name": "Tech Push",
    "schedule": "0 */2 * * *",
    "next_run": "2026-01-13T14:00:00Z07:00",
    "last_run": "2026-01-13T12:00:00Z07:00",
    "status": "active",
    "can_trigger": true,
    "description": "Push tech news to configured groups"
  },
  {
    "id": "rss_push",
    "name": "RSS Push",
    "schedule": "0 */4 * * *",
    "next_run": "2026-01-13T16:00:00Z07:00",
    "last_run": "2026-01-13T12:00:00Z07:00",
    "status": "active",
    "can_trigger": true,
    "description": "Push RSS feeds to configured groups"
  }
]
```

**Response** (500 Internal Server Error):
```json
{
  "error": "scheduler not available"
}
```

**Fields**:
- `id`: Unique task identifier
- `name`: Human-readable task name
- `schedule`: Cron expression for the task schedule
- `next_run`: ISO 8601 timestamp of next scheduled run
- `last_run`: ISO 8601 timestamp of last run (empty if never run)
- `status`: Task status (always "active")
- `can_trigger`: Whether the task can be manually triggered
- `description`: Task description

### POST /api/tasks/:id/trigger

Manually trigger a scheduled task.

**Authentication**: Required

**URL Parameters**:
- `id` (path parameter): The ID of the task to trigger

**Response** (200 OK):
```json
{
  "status": "triggered",
  "task_id": "tech_push"
}
```

**Response** (400 Bad Request):
```json
{
  "error": "task not found or cannot be triggered",
  "id": "invalid_task"
}
```

**Response** (500 Internal Server Error):
```json
{
  "error": "scheduler not available"
}
```

**Notes**:
- WebSocket clients will receive a broadcast with updated task status
- Not all tasks support manual triggering

---

## Configuration Management

### GET /api/config

Get all configuration keys and their current values.

**Authentication**: Required

**Response** (200 OK):
```json
[
  {
    "key": "NAPCAT_HTTP_URL",
    "value": "http://localhost:3000",
    "description": "NapCat HTTP API 地址"
  },
  {
    "key": "AI_ENABLED",
    "value": "true",
    "description": "是否启用 AI 功能"
  },
  {
    "key": "AI_USE_UNIFIED",
    "value": "false",
    "description": "是否使用统一模型"
  },
  {
    "key": "AI_UNIFIED_MODEL",
    "value": "",
    "description": "统一模型名称"
  },
  {
    "key": "AI_TEXT_MODEL",
    "value": "deepseek/deepseek-r1-turbo",
    "description": "Text 模型名称"
  },
  {
    "key": "AI_VISION_MODEL",
    "value": "qwen/qwen3-vl-235b-a22b-thinking",
    "description": "Vision 模型名称"
  }
]
```

**Available Configuration Keys**:

| Key | Description |
|------|-------------|
| `NAPCAT_HTTP_URL` | NapCat HTTP API address |
| `NAPCAT_WS_URL` | NapCat WebSocket address |
| `NAPCAT_HTTP_TOKEN` | NapCat HTTP authentication token |
| `NAPCAT_WS_TOKEN` | NapCat WebSocket authentication token |
| `SERVER_PORT` | Admin backend port |
| `SERVER_ENABLED` | Whether admin backend is enabled |
| `DEBUG` | DEBUG mode |
| `COMMAND_PREFIX` | Command prefix |
| `AI_ENABLED` | Whether AI features are enabled |
| `AI_URL` | AI API address |
| `AI_KEY` | AI API key |
| `AI_USE_UNIFIED` | Whether to use unified model mode |
| `AI_UNIFIED_MODEL` | Unified model name |
| `AI_TEXT_MODEL` | Text model name |
| `AI_VISION_MODEL` | Vision model name |
| `AI_TEMPERATURE` | AI temperature parameter |
| `AI_TOP_P` | AI Top-P parameter |
| `AI_MAX_TOKENS` | AI max tokens |
| `AI_MAX_HISTORY` | AI max history records |
| `AI_IMAGE_DETAIL` | Image processing detail (high/low/auto) |
| `AI_TEXT_MODEL_TOOLS` | Enabled tools for text model (comma-separated) |
| `AI_VISION_MODEL_TOOLS` | Enabled tools for vision model (comma-separated) |
| `SYSTEM_PROMPT_PATH` | System prompt path |
| `ANALYZER_PROMPT_PATH` | Analyzer prompt path |
| `HOT_API_HOST` | Hot API URL |
| `HOT_API_KEY` | Hot API key |
| `RSS_API_HOST` | RSS API URL |
| `TECH_PUSH_GROUPS` | Tech push group IDs (comma-separated) |
| `TECH_PUSH_USERS` | Tech push user IDs (comma-separated) |
| `RSS_PUSH_GROUPS` | RSS push group IDs (comma-separated) |
| `RSS_PUSH_USERS` | RSS push user IDs (comma-separated) |
| `ALLOWED_USERS` | Allowed user IDs (comma-separated) |
| `ALLOWED_GROUPS` | Allowed group IDs (comma-separated) |
| `ADMIN_IDS` | Admin user IDs (comma-separated) |
| `SERP_API_KEY` | SerpAPI key (web search) |
| `WEATHER_API_KEY` | WeatherAPI key (weather query) |

### POST /api/config

Update one or more configuration values.

**Authentication**: Required

**Request Body**:
```json
[
  {
    "key": "AI_ENABLED",
    "value": "true"
  },
  {
    "key": "AI_USE_UNIFIED",
    "value": "true"
  },
  {
    "key": "AI_UNIFIED_MODEL",
    "value": "qwen/qwen2.5-72b-instruct"
  }
]
```

**Response** (200 OK):
```json
{
  "status": "updated",
  "updated_count": 3,
  "updated_keys": {
    "AI_ENABLED": true,
    "AI_USE_UNIFIED": true,
    "AI_UNIFIED_MODEL": true
  }
}
```

**Response** (400 Bad Request):
```json
{
  "error": "invalid request"
}
```

**Response** (500 Internal Server Error):
```json
{
  "error": "failed to read .env"
}
```

**Notes**:
- Updates are written to the `.env` file (or file specified by `APP_ENV_FILE`)
- Environment variables are updated in memory immediately
- Some changes may require a bot restart to take effect

---

## Logging

### GET /api/logs

Get recent log entries.

**Authentication**: Required

**Query Parameters**:
- `level` (optional): Filter logs by level (e.g., "INFO", "WARN", "ERROR")
- `limit` (optional): Maximum number of log entries to return (default: 100)

**Example Request**:
```
GET /api/logs?level=INFO&limit=50
```

**Response** (200 OK):
```json
[
  {
    "timestamp": "2026-01-13T10:30:00Z07:00",
    "level": "INFO",
    "message": "Bot engine initialized",
    "context": {
      "module": "engine"
    }
  },
  {
    "timestamp": "2026-01-13T10:30:05Z07:00",
    "level": "INFO",
    "message": "WebSocket connected to ws://localhost:3001",
    "context": {
      "module": "websocket"
    }
  },
  {
    "timestamp": "2026-01-13T10:31:00Z07:00",
    "level": "ERROR",
    "message": "API request failed: connection timeout",
    "context": {
      "module": "provider",
      "endpoint": "/v1/chat/completions"
    }
  }
]
```

**Fields**:
- `timestamp`: ISO 8601 timestamp
- `level`: Log level (DEBUG, INFO, WARN, ERROR)
- `message`: Log message
- `context`: Optional additional context information

---

## Presets Management

### GET /api/presets

Get all available preset files.

**Authentication**: Required

**Response** (200 OK):
```json
[
  {
    "name": "system_prompt.md",
    "path": "configs/system_prompt.md",
    "content": "You are a helpful AI assistant..."
  },
  {
    "name": "analyzer_prompt.md",
    "path": "configs/analyzer_prompt.md",
    "content": "Analyze the following content..."
  }
]
```

**Notes**:
- Presets are markdown files containing system prompts or configuration templates
- Available presets are configured via environment variables

### GET /api/presets/:filename

Get a specific preset file content.

**Authentication**: Required

**URL Parameters**:
- `filename` (path parameter): The name of the preset file (e.g., "system_prompt.md")

**Response** (200 OK):
```json
{
  "name": "system_prompt.md",
  "path": "configs/system_prompt.md",
  "content": "You are a helpful AI assistant..."
}
```

**Response** (400 Bad Request):
```json
{
  "error": "Invalid filename"
}
```

**Response** (404 Not Found):
```json
{
  "error": "Preset not found"
}
```

**Notes**:
- Filename must end with `.md` extension
- Only files configured in environment variables are accessible

### PUT /api/presets/:filename

Update a preset file content.

**Authentication**: Required

**URL Parameters**:
- `filename` (path parameter): The name of the preset file to update

**Request Body**:
```json
{
  "content": "You are a helpful AI assistant that provides accurate and concise responses."
}
```

**Response** (200 OK):
```json
{
  "message": "Preset updated successfully"
}
```

**Response** (400 Bad Request):
```json
{
  "error": "Invalid filename"
}
```

**Response** (404 Not Found):
```json
{
  "error": "Preset not found"
}
```

**Response** (500 Internal Server Error):
```json
{
  "error": "Failed to update preset: permission denied"
}
```

**Notes**:
- Creates parent directories if they don't exist
- Overwrites existing file content
- Changes take effect immediately for new conversations

---

## AI Features

### GET /api/ai/tools

Get all available AI tools and their status.

**Authentication**: Required

**Response** (200 OK):
```json
[
  {
    "name": "get_current_time",
    "description": "获取当前时间",
    "category": "both",
    "text_enabled": true,
    "vision_enabled": true
  },
  {
    "name": "parse_relative_time",
    "description": "解析相对时间表达式（如'3天后'）",
    "category": "both",
    "text_enabled": true,
    "vision_enabled": true
  },
  {
    "name": "web_search",
    "description": "网络搜索",
    "category": "both",
    "text_enabled": true,
    "vision_enabled": true
  },
  {
    "name": "get_weather",
    "description": "获取当前天气",
    "category": "both",
    "text_enabled": true,
    "vision_enabled": true
  },
  {
    "name": "get_weather_forecast",
    "description": "获取天气预报",
    "category": "both",
    "text_enabled": true,
    "vision_enabled": true
  }
]
```

**Fields**:
- `name`: Tool name
- `description`: Tool description
- `category`: Tool category (`text`, `vision`, or `both`)
- `text_enabled`: Whether the tool is enabled for text model
- `vision_enabled`: Whether the tool is enabled for vision model

**Note**: Tool enablement is controlled by the `AI_TEXT_MODEL_TOOLS` and `AI_VISION_MODEL_TOOLS` environment variables. If these variables are empty, all tools are enabled by default.

**Tool Categories**:
- `text`: Tools available for text-only conversations
- `vision`: Tools available for vision conversations
- `both`: Tools available for both text and vision conversations

**Notes**:
- Tools are automatically matched to models based on their category
- Tool categories are defined in tool registration
- `text_enabled` indicates whether tool is enabled for text model
- `vision_enabled` indicates whether tool is enabled for vision model

### GET /api/ai/stats

Get AI usage statistics.

**Authentication**: Required

**Response** (200 OK):
```json
{
  "total_calls": 1250,
  "tool_usage": {
    "web_search": 450,
    "get_weather": 320,
    "get_current_time": 480
  },
  "success_rate": 98.5,
  "recent_calls": [
    {
      "timestamp": "2026-01-13T10:30:00Z07:00",
      "model": "deepseek/deepseek-r1-turbo",
      "tools_used": ["web_search"],
      "success": true,
      "duration_ms": 1250
    },
    {
      "timestamp": "2026-01-13T10:31:00Z07:00",
      "model": "qwen/qwen3-vl-235b-a22b-thinking",
      "tools_used": [],
      "success": true,
      "duration_ms": 890
    }
  ]
}
```

**Fields**:
- `total_calls`: Total number of AI API calls
- `tool_usage`: Map of tool names to usage count
- `success_rate`: Percentage of successful calls (0-100)
- `recent_calls`: Array of recent call records

**Notes**:
- Statistics are reset on bot restart
- Recent calls typically include the last 100 calls

### POST /api/ai/chat

Test AI chat functionality with text-only input.

**Authentication**: Required

**Request Body**:
```json
{
  "message": "What is the current time?",
  "conversation_id": "test-conversation-123",
  "model": "deepseek/deepseek-r1-turbo"
}
```

**Request Fields**:
- `message` (required): The user message to send to the AI
- `conversation_id` (optional): Conversation ID for context (auto-generated if not provided)
- `model` (optional): Override the default model for this request

**Response** (200 OK):
```json
{
  "response": "The current time is 2026-01-13 10:30:00",
  "conversation_id": "test-conversation-123",
  "tool_calls": [
    {
      "name": "get_current_time",
      "arguments": "{}"
    }
  ]
}
```

**Response** (400 Bad Request):
```json
{
  "error": "invalid request"
}
```

**Response** (500 Internal Server Error):
```json
{
  "error": "AI API request failed: rate limit exceeded"
}
```

**Notes**:
- This endpoint uses the text model configuration
- Tool calls are included in the response if the AI used any tools
- Conversation history is maintained for the provided conversation_id

### POST /api/ai/chat/image

Test AI chat functionality with image input.

**Authentication**: Required

**Request Body**:
```json
{
  "message": "Describe this image",
  "images": [
    "https://example.com/image1.jpg",
    "https://example.com/image2.jpg"
  ],
  "conversation_id": "test-vision-456",
  "model": "qwen/qwen3-vl-235b-a22b-thinking"
}
```

**Request Fields**:
- `message` (required): The user message to send to the AI
- `images` (required): Array of image URLs to analyze
- `conversation_id` (optional): Conversation ID for context (auto-generated if not provided)
- `model` (optional): Override the default model for this request

**Response** (200 OK):
```json
{
  "response": "The image shows a beautiful sunset over the ocean with orange and purple colors...",
  "conversation_id": "test-vision-456",
  "tool_calls": []
}
```

**Response** (400 Bad Request):
```json
{
  "error": "no images provided"
}
```

**Response** (500 Internal Server Error):
```json
{
  "error": "AI API request failed: invalid image URL"
}
```

**Notes**:
- This endpoint uses the vision model configuration
- Multiple images can be provided in a single request
- Image URLs must be publicly accessible

---

## WebSocket

### GET /ws

WebSocket endpoint for real-time updates.

**Authentication**: Required (via query parameter or token in connection)

**Query Parameters**:
- `token` (optional): JWT token for authentication

**Connection Example**:
```javascript
const ws = new WebSocket('ws://localhost:8080/ws?token=your_jwt_token');

ws.onmessage = (event) => {
  const data = JSON.parse(event.data);
  console.log('Received:', data);
};
```

**Broadcast Events**:

#### Status Update
```json
{
  "type": "status",
  "data": {
    "running": true,
    "uptime": 3600,
    "version": "1.0.0",
    "go_version": "go1.23.3"
  }
}
```

#### Plugin Update
```json
{
  "type": "plugins",
  "data": [
    {
      "name": "ai_chat",
      "enabled": true,
      "description": "AI conversation plugin"
    }
  ]
}
```

#### Task Update
```json
{
  "type": "tasks",
  "data": [
    {
      "id": "tech_push",
      "name": "Tech Push",
      "schedule": "0 */2 * * *",
      "next_run": "2026-01-13T14:00:00Z07:00",
      "last_run": "2026-01-13T12:00:00Z07:00",
      "status": "active",
      "can_trigger": true,
      "description": "Push tech news to configured groups"
    }
  ]
}
```

#### Log Entry
```json
{
  "type": "log",
  "data": {
    "timestamp": "2026-01-13T10:30:00Z07:00",
    "level": "INFO",
    "message": "WebSocket connected",
    "context": {
      "module": "websocket"
    }
  }
}
```

**Notes**:
- WebSocket connection requires valid JWT token
- All broadcast events are sent to all connected clients
- Status broadcasts occur every 3 seconds by default

---

## Error Responses

All endpoints may return the following error responses:

### 400 Bad Request
```json
{
  "error": "invalid request"
}
```

### 401 Unauthorized
```json
{
  "error": "invalid credentials"
}
```

### 404 Not Found
```json
{
  "error": "resource not found"
}
```

### 500 Internal Server Error
```json
{
  "error": "internal server error"
}
```

### 503 Service Unavailable
```json
{
  "error": "service not available"
}
```

---

## Rate Limiting

Currently, there are no explicit rate limits on the API endpoints. However, the underlying AI API may have rate limits that will be reflected in error responses.

---

## CORS

The API supports CORS (Cross-Origin Resource Sharing) for all endpoints. This allows web applications running on different domains to access the API.

---

## Examples

### Complete Authentication Flow

```bash
# 1. Login to get token
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"your_password"}'

# Response: {"token":"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}

# 2. Use token to access protected endpoints
curl http://localhost:8080/api/status \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
```

### Update AI Configuration

```bash
curl -X POST http://localhost:8080/api/config \
  -H "Authorization: Bearer your_token" \
  -H "Content-Type: application/json" \
  -d '[
    {"key":"AI_USE_UNIFIED","value":"true"},
    {"key":"AI_UNIFIED_MODEL","value":"qwen/qwen2.5-72b-instruct"}
  ]'
```

### Test AI Chat

```bash
curl -X POST http://localhost:8080/api/ai/chat \
  -H "Authorization: Bearer your_token" \
  -H "Content-Type: application/json" \
  -d '{
    "message":"What is the capital of France?",
    "conversation_id":"test-123"
  }'
```

### Trigger a Task

```bash
curl -X POST http://localhost:8080/api/tasks/tech_push/trigger \
  -H "Authorization: Bearer your_token"
```

---

## Versioning

The current API version is **v1**. Future versions may introduce breaking changes and will be indicated by URL path (e.g., `/api/v2/...`).

---

## Support

For issues or questions about the API, please refer to the main project documentation or create an issue in the project repository.
