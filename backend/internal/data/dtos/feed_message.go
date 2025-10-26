package dtos

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/hay-kot/hookfeed/backend/internal/data/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type MessageState string

const (
	MessageStateNew          MessageState = "new"
	MessageStateAcknowledged MessageState = "acknowledged"
	MessageStateResolved     MessageState = "resolved"
	MessageStateArchived     MessageState = "archived"
)

type FeedMessage struct {
	ID             uuid.UUID       `json:"id"`
	FeedSlug       string          `json:"feedSlug"`
	RawRequest     json.RawMessage `json:"rawRequest,omitempty"`
	RawHeaders     json.RawMessage `json:"rawHeaders,omitempty"`
	RawQueryParams json.RawMessage `json:"rawQueryParams,omitempty"`
	Title          *string         `json:"title"`
	Message        *string         `json:"message"`
	Priority       int32           `json:"priority"`
	Logs           []string        `json:"logs"`
	Metadata       json.RawMessage `json:"metadata"`
	State          string          `json:"state"`
	StateChangedAt *time.Time      `json:"stateChangedAt,omitempty"`
	ReceivedAt     time.Time       `json:"receivedAt"`
	ProcessedAt    *time.Time      `json:"processedAt,omitempty"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
	SearchVector   *string         `json:"-"`
}

type FeedMessageCreate struct {
	FeedID         string          `json:"feedSlug"       validate:"required"`
	RawRequest     json.RawMessage `json:"rawRequest"     validate:"required"`
	RawHeaders     json.RawMessage `json:"rawHeaders"     validate:"required"`
	RawQueryParams json.RawMessage `json:"rawQueryParams" validate:"required"`
	Metadata       json.RawMessage `json:"metadata"`
	Title          string          `json:"title"`
	Message        string          `json:"message"`
	Priority       int32           `json:"priority"       validate:"omitempty,min=1,max=5"`
	Logs           []string        `json:"logs"`
	State          string          `json:"state"          validate:"omitempty,oneof=new acknowledged resolved archived"`
	ReceivedAt     time.Time       `json:"receivedAt"`
	ProcessedAt    *time.Time      `json:"processedAt"`
}

// FeedMessageCreateNew creates a base example of the FeedMessageCreate. We do this to ensure
// that some default fields are set like the RawRequest properties. If we don't set these
// We'll get an error in the database
func FeedMessageCreateNew() FeedMessageCreate {
	return FeedMessageCreate{
		RawRequest:     json.RawMessage("{}"),
		RawHeaders:     json.RawMessage("{}"),
		RawQueryParams: json.RawMessage("{}"),
		Metadata:       json.RawMessage("{}"),
		Title:          "Untitled Message",
		Message:        "no message provided",
		Priority:       1,
		Logs:           []string{},
		State:          "new",
		ReceivedAt:     time.Now(),
		ProcessedAt:    &time.Time{},
	}
}

// NewFeedMessageCreateFromHTTP creates a FeedMessageCreate from HTTP request components.
// This ensures proper initialization of all required fields with consistent defaults.
func NewFeedMessageCreateFromHTTP(
	feedID string,
	body any,
	headers http.Header,
	queryParams url.Values,
) (FeedMessageCreate, error) {
	// Marshal body
	rawRequest, err := json.Marshal(body)
	if err != nil {
		return FeedMessageCreate{}, fmt.Errorf("marshal body: %w", err)
	}

	// Marshal headers (ensure not nil)
	if headers == nil {
		headers = make(http.Header)
	}
	rawHeaders, err := json.Marshal(headers)
	if err != nil {
		return FeedMessageCreate{}, fmt.Errorf("marshal headers: %w", err)
	}

	// Marshal query params (ensure not nil)
	if queryParams == nil {
		queryParams = make(url.Values)
	}
	rawQueryParams, err := json.Marshal(queryParams)
	if err != nil {
		return FeedMessageCreate{}, fmt.Errorf("marshal query params: %w", err)
	}

	return FeedMessageCreate{
		FeedID:         feedID,
		RawRequest:     rawRequest,
		RawHeaders:     rawHeaders,
		RawQueryParams: rawQueryParams,
		Logs:           []string{},
		Metadata:       json.RawMessage("{}"),
	}, nil
}

type FeedMessageUpdate struct {
	Title    *string         `json:"title"`
	Message  *string         `json:"message"`
	Priority *int32          `json:"priority" validate:"omitempty,min=1,max=5"`
	State    *string         `json:"state"    validate:"omitempty,oneof=new acknowledged resolved archived"`
	Logs     []string        `json:"logs"`
	Metadata json.RawMessage `json:"metadata"`
}

type FeedMessageUpdateState struct {
	State string `json:"state" validate:"required,oneof=new acknowledged resolved archived"`
}

type FeedMessageBulkUpdateState struct {
	MessageIDs []uuid.UUID `json:"messageIds" validate:"required,min=1"`
	State      string      `json:"state"      validate:"required,oneof=new acknowledged resolved archived"`
}

type FeedMessageBulkDelete struct {
	MessageIDs []uuid.UUID              `json:"messageIds,omitempty"`
	Filter     *FeedMessageDeleteFilter `json:"filter,omitempty"`
}

type FeedMessageDeleteFilter struct {
	Priority  *int32     `json:"priority"  validate:"omitempty,min=1,max=5"`
	OlderThan *time.Time `json:"olderThan"`
}

type FeedMessageQuery struct {
	Pagination
	FeedSlug *string    `json:"feedSlug" query:"feedSlug"`
	Priority *int32     `json:"priority" validate:"omitempty,min=1,max=5"                              query:"priority"`
	State    *string    `json:"state"    validate:"omitempty,oneof=new acknowledged resolved archived" query:"state"`
	Since    *time.Time `json:"since"    query:"since"`
	Until    *time.Time `json:"until"    query:"until"`
	Query    *string    `json:"q"        query:"q"`
}

func MapFeedMessage(d db.FeedMessage) FeedMessage {
	priority := int32(3)
	if d.Priority != nil {
		priority = *d.Priority
	}

	state := "new"
	if d.State != nil {
		state = *d.State
	}

	//exhaustruct:enforce
	return FeedMessage{
		ID:             d.ID,
		FeedSlug:       d.FeedSlug,
		RawRequest:     json.RawMessage(d.RawRequest),
		RawHeaders:     json.RawMessage(d.RawHeaders),
		RawQueryParams: json.RawMessage(d.RawQueryParams),
		Title:          d.Title,
		Message:        d.Message,
		Priority:       priority,
		Logs:           d.Logs,
		Metadata:       json.RawMessage(d.Metadata),
		State:          state,
		StateChangedAt: pgTimestampToTimePtr(d.StateChangedAt),
		ReceivedAt:     d.ReceivedAt,
		ProcessedAt:    pgTimestampToTimePtr(d.ProcessedAt),
		CreatedAt:      d.CreatedAt,
		UpdatedAt:      d.UpdatedAt,
		SearchVector:   nil,
	}
}

func MapFeedMessageView(d db.FeedMessagesView) FeedMessage {
	priority := int32(3)
	if d.Priority != nil {
		priority = *d.Priority
	}

	state := "new"
	if d.State != nil {
		state = *d.State
	}

	//exhaustruct:enforce
	return FeedMessage{
		ID:             d.ID,
		FeedSlug:       d.FeedSlug,
		RawRequest:     json.RawMessage(d.RawRequest),
		RawHeaders:     json.RawMessage(d.RawHeaders),
		RawQueryParams: json.RawMessage(d.RawQueryParams),
		Title:          d.Title,
		Message:        d.Message,
		Priority:       priority,
		Logs:           d.Logs,
		Metadata:       json.RawMessage(d.Metadata),
		State:          state,
		StateChangedAt: pgTimestampToTimePtr(d.StateChangedAt),
		ReceivedAt:     d.ReceivedAt,
		ProcessedAt:    pgTimestampToTimePtr(d.ProcessedAt),
		CreatedAt:      d.CreatedAt,
		UpdatedAt:      d.UpdatedAt,
		SearchVector:   nil,
	}
}

// Helper functions for type conversions
func pgTimestampToTimePtr(ts pgtype.Timestamp) *time.Time {
	if !ts.Valid {
		return nil
	}
	return &ts.Time
}
