package services

import (
	"fmt"
	"io"

	"github.com/goccy/go-yaml"
)

// Config represents the complete HookFeed configuration
type FeedFile struct {
	Middleware []string `yaml:"middleware"` // Filenames in execution order
	Feeds      []Feed   `yaml:"feeds"`
}

// Feed represents a webhook feed configuration
type Feed struct {
	Name        string     `yaml:"name"`
	Slug        string     `yaml:"slug"` // Used as the key for webhook URLs
	Description string     `yaml:"description"`
	Middleware  []string   `yaml:"middleware"` // Filenames of middleware scripts
	Adapters    *[]string  `yaml:"adapters"`   // pointer to distinguish between null, empty array, and populated array
	Retention   *Retention `yaml:"retention"`
}

func (f Feed) UseAdapters() (bool, []string) {
	shouldUse := f.Adapters != nil

	if !shouldUse {
		return false, []string{}
	}

	return true, *f.Adapters
}

func (f Feed) RetentionOrDefault() Retention {
	if f.Retention != nil {
		return *f.Retention
	}

	return Retention{
		MaxCount:   10_000,
		MaxAgeDays: 10_000,
	}
}

// Retention defines message retention policies
type Retention struct {
	MaxCount   int `yaml:"max_count"`
	MaxAgeDays int `yaml:"max_age_days"`
}

// FeedFileLoad reads and parses a YAML configuration file
func FeedFileLoad(reader io.Reader) (*FeedFile, error) {
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg FeedFile
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// Validate checks if the configuration is valid
func (c *FeedFile) Validate() error {
	if len(c.Feeds) == 0 {
		return fmt.Errorf("at least one feed must be defined")
	}

	// Validate feeds
	slugs := make(map[string]bool)
	for i, feed := range c.Feeds {
		if feed.Name == "" {
			return fmt.Errorf("feed[%d]: name is required", i)
		}
		if feed.Slug == "" {
			return fmt.Errorf("feed[%d]: slug is required", i)
		}
		if slugs[feed.Slug] {
			return fmt.Errorf("feed[%d]: duplicate slug '%s'", i, feed.Slug)
		}
		slugs[feed.Slug] = true
	}

	return nil
}
