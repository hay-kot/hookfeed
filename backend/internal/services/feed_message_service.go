package services

import (
	"context"

	"github.com/google/uuid"

	"github.com/hay-kot/hookfeed/backend/internal/data/db"
	"github.com/hay-kot/hookfeed/backend/internal/data/dtos"
	"github.com/rs/zerolog"
)

type FeedMessageQuery struct {
	dtos.Pagination
	OrderBy string `json:"orderBy" query:"orderBy" validate:"omitempty,oneof=created_at"`
}

type FeedMessageService struct {
	l      zerolog.Logger
	db     *db.QueriesExt
	mapper dtos.MapFunc[db.FeedMessage, dtos.FeedMessage]
}

func NewFeedMessageService(l zerolog.Logger, db *db.QueriesExt) *FeedMessageService {
	return &FeedMessageService{
		l:      l,
		db:     db,
		mapper: dtos.MapFeedMessage,
	}
}

func (s *FeedMessageService) Get(ctx context.Context, id uuid.UUID) (dtos.FeedMessage, error) {
	entity, err := s.db.FeedMessageByID(ctx, id)
	if err != nil {
		return dtos.FeedMessage{}, err
	}

	return s.mapper(entity), nil
}

func (s *FeedMessageService) GetAll(ctx context.Context, page FeedMessageQuery) (dtos.PaginationResponse[dtos.FeedMessage], error) {
	count, err := s.db.FeedMessageGetAllCount(ctx)
	if err != nil {
		return dtos.PaginationResponse[dtos.FeedMessage]{}, err
	}

	entities, err := s.db.FeedMessageGetAll(ctx, db.FeedMessageGetAllParams{
		Limit:   int32(page.Limit),
		Offset:  int32(page.Skip),
		OrderBy: page.OrderBy,
	})
	if err != nil {
		return dtos.PaginationResponse[dtos.FeedMessage]{}, err
	}

	return dtos.PaginationResponse[dtos.FeedMessage]{
		Total: int(count),
		Items: s.mapper.Slice(entities),
	}, nil
}

func (s *FeedMessageService) Create(ctx context.Context, data dtos.FeedMessageCreate) (dtos.FeedMessage, error) {
	panic("not implemented")
}

func (s *FeedMessageService) Update(ctx context.Context, id uuid.UUID, data dtos.FeedMessageUpdate) (dtos.FeedMessage, error) {
	panic("not implemented")
}

func (s *FeedMessageService) Delete(ctx context.Context, id uuid.UUID) error {
	err := s.db.FeedMessageDeleteByID(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
