// Package lib defines helper functions
package lib

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

// for compile flags
var (
	version = "dev"
	commit  string
	date    = "---"
)

// Config can be set via environment variables
type config struct {
	APIEndpoint string `envconfig:"API_ENDPOINT" default:"http://host.docker.internal:8080"`
	Version     string `envconfig:"VERSION" default:"dev"`
	Port        string `envconfig:"PORT" default:"9000"`
	LogLevel    string `envconfig:"LOG_LEVEL" default:"debug"`
}

// Config represents its configurations
var Config *config

func init() {
	cfg := &config{}
	envconfig.MustProcess("lottery-web", cfg)
	if len(version) > 0 && len(commit) > 0 && len(date) > 0 {
		cfg.Version = fmt.Sprintf("%s-%s (built at %s)", version, commit, date)
	}
	Config = cfg
}
