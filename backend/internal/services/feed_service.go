package services

import (
	"github.com/hay-kot/hookfeed/backend/internal/core/feeds"
	"github.com/hay-kot/hookfeed/backend/internal/data/dtos"
	"github.com/hay-kot/hookfeed/backend/pkgs/utils"
)

type FeedService struct {
	cache *feeds.Cache // exported for webhook service access
}

func (f *FeedService) GetCache() *feeds.Cache {
	return f.cache
}

func NewFeedService(cache *feeds.Cache) *FeedService {
	return &FeedService{
		cache: cache,
	}
}

func (f *FeedService) GetByKey(key string) (dtos.Feed, bool) {
	ok, feed := f.cache.GetByKey(key)
	if !ok {
		return dtos.Feed{}, false
	}

	return dtos.Feed{
		Name:        feed.Name,
		ID:          feed.ID,
		Keys:        feed.Keys,
		Description: feed.Description,
		Middleware:  feed.Middleware,
		Adapters:    feed.Adapters,
		Retention: dtos.Retention{
			MaxCount:   feed.Retention.MaxCount,
			MaxAgeDays: feed.Retention.MaxAgeDays,
		},
	}, true
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
