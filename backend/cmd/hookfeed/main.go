package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hay-kot/hookfeed/backend/internal/console"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
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

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	ctx := context.Background()

	err := run(ctx, os.Args)

	exitCode := 0
	if err != nil {
		fmt.Print(console.FatalError(err))
		exitCode = 1
	}

	os.Exit(exitCode)
}

func run(ctx context.Context, args []string) error {
	app := &cli.Command{
		Name:     "hookfeed",
		Usage:    "Transform webhooks using Lua scripts",
		Commands: nil,
		// Propigate usage errors to fatal printer
		OnUsageError: func(ctx context.Context, cmd *cli.Command, err error, isSubcommand bool) error {
			return err
		},
	}

	app = NewValidateCommand().Register(app)
	app = NewServeCommand().Register(app)

	return app.Run(ctx, args)
}
