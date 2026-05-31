package notification_dto

import "time"

// ListNotificationsInput controls whether resolved notifications are included.
type ListNotificationsInput struct {
	IncludeResolved bool `json:"include_resolved"`
}

// NotificationItem is the frontend-facing notification DTO.
type NotificationItem struct {
	ID         uint       `json:"id"`
	SessionID  uint       `json:"session_id"`
	Kind       string     `json:"kind"`
	Title      string     `json:"title"`
	Message    string     `json:"message"`
	CreatedAt  time.Time  `json:"created_at"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
}

// ListNotificationsOutput wraps the notification list.
type ListNotificationsOutput struct {
	Items []NotificationItem `json:"items"`
}
