package plugin_dto

// ResetCliDataInput identifies which CLI plugin's data directory to clear.
type ResetCliDataInput struct {
	ID string `json:"id"`
}

// ResetCliDataOutput is empty; clients re-fetch the extension if they need updated state.
type ResetCliDataOutput struct{}
