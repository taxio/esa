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

func NewEditSubCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "edit post",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			// Validation
			if len(args) > 1 {
				return fail.New("Invalid arguments")
			}

			fs := afero.NewOsFs()
			diApp, err := NewDiApp(ctx, fs)
			if err != nil {
				return fail.Wrap(err)
			}
			postSrv := diApp.PostService

			// Get Post ID
			var postId int
			if len(args) == 0 {
				// Search for post incrementally.
				pId, err := searchPostId()
				if err != nil {
					return fail.Wrap(err)
				}
				postId = pId
			} else {
				pId, err := ParsePostIdFromArg(args[0])
				if err != nil {
					return fail.Wrap(err)
				}
				postId = pId
			}

			err = postSrv.EditPost(ctx, postId)
			if err != nil {
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

func searchPostId() (int, error) {
	return 0, fail.New("Unimplemented")
}
