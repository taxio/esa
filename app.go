package main

import (
	"context"
	"os"
	"path"

	"github.com/spf13/afero"
	"github.com/srvc/fail/v4"
)

type DiApp struct {
	Config        *Config
	Client        *Client
	PostService   *PostService
	ConfigManager *ConfigManager
	Editor        Editor
}

func NewDiApp(ctx context.Context, fs afero.Fs) (*DiApp, error) {
	configDirPath, err := getConfigDirPath()
	if err != nil {
		return nil, fail.Wrap(err)
	}
	cacheDirPath, err := getCacheDirPath()
	if err != nil {
		return nil, fail.Wrap(err)
	}

	defaultEditor := os.Getenv("EDITOR")
	if defaultEditor == "" {
		defaultEditor = "vim"
	}
	configManager := NewConfigManager(fs, configDirPath, NewEditor(defaultEditor))
	cfg, err := configManager.Load(ctx)
	if err != nil {
		return nil, fail.Wrap(err)
	}
	editor := NewEditor(cfg.Editor)

	client, err := NewClient(cfg.AccessToken, cfg.TeamName)
	if err != nil {
		return nil, fail.Wrap(err)
	}

	postSrv := NewPostService(fs, client, cacheDirPath, editor)

	return &DiApp{
		Config:        cfg,
		Client:        client,
		PostService:   postSrv,
		ConfigManager: configManager,
		Editor:        editor,
	}, nil
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
