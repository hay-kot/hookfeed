package services

type FeedService struct {
	cfg FeedFile
}

func NewFeedService(feedfile FeedFile) *FeedService {
	return &FeedService{
		cfg: feedfile,
	}
}
