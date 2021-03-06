package main

import "github.com/spf13/cobra"

func NewEditSubCmd(cfg *Config, client *Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "edit post",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return cmd
}
