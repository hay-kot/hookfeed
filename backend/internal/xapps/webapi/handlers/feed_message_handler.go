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
//	@Param			orderBy	query		string										false	"order by"						Enums(created_at)
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

// Update godoc
//
//	@Tags			FeedMessage
//	@Summary		Update a FeedMessage
//	@Description	Update a FeedMessage
//	@Accept			json
//	@Produce		json
//	@Param			FeedMessage	body		dtos.FeedMessageUpdate	true	"The FeedMessage details"
//	@Param			id			path		string					true	"The FeedMessage ID"
//	@Success		200			{object}	dtos.FeedMessage
//	@Router			/v1/feed-messages/{id} [PUT]
//	@Security		Bearer
func (uc *FeedMessageController) Update(w http.ResponseWriter, r *http.Request) error {
	id, body, err := extractors.BodyWithID[dtos.FeedMessageUpdate](r, "id")
	if err != nil {
		return err
	}

	entity, err := uc.service.Update(r.Context(), id, body)
	if err != nil {
		return err
	}

	return server.JSON(w, http.StatusOK, entity)
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
