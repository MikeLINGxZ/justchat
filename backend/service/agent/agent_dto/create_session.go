package agent_dto

type CreateSessionInput struct {
	Title string   `json:"title"`
	Tags  []string `json:"tags,omitempty"`
}

type CreateSessionOutput struct {
	SessionID uint     `json:"session_id"`
	Title     string   `json:"title"`
	Tags      []string `json:"tags"`
}
