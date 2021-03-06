package main

import (
	"github.com/spf13/cobra"
	"github.com/taxio/esa/api"
)

func NewRootCmd(client *api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "esa",
		Short: "A cli tool for esa",
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
