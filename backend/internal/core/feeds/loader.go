package feeds

import (
	"fmt"
	"io"

	"github.com/goccy/go-yaml"
)

func Load(reader io.Reader) (*Config, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	/* if err := cfg.Validate(); err != nil { */
	/* 	return nil, fmt.Errorf("invalid configuration: %w", err) */
	/* } */

	return &cfg, nil
}
