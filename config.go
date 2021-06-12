package main

import (
	_ "embed"
)

//go:embed config.tmpl
var configTmpl []byte

const (
	ConfigName = "config"
	ConfigType = "toml"
)

type Config struct {
	AppName string
	Version string

	AccessToken string `mapstructure:"access_token"`
	TeamName    string `mapstructure:"team_name"`
	Editor      string `mapstructure:"editor"`
	SelectCmd   string `mapstructure:"select_cmd"`
}
