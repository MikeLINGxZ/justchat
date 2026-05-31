package agent_dto

type RenameSessionInput struct {
	SessionID uint   `json:"session_id"`
	Title     string `json:"title"`
}

type RenameSessionOutput struct {
}
