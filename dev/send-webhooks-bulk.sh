#!/bin/bash

# Bulk webhook testing script for HookFeed
# Sends multiple webhooks with different levels

set -e

# Configuration
API_URL="${HOOKFEED_API_URL:-http://localhost:9990}"
DEFAULT_FEED_KEY="3ftyjSfiZott983g1rWykD"  # dev-test feed from dev.feeds.yml

# Parse arguments
FEED_KEY="${1:-$DEFAULT_FEED_KEY}"
COUNT="${2:-5}"

# Color codes
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}  Sending $COUNT bulk webhooks to HookFeed${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

LEVELS=("info" "warning" "error" "success" "debug")
MESSAGES=(
  "Application started successfully"
  "High memory usage detected"
  "Failed to connect to database"
  "Deployment completed"
  "Debug: Processing batch job"
)

for i in $(seq 1 "$COUNT"); do
  # Rotate through levels
  LEVEL_INDEX=$(( (i - 1) % 5 ))
  LEVEL="${LEVELS[$LEVEL_INDEX]}"
  MESSAGE="${MESSAGES[$LEVEL_INDEX]} #$i"

  echo "[$i/$COUNT] Sending $LEVEL webhook..."

  PAYLOAD=$(cat <<EOF
{
  "message": "$MESSAGE",
  "level": "$LEVEL",
  "timestamp": "$(date -u +"%Y-%m-%dT%H:%M:%SZ")",
  "source": "bulk-test-script",
  "iteration": $i,
  "metadata": {
    "hostname": "$(hostname)",
    "batch_id": "$(uuidgen | tr '[:upper:]' '[:lower:]')"
  }
}
EOF
)

  RESPONSE=$(curl -s -w "\n%{http_code}" -X POST \
    "$API_URL/hooks/$FEED_KEY" \
    -H "Content-Type: application/json" \
    -d "$PAYLOAD")

  HTTP_CODE=$(echo "$RESPONSE" | tail -n 1)

  if [ "$HTTP_CODE" -eq 202 ]; then
    echo "  ✓ Success"
  else
    echo "  ✗ Failed (HTTP $HTTP_CODE)"
  fi

  # Small delay to avoid overwhelming the server
  sleep 0.1
done

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo "✓ Sent $COUNT webhooks"
echo ""
