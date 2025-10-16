package handlers

import (
	"net/http"

	"github.com/hay-kot/hookfeed/backend/internal/data/dtos"
	"github.com/hay-kot/hookfeed/backend/internal/services"
	"github.com/hay-kot/hookfeed/backend/internal/web/extractors"
	"github.com/hay-kot/httpkit/server"
)

type PackingListController struct {
	service *services.PackingListService
}

func NewPackingListController(service *services.PackingListService) *PackingListController {
	return &PackingListController{
		service: service,
	}
}

// GetAll godoc
//
//	@Tags			PackingList
//	@Summary		List all PackingLists
//	@Description	List all PackingLists
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	dtos.PaginationResponse[dtos.PackingListSummary]	"A list of PackingLists"
//	@Param			orderBy	query		string												false	"order by"						Enums(created_at)
//	@Param			skip	query		int													false	"The number of items to skip"	default(0)
//	@Param			limit	query		int													false	"The number of items to return"	default(100)
//	@Router			/v1/packing-lists [GET]
//	@Security		Bearer
func (uc *PackingListController) GetAll(w http.ResponseWriter, r *http.Request) error {
	user := services.UserFrom(r.Context())

	page, err := extractors.QueryT[dtos.PackingListQuery](r)
	if err != nil {
		return err
	}
	page.Pagination = page.WithDefaults()

	entities, err := uc.service.GetAllByUser(r.Context(), user.ID, page)
	if err != nil {
		return err
	}

	return server.JSON(w, http.StatusOK, entities)
}

// Get godoc
//
//	@Tags			PackingList
//	@Summary		Get a PackingList
//	@Description	Get a PackingList
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"The PackingList ID"
//	@Success		200	{object}	dtos.PackingList
//	@Router			/v1/packing-lists/{id} [GET]
//	@Security		Bearer
func (uc *PackingListController) Get(w http.ResponseWriter, r *http.Request) error {
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
//	@Tags			PackingList
//	@Summary		Create a new PackingList
//	@Description	Create a new PackingList
//	@Accept			json
//	@Produce		json
//	@Param			PackingList	body		dtos.PackingListCreate	true	"The PackingList details"
//	@Success		201			{object}	dtos.PackingList
//	@Router			/v1/packing-lists [POST]
//	@Security		Bearer
func (uc *PackingListController) Create(w http.ResponseWriter, r *http.Request) error {
	user := services.UserFrom(r.Context())

	body, err := extractors.Body[dtos.PackingListCreate](r)
	if err != nil {
		return err
	}

	entity, err := uc.service.Create(r.Context(), user.ID, body)
	if err != nil {
		return err
	}

	return server.JSON(w, http.StatusCreated, entity)
}

// Update godoc
//
//	@Tags			PackingList
//	@Summary		Update a PackingList
//	@Description	Update a PackingList
//	@Accept			json
//	@Produce		json
//	@Param			PackingList	body		dtos.PackingListUpdate	true	"The PackingList details"
//	@Param			id			path		string					true	"The PackingList ID"
//	@Success		200			{object}	dtos.PackingList
//	@Router			/v1/packing-lists/{id} [PUT]
//	@Security		Bearer
func (uc *PackingListController) Update(w http.ResponseWriter, r *http.Request) error {
	user := services.UserFrom(r.Context())

	id, body, err := extractors.BodyWithID[dtos.PackingListUpdate](r, "id")
	if err != nil {
		return err
	}

	entity, err := uc.service.Update(r.Context(), user.ID, id, body)
	if err != nil {
		return err
	}

	return server.JSON(w, http.StatusOK, entity)
}

// Delete godoc
//
//	@Tags			PackingList
//	@Summary		Delete a PackingList
//	@Description	Delete a PackingList
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"The PackingList ID"
//	@Success		204
//	@Router			/v1/packing-lists/{id} [DELETE]
//	@Security		Bearer
func (uc *PackingListController) Delete(w http.ResponseWriter, r *http.Request) error {
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
