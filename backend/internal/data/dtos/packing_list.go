package dtos

import (
	"time"

	"cloud.google.com/go/civil"
	"github.com/google/uuid"
	"github.com/hay-kot/hookfeed/backend/internal/data/db"
)

type PackingListQuery struct {
	Pagination
	OrderBy string `json:"orderBy" validate:"omitempty,oneof=created_at" query:"orderBy"`
}

type PackingListBase struct {
	ID          uuid.UUID  `json:"id"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	UserID      uuid.UUID  `json:"userId"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	DueDate     civil.Date `json:"dueDate,omitzero" swaggertype:"string" x-nullable:"true"`
	Days        int        `json:"days,omitzero"                         x-nullable:"true"`
	Tags        []string   `json:"tags"`
}

type PackingListSummary struct {
	PackingListBase
	ItemCount       int `json:"itemCount"`
	ItemPackedCount int `json:"itemPackedCount"`
}

type PackingList struct {
	PackingListBase
	Items []PackingListItem `json:"items"`
}

type PackingListCreate struct {
	Name        string      `json:"name"        validate:"required"`
	Description string      `json:"description"`
	DueDate     *civil.Date `json:"dueDate"     extensions:"x-nullable" swaggertype:"string"`
	Days        int         `json:"days"        extensions:"x-nullable"`
	Tags        []string    `json:"tags"        extensions:"x-nullable"`
}

type PackingListUpdate struct {
	Name        *string     `json:"name"        extensions:"x-nullable"`
	Description *string     `json:"description" extensions:"x-nullable"`
	Status      *string     `json:"status"      extensions:"x-nullable,oneof=not-started in-progress completed"`
	DueDate     *civil.Date `json:"dueDate"     extensions:"x-nullable"                                         swaggertype:"string"`
	Days        *int        `json:"days"        extensions:"x-nullable"`
	Tags        *[]string   `json:"tags"        extensions:"x-nullable"`
}

func MapPackingList(d db.PackingList) PackingListBase {
	return PackingListBase{
		ID:          d.ID,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
		UserID:      d.UserID,
		Name:        d.Name,
		Description: d.Description,
		Status:      d.Status,
		DueDate:     db.IntoCivilDate(d.DueDate),
		Days:        int(d.Days),
		Tags:        d.Tags,
	}
}

func MapPackingListRow(d db.PackingListGetAllByUserRow) PackingListSummary {
	return PackingListSummary{
		PackingListBase: MapPackingList(d.PackingList),
		ItemCount:       int(d.ItemCount),
		ItemPackedCount: int(d.PackedCount),
	}
}
