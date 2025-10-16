package handlers

import (
	"net/http"

	"github.com/hay-kot/hookfeed/backend/internal/services"
	"github.com/hay-kot/httpkit/server"
)

type FeedController struct {
	feedService *services.FeedService
}

func NewFeedController(feedService *services.FeedService) *FeedController {
	return &FeedController{
		feedService: feedService,
	}
}

// GetAll godoc
//
//	@Tags			Feeds
//	@Summary		Get all feeds
//	@Description	Get all feeds configured in the system
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}	dtos.Feed
//	@Router			/v1/feeds [GET]
//	@Security		Bearer
func (fc *FeedController) GetAll(w http.ResponseWriter, r *http.Request) error {
	if fc.feedService == nil {
		return server.JSON(w, http.StatusOK, []interface{}{})
	}
	feeds := fc.feedService.GetAllFeeds()
	return server.JSON(w, http.StatusOK, feeds)
}
