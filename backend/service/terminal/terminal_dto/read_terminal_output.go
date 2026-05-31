package terminal_dto

import pkgterminal "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/terminal"

// ReadTerminalOutputInput identifies the terminal and cursor to read from.
type ReadTerminalOutputInput struct {
	TerminalID string `json:"terminal_id"`
	Cursor     int64  `json:"cursor"`
}

// ReadTerminalOutputOutput returns terminal output chunks after the cursor.
type ReadTerminalOutputOutput struct {
	Chunks []pkgterminal.TerminalOutputChunk `json:"chunks"`
}
