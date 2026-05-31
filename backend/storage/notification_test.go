package storage

import (
	"testing"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
)

// TestNotificationCRUD verifies notification create/list/resolve behavior.
func TestNotificationCRUD(t *testing.T) {
	stor := newTestStorage(t)

	n, err := stor.CreateNotification(data_models.Notification{
		SessionID: 1,
		Kind:      "needs_attention",
		Title:     "T",
		Message:   "M",
	})
	if err != nil {
		t.Fatal(err)
	}

	list, err := stor.ListNotifications(false)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("want 1 got %d", len(list))
	}

	if err := stor.ResolveNotification(n.ID); err != nil {
		t.Fatal(err)
	}

	list, err = stor.ListNotifications(false)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 0 {
		t.Fatalf("expected resolved notifications to be hidden, got %d", len(list))
	}
}
