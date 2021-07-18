package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
	"github.com/srvc/fail/v4"

	"github.com/taxio/esa/log"
)

var ErrPostCacheAlreadyExists = errors.New("post cache already exists")

type postCachePath struct {
	BaseDir string
	Meta    string
	Body    string
}

func newPostCachePath(cacheBaseDir string, postId int) *postCachePath {
	baseDir := filepath.Join(cacheBaseDir, "posts", fmt.Sprintf("%d", postId))
	return &postCachePath{
		BaseDir: baseDir,
		Meta:    filepath.Join(baseDir, "meta.json"),
		Body:    filepath.Join(baseDir, "body.md"),
	}
}

func NewPostService(fs afero.Fs, client *Client, cacheBaseDir string, editor Editor) *PostService {
	return &PostService{
		af:           &afero.Afero{Fs: fs},
		client:       client,
		cacheBaseDir: cacheBaseDir,
		editor:       editor,
	}
}

type PostService struct {
	af           *afero.Afero
	client       *Client
	cacheBaseDir string
	editor       Editor
}

func (s *PostService) EditPost(ctx context.Context, postId int) error {
	// Get Post Detail
	post, err := s.client.GetPost(ctx, postId)
	if err != nil {
		return fail.Wrap(err)
	}
	log.Printf("%#v\n", post)

	// Write post data to cache file
	cachePath, err := s.savePostTemporary(ctx, post)
	if err != nil {
		return fail.Wrap(err)
	}

	// open cache file by editor
	if err := s.editor.Exec(ctx, cachePath.Body); err != nil {
		return fail.Wrap(err)
	}

	// save post to esa.io
	if err := s.updatePost(ctx, cachePath); err != nil {
		return fail.Wrap(err)
	}

	// remove cache file
	if err := s.af.RemoveAll(cachePath.BaseDir); err != nil {
		return fail.Wrap(err)
	}
	return nil
}

func (s *PostService) EditPostFromCache(ctx context.Context, postId int) error {
	cachePath := s.postCachePath(postId)

	// open cache file by editor
	if err := s.editor.Exec(ctx, cachePath.Body); err != nil {
		return fail.Wrap(err)
	}

	// save post to esa.io
	if err := s.updatePost(ctx, cachePath); err != nil {
		return fail.Wrap(err)
	}

	// remove cache file
	if err := s.af.RemoveAll(cachePath.BaseDir); err != nil {
		return fail.Wrap(err)
	}
	return nil
}

func (s *PostService) DeletePostCache(ctx context.Context, postId int) error {
	cachePath := s.postCachePath(postId)
	if err := s.af.RemoveAll(cachePath.BaseDir); err != nil {
		return fail.Wrap(err)
	}
	return nil
}

func (s *PostService) CacheDirPath(ctx context.Context, postId int) string {
	cachePath := s.postCachePath(postId)
	return cachePath.BaseDir
}

func (s *PostService) savePostTemporary(ctx context.Context, post *Post) (*postCachePath, error) {
	cachePath := newPostCachePath(s.cacheBaseDir, post.Number)

	// Check for the existence of a cache directory.
	ok, err := s.af.DirExists(cachePath.BaseDir)
	if err != nil {
		return nil, fail.Wrap(err)
	}
	if ok {
		return nil, fail.Wrap(ErrPostCacheAlreadyExists, fail.WithMessage(cachePath.BaseDir))
	}

	// Create a cache directory.
	if err := s.af.MkdirAll(cachePath.BaseDir, 0755); err != nil {
		return nil, fail.Wrap(err)
	}

	// Save meta of post
	jsonBytes, err := json.Marshal(post)
	if err != nil {
		return nil, fail.Wrap(err)
	}
	if err := s.af.WriteFile(cachePath.Meta, jsonBytes, 0644); err != nil {
		return nil, fail.Wrap(err)
	}

	// Save body markdown data
	if err := s.af.WriteFile(cachePath.Body, []byte(post.OriginalRevision.BodyMd), 0644); err != nil {
		return nil, fail.Wrap(err)
	}

	return cachePath, nil
}

func (s *PostService) updatePost(ctx context.Context, cachePath *postCachePath) error {
	var post Post

	metaBytes, err := s.af.ReadFile(cachePath.Meta)
	if err != nil {
		return fail.Wrap(err)
	}
	if err := json.Unmarshal(metaBytes, &post); err != nil {
		return fail.Wrap(err)
	}

	bodyBytes, err := s.af.ReadFile(cachePath.Body)
	if err != nil {
		return fail.Wrap(err)
	}
	post.BodyMd = string(bodyBytes)

	if err := s.client.UpdatePost(ctx, &post); err != nil {
		return fail.Wrap(err)
	}

	return nil
}

func (s *PostService) postCachePath(postId int) *postCachePath {
	return newPostCachePath(s.cacheBaseDir, postId)
}
