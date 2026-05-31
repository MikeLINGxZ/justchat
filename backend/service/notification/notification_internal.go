package notification

import (
	"context"
	"errors"
	"sync"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/notification/notification_dto"
)

type waitResult struct {
	err error
}

var (
	waitersMu sync.Mutex
	waiters   = map[uint]chan waitResult{}
)

// NotifyAttention creates a needs_attention notification and registers a waiter for user action.
func (n *Notification) NotifyAttention(ctx context.Context, sessionID uint, title, message string) (uint, error) {
	out, err := n.Notify(ctx, notification_dto.NotifyInput{
		SessionID: sessionID,
		Kind:      "needs_attention",
		Title:     title,
		Message:   message,
	})
	if err != nil {
		return 0, err
	}

	waitersMu.Lock()
	waiters[out.NotificationID] = make(chan waitResult, 1)
	waitersMu.Unlock()
	return out.NotificationID, nil
}

// WaitForResolution blocks until the notification is resolved, rejected, dismissed, or times out.
func (n *Notification) WaitForResolution(ctx context.Context, notificationID uint) error {
	waitersMu.Lock()
	ch, ok := waiters[notificationID]
	waitersMu.Unlock()
	if !ok {
		return errors.New("no waiter registered")
	}

	select {
	case result := <-ch:
		return result.err
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(30 * time.Minute):
		return errors.New("attention request timed out")
	}
}

// emit sends a global event when the Wails app is available.
func (n *Notification) emit(name string, payload any) {
	if n.wailsApp != nil {
		n.wailsApp.Event.Emit(name, payload)
	}
}

// dtoFromModel converts the storage model to the public DTO.
func dtoFromModel(model data_models.Notification) notification_dto.NotificationItem {
	return notification_dto.NotificationItem{
		ID:         model.ID,
		SessionID:  model.SessionID,
		Kind:       model.Kind,
		Title:      model.Title,
		Message:    model.Message,
		CreatedAt:  model.CreatedAt,
		ResolvedAt: model.ResolvedAt,
	}
}

// wakeWaiter resolves one waiting notification request and clears its registration.
func wakeWaiter(notificationID uint, err error) {
	waitersMu.Lock()
	ch, ok := waiters[notificationID]
	if ok {
		delete(waiters, notificationID)
	}
	waitersMu.Unlock()
	if ok {
		ch <- waitResult{err: err}
		close(ch)
	}
}
