package main

import (
	_ "embed"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"text/template"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/srvc/fail/v4"
)

//go:embed config.tmpl
var configTmpl []byte

const (
	AppName = "esa"
	Version = "0.0.1"

	ConfigName = "config"
	ConfigType = "toml"
)

type Config struct {
	af   *afero.Afero
	Path string

	AppName       string
	Version       string
	ConfigDirPath string
	CacheDirPath  string

	AccessToken string
	TeamName    string
	Editor      string
	SelectCmd   string
}

func LoadConfig(fs afero.Fs) (*Config, error) {
	var cfg Config

	cfg.af = &afero.Afero{Fs: fs}
	cfg.AppName = AppName
	cfg.Version = Version

	cacheDirPath, err := getCacheDirPath()
	if err != nil {
		return nil, fail.Wrap(err)
	}
	cfg.CacheDirPath = cacheDirPath

	configDirPath, err := getConfigDirPath()
	if err != nil {
		return nil, fail.Wrap(err)
	}
	cfg.ConfigDirPath = configDirPath

	// check config file existence
	cfg.Path = path.Join(configDirPath, fmt.Sprintf("%s.%s", ConfigName, ConfigType))
	if err := cfg.initIfNotExists(); err != nil {
		return nil, fail.Wrap(err)
	}

	if err := cfg.Reload(); err != nil {
		return nil, fail.Wrap(err)
	}

	return &cfg, nil
}

func (c *Config) initIfNotExists() error {
	// check cache existence
	if ok, err := c.af.DirExists(c.CacheDirPath); !ok || err != nil {
		if err != nil {
			return fail.Wrap(err)
		}
		log.Printf("mkdir %s\n", c.CacheDirPath)
		if err := c.af.MkdirAll(c.CacheDirPath, 0755); err != nil {
			return fail.Wrap(err)
		}
	}

	// check config existence
	if ok, err := c.af.Exists(c.Path); ok || err != nil {
		if err != nil {
			return fail.Wrap(err)
		}
		return nil
	}

	// create config dir if not exists
	if ok, err := c.af.DirExists(c.ConfigDirPath); !ok || err != nil {
		if err != nil {
			return fail.Wrap(err)
		}
		log.Printf("mkdir %s\n", c.ConfigDirPath)
		if err := c.af.MkdirAll(c.ConfigDirPath, 0755); err != nil {
			return fail.Wrap(err)
		}
	}

	// touch
	log.Printf("touch %s\n", c.Path)
	f, err := c.af.Create(c.Path)
	if err != nil {
		return fail.Wrap(err)
	}
	defer f.Close()

	// set default value
	c.AccessToken = ""
	c.TeamName = ""
	c.Editor = "vim"
	c.SelectCmd = "peco"

	fmt.Println(string(configTmpl))
	t, err := template.New(fmt.Sprintf("%s.%s", ConfigName, ConfigType)).Parse(string(configTmpl))
	if err != nil {
		return fail.Wrap(err)
	}
	log.Printf("write data to %s\n", c.Path)
	if err := t.Execute(io.Writer(f), c); err != nil {
		return fail.Wrap(err)
	}

	return nil
}

func getConfigDirPath() (string, error) {
	// ESA_CONFIG_DIR
	cfgDirPath := os.Getenv("ESA_CONFIG_DIR")
	if cfgDirPath != "" {
		return path.Join(cfgDirPath, AppName), nil
	}

	// XDG
	cfgDirPath = os.Getenv("XDG_CONFIG_HOME")
	if cfgDirPath != "" {
		return path.Join(cfgDirPath, AppName), nil
	}
	cfgDirPath = os.Getenv("HOME")
	if cfgDirPath != "" {
		return path.Join(cfgDirPath, ".config", AppName), nil
	}

	// Default
	cfgDirPath, err := os.UserConfigDir()
	if err != nil {
		return "", fail.Wrap(err)
	}

	return path.Join(cfgDirPath, AppName), nil
}

func getCacheDirPath() (string, error) {
	// ESA_CONFIG_DIR
	cacheDirPath := os.Getenv("ESA_CACHE_DIR")
	if cacheDirPath != "" {
		return path.Join(cacheDirPath, AppName), nil
	}

	// XDG
	cacheDirPath = os.Getenv("XDG_CACHE_HOME")
	if cacheDirPath != "" {
		return path.Join(cacheDirPath, AppName), nil
	}
	cacheDirPath = os.Getenv("HOME")
	if cacheDirPath != "" {
		return path.Join(cacheDirPath, ".cache", AppName), nil
	}

	// Default
	cacheDirPath, err := os.UserCacheDir()
	if err != nil {
		return "", fail.Wrap(err)
	}

	return path.Join(cacheDirPath, AppName), nil
}

func (c *Config) Reload() error {
	viper.SetConfigName(ConfigName)
	viper.SetConfigType(ConfigType)
	viper.AddConfigPath(c.ConfigDirPath)
	if err := viper.ReadInConfig(); err != nil {
		return fail.Wrap(err)
	}
	if err := viper.Unmarshal(c); err != nil {
		return fail.Wrap(err)
	}
	log.Printf("Load Config: %#v\n", c)
	return nil
}
