package main

import (
	"context"
	"fmt"
	"github.com/izumin5210/hx"
	"github.com/srvc/fail/v4"
	"net/url"
	"time"
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

func (c *Client) GetPosts(ctx context.Context, page, perPage int) (*GetPostsResponse, error) {
	var res GetPostsResponse

	err := c.client.Get(
		ctx,
		hx.Path("posts"),
		hx.Query("page", fmt.Sprint(page)),
		hx.Query("per_page", fmt.Sprint(perPage)),
		hx.WhenSuccess(hx.AsJSON(&res)),
		hx.WhenFailure(hx.AsError()),
	)
	if err != nil {
		return nil, fail.Wrap(err)
	}

	return &res, nil
}

func (c *Client) GetAllPosts(ctx context.Context) (int, []*GetPostsResponsePost, error) {
	var posts []*GetPostsResponsePost
	var totalCount int

	page := 1
	for {
		res, err := c.GetPosts(ctx, page, MaxPostsPerPage)
		if err != nil {
			return 0, nil, fail.Wrap(err)
		}
		posts = append(posts, res.Posts...)
		if res.NextPage == 0 {
			totalCount = res.TotalCount
			break
		}
		page++
	}

	return totalCount, posts, nil
}
