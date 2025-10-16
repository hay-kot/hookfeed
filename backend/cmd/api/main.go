package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hay-kot/hookfeed/backend/internal/data/db"
	"github.com/hay-kot/hookfeed/backend/internal/data/db/migrations"
	"github.com/hay-kot/hookfeed/backend/internal/services"
	"github.com/hay-kot/hookfeed/backend/internal/xapps/intervalbot"
	"github.com/hay-kot/hookfeed/backend/internal/xapps/tasker"
	"github.com/hay-kot/hookfeed/backend/internal/xapps/webapi"
	"github.com/hay-kot/plugs/plugs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// Build information. Populated at build-time via -ldflags flag.
	version = "dev"
	commit  = "HEAD"
	date    = "now"
)

func build() string {
	short := commit
	if len(commit) > 7 {
		short = commit[:7]
	}

	return fmt.Sprintf("%s (%s) %s", version, short, date)
}

// @title                      hookfeed
// @version                    0.1
// @description                This is a standard Rest API template
// @BasePath                   /api
// @securityDefinitions.apikey Bearer
// @in                         header
// @name                       Authorization
// @description                "Type 'Bearer TOKEN' to correctly set the API Key"
func main() {
	err := run()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to run program")
	}
}

func run() error {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	migrations.SetLogger(log.Logger)

	cfg := Load()

	queries, err := db.NewExt(context.Background(), log.Logger, cfg.Postgres, true)
	if err != nil {
		return fmt.Errorf("failed to create database connection: %w", nil)
	}

	taskRunner := tasker.New()

	services := services.NewService(cfg.ServiceCfg, log.Logger, queries, taskRunner)

	webAPI := webapi.New(log.Logger, build(), cfg.Web, services)

	intervalBot := intervalbot.New(log.Logger)

	mgr := plugs.New(
		plugs.WithPrintln(log.Logger.Print),
	)

	mgr.AddFunc("interval_bot", intervalBot.Start)
	mgr.AddFunc("task_runner", taskRunner.Start)
	mgr.AddFunc("web_api", webAPI.Start)

	return mgr.Start(context.Background())
}
