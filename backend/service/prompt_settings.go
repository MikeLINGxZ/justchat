package service

import (
	"fmt"
	"os"
	"strings"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/prompts"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils/ierror"
)

func (s *Service) reloadPromptSet() error {
	promptSet, err := prompts.LoadPromptSet()
	if err != nil {
		s.prompts = promptSet
		return err
	}
	s.prompts = promptSet
	return nil
}

func (s *Service) ListPromptFiles() ([]view_models.PromptFileSummary, error) {
	items := prompts.PromptFiles()
	result := make([]view_models.PromptFileSummary, 0, len(items))
	for _, item := range items {
		path, err := prompts.PromptPath(item.Name)
		if err != nil {
			return nil, ierror.NewError(err)
		}
		var updatedAt *time.Time
		info, statErr := os.Stat(path)
		if statErr == nil {
			modTime := info.ModTime()
			updatedAt = &modTime
		}
		result = append(result, view_models.PromptFileSummary{
			Name:        item.Name,
			Title:       item.Title,
			Description: item.Description,
			IsSystem:    item.IsSystem,
			UpdatedAt:   updatedAt,
		})
	}
	return result, nil
}

func (s *Service) GetPromptFile(name string) (*view_models.PromptFileDetail, error) {
	meta, ok := prompts.FindPromptMetadata(name)
	if !ok {
		return nil, ierror.NewError(fmt.Errorf("unsupported prompt file: %s", name))
	}
	fallback, ok := prompts.DefaultPromptContent(name)
	if !ok {
		return nil, ierror.NewError(fmt.Errorf("default prompt not found: %s", name))
	}
	content, err := prompts.LoadPrompt(name, fallback)
	if err != nil {
		// Preserve fallback content while surfacing the error to the caller only as a warning-level log via caller handling.
	}
	path, pathErr := prompts.PromptPath(name)
	if pathErr != nil {
		return nil, ierror.NewError(pathErr)
	}
	var updatedAt *time.Time
	info, statErr := os.Stat(path)
	if statErr == nil {
		modTime := info.ModTime()
		updatedAt = &modTime
	}
	return &view_models.PromptFileDetail{
		Name:        meta.Name,
		Title:       meta.Title,
		Description: meta.Description,
		IsSystem:    meta.IsSystem,
		Content:     content,
		UpdatedAt:   updatedAt,
	}, nil
}

func (s *Service) UpdatePromptFile(name string, content string) (*view_models.PromptFileDetail, error) {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil, ierror.NewError(fmt.Errorf("prompt content cannot be empty"))
	}
	if err := prompts.SavePromptContent(name, content); err != nil {
		return nil, ierror.NewError(err)
	}
	if err := s.reloadPromptSet(); err != nil {
		logger.Warm("reload prompt set fallback:", err)
	}
	return s.GetPromptFile(name)
}

func (s *Service) ResetPromptFile(name string) (*view_models.PromptFileDetail, error) {
	content, ok := prompts.DefaultPromptContent(name)
	if !ok {
		return nil, ierror.NewError(fmt.Errorf("default prompt not found: %s", name))
	}
	return s.UpdatePromptFile(name, content)
}
