package dtos

import (
	"time"

	"github.com/google/uuid"
	"github.com/hay-kot/hookfeed/backend/internal/data/db"
)

type FeedMessage struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type FeedMessageCreate struct{}

type FeedMessageUpdate struct{}

func MapFeedMessage(d db.FeedMessage) FeedMessage {
	//exhaustruct:enforce
	return FeedMessage{
		ID:        d.ID,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}
