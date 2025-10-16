package feeds

// Cache is a readonly cache
type Cache struct {
	allFeeds   []FeedParsed          // stored copy of the original feeds to ensure consistent ordering
	cacheByID  map[string]FeedParsed // id => Feed
	cacheByKey map[string]string     // key => id
}

func NewCache(config *Config) *Cache {
	cache := &Cache{
		allFeeds:   make([]FeedParsed, 0, len(config.Feeds)),
		cacheByID:  make(map[string]FeedParsed),
		cacheByKey: make(map[string]string),
	}

	for _, feed := range config.Feeds {
		parsed := feed.IntoParsed()
		cache.allFeeds = append(cache.allFeeds, parsed)
		cache.cacheByID[parsed.ID] = parsed

		for _, key := range parsed.Keys {
			cache.cacheByKey[key] = parsed.ID
		}
	}

	return cache
}

func (c *Cache) GetByKey(key string) (ok bool, feed FeedParsed) {
	id, exists := c.cacheByKey[key]
	if !exists {
		return false, FeedParsed{}
	}
	return c.GetByID(id)
}

func (c *Cache) GetByID(id string) (ok bool, feed FeedParsed) {
	feed, exists := c.cacheByID[id]
	return exists, feed
}

func (c *Cache) GetAll() []FeedParsed {
	return c.allFeeds
}
