package notification

import (
	"context"
	"fmt"
	"testing"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/notification/notification_dto"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/storage"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// newTestNotificationService creates a notification service backed by isolated in-memory SQLite.
func newTestNotificationService(t *testing.T) *Notification {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}

	stor, err := storage.NewStorageFromDB(db)
	if err != nil {
		t.Fatal(err)
	}

	return NewNotification(stor)
}

// TestNotifyAndListNotifications verifies notifications are persisted and returned through the service.
func TestNotifyAndListNotifications(t *testing.T) {
	svc := newTestNotificationService(t)

	_, err := svc.Notify(context.Background(), notification_dto.NotifyInput{
		SessionID: 3,
		Kind:      "needs_attention",
		Title:     "Need help",
		Message:   "Please reply",
	})
	if err != nil {
		t.Fatalf("Notify: %v", err)
	}

	listed, err := svc.ListNotifications(context.Background(), notification_dto.ListNotificationsInput{})
	if err != nil {
		t.Fatalf("ListNotifications: %v", err)
	}
	if len(listed.Items) != 1 {
		t.Fatalf("expected 1 notification, got %d", len(listed.Items))
	}
	if listed.Items[0].Title != "Need help" {
		t.Fatalf("expected title to round-trip, got %q", listed.Items[0].Title)
	}
}

// TestResolveAndDismissNotification verifies resolve filters items and dismiss deletes them.
func TestResolveAndDismissNotification(t *testing.T) {
	svc := newTestNotificationService(t)

	created, err := svc.Notify(context.Background(), notification_dto.NotifyInput{
		SessionID: 7,
		Kind:      "info",
		Title:     "Done",
		Message:   "Background task completed",
	})
	if err != nil {
		t.Fatalf("Notify: %v", err)
	}

	if _, err := svc.ResolveNotification(context.Background(), notification_dto.ResolveNotificationInput{ID: created.NotificationID}); err != nil {
		t.Fatalf("ResolveNotification: %v", err)
	}

	listed, err := svc.ListNotifications(context.Background(), notification_dto.ListNotificationsInput{})
	if err != nil {
		t.Fatalf("ListNotifications: %v", err)
	}
	if len(listed.Items) != 0 {
		t.Fatalf("expected resolved items to be hidden, got %d", len(listed.Items))
	}

	if _, err := svc.DismissNotification(context.Background(), notification_dto.DismissNotificationInput{ID: created.NotificationID}); err != nil {
		t.Fatalf("DismissNotification: %v", err)
	}
}
