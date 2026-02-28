#!/bin/bash

# Configuration
BASE_URL="https://speed-run.onrender.com"
echo "========================================"
echo "🚀 Starting Full API Tests for Speedrun"
echo "========================================"
echo "Target URL: $BASE_URL"
echo ""

# Level 1: Ping
echo "▶️ Level 1: GET /ping"
curl -s -X GET "$BASE_URL/ping" | jq .
echo ""

# Level 2: Echo
echo "▶️ Level 2: POST /echo"
curl -s -X POST "$BASE_URL/echo" \
  -H "Content-Type: application/json" \
  -d '{"message": "hello", "nested": {"key": "val"}, "number": 42}' | jq .
echo ""

# Level 3: CRUD Create (POST /books) -> Now public, returns Raw Object
echo "▶️ Level 3: POST /books (Create Book 1)"
CREATE_RESPONSE=$(curl -s -X POST "$BASE_URL/books" \
  -H "Content-Type: application/json" \
  -d '{"title": "The Go Programming Language", "author": "Alan A. A. Donovan", "year": 2015}')
echo $CREATE_RESPONSE | jq .
BOOK_ID=$(echo $CREATE_RESPONSE | jq -r .id)
echo "Created Book ID: $BOOK_ID"
echo ""

echo "▶️ Level 3: POST /books (Create Book 2)"
curl -s -X POST "$BASE_URL/books" \
  -H "Content-Type: application/json" \
  -d '{"title": "Black Hat Go", "author": "Tom Steele", "year": 2020}' | jq .
echo ""

# Level 5: Auth Guard (Get Token)
echo "▶️ Level 5 (Prep): POST /auth/token (Login)"
TOKEN_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/token" \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "password"}')
TOKEN=$(echo $TOKEN_RESPONSE | jq -r .token)
echo "Received Token: ${TOKEN:0:15}..."
echo ""

# Level 4: CRUD Read All (Protected now because Level 5 was triggered)
echo "▶️ Level 4 & 5: GET /books (Protected List)"
curl -s -X GET "$BASE_URL/books" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

# Level 4: CRUD Read One (Public)
echo "▶️ Level 4: GET /books/:id (Read One)"
curl -s -X GET "$BASE_URL/books/$BOOK_ID" | jq .
echo ""

# Level 6: Search & Paginate (Protected)
echo "▶️ Level 6: GET /books?author=Alan (Search)"
curl -s -X GET "$BASE_URL/books?author=Alan" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

echo "▶️ Level 6: GET /books?page=1&limit=1 (Paginate)"
curl -s -X GET "$BASE_URL/books?page=1&limit=1" \
  -H "Authorization: Bearer $TOKEN" | jq .
echo ""

# Level 4: CRUD Update
echo "▶️ Level 4: PUT /books/:id (Update)"
curl -s -X PUT "$BASE_URL/books/$BOOK_ID" \
  -H "Content-Type: application/json" \
  -d '{"title": "The Go Programming Language (Updated)", "author": "Alan A. A. Donovan", "year": 2024}' | jq .
echo ""

# Level 7: Error Handling (Bad Request)
echo "▶️ Level 7: POST /books (Invalid Payload - 400 Bad Request)"
curl -s -w "\nHTTP Status: %{http_code}\n" -X POST "$BASE_URL/books" \
  -H "Content-Type: application/json" \
  -d '{"year": 2024}' | grep -v 'HTTP Status: 100'
echo ""

# Level 7: Error Handling (Not Found)
echo "▶️ Level 7: GET /books/:id (Not Found - 404)"
curl -s -w "\nHTTP Status: %{http_code}\n" -X GET "$BASE_URL/books/invalid-uuid-0000" | grep -v 'HTTP Status: 100'
echo ""

echo "▶️ Level 7: PUT /books/nonexistent (Not Found - 404)"
curl -s -w "\nHTTP Status: %{http_code}\n" -X PUT "$BASE_URL/books/nonexistent-bot" \
  -H "Content-Type: application/json" \
  -d '{"title": "Ghost Book"}' | grep -v 'HTTP Status: 100'
echo ""

# Level 4: CRUD Delete
echo "▶️ Level 4: DELETE /books/:id (Delete)"
curl -s -i -X DELETE "$BASE_URL/books/$BOOK_ID"
echo ""

echo "========================================"
echo "✅ All automated local tests completed!"
echo "========================================"
