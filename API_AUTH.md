# Authentication & Authorization API Documentation

## Overview

The playground backend now includes a comprehensive authentication and authorization system with:
- JWT-based authentication
- Role-Based Access Control (RBAC)
- User groups with permission inheritance
- Password management
- Admin-only operations

## Authentication Endpoints

### Register a New User
```bash
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "alice",
  "email": "alice@example.com",
  "password": "password123"
}
```

### Login
```bash
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "alice",  # Can also be email
  "password": "password123"
}

Response:
{
  "token": "eyJhbGci...",
  "user": {...}
}
```

### Get Current User Info
```bash
GET /api/v1/auth/me
Authorization: Bearer <token>
```

### Change Password
```bash
POST /api/v1/auth/change-password
Authorization: Bearer <token>
Content-Type: application/json

{
  "old_password": "password123",
  "new_password": "newpassword456"
}
```

### Logout
```bash
POST /api/v1/auth/logout
Authorization: Bearer <token>
```

## User Management Endpoints

### List Users (Authenticated)
```bash
GET /api/v1/users
Authorization: Bearer <token>
```

### Get User by ID (Authenticated)
```bash
GET /api/v1/users/{id}
Authorization: Bearer <token>
```

### Update User (Self or Admin)
```bash
PUT /api/v1/users/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "username": "newalice",
  "email": "newalice@example.com"
}
```

### Delete User (Admin Only)
```bash
DELETE /api/v1/users/{id}
Authorization: Bearer <admin-token>
```

## Issue/Topic Management Endpoints

### List Issues (Public)
```bash
GET /api/v1/issues
```

### Get Issue by ID (Public)
```bash
GET /api/v1/issues/{id}
```

### Create Issue (Authenticated)
```bash
POST /api/v1/issues
Authorization: Bearer <token>
Content-Type: application/json

{
  "user_id": 1,
  "summary": "Implement dark mode",
  "description": "Users want a dark mode option"
}
```

### Update Issue (Owner or Admin)
```bash
PUT /api/v1/issues/{id}
Authorization: Bearer <token>
Content-Type: application/json

{
  "summary": "Updated summary",
  "description": "Updated description"
}
```

### Delete Issue (Owner or Admin)
```bash
DELETE /api/v1/issues/{id}
Authorization: Bearer <token>
```

### Vote on Issue (Authenticated)
```bash
POST /api/v1/issues/{id}/vote
Authorization: Bearer <token>
Content-Type: application/json

{
  "user_id": 1
}
```

### Unvote Issue (Authenticated)
```bash
POST /api/v1/issues/{id}/unvote
Authorization: Bearer <token>
Content-Type: application/json

{
  "user_id": 1
}
```

## Message/Discussion Endpoints

### Create Message (Authenticated)
```bash
POST /api/v1/messages
Authorization: Bearer <token>
Content-Type: application/json

{
  "user_id": 1,
  "content": "This is a message",
  "issue_id": 1  # Optional - links message to an issue/topic
}
```

### List Messages (Authenticated)
```bash
GET /api/v1/messages?issue_id=1&limit=100
Authorization: Bearer <token>
```

### Get Message by ID (Authenticated)
```bash
GET /api/v1/messages/{id}
Authorization: Bearer <token>
```

## Permission Management (Admin Only)

### Create Permission
```bash
POST /api/v1/admin/permissions
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "name": "delete_issue",
  "description": "Can delete issues",
  "resource": "issue",
  "action": "delete"
}
```

### Create Role
```bash
POST /api/v1/admin/roles
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "name": "moderator",
  "description": "Can moderate content"
}
```

### Create User Group
```bash
POST /api/v1/admin/groups
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "name": "moderators",
  "description": "Moderator group"
}
```

### Grant Permission to User
```bash
POST /api/v1/admin/users/{user_id}/permissions
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "permission_id": 1
}
```

### Revoke Permission from User
```bash
DELETE /api/v1/admin/users/{user_id}/permissions
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "permission_id": 1
}
```

### Assign Role to User
```bash
POST /api/v1/admin/users/{user_id}/roles
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "role_id": 1
}
```

### Remove Role from User
```bash
DELETE /api/v1/admin/users/{user_id}/roles
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "role_id": 1
}
```

### Add User to Group
```bash
POST /api/v1/admin/users/{user_id}/groups
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "group_id": 1
}
```

### Remove User from Group
```bash
DELETE /api/v1/admin/users/{user_id}/groups
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "group_id": 1
}
```

### Grant Permission to Role
```bash
POST /api/v1/admin/roles/{role_id}/permissions
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "permission_id": 1
}
```

### Grant Permission to Group
```bash
POST /api/v1/admin/groups/{group_id}/permissions
Authorization: Bearer <admin-token>
Content-Type: application/json

{
  "permission_id": 1
}
```

## Permission System

### Permission Inheritance

Users can have permissions through three mechanisms:
1. **Direct Assignment**: Permissions granted directly to the user
2. **Role Assignment**: Permissions inherited from assigned roles
3. **Group Membership**: Permissions inherited from groups the user belongs to

### Admin Users

Users with `is_admin=true` have all permissions automatically, bypassing individual permission checks.

### Permission Checking

The system checks permissions in this order:
1. Is user an admin? → Grant access
2. Does user have direct permission? → Grant access
3. Do any of user's roles have permission? → Grant access
4. Do any of user's groups have permission? → Grant access
5. Otherwise → Deny access

## Configuration

Add JWT secret to config file:

```yaml
server:
  host: "0.0.0.0"
  port: 8080

auth:
  jwt_secret: "your-secret-key-here"  # Change in production!

database:
  type: "sqlite"
  path: "./playground.db"
```

Or use environment variable:
```bash
JWT_SECRET=your-secret-key ./bin/server
```

## Testing

A test script is provided: `test_auth_system.sh`

```bash
chmod +x test_auth_system.sh
./test_auth_system.sh
```

This tests:
- ✅ User registration
- ✅ User login with JWT
- ✅ Protected endpoints
- ✅ Issue CRUD operations
- ✅ Admin-only operations
- ✅ Permission/role/group management
- ✅ Password changes
- ✅ Authorization checks
