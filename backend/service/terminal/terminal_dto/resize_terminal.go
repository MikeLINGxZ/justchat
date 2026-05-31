package terminal_dto

// ResizeTerminalInput carries the current frontend terminal dimensions.
type ResizeTerminalInput struct {
	TerminalID string `json:"terminal_id"`
	Rows       uint16 `json:"rows"`
	Cols       uint16 `json:"cols"`
}

// ResizeTerminalOutput acknowledges that the resize request was handled.
type ResizeTerminalOutput struct{}
