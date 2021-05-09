package main

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/srvc/fail/v4"
	"github.com/taxio/esa/log"
	"os"
)

func NewRootCmd() *cobra.Command {
	var cfg *Config
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

			// Load global config
			c, err := LoadConfig(afero.NewOsFs())
			if err != nil {
				return fail.Wrap(err)
			}
			cfg = &c
			return nil
		},
	}

	cmd.PersistentFlags().Bool("verbose", false, "print log for developers")

	//subCmds := []*cobra.Command{}
	//cmd.AddCommand(subCmds...)

	return cmd
}
