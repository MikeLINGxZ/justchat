package notification

import (
	"context"
	"errors"

	"github.com/wailsapp/wails/v3/pkg/application"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/ierror"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/notification/notification_dto"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/storage"
)

const (
	eventNotificationCreated  = "notification.created"
	eventNotificationResolved = "notification.resolved"
)

// Notification exposes global notification CRUD and waiter coordination to Wails.
type Notification struct {
	store    *storage.Storage
	wailsApp *application.App
}

// NewNotification creates a notification service bound to the shared storage layer.
func NewNotification(store *storage.Storage) *Notification {
	return &Notification{store: store}
}

// Notify creates a notification record and broadcasts a created event.
func (n *Notification) Notify(ctx context.Context, input notification_dto.NotifyInput) (*notification_dto.NotifyOutput, error) {
	created, err := n.store.CreateNotification(data_models.Notification{
		SessionID: input.SessionID,
		Kind:      input.Kind,
		Title:     input.Title,
		Message:   input.Message,
	})
	if err != nil {
		return nil, ierror.Error(ierror.ErrNotificationCreate, err)
	}
	n.emit(eventNotificationCreated, dtoFromModel(*created))
	return &notification_dto.NotifyOutput{NotificationID: created.ID}, nil
}

// ListNotifications returns the current notification list ordered by creation time.
func (n *Notification) ListNotifications(ctx context.Context, input notification_dto.ListNotificationsInput) (*notification_dto.ListNotificationsOutput, error) {
	items, err := n.store.ListNotifications(input.IncludeResolved)
	if err != nil {
		return nil, ierror.Error(ierror.ErrNotificationList, err)
	}

	result := make([]notification_dto.NotificationItem, 0, len(items))
	for _, item := range items {
		result = append(result, dtoFromModel(item))
	}
	return &notification_dto.ListNotificationsOutput{Items: result}, nil
}

// ResolveNotification marks a notification resolved, emits an event, and wakes waiters as success.
func (n *Notification) ResolveNotification(ctx context.Context, input notification_dto.ResolveNotificationInput) (*notification_dto.ResolveNotificationOutput, error) {
	if err := n.store.ResolveNotification(input.ID); err != nil {
		return nil, ierror.Error(ierror.ErrNotificationResolve, err)
	}
	wakeWaiter(input.ID, nil)
	n.emit(eventNotificationResolved, map[string]uint{"id": input.ID})
	return &notification_dto.ResolveNotificationOutput{}, nil
}

// RejectNotification marks a notification resolved and wakes any waiter with a rejection error.
func (n *Notification) RejectNotification(ctx context.Context, input notification_dto.ResolveNotificationInput) (*notification_dto.ResolveNotificationOutput, error) {
	if err := n.store.ResolveNotification(input.ID); err != nil {
		return nil, ierror.Error(ierror.ErrNotificationResolve, err)
	}
	wakeWaiter(input.ID, errors.New("user rejected"))
	n.emit(eventNotificationResolved, map[string]uint{"id": input.ID})
	return &notification_dto.ResolveNotificationOutput{}, nil
}

// DismissNotification permanently removes a notification row.
func (n *Notification) DismissNotification(ctx context.Context, input notification_dto.DismissNotificationInput) (*notification_dto.DismissNotificationOutput, error) {
	if err := n.store.DeleteNotification(input.ID); err != nil {
		return nil, ierror.Error(ierror.ErrNotificationDismiss, err)
	}
	wakeWaiter(input.ID, errors.New("notification dismissed"))
	n.emit(eventNotificationResolved, map[string]uint{"id": input.ID})
	return &notification_dto.DismissNotificationOutput{}, nil
}
