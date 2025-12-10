package service

import (
	"mime"
	"os"
	"path/filepath"

	"github.com/cloudwego/eino/schema"
	"github.com/wailsapp/wails/v3/pkg/application"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils/ierror"
)

func (s *Service) SelectFiles(fileType view_models.FileType) ([]view_models.File, error) {
	displayName := ""
	pattern := ""
	switch fileType {
	case view_models.FileTypeImg:
		displayName = "选择图片"
		pattern = "*.jpg;*.png;*.jpeg;*.bmp"
	case view_models.FileTypeText:
		displayName = "选择文本"
		pattern = "*.txt;*.log;*.text;*.json;*.html;*.css;*.scss"
	default:
		return nil, ierror.New(ierror.ErrCodeUnsupportedFileType)
	}

	paths, err := application.OpenFileDialog().SetTitle("").AddFilter(displayName, pattern).PromptForMultipleSelection()
	if err != nil {
		return nil, ierror.NewError(err)
	}
	if len(paths) == 0 {
		return []view_models.File{}, nil
	}

	files := make([]view_models.File, 0, len(paths))
	for _, path := range paths {
		// 获取文件MIMEType
		ext := filepath.Ext(path)
		mimeType := mime.TypeByExtension(ext)
		if mimeType == "" {
			// 如果无法通过扩展名获取，使用默认值
			mimeType = "application/octet-stream"
		}
		// local path
		localPath := "local:" + path
		// 组装
		var part schema.MessageInputPart
		switch fileType {
		case view_models.FileTypeImg:
			part.Type = schema.ChatMessagePartTypeImageURL
			part.Image = &schema.MessageInputImage{
				MessagePartCommon: schema.MessagePartCommon{
					URL:      &localPath,
					MIMEType: mimeType,
				},
				Detail: schema.ImageURLDetailHigh,
			}
		case view_models.FileTypeText:
			fileBytes, err := os.ReadFile(path)
			if err != nil {
				return nil, ierror.NewError(err)
			}
			part.Type = schema.ChatMessagePartTypeText
			part.Text = string(fileBytes)
		default:
			return nil, ierror.New(ierror.ErrCodeUnsupportedFileType)
		}

		file := view_models.File{
			FilePath: path,
			Name:     filepath.Base(path),
			Part:     part,
		}

		// 如果是图片类型，生成预览
		if fileType == view_models.FileTypeImg {
			preview, err := utils.GenerateImagePreview(path)
			if err == nil {
				// 只有成功生成预览时才设置
				file.PreviewImg = &preview
			}
			// 如果生成预览失败，仍然添加文件，但不设置预览
		}

		files = append(files, file)
	}

	return files, nil
}
