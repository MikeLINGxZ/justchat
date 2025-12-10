package view_models

import "github.com/cloudwego/eino/schema"

type File struct {
	PreviewImg *string                 `json:"preview"`
	Name       string                  `json:"name"`
	FilePath   string                  `json:"file_path"`
	Part       schema.MessageInputPart `json:"part"`
}

type FileType int

const (
	FileTypeUnknown FileType = iota
	FileTypeImg
	FileTypeText
)
