package db

import (
	"fmt"
	"strings"
)

type Config struct {
	Host       string
	Port       int
	Username   string
	Password   string
	Database   string
	Schema     string
	Insecure   bool
	Migrations bool
}

func (cfg *Config) DSN() string {
	conn := fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)
	var params []string
	if cfg.Insecure {
		params = append(params, "sslmode=disable")
	}
	if cfg.Schema != "" {
		params = append(params, "search_path="+cfg.Schema)
	}
	if len(params) > 0 {
		conn = conn + "?" + strings.Join(params, "&")
	}
	return conn
}
