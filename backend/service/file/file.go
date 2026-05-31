package file

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/dir"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/i18n"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/ierror"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/file/file_dto"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type File struct {
	wailsApp *application.App
}

// SelectFolder select a folder
func (f *File) SelectFolder(ctx context.Context, input file_dto.SelectFolderInput) (*file_dto.SelectFolderOutput, error) {
	if input.FolderPath == "" {
		input.FolderPath = dir.HomeDir()
	}
	path, err := f.wailsApp.Dialog.OpenFile().
		SetTitle(i18n.TCurrent(i18n.TCurrent("select.folder", nil), nil)).
		SetDirectory(input.FolderPath).
		CanChooseDirectories(true).
		CanChooseFiles(false).
		PromptForSingleSelection()

	if err != nil || path == "" {
		return nil, ierror.Error(ierror.ErrFileSelectFolder, err)
	}

	return &file_dto.SelectFolderOutput{FolderPath: path}, nil
}

// SelectFile select a file
func (f *File) SelectFile(ctx context.Context, input file_dto.SelectFileInput) (*file_dto.SelectFileOutput, error) {
	if input.DefaultFolderPath == "" {
		input.DefaultFolderPath = dir.HomeDir()
	}
	path, err := f.wailsApp.Dialog.OpenFile().
		SetTitle(i18n.TCurrent("select.file", nil)).
		SetDirectory(input.DefaultFolderPath).
		CanChooseFiles(true).
		PromptForSingleSelection()
	if err != nil || path == "" {
		return nil, ierror.Error(ierror.ErrFileSelectFile, err)
	}

	return &file_dto.SelectFileOutput{FilePath: path}, nil
}

// OpenFile open a file
func (f *File) OpenFile(ctx context.Context, input file_dto.OpenFileInput) (*file_dto.OpenFileOutput, error) {
	err := f.wailsApp.Browser.OpenFile(input.Path)
	if err != nil {
		return nil, ierror.Error(ierror.ErrFileOpen, err)
	}
	return &file_dto.OpenFileOutput{}, nil
}

// SaveTempFile decodes base64 data and writes it to the app temp directory.
// Returns the absolute file path on success.
func (f *File) SaveTempFile(ctx context.Context, input file_dto.SaveTempFileInput) (*file_dto.SaveTempFileOutput, error) {
	data, err := base64.StdEncoding.DecodeString(input.Data)
	if err != nil {
		return nil, ierror.Error(ierror.ErrFileSaveTempFile, err)
	}

	tmpDir := filepath.Join(os.TempDir(), "lemontea")
	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		return nil, ierror.Error(ierror.ErrFileSaveTempFile, err)
	}

	safeName := filepath.Base(input.Name)
	if safeName == "" || safeName == "." {
		safeName = "file"
	}
	name := fmt.Sprintf("%d-%s", time.Now().UnixNano(), safeName)
	path := filepath.Join(tmpDir, name)

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return nil, ierror.Error(ierror.ErrFileSaveTempFile, err)
	}

	return &file_dto.SaveTempFileOutput{FilePath: path}, nil
}
