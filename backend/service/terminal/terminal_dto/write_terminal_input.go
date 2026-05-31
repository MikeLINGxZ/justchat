package terminal_dto

// WriteTerminalInputInput carries user input destined for an active terminal.
type WriteTerminalInputInput struct {
	TerminalID string `json:"terminal_id"`
	Data       string `json:"data"`
}

// WriteTerminalInputOutput acknowledges that terminal input was accepted.
type WriteTerminalInputOutput struct{}
