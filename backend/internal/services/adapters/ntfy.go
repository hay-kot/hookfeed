package adapters

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/hay-kot/hookfeed/backend/internal/data/dtos"
)

type Adapter interface {
	UnmarshalRequest(r *http.Request) error
	AsFeedMessage() dtos.FeedMessageCreate
}

// NtfyMessage represents a parsed ntfy-compatible message
type NtfyMessage struct {
	Topic    string            `json:"topic,omitempty"`
	Message  string            `json:"message,omitempty"`
	Title    string            `json:"title,omitempty"`
	Priority int32             `json:"priority,omitempty"`
	Tags     []string          `json:"tags,omitempty"`
	Click    string            `json:"click,omitempty"`
	Icon     string            `json:"icon,omitempty"`
	Actions  []json.RawMessage `json:"actions,omitempty"`
	Markdown bool              `json:"markdown,omitempty"`
}

// NtfyAdapter adapts ntfy-style requests to FeedMessage
type NtfyAdapter struct {
	Message        NtfyMessage
	RawBody        []byte
	RawHeaders     map[string]string
	RawQueryParams map[string][]string
}

// UnmarshalRequest parses an ntfy-compatible HTTP request
func (na *NtfyAdapter) UnmarshalRequest(r *http.Request) error {
	topic := chi.URLParam(r, "topic")
	na.Message = NtfyMessage{
		Topic:    topic,
		Priority: 3, // default priority
	}

	// Read body once
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	na.RawBody = bodyBytes
	// Restore body for potential later use
	r.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

	// Capture raw headers
	na.RawHeaders = make(map[string]string)
	for k, v := range r.Header {
		if len(v) > 0 {
			na.RawHeaders[k] = v[0]
		}
	}

	// Try to parse as JSON first
	if r.Header.Get("Content-Type") == "application/json" {
		var jsonMsg NtfyMessage
		if err := json.Unmarshal(bodyBytes, &jsonMsg); err == nil {
			// JSON parsed successfully
			if jsonMsg.Topic != "" {
				na.Message.Topic = jsonMsg.Topic
			}
			if jsonMsg.Message != "" {
				na.Message.Message = jsonMsg.Message
			}
			if jsonMsg.Title != "" {
				na.Message.Title = jsonMsg.Title
			}
			if jsonMsg.Priority > 0 {
				na.Message.Priority = jsonMsg.Priority
			}
			if len(jsonMsg.Tags) > 0 {
				na.Message.Tags = jsonMsg.Tags
			}
			na.Message.Click = jsonMsg.Click
			na.Message.Icon = jsonMsg.Icon
			na.Message.Actions = jsonMsg.Actions
			na.Message.Markdown = jsonMsg.Markdown
		}
	}

	// Parse query parameters (override JSON if present)
	query := r.URL.Query()
	na.RawQueryParams = query

	if title := GetQueryParam(query, "title", "t"); title != "" {
		na.Message.Title = title
	}

	if msgText := GetQueryParam(query, "message", "m"); msgText != "" {
		na.Message.Message = msgText
	}

	if priorityStr := GetQueryParam(query, "priority", "p"); priorityStr != "" {
		if p, err := ParsePriority(priorityStr); err == nil {
			na.Message.Priority = p
		}
	}

	if tagsStr := GetQueryParam(query, "tags", "ta"); tagsStr != "" {
		na.Message.Tags = SplitAndTrim(tagsStr)
	}

	if click := GetQueryParam(query, "click"); click != "" {
		na.Message.Click = click
	}

	if icon := GetQueryParam(query, "icon"); icon != "" {
		na.Message.Icon = icon
	}

	if markdown := GetQueryParam(query, "markdown", "md"); markdown != "" {
		na.Message.Markdown = ParseBool(markdown)
	}

	// Parse headers (headers override query parameters if present)
	if title := GetHeader(r, "X-Title", "Title"); title != "" {
		na.Message.Title = title
	}

	if msgText := GetHeader(r, "X-Message", "Message"); msgText != "" {
		na.Message.Message = msgText
	}

	if priorityStr := GetHeader(r, "X-Priority", "Priority"); priorityStr != "" {
		if p, err := ParsePriority(priorityStr); err == nil {
			na.Message.Priority = p
		}
	}

	if tagsStr := GetHeader(r, "X-Tags", "Tags"); tagsStr != "" {
		na.Message.Tags = SplitAndTrim(tagsStr)
	}

	if click := GetHeader(r, "X-Click", "Click"); click != "" {
		na.Message.Click = click
	}

	if icon := GetHeader(r, "X-Icon", "Icon"); icon != "" {
		na.Message.Icon = icon
	}

	if markdown := GetHeader(r, "X-Markdown", "Markdown"); markdown != "" {
		na.Message.Markdown = ParseBool(markdown)
	}

	// If message is still empty, use body as plain text
	if na.Message.Message == "" {
		na.Message.Message = string(bodyBytes)
	}

	return nil
}

// AsFeedMessage converts the ntfy message to a FeedMessageCreate DTO
func (na *NtfyAdapter) AsFeedMessage() dtos.FeedMessageCreate {
	// Build the body structure for RawRequest
	var bodyData any
	if json.Valid(na.RawBody) {
		// Body is already valid JSON, unmarshal it
		json.Unmarshal(na.RawBody, &bodyData)
	} else {
		// Body is plain text, wrap it
		bodyData = map[string]string{"body": string(na.RawBody)}
	}

	// Convert map[string]string headers to http.Header format
	headers := make(map[string][]string)
	for k, v := range na.RawHeaders {
		headers[k] = []string{v}
	}

	// Use the constructor to ensure proper initialization of raw fields
	createDTO, err := dtos.NewFeedMessageCreateFromHTTP(
		na.Message.Topic,
		bodyData,
		headers,
		na.RawQueryParams,
	)
	if err != nil {
		// Fallback to empty DTO if constructor fails (shouldn't happen)
		createDTO = dtos.FeedMessageCreate{
			FeedID:         na.Message.Topic,
			RawRequest:     json.RawMessage("{}"),
			RawHeaders:     json.RawMessage("{}"),
			RawQueryParams: json.RawMessage("{}"),
			Logs:           []string{},
			Metadata:       json.RawMessage("{}"),
		}
	}

	// Build metadata from ntfy-specific fields
	metadata := make(map[string]interface{})
	if len(na.Message.Tags) > 0 {
		metadata["tags"] = na.Message.Tags
	}
	if na.Message.Click != "" {
		metadata["click"] = na.Message.Click
	}
	if na.Message.Icon != "" {
		metadata["icon"] = na.Message.Icon
	}
	if len(na.Message.Actions) > 0 {
		metadata["actions"] = na.Message.Actions
	}
	if na.Message.Markdown {
		metadata["markdown"] = true
	}
	metadataJSON, _ := json.Marshal(metadata)

	// Override with ntfy-specific parsed fields
	title := na.Message.Title
	message := na.Message.Message
	priority := na.Message.Priority

	createDTO.Title = &title
	createDTO.Message = &message
	createDTO.Priority = &priority
	createDTO.Metadata = json.RawMessage(metadataJSON)

	return createDTO
}
