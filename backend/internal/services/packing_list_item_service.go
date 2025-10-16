package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/hay-kot/hookfeed/backend/internal/data/db"
	"github.com/hay-kot/hookfeed/backend/internal/data/dtos"
	"github.com/hay-kot/hookfeed/backend/pkgs/utils"
	"github.com/rs/zerolog"
)

type PackingListItemQuery struct {
	dtos.Pagination

	OrderBy string `json:"orderBy" validate:"omitempty,oneof=created_at" query:"orderBy"`
}

type PackingListItemService struct {
	l      zerolog.Logger
	db     *db.QueriesExt
	mapper dtos.MapFunc[db.PackingListItem, dtos.PackingListItem]
}

func NewPackingListItemService(l zerolog.Logger, db *db.QueriesExt) *PackingListItemService {
	return &PackingListItemService{
		l:      l,
		db:     db,
		mapper: dtos.MapPackingListItem,
	}
}

func (s *PackingListItemService) Get(ctx context.Context, id uuid.UUID) (dtos.PackingListItem, error) {
	entity, err := s.db.PackingListItemByID(ctx, id)
	if err != nil {
		return dtos.PackingListItem{}, err
	}

	return s.mapper(entity), nil
}

func (s *PackingListItemService) Create(ctx context.Context, listID uuid.UUID, data dtos.PackingListItemCreate) (dtos.PackingListItem, error) {
	e, err := s.db.PackingListItemCreate(ctx, db.PackingListItemCreateParams{
		PackingListID: listID,
		Name:          data.Name,
		Category:      data.Category,
		Notes:         data.Notes,
		Quantity:      int32(data.Quantity),
	})
	if err != nil {
		return dtos.PackingListItem{}, err
	}

	return s.mapper(e), nil
}

func (s *PackingListItemService) Update(ctx context.Context, listId, id uuid.UUID, data dtos.PackingListItemUpdate) (dtos.PackingListItem, error) {
	e, err := s.db.PackingListItemUpdate(ctx, db.PackingListItemUpdateParams{
		Name:     data.Name,
		Category: data.Category,
		Notes:    data.Notes,
		Quantity: utils.PtrInt32(data.Quantity),
		IsPacked: data.IsPacked,
		ID:       id,
	})
	if err != nil {
		return dtos.PackingListItem{}, err
	}

	return s.mapper(e), nil
}

func (s *PackingListItemService) Delete(ctx context.Context, listId, id uuid.UUID) error {
	err := s.db.PackingListItemDeleteByID(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
