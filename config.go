package main

import (
	"os"
	"path"

	"github.com/kelseyhightower/envconfig"
	"github.com/srvc/fail/v4"
)

const (
	AppName = "esa"
	Version = "0.0.1"
)

type Config struct {
	AccessToken  string `envconfig:"esa_access_token" required:"true"`
	TeamName     string `envconfig:"esa_team_name" required:"true"`
	CacheDirPath string
}

func LoadConfig() (*Config, error) {
	cfg := &Config{}

	// load environmental values
	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, fail.Wrap(err)
	}

	cacheHomeDir, err := os.UserCacheDir()
	if err != nil {
		return nil, fail.Wrap(err)
	}
	cfg.CacheDirPath = path.Join(cacheHomeDir, AppName)

	return cfg, nil
}
