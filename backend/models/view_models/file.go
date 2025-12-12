package view_models

type File struct {
	PreviewImg *string `json:"preview"`
	Name       string  `json:"name"`
	FilePath   string  `json:"file_path"`
	MineType   string  `json:"mine_type"`
}
