package main

import (
	"context"
	"os"

	"github.com/hay-kot/hookfeed/backend/pkgs/logutils"
	"github.com/rs/zerolog/log"

	_ "github.com/lib/pq"
	"github.com/urfave/cli/v3"
)

func main() {
	var err error
	log.Logger, err = logutils.Factory("info", "console")
	if err != nil {
		panic(err)
	}

	app := &cli.Command{
		Name:  "dbtool",
		Usage: "database migration tool",
		Commands: []*cli.Command{
			MigrateCommand(),
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal().Err(err).Msg("failed to run cli")
	}
}
