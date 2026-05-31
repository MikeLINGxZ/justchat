package notification_dto

// ResolveNotificationInput identifies the notification to resolve.
type ResolveNotificationInput struct {
	ID uint `json:"id"`
}

// ResolveNotificationOutput is returned after a resolve or reject action completes.
type ResolveNotificationOutput struct{}
