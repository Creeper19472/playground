# Playground - Real-Time Collaboration Backend

A high-performance Go backend that enables real-time user interaction with messaging, issue tracking, voting, and reference management. Designed to handle 100,000 concurrent users with approximately 30,000 actively participating.

## Features

- **Real-time Messaging**: WebSocket-based instant messaging with broadcast capabilities
- **Issue Management**: Users can create, discuss, and propose issue summaries
- **Voting System**: Democratic voting on issue summaries with automatic ranking
- **Reference Support**: Add URLs and references to support arguments
- **High Scalability**: Optimized for 100k concurrent connections
- **RESTful API**: Clean HTTP API for all operations
- **WebSocket Events**: Real-time updates for new messages, issues, and votes

## Architecture

### Components

- **WebSocket Hub**: Manages all active WebSocket connections with efficient broadcasting
- **RESTful API**: HTTP endpoints for CRUD operations
- **SQLite Database**: Persistent storage with optimized indexes (easily swappable to PostgreSQL)
- **Service Layer**: Business logic separation for maintainability
- **Handler Layer**: HTTP/WebSocket request handling

### Tech Stack

- **Language**: Go 1.24+
- **WebSocket**: gorilla/websocket
- **HTTP Router**: gorilla/mux
- **Database**: SQLite (mattn/go-sqlite3)
- **CORS**: rs/cors
- **UUID**: google/uuid

## Quick Start

### Prerequisites

- Go 1.24 or higher
- GCC (for SQLite compilation)

### Installation

```bash
# Clone the repository
git clone https://github.com/Creeper19472/playground.git
cd playground

# Install dependencies
go mod download

# Build the server
go build -o bin/server ./cmd/server

# Run the server
./bin/server
```

The server will start on port 8080 by default. You can override this with the `PORT` environment variable:

```bash
PORT=3000 ./bin/server
```

## API Documentation

### Base URL

```
http://localhost:8080/api/v1
```

### Endpoints

#### Users

- `POST /users` - Create a new user
  ```json
  {
    "username": "john_doe",
    "email": "john@example.com"
  }
  ```

- `GET /users/{id}` - Get user by ID
- `GET /users` - List all users

#### Messages

- `POST /messages` - Send a message
  ```json
  {
    "user_id": 1,
    "content": "Hello, world!",
    "issue_id": 1  // optional
  }
  ```

- `GET /messages/{id}` - Get message by ID
- `GET /messages?issue_id={id}&limit={n}` - List messages (optionally filtered by issue)

#### Issues

- `POST /issues` - Create a new issue
  ```json
  {
    "user_id": 1,
    "summary": "Feature request: Dark mode",
    "description": "We need a dark mode for better UX"
  }
  ```

- `GET /issues/{id}` - Get issue by ID
- `GET /issues` - List all issues (sorted by vote count)
- `POST /issues/{id}/vote` - Vote on an issue
  ```json
  {
    "user_id": 1
  }
  ```
- `POST /issues/{id}/unvote` - Remove vote from an issue

#### References

- `POST /references` - Add a reference to an issue
  ```json
  {
    "user_id": 1,
    "issue_id": 1,
    "url": "https://example.com/article",
    "title": "Supporting Evidence",
    "description": "This article supports our argument"
  }
  ```

- `GET /references/{id}` - Get reference by ID
- `GET /issues/{issue_id}/references` - List all references for an issue

### WebSocket

Connect to the WebSocket endpoint:

```
ws://localhost:8080/ws?user_id={id}&username={name}
```

#### WebSocket Events

The server broadcasts the following event types:

- `connected` - Welcome message when client connects
- `new_message` - New message created
- `new_issue` - New issue created
- `issue_voted` - Issue received a vote
- `issue_unvoted` - Vote removed from issue

Event format:
```json
{
  "type": "event_type",
  "data": { /* event-specific data */ }
}
```

## Performance & Scalability

### Concurrent Connections

The system is designed to handle:
- 100,000 concurrent WebSocket connections
- ~30,000 actively participating users
- Efficient message broadcasting with minimal latency

### Optimizations

1. **Connection Pooling**: Database connection pool configured for high concurrency
2. **Indexed Queries**: Strategic database indexes on frequently queried columns
3. **Efficient Broadcasting**: Non-blocking channel operations prevent slow clients from affecting others
4. **Transaction Management**: ACID-compliant vote counting to prevent race conditions

### Production Deployment

For production use with higher loads, consider:

1. **Database Migration**: Switch from SQLite to PostgreSQL
   ```go
   // Update db/database.go to use PostgreSQL driver
   _ "github.com/lib/pq"
   ```

2. **Redis Caching**: Add Redis for vote count caching
3. **Load Balancing**: Deploy multiple instances behind a load balancer
4. **Message Queue**: Use Redis Pub/Sub or RabbitMQ for distributed broadcasting

## Development

### Running Tests

```bash
go test ./...
```

### Project Structure

```
.
├── cmd/
│   └── server/          # Application entry point
├── internal/
│   ├── db/              # Database layer
│   ├── handlers/        # HTTP/WebSocket handlers
│   ├── models/          # Data models
│   └── services/        # Business logic
├── pkg/
│   └── websocket/       # WebSocket hub implementation
└── go.mod               # Go module definition
```

## Examples

### Creating a User and Posting a Message

```bash
# Create a user
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","email":"alice@example.com"}'

# Post a message
curl -X POST http://localhost:8080/api/v1/messages \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"content":"Hello everyone!"}'
```

### Creating an Issue and Voting

```bash
# Create an issue
curl -X POST http://localhost:8080/api/v1/issues \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"summary":"Add export feature","description":"Users need to export data"}'

# Vote on the issue
curl -X POST http://localhost:8080/api/v1/issues/1/vote \
  -H "Content-Type: application/json" \
  -d '{"user_id":2}'

# Get ranked issues
curl http://localhost:8080/api/v1/issues
```

### WebSocket Connection (JavaScript)

```javascript
const ws = new WebSocket('ws://localhost:8080/ws?user_id=1&username=alice');

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('Received:', message.type, message.data);
};

ws.onopen = () => {
  console.log('Connected to WebSocket');
};
```

## Health Check

```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "ok",
  "connected_clients": 42
}
```

## License

MIT License - See LICENSE file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
