package notification_dto

// DismissNotificationInput identifies the notification to remove.
type DismissNotificationInput struct {
	ID uint `json:"id"`
}

// DismissNotificationOutput is returned after deletion.
type DismissNotificationOutput struct{}
