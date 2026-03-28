package event

import "fmt"

type EventType string

const (
	EventTypeMsg  EventType = "msg"
	EventTypeTask EventType = "task"
)

func GenEventsKey(eventType EventType, info string) string {
	return fmt.Sprintf("event:%s:%s", eventType, info)
}
