package dtos

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/hay-kot/hookfeed/backend/internal/data/db"
	"github.com/jackc/pgx/v5/pgtype"
)

type MessageLevel string

const (
	MessageLevelInfo    MessageLevel = "info"
	MessageLevelWarning MessageLevel = "warning"
	MessageLevelError   MessageLevel = "error"
	MessageLevelSuccess MessageLevel = "success"
	MessageLevelDebug   MessageLevel = "debug"
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
	Title          *string         `json:"title"`
	Message        *string         `json:"message"`
	Level          string          `json:"level"`
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
	FeedID      string          `json:"feedSlug"    validate:"required"`
	RawRequest  json.RawMessage `json:"rawRequest"  validate:"required"`
	RawHeaders  json.RawMessage `json:"rawHeaders"  validate:"required"`
	Title       *string         `json:"title"`
	Message     *string         `json:"message"`
	Level       *string         `json:"level"       validate:"omitempty,oneof=info warning error success debug"`
	Logs        []string        `json:"logs"`
	Metadata    json.RawMessage `json:"metadata"`
	State       *string         `json:"state"       validate:"omitempty,oneof=new acknowledged resolved archived"`
	ReceivedAt  *time.Time      `json:"receivedAt"`
	ProcessedAt *time.Time      `json:"processedAt"`
}

type FeedMessageUpdate struct {
	Title    *string         `json:"title"`
	Message  *string         `json:"message"`
	Level    *string         `json:"level"    validate:"omitempty,oneof=info warning error success debug"`
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
	Level     *string    `json:"level"     validate:"omitempty,oneof=info warning error success debug"`
	OlderThan *time.Time `json:"olderThan"`
}

type FeedMessageQuery struct {
	Pagination
	FeedSlug *string    `json:"feedSlug" query:"feedSlug"`
	Level    *string    `json:"level"    validate:"omitempty,oneof=info warning error success debug"   query:"level"`
	State    *string    `json:"state"    validate:"omitempty,oneof=new acknowledged resolved archived" query:"state"`
	Since    *time.Time `json:"since"    query:"since"`
	Until    *time.Time `json:"until"    query:"until"`
	Query    *string    `json:"q"        query:"q"`
}

func MapFeedMessage(d db.FeedMessage) FeedMessage {
	level := "info"
	if d.Level != nil {
		level = *d.Level
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
		Title:          d.Title,
		Message:        d.Message,
		Level:          level,
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
	level := "info"
	if d.Level != nil {
		level = *d.Level
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
		Title:          d.Title,
		Message:        d.Message,
		Level:          level,
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

func timePtrToPgTimestamp(t *time.Time) pgtype.Timestamp {
	if t == nil {
		return pgtype.Timestamp{}
	}
	return pgtype.Timestamp{
		Time:  *t,
		Valid: true,
	}
}
