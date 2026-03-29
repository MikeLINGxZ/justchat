package view_models

import "time"

type PromptFileSummary struct {
	Name        string     `json:"name"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	IsSystem    bool       `json:"is_system"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type PromptFileDetail struct {
	Name        string     `json:"name"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	IsSystem    bool       `json:"is_system"`
	Content     string     `json:"content"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
