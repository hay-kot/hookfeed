package handlers

import (
	"net/http"

	"github.com/hay-kot/hookfeed/backend/internal/data/dtos"
	"github.com/hay-kot/hookfeed/backend/internal/services"
	"github.com/hay-kot/hookfeed/backend/internal/web/extractors"
	"github.com/hay-kot/httpkit/server"
)

type WebhookController struct {
	webhookService *services.WebhookService
}

func NewWebhookController(webhookService *services.WebhookService) *WebhookController {
	return &WebhookController{
		webhookService: webhookService,
	}
}

// HandleWebhook godoc
//
//	@Tags			Webhooks
//	@Summary		Receive webhook
//	@Description	Accepts webhooks in any format and processes them according to feed configuration
//	@Accept			json
//	@Produce		json
//	@Param			key		path		string	true	"Feed key (from feed configuration)"
//	@Param			body	body		object	true	"Webhook payload (any JSON)"
//	@Success		202		{object}	dtos.WebhookResponse
//	@Failure		400		{object}	server.ErrorResp
//	@Failure		401		{object}	server.ErrorResp
//	@Failure		404		{object}	server.ErrorResp
//	@Failure		500		{object}	server.ErrorResp
//	@Router			/hooks/{slug} [POST]
func (wc *WebhookController) HandleWebhook(w http.ResponseWriter, r *http.Request) error {
	// Extract key from URL
	key, err := extractors.Slug(r, "key")
	if err != nil {
		return err
	}

	val := map[string]any{}
	err = server.Decode(r, &val)
	if err != nil {
		return err
	}

	// Build the webhook request
	webhookReq := dtos.WebhookRequest{
		FeedKey: key,
		Headers: r.Header,
		Body:    val,
	}

	// Process the webhook
	response, err := wc.webhookService.ProcessWebhook(webhookReq)
	if err != nil {
		return err
	}

	// Return 202 Accepted with the response
	return server.JSON(w, http.StatusAccepted, response)
}
