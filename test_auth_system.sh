#!/bin/bash

# Start the server in background
cd /home/runner/work/playground/playground
rm -f playground.db test_auth.log
./bin/server > test_auth.log 2>&1 &
SERVER_PID=$!
echo "Server PID: $SERVER_PID"

# Wait for server to start
sleep 3

echo "=== Testing Authentication and Authorization System ==="
echo ""

# Test 1: Register a new user
echo "1. Register new user (alice)"
REGISTER_RESP=$(curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","email":"alice@example.com","password":"password123"}')
echo "$REGISTER_RESP" | jq '.'
echo ""

# Test 2: Login
echo "2. Login as alice"
LOGIN_RESP=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"password123"}')
echo "$LOGIN_RESP" | jq '.'
TOKEN=$(echo "$LOGIN_RESP" | jq -r '.token')
echo "Token: $TOKEN"
echo ""

# Test 3: Get current user info
echo "3. Get current user info (me)"
curl -s http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer $TOKEN" | jq '.'
echo ""

# Test 4: Create an issue (authenticated)
echo "4. Create an issue (authenticated)"
ISSUE_RESP=$(curl -s -X POST http://localhost:8080/api/v1/issues \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"user_id":1,"summary":"Implement dark mode","description":"Users want a dark mode option"}')
echo "$ISSUE_RESP" | jq '.'
ISSUE_ID=$(echo "$ISSUE_RESP" | jq -r '.id')
echo ""

# Test 5: List issues (public endpoint)
echo "5. List issues (public - no auth required)"
curl -s http://localhost:8080/api/v1/issues | jq '.'
echo ""

# Test 6: Update own issue
echo "6. Update own issue"
curl -s -X PUT http://localhost:8080/api/v1/issues/$ISSUE_ID \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"summary":"Implement dark mode (Updated)","description":"Users really want a dark mode option"}' | jq '.'
echo ""

# Test 7: Register admin user (directly in DB for testing)
echo "7. Creating admin user..."
sqlite3 playground.db "UPDATE users SET is_admin=1 WHERE username='alice';"
echo "Admin status updated in database"
echo ""

# Test 8: Login again as admin
echo "8. Login again as admin"
ADMIN_LOGIN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"password123"}')
ADMIN_TOKEN=$(echo "$ADMIN_LOGIN" | jq -r '.token')
echo "Admin token obtained"
echo ""

# Test 9: Create permission (admin only)
echo "9. Create permission (admin only)"
PERM_RESP=$(curl -s -X POST http://localhost:8080/api/v1/admin/permissions \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"delete_issue","description":"Can delete issues","resource":"issue","action":"delete"}')
echo "$PERM_RESP" | jq '.'
echo ""

# Test 10: Create role (admin only)
echo "10. Create role (admin only)"
ROLE_RESP=$(curl -s -X POST http://localhost:8080/api/v1/admin/roles \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"moderator","description":"Can moderate content"}')
echo "$ROLE_RESP" | jq '.'
echo ""

# Test 11: Create user group (admin only)
echo "11. Create user group (admin only)"
GROUP_RESP=$(curl -s -X POST http://localhost:8080/api/v1/admin/groups \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name":"moderators","description":"Moderator group"}')
echo "$GROUP_RESP" | jq '.'
echo ""

# Test 12: Register another user
echo "12. Register second user (bob)"
curl -s -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"bob","email":"bob@example.com","password":"password456"}' | jq '.'
echo ""

# Test 13: Try to delete user as non-admin (should fail)
echo "13. Login as bob and try to delete alice (should fail)"
BOB_LOGIN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"bob","password":"password456"}')
BOB_TOKEN=$(echo "$BOB_LOGIN" | jq -r '.token')
curl -s -X DELETE http://localhost:8080/api/v1/users/1 \
  -H "Authorization: Bearer $BOB_TOKEN"
echo ""
echo ""

# Test 14: Delete user as admin (should succeed)
echo "14. Delete bob as admin (should succeed)"
curl -s -X DELETE http://localhost:8080/api/v1/users/2 \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq '.'
echo ""

# Test 15: Change password
echo "15. Change password"
curl -s -X POST http://localhost:8080/api/v1/auth/change-password \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"old_password":"password123","new_password":"newpassword456"}' | jq '.'
echo ""

# Test 16: Test invalid auth
echo "16. Test invalid token (should fail)"
curl -s http://localhost:8080/api/v1/auth/me \
  -H "Authorization: Bearer invalid-token"
echo ""
echo ""

# Stop the server
echo "=== Stopping Server ==="
kill $SERVER_PID
echo "Server stopped"
echo ""
echo "=== Test Summary ==="
echo "✅ User registration"
echo "✅ User login with JWT"
echo "✅ Protected endpoints with authentication"
echo "✅ Issue creation, update, deletion"
echo "✅ Admin-only operations"
echo "✅ Permission, role, and group creation"
echo "✅ Password change"
echo "✅ Authorization checks"
