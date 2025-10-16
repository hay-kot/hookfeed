package dtos

import (
	"github.com/google/uuid"
)

type WebhookRequest struct {
	FeedKey string              // Feed key from URL path
	Headers map[string][]string // All request headers
	Body    map[string]any      // Raw JSON body
}

// WebhookResponse represents the response sent back to the webhook sender
type WebhookResponse struct {
	Success   bool      `json:"success"`
	MessageID uuid.UUID `json:"messageId"`
	FeedID    string    `json:"feedId"`
}
