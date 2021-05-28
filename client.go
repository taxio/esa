package main

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/izumin5210/hx"
	"github.com/srvc/fail/v4"
	"github.com/taxio/esa/log"
)

const (
	BaseUrl = "https://api.esa.io/v1"
)

type Client struct {
	client *hx.Client
}

func NewClient(token, team string) (*Client, error) {
	baseUrl, err := url.Parse(hx.Path(BaseUrl, "teams", team) + "/")
	if err != nil {
		return nil, fail.Wrap(err)
	}

	return &Client{
		client: hx.NewClient(
			hx.BaseURL(baseUrl),
			hx.Bearer(token),
		),
	}, nil
}

type GetPostsResponsePost struct {
	Number    int       `json:"number"`
	Name      string    `json:"name"`
	FullName  string    `json:"full_name"`
	Wip       bool      `json:"wip"`
	BodyMd    string    `json:"body_md"`
	Url       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Tags      []string  `json:"tags"`
}

type GetPostsResponse struct {
	Posts      []*GetPostsResponsePost `json:"posts"`
	PrevPage   int                     `json:"prev_page"`
	NextPage   int                     `json:"next_page"`
	TotalCount int                     `json:"total_count"`
	Page       int                     `json:"page"`
	PerPage    int                     `json:"per_page"`
	MaxPerPage int                     `json:"max_per_page"`
}

const MaxPostsPerPage = 100

func (c *Client) GetPosts(ctx context.Context, max int, sortKey SortKey) ([]*Post, error) {
	log.Printf("GetPosts. max: %d, sortKey: %s\n", max, sortKey)
	var posts []*Post

	remain := max
	page := 0
	for remain > 0 {
		perPage := minInt(remain, MaxPostsPerPage)
		remain = remain - perPage

		log.Printf("Get post. page: %d, per_page: %d, sort: %s\n", page, perPage, sortKey)
		var res GetPostsResponse
		err := c.client.Get(
			ctx,
			hx.Path("posts"),
			hx.Query("page", fmt.Sprint(page)),
			hx.Query("per_page", fmt.Sprint(perPage)),
			hx.Query("sort", sortKey.String()),
			hx.WhenSuccess(hx.AsJSON(&res)),
			hx.WhenFailure(hx.AsError()),
		)
		if err != nil {
			return nil, fail.Wrap(err)
		}

		for _, p := range res.Posts {
			posts = append(posts, &Post{
				Number:    p.Number,
				Name:      p.Name,
				FullName:  p.FullName,
				BodyMd:    p.BodyMd,
				CreatedAt: p.CreatedAt,
				UpdatedAt: p.UpdatedAt,
				Url:       p.Url,
				Tags:      p.Tags,
			})
		}

		if res.NextPage == 0 {
			log.Println("there is no next page")
			break
		}
		page = res.NextPage
	}

	return posts, nil
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type UserInfo struct {
	MySelf     bool   `json:"myself"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
	Icon       string `json:"icon"`
}

type GetPostResponse struct {
	Number         int       `json:"number"`
	Name           string    `json:"name"`
	FullName       string    `json:"full_name"`
	Wip            bool      `json:"wip"`
	BodyMd         string    `json:"body_md"`
	CreatedAt      time.Time `json:"created_at"`
	Url            string    `json:"url"`
	UpdatedAt      time.Time `json:"updated_at"`
	Tags           []string  `json:"tags"`
	Category       string    `json:"category"`
	RevisionNumber int       `json:"revision_number"`
	CreatedBy      *UserInfo `json:"created_by"`
	UpdatedBy      *UserInfo `json:"updated_by"`
}

func (c *Client) GetPost(ctx context.Context, postId int) (*Post, error) {
	log.Println("Get Post Detail")
	var res GetPostResponse

	err := c.client.Get(
		ctx,
		hx.Path("posts", postId),
		hx.WhenSuccess(hx.AsJSON(&res)),
		hx.WhenFailure(hx.AsError()),
	)
	if err != nil {
		return nil, fail.Wrap(err)
	}

	post := Post{
		Number:    res.Number,
		Name:      res.Name,
		FullName:  res.FullName,
		CreatedAt: res.CreatedAt,
		UpdatedAt: res.UpdatedAt,
		Message:   "",
		Url:       res.Url,
		Tags:      res.Tags,
		Category:  res.Category,
		OriginalRevision: PostRevision{
			BodyMd: res.BodyMd,
			Number: res.RevisionNumber,
			User:   res.UpdatedBy.ScreenName,
		},
	}

	return &post, nil
}

func (c *Client) GetTemplatePosts(ctx context.Context, max int, sortKey SortKey) ([]*Post, error) {
	log.Printf("GetTemplatePosts. max: %d, sortKey: %s\n", max, sortKey)
	var posts []*Post

	remain := max
	page := 0
	for remain > 0 {
		perPage := minInt(remain, MaxPostsPerPage)
		remain = remain - perPage

		log.Printf("Get post. page: %d, per_page: %d, sort: %s\n", page, perPage, sortKey)
		var res GetPostsResponse
		err := c.client.Get(
			ctx,
			hx.Path("posts"),
			hx.Query("page", fmt.Sprint(page)),
			hx.Query("per_page", fmt.Sprint(perPage)),
			hx.Query("sort", sortKey.String()),
			hx.Query("q", "category:Templates"),
			hx.WhenSuccess(hx.AsJSON(&res)),
			hx.WhenFailure(hx.AsError()),
		)
		if err != nil {
			return nil, fail.Wrap(err)
		}

		for _, p := range res.Posts {
			posts = append(posts, &Post{
				Number:    p.Number,
				Name:      p.Name,
				FullName:  p.FullName,
				BodyMd:    p.BodyMd,
				CreatedAt: p.CreatedAt,
				UpdatedAt: p.UpdatedAt,
				Url:       p.Url,
				Tags:      p.Tags,
			})
		}

		if res.NextPage == 0 {
			log.Println("there is no next page")
			break
		}
		page = res.NextPage
	}

	return posts, nil
}

type PatchPostRequestOriginalRevision struct {
	BodyMd string `json:"body_md"`
	Number int    `json:"number"`
	User   string `json:"user"`
}

type PatchPostRequest struct {
	Name             string                            `json:"name"`
	BodyMd           string                            `json:"body_md"`
	Tags             []string                          `json:"tags"`
	Category         string                            `json:"category"`
	Wip              bool                              `json:"wip"`
	Message          string                            `json:"message"`
	OriginalRevision *PatchPostRequestOriginalRevision `json:"original_revision"`
}

func (c *Client) UpdatePost(ctx context.Context, post *Post) error {
	log.Println("Update Post")
	req := PatchPostRequest{
		Name:     post.Name,
		BodyMd:   post.BodyMd,
		Tags:     post.Tags,
		Category: post.Category,
		Wip:      false,
		Message:  post.Message,
		OriginalRevision: &PatchPostRequestOriginalRevision{
			BodyMd: post.OriginalRevision.BodyMd,
			Number: post.OriginalRevision.Number,
			User:   post.OriginalRevision.User,
		},
	}

	err := c.client.Patch(
		ctx,
		hx.Path("posts", post.Number),
		hx.JSON(req),
		hx.WhenFailure(hx.AsError()),
	)
	if err != nil {
		return fail.Wrap(err)
	}

	return nil
}
