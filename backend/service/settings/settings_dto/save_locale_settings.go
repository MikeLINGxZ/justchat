package settings_dto

type SaveLocaleSettingsInput struct {
	Locale   string `json:"locale"`
	Language string `json:"language"`
}

type SaveLocaleSettingsOutput struct {
}
