package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/hay-kot/hookfeed/backend/internal/data/db"
	"github.com/hay-kot/hookfeed/backend/internal/data/dtos"
	"github.com/hay-kot/hookfeed/backend/pkgs/utils"
	"github.com/rs/zerolog"
)

type PackingListService struct {
	l         zerolog.Logger
	db        *db.QueriesExt
	mapper    dtos.MapFunc[db.PackingList, dtos.PackingListBase]
	mapperRow dtos.MapFunc[db.PackingListGetAllByUserRow, dtos.PackingListSummary]
}

func NewPackingListService(l zerolog.Logger, db *db.QueriesExt) *PackingListService {
	return &PackingListService{
		l:         l,
		db:        db,
		mapper:    dtos.MapPackingList,
		mapperRow: dtos.MapPackingListRow,
	}
}

func (s *PackingListService) Get(ctx context.Context, id uuid.UUID) (dtos.PackingList, error) {
	entity, err := s.db.PackingListByID(ctx, id)
	if err != nil {
		return dtos.PackingList{}, err
	}

	items, err := s.db.PackingListItemGetAll(ctx, id)
	if err != nil {
		return dtos.PackingList{}, err
	}

	mapped := dtos.PackingList{
		PackingListBase: s.mapper(entity),
		Items:           dtos.MapFunc[db.PackingListItem, dtos.PackingListItem](dtos.MapPackingListItem).Slice(items),
	}

	return mapped, nil
}

func (s *PackingListService) GetAllByUser(ctx context.Context, userID uuid.UUID, page dtos.PackingListQuery) (dtos.PaginationResponse[dtos.PackingListSummary], error) {
	count, err := s.db.PackingListGetAllByUserCount(ctx, userID)
	if err != nil {
		return dtos.PaginationResponse[dtos.PackingListSummary]{}, err
	}
	orderBy, err := db.OrderByWithDefaults(&page.OrderBy, "created_at", utils.Ptr("desc"), "desc")
	if err != nil {
		return dtos.PaginationResponse[dtos.PackingListSummary]{}, err
	}

	entities, err := s.db.PackingListGetAllByUser(ctx, db.PackingListGetAllByUserParams{
		UserID:  userID,
		Limit:   int32(page.Limit),
		Offset:  int32(page.Skip),
		OrderBy: orderBy,
	})
	if err != nil {
		return dtos.PaginationResponse[dtos.PackingListSummary]{}, err
	}

	return dtos.PaginationResponse[dtos.PackingListSummary]{
		Total: int(count),
		Items: s.mapperRow.Slice(entities),
	}, nil
}

func (s *PackingListService) Create(ctx context.Context, userID uuid.UUID, data dtos.PackingListCreate) (dtos.PackingList, error) {
	e, err := s.db.PackingListCreate(ctx, db.PackingListCreateParams{
		UserID:      userID,
		Name:        data.Name,
		Description: data.Description,
		DueDate:     db.IntoPgDate(data.DueDate),
		Days:        int32(data.Days),
		Tags:        utils.ToNotNil(data.Tags),
	})
	if err != nil {
		return dtos.PackingList{}, err
	}

	return s.Get(ctx, e.ID)
}

func (s *PackingListService) Update(ctx context.Context, userID, id uuid.UUID, data dtos.PackingListUpdate) (dtos.PackingList, error) {
	e, err := s.db.PackingListUpdate(ctx, db.PackingListUpdateParams{
		Name:        data.Name,
		Description: data.Description,
		UserID:      userID,
		ID:          id,
		DueDate:     db.IntoPgDate(data.DueDate),
		Status:      data.Status,
		Days:        utils.PtrInt32(data.Days),
		UpdateTags:  data.Tags != nil, // only update it ptr is not nil
		Tags:        utils.PtrToSlice(data.Tags),
	})
	if err != nil {
		return dtos.PackingList{}, err
	}

	return s.Get(ctx, e.ID)
}

func (s *PackingListService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.db.PackingListDeleteByID(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
