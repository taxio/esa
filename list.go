package main

import (
	"fmt"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/srvc/fail/v4"
)

func NewListSubCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list posts",
		RunE: func(cmd *cobra.Command, args []string) error {
			fs := afero.NewOsFs()
			cfg, err := LoadConfig(fs)
			if err != nil {
				return fail.Wrap(err)
			}
			client, err := NewClient(cfg.AccessToken, cfg.TeamName)
			if err != nil {
				return fail.Wrap(err)
			}

			_, posts, err := client.GetAllPosts(cmd.Context())
			if err != nil {
				return fail.Wrap(err)
			}
			for _, post := range posts {
				fmt.Printf("%s: %s\n", post.FullName, post.Url)
			}
			return nil
		},
	}
	return cmd
}
