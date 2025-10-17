package handlers

import (
	"net/http"

	"github.com/hay-kot/httpkit/server"
	"github.com/hay-kot/hookfeed/backend/internal/data/dtos"
	"github.com/hay-kot/hookfeed/backend/internal/services"
	"github.com/hay-kot/hookfeed/backend/internal/web/extractors"
)

type FeedMessageController struct {
	service *services.FeedMessageService
}

func NewFeedMessageController(service *services.FeedMessageService) *FeedMessageController {
	return &FeedMessageController{
		service: service,
	}
}

// GetAll godoc
//
//	@Tags			FeedMessage
//	@Summary		List all FeedMessages
//	@Description	List all FeedMessages
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	dtos.PaginationResponse[dtos.FeedMessage]	"A list of FeedMessages"
//	@Param			orderBy	query		string										false	"order by"						Enums(created_at,received_at)
//	@Param			skip	query		int											false	"The number of items to skip"	default(0)
//	@Param			limit	query		int											false	"The number of items to return"	default(100)
//	@Router			/v1/feed-messages [GET]
//	@Security		Bearer
func (uc *FeedMessageController) GetAll(w http.ResponseWriter, r *http.Request) error {
	page, err := extractors.QueryT[services.FeedMessageQuery](r)
	if err != nil {
		return err
	}
	page.Pagination = page.WithDefaults()

	entities, err := uc.service.GetAll(r.Context(), page)
	if err != nil {
		return err
	}

	return server.JSON(w, http.StatusOK, entities)
}

// Get godoc
//
//	@Tags			FeedMessage
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

// GetByFeedSlug godoc
//
//	@Tags			FeedMessage
//	@Summary		List messages by feed slug
//	@Description	List messages for a specific feed with optional filters
//	@Accept			json
//	@Produce		json
//	@Param			feed-slug	path		string	true	"The Feed Slug"
//	@Param			level		query		string	false	"Filter by level"	Enums(info,warning,error,success,debug)
//	@Param			state		query		string	false	"Filter by state"	Enums(new,acknowledged,resolved,archived)
//	@Param			since		query		string	false	"Filter by received date (ISO 8601)"
//	@Param			until		query		string	false	"Filter by received date (ISO 8601)"
//	@Param			skip		query		int		false	"The number of items to skip"	default(0)
//	@Param			limit		query		int		false	"The number of items to return"	default(100)
//	@Success		200			{object}	dtos.PaginationResponse[dtos.FeedMessage]
//	@Router			/v1/feeds/{feed-slug}/messages [GET]
//	@Security		Bearer
func (uc *FeedMessageController) GetByFeedSlug(w http.ResponseWriter, r *http.Request) error {
	feedSlug := r.PathValue("feed-slug")

	query, err := extractors.QueryT[dtos.FeedMessageQuery](r)
	if err != nil {
		return err
	}
	query.Pagination = query.WithDefaults()

	entities, err := uc.service.GetByFeedSlug(r.Context(), feedSlug, query)
	if err != nil {
		return err
	}

	return server.JSON(w, http.StatusOK, entities)
}

// Search godoc
//
//	@Tags			FeedMessage
//	@Summary		Search messages
//	@Description	Search messages with optional filters
//	@Accept			json
//	@Produce		json
//	@Param			feedSlug	query		string	false	"Filter by feed slug"
//	@Param			level		query		string	false	"Filter by level"	Enums(info,warning,error,success,debug)
//	@Param			state		query		string	false	"Filter by state"	Enums(new,acknowledged,resolved,archived)
//	@Param			since		query		string	false	"Filter by received date (ISO 8601)"
//	@Param			until		query		string	false	"Filter by received date (ISO 8601)"
//	@Param			q			query		string	false	"Search query"
//	@Param			skip		query		int		false	"The number of items to skip"	default(0)
//	@Param			limit		query		int		false	"The number of items to return"	default(100)
//	@Success		200			{object}	dtos.PaginationResponse[dtos.FeedMessage]
//	@Router			/v1/messages/search [GET]
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
//	@Tags			FeedMessage
//	@Summary		Create a new FeedMessage
//	@Description	Create a new FeedMessage
//	@Accept			json
//	@Produce		json
//	@Param			FeedMessage	body		dtos.FeedMessageCreate	true	"The FeedMessage details"
//	@Success		201			{object}	dtos.FeedMessage
//	@Router			/v1/feed-messages [POST]
//	@Security		Bearer
func (uc *FeedMessageController) Create(w http.ResponseWriter, r *http.Request) error {
	body, err := extractors.Body[dtos.FeedMessageCreate](r)
	if err != nil {
		return err
	}

	entity, err := uc.service.Create(r.Context(), body)
	if err != nil {
		return err
	}

	return server.JSON(w, http.StatusCreated, entity)
}

// UpdateState godoc
//
//	@Tags			FeedMessage
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
//	@Tags			FeedMessage
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
//	@Tags			FeedMessage
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
//	@Tags			FeedMessage
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
