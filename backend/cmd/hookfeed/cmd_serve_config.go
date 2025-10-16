package main

import (
	"fmt"

	"github.com/caarlos0/env/v10"
	"github.com/hay-kot/hookfeed/backend/internal/data/db"
	"github.com/hay-kot/hookfeed/backend/internal/services"
	"github.com/hay-kot/hookfeed/backend/internal/xapps/webapi"
)

type Config struct {
	Web        webapi.Config
	Postgres   db.Config
	ServiceCfg services.Config
	FeedFile   string `json:"feed_file" conf:"default:configs/feeds.yml"`
}

var EnvPrefix = "HF_"

func LoadConfig() Config {
	cfg := Config{}

	err := env.ParseWithOptions(&cfg, env.Options{Prefix: EnvPrefix})
	if err != nil {
		panic(fmt.Sprintf("error loading config: %s", err))
	}

	return cfg
}
