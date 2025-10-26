package adapters

import (
	"cmp"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hay-kot/hookfeed/backend/internal/data/dtos"
)

// ntfyMessage represents a parsed ntfy-compatible message
type ntfyMessage struct {
	Topic    string            `json:"topic"`
	Message  string            `json:"message"`
	Title    string            `json:"title"`
	Priority int32             `json:"priority"`
	Tags     []string          `json:"tags"`
	Click    string            `json:"click"`
	Icon     string            `json:"icon"`
	Actions  []json.RawMessage `json:"actions"`
	Markdown bool              `json:"markdown"`
}

// ParseNtfyMessage parses a ntfy compatible http request and transforms it into
// a validated creation object or returns an error.
// Priority order: JSON body < Query Params < Headers
func ParseNtfyMessage(r *http.Request) (dtos.FeedMessageCreate, error) {
	// Copy HTTP request data (raw body, headers, query params)
	data, err := feedMessageFromRequest(r)
	if err != nil {
		return dtos.FeedMessageCreate{}, fmt.Errorf("failed to copy request: %w", err)
	}

	// Set the feed ID from URL parameter
	data.FeedID = chi.URLParam(r, "topic")

	// Parse JSON body if Content-Type is application/json
	var jsonMsg ntfyMessage
	if r.Header.Get("Content-Type") == "application/json" {
		var bodyData map[string]interface{}
		if err := json.Unmarshal(data.RawRequest, &bodyData); err == nil {
			// Try to unmarshal as ntfy message structure
			bodyBytes, _ := json.Marshal(bodyData)
			_ = json.Unmarshal(bodyBytes, &jsonMsg)
		}
	}

	// Extract values from query parameters
	query := r.URL.Query()
	queryTitle := GetQueryParam(query, "title", "t")
	queryMessage := GetQueryParam(query, "message", "m")
	queryPriorityStr := GetQueryParam(query, "priority", "p")

	var queryPriority int32
	if queryPriorityStr != "" {
		queryPriority, _ = ParsePriority(queryPriorityStr)
	}

	// Extract values from headers
	headerTitle := GetHeader(r, "X-Title", "Title")
	headerMessage := GetHeader(r, "X-Message", "Message")
	headerPriorityStr := GetHeader(r, "X-Priority", "Priority")

	var headerPriority int32
	if headerPriorityStr != "" {
		headerPriority, _ = ParsePriority(headerPriorityStr)
	}

	// Use cmp.Or to select first non-empty value (precedence: headers > query > json)
	title := cmp.Or(headerTitle, queryTitle, jsonMsg.Title)
	message := cmp.Or(headerMessage, queryMessage, jsonMsg.Message)
	priority := cmp.Or(headerPriority, queryPriority, jsonMsg.Priority, int32(3))

	// If message is still empty, try to use the raw body as plain text
	if message == "" {
		var bodyData map[string]interface{}
		if err := json.Unmarshal(data.RawRequest, &bodyData); err == nil {
			// Check for $body key (plain text wrapped by copyBody)
			if bodyStr, ok := bodyData["$body"].(string); ok && bodyStr != "" {
				message = bodyStr
			}
		}
	}

	// Set the parsed ntfy fields
	data.Title = title
	data.Message = message
	data.Priority = priority

	return data, nil
}
