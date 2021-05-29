package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/srvc/fail/v4"
	"github.com/taxio/esa/log"
)

func NewEditSubCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "edit post",
		RunE: func(cmd *cobra.Command, args []string) error {
			// Validation
			if len(args) > 1 {
				return fail.New("Invalid arguments")
			}

			fs := afero.NewOsFs()
			cfg, err := LoadConfig(fs)
			if err != nil {
				return fail.Wrap(err)
			}
			client, err := NewClient(cfg.AccessToken, cfg.TeamName)
			if err != nil {
				return fail.Wrap(err)
			}

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

			// Get Post Detail
			post, err := client.GetPost(cmd.Context(), postId)
			if err != nil {
				return fail.Wrap(err)
			}
			log.Printf("%#v\n", post)

			// write post data to temporary file
			cacheDirPath, err := savePostTemporary(fs, cfg.CacheDirPath, post)
			if err != nil {
				return fail.Wrap(err)
			}
			log.Printf("cache dir: %s\n", cacheDirPath)

			// open temporary file by editor
			if err := execEditor(cfg.Editor, filepath.Join(cacheDirPath, "body.md")); err != nil {
				return fail.Wrap(err)
			}

			// save post to esa.io
			if err := updatePost(cmd.Context(), cacheDirPath, fs, client); err != nil {
				return fail.Wrap(err)
			}

			// rm temporary file
			if err := fs.RemoveAll(cacheDirPath); err != nil {
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

var ErrCacheAlreadyExists = errors.New("cache already exists")

func savePostTemporary(fs afero.Fs, cacheDir string, post *Post) (string, error) {
	af := afero.Afero{Fs: fs}
	// CacheDir/posts/:post_number
	cachePath := filepath.Join(cacheDir, "posts", fmt.Sprintf("%d", post.Number))

	// Check for the existence of a chache directory.
	ok, err := af.DirExists(cachePath)
	if err != nil {
		return "", fail.Wrap(err)
	}
	if ok {
		return "", fail.Wrap(ErrCacheAlreadyExists, fail.WithMessage(cachePath))
	}

	// Create a cache directory.
	if err := af.MkdirAll(cachePath, 0755); err != nil {
		return "", fail.Wrap(err)
	}

	// Save meta of post
	jsonBytes, err := json.Marshal(post)
	if err != nil {
		return "", fail.Wrap(err)
	}
	if err := af.WriteFile(filepath.Join(cachePath, "post.json"), jsonBytes, 0644); err != nil {
		return "", fail.Wrap(err)
	}

	// Save body markdown data
	if err := af.WriteFile(filepath.Join(cachePath, "body.md"), []byte(post.OriginalRevision.BodyMd), 0644); err != nil {
		return "", fail.Wrap(err)
	}

	return cachePath, nil
}

func updatePost(ctx context.Context, cacheDirPath string, fs afero.Fs, client *Client) error {
	af := afero.Afero{Fs: fs}
	var post Post
	// load meta data from post.json
	postBytes, err := af.ReadFile(filepath.Join(cacheDirPath, "post.json"))
	if err != nil {
		return fail.Wrap(err)
	}
	if err := json.Unmarshal(postBytes, &post); err != nil {
		return fail.Wrap(err)
	}

	// load body md from body.md
	bodyBytes, err := af.ReadFile(filepath.Join(cacheDirPath, "body.md"))
	if err != nil {
		return fail.Wrap(err)
	}
	post.BodyMd = string(bodyBytes)

	// update post
	if err := client.UpdatePost(ctx, &post); err != nil {
		return fail.Wrap(err)
	}

	return nil
}
