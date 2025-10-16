package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

type GenerateFunc func(c *Config) error

func main() {
	app := &cli.Command{
		Name:  "codegen",
		Usage: "cli tool for admin tasks",
		Before: func(ctx context.Context, c *cli.Command) (context.Context, error) {
			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
			return ctx, nil
		},
		Commands: []*cli.Command{
			{
				Name:        "generate",
				Description: "generate code from flags",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "config",
						Aliases:  []string{"c"},
						Usage:    "config file",
						Required: true,
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					conf := &Config{}

					if err := conf.Load(c.String("config")); err != nil {
						return err
					}

					fmt.Println(conf.Dump())

					generators := map[string]GenerateFunc{
						"Typescript": GenerateTypescript,
					}

					for name, generator := range generators {
						log.Info().Str("generator", name).Msg("running generator")

						err := generator(conf)
						if err != nil {
							log.Error().Err(err).Str("generator", name).Msg("failed to run generator")
							return err
						}

						log.Info().Str("generator", name).Msg("generator ran successfully")
					}

					return nil
				},
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		log.Fatal().Err(err).Msg("failed to run cli")
	}
}
