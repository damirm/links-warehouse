package postgres

import (
	"fmt"
	"time"
)

type Config struct {
	Host           string
	Port           uint
	User           string
	Schema         string
	Database       string
	Password       string
	TimeZone       string
	MigrationsPath string
	ConnectTimeout time.Duration
}

func (c *Config) ConnString() string {
	return fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%d sslmode=disable timezone=%s search_path=%s connect_timeout=%d",
		c.User, c.Password, c.Database, c.Host, c.Port, c.TimeZone, c.Schema, c.ConnectTimeout,
	)
}
