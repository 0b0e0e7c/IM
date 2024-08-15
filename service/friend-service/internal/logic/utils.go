package logic

import "fmt"

type FriendServiceError struct {
	Message string
}

func (e *FriendServiceError) Error() string {
	return fmt.Sprintf("friend service error: %s", e.Message)
}
