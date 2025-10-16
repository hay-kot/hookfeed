package services

import (
	"github.com/hay-kot/hookfeed/backend/internal/core/feeds"
	"github.com/hay-kot/hookfeed/backend/internal/data/dtos"
	"github.com/hay-kot/hookfeed/backend/pkgs/utils"
)

type FeedService struct {
	cache *feeds.Cache
}

func NewFeedService(cache *feeds.Cache) *FeedService {
	return &FeedService{
		cache: cache,
	}
}

func (f *FeedService) GetAllFeeds() []dtos.Feed {
	parsed := f.cache.GetAll()

	return utils.Map(parsed, func(f feeds.FeedParsed) dtos.Feed {
		return dtos.Feed{
			Name:        f.Name,
			ID:          f.ID,
			Keys:        f.Keys,
			Description: f.Description,
			Middleware:  f.Middleware,
			Adapters:    f.Adapters,
			Retention: dtos.Retention{
				MaxCount:   f.Retention.MaxCount,
				MaxAgeDays: f.Retention.MaxAgeDays,
			},
		}
	})
}
