package main

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/srvc/fail/v4"
)

func NewConfigSubCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "configure esa tool",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := editConfigFile(afero.NewOsFs())
			return fail.Wrap(err)
		},
	}
	return cmd
}

func editConfigFile(fs afero.Fs) error {
	// Load config
	cfg, err := LoadConfig(fs)
	if err != nil {
		return fail.Wrap(err)
	}

	// open config by editor
	if err := execEditor(cfg.Editor, cfg.Path); err != nil {
		return fail.Wrap(err)
	}

	// reload config
	if err := cfg.Reload(); err != nil {
		return fail.Wrap(err)
	}
	return nil
}
