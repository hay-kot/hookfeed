package main

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

// Config represents the complete HookFeed configuration
type Config struct {
	Version    int      `yaml:"version"`
	Middleware []string `yaml:"middleware"` // Filenames in execution order
	Feeds      []Feed   `yaml:"feeds"`
	Database   Database `yaml:"database"`
	Server     Server   `yaml:"server"`
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

// Database configuration
type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	SSLMode  string `yaml:"sslMode"`
}

// Server configuration
type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// Load reads and parses a YAML configuration file
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Version != 1 {
		return fmt.Errorf("unsupported config version: %d (expected 1)", c.Version)
	}

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

// ConnectionString returns a PostgreSQL connection string
func (d *Database) ConnectionString() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		d.Host, d.Port, d.User, d.Password, d.Database, d.SSLMode,
	)
}
