#!/bin/bash

# Ntfy message testing script for HookFeed
# Sends various test messages using ntfy CLI to demonstrate priority levels

set -e

# Configuration
API_URL="${HOOKFEED_API_URL:-http://localhost:9990}"
TOPIC="${1:-ntfy-test}"  # Default to ntfy-test feed

# Color codes for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}  HookFeed Ntfy Message Testing${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${YELLOW}API URL:${NC} $API_URL"
echo -e "${YELLOW}Topic:${NC} $TOPIC"
echo ""

# Check if ntfy is installed
if ! command -v ntfy &> /dev/null; then
    echo -e "${RED}✗ ntfy CLI not found. Please install it first:${NC}"
    echo "  brew install ntfy (macOS)"
    echo "  or visit: https://ntfy.sh/"
    exit 1
fi

echo -e "${YELLOW}Sending test messages with different priorities...${NC}"
echo ""

# Extract host from API_URL (remove http:// or https://)
SERVER_HOST=$(echo "$API_URL" | sed -e 's|^https\?://||')

# Message 1: Min priority (debug-level)
echo -e "${BLUE}[1/5]${NC} Sending min priority message..."
echo -e "${YELLOW}$ ntfy publish $SERVER_HOST/$TOPIC -p min -t 'Debug Info' --tags='bug,wrench' 'This is a debug message...'${NC}"
ntfy publish \
    "$SERVER_HOST/$TOPIC" \
    "This is a debug message with minimum priority" \
    -p min \
    -t "Debug Info" \
    --tags="bug,wrench" | jq
echo ""

# Message 2: Low priority
echo -e "${BLUE}[2/5]${NC} Sending low priority message..."
echo -e "${YELLOW}$ ntfy publish $SERVER_HOST/$TOPIC -p low -t 'Info Update' --tags='information_source' 'Routine information update'${NC}"
ntfy publish \
    "$SERVER_HOST/$TOPIC" \
    "Routine information update" \
    -p low \
    -t "Info Update" \
    --tags="information_source" | jq
echo ""

# Message 3: Default priority
echo -e "${BLUE}[3/5]${NC} Sending default priority message..."
echo -e "${YELLOW}$ ntfy publish $SERVER_HOST/$TOPIC -p default -t 'Standard Notification' --tags='bell' 'This is a standard notification...'${NC}"
ntfy publish \
    "$SERVER_HOST/$TOPIC" \
    "This is a standard notification with default priority" \
    -p default \
    -t "Standard Notification" \
    --tags="bell" | jq
echo ""

# Message 4: High priority (warning)
echo -e "${BLUE}[4/5]${NC} Sending high priority message..."
echo -e "${YELLOW}$ ntfy publish $SERVER_HOST/$TOPIC -p high -t 'Warning Alert' --tags='warning,exclamation' 'High priority warning...'${NC}"
ntfy publish \
    "$SERVER_HOST/$TOPIC" \
    "High priority warning - something needs attention" \
    -p high \
    -t "Warning Alert" \
    --tags="warning,exclamation" | jq
echo ""

# Message 5: Max priority (urgent)
echo -e "${BLUE}[5/5]${NC} Sending max priority message..."
echo -e "${YELLOW}$ ntfy publish $SERVER_HOST/$TOPIC -p max -t '🚨 Critical Alert' --tags='rotating_light,fire' 'URGENT: Critical system alert...'${NC}"
ntfy publish \
    "$SERVER_HOST/$TOPIC" \
    "URGENT: Critical system alert requiring immediate action!" \
    -p max \
    -t "🚨 Critical Alert" \
    --tags="rotating_light,fire" | jq
echo ""

echo ""
echo -e "${GREEN}✓ All test messages sent successfully!${NC}"
echo ""

# Bonus: Send a message with JSON body via curl
echo -e "${YELLOW}Sending JSON message via curl...${NC}"
curl -s -X POST "$API_URL/$TOPIC" \
    -H "Content-Type: application/json" \
    -d '{
        "message": "This message was sent as JSON",
        "title": "JSON Test",
        "priority": 4,
        "tags": ["package", "sparkles"]
    }' | jq . 2>/dev/null || echo "Message sent"

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}Testing complete!${NC}"
echo ""
echo -e "${YELLOW}Priority mapping:${NC}"
echo "  min (1)     → Debug messages"
echo "  low (2)     → Low priority info"
echo "  default (3) → Normal notifications"
echo "  high (4)    → Warnings"
echo "  max (5)     → Critical alerts"
echo ""
