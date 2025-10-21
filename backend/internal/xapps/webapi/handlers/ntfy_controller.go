package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/hay-kot/hookfeed/backend/internal/data/dtos"
	"github.com/hay-kot/hookfeed/backend/internal/services"
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

// ntfyMessage represents a parsed ntfy-compatible message
type ntfyMessage struct {
	Topic    string            `json:"topic,omitempty"`
	Message  string            `json:"message,omitempty"`
	Title    string            `json:"title,omitempty"`
	Priority int32             `json:"priority,omitempty"`
	Tags     []string          `json:"tags,omitempty"`
	Click    string            `json:"click,omitempty"`
	Icon     string            `json:"icon,omitempty"`
	Actions  []json.RawMessage `json:"actions,omitempty"`
	Markdown bool              `json:"markdown,omitempty"`
}

// Publish godoc
//
//	@Tags			Ntfy
//	@Summary		Publish ntfy-compatible notification
//	@Description	Accepts ntfy-style POST/PUT requests for publishing notifications
//	@Accept			json,text/plain
//	@Produce		json
//	@Param			topic		path		string	true	"Topic/Feed name"
//	@Param			X-Title		header		string	false	"Notification title"
//	@Param			X-Priority	header		int		false	"Priority (1-5, default 3)"
//	@Param			X-Tags		header		string	false	"Comma-separated tags"
//	@Param			X-Message	header		string	false	"Message content"
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
	if nc.feedService != nil {
		cache := nc.feedService.GetCache()
		if cache != nil {
			ok, _ := cache.GetByKey(topic)
			if !ok {
				nc.logger.Warn().
					Str("topic", topic).
					Msg("feed not found for ntfy message")
				return server.Error().
					Status(http.StatusNotFound).
					Msg("feed not found").
					Write(r.Context(), w)
			}
		}
	}

	// Parse the ntfy message from headers and body
	msg, err := nc.parseNtfyMessage(r, topic)
	if err != nil {
		nc.logger.Error().
			Err(err).
			Str("topic", topic).
			Msg("failed to parse ntfy message")
		return server.Error().
			Status(http.StatusBadRequest).
			Msg(err.Error()).
			Write(r.Context(), w)
	}

	nc.logger.Info().
		Str("topic", topic).
		Str("title", msg.Title).
		Int32("priority", msg.Priority).
		Strs("tags", msg.Tags).
		Msg("parsed ntfy message")

	// Read raw request body
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	// Ensure RawRequest is valid JSON
	var rawRequestJSON json.RawMessage
	if json.Valid(bodyBytes) {
		// Body is already valid JSON
		rawRequestJSON = json.RawMessage(bodyBytes)
	} else {
		// Body is plain text, wrap it in JSON
		wrapped := map[string]string{"body": string(bodyBytes)}
		rawRequestJSON, _ = json.Marshal(wrapped)
	}

	// Capture raw headers
	headers := make(map[string]string)
	for k, v := range r.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}
	headersJSON, _ := json.Marshal(headers)

	// Build metadata from ntfy-specific fields
	metadata := make(map[string]interface{})
	if len(msg.Tags) > 0 {
		metadata["tags"] = msg.Tags
	}
	if msg.Click != "" {
		metadata["click"] = msg.Click
	}
	if msg.Icon != "" {
		metadata["icon"] = msg.Icon
	}
	if len(msg.Actions) > 0 {
		metadata["actions"] = msg.Actions
	}
	if msg.Markdown {
		metadata["markdown"] = true
	}
	metadataJSON, _ := json.Marshal(metadata)

	// Create feed message
	title := msg.Title
	message := msg.Message
	priority := msg.Priority

	createDTO := dtos.FeedMessageCreate{
		FeedID:     topic,
		RawRequest: rawRequestJSON,
		RawHeaders: json.RawMessage(headersJSON),
		Title:      &title,
		Message:    &message,
		Priority:   &priority,
		Metadata:   json.RawMessage(metadataJSON),
	}

	feedMessage, err := nc.feedMessageService.Create(r.Context(), createDTO)
	if err != nil {
		nc.logger.Error().
			Err(err).
			Str("topic", topic).
			Msg("failed to create feed message from ntfy")
		return err
	}

	nc.logger.Info().
		Str("message_id", feedMessage.ID.String()).
		Str("topic", topic).
		Str("feed_slug", feedMessage.FeedSlug).
		Msg("ntfy message saved successfully")

	return server.JSON(w, http.StatusOK, feedMessage)
}

// parseNtfyMessage extracts ntfy message fields from headers and body
func (nc *NtfyController) parseNtfyMessage(r *http.Request, topic string) (*ntfyMessage, error) {
	msg := &ntfyMessage{
		Topic:    topic,
		Priority: 3, // default priority
	}

	// Try to parse as JSON first
	if r.Header.Get("Content-Type") == "application/json" {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		// Restore body for later reading
		r.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

		var jsonMsg ntfyMessage
		if err := json.Unmarshal(bodyBytes, &jsonMsg); err == nil {
			// JSON parsed successfully
			if jsonMsg.Topic != "" {
				msg.Topic = jsonMsg.Topic
			}
			if jsonMsg.Message != "" {
				msg.Message = jsonMsg.Message
			}
			if jsonMsg.Title != "" {
				msg.Title = jsonMsg.Title
			}
			if jsonMsg.Priority > 0 {
				msg.Priority = jsonMsg.Priority
			}
			if len(jsonMsg.Tags) > 0 {
				msg.Tags = jsonMsg.Tags
			}
			msg.Click = jsonMsg.Click
			msg.Icon = jsonMsg.Icon
			msg.Actions = jsonMsg.Actions
			msg.Markdown = jsonMsg.Markdown
		}
	}

	// Parse headers (headers override JSON if present)
	if title := nc.getHeader(r, "X-Title", "Title"); title != "" {
		msg.Title = title
	}

	if msgText := nc.getHeader(r, "X-Message", "Message"); msgText != "" {
		msg.Message = msgText
	}

	if priorityStr := nc.getHeader(r, "X-Priority", "Priority"); priorityStr != "" {
		if p, err := nc.parsePriority(priorityStr); err == nil {
			msg.Priority = p
		}
	}

	if tagsStr := nc.getHeader(r, "X-Tags", "Tags"); tagsStr != "" {
		msg.Tags = strings.Split(tagsStr, ",")
		// Trim whitespace from tags
		for i := range msg.Tags {
			msg.Tags[i] = strings.TrimSpace(msg.Tags[i])
		}
	}

	if click := nc.getHeader(r, "X-Click", "Click"); click != "" {
		msg.Click = click
	}

	if icon := nc.getHeader(r, "X-Icon", "Icon"); icon != "" {
		msg.Icon = icon
	}

	if markdown := nc.getHeader(r, "X-Markdown", "Markdown"); markdown != "" {
		msg.Markdown = markdown == "true" || markdown == "1" || markdown == "yes"
	}

	// If message is still empty, read from body as plain text
	if msg.Message == "" {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		// Restore body for later reading
		r.Body = io.NopCloser(strings.NewReader(string(bodyBytes)))

		msg.Message = string(bodyBytes)
	}

	return msg, nil
}

// getHeader retrieves a header value, trying multiple possible keys
func (nc *NtfyController) getHeader(r *http.Request, keys ...string) string {
	for _, key := range keys {
		if val := r.Header.Get(key); val != "" {
			return val
		}
	}
	return ""
}

// parsePriority converts priority string to int32 (1-5)
func (nc *NtfyController) parsePriority(s string) (int32, error) {
	// Handle named priorities (ntfy compatibility)
	switch strings.ToLower(s) {
	case "min", "1":
		return 1, nil
	case "low", "2":
		return 2, nil
	case "default", "3", "":
		return 3, nil
	case "high", "4":
		return 4, nil
	case "max", "urgent", "5":
		return 5, nil
	}

	// Try parsing as integer
	p, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 3, err
	}

	// Clamp to valid range
	if p < 1 {
		p = 1
	}
	if p > 5 {
		p = 5
	}

	return int32(p), nil
}
