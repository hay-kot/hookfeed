// Package services contains the main business logic of the application
package services

import (
	"fmt"
	"os"

	"github.com/hay-kot/hookfeed/backend/internal/core/feeds"
	"github.com/hay-kot/hookfeed/backend/internal/core/tasks"
	"github.com/hay-kot/hookfeed/backend/internal/data/db"
	"github.com/rs/zerolog"
)

type Config struct {
	CompanyName string `json:"company_name" conf:"default:Gottl Inc."`
	WebURL      string `json:"web_url"      conf:"default:http://localhost:8080"`
	FeedFile    string `json:"feed_file"    conf:"default:"`
}

// Service is a collection of all services in the application
type Service struct {
	Admin        *AdminService
	Users        *UserService
	Passwords    *PasswordService
	Feeds        *FeedService
	Webhooks     *WebhookService
	FeedMessages *FeedMessageService
	// $scaffold_inject_service
}

func NewService(
	cfg Config,
	l zerolog.Logger,
	db *db.QueriesExt,
	queue tasks.Queue,
) (*Service, error) {
	// Load feed file if path is provided
	var feedService *FeedService
	if cfg.FeedFile != "" {
		file, err := os.Open(cfg.FeedFile)
		if err != nil {
			return nil, fmt.Errorf("failed to open feed file: %w", err)
		}
		defer func() { _ = file.Close() }()

		feedFile, err := feeds.Load(file)
		if err != nil {
			return nil, fmt.Errorf("failed to parse feed file: %w", err)
		}

		l.Info().
			Int("feeds", len(feedFile.Feeds)).
			Int("middleware", len(feedFile.Middleware)).
			Msg("loaded feed configuration")

		cache := feeds.NewCache(feedFile)

		feedService = NewFeedService(cache)
	}

	webhookService := NewWebhookService(l, feedService)

	return &Service{
		Admin:        NewAdminService(l, db),
		Users:        NewUserService(l, db),
		Passwords:    NewPasswordService(cfg, l, db, queue),
		Feeds:        feedService,
		Webhooks:     webhookService,
		FeedMessages: NewFeedMessageService(l, db),
		// $scaffold_inject_constructor
	}, nil
}
