package main

import (
	"github.com/spf13/cobra"
)

func NewRootCmd(cfg *Config, client *Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "esa",
		Short: "A cli tool for esa",
		Version: cfg.Version,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	subCmds := []*cobra.Command{
		NewListSubCmd(client),
	}
	cmd.AddCommand(subCmds...)

	return cmd
}
