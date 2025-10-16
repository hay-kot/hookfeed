package main

import (
	"context"
	"encoding/json"
	"net/http"
	"os"

	"github.com/hay-kot/hookfeed/backend/cmd/tools/seed/client"
	"github.com/hay-kot/hookfeed/backend/internal/data/dtos"
	"github.com/hay-kot/hookfeed/backend/pkgs/logutils"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v3"
)

type Config struct {
	Data []Data `json:"data"`
}

type Data struct {
	User  dtos.UserRegister  `json:"user"`
	Lists []dtos.PackingList `json:"lists"`
}

func main() {
	var err error
	log.Logger, err = logutils.Factory("debug", "console")
	if err != nil {
		panic(err)
	}

	log.Info().Msg("hookfeed Seeder Tool")

	app := &cli.Command{
		Name:  "seeder",
		Usage: "seeder tool for hookfeed development",
		Commands: []*cli.Command{
			{
				Name: "seed",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:        "config",
						HideDefault: false,
						Usage:       "Path to the config file",
						Required:    true,
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					f, err := os.Open(c.String("config"))
					if err != nil {
						return err
					}

					cfg := Config{}
					err = json.NewDecoder(f).Decode(&cfg)
					if err != nil {
						return err
					}

					cl := &http.Client{
						Transport: logutils.NewRoundTripper("client", log.Logger, http.DefaultTransport),
					}

					hookfeed := client.New(cl, "http://localhost:9990")

					for _, data := range cfg.Data {
						err := seed(ctx, log.Logger, hookfeed, data)
						if err != nil {
							return err
						}
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
