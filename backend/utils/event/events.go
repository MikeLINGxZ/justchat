package event

import "fmt"

type EventType string

const (
	EventTypeMsg       EventType = "msg"
	EventTypeTask      EventType = "task"
	EventTypeChatTitle EventType = "chat_title"
)

func GenEventsKey(eventType EventType, info string) string {
	return fmt.Sprintf("event:%s:%s", eventType, info)
}
