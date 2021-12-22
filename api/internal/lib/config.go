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
	GcloudCreds     string `envconfig:"GCLOUD_CREDENTIALS" default:"creds.json"`
	ProjectID       string `envconfig:"PROJECT_ID"`
	SpannerInstance string `envconfig:"SPANNER_INSTANCE" default:"jaguer-cn-lottery"`
	SpannerDatabase string `envconfig:"SPANNER_DATABASE" default:"app"`
	SheetID         string `envconfig:"SPREAD_SHEET_ID"`
	SheetTabName    string `envconfig:"SPREAD_SHEET_TAB_NAME" default:"swags"`
	Version         string `envconfig:"VERSION" default:"dev"`
	Port            string `envconfig:"PORT" default:"8080"`
	LogLevel        string `envconfig:"LOG_LEVEL" default:"debug"`
}

// Config represents its configurations
var Config *config

func init() {
	cfg := &config{}
	envconfig.MustProcess("lottery", cfg)
	if len(version) > 0 && len(commit) > 0 && len(date) > 0 {
		cfg.Version = fmt.Sprintf("%s-%s (built at %s)", version, commit, date)
	}
	Config = cfg
}
