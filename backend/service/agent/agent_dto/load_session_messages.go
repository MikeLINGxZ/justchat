package agent_dto

type LoadSessionMessagesInput struct {
	SessionID uint `json:"session_id"`
	Offset    int  `json:"offset"`
	Limit     int  `json:"limit"`
}

type MessageItem struct {
	ID          uint   `json:"id"`
	SessionID   uint   `json:"session_id"`
	ParentID    *uint  `json:"parent_id"`
	Role        string `json:"role"`
	ContentType string `json:"content_type"`
	Content     string `json:"content"`
	ModelName   string `json:"model_name"`
	AgentName   string `json:"agent_name"`
	TokensIn    int    `json:"tokens_in"`
	TokensOut   int    `json:"tokens_out"`
	Extra       string `json:"extra"`
	Attachments string `json:"attachments"`
	Created     string `json:"created"`
}

type LoadSessionMessagesOutput struct {
	Messages []MessageItem `json:"messages"`
	Total    int64         `json:"total"`
	HasMore  bool          `json:"has_more"`
}
