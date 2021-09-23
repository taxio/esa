package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/srvc/fail/v4"
)

func NewEditLatestSubCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "latest",
		Short: "edit latest post",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			fs := afero.NewOsFs()
			diApp, err := NewDiApp(ctx, fs)
			if err != nil {
				return fail.Wrap(err)
			}
			client := diApp.Client

			posts, err := client.GetPosts(cmd.Context(), 1, Updated)
			if err != nil {
				return fail.Wrap(err)
			}
			if len(posts) == 0 {
				return fail.New("no posts")
			}
			postId := posts[0].Number

			postSrv := diApp.PostService
			if err := postSrv.EditPost(ctx, postId); err != nil {
				if errors.Is(err, ErrPostCacheAlreadyExists) {
					scanner := bufio.NewScanner(os.Stdin)
					fmt.Printf("The file being edited exists. %s\n", postSrv.CacheDirPath(ctx, postId))
					fmt.Printf("Do you want to edit thid file? (y/n):")
					scanner.Scan()
					ans := scanner.Text()
					if ans == "y" {
						err = postSrv.EditPostFromCache(ctx, postId)
						if err != nil {
							return fail.Wrap(err)
						}
					} else if ans == "n" {
						if err := postSrv.DeletePostCache(ctx, postId); err != nil {
							return fail.Wrap(err)
						}
						if err := postSrv.EditPost(ctx, postId); err != nil {
							return fail.Wrap(err)
						}
					} else {
						return fail.New("please input y or n.")
					}
				} else {
					return fail.Wrap(err)
				}
			}
			return nil
		},
	}
	return cmd
}
