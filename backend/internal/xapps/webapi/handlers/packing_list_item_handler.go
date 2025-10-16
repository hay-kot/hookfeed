package handlers

import (
	"net/http"

	"github.com/hay-kot/hookfeed/backend/internal/data/dtos"
	"github.com/hay-kot/hookfeed/backend/internal/services"
	"github.com/hay-kot/hookfeed/backend/internal/web/extractors"
	"github.com/hay-kot/httpkit/server"
)

type PackingListItemController struct {
	service *services.PackingListItemService
}

func NewPackingListItemController(service *services.PackingListItemService) *PackingListItemController {
	return &PackingListItemController{
		service: service,
	}
}

// Create godoc
//
//	@Tags			PackingListItem
//	@Summary		Create a new PackingListItem
//	@Description	Create a new PackingListItem
//	@Accept			json
//	@Produce		json
//	@Param			PackingListItem	body		dtos.PackingListItemCreate	true	"The PackingListItem details"
//	@Param			id				path		string						true	"The PackingList ID"
//	@Success		201				{object}	dtos.PackingListItem
//	@Router			/v1/packing-lists/{id}/items [POST]
//	@Security		Bearer
func (uc *PackingListItemController) Create(w http.ResponseWriter, r *http.Request) error {
	id, body, err := extractors.BodyWithID[dtos.PackingListItemCreate](r, "id")
	if err != nil {
		return err
	}

	entity, err := uc.service.Create(r.Context(), id, body)
	if err != nil {
		return err
	}

	return server.JSON(w, http.StatusCreated, entity)
}

// Update godoc
//
//	@Tags			PackingListItem
//	@Summary		Update a PackingListItem
//	@Description	Update a PackingListItem
//	@Accept			json
//	@Produce		json
//	@Param			PackingListItem	body		dtos.PackingListItemUpdate	true	"The PackingListItem details"
//	@Param			id				path		string						true	"The PackingList ID"
//	@Param			item-id			path		string						true	"The PackingListItem ID"
//	@Success		200				{object}	dtos.PackingListItem
//	@Router			/v1/packing-lists/{id}/items/{item-id} [PUT]
//	@Security		Bearer
func (uc *PackingListItemController) Update(w http.ResponseWriter, r *http.Request) error {
	listId, id, err := extractors.ID2(r, "id", "item-id")
	if err != nil {
		return err
	}

	body, err := extractors.Body[dtos.PackingListItemUpdate](r)
	if err != nil {
		return err
	}

	entity, err := uc.service.Update(r.Context(), listId, id, body)
	if err != nil {
		return err
	}

	return server.JSON(w, http.StatusOK, entity)
}

// Delete godoc
//
//	@Tags			PackingListItem
//	@Summary		Delete a PackingListItem
//	@Description	Delete a PackingListItem
//	@Accept			json
//	@Produce		json
//	@Param			id		path	string	true	"The PackingList ID"
//	@Param			item-id	path	string	true	"The PackingListItem ID"
//	@Success		204
//	@Router			/v1/packing-lists/{id}/items/{item-id} [DELETE]
//	@Security		Bearer
func (uc *PackingListItemController) Delete(w http.ResponseWriter, r *http.Request) error {
	listId, id, err := extractors.ID2(r, "id", "item-id")
	if err != nil {
		return err
	}

	err = uc.service.Delete(r.Context(), listId, id)
	if err != nil {
		return err
	}

	return server.JSON(w, http.StatusNoContent, nil)
}
