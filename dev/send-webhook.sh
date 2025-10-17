#!/bin/bash

# Webhook testing script for HookFeed
# Usage: ./dev/send-webhook.sh [feed-key] [message]

set -e

# Configuration
API_URL="${HOOKFEED_API_URL:-http://localhost:9990}"
DEFAULT_FEED_KEY="3ftyjSfiZott983g1rWykD"  # dev-test feed from dev.feeds.yml

# Parse arguments
FEED_KEY="${1:-$DEFAULT_FEED_KEY}"
MESSAGE="${2:-Test webhook message}"

# Color codes for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}  HookFeed Webhook Testing Script${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${YELLOW}API URL:${NC} $API_URL"
echo -e "${YELLOW}Feed Key:${NC} $FEED_KEY"
echo -e "${YELLOW}Message:${NC} $MESSAGE"
echo ""

# Create payload
PAYLOAD=$(cat <<EOF
{
  "message": "$MESSAGE",
  "level": "info",
  "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "source": "test-script",
  "metadata": {
    "hostname": "$(hostname)",
    "user": "$(whoami)"
  }
}
EOF
)

echo -e "${YELLOW}Payload:${NC}"
echo "$PAYLOAD" | jq . 2>/dev/null || echo "$PAYLOAD"
echo ""

# Send webhook
echo -e "${YELLOW}Sending webhook...${NC}"
RESPONSE=$(curl -s -w "\n%{http_code}" -X POST \
  "$API_URL/hooks/$FEED_KEY" \
  -H "Content-Type: application/json" \
  -d "$PAYLOAD")

# Extract HTTP status code (last line) and response body (everything else)
HTTP_CODE=$(echo "$RESPONSE" | tail -n 1)
BODY=$(echo "$RESPONSE" | sed '$d')

echo ""
echo -e "${YELLOW}Response (HTTP $HTTP_CODE):${NC}"
echo "$BODY" | jq . 2>/dev/null || echo "$BODY"
echo ""

# Check if successful
if [ "$HTTP_CODE" -eq 202 ]; then
  echo -e "${GREEN}✓ Webhook successfully delivered!${NC}"

  # Extract message ID and feed ID from response
  MESSAGE_ID=$(echo "$BODY" | jq -r '.messageId' 2>/dev/null || echo "")
  FEED_ID=$(echo "$BODY" | jq -r '.feedId' 2>/dev/null || echo "")

  if [ -n "$MESSAGE_ID" ] && [ "$MESSAGE_ID" != "null" ]; then
    echo -e "${GREEN}  Message ID: $MESSAGE_ID${NC}"
  fi

  if [ -n "$FEED_ID" ] && [ "$FEED_ID" != "null" ]; then
    echo -e "${GREEN}  Feed ID: $FEED_ID${NC}"
  fi
else
  echo -e "\033[0;31m✗ Webhook delivery failed with HTTP $HTTP_CODE${NC}"
  exit 1
fi

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
