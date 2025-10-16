// Package services contains the main business logic of the application
package services

import (
	"github.com/hay-kot/hookfeed/backend/internal/core/tasks"
	"github.com/hay-kot/hookfeed/backend/internal/data/db"
	"github.com/rs/zerolog"
)

type Config struct {
	CompanyName string `json:"company_name" conf:"default:Gottl Inc."`
	WebURL      string `json:"web_url"      conf:"default:http://localhost:8080"`
}

// Service is a collection of all services in the application
type Service struct {
	Admin            *AdminService
	Users            *UserService
	Passwords        *PasswordService
	PackingLists     *PackingListService
	PackingListItems *PackingListItemService
	// $scaffold_inject_service
}

func NewService(
	cfg Config,
	l zerolog.Logger,
	db *db.QueriesExt,
	queue tasks.Queue,
) *Service {
	return &Service{
		Admin:            NewAdminService(l, db),
		Users:            NewUserService(l, db),
		Passwords:        NewPasswordService(cfg, l, db, queue),
		PackingLists:     NewPackingListService(l, db),
		PackingListItems: NewPackingListItemService(l, db),
		// $scaffold_inject_constructor
	}
}
