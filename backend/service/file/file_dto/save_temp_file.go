package file_dto

type SaveTempFileInput struct {
	Name string `json:"name"` // original filename, e.g. "screenshot.png"
	Data string `json:"data"` // base64-encoded file content
	Mime string `json:"mime"` // e.g. "image/png"
}

type SaveTempFileOutput struct {
	FilePath string `json:"file_path"`
}
