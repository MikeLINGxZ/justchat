package utils

import "fmt"

func GenEventsKey(in string) string {
	return fmt.Sprintf("user:%s", in)
}
