package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"text/template"

	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/srvc/fail/v4"

	"github.com/taxio/esa/log"
)

type ConfigManager struct {
	af       *afero.Afero
	dirPath  string
	filePath string
	editor   Editor
}

func NewConfigManager(fs afero.Fs, dirPath string, editor Editor) *ConfigManager {
	return &ConfigManager{
		af:       &afero.Afero{Fs: fs},
		dirPath:  dirPath,
		filePath: filepath.Join(dirPath, fmt.Sprintf("%s.%s", ConfigName, ConfigType)),
		editor:   editor,
	}
}

func (c *ConfigManager) Load(ctx context.Context) (*Config, error) {
	log.Println("Load Config...")
	if err := c.initIfNotExists(ctx); err != nil {
		return nil, fail.Wrap(err)
	}

	var cfg Config

	viper.SetConfigName(ConfigName)
	viper.SetConfigType(ConfigType)
	viper.AddConfigPath(c.dirPath)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fail.Wrap(err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fail.Wrap(err)
	}
	cfg.AppName = appName
	cfg.Version = version

	c.editor.SetEditor(cfg.Editor)

	log.Println("Configured.")

	return &cfg, nil
}

func (c *ConfigManager) initIfNotExists(ctx context.Context) error {
	// check config existence
	if ok, err := c.af.Exists(c.filePath); ok || err != nil {
		if err != nil {
			return fail.Wrap(err)
		}
		return nil
	}

	// create config dir if not exists
	if ok, err := c.af.DirExists(c.dirPath); !ok || err != nil {
		if err != nil {
			return fail.Wrap(err)
		}
		log.Printf("mkdir %s\n", c.dirPath)
		if err := c.af.MkdirAll(c.dirPath, 0755); err != nil {
			return fail.Wrap(err)
		}
	}

	// touch
	log.Printf("touch %s\n", c.filePath)
	f, err := c.af.Create(c.filePath)
	if err != nil {
		return fail.Wrap(err)
	}
	defer mustClose(f)

	var cfg Config
	// set default value
	cfg.AccessToken = ""
	cfg.TeamName = ""
	cfg.Editor = os.Getenv("EDITOR")
	cfg.SelectCmd = "peco"

	t, err := template.New(fmt.Sprintf("%s.%s", ConfigName, ConfigType)).Parse(string(configTmpl))
	if err != nil {
		return fail.Wrap(err)
	}
	log.Printf("write data to %s\n", c.filePath)
	if err := t.Execute(io.Writer(f), c); err != nil {
		return fail.Wrap(err)
	}

	return nil
}

func (c *ConfigManager) Edit(ctx context.Context) error {
	if err := c.editor.Exec(ctx, c.filePath); err != nil {
		return fail.Wrap(err)
	}
	return nil
}
