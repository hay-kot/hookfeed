package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"

	"github.com/hay-kot/plugs/plugs"
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
		UsageText: "hookfeed serve --config <config.yml> --middleware-dir <dir>",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Path to configuration file",
				Required:    true,
				Destination: &s.flags.config,
			},
			&cli.StringFlag{
				Name:        "middleware-dir",
				Aliases:     []string{"m"},
				Usage:       "Path to directory containing middleware Lua scripts",
				Required:    true,
				Destination: &s.flags.middlewareDir,
			},
		},
		Action: s.serve,
	}

	app.Commands = append(app.Commands, cmd)
	return app
}

func (s *ServeCmd) serve(ctx context.Context, cmd *cli.Command) error {
	log.Info().
		Str("config", s.flags.config).
		Str("middlewareDir", s.flags.middlewareDir).
		Msg("starting HookFeed server")

	// TODO: Implement server startup
	fmt.Println("Server command not yet implemented")

	runner := plugs.New(
		plugs.WithTimeout(10*time.Second),
		plugs.WithPrintln(func(a ...any) {
			log.Info().Msg(fmt.Sprint(a...))
		}),
	)

	return runner.Start(ctx)
}
