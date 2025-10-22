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
echo -e "${BLUE}  HookFeed Webhook Endpoint Examples${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${YELLOW}Testing webhook endpoints with different feeds...${NC}"
echo ""

# Production Alerts - Critical system alert
echo -e "${BLUE}[1/6]${NC} Sending to Production Alerts feed..."
echo -e "${YELLOW}$ curl -X POST $API_URL/hooks/1ftyjSfiZott986g1rWykD${NC}"
curl -s -X POST "$API_URL/hooks/1ftyjSfiZott986g1rWykD" \
    -H "Content-Type: application/json" \
    -d '{
        "alert": "high_cpu_usage",
        "severity": "critical",
        "host": "prod-server-01",
        "message": "CPU usage at 95% for 5 minutes",
        "metrics": {
            "cpu_percent": 95.3,
            "memory_percent": 78.2,
            "disk_usage": 82.1
        },
        "timestamp": "2024-01-15T10:30:00Z"
    }' | jq . 2>/dev/null || echo "✓ Sent"
echo ""

# Production Alerts - Database connection issue
echo -e "${BLUE}[2/6]${NC} Sending database alert to Production Alerts..."
echo -e "${YELLOW}$ curl -X POST $API_URL/hooks/KcvgA8pwwmY3XZgEAO4KX3${NC}"
curl -s -X POST "$API_URL/hooks/KcvgA8pwwmY3XZgEAO4KX3" \
    -H "Content-Type: application/json" \
    -d '{
        "alert": "database_connection_pool_exhausted",
        "severity": "critical",
        "service": "api-gateway",
        "message": "Database connection pool exhausted - 0/100 connections available",
        "impact": "Users experiencing 500 errors",
        "runbook": "https://wiki.company.com/runbooks/db-connections"
    }' | jq . 2>/dev/null || echo "✓ Sent"
echo ""

# GitHub Events - Simulated GitHub webhook
echo -e "${BLUE}[3/6]${NC} Sending GitHub push event..."
echo -e "${YELLOW}$ curl -X POST $API_URL/hooks/1ftyjSfiZott983g1rWykD${NC}"
curl -s -X POST "$API_URL/hooks/1ftyjSfiZott983g1rWykD" \
    -H "Content-Type: application/json" \
    -H "X-GitHub-Event: push" \
    -d '{
        "ref": "refs/heads/main",
        "repository": {
            "name": "hookfeed",
            "full_name": "hay-kot/hookfeed",
            "html_url": "https://github.com/hay-kot/hookfeed"
        },
        "pusher": {
            "name": "hay-kot",
            "email": "hay-kot@users.noreply.github.com"
        },
        "commits": [
            {
                "id": "abc123def456",
                "message": "feat: add webhook processing middleware",
                "author": {
                    "name": "Hayden K",
                    "email": "hay-kot@users.noreply.github.com"
                },
                "url": "https://github.com/hay-kot/hookfeed/commit/abc123def456"
            }
        ]
    }' | jq . 2>/dev/null || echo "✓ Sent"
echo ""

# GitHub Events - Pull request opened
echo -e "${BLUE}[4/6]${NC} Sending GitHub pull request event..."
echo -e "${YELLOW}$ curl -X POST $API_URL/hooks/KcvgA8pwwmY3XZ1EAO4KX3${NC}"
curl -s -X POST "$API_URL/hooks/KcvgA8pwwmY3XZ1EAO4KX3" \
    -H "Content-Type: application/json" \
    -H "X-GitHub-Event: pull_request" \
    -d '{
        "action": "opened",
        "pull_request": {
            "number": 42,
            "title": "Add ntfy adapter support",
            "user": {
                "login": "hay-kot"
            },
            "html_url": "https://github.com/hay-kot/hookfeed/pull/42",
            "body": "This PR adds support for the ntfy notification protocol",
            "labels": ["enhancement", "feature"]
        },
        "repository": {
            "full_name": "hay-kot/hookfeed"
        }
    }' | jq . 2>/dev/null || echo "✓ Sent"
echo ""

# Development Testing - Custom webhook with nested data
echo -e "${BLUE}[5/6]${NC} Sending custom webhook to Development Testing feed..."
echo -e "${YELLOW}$ curl -X POST $API_URL/hooks/3ftyjSfiZott983g1rWykD${NC}"
curl -s -X POST "$API_URL/hooks/3ftyjSfiZott983g1rWykD" \
    -H "Content-Type: application/json" \
    -d '{
        "event_type": "deployment",
        "environment": "staging",
        "service": "api-v2",
        "version": "1.2.3",
        "status": "success",
        "deployment": {
            "started_at": "2024-01-15T10:30:00Z",
            "completed_at": "2024-01-15T10:35:00Z",
            "duration_seconds": 300,
            "deployed_by": "github-actions[bot]"
        },
        "health_checks": {
            "http_status": 200,
            "database": "connected",
            "redis": "connected"
        }
    }' | jq . 2>/dev/null || echo "✓ Sent"
echo ""

# Development Testing - Error report
echo -e "${BLUE}[6/6]${NC} Sending error report to Development Testing..."
echo -e "${YELLOW}$ curl -X POST $API_URL/hooks/fcvgA8pwwmY3XZ1EAO4KX3${NC}"
curl -s -X POST "$API_URL/hooks/fcvgA8pwwmY3XZ1EAO4KX3" \
    -H "Content-Type: application/json" \
    -d '{
        "error": {
            "type": "UnhandledException",
            "message": "Failed to process webhook: connection timeout",
            "stack_trace": "Error: connection timeout\n  at processWebhook (webhook.go:42)\n  at handleRequest (server.go:123)",
            "request_id": "req_abc123xyz",
            "user_id": "user_789"
        },
        "context": {
            "url": "/api/hooks/test",
            "method": "POST",
            "ip_address": "192.168.1.100"
        },
        "severity": "error",
        "timestamp": "2024-01-15T10:30:00Z"
    }' | jq . 2>/dev/null || echo "✓ Sent"
echo ""

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
