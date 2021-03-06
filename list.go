package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/srvc/fail/v4"
	"github.com/taxio/esa/api"
)

func NewListSubCmd(client *api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list posts",
		RunE: func(cmd *cobra.Command, args []string) error {
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
