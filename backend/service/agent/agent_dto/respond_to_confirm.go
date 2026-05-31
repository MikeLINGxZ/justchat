package agent_dto

type RespondToConfirmInput struct {
	SessionID uint   `json:"session_id"`
	Approved  bool   `json:"approved"`
	Message   string `json:"message"`
	Action    string `json:"action"`
}

type RespondToConfirmOutput struct {
}
