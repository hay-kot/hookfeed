package dtos

import (
	"time"

	"github.com/google/uuid"
	"github.com/hay-kot/hookfeed/backend/internal/data/db"
)

type PackingListItem struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `json:"name"`
	Category  string    `json:"category"`
	Quantity  int       `json:"quantity"`
	IsPacked  bool      `json:"isPacked"`
	Notes     string    `json:"notes"`
}

type PackingListItemCreate struct {
	Name     string `json:"name"     validate:"required"`
	Category string `json:"category"`
	Quantity int    `json:"quantity" validate:"min=1"`
	Notes    string `json:"notes"`
}

type PackingListItemUpdate struct {
	Name     *string `json:"name"     extensions:"x-nullable"`
	Category *string `json:"category" extensions:"x-nullable"`
	Quantity *int    `json:"quantity" extensions:"x-nullable"`
	IsPacked *bool   `json:"isPacked" extensions:"x-nullable"`
	Notes    *string `json:"notes"    extensions:"x-nullable"`
}

func MapPackingListItem(d db.PackingListItem) PackingListItem {
	return PackingListItem{
		ID:        d.ID,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
		Name:      d.Name,
		Category:  d.Category,
		Quantity:  int(d.Quantity),
		IsPacked:  d.IsPacked,
		Notes:     d.Notes,
	}
}
