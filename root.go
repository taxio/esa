package main

import "github.com/spf13/cobra"

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Short: "esa",
		Long:  "A cli tool for esa",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	subCmds := []*cobra.Command{
		NewListSubCmd(),
	}
	cmd.AddCommand(subCmds...)

	return cmd
}
