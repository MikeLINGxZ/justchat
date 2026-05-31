package agent_dto

type ListSessionsInput struct {
	Cursor        uint `json:"cursor"`
	Limit         int  `json:"limit"`
	StarredOnly   bool `json:"starred_only"`
	IncludeHidden bool `json:"include_hidden"`
}

type SessionItem struct {
	ID      uint     `json:"id"`
	Title   string   `json:"title"`
	Kind    string   `json:"kind"`
	Tags    []string `json:"tags"`
	Starred bool     `json:"starred"`
	Status  string   `json:"status"`
	Created string   `json:"created"`
	Updated string   `json:"updated"`
}

type ListSessionsOutput struct {
	Sessions   []SessionItem `json:"sessions"`
	NextCursor uint          `json:"next_cursor"`
	HasMore    bool          `json:"has_more"`
}
