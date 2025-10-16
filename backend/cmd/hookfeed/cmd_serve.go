package main

import (
	"context"
	"fmt"

	"github.com/hay-kot/hookfeed/backend/internal/data/db"
	"github.com/hay-kot/hookfeed/backend/internal/data/db/migrations"
	"github.com/hay-kot/hookfeed/backend/internal/services"
	"github.com/hay-kot/hookfeed/backend/internal/xapps/intervalbot"
	"github.com/hay-kot/hookfeed/backend/internal/xapps/tasker"
	"github.com/hay-kot/hookfeed/backend/internal/xapps/webapi"
	"github.com/hay-kot/plugs/plugs"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

type ServeCmd struct {
	flags struct {
		config        string
		middlewareDir string
	}
}

func NewServeCommand() *ServeCmd {
	return &ServeCmd{}
}

func (s *ServeCmd) Register(app *cli.Command) *cli.Command {
	cmd := &cli.Command{
		Name:      "serve",
		Usage:     "Start the HookFeed server",
		UsageText: "hookfeed serve [options]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Path to configuration file",
				Destination: &s.flags.config,
			},
			&cli.StringFlag{
				Name:        "middleware-dir",
				Aliases:     []string{"m"},
				Usage:       "Path to directory containing middleware Lua scripts",
				Destination: &s.flags.middlewareDir,
			},
		},
		Action: s.serve,
	}

	app.Commands = append(app.Commands, cmd)
	return app
}

func (s *ServeCmd) serve(ctx context.Context, cmd *cli.Command) error {
	log.Info().Msg("starting HookFeed server")

	// Set migration logger
	migrations.SetLogger(log.Logger)

	// Load configuration from environment
	cfg := LoadConfig()

	// Override feed file path from command line flag if provided
	if s.flags.config != "" {
		cfg.ServiceCfg.FeedFile = s.flags.config
	}

	log.Info().
		Str("host", cfg.Web.Host).
		Str("port", cfg.Web.Port).
		Str("database", cfg.Postgres.Host).
		Str("feedFile", cfg.ServiceCfg.FeedFile).
		Msg("configuration loaded")

	// Initialize database connection with migrations
	queries, err := db.NewExt(ctx, log.Logger, cfg.Postgres, true)
	if err != nil {
		return fmt.Errorf("failed to create database connection: %w", err)
	}

	// Initialize task runner
	taskRunner := tasker.New()

	// Initialize services
	svcs, err := services.NewService(cfg.ServiceCfg, log.Logger, queries, taskRunner)
	if err != nil {
		return fmt.Errorf("failed to initialize services: %w", err)
	}

	// Initialize web API
	webAPI := webapi.New(log.Logger, build(), cfg.Web, svcs)

	// Initialize interval bot for scheduled tasks
	intervalBot := intervalbot.New(log.Logger)

	// Create plugs manager to orchestrate all services
	mgr := plugs.New(
		plugs.WithPrintln(log.Logger.Print),
	)

	// Register all services
	mgr.AddFunc("interval_bot", intervalBot.Start)
	mgr.AddFunc("task_runner", taskRunner.Start)
	mgr.AddFunc("web_api", webAPI.Start)

	log.Info().Msg("starting all services")

	// Start all services and block until context is cancelled
	return mgr.Start(ctx)
}
