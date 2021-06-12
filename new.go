package main

import (
	"fmt"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/srvc/fail/v4"
)

func NewNewSubCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new",
		Short: "create new post from template",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validation
			if len(args) > 1 {
				return fail.New("Invalid arguments")
			}

			ctx := cmd.Context()
			fs := afero.NewOsFs()
			diApp, err := NewDiApp(ctx, fs)
			if err != nil {
				return fail.Wrap(err)
			}
			client := diApp.Client

			template, err := cmd.Flags().GetString("template")
			if err != nil {
				return fail.Wrap(err)
			}

			if template != "" {
				templatePostId, err := ParsePostIdFromArg(template)
				if err != nil {
					return fail.Wrap(err)
				}

				post, err := client.CreatePostFromTemplate(ctx, templatePostId)
				if err != nil {
					return fail.Wrap(err)
				}

				fmt.Printf("Created: %s\n", post.FullName)
				fmt.Println(post.Url)
			} else {
				panic("Unimplemented")
			}

			return nil
		},
	}

	cmd.Flags().StringP("template", "t", "", "the id of template post")

	return cmd
}
