package logic

import "fmt"

type MessageServiceError struct {
	Message string
}

func (e *MessageServiceError) Error() string {
	return fmt.Sprintf("friend service error: %s", e.Message)
}
