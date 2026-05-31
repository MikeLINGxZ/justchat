package file_dto

type SelectFileInput struct {
	DefaultFolderPath string `json:"default_folder_path"`
}

type SelectFileOutput struct {
	FilePath string `json:"file_path"`
}
