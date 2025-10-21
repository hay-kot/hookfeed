package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/hay-kot/hookfeed/backend/internal/data/db"
	"github.com/hay-kot/hookfeed/backend/internal/data/dtos"
	"github.com/rs/zerolog"
)

type FeedMessageQuery struct {
	dtos.Pagination
	OrderBy string `json:"orderBy" validate:"omitempty,oneof=created_at received_at" query:"orderBy"`
}

type FeedMessageService struct {
	l      zerolog.Logger
	db     *db.QueriesExt
	mapper dtos.MapFunc[db.FeedMessagesView, dtos.FeedMessage]
}

func NewFeedMessageService(l zerolog.Logger, db *db.QueriesExt) *FeedMessageService {
	return &FeedMessageService{
		l:      l,
		db:     db,
		mapper: dtos.MapFeedMessageView,
	}
}

func (s *FeedMessageService) Get(ctx context.Context, id uuid.UUID) (dtos.FeedMessage, error) {
	row, err := s.db.FeedMessageByID(ctx, id)
	if err != nil {
		return dtos.FeedMessage{}, err
	}

	return s.mapper(row.FeedMessagesView), nil
}

func (s *FeedMessageService) GetAll(ctx context.Context, page FeedMessageQuery) (dtos.PaginationResponse[dtos.FeedMessage], error) {
	count, err := s.db.FeedMessageGetAllCount(ctx)
	if err != nil {
		return dtos.PaginationResponse[dtos.FeedMessage]{}, err
	}

	orderBy := page.OrderBy
	if orderBy == "" {
		orderBy = "received_at:desc"
	}

	rows, err := s.db.FeedMessageGetAll(ctx, db.FeedMessageGetAllParams{
		Limit:   int32(page.Limit),
		Offset:  int32(page.Skip),
		OrderBy: orderBy,
	})
	if err != nil {
		return dtos.PaginationResponse[dtos.FeedMessage]{}, err
	}

	// Extract views from row wrappers
	views := make([]db.FeedMessagesView, len(rows))
	for i, row := range rows {
		views[i] = row.FeedMessagesView
	}

	return dtos.PaginationResponse[dtos.FeedMessage]{
		Total: int(count),
		Items: s.mapper.Slice(views),
	}, nil
}

func (s *FeedMessageService) GetByFeedSlug(ctx context.Context, feedSlug string, query dtos.FeedMessageQuery) (dtos.PaginationResponse[dtos.FeedMessage], error) {
	count, err := s.db.FeedMessagesByFeedSlugCount(ctx, db.FeedMessagesByFeedSlugCountParams{
		FeedSlug: feedSlug,
		Priority: query.Priority,
		State:    query.State,
		Since:    timePtrToPgTimestamp(query.Since),
		Until:    timePtrToPgTimestamp(query.Until),
	})
	if err != nil {
		return dtos.PaginationResponse[dtos.FeedMessage]{}, err
	}

	rows, err := s.db.FeedMessagesByFeedSlug(ctx, db.FeedMessagesByFeedSlugParams{
		FeedSlug: feedSlug,
		Priority: query.Priority,
		State:    query.State,
		Since:    timePtrToPgTimestamp(query.Since),
		Until:    timePtrToPgTimestamp(query.Until),
		Limit:    int32(query.Limit),
		Offset:   int32(query.Skip),
	})
	if err != nil {
		return dtos.PaginationResponse[dtos.FeedMessage]{}, err
	}

	// Extract views from row wrappers
	views := make([]db.FeedMessagesView, len(rows))
	for i, row := range rows {
		views[i] = row.FeedMessagesView
	}

	return dtos.PaginationResponse[dtos.FeedMessage]{
		Total: int(count),
		Items: s.mapper.Slice(views),
	}, nil
}

func (s *FeedMessageService) Search(ctx context.Context, query dtos.FeedMessageQuery) (dtos.PaginationResponse[dtos.FeedMessage], error) {
	count, err := s.db.FeedMessageSearchCount(ctx, db.FeedMessageSearchCountParams{
		FeedSlug: query.FeedSlug,
		Priority: query.Priority,
		State:    query.State,
		Since:    timePtrToPgTimestamp(query.Since),
		Until:    timePtrToPgTimestamp(query.Until),
		Query:    query.Query,
	})
	if err != nil {
		return dtos.PaginationResponse[dtos.FeedMessage]{}, err
	}

	rows, err := s.db.FeedMessageSearch(ctx, db.FeedMessageSearchParams{
		FeedSlug: query.FeedSlug,
		Priority: query.Priority,
		State:    query.State,
		Since:    timePtrToPgTimestamp(query.Since),
		Until:    timePtrToPgTimestamp(query.Until),
		Query:    query.Query,
		Limit:    int32(query.Limit),
		Offset:   int32(query.Skip),
	})
	if err != nil {
		return dtos.PaginationResponse[dtos.FeedMessage]{}, err
	}

	// Extract views from row wrappers
	views := make([]db.FeedMessagesView, len(rows))
	for i, row := range rows {
		views[i] = row.FeedMessagesView
	}

	return dtos.PaginationResponse[dtos.FeedMessage]{
		Total: int(count),
		Items: s.mapper.Slice(views),
	}, nil
}

func (s *FeedMessageService) Create(ctx context.Context, data dtos.FeedMessageCreate) (dtos.FeedMessage, error) {
	priority := int32(3)
	if data.Priority != nil {
		priority = *data.Priority
	}

	state := "new"
	if data.State != nil {
		state = *data.State
	}

	receivedAt := time.Now()
	if data.ReceivedAt != nil {
		receivedAt = *data.ReceivedAt
	}

	row, err := s.db.FeedMessageCreate(ctx, db.FeedMessageCreateParams{
		FeedSlug:    data.FeedID,
		RawRequest:  []byte(data.RawRequest),
		RawHeaders:  []byte(data.RawHeaders),
		Title:       data.Title,
		Message:     data.Message,
		Priority:    &priority,
		Logs:        data.Logs,
		Metadata:    []byte(data.Metadata),
		State:       &state,
		ReceivedAt:  receivedAt,
		ProcessedAt: timePtrToPgTimestamp(data.ProcessedAt),
	})
	if err != nil {
		return dtos.FeedMessage{}, err
	}

	// Convert row to view type
	view := db.FeedMessagesView(row)

	return s.mapper(view), nil
}

func (s *FeedMessageService) UpdateState(ctx context.Context, id uuid.UUID, state string) (dtos.FeedMessage, error) {
	row, err := s.db.FeedMessageUpdateState(ctx, db.FeedMessageUpdateStateParams{
		ID:    id,
		State: &state,
	})
	if err != nil {
		return dtos.FeedMessage{}, err
	}

	// Convert row to view type
	view := db.FeedMessagesView(row)

	return s.mapper(view), nil
}

func (s *FeedMessageService) BulkUpdateState(ctx context.Context, data dtos.FeedMessageBulkUpdateState) error {
	return s.db.FeedMessageBulkUpdateState(ctx, db.FeedMessageBulkUpdateStateParams{
		Column1: data.MessageIDs,
		State:   &data.State,
	})
}

func (s *FeedMessageService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.db.FeedMessageDeleteByID(ctx, id)
}

func (s *FeedMessageService) BulkDelete(ctx context.Context, data dtos.FeedMessageBulkDelete) (int, error) {
	if len(data.MessageIDs) > 0 {
		count, err := s.db.FeedMessageBulkDelete(ctx, data.MessageIDs)
		return int(count), err
	}

	if data.Filter != nil && data.Filter.OlderThan != nil {
		// This assumes we have a feedID - we might need to adjust this
		// based on how the bulk delete API is designed
		return 0, nil
	}

	return 0, nil
}

// Helper functions
func timePtrToPgTimestamp(t *time.Time) pgtype.Timestamp {
	if t == nil {
		return pgtype.Timestamp{}
	}
	return pgtype.Timestamp{
		Time:  *t,
		Valid: true,
	}
}
