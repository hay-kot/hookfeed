package db

import "fmt"

type Config struct {
	Database     string `toml:"database"       env:"POSTGRES_DATABASE"       envDefault:"hookfeed"`
	Username     string `toml:"user"           env:"POSTGRES_USERNAME"       envDefault:"postgres"`
	Password     string `toml:"password"       env:"POSTGRES_PASSWORD"       envDefault:"postgres"`
	Host         string `toml:"host"           env:"POSTGRES_HOST"           envDefault:"127.0.0.1"`
	Port         string `toml:"port"           env:"POSTGRES_PORT"           envDefault:"8861"`
	MaxIdleConns int    `toml:"max_idle_conns" env:"POSTGRES_MAX_IDLE_CONNS" envDefault:"2"`
	MaxOpenConns int    `toml:"max_open_conns" env:"POSTGRES_MAX_OPEN_CONNS" envDefault:"0"`
	SSLMode      string `toml:"ssl_mode"       env:"POSTGRES_SSL_MODE"       envDefault:"disable"`
}

func (d Config) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.Username, d.Password,
		d.Host, d.Port,
		d.Database,
		d.SSLMode,
	)
}
