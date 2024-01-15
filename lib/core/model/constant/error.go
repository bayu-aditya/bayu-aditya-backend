package constant

import "errors"

var (
	ErrRoomNotFound    = errors.New("room not found")
	ErrUnauthenticated = errors.New("unauthenticated")
)
