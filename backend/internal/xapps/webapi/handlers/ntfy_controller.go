package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/hay-kot/hookfeed/backend/internal/services"
	"github.com/hay-kot/hookfeed/backend/internal/services/adapters"
	"github.com/hay-kot/httpkit/server"
	"github.com/rs/zerolog"
)

type NtfyController struct {
	logger             zerolog.Logger
	feedMessageService *services.FeedMessageService
	feedService        *services.FeedService
}

func NewNtfyController(logger zerolog.Logger, feedMessageService *services.FeedMessageService, feedService *services.FeedService) *NtfyController {
	return &NtfyController{
		logger:             logger.With().Str("controller", "ntfy").Logger(),
		feedMessageService: feedMessageService,
		feedService:        feedService,
	}
}

// Publish godoc
//
//	@Tags			Ntfy
//	@Summary		Publish ntfy-compatible notification
//	@Description	Accepts ntfy-style POST/PUT requests for publishing notifications. Supports ntfy headers (X-Title, X-Message, etc.), query parameters (title, message, priority, etc.), JSON body, and plain text body.
//	@Accept			json,text/plain
//	@Produce		json
//	@Param			topic		path		string	true	"Topic/Feed name"
//	@Param			title		query		string	false	"Notification title (alias: t)"
//	@Param			message		query		string	false	"Notification message (alias: m)"
//	@Param			priority	query		int		false	"Priority 1-5, default 3 (alias: p)"
//	@Param			tags		query		string	false	"Comma-separated tags (alias: ta)"
//	@Param			click		query		string	false	"URL opened when notification is clicked"
//	@Param			icon		query		string	false	"URL for custom notification icon"
//	@Param			markdown	query		bool	false	"Enable Markdown rendering (alias: md)"
//	@Param			X-Title		header		string	false	"Notification title (overrides query param)"
//	@Param			X-Priority	header		int		false	"Priority (1-5, default 3, overrides query param)"
//	@Param			X-Tags		header		string	false	"Comma-separated tags (overrides query param)"
//	@Param			X-Message	header		string	false	"Message content (overrides query param)"
//	@Param			X-Click		header		string	false	"Click URL (overrides query param)"
//	@Param			X-Icon		header		string	false	"Icon URL (overrides query param)"
//	@Param			X-Markdown	header		bool	false	"Enable Markdown (overrides query param)"
//	@Param			body		body		string	false	"Message body (plain text or JSON)"
//	@Success		200			{object}	dtos.FeedMessage
//	@Router			/{topic} [POST]
//	@Router			/{topic} [PUT]
func (nc *NtfyController) Publish(w http.ResponseWriter, r *http.Request) error {
	topic := chi.URLParam(r, "topic")
	if topic == "" {
		return server.Error().Msg("topic is required").Write(r.Context(), w)
	}

	nc.logger.Info().
		Str("topic", topic).
		Str("content_type", r.Header.Get("Content-Type")).
		Msg("received ntfy message")

		// Verify feed exists
	cache := nc.feedService.GetCache()
	if cache != nil {
		ok, _ := cache.GetByKey(topic)
		if !ok {
			nc.logger.Warn().Str("topic", topic).Msg("feed not found for ntfy message")
			return server.Error().
				Status(http.StatusNotFound).
				Msg("feed not found").
				Write(r.Context(), w)
		}
	}

	// Use adapter to parse the request
	adapter := &adapters.NtfyAdapter{}
	if err := adapter.UnmarshalRequest(r); err != nil {
		nc.logger.Error().Err(err).Str("topic", topic).Msg("failed to parse ntfy message")
		return server.Error().
			Status(http.StatusBadRequest).
			Msg(err.Error()).
			Write(r.Context(), w)
	}

	nc.logger.Info().
		Str("topic", topic).
		Str("title", adapter.Message.Title).
		Int32("priority", adapter.Message.Priority).
		Strs("tags", adapter.Message.Tags).
		Msg("parsed ntfy message")

	// Convert to DTO and create feed message
	createDTO := adapter.AsFeedMessage()
	feedMessage, err := nc.feedMessageService.Create(r.Context(), createDTO)
	if err != nil {
		nc.logger.Error().Err(err).Str("topic", topic).Msg("failed to create feed message from ntfy")
		return err
	}

	nc.logger.Info().
		Str("message_id", feedMessage.ID.String()).
		Str("topic", topic).
		Str("feed_slug", feedMessage.FeedSlug).
		Msg("ntfy message saved successfully")

	return server.JSON(w, http.StatusOK, feedMessage)
}
