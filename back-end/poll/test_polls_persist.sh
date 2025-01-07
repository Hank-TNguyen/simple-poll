#!/usr/bin/env bash
#
# test_polls_persist.sh
#
# Creates a Poll and then two Questions (with separate POST calls).
# Afterward, creates Choices (4 for Question #1 and 2 for Question #2).
# This script does NOT delete the created Poll, Questions, or Choices.
#
# Note: The API does not support batch creation. We must call each endpoint
# (Poll, Question, Choice) individually.

set -euo pipefail

# Adjust as needed; no trailing slash unless your API requires it.
API_BASE_URL="http://localhost:3000"

#
# 1) Create a Poll
#
echo "=================================="
echo "STEP 1: Create a Poll via POST"
echo "=================================="
POLL_PAYLOAD='{
  "title": "Sample Poll (Persistent)",
  "description": "This poll will remain in the database",
  "created_by": 100
}'

# DEBUG: print the curl command about to run
echo "DEBUG: About to run curl command for creating poll:"
echo "curl -s -w \"\\nHTTP_CODE:%{http_code}\" \\
  -X POST \"${API_BASE_URL}/api/polls/\" \\
  -H \"Content-Type: application/json\" \\
  -d \"${POLL_PAYLOAD}\""

RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" \
  -X POST "${API_BASE_URL}/api/polls/" \
  -H "Content-Type: application/json" \
  -d "${POLL_PAYLOAD}"
)

BODY=$(echo "$RESPONSE" | sed -e '/HTTP_CODE:/d')
HTTP_CODE=$(echo "$RESPONSE" | sed -n 's/.*HTTP_CODE:\([0-9]*\).*/\1/p')

echo "HTTP code: $HTTP_CODE"
echo "Response body: $BODY"
echo

if [[ "$HTTP_CODE" -ne 200 && "$HTTP_CODE" -ne 201 ]]; then
  echo "ERROR: Failed to create poll."
  exit 1
fi

# Extract poll ID via jq (install jq if you haven't already: https://stedolan.github.io/jq/)
POLL_ID=$(echo "$BODY" | jq -r '.id' 2>/dev/null || true)
if [[ -z "$POLL_ID" || "$POLL_ID" == "null" ]]; then
  echo "ERROR: Could not parse poll ID."
  exit 1
fi

echo "Created poll with ID: $POLL_ID"
echo


#
# 2) Create Question #1
#
echo "=================================="
echo "STEP 2: Create Question #1"
echo "=================================="

QUESTION1_PAYLOAD=$(cat <<EOF
{
  "poll_id": $POLL_ID,
  "text": "give me one"
}
EOF
)

# DEBUG: print the curl command about to run
echo "DEBUG: About to run curl command for creating Question #1:"
echo "curl -s -w \"\\nHTTP_CODE:%{http_code}\" \\
  -X POST \"${API_BASE_URL}/api/questions/\" \\
  -H \"Content-Type: application/json\" \\
  -d \"${QUESTION1_PAYLOAD}\""

RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" \
  -X POST "${API_BASE_URL}/api/questions/" \
  -H "Content-Type: application/json" \
  -d "${QUESTION1_PAYLOAD}"
)

BODY=$(echo "$RESPONSE" | sed -e '/HTTP_CODE:/d')
HTTP_CODE=$(echo "$RESPONSE" | sed -n 's/.*HTTP_CODE:\([0-9]*\).*/\1/p')

echo "HTTP code: $HTTP_CODE"
echo "Response body: $BODY"
echo

if [[ "$HTTP_CODE" -ne 200 && "$HTTP_CODE" -ne 201 ]]; then
  echo "ERROR: Failed to create Question #1."
  exit 1
fi

QUESTION_1_ID=$(echo "$BODY" | jq -r '.id' 2>/dev/null || true)
if [[ -z "$QUESTION_1_ID" || "$QUESTION_1_ID" == "null" ]]; then
  echo "ERROR: Could not parse Question #1 ID."
  exit 1
fi

echo "Created Question #1 with ID: $QUESTION_1_ID"
echo


#
# 3) Create Question #2
#
echo "=================================="
echo "STEP 3: Create Question #2"
echo "=================================="

QUESTION2_PAYLOAD=$(cat <<EOF
{
  "poll_id": $POLL_ID,
  "text": "give me one"
}
EOF
)

# DEBUG: print the curl command about to run
echo "DEBUG: About to run curl command for creating Question #2:"
echo "curl -s -w \"\\nHTTP_CODE:%{http_code}\" \\
  -X POST \"${API_BASE_URL}/api/questions/\" \\
  -H \"Content-Type: application/json\" \\
  -d \"${QUESTION2_PAYLOAD}\""

RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" \
  -X POST "${API_BASE_URL}/api/questions/" \
  -H "Content-Type: application/json" \
  -d "${QUESTION2_PAYLOAD}"
)

BODY=$(echo "$RESPONSE" | sed -e '/HTTP_CODE:/d')
HTTP_CODE=$(echo "$RESPONSE" | sed -n 's/.*HTTP_CODE:\([0-9]*\).*/\1/p')

echo "HTTP code: $HTTP_CODE"
echo "Response body: $BODY"
echo

if [[ "$HTTP_CODE" -ne 200 && "$HTTP_CODE" -ne 201 ]]; then
  echo "ERROR: Failed to create Question #2."
  exit 1
fi

QUESTION_2_ID=$(echo "$BODY" | jq -r '.id' 2>/dev/null || true)
if [[ -z "$QUESTION_2_ID" || "$QUESTION_2_ID" == "null" ]]; then
  echo "ERROR: Could not parse Question #2 ID."
  exit 1
fi

echo "Created Question #2 with ID: $QUESTION_2_ID"
echo


#
# 4) Create 4 Choices for Question #1
#
echo "=================================="
echo "STEP 4: Create 4 Choices for Question #1"
echo "=================================="
for choice_text in "Choice A" "Choice B" "Choice C" "Choice D"; do
  CHOICE_PAYLOAD=$(cat <<EOF
{
  "question_id": $QUESTION_1_ID,
  "text": "$choice_text"
}
EOF
  )

  # DEBUG: print the curl command about to run
  echo "DEBUG: About to run curl command for creating choice \"$choice_text\" for Question #1:"
  echo "curl -s -w \"\\nHTTP_CODE:%{http_code}\" \\
    -X POST \"${API_BASE_URL}/api/choices/\" \\
    -H \"Content-Type: application/json\" \\
    -d \"${CHOICE_PAYLOAD}\""

  RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" \
    -X POST "${API_BASE_URL}/api/choices/" \
    -H "Content-Type: application/json" \
    -d "${CHOICE_PAYLOAD}"
  )

  BODY=$(echo "$RESPONSE" | sed -e '/HTTP_CODE:/d')
  HTTP_CODE=$(echo "$RESPONSE" | sed -n 's/.*HTTP_CODE:\([0-9]*\).*/\1/p')

  if [[ "$HTTP_CODE" -ne 200 && "$HTTP_CODE" -ne 201 ]]; then
    echo "ERROR: Failed to create choice '$choice_text' for Question #1."
    exit 1
  fi

  echo "Created choice '$choice_text' (HTTP code: $HTTP_CODE)"
done

echo


#
# 5) Create 2 Choices for Question #2
#
echo "=================================="
echo "STEP 5: Create 2 Choices for Question #2"
echo "=================================="
for choice_text in "Choice X" "Choice Y"; do
  CHOICE_PAYLOAD=$(cat <<EOF
{
  "question_id": $QUESTION_2_ID,
  "text": "$choice_text"
}
EOF
  )

  # DEBUG: print the curl command about to run
  echo "DEBUG: About to run curl command for creating choice \"$choice_text\" for Question #2:"
  echo "curl -s -w \"\\nHTTP_CODE:%{http_code}\" \\
    -X POST \"${API_BASE_URL}/api/choices/\" \\
    -H \"Content-Type: application/json\" \\
    -d \"${CHOICE_PAYLOAD}\""

  RESPONSE=$(curl -s -w "\nHTTP_CODE:%{http_code}" \
    -X POST "${API_BASE_URL}/api/choices/" \
    -H "Content-Type: application/json" \
    -d "${CHOICE_PAYLOAD}"
  )

  BODY=$(echo "$RESPONSE" | sed -e '/HTTP_CODE:/d')
  HTTP_CODE=$(echo "$RESPONSE" | sed -n 's/.*HTTP_CODE:\([0-9]*\).*/\1/p')

  if [[ "$HTTP_CODE" -ne 200 && "$HTTP_CODE" -ne 201 ]]; then
    echo "ERROR: Failed to create choice '$choice_text' for Question #2."
    exit 1
  fi

  echo "Created choice '$choice_text' (HTTP code: $HTTP_CODE)"
done

echo
echo "=========================================="
echo "All steps completed successfully."
echo "Poll, Questions, and Choices have been created and kept."
echo "=========================================="
exit 0