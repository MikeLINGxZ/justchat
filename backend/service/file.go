package service

import (
	"mime"
	"os"
	"path/filepath"
	"strings"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils/ierror"
)

func (s *Service) SelectFiles() ([]view_models.FileInfo, error) {

	// 调用系统接口选择文件
	pattern := "*.txt;*.log;*.text;*.json;*.html;*.css;*.scss;*.jpg;*.png;*.jpeg;*.bmp"
	paths, err := s.app.Dialog.OpenFile().SetTitle("").AddFilter("选择文件", pattern).PromptForMultipleSelection()
	if err != nil {
		return nil, ierror.NewError(err)
	}
	if len(paths) == 0 {
		return []view_models.FileInfo{}, nil
	}

	return s.fileInfo(paths)
}

func (s *Service) OpenFile(path string) error {
	return s.app.Browser.OpenFile(path)
}

func (s *Service) fileInfo(paths []string) ([]view_models.FileInfo, error) {
	result := make([]view_models.FileInfo, 0, len(paths))

	for _, p := range paths {
		stat, err := os.Stat(p)
		if err != nil {
			return nil, ierror.NewError(err)
		}

		mimeType := mime.TypeByExtension(filepath.Ext(p))
		if mimeType == "" {
			mimeType = "application/octet-stream"
		}

		info := view_models.FileInfo{
			Name:     stat.Name(),
			Path:     p,
			MineType: mimeType,
			Size:     stat.Size(),
		}

		if strings.HasPrefix(mimeType, "image/") {
			preview, err := utils.GenerateImagePreview(p)
			if err == nil {
				info.Preview = &preview
			}
		}

		result = append(result, info)
	}

	return result, nil
}
