package main

import "time"

type Post struct {
	Number           int          `json:"number,omitempty"`
	Name             string       `json:"name,omitempty"`
	FullName         string       `json:"full_name,omitempty"`
	BodyMd           string       `json:"body_md,omitempty"`
	CreatedAt        time.Time    `json:"created_at"`
	UpdatedAt        time.Time    `json:"updated_at"`
	Message          string       `json:"message,omitempty"`
	Url              string       `json:"url,omitempty"`
	Tags             []string     `json:"tags,omitempty"`
	Category         string       `json:"category,omitempty"`
	OriginalRevision PostRevision `json:"original_revision,omitempty"`
}

type PostRevision struct {
	BodyMd string `json:"body_md"`
	Number int    `json:"number"`
	User   string `json:"user"`
}
