//go:build windows
// +build windows

package win

import "errors"

// ErrNotFound is returned when no matching window is found.
var ErrNotFound = errors.New("robotgo: window not found")

// ErrNotSupported is returned when an operation is not supported.
var ErrNotSupported = errors.New("robotgo: operation not supported")
