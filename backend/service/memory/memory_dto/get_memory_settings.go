package memory_dto

type GetMemorySettingsInput struct {
}

type GetMemorySettingsOutput struct {
	Enabled bool `json:"enabled"`
}
