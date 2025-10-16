package main

import (
	"context"

	"github.com/hay-kot/hookfeed/backend/internal/data/db/migrations"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

func MigrateCommand() *cli.Command {
	migrations.SetLogger(log.Logger)

	return &cli.Command{
		Name:  "migrate",
		Usage: "run migraton related functions",
		Commands: []*cli.Command{
			{
				Name:  "up",
				Usage: "runs all migrations in order",
				Action: func(ctx context.Context, c *cli.Command) error {
					return migrations.Up(log.Logger, nil)
				},
			},
			{
				Name:  "down",
				Usage: "rolls back a single migration",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "all",
						Usage: "when specified rolls back all migrations",
						Value: false,
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					return migrations.DownByOne(ctx, log.Logger, nil)
				},
			},
		},
	}
}
