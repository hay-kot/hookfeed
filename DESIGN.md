# HookFeed Design Document

## Table of Contents

1. [System Overview](#system-overview)
2. [Core Concepts](#core-concepts)
3. [Data Models](#data-models)
4. [YAML Configuration](#yaml-configuration)
5. [Middleware System](#middleware-system)
6. [Adapter System](#adapter-system)
7. [Message Processing Pipeline](#message-processing-pipeline)
8. [API Endpoints](#api-endpoints)
9. [WebSocket Protocol](#websocket-protocol)
10. [Database Schema](#database-schema)

---

## System Overview

HookFeed is a self-hosted webhook aggregation and display system that accepts webhooks in any format and displays them in organized feeds. The system uses Lua scripting for flexible message transformation and supports multiple webhook formats through an adapter system.

### Key Features

- **Flexible webhook ingestion** - Accept any webhook format
- **Lua middleware pipeline** - Transform and enrich messages
- **Adapter system** - Built-in support for popular webhook formats (Discord, Ntfy)
- **Infrastructure as Code** - All configuration defined in YAML
- **Real-time updates** - WebSocket-based live feed updates
- **Message management** - Search, filter, and manage message state
- **Retention policies** - Automatic cleanup based on count and age

### Architecture

```
┌─────────────────┐
│  Webhook POST   │
│  /hooks/:slug   │
└────────┬────────┘
         │
         ▼
┌─────────────────────────────────────────┐
│         Processing Pipeline             │
│                                         │
│  1. Load Feed Config                    │
│  2. Execute Global Middleware (ordered) │
│  3. Execute Feed Middleware (ordered)   │
│  4. Apply Adapter (if configured)       │
│  5. Validate & Save Message             │
│  6. Broadcast via WebSocket             │
└─────────────────────────────────────────┘
         │
         ▼
┌─────────────────┐      ┌──────────────┐
│   PostgreSQL    │◄─────┤  WebSocket   │
│   (Messages)    │      │    Hub       │
└─────────────────┘      └──────────────┘
```

### Design Goals

1. **Infrastructure as Code First** - All configuration managed via YAML files
2. **Flexibility** - Support any webhook format through Lua scripting
3. **Simplicity** - Easy to understand and deploy
4. **Performance** - Handle high-throughput webhook ingestion
5. **Observability** - Clear logging and debugging capabilities
6. **Extensibility** - Easy to add new adapters and middleware

---

## Core Concepts

### Feed

A **Feed** is a named collection of webhook messages. Each feed has a unique slug used in the webhook URL and can be configured with specific middleware, adapters, and retention policies.

### Middleware

**Middleware** are Lua scripts that transform webhook payloads. They execute in order and can:

- Modify the payload
- Add metadata
- Control processing flow (abort, skip, bypass remaining middleware)
- Add logs for debugging

### Adapter

An **Adapter** recognizes and transforms specific webhook formats (e.g., Discord, Ntfy) into HookFeed's message format. Adapters can be versioned and explicitly configured per feed.

### Message

A **Message** is a processed webhook stored in a feed. Messages contain:

- Raw webhook data (immutable)
- Processed fields (title, message, level, logs)
- State tracking (new, acknowledged, resolved, archived)
- Metadata from middleware and adapters

### Infrastructure as Code (IaC)

All feeds and configuration are defined in YAML files and synced to the database via CLI. The UI is read-only for viewing messages and managing message state.

---

## Data Models

### Enums

```go
type MessageLevel string
const (
    LevelInfo    MessageLevel = "info"
    LevelWarning MessageLevel = "warning"
    LevelError   MessageLevel = "error"
    LevelSuccess MessageLevel = "success"
    LevelDebug   MessageLevel = "debug"
)

type MessageState string
const (
    StateNew          MessageState = "new"
    StateAcknowledged MessageState = "acknowledged"
    StateResolved     MessageState = "resolved"
    StateArchived     MessageState = "archived"
)

type MiddlewareAction string
const (
    ActionContinue MiddlewareAction = "continue"  // Proceed normally
    ActionAbort    MiddlewareAction = "abort"     // Stop, don't save
    ActionSkip     MiddlewareAction = "skip"      // Skip remaining middleware
    ActionBypass   MiddlewareAction = "bypass"    // Skip adapters
)
```

### Feed

```go
type Feed struct {
    ID                  uuid.UUID  `json:"id"`
    Name                string     `json:"name"`
    Slug                string     `json:"slug"`
    Description         string     `json:"description"`
    Key                 string     `json:"key"`
    MiddlewareScripts   []string   `json:"middlewareScripts"`
    Adapters            []string   `json:"adapters"`
    RetentionMaxCount   *int       `json:"retentionMaxCount"`
    RetentionMaxAgeDays *int       `json:"retentionMaxAgeDays"`
    CreatedAt           time.Time  `json:"createdAt"`
    UpdatedAt           time.Time  `json:"updatedAt"`
}
```

### Message

```go
type Message struct {
    ID           uuid.UUID       `json:"id"`
    FeedID       uuid.UUID       `json:"feedId"`
    RawRequest   json.RawMessage `json:"rawRequest"`
    RawHeaders   json.RawMessage `json:"rawHeaders"`
    Title        *string         `json:"title"`
    Message      *string         `json:"message"`
    Level        MessageLevel    `json:"level"`
    Logs         []string        `json:"logs"`
    Metadata     json.RawMessage `json:"metadata"`
    State        MessageState    `json:"state"`
    StateChanged *time.Time      `json:"stateChangedAt"`
    ReceivedAt   time.Time       `json:"receivedAt"`
    ProcessedAt  *time.Time      `json:"processedAt"`
    CreatedAt    time.Time       `json:"createdAt"`
}
```

### Global Middleware

```go
type GlobalMiddleware struct {
    ID             uuid.UUID `json:"id"`
    Name           string    `json:"name"`
    Description    string    `json:"description"`
    ScriptPath     string    `json:"scriptPath"`
    ExecutionOrder int       `json:"executionOrder"`
    IsEnabled      bool      `json:"isEnabled"`
    CreatedAt      time.Time `json:"createdAt"`
    UpdatedAt      time.Time `json:"updatedAt"`
}
```

### Lua Processing Context

```go
type LuaContext struct {
    Action  MiddlewareAction `json:"action"`
    Error   *string          `json:"error"`
    Payload LuaPayload       `json:"payload"`
}

type LuaPayload struct {
    Raw      map[string]interface{} `json:"raw"`
    Headers  map[string]string      `json:"headers"`
    Title    *string                `json:"title"`
    Message  *string                `json:"message"`
    Level    MessageLevel           `json:"level"`
    Logs     []string               `json:"logs"`
    Metadata map[string]interface{} `json:"metadata"`
}
```

---

## YAML Configuration

### Main Configuration File

```yaml
# hookfeed.yaml
version: 1

# Global middleware (applies to all feeds)
globalMiddleware:
  - name: "Request Logger"
    script: "middleware/logger.lua"
    order: 1
    enabled: true

  - name: "Rate Limiter"
    script: "middleware/rate_limit.lua"
    order: 2
    enabled: true

# Feed definitions
feeds:
  - name: "Production Alerts"
    slug: "prod-alerts"
    key: "sk_prod_abc123xyz789" # Optional, auto-generated if omitted
    description: "Critical production alerts"

    middleware:
      - "middleware/enrich_alerts.lua"
      - "middleware/severity_mapper.lua"

    adapters:
      - "discord@v2"

    retention:
      maxCount: 10000
      maxAgeDays: 90

  - name: "GitHub Events"
    slug: "github"
    description: "GitHub webhook events"

    middleware:
      - "middleware/github_formatter.lua"

    # Empty array = auto-detect adapter
    adapters: []

    retention:
      maxCount: 5000
      maxAgeDays: 30

  - name: "Internal Notifications"
    slug: "internal"

    # Null = no adapters (raw mode)
    adapters: null

    retention:
      maxAgeDays: 7
```

### Configuration Sync

```bash
# Sync configuration from YAML to database
hookfeed sync --config hookfeed.yaml

# Dry run (show what would change)
hookfeed sync --config hookfeed.yaml --dry-run

# Validate configuration without applying
hookfeed validate --config hookfeed.yaml

# Export current configuration to YAML
hookfeed export --output current-config.yaml
```

---

## Middleware System

### Execution Pipeline

```
Webhook POST
    ↓
Initialize Context
    ↓
Global Middleware (ordered)
    ↓
Feed Middleware (ordered)
    ↓
Apply Adapter (if not bypassed)
    ↓
Save Message
```

### Lua Script Structure

All middleware scripts must implement a `process` function:

```lua
function process(context)
    -- context.action: "continue" | "abort" | "skip" | "bypass"
    -- context.error: string | nil
    -- context.payload: table with raw, headers, title, message, level, logs, metadata

    local payload = context.payload

    -- Access raw webhook data
    local eventType = payload.raw.type

    -- Modify message fields
    payload.title = "My Title"
    payload.message = "My Message"
    payload.level = "info"

    -- Add logs
    table.insert(payload.logs, "Processing completed")

    -- Add metadata
    payload.metadata.customField = "value"

    -- Control flow (optional)
    -- context.action = "continue"
    -- context.action = "abort"
    -- context.action = "skip"
    -- context.action = "bypass"

    -- Report errors (optional)
    -- context.error = "Something went wrong"

    return context
end
```

### Built-in Lua Functions

```lua
-- JSON operations
local obj = json_decode(jsonString)
local str = json_encode(table)

-- Logging helpers
add_log("Debug message")

-- Level helpers
set_level("error")

-- Standard Lua
os.time()
os.date()
string.match()
string.gsub()
string.format()
```

### Control Flow Actions

| Action     | Behavior                                        |
| ---------- | ----------------------------------------------- |
| `continue` | Proceed to next middleware/stage (default)      |
| `abort`    | Stop processing immediately, don't save message |
| `skip`     | Skip remaining middleware in current stage      |
| `bypass`   | Skip adapter processing, proceed to save        |

---

## Adapter System

### Adapter Interface

```go
type Adapter interface {
    Name() string
    Version() string
    Detect(payload map[string]interface{}, headers map[string]string) bool
    Transform(payload map[string]interface{}, headers map[string]string) (*LuaPayload, error)
}
```

### Built-in Adapters

#### Discord (`discord@v2`)

**Detection:**

- User-Agent contains "Discord"
- Payload has `content` or `embeds` fields
- Payload has `username` or `avatarUrl` fields

**Transformation:**

- `content` → `message`
- `embeds[0].title` → `title`
- `embeds[0].description` → `message` (if no content)
- `embeds[0].color` → `level` (via color mapping)
- `username` → `metadata.discordUsername`

**Color Mappings:**

- `15158332` (red) → `error`
- `16776960` (yellow) → `warning`
- `3066993` (green) → `success`
- `3447003` (blue) → `info`

#### Ntfy (`ntfy@v1`)

**Detection:**

- Header `X-Ntfy-ID` exists
- Header `X-Priority` exists

**Transformation:**

- `title` → `title`
- `message` → `message`
- `priority` (1-5) → `level` (via priority mapping)
- `tags[]` → `metadata.tags`

**Priority Mappings:**

- `5` (max) → `error`
- `4` (high) → `warning`
- `3` (default) → `info`
- `2` (low) → `info`
- `1` (min) → `debug`

### Adapter Configuration

```yaml
feeds:
  - name: "Example"
    slug: "example"

    # Option 1: Explicit adapter
    adapters:
      - "discord@v2"

    # Option 2: Multiple adapters (tried in order)
    adapters:
      - "discord@v2"
      - "ntfy@v1"

    # Option 3: Auto-detect
    adapters: []

    # Option 4: No adapters (raw mode)
    adapters: null
```

---

## Message Processing Pipeline

### Processing Flow

1. **Load Feed** - Get feed configuration by slug
2. **Initialize Context** - Create processing context with raw data
3. **Global Middleware** - Execute in order (check for abort/skip)
4. **Feed Middleware** - Execute in order (check for abort/skip)
5. **Apply Adapter** - Transform using configured/detected adapter (unless bypassed)
6. **Save Message** - Persist to database
7. **Broadcast** - Send to WebSocket subscribers
8. **Enforce Retention** - Check and apply retention policies (async)

### Error Handling

- Middleware errors are logged to `message.logs`
- Processing continues unless action is `abort`
- Final message includes all errors for debugging
- Adapter errors are logged but don't stop processing

---

## API Endpoints

### Webhook Ingestion

```
POST /hooks/:slug
```

**Headers:**

- `Content-Type: application/json`
- `X-Hook-Key: <feed_key>` (optional)

**Request:** Any JSON payload

**Response:**

```json
{
  "success": true,
  "messageId": "550e8400-e29b-41d4-a716-446655440000",
  "feedId": "650e8400-e29b-41d4-a716-446655440000"
}
```

**Status Codes:**

- `202` - Accepted
- `400` - Invalid JSON
- `401` - Invalid/missing key
- `404` - Feed not found
- `500` - Processing error

---

### Feed Management (Read-Only)

#### List Feeds

```
GET /api/feeds
```

**Query Parameters:**

- `limit` (default: 50)
- `offset` (default: 0)

**Response:**

```json
{
  "feeds": [
    {
      "id": "650e8400-e29b-41d4-a716-446655440000",
      "name": "Production Alerts",
      "slug": "prod-alerts",
      "description": "Critical production alerts",
      "adapters": ["discord@v2"],
      "retentionMaxCount": 10000,
      "retentionMaxAgeDays": 90,
      "createdAt": "2025-01-15T10:30:00Z",
      "updatedAt": "2025-01-15T10:30:00Z"
    }
  ],
  "total": 15,
  "limit": 50,
  "offset": 0
}
```

#### Get Feed

```
GET /api/feeds/:id
```

**Response:**

```json
{
  "id": "650e8400-e29b-41d4-a716-446655440000",
  "name": "Production Alerts",
  "slug": "prod-alerts",
  "description": "Critical production alerts",
  "key": "sk_prod_***440000",
  "middlewareScripts": ["middleware/enrich_alerts.lua"],
  "adapters": ["discord@v2"],
  "retentionMaxCount": 10000,
  "retentionMaxAgeDays": 90,
  "messageCount": 8437,
  "lastMessageAt": "2025-10-16T14:22:15Z",
  "createdAt": "2025-01-15T10:30:00Z",
  "updatedAt": "2025-01-15T10:30:00Z"
}
```

---

### Message Management

#### List Messages

```
GET /api/feeds/:feedId/messages
```

**Query Parameters:**

- `limit` (default: 50, max: 100)
- `offset` (default: 0)
- `level` (optional)
- `state` (optional)
- `since` (ISO 8601, optional)
- `until` (ISO 8601, optional)
- `search` (optional)

**Response:**

```json
{
  "messages": [
    {
      "id": "750e8400-e29b-41d4-a716-446655440000",
      "feedId": "650e8400-e29b-41d4-a716-446655440000",
      "title": "Deploy to production",
      "message": "Deployment successful",
      "level": "success",
      "state": "new",
      "logs": ["Processing repository: my-app"],
      "metadata": {
        "branch": "main",
        "commitCount": 3
      },
      "receivedAt": "2025-10-16T14:22:15Z",
      "processedAt": "2025-10-16T14:22:15.234Z",
      "createdAt": "2025-10-16T14:22:15.250Z"
    }
  ],
  "total": 8437,
  "limit": 50,
  "offset": 0
}
```

#### Get Message

```
GET /api/messages/:id
```

**Query Parameters:**

- `includeRaw` (bool, default: false)

**Response:**

```json
{
  "id": "750e8400-e29b-41d4-a716-446655440000",
  "feedId": "650e8400-e29b-41d4-a716-446655440000",
  "title": "Deploy to production",
  "message": "Deployment successful",
  "level": "success",
  "state": "new",
  "logs": ["..."],
  "metadata": {},
  "rawRequest": {
    "ref": "refs/heads/main",
    "repository": {
      "name": "my-app"
    }
  },
  "rawHeaders": {
    "Content-Type": "application/json",
    "X-GitHub-Event": "push"
  },
  "receivedAt": "2025-10-16T14:22:15Z",
  "processedAt": "2025-10-16T14:22:15.234Z",
  "createdAt": "2025-10-16T14:22:15.250Z"
}
```

#### Update Message State

```
PATCH /api/messages/:id/state
```

**Request:**

```json
{
  "state": "acknowledged"
}
```

**Response:** Updated message object

**Valid Transitions:**

```
new → acknowledged → resolved → archived
  ↓         ↓           ↓
  └─────────┴───────────┴─────────→ archived
```

#### Delete Message

```
DELETE /api/messages/:id
```

**Response:** `204 No Content`

#### Bulk Update State

```
POST /api/feeds/:feedId/messages/bulk-state
```

**Request:**

```json
{
  "messageIds": ["750e8400-...", "850e8400-..."],
  "state": "resolved"
}
```

**Response:**

```json
{
  "updated": 2,
  "failed": 0
}
```

#### Bulk Delete

```
POST /api/feeds/:feedId/messages/bulk-delete
```

**Request:**

```json
{
  "messageIds": ["..."],
  "filter": {
    "level": "debug",
    "olderThan": "2025-09-01T00:00:00Z"
  }
}
```

**Response:**

```json
{
  "deleted": 150
}
```

---

### Search

```
GET /api/search
```

**Query Parameters:**

- `q` (required) - Search query
- `feedId` (optional)
- `level` (optional)
- `state` (optional)
- `since` (optional)
- `until` (optional)
- `limit` (default: 50)
- `offset` (default: 0)

**Response:** Same as List Messages

---

### Global Middleware (Read-Only)

```
GET /api/middleware
```

**Response:**

```json
{
  "middleware": [
    {
      "id": "850e8400-e29b-41d4-a716-446655440000",
      "name": "Request Logger",
      "description": "Logs all incoming requests",
      "scriptPath": "middleware/logger.lua",
      "executionOrder": 1,
      "isEnabled": true,
      "createdAt": "2025-01-15T10:30:00Z",
      "updatedAt": "2025-01-15T10:30:00Z"
    }
  ]
}
```

---

## WebSocket Protocol

### Connection

```
ws://localhost:8080/ws
```

### Subscribe to Feed

```json
{
  "type": "subscribe",
  "feedId": "650e8400-e29b-41d4-a716-446655440000"
}
```

### Unsubscribe from Feed

```json
{
  "type": "unsubscribe",
  "feedId": "650e8400-e29b-41d4-a716-446655440000"
}
```

### Server Messages

#### New Message

```json
{
  "type": "message",
  "feedId": "650e8400-e29b-41d4-a716-446655440000",
  "data": {
    "id": "750e8400-e29b-41d4-a716-446655440000",
    "title": "New deployment",
    "message": "Deployed to production",
    "level": "success",
    "state": "new",
    "receivedAt": "2025-10-16T14:22:15Z"
  }
}
```

#### Message State Changed

```json
{
  "type": "messageState",
  "feedId": "650e8400-e29b-41d4-a716-446655440000",
  "data": {
    "id": "750e8400-e29b-41d4-a716-446655440000",
    "state": "acknowledged",
    "stateChangedAt": "2025-10-16T14:25:00Z"
  }
}
```

#### Message Deleted

```json
{
  "type": "messageDeleted",
  "feedId": "650e8400-e29b-41d4-a716-446655440000",
  "data": {
    "id": "750e8400-e29b-41d4-a716-446655440000"
  }
}
```

---

## Database Schema

```sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS pg_trgm;

-- Feeds
CREATE TABLE feeds (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT DEFAULT '',
    key VARCHAR(255) UNIQUE NOT NULL,
    middleware_scripts TEXT[] DEFAULT '{}',
    adapters TEXT[] DEFAULT '{}',
    retention_max_count INTEGER,
    retention_max_age_days INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_feeds_slug ON feeds(slug);
CREATE INDEX idx_feeds_key ON feeds(key);

-- Messages
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    raw_request JSONB NOT NULL,
    raw_headers JSONB NOT NULL,
    title VARCHAR(500),
    message TEXT,
    level VARCHAR(20) DEFAULT 'info',
    logs TEXT[] DEFAULT '{}',
    metadata JSONB DEFAULT '{}'::jsonb,
    state VARCHAR(20) DEFAULT 'new',
    state_changed_at TIMESTAMP WITH TIME ZONE,
    received_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    processed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT chk_level CHECK (level IN ('info', 'warning', 'error', 'success', 'debug')),
    CONSTRAINT chk_state CHECK (state IN ('new', 'acknowledged', 'resolved', 'archived'))
);

CREATE INDEX idx_messages_feed_id ON messages(feed_id);
CREATE INDEX idx_messages_received_at ON messages(received_at DESC);
CREATE INDEX idx_messages_level ON messages(level);
CREATE INDEX idx_messages_state ON messages(state);
CREATE INDEX idx_messages_feed_received ON messages(feed_id, received_at DESC);
CREATE INDEX idx_messages_feed_state ON messages(feed_id, state);

-- Full-text search
ALTER TABLE messages ADD COLUMN search_vector tsvector;
CREATE INDEX messages_search_idx ON messages USING GIN(search_vector);

CREATE OR REPLACE FUNCTION messages_search_trigger() RETURNS trigger AS $$
BEGIN
    NEW.search_vector :=
        setweight(to_tsvector('english', COALESCE(NEW.title, '')), 'A') ||
        setweight(to_tsvector('english', COALESCE(NEW.message, '')), 'B') ||
        setweight(to_tsvector('english', COALESCE(array_to_string(NEW.logs, ' '), '')), 'C');
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER messages_search_update
    BEFORE INSERT OR UPDATE ON messages
    FOR EACH ROW EXECUTE FUNCTION messages_search_trigger();

-- Global middleware
CREATE TABLE global_middleware (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT DEFAULT '',
    script_path VARCHAR(500) NOT NULL,
    execution_order INTEGER NOT NULL DEFAULT 0,
    is_enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_middleware_order ON global_middleware(execution_order, is_enabled);

-- Updated_at trigger
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_feeds_updated_at BEFORE UPDATE ON feeds
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_middleware_updated_at BEFORE UPDATE ON global_middleware
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

---

## Future Considerations (v2+)

- Statistics & analytics dashboards
- Advanced retention policies (level-based, archival)
- Multi-user support with RBAC
- Enhanced Lua features (HTTP calls, shared libraries, debugging)
- Message threading/correlation
- File attachments support
- Interactive message actions
