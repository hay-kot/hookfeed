package handlers

import (
	"net/http"

	"github.com/hay-kot/hookfeed/backend/internal/data/dtos"
	"github.com/hay-kot/hookfeed/backend/internal/services"
	"github.com/hay-kot/hookfeed/backend/internal/services/adapters"
	"github.com/hay-kot/hookfeed/backend/internal/web/extractors"
	"github.com/hay-kot/httpkit/server"
)

type FeedMessageController struct {
	service *services.FeedMessageService
}

func NewFeedMessageController(service *services.FeedMessageService) *FeedMessageController {
	return &FeedMessageController{
		service: service,
	}
}

// Get godoc
//
//	@Tags			Feed Messages
//	@Summary		Get a FeedMessage
//	@Description	Get a FeedMessage
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"The FeedMessage ID"
//	@Success		200	{object}	dtos.FeedMessage
//	@Router			/v1/feed-messages/{id} [GET]
//	@Security		Bearer
func (uc *FeedMessageController) Get(w http.ResponseWriter, r *http.Request) error {
	id, err := extractors.ID(r, "id")
	if err != nil {
		return err
	}

	entity, err := uc.service.Get(r.Context(), id)
	if err != nil {
		return err
	}

	return server.JSON(w, http.StatusOK, entity)
}

// Search godoc
//
//	@Tags			Feed Messages
//	@Summary		Search messages
//	@Description	Search messages with optional filters
//	@Accept			json
//	@Produce		json
//	@Param			feedSlug	query		string	false	"Filter by feed slug"
//	@Param			priority	query		int		false	"Filter by priority (1-5)"	minimum(1)	maximum(5)
//	@Param			state		query		string	false	"Filter by state"			Enums(new,acknowledged,resolved,archived)
//	@Param			since		query		string	false	"Filter by received date (ISO 8601)"
//	@Param			until		query		string	false	"Filter by received date (ISO 8601)"
//	@Param			q			query		string	false	"Search query"
//	@Param			skip		query		int		false	"The number of items to skip"	default(0)
//	@Param			limit		query		int		false	"The number of items to return"	default(100)
//	@Success		200			{object}	dtos.PaginationResponse[dtos.FeedMessage]
//	@Router			/v1/feed-messages [GET]
//	@Security		Bearer
func (uc *FeedMessageController) Search(w http.ResponseWriter, r *http.Request) error {
	query, err := extractors.QueryT[dtos.FeedMessageQuery](r)
	if err != nil {
		return err
	}
	query.Pagination = query.WithDefaults()

	entities, err := uc.service.Search(r.Context(), query)
	if err != nil {
		return err
	}

	return server.JSON(w, http.StatusOK, entities)
}

// Create godoc
//
//	@Tags			Feed Messages
//	@Summary		Create a new FeedMessage
//	@Description	Create a new FeedMessage
//	@Accept			json
//	@Produce		json
//	@Param			FeedMessage	body		dtos.FeedMessageCreate	true	"The FeedMessage details"
//	@Success		201			{object}	dtos.FeedMessage
//	@Router			/v1/feed-messages [POST]
//	@Security		Bearer
func (uc *FeedMessageController) Create(w http.ResponseWriter, r *http.Request) error {
	// Use adapter to parse the request
	adapter := &adapters.RawAdapter{}
	if err := adapter.UnmarshalRequest(r); err != nil {
		return server.Error().
			Status(http.StatusBadRequest).
			Msg(err.Error()).
			Write(r.Context(), w)
	}

	// Convert to DTO and create feed message
	createDTO := adapter.AsFeedMessage()
	entity, err := uc.service.Create(r.Context(), createDTO)
	if err != nil {
		return err
	}

	return server.JSON(w, http.StatusCreated, entity)
}

// UpdateState godoc
//
//	@Tags			Feed Messages
//	@Summary		Update a FeedMessage state
//	@Description	Update a FeedMessage state
//	@Accept			json
//	@Produce		json
//	@Param			body	body		dtos.FeedMessageUpdateState	true	"The state update"
//	@Param			id		path		string						true	"The FeedMessage ID"
//	@Success		200		{object}	dtos.FeedMessage
//	@Router			/v1/feed-messages/{id}/state [PATCH]
//	@Security		Bearer
func (uc *FeedMessageController) UpdateState(w http.ResponseWriter, r *http.Request) error {
	id, err := extractors.ID(r, "id")
	if err != nil {
		return err
	}

	body, err := extractors.Body[dtos.FeedMessageUpdateState](r)
	if err != nil {
		return err
	}

	entity, err := uc.service.UpdateState(r.Context(), id, body.State)
	if err != nil {
		return err
	}

	return server.JSON(w, http.StatusOK, entity)
}

// BulkUpdateState godoc
//
//	@Tags			Feed Messages
//	@Summary		Bulk update message states
//	@Description	Bulk update message states for a feed
//	@Accept			json
//	@Produce		json
//	@Param			feed-slug	path	string							true	"The Feed Slug"
//	@Param			body		body	dtos.FeedMessageBulkUpdateState	true	"The bulk update request"
//	@Success		204
//	@Router			/v1/feeds/{feed-slug}/messages/bulk-state [POST]
//	@Security		Bearer
func (uc *FeedMessageController) BulkUpdateState(w http.ResponseWriter, r *http.Request) error {
	body, err := extractors.Body[dtos.FeedMessageBulkUpdateState](r)
	if err != nil {
		return err
	}

	err = uc.service.BulkUpdateState(r.Context(), body)
	if err != nil {
		return err
	}

	return server.JSON(w, http.StatusNoContent, nil)
}

// BulkDelete godoc
//
//	@Tags			Feed Messages
//	@Summary		Bulk delete messages
//	@Description	Bulk delete messages for a feed
//	@Accept			json
//	@Produce		json
//	@Param			feed-slug	path		string						true	"The Feed Slug"
//	@Param			body		body		dtos.FeedMessageBulkDelete	true	"The bulk delete request"
//	@Success		200			{object}	map[string]int
//	@Router			/v1/feeds/{feed-slug}/messages/bulk-delete [POST]
//	@Security		Bearer
func (uc *FeedMessageController) BulkDelete(w http.ResponseWriter, r *http.Request) error {
	body, err := extractors.Body[dtos.FeedMessageBulkDelete](r)
	if err != nil {
		return err
	}

	count, err := uc.service.BulkDelete(r.Context(), body)
	if err != nil {
		return err
	}

	return server.JSON(w, http.StatusOK, map[string]int{"deleted": count})
}

// Delete godoc
//
//	@Tags			Feed Messages
//	@Summary		Delete a FeedMessage
//	@Description	Delete a FeedMessage
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"The FeedMessage ID"
//	@Success		204
//	@Router			/v1/feed-messages/{id} [DELETE]
//	@Security		Bearer
func (uc *FeedMessageController) Delete(w http.ResponseWriter, r *http.Request) error {
	id, err := extractors.ID(r, "id")
	if err != nil {
		return err
	}

	err = uc.service.Delete(r.Context(), id)
	if err != nil {
		return err
	}

	return server.JSON(w, http.StatusNoContent, nil)
}
