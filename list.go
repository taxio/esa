package main

import "github.com/spf13/cobra"

func NewListSubCmd() *cobra.Command {
	cmd := &cobra.Command{
		Short: "list",
		Long:  "Show posts",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	return cmd
}
