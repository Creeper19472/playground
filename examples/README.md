# Demo Application

This directory contains a demo HTML application that showcases all features of the Playground backend.

## Features Demonstrated

### Authentication & Authorization
- User registration with email and password
- JWT-based login/logout
- Password change functionality
- Admin vs regular user roles
- Protected routes and actions

### Real-time Collaboration
- WebSocket connection for live updates
- Real-time messaging
- Live issue updates (creation, voting, editing, deletion)
- Event stream showing all real-time activities

### Issue/Topic Management
- Create issues with summary and description
- Vote and unvote on issues
- Issues ranked by vote count
- Edit issues (owner and admin only)
- Delete issues (owner and admin only)
- Discussion threads per issue

### Admin Features
- User management (view all users)
- Role management (create and view roles)
- Permission management (create and view permissions)
- Admin panel only visible to admin users

## How to Use

### 1. Start the Backend Server

```bash
# Build the server
go build -o bin/server ./cmd/server

# Run with default configuration (SQLite)
./bin/server

# Or run with custom configuration
./bin/server -config config.yaml
```

### 2. Open the Demo

Open `demo.html` in your web browser:

```bash
# Option 1: Direct file access
open examples/demo.html

# Option 2: Via HTTP server (recommended for CORS)
cd examples
python3 -m http.server 8081
# Then navigate to http://localhost:8081/demo.html
```

### 3. Register and Login

1. Click the "Register" tab
2. Enter username, email, and password (confirm password)
3. Click "Register"
4. Switch to "Login" tab
5. Enter your credentials and login

### 4. Try the Features

**Basic User Features:**
- Send messages in the message panel
- Create issues/topics with summary and description
- Vote on issues to increase their ranking
- Edit or delete your own issues

**Admin Features (if you're an admin):**
- View the Admin Panel at the bottom
- See all users in the system
- Create roles and permissions
- Manage user access levels

### 5. Observe Real-time Updates

- Open the demo in multiple browser windows
- Actions in one window appear instantly in others
- Check the "Real-time Events" panel to see WebSocket messages

## API Endpoints Used

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login and get JWT token
- `POST /api/v1/auth/logout` - Logout
- `POST /api/v1/auth/change-password` - Change password

### Issues
- `GET /api/v1/issues` - List all issues (public)
- `POST /api/v1/issues` - Create issue (authenticated)
- `PUT /api/v1/issues/:id` - Update issue (owner/admin)
- `DELETE /api/v1/issues/:id` - Delete issue (owner/admin)
- `POST /api/v1/issues/:id/vote` - Vote on issue
- `DELETE /api/v1/issues/:id/vote` - Remove vote

### Messages
- `GET /api/v1/messages` - List messages (authenticated)
- `POST /api/v1/messages` - Send message (authenticated)

### Admin
- `GET /api/v1/users` - List users (admin only)
- `GET /api/v1/roles` - List roles (admin only)
- `POST /api/v1/roles` - Create role (admin only)
- `GET /api/v1/permissions` - List permissions (admin only)
- `POST /api/v1/permissions` - Create permission (admin only)

### WebSocket
- `ws://localhost:8080/ws?user_id={id}&username={name}` - Real-time connection

## Configuration

The demo is configured to connect to:
- **Backend API**: `http://localhost:8080/api/v1`
- **WebSocket**: `ws://localhost:8080/ws`

To connect to a different backend, edit the constants at the top of `demo.html`:

```javascript
const API_URL = 'http://your-server:port/api/v1';
const AUTH_URL = 'http://your-server:port/api/v1/auth';
const WS_URL = 'ws://your-server:port/ws';
```

## Creating an Admin User

By default, the first registered user is not an admin. To create an admin user:

1. Register a normal user via the demo
2. Update the database directly:

```bash
# SQLite
sqlite3 playground.db "UPDATE users SET is_admin = 1 WHERE username = 'yourusername';"

# PostgreSQL
psql -d playground -c "UPDATE users SET is_admin = true WHERE username = 'yourusername';"

# MySQL
mysql -D playground -e "UPDATE users SET is_admin = 1 WHERE username = 'yourusername';"
```

3. Logout and login again to see the admin panel

## Troubleshooting

### Can't connect to backend
- Ensure the backend server is running on port 8080
- Check browser console for CORS errors
- Verify API_URL matches your server address

### WebSocket not connecting
- Check that the server is running
- Look for WebSocket errors in browser console
- Verify the user_id and username are valid

### 404 errors
- Confirm you're using the correct API paths (`/api/v1/auth/*` not `/api/v1/*`)
- Check server logs for routing issues
- Verify the endpoint exists in the API

### Admin panel not showing
- Ensure you're logged in as an admin user
- Check the `is_admin` field in the database
- Try logging out and back in

## Next Steps

- Explore the backend API documentation in `API_AUTH.md`
- Try the test script: `./test_auth_system.sh`
- Check the configuration examples in `config.example.*.yaml`
- Review the architecture in `README.md`
