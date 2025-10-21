package handlers

import (
	"net/http"

	"github.com/hay-kot/hookfeed/backend/internal/data/dtos"
	"github.com/hay-kot/httpkit/server"
)

// Info godoc
//
//	@Summary		Get the status of the service
//	@Description	Get the status of the service
//	@Tags			Status
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	dtos.StatusResponse
//	@Router			/v1/info [GET]
func Info(resp dtos.StatusResponse) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		return server.JSON(w, http.StatusOK, resp)
	}
}
