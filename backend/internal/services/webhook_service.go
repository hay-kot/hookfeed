package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/hay-kot/hookfeed/backend/internal/core/feeds"
	"github.com/hay-kot/hookfeed/backend/internal/data/dtos"
	"github.com/rs/zerolog"
)

var (
	ErrFeedNotFound  = errors.New("feed not found")
	ErrInvalidAPIKey = errors.New("invalid API key")
	ErrInvalidJSON   = errors.New("invalid JSON body")
	ErrFeedNotInit   = errors.New("feed service not initialized")
)

type WebhookService struct {
	logger             zerolog.Logger
	feedService        *FeedService
	feedMessageService *FeedMessageService
}

func NewWebhookService(logger zerolog.Logger, feedService *FeedService, feedMessageService *FeedMessageService) *WebhookService {
	return &WebhookService{
		logger:             logger.With().Str("service", "webhook").Logger(),
		feedService:        feedService,
		feedMessageService: feedMessageService,
	}
}

// ProcessWebhook handles the incoming webhook request
func (w *WebhookService) ProcessWebhook(ctx context.Context, req dtos.WebhookRequest) (*dtos.WebhookResponse, error) {
	// Log the incoming request
	w.logger.Info().
		Str("slug", req.FeedKey).
		Int("body_size", len(req.Body)).
		Msg("received webhook request")

	// Find the feed by slug
	feed, err := w.findFeedBySlug(req.FeedKey)
	if err != nil {
		w.logger.Error().
			Err(err).
			Str("slug", req.FeedKey).
			Msg("feed not found")
		return nil, err
	}

	w.logger.Info().
		Str("feed_id", feed.ID).
		Str("feed_name", feed.Name).
		Msg("matched feed")

	w.logger.Info().
		Interface("body", req.Body).
		Interface("headers", req.Headers).
		Msg("webhook payload")

	// Convert request body and headers to JSON
	rawRequest, err := json.Marshal(req.Body)
	if err != nil {
		w.logger.Error().
			Err(err).
			Msg("failed to marshal request body")
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	rawHeaders, err := json.Marshal(req.Headers)
	if err != nil {
		w.logger.Error().
			Err(err).
			Msg("failed to marshal request headers")
		return nil, fmt.Errorf("failed to marshal request headers: %w", err)
	}

	// Create the feed message
	receivedAt := time.Now()
	createMsg := dtos.FeedMessageCreate{
		FeedID:   feed.ID,
		RawRequest: json.RawMessage(rawRequest),
		RawHeaders: json.RawMessage(rawHeaders),
		ReceivedAt: &receivedAt,
		Logs:       []string{},
		Metadata:   json.RawMessage("{}"),
	}

	// Save message to database
	message, err := w.feedMessageService.Create(ctx, createMsg)
	if err != nil {
		w.logger.Error().
			Err(err).
			Str("feed_id", feed.ID).
			Msg("failed to save message to database")
		return nil, fmt.Errorf("failed to save message: %w", err)
	}

	w.logger.Info().
		Str("message_id", message.ID.String()).
		Str("feed_id", feed.ID).
		Msg("webhook processed and saved successfully")

	// TODO: In future iterations, we'll:
	// - Execute global middleware
	// - Execute feed middleware
	// - Apply adapters
	// - Broadcast via WebSocket
	// - Enforce retention policies

	return &dtos.WebhookResponse{
		Success:   true,
		MessageID: message.ID,
		FeedID:    feed.ID,
	}, nil
}

// findFeedBySlug looks up a feed by its key
func (w *WebhookService) findFeedBySlug(slug string) (feeds.FeedParsed, error) {
	if w.feedService == nil {
		return feeds.FeedParsed{}, ErrFeedNotInit
	}

	// Try to look up by key using the cache
	cache := w.feedService.GetCache()
	if cache == nil {
		return feeds.FeedParsed{}, ErrFeedNotInit
	}

	// Look up by key only
	ok, feed := cache.GetByKey(slug)
	if ok {
		return feed, nil
	}

	return feeds.FeedParsed{}, fmt.Errorf("%w: %s", ErrFeedNotFound, slug)
}
