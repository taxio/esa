package main

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/srvc/fail/v4"
	"github.com/taxio/esa/log"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     AppName,
		Short:   "A cli tool for esa",
		Version: Version,
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
		NewConfigSubCmd(),
		NewListSubCmd(),
		NewEditSubCmd(),
		NewNewSubCmd(),
	}
	cmd.AddCommand(subCmds...)

	return cmd
}
