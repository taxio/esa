package main

import (
	"fmt"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/srvc/fail/v4"
)

type SortKey int

const (
	_ SortKey = iota
	Updated
	Created
	Number
	Stars
	Watches
	Comments
	BestMatch
)

func sortKeyFromString(key string) (SortKey, error) {
	switch key {
	case "updated":
		return Updated, nil
	case "created":
		return Created, nil
	case "number":
		return Number, nil
	case "stars":
		return Stars, nil
	case "watches":
		return Watches, nil
	case "comments":
		return Comments, nil
	case "best_match":
		return BestMatch, nil
	}

	return 0, fail.New(fmt.Sprintf("Unknown SortKey: %s", key))
}

func NewListSubCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list posts",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Setup
			fs := afero.NewOsFs()
			cfg, err := LoadConfig(fs)
			if err != nil {
				return fail.Wrap(err)
			}
			client, err := NewClient(cfg.AccessToken, cfg.TeamName)
			if err != nil {
				return fail.Wrap(err)
			}

			count, err := cmd.Flags().GetInt("count")
			if err != nil {
				return fail.Wrap(err)
			}
			sortKeyName, err := cmd.Flags().GetString("sort")
			if err != nil {
				return fail.Wrap(err)
			}
			sortKey, err := sortKeyFromString(sortKeyName)
			if err != nil {
				return fail.Wrap(err)
			}

			posts, err := client.GetPosts(cmd.Context(), count, sortKey)
			if err != nil {
				return fail.Wrap(err)
			}
			for _, post := range posts {
				fmt.Printf("%s: %s\n", post.FullName, post.Url)
			}
			return nil
		},
	}

	cmd.Flags().IntP("count", "c", 20, "Only print the number of posts")
	cmd.Flags().String("sort", "updated", "Sort key [updated:created:number:stars:watches:comments:best_match]")

	return cmd
}
