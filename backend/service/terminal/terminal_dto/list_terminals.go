package terminal_dto

import pkgterminal "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/terminal"

// ListTerminalsInput identifies the chat session whose terminals should be loaded.
type ListTerminalsInput struct {
	SessionID uint `json:"session_id"`
}

// ListTerminalsOutput returns terminal metadata for a chat session.
type ListTerminalsOutput struct {
	Items []pkgterminal.Info `json:"items"`
}
