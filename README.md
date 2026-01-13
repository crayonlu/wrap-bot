> A bot framework using golang
* Introduction
The framework is a highly available and scalable QQ robot framework
* Dev
It is developed using golang which is known for its high concurrency
* modal
Use onion modal pattern,which is a intuitive way
---------
- **Trace:** Start from 2025.10.12...
- **Status:** 🚧 In Active Development

![Architexture](/docs/Architecture.svg)

> Project Structure

```
bot/
├── cmd/bot/                  # Main application entry point
├── internal/
│   ├── config/               # Configuration management
│   └── admin/               # Admin API and WebSocket
├── pkgs/
│   ├── bot/                  # Core bot engine
│   │   ├── engine.go         # Bot engine with middleware support
│   │   ├── context.go        # Request context
│   │   ├── event.go          # Event types and handlers
│   │   ├── middleware.go     # Built-in middlewares
│   │   └── helper.go         # Helper functions
│   ├── napcat/               # NapCat API client
│   │   ├── client.go         # HTTP client
│   │   ├── websocket.go      # WebSocket client
│   │   └── message.go        # Message operations
│   └── feature/              # Utility features
│       └── scheduler.go      # Task scheduler
├── plugins/                  # Bot plugins
│   ├── plugin.go             # Plugin registration
│   ├── ping.go               # Ping command
│   ├── echo.go               # Echo command
│   ├── help.go               # Help command
│   ├── ai_chat.go            # AI chat integration
│   └── tech_push/            # Tech news push
│       ├── tech_push.go
│       └── handlers/         # Data source handlers
├── configs/                  # Configuration files
├── wrap-web/                # Admin web panel (React + shadcn/ui)
│   ├── src/
│   │   ├── components/       # React components
│   │   ├── pages/           # Page components
│   │   ├── lib/             # Utilities (API, WebSocket)
│   │   ├── stores/          # Zustand state management
│   │   └── hooks/           # Custom React hooks
│   └── package.json
└── .github/workflows/        # CI/CD workflows
```

> Architecture Overview

**Core Components:**
- Engine: Middleware-based request handling with onion model pattern
- Context: Request context carrying event data and response methods
- Event: WebSocket event processing from NapCat
- Middleware: Recovery, Logger, Authentication, RateLimit, etc.

**Integration:**
- NapCat Client: HTTP API calls for sending messages and querying data
- WebSocket Client: Real-time event receiving from NapCat server
- Scheduler: Cron-like task scheduling for periodic operations

**Plugin System:**
- Registration-based plugin loading
- Middleware chain execution
- Support for commands, event handlers, and scheduled tasks

![System](/docs/System.svg)