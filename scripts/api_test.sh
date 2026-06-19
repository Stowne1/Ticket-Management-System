#!/bin/bash
# API Smoke Test Script for Ticket Management System
# Requires a running server on localhost:8080.

BASE="http://localhost:8080/api"
TICKETS="$BASE/tickets"
PASS=0
FAIL=0

check() {
  local label="$1" got="$2" want="$3"
  if [ "$got" = "$want" ]; then
    echo "  PASS: $label (HTTP $got)"
    PASS=$((PASS + 1))
  else
    echo "  FAIL: $label — expected HTTP $want, got HTTP $got"
    FAIL=$((FAIL + 1))
  fi
}

# --- Auth setup ---
# Register a dedicated smoke-test user (may already exist on re-runs; that's fine).
printf "\n=== Register test user ===\n"
curl -s -o /dev/null -X POST "$BASE/register" \
  -H "Content-Type: application/json" \
  -d '{"email":"smoketest@example.com","password":"smokepass123"}'

printf "\n=== Login ===\n"
LOGIN=$(curl -s -X POST "$BASE/login" \
  -H "Content-Type: application/json" \
  -d '{"email":"smoketest@example.com","password":"smokepass123"}')
echo "$LOGIN"

# Extract token without requiring jq
TOKEN=$(echo "$LOGIN" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
if [ -z "$TOKEN" ]; then
  echo "ERROR: could not get auth token — aborting smoke test"
  exit 1
fi
AUTH="Authorization: Bearer $TOKEN"

# --- Ticket CRUD ---
printf "\n=== Create Ticket ===\n"
CREATE=$(curl -s -w "\nHTTP_STATUS:%{http_code}\n" -X POST "$TICKETS" \
  -H "Content-Type: application/json" -H "$AUTH" \
  -d '{"title":"Smoke Test Ticket","description":"Created by smoke test","status":"open"}')
echo "$CREATE"
check "Create ticket" "$(echo "$CREATE" | grep -o 'HTTP_STATUS:[0-9]*' | cut -d: -f2)" "201"

TICKET_ID=$(echo "$CREATE" | grep -o '"id":[0-9]*' | head -1 | cut -d: -f2)
if [ -z "$TICKET_ID" ]; then TICKET_ID=1; fi

printf "\n=== Get Ticket ===\n"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" -H "$AUTH" "$TICKETS/$TICKET_ID")
check "Get ticket" "$STATUS" "200"

printf "\n=== Update Ticket ===\n"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X PUT "$TICKETS/$TICKET_ID" \
  -H "Content-Type: application/json" -H "$AUTH" \
  -d '{"title":"Updated Title","description":"Updated desc","status":"closed"}')
check "Update ticket" "$STATUS" "200"

printf "\n=== Delete Ticket ===\n"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X DELETE -H "$AUTH" "$TICKETS/$TICKET_ID")
check "Delete ticket" "$STATUS" "200"

# --- Error cases ---
printf "\n=== Create ticket with missing fields ===\n"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$TICKETS" \
  -H "Content-Type: application/json" -H "$AUTH" \
  -d '{"title":"","description":"","status":""}')
check "Missing fields → 400" "$STATUS" "400"

printf "\n=== Get non-existent ticket ===\n"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" -H "$AUTH" "$TICKETS/999999")
check "Not found → 404" "$STATUS" "404"

printf "\n=== Request without token ===\n"
STATUS=$(curl -s -o /dev/null -w "%{http_code}" "$TICKETS")
check "No token → 401" "$STATUS" "401"

printf "\n=== Results: $PASS passed, $FAIL failed ===\n"
[ "$FAIL" -eq 0 ]
