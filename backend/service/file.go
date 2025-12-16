package service

import (
	"mime"
	"path/filepath"

	"github.com/wailsapp/wails/v3/pkg/application"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils/ierror"
)

func (s *Service) SelectFiles() ([]view_models.File, error) {

	// 调用系统接口选择文件
	pattern := "*.txt;*.log;*.text;*.json;*.html;*.css;*.scss;*.jpg;*.png;*.jpeg;*.bmp"
	paths, err := application.OpenFileDialog().SetTitle("").AddFilter("选择文件", pattern).PromptForMultipleSelection()
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
			mimeType = "application/octet-stream"
		}

		// 通过mineType获取消息类型
		chatMessagePartType, err := utils.MimeType2ChatMessagePartType(mimeType)
		if err != nil {
			return nil, ierror.NewError(err)
		}

		// todo 如果为图像，则设置预览base64 200x200
		var previewImg *string

		file := view_models.File{
			ChatMessagePartType: chatMessagePartType,
			PreviewImg:          previewImg,
			Name:                filepath.Base(path),
			FilePath:            path,
			MineType:            mimeType,
		}

		files = append(files, file)
	}

	return files, nil
}

func (s *Service) OpenFile(path string) error {
	return s.app.Browser.OpenFile(path)
}
