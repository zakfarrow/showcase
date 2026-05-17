package model

import "time"

type Project struct {
	ID          int       `json:"id"`
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	TechStack   []string  `json:"tech_stack"`
	Status      string    `json:"status"`
	Features    []string  `json:"features"`
	Challenges  string    `json:"challenges"`
	Learnings   string    `json:"learnings"`
	FuturePlans string    `json:"future_plans"`
	GitHubURL   *string   `json:"github_url,omitempty"`
	LiveURL     *string   `json:"live_url,omitempty"`
	ImageURL    *string   `json:"image_url,omitempty"`
	Featured    bool      `json:"featured"`
	SortOrder   int       `json:"sort_order"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
