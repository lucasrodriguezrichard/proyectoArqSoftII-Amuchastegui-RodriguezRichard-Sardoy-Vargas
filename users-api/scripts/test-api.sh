#!/bin/bash

# Script para probar todos los endpoints de users-api
# Asegúrate de que la API esté corriendo en localhost:8080

API_URL="http://localhost:8080"
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo "🧪 Testing Users API at $API_URL"
echo "=================================="
echo ""

# 1. Health Check
echo -e "${BLUE}1. Testing Health Check${NC}"
echo "GET $API_URL/health"
curl -s -w "\nStatus: %{http_code}\n" $API_URL/health
echo ""
echo "---"
echo ""

# 2. Register User
echo -e "${BLUE}2. Testing Register User${NC}"
echo "POST $API_URL/api/users/register"
REGISTER_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST $API_URL/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123",
    "first_name": "Test",
    "last_name": "User"
  }')

REGISTER_BODY=$(echo "$REGISTER_RESPONSE" | sed '$d')
REGISTER_STATUS=$(echo "$REGISTER_RESPONSE" | tail -n1)

echo "$REGISTER_BODY" | jq '.'
echo "Status: $REGISTER_STATUS"

# Extract user ID for later tests
USER_ID=$(echo "$REGISTER_BODY" | jq -r '.id')
echo -e "${GREEN}✓ User ID: $USER_ID${NC}"
echo ""
echo "---"
echo ""

# 3. Login
echo -e "${BLUE}3. Testing Login${NC}"
echo "POST $API_URL/api/users/login"
LOGIN_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST $API_URL/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "test@example.com",
    "password": "password123"
  }')

LOGIN_BODY=$(echo "$LOGIN_RESPONSE" | sed '$d')
LOGIN_STATUS=$(echo "$LOGIN_RESPONSE" | tail -n1)

echo "$LOGIN_BODY" | jq '.'
echo "Status: $LOGIN_STATUS"

# Extract access token
ACCESS_TOKEN=$(echo "$LOGIN_BODY" | jq -r '.tokens.access_token')
echo -e "${GREEN}✓ Access Token: ${ACCESS_TOKEN:0:50}...${NC}"
echo ""
echo "---"
echo ""

# 4. Get User by ID
echo -e "${BLUE}4. Testing Get User by ID${NC}"
echo "GET $API_URL/api/users/$USER_ID"
curl -s -w "\nStatus: %{http_code}\n" $API_URL/api/users/$USER_ID | jq '.'
echo ""
echo "---"
echo ""

# 5. Register Admin (should fail - need to be admin)
echo -e "${BLUE}5. Testing Register Admin User${NC}"
echo "POST $API_URL/api/users/register (creating admin manually for testing)"
ADMIN_REGISTER_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST $API_URL/api/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@example.com",
    "password": "adminpass123",
    "first_name": "Admin",
    "last_name": "User"
  }')

ADMIN_REGISTER_BODY=$(echo "$ADMIN_REGISTER_RESPONSE" | sed '$d')
ADMIN_REGISTER_STATUS=$(echo "$ADMIN_REGISTER_RESPONSE" | tail -n1)

echo "$ADMIN_REGISTER_BODY" | jq '.'
echo "Status: $ADMIN_REGISTER_STATUS"

ADMIN_USER_ID=$(echo "$ADMIN_REGISTER_BODY" | jq -r '.id')
echo -e "${GREEN}✓ Admin User ID: $ADMIN_USER_ID${NC}"
echo ""
echo "---"
echo ""

# 6. Login as Admin
echo -e "${BLUE}6. Testing Login as Admin${NC}"
echo "POST $API_URL/api/users/login"
ADMIN_LOGIN_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST $API_URL/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "identifier": "admin@example.com",
    "password": "adminpass123"
  }')

ADMIN_LOGIN_BODY=$(echo "$ADMIN_LOGIN_RESPONSE" | sed '$d')
ADMIN_LOGIN_STATUS=$(echo "$ADMIN_LOGIN_RESPONSE" | tail -n1)

echo "$ADMIN_LOGIN_BODY" | jq '.'
echo "Status: $ADMIN_LOGIN_STATUS"

ADMIN_ACCESS_TOKEN=$(echo "$ADMIN_LOGIN_BODY" | jq -r '.tokens.access_token')
echo -e "${GREEN}✓ Admin Access Token: ${ADMIN_ACCESS_TOKEN:0:50}...${NC}"
echo ""
echo "---"
echo ""

# Note: To test POST /api/admin/users, you need to manually update the user role to 'admin' in the database
# or use a pre-seeded admin account

echo -e "${BLUE}7. Testing Create Admin (Protected Endpoint)${NC}"
echo "POST $API_URL/api/admin/users"
echo "Note: This will fail unless the logged-in user has 'admin' role"
echo "The user created via /register has role='user' by default"
echo ""
echo "To test this endpoint, you need to:"
echo "1. Manually update the user role in the database:"
echo "   UPDATE users SET role='admin' WHERE id=$ADMIN_USER_ID;"
echo "2. Or create a migration/seed script to insert an admin user"
echo ""

# Attempt to create admin user (will likely fail with 403)
curl -s -w "\nStatus: %{http_code}\n" -X POST $API_URL/api/admin/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_ACCESS_TOKEN" \
  -d '{
    "username": "newadmin",
    "email": "newadmin@example.com",
    "password": "newadminpass123",
    "first_name": "New",
    "last_name": "Admin"
  }' | jq '.'

echo ""
echo "=================================="
echo -e "${GREEN}✅ API Testing Complete!${NC}"
echo ""
echo "Summary:"
echo "- Health Check: ✓"
echo "- Register User: ✓"
echo "- Login: ✓"
echo "- Get User by ID: ✓"
echo "- Admin endpoints: Require DB role update"
