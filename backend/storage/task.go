package storage

import (
	"context"
	"errors"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gorm.io/gorm"
)

func (s *Storage) CreateTask(ctx context.Context, task data_models.Task) error {
	return s.sqliteDB.WithContext(ctx).Create(&task).Error
}

func (s *Storage) SaveTask(ctx context.Context, task data_models.Task) error {
	var existing data_models.Task
	err := s.sqliteDB.WithContext(ctx).Where("task_uuid = ?", task.TaskUuid).First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.sqliteDB.WithContext(ctx).Create(&task).Error
		}
		return err
	}
	return s.sqliteDB.WithContext(ctx).Where("task_uuid = ?", task.TaskUuid).Updates(&task).Error
}

func (s *Storage) GetTask(ctx context.Context, taskUuid string) (*data_models.Task, error) {
	var task data_models.Task
	err := s.sqliteDB.WithContext(ctx).Where("task_uuid = ?", taskUuid).First(&task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

func (s *Storage) GetChatActiveTask(ctx context.Context, chatUuid string) (*data_models.Task, error) {
	var task data_models.Task
	err := s.sqliteDB.WithContext(ctx).
		Where("chat_uuid = ? AND status IN ?", chatUuid, []data_models.TaskStatus{data_models.TaskStatusPending, data_models.TaskStatusRunning}).
		Order("created_at DESC").
		First(&task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

func (s *Storage) GetRunningTasks(ctx context.Context) ([]data_models.Task, error) {
	return s.GetTasksByStatus(ctx, []data_models.TaskStatus{data_models.TaskStatusPending, data_models.TaskStatusRunning})
}

func (s *Storage) GetTasksByStatus(ctx context.Context, statuses []data_models.TaskStatus) ([]data_models.Task, error) {
	var tasks []data_models.Task
	err := s.sqliteDB.WithContext(ctx).
		Where("status IN ?", statuses).
		Order("created_at DESC").
		Find(&tasks).Error
	if err != nil {
		return nil, err
	}
	return tasks, nil
}
