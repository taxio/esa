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
			ctx := cmd.Context()
			fs := afero.NewOsFs()
			diApp, err := NewDiApp(ctx, fs)
			if err != nil {
				return fail.Wrap(err)
			}

			cfgManager := diApp.ConfigManager
			if err := cfgManager.Edit(ctx); err != nil {
				return fail.Wrap(err)
			}
			return nil
		},
	}
	return cmd
}
