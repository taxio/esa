package main

import (
	"errors"

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
			if !errors.Is(err, ErrPostCacheAlreadyExists) && err != nil {
				return fail.Wrap(err)
			}
			err = postSrv.EditPostFromCache(postId)
			if err != nil {
				return fail.Wrap(err)
			}

			return nil
		},
	}

	return cmd
}

func searchPostId() (int, error) {
	return 0, fail.New("Unimplemented")
}
