package main

import (
	"github.com/spf13/cobra"
	"github.com/srvc/fail/v4"
)

func NewConfigSubCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "configure esa tool",
		RunE: func(cmd *cobra.Command, args []string) error {
			// open config by editor
			if err := execEditor(cfg.Editor, cfg.Path); err != nil {
				return fail.Wrap(err)
			}

			// reload config
			if err := cfg.Reload(); err != nil {
				return fail.Wrap(err)
			}
			return nil
		},
	}
	return cmd
}
