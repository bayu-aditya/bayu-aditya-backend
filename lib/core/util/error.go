package util

import (
	"fmt"
)

// ErrWrap if err is nil, then return nil
func ErrWrap(prefix string, err error, message string) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf("[%s] %s: %w", prefix, message, err)
}
