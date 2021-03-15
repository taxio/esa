package main

import (
	"github.com/spf13/cobra"
	"github.com/srvc/fail/v4"
	"github.com/taxio/esa/log"
	"os"
)

func NewRootCmd(cfg *Config, client *Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "esa",
		Short:   "A cli tool for esa",
		Version: cfg.Version,
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			verbose, err := cmd.Flags().GetBool("verbose")
			if err != nil {
				return fail.Wrap(err)
			}
			if verbose {
				log.SetVerboseLogger(os.Stdout)
				log.Println("verbose on")
			}
			return nil
		},
	}

	cmd.PersistentFlags().Bool("verbose", false, "print log for developers")

	subCmds := []*cobra.Command{
		NewListSubCmd(client),
		NewEditSubCmd(cfg, client),
		NewConfigSubCmd(cfg),
	}
	cmd.AddCommand(subCmds...)

	return cmd
}
