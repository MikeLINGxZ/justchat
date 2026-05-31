package plugin_dto

// LoginCliInput identifies the CLI plugin to start a login session for.
type LoginCliInput struct {
	ID string `json:"id"`
}

// LoginCliOutput is empty; actual output is delivered via "cli.login.output" and "cli.login.done" events.
type LoginCliOutput struct{}

// SendLoginStdinInput sends raw bytes to a running login session's stdin.
type SendLoginStdinInput struct {
	ID   string `json:"id"`
	Data string `json:"data"`
}

// CancelLoginCliInput terminates a running login session.
type CancelLoginCliInput struct {
	ID string `json:"id"`
}

// ResizeLoginCliInput resizes the PTY of a running login session.
type ResizeLoginCliInput struct {
	ID   string `json:"id"`
	Rows uint16 `json:"rows"`
	Cols uint16 `json:"cols"`
}
