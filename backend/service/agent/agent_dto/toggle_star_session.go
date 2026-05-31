package agent_dto

type ToggleStarSessionInput struct {
	SessionID uint `json:"session_id"`
	Starred   bool `json:"starred"`
}

type ToggleStarSessionOutput struct {
}
