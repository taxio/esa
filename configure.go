package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/srvc/fail/v4"
	"os"
	"os/exec"
)

func NewConfigSubCmd(cfg *Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "configure esa tool",
		RunE: func(cmd *cobra.Command, args []string) error {
			// open config by editor
			c := exec.Command("sh", "-c", fmt.Sprintf("%s %s", cfg.Editor, cfg.Path))
			c.Stderr = os.Stderr
			c.Stdout = os.Stdout
			c.Stdin = os.Stdin
			if err := c.Run(); err != nil {
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
