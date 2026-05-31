package notification_dto

// NotifyInput defines the payload for creating a notification.
type NotifyInput struct {
	SessionID uint   `json:"session_id"`
	Kind      string `json:"kind"`
	Title     string `json:"title"`
	Message   string `json:"message"`
}

// NotifyOutput returns the created notification identifier.
type NotifyOutput struct {
	NotificationID uint `json:"notification_id"`
}
