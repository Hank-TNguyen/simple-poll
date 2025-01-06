#!/usr/bin/env bash
#
# test_polls.sh - A script to test the Poll endpoints in your Go backend.
# 
# Usage:
#   1. Make it executable: chmod +x test_polls.sh
#   2. Run it: ./test_polls.sh

set -euo pipefail

# Update BASE_URL as needed. 
# If your Go code expects the trailing slash, include it here.
BASE_URL="http://localhost:3000/api/polls/"

echo "=================================="
echo "STEP 1: Create a new Poll via POST"
echo "=================================="

echo "About to run:"
echo "curl -s -w \"\\nHTTP_CODE:%{http_code}\" -X POST \"${BASE_URL}\" \\"
echo "  -H \"Content-Type: application/json\" \\"
echo "  -d '{\"title\":\"Test Poll from bash script\",\"description\":\"Testing created_by field\",\"created_by\":100}'"

RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X POST "${BASE_URL}" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Poll from bash script",
    "description": "Testing created_by field",
    "created_by": 100
  }'
)

# Separate the body and the HTTP code
BODY=$(echo "$RESPONSE" | sed -e '/HTTP_CODE:/d')
HTTP_CODE=$(echo "$RESPONSE" | sed -n 's/.*HTTP_CODE:\([0-9]*\).*/\1/p')

echo "HTTP code: $HTTP_CODE"
echo "Response body:"
echo "$BODY"
echo

# If you have jq installed, parse the 'id' field from JSON
POLL_ID=$(echo "$BODY" | jq -r '.id' 2>/dev/null || true)
if [[ -z "${POLL_ID}" || "${POLL_ID}" == "null" ]]; then
  echo "ERROR: Could not parse new poll ID (did the request fail?)."
  exit 1
fi

echo "Created poll with ID: $POLL_ID"
echo

echo "=================================="
echo "STEP 2: List all Polls via GET"
echo "=================================="
echo "About to run:"
echo "curl -s -w \"\\nHTTP_CODE:%{http_code}\" -X GET \"${BASE_URL}\""

RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X GET "${BASE_URL}")
BODY=$(echo "$RESPONSE" | sed -e '/HTTP_CODE:/d')
HTTP_CODE=$(echo "$RESPONSE" | sed -n 's/.*HTTP_CODE:\([0-9]*\).*/\1/p')

echo "HTTP code: $HTTP_CODE"
echo "All polls (JSON):"
echo "$BODY"
echo

echo "============================================"
echo "STEP 3: Retrieve the newly created Poll by ID"
echo "============================================"

ENDPOINT="${BASE_URL}${POLL_ID}"
echo "About to run:"
echo "curl -s -w \"\\nHTTP_CODE:%{http_code}\" -X GET \"${ENDPOINT}\""

RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X GET "${ENDPOINT}")
BODY=$(echo "$RESPONSE" | sed -e '/HTTP_CODE:/d')
HTTP_CODE=$(echo "$RESPONSE" | sed -n 's/.*HTTP_CODE:\([0-9]*\).*/\1/p')

echo "HTTP code: $HTTP_CODE"
echo "Poll details (JSON):"
echo "$BODY"
echo

echo "==================================="
echo "STEP 4: Delete the newly created Poll"
echo "==================================="

echo "About to run:"
echo "curl -s -w \"\\nHTTP_CODE:%{http_code}\" -X DELETE \"${ENDPOINT}\""

RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" -X DELETE "${ENDPOINT}")
BODY=$(echo "$RESPONSE" | sed -e '/HTTP_CODE:/d')
HTTP_CODE=$(echo "$RESPONSE" | sed -n 's/.*HTTP_CODE:\([0-9]*\).*/\1/p')

echo "HTTP code: $HTTP_CODE"
echo "Delete response (JSON):"
echo "$BODY"
echo

if [[ "$HTTP_CODE" == "200" ]]; then
  echo "Successfully deleted poll with ID: $POLL_ID"
else
  echo "ERROR: Could not delete poll with ID: $POLL_ID"
  exit 1
fi

echo "=========================================="
echo "All steps completed successfully."
echo "=========================================="
exit 0