package main

import "time"

type (
	GhostPost struct {
		ID                 string   `json:"id"`
		UUID               string   `json:"uuid"`
		Title              string   `json:"title"`
		Slug               string   `json:"slug"`
		HTML               string   `json:"html"`
		CommentID          string   `json:"comment_id"`
		FeatureImage       string   `json:"feature_image"`
		Featured           bool     `json:"featured"`
		Page               bool     `json:"page"`
		MetaTitle          string   `json:"meta_title"`
		MetaDescription    string   `json:"meta_description"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		PublishedAt        time.Time `json:"published_at"`
		CustomExcerpt      string   `json:"custom_excerpt"`
		OGImage            string   `json:"og_image"`
		OGTitle            string   `json:"og_title"`
		OGDescription      string   `json:"og_description"`
		TwitterImage       string   `json:"twitter_image"`
		TwitterTitle       string   `json:"twitter_title"`
		TwitterDescription string   `json:"twitter_description"`
		CustomTemplate     string   `json:"custom_template"`
		URL                string   `json:"url"`
		Excerpt            string   `json:"excerpt"`
	}

	GhostPagination struct {
		Page	int `json:"page"`
		Limit	int `json:"limit"`
		Pages	int `json:"pages"`
		Total	int `json:"total"`
		Next	int `json:"next"`
		Prev	int `json:"prev"`
	}
	
	GhostMeta struct {
		Pagination GhostPagination `json:"pagination"`
	}

	GhostPostsResponse struct {
		Posts []GhostPost `json:"posts"`
		Meta  GhostMeta   `json:"meta"`
	}
)