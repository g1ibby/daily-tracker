// Package config provides configurations from environment variables.
package config

import (
	"github.com/powerman/appcfg"
)

// All configurable values of the service.
//
// If microservice may runs in different ways (e.g. using CLI subcommands)
// then these subcommands may use subset of these values.
var all = &struct { //nolint:gochecknoglobals // Config is global anyway.
	ApiSecret   appcfg.String         `env:"API_SECRET"`
	SheetSecret appcfg.String         `env:"SHEET_SECRET"`
	SheetID     appcfg.String         `env:"SHEET_ID"`
	TgBotToken  appcfg.NotEmptyString `env:"TG_BOT_TOKEN"`
	TgUserID    appcfg.Int64          `env:"TG_USER_ID"`
	Domain      appcfg.NotEmptyString `env:"DOMAIN"`
}{ // Defaults, if any:
	ApiSecret:   appcfg.MustString(""),
	SheetSecret: appcfg.MustString(""),
	SheetID:     appcfg.MustString(""),
	TgUserID:    appcfg.MustInt64("0"),
}

// ServeConfig contains configuration for subcommand.
type AppConfig struct {
	ApiSecret   string
	SheetSecret string
	SheetID     string
	TgBotToken  string
	TgUserID    int64
	Domain      string
}

// GetApp validates and returns configuration for subcommand.
func GetApp() (c *AppConfig, err error) {
	fromEnv := appcfg.NewFromEnv("")
	err = appcfg.ProvideStruct(all, fromEnv)
	if err != nil {
		return nil, err
	}

	c = &AppConfig{
		ApiSecret:   all.ApiSecret.Value(&err),
		SheetSecret: all.SheetSecret.Value(&err),
		SheetID:     all.SheetID.Value(&err),
		TgBotToken:  all.TgBotToken.Value(&err),
		TgUserID:    all.TgUserID.Value(&err),
		Domain:      all.Domain.Value(&err),
	}
	if err != nil {
		return nil, err
	}
	return c, nil
}
