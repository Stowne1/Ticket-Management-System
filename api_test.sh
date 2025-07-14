#!/bin/bash

API_URL="http://localhost:8080/tickets"

# 1. Create a ticket
printf "\n=== Create Ticket ===\n"
CREATE_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}\n" -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d '{"title":"Test Ticket","description":"Test Desc","status":"open"}')
echo "$CREATE_RESPONSE"

# Extract ID (assume ID is 1 for this test, or adjust as needed)
TICKET_ID=1

# 2. Get the ticket by ID
printf "\n=== Get Ticket ===\n"
curl -s -w "\nHTTP_STATUS:%{http_code}\n" "$API_URL/$TICKET_ID"

# 3. Update the ticket
printf "\n=== Update Ticket ===\n"
curl -s -w "\nHTTP_STATUS:%{http_code}\n" -X PUT "$API_URL/$TICKET_ID" \
  -H "Content-Type: application/json" \
  -d '{"title":"Test Ticket","description":"Updated Desc","status":"closed"}'

# 4. Delete the ticket
printf "\n=== Delete Ticket ===\n"
curl -s -w "\nHTTP_STATUS:%{http_code}\n" -X DELETE "$API_URL/$TICKET_ID"

# 5. Error: Create with missing fields
printf "\n=== Create Ticket (Missing Fields) ===\n"
curl -s -w "\nHTTP_STATUS:%{http_code}\n" -X POST "$API_URL" \
  -H "Content-Type: application/json" \
  -d '{"title":"","description":"","status":""}'

# 6. Error: Get non-existent ticket
printf "\n=== Get Non-existent Ticket ===\n"
curl -s -w "\nHTTP_STATUS:%{http_code}\n" "$API_URL/9999"

# 7. Error: Update with invalid ID
printf "\n=== Update Invalid Ticket ===\n"
curl -s -w "\nHTTP_STATUS:%{http_code}\n" -X PUT "$API_URL/9999" \
  -H "Content-Type: application/json" \
  -d '{"title":"Test Ticket","description":"Updated Desc","status":"closed"}'

# 8. Error: Delete with invalid ID
printf "\n=== Delete Invalid Ticket ===\n"
curl -s -w "\nHTTP_STATUS:%{http_code}\n" -X DELETE "$API_URL/9999" 